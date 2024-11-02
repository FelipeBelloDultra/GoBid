package services

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type MessageKind = int

const (
	PlaceBid MessageKind = iota
	SuccessfullyPlacedBid
	FailedToPlaceBid
	NewBidPlaced
	AuctionFinished
)

type Message struct {
	Message string
	Amount  float64
	Kind    MessageKind
	UserID  uuid.UUID
}

type AuctionLobby struct {
	sync.Mutex
	Rooms map[uuid.UUID]*AuctionRoom
}

type AuctionRoom struct {
	ID          uuid.UUID
	Context     context.Context
	Broadcast   chan Message
	Unregister  chan *Client
	Register    chan *Client
	Clients     map[uuid.UUID]*Client
	BidsService BidsService
}

func (r *AuctionRoom) registerClient(c *Client) {
	slog.Info("New user connected", "Client", c)
	r.Clients[c.UserID] = c
}

func (r *AuctionRoom) unregisterClient(c *Client) {
	slog.Info("User disconnected", "Client", c)
	delete(r.Clients, c.UserID)
}

func (r *AuctionRoom) broadcastMessage(m Message) {
	slog.Info("New message received", "RoomID", r.ID, "Message", m.Message, "UserID", m.UserID)
	switch m.Kind {
	case PlaceBid:
		bid, err := r.BidsService.PlaceBid(r.Context, r.ID, m.UserID, m.Amount)
		if err != nil {
			if errors.Is(err, ErrBidIsTooLow) {
				if client, ok := r.Clients[m.UserID]; ok {
					client.Send <- Message{Kind: FailedToPlaceBid, Message: ErrBidIsTooLow.Error()}
				}
				return
			}
		}

		if client, ok := r.Clients[m.UserID]; ok {
			client.Send <- Message{Kind: SuccessfullyPlacedBid, Message: "your bid was successfully placed"}
		}

		for id, client := range r.Clients {
			newBidMessage := Message{
				Kind:    NewBidPlaced,
				Message: "a new bid was placed",
				Amount:  bid.BidAmount,
			}
			if id == m.UserID {
				continue
			}
			client.Send <- newBidMessage
		}
	}
}

func (r *AuctionRoom) Run() {
	slog.Info("Auction has begun", "auctionID", r.ID)
	defer func() {
		close(r.Broadcast)
		close(r.Register)
		close(r.Unregister)
	}()

	for {
		select {
		case client := <-r.Register:
			r.registerClient(client)
		case client := <-r.Unregister:
			r.unregisterClient(client)
		case message := <-r.Broadcast:
			r.broadcastMessage(message)
		case <-r.Context.Done():
			slog.Info("Auction has ended", "auctionID", r.ID)

			for _, client := range r.Clients {
				client.Send <- Message{
					Kind:    AuctionFinished,
					Message: "auction has been finished",
				}
			}
			return
		}
	}
}

func NewAuctionRoom(ctx context.Context, id uuid.UUID, bidsService BidsService) *AuctionRoom {
	return &AuctionRoom{
		ID:          id,
		Broadcast:   make(chan Message),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Clients:     make(map[uuid.UUID]*Client),
		Context:     ctx,
		BidsService: bidsService,
	}
}

type Client struct {
	Room   *AuctionRoom
	Conn   *websocket.Conn
	UserID uuid.UUID
	Send   chan Message
}

func NewClient(room *AuctionRoom, conn *websocket.Conn, userId uuid.UUID) *Client {
	return &Client{
		Room:   room,
		Conn:   conn,
		Send:   make(chan Message, 512),
		UserID: userId,
	}
}

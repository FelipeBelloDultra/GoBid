package services

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type MessageKind = int

const (
	PlaceBid MessageKind = iota
)

type Message struct {
	Message string
	Kind    MessageKind
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

func NewAuctionRoom(ctx context.Context, id uuid.UUID, bidsService BidsService) *AuctionRoom {
	return &AuctionRoom{
		ID:          id,
		Broadcast:   make(chan Message),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
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

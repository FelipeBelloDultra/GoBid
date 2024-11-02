# GoBid

GoBid is an auction application built in Go, designed for `real-time` bidding on products. It utilizes GoLang for server-side logic, `chi` for routing, and SQL-based migrations and queries through `tern` and `sqlc`, respectively. The application supports user sign-ups, sign-ins, and product creation, where each new product initiates a unique auction room for users to participate.

## Features

- User Authentication: Sign up, sign in, and logout functionalities.
- Auction Rooms: Each product has a dedicated auction room for users to place bids in real-time.
- WebSocket Support: Enables real-time bidding and notifications for bid updates.
- Bid Management: Handles bid placements with validation for bid amounts and informs all clients in the room of new bids.
- Auction Lifecycle: Starts a new auction upon product creation and manages the auction end based on specified duration.

## Tech Stack

- Language: Go (GoLang)
- Routing: chi
- Database: PostgreSQL (via Docker)
- ORM/Query Generation: sqlc for generating type-safe queries
- Migrations: tern for database migrations
- WebSocket: gorilla/websocket for real-time communication

## Installation

1 Clone the repository:

```bash
git clone git@github.com:FelipeBelloDultra/GoBid.git
cd GoBid
```

2 Copy the `.env.example` to `.env` and update with your configuration

3 Start the PostgreSQL container:

```bash
docker compose up -d
```

4 Run database migrations:

```bash
go run cmd/tern-dotenv/main.go
```

5 Start the application:

```bash
go run cmd/api/main.go
```

## API Routes

### User Routes

- `POST /api/v1/users/sign-up` - Sign up a new user.
- `POST /api/v1/users/sign-in` - Sign in a user.
- `POST /api/v1/users/logout` - Logout (requires authentication).

### Product Routes

- `POST /api/v1/products` - Create a new product and initiate an auction room (requires authentication).
- `GET /api/v1/products/ws/subscribe/{product_id}` - WebSocket endpoint for subscribing to auction updates (requires authentication).

### Usage

Upon creating a product, an auction room is generated where users can place bids via WebSocket connections. The auction room manages clients, processes bids, and broadcasts bid updates and auction events to all participants.

### AuctionRoom Logic

- registerClient: Adds a client to the room.
- unregisterClient: Removes a client from the room.
- broadcastMessage: Distributes messages to clients, handling bid placements, success/failure notifications, and auction-end events.

### WebSocket Events

- PlaceBid: Triggered when a user places a bid.
- SuccessfullyPlacedBid: Sent to users when their bid is accepted.
- NewBidPlaced: Broadcasted to all users when a new bid is placed.
- AuctionFinished: Notifies all users that the auction has ended.

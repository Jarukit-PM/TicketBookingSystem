package ws

import (
	"time"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/inventory"
)

// Event type names sent to WebSocket clients.
const (
	EventSnapshot     = "snapshot"
	EventSeatHeld     = "seat_held"
	EventSeatReleased = "seat_released"
	EventSeatSold     = "seat_sold"
)

// Message is the wire format for server → client events.
type Message struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

// SeatHeldPayload is sent when a seat is held.
type SeatHeldPayload struct {
	SeatID    string    `json:"seatId"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// SeatReleasedPayload is sent when a hold expires or is released.
type SeatReleasedPayload struct {
	SeatID string `json:"seatId"`
}

// SeatSoldPayload is sent when a seat is confirmed sold.
type SeatSoldPayload struct {
	SeatID string `json:"seatId"`
}

// SnapshotPayload wraps a full inventory snapshot on connect.
type SnapshotPayload struct {
	Snapshot *inventory.Snapshot `json:"snapshot"`
}

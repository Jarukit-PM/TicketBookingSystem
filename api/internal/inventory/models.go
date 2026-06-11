package inventory

import (
	"time"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/catalog"
)

const (
	StatusAvailable = "AVAILABLE"
	StatusHeld      = "HELD"
	StatusSold      = "SOLD"
	StatusBlocked   = "BLOCKED"
)

// Seat is a layout seat with computed inventory status.
type Seat struct {
	SeatID string `json:"seatId"`
	Row    int    `json:"row"`
	Col    int    `json:"col"`
	Type   string `json:"type"`
	Status string `json:"status"`
}

// Snapshot is the public seat map state for one showtime.
type Snapshot struct {
	ShowtimeID string             `json:"showtimeId"`
	ScreenID   string             `json:"screenId"`
	ScreenName string             `json:"screenName"`
	MovieID    string             `json:"movieId"`
	StartsAt   time.Time          `json:"startsAt"`
	PriceTiers catalog.PriceTiers `json:"priceTiers"`
	Seats      []Seat             `json:"seats"`
}

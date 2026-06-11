package booking

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CollectionName = "bookings"

const StatusConfirmed = "CONFIRMED"

// Booking is a confirmed purchase persisted after checkout.
type Booking struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID `bson:"userId" json:"userId"`
	ShowtimeID  primitive.ObjectID `bson:"showtimeId" json:"showtimeId"`
	Seats       []string           `bson:"seats" json:"seats"`
	Total       int64              `bson:"total" json:"total"`
	BookingRef  string             `bson:"bookingRef" json:"bookingRef"`
	TicketToken string             `bson:"ticketToken" json:"-"`
	Status      string             `bson:"status" json:"status"`
	ConfirmedAt time.Time          `bson:"confirmedAt" json:"confirmedAt"`
}

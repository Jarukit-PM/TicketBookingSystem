package booking

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Repository persists confirmed bookings.
type Repository interface {
	Insert(ctx context.Context, b *Booking) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*Booking, error)
	FindByBookingRef(ctx context.Context, ref string) (*Booking, error)
	ListByUser(ctx context.Context, userID primitive.ObjectID) ([]Booking, error)
	ListConfirmedByUser(ctx context.Context, userID primitive.ObjectID) ([]Booking, error)
	ListConfirmedByShowtime(ctx context.Context, showtimeID primitive.ObjectID) ([]Booking, error)
	CountConfirmedBetween(ctx context.Context, from, to time.Time) (int, error)
	ListRecentConfirmed(ctx context.Context, limit int) ([]Booking, error)
}

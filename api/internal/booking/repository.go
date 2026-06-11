package booking

import (
	"context"

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
}

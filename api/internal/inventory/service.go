package inventory

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/booking"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/catalog"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/hold"
)

// Service computes seat inventory snapshots for showtimes.
type Service struct {
	showtimes catalog.ShowtimeRepository
	screens   catalog.ScreenRepository
	bookings  booking.Repository
	redis     *redis.Client
}

// NewService returns an inventory service backed by catalog repos and Redis holds.
func NewService(
	showtimes catalog.ShowtimeRepository,
	screens catalog.ScreenRepository,
	bookings booking.Repository,
	rdb *redis.Client,
) *Service {
	return &Service{
		showtimes: showtimes,
		screens:   screens,
		bookings:  bookings,
		redis:     rdb,
	}
}

// Snapshot returns the read-only seat map for a showtime.
func (s *Service) Snapshot(ctx context.Context, showtimeID primitive.ObjectID) (*Snapshot, error) {
	showtime, err := s.showtimes.FindShowtimeByID(ctx, showtimeID)
	if err != nil {
		return nil, fmt.Errorf("find showtime: %w", err)
	}
	if showtime == nil {
		return nil, ErrShowtimeNotFound
	}

	screen, err := s.screens.FindScreenByID(ctx, showtime.ScreenID)
	if err != nil {
		return nil, fmt.Errorf("find screen: %w", err)
	}
	if screen == nil {
		return nil, ErrScreenNotFound
	}

	bookings, err := s.bookings.ListConfirmedByShowtime(ctx, showtimeID)
	if err != nil {
		return nil, fmt.Errorf("list bookings: %w", err)
	}

	held, err := hold.ListHeldSeatIDs(ctx, s.redis, showtimeID.Hex())
	if err != nil {
		return nil, fmt.Errorf("list holds: %w", err)
	}

	return &Snapshot{
		ShowtimeID: showtime.ID.Hex(),
		ScreenID:   screen.ID.Hex(),
		ScreenName: screen.Name,
		MovieID:    showtime.MovieID.Hex(),
		StartsAt:   showtime.StartsAt,
		PriceTiers: showtime.PriceTiers,
		Seats:      ComputeSeats(screen.Layout, SoldSeatSet(bookings), held),
	}, nil
}

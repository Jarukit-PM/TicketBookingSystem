package admin

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/booking"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/catalog"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/user"
)

// ErrInvalidQuery is returned when a search filter value is malformed.
var ErrInvalidQuery = errors.New("invalid search query")

const adminBookingSearchLimit = 100

// BookingSearchQuery holds optional admin booking search filters.
// Priority when multiple are set: bookingRef > userId > email > showtimeId.
type BookingSearchQuery struct {
	Email      string
	BookingRef string
	UserID     string
	ShowtimeID string
}

// BookingSummary is a confirmed booking row for admin lookup.
type BookingSummary struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId,omitempty"`
	UserEmail   string    `json:"userEmail,omitempty"`
	BookingRef  string    `json:"bookingRef"`
	ShowtimeID  string    `json:"showtimeId"`
	MovieTitle  string    `json:"movieTitle"`
	Seats       []string  `json:"seats"`
	Total       int64     `json:"total"`
	ConfirmedAt time.Time `json:"confirmedAt"`
}

// BookingsService provides read-only admin booking lookup.
type BookingsService struct {
	Bookings  booking.Repository
	Showtimes catalog.ShowtimeRepository
	Movies    catalog.MovieRepository
	Users     user.Repository
}

// Search returns confirmed bookings matching the first non-empty filter.
func (s *BookingsService) Search(ctx context.Context, q BookingSearchQuery) ([]BookingSummary, error) {
	bookings, err := s.resolveSearch(ctx, q)
	if err != nil {
		return nil, err
	}
	return s.summarizeBookings(ctx, bookings)
}

// ListUserBookings returns confirmed booking history for a user.
func (s *BookingsService) ListUserBookings(ctx context.Context, userID primitive.ObjectID) ([]BookingSummary, error) {
	bookings, err := s.Bookings.ListByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list user bookings: %w", err)
	}
	return s.summarizeBookings(ctx, bookings)
}

func (s *BookingsService) resolveSearch(ctx context.Context, q BookingSearchQuery) ([]booking.Booking, error) {
	ref := strings.TrimSpace(q.BookingRef)
	userID := strings.TrimSpace(q.UserID)
	email := strings.TrimSpace(strings.ToLower(q.Email))
	showtimeID := strings.TrimSpace(q.ShowtimeID)

	switch {
	case ref != "":
		b, err := s.Bookings.FindByBookingRef(ctx, ref)
		if err != nil {
			return nil, fmt.Errorf("find booking by ref: %w", err)
		}
		if b == nil || b.Status != booking.StatusConfirmed {
			return []booking.Booking{}, nil
		}
		return []booking.Booking{*b}, nil
	case userID != "":
		oid, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			return nil, ErrInvalidQuery
		}
		bookings, err := s.Bookings.ListByUser(ctx, oid)
		if err != nil {
			return nil, fmt.Errorf("list bookings by user: %w", err)
		}
		return limitBookings(bookings, adminBookingSearchLimit), nil
	case email != "":
		u, err := s.Users.FindByEmail(ctx, email)
		if err != nil {
			return nil, fmt.Errorf("find user by email: %w", err)
		}
		if u == nil {
			return []booking.Booking{}, nil
		}
		bookings, err := s.Bookings.ListByUser(ctx, u.ID)
		if err != nil {
			return nil, fmt.Errorf("list bookings by email: %w", err)
		}
		return limitBookings(bookings, adminBookingSearchLimit), nil
	case showtimeID != "":
		oid, err := primitive.ObjectIDFromHex(showtimeID)
		if err != nil {
			return nil, ErrInvalidQuery
		}
		bookings, err := s.Bookings.ListConfirmedByShowtime(ctx, oid)
		if err != nil {
			return nil, fmt.Errorf("list bookings by showtime: %w", err)
		}
		return limitBookings(bookings, adminBookingSearchLimit), nil
	default:
		return []booking.Booking{}, nil
	}
}

func limitBookings(bookings []booking.Booking, limit int) []booking.Booking {
	if len(bookings) <= limit {
		return bookings
	}
	return bookings[:limit]
}

func (s *BookingsService) summarizeBookings(ctx context.Context, bookings []booking.Booking) ([]BookingSummary, error) {
	movieCache := make(map[primitive.ObjectID]string)
	showtimeCache := make(map[primitive.ObjectID]primitive.ObjectID)
	userCache := make(map[primitive.ObjectID]string)
	out := make([]BookingSummary, 0, len(bookings))

	for _, b := range bookings {
		movieID, err := s.movieIDForShowtime(ctx, showtimeCache, b.ShowtimeID)
		if err != nil {
			return nil, err
		}

		title, err := s.movieTitle(ctx, movieCache, movieID)
		if err != nil {
			return nil, err
		}

		email, err := s.userEmail(ctx, userCache, b.UserID)
		if err != nil {
			return nil, err
		}

		seats := b.Seats
		if seats == nil {
			seats = []string{}
		}

		out = append(out, BookingSummary{
			ID:          b.ID.Hex(),
			UserID:      b.UserID.Hex(),
			UserEmail:   email,
			BookingRef:  b.BookingRef,
			ShowtimeID:  b.ShowtimeID.Hex(),
			MovieTitle:  title,
			Seats:       seats,
			Total:       b.Total,
			ConfirmedAt: b.ConfirmedAt,
		})
	}

	return out, nil
}

func (s *BookingsService) movieIDForShowtime(
	ctx context.Context,
	cache map[primitive.ObjectID]primitive.ObjectID,
	showtimeID primitive.ObjectID,
) (primitive.ObjectID, error) {
	if movieID, ok := cache[showtimeID]; ok {
		return movieID, nil
	}
	showtime, err := s.Showtimes.FindShowtimeByID(ctx, showtimeID)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("find showtime %s: %w", showtimeID.Hex(), err)
	}
	if showtime == nil {
		return primitive.NilObjectID, nil
	}
	cache[showtimeID] = showtime.MovieID
	return showtime.MovieID, nil
}

func (s *BookingsService) movieTitle(
	ctx context.Context,
	cache map[primitive.ObjectID]string,
	movieID primitive.ObjectID,
) (string, error) {
	if movieID.IsZero() {
		return "Unknown movie", nil
	}
	if title, ok := cache[movieID]; ok {
		return title, nil
	}
	movie, err := s.Movies.FindMovieByID(ctx, movieID)
	if err != nil {
		return "", fmt.Errorf("find movie %s: %w", movieID.Hex(), err)
	}
	if movie == nil {
		cache[movieID] = "Unknown movie"
		return "Unknown movie", nil
	}
	cache[movieID] = movie.Title
	return movie.Title, nil
}

func (s *BookingsService) userEmail(
	ctx context.Context,
	cache map[primitive.ObjectID]string,
	userID primitive.ObjectID,
) (string, error) {
	if userID.IsZero() {
		return "", nil
	}
	if email, ok := cache[userID]; ok {
		return email, nil
	}
	u, err := s.Users.FindByID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("find user %s: %w", userID.Hex(), err)
	}
	if u == nil {
		cache[userID] = ""
		return "", nil
	}
	cache[userID] = u.Email
	return u.Email, nil
}

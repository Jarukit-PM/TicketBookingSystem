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

const (
	adminBookingSearchLimit  = 100
	defaultBookingPageLimit  = 20
)

// BookingSearchResult is a paginated admin booking search response.
type BookingSearchResult struct {
	Bookings []BookingSummary `json:"bookings"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	Limit    int              `json:"limit"`
}

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
	Locale      string    `json:"locale,omitempty"`
	ConfirmedAt time.Time `json:"confirmedAt"`
}

// BookingsService provides read-only admin booking lookup.
type BookingsService struct {
	Bookings  booking.Repository
	Showtimes catalog.ShowtimeRepository
	Movies    catalog.MovieRepository
	Users     user.Repository
}

// Search returns confirmed bookings matching the first non-empty filter, or all
// confirmed bookings when no filter is set. Results are paginated newest first.
func (s *BookingsService) Search(ctx context.Context, q BookingSearchQuery, page, limit int) (*BookingSearchResult, error) {
	limit = clampBookingLimit(limit)
	if page <= 0 {
		page = 1
	}

	bookings, total, err := s.resolveSearch(ctx, q, page, limit)
	if err != nil {
		return nil, err
	}
	summaries, err := s.summarizeBookings(ctx, bookings)
	if err != nil {
		return nil, err
	}
	return &BookingSearchResult{
		Bookings: summaries,
		Total:    total,
		Page:     page,
		Limit:    limit,
	}, nil
}

// ListUserBookings returns confirmed booking history for a user.
func (s *BookingsService) ListUserBookings(ctx context.Context, userID primitive.ObjectID) ([]BookingSummary, error) {
	bookings, err := s.Bookings.ListByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list user bookings: %w", err)
	}
	return s.summarizeBookings(ctx, bookings)
}

func (s *BookingsService) resolveSearch(
	ctx context.Context,
	q BookingSearchQuery,
	page, limit int,
) ([]booking.Booking, int64, error) {
	ref := strings.TrimSpace(q.BookingRef)
	userID := strings.TrimSpace(q.UserID)
	email := strings.TrimSpace(strings.ToLower(q.Email))
	showtimeID := strings.TrimSpace(q.ShowtimeID)

	switch {
	case ref != "":
		b, err := s.Bookings.FindByBookingRef(ctx, ref)
		if err != nil {
			return nil, 0, fmt.Errorf("find booking by ref: %w", err)
		}
		if b == nil || b.Status != booking.StatusConfirmed {
			return []booking.Booking{}, 0, nil
		}
		return []booking.Booking{*b}, 1, nil
	case userID != "":
		oid, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			return nil, 0, ErrInvalidQuery
		}
		bookings, err := s.Bookings.ListByUser(ctx, oid)
		if err != nil {
			return nil, 0, fmt.Errorf("list bookings by user: %w", err)
		}
		return paginateBookings(bookings, page, limit)
	case email != "":
		u, err := s.Users.FindByEmail(ctx, email)
		if err != nil {
			return nil, 0, fmt.Errorf("find user by email: %w", err)
		}
		if u == nil {
			return []booking.Booking{}, 0, nil
		}
		bookings, err := s.Bookings.ListByUser(ctx, u.ID)
		if err != nil {
			return nil, 0, fmt.Errorf("list bookings by email: %w", err)
		}
		return paginateBookings(bookings, page, limit)
	case showtimeID != "":
		oid, err := primitive.ObjectIDFromHex(showtimeID)
		if err != nil {
			return nil, 0, ErrInvalidQuery
		}
		bookings, err := s.Bookings.ListConfirmedByShowtime(ctx, oid)
		if err != nil {
			return nil, 0, fmt.Errorf("list bookings by showtime: %w", err)
		}
		return paginateBookings(bookings, page, limit)
	default:
		total, err := s.Bookings.CountConfirmed(ctx)
		if err != nil {
			return nil, 0, fmt.Errorf("count confirmed bookings: %w", err)
		}
		bookings, err := s.Bookings.ListConfirmedPage(ctx, SkipFor(page, limit), limit)
		if err != nil {
			return nil, 0, fmt.Errorf("list confirmed bookings: %w", err)
		}
		return bookings, total, nil
	}
}

func paginateBookings(bookings []booking.Booking, page, limit int) ([]booking.Booking, int64, error) {
	total := int64(len(bookings))
	if total > adminBookingSearchLimit {
		bookings = bookings[:adminBookingSearchLimit]
		total = adminBookingSearchLimit
	}
	start := SkipFor(page, limit)
	if start >= len(bookings) {
		return []booking.Booking{}, total, nil
	}
	end := start + limit
	if end > len(bookings) {
		end = len(bookings)
	}
	return bookings[start:end], total, nil
}

func clampBookingLimit(limit int) int {
	if limit <= 0 {
		return defaultBookingPageLimit
	}
	if limit > maxPageLimit {
		return maxPageLimit
	}
	return limit
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
			Locale:      booking.ParseLocale(b.Locale),
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

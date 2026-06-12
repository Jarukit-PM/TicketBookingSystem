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

const defaultBookingPageLimit = 20

// BookingSearchResult is a paginated admin booking search response.
type BookingSearchResult struct {
	Bookings []BookingSummary `json:"bookings"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	Limit    int              `json:"limit"`
}

// BookingSearchQuery holds optional admin booking search filters.
// Multiple filters are combined (AND). Email resolves to a user ID.
type BookingSearchQuery struct {
	Email         string
	BookingRef    string
	UserID        string
	ShowtimeID    string
	MovieID       string
	Locale        string
	ConfirmedFrom *time.Time
	ConfirmedTo   *time.Time // exclusive upper bound
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

// BookingDetail is a full confirmed booking for admin support lookup.
type BookingDetail struct {
	BookingSummary
	StartsAt    time.Time `json:"startsAt"`
	CinemaID    string    `json:"cinemaId"`
	CinemaName  string    `json:"cinemaName"`
	ScreenID    string    `json:"screenId"`
	ScreenName  string    `json:"screenName"`
	PosterURL   string    `json:"posterUrl"`
	Status      string    `json:"status"`
}

// BookingsService provides read-only admin booking lookup.
type BookingsService struct {
	Bookings  booking.Repository
	Showtimes catalog.ShowtimeRepository
	Screens   catalog.ScreenRepository
	Cinemas   catalog.CinemaRepository
	Movies    catalog.MovieRepository
	Users     user.Repository
}

// Search returns confirmed bookings matching the combined filters, or all
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

// GetByID returns a confirmed booking with venue and showtime context for admin detail views.
func (s *BookingsService) GetByID(ctx context.Context, bookingID primitive.ObjectID) (*BookingDetail, error) {
	b, err := s.Bookings.FindByID(ctx, bookingID)
	if err != nil {
		return nil, fmt.Errorf("find booking: %w", err)
	}
	if b == nil || b.Status != booking.StatusConfirmed {
		return nil, booking.ErrBookingNotFound
	}

	summaries, err := s.summarizeBookings(ctx, []booking.Booking{*b})
	if err != nil {
		return nil, err
	}
	if len(summaries) == 0 {
		return nil, booking.ErrBookingNotFound
	}

	showtime, err := s.Showtimes.FindShowtimeByID(ctx, b.ShowtimeID)
	if err != nil {
		return nil, fmt.Errorf("find showtime: %w", err)
	}
	if showtime == nil {
		return nil, booking.ErrBookingNotFound
	}

	screen, err := s.Screens.FindScreenByID(ctx, showtime.ScreenID)
	if err != nil {
		return nil, fmt.Errorf("find screen: %w", err)
	}
	if screen == nil {
		return nil, booking.ErrBookingNotFound
	}

	cinema, err := s.Cinemas.FindCinemaByID(ctx, screen.CinemaID)
	if err != nil {
		return nil, fmt.Errorf("find cinema: %w", err)
	}
	if cinema == nil {
		return nil, booking.ErrBookingNotFound
	}

	movie, err := s.Movies.FindMovieByID(ctx, showtime.MovieID)
	if err != nil {
		return nil, fmt.Errorf("find movie: %w", err)
	}

	posterURL := ""
	if movie != nil {
		posterURL = movie.PosterURL
	}

	return &BookingDetail{
		BookingSummary: summaries[0],
		StartsAt:       showtime.StartsAt,
		CinemaID:       cinema.ID.Hex(),
		CinemaName:     cinema.Name,
		ScreenID:       screen.ID.Hex(),
		ScreenName:     screen.Name,
		PosterURL:      posterURL,
		Status:         b.Status,
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
	filter, empty, err := s.buildConfirmedFilter(ctx, q)
	if err != nil {
		return nil, 0, err
	}
	if empty {
		return []booking.Booking{}, 0, nil
	}

	total, err := s.Bookings.CountConfirmedFiltered(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("count filtered bookings: %w", err)
	}
	bookings, err := s.Bookings.ListConfirmedFiltered(ctx, filter, SkipFor(page, limit), limit)
	if err != nil {
		return nil, 0, fmt.Errorf("list filtered bookings: %w", err)
	}
	return bookings, total, nil
}

func (s *BookingsService) buildConfirmedFilter(
	ctx context.Context,
	q BookingSearchQuery,
) (booking.ConfirmedFilter, bool, error) {
	filter := booking.ConfirmedFilter{
		BookingRef:    strings.TrimSpace(q.BookingRef),
		ConfirmedFrom: q.ConfirmedFrom,
		ConfirmedTo:   q.ConfirmedTo,
	}

	userID := strings.TrimSpace(q.UserID)
	email := strings.TrimSpace(strings.ToLower(q.Email))
	switch {
	case userID != "":
		oid, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			return filter, false, ErrInvalidQuery
		}
		filter.UserID = oid
	case email != "":
		u, err := s.Users.FindByEmail(ctx, email)
		if err != nil {
			return filter, false, fmt.Errorf("find user by email: %w", err)
		}
		if u == nil {
			return filter, true, nil
		}
		filter.UserID = u.ID
	}

	showtimeID := strings.TrimSpace(q.ShowtimeID)
	movieID := strings.TrimSpace(q.MovieID)
	if showtimeID != "" {
		oid, err := primitive.ObjectIDFromHex(showtimeID)
		if err != nil {
			return filter, false, ErrInvalidQuery
		}
		filter.ShowtimeID = oid
	} else if movieID != "" {
		oid, err := primitive.ObjectIDFromHex(movieID)
		if err != nil {
			return filter, false, ErrInvalidQuery
		}
		showtimes, err := s.Showtimes.ListShowtimesByMovie(ctx, oid)
		if err != nil {
			return filter, false, fmt.Errorf("list showtimes by movie: %w", err)
		}
		if len(showtimes) == 0 {
			return filter, true, nil
		}
		filter.ShowtimeIDs = make([]primitive.ObjectID, 0, len(showtimes))
		for _, st := range showtimes {
			filter.ShowtimeIDs = append(filter.ShowtimeIDs, st.ID)
		}
	}

	locale := strings.TrimSpace(strings.ToLower(q.Locale))
	if locale != "" {
		switch locale {
		case "en", "th":
			filter.Locale = locale
		default:
			return filter, false, ErrInvalidQuery
		}
	}

	return filter, false, nil
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

package admin

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/booking"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/catalog"
)

const recentBookingsLimit = 10

// Dashboard is the admin operational summary for today.
type Dashboard struct {
	BookingsToday   int              `json:"bookingsToday"`
	ShowtimesToday  int              `json:"showtimesToday"`
	AvgOccupancyPct float64          `json:"avgOccupancyPct"`
	RecentBookings  []BookingSummary `json:"recentBookings"`
}

// DashboardService aggregates catalog and booking data for the admin console.
type DashboardService struct {
	Showtimes catalog.ShowtimeRepository
	Screens   catalog.ScreenRepository
	Movies    catalog.MovieRepository
	Bookings  booking.Repository
}

// GetDashboard returns today's metrics and recent bookings (UTC day boundaries).
func (s *DashboardService) GetDashboard(ctx context.Context) (*Dashboard, error) {
	now := time.Now().UTC()
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	dayEnd := dayStart.Add(24 * time.Hour)

	bookingsToday, err := s.Bookings.CountConfirmedBetween(ctx, dayStart, dayEnd)
	if err != nil {
		return nil, fmt.Errorf("count bookings today: %w", err)
	}

	showtimes, err := s.Showtimes.ListAdminShowtimes(ctx, catalog.AdminShowtimeFilter{
		From: &dayStart,
		To:   &dayEnd,
	})
	if err != nil {
		return nil, fmt.Errorf("list showtimes today: %w", err)
	}

	avgOccupancy, err := s.avgOccupancyForShowtimes(ctx, showtimes)
	if err != nil {
		return nil, err
	}

	recent, err := s.Bookings.ListRecentConfirmed(ctx, recentBookingsLimit)
	if err != nil {
		return nil, fmt.Errorf("list recent bookings: %w", err)
	}

	summaries, err := s.summarizeBookings(ctx, recent)
	if err != nil {
		return nil, err
	}

	return &Dashboard{
		BookingsToday:   bookingsToday,
		ShowtimesToday:  len(showtimes),
		AvgOccupancyPct: avgOccupancy,
		RecentBookings:  summaries,
	}, nil
}

func (s *DashboardService) avgOccupancyForShowtimes(ctx context.Context, showtimes []catalog.Showtime) (float64, error) {
	if len(showtimes) == 0 {
		return 0, nil
	}

	screenCache := make(map[primitive.ObjectID]*catalog.Screen)
	var totalPct float64
	var counted int

	for _, showtime := range showtimes {
		screen, err := s.screenForShowtime(ctx, screenCache, showtime.ScreenID)
		if err != nil {
			return 0, err
		}
		if screen == nil {
			continue
		}

		sellable := catalog.SellableSeatCount(screen.Layout)
		if sellable == 0 {
			continue
		}

		bookings, err := s.Bookings.ListConfirmedByShowtime(ctx, showtime.ID)
		if err != nil {
			return 0, fmt.Errorf("list bookings for showtime %s: %w", showtime.ID.Hex(), err)
		}

		sold := 0
		for _, b := range bookings {
			sold += len(b.Seats)
		}

		totalPct += float64(sold) / float64(sellable) * 100
		counted++
	}

	if counted == 0 {
		return 0, nil
	}
	return totalPct / float64(counted), nil
}

func (s *DashboardService) screenForShowtime(
	ctx context.Context,
	cache map[primitive.ObjectID]*catalog.Screen,
	screenID primitive.ObjectID,
) (*catalog.Screen, error) {
	if screen, ok := cache[screenID]; ok {
		return screen, nil
	}
	screen, err := s.Screens.FindScreenByID(ctx, screenID)
	if err != nil {
		return nil, fmt.Errorf("find screen %s: %w", screenID.Hex(), err)
	}
	cache[screenID] = screen
	return screen, nil
}

func (s *DashboardService) summarizeBookings(ctx context.Context, bookings []booking.Booking) ([]BookingSummary, error) {
	movieCache := make(map[primitive.ObjectID]string)
	showtimeCache := make(map[primitive.ObjectID]primitive.ObjectID)
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

		seats := b.Seats
		if seats == nil {
			seats = []string{}
		}

		out = append(out, BookingSummary{
			ID:          b.ID.Hex(),
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

func (s *DashboardService) movieIDForShowtime(
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

func (s *DashboardService) movieTitle(
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

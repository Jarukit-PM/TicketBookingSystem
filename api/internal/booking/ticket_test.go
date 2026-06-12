package booking_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/booking"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/catalog"
)

type ticketBookingsRepo struct {
	byRef map[string]*booking.Booking
}

func (r *ticketBookingsRepo) Insert(context.Context, *booking.Booking) error { return nil }
func (r *ticketBookingsRepo) FindByID(context.Context, primitive.ObjectID) (*booking.Booking, error) {
	return nil, nil
}
func (r *ticketBookingsRepo) FindByBookingRef(_ context.Context, ref string) (*booking.Booking, error) {
	return r.byRef[ref], nil
}
func (r *ticketBookingsRepo) ListByUser(context.Context, primitive.ObjectID) ([]booking.Booking, error) {
	return nil, nil
}
func (r *ticketBookingsRepo) ListConfirmedByUser(context.Context, primitive.ObjectID) ([]booking.Booking, error) {
	return nil, nil
}
func (r *ticketBookingsRepo) ListConfirmedByShowtime(context.Context, primitive.ObjectID) ([]booking.Booking, error) {
	return nil, nil
}
func (r *ticketBookingsRepo) CountConfirmedBetween(context.Context, time.Time, time.Time) (int, error) {
	return 0, nil
}
func (r *ticketBookingsRepo) ListRecentConfirmed(context.Context, int) ([]booking.Booking, error) {
	return nil, nil
}
func (r *ticketBookingsRepo) CountConfirmed(context.Context) (int64, error) { return 0, nil }
func (r *ticketBookingsRepo) CountConfirmedFiltered(ctx context.Context, _ booking.ConfirmedFilter) (int64, error) {
	return 0, nil
}
func (r *ticketBookingsRepo) ListConfirmedPage(context.Context, int, int) ([]booking.Booking, error) {
	return nil, nil
}
func (r *ticketBookingsRepo) ListConfirmedFiltered(ctx context.Context, _ booking.ConfirmedFilter, _, _ int) ([]booking.Booking, error) {
	return nil, nil
}
func (r *ticketBookingsRepo) UpdateTicketToken(context.Context, primitive.ObjectID, string) error {
	return nil
}

type ticketMovies struct{ movie *catalog.Movie }

func (m ticketMovies) InsertMovie(context.Context, *catalog.Movie) error { return nil }
func (m ticketMovies) FindMovieByID(_ context.Context, id primitive.ObjectID) (*catalog.Movie, error) {
	if m.movie != nil && m.movie.ID == id {
		return m.movie, nil
	}
	return nil, nil
}
func (ticketMovies) ListMoviesByStatus(context.Context, string) ([]catalog.Movie, error) { return nil, nil }
func (ticketMovies) ListComingSoonMovies(context.Context) ([]catalog.Movie, error)      { return nil, nil }
func (ticketMovies) ListNonArchivedMovies(context.Context) ([]catalog.Movie, error)     { return nil, nil }
func (ticketMovies) ListMovies(context.Context) ([]catalog.Movie, error)                { return nil, nil }
func (ticketMovies) UpdateMovie(context.Context, *catalog.Movie) error                  { return nil }
func (ticketMovies) DeleteMovie(context.Context, primitive.ObjectID) error              { return nil }

func TestGetTicketByRefValidToken(t *testing.T) {
	const secret = "test-ticket-secret"
	bookingID := primitive.NewObjectID()
	showtimeID := primitive.NewObjectID()
	screenID := primitive.NewObjectID()
	cinemaID := primitive.NewObjectID()
	movieID := primitive.NewObjectID()
	ref := "TBS-ABC123"
	token := booking.SignTicketToken(secret, ref, bookingID.Hex())

	svc := booking.NewService(
		stubShowtimes{showtime: &catalog.Showtime{
			ID:       showtimeID,
			ScreenID: screenID,
			MovieID:  movieID,
			StartsAt: time.Date(2026, 6, 12, 19, 0, 0, 0, time.UTC),
		}},
		stubScreens{screen: &catalog.Screen{ID: screenID, CinemaID: cinemaID, Name: "Screen 1"}},
		stubCinemas{cinema: &catalog.Cinema{ID: cinemaID, Name: "Major Cineplex"}},
		ticketMovies{movie: &catalog.Movie{ID: movieID, Title: "Test Movie"}},
		&ticketBookingsRepo{byRef: map[string]*booking.Booking{
			ref: {
				ID:          bookingID,
				ShowtimeID:  showtimeID,
				BookingRef:  ref,
				TicketToken: token,
				Status:      booking.StatusConfirmed,
				Seats:       []string{"A1"},
				Total:       35000,
			},
		}},
		nil,
		nil,
		nil,
		booking.WithTicketConfig(secret, "http://localhost:5173"),
	)

	got, err := svc.GetTicketByRef(context.Background(), ref, token)
	if err != nil {
		t.Fatalf("GetTicketByRef: %v", err)
	}
	if got.BookingRef != ref {
		t.Fatalf("bookingRef = %q, want %q", got.BookingRef, ref)
	}
	if got.QRPngBase64 == "" {
		t.Fatalf("expected qr png payload")
	}
}

func TestGetTicketByRefRejectsTamperedToken(t *testing.T) {
	const secret = "test-ticket-secret"
	bookingID := primitive.NewObjectID()
	ref := "TBS-ABC123"
	token := booking.SignTicketToken(secret, ref, bookingID.Hex())

	svc := booking.NewService(
		stubShowtimes{},
		stubScreens{},
		stubCinemas{},
		ticketMovies{},
		&ticketBookingsRepo{byRef: map[string]*booking.Booking{
			ref: {
				ID:          bookingID,
				BookingRef:  ref,
				TicketToken: token,
				Status:      booking.StatusConfirmed,
			},
		}},
		nil,
		nil,
		nil,
		booking.WithTicketConfig(secret, "http://localhost:5173"),
	)

	_, err := svc.GetTicketByRef(context.Background(), ref, token+"x")
	if !errors.Is(err, booking.ErrInvalidTicket) {
		t.Fatalf("expected ErrInvalidTicket, got %v", err)
	}
}

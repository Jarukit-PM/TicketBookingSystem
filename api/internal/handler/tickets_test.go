package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/booking"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/catalog"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/handler"
)

const publicTicketSecret = "test-ticket-secret"

type publicTicketBookingsRepo struct {
	byRef map[string]*booking.Booking
}

func (r *publicTicketBookingsRepo) Insert(context.Context, *booking.Booking) error { return nil }
func (r *publicTicketBookingsRepo) FindByID(context.Context, primitive.ObjectID) (*booking.Booking, error) {
	return nil, nil
}
func (r *publicTicketBookingsRepo) FindByBookingRef(_ context.Context, ref string) (*booking.Booking, error) {
	if r.byRef == nil {
		return nil, nil
	}
	return r.byRef[ref], nil
}
func (r *publicTicketBookingsRepo) ListByUser(context.Context, primitive.ObjectID) ([]booking.Booking, error) {
	return nil, nil
}
func (r *publicTicketBookingsRepo) ListConfirmedByUser(context.Context, primitive.ObjectID) ([]booking.Booking, error) {
	return nil, nil
}
func (r *publicTicketBookingsRepo) ListConfirmedByShowtime(context.Context, primitive.ObjectID) ([]booking.Booking, error) {
	return nil, nil
}
func (r *publicTicketBookingsRepo) CountConfirmedBetween(context.Context, time.Time, time.Time) (int, error) {
	return 0, nil
}
func (r *publicTicketBookingsRepo) ListRecentConfirmed(context.Context, int) ([]booking.Booking, error) {
	return nil, nil
}
func (r *publicTicketBookingsRepo) CountConfirmed(context.Context) (int64, error) { return 0, nil }
func (r *publicTicketBookingsRepo) CountConfirmedFiltered(context.Context, booking.ConfirmedFilter) (int64, error) {
	return 0, nil
}
func (r *publicTicketBookingsRepo) ListConfirmedPage(context.Context, int, int) ([]booking.Booking, error) {
	return nil, nil
}
func (r *publicTicketBookingsRepo) ListConfirmedFiltered(context.Context, booking.ConfirmedFilter, int, int) ([]booking.Booking, error) {
	return nil, nil
}
func (r *publicTicketBookingsRepo) UpdateTicketToken(context.Context, primitive.ObjectID, string) error {
	return nil
}

type publicTicketMovies struct{ movie *catalog.Movie }

func (m publicTicketMovies) InsertMovie(context.Context, *catalog.Movie) error { return nil }
func (m publicTicketMovies) FindMovieByID(_ context.Context, id primitive.ObjectID) (*catalog.Movie, error) {
	if m.movie != nil && m.movie.ID == id {
		return m.movie, nil
	}
	return nil, nil
}
func (publicTicketMovies) ListMoviesByStatus(context.Context, string) ([]catalog.Movie, error) {
	return nil, nil
}
func (publicTicketMovies) ListComingSoonMovies(context.Context) ([]catalog.Movie, error) { return nil, nil }
func (publicTicketMovies) ListNonArchivedMovies(context.Context) ([]catalog.Movie, error) {
	return nil, nil
}
func (publicTicketMovies) ListMovies(context.Context) ([]catalog.Movie, error) { return nil, nil }
func (publicTicketMovies) UpdateMovie(context.Context, *catalog.Movie) error   { return nil }
func (publicTicketMovies) DeleteMovie(context.Context, primitive.ObjectID) error { return nil }

type publicTicketShowtimes struct{ showtime *catalog.Showtime }

func (s publicTicketShowtimes) InsertShowtime(context.Context, *catalog.Showtime) error { return nil }
func (s publicTicketShowtimes) FindShowtimeByID(_ context.Context, id primitive.ObjectID) (*catalog.Showtime, error) {
	if s.showtime != nil && s.showtime.ID == id {
		return s.showtime, nil
	}
	return nil, nil
}
func (publicTicketShowtimes) ListShowtimesByScreen(context.Context, primitive.ObjectID, time.Time) ([]catalog.Showtime, error) {
	return nil, nil
}
func (publicTicketShowtimes) ListShowtimesByMovie(context.Context, primitive.ObjectID) ([]catalog.Showtime, error) {
	return nil, nil
}
func (publicTicketShowtimes) ListShowtimesByScreens(context.Context, []primitive.ObjectID) ([]catalog.Showtime, error) {
	return nil, nil
}
func (publicTicketShowtimes) ListShowtimesByCinemaMovie(context.Context, []primitive.ObjectID, primitive.ObjectID) ([]catalog.Showtime, error) {
	return nil, nil
}
func (publicTicketShowtimes) ListAdminShowtimes(context.Context, catalog.AdminShowtimeFilter) ([]catalog.Showtime, error) {
	return nil, nil
}
func (publicTicketShowtimes) UpdateShowtime(context.Context, *catalog.Showtime) error { return nil }
func (publicTicketShowtimes) DeleteShowtime(context.Context, primitive.ObjectID) error  { return nil }

type publicTicketScreens struct{ screen *catalog.Screen }

func (s publicTicketScreens) InsertScreen(context.Context, *catalog.Screen) error { return nil }
func (s publicTicketScreens) FindScreenByID(_ context.Context, id primitive.ObjectID) (*catalog.Screen, error) {
	if s.screen != nil && s.screen.ID == id {
		return s.screen, nil
	}
	return nil, nil
}
func (publicTicketScreens) ListScreensByCinema(context.Context, primitive.ObjectID) ([]catalog.Screen, error) {
	return nil, nil
}
func (publicTicketScreens) ListScreens(context.Context, *primitive.ObjectID) ([]catalog.Screen, error) {
	return nil, nil
}
func (publicTicketScreens) UpdateScreen(context.Context, *catalog.Screen) error { return nil }
func (publicTicketScreens) DeleteScreen(context.Context, primitive.ObjectID) error { return nil }

type publicTicketCinemas struct{ cinema *catalog.Cinema }

func (c publicTicketCinemas) InsertCinema(context.Context, *catalog.Cinema) error { return nil }
func (c publicTicketCinemas) FindCinemaByID(_ context.Context, id primitive.ObjectID) (*catalog.Cinema, error) {
	if c.cinema != nil && c.cinema.ID == id {
		return c.cinema, nil
	}
	return nil, nil
}
func (publicTicketCinemas) ListCinemas(context.Context) ([]catalog.Cinema, error) { return nil, nil }
func (publicTicketCinemas) UpdateCinema(context.Context, *catalog.Cinema) error  { return nil }
func (publicTicketCinemas) DeleteCinema(context.Context, primitive.ObjectID) error { return nil }

func setupPublicTicketRouter(t *testing.T, repo *publicTicketBookingsRepo) (*gin.Engine, primitive.ObjectID) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	showtimeID := primitive.NewObjectID()
	screenID := primitive.NewObjectID()
	cinemaID := primitive.NewObjectID()
	movieID := primitive.NewObjectID()

	svc := booking.NewService(
		publicTicketShowtimes{showtime: &catalog.Showtime{
			ID:       showtimeID,
			ScreenID: screenID,
			MovieID:  movieID,
			StartsAt: time.Date(2026, 6, 12, 19, 0, 0, 0, time.UTC),
		}},
		publicTicketScreens{screen: &catalog.Screen{ID: screenID, CinemaID: cinemaID, Name: "Screen 1"}},
		publicTicketCinemas{cinema: &catalog.Cinema{ID: cinemaID, Name: "Major Cineplex"}},
		publicTicketMovies{movie: &catalog.Movie{ID: movieID, Title: "Test Movie"}},
		repo,
		nil,
		nil,
		nil,
		booking.WithTicketConfig(publicTicketSecret, "http://localhost:5173"),
	)

	r := gin.New()
	r.GET("/api/tickets/:ref", handler.GetPublicTicket(handler.BookingsDeps{Bookings: svc}))
	return r, showtimeID
}

func TestGetPublicTicketValidToken(t *testing.T) {
	t.Parallel()

	bookingID := primitive.NewObjectID()
	ref := "TBS-PUBLIC01"
	token := booking.SignTicketToken(publicTicketSecret, ref, bookingID.Hex())

	repo := &publicTicketBookingsRepo{byRef: map[string]*booking.Booking{}}
	r, showtimeID := setupPublicTicketRouter(t, repo)
	repo.byRef[ref] = &booking.Booking{
		ID:          bookingID,
		ShowtimeID:  showtimeID,
		BookingRef:  ref,
		TicketToken: token,
		Status:      booking.StatusConfirmed,
		Seats:       []string{"A1"},
		Total:       35000,
	}

	req := httptest.NewRequest(http.MethodGet, "/api/tickets/"+ref+"?t="+token, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body %s", w.Code, w.Body.String())
	}
	body := w.Body.String()
	for _, want := range []string{ref, "Test Movie", "Major Cineplex", "qrPngBase64"} {
		if !strings.Contains(body, want) {
			t.Fatalf("body = %s, want substring %q", body, want)
		}
	}
}

func TestGetPublicTicketRejectsInvalidToken(t *testing.T) {
	t.Parallel()

	bookingID := primitive.NewObjectID()
	ref := "TBS-PUBLIC02"
	token := booking.SignTicketToken(publicTicketSecret, ref, bookingID.Hex())

	repo := &publicTicketBookingsRepo{byRef: map[string]*booking.Booking{
		ref: {
			ID:          bookingID,
			BookingRef:  ref,
			TicketToken: token,
			Status:      booking.StatusConfirmed,
		},
	}}
	r, _ := setupPublicTicketRouter(t, repo)

	req := httptest.NewRequest(http.MethodGet, "/api/tickets/"+ref+"?t="+token+"x", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
	if !strings.Contains(w.Body.String(), "INVALID_TICKET") {
		t.Fatalf("body = %s, want INVALID_TICKET", w.Body.String())
	}
}

func TestGetPublicTicketRejectsMissingToken(t *testing.T) {
	t.Parallel()

	ref := "TBS-PUBLIC03"
	r, _ := setupPublicTicketRouter(t, &publicTicketBookingsRepo{})

	req := httptest.NewRequest(http.MethodGet, "/api/tickets/"+ref, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
}

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

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/admin"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/auth"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/booking"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/catalog"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/handler"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/user"
)

type adminBookingsRepo struct {
	byRef   *booking.Booking
	allPage []booking.Booking
	allTotal int64
}

func (r *adminBookingsRepo) Insert(_ context.Context, _ *booking.Booking) error { return nil }
func (r *adminBookingsRepo) FindByID(_ context.Context, _ primitive.ObjectID) (*booking.Booking, error) {
	return nil, nil
}
func (r *adminBookingsRepo) FindByBookingRef(_ context.Context, ref string) (*booking.Booking, error) {
	if r.byRef != nil && r.byRef.BookingRef == ref {
		return r.byRef, nil
	}
	return nil, nil
}
func (r *adminBookingsRepo) ListByUser(_ context.Context, _ primitive.ObjectID) ([]booking.Booking, error) {
	return nil, nil
}
func (r *adminBookingsRepo) ListConfirmedByUser(_ context.Context, _ primitive.ObjectID) ([]booking.Booking, error) {
	return nil, nil
}
func (r *adminBookingsRepo) ListConfirmedByShowtime(_ context.Context, _ primitive.ObjectID) ([]booking.Booking, error) {
	return nil, nil
}
func (r *adminBookingsRepo) CountConfirmedBetween(_ context.Context, _, _ time.Time) (int, error) {
	return 0, nil
}
func (r *adminBookingsRepo) ListRecentConfirmed(_ context.Context, _ int) ([]booking.Booking, error) {
	return nil, nil
}
func (r *adminBookingsRepo) CountConfirmed(_ context.Context) (int64, error) {
	return r.allTotal, nil
}
func (r *adminBookingsRepo) ListConfirmedPage(_ context.Context, _, _ int) ([]booking.Booking, error) {
	return r.allPage, nil
}

type adminBookingsShowtimes struct{}

func (r *adminBookingsShowtimes) InsertShowtime(_ context.Context, _ *catalog.Showtime) error {
	return nil
}
func (r *adminBookingsShowtimes) FindShowtimeByID(_ context.Context, id primitive.ObjectID) (*catalog.Showtime, error) {
	return &catalog.Showtime{ID: id, MovieID: primitive.NewObjectID()}, nil
}
func (r *adminBookingsShowtimes) ListShowtimesByScreen(_ context.Context, _ primitive.ObjectID, _ time.Time) ([]catalog.Showtime, error) {
	return nil, nil
}
func (r *adminBookingsShowtimes) ListShowtimesByMovie(_ context.Context, _ primitive.ObjectID) ([]catalog.Showtime, error) {
	return nil, nil
}
func (r *adminBookingsShowtimes) ListShowtimesByScreens(_ context.Context, _ []primitive.ObjectID) ([]catalog.Showtime, error) {
	return nil, nil
}
func (r *adminBookingsShowtimes) ListShowtimesByCinemaMovie(_ context.Context, _ []primitive.ObjectID, _ primitive.ObjectID) ([]catalog.Showtime, error) {
	return nil, nil
}
func (r *adminBookingsShowtimes) ListAdminShowtimes(_ context.Context, _ catalog.AdminShowtimeFilter) ([]catalog.Showtime, error) {
	return nil, nil
}
func (r *adminBookingsShowtimes) UpdateShowtime(_ context.Context, _ *catalog.Showtime) error { return nil }
func (r *adminBookingsShowtimes) DeleteShowtime(_ context.Context, _ primitive.ObjectID) error { return nil }

type adminBookingsMovies struct{}

func (r *adminBookingsMovies) InsertMovie(_ context.Context, _ *catalog.Movie) error { return nil }
func (r *adminBookingsMovies) FindMovieByID(_ context.Context, _ primitive.ObjectID) (*catalog.Movie, error) {
	return &catalog.Movie{Title: "Test Movie"}, nil
}
func (r *adminBookingsMovies) ListMoviesByStatus(_ context.Context, _ string) ([]catalog.Movie, error) {
	return nil, nil
}
func (r *adminBookingsMovies) ListComingSoonMovies(_ context.Context) ([]catalog.Movie, error) {
	return nil, nil
}
func (r *adminBookingsMovies) ListNonArchivedMovies(_ context.Context) ([]catalog.Movie, error) {
	return nil, nil
}
func (r *adminBookingsMovies) ListMovies(_ context.Context) ([]catalog.Movie, error) { return nil, nil }
func (r *adminBookingsMovies) UpdateMovie(_ context.Context, _ *catalog.Movie) error { return nil }
func (r *adminBookingsMovies) DeleteMovie(_ context.Context, _ primitive.ObjectID) error {
	return nil
}

type adminBookingsUsers struct{}

func (r *adminBookingsUsers) Insert(_ context.Context, _ *user.User) error { return nil }
func (r *adminBookingsUsers) FindByID(_ context.Context, _ primitive.ObjectID) (*user.User, error) {
	return &user.User{Email: "customer@example.com"}, nil
}
func (r *adminBookingsUsers) FindByEmail(_ context.Context, _ string) (*user.User, error) {
	return nil, nil
}
func (r *adminBookingsUsers) FindByGoogleID(_ context.Context, _ string) (*user.User, error) {
	return nil, nil
}
func (r *adminBookingsUsers) Update(_ context.Context, _ *user.User) error { return nil }

func setupAdminBookingsRouter(t *testing.T, role string, repo *adminBookingsRepo) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	svc := &admin.BookingsService{
		Bookings:  repo,
		Showtimes: &adminBookingsShowtimes{},
		Movies:    &adminBookingsMovies{},
		Users:     &adminBookingsUsers{},
	}
	deps := handler.AdminBookingsDeps{Service: svc}

	tokens := auth.NewTokenService("test-secret", time.Hour)
	authSvc := auth.NewService(nil, tokens, auth.NewLoginRateLimiter(nil), "")
	authMw := auth.MiddlewareDeps{Service: authSvc}

	userID := primitive.NewObjectID()
	token, _, err := tokens.Issue(userID, role)
	if err != nil {
		t.Fatal(err)
	}

	r := gin.New()
	adminGroup := r.Group("/api/admin")
	adminGroup.Use(func(c *gin.Context) {
		c.Request.AddCookie(&http.Cookie{Name: auth.CookieName, Value: token})
		c.Next()
	})
	adminGroup.Use(auth.RequireAuth(authMw), auth.RequireAdmin(authMw))
	adminGroup.GET("/bookings", handler.SearchAdminBookings(deps))

	return r
}

func TestAdminBookingsRejectsCustomer(t *testing.T) {
	r := setupAdminBookingsRouter(t, user.RoleCustomer, &adminBookingsRepo{})

	req := httptest.NewRequest(http.MethodGet, "/api/admin/bookings?bookingRef=ABC123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", w.Code)
	}
}

func TestAdminBookingsListAllPaginated(t *testing.T) {
	repo := &adminBookingsRepo{
		allPage: []booking.Booking{
			{
				ID:          primitive.NewObjectID(),
				UserID:      primitive.NewObjectID(),
				ShowtimeID:  primitive.NewObjectID(),
				BookingRef:  "BK-PAGE-001",
				Status:      booking.StatusConfirmed,
				Seats:       []string{"A1"},
				Total:       1200,
				ConfirmedAt: time.Now().UTC(),
			},
		},
		allTotal: 1,
	}
	r := setupAdminBookingsRouter(t, user.RoleAdmin, repo)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/bookings?page=1&limit=20", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body %s", w.Code, w.Body.String())
	}
	body := w.Body.String()
	if !strings.Contains(body, `"total":1`) || !strings.Contains(body, "BK-PAGE-001") {
		t.Fatalf("body = %s, want paginated booking list", body)
	}
}

func TestAdminBookingsSearchByRef(t *testing.T) {
	ref := "BK-TEST-001"
	repo := &adminBookingsRepo{
		byRef: &booking.Booking{
			ID:          primitive.NewObjectID(),
			UserID:      primitive.NewObjectID(),
			ShowtimeID:  primitive.NewObjectID(),
			BookingRef:  ref,
			Status:      booking.StatusConfirmed,
			Seats:       []string{"A1"},
			Total:       1200,
			ConfirmedAt: time.Now().UTC(),
		},
	}
	r := setupAdminBookingsRouter(t, user.RoleAdmin, repo)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/bookings?bookingRef="+ref, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body %s", w.Code, w.Body.String())
	}
	body := w.Body.String()
	if !strings.Contains(body, ref) || !strings.Contains(body, "Test Movie") {
		t.Fatalf("body = %s, want booking ref and movie title", body)
	}
}

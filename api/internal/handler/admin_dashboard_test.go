package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
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

type memBookingRepo struct{}

func (r *memBookingRepo) Insert(_ context.Context, _ *booking.Booking) error { return nil }
func (r *memBookingRepo) FindByID(_ context.Context, _ primitive.ObjectID) (*booking.Booking, error) {
	return nil, nil
}
func (r *memBookingRepo) FindByBookingRef(_ context.Context, _ string) (*booking.Booking, error) {
	return nil, nil
}
func (r *memBookingRepo) ListByUser(_ context.Context, _ primitive.ObjectID) ([]booking.Booking, error) {
	return nil, nil
}
func (r *memBookingRepo) ListConfirmedByShowtime(_ context.Context, _ primitive.ObjectID) ([]booking.Booking, error) {
	return nil, nil
}
func (r *memBookingRepo) CountConfirmedBetween(_ context.Context, _, _ time.Time) (int, error) {
	return 0, nil
}
func (r *memBookingRepo) ListRecentConfirmed(_ context.Context, _ int) ([]booking.Booking, error) {
	return []booking.Booking{}, nil
}

type memShowtimeRepo struct{}

func (r *memShowtimeRepo) InsertShowtime(_ context.Context, _ *catalog.Showtime) error { return nil }
func (r *memShowtimeRepo) FindShowtimeByID(_ context.Context, _ primitive.ObjectID) (*catalog.Showtime, error) {
	return nil, nil
}
func (r *memShowtimeRepo) ListShowtimesByScreen(_ context.Context, _ primitive.ObjectID, _ time.Time) ([]catalog.Showtime, error) {
	return nil, nil
}
func (r *memShowtimeRepo) ListShowtimesByMovie(_ context.Context, _ primitive.ObjectID) ([]catalog.Showtime, error) {
	return nil, nil
}
func (r *memShowtimeRepo) ListAdminShowtimes(_ context.Context, _ catalog.AdminShowtimeFilter) ([]catalog.Showtime, error) {
	return nil, nil
}
func (r *memShowtimeRepo) ListShowtimesByScreens(_ context.Context, _ []primitive.ObjectID) ([]catalog.Showtime, error) {
	return nil, nil
}
func (r *memShowtimeRepo) ListShowtimesByCinemaMovie(_ context.Context, _ []primitive.ObjectID, _ primitive.ObjectID) ([]catalog.Showtime, error) {
	return nil, nil
}
func (r *memShowtimeRepo) UpdateShowtime(_ context.Context, _ *catalog.Showtime) error { return nil }
func (r *memShowtimeRepo) DeleteShowtime(_ context.Context, _ primitive.ObjectID) error { return nil }

type memScreenRepo struct{}

func (r *memScreenRepo) InsertScreen(_ context.Context, _ *catalog.Screen) error { return nil }
func (r *memScreenRepo) FindScreenByID(_ context.Context, _ primitive.ObjectID) (*catalog.Screen, error) {
	return nil, nil
}
func (r *memScreenRepo) ListScreensByCinema(_ context.Context, _ primitive.ObjectID) ([]catalog.Screen, error) {
	return nil, nil
}
func (r *memScreenRepo) ListScreens(_ context.Context, _ *primitive.ObjectID) ([]catalog.Screen, error) {
	return nil, nil
}
func (r *memScreenRepo) UpdateScreen(_ context.Context, _ *catalog.Screen) error { return nil }
func (r *memScreenRepo) DeleteScreen(_ context.Context, _ primitive.ObjectID) error { return nil }

type memMovieRepo struct{}

func (r *memMovieRepo) InsertMovie(_ context.Context, _ *catalog.Movie) error { return nil }
func (r *memMovieRepo) FindMovieByID(_ context.Context, _ primitive.ObjectID) (*catalog.Movie, error) {
	return nil, nil
}
func (r *memMovieRepo) ListMoviesByStatus(_ context.Context, _ string) ([]catalog.Movie, error) {
	return nil, nil
}
func (r *memMovieRepo) ListMovies(_ context.Context) ([]catalog.Movie, error) { return nil, nil }
func (r *memMovieRepo) ListComingSoonMovies(_ context.Context) ([]catalog.Movie, error) { return nil, nil }
func (r *memMovieRepo) ListNonArchivedMovies(_ context.Context) ([]catalog.Movie, error) { return nil, nil }
func (r *memMovieRepo) UpdateMovie(_ context.Context, _ *catalog.Movie) error { return nil }
func (r *memMovieRepo) DeleteMovie(_ context.Context, _ primitive.ObjectID) error { return nil }

func setupDashboardRouter(t *testing.T, role string) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	dashboardSvc := &admin.DashboardService{
		Showtimes: &memShowtimeRepo{},
		Screens:   &memScreenRepo{},
		Movies:    &memMovieRepo{},
		Bookings:  &memBookingRepo{},
	}
	deps := handler.AdminDashboardDeps{Service: dashboardSvc}

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
	adminGroup.GET("/dashboard", handler.GetAdminDashboard(deps))

	return r
}

func TestAdminDashboardRejectsCustomer(t *testing.T) {
	r := setupDashboardRouter(t, user.RoleCustomer)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/dashboard", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", w.Code)
	}
}

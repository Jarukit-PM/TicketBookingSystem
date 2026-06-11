package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/audit"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/auth"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/catalog"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/handler"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/user"
)

type memAuditRepo struct {
	logs []audit.AuditLog
}

func (r *memAuditRepo) InsertAuditLog(_ context.Context, log *audit.AuditLog) error {
	log.ID = primitive.NewObjectID()
	r.logs = append(r.logs, *log)
	return nil
}

func (r *memAuditRepo) ListAuditLogs(_ context.Context, _ int64) ([]audit.AuditLog, error) {
	return r.logs, nil
}

type memCatalogRepos struct {
	cinemas map[primitive.ObjectID]*catalog.Cinema
	movies  map[primitive.ObjectID]*catalog.Movie
}

func newMemCatalogRepos() *memCatalogRepos {
	return &memCatalogRepos{
		cinemas: make(map[primitive.ObjectID]*catalog.Cinema),
		movies:  make(map[primitive.ObjectID]*catalog.Movie),
	}
}

func (r *memCatalogRepos) asMongoRepositories() catalog.MongoRepositories {
	return catalog.MongoRepositories{
		Cinemas:   r,
		Screens:   r,
		Movies:    r,
		Showtimes: r,
	}
}

func (r *memCatalogRepos) InsertCinema(_ context.Context, cinema *catalog.Cinema) error {
	cinema.ID = primitive.NewObjectID()
	r.cinemas[cinema.ID] = cinema
	return nil
}

func (r *memCatalogRepos) FindCinemaByID(_ context.Context, id primitive.ObjectID) (*catalog.Cinema, error) {
	c, ok := r.cinemas[id]
	if !ok {
		return nil, nil
	}
	return c, nil
}

func (r *memCatalogRepos) ListCinemas(_ context.Context) ([]catalog.Cinema, error) {
	out := make([]catalog.Cinema, 0, len(r.cinemas))
	for _, c := range r.cinemas {
		out = append(out, *c)
	}
	return out, nil
}

func (r *memCatalogRepos) UpdateCinema(_ context.Context, cinema *catalog.Cinema) error {
	r.cinemas[cinema.ID] = cinema
	return nil
}

func (r *memCatalogRepos) DeleteCinema(_ context.Context, id primitive.ObjectID) error {
	delete(r.cinemas, id)
	return nil
}

func (r *memCatalogRepos) InsertScreen(_ context.Context, screen *catalog.Screen) error { return nil }
func (r *memCatalogRepos) FindScreenByID(_ context.Context, _ primitive.ObjectID) (*catalog.Screen, error) {
	return nil, nil
}
func (r *memCatalogRepos) ListScreensByCinema(_ context.Context, _ primitive.ObjectID) ([]catalog.Screen, error) {
	return nil, nil
}
func (r *memCatalogRepos) ListScreens(_ context.Context, _ *primitive.ObjectID) ([]catalog.Screen, error) {
	return nil, nil
}
func (r *memCatalogRepos) UpdateScreen(_ context.Context, _ *catalog.Screen) error { return nil }
func (r *memCatalogRepos) DeleteScreen(_ context.Context, _ primitive.ObjectID) error { return nil }

func (r *memCatalogRepos) InsertMovie(_ context.Context, movie *catalog.Movie) error {
	movie.ID = primitive.NewObjectID()
	r.movies[movie.ID] = movie
	return nil
}

func (r *memCatalogRepos) FindMovieByID(_ context.Context, id primitive.ObjectID) (*catalog.Movie, error) {
	m, ok := r.movies[id]
	if !ok {
		return nil, nil
	}
	return m, nil
}

func (r *memCatalogRepos) ListMoviesByStatus(_ context.Context, _ string) ([]catalog.Movie, error) {
	return nil, nil
}

func (r *memCatalogRepos) ListMovies(_ context.Context) ([]catalog.Movie, error) {
	out := make([]catalog.Movie, 0, len(r.movies))
	for _, m := range r.movies {
		out = append(out, *m)
	}
	return out, nil
}

func (r *memCatalogRepos) ListComingSoonMovies(ctx context.Context) ([]catalog.Movie, error) {
	return r.ListMoviesByStatus(ctx, catalog.MovieStatusComingSoon)
}

func (r *memCatalogRepos) ListNonArchivedMovies(_ context.Context) ([]catalog.Movie, error) {
	return nil, nil
}

func (r *memCatalogRepos) UpdateMovie(_ context.Context, movie *catalog.Movie) error {
	r.movies[movie.ID] = movie
	return nil
}

func (r *memCatalogRepos) DeleteMovie(_ context.Context, id primitive.ObjectID) error {
	delete(r.movies, id)
	return nil
}

func (r *memCatalogRepos) InsertShowtime(_ context.Context, _ *catalog.Showtime) error { return nil }
func (r *memCatalogRepos) FindShowtimeByID(_ context.Context, _ primitive.ObjectID) (*catalog.Showtime, error) {
	return nil, nil
}
func (r *memCatalogRepos) ListShowtimesByScreen(_ context.Context, _ primitive.ObjectID, _ time.Time) ([]catalog.Showtime, error) {
	return nil, nil
}
func (r *memCatalogRepos) ListShowtimesByMovie(_ context.Context, _ primitive.ObjectID) ([]catalog.Showtime, error) {
	return nil, nil
}
func (r *memCatalogRepos) ListAdminShowtimes(_ context.Context, _ catalog.AdminShowtimeFilter) ([]catalog.Showtime, error) {
	return nil, nil
}
func (r *memCatalogRepos) ListShowtimesByScreens(_ context.Context, _ []primitive.ObjectID) ([]catalog.Showtime, error) {
	return nil, nil
}
func (r *memCatalogRepos) ListShowtimesByCinemaMovie(_ context.Context, _ []primitive.ObjectID, _ primitive.ObjectID) ([]catalog.Showtime, error) {
	return nil, nil
}
func (r *memCatalogRepos) UpdateShowtime(_ context.Context, _ *catalog.Showtime) error { return nil }
func (r *memCatalogRepos) DeleteShowtime(_ context.Context, _ primitive.ObjectID) error { return nil }

func setupAdminRouter(t *testing.T, role string) (*gin.Engine, *memAuditRepo) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	repos := newMemCatalogRepos()
	auditRepo := &memAuditRepo{}
	svc := catalog.NewAdminService(repos.asMongoRepositories(), auditRepo)
	deps := handler.AdminCatalogDeps{Service: svc}

	tokens := auth.NewTokenService("test-secret", time.Hour)
	authSvc := auth.NewService(nil, tokens, auth.NewLoginRateLimiter(nil), "")
	authMw := auth.MiddlewareDeps{Service: authSvc}

	userID := primitive.NewObjectID()
	token, _, err := tokens.Issue(userID, role)
	if err != nil {
		t.Fatal(err)
	}

	r := gin.New()
	admin := r.Group("/api/admin")
	admin.Use(func(c *gin.Context) {
		c.Request.AddCookie(&http.Cookie{Name: auth.CookieName, Value: token})
		c.Next()
	})
	admin.Use(auth.RequireAuth(authMw), auth.RequireAdmin(authMw))
	admin.POST("/movies", handler.CreateAdminMovie(deps))
	admin.GET("/movies", handler.ListAdminMovies(deps))

	return r, auditRepo
}

func TestAdminMoviesRejectsCustomer(t *testing.T) {
	r, _ := setupAdminRouter(t, user.RoleCustomer)

	body, _ := json.Marshal(map[string]any{
		"title":  "Blocked",
		"status": catalog.MovieStatusNowShowing,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/admin/movies", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", w.Code)
	}
}

func TestAdminCreateMovieWritesAuditLog(t *testing.T) {
	r, auditRepo := setupAdminRouter(t, user.RoleAdmin)

	body, _ := json.Marshal(map[string]any{
		"title":  "Admin Film",
		"status": catalog.MovieStatusNowShowing,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/admin/movies", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body=%s", w.Code, w.Body.String())
	}
	if len(auditRepo.logs) != 1 {
		t.Fatalf("audit logs = %d, want 1", len(auditRepo.logs))
	}
	if auditRepo.logs[0].Action != catalog.AuditActionCreate {
		t.Fatalf("action = %q, want create", auditRepo.logs[0].Action)
	}
}

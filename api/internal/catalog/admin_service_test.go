package catalog_test

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/audit"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/catalog"
)

type memAuditRepo struct {
	logs []audit.AuditLog
}

func (r *memAuditRepo) InsertAuditLog(_ context.Context, log *audit.AuditLog) error {
	log.ID = primitive.NewObjectID()
	r.logs = append(r.logs, *log)
	return nil
}

func (r *memAuditRepo) ListAuditLogs(_ context.Context, _ audit.LogPage) ([]audit.AuditLog, error) {
	return r.logs, nil
}

type memCatalogRepos struct {
	cinemas   map[primitive.ObjectID]*catalog.Cinema
	screens   map[primitive.ObjectID]*catalog.Screen
	movies    map[primitive.ObjectID]*catalog.Movie
	showtimes map[primitive.ObjectID]*catalog.Showtime
}

func newMemCatalogRepos() *memCatalogRepos {
	return &memCatalogRepos{
		cinemas:   make(map[primitive.ObjectID]*catalog.Cinema),
		screens:   make(map[primitive.ObjectID]*catalog.Screen),
		movies:    make(map[primitive.ObjectID]*catalog.Movie),
		showtimes: make(map[primitive.ObjectID]*catalog.Showtime),
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

func (r *memCatalogRepos) InsertScreen(_ context.Context, screen *catalog.Screen) error {
	screen.ID = primitive.NewObjectID()
	r.screens[screen.ID] = screen
	return nil
}

func (r *memCatalogRepos) FindScreenByID(_ context.Context, id primitive.ObjectID) (*catalog.Screen, error) {
	s, ok := r.screens[id]
	if !ok {
		return nil, nil
	}
	return s, nil
}

func (r *memCatalogRepos) ListScreensByCinema(ctx context.Context, cinemaID primitive.ObjectID) ([]catalog.Screen, error) {
	return r.ListScreens(ctx, &cinemaID)
}

func (r *memCatalogRepos) ListScreens(_ context.Context, cinemaID *primitive.ObjectID) ([]catalog.Screen, error) {
	out := make([]catalog.Screen, 0)
	for _, s := range r.screens {
		if cinemaID == nil || s.CinemaID == *cinemaID {
			out = append(out, *s)
		}
	}
	return out, nil
}

func (r *memCatalogRepos) UpdateScreen(_ context.Context, screen *catalog.Screen) error {
	r.screens[screen.ID] = screen
	return nil
}

func (r *memCatalogRepos) DeleteScreen(_ context.Context, id primitive.ObjectID) error {
	delete(r.screens, id)
	return nil
}

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

func (r *memCatalogRepos) ListMoviesByStatus(_ context.Context, status string) ([]catalog.Movie, error) {
	out := make([]catalog.Movie, 0)
	for _, m := range r.movies {
		if m.Status == status {
			out = append(out, *m)
		}
	}
	return out, nil
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
	out := make([]catalog.Movie, 0)
	for _, m := range r.movies {
		if m.Status != catalog.MovieStatusArchived {
			out = append(out, *m)
		}
	}
	return out, nil
}

func (r *memCatalogRepos) UpdateMovie(_ context.Context, movie *catalog.Movie) error {
	r.movies[movie.ID] = movie
	return nil
}

func (r *memCatalogRepos) DeleteMovie(_ context.Context, id primitive.ObjectID) error {
	delete(r.movies, id)
	return nil
}

func (r *memCatalogRepos) InsertShowtime(_ context.Context, showtime *catalog.Showtime) error {
	showtime.ID = primitive.NewObjectID()
	r.showtimes[showtime.ID] = showtime
	return nil
}

func (r *memCatalogRepos) FindShowtimeByID(_ context.Context, id primitive.ObjectID) (*catalog.Showtime, error) {
	s, ok := r.showtimes[id]
	if !ok {
		return nil, nil
	}
	return s, nil
}

func (r *memCatalogRepos) ListShowtimesByScreen(_ context.Context, screenID primitive.ObjectID, from time.Time) ([]catalog.Showtime, error) {
	out := make([]catalog.Showtime, 0)
	for _, s := range r.showtimes {
		if s.ScreenID == screenID && !s.StartsAt.Before(from) {
			out = append(out, *s)
		}
	}
	return out, nil
}

func (r *memCatalogRepos) ListShowtimesByMovie(_ context.Context, movieID primitive.ObjectID) ([]catalog.Showtime, error) {
	out := make([]catalog.Showtime, 0)
	for _, s := range r.showtimes {
		if s.MovieID == movieID {
			out = append(out, *s)
		}
	}
	return out, nil
}

func (r *memCatalogRepos) ListAdminShowtimes(_ context.Context, filter catalog.AdminShowtimeFilter) ([]catalog.Showtime, error) {
	out := make([]catalog.Showtime, 0)
	for _, s := range r.showtimes {
		if filter.MovieID != nil && s.MovieID != *filter.MovieID {
			continue
		}
		if filter.ScreenID != nil && s.ScreenID != *filter.ScreenID {
			continue
		}
		if filter.From != nil && s.StartsAt.Before(*filter.From) {
			continue
		}
		if filter.To != nil && s.StartsAt.After(*filter.To) {
			continue
		}
		if filter.CinemaID != nil {
			screen, ok := r.screens[s.ScreenID]
			if !ok || screen.CinemaID != *filter.CinemaID {
				continue
			}
		}
		out = append(out, *s)
	}
	return out, nil
}

func (r *memCatalogRepos) ListShowtimesByScreens(_ context.Context, screenIDs []primitive.ObjectID) ([]catalog.Showtime, error) {
	allowed := make(map[primitive.ObjectID]struct{}, len(screenIDs))
	for _, id := range screenIDs {
		allowed[id] = struct{}{}
	}
	out := make([]catalog.Showtime, 0)
	for _, s := range r.showtimes {
		if _, ok := allowed[s.ScreenID]; ok {
			out = append(out, *s)
		}
	}
	return out, nil
}

func (r *memCatalogRepos) ListShowtimesByCinemaMovie(
	_ context.Context,
	screenIDs []primitive.ObjectID,
	movieID primitive.ObjectID,
) ([]catalog.Showtime, error) {
	allowed := make(map[primitive.ObjectID]struct{}, len(screenIDs))
	for _, id := range screenIDs {
		allowed[id] = struct{}{}
	}
	out := make([]catalog.Showtime, 0)
	for _, s := range r.showtimes {
		if s.MovieID == movieID {
			if _, ok := allowed[s.ScreenID]; ok {
				out = append(out, *s)
			}
		}
	}
	return out, nil
}

func (r *memCatalogRepos) UpdateShowtime(_ context.Context, showtime *catalog.Showtime) error {
	r.showtimes[showtime.ID] = showtime
	return nil
}

func (r *memCatalogRepos) DeleteShowtime(_ context.Context, id primitive.ObjectID) error {
	delete(r.showtimes, id)
	return nil
}

func TestAdminCreateChainWritesAuditLogs(t *testing.T) {
	repos := newMemCatalogRepos()
	auditRepo := &memAuditRepo{}
	svc := catalog.NewAdminService(repos.asMongoRepositories(), auditRepo)
	ctx := context.Background()
	actorID := primitive.NewObjectID()

	cinema, err := svc.CreateCinema(ctx, actorID, &catalog.Cinema{
		Name:     "Downtown",
		Address:  "1 Main St",
		Timezone: "Asia/Bangkok",
	})
	if err != nil {
		t.Fatal(err)
	}

	screen, err := svc.CreateScreen(ctx, actorID, &catalog.Screen{
		CinemaID: cinema.ID,
		Name:     "Hall 1",
		Layout: catalog.ScreenLayout{
			Seats: []catalog.LayoutSeat{
				{SeatID: "A-1", Row: 1, Col: 1, Type: catalog.SeatTypeStandard},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	movie, err := svc.CreateMovie(ctx, actorID, &catalog.Movie{
		Title:       "Test Film",
		PosterURL:   "https://example.com/poster.jpg",
		DurationMin: 120,
		Rating:      "PG",
		Synopsis:    "A test movie",
		Status:      catalog.MovieStatusNowShowing,
	})
	if err != nil {
		t.Fatal(err)
	}

	showtime, err := svc.CreateShowtime(ctx, actorID, &catalog.Showtime{
		MovieID:  movie.ID,
		ScreenID: screen.ID,
		StartsAt: time.Now().UTC().Add(24 * time.Hour),
		PriceTiers: catalog.PriceTiers{
			Standard:   1000,
			VIP:        1500,
			Wheelchair: 1000,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(auditRepo.logs) != 4 {
		t.Fatalf("audit logs = %d, want 4", len(auditRepo.logs))
	}
	if auditRepo.logs[0].Action != catalog.AuditActionCreate || auditRepo.logs[0].Entity != "cinema" {
		t.Fatalf("unexpected first audit log: %+v", auditRepo.logs[0])
	}
	if showtime.ID.IsZero() {
		t.Fatal("expected showtime id")
	}
}

func TestAdminDeleteMovieWritesAuditLog(t *testing.T) {
	repos := newMemCatalogRepos()
	auditRepo := &memAuditRepo{}
	svc := catalog.NewAdminService(repos.asMongoRepositories(), auditRepo)
	ctx := context.Background()
	actorID := primitive.NewObjectID()

	movie, err := svc.CreateMovie(ctx, actorID, &catalog.Movie{
		Title:  "Delete Me",
		Status: catalog.MovieStatusNowShowing,
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := svc.DeleteMovie(ctx, actorID, movie.ID); err != nil {
		t.Fatal(err)
	}

	last := auditRepo.logs[len(auditRepo.logs)-1]
	if last.Action != catalog.AuditActionDelete || last.Entity != "movie" {
		t.Fatalf("unexpected delete audit log: %+v", last)
	}
}

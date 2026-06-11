package main

import (
	"context"
	"log"
	"time"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/auth"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/catalog"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/config"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/db"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/user"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	cfg := config.MustLoad()
	ctx := context.Background()

	client, err := db.ConnectMongo(ctx, cfg.MongoURI)
	if err != nil {
		log.Fatalf("mongo connect: %v", err)
	}
	defer func() {
		_ = db.DisconnectMongo(context.Background(), client)
	}()

	database := db.Database(client, cfg.MongoURI)
	if err := db.EnsureIndexes(ctx, database); err != nil {
		log.Fatalf("ensure indexes: %v", err)
	}

	if err := seedAdmin(ctx, cfg, database); err != nil {
		log.Fatalf("seed admin: %v", err)
	}

	repos := catalog.NewMongoRepositories(database)

	existing, err := repos.Cinemas.ListCinemas(ctx)
	if err != nil {
		log.Fatalf("list cinemas: %v", err)
	}
	if len(existing) > 0 {
		log.Println("seed: catalog data already present, skipping")
		return
	}

	cinema := &catalog.Cinema{
		Name:     "TBS Seed Cinema",
		Address:  "123 Main Street, Bangkok",
		Timezone: "Asia/Bangkok",
	}
	if err := repos.Cinemas.InsertCinema(ctx, cinema); err != nil {
		log.Fatalf("insert cinema: %v", err)
	}

	screenA := &catalog.Screen{
		CinemaID: cinema.ID,
		Name:     "Hall A",
		Layout:   sampleLayout("A"),
	}
	screenB := &catalog.Screen{
		CinemaID: cinema.ID,
		Name:     "Hall B",
		Layout:   sampleLayout("B"),
	}
	for _, s := range []*catalog.Screen{screenA, screenB} {
		if err := repos.Screens.InsertScreen(ctx, s); err != nil {
			log.Fatalf("insert screen %s: %v", s.Name, err)
		}
	}

	movie1 := &catalog.Movie{
		Title:       "Neon Horizon",
		PosterURL:   "https://example.com/posters/neon-horizon.jpg",
		DurationMin: 128,
		Rating:      "PG-13",
		Synopsis:    "A sci-fi thriller set in a rain-soaked megacity.",
		Status:      catalog.MovieStatusNowShowing,
	}
	movie2 := &catalog.Movie{
		Title:       "Midnight Express",
		PosterURL:   "https://example.com/posters/midnight-express.jpg",
		DurationMin: 102,
		Rating:      "R",
		Synopsis:    "An action drama racing against the clock.",
		Status:      catalog.MovieStatusNowShowing,
	}
	for _, m := range []*catalog.Movie{movie1, movie2} {
		if err := repos.Movies.InsertMovie(ctx, m); err != nil {
			log.Fatalf("insert movie %s: %v", m.Title, err)
		}
	}

	defaultTiers := catalog.PriceTiers{
		Standard:   1200,
		VIP:        1800,
		Wheelchair: 1200,
	}

	base := time.Now().UTC().Add(24 * time.Hour).Truncate(time.Hour)
	showtimes := []struct {
		movie  *catalog.Movie
		screen *catalog.Screen
		offset time.Duration
	}{
		{movie: movie1, screen: screenA, offset: 0},
		{movie: movie1, screen: screenB, offset: 3 * time.Hour},
		{movie: movie2, screen: screenA, offset: 6 * time.Hour},
		{movie: movie2, screen: screenB, offset: 9 * time.Hour},
		{movie: movie1, screen: screenA, offset: 30 * time.Hour},
	}

	for _, st := range showtimes {
		showtime := &catalog.Showtime{
			MovieID:    st.movie.ID,
			ScreenID:   st.screen.ID,
			StartsAt:   base.Add(st.offset),
			PriceTiers: defaultTiers,
			Status:     catalog.ShowtimeStatusOpen,
		}
		if err := repos.Showtimes.InsertShowtime(ctx, showtime); err != nil {
			log.Fatalf("insert showtime: %v", err)
		}
	}

	log.Printf("seed: created cinema %s, 2 screens, 2 movies, %d showtimes", cinema.Name, len(showtimes))
}

func seedAdmin(ctx context.Context, cfg config.Config, database *mongo.Database) error {
	if cfg.AdminEmail == "" || cfg.AdminSeedPassword == "" {
		return nil
	}

	userRepo := user.NewMongoRepository(database)
	tokenService := auth.NewTokenService(cfg.JWTSecret, cfg.JWTExpiryDuration())
	authService := auth.NewService(userRepo, tokenService, auth.NewLoginRateLimiter(nil), cfg.AdminEmail)
	return authService.BootstrapAdmin(ctx, cfg.AdminEmail, cfg.AdminSeedPassword)
}

func sampleLayout(rowPrefix string) catalog.ScreenLayout {
	seats := []catalog.LayoutSeat{
		{SeatID: rowPrefix + "-1", Row: 1, Col: 1, Type: catalog.SeatTypeStandard},
		{SeatID: rowPrefix + "-2", Row: 1, Col: 2, Type: catalog.SeatTypeVIP},
		{SeatID: rowPrefix + "-3", Row: 1, Col: 3, Type: catalog.SeatTypeWheelchair},
		{SeatID: rowPrefix + "-4", Row: 2, Col: 1, Type: catalog.SeatTypeStandard},
		{SeatID: rowPrefix + "-5", Row: 2, Col: 2, Type: catalog.SeatTypeStandard},
		{SeatID: rowPrefix + "-X", Row: 2, Col: 3, Type: catalog.SeatTypeBlocked},
	}
	return catalog.ScreenLayout{Seats: seats}
}

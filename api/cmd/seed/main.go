package main

import (
	"context"
	"log"
	"time"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/auth"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/booking"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/catalog"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/config"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/db"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/user"
	"go.mongodb.org/mongo-driver/bson"
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
		log.Println("seed: catalog data already present, skipping catalog")
		if err := seedDashboardDemo(ctx, cfg, database, repos); err != nil {
			log.Fatalf("seed dashboard demo: %v", err)
		}
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

	todayBase := time.Now().UTC().Truncate(time.Hour).Add(2 * time.Hour)
	tomorrowBase := time.Now().UTC().Add(24 * time.Hour).Truncate(time.Hour)
	showtimes := []struct {
		movie   *catalog.Movie
		screen  *catalog.Screen
		startsAt time.Time
	}{
		{movie: movie1, screen: screenA, startsAt: todayBase},
		{movie: movie2, screen: screenB, startsAt: todayBase.Add(4 * time.Hour)},
		{movie: movie1, screen: screenA, startsAt: tomorrowBase},
		{movie: movie1, screen: screenB, startsAt: tomorrowBase.Add(3 * time.Hour)},
		{movie: movie2, screen: screenA, startsAt: tomorrowBase.Add(6 * time.Hour)},
		{movie: movie2, screen: screenB, startsAt: tomorrowBase.Add(9 * time.Hour)},
		{movie: movie1, screen: screenA, startsAt: tomorrowBase.Add(30 * time.Hour)},
	}

	for _, st := range showtimes {
		showtime := &catalog.Showtime{
			MovieID:    st.movie.ID,
			ScreenID:   st.screen.ID,
			StartsAt:   st.startsAt,
			PriceTiers: defaultTiers,
			Status:     catalog.ShowtimeStatusOpen,
		}
		if err := repos.Showtimes.InsertShowtime(ctx, showtime); err != nil {
			log.Fatalf("insert showtime: %v", err)
		}
	}

	log.Printf("seed: created cinema %s, 2 screens, 2 movies, %d showtimes", cinema.Name, len(showtimes))

	if err := seedDashboardDemo(ctx, cfg, database, repos); err != nil {
		log.Fatalf("seed dashboard demo: %v", err)
	}
}

func seedDashboardDemo(ctx context.Context, cfg config.Config, database *mongo.Database, repos catalog.MongoRepositories) error {
	if err := ensureTodayShowtimes(ctx, repos); err != nil {
		return err
	}
	return seedDemoBookingIfEmpty(ctx, cfg, database, repos)
}

func ensureTodayShowtimes(ctx context.Context, repos catalog.MongoRepositories) error {
	now := time.Now().UTC()
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	dayEnd := dayStart.Add(24 * time.Hour)

	existing, err := repos.Showtimes.ListAdminShowtimes(ctx, catalog.AdminShowtimeFilter{
		From: &dayStart,
		To:   &dayEnd,
	})
	if err != nil {
		return err
	}
	if len(existing) >= 2 {
		return nil
	}

	movies, err := repos.Movies.ListMovies(ctx)
	if err != nil {
		return err
	}
	screens, err := repos.Screens.ListScreens(ctx, nil)
	if err != nil {
		return err
	}
	if len(movies) == 0 || len(screens) == 0 {
		log.Println("seed: skip today showtimes — no movies or screens")
		return nil
	}

	defaultTiers := catalog.PriceTiers{Standard: 1200, VIP: 1800, Wheelchair: 1200}
	base := now.Truncate(time.Hour).Add(2 * time.Hour)
	needed := 2 - len(existing)

	for i := 0; i < needed; i++ {
		movie := movies[i%len(movies)]
		screen := screens[i%len(screens)]
		showtime := &catalog.Showtime{
			MovieID:    movie.ID,
			ScreenID:   screen.ID,
			StartsAt:   base.Add(time.Duration(i) * 3 * time.Hour),
			PriceTiers: defaultTiers,
			Status:     catalog.ShowtimeStatusOpen,
		}
		if err := repos.Showtimes.InsertShowtime(ctx, showtime); err != nil {
			return err
		}
		log.Printf("seed: added today showtime %s at %s", movie.Title, showtime.StartsAt.Format(time.RFC3339))
	}
	return nil
}

func seedDemoBookingIfEmpty(ctx context.Context, cfg config.Config, database *mongo.Database, repos catalog.MongoRepositories) error {
	count, err := database.Collection(booking.CollectionName).CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	userRepo := user.NewMongoRepository(database)
	seedUser, err := userRepo.FindByEmail(ctx, cfg.AdminEmail)
	if err != nil {
		return err
	}
	if seedUser == nil {
		log.Println("seed: skip demo booking — no admin user to attach booking")
		return nil
	}

	now := time.Now().UTC()
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	dayEnd := dayStart.Add(24 * time.Hour)
	showtimes, err := repos.Showtimes.ListAdminShowtimes(ctx, catalog.AdminShowtimeFilter{
		From: &dayStart,
		To:   &dayEnd,
	})
	if err != nil {
		return err
	}
	if len(showtimes) == 0 {
		log.Println("seed: skip demo booking — no showtimes today")
		return nil
	}

	showtime := showtimes[0]
	screen, err := repos.Screens.FindScreenByID(ctx, showtime.ScreenID)
	if err != nil {
		return err
	}
	if screen == nil {
		return nil
	}

	seatIDs := []string{screen.Layout.Seats[0].SeatID}
	total, err := catalog.TotalForSeats(showtime.PriceTiers, screen.Layout.Seats, seatIDs)
	if err != nil {
		return err
	}

	ref, err := booking.GenerateBookingRef()
	if err != nil {
		return err
	}

	bookingRepo := booking.NewMongoRepository(database)
	demo := &booking.Booking{
		UserID:      seedUser.ID,
		ShowtimeID:  showtime.ID,
		Seats:       seatIDs,
		Total:       total,
		BookingRef:  ref,
		TicketToken: "seed-demo-token",
		Status:      booking.StatusConfirmed,
		ConfirmedAt: now,
	}
	if err := bookingRepo.Insert(ctx, demo); err != nil {
		return err
	}
	log.Printf("seed: created demo booking %s", ref)
	return nil
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

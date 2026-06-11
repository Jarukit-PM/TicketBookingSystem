package main

import (
	"context"
	"flag"
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
	resetCatalog := flag.Bool("reset-catalog", false, "clear cinemas/movies/screens/showtimes and re-seed Thai catalog")
	flag.Parse()

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

	switch {
	case *resetCatalog:
		if err := clearCatalog(ctx, database); err != nil {
			log.Fatalf("clear catalog: %v", err)
		}
		if err := seedThaiCatalog(ctx, repos); err != nil {
			log.Fatalf("seed thai catalog: %v", err)
		}
	case len(existing) > 0:
		log.Println("seed: catalog data already present, skipping catalog (use -reset-catalog to replace with Thai data)")
	default:
		if err := seedThaiCatalog(ctx, repos); err != nil {
			log.Fatalf("seed thai catalog: %v", err)
		}
	}

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

	defaultTiers := thaiPriceTiers
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
	userRepo := user.NewMongoRepository(database)
	return auth.BootstrapConfiguredAdmin(ctx, cfg, userRepo)
}

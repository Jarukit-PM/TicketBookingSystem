package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "time/tzdata"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/audit"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/booking"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/catalog"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/config"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/db"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/email"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/tasks"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/user"
)

func main() {
	cfg := config.MustLoad()
	ctx := context.Background()
	mongoClient := db.MustConnectMongo(ctx, cfg.MongoURI)
	redisClient := db.MustConnectRedis(cfg.RedisURL)
	defer func() {
		_ = db.DisconnectMongo(context.Background(), mongoClient)
		_ = redisClient.Close()
	}()

	database := db.Database(mongoClient, cfg.MongoURI)
	catalogRepos := catalog.NewMongoRepositories(database)
	auditRepos := audit.NewMongoRepositories(database)
	bookingRepo := booking.NewMongoRepository(database)
	userRepo := user.NewMongoRepository(database)

	emailSvc := email.NewService(
		bookingRepo,
		userRepo,
		email.CatalogReader{
			Showtimes: catalogRepos.Showtimes,
			Screens:   catalogRepos.Screens,
			Movies:    catalogRepos.Movies,
			Cinemas:   catalogRepos.Cinemas,
		},
		auditRepos.EmailLogs,
		email.NewBrevoClient(cfg.BrevoAPIKey, cfg.EmailFrom),
		cfg.AppURL,
	)

	asynqRedisOpt, err := asynqRedisOptFromURL(cfg.RedisURL)
	if err != nil {
		log.Fatalf("asynq redis: %v", err)
	}

	srv := asynq.NewServer(asynqRedisOpt, asynq.Config{})
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeEmailSend, emailSvc.HandleEmailSend)

	go func() {
		log.Println("worker ready")
		if err := srv.Run(mux); err != nil {
			log.Fatalf("worker run: %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("worker shutting down")
	srv.Shutdown()
}

func asynqRedisOptFromURL(redisURL string) (asynq.RedisClientOpt, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return asynq.RedisClientOpt{}, err
	}
	return asynq.RedisClientOpt{Addr: opts.Addr, Password: opts.Password, DB: opts.DB}, nil
}

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/config"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/db"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/tasks"
)

func main() {
	cfg := config.MustLoad()

	ctx := context.Background()
	mongoClient := db.MustConnectMongo(ctx, cfg.MongoURI)
	redisClient := db.MustConnectRedis(cfg.RedisURL)

	defer func() {
		if err := db.DisconnectMongo(context.Background(), mongoClient); err != nil {
			log.Printf("mongo disconnect: %v", err)
		}
		if err := redisClient.Close(); err != nil {
			log.Printf("redis close: %v", err)
		}
	}()

	asynqRedisOpt, err := asynqRedisOptFromURL(cfg.RedisURL)
	if err != nil {
		log.Fatalf("asynq redis: %v", err)
	}

	srv := asynq.NewServer(asynqRedisOpt, asynq.Config{})
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeEmailSend, handleEmailSend)

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

func handleEmailSend(ctx context.Context, t *asynq.Task) error {
	log.Printf("email:send stub: %s", string(t.Payload()))
	return nil
}

func asynqRedisOptFromURL(redisURL string) (asynq.RedisClientOpt, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return asynq.RedisClientOpt{}, err
	}

	return asynq.RedisClientOpt{
		Addr:     opts.Addr,
		Password: opts.Password,
		DB:       opts.DB,
	}, nil
}

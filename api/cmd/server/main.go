package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/config"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/db"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/server"
)

func main() {
	cfg := config.MustLoad()

	ctx := context.Background()
	mongoClient := db.MustConnectMongo(ctx, cfg.MongoURI)
	if err := db.EnsureIndexes(ctx, db.Database(mongoClient, cfg.MongoURI)); err != nil {
		log.Fatalf("ensure indexes: %v", err)
	}
	redisClient := db.MustConnectRedis(cfg.RedisURL)

	defer func() {
		if err := db.DisconnectMongo(context.Background(), mongoClient); err != nil {
			log.Printf("mongo disconnect: %v", err)
		}
		if err := redisClient.Close(); err != nil {
			log.Printf("redis close: %v", err)
		}
	}()

	router := server.NewRouter(server.Deps{
		Mongo: mongoClient,
		Redis: redisClient,
	})

	addr := "0.0.0.0:" + cfg.Port
	srv := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("api server listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("api server shutting down")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}

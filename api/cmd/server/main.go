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

	_ "time/tzdata" // embed IANA zones for Alpine/scratch images (Asia/Bangkok, etc.)

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/auth"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/config"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/db"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/server"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/user"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/tasks"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/ws"
)

func main() {
	cfg := config.MustLoad()

	ctx := context.Background()
	mongoClient := db.MustConnectMongo(ctx, cfg.MongoURI)
	database := db.Database(mongoClient, cfg.MongoURI)
	if err := db.EnsureIndexes(ctx, database); err != nil {
		log.Fatalf("ensure indexes: %v", err)
	}
	userRepo := user.NewMongoRepository(database)
	if err := auth.BootstrapConfiguredAdmin(ctx, cfg, userRepo); err != nil {
		log.Fatalf("bootstrap admin: %v", err)
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

	hub := ws.NewHub(redisClient)
	hubCtx, hubCancel := context.WithCancel(ctx)
	defer hubCancel()
	go hub.Run(hubCtx)

	taskClient, err := tasks.NewClient(cfg.RedisURL)
	if err != nil {
		log.Fatalf("task client: %v", err)
	}
	defer func() {
		if err := taskClient.Close(); err != nil {
			log.Printf("task client close: %v", err)
		}
	}()

	router := server.NewRouter(server.Deps{
		Config:     cfg,
		Mongo:      mongoClient,
		Redis:      redisClient,
		Hub:        hub,
		TaskClient: taskClient,
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

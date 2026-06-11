package server

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/handler"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/middleware"
)

// Deps holds shared dependencies for HTTP routes.
type Deps struct {
	Mongo *mongo.Client
	Redis *redis.Client
}

// NewRouter builds the Gin engine with middleware and routes.
func NewRouter(deps Deps) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery(), middleware.RequestID())

	healthDeps := handler.HealthDeps{
		Mongo: deps.Mongo,
		Redis: deps.Redis,
	}
	r.GET("/api/health", handler.Health(healthDeps))

	return r
}

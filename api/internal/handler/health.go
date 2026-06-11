package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/db"
	"github.com/Jarukit-PM/TicketBookingSystem/api/pkg/httputil"
)

const deepCheckTimeout = 3 * time.Second

// HealthDeps holds dependencies for the health endpoint.
type HealthDeps struct {
	Mongo *mongo.Client
	Redis *redis.Client
}

// Health returns a Gin handler for GET /api/health.
func Health(deps HealthDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Query("deep") != "1" {
			httputil.OK(c, gin.H{"status": "ok"})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), deepCheckTimeout)
		defer cancel()

		if err := db.PingMongo(ctx, deps.Mongo); err != nil {
			httputil.Error(c, http.StatusServiceUnavailable, "MONGO_UNAVAILABLE", err.Error())
			return
		}

		if err := db.PingRedis(ctx, deps.Redis); err != nil {
			httputil.Error(c, http.StatusServiceUnavailable, "REDIS_UNAVAILABLE", err.Error())
			return
		}

		httputil.OK(c, gin.H{"status": "ok"})
	}
}

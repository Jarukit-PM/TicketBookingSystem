package server

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/auth"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/config"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/db"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/handler"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/middleware"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/user"
)

// Deps holds shared dependencies for HTTP routes.
type Deps struct {
	Config config.Config
	Mongo  *mongo.Client
	Redis  *redis.Client
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

	database := db.Database(deps.Mongo, deps.Config.MongoURI)
	userRepo := user.NewMongoRepository(database)
	tokenSvc := auth.NewTokenService(deps.Config.JWTSecret, deps.Config.JWTExpiryDuration())
	rateLimiter := auth.NewLoginRateLimiter(deps.Redis)
	authSvc := auth.NewService(userRepo, tokenSvc, rateLimiter, deps.Config.AdminEmail)

	googleOAuth := auth.NewGoogleOAuth(
		deps.Config.GoogleClientID,
		deps.Config.GoogleClientSecret,
		deps.Config.GoogleRedirectURL(),
	)

	cookieOpts := auth.CookieOptions{Secure: deps.Config.CookieSecure()}
	authDeps := handler.AuthDeps{
		Service:       authSvc,
		Google:        googleOAuth,
		AppURL:        deps.Config.AppURL,
		CookieOptions: cookieOpts,
	}
	authMw := auth.MiddlewareDeps{Service: authSvc}

	api := r.Group("/api")
	authRoutes := api.Group("/auth")
	authRoutes.POST("/register", handler.Register(authDeps))
	authRoutes.POST("/login", handler.Login(authDeps))
	authRoutes.POST("/logout", handler.Logout(authDeps))
	authRoutes.GET("/me", auth.RequireAuth(authMw), handler.Me(authDeps))
	authRoutes.GET("/google", handler.GoogleStart(authDeps))
	authRoutes.GET("/google/callback", handler.GoogleCallback(authDeps))

	return r
}

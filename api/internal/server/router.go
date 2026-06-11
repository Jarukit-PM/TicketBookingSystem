package server

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/auth"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/booking"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/catalog"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/config"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/db"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/handler"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/hold"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/inventory"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/middleware"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/user"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/tasks"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/ws"
)

// Deps holds shared dependencies for HTTP routes.
type Deps struct {
	Config      config.Config
	Mongo       *mongo.Client
	Redis       *redis.Client
	Hub         *ws.Hub
	TaskClient  *tasks.Client
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

	cookieOpts := auth.CookieOptions{Secure: deps.Config.CookieSecure()}
	authDeps := handler.AuthDeps{Service: authSvc, CookieOptions: cookieOpts}
	authMw := auth.MiddlewareDeps{Service: authSvc}

	api := r.Group("/api")
	authRoutes := api.Group("/auth")
	authRoutes.POST("/register", handler.Register(authDeps))
	authRoutes.POST("/login", handler.Login(authDeps))
	authRoutes.POST("/logout", handler.Logout(authDeps))
	authRoutes.GET("/me", auth.RequireAuth(authMw), handler.Me(authDeps))


	catalogRepos := catalog.NewMongoRepositories(database)
	catalogSvc := &catalog.Service{
		Cinemas:   catalogRepos.Cinemas,
		Screens:   catalogRepos.Screens,
		Movies:    catalogRepos.Movies,
		Showtimes: catalogRepos.Showtimes,
	}
	catalogDeps := handler.CatalogDeps{Service: catalogSvc}
	api.GET("/cinemas", handler.ListCinemas(catalogDeps))
	api.GET("/movies", handler.ListMovies(catalogDeps))
	api.GET("/movies/:id", handler.GetMovie(catalogDeps))
	api.GET("/showtimes", handler.ListShowtimes(catalogDeps))

	bookingRepo := booking.NewMongoRepository(database)
	inventorySvc := inventory.NewService(
		catalogRepos.Showtimes,
		catalogRepos.Screens,
		bookingRepo,
		deps.Redis,
	)
	seatsDeps := handler.SeatsDeps{Inventory: inventorySvc}
	api.GET("/showtimes/:id/seats", handler.GetShowtimeSeats(seatsDeps))

	holdSvc := hold.NewService(
		catalogRepos.Showtimes,
		catalogRepos.Screens,
		catalogRepos.Cinemas,
		bookingRepo,
		deps.Redis,
	)
	holdsDeps := handler.HoldsDeps{Holds: holdSvc, Publisher: deps.Hub}
	showtimeHolds := api.Group("/showtimes/:id/holds")
	showtimeHolds.POST("", auth.RequireAuth(authMw), handler.AddShowtimeHolds(holdsDeps))
	showtimeHolds.DELETE("", auth.RequireAuth(authMw), handler.RemoveShowtimeHolds(holdsDeps))

	idempotency := booking.NewIdempotencyStore(deps.Redis, 0)
	bookingSvc := booking.NewService(
		catalogRepos.Showtimes,
		catalogRepos.Screens,
		catalogRepos.Cinemas,
		catalogRepos.Movies,
		bookingRepo,
		holdSvc,
		deps.Redis,
		idempotency,
		booking.WithTicketConfig(deps.Config.TicketHMACSecret(), deps.Config.AppURL),
	)
	bookingsDeps := handler.BookingsDeps{
		Bookings:  bookingSvc,
		Tasks:     deps.TaskClient,
		Publisher: deps.Hub,
	}
	api.POST("/bookings/confirm", auth.RequireAuth(authMw), handler.ConfirmBooking(bookingsDeps))
	api.GET("/bookings/:id/ticket", auth.RequireAuth(authMw), handler.GetBookingTicket(bookingsDeps))

	wsDeps := ws.HandlerDeps{Hub: deps.Hub, Inventory: inventorySvc}
	r.GET("/ws/showtimes/:id", ws.Showtime(wsDeps))

	return r
}

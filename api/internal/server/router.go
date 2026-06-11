package server

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"

	adminpkg "github.com/Jarukit-PM/TicketBookingSystem/api/internal/admin"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/audit"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/auth"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/booking"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/catalog"
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

	auditRepos := audit.NewMongoRepositories(database)
	adminCatalogSvc := catalog.NewAdminService(catalogRepos, auditRepos.AuditLogs)
	adminCatalogDeps := handler.AdminCatalogDeps{Service: adminCatalogSvc}

	admin := api.Group("/admin")
	admin.Use(auth.RequireAuth(authMw), auth.RequireAdmin(authMw))

	bookingRepo := booking.NewMongoRepository(database)
	dashboardSvc := &adminpkg.DashboardService{
		Showtimes: catalogRepos.Showtimes,
		Screens:   catalogRepos.Screens,
		Movies:    catalogRepos.Movies,
		Bookings:  bookingRepo,
	}
	dashboardDeps := handler.AdminDashboardDeps{Service: dashboardSvc}
	admin.GET("/dashboard", handler.GetAdminDashboard(dashboardDeps))

	movies := admin.Group("/movies")
	movies.GET("", handler.ListAdminMovies(adminCatalogDeps))
	movies.POST("", handler.CreateAdminMovie(adminCatalogDeps))
	movies.GET("/:id", handler.GetAdminMovie(adminCatalogDeps))
	movies.PUT("/:id", handler.UpdateAdminMovie(adminCatalogDeps))
	movies.DELETE("/:id", handler.DeleteAdminMovie(adminCatalogDeps))

	cinemas := admin.Group("/cinemas")
	cinemas.GET("", handler.ListAdminCinemas(adminCatalogDeps))
	cinemas.POST("", handler.CreateAdminCinema(adminCatalogDeps))
	cinemas.GET("/:id", handler.GetAdminCinema(adminCatalogDeps))
	cinemas.PUT("/:id", handler.UpdateAdminCinema(adminCatalogDeps))
	cinemas.DELETE("/:id", handler.DeleteAdminCinema(adminCatalogDeps))

	screens := admin.Group("/screens")
	screens.GET("", handler.ListAdminScreens(adminCatalogDeps))
	screens.POST("", handler.CreateAdminScreen(adminCatalogDeps))
	screens.GET("/:id", handler.GetAdminScreen(adminCatalogDeps))
	screens.PUT("/:id", handler.UpdateAdminScreen(adminCatalogDeps))
	screens.DELETE("/:id", handler.DeleteAdminScreen(adminCatalogDeps))

	showtimes := admin.Group("/showtimes")
	showtimes.GET("", handler.ListAdminShowtimes(adminCatalogDeps))
	showtimes.POST("", handler.CreateAdminShowtime(adminCatalogDeps))
	showtimes.GET("/:id", handler.GetAdminShowtime(adminCatalogDeps))
	showtimes.PUT("/:id", handler.UpdateAdminShowtime(adminCatalogDeps))
	showtimes.DELETE("/:id", handler.DeleteAdminShowtime(adminCatalogDeps))

	return r
}

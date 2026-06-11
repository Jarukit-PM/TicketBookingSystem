package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/auth"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/catalog"
	"github.com/Jarukit-PM/TicketBookingSystem/api/pkg/httputil"
)

// AdminCatalogDeps holds admin catalog handler dependencies.
type AdminCatalogDeps struct {
	Service *catalog.AdminService
}

// ListAdminMovies handles GET /api/admin/movies.
func ListAdminMovies(deps AdminCatalogDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		movies, err := deps.Service.ListMovies(c.Request.Context())
		if err != nil {
			httputil.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
			return
		}
		httputil.OK(c, gin.H{"movies": movies})
	}
}

// GetAdminMovie handles GET /api/admin/movies/:id.
func GetAdminMovie(deps AdminCatalogDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := parseObjectID(c, "id")
		if !ok {
			return
		}
		movie, err := deps.Service.GetMovie(c.Request.Context(), id)
		if err != nil {
			writeCatalogError(c, err)
			return
		}
		httputil.OK(c, gin.H{"movie": movie})
	}
}

// CreateAdminMovie handles POST /api/admin/movies.
func CreateAdminMovie(deps AdminCatalogDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		actorID, ok := actorIDFromContext(c)
		if !ok {
			return
		}
		var movie catalog.Movie
		if err := c.ShouldBindJSON(&movie); err != nil {
			httputil.Error(c, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
			return
		}
		created, err := deps.Service.CreateMovie(c.Request.Context(), actorID, &movie)
		if err != nil {
			writeCatalogError(c, err)
			return
		}
		c.JSON(http.StatusCreated, gin.H{"movie": created})
	}
}

// UpdateAdminMovie handles PUT /api/admin/movies/:id.
func UpdateAdminMovie(deps AdminCatalogDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		actorID, ok := actorIDFromContext(c)
		if !ok {
			return
		}
		id, ok := parseObjectID(c, "id")
		if !ok {
			return
		}
		var movie catalog.Movie
		if err := c.ShouldBindJSON(&movie); err != nil {
			httputil.Error(c, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
			return
		}
		updated, err := deps.Service.UpdateMovie(c.Request.Context(), actorID, id, &movie)
		if err != nil {
			writeCatalogError(c, err)
			return
		}
		httputil.OK(c, gin.H{"movie": updated})
	}
}

// DeleteAdminMovie handles DELETE /api/admin/movies/:id.
func DeleteAdminMovie(deps AdminCatalogDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		actorID, ok := actorIDFromContext(c)
		if !ok {
			return
		}
		id, ok := parseObjectID(c, "id")
		if !ok {
			return
		}
		if err := deps.Service.DeleteMovie(c.Request.Context(), actorID, id); err != nil {
			writeCatalogError(c, err)
			return
		}
		c.Status(http.StatusNoContent)
	}
}

// ListAdminCinemas handles GET /api/admin/cinemas.
func ListAdminCinemas(deps AdminCatalogDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		cinemas, err := deps.Service.ListCinemas(c.Request.Context())
		if err != nil {
			httputil.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
			return
		}
		httputil.OK(c, gin.H{"cinemas": cinemas})
	}
}

// GetAdminCinema handles GET /api/admin/cinemas/:id.
func GetAdminCinema(deps AdminCatalogDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := parseObjectID(c, "id")
		if !ok {
			return
		}
		cinema, err := deps.Service.GetCinema(c.Request.Context(), id)
		if err != nil {
			writeCatalogError(c, err)
			return
		}
		httputil.OK(c, gin.H{"cinema": cinema})
	}
}

// CreateAdminCinema handles POST /api/admin/cinemas.
func CreateAdminCinema(deps AdminCatalogDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		actorID, ok := actorIDFromContext(c)
		if !ok {
			return
		}
		var cinema catalog.Cinema
		if err := c.ShouldBindJSON(&cinema); err != nil {
			httputil.Error(c, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
			return
		}
		created, err := deps.Service.CreateCinema(c.Request.Context(), actorID, &cinema)
		if err != nil {
			writeCatalogError(c, err)
			return
		}
		c.JSON(http.StatusCreated, gin.H{"cinema": created})
	}
}

// UpdateAdminCinema handles PUT /api/admin/cinemas/:id.
func UpdateAdminCinema(deps AdminCatalogDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		actorID, ok := actorIDFromContext(c)
		if !ok {
			return
		}
		id, ok := parseObjectID(c, "id")
		if !ok {
			return
		}
		var cinema catalog.Cinema
		if err := c.ShouldBindJSON(&cinema); err != nil {
			httputil.Error(c, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
			return
		}
		updated, err := deps.Service.UpdateCinema(c.Request.Context(), actorID, id, &cinema)
		if err != nil {
			writeCatalogError(c, err)
			return
		}
		httputil.OK(c, gin.H{"cinema": updated})
	}
}

// DeleteAdminCinema handles DELETE /api/admin/cinemas/:id.
func DeleteAdminCinema(deps AdminCatalogDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		actorID, ok := actorIDFromContext(c)
		if !ok {
			return
		}
		id, ok := parseObjectID(c, "id")
		if !ok {
			return
		}
		if err := deps.Service.DeleteCinema(c.Request.Context(), actorID, id); err != nil {
			writeCatalogError(c, err)
			return
		}
		c.Status(http.StatusNoContent)
	}
}

// ListAdminScreens handles GET /api/admin/screens.
func ListAdminScreens(deps AdminCatalogDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		var cinemaID *primitive.ObjectID
		if raw := c.Query("cinemaId"); raw != "" {
			id, err := primitive.ObjectIDFromHex(raw)
			if err != nil {
				httputil.Error(c, http.StatusBadRequest, "INVALID_QUERY", "invalid cinemaId")
				return
			}
			cinemaID = &id
		}
		screens, err := deps.Service.ListScreens(c.Request.Context(), cinemaID)
		if err != nil {
			httputil.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
			return
		}
		httputil.OK(c, gin.H{"screens": screens})
	}
}

// GetAdminScreen handles GET /api/admin/screens/:id.
func GetAdminScreen(deps AdminCatalogDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := parseObjectID(c, "id")
		if !ok {
			return
		}
		screen, err := deps.Service.GetScreen(c.Request.Context(), id)
		if err != nil {
			writeCatalogError(c, err)
			return
		}
		httputil.OK(c, gin.H{"screen": screen})
	}
}

// CreateAdminScreen handles POST /api/admin/screens.
func CreateAdminScreen(deps AdminCatalogDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		actorID, ok := actorIDFromContext(c)
		if !ok {
			return
		}
		var screen catalog.Screen
		if err := c.ShouldBindJSON(&screen); err != nil {
			httputil.Error(c, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
			return
		}
		created, err := deps.Service.CreateScreen(c.Request.Context(), actorID, &screen)
		if err != nil {
			writeCatalogError(c, err)
			return
		}
		c.JSON(http.StatusCreated, gin.H{"screen": created})
	}
}

// UpdateAdminScreen handles PUT /api/admin/screens/:id.
func UpdateAdminScreen(deps AdminCatalogDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		actorID, ok := actorIDFromContext(c)
		if !ok {
			return
		}
		id, ok := parseObjectID(c, "id")
		if !ok {
			return
		}
		var screen catalog.Screen
		if err := c.ShouldBindJSON(&screen); err != nil {
			httputil.Error(c, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
			return
		}
		updated, err := deps.Service.UpdateScreen(c.Request.Context(), actorID, id, &screen)
		if err != nil {
			writeCatalogError(c, err)
			return
		}
		httputil.OK(c, gin.H{"screen": updated})
	}
}

// DeleteAdminScreen handles DELETE /api/admin/screens/:id.
func DeleteAdminScreen(deps AdminCatalogDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		actorID, ok := actorIDFromContext(c)
		if !ok {
			return
		}
		id, ok := parseObjectID(c, "id")
		if !ok {
			return
		}
		if err := deps.Service.DeleteScreen(c.Request.Context(), actorID, id); err != nil {
			writeCatalogError(c, err)
			return
		}
		c.Status(http.StatusNoContent)
	}
}

// ListAdminShowtimes handles GET /api/admin/showtimes.
func ListAdminShowtimes(deps AdminCatalogDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		filter, ok := parseShowtimeFilter(c)
		if !ok {
			return
		}
		showtimes, err := deps.Service.ListShowtimes(c.Request.Context(), filter)
		if err != nil {
			httputil.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
			return
		}
		httputil.OK(c, gin.H{"showtimes": showtimes})
	}
}

// GetAdminShowtime handles GET /api/admin/showtimes/:id.
func GetAdminShowtime(deps AdminCatalogDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := parseObjectID(c, "id")
		if !ok {
			return
		}
		showtime, err := deps.Service.GetShowtime(c.Request.Context(), id)
		if err != nil {
			writeCatalogError(c, err)
			return
		}
		httputil.OK(c, gin.H{"showtime": showtime})
	}
}

// CreateAdminShowtime handles POST /api/admin/showtimes.
func CreateAdminShowtime(deps AdminCatalogDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		actorID, ok := actorIDFromContext(c)
		if !ok {
			return
		}
		var showtime catalog.Showtime
		if err := c.ShouldBindJSON(&showtime); err != nil {
			httputil.Error(c, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
			return
		}
		created, err := deps.Service.CreateShowtime(c.Request.Context(), actorID, &showtime)
		if err != nil {
			writeCatalogError(c, err)
			return
		}
		c.JSON(http.StatusCreated, gin.H{"showtime": created})
	}
}

// UpdateAdminShowtime handles PUT /api/admin/showtimes/:id.
func UpdateAdminShowtime(deps AdminCatalogDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		actorID, ok := actorIDFromContext(c)
		if !ok {
			return
		}
		id, ok := parseObjectID(c, "id")
		if !ok {
			return
		}
		var showtime catalog.Showtime
		if err := c.ShouldBindJSON(&showtime); err != nil {
			httputil.Error(c, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
			return
		}
		updated, err := deps.Service.UpdateShowtime(c.Request.Context(), actorID, id, &showtime)
		if err != nil {
			writeCatalogError(c, err)
			return
		}
		httputil.OK(c, gin.H{"showtime": updated})
	}
}

// DeleteAdminShowtime handles DELETE /api/admin/showtimes/:id.
func DeleteAdminShowtime(deps AdminCatalogDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		actorID, ok := actorIDFromContext(c)
		if !ok {
			return
		}
		id, ok := parseObjectID(c, "id")
		if !ok {
			return
		}
		if err := deps.Service.DeleteShowtime(c.Request.Context(), actorID, id); err != nil {
			writeCatalogError(c, err)
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func actorIDFromContext(c *gin.Context) (primitive.ObjectID, bool) {
	sessionUser, ok := auth.UserFromContext(c)
	if !ok {
		httputil.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
		return primitive.NilObjectID, false
	}
	return sessionUser.ID, true
}

func parseObjectID(c *gin.Context, param string) (primitive.ObjectID, bool) {
	id, err := primitive.ObjectIDFromHex(c.Param(param))
	if err != nil {
		httputil.Error(c, http.StatusBadRequest, "INVALID_ID", "invalid id")
		return primitive.NilObjectID, false
	}
	return id, true
}

func parseShowtimeFilter(c *gin.Context) (catalog.AdminShowtimeFilter, bool) {
	filter := catalog.AdminShowtimeFilter{}
	if raw := c.Query("cinemaId"); raw != "" {
		id, err := primitive.ObjectIDFromHex(raw)
		if err != nil {
			httputil.Error(c, http.StatusBadRequest, "INVALID_QUERY", "invalid cinemaId")
			return filter, false
		}
		filter.CinemaID = &id
	}
	if raw := c.Query("movieId"); raw != "" {
		id, err := primitive.ObjectIDFromHex(raw)
		if err != nil {
			httputil.Error(c, http.StatusBadRequest, "INVALID_QUERY", "invalid movieId")
			return filter, false
		}
		filter.MovieID = &id
	}
	if raw := c.Query("screenId"); raw != "" {
		id, err := primitive.ObjectIDFromHex(raw)
		if err != nil {
			httputil.Error(c, http.StatusBadRequest, "INVALID_QUERY", "invalid screenId")
			return filter, false
		}
		filter.ScreenID = &id
	}
	if raw := c.Query("from"); raw != "" {
		t, err := time.Parse(time.RFC3339, raw)
		if err != nil {
			httputil.Error(c, http.StatusBadRequest, "INVALID_QUERY", "invalid from")
			return filter, false
		}
		filter.From = &t
	}
	if raw := c.Query("to"); raw != "" {
		t, err := time.Parse(time.RFC3339, raw)
		if err != nil {
			httputil.Error(c, http.StatusBadRequest, "INVALID_QUERY", "invalid to")
			return filter, false
		}
		filter.To = &t
	}
	return filter, true
}


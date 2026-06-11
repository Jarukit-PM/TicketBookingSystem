package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/catalog"
	"github.com/Jarukit-PM/TicketBookingSystem/api/pkg/httputil"
)

// CatalogDeps holds dependencies for public catalog routes.
type CatalogDeps struct {
	Service *catalog.Service
}

// ListCinemas handles GET /api/cinemas.
func ListCinemas(deps CatalogDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		cinemas, err := deps.Service.ListCinemas(c.Request.Context())
		if err != nil {
			httputil.Error(c, http.StatusInternalServerError, "CATALOG_ERROR", err.Error())
			return
		}
		if cinemas == nil {
			cinemas = []catalog.Cinema{}
		}
		httputil.OK(c, cinemas)
	}
}

// ListMovies handles GET /api/movies.
func ListMovies(deps CatalogDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		cinemaID, err := parseObjectIDQuery(c, "cinemaId")
		if err != nil {
			httputil.Error(c, http.StatusBadRequest, "INVALID_CINEMA_ID", err.Error())
			return
		}

		tab, ok := catalog.ParseBrowseTab(c.Query("tab"))
		if !ok {
			httputil.Error(c, http.StatusBadRequest, "INVALID_TAB", "tab must be now_showing or coming_soon")
			return
		}

		movies, err := deps.Service.BrowseMovies(c.Request.Context(), cinemaID, tab)
		if err != nil {
			writeCatalogError(c, err)
			return
		}
		if movies == nil {
			movies = []catalog.Movie{}
		}
		httputil.OK(c, movies)
	}
}

// GetMovie handles GET /api/movies/:id.
func GetMovie(deps CatalogDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		movieID, err := parseObjectIDParam(c, "id")
		if err != nil {
			httputil.Error(c, http.StatusBadRequest, "INVALID_MOVIE_ID", err.Error())
			return
		}

		cinemaID, err := parseObjectIDQuery(c, "cinemaId")
		if err != nil {
			httputil.Error(c, http.StatusBadRequest, "INVALID_CINEMA_ID", err.Error())
			return
		}

		detail, err := deps.Service.GetMovieDetail(c.Request.Context(), movieID, cinemaID)
		if err != nil {
			writeCatalogError(c, err)
			return
		}
		httputil.OK(c, detail)
	}
}

// ListShowtimes handles GET /api/showtimes.
func ListShowtimes(deps CatalogDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		cinemaID, err := parseObjectIDQuery(c, "cinemaId")
		if err != nil {
			httputil.Error(c, http.StatusBadRequest, "INVALID_CINEMA_ID", err.Error())
			return
		}

		movieID, err := parseObjectIDQuery(c, "movieId")
		if err != nil {
			httputil.Error(c, http.StatusBadRequest, "INVALID_MOVIE_ID", err.Error())
			return
		}

		var date *time.Time
		if raw := c.Query("date"); raw != "" {
			parsed, err := time.Parse("2006-01-02", raw)
			if err != nil {
				httputil.Error(c, http.StatusBadRequest, "INVALID_DATE", "date must be YYYY-MM-DD")
				return
			}
			date = &parsed
		}

		showtimes, err := deps.Service.ListShowtimes(c.Request.Context(), cinemaID, movieID, date)
		if err != nil {
			writeCatalogError(c, err)
			return
		}
		if showtimes == nil {
			showtimes = []catalog.ShowtimeView{}
		}
		httputil.OK(c, showtimes)
	}
}

func parseObjectIDQuery(c *gin.Context, key string) (primitive.ObjectID, error) {
	raw := c.Query(key)
	if raw == "" {
		return primitive.NilObjectID, errors.New(key + " is required")
	}
	return primitive.ObjectIDFromHex(raw)
}

func parseObjectIDParam(c *gin.Context, key string) (primitive.ObjectID, error) {
	raw := c.Param(key)
	if raw == "" {
		return primitive.NilObjectID, errors.New(key + " is required")
	}
	return primitive.ObjectIDFromHex(raw)
}


func writeCatalogError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, catalog.ErrCinemaNotFound):
		httputil.Error(c, http.StatusNotFound, "CINEMA_NOT_FOUND", "cinema not found")
	case errors.Is(err, catalog.ErrMovieNotFound):
		httputil.Error(c, http.StatusNotFound, "MOVIE_NOT_FOUND", "movie not found")
	case errors.Is(err, catalog.ErrNotFound):
		httputil.Error(c, http.StatusNotFound, "NOT_FOUND", err.Error())
	case errors.Is(err, catalog.ErrInvalidInput),
		errors.Is(err, catalog.ErrInvalidStatus),
		errors.Is(err, catalog.ErrInvalidSeat),
		errors.Is(err, catalog.ErrInvalidShowtime):
		httputil.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	default:
		httputil.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}
}

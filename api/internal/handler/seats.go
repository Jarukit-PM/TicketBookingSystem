package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/inventory"
	"github.com/Jarukit-PM/TicketBookingSystem/api/pkg/httputil"
)

// SeatsDeps holds dependencies for seat map handlers.
type SeatsDeps struct {
	Inventory *inventory.Service
}

// GetShowtimeSeats handles GET /api/showtimes/:id/seats.
func GetShowtimeSeats(deps SeatsDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		showtimeID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			httputil.Error(c, http.StatusBadRequest, "INVALID_SHOWTIME_ID", "invalid showtime id")
			return
		}

		snapshot, err := deps.Inventory.Snapshot(c.Request.Context(), showtimeID)
		if err != nil {
			if errors.Is(err, inventory.ErrShowtimeNotFound) {
				httputil.Error(c, http.StatusNotFound, "SHOWTIME_NOT_FOUND", "showtime not found")
				return
			}
			if errors.Is(err, inventory.ErrScreenNotFound) {
				httputil.Error(c, http.StatusNotFound, "SCREEN_NOT_FOUND", "screen not found")
				return
			}
			httputil.Error(c, http.StatusInternalServerError, "INVENTORY_ERROR", "failed to load seat map")
			return
		}

		httputil.OK(c, snapshot)
	}
}

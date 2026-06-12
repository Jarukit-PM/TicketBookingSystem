package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/booking"
	"github.com/Jarukit-PM/TicketBookingSystem/api/pkg/httputil"
)

// GetPublicTicket handles GET /api/tickets/:ref?t=.
func GetPublicTicket(deps BookingsDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ticket, err := deps.Bookings.GetTicketByRef(
			c.Request.Context(),
			c.Param("ref"),
			c.Query("t"),
		)
		if err != nil {
			writePublicTicketError(c, err)
			return
		}
		httputil.OK(c, ticket)
	}
}

func writePublicTicketError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, booking.ErrInvalidTicket):
		httputil.Error(c, http.StatusNotFound, "INVALID_TICKET", "invalid or unknown ticket")
	default:
		httputil.Error(c, http.StatusInternalServerError, "TICKET_ERROR", "failed to load ticket")
	}
}

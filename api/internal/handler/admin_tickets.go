package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/admin"
	"github.com/Jarukit-PM/TicketBookingSystem/api/pkg/httputil"
)

// AdminTicketsDeps holds admin ticket handler dependencies.
type AdminTicketsDeps struct {
	Service *admin.TicketsService
}

// ResolveAdminTicket handles GET /api/admin/tickets/resolve.
func ResolveAdminTicket(deps AdminTicketsDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := deps.Service.Resolve(c.Request.Context(), c.Query("ref"), c.Query("t"))
		if err != nil {
			if errors.Is(err, admin.ErrInvalidTicket) {
				httputil.Error(c, http.StatusNotFound, "INVALID_TICKET", "invalid or unknown ticket")
				return
			}
			httputil.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
			return
		}
		httputil.OK(c, result)
	}
}

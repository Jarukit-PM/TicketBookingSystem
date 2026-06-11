package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/admin"
	"github.com/Jarukit-PM/TicketBookingSystem/api/pkg/httputil"
)

// AdminDashboardDeps holds admin dashboard handler dependencies.
type AdminDashboardDeps struct {
	Service *admin.DashboardService
}

// GetAdminDashboard handles GET /api/admin/dashboard.
func GetAdminDashboard(deps AdminDashboardDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		dashboard, err := deps.Service.GetDashboard(c.Request.Context())
		if err != nil {
			httputil.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
			return
		}
		httputil.OK(c, dashboard)
	}
}

package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/admin"
	"github.com/Jarukit-PM/TicketBookingSystem/api/pkg/httputil"
)

// AdminLogsDeps holds admin log handler dependencies.
type AdminLogsDeps struct {
	Service *admin.LogsService
}

// ListAdminAuditLogs handles GET /api/admin/audit-logs.
func ListAdminAuditLogs(deps AdminLogsDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, limit := parseLogPagination(c)
		logs, err := deps.Service.ListAuditLogs(c.Request.Context(), page, limit)
		if err != nil {
			httputil.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
			return
		}
		httputil.OK(c, gin.H{"logs": httputil.JSONSlice(logs), "page": page, "limit": limit})
	}
}

// ListAdminEmailLogs handles GET /api/admin/email-logs.
func ListAdminEmailLogs(deps AdminLogsDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, limit := parseLogPagination(c)
		var bookingID *primitive.ObjectID
		if raw := c.Query("bookingId"); raw != "" {
			id, err := primitive.ObjectIDFromHex(raw)
			if err != nil {
				httputil.Error(c, http.StatusBadRequest, "INVALID_QUERY", "invalid bookingId")
				return
			}
			bookingID = &id
		}
		logs, err := deps.Service.ListEmailLogs(c.Request.Context(), page, limit, bookingID)
		if err != nil {
			httputil.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
			return
		}
		httputil.OK(c, gin.H{"logs": httputil.JSONSlice(logs), "page": page, "limit": limit})
	}
}

func parseLogPagination(c *gin.Context) (page int, limit int) {
	page = 1
	limit = 50
	if raw := c.Query("page"); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil && v > 0 {
			page = v
		}
	}
	if raw := c.Query("limit"); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil && v > 0 {
			limit = v
		}
	}
	return page, limit
}

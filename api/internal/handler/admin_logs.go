package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/admin"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/audit"
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
		filter, err := parseAuditLogFilter(c)
		if err != nil {
			httputil.Error(c, http.StatusBadRequest, "INVALID_QUERY", err.Error())
			return
		}
		logs, err := deps.Service.ListAuditLogs(c.Request.Context(), page, limit, filter)
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
		filter, err := parseEmailLogFilter(c)
		if err != nil {
			httputil.Error(c, http.StatusBadRequest, "INVALID_QUERY", err.Error())
			return
		}
		logs, err := deps.Service.ListEmailLogs(c.Request.Context(), page, limit, filter)
		if err != nil {
			httputil.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
			return
		}
		httputil.OK(c, gin.H{"logs": httputil.JSONSlice(logs), "page": page, "limit": limit})
	}
}

func parseEmailLogFilter(c *gin.Context) (audit.EmailLogFilter, error) {
	filter := audit.EmailLogFilter{
		To:     strings.TrimSpace(c.Query("to")),
		Type:   c.Query("type"),
		Status: c.Query("status"),
	}
	if raw := c.Query("bookingId"); raw != "" {
		id, err := primitive.ObjectIDFromHex(raw)
		if err != nil {
			return filter, errors.New("invalid bookingId")
		}
		filter.BookingID = &id
	}
	if from, err := parseLogDateQuery(c, "sentFrom"); err != nil {
		return filter, err
	} else if from != nil {
		filter.SentFrom = from
	}
	if to, err := parseLogDateEndQuery(c, "sentTo"); err != nil {
		return filter, err
	} else if to != nil {
		filter.SentTo = to
	}
	return filter, nil
}

func parseLogDateQuery(c *gin.Context, key string) (*time.Time, error) {
	raw := strings.TrimSpace(c.Query(key))
	if raw == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return nil, errors.New("invalid " + key)
	}
	utc := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	return &utc, nil
}

func parseLogDateEndQuery(c *gin.Context, key string) (*time.Time, error) {
	raw := strings.TrimSpace(c.Query(key))
	if raw == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return nil, errors.New("invalid " + key)
	}
	utc := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC).Add(24 * time.Hour)
	return &utc, nil
}

func parseAuditLogFilter(c *gin.Context) (audit.AuditLogFilter, error) {
	filter := audit.AuditLogFilter{
		Action:     c.Query("action"),
		Entity:     c.Query("entity"),
		EntityID:   c.Query("entityId"),
		BookingRef: c.Query("bookingRef"),
	}
	if raw := c.Query("actorId"); raw != "" {
		id, err := primitive.ObjectIDFromHex(raw)
		if err != nil {
			return filter, errors.New("invalid actorId")
		}
		filter.ActorID = &id
	}
	return filter, nil
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

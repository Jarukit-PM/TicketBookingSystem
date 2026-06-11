package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/admin"
	"github.com/Jarukit-PM/TicketBookingSystem/api/pkg/httputil"
)

// AdminBookingsDeps holds admin booking handler dependencies.
type AdminBookingsDeps struct {
	Service *admin.BookingsService
}

// SearchAdminBookings handles GET /api/admin/bookings.
func SearchAdminBookings(deps AdminBookingsDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		q := admin.BookingSearchQuery{
			Email:      c.Query("email"),
			BookingRef: c.Query("bookingRef"),
			UserID:     c.Query("userId"),
			ShowtimeID: c.Query("showtimeId"),
		}
		page, limit := parseBookingPagination(c)
		result, err := deps.Service.Search(c.Request.Context(), q, page, limit)
		if err != nil {
			if errors.Is(err, admin.ErrInvalidQuery) {
				httputil.Error(c, http.StatusBadRequest, "INVALID_QUERY", "invalid search filter")
				return
			}
			httputil.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
			return
		}
		httputil.OK(c, gin.H{
			"bookings": httputil.JSONSlice(result.Bookings),
			"total":    result.Total,
			"page":     result.Page,
			"limit":    result.Limit,
		})
	}
}

func parseBookingPagination(c *gin.Context) (page int, limit int) {
	page = 1
	limit = 0
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

// ListAdminUserBookings handles GET /api/admin/users/:userId/bookings.
func ListAdminUserBookings(deps AdminBookingsDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := parseObjectID(c, "userId")
		if !ok {
			return
		}
		bookings, err := deps.Service.ListUserBookings(c.Request.Context(), userID)
		if err != nil {
			httputil.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
			return
		}
		httputil.OK(c, gin.H{"bookings": httputil.JSONSlice(bookings)})
	}
}

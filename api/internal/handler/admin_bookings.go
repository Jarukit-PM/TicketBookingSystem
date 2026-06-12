package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/admin"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/booking"
	"github.com/Jarukit-PM/TicketBookingSystem/api/pkg/httputil"
)

// AdminBookingsDeps holds admin booking handler dependencies.
type AdminBookingsDeps struct {
	Service *admin.BookingsService
}

// SearchAdminBookings handles GET /api/admin/bookings.
func SearchAdminBookings(deps AdminBookingsDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		q, ok := parseBookingSearchQuery(c)
		if !ok {
			return
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

func parseBookingSearchQuery(c *gin.Context) (admin.BookingSearchQuery, bool) {
	q := admin.BookingSearchQuery{
		Email:      c.Query("email"),
		BookingRef: c.Query("bookingRef"),
		UserID:     c.Query("userId"),
		ShowtimeID: c.Query("showtimeId"),
		MovieID:    c.Query("movieId"),
		Locale:     c.Query("locale"),
	}
	if from, ok := parseBookingDateQuery(c, "confirmedFrom"); !ok {
		return q, false
	} else if from != nil {
		q.ConfirmedFrom = from
	}
	if to, ok := parseBookingDateEndQuery(c, "confirmedTo"); !ok {
		return q, false
	} else if to != nil {
		q.ConfirmedTo = to
	}
	return q, true
}

func parseBookingDateQuery(c *gin.Context, key string) (*time.Time, bool) {
	raw := strings.TrimSpace(c.Query(key))
	if raw == "" {
		return nil, true
	}
	t, err := time.Parse("2006-01-02", raw)
	if err != nil {
		httputil.Error(c, http.StatusBadRequest, "INVALID_QUERY", "invalid "+key)
		return nil, false
	}
	utc := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	return &utc, true
}

func parseBookingDateEndQuery(c *gin.Context, key string) (*time.Time, bool) {
	raw := strings.TrimSpace(c.Query(key))
	if raw == "" {
		return nil, true
	}
	t, err := time.Parse("2006-01-02", raw)
	if err != nil {
		httputil.Error(c, http.StatusBadRequest, "INVALID_QUERY", "invalid "+key)
		return nil, false
	}
	utc := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC).Add(24 * time.Hour)
	return &utc, true
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

// GetAdminBooking handles GET /api/admin/bookings/:id.
func GetAdminBooking(deps AdminBookingsDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		bookingID, ok := parseObjectID(c, "id")
		if !ok {
			return
		}
		detail, err := deps.Service.GetByID(c.Request.Context(), bookingID)
		if err != nil {
			if errors.Is(err, booking.ErrBookingNotFound) {
				httputil.Error(c, http.StatusNotFound, "BOOKING_NOT_FOUND", "booking not found")
				return
			}
			httputil.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
			return
		}
		httputil.OK(c, detail)
	}
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

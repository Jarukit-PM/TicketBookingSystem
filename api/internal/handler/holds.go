package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/auth"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/hold"
	"github.com/Jarukit-PM/TicketBookingSystem/api/pkg/httputil"
)

// HoldsDeps holds dependencies for seat hold handlers.
type HoldsDeps struct {
	Holds *hold.Service
}

type holdsRequest struct {
	SeatIDs []string `json:"seatIds"`
}

// AddShowtimeHolds handles POST /api/showtimes/:id/holds.
func AddShowtimeHolds(deps HoldsDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := auth.UserFromContext(c)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var req holdsRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			httputil.Error(c, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
			return
		}

		result, err := deps.Holds.AddSeats(c.Request.Context(), user.ID.Hex(), c.Param("id"), req.SeatIDs)
		if err != nil {
			writeHoldError(c, err)
			return
		}

		httputil.OK(c, result)
	}
}

// RemoveShowtimeHolds handles DELETE /api/showtimes/:id/holds.
func RemoveShowtimeHolds(deps HoldsDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := auth.UserFromContext(c)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var req holdsRequest
		_ = c.ShouldBindJSON(&req)

		result, err := deps.Holds.RemoveSeats(c.Request.Context(), user.ID.Hex(), c.Param("id"), req.SeatIDs)
		if err != nil {
			writeHoldError(c, err)
			return
		}

		httputil.OK(c, result)
	}
}

func writeHoldError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, hold.ErrShowtimeNotFound):
		httputil.Error(c, http.StatusNotFound, "SHOWTIME_NOT_FOUND", "showtime not found")
	case errors.Is(err, hold.ErrShowtimeStarted):
		httputil.Error(c, http.StatusConflict, "SHOWTIME_STARTED", "showtime already started")
	case errors.Is(err, hold.ErrSeatNotFound):
		httputil.Error(c, http.StatusBadRequest, "SEAT_NOT_FOUND", "seat not found in layout")
	case errors.Is(err, hold.ErrSeatBlocked):
		httputil.Error(c, http.StatusConflict, "SEAT_BLOCKED", "seat is blocked")
	case errors.Is(err, hold.ErrSeatSold):
		httputil.Error(c, http.StatusConflict, "SEAT_SOLD", "seat is sold")
	case errors.Is(err, hold.ErrSeatHeldByOther):
		httputil.Error(c, http.StatusConflict, "SEAT_HELD", "seat held by another user")
	case errors.Is(err, hold.ErrSeatLimitExceeded):
		httputil.Error(c, http.StatusConflict, "SEAT_LIMIT_EXCEEDED", "maximum 10 seats per hold")
	case errors.Is(err, hold.ErrSeatNotHeld):
		httputil.Error(c, http.StatusConflict, "SEAT_NOT_HELD", "seat not in your holds")
	default:
		httputil.Error(c, http.StatusInternalServerError, "HOLD_ERROR", "failed to update holds")
	}
}

package handler

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/audit"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/auth"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/booking"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/tasks"
	"github.com/Jarukit-PM/TicketBookingSystem/api/pkg/httputil"
)

// BookingEventPublisher broadcasts seat_sold events after confirm.
type BookingEventPublisher interface {
	PublishSeatSold(ctx context.Context, showtimeID, seatID string) error
}

// BookingsDeps holds dependencies for booking handlers.
type BookingsDeps struct {
	Bookings  *booking.Service
	Tasks     tasks.Enqueuer
	Publisher BookingEventPublisher
	Audit     *audit.Logger
}

type confirmRequest struct {
	ShowtimeID string `json:"showtimeId"`
}

// ConfirmBooking handles POST /api/bookings/confirm.
func ConfirmBooking(deps BookingsDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := auth.UserFromContext(c)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		idempotencyKey := c.GetHeader("Idempotency-Key")
		if idempotencyKey == "" {
			httputil.Error(c, http.StatusBadRequest, "IDEMPOTENCY_KEY_REQUIRED", "Idempotency-Key header is required")
			return
		}

		var req confirmRequest
		if err := c.ShouldBindJSON(&req); err != nil || req.ShowtimeID == "" {
			httputil.Error(c, http.StatusBadRequest, "INVALID_BODY", "showtimeId is required")
			return
		}

		result, err := deps.Bookings.Confirm(c.Request.Context(), user.ID.Hex(), req.ShowtimeID, idempotencyKey)
		if err != nil {
			writeBookingError(c, deps.Audit, user.ID, req.ShowtimeID, err)
			return
		}

		if deps.Audit != nil {
			deps.Audit.BookingSuccess(
				c.Request.Context(),
				user.ID,
				result.ID.Hex(),
				req.ShowtimeID,
				result.BookingRef,
				result.Seats,
				result.Total,
			)
		}

		publishSeatSoldEvents(c.Request.Context(), deps.Publisher, req.ShowtimeID, result.Seats)

		if deps.Tasks != nil {
			if err := deps.Tasks.EnqueueEmailSend(c.Request.Context(), result.ID.Hex()); err != nil {
				log.Printf("enqueue email:send booking=%s: %v", result.ID.Hex(), err)
				if deps.Audit != nil {
					deps.Audit.SystemError(
						c.Request.Context(),
						user.ID,
						"booking",
						result.ID.Hex(),
						"EMAIL_ENQUEUE_FAILED",
						err.Error(),
					)
				}
			}
		}

		httputil.OK(c, toConfirmResponse(result))
	}
}

// ListMyBookings handles GET /api/bookings/mine.
func ListMyBookings(deps BookingsDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := auth.UserFromContext(c)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		upcoming := c.DefaultQuery("upcoming", "true") == "true"
		items, err := deps.Bookings.ListMine(c.Request.Context(), user.ID.Hex(), upcoming)
		if err != nil {
			writeBookingQueryError(c, err)
			return
		}
		httputil.OK(c, gin.H{"bookings": items})
	}
}

// GetBooking handles GET /api/bookings/:id.
func GetBooking(deps BookingsDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := auth.UserFromContext(c)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		item, err := deps.Bookings.GetDetail(c.Request.Context(), user.ID.Hex(), user.Role, c.Param("id"))
		if err != nil {
			writeBookingQueryError(c, err)
			return
		}
		httputil.OK(c, item)
	}
}

// GetBookingTicket handles GET /api/bookings/:id/ticket.
func GetBookingTicket(deps BookingsDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := auth.UserFromContext(c)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ticket, err := deps.Bookings.GetTicket(c.Request.Context(), user.ID.Hex(), c.Param("id"))
		if err != nil {
			writeTicketError(c, err)
			return
		}
		httputil.OK(c, ticket)
	}
}

func writeTicketError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, booking.ErrBookingNotFound):
		httputil.Error(c, http.StatusNotFound, "BOOKING_NOT_FOUND", "booking not found")
	case errors.Is(err, booking.ErrForbidden):
		httputil.Error(c, http.StatusForbidden, "FORBIDDEN", "not allowed to view this ticket")
	default:
		httputil.Error(c, http.StatusInternalServerError, "TICKET_ERROR", "failed to load ticket")
	}
}

func toConfirmResponse(b *booking.Booking) gin.H {
	return gin.H{
		"id":          b.ID.Hex(),
		"bookingRef":  b.BookingRef,
		"showtimeId":  b.ShowtimeID.Hex(),
		"seats":       b.Seats,
		"total":       b.Total,
		"status":      b.Status,
		"confirmedAt": b.ConfirmedAt,
	}
}

func publishSeatSoldEvents(ctx context.Context, pub BookingEventPublisher, showtimeID string, seatIDs []string) {
	if pub == nil {
		return
	}
	for _, seatID := range seatIDs {
		_ = pub.PublishSeatSold(ctx, showtimeID, seatID)
	}
}

func writeBookingError(c *gin.Context, auditLog *audit.Logger, userID primitive.ObjectID, showtimeID string, err error) {
	switch {
	case errors.Is(err, booking.ErrIdempotencyRequired):
		httputil.Error(c, http.StatusBadRequest, "IDEMPOTENCY_KEY_REQUIRED", "Idempotency-Key header is required")
	case errors.Is(err, booking.ErrNoActiveHolds):
		logBookingFailed(auditLog, c, userID, showtimeID, "NO_ACTIVE_HOLDS", "no active seat holds for this showtime")
		httputil.Error(c, http.StatusConflict, "NO_ACTIVE_HOLDS", "no active seat holds for this showtime")
	case errors.Is(err, booking.ErrSeatConflict):
		logBookingFailed(auditLog, c, userID, showtimeID, "SEAT_CONFLICT", "one or more seats are no longer available")
		httputil.Error(c, http.StatusConflict, "SEAT_CONFLICT", "one or more seats are no longer available")
	case errors.Is(err, booking.ErrShowtimeNotFound):
		httputil.Error(c, http.StatusNotFound, "SHOWTIME_NOT_FOUND", "showtime not found")
	case errors.Is(err, booking.ErrShowtimeStarted):
		logBookingFailed(auditLog, c, userID, showtimeID, "SHOWTIME_STARTED", "showtime already started")
		httputil.Error(c, http.StatusConflict, "SHOWTIME_STARTED", "showtime already started")
	case errors.Is(err, booking.ErrSeatLimitExceeded):
		logBookingFailed(auditLog, c, userID, showtimeID, "SEAT_LIMIT_EXCEEDED", "maximum 10 seats per booking")
		httputil.Error(c, http.StatusConflict, "SEAT_LIMIT_EXCEEDED", "maximum 10 seats per booking")
	default:
		if auditLog != nil {
			auditLog.SystemError(c.Request.Context(), userID, "booking", showtimeID, "CONFIRM_ERROR", err.Error())
		}
		httputil.Error(c, http.StatusInternalServerError, "CONFIRM_ERROR", "failed to confirm booking")
	}
}

func logBookingFailed(auditLog *audit.Logger, c *gin.Context, userID primitive.ObjectID, showtimeID, code, message string) {
	if auditLog == nil {
		return
	}
	auditLog.BookingFailed(c.Request.Context(), userID, showtimeID, code, message)
}

func writeBookingQueryError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, booking.ErrBookingNotFound):
		httputil.Error(c, http.StatusNotFound, "BOOKING_NOT_FOUND", "booking not found")
	case errors.Is(err, booking.ErrForbidden):
		httputil.Error(c, http.StatusForbidden, "FORBIDDEN", "you do not have access to this booking")
	default:
		httputil.Error(c, http.StatusInternalServerError, "BOOKING_QUERY_ERROR", "failed to load booking")
	}
}

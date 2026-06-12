package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/booking"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/tasks"
	"github.com/Jarukit-PM/TicketBookingSystem/api/pkg/httputil"
)

// AdminEmailDeps holds dependencies for admin email actions.
type AdminEmailDeps struct {
	Bookings booking.Repository
	Tasks    tasks.Enqueuer
}

// ResendBookingEmail handles POST /api/admin/bookings/:id/resend-email.
func ResendBookingEmail(deps AdminEmailDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			httputil.Error(c, http.StatusBadRequest, "INVALID_ID", "invalid booking id")
			return
		}
		b, err := deps.Bookings.FindByID(c.Request.Context(), id)
		if err != nil {
			httputil.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
			return
		}
		if b == nil || b.Status != booking.StatusConfirmed {
			httputil.Error(c, http.StatusNotFound, "BOOKING_NOT_FOUND", "confirmed booking not found")
			return
		}
		if deps.Tasks == nil {
			httputil.Error(c, http.StatusServiceUnavailable, "TASKS_UNAVAILABLE", "email queue unavailable")
			return
		}
		if err := deps.Tasks.EnqueueEmailSend(c.Request.Context(), id.Hex()); err != nil {
			log.Printf("admin resend email booking=%s: %v", id.Hex(), err)
			httputil.Error(c, http.StatusInternalServerError, "EMAIL_ENQUEUE_FAILED", "failed to queue confirmation email")
			return
		}
		httputil.OK(c, gin.H{"queued": true, "bookingId": id.Hex()})
	}
}

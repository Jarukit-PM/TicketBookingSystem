package audit

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Logger writes append-only audit records. Failures are logged but never block callers.
type Logger struct {
	repo AuditRepository
}

// NewLogger returns an audit logger backed by the given repository.
func NewLogger(repo AuditRepository) *Logger {
	return &Logger{repo: repo}
}

// BookingSuccess records a confirmed booking.
func (l *Logger) BookingSuccess(
	ctx context.Context,
	userID primitive.ObjectID,
	bookingID, showtimeID, bookingRef string,
	seats []string,
	total int64,
) {
	l.insert(ctx, userID, ActionBookingSuccess, "booking", bookingID, map[string]any{
		"showtimeId": showtimeID,
		"bookingRef": bookingRef,
		"seats":      seats,
		"total":      total,
	})
}

// BookingTimeout records a hold that expired without confirm.
func (l *Logger) BookingTimeout(ctx context.Context, userID primitive.ObjectID, showtimeID, seatID string) {
	l.insert(ctx, userID, ActionBookingTimeout, "showtime", showtimeID, map[string]any{
		"seatId": seatID,
		"reason": "hold_ttl_expired",
	})
}

// SeatReleased records an explicit user abandon or deselect.
func (l *Logger) SeatReleased(
	ctx context.Context,
	userID primitive.ObjectID,
	showtimeID string,
	seatIDs []string,
) {
	l.insert(ctx, userID, ActionSeatReleased, "showtime", showtimeID, map[string]any{
		"seatIds": seatIDs,
		"reason":  "user_released",
	})
}

// BookingFailed records a rejected confirm attempt (business rule, not server fault).
func (l *Logger) BookingFailed(
	ctx context.Context,
	userID primitive.ObjectID,
	showtimeID, code, message string,
) {
	l.insert(ctx, userID, ActionBookingFailed, "showtime", showtimeID, map[string]any{
		"code":    code,
		"message": message,
	})
}

// SystemError records an unexpected failure in booking or hold flows.
func (l *Logger) SystemError(
	ctx context.Context,
	userID primitive.ObjectID,
	entity, entityID, code, message string,
) {
	l.insert(ctx, userID, ActionSystemError, entity, entityID, map[string]any{
		"code":    code,
		"message": message,
	})
}

func (l *Logger) insert(
	ctx context.Context,
	actorID primitive.ObjectID,
	action, entity, entityID string,
	meta map[string]any,
) {
	if l == nil || l.repo == nil {
		return
	}
	if err := l.repo.InsertAuditLog(ctx, &AuditLog{
		ActorID:   ActorIDPtr(actorID),
		Action:    action,
		Entity:    entity,
		EntityID:  entityID,
		Meta:      meta,
		CreatedAt: time.Now().UTC(),
	}); err != nil {
		log.Printf("audit log %s %s/%s: %v", action, entity, entityID, err)
	}
}

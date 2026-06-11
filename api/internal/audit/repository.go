package audit

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuditRepository persists audit logs.
type AuditRepository interface {
	InsertAuditLog(ctx context.Context, log *AuditLog) error
	ListAuditLogs(ctx context.Context, limit int64) ([]AuditLog, error)
}

// EmailLogRepository persists email delivery logs.
type EmailLogRepository interface {
	InsertEmailLog(ctx context.Context, log *EmailLog) error
	ListByBooking(ctx context.Context, bookingID primitive.ObjectID) ([]EmailLog, error)
}

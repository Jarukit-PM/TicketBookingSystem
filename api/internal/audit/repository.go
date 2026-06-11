package audit

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// LogPage describes offset pagination for admin log lists.
type LogPage struct {
	Limit int64
	Skip  int64
}

// AuditRepository persists audit logs.
type AuditRepository interface {
	InsertAuditLog(ctx context.Context, log *AuditLog) error
	ListAuditLogs(ctx context.Context, page LogPage) ([]AuditLog, error)
}

// EmailLogRepository persists email delivery logs.
type EmailLogRepository interface {
	InsertEmailLog(ctx context.Context, log *EmailLog) error
	ListByBooking(ctx context.Context, bookingID primitive.ObjectID) ([]EmailLog, error)
	ListEmailLogs(ctx context.Context, page LogPage, bookingID *primitive.ObjectID) ([]EmailLog, error)
}

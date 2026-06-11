package admin

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/audit"
)

const (
	defaultLogLimit = 50
	maxLogLimit     = 100
)

// LogsService provides read-only admin log lookup.
type LogsService struct {
	AuditLogs audit.AuditRepository
	EmailLogs audit.EmailLogRepository
}

// ListAuditLogs returns paginated audit logs, newest first.
func (s *LogsService) ListAuditLogs(ctx context.Context, page, limit int) ([]audit.AuditLog, error) {
	logPage := audit.LogPage{Limit: int64(clampLimit(limit)), Skip: int64(skipFor(page, limit))}
	logs, err := s.AuditLogs.ListAuditLogs(ctx, logPage)
	if err != nil {
		return nil, fmt.Errorf("list audit logs: %w", err)
	}
	return logs, nil
}

// ListEmailLogs returns paginated email logs, optionally filtered by booking.
func (s *LogsService) ListEmailLogs(ctx context.Context, page, limit int, bookingID *primitive.ObjectID) ([]audit.EmailLog, error) {
	logPage := audit.LogPage{Limit: int64(clampLimit(limit)), Skip: int64(skipFor(page, limit))}
	logs, err := s.EmailLogs.ListEmailLogs(ctx, logPage, bookingID)
	if err != nil {
		return nil, fmt.Errorf("list email logs: %w", err)
	}
	return logs, nil
}

func clampLimit(limit int) int {
	if limit <= 0 {
		return defaultLogLimit
	}
	if limit > maxLogLimit {
		return maxLogLimit
	}
	return limit
}

func skipFor(page, limit int) int {
	if page <= 1 {
		return 0
	}
	return (page - 1) * clampLimit(limit)
}

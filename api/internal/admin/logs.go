package admin

import (
	"context"
	"fmt"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/audit"
)

// LogsService provides read-only admin log lookup.
type LogsService struct {
	AuditLogs audit.AuditRepository
	EmailLogs audit.EmailLogRepository
}

// ListAuditLogs returns paginated audit logs, newest first.
func (s *LogsService) ListAuditLogs(ctx context.Context, page, limit int, filter audit.AuditLogFilter) ([]audit.AuditLog, error) {
	logPage := audit.LogPage{Limit: int64(ClampPageLimit(limit)), Skip: int64(SkipFor(page, limit))}
	logs, err := s.AuditLogs.ListAuditLogs(ctx, logPage, filter)
	if err != nil {
		return nil, fmt.Errorf("list audit logs: %w", err)
	}
	return logs, nil
}

// ListEmailLogs returns paginated email logs, optionally filtered.
func (s *LogsService) ListEmailLogs(ctx context.Context, page, limit int, filter audit.EmailLogFilter) ([]audit.EmailLog, error) {
	logPage := audit.LogPage{Limit: int64(ClampPageLimit(limit)), Skip: int64(SkipFor(page, limit))}
	logs, err := s.EmailLogs.ListEmailLogs(ctx, logPage, filter)
	if err != nil {
		return nil, fmt.Errorf("list email logs: %w", err)
	}
	return logs, nil
}


package audit

import (
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EmailLogFilter narrows admin email log queries.
type EmailLogFilter struct {
	BookingID *primitive.ObjectID
	To        string
	Type      string
	Status    string
	SentFrom  *time.Time
	SentTo    *time.Time
}

// Matches reports whether log satisfies all non-empty filter fields.
func (f EmailLogFilter) Matches(log EmailLog) bool {
	if f.BookingID != nil && !f.BookingID.IsZero() && log.BookingID != *f.BookingID {
		return false
	}
	if f.To != "" && !strings.Contains(strings.ToLower(log.To), strings.ToLower(f.To)) {
		return false
	}
	if f.Type != "" && log.Type != f.Type {
		return false
	}
	if f.Status != "" && log.Status != f.Status {
		return false
	}
	if f.SentFrom != nil && log.CreatedAt.Before(*f.SentFrom) {
		return false
	}
	if f.SentTo != nil && !log.CreatedAt.Before(*f.SentTo) {
		return false
	}
	return true
}

package audit

import "go.mongodb.org/mongo-driver/bson/primitive"

// AuditLogFilter narrows admin audit log queries.
type AuditLogFilter struct {
	Action     string
	Entity     string
	EntityID   string
	ActorID    *primitive.ObjectID
	BookingRef string
}

// Matches reports whether log satisfies all non-empty filter fields.
func (f AuditLogFilter) Matches(log AuditLog) bool {
	if f.Action != "" && log.Action != f.Action {
		return false
	}
	if f.Entity != "" && log.Entity != f.Entity {
		return false
	}
	if f.EntityID != "" && log.EntityID != f.EntityID {
		return false
	}
	if f.ActorID != nil && !f.ActorID.IsZero() {
		if log.ActorID == nil || *log.ActorID != *f.ActorID {
			return false
		}
	}
	if f.BookingRef != "" {
		ref, _ := log.Meta["bookingRef"].(string)
		if ref != f.BookingRef {
			return false
		}
	}
	return true
}

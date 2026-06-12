package audit

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionAuditLogs = "audit_logs"
	CollectionEmailLogs = "email_logs"
)

const EmailTypeConfirmation = "CONFIRMATION"

// AuditLog is an append-only admin or system action record.
type AuditLog struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ActorID   *primitive.ObjectID `bson:"actorId,omitempty" json:"actorId,omitempty"`
	Action    string             `bson:"action" json:"action"`
	Entity    string             `bson:"entity" json:"entity"`
	EntityID  string             `bson:"entityId" json:"entityId"`
	Meta      map[string]any     `bson:"meta,omitempty" json:"meta,omitempty"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}

// EmailLog records an email provider send attempt for a booking event.
type EmailLog struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BookingID  primitive.ObjectID `bson:"bookingId" json:"bookingId"`
	Type       string             `bson:"type" json:"type"`
	To         string             `bson:"to" json:"to"`
	ProviderID string             `bson:"providerId,omitempty" json:"providerId,omitempty"`
	Status     string             `bson:"status" json:"status"`
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
}

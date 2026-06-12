package audit

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestAuditLogFilterMatches(t *testing.T) {
	actorID := primitive.NewObjectID()
	log := AuditLog{
		Action:   ActionBookingSuccess,
		Entity:   "booking",
		EntityID: "booking-1",
		ActorID:  ActorIDPtr(actorID),
		Meta:     map[string]any{"bookingRef": "TBS-123"},
	}

	actionFilter := AuditLogFilter{Action: ActionBookingSuccess}
	if !actionFilter.Matches(log) {
		t.Fatal("expected action filter to match")
	}
	wrongAction := AuditLogFilter{Action: ActionCreate}
	if wrongAction.Matches(log) {
		t.Fatal("expected mismatched action to fail")
	}
	refFilter := AuditLogFilter{BookingRef: "TBS-123"}
	if !refFilter.Matches(log) {
		t.Fatal("expected bookingRef filter to match meta")
	}
	actorFilter := AuditLogFilter{ActorID: &actorID}
	if !actorFilter.Matches(log) {
		t.Fatal("expected actorId filter to match")
	}
}

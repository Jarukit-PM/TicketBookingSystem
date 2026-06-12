package audit

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestEmailLogFilterMatches(t *testing.T) {
	bookingID := primitive.NewObjectID()
	sentAt := time.Date(2026, 6, 12, 10, 0, 0, 0, time.UTC)
	log := EmailLog{
		BookingID: bookingID,
		Type:      EmailTypeConfirmation,
		To:        "user@example.com",
		Status:    "SENT",
		CreatedAt: sentAt,
	}

	statusFilter := EmailLogFilter{Status: "SENT"}
	if !statusFilter.Matches(log) {
		t.Fatal("expected status filter to match")
	}
	toFilter := EmailLogFilter{To: "example.com"}
	if !toFilter.Matches(log) {
		t.Fatal("expected partial to filter to match")
	}
	bookingFilter := EmailLogFilter{BookingID: &bookingID}
	if !bookingFilter.Matches(log) {
		t.Fatal("expected bookingId filter to match")
	}
}

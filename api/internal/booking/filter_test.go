package booking_test

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/booking"
)

func TestConfirmedFilterMatches(t *testing.T) {
	userID := primitive.NewObjectID()
	showtimeID := primitive.NewObjectID()
	confirmedAt := time.Date(2026, 6, 12, 14, 30, 0, 0, time.UTC)
	from := time.Date(2026, 6, 12, 0, 0, 0, 0, time.UTC)
	to := from.Add(24 * time.Hour)

	b := booking.Booking{
		UserID:      userID,
		ShowtimeID:  showtimeID,
		BookingRef:  "TBS-TEST",
		Status:      booking.StatusConfirmed,
		Locale:      "th",
		ConfirmedAt: confirmedAt,
	}

	filter := booking.ConfirmedFilter{
		BookingRef:    "TBS-TEST",
		UserID:        userID,
		ShowtimeID:    showtimeID,
		Locale:        "th",
		ConfirmedFrom: &from,
		ConfirmedTo:   &to,
	}
	if !filter.Matches(b) {
		t.Fatal("expected booking to match combined filter")
	}

	filter.Locale = "en"
	if filter.Matches(b) {
		t.Fatal("expected locale mismatch to reject booking")
	}
}

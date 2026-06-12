package booking_test

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/booking"
)

func confirmedBookingFixture() (booking.Booking, primitive.ObjectID, primitive.ObjectID) {
	userID := primitive.NewObjectID()
	showtimeID := primitive.NewObjectID()
	confirmedAt := time.Date(2026, 6, 12, 14, 30, 0, 0, time.UTC)
	return booking.Booking{
		UserID:      userID,
		ShowtimeID:  showtimeID,
		BookingRef:  "TBS-TEST",
		Status:      booking.StatusConfirmed,
		Locale:      "th",
		ConfirmedAt: confirmedAt,
	}, userID, showtimeID
}

func TestConfirmedFilterMatches(t *testing.T) {
	b, userID, showtimeID := confirmedBookingFixture()
	from := time.Date(2026, 6, 12, 0, 0, 0, 0, time.UTC)
	to := from.Add(24 * time.Hour)

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

func TestConfirmedFilterRejectsNonConfirmed(t *testing.T) {
	b, userID, showtimeID := confirmedBookingFixture()
	b.Status = "PENDING"

	filter := booking.ConfirmedFilter{
		BookingRef: "TBS-TEST",
		UserID:     userID,
		ShowtimeID: showtimeID,
	}
	if filter.Matches(b) {
		t.Fatal("expected non-confirmed booking to be rejected")
	}
}

func TestConfirmedFilterShowtimeIDs(t *testing.T) {
	b, userID, showtimeID := confirmedBookingFixture()
	otherShowtime := primitive.NewObjectID()

	match := booking.ConfirmedFilter{
		UserID:      userID,
		ShowtimeIDs: []primitive.ObjectID{otherShowtime, showtimeID},
	}
	if !match.Matches(b) {
		t.Fatal("expected booking to match showtime ID list")
	}

	noMatch := booking.ConfirmedFilter{
		UserID:      userID,
		ShowtimeIDs: []primitive.ObjectID{otherShowtime},
	}
	if noMatch.Matches(b) {
		t.Fatal("expected booking outside showtime ID list to be rejected")
	}
}

func TestConfirmedFilterDateRangeExclusiveEnd(t *testing.T) {
	b, _, _ := confirmedBookingFixture()
	dayStart := time.Date(2026, 6, 12, 0, 0, 0, 0, time.UTC)
	dayEnd := dayStart.Add(24 * time.Hour)

	inRange := booking.ConfirmedFilter{
		ConfirmedFrom: &dayStart,
		ConfirmedTo:   &dayEnd,
	}
	if !inRange.Matches(b) {
		t.Fatal("expected booking within date range")
	}

	atEnd := b.ConfirmedAt
	onEnd := booking.ConfirmedFilter{ConfirmedTo: &atEnd}
	if onEnd.Matches(b) {
		t.Fatal("expected confirmedTo to be exclusive upper bound")
	}

	beforeStart := booking.ConfirmedFilter{
		ConfirmedFrom: ptrTime(dayStart.Add(24 * time.Hour)),
	}
	if beforeStart.Matches(b) {
		t.Fatal("expected booking before confirmedFrom to be rejected")
	}
}

func ptrTime(t time.Time) *time.Time {
	return &t
}

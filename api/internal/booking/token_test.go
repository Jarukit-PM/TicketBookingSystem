package booking_test

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/booking"
)

func TestSignAndValidateTicketToken(t *testing.T) {
	t.Parallel()

	secret := "test-ticket-secret"
	id := primitive.NewObjectID()
	ref := "TBS-ABC123"
	token := booking.SignTicketToken(secret, ref, id.Hex())

	b := &booking.Booking{
		ID:          id,
		BookingRef:  ref,
		TicketToken: token,
		Status:      booking.StatusConfirmed,
	}

	if !booking.ValidateTicketToken(ref, token, b, secret) {
		t.Fatal("expected valid token")
	}
	if booking.ValidateTicketToken(ref, token+"x", b, secret) {
		t.Fatal("expected tampered token to be rejected")
	}
	if booking.ValidateTicketToken("TBS-WRONG", token, b, secret) {
		t.Fatal("expected wrong ref to be rejected")
	}
	if booking.ValidateTicketToken(ref, token, b, "wrong-secret") {
		t.Fatal("expected wrong secret to be rejected")
	}
}

func TestTicketURL(t *testing.T) {
	t.Parallel()

	got := booking.TicketURL("http://localhost:5173", "TBS-XYZ", "abc123")
	if got != "http://localhost:5173/ticket/TBS-XYZ?t=abc123" {
		t.Fatalf("unexpected ticket url: %s", got)
	}
}

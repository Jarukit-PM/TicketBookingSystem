package email

import (
	"strings"
	"testing"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/booking"
)

func TestFormatTHBAmount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		satang int64
		want   string
	}{
		{2500, "25"},
		{22000, "220"},
		{1250, "12.50"},
	}
	for _, tc := range tests {
		if got := formatTHBAmount(tc.satang); got != tc.want {
			t.Fatalf("formatTHBAmount(%d) = %q, want %q", tc.satang, got, tc.want)
		}
	}
}

func TestRenderConfirmationEnglish(t *testing.T) {
	t.Parallel()

	html, text, err := renderConfirmation(booking.LocaleEN, confirmationData{
		BookingRef: "TBS-ABC",
		MovieTitle: "Film",
		CinemaName: "Cinema",
		ScreenName: "Hall 1",
		StartsAt:   "Mon, 01 Jan 2026 19:00:00 UTC",
		Seats:      "A-1, A-2",
		Total:      "500",
		TicketURL:  "http://localhost/ticket/TBS-ABC?t=1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "Booking confirmed") {
		t.Fatalf("expected English html heading, got %q", html)
	}
	if !strings.Contains(text, "Booking TBS-ABC") {
		t.Fatalf("expected English text, got %q", text)
	}
	if got := confirmationSubject(booking.LocaleEN, "Film"); got != "Your tickets — Film" {
		t.Fatalf("subject = %q", got)
	}
}

func TestRenderConfirmationThai(t *testing.T) {
	t.Parallel()

	html, text, err := renderConfirmation(booking.LocaleTH, confirmationData{
		BookingRef: "TBS-ABC",
		MovieTitle: "Film",
		CinemaName: "Cinema",
		ScreenName: "Hall 1",
		StartsAt:   "Mon, 01 Jan 2026 19:00:00 UTC",
		Seats:      "A-1",
		Total:      "500",
		TicketURL:  "http://localhost/ticket/TBS-ABC?t=1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "ยืนยันการจองแล้ว") {
		t.Fatalf("expected Thai html heading, got %q", html)
	}
	if !strings.Contains(text, "การจอง TBS-ABC") {
		t.Fatalf("expected Thai text, got %q", text)
	}
	if got := confirmationSubject(booking.LocaleTH, "Film"); got != "ตั๋วของคุณ — Film" {
		t.Fatalf("subject = %q", got)
	}
}

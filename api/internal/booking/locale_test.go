package booking_test

import (
	"testing"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/booking"
)

func TestParseLocale(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   string
		want string
	}{
		{"th", booking.LocaleTH},
		{"TH", booking.LocaleTH},
		{" th ", booking.LocaleTH},
		{"en", booking.LocaleEN},
		{"", booking.LocaleEN},
		{"fr", booking.LocaleEN},
	}
	for _, tc := range tests {
		if got := booking.ParseLocale(tc.in); got != tc.want {
			t.Errorf("ParseLocale(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

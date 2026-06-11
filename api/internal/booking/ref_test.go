package booking_test

import (
	"strings"
	"testing"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/booking"
)

func TestGenerateBookingRef_Format(t *testing.T) {
	for i := 0; i < 50; i++ {
		ref, err := booking.GenerateBookingRef()
		if err != nil {
			t.Fatalf("iteration %d: %v", i, err)
		}
		if !booking.ValidateBookingRefFormat(ref) {
			t.Fatalf("invalid format: %q", ref)
		}
		if !strings.HasPrefix(ref, "TBS-") {
			t.Fatalf("missing prefix: %q", ref)
		}
		suffix := strings.TrimPrefix(ref, "TBS-")
		if len(suffix) < 6 || len(suffix) > 8 {
			t.Fatalf("suffix length %d out of range: %q", len(suffix), ref)
		}
		for _, ch := range suffix {
			if strings.ContainsRune("01IO", ch) {
				t.Fatalf("ambiguous char %q in %q", ch, ref)
			}
		}
	}
}

func TestValidateBookingRefFormat(t *testing.T) {
	tests := []struct {
		name string
		ref  string
		want bool
	}{
		{name: "valid min length", ref: "TBS-ABCDEF", want: true},
		{name: "valid max length", ref: "TBS-ABCDEFGH", want: true},
		{name: "missing prefix", ref: "ABC-ABCDEF", want: false},
		{name: "too short", ref: "TBS-ABCDE", want: false},
		{name: "too long", ref: "TBS-ABCDEFGHI", want: false},
		{name: "ambiguous zero", ref: "TBS-ABCD0F", want: false},
		{name: "ambiguous O", ref: "TBS-ABCDOF", want: false},
		{name: "ambiguous one", ref: "TBS-ABCD1F", want: false},
		{name: "ambiguous I", ref: "TBS-ABCDIF", want: false},
		{name: "lowercase rejected", ref: "TBS-abcdef", want: false},
		{name: "empty", ref: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := booking.ValidateBookingRefFormat(tt.ref)
			if got != tt.want {
				t.Fatalf("ValidateBookingRefFormat(%q) = %v, want %v", tt.ref, got, tt.want)
			}
		})
	}
}

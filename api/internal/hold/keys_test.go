package hold

import "testing"

func TestParseHoldKey(t *testing.T) {
	t.Parallel()

	showtimeID := "507f1f77bcf86cd799439011"
	seatID := "A-12"
	key := SeatKey(showtimeID, seatID)

	gotShowtime, gotSeat, ok := ParseHoldKey(key)
	if !ok {
		t.Fatal("ParseHoldKey() ok = false, want true")
	}
	if gotShowtime != showtimeID || gotSeat != seatID {
		t.Fatalf("ParseHoldKey() = (%q, %q), want (%q, %q)", gotShowtime, gotSeat, showtimeID, seatID)
	}

	if _, _, ok := ParseHoldKey("user_holds:abc:def"); ok {
		t.Fatal("ParseHoldKey(user_holds) ok = true, want false")
	}
}

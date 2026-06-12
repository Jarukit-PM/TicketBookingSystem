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

func TestParseUserHoldsKey(t *testing.T) {
	t.Parallel()

	userID := "507f1f77bcf86cd799439011"
	showtimeID := "6a2ae94e92c8e0854b5a8e68"
	key := UserHoldsKey(userID, showtimeID)

	gotUser, gotShowtime, ok := ParseUserHoldsKey(key)
	if !ok {
		t.Fatal("ParseUserHoldsKey() ok = false, want true")
	}
	if gotUser != userID || gotShowtime != showtimeID {
		t.Fatalf("ParseUserHoldsKey() = (%q, %q), want (%q, %q)", gotUser, gotShowtime, userID, showtimeID)
	}

	if _, _, ok := ParseUserHoldsKey("hold:abc:def"); ok {
		t.Fatal("ParseUserHoldsKey(hold) ok = true, want false")
	}
}

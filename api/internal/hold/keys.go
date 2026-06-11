package hold

import (
	"fmt"
	"strings"
	"time"
)

const (
	keyPrefix       = "hold:"
	userHoldsPrefix = "user_holds:"
	defaultHoldTTL  = 5 * time.Minute
	maxSeatsPerHold = 10
)

// DefaultTTL is the hold lifetime (5 minutes).
const DefaultTTL = defaultHoldTTL

// MaxSeatsPerHold is the per-user per-showtime seat cap.
const MaxSeatsPerHold = maxSeatsPerHold

// SeatKey returns the Redis key for a seat hold on a showtime.
func SeatKey(showtimeID, seatID string) string {
	return fmt.Sprintf("%s%s:%s", keyPrefix, showtimeID, seatID)
}

// UserHoldsKey returns the Redis set key tracking a user's held seats on a showtime.
func UserHoldsKey(userID, showtimeID string) string {
	return fmt.Sprintf("%s%s:%s", userHoldsPrefix, userID, showtimeID)
}

// SeatKeyPattern returns a SCAN pattern for all holds on a showtime.
func SeatKeyPattern(showtimeID string) string {
	return SeatKey(showtimeID, "*")
}

// ParseSeatIDFromKey extracts the seat ID from a hold key for the given showtime.
func ParseSeatIDFromKey(key, showtimeID string) (string, bool) {
	prefix := keyPrefix + showtimeID + ":"
	if !strings.HasPrefix(key, prefix) {
		return "", false
	}
	seatID := strings.TrimPrefix(key, prefix)
	if seatID == "" {
		return "", false
	}
	return seatID, true
}

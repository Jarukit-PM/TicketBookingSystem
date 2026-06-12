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

// ParseUserHoldsKey parses user_holds:{userId}:{showtimeId}.
func ParseUserHoldsKey(key string) (userID, showtimeID string, ok bool) {
	if !strings.HasPrefix(key, userHoldsPrefix) {
		return "", "", false
	}
	rest := strings.TrimPrefix(key, userHoldsPrefix)
	i := strings.LastIndex(rest, ":")
	if i <= 0 || i >= len(rest)-1 {
		return "", "", false
	}
	userID = rest[:i]
	showtimeID = rest[i+1:]
	if userID == "" || showtimeID == "" {
		return "", "", false
	}
	return userID, showtimeID, true
}

// UserHoldsKey returns the Redis set key tracking a user's held seats on a showtime.
func UserHoldsKey(userID, showtimeID string) string {
	return fmt.Sprintf("%s%s:%s", userHoldsPrefix, userID, showtimeID)
}

// UserHoldsKeyPattern returns a SCAN pattern for all user hold sets on a showtime.
func UserHoldsKeyPattern(showtimeID string) string {
	return fmt.Sprintf("%s*:%s", userHoldsPrefix, showtimeID)
}

// SeatKeyPattern returns a SCAN pattern for all holds on a showtime.
func SeatKeyPattern(showtimeID string) string {
	return SeatKey(showtimeID, "*")
}

// ParseHoldKey parses hold:{showtimeId}:{seatId}.
func ParseHoldKey(key string) (showtimeID, seatID string, ok bool) {
	if !strings.HasPrefix(key, keyPrefix) {
		return "", "", false
	}
	rest := strings.TrimPrefix(key, keyPrefix)
	i := strings.Index(rest, ":")
	if i <= 0 || i >= len(rest)-1 {
		return "", "", false
	}
	showtimeID = rest[:i]
	seatID = rest[i+1:]
	if showtimeID == "" || seatID == "" {
		return "", "", false
	}
	return showtimeID, seatID, true
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

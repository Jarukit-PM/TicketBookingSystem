package catalog

import (
	"fmt"
	"time"
)

// NowInTimezone returns the current instant expressed in the given IANA timezone.
// Comparison with stored UTC showtime startsAt values uses the same instant.
func NowInTimezone(timezone string) (time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("load timezone %q: %w", timezone, err)
	}
	return time.Now().In(loc), nil
}

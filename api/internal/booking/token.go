package booking

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

const ticketTokenBytes = 32

// GenerateTicketToken returns an opaque secret for QR validation (distinct from bookingRef).
func GenerateTicketToken() (string, error) {
	buf := make([]byte, ticketTokenBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("random bytes: %w", err)
	}
	return hex.EncodeToString(buf), nil
}

package booking

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SignTicketToken returns an HMAC-SHA256 hex token for a confirmed booking.
func SignTicketToken(secret, bookingRef, bookingID string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(bookingRef))
	_, _ = mac.Write([]byte("|"))
	_, _ = mac.Write([]byte(bookingID))
	return hex.EncodeToString(mac.Sum(nil))
}

// GenerateTicketToken signs a ticket token for a booking about to be inserted.
func GenerateTicketToken(secret, bookingRef string, bookingID primitive.ObjectID) string {
	return SignTicketToken(secret, bookingRef, bookingID.Hex())
}

// ValidateTicketToken reports whether token matches the booking for ref.
func ValidateTicketToken(ref, token string, b *Booking, secret string) bool {
	if b == nil || ref == "" || token == "" || ref != b.BookingRef || b.Status != StatusConfirmed {
		return false
	}
	expected := SignTicketToken(secret, ref, b.ID.Hex())
	if subtle.ConstantTimeCompare([]byte(token), []byte(expected)) == 1 {
		return true
	}
	// Accept the persisted token (email links use this value).
	if b.TicketToken != "" {
		return subtle.ConstantTimeCompare([]byte(token), []byte(b.TicketToken)) == 1
	}
	return false
}

// TicketURL builds the public ticket link encoded in QR codes.
func TicketURL(appURL, bookingRef, ticketToken string) string {
	base := strings.TrimRight(appURL, "/")
	u, err := url.Parse(fmt.Sprintf("%s/ticket/%s", base, url.PathEscape(bookingRef)))
	if err != nil {
		return fmt.Sprintf("%s/ticket/%s?t=%s", base, bookingRef, url.QueryEscape(ticketToken))
	}
	q := u.Query()
	q.Set("t", ticketToken)
	u.RawQuery = q.Encode()
	return u.String()
}

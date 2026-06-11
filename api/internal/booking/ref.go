package booking

import (
	"crypto/rand"
	"fmt"
	"strings"
)

const bookingRefPrefix = "TBS-"

// bookingRefAlphabet excludes ambiguous characters 0/O and 1/I.
const bookingRefAlphabet = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"

const (
	bookingRefMinLen = 6
	bookingRefMaxLen = 8
)

// GenerateBookingRef returns a human-readable booking reference (TBS- + 6–8 unambiguous alphanumerics).
func GenerateBookingRef() (string, error) {
	b, err := secureRandomByte()
	if err != nil {
		return "", err
	}
	length := bookingRefMinLen + int(b)%(bookingRefMaxLen-bookingRefMinLen+1)
	suffix, err := randomString(length)
	if err != nil {
		return "", err
	}
	return bookingRefPrefix + suffix, nil
}

func randomString(n int) (string, error) {
	buf := make([]byte, n)
	alphabetLen := byte(len(bookingRefAlphabet))
	for i := range buf {
		b, err := secureRandomByte()
		if err != nil {
			return "", fmt.Errorf("random byte: %w", err)
		}
		buf[i] = bookingRefAlphabet[b%alphabetLen]
	}
	return string(buf), nil
}

func secureRandomByte() (byte, error) {
	var b [1]byte
	if _, err := rand.Read(b[:]); err != nil {
		return 0, err
	}
	return b[0], nil
}

// ValidateBookingRefFormat reports whether ref matches the TBS- customer-facing format.
func ValidateBookingRefFormat(ref string) bool {
	if !strings.HasPrefix(ref, bookingRefPrefix) {
		return false
	}
	suffix := strings.TrimPrefix(ref, bookingRefPrefix)
	if len(suffix) < bookingRefMinLen || len(suffix) > bookingRefMaxLen {
		return false
	}
	for _, ch := range suffix {
		if !strings.ContainsRune(bookingRefAlphabet, ch) {
			return false
		}
	}
	return true
}

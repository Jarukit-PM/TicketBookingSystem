package booking

import "strings"

const (
	LocaleEN = "en"
	LocaleTH = "th"
)

// ParseLocale normalizes a client locale header value to en or th.
func ParseLocale(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case LocaleTH:
		return LocaleTH
	default:
		return LocaleEN
	}
}

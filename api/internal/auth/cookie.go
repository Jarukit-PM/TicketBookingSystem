package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CookieOptions configures the session cookie.
type CookieOptions struct {
	Secure bool
}

// SetSessionCookie writes the httpOnly JWT session cookie.
func SetSessionCookie(c *gin.Context, token string, expiresAt time.Time, opts CookieOptions) {
	maxAge := int(time.Until(expiresAt).Seconds())
	if maxAge < 0 {
		maxAge = 0
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(CookieName, token, maxAge, "/", "", opts.Secure, true)
}

// ClearSessionCookie removes the session cookie.
func ClearSessionCookie(c *gin.Context, opts CookieOptions) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(CookieName, "", -1, "/", "", opts.Secure, true)
}

// SessionTokenFromRequest reads the session JWT from the request cookie.
func SessionTokenFromRequest(c *gin.Context) (string, error) {
	return c.Cookie(CookieName)
}

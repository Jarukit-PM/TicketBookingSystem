package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/user"
)

// MiddlewareDeps holds dependencies for auth middleware.
type MiddlewareDeps struct {
	Service *Service
}

// RequireAuth validates the session cookie and attaches the user to context.
func RequireAuth(deps MiddlewareDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := SessionTokenFromRequest(c)
		if err != nil || token == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		sessionUser, err := deps.Service.ParseSession(token)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		SetUser(c, sessionUser)
		c.Next()
	}
}

// RequireAdmin requires an authenticated admin user.
func RequireAdmin(deps MiddlewareDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionUser, ok := UserFromContext(c)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if sessionUser.Role != user.RoleAdmin {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}

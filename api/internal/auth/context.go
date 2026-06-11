package auth

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type contextKey string

const userContextKey contextKey = "authUser"

// SessionUser is the authenticated user attached to the request context.
type SessionUser struct {
	ID   primitive.ObjectID
	Role string
}

// SetUser stores the authenticated user on the Gin context.
func SetUser(c *gin.Context, user SessionUser) {
	c.Set(string(userContextKey), user)
}

// UserFromContext returns the authenticated user from the Gin context.
func UserFromContext(c *gin.Context) (SessionUser, bool) {
	value, ok := c.Get(string(userContextKey))
	if !ok {
		return SessionUser{}, false
	}
	user, ok := value.(SessionUser)
	return user, ok
}

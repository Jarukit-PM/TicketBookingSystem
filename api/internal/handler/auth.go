package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/auth"
	"github.com/Jarukit-PM/TicketBookingSystem/api/pkg/httputil"
)

// AuthDeps holds auth handler dependencies.
type AuthDeps struct {
	Service       *auth.Service
	CookieOptions auth.CookieOptions
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register handles POST /api/auth/register.
func Register(deps AuthDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req registerRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			httputil.Error(c, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
			return
		}

		result, err := deps.Service.Register(c.Request.Context(), auth.RegisterInput{
			Email:    req.Email,
			Password: req.Password,
			Name:     req.Name,
		})
		if err != nil {
			writeAuthError(c, err)
			return
		}

		auth.SetSessionCookie(c, result.Token, result.Exp, deps.CookieOptions)
		c.JSON(http.StatusCreated, gin.H{"user": result.User})
	}
}

// Login handles POST /api/auth/login.
func Login(deps AuthDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req loginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			httputil.Error(c, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
			return
		}

		result, err := deps.Service.Login(c.Request.Context(), auth.LoginInput{
			Email:    req.Email,
			Password: req.Password,
		})
		if err != nil {
			writeAuthError(c, err)
			return
		}

		auth.SetSessionCookie(c, result.Token, result.Exp, deps.CookieOptions)
		httputil.OK(c, gin.H{"user": result.User})
	}
}

// Logout handles POST /api/auth/logout.
func Logout(deps AuthDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth.ClearSessionCookie(c, deps.CookieOptions)
		httputil.OK(c, gin.H{"ok": true})
	}
}

// Me handles GET /api/auth/me.
func Me(deps AuthDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionUser, ok := auth.UserFromContext(c)
		if !ok {
			httputil.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
			return
		}

		user, err := deps.Service.Me(c.Request.Context(), sessionUser.ID)
		if err != nil {
			writeAuthError(c, err)
			return
		}

		httputil.OK(c, gin.H{"user": user})
	}
}

func writeAuthError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, auth.ErrInvalidEmail),
		errors.Is(err, auth.ErrInvalidPassword),
		errors.Is(err, auth.ErrInvalidName):
		httputil.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	case errors.Is(err, auth.ErrEmailTaken):
		httputil.Error(c, http.StatusConflict, "EMAIL_TAKEN", err.Error())
	case errors.Is(err, auth.ErrInvalidCredentials):
		httputil.Error(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", err.Error())
	case errors.Is(err, auth.ErrTooManyAttempts):
		httputil.Error(c, http.StatusTooManyRequests, "TOO_MANY_ATTEMPTS", err.Error())
	case errors.Is(err, auth.ErrUnauthorized):
		httputil.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
	default:
		httputil.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
	}
}

package handler

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/auth"
	"github.com/Jarukit-PM/TicketBookingSystem/api/pkg/httputil"
)

// AuthDeps holds auth handler dependencies.
type AuthDeps struct {
	Service       *auth.Service
	Google        *auth.GoogleOAuth
	AppURL        string
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

// GoogleStart handles GET /api/auth/google.
func GoogleStart(deps AuthDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		if deps.Google == nil || !deps.Google.Enabled() {
			httputil.Error(c, http.StatusServiceUnavailable, "OAUTH_NOT_CONFIGURED", auth.ErrGoogleNotConfigured.Error())
			return
		}

		state, err := auth.NewOAuthState()
		if err != nil {
			httputil.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
			return
		}

		auth.SetOAuthStateCookie(c, state, deps.CookieOptions)
		c.Redirect(http.StatusFound, deps.Google.AuthCodeURL(state))
	}
}

// GoogleCallback handles GET /api/auth/google/callback.
func GoogleCallback(deps AuthDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		if deps.Google == nil || !deps.Google.Enabled() {
			redirectOAuthError(c, deps.AppURL, "oauth_not_configured")
			return
		}

		if c.Query("error") != "" {
			redirectOAuthError(c, deps.AppURL, "oauth_denied")
			return
		}

		state, err := auth.OAuthStateFromRequest(c)
		if err != nil || state != c.Query("state") {
			redirectOAuthError(c, deps.AppURL, "oauth_state_invalid")
			return
		}
		auth.ClearOAuthStateCookie(c, deps.CookieOptions)

		code := c.Query("code")
		if code == "" {
			redirectOAuthError(c, deps.AppURL, "oauth_code_missing")
			return
		}

		profile, err := deps.Google.ExchangeCode(c.Request.Context(), code)
		if err != nil {
			redirectOAuthError(c, deps.AppURL, "oauth_exchange_failed")
			return
		}

		result, err := deps.Service.LoginWithGoogle(c.Request.Context(), profile)
		if err != nil {
			redirectOAuthError(c, deps.AppURL, oauthErrorCode(err))
			return
		}

		auth.SetSessionCookie(c, result.Token, result.Exp, deps.CookieOptions)
		c.Redirect(http.StatusFound, deps.AppURL)
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

func redirectOAuthError(c *gin.Context, appURL, code string) {
	target := strings.TrimRight(appURL, "/") + "/login?error=" + url.QueryEscape(code)
	c.Redirect(http.StatusFound, target)
}

func oauthErrorCode(err error) string {
	switch {
	case errors.Is(err, auth.ErrGoogleEmailNotVerified):
		return "oauth_email_unverified"
	case errors.Is(err, auth.ErrGoogleAccountConflict):
		return "oauth_account_conflict"
	default:
		return "oauth_failed"
	}
}

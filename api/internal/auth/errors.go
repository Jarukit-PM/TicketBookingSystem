package auth

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailTaken         = errors.New("email already registered")
	ErrInvalidEmail       = errors.New("invalid email")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrInvalidName        = errors.New("invalid name")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrTooManyAttempts        = errors.New("too many login attempts")
	ErrGoogleNotConfigured    = errors.New("google oauth is not configured")
	ErrInvalidOAuthState      = errors.New("invalid oauth state")
	ErrGoogleEmailNotVerified = errors.New("google email is not verified")
	ErrGoogleProfileInvalid   = errors.New("invalid google profile")
	ErrGoogleAccountConflict  = errors.New("google account conflict")
)

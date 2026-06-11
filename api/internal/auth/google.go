package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/user"
)

const oauthStateCookieName = "oauth_state"

// GoogleProfile is the verified identity returned by Google after OAuth.
type GoogleProfile struct {
	ID            string
	Email         string
	Name          string
	EmailVerified bool
}

// GoogleOAuth configures the Google authorization code flow.
type GoogleOAuth struct {
	config  *oauth2.Config
	enabled bool
}

// NewGoogleOAuth returns a Google OAuth client. When client ID or secret is empty,
// OAuth routes respond with service unavailable.
func NewGoogleOAuth(clientID, clientSecret, redirectURL string) *GoogleOAuth {
	if strings.TrimSpace(clientID) == "" || strings.TrimSpace(clientSecret) == "" {
		return &GoogleOAuth{enabled: false}
	}

	return &GoogleOAuth{
		enabled: true,
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
		},
	}
}

// Enabled reports whether Google OAuth credentials are configured.
func (g *GoogleOAuth) Enabled() bool {
	return g != nil && g.enabled
}

// AuthCodeURL returns the Google authorization URL for the given CSRF state.
func (g *GoogleOAuth) AuthCodeURL(state string) string {
	return g.config.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

// ExchangeCode trades an authorization code for a Google profile.
func (g *GoogleOAuth) ExchangeCode(ctx context.Context, code string) (GoogleProfile, error) {
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return GoogleProfile{}, fmt.Errorf("exchange oauth code: %w", err)
	}

	client := g.config.Client(ctx, token)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return GoogleProfile{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return GoogleProfile{}, fmt.Errorf("fetch google userinfo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return GoogleProfile{}, fmt.Errorf("google userinfo status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var info struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		Name          string `json:"name"`
		EmailVerified bool   `json:"verified_email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return GoogleProfile{}, fmt.Errorf("decode google userinfo: %w", err)
	}

	return GoogleProfile{
		ID:            info.ID,
		Email:         info.Email,
		Name:          strings.TrimSpace(info.Name),
		EmailVerified: info.EmailVerified,
	}, nil
}

// NewOAuthState returns a URL-safe CSRF token for the OAuth redirect.
func NewOAuthState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// SetOAuthStateCookie stores the OAuth CSRF state in a short-lived httpOnly cookie.
func SetOAuthStateCookie(c CookieWriter, state string, opts CookieOptions) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(oauthStateCookieName, state, 600, "/", "", opts.Secure, true)
}

// ClearOAuthStateCookie removes the OAuth CSRF cookie.
func ClearOAuthStateCookie(c CookieWriter, opts CookieOptions) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(oauthStateCookieName, "", -1, "/", "", opts.Secure, true)
}

// OAuthStateFromRequest reads the OAuth state cookie.
func OAuthStateFromRequest(c CookieReader) (string, error) {
	state, err := c.Cookie(oauthStateCookieName)
	if err != nil || state == "" {
		return "", ErrInvalidOAuthState
	}
	return state, nil
}

// CookieWriter sets cookies on a response.
type CookieWriter interface {
	SetSameSite(http.SameSite)
	SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool)
}

// CookieReader reads cookies from a request.
type CookieReader interface {
	Cookie(name string) (string, error)
}

// LoginWithGoogle signs in via Google OAuth, auto-linking by verified email when needed.
func (s *Service) LoginWithGoogle(ctx context.Context, profile GoogleProfile) (SessionResult, error) {
	if !profile.EmailVerified {
		return SessionResult{}, ErrGoogleEmailNotVerified
	}

	email := normalizeEmail(profile.Email)
	if err := validateEmail(email); err != nil {
		return SessionResult{}, err
	}
	if strings.TrimSpace(profile.ID) == "" {
		return SessionResult{}, ErrGoogleProfileInvalid
	}

	if byGoogle, err := s.users.FindByGoogleID(ctx, profile.ID); err != nil {
		return SessionResult{}, err
	} else if byGoogle != nil {
		return s.issueSession(byGoogle)
	}

	existing, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return SessionResult{}, err
	}
	if existing != nil {
		if existing.GoogleID != "" && existing.GoogleID != profile.ID {
			return SessionResult{}, ErrGoogleAccountConflict
		}
		existing.GoogleID = profile.ID
		if strings.TrimSpace(existing.Name) == "" && profile.Name != "" {
			existing.Name = profile.Name
		}
		if s.adminEmail != "" && email == s.adminEmail && existing.Role != user.RoleAdmin {
			existing.Role = user.RoleAdmin
		}
		if err := s.users.Update(ctx, existing); err != nil {
			return SessionResult{}, err
		}
		return s.issueSession(existing)
	}

	name := profile.Name
	if name == "" {
		name = strings.Split(email, "@")[0]
	}

	role := user.RoleCustomer
	if s.adminEmail != "" && email == s.adminEmail {
		role = user.RoleAdmin
	}

	u := &user.User{
		Email:    email,
		GoogleID: profile.ID,
		Name:     name,
		Role:     role,
	}
	if err := s.users.Insert(ctx, u); err != nil {
		return SessionResult{}, err
	}

	return s.issueSession(u)
}

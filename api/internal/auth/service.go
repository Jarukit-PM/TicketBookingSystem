package auth

import (
	"context"
	"net/mail"
	"strings"
	"time"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const minPasswordLength = 8

// UserDTO is the public user shape returned by auth endpoints.
type UserDTO struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

// RegisterInput is the register request payload.
type RegisterInput struct {
	Email    string
	Password string
	Name     string
}

// LoginInput is the login request payload.
type LoginInput struct {
	Email    string
	Password string
}

// SessionResult is returned after login or register.
type SessionResult struct {
	User  UserDTO
	Token string
	Exp   time.Time
}

// Service handles registration, login, and session introspection.
type Service struct {
	users       user.Repository
	tokens      *TokenService
	rateLimiter *LoginRateLimiter
	adminEmail  string
}

// NewService returns an auth service.
func NewService(
	users user.Repository,
	tokens *TokenService,
	rateLimiter *LoginRateLimiter,
	adminEmail string,
) *Service {
	return &Service{
		users:       users,
		tokens:      tokens,
		rateLimiter: rateLimiter,
		adminEmail:  normalizeEmail(adminEmail),
	}
}

// Register creates a user account and returns a session token.
func (s *Service) Register(ctx context.Context, input RegisterInput) (SessionResult, error) {
	email := normalizeEmail(input.Email)
	if err := validateEmail(email); err != nil {
		return SessionResult{}, err
	}
	if err := validatePassword(input.Password); err != nil {
		return SessionResult{}, err
	}
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return SessionResult{}, ErrInvalidName
	}

	existing, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return SessionResult{}, err
	}
	if existing != nil {
		return SessionResult{}, ErrEmailTaken
	}

	hash, err := HashPassword(input.Password)
	if err != nil {
		return SessionResult{}, err
	}

	role := user.RoleCustomer
	if s.adminEmail != "" && email == s.adminEmail {
		role = user.RoleAdmin
	}

	u := &user.User{
		Email:        email,
		PasswordHash: hash,
		Name:         name,
		Role:         role,
	}
	if err := s.users.Insert(ctx, u); err != nil {
		return SessionResult{}, err
	}

	return s.issueSession(u)
}

// Login validates credentials and returns a session token.
func (s *Service) Login(ctx context.Context, input LoginInput) (SessionResult, error) {
	email := normalizeEmail(input.Email)
	if err := validateEmail(email); err != nil {
		return SessionResult{}, ErrInvalidCredentials
	}
	if strings.TrimSpace(input.Password) == "" {
		return SessionResult{}, ErrInvalidCredentials
	}

	allowed, err := s.rateLimiter.Allow(ctx, email)
	if err != nil {
		return SessionResult{}, err
	}
	if !allowed {
		return SessionResult{}, ErrTooManyAttempts
	}

	u, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return SessionResult{}, err
	}
	if u == nil || u.PasswordHash == "" || !CheckPassword(u.PasswordHash, input.Password) {
		if recErr := s.rateLimiter.RecordFailure(ctx, email); recErr != nil {
			return SessionResult{}, recErr
		}
		return SessionResult{}, ErrInvalidCredentials
	}

	if err := s.rateLimiter.Reset(ctx, email); err != nil {
		return SessionResult{}, err
	}

	if s.adminEmail != "" && email == s.adminEmail && u.Role != user.RoleAdmin {
		u.Role = user.RoleAdmin
		if err := s.users.Update(ctx, u); err != nil {
			return SessionResult{}, err
		}
	}

	return s.issueSession(u)
}

// Me returns the current user by ID.
func (s *Service) Me(ctx context.Context, userID primitive.ObjectID) (UserDTO, error) {
	u, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return UserDTO{}, err
	}
	if u == nil {
		return UserDTO{}, ErrUnauthorized
	}
	return toUserDTO(u), nil
}

// ParseSession validates a JWT and returns the session user.
func (s *Service) ParseSession(token string) (SessionUser, error) {
	claims, err := s.tokens.Parse(token)
	if err != nil {
		return SessionUser{}, ErrUnauthorized
	}

	userID, err := primitive.ObjectIDFromHex(claims.Subject)
	if err != nil {
		return SessionUser{}, ErrUnauthorized
	}

	return SessionUser{
		ID:   userID,
		Role: claims.Role,
	}, nil
}

// BootstrapAdmin ensures the configured admin email has admin role.
func (s *Service) BootstrapAdmin(ctx context.Context, email, password string) error {
	if s.adminEmail == "" {
		return nil
	}

	normalized := normalizeEmail(email)
	if normalized != s.adminEmail {
		return nil
	}

	u, err := s.users.FindByEmail(ctx, normalized)
	if err != nil {
		return err
	}

	hash, err := HashPassword(password)
	if err != nil {
		return err
	}

	if u == nil {
		return s.users.Insert(ctx, &user.User{
			Email:        normalized,
			PasswordHash: hash,
			Name:         "Admin",
			Role:         user.RoleAdmin,
		})
	}

	if u.Role != user.RoleAdmin {
		u.Role = user.RoleAdmin
	}
	if u.PasswordHash == "" {
		u.PasswordHash = hash
	}
	return s.users.Update(ctx, u)
}

func (s *Service) issueSession(u *user.User) (SessionResult, error) {
	token, exp, err := s.tokens.Issue(u.ID, u.Role)
	if err != nil {
		return SessionResult{}, err
	}
	return SessionResult{
		User:  toUserDTO(u),
		Token: token,
		Exp:   exp,
	}, nil
}

func toUserDTO(u *user.User) UserDTO {
	return UserDTO{
		ID:    u.ID.Hex(),
		Email: u.Email,
		Name:  u.Name,
		Role:  u.Role,
	}
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func validateEmail(email string) error {
	if email == "" {
		return ErrInvalidEmail
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return ErrInvalidEmail
	}
	return nil
}

func validatePassword(password string) error {
	if len(password) < minPasswordLength {
		return ErrInvalidPassword
	}
	return nil
}

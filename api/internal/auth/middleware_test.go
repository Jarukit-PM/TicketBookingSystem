package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type stubUserRepo struct {
	users map[string]*user.User
}

func (s *stubUserRepo) Insert(ctx context.Context, u *user.User) error {
	if u.ID.IsZero() {
		u.ID = primitive.NewObjectID()
	}
	s.users[u.Email] = u
	return nil
}

func (s *stubUserRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*user.User, error) {
	for _, u := range s.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, nil
}

func (s *stubUserRepo) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	u, ok := s.users[email]
	if !ok {
		return nil, nil
	}
	return u, nil
}

func (s *stubUserRepo) FindByGoogleID(ctx context.Context, googleID string) (*user.User, error) {
	return nil, nil
}

func (s *stubUserRepo) Update(ctx context.Context, u *user.User) error {
	s.users[u.Email] = u
	return nil
}

func newTestService(t *testing.T, adminEmail string) *Service {
	t.Helper()
	return NewService(
		&stubUserRepo{users: map[string]*user.User{}},
		NewTokenService("test-secret", time.Hour),
		NewLoginRateLimiter(nil),
		adminEmail,
	)
}

func TestRegisterAssignsAdminForConfiguredEmail(t *testing.T) {
	svc := newTestService(t, "admin@example.com")

	result, err := svc.Register(context.Background(), RegisterInput{
		Email:    "admin@example.com",
		Password: "password123",
		Name:     "Admin User",
	})
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	if result.User.Role != user.RoleAdmin {
		t.Errorf("Role = %q, want %q", result.User.Role, user.RoleAdmin)
	}
}

func TestRequireAdminRejectsCustomer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := newTestService(t, "")
	userID := primitive.NewObjectID()
	token, _, err := svc.tokens.Issue(userID, user.RoleCustomer)
	if err != nil {
		t.Fatalf("Issue() error = %v", err)
	}

	r := gin.New()
	r.GET("/admin", RequireAuth(MiddlewareDeps{Service: svc}), RequireAdmin(MiddlewareDeps{Service: svc}), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.AddCookie(&http.Cookie{Name: CookieName, Value: token})
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusForbidden)
	}
}

func TestRequireAdminAllowsAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := newTestService(t, "")
	userID := primitive.NewObjectID()
	token, _, err := svc.tokens.Issue(userID, user.RoleAdmin)
	if err != nil {
		t.Fatalf("Issue() error = %v", err)
	}

	r := gin.New()
	r.GET("/admin", RequireAuth(MiddlewareDeps{Service: svc}), RequireAdmin(MiddlewareDeps{Service: svc}), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.AddCookie(&http.Cookie{Name: CookieName, Value: token})
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestRequireAuthRejectsMissingCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := newTestService(t, "")
	r := gin.New()
	r.GET("/protected", RequireAuth(MiddlewareDeps{Service: svc}), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

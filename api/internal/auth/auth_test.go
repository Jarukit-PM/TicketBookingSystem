package auth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/auth"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/user"
)

type memUserRepo struct {
	byEmail map[string]*user.User
}

func newMemUserRepo() *memUserRepo {
	return &memUserRepo{byEmail: make(map[string]*user.User)}
}

func (r *memUserRepo) Insert(_ context.Context, u *user.User) error {
	u.ID = primitive.NewObjectID()
	r.byEmail[u.Email] = u
	return nil
}

func (r *memUserRepo) FindByID(_ context.Context, id primitive.ObjectID) (*user.User, error) {
	for _, u := range r.byEmail {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, nil
}

func (r *memUserRepo) FindByEmail(_ context.Context, email string) (*user.User, error) {
	u, ok := r.byEmail[email]
	if !ok {
		return nil, nil
	}
	return u, nil
}

func (r *memUserRepo) FindByGoogleID(_ context.Context, googleID string) (*user.User, error) {
	for _, u := range r.byEmail {
		if u.GoogleID == googleID {
			return u, nil
		}
	}
	return nil, nil
}

func (r *memUserRepo) Update(_ context.Context, u *user.User) error {
	r.byEmail[u.Email] = u
	return nil
}

func TestPasswordHashAndVerify(t *testing.T) {
	hash, err := auth.HashPassword("secret123")
	if err != nil {
		t.Fatal(err)
	}
	if !auth.CheckPassword(hash, "secret123") {
		t.Fatal("expected password to match")
	}
	if auth.CheckPassword(hash, "wrong") {
		t.Fatal("expected wrong password to fail")
	}
}

func TestJWTIssueAndParse(t *testing.T) {
	tokens := auth.NewTokenService("test-secret", time.Hour)
	id := primitive.NewObjectID()
	token, _, err := tokens.Issue(id, user.RoleCustomer)
	if err != nil {
		t.Fatal(err)
	}
	claims, err := tokens.Parse(token)
	if err != nil {
		t.Fatal(err)
	}
	if claims.Subject != id.Hex() {
		t.Fatalf("subject = %q, want %q", claims.Subject, id.Hex())
	}
	if claims.Role != user.RoleCustomer {
		t.Fatalf("role = %q, want customer", claims.Role)
	}
}

func TestRequireAdminRejectsCustomer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tokens := auth.NewTokenService("test-secret", time.Hour)
	svc := auth.NewService(newMemUserRepo(), tokens, auth.NewLoginRateLimiter(nil), "")

	id := primitive.NewObjectID()
	token, _, err := tokens.Issue(id, user.RoleCustomer)
	if err != nil {
		t.Fatal(err)
	}

	mw := auth.MiddlewareDeps{Service: svc}
	r := gin.New()
	r.GET("/admin/test", auth.RequireAuth(mw), auth.RequireAdmin(mw), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/test", nil)
	req.AddCookie(&http.Cookie{Name: auth.CookieName, Value: token})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", w.Code)
	}
}

func TestRegisterAndLogin(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	repo := newMemUserRepo()
	tokens := auth.NewTokenService("test-secret", time.Hour)
	svc := auth.NewService(repo, tokens, auth.NewLoginRateLimiter(rdb), "admin@example.com")

	ctx := context.Background()
	result, err := svc.Register(ctx, auth.RegisterInput{
		Email:    "user@example.com",
		Password: "password123",
		Name:     "Test User",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.User.Role != user.RoleCustomer {
		t.Fatalf("role = %q, want customer", result.User.Role)
	}

	loggedIn, err := svc.Login(ctx, auth.LoginInput{
		Email:    "user@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatal(err)
	}
	if loggedIn.User.ID != result.User.ID {
		t.Fatal("login returned different user")
	}

	admin, err := svc.Register(ctx, auth.RegisterInput{
		Email:    "admin@example.com",
		Password: "password123",
		Name:     "Admin",
	})
	if err != nil {
		t.Fatal(err)
	}
	if admin.User.Role != user.RoleAdmin {
		t.Fatalf("admin role = %q, want admin", admin.User.Role)
	}
}

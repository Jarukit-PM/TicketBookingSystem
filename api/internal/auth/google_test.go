package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/auth"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/user"
)

func TestLoginWithGoogleAutoLinksExistingEmail(t *testing.T) {
	repo := newMemUserRepo()
	tokens := auth.NewTokenService("test-secret", time.Hour)
	svc := auth.NewService(repo, tokens, auth.NewLoginRateLimiter(nil), "")

	ctx := context.Background()
	existingID := primitive.NewObjectID()
	repo.byEmail["user@example.com"] = &user.User{
		ID:           existingID,
		Email:        "user@example.com",
		PasswordHash: "hashed-password",
		Name:         "Password User",
		Role:         user.RoleCustomer,
	}

	result, err := svc.LoginWithGoogle(ctx, auth.GoogleProfile{
		ID:            "google-sub-123",
		Email:         "user@example.com",
		Name:          "Google Name",
		EmailVerified: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	if result.User.ID != existingID.Hex() {
		t.Fatalf("user id = %q, want %q", result.User.ID, existingID.Hex())
	}

	linked := repo.byEmail["user@example.com"]
	if linked.GoogleID != "google-sub-123" {
		t.Fatalf("google id = %q, want google-sub-123", linked.GoogleID)
	}
	if linked.PasswordHash != "hashed-password" {
		t.Fatal("expected password hash preserved after auto-link")
	}
}

func TestLoginWithGoogleCreatesNewUser(t *testing.T) {
	repo := newMemUserRepo()
	tokens := auth.NewTokenService("test-secret", time.Hour)
	svc := auth.NewService(repo, tokens, auth.NewLoginRateLimiter(nil), "")

	result, err := svc.LoginWithGoogle(context.Background(), auth.GoogleProfile{
		ID:            "google-sub-new",
		Email:         "new@example.com",
		Name:          "New User",
		EmailVerified: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.User.Email != "new@example.com" {
		t.Fatalf("email = %q", result.User.Email)
	}
	if repo.byEmail["new@example.com"].GoogleID != "google-sub-new" {
		t.Fatal("expected google id on new user")
	}
}

func TestLoginWithGoogleRejectsUnverifiedEmail(t *testing.T) {
	svc := auth.NewService(newMemUserRepo(), auth.NewTokenService("test-secret", time.Hour), auth.NewLoginRateLimiter(nil), "")

	_, err := svc.LoginWithGoogle(context.Background(), auth.GoogleProfile{
		ID:            "google-sub-123",
		Email:         "user@example.com",
		Name:          "User",
		EmailVerified: false,
	})
	if !errors.Is(err, auth.ErrGoogleEmailNotVerified) {
		t.Fatalf("err = %v, want ErrGoogleEmailNotVerified", err)
	}
}

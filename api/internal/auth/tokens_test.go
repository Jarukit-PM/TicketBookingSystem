package auth

import (
	"testing"
	"time"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestTokenServiceIssueAndParse(t *testing.T) {
	svc := NewTokenService("test-secret", time.Hour)
	userID := primitive.NewObjectID()

	token, expiresAt, err := svc.Issue(userID, user.RoleCustomer)
	if err != nil {
		t.Fatalf("Issue() error = %v", err)
	}
	if token == "" {
		t.Fatal("Issue() returned empty token")
	}
	if expiresAt.Before(time.Now().UTC()) {
		t.Fatalf("expiresAt = %v, want future time", expiresAt)
	}

	claims, err := svc.Parse(token)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if claims.Subject != userID.Hex() {
		t.Errorf("Subject = %q, want %q", claims.Subject, userID.Hex())
	}
	if claims.Role != user.RoleCustomer {
		t.Errorf("Role = %q, want %q", claims.Role, user.RoleCustomer)
	}
}

func TestTokenServiceParseRejectsInvalidToken(t *testing.T) {
	svc := NewTokenService("test-secret", time.Hour)
	if _, err := svc.Parse("not-a-token"); err == nil {
		t.Fatal("Parse() error = nil, want error for invalid token")
	}
}

func TestTokenServiceParseRejectsWrongSecret(t *testing.T) {
	issuer := NewTokenService("issuer-secret", time.Hour)
	parser := NewTokenService("other-secret", time.Hour)

	token, _, err := issuer.Issue(primitive.NewObjectID(), user.RoleAdmin)
	if err != nil {
		t.Fatalf("Issue() error = %v", err)
	}
	if _, err := parser.Parse(token); err == nil {
		t.Fatal("Parse() error = nil, want error for wrong secret")
	}
}

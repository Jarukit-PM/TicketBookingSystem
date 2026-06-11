package auth

import (
	"testing"
)

func TestHashPasswordAndCheckPassword(t *testing.T) {
	hash, err := HashPassword("secret-password")
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	if hash == "" {
		t.Fatal("HashPassword() returned empty hash")
	}
	if !CheckPassword(hash, "secret-password") {
		t.Error("CheckPassword() = false, want true for matching password")
	}
	if CheckPassword(hash, "wrong-password") {
		t.Error("CheckPassword() = true, want false for wrong password")
	}
}

func TestHashPasswordUsesBcryptCost12(t *testing.T) {
	hash, err := HashPassword("secret-password")
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	if len(hash) < 60 {
		t.Fatalf("hash length = %d, want bcrypt hash length", len(hash))
	}
}

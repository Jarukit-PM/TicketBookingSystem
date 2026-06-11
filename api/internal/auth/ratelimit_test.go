package auth

import (
	"context"
	"testing"
)

func TestLoginRateLimiterBlocksAfterMaxAttempts(t *testing.T) {
	client := startTestRedis(t)
	limiter := NewLoginRateLimiter(client)
	ctx := context.Background()
	email := "user@example.com"

	for i := 0; i < loginRateLimitMax; i++ {
		allowed, err := limiter.Allow(ctx, email)
		if err != nil {
			t.Fatalf("Allow() error = %v", err)
		}
		if !allowed {
			t.Fatalf("Allow() = false on attempt %d, want true", i+1)
		}
		if err := limiter.RecordFailure(ctx, email); err != nil {
			t.Fatalf("RecordFailure() error = %v", err)
		}
	}

	allowed, err := limiter.Allow(ctx, email)
	if err != nil {
		t.Fatalf("Allow() error = %v", err)
	}
	if allowed {
		t.Fatal("Allow() = true, want false after max attempts")
	}
}

func TestLoginRateLimiterResetClearsAttempts(t *testing.T) {
	client := startTestRedis(t)
	limiter := NewLoginRateLimiter(client)
	ctx := context.Background()
	email := "reset@example.com"

	if err := limiter.RecordFailure(ctx, email); err != nil {
		t.Fatalf("RecordFailure() error = %v", err)
	}
	if err := limiter.Reset(ctx, email); err != nil {
		t.Fatalf("Reset() error = %v", err)
	}

	allowed, err := limiter.Allow(ctx, email)
	if err != nil {
		t.Fatalf("Allow() error = %v", err)
	}
	if !allowed {
		t.Fatal("Allow() = false, want true after reset")
	}
}

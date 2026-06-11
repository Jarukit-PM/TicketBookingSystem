package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	loginRateLimitMax    = 5
	loginRateLimitWindow = 15 * time.Minute
)

// LoginRateLimiter tracks failed login attempts in Redis.
type LoginRateLimiter struct {
	redis *redis.Client
}

// NewLoginRateLimiter returns a Redis-backed login rate limiter.
func NewLoginRateLimiter(client *redis.Client) *LoginRateLimiter {
	return &LoginRateLimiter{redis: client}
}

func loginRateLimitKey(email string) string {
	return fmt.Sprintf("auth:login:attempts:%s", email)
}

// Allow reports whether another login attempt is permitted for the email.
func (l *LoginRateLimiter) Allow(ctx context.Context, email string) (bool, error) {
	if l.redis == nil {
		return true, nil
	}

	key := loginRateLimitKey(email)
	count, err := l.redis.Get(ctx, key).Int()
	if err == redis.Nil {
		return true, nil
	}
	if err != nil {
		return false, fmt.Errorf("get login attempts: %w", err)
	}
	return count < loginRateLimitMax, nil
}

// RecordFailure increments the failed login counter for the email.
func (l *LoginRateLimiter) RecordFailure(ctx context.Context, email string) error {
	if l.redis == nil {
		return nil
	}

	key := loginRateLimitKey(email)
	pipe := l.redis.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, loginRateLimitWindow)
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("record login failure: %w", err)
	}
	if incr.Err() != nil {
		return fmt.Errorf("incr login attempts: %w", incr.Err())
	}
	return nil
}

// Reset clears failed login attempts after a successful login.
func (l *LoginRateLimiter) Reset(ctx context.Context, email string) error {
	if l.redis == nil {
		return nil
	}

	if err := l.redis.Del(ctx, loginRateLimitKey(email)).Err(); err != nil {
		return fmt.Errorf("reset login attempts: %w", err)
	}
	return nil
}

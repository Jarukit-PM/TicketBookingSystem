package config

import (
	"testing"
)

func TestLoadEnvOverridesDefaults(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("MONGO_URI", "mongodb://custom:27017/db")
	t.Setenv("REDIS_URL", "redis://custom:6379/1")
	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("APP_URL", "http://example.test")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Port != "9090" {
		t.Errorf("Port = %q, want %q", cfg.Port, "9090")
	}
	if cfg.MongoURI != "mongodb://custom:27017/db" {
		t.Errorf("MongoURI = %q, want custom URI", cfg.MongoURI)
	}
	if cfg.RedisURL != "redis://custom:6379/1" {
		t.Errorf("RedisURL = %q, want custom URL", cfg.RedisURL)
	}
	if cfg.JWTSecret != "test-secret" {
		t.Errorf("JWTSecret = %q, want test-secret", cfg.JWTSecret)
	}
	if cfg.AppURL != "http://example.test" {
		t.Errorf("AppURL = %q, want http://example.test", cfg.AppURL)
	}
}

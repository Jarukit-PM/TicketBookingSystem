package auth

import (
	"context"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/config"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/user"
)

// BootstrapConfiguredAdmin creates or promotes the user from ADMIN_EMAIL and
// ADMIN_SEED_PASSWORD when both are set. No-op when either env value is empty.
func BootstrapConfiguredAdmin(ctx context.Context, cfg config.Config, users user.Repository) error {
	if cfg.AdminEmail == "" || cfg.AdminSeedPassword == "" {
		return nil
	}

	tokens := NewTokenService(cfg.JWTSecret, cfg.JWTExpiryDuration())
	svc := NewService(users, tokens, NewLoginRateLimiter(nil), cfg.AdminEmail)
	return svc.BootstrapAdmin(ctx, cfg.AdminEmail, cfg.AdminSeedPassword)
}

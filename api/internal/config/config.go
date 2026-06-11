package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds application settings loaded from config files and environment.
type Config struct {
	Port              string `mapstructure:"port"`
	MongoURI          string `mapstructure:"mongo_uri"`
	RedisURL          string `mapstructure:"redis_url"`
	JWTSecret         string `mapstructure:"jwt_secret"`
	JWTExpiry         string `mapstructure:"jwt_expiry"`
	AdminEmail        string `mapstructure:"admin_email"`
	AdminSeedPassword string `mapstructure:"admin_seed_password"`
	AppURL            string `mapstructure:"app_url"`
	TicketSecret      string `mapstructure:"ticket_secret"`
	SendGridAPIKey    string `mapstructure:"sendgrid_api_key"`
	EmailFrom         string `mapstructure:"email_from"`
	GinMode           string `mapstructure:"gin_mode"`
}

// Load reads configuration from config.yaml (optional) and environment variables.
func Load() (Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/tbs")

	v.SetDefault("port", "8080")
	v.SetDefault("mongo_uri", "mongodb://mongo:27017/tbs")
	v.SetDefault("redis_url", "redis://redis:6379/0")
	v.SetDefault("jwt_secret", "dev-only-change-me")
	v.SetDefault("jwt_expiry", "168h")
	v.SetDefault("app_url", "http://localhost")
	v.SetDefault("gin_mode", "debug")

	_ = v.BindEnv("port", "PORT")
	_ = v.BindEnv("mongo_uri", "MONGO_URI")
	_ = v.BindEnv("redis_url", "REDIS_URL")
	_ = v.BindEnv("jwt_secret", "JWT_SECRET")
	_ = v.BindEnv("jwt_expiry", "JWT_EXPIRY")
	_ = v.BindEnv("admin_email", "ADMIN_EMAIL")
	_ = v.BindEnv("admin_seed_password", "ADMIN_SEED_PASSWORD")
	_ = v.BindEnv("app_url", "APP_URL")
	_ = v.BindEnv("ticket_secret", "TICKET_SECRET")
	_ = v.BindEnv("sendgrid_api_key", "SENDGRID_API_KEY")
	_ = v.BindEnv("email_from", "EMAIL_FROM")
	_ = v.BindEnv("gin_mode", "GIN_MODE")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok && !strings.Contains(err.Error(), "Not Found") {
			return Config{}, fmt.Errorf("read config: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("unmarshal config: %w", err)
	}

	return cfg, nil
}

// MustLoad loads configuration or panics.
func MustLoad() Config {
	cfg, err := Load()
	if err != nil {
		panic(err)
	}
	return cfg
}

// JWTExpiryDuration parses JWT_EXPIRY (default 168h / 7 days).
func (c Config) JWTExpiryDuration() time.Duration {
	if c.JWTExpiry == "" {
		return 168 * time.Hour
	}
	d, err := time.ParseDuration(c.JWTExpiry)
	if err != nil {
		return 168 * time.Hour
	}
	return d
}

// CookieSecure reports whether session cookies should set the Secure flag.
func (c Config) CookieSecure() bool {
	if c.GinMode == "release" {
		return true
	}
	return strings.EqualFold(os.Getenv("COOKIE_SECURE"), "true")
}

// TicketHMACSecret returns the secret used to sign ticket tokens.
func (c Config) TicketHMACSecret() string {
	if c.TicketSecret != "" {
		return c.TicketSecret
	}
	return c.JWTSecret
}

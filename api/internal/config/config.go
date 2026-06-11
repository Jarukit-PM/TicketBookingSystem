package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds application settings loaded from config files and environment.
type Config struct {
	Port      string `mapstructure:"port"`
	MongoURI  string `mapstructure:"mongo_uri"`
	RedisURL  string `mapstructure:"redis_url"`
	JWTSecret string `mapstructure:"jwt_secret"`
	AppURL    string `mapstructure:"app_url"`
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
	v.SetDefault("app_url", "http://localhost")

	_ = v.BindEnv("port", "PORT")
	_ = v.BindEnv("mongo_uri", "MONGO_URI")
	_ = v.BindEnv("redis_url", "REDIS_URL")
	_ = v.BindEnv("jwt_secret", "JWT_SECRET")
	_ = v.BindEnv("app_url", "APP_URL")

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

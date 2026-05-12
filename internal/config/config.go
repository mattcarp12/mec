package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

var _config *Config

type Config struct {
	HTTPPort string

	JWTIssuer         string
	JWTAudience       string
	JWTPrivateKeyPath string

	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration

	CookieSecure bool
}

func Init() {
	cfg := &Config{
		HTTPPort: getEnv("HTTP_PORT", "8080"),

		JWTIssuer:         getEnv("JWT_ISSUER", "mec-auth-service"),
		JWTAudience:       getEnv("JWT_AUDIENCE", "mec-api"),
		JWTPrivateKeyPath: getEnv("JWT_PRIVATE_KEY_PATH", "private.pem"),

		AccessTokenTTL:  getEnvDuration("ACCESS_TOKEN_TTL", 15*time.Minute),
		RefreshTokenTTL: getEnvDuration("REFRESH_TOKEN_TTL", 7*24*time.Hour),

		CookieSecure: getEnvBool("COOKIE_SECURE", false),
	}

	if err := validate(cfg); err != nil {
		log.Fatalf("Unable to validate config: %v", err)
	}

	_config = cfg
}

func Get() *Config {
	if _config == nil {
		log.Fatal("Config not loaded. Call config.Load() before accessing config.Get()")
	}
	return _config
}

func validate(cfg *Config) error {
	if cfg.JWTIssuer == "" {
		return fmt.Errorf("JWT_ISSUER cannot be empty")
	}

	if cfg.JWTAudience == "" {
		return fmt.Errorf("JWT_AUDIENCE cannot be empty")
	}

	if cfg.AccessTokenTTL <= 0 {
		return fmt.Errorf("ACCESS_TOKEN_TTL must be > 0")
	}

	if cfg.RefreshTokenTTL <= 0 {
		return fmt.Errorf("REFRESH_TOKEN_TTL must be > 0")
	}

	return nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	d, err := time.ParseDuration(value)
	if err != nil {
		panic(fmt.Sprintf("invalid duration for %s: %v", key, err))
	}

	return d
}

func getEnvBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	b, err := strconv.ParseBool(value)
	if err != nil {
		panic(fmt.Sprintf("invalid bool for %s: %v", key, err))
	}

	return b
}

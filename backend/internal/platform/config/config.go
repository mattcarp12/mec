package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	Environment string `env:"APP_ENV" envDefault:"development"`
	Port        int    `env:"PORT" envDefault:"8080"`
	DatabaseURL string `env:"DATABASE_URL,required"`
	// We will add Zitadel / Stripe keys here later
}

// Load reads the .env file and parses it into AppConfig
func Load() (*AppConfig, error) {
	// Attempt to load .env file, but ignore errors if it doesn't exist
	// (e.g., in production where env vars are injected directly)
	_ = godotenv.Load()

	var cfg AppConfig
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	log.Printf("Configuration loaded for environment: %s", cfg.Environment)
	return &cfg, nil
}
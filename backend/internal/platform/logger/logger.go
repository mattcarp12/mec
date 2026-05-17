package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// New creates a configured zerolog.Logger
func New(env string) zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339

	if env == "development" {
		// Pretty-print to console in development
		return zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).
			With().
			Timestamp().
			Caller().
			Logger()
	}

	// Fast JSON logging in production
	return zerolog.New(os.Stdout).
		With().
		Timestamp().
		Logger()
}
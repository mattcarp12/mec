package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // Keep the driver registry contained here
)

// InitDB opens a connection pool, pings it to ensure health, 
// and configures optimal pool settings.
func InitDB(databaseURL string) (*sql.DB, *Queries, error) {
	dbConn, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set safe, production-ready connection pool limits
	dbConn.SetMaxOpenConns(25)
	dbConn.SetMaxIdleConns(25)
	dbConn.SetConnMaxLifetime(5 * time.Minute)

	if err := dbConn.Ping(); err != nil {
		dbConn.Close() // Clean up before exiting on failure
		return nil, nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Initialize your compiled sqlc queries wrapper
	queries := New(dbConn)

	return dbConn, queries, nil
}
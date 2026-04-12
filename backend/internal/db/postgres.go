package db

import (
	"context"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // SOTA: pgx is the fastest Postgres driver for Go
	"github.com/jmoiron/sqlx"
)

// NewPostgresDB establishes a connection pool to Postgres.
func NewPostgresDB(ctx context.Context, dsn string) (*sqlx.DB, error) {
	// ConnectContext automatically pings the database to ensure the connection is alive
	db, err := sqlx.ConnectContext(ctx, "pgx", dsn)
	if err != nil {
		return nil, err
	}

	// SOTA Connection Pooling Setup
	// Without this, Go will open infinite connections under load and crash your database.
	db.SetMaxOpenConns(25)                 // Max simultaneous connections
	db.SetMaxIdleConns(25)                 // Max idle connections to keep alive
	db.SetConnMaxLifetime(5 * time.Minute) // Retire connections after 5 mins to prevent stale sockets

	return db, nil
}

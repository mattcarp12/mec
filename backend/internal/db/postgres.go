package db

import (
	"context"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // SOTA: pgx is the fastest Postgres driver for Go
	"github.com/jmoiron/sqlx"
)

var dsn = os.Getenv("DATABASE_URL")
var DB *sqlx.DB

func Init() {
	db, err := newPostgresDB()
	if err != nil {
		log.Fatalf("Fatal: Could not connect to database: %v", err)
	}
	log.Println("Successfully connected to Postgres!")
	DB = db
}

// newPostgresDB establishes a connection pool to Postgres.
func newPostgresDB() (*sqlx.DB, error) {
	if dsn == "" {
		dsn = "postgres://admin:supersecretpassword@localhost:5432/identity?sslmode=disable"
	}

	// ConnectContext automatically pings the database to ensure the connection is alive
	db, err := sqlx.ConnectContext(context.TODO(), "pgx", dsn)
	if err != nil {
		return nil, err
	}

	// Connection Pooling Setup
	// Without this, Go will open infinite connections under load and crash your database.
	db.SetMaxOpenConns(25)                 // Max simultaneous connections
	db.SetMaxIdleConns(25)                 // Max idle connections to keep alive
	db.SetConnMaxLifetime(5 * time.Minute) // Retire connections after 5 mins to prevent stale sockets

	return db, nil
}

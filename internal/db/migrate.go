package db

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/identity/*.sql
var _identityMigrations embed.FS
var identityMigrations, _ = fs.Sub(_identityMigrations, "migrations/identity")

//go:embed migrations/catalog/*.sql
var _catalogMigtations embed.FS
var catalogMigrations, _ = fs.Sub(_catalogMigtations, "migrations/catalog")

// runMigrations applies all pending UP migrations to the database.
func runMigrations(migrationsFS fs.FS) error {
	log.Println("Initializing database migrations...")

	// 1. Load the embedded file system
	d, err := iofs.New(migrationsFS, ".")
	if err != nil {
		return fmt.Errorf("failed to load embedded migrations: %w", err)
	}

	// 2. Initialize the migrator
	m, err := migrate.NewWithSourceInstance("iofs", d, dsn)
	if err != nil {
		return fmt.Errorf("failed to initialize migrator: %w", err)
	}

	// 3. Run the migrations
	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("Database is already up to date. No migrations applied.")
			return nil
		}
		return fmt.Errorf("failed to run migrate up: %w", err)
	}

	log.Println("Migrations applied successfully!")
	return nil
}

func RunMigrations(domain string) error {
	switch domain {
	case "identity":
		return runMigrations(identityMigrations)
	case "catalog":
		return runMigrations(catalogMigrations)
	default:
		return fmt.Errorf("migration for %s not supported", domain)
	}
}

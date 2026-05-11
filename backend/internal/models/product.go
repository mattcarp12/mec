package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID   uuid.UUID `db:"id" json:"id"`
	Name string    `db:"name" json:"name"`

	Description string `json:"description" db:"description"`

	// 2. Money Handling (Always use Integers/Cents to prevent rounding errors)
	// Example: $10.99 becomes 1099
	PriceCents     int64 `json:"price_cents" db:"price_cents"`
	InventoryCount int   `json:"inventory_count" db:"inventory_count"`

	// 3. The JSONB "Flex" Data
	// This maps directly to the Postgres JSONB column.
	// pgx will automatically marshal this into JSON for you.
	Attributes map[string]any `json:"attributes" db:"attributes"`

	// 4. Metadata
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

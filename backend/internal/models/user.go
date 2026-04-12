package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a verified identity in our ecosystem.
type User struct {
	ID    uuid.UUID `db:"id" json:"id"`
	Email string    `db:"email" json:"email"`

	// SOTA Security: The json:"-" tag ensures the password hash is
	// NEVER accidentally leaked in a JSON API response.
	PasswordHash string `db:"password_hash" json:"-"`

	Role      string    `db:"role" json:"role"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

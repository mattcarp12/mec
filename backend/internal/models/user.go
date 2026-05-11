package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a verified identity in our ecosystem.
type User struct {
	ID    uuid.UUID `db:"id" json:"id"`
	Email string    `db:"email" json:"email"`

	// Security: The json:"-" tag ensures the password hash is
	// NEVER accidentally leaked in a JSON API response.
	PasswordHash string `db:"password_hash" json:"-"`

	Role      string    `db:"role" json:"role"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type RefreshToken struct {
	TokenHash string    `db:"token_hash" json:"-"`
	UserID    uuid.UUID `db:"user_id" json:"-"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	Revoked   bool      `db:"revoked" json:"revoked"`
}

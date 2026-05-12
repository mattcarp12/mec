package main

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/mattcarp12/mec/internal/db"
	"github.com/mattcarp12/mec/internal/models"
)

func CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (email, password_hash, role) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at, updated_at`

	// QueryRowContext is crucial for context cancellation (e.g., if the user
	// closes their browser mid-request, we cancel the DB query to save resources).
	err := db.DB.QueryRowContext(ctx, query, user.Email, user.PasswordHash, user.Role).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Printf("error creating user: %v\n", err)
	}

	return err
}

func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, email, password_hash, role, created_at, updated_at FROM users WHERE email = $1`

	err := db.DB.GetContext(ctx, &user, query, email)
	if err != nil {
		log.Printf("error fetching user: %v\n", err)
		return nil, err // Will return sql.ErrNoRows if not found
	}
	return &user, nil
}

func StoreRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	query := `
	INSERT INTO refresh_tokens (token_hash, user_id, expires_at)
	VALUES ($1, $2, $3)`

	_, err := db.DB.ExecContext(ctx, query, tokenHash, userID, expiresAt)
	if err != nil {
		log.Printf("error inserting refresh token: %w", err)
	}
	return err
}

func GetRefreshToken(ctx context.Context, tokenHash string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	query := `SELECT token_hash, user_id, expires_at, revoked FROM refresh_tokens WHERE token_hash = $1`

	err := db.DB.GetContext(ctx, &rt, query, tokenHash)
	if err != nil {
		log.Printf("error fetching refresh token: %w", err)
		return nil, err
	}
	return &rt, nil
}

// RevokeRefreshToken marks a specific token as revoked in the database.
func RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	query := `UPDATE refresh_tokens SET revoked = true WHERE token_hash = $1`

	// ExecContext is used here because we don't need to return any rows[cite: 17].
	_, err := db.DB.ExecContext(ctx, query, tokenHash)
	if err != nil {
		log.Printf("error revoking refresh token: %v\n", err)
	}
	return err
}

package repository

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/mattcarp12/mec/internal/models"
)

// UserRepository defines the strict contract for data access.
type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

// PostgresUserRepository implements UserRepository for PostgreSQL.
type PostgresUserRepository struct {
	db *sqlx.DB
}

func NewPostgresUserRepository(db *sqlx.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (email, password_hash, role) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at, updated_at`

	// QueryRowContext is crucial for context cancellation (e.g., if the user
	// closes their browser mid-request, we cancel the DB query to save resources).
	err := r.db.QueryRowContext(ctx, query, user.Email, user.PasswordHash, user.Role).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Printf("error creating user: %v\n", err)
	}

	return err
}

func (r *PostgresUserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, email, password_hash, role, created_at, updated_at FROM users WHERE email = $1`

	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		log.Printf("error fetching user: %v\n", err)
		return nil, err // Will return sql.ErrNoRows if not found
	}
	return &user, nil
}

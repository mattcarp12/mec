package auth

import (
	"context"
	"errors"

	"github.com/mattcarp12/mec/internal/models"
	"github.com/mattcarp12/mec/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserExists         = errors.New("user with this email already exists")
)

// AuthService handles the business logic for user identities.
type AuthService struct {
	repo repository.UserRepository
}

// NewAuthService creates a new instance of the AuthService.
func NewAuthService(repo repository.UserRepository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

// Register creates a new user after hashing their password.
func (s *AuthService) Register(ctx context.Context, email, password string) (*models.User, error) {
	// 1. Check if user already exists (optional, DB constraint will also catch this,
	// but doing it here allows for a cleaner error message).
	existingUser, _ := s.repo.GetUserByEmail(ctx, email)
	if existingUser != nil {
		return nil, ErrUserExists
	}

	// 2. Hash the password
	hash, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	// 3. Create the user model
	user := &models.User{
		Email:        email,
		PasswordHash: hash,
		Role:         "customer", // Default role
	}

	// 4. Persist to database
	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Login verifies credentials and returns the user if successful.
// Note: We will add PASETO token generation to this in the next step.
func (s *AuthService) Login(ctx context.Context, email, password string) (*models.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		// We don't expose if the user exists or not to prevent enumeration attacks
		return nil, ErrInvalidCredentials
	}

	match, err := VerifyPassword(password, user.PasswordHash)
	if err != nil || !match {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

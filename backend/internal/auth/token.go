package auth

import (
	"time"
	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
)

// GeneratePASETO creates a secure, encrypted V4 PASETO token.
func GeneratePASETO(userID uuid.UUID, secretKey string) (string, error) {
	token := paseto.NewToken()

	// SOTA Token Claims
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(24 * time.Hour))
	token.SetString("user_id", userID.String())

	// PASETO requires exactly a 32-byte key for V4 symmetric encryption
	key, err := paseto.V4SymmetricKeyFromBytes([]byte(secretKey))
	if err != nil {
		return "", err
	}

	// Encrypt the token (this ensures both authenticity and confidentiality)
	encryptedToken := token.V4Encrypt(key, nil)
	return encryptedToken, nil
}

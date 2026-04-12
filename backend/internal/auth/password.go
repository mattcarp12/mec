package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Parameters for Argon2id.
// These should ideally be configurable via environment variables in production.
const (
	iterations = 3         // Number of iterations
	memory     = 16 * 1024 // 16MB
	threads    = 2         // Number of parallel threads
	keyLen     = 32        // Length of the generated hash
	saltLen    = 16        // Length of the random salt
)

// HashPassword generates an Argon2id hash from a plaintext password.
func HashPassword(password string) (string, error) {
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, iterations, memory, threads, keyLen)

	// We encode the parameters into the final string so we can verify it later
	// Format: $argon2id$v=19$m=65536,t=3,p=2$<salt>$<hash>
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, memory, iterations, threads, b64Salt, b64Hash)

	return encodedHash, nil
}

// VerifyPassword compares a plaintext password against an encoded Argon2id hash.
func VerifyPassword(password, encodedHash string) (bool, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, errors.New("invalid hash format")
	}

	var version int
	var mem, t uint32
	var p uint8
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return false, err
	}
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &mem, &t, &p)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	hash := argon2.IDKey([]byte(password), salt, t, mem, p, uint32(len(decodedHash)))

	// SOTA Security: Use subtle.ConstantTimeCompare to prevent timing attacks
	if subtle.ConstantTimeCompare(decodedHash, hash) == 1 {
		return true, nil
	}
	return false, nil
}

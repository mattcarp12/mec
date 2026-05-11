package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

var privateKey *rsa.PrivateKey

func init() {
	loadPrivateKey()
}

func loadPrivateKey() {
	path := os.Getenv("JWT_PRIVATE_KEY_PATH")
	if path == "" {
		log.Println("JWT_PRIVATE_KEY_PATH not set, using default path 'private.pem'")
		path = "private.pem" // Default path for development
	}

	// read file at path
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read private key: %v", err)
	}

	// parse PEM block
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(data)
	if err != nil {
		log.Fatalf("failed to parse private key: %v", err)
	}
}

// Parameters for Argon2id.
// These should ideally be configurable via environment variables in production.
const (
	iterations = 3         // Number of iterations
	memory     = 16 * 1024 // 16MB
	threads    = 2         // Number of parallel threads
	keyLen     = 32        // Length of the generated hash
	saltLen    = 16        // Length of the random salt
)

// hashPassword generates an Argon2id hashPassword from a plaintext password.
func hashPassword(pass string) (string, error) {
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(pass), salt, iterations, memory, threads, keyLen)

	// We encode the parameters into the final string so we can verify it later
	// Format: $argon2id$v=19$m=65536,t=3,p=2$<salt>$<hash>
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, memory, iterations, threads, b64Salt, b64Hash)

	return encodedHash, nil
}

// verifyPassword compares a plaintext password against an encoded Argon2id hash.
func verifyPassword(password, encodedHash string) (bool, error) {
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

	// Security: Use subtle.ConstantTimeCompare to prevent timing attacks
	if subtle.ConstantTimeCompare(decodedHash, hash) == 1 {
		return true, nil
	}
	return false, nil
}

func generateAccessToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"iss": "mec-auth-service",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func generateRefreshToken() string {
	return rand.Text()
}

// hashToken generates a SHA-256 hash of a plaintext token string.
// We use hex encoding to make it database-friendly (VARCHAR).
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

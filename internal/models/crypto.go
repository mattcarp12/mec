package models

// JWK represents a single JSON Web Key
type JWK struct {
	Kty string `json:"kty"` // Key Type (e.g., "RSA")
	Alg string `json:"alg"` // Algorithm (e.g., "RS256")
	Use string `json:"use"` // Intended Use (e.g., "sig" for signature)
	Kid string `json:"kid"` // Key ID
	N   string `json:"n"`   // Modulus
	E   string `json:"e"`   // Exponent
}

// JWKS represents a set of JSON Web Keys
type JWKS struct {
	Keys []JWK `json:"keys"`
}

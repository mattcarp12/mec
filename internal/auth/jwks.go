package auth

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mattcarp12/mec/internal/models"
)

type JWKSClient struct {
	jwksURL string

	mu   sync.RWMutex
	keys map[string]*rsa.PublicKey

	httpClient *http.Client
}

func NewJWKSClient(jwksURL string) *JWKSClient {
	return &JWKSClient{
		jwksURL: jwksURL,
		keys:    make(map[string]*rsa.PublicKey),
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// Refresh fetches the JWKS from the configured URL and updates the client's key set.
func (c *JWKSClient) Refresh(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.jwksURL, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	var jwks models.JWKS

	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return err
	}

	newKeys := make(map[string]*rsa.PublicKey)

	for _, jwk := range jwks.Keys {
		pubKey, err := jwkToPublicKey(jwk)
		if err != nil {
			continue
		}

		newKeys[jwk.Kid] = pubKey
	}

	c.mu.Lock()
	c.keys = newKeys
	c.mu.Unlock()

	return nil
}

func (c *JWKSClient) GetKey(kid string) (*rsa.PublicKey, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	key, ok := c.keys[kid]
	return key, ok
}

func jwkToPublicKey(jwk models.JWK) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, err
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, err
	}

	n := new(big.Int).SetBytes(nBytes)

	e := 0
	for _, b := range eBytes {
		e = e<<8 + int(b)
	}

	return &rsa.PublicKey{
		N: n,
		E: e,
	}, nil
}

type Claims struct {
	Subject string   `json:"sub"`
	Roles   []string `json:"roles"`
	Scope   string   `json:"scope"`

	jwt.RegisteredClaims
}

func (c *JWKSClient) VerifyToken(
	tokenString string,
	issuer string,
	audience string,
) (*Claims, error) {

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			kidRaw, ok := token.Header["kid"]
			if !ok {
				return nil, errors.New("missing kid")
			}

			kid, ok := kidRaw.(string)
			if !ok {
				return nil, errors.New("invalid kid")
			}

			key, ok := c.GetKey(kid)
			if !ok {
				return nil, errors.New("unknown kid")
			}

			return key, nil
		},
		jwt.WithIssuer(issuer),
		jwt.WithAudience(audience),
		jwt.WithValidMethods([]string{"RS256"}),
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
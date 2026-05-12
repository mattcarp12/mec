package main

import (
	"context"
	"crypto/rsa"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mattcarp12/mec/internal/config"
)

type contextKey string

const userContextKey contextKey = "user"

type AuthenticatedUser struct {
	Subject string
	Roles   []string
	Scopes  []string
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			writeJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "missing authorization header",
			})
			return
		}

		tokenString, err := extractBearerToken(authHeader)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "invalid authorization header",
			})
			return
		}

		token, err := jwt.Parse(tokenString, keyFunc,
			jwt.WithIssuer(config.Get().JWTIssuer),
			jwt.WithAudience(config.Get().JWTAudience),
			jwt.WithValidMethods([]string{"RS256"}),
		)

		if err != nil || !token.Valid {
			writeJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "invalid token",
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			writeJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "invalid claims",
			})
			return
		}

		sub, _ := claims["sub"].(string)

		user := &AuthenticatedUser{
			Subject: sub,
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractBearerToken(header string) (string, error) {
	parts := strings.Split(header, " ")

	if len(parts) != 2 {
		return "", errors.New("invalid authorization header")
	}

	if parts[0] != "Bearer" {
		return "", errors.New("invalid authorization type")
	}

	return parts[1], nil
}

func keyFunc(token *jwt.Token) (any, error) {
	kidRaw, ok := token.Header["kid"]
	if !ok {
		return nil, errors.New("missing kid")
	}

	kid, ok := kidRaw.(string)
	if !ok {
		return nil, errors.New("invalid kid")
	}

	key, ok := keystore.Load(kid)
	if !ok {
		return nil, errors.New("unknown kid")
	}

	privateKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("invalid key type")
	}

	return &privateKey.PublicKey, nil
}

func GetAuthenticatedUser(ctx context.Context) (*AuthenticatedUser, bool) {
	user, ok := ctx.Value(userContextKey).(*AuthenticatedUser)
	return user, ok
}

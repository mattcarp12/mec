package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type PBUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type PBAuthResponse struct {
	Record PBUser `json:"record"`
}

type contextKey string
const UserIDKey contextKey = "user_id"

// PocketBaseAuth validates the incoming JWT against the PocketBase instance
func PocketBaseAuth(pbURL string) func(http.Handler) http.Handler {
	// Short timeout so auth checks don't block indefinitely
	client := &http.Client{Timeout: 2 * time.Second}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")

			// Token introspection call to PocketBase
			req, err := http.NewRequest("POST", pbURL+"/api/collections/users/auth-refresh", nil)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			req.Header.Set("Authorization", "Bearer "+token)

			resp, err := client.Do(req)
			if err != nil || resp.StatusCode != http.StatusOK {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			defer resp.Body.Close()

			var pbResp PBAuthResponse
			if err := json.NewDecoder(resp.Body).Decode(&pbResp); err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Embed the validated user identity into the context
			ctx := context.WithValue(r.Context(), UserIDKey, pbResp.Record.ID)
			
			// Pass the request down the chain
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
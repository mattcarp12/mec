package main

import (
	"encoding/json"
	"net/http"
	"time"
)

// writeJSON is a helper to centralize JSON responses.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Register handles POST /api/v1/auth/register
func Register_Handler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// SOTA: Use Decode instead of Unmarshal for memory efficiency on streams
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request payload"})
		return
	}

	user, err := RegisterNewUser(r.Context(), payload.Email, payload.Password)
	if err != nil {
		if err == ErrUserExists {
			writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

// Login handles POST /api/v1/auth/login
func Login_Handler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request payload"})
		return
	}

	user, err := Login(r.Context(), payload.Email, payload.Password)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}

	// Generate PASETO Token
	token, err := generatePASETO(user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not generate token"})
		return
	}

	// Security: Set the token in an HttpOnly cookie to prevent XSS theft.
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,                 // JS cannot read this
		Secure:   false,                // Set to TRUE in production (requires HTTPS)
		SameSite: http.SameSiteLaxMode, // Protects against CSRF
	})

	writeJSON(w, http.StatusOK, map[string]string{"message": "login successful"})
}

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

	// Generate JWT
	accessToken, err := generateAccessToken(user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not generate token"})
		return
	}

	// Generate Refresh Token
	refreshToken := generateRefreshToken()
	err = StoreRefreshToken(r.Context(), user.ID, hashToken(refreshToken), time.Now().Add(7*24*time.Hour))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error storing token"})
		return
	}

	// Security: Set the token in an HttpOnly cookie to prevent XSS theft.
	// Set the cookie (Refresh Token)
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken, // Store the RAW token in the browser, not the hash!
		Path:     "/",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   false, // Set to true in prod
		SameSite: http.SameSiteLaxMode,
	})

	// Return the Access Token as JSON
	writeJSON(w, http.StatusOK, map[string]any{
		"access_token": accessToken,
		"expires_in":   900,
	})
}

// Refresh_Handler handles POST /api/v1/auth/refresh
func Refresh_Handler(w http.ResponseWriter, r *http.Request) {
	// 1. Extract the Refresh Token from the HttpOnly cookie
	cookie, err := r.Cookie("refresh_token") // Make sure this matches the cookie name set in Login!
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "missing refresh token"})
		return
	}
	rawToken := cookie.Value

	// 2. Hash the raw token
	tokenHash := hashToken(rawToken)

	// 3. Fetch the token record from the database
	rt, err := GetRefreshToken(r.Context(), tokenHash)
	if err != nil {
		// If it's sql.ErrNoRows or any other error, deny access.
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid refresh token"})
		return
	}

	// 4. Validate the token state
	if rt.Revoked {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "token has been revoked"})
		return
	}
	if time.Now().After(rt.ExpiresAt) {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "refresh token has expired, please log in again"})
		return
	}

	// 5. Generate a NEW 15-minute Access Token using the UserID from the database
	newAccessToken, err := generateAccessToken(rt.UserID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to generate access token"})
		return
	}

	// 6. Return the new Access Token in the JSON body
	writeJSON(w, http.StatusOK, map[string]any{
		"access_token": newAccessToken,
		"expires_in":   900, // 15 minutes in seconds
	})
}

// Logout_Handler handles POST /api/v1/auth/logout
func Logout_Handler(w http.ResponseWriter, r *http.Request) {
	// 1. Extract the Refresh Token from the HttpOnly cookie
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		// If there's no cookie, from our perspective, they are already logged out.
		// We return 200 OK so the frontend can clear its local state without erroring.
		writeJSON(w, http.StatusOK, map[string]string{"message": "logged out"})
		return
	}

	// 2. Hash the raw token to match the database column
	tokenHash := hashToken(cookie.Value)

	// 3. Mark as revoked in the database
	err = RevokeRefreshToken(r.Context(), tokenHash)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error during logout"})
		return
	}

	// 4. Instruct the browser to delete the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0), // Set expiration to the past
		MaxAge:   -1,              // Tell the browser to delete it immediately
		HttpOnly: true,            // Maintain the same security flags [cite: 21]
		Secure:   false,           // Set to true in prod (requires HTTPS) [cite: 21]
		SameSite: http.SameSiteLaxMode, // Protects against CSRF [cite: 21]
	})

	writeJSON(w, http.StatusOK, map[string]string{"message": "successfully logged out"})
}
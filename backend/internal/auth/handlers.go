package auth

import (
	"encoding/json"
	"net/http"
	"time"
)

type Handler struct {
	service     *AuthService
	tokenSecret string
}

func NewHandler(service *AuthService, tokenSecret string) *Handler {
	return &Handler{
		service:     service,
		tokenSecret: tokenSecret,
	}
}

// writeJSON is a helper to centralize JSON responses.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Register handles POST /api/v1/auth/register
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// SOTA: Use Decode instead of Unmarshal for memory efficiency on streams
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request payload"})
		return
	}

	user, err := h.service.Register(r.Context(), payload.Email, payload.Password)
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
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request payload"})
		return
	}

	user, err := h.service.Login(r.Context(), payload.Email, payload.Password)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}

	// Generate PASETO Token
	token, err := GeneratePASETO(user.ID, h.tokenSecret)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not generate token"})
		return
	}

	// SOTA Security: Set the token in an HttpOnly cookie to prevent XSS theft.
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

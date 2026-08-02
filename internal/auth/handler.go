package auth

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/qper/hflow/internal/domain"
)

// RegisterRequest is the initial registration payload.
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest is the initial login payload.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RefreshRequest is the basic refresh payload.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// LogoutRequest allows the caller to revoke a specific refresh token.
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Handler provides minimal auth endpoints for manual verification.
type Handler struct {
	db           *sql.DB
	passwordHash PasswordHasher
	tokenManager TokenManager
}

func NewHandler(db *sql.DB, accessSecret, refreshSecret string) Handler {
	return Handler{
		db:           db,
		passwordHash: NewPasswordHasher(),
		tokenManager: NewTokenManager(accessSecret, refreshSecret),
	}
}

func (h Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "username, email, and password are required"})
		return
	}
	if err := ValidatePassword(req.Password); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	passwordHash, err := h.passwordHash.HashPassword(req.Password)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	user := domain.User{Username: req.Username, Email: req.Email, Theme: "system", Timezone: "UTC"}
	_, err = h.db.Exec(`
		INSERT INTO users (username, email, display_name, timezone, theme, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
	`, user.Username, user.Email, user.DisplayName, user.Timezone, user.Theme, passwordHash)
	if err != nil {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "user already exists"})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"message": "registered"})
}

func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	var storedHash string
	var userID string
	err := h.db.QueryRow(`SELECT id, password_hash FROM users WHERE username = $1 AND deleted_at IS NULL`, req.Username).Scan(&userID, &storedHash)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}
	if !h.passwordHash.VerifyPassword(storedHash, req.Password) {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}
	pair, err := h.tokenManager.IssueTokenPair(userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if err := h.storeRefreshToken(userID, pair.RefreshToken, pair.RefreshExpiresAt); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to persist refresh token"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"access_token":  pair.AccessToken,
		"refresh_token": pair.RefreshToken,
		"token_type":    "bearer",
		"expires_at":    pair.ExpiresAt.UTC().Format(time.RFC3339),
	})
}

func (h Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	claims, err := h.tokenManager.ParseToken(req.RefreshToken, "refresh")
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid refresh token"})
		return
	}
	if valid, err := h.hasRefreshToken(claims.UserID, req.RefreshToken); err != nil || !valid {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "refresh token has been revoked or expired"})
		return
	}
	if err := h.revokeRefreshToken(claims.UserID, req.RefreshToken); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to revoke refresh token"})
		return
	}
	pair, err := h.tokenManager.IssueTokenPair(claims.UserID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if err := h.storeRefreshToken(claims.UserID, pair.RefreshToken, pair.RefreshExpiresAt); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to persist refresh token"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"access_token":  pair.AccessToken,
		"refresh_token": pair.RefreshToken,
		"token_type":    "bearer",
		"expires_at":    pair.ExpiresAt.UTC().Format(time.RFC3339),
	})
}

func (h Handler) Logout(w http.ResponseWriter, r *http.Request) {
	var req LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err.Error() != "EOF" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.RefreshToken == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "refresh token is required"})
		return
	}
	claims, err := h.tokenManager.ParseToken(req.RefreshToken, "refresh")
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid refresh token"})
		return
	}
	if err := h.revokeRefreshToken(claims.UserID, req.RefreshToken); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to revoke refresh token"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}

func (h Handler) storeRefreshToken(userID, rawToken string, expiresAt time.Time) error {
	tokenHash := hashToken(rawToken)
	_, err := h.db.Exec(`
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
	`, userID, tokenHash, expiresAt)
	return err
}

func (h Handler) hasRefreshToken(userID, rawToken string) (bool, error) {
	tokenHash := hashToken(rawToken)
	var valid bool
	err := h.db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM refresh_tokens
			WHERE user_id = $1 AND token_hash = $2 AND revoked_at IS NULL AND expires_at > NOW()
		)
	`, userID, tokenHash).Scan(&valid)
	return valid, err
}

func (h Handler) revokeRefreshToken(userID, rawToken string) error {
	tokenHash := hashToken(rawToken)
	_, err := h.db.Exec(`
		UPDATE refresh_tokens
		SET revoked_at = NOW(), updated_at = NOW()
		WHERE user_id = $1 AND token_hash = $2 AND revoked_at IS NULL
	`, userID, tokenHash)
	return err
}

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

package auth

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// PasswordHasher provides password hashing and verification using bcrypt.
type PasswordHasher struct{}

func NewPasswordHasher() PasswordHasher {
	return PasswordHasher{}
}

func (PasswordHasher) HashPassword(password string) (string, error) {
	if err := ValidatePassword(password); err != nil {
		return "", err
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return string(bytes), nil
}

func (PasswordHasher) VerifyPassword(hash, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}
	return true
}

// ValidatePassword enforces a reasonable default policy for the initial release.
func ValidatePassword(password string) error {
	if len(password) < 12 {
		return errors.New("password must be at least 12 characters long")
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case strings.ContainsRune(`!@#$%^&*()-_=+[]{}<>/?`, r):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return errors.New("password must contain uppercase, lowercase, number, and special character")
	}
	return nil
}

// TokenPair contains access and refresh tokens.
type TokenPair struct {
	AccessToken      string
	RefreshToken     string
	ExpiresAt        time.Time
	RefreshExpiresAt time.Time
}

// TokenManager signs and validates access and refresh tokens.
type TokenManager struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewTokenManager(accessSecret, refreshSecret string) TokenManager {
	return TokenManager{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		accessTTL:     15 * time.Minute,
		refreshTTL:    30 * 24 * time.Hour,
	}
}

func (m TokenManager) AccessTTL() time.Duration  { return m.accessTTL }
func (m TokenManager) RefreshTTL() time.Duration { return m.refreshTTL }

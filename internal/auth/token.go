package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Claims is a minimal token payload for the initial auth implementation.
type Claims struct {
	UserID string `json:"user_id"`
	Type   string `json:"type"`
	Exp    int64  `json:"exp"`
}

func (m TokenManager) issueToken(userID, tokenType string, ttl time.Duration) (string, time.Time, error) {
	issuedAt := time.Now().UTC()
	expiresAt := issuedAt.Add(ttl)
	claims := Claims{UserID: userID, Type: tokenType, Exp: expiresAt.Unix()}
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("marshal claims: %w", err)
	}

	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)
	secret := []byte(m.accessSecret)
	if tokenType == "refresh" {
		secret = []byte(m.refreshSecret)
	}
	signature := hmac.New(sha256.New, secret)
	_, _ = signature.Write([]byte(encodedPayload))
	signed := encodedPayload + "." + base64.RawURLEncoding.EncodeToString(signature.Sum(nil))
	return signed, expiresAt, nil
}

func (m TokenManager) IssueTokenPair(userID string) (TokenPair, error) {
	accessToken, accessExpiresAt, err := m.issueToken(userID, "access", m.accessTTL)
	if err != nil {
		return TokenPair{}, err
	}
	refreshToken, refreshExpiresAt, err := m.issueToken(userID, "refresh", m.refreshTTL)
	if err != nil {
		return TokenPair{}, err
	}
	if refreshExpiresAt.Before(accessExpiresAt) {
		return TokenPair{}, fmt.Errorf("refresh token expiry must be after access token expiry")
	}
	return TokenPair{AccessToken: accessToken, RefreshToken: refreshToken, ExpiresAt: accessExpiresAt, RefreshExpiresAt: refreshExpiresAt}, nil
}

func (m TokenManager) ParseToken(token string, expectedType string) (Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return Claims{}, fmt.Errorf("invalid token format")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return Claims{}, fmt.Errorf("decode payload: %w", err)
	}
	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return Claims{}, fmt.Errorf("decode claims: %w", err)
	}
	if claims.Type != expectedType {
		return Claims{}, fmt.Errorf("unexpected token type")
	}
	if time.Now().UTC().Unix() > claims.Exp {
		return Claims{}, fmt.Errorf("token expired")
	}
	secret := []byte(m.accessSecret)
	if expectedType == "refresh" {
		secret = []byte(m.refreshSecret)
	}
	signature := hmac.New(sha256.New, secret)
	_, _ = signature.Write([]byte(parts[0]))
	expectedSignature := base64.RawURLEncoding.EncodeToString(signature.Sum(nil))
	if !hmac.Equal([]byte(parts[1]), []byte(expectedSignature)) {
		return Claims{}, fmt.Errorf("invalid token signature")
	}
	return claims, nil
}

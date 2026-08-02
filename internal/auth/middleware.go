package auth

import (
	"context"
	"net/http"
	"strings"
)

// Context key type for injecting the authenticated user ID into request context.
type contextKey string

const userIDContextKey contextKey = "user_id"

// RequireAuth validates a bearer access token and places the user ID into the request context.
func RequireAuth(tokenManager TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			claims, err := tokenManager.ParseToken(strings.TrimSpace(strings.TrimPrefix(header, "Bearer ")), "access")
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := WithUserID(r.Context(), claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Middleware keeps the older helper name available for future compatibility.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

// WithUserID adds a user ID to the request context.
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}

// UserIDFromContext extracts the user ID from the request context.
func UserIDFromContext(ctx context.Context) string {
	if userID, ok := ctx.Value(userIDContextKey).(string); ok {
		return userID
	}
	return ""
}

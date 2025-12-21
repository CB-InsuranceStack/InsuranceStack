package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/CB-AccountStack/AccountStack/apps/api-accounts/internal/auth"
	"github.com/sirupsen/logrus"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const userIDKey contextKey = "userID"

// AuthMiddleware validates JWT tokens and extracts user information
func AuthMiddleware(logger *logrus.Logger) func(http.Handler) http.Handler {
	// Get JWT secret from environment
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-key-change-in-production"
		logger.Warn("JWT_SECRET not set, using default (not secure for production)")
	}

	jwtManager := auth.NewJWTManager(jwtSecret, 24*time.Hour)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip auth for health check and login endpoints
			if r.URL.Path == "/healthz" || r.URL.Path == "/login" {
				next.ServeHTTP(w, r)
				return
			}

			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.Warn("No authorization header provided")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check Bearer token format
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				logger.Warn("Invalid authorization header format")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Verify token
			claims, err := jwtManager.Verify(parts[1])
			if err != nil {
				logger.WithError(err).Warn("Invalid token")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Add user ID to request context
			ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)

			logger.WithField("userId", claims.UserID).Debug("User authenticated")

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID extracts the user ID from the request context
func GetUserID(r *http.Request) string {
	userID, ok := r.Context().Value(userIDKey).(string)
	if !ok {
		return "user-001" // Default fallback
	}
	return userID
}

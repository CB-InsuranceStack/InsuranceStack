package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/CB-AccountStack/AccountStack/apps/api-accounts/internal/auth"
	"github.com/CB-AccountStack/AccountStack/apps/api-accounts/internal/repository"
	"github.com/sirupsen/logrus"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	jwtManager   *auth.JWTManager
	logger       *logrus.Logger
	repo         *repository.Repository
	demoPassword string // Hashed password - same for all users in demo mode
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(repo *repository.Repository, logger *logrus.Logger) *AuthHandler {
	// Get JWT secret from environment
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-key-change-in-production"
		logger.Warn("JWT_SECRET not set, using default (not secure for production)")
	}

	// Get demo password from environment (all users use this in demo mode)
	password := os.Getenv("AUTH_PASSWORD")
	if password == "" {
		password = "demo123"
	}

	// Hash the password for comparison
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		logger.WithError(err).Fatal("Failed to hash password")
	}

	return &AuthHandler{
		jwtManager:   auth.NewJWTManager(jwtSecret, 24*time.Hour),
		logger:       logger,
		repo:         repo,
		demoPassword: hashedPassword,
	}
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expiresIn"` // seconds
	User      User   `json:"user"`
}

// User represents basic user info
type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode login request")
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Look up user by email in the repository
	user, err := h.repo.GetUserByEmail(req.Username)
	if err != nil {
		h.logger.WithField("username", req.Username).Warn("User not found")
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Validate password (all users use demo password in demo mode)
	if err := auth.VerifyPassword(h.demoPassword, req.Password); err != nil {
		h.logger.WithField("username", req.Username).Warn("Invalid password")
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := h.jwtManager.Generate(user.ID, req.Username)
	if err != nil {
		h.logger.WithError(err).Error("Failed to generate token")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := LoginResponse{
		Token:     token,
		ExpiresIn: 86400, // 24 hours in seconds
		User: User{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	h.logger.WithField("username", req.Username).Info("User logged in successfully")
}

package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/CB-AccountStack/AccountStack/apps/api-accounts/internal/middleware"
	"github.com/CB-AccountStack/AccountStack/apps/api-accounts/internal/services"
	"github.com/sirupsen/logrus"
)

// UserHandler handles user-related requests
type UserHandler struct {
	userService *services.UserService
	logger      *logrus.Logger
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *services.UserService, logger *logrus.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// GetMe handles GET /me - returns current user info
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		h.logger.WithError(err).WithField("userId", userID).Error("Failed to get user")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "not_found",
			Message: "User not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/CB-AccountStack/AccountStack/apps/api-insights/internal/middleware"
	"github.com/CB-AccountStack/AccountStack/apps/api-insights/internal/services"
	"github.com/sirupsen/logrus"
)

// AlertsHandler handles alerts-related requests
type AlertsHandler struct {
	service *services.AlertsService
	logger  *logrus.Logger
}

// NewAlertsHandler creates a new alerts handler
func NewAlertsHandler(service *services.AlertsService, logger *logrus.Logger) *AlertsHandler {
	return &AlertsHandler{
		service: service,
		logger:  logger,
	}
}

// GetAlerts handles GET /alerts - list all alerts for the authenticated user
// Returns 503 Service Unavailable if the alerts feature is disabled
func (h *AlertsHandler) GetAlerts(w http.ResponseWriter, r *http.Request) {
	// Check if alerts feature is enabled
	if !h.service.IsAlertsEnabled() {
		h.logger.Warn("Alerts endpoint accessed but feature is disabled")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "Service Unavailable",
			"message": "Alerts feature is currently disabled",
		})
		return
	}

	userID := middleware.GetUserID(r)

	alerts, err := h.service.GetAlertsByUserID(userID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get alerts")
		http.Error(w, "Failed to retrieve alerts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(alerts)
}

package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/CB-AccountStack/AccountStack/apps/api-insights/internal/middleware"
	"github.com/CB-AccountStack/AccountStack/apps/api-insights/internal/services"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// InsightsHandler handles insights-related requests
type InsightsHandler struct {
	service *services.InsightsService
	logger  *logrus.Logger
}

// NewInsightsHandler creates a new insights handler
func NewInsightsHandler(service *services.InsightsService, logger *logrus.Logger) *InsightsHandler {
	return &InsightsHandler{
		service: service,
		logger:  logger,
	}
}

// GetInsights handles GET /insights - list all insights for the authenticated user
func (h *InsightsHandler) GetInsights(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	insights, err := h.service.GetInsightsByUserID(userID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get insights")
		http.Error(w, "Failed to retrieve insights", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(insights)
}

// GetInsightByID handles GET /insights/{id} - get a specific insight by ID
func (h *InsightsHandler) GetInsightByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	insightID := vars["id"]
	userID := middleware.GetUserID(r)

	insight, err := h.service.GetInsightByID(insightID)
	if err != nil {
		h.logger.WithError(err).WithField("insightId", insightID).Error("Failed to get insight")
		http.Error(w, "Insight not found", http.StatusNotFound)
		return
	}

	// Verify the insight belongs to the authenticated user
	if insight.UserID != userID {
		h.logger.WithFields(logrus.Fields{
			"insightId":  insightID,
			"userId":     userID,
			"ownerId":    insight.UserID,
		}).Warn("User attempted to access another user's insight")
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(insight)
}

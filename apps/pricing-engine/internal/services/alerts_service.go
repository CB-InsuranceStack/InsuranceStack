package services

import (
	"fmt"

	"github.com/CB-AccountStack/AccountStack/apps/api-insights/internal/features"
	"github.com/CB-AccountStack/AccountStack/apps/api-insights/internal/models"
	"github.com/CB-AccountStack/AccountStack/apps/api-insights/internal/repository"
	"github.com/sirupsen/logrus"
)

// AlertsService handles business logic for alerts
type AlertsService struct {
	repo   *repository.Repository
	flags  *features.Flags
	logger *logrus.Logger
}

// NewAlertsService creates a new alerts service
func NewAlertsService(repo *repository.Repository, flags *features.Flags, logger *logrus.Logger) *AlertsService {
	return &AlertsService{
		repo:   repo,
		flags:  flags,
		logger: logger,
	}
}

// GetAlertsByUserID retrieves alerts for a user
// Returns an error if alerts are disabled via feature flag
func (s *AlertsService) GetAlertsByUserID(userID string) ([]*models.Alert, error) {
	// Check if alerts feature is enabled
	if !s.flags.IsAlertsEnabled() {
		s.logger.WithField("userId", userID).Warn("Alerts feature is disabled")
		return nil, fmt.Errorf("alerts feature is currently disabled")
	}

	alerts, err := s.repo.GetAlertsByUserID(userID)
	if err != nil {
		s.logger.WithError(err).Error("Failed to retrieve alerts")
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"userId":      userID,
		"alertCount":  len(alerts),
	}).Debug("Retrieved alerts for user")

	return alerts, nil
}

// IsAlertsEnabled returns whether the alerts feature is currently enabled
func (s *AlertsService) IsAlertsEnabled() bool {
	return s.flags.IsAlertsEnabled()
}

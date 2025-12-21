package services

import (
	"github.com/CB-AccountStack/AccountStack/apps/api-insights/internal/features"
	"github.com/CB-AccountStack/AccountStack/apps/api-insights/internal/models"
	"github.com/CB-AccountStack/AccountStack/apps/api-insights/internal/repository"
	"github.com/sirupsen/logrus"
)

// InsightsService handles business logic for insights
type InsightsService struct {
	repo   *repository.Repository
	flags  *features.Flags
	logger *logrus.Logger
}

// NewInsightsService creates a new insights service
func NewInsightsService(repo *repository.Repository, flags *features.Flags, logger *logrus.Logger) *InsightsService {
	return &InsightsService{
		repo:   repo,
		flags:  flags,
		logger: logger,
	}
}

// GetInsightsByUserID retrieves insights for a user, applying V2 modifications if enabled
func (s *InsightsService) GetInsightsByUserID(userID string) ([]*models.Insight, error) {
	insights, err := s.repo.GetInsightsByUserID(userID)
	if err != nil {
		s.logger.WithError(err).Error("Failed to retrieve insights")
		return nil, err
	}

	// Apply V2 algorithm modifications if flag is enabled
	if s.flags.IsInsightsV2Enabled() {
		insights = s.applyV2Algorithm(insights)
		s.logger.WithField("userId", userID).Debug("Applied V2 insights algorithm")
	}

	return insights, nil
}

// GetInsightByID retrieves a specific insight by ID, applying V2 modifications if enabled
func (s *InsightsService) GetInsightByID(insightID string) (*models.Insight, error) {
	insight, err := s.repo.GetInsightByID(insightID)
	if err != nil {
		s.logger.WithError(err).Error("Failed to retrieve insight")
		return nil, err
	}

	// Apply V2 algorithm modifications if flag is enabled
	if s.flags.IsInsightsV2Enabled() {
		insight = s.applyV2ToSingleInsight(insight)
		s.logger.WithField("insightId", insightID).Debug("Applied V2 insights algorithm")
	}

	return insight, nil
}

// applyV2Algorithm applies the V2 insights calculation algorithm
// In this demo, we append "(V2)" to titles to show the new algorithm is active
func (s *InsightsService) applyV2Algorithm(insights []*models.Insight) []*models.Insight {
	modifiedInsights := make([]*models.Insight, len(insights))

	for i, insight := range insights {
		// Create a copy to avoid modifying the original
		modified := *insight
		modified.Title = insight.Title + " (V2)"
		modifiedInsights[i] = &modified
	}

	return modifiedInsights
}

// applyV2ToSingleInsight applies V2 algorithm to a single insight
func (s *InsightsService) applyV2ToSingleInsight(insight *models.Insight) *models.Insight {
	// Create a copy to avoid modifying the original
	modified := *insight
	modified.Title = insight.Title + " (V2)"
	return &modified
}

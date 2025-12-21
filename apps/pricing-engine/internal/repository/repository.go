package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/CB-AccountStack/AccountStack/apps/api-insights/internal/models"
	"github.com/sirupsen/logrus"
)

// Repository provides data access for insights and alerts
type Repository struct {
	insights map[string]*models.Insight
	alerts   map[string]*models.Alert
	mu       sync.RWMutex
	logger   *logrus.Logger
}

// NewRepository creates a new repository and loads data from JSON files
func NewRepository(dataPath string, logger *logrus.Logger) (*Repository, error) {
	repo := &Repository{
		insights: make(map[string]*models.Insight),
		alerts:   make(map[string]*models.Alert),
		logger:   logger,
	}

	// Load insights
	if err := repo.loadInsights(filepath.Join(dataPath, "insights.json")); err != nil {
		return nil, fmt.Errorf("failed to load insights: %w", err)
	}

	// Generate alerts from insights (alerts are derived from high-priority insights)
	repo.generateAlerts()

	logger.Infof("Loaded %d insights and generated %d alerts from %s", len(repo.insights), len(repo.alerts), dataPath)

	return repo, nil
}

// loadInsights loads insights from a JSON file
func (r *Repository) loadInsights(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var insights []*models.Insight
	if err := json.Unmarshal(data, &insights); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, insight := range insights {
		r.insights[insight.ID] = insight
	}

	return nil
}

// generateAlerts creates alerts from high-priority actionable insights
func (r *Repository) generateAlerts() {
	r.mu.Lock()
	defer r.mu.Unlock()

	alertCounter := 1
	for _, insight := range r.insights {
		// Create alerts for medium and high severity actionable insights
		if insight.Actionable && (insight.Severity == "medium" || insight.Severity == "high") {
			alert := &models.Alert{
				ID:        fmt.Sprintf("alert-%03d", alertCounter),
				UserID:    insight.UserID,
				Type:      insight.Type,
				Title:     insight.Title,
				Message:   insight.Description,
				Priority:  mapSeverityToPriority(insight.Severity),
				CreatedAt: insight.CreatedAt,
				Read:      false,
			}
			r.alerts[alert.ID] = alert
			alertCounter++
		}
	}
}

// mapSeverityToPriority converts insight severity to alert priority
func mapSeverityToPriority(severity string) string {
	switch severity {
	case "high":
		return "critical"
	case "medium":
		return "high"
	default:
		return "medium"
	}
}

// GetInsightByID retrieves an insight by ID
func (r *Repository) GetInsightByID(insightID string) (*models.Insight, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	insight, exists := r.insights[insightID]
	if !exists {
		return nil, fmt.Errorf("insight not found")
	}

	return insight, nil
}

// GetInsightsByUserID retrieves all insights for a specific user
func (r *Repository) GetInsightsByUserID(userID string) ([]*models.Insight, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var userInsights []*models.Insight
	for _, insight := range r.insights {
		if insight.UserID == userID {
			userInsights = append(userInsights, insight)
		}
	}

	return userInsights, nil
}

// GetAlertsByUserID retrieves all alerts for a specific user
func (r *Repository) GetAlertsByUserID(userID string) ([]*models.Alert, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var userAlerts []*models.Alert
	for _, alert := range r.alerts {
		if alert.UserID == userID {
			userAlerts = append(userAlerts, alert)
		}
	}

	return userAlerts, nil
}

// GetAllInsights returns all insights (for testing/admin purposes)
func (r *Repository) GetAllInsights() []*models.Insight {
	r.mu.RLock()
	defer r.mu.RUnlock()

	insights := make([]*models.Insight, 0, len(r.insights))
	for _, insight := range r.insights {
		insights = append(insights, insight)
	}

	return insights
}

// GetAllAlerts returns all alerts (for testing/admin purposes)
func (r *Repository) GetAllAlerts() []*models.Alert {
	r.mu.RLock()
	defer r.mu.RUnlock()

	alerts := make([]*models.Alert, 0, len(r.alerts))
	for _, alert := range r.alerts {
		alerts = append(alerts, alert)
	}

	return alerts
}

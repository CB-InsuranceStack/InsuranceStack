package models

import (
	"testing"
	"time"
)

func TestInsightCreation(t *testing.T) {
	rec := "Consider reducing spending"

	tests := []struct {
		name     string
		insight  *Insight
		wantType string
	}{
		{
			name: "spending insight",
			insight: &Insight{
				ID:             "ins-001",
				UserID:         "user-001",
				Type:           "spending",
				Category:       "budget",
				Title:          "High Spending Alert",
				Description:    "Your spending is above average",
				Severity:       "warning",
				CreatedAt:      time.Now(),
				Actionable:     true,
				Recommendation: &rec,
			},
			wantType: "spending",
		},
		{
			name: "savings insight",
			insight: &Insight{
				ID:          "ins-002",
				UserID:      "user-002",
				Type:        "savings",
				Category:    "growth",
				Title:       "Savings Milestone",
				Description: "You've reached 10k savings",
				Severity:    "info",
				CreatedAt:   time.Now(),
				Actionable:  false,
			},
			wantType: "savings",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.insight.Type != tt.wantType {
				t.Errorf("Type mismatch: got %v, want %v", tt.insight.Type, tt.wantType)
			}
			if tt.insight.ID == "" {
				t.Error("ID should not be empty")
			}
		})
	}
}

func TestInsightTypes(t *testing.T) {
	tests := []struct {
		name         string
		insightType  string
	}{
		{"spending insight", "spending"},
		{"savings insight", "savings"},
		{"budget insight", "budget"},
		{"goal insight", "goal"},
		{"trend insight", "trend"},
		{"recommendation insight", "recommendation"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			insight := &Insight{
				ID:   "test-id",
				Type: tt.insightType,
			}
			if insight.Type != tt.insightType {
				t.Errorf("Type mismatch: got %v, want %v", insight.Type, tt.insightType)
			}
		})
	}
}

func TestInsightSeverity(t *testing.T) {
	tests := []struct {
		name     string
		severity string
	}{
		{"info severity", "info"},
		{"warning severity", "warning"},
		{"critical severity", "critical"},
		{"low severity", "low"},
		{"medium severity", "medium"},
		{"high severity", "high"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			insight := &Insight{
				ID:       "test-id",
				Severity: tt.severity,
			}
			if insight.Severity != tt.severity {
				t.Errorf("Severity mismatch: got %v, want %v", insight.Severity, tt.severity)
			}
		})
	}
}

func TestInsightCategory(t *testing.T) {
	tests := []struct {
		name     string
		category string
	}{
		{"budget category", "budget"},
		{"investment category", "investment"},
		{"debt category", "debt"},
		{"income category", "income"},
		{"expense category", "expense"},
		{"goal category", "goal"},
		{"saving category", "saving"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			insight := &Insight{
				ID:       "test-id",
				Category: tt.category,
			}
			if insight.Category != tt.category {
				t.Errorf("Category mismatch: got %v, want %v", insight.Category, tt.category)
			}
		})
	}
}

func TestInsightActionable(t *testing.T) {
	tests := []struct {
		name       string
		actionable bool
		hasRec     bool
	}{
		{"actionable with recommendation", true, true},
		{"actionable without recommendation", true, false},
		{"not actionable", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rec *string
			if tt.hasRec {
				r := "Take action"
				rec = &r
			}

			insight := &Insight{
				ID:             "test-id",
				Actionable:     tt.actionable,
				Recommendation: rec,
			}

			if insight.Actionable != tt.actionable {
				t.Errorf("Actionable mismatch: got %v, want %v", insight.Actionable, tt.actionable)
			}
			if tt.hasRec && insight.Recommendation == nil {
				t.Error("Expected recommendation to be set")
			}
			if !tt.hasRec && insight.Recommendation != nil {
				t.Error("Expected recommendation to be nil")
			}
		})
	}
}

func TestInsightTimestamps(t *testing.T) {
	now := time.Now()
	past := now.AddDate(0, 0, -7)
	future := now.AddDate(0, 0, 7)

	tests := []struct {
		name      string
		createdAt time.Time
		valid     bool
	}{
		{"current time", now, true},
		{"past time", past, true},
		{"future time", future, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			insight := &Insight{
				ID:        "test-id",
				CreatedAt: tt.createdAt,
			}
			if insight.CreatedAt.IsZero() {
				t.Error("CreatedAt should not be zero")
			}
		})
	}
}

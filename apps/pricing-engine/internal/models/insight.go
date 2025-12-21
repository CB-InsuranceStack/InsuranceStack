package models

import "time"

// Insight represents a financial insight or recommendation for a user
type Insight struct {
	ID             string     `json:"id"`
	UserID         string     `json:"userId"`
	Type           string     `json:"type"`
	Category       string     `json:"category"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	Severity       string     `json:"severity"`
	CreatedAt      time.Time  `json:"createdAt"`
	Actionable     bool       `json:"actionable"`
	Recommendation *string    `json:"recommendation"`
}

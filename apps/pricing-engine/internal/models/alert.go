package models

import "time"

// Alert represents a real-time alert for a user
type Alert struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	Message     string    `json:"message"`
	Priority    string    `json:"priority"`
	CreatedAt   time.Time `json:"createdAt"`
	Read        bool      `json:"read"`
	ActionURL   *string   `json:"actionUrl,omitempty"`
}

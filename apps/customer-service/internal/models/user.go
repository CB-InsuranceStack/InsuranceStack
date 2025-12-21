package models

import "time"

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Country   string    `json:"country"` // ISO 3166-1 alpha-2 country code (US, UK, FR, etc.)
	CreatedAt time.Time `json:"createdAt"`
	LastLogin time.Time `json:"lastLogin"`
}

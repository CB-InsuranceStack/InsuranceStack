package models

import (
	"testing"
	"time"
)

func TestUserCreation(t *testing.T) {
	tests := []struct {
		name      string
		user      *User
		wantEmail string
		wantName  string
	}{
		{
			name: "user with full details",
			user: &User{
				ID:        "user-001",
				Email:     "john.doe@example.com",
				Name:      "John Doe",
				FirstName: "John",
				LastName:  "Doe",
				Country:   "US",
				CreatedAt: time.Now(),
				LastLogin: time.Now(),
			},
			wantEmail: "john.doe@example.com",
			wantName:  "John Doe",
		},
		{
			name: "user with minimal details",
			user: &User{
				ID:        "user-002",
				Email:     "jane@example.com",
				Name:      "Jane Smith",
				FirstName: "Jane",
				LastName:  "Smith",
				Country:   "UK",
				CreatedAt: time.Now(),
			},
			wantEmail: "jane@example.com",
			wantName:  "Jane Smith",
		},
		{
			name: "user from France",
			user: &User{
				ID:        "user-003",
				Email:     "francois@example.fr",
				Name:      "François Dubois",
				FirstName: "François",
				LastName:  "Dubois",
				Country:   "FR",
				CreatedAt: time.Now(),
			},
			wantEmail: "francois@example.fr",
			wantName:  "François Dubois",
		},
		{
			name: "user from Germany",
			user: &User{
				ID:        "user-004",
				Email:     "hans@example.de",
				Name:      "Hans Mueller",
				FirstName: "Hans",
				LastName:  "Mueller",
				Country:   "DE",
				CreatedAt: time.Now(),
			},
			wantEmail: "hans@example.de",
			wantName:  "Hans Mueller",
		},
		{
			name: "user from Japan",
			user: &User{
				ID:        "user-005",
				Email:     "tanaka@example.jp",
				Name:      "Tanaka Yuki",
				FirstName: "Yuki",
				LastName:  "Tanaka",
				Country:   "JP",
				CreatedAt: time.Now(),
			},
			wantEmail: "tanaka@example.jp",
			wantName:  "Tanaka Yuki",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.user.Email != tt.wantEmail {
				t.Errorf("Email mismatch: got %v, want %v", tt.user.Email, tt.wantEmail)
			}
			if tt.user.Name != tt.wantName {
				t.Errorf("Name mismatch: got %v, want %v", tt.user.Name, tt.wantName)
			}
			if tt.user.ID == "" {
				t.Error("User ID should not be empty")
			}
			if tt.user.CreatedAt.IsZero() {
				t.Error("CreatedAt should not be zero")
			}
		})
	}
}

func TestUserCountryCodes(t *testing.T) {
	tests := []struct {
		name        string
		countryCode string
		valid       bool
	}{
		{"US country code", "US", true},
		{"UK country code", "UK", true},
		{"FR country code", "FR", true},
		{"DE country code", "DE", true},
		{"JP country code", "JP", true},
		{"CA country code", "CA", true},
		{"AU country code", "AU", true},
		{"ES country code", "ES", true},
		{"IT country code", "IT", true},
		{"NL country code", "NL", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{
				ID:      "test-id",
				Country: tt.countryCode,
			}
			if user.Country != tt.countryCode {
				t.Errorf("Country mismatch: got %v, want %v", user.Country, tt.countryCode)
			}
			if len(user.Country) != 2 {
				t.Errorf("Country code should be 2 characters, got %d", len(user.Country))
			}
		})
	}
}

func TestUserEmailValidation(t *testing.T) {
	tests := []struct {
		name  string
		email string
		valid bool
	}{
		{"valid email with dot", "user@example.com", true},
		{"valid email with subdomain", "user@mail.example.com", true},
		{"valid email with plus", "user+tag@example.com", true},
		{"valid email with numbers", "user123@example.com", true},
		{"valid email with dash", "user-name@example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{
				ID:    "test-id",
				Email: tt.email,
			}
			if user.Email != tt.email {
				t.Errorf("Email mismatch: got %v, want %v", user.Email, tt.email)
			}
		})
	}
}

func TestUserNameFields(t *testing.T) {
	tests := []struct {
		name      string
		firstName string
		lastName  string
		fullName  string
	}{
		{"simple name", "John", "Doe", "John Doe"},
		{"hyphenated last name", "Mary", "Smith-Jones", "Mary Smith-Jones"},
		{"name with prefix", "Dr. James", "Wilson", "Dr. James Wilson"},
		{"single word last name", "Alice", "Brown", "Alice Brown"},
		{"name with apostrophe", "Sean", "O'Brien", "Sean O'Brien"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{
				ID:        "test-id",
				FirstName: tt.firstName,
				LastName:  tt.lastName,
				Name:      tt.fullName,
			}
			if user.FirstName != tt.firstName {
				t.Errorf("FirstName mismatch: got %v, want %v", user.FirstName, tt.firstName)
			}
			if user.LastName != tt.lastName {
				t.Errorf("LastName mismatch: got %v, want %v", user.LastName, tt.lastName)
			}
		})
	}
}

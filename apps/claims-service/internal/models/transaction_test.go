package models

import (
	"testing"
	"time"
)

func TestTransactionMatches(t *testing.T) {
	// Create a sample transaction
	txn := &Transaction{
		ID:          "txn-001",
		AccountID:   "acc-001",
		Date:        time.Date(2024, 12, 10, 10, 0, 0, 0, time.UTC),
		Description: "Test Transaction",
		Amount:      -50.00,
		Category:    "shopping",
		Merchant:    "Test Store",
		Status:      "completed",
		Type:        "debit",
	}

	tests := []struct {
		name     string
		filters  *TransactionFilters
		expected bool
	}{
		{
			name: "matches account ID",
			filters: &TransactionFilters{
				AccountID: "acc-001",
			},
			expected: true,
		},
		{
			name: "does not match account ID",
			filters: &TransactionFilters{
				AccountID: "acc-002",
			},
			expected: false,
		},
		{
			name: "matches date range",
			filters: &TransactionFilters{
				AccountID: "acc-001",
				StartDate: timePtr(time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC)),
				EndDate:   timePtr(time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)),
			},
			expected: true,
		},
		{
			name: "does not match date range - before start",
			filters: &TransactionFilters{
				AccountID: "acc-001",
				StartDate: timePtr(time.Date(2024, 12, 15, 0, 0, 0, 0, time.UTC)),
			},
			expected: false,
		},
		{
			name: "matches category",
			filters: &TransactionFilters{
				AccountID: "acc-001",
				Category:  "shopping",
			},
			expected: true,
		},
		{
			name: "does not match category",
			filters: &TransactionFilters{
				AccountID: "acc-001",
				Category:  "food_dining",
			},
			expected: false,
		},
		{
			name: "matches amount range",
			filters: &TransactionFilters{
				AccountID: "acc-001",
				MinAmount: floatPtr(-100.0),
				MaxAmount: floatPtr(-10.0),
			},
			expected: true,
		},
		{
			name: "does not match amount range - too small",
			filters: &TransactionFilters{
				AccountID: "acc-001",
				MinAmount: floatPtr(-40.0),
			},
			expected: false,
		},
		{
			name: "matches complex filter",
			filters: &TransactionFilters{
				AccountID: "acc-001",
				Category:  "shopping",
				MinAmount: floatPtr(-100.0),
				MaxAmount: floatPtr(0.0),
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := txn.Matches(tt.filters)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// Helper function to create time pointers
func timePtr(t time.Time) *time.Time {
	return &t
}

// Helper function to create float64 pointers
func floatPtr(f float64) *float64 {
	return &f
}

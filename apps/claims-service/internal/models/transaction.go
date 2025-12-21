package models

import "time"

// Transaction represents a financial transaction
type Transaction struct {
	ID          string    `json:"id"`
	AccountID   string    `json:"accountId"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Merchant    string    `json:"merchant"`
	Status      string    `json:"status"`
	Type        string    `json:"type"`
}

// TransactionFilters represents filters for transaction queries
type TransactionFilters struct {
	AccountID string
	StartDate *time.Time
	EndDate   *time.Time
	Category  string
	MinAmount *float64
	MaxAmount *float64
}

// Matches checks if a transaction matches the given filters
func (t *Transaction) Matches(filters *TransactionFilters) bool {
	// Account ID filter (always allowed)
	if filters.AccountID != "" && t.AccountID != filters.AccountID {
		return false
	}

	// Date range filter
	if filters.StartDate != nil && t.Date.Before(*filters.StartDate) {
		return false
	}
	if filters.EndDate != nil && t.Date.After(*filters.EndDate) {
		return false
	}

	// Category filter
	if filters.Category != "" && t.Category != filters.Category {
		return false
	}

	// Amount range filter
	if filters.MinAmount != nil && t.Amount < *filters.MinAmount {
		return false
	}
	if filters.MaxAmount != nil && t.Amount > *filters.MaxAmount {
		return false
	}

	return true
}

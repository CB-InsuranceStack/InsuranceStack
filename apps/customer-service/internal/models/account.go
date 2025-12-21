package models

import "time"

// Account represents a bank account in the system
type Account struct {
	ID            string    `json:"id"`
	UserID        string    `json:"userId"`
	AccountNumber string    `json:"accountNumber"`
	AccountType   string    `json:"accountType"`
	AccountName   string    `json:"accountName"`
	Balance       float64   `json:"balance"`
	Currency      string    `json:"currency"`
	CreditLimit   *float64  `json:"creditLimit,omitempty"`
	Status        string    `json:"status"`
	OpenedDate    time.Time `json:"openedDate"`
	LastActivity  time.Time `json:"lastActivity"`
}

// AccountResponse represents an account in API responses with optional masking
type AccountResponse struct {
	ID            string    `json:"id"`
	UserID        string    `json:"userId"`
	AccountNumber string    `json:"accountNumber"`
	AccountType   string    `json:"accountType"`
	AccountName   string    `json:"accountName"`
	Balance       any       `json:"balance"` // Can be float64 or string (masked)
	Currency      string    `json:"currency"`
	CreditLimit   any       `json:"creditLimit,omitempty"` // Can be float64 or string (masked)
	Status        string    `json:"status"`
	OpenedDate    time.Time `json:"openedDate"`
	LastActivity  time.Time `json:"lastActivity"`
}

// ToResponse converts an Account to AccountResponse with optional masking and currency override
func (a *Account) ToResponse(maskAmounts bool, currency string) AccountResponse {
	resp := AccountResponse{
		ID:            a.ID,
		UserID:        a.UserID,
		AccountNumber: a.AccountNumber,
		AccountType:   a.AccountType,
		AccountName:   a.AccountName,
		Currency:      currency, // Use feature flag currency
		Status:        a.Status,
		OpenedDate:    a.OpenedDate,
		LastActivity:  a.LastActivity,
	}

	if maskAmounts {
		resp.Balance = "***.**"
		if a.CreditLimit != nil {
			resp.CreditLimit = "***.**"
		}
	} else {
		resp.Balance = a.Balance
		if a.CreditLimit != nil {
			resp.CreditLimit = *a.CreditLimit
		}
	}

	return resp
}

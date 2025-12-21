package models

import (
	"testing"
	"time"
)

func TestAccountToResponse(t *testing.T) {
	creditLimit := 5000.0
	now := time.Now()

	tests := []struct {
		name        string
		account     *Account
		maskAmounts bool
		currency    string
		wantMasked  bool
	}{
		{
			name: "unmasked account with credit limit",
			account: &Account{
				ID:            "acc-001",
				UserID:        "user-001",
				AccountNumber: "1234567890",
				AccountType:   "checking",
				AccountName:   "Main Checking",
				Balance:       1500.50,
				Currency:      "USD",
				CreditLimit:   &creditLimit,
				Status:        "active",
				OpenedDate:    now,
				LastActivity:  now,
			},
			maskAmounts: false,
			currency:    "USD",
			wantMasked:  false,
		},
		{
			name: "masked account with credit limit",
			account: &Account{
				ID:            "acc-002",
				UserID:        "user-002",
				AccountNumber: "9876543210",
				AccountType:   "savings",
				AccountName:   "Savings Account",
				Balance:       25000.00,
				Currency:      "EUR",
				CreditLimit:   &creditLimit,
				Status:        "active",
				OpenedDate:    now,
				LastActivity:  now,
			},
			maskAmounts: true,
			currency:    "EUR",
			wantMasked:  true,
		},
		{
			name: "unmasked account without credit limit",
			account: &Account{
				ID:            "acc-003",
				UserID:        "user-003",
				AccountNumber: "1111222233",
				AccountType:   "checking",
				AccountName:   "Basic Checking",
				Balance:       500.00,
				Currency:      "GBP",
				CreditLimit:   nil,
				Status:        "active",
				OpenedDate:    now,
				LastActivity:  now,
			},
			maskAmounts: false,
			currency:    "GBP",
			wantMasked:  false,
		},
		{
			name: "masked account without credit limit",
			account: &Account{
				ID:            "acc-004",
				UserID:        "user-004",
				AccountNumber: "4444555566",
				AccountType:   "savings",
				AccountName:   "Emergency Fund",
				Balance:       10000.00,
				Currency:      "USD",
				CreditLimit:   nil,
				Status:        "active",
				OpenedDate:    now,
				LastActivity:  now,
			},
			maskAmounts: true,
			currency:    "USD",
			wantMasked:  true,
		},
		{
			name: "currency override",
			account: &Account{
				ID:            "acc-005",
				UserID:        "user-005",
				AccountNumber: "7777888899",
				AccountType:   "checking",
				AccountName:   "Multi-Currency",
				Balance:       3000.00,
				Currency:      "USD",
				CreditLimit:   nil,
				Status:        "active",
				OpenedDate:    now,
				LastActivity:  now,
			},
			maskAmounts: false,
			currency:    "JPY", // Override currency
			wantMasked:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := tt.account.ToResponse(tt.maskAmounts, tt.currency)

			// Check basic fields are preserved
			if resp.ID != tt.account.ID {
				t.Errorf("ID mismatch: got %v, want %v", resp.ID, tt.account.ID)
			}
			if resp.UserID != tt.account.UserID {
				t.Errorf("UserID mismatch: got %v, want %v", resp.UserID, tt.account.UserID)
			}
			if resp.Currency != tt.currency {
				t.Errorf("Currency mismatch: got %v, want %v", resp.Currency, tt.currency)
			}

			// Check balance masking
			if tt.wantMasked {
				if resp.Balance != "***.**" {
					t.Errorf("Expected masked balance, got %v", resp.Balance)
				}
				if tt.account.CreditLimit != nil && resp.CreditLimit != "***.**" {
					t.Errorf("Expected masked credit limit, got %v", resp.CreditLimit)
				}
			} else {
				if balFloat, ok := resp.Balance.(float64); !ok || balFloat != tt.account.Balance {
					t.Errorf("Balance mismatch: got %v, want %v", resp.Balance, tt.account.Balance)
				}
				if tt.account.CreditLimit != nil {
					if clFloat, ok := resp.CreditLimit.(float64); !ok || clFloat != *tt.account.CreditLimit {
						t.Errorf("CreditLimit mismatch: got %v, want %v", resp.CreditLimit, *tt.account.CreditLimit)
					}
				}
			}
		})
	}
}

func TestAccountTypes(t *testing.T) {
	tests := []struct {
		name        string
		accountType string
		valid       bool
	}{
		{"checking account", "checking", true},
		{"savings account", "savings", true},
		{"credit account", "credit", true},
		{"investment account", "investment", true},
		{"loan account", "loan", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc := &Account{
				ID:          "test-id",
				AccountType: tt.accountType,
			}
			if acc.AccountType != tt.accountType {
				t.Errorf("AccountType mismatch: got %v, want %v", acc.AccountType, tt.accountType)
			}
		})
	}
}

func TestAccountStatus(t *testing.T) {
	tests := []struct {
		name   string
		status string
		valid  bool
	}{
		{"active status", "active", true},
		{"inactive status", "inactive", true},
		{"closed status", "closed", true},
		{"frozen status", "frozen", true},
		{"suspended status", "suspended", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc := &Account{
				ID:     "test-id",
				Status: tt.status,
			}
			if acc.Status != tt.status {
				t.Errorf("Status mismatch: got %v, want %v", acc.Status, tt.status)
			}
		})
	}
}

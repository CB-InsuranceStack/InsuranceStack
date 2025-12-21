package services

import (
	"testing"
	"time"

	"github.com/CB-AccountStack/AccountStack/apps/api-transactions/internal/models"
)

func TestFilterTransactionsByAccount(t *testing.T) {
	txns := []*models.Transaction{
		{ID: "tx1", AccountID: "acc1", Amount: 100.0},
		{ID: "tx2", AccountID: "acc1", Amount: -50.0},
		{ID: "tx3", AccountID: "acc2", Amount: 200.0},
		{ID: "tx4", AccountID: "acc1", Amount: -25.0},
	}

	tests := []struct {
		name        string
		accountID   string
		wantCount   int
	}{
		{"account 1 transactions", "acc1", 3},
		{"account 2 transactions", "acc2", 1},
		{"non-existent account", "acc3", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered := filterByAccount(txns, tt.accountID)
			if len(filtered) != tt.wantCount {
				t.Errorf("Expected %d transactions, got %d", tt.wantCount, len(filtered))
			}
		})
	}
}

func filterByAccount(txns []*models.Transaction, accountID string) []*models.Transaction {
	var result []*models.Transaction
	for _, tx := range txns {
		if tx.AccountID == accountID {
			result = append(result, tx)
		}
	}
	return result
}

func TestFilterTransactionsByType(t *testing.T) {
	txns := []*models.Transaction{
		{ID: "tx1", Type: "debit", Amount: -50.0},
		{ID: "tx2", Type: "credit", Amount: 100.0},
		{ID: "tx3", Type: "debit", Amount: -25.0},
		{ID: "tx4", Type: "credit", Amount: 200.0},
		{ID: "tx5", Type: "debit", Amount: -75.0},
	}

	tests := []struct {
		name      string
		txType    string
		wantCount int
	}{
		{"debit transactions", "debit", 3},
		{"credit transactions", "credit", 2},
		{"transfer transactions", "transfer", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered := filterByType(txns, tt.txType)
			if len(filtered) != tt.wantCount {
				t.Errorf("Expected %d transactions, got %d", tt.wantCount, len(filtered))
			}
		})
	}
}

func filterByType(txns []*models.Transaction, txType string) []*models.Transaction {
	var result []*models.Transaction
	for _, tx := range txns {
		if tx.Type == txType {
			result = append(result, tx)
		}
	}
	return result
}

func TestFilterTransactionsByCategory(t *testing.T) {
	txns := []*models.Transaction{
		{ID: "tx1", Category: "groceries"},
		{ID: "tx2", Category: "entertainment"},
		{ID: "tx3", Category: "groceries"},
		{ID: "tx4", Category: "transport"},
		{ID: "tx5", Category: "groceries"},
	}

	tests := []struct {
		name      string
		category  string
		wantCount int
	}{
		{"groceries", "groceries", 3},
		{"entertainment", "entertainment", 1},
		{"transport", "transport", 1},
		{"utilities", "utilities", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered := filterByCategory(txns, tt.category)
			if len(filtered) != tt.wantCount {
				t.Errorf("Expected %d transactions, got %d", tt.wantCount, len(filtered))
			}
		})
	}
}

func filterByCategory(txns []*models.Transaction, category string) []*models.Transaction {
	var result []*models.Transaction
	for _, tx := range txns {
		if tx.Category == category {
			result = append(result, tx)
		}
	}
	return result
}

func TestFilterTransactionsByDateRange(t *testing.T) {
	now := time.Now()
	txns := []*models.Transaction{
		{ID: "tx1", Date: now.AddDate(0, 0, -10)},
		{ID: "tx2", Date: now.AddDate(0, 0, -5)},
		{ID: "tx3", Date: now.AddDate(0, 0, -2)},
		{ID: "tx4", Date: now.AddDate(0, 0, -1)},
	}

	tests := []struct {
		name       string
		startDays  int
		endDays    int
		wantCount  int
	}{
		{"last 7 days", -7, 0, 3},
		{"last 3 days", -3, 0, 2},
		{"last 1 day", -1, 0, 1},
		{"last 30 days", -30, 0, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := now.AddDate(0, 0, tt.startDays)
			end := now.AddDate(0, 0, tt.endDays)
			filtered := filterByDateRange(txns, start, end)
			if len(filtered) != tt.wantCount {
				t.Errorf("Expected %d transactions, got %d", tt.wantCount, len(filtered))
			}
		})
	}
}

func filterByDateRange(txns []*models.Transaction, start, end time.Time) []*models.Transaction {
	var result []*models.Transaction
	for _, tx := range txns {
		if (tx.Date.After(start) || tx.Date.Equal(start)) && (tx.Date.Before(end) || tx.Date.Equal(end)) {
			result = append(result, tx)
		}
	}
	return result
}

func TestFilterTransactionsByStatus(t *testing.T) {
	txns := []*models.Transaction{
		{ID: "tx1", Status: "completed"},
		{ID: "tx2", Status: "pending"},
		{ID: "tx3", Status: "completed"},
		{ID: "tx4", Status: "failed"},
		{ID: "tx5", Status: "completed"},
	}

	tests := []struct {
		name      string
		status    string
		wantCount int
	}{
		{"completed transactions", "completed", 3},
		{"pending transactions", "pending", 1},
		{"failed transactions", "failed", 1},
		{"cancelled transactions", "cancelled", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered := filterByStatus(txns, tt.status)
			if len(filtered) != tt.wantCount {
				t.Errorf("Expected %d transactions, got %d", tt.wantCount, len(filtered))
			}
		})
	}
}

func filterByStatus(txns []*models.Transaction, status string) []*models.Transaction {
	var result []*models.Transaction
	for _, tx := range txns {
		if tx.Status == status {
			result = append(result, tx)
		}
	}
	return result
}

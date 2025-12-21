package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/CB-AccountStack/AccountStack/apps/api-transactions/internal/models"
	"github.com/sirupsen/logrus"
)

// Account represents a user's account (minimal structure needed for filtering)
type Account struct {
	ID     string `json:"id"`
	UserID string `json:"userId"`
}

// Repository provides data access for transactions
type Repository struct {
	transactions map[string]*models.Transaction
	accounts     map[string]*Account // accountID -> Account
	mu           sync.RWMutex
	logger       *logrus.Logger
}

// NewRepository creates a new repository and loads data from JSON files
func NewRepository(dataPath string, logger *logrus.Logger) (*Repository, error) {
	repo := &Repository{
		transactions: make(map[string]*models.Transaction),
		accounts:     make(map[string]*Account),
		logger:       logger,
	}

	// Load accounts first (needed for user isolation)
	if err := repo.loadAccounts(filepath.Join(dataPath, "accounts.json")); err != nil {
		return nil, fmt.Errorf("failed to load accounts: %w", err)
	}

	// Load transactions
	if err := repo.loadTransactions(filepath.Join(dataPath, "transactions.json")); err != nil {
		return nil, fmt.Errorf("failed to load transactions: %w", err)
	}

	logger.Infof("Loaded %d accounts and %d transactions from %s", len(repo.accounts), len(repo.transactions), dataPath)

	return repo, nil
}

// loadAccounts loads accounts from a JSON file
func (r *Repository) loadAccounts(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var accounts []*Account
	if err := json.Unmarshal(data, &accounts); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, acc := range accounts {
		r.accounts[acc.ID] = acc
	}

	return nil
}

// loadTransactions loads transactions from a JSON file
func (r *Repository) loadTransactions(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var transactions []*models.Transaction
	if err := json.Unmarshal(data, &transactions); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, txn := range transactions {
		r.transactions[txn.ID] = txn
	}

	return nil
}

// GetTransactionByID retrieves a transaction by ID
func (r *Repository) GetTransactionByID(txnID string) (*models.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	txn, exists := r.transactions[txnID]
	if !exists {
		return nil, fmt.Errorf("transaction not found")
	}

	return txn, nil
}

// GetAllTransactions returns all transactions
func (r *Repository) GetAllTransactions() []*models.Transaction {
	r.mu.RLock()
	defer r.mu.RUnlock()

	transactions := make([]*models.Transaction, 0, len(r.transactions))
	for _, txn := range r.transactions {
		transactions = append(transactions, txn)
	}

	return transactions
}

// GetAccountIDsByUserID retrieves all account IDs for a given user
func (r *Repository) GetAccountIDsByUserID(userID string) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var accountIDs []string
	for _, acc := range r.accounts {
		if acc.UserID == userID {
			accountIDs = append(accountIDs, acc.ID)
		}
	}

	return accountIDs
}

// GetTransactionsByFilter retrieves transactions matching the given filters
func (r *Repository) GetTransactionsByFilter(filters *models.TransactionFilters) []*models.Transaction {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []*models.Transaction
	for _, txn := range r.transactions {
		if txn.Matches(filters) {
			filtered = append(filtered, txn)
		}
	}

	return filtered
}

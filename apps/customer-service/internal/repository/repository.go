package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/CB-AccountStack/AccountStack/apps/api-accounts/internal/models"
	"github.com/sirupsen/logrus"
)

// Repository provides data access for users and accounts
type Repository struct {
	users    map[string]*models.User
	accounts map[string]*models.Account
	mu       sync.RWMutex
	logger   *logrus.Logger
}

// NewRepository creates a new repository and loads data from JSON files
func NewRepository(dataPath string, logger *logrus.Logger) (*Repository, error) {
	repo := &Repository{
		users:    make(map[string]*models.User),
		accounts: make(map[string]*models.Account),
		logger:   logger,
	}

	// Load users
	if err := repo.loadUsers(filepath.Join(dataPath, "users.json")); err != nil {
		return nil, fmt.Errorf("failed to load users: %w", err)
	}

	// Load accounts
	if err := repo.loadAccounts(filepath.Join(dataPath, "accounts.json")); err != nil {
		return nil, fmt.Errorf("failed to load accounts: %w", err)
	}

	logger.Infof("Loaded %d users and %d accounts from %s", len(repo.users), len(repo.accounts), dataPath)

	return repo, nil
}

// loadUsers loads users from a JSON file
func (r *Repository) loadUsers(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var users []*models.User
	if err := json.Unmarshal(data, &users); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, user := range users {
		r.users[user.ID] = user
	}

	return nil
}

// loadAccounts loads accounts from a JSON file
func (r *Repository) loadAccounts(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var accounts []*models.Account
	if err := json.Unmarshal(data, &accounts); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, account := range accounts {
		r.accounts[account.ID] = account
	}

	return nil
}

// GetUserByID retrieves a user by ID
func (r *Repository) GetUserByID(userID string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[userID]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email address
func (r *Repository) GetUserByEmail(email string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, fmt.Errorf("user not found")
}

// GetAccountByID retrieves an account by ID
func (r *Repository) GetAccountByID(accountID string) (*models.Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	account, exists := r.accounts[accountID]
	if !exists {
		return nil, fmt.Errorf("account not found")
	}

	return account, nil
}

// GetAccountsByUserID retrieves all accounts for a specific user
func (r *Repository) GetAccountsByUserID(userID string) ([]*models.Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var userAccounts []*models.Account
	for _, account := range r.accounts {
		if account.UserID == userID {
			userAccounts = append(userAccounts, account)
		}
	}

	return userAccounts, nil
}

// GetAllUsers returns all users (for testing/admin purposes)
func (r *Repository) GetAllUsers() []*models.User {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*models.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	return users
}

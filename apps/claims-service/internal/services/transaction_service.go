package services

import (
	"sort"
	"time"

	"github.com/CB-AccountStack/AccountStack/apps/api-transactions/internal/features"
	"github.com/CB-AccountStack/AccountStack/apps/api-transactions/internal/models"
	"github.com/CB-AccountStack/AccountStack/apps/api-transactions/internal/repository"
	"github.com/sirupsen/logrus"
)

// TransactionService handles business logic for transactions
type TransactionService struct {
	repo   *repository.Repository
	flags  *features.Flags
	logger *logrus.Logger
}

// NewTransactionService creates a new transaction service
func NewTransactionService(repo *repository.Repository, flags *features.Flags, logger *logrus.Logger) *TransactionService {
	return &TransactionService{
		repo:   repo,
		flags:  flags,
		logger: logger,
	}
}

// GetTransactionByID retrieves a transaction by ID
func (s *TransactionService) GetTransactionByID(txnID string) (*models.Transaction, error) {
	return s.repo.GetTransactionByID(txnID)
}

// GetTransactions retrieves transactions with optional filters
// Advanced filters (date range, amount range, category) are only applied if feature flag is enabled
// IMPORTANT: This enforces user isolation - only returns transactions for the specified user's accounts
func (s *TransactionService) GetTransactions(userID string, filters *models.TransactionFilters) ([]*models.Transaction, error) {
	// Get all account IDs for this user (enforces user isolation)
	userAccountIDs := s.repo.GetAccountIDsByUserID(userID)
	if len(userAccountIDs) == 0 {
		s.logger.WithField("userId", userID).Warn("No accounts found for user")
		return []*models.Transaction{}, nil
	}

	s.logger.WithFields(logrus.Fields{
		"userId":     userID,
		"accountIds": userAccountIDs,
	}).Debug("Filtering transactions by user accounts")

	// Check if advanced filters are enabled
	advancedFiltersEnabled := s.flags.IsAdvancedFiltersEnabled()

	// Get transactions for all user's accounts
	var allTransactions []*models.Transaction

	for _, accountID := range userAccountIDs {
		// Create filter for this account
		effectiveFilters := &models.TransactionFilters{
			AccountID: accountID,
		}

		// Only apply advanced filters if the feature flag is enabled
		if advancedFiltersEnabled {
			effectiveFilters.StartDate = filters.StartDate
			effectiveFilters.EndDate = filters.EndDate
			effectiveFilters.Category = filters.Category
			effectiveFilters.MinAmount = filters.MinAmount
			effectiveFilters.MaxAmount = filters.MaxAmount
		} else {
			// Log that advanced filters are ignored
			if filters.StartDate != nil || filters.EndDate != nil || filters.Category != "" ||
				filters.MinAmount != nil || filters.MaxAmount != nil {
				s.logger.Info("Advanced filters requested but feature flag is disabled, only accountId filter will be applied")
			}
		}

		// Get transactions for this account
		accountTransactions := s.repo.GetTransactionsByFilter(effectiveFilters)
		allTransactions = append(allTransactions, accountTransactions...)
	}

	// Sort by date descending (most recent first)
	sort.Slice(allTransactions, func(i, j int) bool {
		return allTransactions[i].Date.After(allTransactions[j].Date)
	})

	s.logger.WithFields(logrus.Fields{
		"userId": userID,
		"count":  len(allTransactions),
	}).Info("Retrieved transactions for user")

	return allTransactions, nil
}

// ParseDateParam parses a date string in ISO 8601 format
func ParseDateParam(dateStr string) (*time.Time, error) {
	if dateStr == "" {
		return nil, nil
	}

	// Try parsing as RFC3339 (ISO 8601)
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		// Try parsing as date only (YYYY-MM-DD)
		t, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, err
		}
	}

	return &t, nil
}

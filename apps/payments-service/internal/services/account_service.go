package services

import (
	"fmt"

	"github.com/CB-AccountStack/AccountStack/apps/api-accounts/internal/features"
	"github.com/CB-AccountStack/AccountStack/apps/api-accounts/internal/models"
	"github.com/CB-AccountStack/AccountStack/apps/api-accounts/internal/repository"
	"github.com/sirupsen/logrus"
)

// AccountService handles business logic for accounts
type AccountService struct {
	repo   *repository.Repository
	flags  *features.Flags
	logger *logrus.Logger
}

// NewAccountService creates a new account service
func NewAccountService(repo *repository.Repository, flags *features.Flags, logger *logrus.Logger) *AccountService {
	return &AccountService{
		repo:   repo,
		flags:  flags,
		logger: logger,
	}
}

// GetAccountByID retrieves an account by ID and applies masking if needed
func (s *AccountService) GetAccountByID(accountID string, userID string) (*models.AccountResponse, error) {
	account, err := s.repo.GetAccountByID(accountID)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"accountId": accountID,
			"userId":    userID,
		}).Warn("Account not found")
		return nil, err
	}

	// Verify the account belongs to the requesting user
	if account.UserID != userID {
		s.logger.WithFields(logrus.Fields{
			"accountId": accountID,
			"userId":    userID,
			"ownerId":   account.UserID,
		}).Warn("Unauthorized access attempt")
		return nil, fmt.Errorf("unauthorized")
	}

	// Get user to determine currency based on country
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		s.logger.WithField("userId", userID).Warn("User not found")
		return nil, err
	}

	// Apply masking and currency based on feature flags and user context
	maskAmounts := s.flags.ShouldMaskAmounts()
	currency := s.flags.GetCurrencyForUser(user.Country)
	s.logger.WithFields(logrus.Fields{
		"accountId":   accountID,
		"userId":      userID,
		"userCountry": user.Country,
		"maskAmounts": maskAmounts,
		"currency":    currency,
	}).Debug("Retrieving account")

	response := account.ToResponse(maskAmounts, currency)
	return &response, nil
}

// GetAccountsByUserID retrieves all accounts for a user with optional masking
func (s *AccountService) GetAccountsByUserID(userID string) ([]models.AccountResponse, error) {
	accounts, err := s.repo.GetAccountsByUserID(userID)
	if err != nil {
		s.logger.WithField("userId", userID).Error("Failed to retrieve accounts")
		return nil, err
	}

	// Get user to determine currency based on country
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		s.logger.WithField("userId", userID).Warn("User not found")
		return nil, err
	}

	// Apply masking and currency based on feature flags and user context
	maskAmounts := s.flags.ShouldMaskAmounts()
	currency := s.flags.GetCurrencyForUser(user.Country)
	s.logger.WithFields(logrus.Fields{
		"userId":      userID,
		"userCountry": user.Country,
		"count":       len(accounts),
		"maskAmounts": maskAmounts,
		"currency":    currency,
	}).Debug("Retrieving accounts")

	responses := make([]models.AccountResponse, len(accounts))
	for i, account := range accounts {
		responses[i] = account.ToResponse(maskAmounts, currency)
	}

	return responses, nil
}

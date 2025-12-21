package services

import (
	"github.com/CB-AccountStack/AccountStack/apps/api-accounts/internal/models"
	"github.com/CB-AccountStack/AccountStack/apps/api-accounts/internal/repository"
	"github.com/sirupsen/logrus"
)

// UserService handles business logic for users
type UserService struct {
	repo   *repository.Repository
	logger *logrus.Logger
}

// NewUserService creates a new user service
func NewUserService(repo *repository.Repository, logger *logrus.Logger) *UserService {
	return &UserService{
		repo:   repo,
		logger: logger,
	}
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(userID string) (*models.User, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		s.logger.WithField("userId", userID).Warn("User not found")
		return nil, err
	}

	s.logger.WithField("userId", userID).Debug("User retrieved")
	return user, nil
}

package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/CB-AccountStack/AccountStack/apps/api-accounts/internal/middleware"
	"github.com/CB-AccountStack/AccountStack/apps/api-accounts/internal/services"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// AccountHandler handles account-related requests
type AccountHandler struct {
	accountService *services.AccountService
	logger         *logrus.Logger
}

// NewAccountHandler creates a new account handler
func NewAccountHandler(accountService *services.AccountService, logger *logrus.Logger) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
		logger:         logger,
	}
}

// GetAccounts handles GET /accounts - returns all accounts for current user
func (h *AccountHandler) GetAccounts(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	accounts, err := h.accountService.GetAccountsByUserID(userID)
	if err != nil {
		h.logger.WithError(err).WithField("userId", userID).Error("Failed to get accounts")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to retrieve accounts",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accounts)
}

// GetAccountByID handles GET /accounts/{id} - returns a specific account
func (h *AccountHandler) GetAccountByID(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	vars := mux.Vars(r)
	accountID := vars["id"]

	if accountID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "bad_request",
			Message: "Account ID is required",
		})
		return
	}

	account, err := h.accountService.GetAccountByID(accountID, userID)
	if err != nil {
		if err.Error() == "unauthorized" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(ErrorResponse{
				Error:   "forbidden",
				Message: "You do not have access to this account",
			})
			return
		}

		h.logger.WithError(err).WithFields(logrus.Fields{
			"userId":    userID,
			"accountId": accountID,
		}).Error("Failed to get account")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "not_found",
			Message: "Account not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(account)
}

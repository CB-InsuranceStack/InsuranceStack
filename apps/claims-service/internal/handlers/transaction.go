package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/CB-AccountStack/AccountStack/apps/api-transactions/internal/middleware"
	"github.com/CB-AccountStack/AccountStack/apps/api-transactions/internal/models"
	"github.com/CB-AccountStack/AccountStack/apps/api-transactions/internal/services"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// TransactionHandler handles transaction-related HTTP requests
type TransactionHandler struct {
	service *services.TransactionService
	logger  *logrus.Logger
}

// NewTransactionHandler creates a new transaction handler
func NewTransactionHandler(service *services.TransactionService, logger *logrus.Logger) *TransactionHandler {
	return &TransactionHandler{
		service: service,
		logger:  logger,
	}
}

// GetTransactions handles GET /transactions
// Supports query parameters:
// - accountId: filter by account ID (always allowed)
// - startDate: filter by start date (ISO 8601 format, requires advancedFilters flag)
// - endDate: filter by end date (ISO 8601 format, requires advancedFilters flag)
// - category: filter by category (requires advancedFilters flag)
// - minAmount: filter by minimum amount (requires advancedFilters flag)
// - maxAmount: filter by maximum amount (requires advancedFilters flag)
//
// IMPORTANT: Enforces user isolation - only returns transactions for the authenticated user's accounts
func (h *TransactionHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserID(r)
	if userID == "" {
		h.logger.Warn("User ID not found in context")
		h.respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse query parameters
	query := r.URL.Query()

	filters := &models.TransactionFilters{
		AccountID: query.Get("accountId"),
	}

	// Parse date filters
	if startDateStr := query.Get("startDate"); startDateStr != "" {
		startDate, err := services.ParseDateParam(startDateStr)
		if err != nil {
			h.logger.WithError(err).Warn("Invalid startDate parameter")
			h.respondError(w, http.StatusBadRequest, "Invalid startDate format. Use ISO 8601 (YYYY-MM-DD or RFC3339)")
			return
		}
		filters.StartDate = startDate
	}

	if endDateStr := query.Get("endDate"); endDateStr != "" {
		endDate, err := services.ParseDateParam(endDateStr)
		if err != nil {
			h.logger.WithError(err).Warn("Invalid endDate parameter")
			h.respondError(w, http.StatusBadRequest, "Invalid endDate format. Use ISO 8601 (YYYY-MM-DD or RFC3339)")
			return
		}
		filters.EndDate = endDate
	}

	// Parse category filter
	if category := query.Get("category"); category != "" {
		filters.Category = category
	}

	// Parse amount filters
	if minAmountStr := query.Get("minAmount"); minAmountStr != "" {
		minAmount, err := strconv.ParseFloat(minAmountStr, 64)
		if err != nil {
			h.logger.WithError(err).Warn("Invalid minAmount parameter")
			h.respondError(w, http.StatusBadRequest, "Invalid minAmount format. Must be a number")
			return
		}
		filters.MinAmount = &minAmount
	}

	if maxAmountStr := query.Get("maxAmount"); maxAmountStr != "" {
		maxAmount, err := strconv.ParseFloat(maxAmountStr, 64)
		if err != nil {
			h.logger.WithError(err).Warn("Invalid maxAmount parameter")
			h.respondError(w, http.StatusBadRequest, "Invalid maxAmount format. Must be a number")
			return
		}
		filters.MaxAmount = &maxAmount
	}

	// Get transactions with filters (user isolation enforced in service layer)
	transactions, err := h.service.GetTransactions(userID, filters)
	if err != nil {
		h.logger.WithError(err).Error("Failed to retrieve transactions")
		h.respondError(w, http.StatusInternalServerError, "Failed to retrieve transactions")
		return
	}

	h.respondJSON(w, http.StatusOK, transactions)
}

// GetTransactionByID handles GET /transactions/{id}
func (h *TransactionHandler) GetTransactionByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	txnID := vars["id"]

	if txnID == "" {
		h.respondError(w, http.StatusBadRequest, "Transaction ID is required")
		return
	}

	transaction, err := h.service.GetTransactionByID(txnID)
	if err != nil {
		h.logger.WithError(err).WithField("txnId", txnID).Warn("Transaction not found")
		h.respondError(w, http.StatusNotFound, "Transaction not found")
		return
	}

	h.respondJSON(w, http.StatusOK, transaction)
}

// respondJSON sends a JSON response
func (h *TransactionHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.WithError(err).Error("Failed to encode response")
	}
}

// respondError sends an error response
func (h *TransactionHandler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, map[string]string{"error": message})
}

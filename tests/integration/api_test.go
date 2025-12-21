package integration

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	accountsAPIBaseURL     = "http://localhost:8001"
	transactionsAPIBaseURL = "http://localhost:8002"
	insightsAPIBaseURL     = "http://localhost:8003"
)

// TestAccountsAPIHealth tests the health endpoint of accounts API
func TestAccountsAPIHealth(t *testing.T) {
	resp, err := http.Get(accountsAPIBaseURL + "/healthz")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var health map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&health)
	require.NoError(t, err)

	assert.Equal(t, "ok", health["status"])
	assert.NotEmpty(t, health["timestamp"])
}

// TestAccountsAPIGetAccounts tests listing accounts
func TestAccountsAPIGetAccounts(t *testing.T) {
	req, err := http.NewRequest("GET", accountsAPIBaseURL+"/accounts", nil)
	require.NoError(t, err)
	req.Header.Set("X-User-ID", "user-001")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var accounts []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&accounts)
	require.NoError(t, err)

	assert.Greater(t, len(accounts), 0, "Should return at least one account")

	// Verify account structure
	if len(accounts) > 0 {
		account := accounts[0]
		assert.NotEmpty(t, account["id"])
		assert.NotEmpty(t, account["accountNumber"])
		assert.NotEmpty(t, account["accountType"])
		assert.NotNil(t, account["balance"])
	}
}

// TestAccountsAPIGetSpecificAccount tests getting a specific account
func TestAccountsAPIGetSpecificAccount(t *testing.T) {
	req, err := http.NewRequest("GET", accountsAPIBaseURL+"/accounts/acc-001", nil)
	require.NoError(t, err)
	req.Header.Set("X-User-ID", "user-001")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var account map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&account)
	require.NoError(t, err)

	assert.Equal(t, "acc-001", account["id"])
	assert.NotEmpty(t, account["accountNumber"])
}

// TestAccountsAPIUnauthorizedAccess tests accessing another user's account
func TestAccountsAPIUnauthorizedAccess(t *testing.T) {
	req, err := http.NewRequest("GET", accountsAPIBaseURL+"/accounts/acc-004", nil)
	require.NoError(t, err)
	req.Header.Set("X-User-ID", "user-001") // user-001 trying to access user-002's account

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}

// TestTransactionsAPIHealth tests the health endpoint of transactions API
func TestTransactionsAPIHealth(t *testing.T) {
	resp, err := http.Get(transactionsAPIBaseURL + "/healthz")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var health map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&health)
	require.NoError(t, err)

	assert.Equal(t, "ok", health["status"])
}

// TestTransactionsAPIGetTransactions tests listing transactions
func TestTransactionsAPIGetTransactions(t *testing.T) {
	req, err := http.NewRequest("GET", transactionsAPIBaseURL+"/transactions?accountId=acc-001", nil)
	require.NoError(t, err)
	req.Header.Set("X-User-ID", "user-001")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var transactions []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&transactions)
	require.NoError(t, err)

	assert.Greater(t, len(transactions), 0, "Should return at least one transaction")

	// Verify all transactions belong to the requested account
	for _, txn := range transactions {
		assert.Equal(t, "acc-001", txn["accountId"])
	}
}

// TestInsightsAPIHealth tests the health endpoint of insights API
func TestInsightsAPIHealth(t *testing.T) {
	resp, err := http.Get(insightsAPIBaseURL + "/healthz")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var health map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&health)
	require.NoError(t, err)

	assert.Equal(t, "ok", health["status"])
}

// TestInsightsAPIGetInsights tests listing insights
func TestInsightsAPIGetInsights(t *testing.T) {
	req, err := http.NewRequest("GET", insightsAPIBaseURL+"/insights", nil)
	require.NoError(t, err)
	req.Header.Set("X-User-ID", "user-001")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var insights []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&insights)
	require.NoError(t, err)

	assert.Greater(t, len(insights), 0, "Should return at least one insight")

	// Verify insight structure
	if len(insights) > 0 {
		insight := insights[0]
		assert.NotEmpty(t, insight["id"])
		assert.NotEmpty(t, insight["type"])
		assert.NotEmpty(t, insight["title"])
		assert.NotEmpty(t, insight["description"])
	}
}

// TestCORSHeaders tests that all APIs return proper CORS headers
func TestCORSHeaders(t *testing.T) {
	endpoints := []string{
		accountsAPIBaseURL + "/healthz",
		transactionsAPIBaseURL + "/healthz",
		insightsAPIBaseURL + "/healthz",
	}

	for _, endpoint := range endpoints {
		resp, err := http.Get(endpoint)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.NotEmpty(t, resp.Header.Get("Access-Control-Allow-Origin"), "CORS header missing for "+endpoint)
	}
}

// TestEndToEndFlow tests a complete user flow across all APIs
func TestEndToEndFlow(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}

	// 1. Get user accounts
	req, err := http.NewRequest("GET", accountsAPIBaseURL+"/accounts", nil)
	require.NoError(t, err)
	req.Header.Set("X-User-ID", "user-001")

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	var accounts []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&accounts)
	require.NoError(t, err)
	require.Greater(t, len(accounts), 0)

	accountID := accounts[0]["id"].(string)

	// 2. Get transactions for first account
	req, err = http.NewRequest("GET", transactionsAPIBaseURL+"/transactions?accountId="+accountID, nil)
	require.NoError(t, err)
	req.Header.Set("X-User-ID", "user-001")

	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	var transactions []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&transactions)
	require.NoError(t, err)
	assert.Greater(t, len(transactions), 0)

	// 3. Get insights for user
	req, err = http.NewRequest("GET", insightsAPIBaseURL+"/insights", nil)
	require.NoError(t, err)
	req.Header.Set("X-User-ID", "user-001")

	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	var insights []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&insights)
	require.NoError(t, err)
	assert.Greater(t, len(insights), 0)
}

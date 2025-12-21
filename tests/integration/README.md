# Integration Tests

Integration tests for AccountStack APIs. These tests verify that all services work correctly together.

## Prerequisites

- Go 1.21+
- All API services running (accounts, transactions, insights)
- Docker Compose (recommended for running services)

## Running Services

Before running tests, start all services:

```bash
# From repository root
docker compose up
```

Or run services individually:

```bash
# Terminal 1 - Accounts API
cd apps/api-accounts
DATA_PATH=../../data/seed go run cmd/server/main.go

# Terminal 2 - Transactions API
cd apps/api-transactions
DATA_PATH=../../data/seed go run cmd/server/main.go

# Terminal 3 - Insights API
cd apps/api-insights
DATA_PATH=../../data/seed go run cmd/server/main.go
```

## Running Tests

```bash
# From this directory
go test -v ./...

# With coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Test Suites

### Health Checks
- Verifies all services are running and healthy
- Tests: `TestAccountsAPIHealth`, `TestTransactionsAPIHealth`, `TestInsightsAPIHealth`

### Accounts API
- List accounts for a user
- Get specific account by ID
- Verify unauthorized access is blocked
- Tests: `TestAccountsAPIGetAccounts`, `TestAccountsAPIGetSpecificAccount`, `TestAccountsAPIUnauthorizedAccess`

### Transactions API
- List transactions with filters
- Verify transaction ownership
- Test advanced filtering (when feature flag enabled)
- Tests: `TestTransactionsAPIGetTransactions`

### Insights API
- List insights for a user
- Get specific insight
- Test alerts endpoint
- Tests: `TestInsightsAPIGetInsights`

### Cross-Service Tests
- CORS headers verification
- End-to-end user flow across all APIs
- Tests: `TestCORSHeaders`, `TestEndToEndFlow`

## Expected Results

All tests should pass when services are running correctly:

```
=== RUN   TestAccountsAPIHealth
--- PASS: TestAccountsAPIHealth (0.01s)
=== RUN   TestAccountsAPIGetAccounts
--- PASS: TestAccountsAPIGetAccounts (0.01s)
=== RUN   TestTransactionsAPIHealth
--- PASS: TestTransactionsAPIHealth (0.01s)
...
PASS
ok      github.com/CB-AccountStack/AccountStack/tests/integration    0.234s
```

## Troubleshooting

### Services not running
```
Error: Get "http://localhost:8001/healthz": dial tcp [::1]:8001: connect: connection refused
```
**Solution**: Start the services with `docker compose up`

### Wrong user ID
```
Error: Status code 403 (expected 200)
```
**Solution**: Verify X-User-ID header matches test user IDs (user-001, user-002)

### Data not loaded
```
Error: Should return at least one account
```
**Solution**: Verify DATA_PATH environment variable points to `../../data/seed`

## CI/CD Integration

These tests run automatically in CloudBees workflows:

```yaml
- name: Run integration tests
  run: |
    docker compose up -d
    sleep 5
    cd tests/integration
    go test -v ./...
```

## Adding New Tests

1. Create test function with `Test` prefix
2. Use `testify/assert` and `testify/require` for assertions
3. Set proper X-User-ID headers
4. Clean up resources (defer resp.Body.Close())
5. Test both success and error cases

Example:

```go
func TestMyNewEndpoint(t *testing.T) {
    req, err := http.NewRequest("GET", accountsAPIBaseURL+"/my-endpoint", nil)
    require.NoError(t, err)
    req.Header.Set("X-User-ID", "user-001")

    client := &http.Client{Timeout: 5 * time.Second}
    resp, err := client.Do(req)
    require.NoError(t, err)
    defer resp.Body.Close()

    assert.Equal(t, http.StatusOK, resp.StatusCode)
}
```

# API Examples

This document provides practical examples for testing the Transactions API.

## Prerequisites

Make sure the service is running:
```bash
make run-dev
```

Or with Docker:
```bash
make docker-build
make docker-run
```

## Health Check

```bash
# Simple health check
curl http://localhost:8002/healthz

# Expected response:
# {"status":"healthy","service":"api-transactions"}
```

## Basic Transaction Queries

### Get All Transactions for an Account

```bash
curl "http://localhost:8002/transactions?accountId=acc-001"
```

### Get Single Transaction by ID

```bash
curl http://localhost:8002/transactions/txn-001

# Pretty print with jq
curl -s http://localhost:8002/transactions/txn-001 | jq
```

### Get Transactions with Authentication Header

```bash
curl -H "X-User-ID: user-001" \
  "http://localhost:8002/transactions?accountId=acc-001"
```

## Advanced Filtering Examples

**Note:** These require the `api.advancedFilters` feature flag to be enabled in CloudBees.

### Filter by Date Range

```bash
# Get transactions between Dec 1-10, 2024
curl "http://localhost:8002/transactions?accountId=acc-001&startDate=2024-12-01&endDate=2024-12-10"

# Get transactions from Dec 5 onwards
curl "http://localhost:8002/transactions?accountId=acc-001&startDate=2024-12-05"

# Get transactions up to Dec 10
curl "http://localhost:8002/transactions?accountId=acc-001&endDate=2024-12-10"

# Using RFC3339 format
curl "http://localhost:8002/transactions?accountId=acc-001&startDate=2024-12-01T00:00:00Z&endDate=2024-12-13T23:59:59Z"
```

### Filter by Category

```bash
# Get all food and dining transactions
curl "http://localhost:8002/transactions?accountId=acc-001&category=food_dining"

# Get all shopping transactions
curl "http://localhost:8002/transactions?category=shopping"

# Get all income transactions
curl "http://localhost:8002/transactions?category=income"

# Get all utility bills
curl "http://localhost:8002/transactions?accountId=acc-001&category=utilities"
```

### Filter by Amount Range

```bash
# Get transactions between -100 and -10 (expenses)
curl "http://localhost:8002/transactions?accountId=acc-001&minAmount=-100&maxAmount=-10"

# Get large expenses (more than $100)
curl "http://localhost:8002/transactions?accountId=acc-001&maxAmount=-100"

# Get small expenses (less than $20)
curl "http://localhost:8002/transactions?accountId=acc-001&minAmount=-20&maxAmount=0"

# Get all income (positive amounts)
curl "http://localhost:8002/transactions?accountId=acc-001&minAmount=0"
```

### Complex Queries

```bash
# Shopping transactions over $50 in December 2024
curl "http://localhost:8002/transactions?accountId=acc-001&category=shopping&maxAmount=-50&startDate=2024-12-01&endDate=2024-12-31"

# Food expenses under $50 in the last week
curl "http://localhost:8002/transactions?accountId=acc-001&category=food_dining&minAmount=-50&maxAmount=0&startDate=2024-12-06"

# All entertainment expenses
curl "http://localhost:8002/transactions?accountId=acc-001&category=entertainment&maxAmount=0"

# Large transactions (income or expenses over $500)
curl "http://localhost:8002/transactions?accountId=acc-001&minAmount=500"
curl "http://localhost:8002/transactions?accountId=acc-001&maxAmount=-500"
```

## Testing Feature Flag Behavior

### When advancedFilters is DISABLED (default)

```bash
# This will only use accountId filter, others will be ignored
curl "http://localhost:8002/transactions?accountId=acc-001&category=shopping&minAmount=-100"

# Check logs - you should see:
# "Advanced filters requested but feature flag is disabled, only accountId filter will be applied"
```

### When advancedFilters is ENABLED

Enable the flag in CloudBees dashboard, then:

```bash
# Now all filters work
curl "http://localhost:8002/transactions?accountId=acc-001&category=shopping&minAmount=-100"

# Check logs - you should see:
# "Using advanced filters"
```

## Account-Specific Queries

### Account acc-001 (Checking Account)

```bash
# All transactions
curl "http://localhost:8002/transactions?accountId=acc-001"

# Recent shopping
curl "http://localhost:8002/transactions?accountId=acc-001&category=shopping"

# Food expenses
curl "http://localhost:8002/transactions?accountId=acc-001&category=food_dining"
```

### Account acc-002 (Savings Account)

```bash
# All transactions
curl "http://localhost:8002/transactions?accountId=acc-002"

# Income only
curl "http://localhost:8002/transactions?accountId=acc-002&category=income"
```

### Account acc-003 (Credit Card)

```bash
# All transactions
curl "http://localhost:8002/transactions?accountId=acc-003"

# Pending transactions
curl "http://localhost:8002/transactions?accountId=acc-003" | jq '[.[] | select(.status=="pending")]'
```

### Account acc-004 (Business Account)

```bash
# All transactions
curl "http://localhost:8002/transactions?accountId=acc-004"

# Business expenses
curl "http://localhost:8002/transactions?accountId=acc-004&category=business"
```

## Error Cases

### Invalid Transaction ID

```bash
curl http://localhost:8002/transactions/invalid-id

# Expected: 404 Not Found
# {"error":"Transaction not found"}
```

### Invalid Date Format

```bash
curl "http://localhost:8002/transactions?accountId=acc-001&startDate=invalid"

# Expected: 400 Bad Request
# {"error":"Invalid startDate format. Use ISO 8601 (YYYY-MM-DD or RFC3339)"}
```

### Invalid Amount Format

```bash
curl "http://localhost:8002/transactions?accountId=acc-001&minAmount=not-a-number"

# Expected: 400 Bad Request
# {"error":"Invalid minAmount format. Must be a number"}
```

## Using with jq for Pretty Output

```bash
# Pretty print JSON
curl -s http://localhost:8002/transactions/txn-001 | jq

# Get only transaction amounts
curl -s "http://localhost:8002/transactions?accountId=acc-001" | jq '[.[] | .amount]'

# Get transaction descriptions
curl -s "http://localhost:8002/transactions?accountId=acc-001" | jq '[.[] | .description]'

# Filter in jq (client-side)
curl -s "http://localhost:8002/transactions?accountId=acc-001" | jq '[.[] | select(.category=="shopping")]'

# Calculate total
curl -s "http://localhost:8002/transactions?accountId=acc-001" | jq '[.[] | .amount] | add'

# Group by category
curl -s "http://localhost:8002/transactions?accountId=acc-001" | jq 'group_by(.category) | map({category: .[0].category, count: length})'
```

## Performance Testing

```bash
# Simple load test with Apache Bench
ab -n 1000 -c 10 "http://localhost:8002/transactions?accountId=acc-001"

# With curl in a loop
for i in {1..100}; do
  curl -s "http://localhost:8002/transactions?accountId=acc-001" > /dev/null
  echo "Request $i completed"
done
```

## Debugging

### Check Logs

```bash
# Run with debug logging
LOG_LEVEL=debug make run-dev

# Or in Docker
docker run -p 8002:8002 \
  -e LOG_LEVEL=debug \
  -e DATA_PATH=/data/seed \
  -v $(pwd)/../../data/seed:/data/seed \
  accountstack/api-transactions:latest
```

### Verbose curl

```bash
# Show request and response headers
curl -v http://localhost:8002/transactions/txn-001

# Show only response headers
curl -I http://localhost:8002/healthz
```

## Integration with Other Services

### Combining with Accounts API

```bash
# Get accounts first
ACCOUNTS=$(curl -s http://localhost:8001/accounts)

# Get transactions for each account
echo $ACCOUNTS | jq -r '.[].id' | while read account_id; do
  echo "Transactions for $account_id:"
  curl -s "http://localhost:8002/transactions?accountId=$account_id" | jq length
done
```

## Postman Collection

You can also import this API into Postman:

1. Create a new collection named "Transactions API"
2. Add these endpoints:
   - GET `{{baseUrl}}/healthz`
   - GET `{{baseUrl}}/transactions?accountId=acc-001`
   - GET `{{baseUrl}}/transactions/txn-001`
   - GET `{{baseUrl}}/transactions?accountId=acc-001&category=shopping`
3. Set environment variable: `baseUrl = http://localhost:8002`

## Tips

1. **Always use accountId filter** - Most queries should include an accountId for better performance
2. **Date formats** - Use ISO 8601 format (YYYY-MM-DD or full RFC3339)
3. **Amount ranges** - Negative amounts are debits/expenses, positive are credits/income
4. **Feature flags** - Check CloudBees dashboard to verify flag states
5. **Logging** - Set LOG_LEVEL=debug to see detailed filter application logs

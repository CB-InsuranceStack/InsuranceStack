# Quick Start Guide

Get the Transactions API up and running in 5 minutes.

## Option 1: Run with Go (Recommended for Development)

### Step 1: Install Dependencies
```bash
cd /Users/brown/git_orgs/CB-AccountStack/AccountStack/apps/api-transactions
go mod download
```

### Step 2: Set Environment Variables
```bash
export DATA_PATH="../../data/seed"
export LOG_LEVEL="debug"
export CLOUDBEES_FM_API_KEY="your-api-key-or-use-dev-mode"
```

### Step 3: Run the Service
```bash
make run-dev
```

### Step 4: Test It
```bash
# Health check
curl http://localhost:8002/healthz

# Get transactions
curl "http://localhost:8002/transactions?accountId=acc-001"

# Get specific transaction
curl http://localhost:8002/transactions/txn-001
```

## Option 2: Run with Docker

### Step 1: Build Image
```bash
cd /Users/brown/git_orgs/CB-AccountStack/AccountStack/apps/api-transactions
make docker-build
```

### Step 2: Run Container
```bash
make docker-run
```

Or manually:
```bash
docker run -p 8002:8002 \
  -e DATA_PATH=/data/seed \
  -v $(pwd)/../../data/seed:/data/seed \
  accountstack/api-transactions:latest
```

### Step 3: Test It
```bash
curl http://localhost:8002/healthz
```

## Testing Feature Flags

### Without CloudBees (Development Mode)

The service runs with default flag values:
- `api.advancedFilters` = false (only accountId filtering works)

```bash
# This will only filter by accountId
curl "http://localhost:8002/transactions?accountId=acc-001&category=shopping"
```

### With CloudBees Feature Management

1. **Sign up for CloudBees**: https://www.cloudbees.com/products/feature-management

2. **Create an Application** in the CloudBees dashboard

3. **Create Feature Flag**:
   - Name: `api.advancedFilters`
   - Type: Boolean
   - Default: false

4. **Get API Key** from CloudBees dashboard

5. **Run with API Key**:
```bash
export CLOUDBEES_FM_API_KEY="your-actual-api-key"
make run-dev
```

6. **Enable the Flag** in CloudBees dashboard

7. **Test Advanced Filters**:
```bash
# Now category filtering works!
curl "http://localhost:8002/transactions?accountId=acc-001&category=shopping"

# Date range filtering
curl "http://localhost:8002/transactions?accountId=acc-001&startDate=2024-12-01&endDate=2024-12-13"

# Amount filtering
curl "http://localhost:8002/transactions?accountId=acc-001&minAmount=-100&maxAmount=-10"
```

## Common Test Scenarios

### Scenario 1: Get All Transactions for an Account
```bash
curl "http://localhost:8002/transactions?accountId=acc-001" | jq
```

### Scenario 2: Find Shopping Transactions (Requires Feature Flag)
```bash
curl "http://localhost:8002/transactions?accountId=acc-001&category=shopping" | jq
```

### Scenario 3: Find Large Expenses (Requires Feature Flag)
```bash
curl "http://localhost:8002/transactions?accountId=acc-001&maxAmount=-100" | jq
```

### Scenario 4: Get Recent Transactions (Requires Feature Flag)
```bash
curl "http://localhost:8002/transactions?accountId=acc-001&startDate=2024-12-10" | jq
```

## Verify Installation

Run this command to verify everything is working:

```bash
#!/bin/bash

echo "Testing Transactions API..."
echo ""

# Test health endpoint
echo "1. Health Check:"
curl -s http://localhost:8002/healthz | jq
echo ""

# Test get all transactions
echo "2. Get Transactions for acc-001:"
curl -s "http://localhost:8002/transactions?accountId=acc-001" | jq '. | length'
echo " transactions found"
echo ""

# Test get specific transaction
echo "3. Get Transaction txn-001:"
curl -s http://localhost:8002/transactions/txn-001 | jq '.id, .description, .amount'
echo ""

echo "All tests passed!"
```

Save this as `test.sh`, make it executable (`chmod +x test.sh`), and run it.

## Next Steps

1. **Read the full README**: [README.md](README.md)
2. **Try API examples**: [API_EXAMPLES.md](API_EXAMPLES.md)
3. **Set up CloudBees**: Enable advanced filtering
4. **Integrate with frontend**: Use these endpoints in your UI
5. **Deploy**: Build Docker image and deploy to your environment

## Troubleshooting

### Port Already in Use
```bash
# Find what's using port 8002
lsof -i :8002

# Use a different port
PORT=8003 make run-dev
```

### Data Not Loading
```bash
# Check if data files exist
ls -la ../../data/seed/transactions.json

# Use absolute path
export DATA_PATH="/Users/brown/git_orgs/CB-AccountStack/AccountStack/data/seed"
make run-dev
```

### Feature Flags Not Working
```bash
# Run with debug logging to see flag values
LOG_LEVEL=debug make run-dev

# Check the logs for:
# "Feature flags initialized" with "advancedFilters": true/false
```

## Development Workflow

```bash
# 1. Make code changes
vim internal/handlers/transaction.go

# 2. Format code
make fmt

# 3. Run tests
make test

# 4. Run locally
make run-dev

# 5. Test changes
curl "http://localhost:8002/transactions?accountId=acc-001"

# 6. Build for production
make build

# 7. Build Docker image
make docker-build
```

## Support

- Full documentation: [README.md](README.md)
- API examples: [API_EXAMPLES.md](API_EXAMPLES.md)
- Issues: Contact the development team

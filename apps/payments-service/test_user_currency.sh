#!/bin/bash

cd /Users/brown/git_orgs/CB-AccountStack/AccountStack/apps/api-accounts

# Function to safely kill and wait for the API to stop
kill_api() {
    pkill -f accounts-api
    sleep 1
    for i in {1..5}; do
        if ! lsof -i :8001 > /dev/null 2>&1; then
            return 0
        fi
        sleep 1
    done
}

# Ensure clean start
kill_api

echo "=========================================="
echo "Testing User-Based Currency Targeting"
echo "=========================================="
echo ""

# Start API
DATA_PATH=../../data/seed ./bin/accounts-api > /tmp/api-user-currency.log 2>&1 &
sleep 3

echo "=== Testing Demo User (US) → Should get USD ==="
TOKEN=$(curl -s -X POST http://localhost:8001/login \
  -H "Content-Type: application/json" \
  -d '{"username":"demo@accountstack.com","password":"demo123"}' | python3 -c "import sys, json; print(json.load(sys.stdin)['token'])")
echo "User: demo@accountstack.com (Country: US)"
curl -s http://localhost:8001/accounts -H "Authorization: Bearer $TOKEN" | \
  python3 -c "import sys, json; data=json.load(sys.stdin); print(f'  Currency: {data[0][\"currency\"]}' if data else '  No accounts')"
echo ""

echo "=== Testing Sarah Chen (UK) → Should get GBP ==="
TOKEN=$(curl -s -X POST http://localhost:8001/login \
  -H "Content-Type: application/json" \
  -d '{"username":"sarah.chen@accountstack.com","password":"demo123"}' | python3 -c "import sys, json; print(json.load(sys.stdin)['token'])")
echo "User: sarah.chen@accountstack.com (Country: UK)"
curl -s http://localhost:8001/accounts -H "Authorization: Bearer $TOKEN" | \
  python3 -c "import sys, json; data=json.load(sys.stdin); print(f'  Currency: {data[0][\"currency\"]}' if data else '  No accounts')"
echo ""

echo "=== Testing François Dubois (FR) → Should get EUR ==="
TOKEN=$(curl -s -X POST http://localhost:8001/login \
  -H "Content-Type: application/json" \
  -d '{"username":"francois.dubois@accountstack.com","password":"demo123"}' | python3 -c "import sys, json; print(json.load(sys.stdin)['token'])")
echo "User: francois.dubois@accountstack.com (Country: FR)"
curl -s http://localhost:8001/accounts -H "Authorization: Bearer $TOKEN" | \
  python3 -c "import sys, json; data=json.load(sys.stdin); print(f'  Currency: {data[0][\"currency\"]}' if data else '  No accounts')"
echo ""

kill_api

echo "=========================================="
echo "✓ User-based currency targeting test complete!"
echo "=========================================="

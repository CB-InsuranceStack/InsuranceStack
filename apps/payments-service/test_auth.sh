#!/bin/bash

cd /Users/brown/git_orgs/CB-AccountStack/AccountStack/apps/api-accounts

# Start API in background
DATA_PATH=../../data/seed ./bin/accounts-api > /tmp/api-accounts.log 2>&1 &
API_PID=$!
sleep 3

echo "=== Testing Login Endpoint ==="
curl -X POST http://localhost:8001/login \
  -H "Content-Type: application/json" \
  -d '{"username":"demo@accountstack.com","password":"demo123"}'

echo -e "\n\n=== Testing Protected Endpoint Without Token (should fail with 401) ==="
curl -w "\nHTTP Status: %{http_code}\n" http://localhost:8001/accounts

echo -e "\n=== Getting Token ==="
TOKEN=$(curl -s -X POST http://localhost:8001/login \
  -H "Content-Type: application/json" \
  -d '{"username":"demo@accountstack.com","password":"demo123"}' | python3 -c "import sys, json; print(json.load(sys.stdin)['token'])")
echo "Token received (first 50 chars): ${TOKEN:0:50}..."

echo -e "\n=== Testing Protected Endpoint With Token (should succeed) ==="
curl -s http://localhost:8001/accounts \
  -H "Authorization: Bearer $TOKEN"

echo -e "\n\n=== Testing /me Endpoint With Token ==="
curl -s http://localhost:8001/me \
  -H "Authorization: Bearer $TOKEN"

# Cleanup
echo -e "\n\n=== Cleaning up ==="
kill $API_PID
echo "Done!"

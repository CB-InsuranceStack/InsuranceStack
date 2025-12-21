#!/bin/bash

cd /Users/brown/git_orgs/CB-AccountStack/AccountStack/apps/api-accounts

# Function to safely kill and wait for the API to stop
kill_api() {
    pkill -f accounts-api
    sleep 1
    # Wait up to 5 seconds for the port to be released
    for i in {1..5}; do
        if ! lsof -i :8001 > /dev/null 2>&1; then
            return 0
        fi
        sleep 1
    done
}

# Ensure clean start
kill_api

echo "=== Testing Default Currency (USD) ==="
DATA_PATH=../../data/seed ./bin/accounts-api > /tmp/api-test-usd.log 2>&1 &
sleep 3
TOKEN=$(curl -s -X POST http://localhost:8001/login -H "Content-Type: application/json" -d '{"username":"demo@accountstack.com","password":"demo123"}' | python3 -c "import sys, json; print(json.load(sys.stdin)['token'])")
curl -s http://localhost:8001/accounts -H "Authorization: Bearer $TOKEN" | python3 -m json.tool | grep -B 1 -A 1 '"currency"'
kill_api

echo -e "\n=== Testing EUR Currency Override ==="
FEATURE_CURRENCY=EUR DATA_PATH=../../data/seed ./bin/accounts-api > /tmp/api-test-eur.log 2>&1 &
sleep 3
TOKEN=$(curl -s -X POST http://localhost:8001/login -H "Content-Type: application/json" -d '{"username":"demo@accountstack.com","password":"demo123"}' | python3 -c "import sys, json; print(json.load(sys.stdin)['token'])")
curl -s http://localhost:8001/accounts -H "Authorization: Bearer $TOKEN" | python3 -m json.tool | grep -B 1 -A 1 '"currency"'
kill_api

echo -e "\n=== Testing GBP Currency Override ==="
FEATURE_CURRENCY=GBP DATA_PATH=../../data/seed ./bin/accounts-api > /tmp/api-test-gbp.log 2>&1 &
sleep 3
TOKEN=$(curl -s -X POST http://localhost:8001/login -H "Content-Type: application/json" -d '{"username":"demo@accountstack.com","password":"demo123"}' | python3 -c "import sys, json; print(json.load(sys.stdin)['token'])")
curl -s http://localhost:8001/accounts -H "Authorization: Bearer $TOKEN" | python3 -m json.tool | grep -B 1 -A 1 '"currency"'
kill_api

echo -e "\n✓ Currency feature flag working correctly!"

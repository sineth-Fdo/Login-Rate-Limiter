#!/bin/bash

URL="http://localhost:8080/login"

echo "=== Brute Force Rate Limit Test ==="
echo ""

# 1. Rapid login attempts (triggers per-IP limit at 10 req/10s)
echo "--- Per-IP test: 15 concurrent requests ---"
for i in $(seq 1 15); do
  curl -s -o /dev/null -w "Request $i: %{http_code}\n" -X POST "$URL" \
    -H "Content-Type: application/json" \
    -d '{"username": "admin", "password": "wrong"}' &
done
wait

echo ""

# 2. Different usernames (password guessing)
echo "--- Password guessing: trying common passwords ---"
PASSWORDS=("1234" "admin" "password" "pass" "12345" "qwerty" "letmein" "welcome" "monkey" "dragon" "master" "login")
for PASS in "${PASSWORDS[@]}"; do
  CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$URL" \
    -H "Content-Type: application/json" \
    -d "{\"username\": \"admin\", \"password\": \"$PASS\"}")
  echo "Password '$PASS': $CODE"
done

echo ""

# 3. Burst test (triggers global limit)
echo "--- Burst test: 30 concurrent requests ---"
for i in $(seq 1 30); do
  curl -s -o /dev/null -w "Request $i: %{http_code}\n" -X POST "$URL" \
    -H "Content-Type: application/json" \
    -d '{"username": "user1", "password": "pass"}' &
done
wait

echo ""
echo "=== Done ==="

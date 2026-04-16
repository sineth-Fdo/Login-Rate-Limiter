# Login Rate Limiter

A Go HTTP server with multi-layer rate limiting on a login endpoint, backed by Redis.

## Rate Limiting Strategy

Requests pass through three layers before reaching the handler:

| Layer    | Algorithm      | Limit                 |
| -------- | -------------- | --------------------- |
| Global   | Token Bucket   | 100 burst, 50 req/sec |
| Per-IP   | Sliding Window | 10 requests / 10s     |
| Per-User | Sliding Window | 20 requests / 10s     |

All rate limit state is stored in Redis using atomic Lua scripts.

## Run with Docker

```bash
docker compose up --build
```

This starts both Redis and the app. Server runs on `http://localhost:8080`.

## Run Locally

Start Redis first:

```bash
docker run -d -p 6379:6379 --name ratelimiter-redis redis
```

Then run the app:

```bash
go run .
```

To use a custom Redis address:

```bash
REDIS_ADDR=localhost:6379 go run .
```

## Test

```bash
# login request
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "1234"}'

# response
{"token":"admin_token_143025","message":"Login successful"}
```

### Brute Force Test

Run the test script to simulate attacks and verify rate limiting:

```bash
bash test_ratelimit.sh
```

This runs three tests:
1. **Per-IP** — 15 concurrent requests from the same IP (blocked after 10)
2. **Password guessing** — tries common passwords against `admin`
3. **Burst** — 30 concurrent requests to stress the global limiter

Requests that exceed the limit get `429 Too Many Requests`.


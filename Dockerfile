FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o login-rate-limiter .

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/login-rate-limiter .

EXPOSE 8080

CMD ["./login-rate-limiter"]

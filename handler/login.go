package handler

import (
	"encoding/json"
	"net/http"
	"time"
)

// Request payload
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Response payload
type LoginResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

// Dummy user store (replace with DB)
var users = map[string]string{
	"admin": "1234",
	"user1": "pass",
}

// Fake token generator (replace with JWT)
func generateToken(username string) string {
	return username + "_token_" + time.Now().Format("150405")
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate user
	password, exists := users[req.Username]
	if !exists || password != req.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate token
	token := generateToken(req.Username)

	// IMPORTANT: set user ID header for downstream (rate limiter)
	// In real world → middleware extracts from JWT instead
	r.Header.Set("X-User-ID", req.Username)

	resp := LoginResponse{
		Token:   token,
		Message: "Login successful",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
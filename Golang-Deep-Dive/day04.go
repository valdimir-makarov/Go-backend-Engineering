package main

// Design and implement a token bucket rate limiter in Go. This rate limiter should:

// Allow a maximum number of requests (limit) per minute.

// Use a token bucket algorithm to refill tokens at a fixed rate.

// Be thread-safe.

// Support multiple users (or API keys).

// ✅ Requirements:
// Allow(userID string) bool
// Returns true if the request from userID is allowed, otherwise false.

// Configurable limit (e.g., 60 requests/minute).

// Internal goroutine to refill tokens per user periodically.

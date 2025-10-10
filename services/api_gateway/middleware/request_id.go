package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
)

type requestIDKey string

const RequestIDContextKey requestIDKey = "requestID"

// RequestIDConfig holds configuration for request ID middleware
type RequestIDConfig struct {
	Header    string        // Header name to check for existing request ID
	Generator func() string // Custom ID generator function
}

// RequestIDMiddleware creates a new request ID middleware
func RequestIDMiddleware(next http.Handler) http.Handler {
	config := RequestIDConfig{
		Header:    "X-Request-ID",
		Generator: defaultIDGenerator,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if request ID already exists
		requestID := r.Header.Get(config.Header)
		if requestID == "" {
			// Generate new request ID
			requestID = config.Generator()
		}

		// Set request ID in response header
		w.Header().Set(config.Header, requestID)

		// Store request ID in context for use in handlers and logging
		ctx := context.WithValue(r.Context(), RequestIDContextKey, requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// NewRequestIDMiddleware creates a new request ID middleware with custom config
func NewRequestIDMiddleware(config RequestIDConfig) func(http.Handler) http.Handler {
	if config.Header == "" {
		config.Header = "X-Request-ID"
	}
	if config.Generator == nil {
		config.Generator = defaultIDGenerator
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if request ID already exists
			requestID := r.Header.Get(config.Header)
			if requestID == "" {
				// Generate new request ID
				requestID = config.Generator()
			}

			// Set request ID in response header
			w.Header().Set(config.Header, requestID)

			// Store request ID in context for use in handlers and logging
			ctx := context.WithValue(r.Context(), RequestIDContextKey, requestID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// defaultIDGenerator generates a random request ID
func defaultIDGenerator() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		// Fallback to timestamp-based ID if random generation fails
		return fmt.Sprintf("req-%d", time.Now().UnixNano())
	}
	return "req-" + hex.EncodeToString(bytes)
}

// GetRequestID extracts request ID from context
func GetRequestID(r *http.Request) string {
	if requestID, ok := r.Context().Value(RequestIDContextKey).(string); ok {
		return requestID
	}
	return ""
}

package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"go.uber.org/zap"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
)

// RecoveryMiddleware creates a recovery middleware that handles panics
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				requestID := GetRequestID(r)

				logger.Global().Error("Panic recovered",
					zap.String("requestId", requestID),
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
					zap.Any("error", err),
					zap.String("stack", string(debug.Stack())),
				)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)

				response := map[string]any{
					"error":     "Internal Server Error",
					"message":   "An unexpected error occurred",
					"requestId": requestID,
					"timestamp": time.Now().UTC(),
				}

				if err := json.NewEncoder(w).Encode(response); err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// RecoveryConfig holds configuration for recovery middleware
type RecoveryConfig struct {
	EnableStackTrace bool
	LogStack         bool
	CustomHandler    func(w http.ResponseWriter, r *http.Request, err any)
}

// NewRecoveryMiddleware creates a recovery middleware with config
func NewRecoveryMiddleware(config RecoveryConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					if config.CustomHandler != nil {
						config.CustomHandler(w, r, err)
						return
					}

					requestID := GetRequestID(r)

					fields := []zap.Field{
						zap.String("requestId", requestID),
						zap.String("method", r.Method),
						zap.String("path", r.URL.Path),
						zap.Any("error", err),
					}

					if config.LogStack {
						fields = append(fields, zap.String("stack", string(debug.Stack())))
					}

					logger.Global().Error("Panic recovered", fields...)

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)

					response := map[string]any{
						"error":     "Internal Server Error",
						"message":   "An unexpected error occurred",
						"requestId": requestID,
						"timestamp": time.Now().UTC(),
					}

					if config.EnableStackTrace {
						response["stack"] = fmt.Sprintf("%v", err)
					}

					if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
						http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					}
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

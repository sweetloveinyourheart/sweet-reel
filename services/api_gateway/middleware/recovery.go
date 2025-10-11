package middleware

import (
	"net/http"
	"runtime/debug"

	"go.uber.org/zap"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
	"github.com/sweetloveinyourheart/sweet-reel/services/api_gateway/errors"
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

				writeErrorResponse(w, http.StatusInternalServerError,
					errors.ErrInternalServer.Message,
					errors.ErrInternalServer.Code)
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

					writeErrorResponse(w, http.StatusInternalServerError,
						errors.ErrInternalServer.Message,
						errors.ErrInternalServer.Code)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

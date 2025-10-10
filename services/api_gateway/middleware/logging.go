package middleware

import (
	"bytes"
	"io"
	"maps"
	"net/http"
	"slices"
	"time"

	"go.uber.org/zap"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
)

// LoggingConfig holds configuration for logging middleware
type LoggingConfig struct {
	SkipPaths   []string
	SkipMethods []string
	LogBody     bool
	LogHeaders  bool
	MaxBodySize int
}

// responseWriter wraps http.ResponseWriter to capture response data
type responseWriter struct {
	http.ResponseWriter
	status int
	body   *bytes.Buffer
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.body != nil {
		rw.body.Write(b)
	}
	return rw.ResponseWriter.Write(b)
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// LoggingMiddleware creates a new logging middleware
func LoggingMiddleware(next http.Handler, config LoggingConfig) http.Handler {
	if config.MaxBodySize == 0 {
		config.MaxBodySize = 1024 // 1KB default
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip logging for certain paths
		path := r.URL.Path
		if slices.Contains(config.SkipPaths, path) {
			next.ServeHTTP(w, r)
			return
		}

		// Skip logging for certain methods
		method := r.Method
		if slices.Contains(config.SkipMethods, method) {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()
		requestID := GetRequestID(r)

		// Read and restore request body if needed
		var bodyBytes []byte
		if config.LogBody && r.Body != nil {
			bodyBytes, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Log request
		fields := []zap.Field{
			zap.String("requestId", requestID),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("userAgent", r.Header.Get("User-Agent")),
			zap.String("clientIP", getClientIP(r)),
			zap.Int64("contentLength", r.ContentLength),
		}

		if config.LogHeaders {
			headers := make(map[string][]string)
			maps.Copy(headers, r.Header)
			fields = append(fields, zap.Any("headers", headers))
		}

		if config.LogBody && len(bodyBytes) > 0 && len(bodyBytes) <= config.MaxBodySize {
			fields = append(fields, zap.String("body", string(bodyBytes)))
		}

		logger.Global().Info("HTTP Request", fields...)

		// Wrap response writer to capture response data
		rw := &responseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
			body:           bytes.NewBuffer(nil),
		}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		responseFields := []zap.Field{
			zap.String("requestId", requestID),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status", rw.status),
			zap.Duration("duration", duration),
			zap.Int("responseSize", rw.body.Len()),
		}

		// Log level based on status code
		status := rw.status
		switch {
		case status >= 500:
			logger.Global().Error("HTTP Response", responseFields...)
		case status >= 400:
			logger.Global().Warn("HTTP Response", responseFields...)
		default:
			logger.Global().Info("HTTP Response", responseFields...)
		}
	})
}

// NewLoggingMiddleware creates a new logging middleware with config
func NewLoggingMiddleware(config LoggingConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return LoggingMiddleware(next, config)
	}
}

// AccessLogMiddleware creates a simple access log middleware
func AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap response writer to capture status
		rw := &responseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(rw, r)

		logger.Global().Info("Access Log",
			zap.String("requestId", GetRequestID(r)),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status", rw.status),
			zap.Duration("duration", time.Since(start)),
			zap.String("clientIP", getClientIP(r)),
		)
	})
}

// getClientIP extracts client IP from request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	if xForwardedFor := r.Header.Get("X-Forwarded-For"); xForwardedFor != "" {
		return xForwardedFor
	}

	// Check X-Real-IP header
	if xRealIP := r.Header.Get("X-Real-IP"); xRealIP != "" {
		return xRealIP
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}

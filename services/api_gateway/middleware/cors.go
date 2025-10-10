package middleware

import (
	"net/http"
	"strings"
)

// CORSConfig holds configuration for CORS middleware
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int
}

// CORSMiddleware creates a new CORS middleware
func CORSMiddleware(next http.Handler, config CORSConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Set allowed origins
		if len(config.AllowOrigins) > 0 {
			for _, allowedOrigin := range config.AllowOrigins {
				if allowedOrigin == "*" || allowedOrigin == origin {
					w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
					break
				}
			}
		}

		// Set allowed methods
		if len(config.AllowMethods) > 0 {
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowMethods, ", "))
		}

		// Set allowed headers
		if len(config.AllowHeaders) > 0 {
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowHeaders, ", "))
		}

		// Set exposed headers
		if len(config.ExposeHeaders) > 0 {
			w.Header().Set("Access-Control-Expose-Headers", strings.Join(config.ExposeHeaders, ", "))
		}

		// Set credentials
		if config.AllowCredentials {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// Set max age
		if config.MaxAge > 0 {
			w.Header().Set("Access-Control-Max-Age", string(rune(config.MaxAge)))
		}

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// NewCORSMiddleware creates a CORS middleware with config
func NewCORSMiddleware(config CORSConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return CORSMiddleware(next, config)
	}
}

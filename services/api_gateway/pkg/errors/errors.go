package errors

import (
	"errors"
)

// Authentication errors
var (
	ErrTokenNotFound        = errors.New("token not found")
	ErrInvalidTokenLookup   = errors.New("invalid token lookup configuration")
	ErrInvalidAuthScheme    = errors.New("invalid authentication scheme")
	ErrInvalidSigningMethod = errors.New("invalid signing method")
	ErrInvalidToken         = errors.New("invalid token")
	ErrInvalidClaims        = errors.New("invalid token claims")
	ErrTokenExpired         = errors.New("token expired")
)

// Service proxy errors
var (
	ErrServiceUnavailable = errors.New("service unavailable")
	ErrRequestTimeout     = errors.New("request timeout")
	ErrInvalidResponse    = errors.New("invalid response from service")
	ErrCircuitBreakerOpen = errors.New("circuit breaker is open")
)

// Request/Response errors
var (
	ErrInvalidRequestID   = errors.New("invalid request ID")
	ErrRequestTooLarge    = errors.New("request body too large")
	ErrInvalidContentType = errors.New("invalid content type")
)

// API Gateway specific errors
var (
	ErrRouteNotFound     = errors.New("route not found")
	ErrMethodNotAllowed  = errors.New("method not allowed")
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
)

// APIError represents a structured API error
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Error implements the error interface
func (e APIError) Error() string {
	return e.Message
}

// NewAPIError creates a new APIError
func NewAPIError(code, message, details string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// Common API errors
var (
	ErrUnauthorized = &APIError{
		Code:    "UNAUTHORIZED",
		Message: "Authentication required",
	}
	ErrForbidden = &APIError{
		Code:    "FORBIDDEN",
		Message: "Access denied",
	}
	ErrNotFound = &APIError{
		Code:    "NOT_FOUND",
		Message: "Resource not found",
	}
	ErrBadRequest = &APIError{
		Code:    "BAD_REQUEST",
		Message: "Invalid request",
	}
	ErrInternalServer = &APIError{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "Internal server error",
	}
	ErrServiceDown = &APIError{
		Code:    "SERVICE_UNAVAILABLE",
		Message: "Service temporarily unavailable",
	}
)

// Auth API errors
var (
	ErrAuthenticationTokenRequired = &APIError{
		Code:    "AUTH_TOKEN_MISSING",
		Message: "Authorization token required",
	}

	ErrAuthenticationTokenInvalid = &APIError{
		Code:    "AUTH_TOKEN_INVALID",
		Message: "Invalid or expired token",
	}

	ErrAuthNoRoles = &APIError{
		Code:    "AUTH_NO_ROLES",
		Message: "Access denied: no roles found",
	}

	ErrAuthInsufficientPermissions = &APIError{
		Code:    "AUTH_INSUFFICIENT_PERMISSIONS",
		Message: "Access denied: insufficient permissions",
	}

	ErrAuthNoPermissions = &APIError{
		Code:    "AUTH_NO_PERMISSIONS",
		Message: "Access denied: no permissions found",
	}
)

package errors

import (
	"errors"
	"net/http"
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

// HTTPError represents an HTTP error with status code and details
type HTTPError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
	Code       string `json:"code"`
}

func (e *HTTPError) Error() string {
	return e.Message
}

// NewHTTPError creates a new HTTP error
func NewHTTPError(statusCode int, message, code string) *HTTPError {
	return &HTTPError{
		StatusCode: statusCode,
		Message:    message,
		Code:       code,
	}
}

// HTTP status error variables for helper functions
var (
	ErrHTTPBadRequest = &HTTPError{
		StatusCode: http.StatusBadRequest,
		Message:    "Bad Request",
		Code:       "BAD_REQUEST",
	}

	ErrHTTPUnauthorized = &HTTPError{
		StatusCode: http.StatusUnauthorized,
		Message:    "Unauthorized",
		Code:       "UNAUTHORIZED",
	}

	ErrHTTPForbidden = &HTTPError{
		StatusCode: http.StatusForbidden,
		Message:    "Forbidden",
		Code:       "FORBIDDEN",
	}

	ErrHTTPNotFound = &HTTPError{
		StatusCode: http.StatusNotFound,
		Message:    "Not Found",
		Code:       "NOT_FOUND",
	}

	ErrHTTPMethodNotAllowed = &HTTPError{
		StatusCode: http.StatusMethodNotAllowed,
		Message:    "Method Not Allowed",
		Code:       "METHOD_NOT_ALLOWED",
	}

	ErrHTTPInternalServer = &HTTPError{
		StatusCode: http.StatusInternalServerError,
		Message:    "Internal Server Error",
		Code:       "INTERNAL_ERROR",
	}
)

// Common API errors
var (
	ErrBadRequest = &APIError{
		Code:    "BAD_REQUEST",
		Message: "Bad Request",
	}
	ErrUnauthorized = &APIError{
		Code:    "UNAUTHORIZED",
		Message: "Unauthorized",
	}
	ErrForbidden = &APIError{
		Code:    "FORBIDDEN",
		Message: "Forbidden",
	}
	ErrNotFound = &APIError{
		Code:    "NOT_FOUND",
		Message: "Not Found",
	}
	ErrMethodNotAllowed = &APIError{
		Code:    "METHOD_NOT_ALLOWED",
		Message: "Method Not Allowed",
	}
	ErrInternalServer = &APIError{
		Code:    "INTERNAL_ERROR",
		Message: "Internal Server Error",
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

	ErrOAuthLoginFailed = &APIError{
		Code:    "OAUTH_LOGIN_FAILED",
		Message: "Login failed",
	}
)

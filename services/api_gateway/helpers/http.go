package helpers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
	"github.com/sweetloveinyourheart/sweet-reel/services/api_gateway/errors"
)

// RequestBody interface for request body validation
type RequestBody interface {
	Validate() error
}

// JSONResponse represents a standard JSON response structure
type JSONResponse struct {
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// ParseJSONBody parses and validates JSON request body
func ParseJSONBody[T RequestBody](r *http.Request, body *T) error {
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		return errors.NewHTTPError(
			http.StatusBadRequest,
			"Invalid request body format",
			"INVALID_JSON",
		)
	}

	if err := (*body).Validate(); err != nil {
		return errors.NewHTTPError(
			http.StatusBadRequest,
			err.Error(),
			"VALIDATION_ERROR",
		)
	}

	return nil
}

// WriteJSONResponse writes a successful JSON response
func WriteJSONResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Global().Error("Failed to encode JSON response", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// WriteJSONSuccess writes a successful JSON response with data
func WriteJSONSuccess(w http.ResponseWriter, data any) {
	WriteJSONResponse(w, http.StatusOK, data)
}

// WriteJSONCreated writes a 201 Created JSON response
func WriteJSONCreated(w http.ResponseWriter, data any) {
	WriteJSONResponse(w, http.StatusCreated, data)
}

// WriteErrorResponse writes an error response in JSON format
func WriteErrorResponse(w http.ResponseWriter, err error) {
	var httpErr *errors.HTTPError

	// Check if it's already an HTTPError
	if he, ok := err.(*errors.HTTPError); ok {
		httpErr = he
	} else {
		// Default to internal server error for unknown errors
		httpErr = errors.NewHTTPError(
			http.StatusInternalServerError,
			"Internal Server Error",
			"INTERNAL_ERROR",
		)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpErr.StatusCode)

	response := JSONResponse{
		Error:   httpErr.Message,
		Code:    httpErr.Code,
		Message: httpErr.Message,
	}

	if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
		logger.Global().Error("Failed to encode error response", zap.Error(encodeErr))
		http.Error(w, httpErr.Message, httpErr.StatusCode)
	}
}

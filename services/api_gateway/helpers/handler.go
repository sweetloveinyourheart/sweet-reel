package helpers

import (
	"net/http"

	"github.com/sweetloveinyourheart/sweet-reel/services/api_gateway/errors"
)

// POST creates a handler that only accepts POST requests
func POST(handler http.HandlerFunc) http.Handler {
	return methodHandler(http.MethodPost, handler)
}

// GET creates a handler that only accepts GET requests
func GET(handler http.HandlerFunc) http.Handler {
	return methodHandler(http.MethodGet, handler)
}

// PUT creates a handler that only accepts PUT requests
func PUT(handler http.HandlerFunc) http.Handler {
	return methodHandler(http.MethodPut, handler)
}

// DELETE creates a handler that only accepts DELETE requests
func DELETE(handler http.HandlerFunc) http.Handler {
	return methodHandler(http.MethodDelete, handler)
}

// PATCH creates a handler that only accepts PATCH requests
func PATCH(handler http.HandlerFunc) http.Handler {
	return methodHandler(http.MethodPatch, handler)
}

// methodHandler ensures the request uses the correct HTTP method
func methodHandler(method string, handler http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			WriteErrorResponse(w, errors.ErrHTTPMethodNotAllowed)
			return
		}
		handler(w, r)
	})
}

package oauth2

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"slices"
	"time"

	"golang.org/x/oauth2"
)

// GenerateState generates a random state parameter for OAuth2 flow
func GenerateState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate random state: %w", err)
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

// ExtractCodeFromRequest extracts the authorization code from HTTP request
func ExtractCodeFromRequest(r *http.Request) (string, error) {
	code := r.URL.Query().Get("code")
	if code == "" {
		return "", fmt.Errorf("authorization code not found in request")
	}

	return code, nil
}

// ExtractStateFromRequest extracts the state parameter from HTTP request
func ExtractStateFromRequest(r *http.Request) (string, error) {
	state := r.URL.Query().Get("state")
	if state == "" {
		return "", fmt.Errorf("state parameter not found in request")
	}

	return state, nil
}

// CheckOAuthError checks for OAuth2 error in the request
func CheckOAuthError(r *http.Request) error {
	errorCode := r.URL.Query().Get("error")
	if errorCode == "" {
		return nil
	}

	errorDescription := r.URL.Query().Get("error_description")
	if errorDescription != "" {
		return fmt.Errorf("oauth2 error %s: %s", errorCode, errorDescription)
	}

	return fmt.Errorf("oauth2 error: %s", errorCode)
}

// IsTokenExpired checks if a token is expired or will expire soon
func IsTokenExpired(token *oauth2.Token, bufferTime time.Duration) bool {
	if token.Expiry.IsZero() {
		return false // No expiry time set
	}

	return time.Now().Add(bufferTime).After(token.Expiry)
}

// TokenTimeUntilExpiry returns the duration until token expires
func TokenTimeUntilExpiry(token *oauth2.Token) time.Duration {
	if token.Expiry.IsZero() {
		return time.Duration(0) // No expiry time set
	}

	remaining := time.Until(token.Expiry)
	if remaining < 0 {
		return time.Duration(0)
	}

	return remaining
}

// SanitizeRedirectURL ensures the redirect URL is properly formatted
func SanitizeRedirectURL(redirectURL string) string {
	// Basic sanitization - in production, you might want more robust validation
	if redirectURL == "" {
		return "http://localhost:8080/auth/callback"
	}

	return redirectURL
}

// ValidateCallbackURL validates if a callback URL is allowed
func ValidateCallbackURL(callbackURL string, allowedURLs []string) bool {
	if len(allowedURLs) == 0 {
		return true // No restrictions
	}

	return slices.Contains(allowedURLs, callbackURL)
}

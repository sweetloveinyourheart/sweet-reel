package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"slices"
	"strings"

	"github.com/golang-jwt/jwt/v4"

	"github.com/sweetloveinyourheart/sweet-reel/services/api_gateway/errors"
)

type contextKey string

const (
	AuthTokenKey          contextKey = "authToken"
	UserContextKey        contextKey = "user"
	UserIDContextKey      contextKey = "userID"
	RolesContextKey       contextKey = "roles"
	PermissionsContextKey contextKey = "permissions"
)

// AuthConfig holds authentication middleware configuration
type AuthConfig struct {
	SigningKey  string
	SkipPaths   []string
	TokenLookup string // "header:Authorization"
	AuthScheme  string // "Bearer"
}

// NewAuthMiddleware creates a new JWT authentication middleware
func NewAuthMiddleware(config AuthConfig) func(http.Handler) http.Handler {
	if config.TokenLookup == "" {
		config.TokenLookup = "header:Authorization"
	}
	if config.AuthScheme == "" {
		config.AuthScheme = "Bearer"
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			for _, skipPath := range config.SkipPaths {
				if strings.HasPrefix(path, skipPath) {
					next.ServeHTTP(w, r)
					return
				}
			}

			token, err := extractToken(r, config)
			if err != nil {
				writeErrorResponse(w, http.StatusUnauthorized,
					errors.ErrAuthenticationTokenRequired.Message,
					errors.ErrAuthenticationTokenRequired.Code)
				return
			}

			claims, err := validateToken(token, config.SigningKey)
			if err != nil {
				writeErrorResponse(w, http.StatusUnauthorized,
					errors.ErrAuthenticationTokenInvalid.Message,
					errors.ErrAuthenticationTokenInvalid.Code)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, AuthTokenKey, token)
			ctx = context.WithValue(ctx, UserContextKey, claims)
			if userID, ok := claims["user_id"]; ok {
				ctx = context.WithValue(ctx, UserIDContextKey, userID)
			}
			if roles, ok := claims["roles"]; ok {
				ctx = context.WithValue(ctx, RolesContextKey, roles)
			}
			if permissions, ok := claims["permissions"]; ok {
				ctx = context.WithValue(ctx, PermissionsContextKey, permissions)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// extractToken extracts JWT token from request based on configuration
func extractToken(r *http.Request, config AuthConfig) (string, error) {
	parts := strings.Split(config.TokenLookup, ":")
	if len(parts) != 2 {
		return "", errors.ErrInvalidTokenLookup
	}

	switch parts[0] {
	case "header":
		authHeader := r.Header.Get(parts[1])
		if authHeader == "" {
			return "", errors.ErrTokenNotFound
		}

		// Check if it follows the expected scheme
		if config.AuthScheme != "" {
			prefix := config.AuthScheme + " "
			if !strings.HasPrefix(authHeader, prefix) {
				return "", errors.ErrInvalidAuthScheme
			}
			return strings.TrimPrefix(authHeader, prefix), nil
		}
		return authHeader, nil

	case "query":
		token := r.URL.Query().Get(parts[1])
		if token == "" {
			return "", errors.ErrTokenNotFound
		}
		return token, nil

	case "cookie":
		cookie, err := r.Cookie(parts[1])
		if err != nil || cookie.Value == "" {
			return "", errors.ErrTokenNotFound
		}
		return cookie.Value, nil

	default:
		return "", errors.ErrInvalidTokenLookup
	}
}

// validateToken validates JWT token and returns claims
func validateToken(tokenString, signingKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrInvalidSigningMethod
		}
		return []byte(signingKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.ErrInvalidClaims
	}

	return claims, nil
}

// RequireRole creates middleware that requires specific roles
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRoles, ok := r.Context().Value(RolesContextKey).([]any)
			if !ok {
				writeErrorResponse(w, http.StatusForbidden,
					errors.ErrAuthNoRoles.Message,
					errors.ErrAuthNoRoles.Code)
				return
			}

			roleStrings := make([]string, len(userRoles))
			for i, role := range userRoles {
				if str, ok := role.(string); ok {
					roleStrings[i] = str
				}
			}

			// Check if user has any of the required roles
			for _, requiredRole := range roles {
				if slices.Contains(roleStrings, requiredRole) {
					next.ServeHTTP(w, r)
					return
				}
			}

			writeErrorResponse(w, http.StatusForbidden,
				errors.ErrAuthInsufficientPermissions.Message,
				errors.ErrAuthInsufficientPermissions.Code)
		})
	}
}

// RequirePermission creates middleware that requires specific permissions
func RequirePermission(permissions ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userPermissions, ok := r.Context().Value(PermissionsContextKey).([]any)
			if !ok {
				writeErrorResponse(w, http.StatusForbidden,
					errors.ErrAuthNoPermissions.Message,
					errors.ErrAuthNoPermissions.Code)
				return
			}

			permissionStrings := make([]string, len(userPermissions))
			for i, perm := range userPermissions {
				if str, ok := perm.(string); ok {
					permissionStrings[i] = str
				}
			}

			for _, requiredPerm := range permissions {
				if slices.Contains(permissionStrings, requiredPerm) {
					next.ServeHTTP(w, r)
					return
				}
			}

			writeErrorResponse(w, http.StatusForbidden,
				errors.ErrAuthInsufficientPermissions.Message,
				errors.ErrAuthInsufficientPermissions.Code)
		})
	}
}

// writeErrorResponse writes an error response in JSON format
func writeErrorResponse(w http.ResponseWriter, statusCode int, message, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := map[string]string{
		"error": message,
		"code":  code,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, message, statusCode)
	}
}

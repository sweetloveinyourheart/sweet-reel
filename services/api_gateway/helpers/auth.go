package helpers

import (
	"net/http"

	"github.com/sweetloveinyourheart/sweet-reel/services/api_gateway/middleware"
)

// GetAuthToken returns the authentication token from the request context.
func GetAuthToken(r *http.Request) string {
	if token, ok := r.Context().Value(middleware.AuthTokenKey).(string); ok {
		return token
	}
	return ""
}

// GetClaims returns the full claims map from the request context.
func GetClaims(r *http.Request) map[string]any {
	if claims, ok := r.Context().Value(middleware.UserContextKey).(map[string]any); ok {
		return claims
	}
	return nil
}

// GetUserID returns the user ID from context (string or nil if missing).
func GetUserID(r *http.Request) string {
	if userID, ok := r.Context().Value(middleware.UserIDContextKey).(string); ok {
		return userID
	}
	return ""
}

// GetRoles returns user roles from context (usually []string).
func GetRoles(r *http.Request) []string {
	if roles, ok := r.Context().Value(middleware.RolesContextKey).([]string); ok {
		return roles
	}

	// sometimes JWT claims might store roles as []interface{}
	if rolesIface, ok := r.Context().Value(middleware.RolesContextKey).([]interface{}); ok {
		roles := make([]string, 0, len(rolesIface))
		for _, v := range rolesIface {
			if s, ok := v.(string); ok {
				roles = append(roles, s)
			}
		}
		return roles
	}

	return nil
}

// GetPermissions returns user permissions from context.
func GetPermissions(r *http.Request) []string {
	if permissions, ok := r.Context().Value(middleware.PermissionsContextKey).([]string); ok {
		return permissions
	}

	if permsIface, ok := r.Context().Value(middleware.PermissionsContextKey).([]interface{}); ok {
		perms := make([]string, 0, len(permsIface))
		for _, v := range permsIface {
			if s, ok := v.(string); ok {
				perms = append(perms, s)
			}
		}
		return perms
	}

	return nil
}

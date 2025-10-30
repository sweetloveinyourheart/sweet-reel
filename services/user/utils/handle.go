package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gofrs/uuid"
)

// GenerateHandleFromName creates a URL-friendly handle from a user's name
// Example: "John Doe" -> "@johndoe"
// Example: "Jane Smith" -> "@janesmith"
// Example: "user@email.com" -> "@useremailcom"
func GenerateHandleFromName(name string) string {
	// Convert to lowercase
	handle := strings.ToLower(name)

	// Remove special characters and replace spaces/non-alphanumeric with empty string
	reg := regexp.MustCompile("[^a-z0-9]+")
	handle = reg.ReplaceAllString(handle, "")

	// Ensure handle is not empty (fallback for names with only special characters)
	if handle == "" {
		// Generate a random handle using UUID
		randomID := uuid.Must(uuid.NewV7()).String()[:8]
		handle = fmt.Sprintf("user%s", randomID)
	}

	// Limit handle length to 30 characters (excluding @)
	if len(handle) > 30 {
		handle = handle[:30]
	}

	// Add @ prefix
	return "@" + handle
}

// GenerateUniqueHandle creates a unique handle by appending a suffix if needed
// Example: "@johndoe" -> "@johndoe123" if "@johndoe" is taken
func GenerateUniqueHandle(baseHandle string, suffix int) string {
	// Remove @ prefix if present
	handle := strings.TrimPrefix(baseHandle, "@")

	if suffix > 0 {
		handle = fmt.Sprintf("%s%d", handle, suffix)
	}

	return "@" + handle
}

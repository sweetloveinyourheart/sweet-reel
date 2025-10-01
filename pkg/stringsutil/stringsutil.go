package stringsutil

import (
	"strings"
)

// IsBlank returns true if a string is empty or contains only whitespace.
func IsBlank(s string) bool {
	return strings.TrimSpace(s) == ""
}

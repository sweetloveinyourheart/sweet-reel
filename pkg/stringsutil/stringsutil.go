package stringsutil

import (
	"regexp"
	"strings"
)

const (
	// MachineNameRegex matches acceptable formats for machine names: letters, numbers, and underscores.
	MachineNameRegex         = "^[a-z0-9_]+$"
	MachineNameExpandedRegex = "^[a-z0-9_-]+$"
)

// IsBlank returns true if a string is empty or contains only whitespace.
func IsBlank(s string) bool {
	return strings.TrimSpace(s) == ""
}

// IsMachineName returns true if the passed in string is a valid machine name
func IsMachineName(s string) bool {
	return regexp.MustCompile(MachineNameRegex).MatchString(s)
}

// ToMachineName takes a string and returns that string with all non-alpha or underscor characters removed and spaces
// converted to underscores.
// Ex: "test, tes1t!" => "test_tes1t"
func ToMachineName(s string) string {
	if IsMachineName(s) {
		return s
	}
	working := strings.Trim(regexp.MustCompile("[^a-z0-9_]+").ReplaceAllString(strings.ToLower(s), "_"), "_")
	var last rune

	var sb strings.Builder

	for _, r := range working {
		if r != last || r != '_' {
			sb.WriteRune(r)
			last = r
		} else {
			continue
		}
	}
	return sb.String()
}

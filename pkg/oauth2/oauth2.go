// Package oauth provides multi-provider OAuth2 authentication functionality
// supporting Google, GitHub, and other OAuth2 providers.
package oauth2

const (
	// Common OAuth2 prompts
	PromptConsent       = "consent"
	PromptSelectAccount = "select_account"
	PromptLogin         = "login"
	PromptNone          = "none"
)

// Predefined scope sets for common use cases

// Google Scopes
var (
	// GoogleProfileScopes includes basic profile information
	GoogleProfileScopes = []string{
		"https://www.googleapis.com/auth/userinfo.profile",
	}

	// GoogleEmailScopes includes email information
	GoogleEmailScopes = []string{
		"https://www.googleapis.com/auth/userinfo.email",
	}

	// GoogleBasicScopes includes both profile and email (most common)
	GoogleBasicScopes = []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	}

	// GoogleExtendedScopes includes additional Google services
	GoogleExtendedScopes = []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
		"https://www.googleapis.com/auth/drive.readonly",
		"https://www.googleapis.com/auth/calendar.readonly",
	}
)

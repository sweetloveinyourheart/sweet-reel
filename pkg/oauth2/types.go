package oauth2

// UserInfo represents normalized user information from any OAuth2 provider
type UserInfo struct {
	ID         string `json:"id"`                   // Unique identifier from the provider
	Email      string `json:"email"`                // User's email address
	Name       string `json:"name"`                 // User's full name
	Username   string `json:"username,omitempty"`   // Provider-specific username (if available)
	Picture    string `json:"picture,omitempty"`    // URL to user's profile picture
	Provider   string `json:"provider"`             // OAuth provider name (google, github, etc.)
	ProviderID string `json:"provider_id"`          // Raw ID from the provider
	Verified   bool   `json:"verified,omitempty"`   // Whether email is verified
	Locale     string `json:"locale,omitempty"`     // User's locale (if available)
	FirstName  string `json:"first_name,omitempty"` // User's first name
	LastName   string `json:"last_name,omitempty"`  // User's last name
	Bio        string `json:"bio,omitempty"`        // User's bio/description
	Location   string `json:"location,omitempty"`   // User's location
	Company    string `json:"company,omitempty"`    // User's company
}

// TokenInfo represents token information from OAuth provider
type TokenInfo struct {
	IssuedTo      string `json:"issued_to"`
	Audience      string `json:"audience"`
	Scope         string `json:"scope"`
	ExpiresIn     int    `json:"expires_in"`
	AccessType    string `json:"access_type"`
	VerifiedEmail bool   `json:"verified_email,omitempty"`
	Email         string `json:"email,omitempty"`
}

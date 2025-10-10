package oauth2

import (
	"encoding/json"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GoogleProvider implements OAuth2 for Google
type GoogleProvider struct{}

// googleUserInfoResponse represents the response from Google's userinfo API
type googleUserInfoResponse struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// GetEndpoint returns Google's OAuth2 endpoint
func (g *GoogleProvider) GetEndpoint() oauth2.Endpoint {
	return google.Endpoint
}

// GetUserInfoURL returns Google's user info API URL
func (g *GoogleProvider) GetUserInfoURL() string {
	return "https://www.googleapis.com/oauth2/v2/userinfo"
}

// GetTokenInfoURL returns Google's token info API URL
func (g *GoogleProvider) GetTokenInfoURL() string {
	return "https://www.googleapis.com/oauth2/v1/tokeninfo"
}

// ParseUserInfo parses Google's user info response into normalized UserInfo
func (g *GoogleProvider) ParseUserInfo(data []byte) (*UserInfo, error) {
	var googleResp googleUserInfoResponse
	if err := json.Unmarshal(data, &googleResp); err != nil {
		return nil, fmt.Errorf("failed to parse Google user info: %w", err)
	}

	return &UserInfo{
		ID:         googleResp.ID,
		Email:      googleResp.Email,
		Name:       googleResp.Name,
		Username:   googleResp.Email, // Google doesn't have username, use email
		Picture:    googleResp.Picture,
		Provider:   string(ProviderGoogle),
		ProviderID: googleResp.ID,
		Verified:   googleResp.VerifiedEmail,
		Locale:     googleResp.Locale,
		FirstName:  googleResp.GivenName,
		LastName:   googleResp.FamilyName,
	}, nil
}

// GetDefaultScopes returns Google's default OAuth2 scopes
func (g *GoogleProvider) GetDefaultScopes() []string {
	return []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	}
}

// GetProviderName returns the human-readable name
func (g *GoogleProvider) GetProviderName() string {
	return "Google"
}

// GetProviderType returns the provider type
func (g *GoogleProvider) GetProviderType() ProviderType {
	return ProviderGoogle
}

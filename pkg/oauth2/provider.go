package oauth2

import (
	"golang.org/x/oauth2"
)

// ProviderType represents the type of OAuth provider
type ProviderType string

const (
	ProviderGoogle ProviderType = "google"
)

// Provider interface abstracts OAuth2 provider differences
type Provider interface {
	// GetEndpoint returns the OAuth2 endpoint for this provider
	GetEndpoint() oauth2.Endpoint

	// GetUserInfoURL returns the URL to fetch user information
	GetUserInfoURL() string

	// GetTokenInfoURL returns the URL to validate tokens (if supported)
	GetTokenInfoURL() string

	// ParseUserInfo parses provider-specific user info response into normalized UserInfo
	ParseUserInfo(data []byte) (*UserInfo, error)

	// GetDefaultScopes returns the default scopes for this provider
	GetDefaultScopes() []string

	// GetProviderName returns the human-readable name of the provider
	GetProviderName() string

	// GetProviderType returns the provider type
	GetProviderType() ProviderType
}

// NewProvider creates a new provider instance based on the provider type
func NewProvider(providerType ProviderType) (Provider, error) {
	switch providerType {
	case ProviderGoogle:
		return &GoogleProvider{}, nil
	default:
		return nil, &ErrUnsupportedProvider{Provider: string(providerType)}
	}
}

// ErrUnsupportedProvider is returned when an unsupported provider is requested
type ErrUnsupportedProvider struct {
	Provider string
}

func (e *ErrUnsupportedProvider) Error() string {
	return "unsupported OAuth provider: " + e.Provider
}

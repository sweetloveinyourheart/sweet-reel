package oauth2

import (
	"golang.org/x/oauth2"
)

// Config holds OAuth2 configuration for any provider
type Config struct {
	Provider     Provider
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

// NewConfig creates a new OAuth2 configuration for the specified provider
func NewConfig(provider Provider, clientID, clientSecret, redirectURL string, scopes []string) *Config {
	if scopes == nil {
		scopes = provider.GetDefaultScopes()
	}

	return &Config{
		Provider:     provider,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       scopes,
	}
}

// NewConfigWithType creates a new OAuth2 configuration using provider type
func NewConfigWithType(providerType ProviderType, clientID, clientSecret, redirectURL string, scopes []string) (*Config, error) {
	provider, err := NewProvider(providerType)
	if err != nil {
		return nil, err
	}

	return NewConfig(provider, clientID, clientSecret, redirectURL, scopes), nil
}

// ToOAuth2Config converts to golang.org/x/oauth2.Config
func (c *Config) ToOAuth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		RedirectURL:  c.RedirectURL,
		Scopes:       c.Scopes,
		Endpoint:     c.Provider.GetEndpoint(),
	}
}

// GetProviderName returns the provider's human-readable name
func (c *Config) GetProviderName() string {
	return c.Provider.GetProviderName()
}

// GetProviderType returns the provider type
func (c *Config) GetProviderType() ProviderType {
	return c.Provider.GetProviderType()
}

// DefaultScopes returns commonly used OAuth2 scopes for the provider
func (c *Config) DefaultScopes() []string {
	return c.Provider.GetDefaultScopes()
}

package oauth2

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
)

type IOAuthClient interface {
	GetAuthURL(state string) string
	GetAuthURLWithOptions(state string, options ...oauth2.AuthCodeOption) string
	ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error)
	GetUserInfo(ctx context.Context, token *oauth2.Token) (*UserInfo, error)
	ValidateToken(ctx context.Context, token *oauth2.Token) error
}

// Client provides OAuth2 functionality for any provider
type Client struct {
	config   *oauth2.Config
	provider Provider
}

// NewClient creates a new OAuth2 client
func NewClient(config *Config) IOAuthClient {
	return &Client{
		config:   config.ToOAuth2Config(),
		provider: config.Provider,
	}
}

// GetAuthURL generates the OAuth2 authorization URL
func (c *Client) GetAuthURL(state string) string {
	return c.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// GetAuthURLWithOptions generates the OAuth2 authorization URL with additional options
func (c *Client) GetAuthURLWithOptions(state string, options ...oauth2.AuthCodeOption) string {
	return c.config.AuthCodeURL(state, options...)
}

// ExchangeCode exchanges the authorization code for tokens
func (c *Client) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	return c.config.Exchange(ctx, code)
}

// GetUserInfo retrieves user information using the access token
func (c *Client) GetUserInfo(ctx context.Context, token *oauth2.Token) (*UserInfo, error) {
	client := c.config.Client(ctx, token)

	req, err := http.NewRequestWithContext(ctx, "GET", c.provider.GetUserInfoURL(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	userInfo, err := c.provider.ParseUserInfo(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	return userInfo, nil
}

// ValidateToken validates if the token is still valid
func (c *Client) ValidateToken(ctx context.Context, token *oauth2.Token) error {
	if !token.Valid() {
		return fmt.Errorf("token is invalid or expired")
	}

	// Provider-specific token validation
	client := c.config.Client(ctx, token)

	switch c.provider.GetProviderType() {
	case ProviderGoogle:
		// Google has a specific token info endpoint
		tokenInfoURL := c.provider.GetTokenInfoURL() + "?access_token=" + token.AccessToken
		req, err := http.NewRequestWithContext(ctx, "GET", tokenInfoURL, nil)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to validate token: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("token validation failed: status %d", resp.StatusCode)
		}

	case ProviderGitHub:
		// GitHub doesn't have a token info endpoint, so we validate by making a user API call
		req, err := http.NewRequestWithContext(ctx, "GET", c.provider.GetUserInfoURL(), nil)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to validate token: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("token validation failed: status %d", resp.StatusCode)
		}

	default:
		// For unknown providers, just check if we can make a user info request
		req, err := http.NewRequestWithContext(ctx, "GET", c.provider.GetUserInfoURL(), nil)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to validate token: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("token validation failed: status %d", resp.StatusCode)
		}
	}

	return nil
}

// GetProvider returns the provider instance
func (c *Client) GetProvider() Provider {
	return c.provider
}

// GetProviderName returns the provider's human-readable name
func (c *Client) GetProviderName() string {
	return c.provider.GetProviderName()
}

// GetProviderType returns the provider type
func (c *Client) GetProviderType() ProviderType {
	return c.provider.GetProviderType()
}

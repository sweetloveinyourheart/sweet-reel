package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"

	oauth2Pkg "github.com/sweetloveinyourheart/sweet-reel/pkg/oauth2"
)

// Ensure that MockOAuth2 implements oauth2Pkg.IOAuthClient
var _ oauth2Pkg.IOAuthClient = (*MockOAuthClient)(nil)

// MockOAuthClient is a mock implementation of IOAuthClient.
type MockOAuthClient struct {
	mock.Mock
}

func (m *MockOAuthClient) GetAuthURL(state string) string {
	args := m.Called(state)
	return args.String(0)
}

func (m *MockOAuthClient) GetAuthURLWithOptions(state string, options ...oauth2.AuthCodeOption) string {
	args := m.Called(state, options)
	return args.String(0)
}

func (m *MockOAuthClient) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	args := m.Called(ctx, code)
	token, _ := args.Get(0).(*oauth2.Token)
	return token, args.Error(1)
}

func (m *MockOAuthClient) GetUserInfo(ctx context.Context, token *oauth2.Token) (*oauth2Pkg.UserInfo, error) {
	args := m.Called(ctx, token)
	info, _ := args.Get(0).(*oauth2Pkg.UserInfo)
	return info, args.Error(1)
}

func (m *MockOAuthClient) ValidateToken(ctx context.Context, token *oauth2.Token) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

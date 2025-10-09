package actions_test

import (
	"context"
	"errors"
	"time"

	"connectrpc.com/connect"
	"github.com/gofrs/uuid"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/grpc"
	proto "github.com/sweetloveinyourheart/sweet-reel/proto/code/user/go"
	"github.com/sweetloveinyourheart/sweet-reel/services/user/actions"
	"github.com/sweetloveinyourheart/sweet-reel/services/user/models"
)

func (as *ActionsSuite) TestActions_UpsertOAuthUser_Success_ExistingUser() {
	as.setupEnvironment()

	// Setup test data
	authUserID := uuid.Must(uuid.NewV7())
	existingUserID := uuid.Must(uuid.NewV7())
	provider := "github"
	providerUserID := "github-789012"
	email := "existing@example.com"
	name := "Existing User"
	picture := "https://example.com/existing-avatar.jpg"

	// Setup context with auth token
	ctx := context.WithValue(context.Background(), grpc.AuthToken, authUserID)

	// Create existing user
	existingUser := &models.User{
		ID:        existingUserID,
		Email:     "existing@example.com",
		Name:      "Existing User",
		Picture:   "https://example.com/existing-avatar.jpg",
		CreatedAt: time.Now().Add(-24 * time.Hour),
		UpdatedAt: time.Now().Add(-1 * time.Hour),
	}

	// Mock GetUserWithIdentity to return existing user
	as.mockUserRepository.On("GetUserWithIdentity", ctx, provider, providerUserID).Return(existingUser, nil)

	// Setup request
	request := &connect.Request[proto.UpsertOAuthUserRequest]{
		Msg: &proto.UpsertOAuthUserRequest{
			Provider:       provider,
			ProviderUserId: providerUserID,
			Email:          email,
			Name:           name,
			Picture:        picture,
		},
	}

	// Execute
	actionsInstance := actions.NewActions(ctx, "test-token")
	response, err := actionsInstance.UpsertOAuthUser(ctx, request)

	// Assertions
	as.NoError(err)
	as.NotNil(response)
	as.NotNil(response.Msg)
	as.False(response.Msg.IsNewUser) // Should be false for existing user
	as.NotNil(response.Msg.User)

	// Verify mocks were called
	as.mockUserRepository.AssertExpectations(as.T())
}

func (as *ActionsSuite) TestActions_UpsertOAuthUser_Unauthenticated() {
	as.setupEnvironment()

	// Setup test data without auth token
	provider := "google"
	providerUserID := "google-123456"
	email := "test@example.com"
	name := "Test User"
	picture := "https://example.com/avatar.jpg"

	// Setup context WITHOUT auth token
	ctx := context.Background()

	// Setup request
	request := &connect.Request[proto.UpsertOAuthUserRequest]{
		Msg: &proto.UpsertOAuthUserRequest{
			Provider:       provider,
			ProviderUserId: providerUserID,
			Email:          email,
			Name:           name,
			Picture:        picture,
		},
	}

	// Execute
	actionsInstance := actions.NewActions(ctx, "test-token")
	response, err := actionsInstance.UpsertOAuthUser(ctx, request)

	// Assertions
	as.Error(err)
	as.Nil(response)
	as.Contains(err.Error(), "invalid session")

	// Verify no repository methods were called
	as.mockUserRepository.AssertNotCalled(as.T(), "GetUserWithIdentity")
}

func (as *ActionsSuite) TestActions_UpsertOAuthUser_GetUserWithIdentity_Error() {
	as.setupEnvironment()

	// Setup test data
	authUserID := uuid.Must(uuid.NewV7())
	provider := "google"
	providerUserID := "google-123456"
	email := "test@example.com"
	name := "Test User"
	picture := "https://example.com/avatar.jpg"

	// Setup context with auth token
	ctx := context.WithValue(context.Background(), grpc.AuthToken, authUserID)

	// Mock GetUserWithIdentity to return an error
	dbError := errors.New("database connection failed")
	as.mockUserRepository.On("GetUserWithIdentity", ctx, provider, providerUserID).Return(nil, dbError)

	// Setup request
	request := &connect.Request[proto.UpsertOAuthUserRequest]{
		Msg: &proto.UpsertOAuthUserRequest{
			Provider:       provider,
			ProviderUserId: providerUserID,
			Email:          email,
			Name:           name,
			Picture:        picture,
		},
	}

	// Execute
	actionsInstance := actions.NewActions(ctx, "test-token")
	response, err := actionsInstance.UpsertOAuthUser(ctx, request)

	// Assertions
	as.Error(err)
	as.Nil(response)

	// Verify mocks were called
	as.mockUserRepository.AssertExpectations(as.T())
}

func (as *ActionsSuite) TestActions_UpsertOAuthUser_InvalidAuthTokenType() {
	as.setupEnvironment()

	// Setup test data with invalid auth token type
	provider := "google"
	providerUserID := "google-123456"
	email := "test@example.com"
	name := "Test User"
	picture := "https://example.com/avatar.jpg"

	// Setup context with invalid auth token type (string instead of uuid.UUID)
	ctx := context.WithValue(context.Background(), grpc.AuthToken, "invalid-token")

	// Setup request
	request := &connect.Request[proto.UpsertOAuthUserRequest]{
		Msg: &proto.UpsertOAuthUserRequest{
			Provider:       provider,
			ProviderUserId: providerUserID,
			Email:          email,
			Name:           name,
			Picture:        picture,
		},
	}

	// Execute
	actionsInstance := actions.NewActions(ctx, "test-token")
	response, err := actionsInstance.UpsertOAuthUser(ctx, request)

	// Assertions
	as.Error(err)
	as.Nil(response)
	as.Contains(err.Error(), "invalid session")

	// Verify no repository methods were called
	as.mockUserRepository.AssertNotCalled(as.T(), "GetUserWithIdentity")
}

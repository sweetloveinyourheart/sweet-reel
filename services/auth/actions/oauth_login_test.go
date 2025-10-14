package actions_test

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/mock"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/oauth2"
	proto "github.com/sweetloveinyourheart/sweet-reel/proto/code/auth/go"
	userProto "github.com/sweetloveinyourheart/sweet-reel/proto/code/user/go"
	"github.com/sweetloveinyourheart/sweet-reel/services/auth/actions"
)

func (as *ActionsSuite) TestActions_OAuthLogin_Success() {
	as.setupEnvironment()

	providerUserID := uuid.Must(uuid.NewV7())
	userID := uuid.Must(uuid.NewV7())
	email := "test@email.com"
	name := "test"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	as.googleOAuthClient.On("GetUserInfo", mock.Anything, mock.Anything).Return(
		&oauth2.UserInfo{
			ID:    providerUserID.String(),
			Email: email,
			Name:  name,
		}, nil)

	as.userServiceClient.On("UpsertOAuthUser", mock.Anything, mock.Anything).Return(
		&connect.Response[userProto.UpsertOAuthUserResponse]{
			Msg: &userProto.UpsertOAuthUserResponse{
				User: &userProto.User{
					Id:    userID.String(),
					Email: email,
					Name:  name,
				},
				IsNewUser: false,
			},
		}, nil)

	// Setup request
	request := &connect.Request[proto.OAuthLoginRequest]{
		Msg: &proto.OAuthLoginRequest{
			Provider:    string(oauth2.ProviderGoogle),
			AccessToken: "test token",
		},
	}

	// Execute
	actionsInstance := actions.NewActions(ctx, "test-token")
	response, err := actionsInstance.OAuthLogin(ctx, request)

	// Assertions
	as.NoError(err)
	as.NotNil(response)
	as.NotEmpty(response.Msg.GetJwtToken())
}

func (as *ActionsSuite) TestActions_OAuthLogin_ProviderNotSupported() {
	as.setupEnvironment()

	ctx := context.Background()

	// Setup request
	request := &connect.Request[proto.OAuthLoginRequest]{
		Msg: &proto.OAuthLoginRequest{
			Provider:    "invalid provider",
			AccessToken: "test token",
		},
	}

	// Execute
	actionsInstance := actions.NewActions(ctx, "test-token")
	response, err := actionsInstance.OAuthLogin(ctx, request)

	// Assertions
	as.Nil(response)
	as.Error(err, "provider is not supported")
}

func (as *ActionsSuite) TestActions_RefreshToken_Success() {
	as.setupEnvironment()

	signingToken := "test-token"
	userID := uuid.Must(uuid.NewV7())
	email := "test@email.com"
	name := "test"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"type":    "refresh",
		"exp":     time.Now().Add(1 * time.Hour).Unix(), // 1 hour expiration
	})

	refreshTokenString, err := newToken.SignedString([]byte(signingToken))
	as.NoError(err)

	as.userServiceClient.On("GetUserByID", mock.Anything, mock.Anything).Return(
		&connect.Response[userProto.GetUserByIDResponse]{
			Msg: &userProto.GetUserByIDResponse{
				User: &userProto.User{
					Id:    userID.String(),
					Email: email,
					Name:  name,
				},
			},
		}, nil)

	// Setup request
	request := &connect.Request[proto.RefreshTokenRequest]{
		Msg: &proto.RefreshTokenRequest{
			JwtRefreshToken: refreshTokenString,
		},
	}

	// Execute
	actionsInstance := actions.NewActions(ctx, signingToken)
	response, err := actionsInstance.RefreshToken(ctx, request)

	// Assertions
	as.NoError(err)
	as.NotNil(response)
	as.NotEmpty(response.Msg.GetJwtToken())
}

func (as *ActionsSuite) TestActions_RefreshToken_InvalidTokenType() {
	as.setupEnvironment()

	signingToken := "test-token"
	userID := uuid.Must(uuid.NewV7())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(1 * time.Hour).Unix(), // 1 hour expiration
	})

	refreshTokenString, err := newToken.SignedString([]byte(signingToken))
	as.NoError(err)

	// Setup request
	request := &connect.Request[proto.RefreshTokenRequest]{
		Msg: &proto.RefreshTokenRequest{
			JwtRefreshToken: refreshTokenString,
		},
	}

	// Execute
	actionsInstance := actions.NewActions(ctx, signingToken)
	response, err := actionsInstance.RefreshToken(ctx, request)

	// Assertions
	as.Nil(response)
	as.Error(err, "invalid token type")
}

package actions_test

import (
	"context"

	"connectrpc.com/connect"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/grpc"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/oauth2"
	proto "github.com/sweetloveinyourheart/sweet-reel/proto/code/auth/go"
	userProto "github.com/sweetloveinyourheart/sweet-reel/proto/code/user/go"
	"github.com/sweetloveinyourheart/sweet-reel/services/auth/actions"
)

func (as *ActionsSuite) TestActions_OAuthLogin_Success() {
	as.setupEnvironment()

	authUserID := uuid.Must(uuid.NewV7())
	providerUserID := uuid.Must(uuid.NewV7())
	userID := uuid.Must(uuid.NewV7())
	email := "test@email.com"
	name := "test"

	ctx := context.WithValue(context.Background(), grpc.AuthToken, authUserID)

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

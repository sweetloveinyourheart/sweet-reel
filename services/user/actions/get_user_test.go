package actions_test

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/grpc"
	proto "github.com/sweetloveinyourheart/sweet-reel/proto/code/user/go"
	"github.com/sweetloveinyourheart/sweet-reel/services/user/actions"
	"github.com/sweetloveinyourheart/sweet-reel/services/user/models"
)

func (as *ActionsSuite) TestActions_GetUserByID_Success() {
	as.setupEnvironment()

	// Setup test data
	authUserID := uuid.Must(uuid.NewV7())
	userID := uuid.Must(uuid.NewV7())

	// Setup context with auth token
	ctx := context.WithValue(context.Background(), grpc.AuthToken, authUserID)

	as.mockUserRepository.On("GetUserByID", ctx, mock.Anything).Return(
		&models.User{
			ID:        userID,
			Email:     "email@mail.com",
			Name:      "John",
			Picture:   "",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil)

	// Setup request
	request := &connect.Request[proto.GetUserByIDRequest]{
		Msg: &proto.GetUserByIDRequest{
			UserId: userID.String(),
		},
	}

	// Execute
	actionsInstance := actions.NewActions(ctx, "test-token")
	response, err := actionsInstance.GetUserByID(ctx, request)

	// Assertions
	as.NoError(err)
	as.NotNil(response)
}

func (as *ActionsSuite) TestActions_GetUserByID_NotFound() {
	as.setupEnvironment()

	// Setup test data
	authUserID := uuid.Must(uuid.NewV7())
	userID := uuid.Must(uuid.NewV7())

	// Setup context with auth token
	ctx := context.WithValue(context.Background(), grpc.AuthToken, authUserID)

	as.mockUserRepository.On("GetUserByID", ctx, mock.Anything).Return(nil, nil)

	// Setup request
	request := &connect.Request[proto.GetUserByIDRequest]{
		Msg: &proto.GetUserByIDRequest{
			UserId: userID.String(),
		},
	}

	// Execute
	actionsInstance := actions.NewActions(ctx, "test-token")
	response, err := actionsInstance.GetUserByID(ctx, request)

	// Assertions
	as.Nil(response)
	as.Error(err)
}

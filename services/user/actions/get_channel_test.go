package actions_test

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"connectrpc.com/connect"
	"github.com/gofrs/uuid"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/grpc"
	proto "github.com/sweetloveinyourheart/sweet-reel/proto/code/user/go"
	"github.com/sweetloveinyourheart/sweet-reel/services/user/actions"
	"github.com/sweetloveinyourheart/sweet-reel/services/user/models"
)

func (as *ActionsSuite) TestActions_GetChannelByHandle_Success() {
	as.setupEnvironment()

	// Setup test data
	authUserID := uuid.Must(uuid.NewV7())
	channelID := uuid.Must(uuid.NewV7())
	ownerID := uuid.Must(uuid.NewV7())
	handle := "@johndoe1234"

	// Setup context with auth token
	ctx := context.WithValue(context.Background(), grpc.AuthToken, authUserID)

	// Create test channel
	testChannel := &models.Channel{
		ID:              channelID,
		OwnerID:         ownerID,
		Name:            "John Doe",
		Handle:          handle,
		Description:     stringPtr("Test channel description"),
		BannerURL:       stringPtr("https://example.com/banner.jpg"),
		SubscriberCount: 100,
		TotalViews:      5000,
		TotalVideos:     25,
		CreatedAt:       time.Now().Add(-30 * 24 * time.Hour),
		UpdatedAt:       time.Now(),
	}

	// Create test owner
	testOwner := &models.User{
		ID:        ownerID,
		Email:     "john@example.com",
		Name:      "John Doe",
		Picture:   "https://example.com/avatar.jpg",
		CreatedAt: time.Now().Add(-60 * 24 * time.Hour),
		UpdatedAt: time.Now(),
	}

	// Mock repository calls
	as.mockChannelRepository.On("GetChannelByHandle", ctx, handle).Return(testChannel, nil)
	as.mockUserRepository.On("GetUserByID", ctx, ownerID).Return(testOwner, nil)

	// Setup request
	request := &connect.Request[proto.GetChannelByHandleRequest]{
		Msg: &proto.GetChannelByHandleRequest{
			Handle: handle,
		},
	}

	// Execute
	actionsInstance := actions.NewActions(ctx, "test-token")
	response, err := actionsInstance.GetChannelByHandle(ctx, request)

	// Assertions
	as.NoError(err)
	as.NotNil(response)
	as.NotNil(response.Msg)
	as.NotNil(response.Msg.Channel)
	as.NotNil(response.Msg.Owner)
	as.Equal(channelID.String(), response.Msg.Channel.Id)
	as.Equal(handle, response.Msg.Channel.Handle)
	as.Equal("John Doe", response.Msg.Channel.Name)
	as.Equal(int32(100), response.Msg.Channel.SubscriberCount)
	as.Equal(int64(5000), response.Msg.Channel.TotalViews)
	as.Equal(int32(25), response.Msg.Channel.TotalVideos)
	as.Equal(ownerID.String(), response.Msg.Owner.Id)
	as.Equal("john@example.com", response.Msg.Owner.Email)

	// Verify mocks were called
	as.mockChannelRepository.AssertExpectations(as.T())
	as.mockUserRepository.AssertExpectations(as.T())
}

func (as *ActionsSuite) TestActions_GetChannelByHandle_EmptyHandle() {
	as.setupEnvironment()

	// Setup test data
	authUserID := uuid.Must(uuid.NewV7())

	// Setup context with auth token
	ctx := context.WithValue(context.Background(), grpc.AuthToken, authUserID)

	// Setup request with empty handle
	request := &connect.Request[proto.GetChannelByHandleRequest]{
		Msg: &proto.GetChannelByHandleRequest{
			Handle: "",
		},
	}

	// Execute
	actionsInstance := actions.NewActions(ctx, "test-token")
	response, err := actionsInstance.GetChannelByHandle(ctx, request)

	// Assertions
	as.Error(err)
	as.Nil(response)

	// Verify mocks were not called
	as.mockChannelRepository.AssertNotCalled(as.T(), "GetChannelByHandle")
	as.mockUserRepository.AssertNotCalled(as.T(), "GetUserByID")
}

func (as *ActionsSuite) TestActions_GetChannelByHandle_ChannelNotFound() {
	as.setupEnvironment()

	// Setup test data
	authUserID := uuid.Must(uuid.NewV7())
	handle := "@nonexistent"

	// Setup context with auth token
	ctx := context.WithValue(context.Background(), grpc.AuthToken, authUserID)

	// Mock repository to return not found error
	as.mockChannelRepository.On("GetChannelByHandle", ctx, handle).Return(nil, sql.ErrNoRows)

	// Setup request
	request := &connect.Request[proto.GetChannelByHandleRequest]{
		Msg: &proto.GetChannelByHandleRequest{
			Handle: handle,
		},
	}

	// Execute
	actionsInstance := actions.NewActions(ctx, "test-token")
	response, err := actionsInstance.GetChannelByHandle(ctx, request)

	// Assertions
	as.Error(err)
	as.Nil(response)

	// Verify mocks were called
	as.mockChannelRepository.AssertExpectations(as.T())
	as.mockUserRepository.AssertNotCalled(as.T(), "GetUserByID")
}

func (as *ActionsSuite) TestActions_GetChannelByHandle_DatabaseError() {
	as.setupEnvironment()

	// Setup test data
	authUserID := uuid.Must(uuid.NewV7())
	handle := "@johndoe"

	// Setup context with auth token
	ctx := context.WithValue(context.Background(), grpc.AuthToken, authUserID)

	// Mock repository to return database error
	dbError := errors.New("database connection failed")
	as.mockChannelRepository.On("GetChannelByHandle", ctx, handle).Return(nil, dbError)

	// Setup request
	request := &connect.Request[proto.GetChannelByHandleRequest]{
		Msg: &proto.GetChannelByHandleRequest{
			Handle: handle,
		},
	}

	// Execute
	actionsInstance := actions.NewActions(ctx, "test-token")
	response, err := actionsInstance.GetChannelByHandle(ctx, request)

	// Assertions
	as.Error(err)
	as.Nil(response)

	// Verify mocks were called
	as.mockChannelRepository.AssertExpectations(as.T())
	as.mockUserRepository.AssertNotCalled(as.T(), "GetUserByID")
}

func (as *ActionsSuite) TestActions_GetChannelByHandle_OwnerNotFound() {
	as.setupEnvironment()

	// Setup test data
	authUserID := uuid.Must(uuid.NewV7())
	channelID := uuid.Must(uuid.NewV7())
	ownerID := uuid.Must(uuid.NewV7())
	handle := "@johndoe"

	// Setup context with auth token
	ctx := context.WithValue(context.Background(), grpc.AuthToken, authUserID)

	// Create test channel
	testChannel := &models.Channel{
		ID:              channelID,
		OwnerID:         ownerID,
		Name:            "John Doe",
		Handle:          handle,
		SubscriberCount: 0,
		TotalViews:      0,
		TotalVideos:     0,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Mock repository calls - channel found but owner not found
	as.mockChannelRepository.On("GetChannelByHandle", ctx, handle).Return(testChannel, nil)
	as.mockUserRepository.On("GetUserByID", ctx, ownerID).Return(nil, sql.ErrNoRows)

	// Setup request
	request := &connect.Request[proto.GetChannelByHandleRequest]{
		Msg: &proto.GetChannelByHandleRequest{
			Handle: handle,
		},
	}

	// Execute
	actionsInstance := actions.NewActions(ctx, "test-token")
	response, err := actionsInstance.GetChannelByHandle(ctx, request)

	// Assertions
	as.Error(err)
	as.Nil(response)

	// Verify mocks were called
	as.mockChannelRepository.AssertExpectations(as.T())
	as.mockUserRepository.AssertExpectations(as.T())
}

func (as *ActionsSuite) TestActions_GetChannelByHandle_NilDescriptionAndBanner() {
	as.setupEnvironment()

	// Setup test data
	authUserID := uuid.Must(uuid.NewV7())
	channelID := uuid.Must(uuid.NewV7())
	ownerID := uuid.Must(uuid.NewV7())
	handle := "@janedoe"

	// Setup context with auth token
	ctx := context.WithValue(context.Background(), grpc.AuthToken, authUserID)

	// Create test channel with nil optional fields
	testChannel := &models.Channel{
		ID:              channelID,
		OwnerID:         ownerID,
		Name:            "Jane Doe",
		Handle:          handle,
		Description:     nil, // nil description
		BannerURL:       nil, // nil banner
		SubscriberCount: 0,
		TotalViews:      0,
		TotalVideos:     0,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Create test owner
	testOwner := &models.User{
		ID:        ownerID,
		Email:     "jane@example.com",
		Name:      "Jane Doe",
		Picture:   "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock repository calls
	as.mockChannelRepository.On("GetChannelByHandle", ctx, handle).Return(testChannel, nil)
	as.mockUserRepository.On("GetUserByID", ctx, ownerID).Return(testOwner, nil)

	// Setup request
	request := &connect.Request[proto.GetChannelByHandleRequest]{
		Msg: &proto.GetChannelByHandleRequest{
			Handle: handle,
		},
	}

	// Execute
	actionsInstance := actions.NewActions(ctx, "test-token")
	response, err := actionsInstance.GetChannelByHandle(ctx, request)

	// Assertions
	as.NoError(err)
	as.NotNil(response)
	as.NotNil(response.Msg.Channel)
	as.Equal("", response.Msg.Channel.Description) // Should be empty string
	as.Equal("", response.Msg.Channel.BannerUrl)   // Should be empty string

	// Verify mocks were called
	as.mockChannelRepository.AssertExpectations(as.T())
	as.mockUserRepository.AssertExpectations(as.T())
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}

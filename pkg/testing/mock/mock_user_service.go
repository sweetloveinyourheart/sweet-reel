package mock

import (
	"context"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/mock"

	proto "github.com/sweetloveinyourheart/sweet-reel/proto/code/user/go"
	userConnect "github.com/sweetloveinyourheart/sweet-reel/proto/code/user/go/grpcconnect"
)

// Ensure that MockUserServiceClient implements userConnect.UserServiceClient
var _ userConnect.UserServiceClient = (*MockUserServiceClient)(nil)

// MockUserServiceClient is a mock implementation of IOAuthClient.
type MockUserServiceClient struct {
	mock.Mock
}

func (m *MockUserServiceClient) UpsertOAuthUser(
	ctx context.Context,
	req *connect.Request[proto.UpsertOAuthUserRequest],
) (*connect.Response[proto.UpsertOAuthUserResponse], error) {
	args := m.Called(ctx, req)
	resp, _ := args.Get(0).(*connect.Response[proto.UpsertOAuthUserResponse])
	return resp, args.Error(1)
}

func (m *MockUserServiceClient) GetUserByID(
	ctx context.Context,
	req *connect.Request[proto.GetUserByIDRequest],
) (*connect.Response[proto.GetUserByIDResponse], error) {
	args := m.Called(ctx, req)
	resp, _ := args.Get(0).(*connect.Response[proto.GetUserByIDResponse])
	return resp, args.Error(1)
}

func (m *MockUserServiceClient) GetChannelByHandle(
	ctx context.Context,
	req *connect.Request[proto.GetChannelByHandleRequest],
) (*connect.Response[proto.GetChannelByHandleResponse], error) {
	args := m.Called(ctx, req)
	resp, _ := args.Get(0).(*connect.Response[proto.GetChannelByHandleResponse])
	return resp, args.Error(1)
}

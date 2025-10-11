package actions

import (
	"context"
	"errors"

	"github.com/gofrs/uuid"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/grpc"
	proto "github.com/sweetloveinyourheart/sweet-reel/proto/code/user/go"

	"connectrpc.com/connect"
)

func (a *actions) GetUserByID(ctx context.Context, request *connect.Request[proto.GetUserByIDRequest]) (*connect.Response[proto.GetUserByIDResponse], error) {
	userID := uuid.FromStringOrNil(request.Msg.GetUserId())
	if userID == uuid.Nil {
		return nil, grpc.InvalidArgumentError(errors.New("invalid user id"))
	}

	user, err := a.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, grpc.InternalError(err)
	}

	if user == nil {
		return nil, grpc.NotFoundError(errors.New("user not found"))
	}

	response := &proto.GetUserByIDResponse{
		User: &proto.User{
			Id:        user.ID.String(),
			Email:     user.Email,
			Name:      user.Name,
			Picture:   user.Picture,
			CreatedAt: user.CreatedAt.String(),
			UpdatedAt: user.UpdatedAt.String(),
		},
	}
	return connect.NewResponse(response), nil
}

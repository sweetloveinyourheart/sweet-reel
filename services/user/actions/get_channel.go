package actions

import (
	"context"
	"database/sql"
	"errors"

	"connectrpc.com/connect"

	"github.com/gofrs/uuid"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/grpc"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/stringsutil"
	proto "github.com/sweetloveinyourheart/sweet-reel/proto/code/user/go"
)

func (a *actions) GetChannelByHandle(ctx context.Context, request *connect.Request[proto.GetChannelByHandleRequest]) (*connect.Response[proto.GetChannelByHandleResponse], error) {
	handle := request.Msg.GetHandle()
	if handle == "" {
		return nil, grpc.InvalidArgumentError(errors.New("handle is required"))
	}

	// Get channel by handle
	channel, err := a.channelRepo.GetChannelByHandle(ctx, handle)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, grpc.NotFoundError(errors.New("channel not found"))
		}
		return nil, grpc.InternalError(err)
	}

	// Get owner user information
	owner, err := a.userRepo.GetUserByID(ctx, channel.OwnerID)
	if err != nil {
		return nil, grpc.InternalError(err)
	}

	// Build response
	response := &proto.GetChannelByHandleResponse{
		Channel: &proto.Channel{
			Id:              channel.ID.String(),
			OwnerId:         channel.OwnerID.String(),
			Name:            channel.Name,
			Handle:          channel.Handle,
			Description:     stringsutil.GetStringValue(channel.Description),
			BannerUrl:       stringsutil.GetStringValue(channel.BannerURL),
			SubscriberCount: int32(channel.SubscriberCount),
			TotalViews:      channel.TotalViews,
			TotalVideos:     int32(channel.TotalVideos),
			CreatedAt:       channel.CreatedAt.String(),
			UpdatedAt:       channel.UpdatedAt.String(),
		},
		Owner: &proto.User{
			Id:        owner.ID.String(),
			Email:     owner.Email,
			Name:      owner.Name,
			Picture:   owner.Picture,
			CreatedAt: owner.CreatedAt.String(),
			UpdatedAt: owner.UpdatedAt.String(),
		},
	}

	return connect.NewResponse(response), nil
}

func (a *actions) GetChannelByUser(ctx context.Context, request *connect.Request[proto.GetChannelByUserRequest]) (*connect.Response[proto.GetChannelByUserResponse], error) {
	userID := uuid.FromStringOrNil(request.Msg.GetUserId())
	if userID == uuid.Nil {
		return nil, grpc.InvalidArgumentError(errors.New("userID is required"))
	}

	// Get channel by userID
	channel, err := a.channelRepo.GetChannelByOwnerID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, grpc.NotFoundError(errors.New("channel not found"))
		}
		return nil, grpc.InternalError(err)
	}

	// Get owner user information
	owner, err := a.userRepo.GetUserByID(ctx, channel.OwnerID)
	if err != nil {
		return nil, grpc.InternalError(err)
	}

	// Build response
	response := &proto.GetChannelByUserResponse{
		Channel: &proto.Channel{
			Id:              channel.ID.String(),
			OwnerId:         channel.OwnerID.String(),
			Name:            channel.Name,
			Handle:          channel.Handle,
			Description:     stringsutil.GetStringValue(channel.Description),
			BannerUrl:       stringsutil.GetStringValue(channel.BannerURL),
			SubscriberCount: int32(channel.SubscriberCount),
			TotalViews:      channel.TotalViews,
			TotalVideos:     int32(channel.TotalVideos),
			CreatedAt:       channel.CreatedAt.String(),
			UpdatedAt:       channel.UpdatedAt.String(),
		},
		Owner: &proto.User{
			Id:        owner.ID.String(),
			Email:     owner.Email,
			Name:      owner.Name,
			Picture:   owner.Picture,
			CreatedAt: owner.CreatedAt.String(),
			UpdatedAt: owner.UpdatedAt.String(),
		},
	}

	return connect.NewResponse(response), nil
}

func (a *actions) GetChannelByID(ctx context.Context, request *connect.Request[proto.GetChannelByIDRequest]) (*connect.Response[proto.GetChannelByIDResponse], error) {
	videoID := uuid.FromStringOrNil(request.Msg.GetVideoId())
	if videoID == uuid.Nil {
		return nil, grpc.InvalidArgumentError(errors.New("videoID is required"))
	}

	// Get channel by videoID
	channel, err := a.channelRepo.GetChannelByID(ctx, videoID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, grpc.NotFoundError(errors.New("channel not found"))
		}
		return nil, grpc.InternalError(err)
	}

	// Get owner user information
	owner, err := a.userRepo.GetUserByID(ctx, channel.OwnerID)
	if err != nil {
		return nil, grpc.InternalError(err)
	}

	// Build response
	response := &proto.GetChannelByIDResponse{
		Channel: &proto.Channel{
			Id:              channel.ID.String(),
			OwnerId:         channel.OwnerID.String(),
			Name:            channel.Name,
			Handle:          channel.Handle,
			Description:     stringsutil.GetStringValue(channel.Description),
			BannerUrl:       stringsutil.GetStringValue(channel.BannerURL),
			SubscriberCount: int32(channel.SubscriberCount),
			TotalViews:      channel.TotalViews,
			TotalVideos:     int32(channel.TotalVideos),
			CreatedAt:       channel.CreatedAt.String(),
			UpdatedAt:       channel.UpdatedAt.String(),
		},
		Owner: &proto.User{
			Id:        owner.ID.String(),
			Email:     owner.Email,
			Name:      owner.Name,
			Picture:   owner.Picture,
			CreatedAt: owner.CreatedAt.String(),
			UpdatedAt: owner.UpdatedAt.String(),
		},
	}

	return connect.NewResponse(response), nil
}

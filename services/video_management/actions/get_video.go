package actions

import (
	"context"
	"database/sql"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/grpc"
	proto "github.com/sweetloveinyourheart/sweet-reel/proto/code/video_management/go"
)

func (a *actions) GetVideoMetadataById(ctx context.Context, request *connect.Request[proto.GetVideoMetadataByIdRequest]) (*connect.Response[proto.GetVideoMetadataByIdResponse], error) {
	videoID := uuid.FromStringOrNil(request.Msg.GetVideoId())
	if videoID == uuid.Nil {
		return nil, grpc.InvalidArgumentError(errors.Errorf("video id is not recognized, id: ", videoID.String()))
	}

	metadata, err := a.videoAggregateRepo.GetVideoMetadata(ctx, videoID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, grpc.NotFoundError(errors.New("video not found"))
		}

		return nil, grpc.InternalError(err)
	}

	response := &proto.GetVideoMetadataByIdResponse{
		VideoId:            metadata.GetID().String(),
		ChannelId:          metadata.GetChannelID().String(),
		VideoTitle:         metadata.GetTitle(),
		VideoDescription:   metadata.GetDescription(),
		TotalView:          metadata.GetViewCount(),
		ProcessedAt:        metadata.GetCreatedAt().Unix(),
		AvailableQualities: metadata.GetAvailableQualities(),
	}

	return connect.NewResponse(response), nil
}

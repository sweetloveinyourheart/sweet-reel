package actions

import (
	"context"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/grpc"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/s3"
	proto "github.com/sweetloveinyourheart/sweet-reel/proto/code/video_management/go"
	"go.uber.org/zap"
)

func (a *actions) GetUserVideos(ctx context.Context, request *connect.Request[proto.GetUserVideosRequest]) (*connect.Response[proto.GetUserVideosResponse], error) {
	userID := uuid.FromStringOrNil(request.Msg.GetUserId())
	if userID == uuid.Nil {
		return nil, grpc.InvalidArgumentError(errors.Errorf("user id is not recognized, id: ", request.Msg.GetUserId()))
	}

	videosWithThumbnail, err := a.videoAggregateRepo.GetVideosWithThumbnailByUploaderID(ctx, userID, int(request.Msg.GetLimit()), int(request.Msg.Offset))
	if err != nil {
		return nil, grpc.InternalError(err)
	}

	var userVideos []*proto.UserVideo
	for _, video := range videosWithThumbnail {
		if len(video.Thumbnails) == 0 {
			logger.Global().Warn("no video thumbnail was found")
			continue
		}

		firstThumbnail := video.Thumbnails[0]
		thumbnailUrl, err := a.s3Client.GenerateDownloadPublicUri(firstThumbnail.GetObjectKey(), s3.S3VideoProcessedBucket, s3.UrlExpirationSeconds)
		if err != nil {
			logger.Global().Error("unable to generate download url for video thumbnail", zap.Error(err))
			continue
		}

		userVideos = append(userVideos, &proto.UserVideo{
			VideoId:      video.ID.String(),
			VideoTitle:   video.Title,
			ThumbnailUrl: thumbnailUrl,
		})
	}

	response := &proto.GetUserVideosResponse{
		Videos: userVideos,
	}

	return connect.NewResponse(response), nil
}

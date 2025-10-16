package actions

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"

	"go.uber.org/zap"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/grpc"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/s3"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/stringsutil"
	proto "github.com/sweetloveinyourheart/sweet-reel/proto/code/video_management/go"
	"github.com/sweetloveinyourheart/sweet-reel/services/video_management/models"
)

func (a *actions) PresignedUrl(ctx context.Context, request *connect.Request[proto.PresignedUrlRequest]) (*connect.Response[proto.PresignedUrlResponse], error) {
	uploaderID := uuid.FromStringOrNil(request.Msg.GetUploaderId())
	if uploaderID == uuid.Nil {
		return nil, grpc.InvalidArgumentError(errors.Errorf("uploader id is not recognized, id: ", request.Msg.GetUploaderId()))
	}

	filename, ext := s3.ExtractFilenameAndExt(request.Msg.GetFileName())
	if stringsutil.IsBlank(ext) || stringsutil.IsBlank(filename) {
		return nil, grpc.InvalidArgumentError(errors.New("cannot extract the filename or extension"))
	}

	var description *string
	if !stringsutil.IsBlank(request.Msg.Description) {
		description = &request.Msg.Description
	}

	newVideo := models.Video{
		ID:          uuid.Must(uuid.NewV7()),
		Title:       request.Msg.GetTitle(),
		Description: description,
		UploaderID:  uploaderID,
		Status:      models.VideoStatusProcessing,
	}

	if err := newVideo.Validate(); err != nil {
		return nil, grpc.InvalidArgumentError(err)
	}

	key := fmt.Sprintf("%s/%s%s", time.Now().Format("2006-01-02"), newVideo.GetID(), ext)
	url, err := a.s3Client.GenerateUploadPublicUri(key, s3.S3VideoUploadedBucket, s3.UrlExpirationSeconds)
	if err != nil {
		logger.Global().Error("Failed to generate presigned URL",
			zap.String("key", key),
			zap.String("bucket", s3.S3VideoUploadedBucket),
			zap.String("videoID", newVideo.GetID().String()),
			zap.Error(err))
		return nil, grpc.InternalError(errors.Wrapf(err, "failed to generate presigned URL for bucket %s, key %s", s3.S3VideoUploadedBucket, key))
	}

	if err := a.videoRepo.CreateVideo(ctx, &newVideo); err != nil {
		logger.Global().Error("Failed to create video in database",
			zap.String("videoID", newVideo.GetID().String()),
			zap.String("uploaderID", uploaderID.String()),
			zap.Error(err))
		return nil, grpc.InternalError(err)
	}

	response := &proto.PresignedUrlResponse{
		VideoId:      newVideo.GetID().String(),
		PresignedUrl: url,
		ExpiresIn:    s3.UrlExpirationSeconds,
	}

	return connect.NewResponse(response), nil
}

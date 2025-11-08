package actions

import (
	"context"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/ffmpeg"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/grpc"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/s3"
	proto "github.com/sweetloveinyourheart/sweet-reel/proto/code/video_management/go"
)

const (
	QualityDefault = ffmpeg.QualityDefault
)

func (a *actions) ServePlaylist(ctx context.Context, request *connect.Request[proto.ServePlaylistRequest]) (*connect.Response[proto.ServePlaylistResponse], error) {
	videoID := uuid.FromStringOrNil(request.Msg.GetVideoId())
	if videoID == uuid.Nil {
		return nil, grpc.InvalidArgumentError(errors.Errorf("video id is not recognized, id: ", videoID.String()))
	}

	masterPlaylist := ""
	variantPlaylistsMap := make(map[string]*proto.ServePlaylistVariant)

	manifests, err := a.videoAggregateRepo.GetVideoManifestsByVideoID(ctx, videoID)
	if err != nil {
		return nil, grpc.InternalError(err)
	}

	for _, manifest := range manifests {
		if manifest == nil {
			continue
		}

		if manifest.Quality == QualityDefault {
			masterPlaylist, err = a.s3Client.GenerateDownloadPublicUri(manifest.GetObjectKey(), s3.S3VideoUploadedBucket, s3.UrlExpirationSeconds)
			if err != nil {
				logger.Global().Error("unable to generate master playlist url", zap.Error(err))
				return nil, grpc.InternalError(err)
			}
		} else {
			variantPlaylist, err := a.s3Client.GenerateDownloadPublicUri(manifest.GetObjectKey(), s3.S3VideoUploadedBucket, s3.UrlExpirationSeconds)
			if err != nil {
				logger.Global().Error("unable to generate variant playlist url", zap.Error(err))
				return nil, grpc.InternalError(err)
			}

			if _, exists := variantPlaylistsMap[manifest.Quality]; !exists {
				variantPlaylistsMap[manifest.Quality] = &proto.ServePlaylistVariant{}
			}

			variantPlaylistsMap[manifest.Quality].PlaylistUrl = variantPlaylist
			variantPlaylistsMap[manifest.Quality].Quality = manifest.Quality
		}
	}

	variants, err := a.videoAggregateRepo.GetVideoVariantsByVideoID(ctx, videoID)
	if err != nil {
		return nil, grpc.InternalError(err)
	}

	for _, variant := range variants {
		if variant == nil {
			continue
		}

		variantSegment, err := a.s3Client.GenerateDownloadPublicUri(variant.GetObjectKey(), s3.S3VideoUploadedBucket, s3.UrlExpirationSeconds)
		if err != nil {
			logger.Global().Error("unable to generate variant segment url", zap.Error(err))
			return nil, grpc.InternalError(err)
		}

		if _, exists := variantPlaylistsMap[variant.Quality]; !exists {
			variantPlaylistsMap[variant.Quality] = &proto.ServePlaylistVariant{}
		}

		variantPlaylistsMap[variant.Quality].SegmentUrls = append(variantPlaylistsMap[variant.Quality].SegmentUrls, variantSegment)
	}

	variantPlaylists := make([]*proto.ServePlaylistVariant, 0)
	for _, val := range variantPlaylistsMap {
		variantPlaylists = append(variantPlaylists, val)
	}

	response := &proto.ServePlaylistResponse{
		PlaylistUrl: masterPlaylist,
		Variants:    variantPlaylists,
	}

	return connect.NewResponse(response), nil
}

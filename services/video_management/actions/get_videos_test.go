package actions_test

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
	proto "github.com/sweetloveinyourheart/sweet-reel/proto/code/video_management/go"
	"github.com/sweetloveinyourheart/sweet-reel/services/video_management/actions"
	"github.com/sweetloveinyourheart/sweet-reel/services/video_management/models"
)

func (as *ActionsSuite) TestActions_GetUserVideos_Success() {
	as.setupEnvironment()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Setup test data
	userID := uuid.Must(uuid.NewV7())
	videoID := uuid.Must(uuid.NewV7())
	thumbnailUrl := "https://s3.example.com/download-url"

	// Setup mock expectations
	as.mockS3.On("GenerateDownloadPublicUri", mock.Anything, mock.Anything, mock.Anything).Return(thumbnailUrl, nil)

	as.mockVideoAggregateRepository.On("GetUploadedVideos",
		mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]*models.UploadedVideo{
		{
			Video: models.Video{
				ID:    videoID,
				Title: "Test",
			},
			ThumbnailObjectKey: "test",
		},
	}, nil)

	// Setup request
	request := &connect.Request[proto.GetUserVideosRequest]{
		Msg: &proto.GetUserVideosRequest{
			UserId: userID.String(),
		},
	}

	// Execute
	actionsInstance := actions.NewActions(ctx, "test-token")
	response, err := actionsInstance.GetUserVideos(ctx, request)

	// Assertions
	as.NoError(err)
	as.NotNil(response)
}

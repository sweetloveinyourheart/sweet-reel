package actions_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/s3"
	proto "github.com/sweetloveinyourheart/sweet-reel/proto/code/video_management/go"
	"github.com/sweetloveinyourheart/sweet-reel/services/video_management/actions"
	"github.com/sweetloveinyourheart/sweet-reel/services/video_management/models"
)

func (as *ActionsSuite) TestActions_PresignedUrl_Success() {
	as.setupEnvironment()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Setup test data
	userID := uuid.Must(uuid.NewV7())
	channelID := uuid.Must(uuid.NewV7())
	title := "Test Video"
	description := "Test video description"
	fileName := "test-video.mp4"
	expectedURL := "https://s3.example.com/presigned-url"

	// Setup mock expectations
	as.mockS3.On("GenerateUploadPublicUri",
		mock.MatchedBy(func(key string) bool {
			// Key should contain the folder structure and extension
			return key != "" && key[len(key)-4:] == ".mp4"
		}),
		s3.S3VideoUploadedBucket,
		uint32(s3.UrlExpirationSeconds)).Return(expectedURL, nil)

	as.mockVideoAggregateRepository.On("CreateVideo", ctx, mock.MatchedBy(func(video *models.Video) bool {
		return video.Title == title &&
			video.Description != nil && *video.Description == description &&
			video.UploaderID == userID &&
			video.ID != uuid.Nil
	})).Return(nil)

	// Setup request
	request := &connect.Request[proto.PresignedUrlRequest]{
		Msg: &proto.PresignedUrlRequest{
			ChannelId:   channelID.String(),
			Title:       title,
			Description: description,
			FileName:    fileName,
			UploaderId:  userID.String(),
		},
	}

	// Execute
	actionsInstance := actions.NewActions(ctx, "test-token")
	response, err := actionsInstance.PresignedUrl(ctx, request)

	// Assertions
	as.NoError(err)
	as.NotNil(response)
	as.NotEmpty(response.Msg.VideoId)
	as.Equal(expectedURL, response.Msg.PresignedUrl)
	as.Equal(int32(s3.UrlExpirationSeconds), response.Msg.ExpiresIn)

	// Verify the video ID is a valid UUID
	_, err = uuid.FromString(response.Msg.VideoId)
	as.NoError(err)

	// Verify all mocks were called as expected
	as.mockS3.AssertExpectations(as.T())
	as.mockVideoAggregateRepository.AssertExpectations(as.T())
}

func (as *ActionsSuite) TestActions_PresignedUrl_Success_WithoutDescription() {
	as.setupEnvironment()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Setup test data
	channelID := uuid.Must(uuid.NewV7())
	userID := uuid.Must(uuid.NewV7())
	title := "Test Video Without Description"
	fileName := "test-video.mp4"
	expectedURL := "https://s3.example.com/presigned-url"

	// Setup mock expectations
	as.mockS3.On("GenerateUploadPublicUri",
		mock.AnythingOfType("string"),
		s3.S3VideoUploadedBucket,
		uint32(s3.UrlExpirationSeconds)).Return(expectedURL, nil)

	as.mockVideoAggregateRepository.On("CreateVideo", ctx, mock.MatchedBy(func(video *models.Video) bool {
		return video.Title == title &&
			video.Description == nil &&
			video.UploaderID == userID
	})).Return(nil)

	// Setup request without description
	request := &connect.Request[proto.PresignedUrlRequest]{
		Msg: &proto.PresignedUrlRequest{
			ChannelId:  channelID.String(),
			Title:      title,
			FileName:   fileName,
			UploaderId: userID.String(),
		},
	}

	// Execute
	actionsInstance := actions.NewActions(ctx, "test-token")
	response, err := actionsInstance.PresignedUrl(ctx, request)

	// Assertions
	as.NoError(err)
	as.NotNil(response)
	as.NotEmpty(response.Msg.VideoId)
	as.Equal(expectedURL, response.Msg.PresignedUrl)

	// Verify all mocks were called as expected
	as.mockS3.AssertExpectations(as.T())
	as.mockVideoAggregateRepository.AssertExpectations(as.T())
}

func (as *ActionsSuite) TestActions_PresignedUrl_InvalidFileName_NoExtension() {
	as.setupEnvironment()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Setup request with filename without extension
	request := &connect.Request[proto.PresignedUrlRequest]{
		Msg: &proto.PresignedUrlRequest{
			Title:    "Test Video",
			FileName: "test-video", // No extension
		},
	}

	// Execute
	actionsInstance := actions.NewActions(ctx, "test-token")
	response, err := actionsInstance.PresignedUrl(ctx, request)

	// Assertions
	as.Error(err)
	as.Nil(response)

	// Verify error is invalid argument
	var connectErr *connect.Error
	as.True(errors.As(err, &connectErr))
	as.Equal(connect.CodeInvalidArgument, connectErr.Code())

	// Verify no mocks were called
	as.mockS3.AssertNotCalled(as.T(), "GenerateUploadPublicUri")
	as.mockVideoAggregateRepository.AssertNotCalled(as.T(), "CreateVideo")
}

func (as *ActionsSuite) TestActions_PresignedUrl_InvalidFileName_Empty() {
	as.setupEnvironment()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Setup request with empty filename
	request := &connect.Request[proto.PresignedUrlRequest]{
		Msg: &proto.PresignedUrlRequest{
			Title:    "Test Video",
			FileName: "",
		},
	}

	// Execute
	actionsInstance := actions.NewActions(ctx, "test-token")
	response, err := actionsInstance.PresignedUrl(ctx, request)

	// Assertions
	as.Error(err)
	as.Nil(response)

	// Verify error is invalid argument
	var connectErr *connect.Error
	as.True(errors.As(err, &connectErr))
	as.Equal(connect.CodeInvalidArgument, connectErr.Code())

	// Verify no mocks were called
	as.mockS3.AssertNotCalled(as.T(), "GenerateUploadPublicUri")
	as.mockVideoAggregateRepository.AssertNotCalled(as.T(), "CreateVideo")
}

func (as *ActionsSuite) TestActions_PresignedUrl_VideoValidationError() {
	as.setupEnvironment()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Setup request with empty title (will cause validation error)
	request := &connect.Request[proto.PresignedUrlRequest]{
		Msg: &proto.PresignedUrlRequest{
			Title:    "", // Empty title will cause validation error
			FileName: "test-video.mp4",
		},
	}

	// Execute
	actionsInstance := actions.NewActions(ctx, "test-token")
	response, err := actionsInstance.PresignedUrl(ctx, request)

	// Assertions
	as.Error(err)
	as.Nil(response)

	// Verify error is invalid argument
	var connectErr *connect.Error
	as.True(errors.As(err, &connectErr))
	as.Equal(connect.CodeInvalidArgument, connectErr.Code())

	// Verify no mocks were called
	as.mockS3.AssertNotCalled(as.T(), "GenerateUploadPublicUri")
	as.mockVideoAggregateRepository.AssertNotCalled(as.T(), "CreateVideo")
}

func (as *ActionsSuite) TestActions_PresignedUrl_S3GenerateUrlError() {
	as.setupEnvironment()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Setup test data
	userID := uuid.Must(uuid.NewV7())
	channelID := uuid.Must(uuid.NewV7())
	title := "Test Video"
	fileName := "test-video.mp4"
	s3Error := errors.New("S3 service unavailable")

	// Setup mock expectations - S3 fails
	as.mockS3.On("GenerateUploadPublicUri",
		mock.AnythingOfType("string"),
		s3.S3VideoUploadedBucket,
		uint32(s3.UrlExpirationSeconds)).Return("", s3Error)

	// Setup request
	request := &connect.Request[proto.PresignedUrlRequest]{
		Msg: &proto.PresignedUrlRequest{
			ChannelId:  channelID.String(),
			Title:      title,
			FileName:   fileName,
			UploaderId: userID.String(),
		},
	}

	// Execute
	actionsInstance := actions.NewActions(ctx, "test-token")
	response, err := actionsInstance.PresignedUrl(ctx, request)

	// Assertions
	as.Error(err)
	as.Nil(response)

	// Verify error is internal error
	var connectErr *connect.Error
	as.True(errors.As(err, &connectErr))
	as.Equal(connect.CodeInternal, connectErr.Code())

	// Verify S3 was called but video repository was not
	as.mockS3.AssertExpectations(as.T())
	as.mockVideoAggregateRepository.AssertNotCalled(as.T(), "CreateVideo")
}

func (as *ActionsSuite) TestActions_PresignedUrl_DatabaseError() {
	as.setupEnvironment()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Setup test data
	channelID := uuid.Must(uuid.NewV7())
	userID := uuid.Must(uuid.NewV7())
	title := "Test Video"
	fileName := "test-video.mp4"
	expectedURL := "https://s3.example.com/presigned-url"
	dbError := errors.New("database connection failed")

	// Setup mock expectations
	as.mockS3.On("GenerateUploadPublicUri",
		mock.AnythingOfType("string"),
		s3.S3VideoUploadedBucket,
		uint32(s3.UrlExpirationSeconds)).Return(expectedURL, nil)

	as.mockVideoAggregateRepository.On("CreateVideo", ctx, mock.AnythingOfType("*models.Video")).Return(dbError)

	// Setup request
	request := &connect.Request[proto.PresignedUrlRequest]{
		Msg: &proto.PresignedUrlRequest{
			ChannelId:  channelID.String(),
			Title:      title,
			FileName:   fileName,
			UploaderId: userID.String(),
		},
	}

	// Execute
	actionsInstance := actions.NewActions(ctx, "test-token")
	response, err := actionsInstance.PresignedUrl(ctx, request)

	// Assertions
	as.Error(err)
	as.Nil(response)

	// Verify error is internal error
	var connectErr *connect.Error
	as.True(errors.As(err, &connectErr))
	as.Equal(connect.CodeInternal, connectErr.Code())

	// Verify all mocks were called as expected
	as.mockS3.AssertExpectations(as.T())
	as.mockVideoAggregateRepository.AssertExpectations(as.T())
}

func (as *ActionsSuite) TestActions_PresignedUrl_KeyGeneration() {
	as.setupEnvironment()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Setup test data
	channelID := uuid.Must(uuid.NewV7())
	userID := uuid.Must(uuid.NewV7())
	title := "Test Video"
	fileName := "test-video.mp4"
	expectedURL := "https://s3.example.com/presigned-url"

	// Capture the generated key for validation
	var capturedKey string
	as.mockS3.On("GenerateUploadPublicUri",
		mock.MatchedBy(func(key string) bool {
			capturedKey = key
			return true
		}),
		s3.S3VideoUploadedBucket,
		uint32(s3.UrlExpirationSeconds)).Return(expectedURL, nil)

	as.mockVideoAggregateRepository.On("CreateVideo", ctx, mock.AnythingOfType("*models.Video")).Return(nil)

	// Setup request
	request := &connect.Request[proto.PresignedUrlRequest]{
		Msg: &proto.PresignedUrlRequest{
			ChannelId:  channelID.String(),
			Title:      title,
			FileName:   fileName,
			UploaderId: userID.String(),
		},
	}

	// Execute
	actionsInstance := actions.NewActions(ctx, "test-token")
	response, err := actionsInstance.PresignedUrl(ctx, request)

	// Assertions
	as.NoError(err)
	as.NotNil(response)

	// Validate key format: raw/YYYY/MM/DD/{uuid}.mp4
	as.Contains(capturedKey, time.Now().Format("2006-01-02"))
	as.Contains(capturedKey, ".mp4")

	// Verify all mocks were called as expected
	as.mockS3.AssertExpectations(as.T())
	as.mockVideoAggregateRepository.AssertExpectations(as.T())
}

func (as *ActionsSuite) TestActions_PresignedUrl_DifferentFileExtensions() {
	as.setupEnvironment()

	testCases := []struct {
		fileName      string
		expectedExt   string
		shouldSucceed bool
	}{
		{"video.mp4", ".mp4", true},
		{"video.avi", ".avi", true},
		{"video.mov", ".mov", true},
		{"video.mkv", ".mkv", true},
		{"video.webm", ".webm", true},
		{"video", "", false},    // No extension
		{".mp4", ".mp4", false}, // No filename
	}

	for _, tc := range testCases {
		as.T().Run(fmt.Sprintf("FileName_%s", tc.fileName), func(t *testing.T) {
			// Reset mocks for each test case
			as.mockS3.ExpectedCalls = nil
			as.mockVideoAggregateRepository.ExpectedCalls = nil

			channelID := uuid.Must(uuid.NewV7())
			userID := uuid.Must(uuid.NewV7())

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			if tc.shouldSucceed {
				// Setup successful mocks
				as.mockS3.On("GenerateUploadPublicUri",
					mock.MatchedBy(func(key string) bool {
						return key != "" && key[len(key)-len(tc.expectedExt):] == tc.expectedExt
					}),
					s3.S3VideoUploadedBucket,
					uint32(s3.UrlExpirationSeconds)).Return("https://example.com/url", nil)

				as.mockVideoAggregateRepository.On("CreateVideo", ctx, mock.AnythingOfType("*models.Video")).Return(nil)
			}

			request := &connect.Request[proto.PresignedUrlRequest]{
				Msg: &proto.PresignedUrlRequest{
					ChannelId:  channelID.String(),
					Title:      "Test Video",
					FileName:   tc.fileName,
					UploaderId: userID.String(),
				},
			}

			actionsInstance := actions.NewActions(ctx, "test-token")
			response, err := actionsInstance.PresignedUrl(ctx, request)

			if tc.shouldSucceed {
				as.NoError(err)
				as.NotNil(response)
				as.mockS3.AssertExpectations(t)
				as.mockVideoAggregateRepository.AssertExpectations(t)
			} else {
				as.Error(err)
				as.Nil(response)

				var connectErr *connect.Error
				as.True(errors.As(err, &connectErr))
				as.Equal(connect.CodeInvalidArgument, connectErr.Code())
			}
		})
	}
}

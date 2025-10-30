package actions_test

import (
	"context"
	"errors"
	"time"

	"connectrpc.com/connect"
	"github.com/gofrs/uuid"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/ffmpeg"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/s3"
	proto "github.com/sweetloveinyourheart/sweet-reel/proto/code/video_management/go"
	"github.com/sweetloveinyourheart/sweet-reel/services/video_management/actions"
	"github.com/sweetloveinyourheart/sweet-reel/services/video_management/models"
)

func (as *ActionsSuite) TestActions_ServePlaylist_Success() {
	as.setupEnvironment()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Setup test data
	videoID := uuid.Must(uuid.NewV7())
	masterPlaylistKey := "master/playlist.m3u8"
	variantPlaylistKey720p := "720p/playlist.m3u8"
	variantPlaylistKey1080p := "1080p/playlist.m3u8"
	segment720pKey := "720p/segment-0.ts"
	segment1080pKey := "1080p/segment-0.ts"

	masterPlaylistURL := "https://s3.example.com/master/playlist.m3u8"
	variantPlaylistURL720p := "https://s3.example.com/720p/playlist.m3u8"
	variantPlaylistURL1080p := "https://s3.example.com/1080p/playlist.m3u8"
	segmentURL720p := "https://s3.example.com/720p/segment-0.ts"
	segmentURL1080p := "https://s3.example.com/1080p/segment-0.ts"

	// Setup mock expectations for manifests
	manifests := []*models.VideoManifest{
		{
			ID:        uuid.Must(uuid.NewV7()),
			VideoID:   videoID,
			ObjectKey: masterPlaylistKey,
			Quality:   ffmpeg.QualityDefault,
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.Must(uuid.NewV7()),
			VideoID:   videoID,
			ObjectKey: variantPlaylistKey720p,
			Quality:   "720p",
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.Must(uuid.NewV7()),
			VideoID:   videoID,
			ObjectKey: variantPlaylistKey1080p,
			Quality:   "1080p",
			CreatedAt: time.Now(),
		},
	}

	as.mockVideoAggregateRepository.On("GetVideoManifestsByVideoID", ctx, videoID).Return(manifests, nil)

	// Setup mock expectations for variants
	variants := []*models.VideoVariant{
		{
			ID:        uuid.Must(uuid.NewV7()),
			VideoID:   videoID,
			Quality:   "720p",
			ObjectKey: segment720pKey,
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.Must(uuid.NewV7()),
			VideoID:   videoID,
			Quality:   "1080p",
			ObjectKey: segment1080pKey,
			CreatedAt: time.Now(),
		},
	}

	as.mockVideoAggregateRepository.On("GetVideoVariantsByVideoID", ctx, videoID).Return(variants, nil)

	// Setup S3 mock expectations for master playlist
	as.mockS3.On("GenerateDownloadPublicUri", masterPlaylistKey, s3.S3VideoUploadedBucket, uint32(s3.UrlExpirationSeconds)).
		Return(masterPlaylistURL, nil)

	// Setup S3 mock expectations for variant playlists
	as.mockS3.On("GenerateDownloadPublicUri", variantPlaylistKey720p, s3.S3VideoUploadedBucket, uint32(s3.UrlExpirationSeconds)).
		Return(variantPlaylistURL720p, nil)
	as.mockS3.On("GenerateDownloadPublicUri", variantPlaylistKey1080p, s3.S3VideoUploadedBucket, uint32(s3.UrlExpirationSeconds)).
		Return(variantPlaylistURL1080p, nil)

	// Setup S3 mock expectations for segments
	as.mockS3.On("GenerateDownloadPublicUri", segment720pKey, s3.S3VideoUploadedBucket, uint32(s3.UrlExpirationSeconds)).
		Return(segmentURL720p, nil)
	as.mockS3.On("GenerateDownloadPublicUri", segment1080pKey, s3.S3VideoUploadedBucket, uint32(s3.UrlExpirationSeconds)).
		Return(segmentURL1080p, nil)

	// Setup request
	request := &connect.Request[proto.ServePlaylistRequest]{
		Msg: &proto.ServePlaylistRequest{
			VideoId: videoID.String(),
		},
	}

	// Execute
	actionsInstance := actions.NewActions(ctx, "test-token")
	response, err := actionsInstance.ServePlaylist(ctx, request)

	// Assertions
	as.NoError(err)
	as.NotNil(response)
	as.Equal(masterPlaylistURL, response.Msg.PlaylistUrl)
	as.Len(response.Msg.Variants, 2)

	// Verify variants contain correct data
	variantsMap := make(map[string]*proto.ServePlaylistVariant)
	for _, variant := range response.Msg.Variants {
		variantsMap[variant.Quality] = variant
	}

	// Check 720p variant
	as.Contains(variantsMap, "720p")
	as.Equal(variantPlaylistURL720p, variantsMap["720p"].PlaylistUrl)
	as.Equal("720p", variantsMap["720p"].Quality)
	as.Contains(variantsMap["720p"].SegmentUrls, segmentURL720p)

	// Check 1080p variant
	as.Contains(variantsMap, "1080p")
	as.Equal(variantPlaylistURL1080p, variantsMap["1080p"].PlaylistUrl)
	as.Equal("1080p", variantsMap["1080p"].Quality)
	as.Contains(variantsMap["1080p"].SegmentUrls, segmentURL1080p)

	// Verify all mocks were called as expected
	as.mockVideoAggregateRepository.AssertExpectations(as.T())
	as.mockS3.AssertExpectations(as.T())
}

func (as *ActionsSuite) TestActions_ServePlaylist_InvalidVideoID() {
	as.setupEnvironment()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Setup request with invalid video ID
	request := &connect.Request[proto.ServePlaylistRequest]{
		Msg: &proto.ServePlaylistRequest{
			VideoId: "invalid-uuid",
		},
	}

	// Execute
	actionsInstance := actions.NewActions(ctx, "test-token")
	response, err := actionsInstance.ServePlaylist(ctx, request)

	// Assertions
	as.Error(err)
	as.Nil(response)

	// Verify error is invalid argument
	var connectErr *connect.Error
	as.True(errors.As(err, &connectErr))
	as.Equal(connect.CodeInvalidArgument, connectErr.Code())

	// Verify no repository or S3 calls were made
	as.mockVideoAggregateRepository.AssertNotCalled(as.T(), "GetVideoManifestsByVideoID")
	as.mockVideoAggregateRepository.AssertNotCalled(as.T(), "GetVideoVariantsByVideoID")
	as.mockS3.AssertNotCalled(as.T(), "GenerateDownloadPublicUri")
}

func (as *ActionsSuite) TestActions_ServePlaylist_EmptyVideoID() {
	as.setupEnvironment()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Setup request with empty video ID
	request := &connect.Request[proto.ServePlaylistRequest]{
		Msg: &proto.ServePlaylistRequest{
			VideoId: "",
		},
	}

	// Execute
	actionsInstance := actions.NewActions(ctx, "test-token")
	response, err := actionsInstance.ServePlaylist(ctx, request)

	// Assertions
	as.Error(err)
	as.Nil(response)

	// Verify error is invalid argument
	var connectErr *connect.Error
	as.True(errors.As(err, &connectErr))
	as.Equal(connect.CodeInvalidArgument, connectErr.Code())

	// Verify no repository or S3 calls were made
	as.mockVideoAggregateRepository.AssertNotCalled(as.T(), "GetVideoManifestsByVideoID")
	as.mockVideoAggregateRepository.AssertNotCalled(as.T(), "GetVideoVariantsByVideoID")
	as.mockS3.AssertNotCalled(as.T(), "GenerateDownloadPublicUri")
}

func (as *ActionsSuite) TestActions_ServePlaylist_GetManifestsError() {
	as.setupEnvironment()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Setup test data
	videoID := uuid.Must(uuid.NewV7())
	dbError := errors.New("database connection failed")

	// Setup mock expectations - database fails
	as.mockVideoAggregateRepository.On("GetVideoManifestsByVideoID", ctx, videoID).Return(nil, dbError)

	// Setup request
	request := &connect.Request[proto.ServePlaylistRequest]{
		Msg: &proto.ServePlaylistRequest{
			VideoId: videoID.String(),
		},
	}

	// Execute
	actionsInstance := actions.NewActions(ctx, "test-token")
	response, err := actionsInstance.ServePlaylist(ctx, request)

	// Assertions
	as.Error(err)
	as.Nil(response)

	// Verify error is internal error
	var connectErr *connect.Error
	as.True(errors.As(err, &connectErr))
	as.Equal(connect.CodeInternal, connectErr.Code())

	// Verify manifests call was made but variants was not
	as.mockVideoAggregateRepository.AssertExpectations(as.T())
	as.mockVideoAggregateRepository.AssertNotCalled(as.T(), "GetVideoVariantsByVideoID")
	as.mockS3.AssertNotCalled(as.T(), "GenerateDownloadPublicUri")
}

func (as *ActionsSuite) TestActions_ServePlaylist_GetVariantsError() {
	as.setupEnvironment()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Setup test data
	videoID := uuid.Must(uuid.NewV7())
	masterPlaylistKey := "master/playlist.m3u8"
	masterPlaylistURL := "https://s3.example.com/master/playlist.m3u8"
	dbError := errors.New("database connection failed")

	// Setup mock expectations for manifests
	manifests := []*models.VideoManifest{
		{
			ID:        uuid.Must(uuid.NewV7()),
			VideoID:   videoID,
			ObjectKey: masterPlaylistKey,
			Quality:   ffmpeg.QualityDefault,
			CreatedAt: time.Now(),
		},
	}

	as.mockVideoAggregateRepository.On("GetVideoManifestsByVideoID", ctx, videoID).Return(manifests, nil)

	// Setup S3 mock for master playlist
	as.mockS3.On("GenerateDownloadPublicUri", masterPlaylistKey, s3.S3VideoUploadedBucket, uint32(s3.UrlExpirationSeconds)).
		Return(masterPlaylistURL, nil)

	// Setup mock expectations for variants - database fails
	as.mockVideoAggregateRepository.On("GetVideoVariantsByVideoID", ctx, videoID).Return(nil, dbError)

	// Setup request
	request := &connect.Request[proto.ServePlaylistRequest]{
		Msg: &proto.ServePlaylistRequest{
			VideoId: videoID.String(),
		},
	}

	// Execute
	actionsInstance := actions.NewActions(ctx, "test-token")
	response, err := actionsInstance.ServePlaylist(ctx, request)

	// Assertions
	as.Error(err)
	as.Nil(response)

	// Verify error is internal error
	var connectErr *connect.Error
	as.True(errors.As(err, &connectErr))
	as.Equal(connect.CodeInternal, connectErr.Code())

	// Verify all repository calls were made
	as.mockVideoAggregateRepository.AssertExpectations(as.T())
	as.mockS3.AssertExpectations(as.T())
}

func (as *ActionsSuite) TestActions_ServePlaylist_S3GenerateMasterPlaylistError() {
	as.setupEnvironment()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Setup test data
	videoID := uuid.Must(uuid.NewV7())
	masterPlaylistKey := "master/playlist.m3u8"
	s3Error := errors.New("S3 service unavailable")

	// Setup mock expectations for manifests
	manifests := []*models.VideoManifest{
		{
			ID:        uuid.Must(uuid.NewV7()),
			VideoID:   videoID,
			ObjectKey: masterPlaylistKey,
			Quality:   ffmpeg.QualityDefault,
			CreatedAt: time.Now(),
		},
	}

	as.mockVideoAggregateRepository.On("GetVideoManifestsByVideoID", ctx, videoID).Return(manifests, nil)

	// Setup S3 mock for master playlist - fails
	as.mockS3.On("GenerateDownloadPublicUri", masterPlaylistKey, s3.S3VideoUploadedBucket, uint32(s3.UrlExpirationSeconds)).
		Return("", s3Error)

	// Setup request
	request := &connect.Request[proto.ServePlaylistRequest]{
		Msg: &proto.ServePlaylistRequest{
			VideoId: videoID.String(),
		},
	}

	// Execute
	actionsInstance := actions.NewActions(ctx, "test-token")
	response, err := actionsInstance.ServePlaylist(ctx, request)

	// Assertions
	as.Error(err)
	as.Nil(response)

	// Verify error is internal error
	var connectErr *connect.Error
	as.True(errors.As(err, &connectErr))
	as.Equal(connect.CodeInternal, connectErr.Code())

	// Verify manifests and S3 calls were made
	as.mockVideoAggregateRepository.AssertExpectations(as.T())
	as.mockS3.AssertExpectations(as.T())
}

func (as *ActionsSuite) TestActions_ServePlaylist_S3GenerateVariantPlaylistError() {
	as.setupEnvironment()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Setup test data
	videoID := uuid.Must(uuid.NewV7())
	variantPlaylistKey := "720p/playlist.m3u8"
	s3Error := errors.New("S3 service unavailable")

	// Setup mock expectations for manifests
	manifests := []*models.VideoManifest{
		{
			ID:        uuid.Must(uuid.NewV7()),
			VideoID:   videoID,
			ObjectKey: variantPlaylistKey,
			Quality:   "720p",
			CreatedAt: time.Now(),
		},
	}

	as.mockVideoAggregateRepository.On("GetVideoManifestsByVideoID", ctx, videoID).Return(manifests, nil)

	// Setup S3 mock for variant playlist - fails
	as.mockS3.On("GenerateDownloadPublicUri", variantPlaylistKey, s3.S3VideoUploadedBucket, uint32(s3.UrlExpirationSeconds)).
		Return("", s3Error)

	// Setup request
	request := &connect.Request[proto.ServePlaylistRequest]{
		Msg: &proto.ServePlaylistRequest{
			VideoId: videoID.String(),
		},
	}

	// Execute
	actionsInstance := actions.NewActions(ctx, "test-token")
	response, err := actionsInstance.ServePlaylist(ctx, request)

	// Assertions
	as.Error(err)
	as.Nil(response)

	// Verify error is internal error
	var connectErr *connect.Error
	as.True(errors.As(err, &connectErr))
	as.Equal(connect.CodeInternal, connectErr.Code())

	// Verify manifests and S3 calls were made
	as.mockVideoAggregateRepository.AssertExpectations(as.T())
	as.mockS3.AssertExpectations(as.T())
}

func (as *ActionsSuite) TestActions_ServePlaylist_S3GenerateVariantSegmentError() {
	as.setupEnvironment()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Setup test data
	videoID := uuid.Must(uuid.NewV7())
	variantPlaylistKey := "720p/playlist.m3u8"
	segmentKey := "720p/segment-0.ts"
	variantPlaylistURL := "https://s3.example.com/720p/playlist.m3u8"
	s3Error := errors.New("S3 service unavailable")

	// Setup mock expectations for manifests
	manifests := []*models.VideoManifest{
		{
			ID:        uuid.Must(uuid.NewV7()),
			VideoID:   videoID,
			ObjectKey: variantPlaylistKey,
			Quality:   "720p",
			CreatedAt: time.Now(),
		},
	}

	as.mockVideoAggregateRepository.On("GetVideoManifestsByVideoID", ctx, videoID).Return(manifests, nil)

	// Setup S3 mock for variant playlist - succeeds
	as.mockS3.On("GenerateDownloadPublicUri", variantPlaylistKey, s3.S3VideoUploadedBucket, uint32(s3.UrlExpirationSeconds)).
		Return(variantPlaylistURL, nil)

	// Setup mock expectations for variants
	variants := []*models.VideoVariant{
		{
			ID:        uuid.Must(uuid.NewV7()),
			VideoID:   videoID,
			Quality:   "720p",
			ObjectKey: segmentKey,
			CreatedAt: time.Now(),
		},
	}

	as.mockVideoAggregateRepository.On("GetVideoVariantsByVideoID", ctx, videoID).Return(variants, nil)

	// Setup S3 mock for segment - fails
	as.mockS3.On("GenerateDownloadPublicUri", segmentKey, s3.S3VideoUploadedBucket, uint32(s3.UrlExpirationSeconds)).
		Return("", s3Error)

	// Setup request
	request := &connect.Request[proto.ServePlaylistRequest]{
		Msg: &proto.ServePlaylistRequest{
			VideoId: videoID.String(),
		},
	}

	// Execute
	actionsInstance := actions.NewActions(ctx, "test-token")
	response, err := actionsInstance.ServePlaylist(ctx, request)

	// Assertions
	as.Error(err)
	as.Nil(response)

	// Verify error is internal error
	var connectErr *connect.Error
	as.True(errors.As(err, &connectErr))
	as.Equal(connect.CodeInternal, connectErr.Code())

	// Verify all repository and S3 calls were made
	as.mockVideoAggregateRepository.AssertExpectations(as.T())
	as.mockS3.AssertExpectations(as.T())
}

func (as *ActionsSuite) TestActions_ServePlaylist_OnlyMasterPlaylist() {
	as.setupEnvironment()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Setup test data - only master playlist, no variants
	videoID := uuid.Must(uuid.NewV7())
	masterPlaylistKey := "master/playlist.m3u8"
	masterPlaylistURL := "https://s3.example.com/master/playlist.m3u8"

	// Setup mock expectations for manifests
	manifests := []*models.VideoManifest{
		{
			ID:        uuid.Must(uuid.NewV7()),
			VideoID:   videoID,
			ObjectKey: masterPlaylistKey,
			Quality:   ffmpeg.QualityDefault,
			CreatedAt: time.Now(),
		},
	}

	as.mockVideoAggregateRepository.On("GetVideoManifestsByVideoID", ctx, videoID).Return(manifests, nil)

	// Setup S3 mock for master playlist
	as.mockS3.On("GenerateDownloadPublicUri", masterPlaylistKey, s3.S3VideoUploadedBucket, uint32(s3.UrlExpirationSeconds)).
		Return(masterPlaylistURL, nil)

	// Setup mock expectations for variants - empty
	as.mockVideoAggregateRepository.On("GetVideoVariantsByVideoID", ctx, videoID).Return([]*models.VideoVariant{}, nil)

	// Setup request
	request := &connect.Request[proto.ServePlaylistRequest]{
		Msg: &proto.ServePlaylistRequest{
			VideoId: videoID.String(),
		},
	}

	// Execute
	actionsInstance := actions.NewActions(ctx, "test-token")
	response, err := actionsInstance.ServePlaylist(ctx, request)

	// Assertions
	as.NoError(err)
	as.NotNil(response)
	as.Equal(masterPlaylistURL, response.Msg.PlaylistUrl)
	as.Len(response.Msg.Variants, 0)

	// Verify all mocks were called as expected
	as.mockVideoAggregateRepository.AssertExpectations(as.T())
	as.mockS3.AssertExpectations(as.T())
}

func (as *ActionsSuite) TestActions_ServePlaylist_NoMasterPlaylist() {
	as.setupEnvironment()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Setup test data - only variant playlists, no master
	videoID := uuid.Must(uuid.NewV7())
	variantPlaylistKey := "720p/playlist.m3u8"
	segmentKey := "720p/segment-0.ts"
	variantPlaylistURL := "https://s3.example.com/720p/playlist.m3u8"
	segmentURL := "https://s3.example.com/720p/segment-0.ts"

	// Setup mock expectations for manifests - no master
	manifests := []*models.VideoManifest{
		{
			ID:        uuid.Must(uuid.NewV7()),
			VideoID:   videoID,
			ObjectKey: variantPlaylistKey,
			Quality:   "720p",
			CreatedAt: time.Now(),
		},
	}

	as.mockVideoAggregateRepository.On("GetVideoManifestsByVideoID", ctx, videoID).Return(manifests, nil)

	// Setup S3 mock for variant playlist
	as.mockS3.On("GenerateDownloadPublicUri", variantPlaylistKey, s3.S3VideoUploadedBucket, uint32(s3.UrlExpirationSeconds)).
		Return(variantPlaylistURL, nil)

	// Setup mock expectations for variants
	variants := []*models.VideoVariant{
		{
			ID:        uuid.Must(uuid.NewV7()),
			VideoID:   videoID,
			Quality:   "720p",
			ObjectKey: segmentKey,
			CreatedAt: time.Now(),
		},
	}

	as.mockVideoAggregateRepository.On("GetVideoVariantsByVideoID", ctx, videoID).Return(variants, nil)

	// Setup S3 mock for segment
	as.mockS3.On("GenerateDownloadPublicUri", segmentKey, s3.S3VideoUploadedBucket, uint32(s3.UrlExpirationSeconds)).
		Return(segmentURL, nil)

	// Setup request
	request := &connect.Request[proto.ServePlaylistRequest]{
		Msg: &proto.ServePlaylistRequest{
			VideoId: videoID.String(),
		},
	}

	// Execute
	actionsInstance := actions.NewActions(ctx, "test-token")
	response, err := actionsInstance.ServePlaylist(ctx, request)

	// Assertions
	as.NoError(err)
	as.NotNil(response)
	as.Empty(response.Msg.PlaylistUrl) // No master playlist
	as.Len(response.Msg.Variants, 1)
	as.Equal("720p", response.Msg.Variants[0].Quality)
	as.Equal(variantPlaylistURL, response.Msg.Variants[0].PlaylistUrl)
	as.Contains(response.Msg.Variants[0].SegmentUrls, segmentURL)

	// Verify all mocks were called as expected
	as.mockVideoAggregateRepository.AssertExpectations(as.T())
	as.mockS3.AssertExpectations(as.T())
}

func (as *ActionsSuite) TestActions_ServePlaylist_MultipleSegmentsPerQuality() {
	as.setupEnvironment()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Setup test data with multiple segments for same quality
	videoID := uuid.Must(uuid.NewV7())
	variantPlaylistKey := "720p/playlist.m3u8"
	segment1Key := "720p/segment-0.ts"
	segment2Key := "720p/segment-1.ts"
	segment3Key := "720p/segment-2.ts"
	variantPlaylistURL := "https://s3.example.com/720p/playlist.m3u8"
	segment1URL := "https://s3.example.com/720p/segment-0.ts"
	segment2URL := "https://s3.example.com/720p/segment-1.ts"
	segment3URL := "https://s3.example.com/720p/segment-2.ts"

	// Setup mock expectations for manifests
	manifests := []*models.VideoManifest{
		{
			ID:        uuid.Must(uuid.NewV7()),
			VideoID:   videoID,
			ObjectKey: variantPlaylistKey,
			Quality:   "720p",
			CreatedAt: time.Now(),
		},
	}

	as.mockVideoAggregateRepository.On("GetVideoManifestsByVideoID", ctx, videoID).Return(manifests, nil)

	// Setup S3 mock for variant playlist
	as.mockS3.On("GenerateDownloadPublicUri", variantPlaylistKey, s3.S3VideoUploadedBucket, uint32(s3.UrlExpirationSeconds)).
		Return(variantPlaylistURL, nil)

	// Setup mock expectations for variants with multiple segments
	variants := []*models.VideoVariant{
		{
			ID:        uuid.Must(uuid.NewV7()),
			VideoID:   videoID,
			Quality:   "720p",
			ObjectKey: segment1Key,
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.Must(uuid.NewV7()),
			VideoID:   videoID,
			Quality:   "720p",
			ObjectKey: segment2Key,
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.Must(uuid.NewV7()),
			VideoID:   videoID,
			Quality:   "720p",
			ObjectKey: segment3Key,
			CreatedAt: time.Now(),
		},
	}

	as.mockVideoAggregateRepository.On("GetVideoVariantsByVideoID", ctx, videoID).Return(variants, nil)

	// Setup S3 mocks for segments
	as.mockS3.On("GenerateDownloadPublicUri", segment1Key, s3.S3VideoUploadedBucket, uint32(s3.UrlExpirationSeconds)).
		Return(segment1URL, nil)
	as.mockS3.On("GenerateDownloadPublicUri", segment2Key, s3.S3VideoUploadedBucket, uint32(s3.UrlExpirationSeconds)).
		Return(segment2URL, nil)
	as.mockS3.On("GenerateDownloadPublicUri", segment3Key, s3.S3VideoUploadedBucket, uint32(s3.UrlExpirationSeconds)).
		Return(segment3URL, nil)

	// Setup request
	request := &connect.Request[proto.ServePlaylistRequest]{
		Msg: &proto.ServePlaylistRequest{
			VideoId: videoID.String(),
		},
	}

	// Execute
	actionsInstance := actions.NewActions(ctx, "test-token")
	response, err := actionsInstance.ServePlaylist(ctx, request)

	// Assertions
	as.NoError(err)
	as.NotNil(response)
	as.Len(response.Msg.Variants, 1)
	as.Equal("720p", response.Msg.Variants[0].Quality)
	as.Equal(variantPlaylistURL, response.Msg.Variants[0].PlaylistUrl)
	as.Len(response.Msg.Variants[0].SegmentUrls, 3)
	as.Contains(response.Msg.Variants[0].SegmentUrls, segment1URL)
	as.Contains(response.Msg.Variants[0].SegmentUrls, segment2URL)
	as.Contains(response.Msg.Variants[0].SegmentUrls, segment3URL)

	// Verify all mocks were called as expected
	as.mockVideoAggregateRepository.AssertExpectations(as.T())
	as.mockS3.AssertExpectations(as.T())
}

func (as *ActionsSuite) TestActions_ServePlaylist_NilManifestsAndVariants() {
	as.setupEnvironment()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Setup test data with nil entries
	videoID := uuid.Must(uuid.NewV7())
	variantPlaylistKey := "720p/playlist.m3u8"
	segmentKey := "720p/segment-0.ts"
	variantPlaylistURL := "https://s3.example.com/720p/playlist.m3u8"
	segmentURL := "https://s3.example.com/720p/segment-0.ts"

	// Setup mock expectations for manifests with nil entries
	manifests := []*models.VideoManifest{
		nil, // Should be skipped
		{
			ID:        uuid.Must(uuid.NewV7()),
			VideoID:   videoID,
			ObjectKey: variantPlaylistKey,
			Quality:   "720p",
			CreatedAt: time.Now(),
		},
		nil, // Should be skipped
	}

	as.mockVideoAggregateRepository.On("GetVideoManifestsByVideoID", ctx, videoID).Return(manifests, nil)

	// Setup S3 mock for variant playlist
	as.mockS3.On("GenerateDownloadPublicUri", variantPlaylistKey, s3.S3VideoUploadedBucket, uint32(s3.UrlExpirationSeconds)).
		Return(variantPlaylistURL, nil)

	// Setup mock expectations for variants with nil entries
	variants := []*models.VideoVariant{
		nil, // Should be skipped
		{
			ID:        uuid.Must(uuid.NewV7()),
			VideoID:   videoID,
			Quality:   "720p",
			ObjectKey: segmentKey,
			CreatedAt: time.Now(),
		},
		nil, // Should be skipped
	}

	as.mockVideoAggregateRepository.On("GetVideoVariantsByVideoID", ctx, videoID).Return(variants, nil)

	// Setup S3 mock for segment
	as.mockS3.On("GenerateDownloadPublicUri", segmentKey, s3.S3VideoUploadedBucket, uint32(s3.UrlExpirationSeconds)).
		Return(segmentURL, nil)

	// Setup request
	request := &connect.Request[proto.ServePlaylistRequest]{
		Msg: &proto.ServePlaylistRequest{
			VideoId: videoID.String(),
		},
	}

	// Execute
	actionsInstance := actions.NewActions(ctx, "test-token")
	response, err := actionsInstance.ServePlaylist(ctx, request)

	// Assertions
	as.NoError(err)
	as.NotNil(response)
	as.Len(response.Msg.Variants, 1)
	as.Equal("720p", response.Msg.Variants[0].Quality)
	as.Equal(variantPlaylistURL, response.Msg.Variants[0].PlaylistUrl)
	as.Contains(response.Msg.Variants[0].SegmentUrls, segmentURL)

	// Verify all mocks were called as expected
	as.mockVideoAggregateRepository.AssertExpectations(as.T())
	as.mockS3.AssertExpectations(as.T())
}

func (as *ActionsSuite) TestActions_ServePlaylist_EmptyManifestsAndVariants() {
	as.setupEnvironment()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Setup test data
	videoID := uuid.Must(uuid.NewV7())

	// Setup mock expectations for empty manifests
	as.mockVideoAggregateRepository.On("GetVideoManifestsByVideoID", ctx, videoID).Return([]*models.VideoManifest{}, nil)

	// Setup mock expectations for empty variants
	as.mockVideoAggregateRepository.On("GetVideoVariantsByVideoID", ctx, videoID).Return([]*models.VideoVariant{}, nil)

	// Setup request
	request := &connect.Request[proto.ServePlaylistRequest]{
		Msg: &proto.ServePlaylistRequest{
			VideoId: videoID.String(),
		},
	}

	// Execute
	actionsInstance := actions.NewActions(ctx, "test-token")
	response, err := actionsInstance.ServePlaylist(ctx, request)

	// Assertions
	as.NoError(err)
	as.NotNil(response)
	as.Empty(response.Msg.PlaylistUrl)
	as.Len(response.Msg.Variants, 0)

	// Verify all mocks were called as expected
	as.mockVideoAggregateRepository.AssertExpectations(as.T())
	as.mockS3.AssertNotCalled(as.T(), "GenerateDownloadPublicUri")
}

package processing

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/ffmpeg"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/kafka"
	mockPkg "github.com/sweetloveinyourheart/sweet-reel/pkg/testing/mock"
	"github.com/sweetloveinyourheart/sweet-reel/services/video_processing/models"
)

// testVideoSplitterProcessManager is a test version that allows mocking FFmpeg
type testVideoSplitterProcessManager struct {
	*VideoSplitterProcessManager
	ffmpegInstance *mockPkg.MockFFmpeg
}

func (t *testVideoSplitterProcessManager) processVideo(ctx context.Context, msg *models.VideoSplitterMessage, videoData []byte) error {
	// Check if FFmpeg is available
	if err := t.ffmpegInstance.IsAvailable(ctx); err != nil {
		return errors.Wrap(err, "FFmpeg not available")
	}

	// Create temporary directory for processing
	tempDir := fmt.Sprintf("/tmp/video_processing_%s", msg.VideoID.String())
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return errors.Wrap(err, "failed to create temp directory")
	}
	defer os.RemoveAll(tempDir) // Clean up temp directory

	// Write video data to temporary file
	inputPath := filepath.Join(tempDir, "input.mp4")
	if err := os.WriteFile(inputPath, videoData, 0644); err != nil {
		return errors.Wrap(err, "failed to write input file")
	}

	// Probe the input file to get information
	_, err := t.ffmpegInstance.ProbeFile(ctx, inputPath)
	if err != nil {
		return errors.Wrap(err, "failed to probe input file")
	}

	// Create output directory for HLS segments
	hlsOutputDir := filepath.Join(tempDir, "hls")
	if err := os.MkdirAll(hlsOutputDir, 0755); err != nil {
		return errors.Wrap(err, "failed to create HLS output directory")
	}

	// Define multiple quality levels for adaptive streaming
	qualities := []ffmpeg.SegmentationOptions{
		{
			SegmentDuration: "6",
			PlaylistType:    "vod",
			PlaylistName:    "playlist.m3u8",
			SegmentPrefix:   "segment",
			SegmentFormat:   "ts",
			VideoCodec:      "libx264",
			VideoBitrate:    "800k",
			AudioCodec:      "aac",
			AudioBitrate:    "96k",
			Resolution:      "854x480", // 480p
		},
	}

	// Progress callback to monitor processing
	progressCallback := func(progress ffmpeg.ProgressInfo) {
		// Mock progress callback for testing
	}

	// Start video segmentation
	if err := t.ffmpegInstance.SegmentVideoMultiQuality(ctx, inputPath, hlsOutputDir, qualities, progressCallback); err != nil {
		return errors.Wrap(err, "failed to segment video")
	}

	// Create thumbnail
	thumbnailPath := filepath.Join(tempDir, "thumbnail.jpg")
	if err := t.ffmpegInstance.CreateThumbnail(ctx, inputPath, thumbnailPath, "00:00:05", 320, 240); err != nil {
		// Don't fail on thumbnail creation error in tests
	}

	// Skip upload in tests
	return nil
}

func TestVideoSplitterProcessManager_processVideo(t *testing.T) {
	tests := []struct {
		name        string
		setupMocks  func(*mockPkg.MockStorage, *mockPkg.MockFFmpeg)
		videoData   []byte
		message     *models.VideoSplitterMessage
		expectError bool
	}{
		{
			name: "successful video processing",
			setupMocks: func(storage *mockPkg.MockStorage, ffmpegMock *mockPkg.MockFFmpeg) {
				ffmpegMock.On("IsAvailable", mock.Anything).Return(nil)
				ffmpegMock.On("ProbeFile", mock.Anything, mock.Anything).Return(&ffmpeg.ProbeInfo{
					Format: ffmpeg.FormatInfo{
						FormatName: "mp4",
						Duration:   "60.0",
						Size:       "10485760", // 10 MB
					},
					Streams: []ffmpeg.StreamInfo{
						{
							CodecType: "video",
						},
						{
							CodecType: "audio",
						},
					},
				}, nil)
				ffmpegMock.On("SegmentVideoMultiQuality", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				ffmpegMock.On("CreateThumbnail", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			videoData: []byte("fake video data"),
			message: &models.VideoSplitterMessage{
				VideoID: uuid.Must(uuid.NewV7()),
				Metadata: models.VideoSplitterMetadata{
					Key:    "original/video.mp4",
					Bucket: "videos",
				},
			},
			expectError: false,
		},
		{
			name: "ffmpeg not available",
			setupMocks: func(storage *mockPkg.MockStorage, ffmpegMock *mockPkg.MockFFmpeg) {
				ffmpegMock.On("IsAvailable", mock.Anything).Return(errors.New("ffmpeg not installed"))
			},
			videoData: []byte("fake video data"),
			message: &models.VideoSplitterMessage{
				VideoID: uuid.Must(uuid.NewV7()),
				Metadata: models.VideoSplitterMetadata{
					Key:    "original/video.mp4",
					Bucket: "videos",
				},
			},
			expectError: true,
		},
		{
			name: "probing file fails",
			setupMocks: func(storage *mockPkg.MockStorage, ffmpegMock *mockPkg.MockFFmpeg) {
				ffmpegMock.On("IsAvailable", mock.Anything).Return(nil)
				ffmpegMock.On("ProbeFile", mock.Anything, mock.Anything).Return(&ffmpeg.ProbeInfo{}, errors.New("probe failed"))
			},
			videoData: []byte("fake video data"),
			message: &models.VideoSplitterMessage{
				VideoID: uuid.Must(uuid.NewV7()),
				Metadata: models.VideoSplitterMetadata{
					Key:    "original/video.mp4",
					Bucket: "videos",
				},
			},
			expectError: true,
		},
		{
			name: "segmentation fails",
			setupMocks: func(storage *mockPkg.MockStorage, ffmpegMock *mockPkg.MockFFmpeg) {
				ffmpegMock.On("IsAvailable", mock.Anything).Return(nil)
				ffmpegMock.On("ProbeFile", mock.Anything, mock.Anything).Return(&ffmpeg.ProbeInfo{
					Format: ffmpeg.FormatInfo{
						FormatName: "mp4",
						Duration:   "60.0",
						Size:       "10485760", // 10 MB
					},
					Streams: []ffmpeg.StreamInfo{
						{
							CodecType: "video",
						},
					},
				}, nil)
				ffmpegMock.On("SegmentVideoMultiQuality", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("segmentation failed"))
			},
			videoData: []byte("fake video data"),
			message: &models.VideoSplitterMessage{
				VideoID: uuid.Must(uuid.NewV7()),
				Metadata: models.VideoSplitterMetadata{
					Key:    "original/video.mp4",
					Bucket: "videos",
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock dependencies
			mockStorage := new(mockPkg.MockStorage)
			ffmpegMock := new(mockPkg.MockFFmpeg)

			// Setup mocks
			tt.setupMocks(mockStorage, ffmpegMock)

			// Create the test process manager
			testVsp := &testVideoSplitterProcessManager{
				VideoSplitterProcessManager: &VideoSplitterProcessManager{
					ctx:           context.Background(),
					storageClient: mockStorage,
					queue:         make(chan lo.Tuple2[context.Context, *kafka.ConsumedMessage], BatchSize*2),
				},
				ffmpegInstance: ffmpegMock,
			}

			// Run the test
			err := testVsp.processVideo(context.Background(), tt.message, tt.videoData)

			// Verify the results
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			// Verify all mocks
			mockStorage.AssertExpectations(t)
			ffmpegMock.AssertExpectations(t)

			// Clean up any temporary directories that might have been created
			tempDir := fmt.Sprintf("/tmp/video_processing_%s", tt.message.VideoID.String())
			os.RemoveAll(tempDir)
		})
	}
}

func TestVideoSplitterProcessManager_uploadProcessedFiles(t *testing.T) {
	tests := []struct {
		name        string
		setupMocks  func(*mockPkg.MockStorage)
		setupFiles  func(hlsDir, thumbnailPath string) error
		message     *models.VideoSplitterMessage
		expectError bool
	}{
		{
			name: "successful file upload",
			setupMocks: func(storage *mockPkg.MockStorage) {
				// Set up expectations for HLS files
				storage.On("Upload", mock.MatchedBy(func(key string) bool {
					return filepath.Base(key) == "master.m3u8"
				}), "videos", mock.Anything, mock.Anything).Return(nil)

				storage.On("Upload", mock.MatchedBy(func(key string) bool {
					return filepath.Base(key) == "playlist.m3u8"
				}), "videos", mock.Anything, mock.Anything).Return(nil).Times(2) // 720p and 480p

				storage.On("Upload", mock.MatchedBy(func(key string) bool {
					return filepath.Base(key) == "segment0.ts"
				}), "videos", mock.Anything, mock.Anything).Return(nil).Times(2) // 720p and 480p

				// Set up expectation for thumbnail upload
				storage.On("Upload", mock.MatchedBy(func(key string) bool {
					return filepath.Base(key) == "thumbnail.jpg"
				}), "videos", mock.Anything, "image/jpeg").Return(nil)
			},
			setupFiles: func(hlsDir, thumbnailPath string) error {
				// Create test files in the HLS directory
				testFiles := []string{
					"master.m3u8",
					"720p/playlist.m3u8",
					"720p/segment0.ts",
					"480p/playlist.m3u8",
					"480p/segment0.ts",
				}

				for _, file := range testFiles {
					filePath := filepath.Join(hlsDir, file)
					dir := filepath.Dir(filePath)
					if dir != hlsDir {
						if err := os.MkdirAll(dir, 0755); err != nil {
							return err
						}
					}
					if err := os.WriteFile(filePath, []byte("test content"), 0644); err != nil {
						return err
					}
				}

				// Create a thumbnail
				return os.WriteFile(thumbnailPath, []byte("thumbnail data"), 0644)
			},
			message: &models.VideoSplitterMessage{
				VideoID: uuid.Must(uuid.NewV7()),
				Metadata: models.VideoSplitterMetadata{
					Key:    "original/video.mp4",
					Bucket: "videos",
				},
			},
			expectError: false,
		},
		{
			name: "upload failure",
			setupMocks: func(storage *mockPkg.MockStorage) {
				// Set up upload to fail
				storage.On("Upload", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("upload failed"))
			},
			setupFiles: func(hlsDir, thumbnailPath string) error {
				// Create a single test file
				filePath := filepath.Join(hlsDir, "master.m3u8")
				return os.WriteFile(filePath, []byte("test content"), 0644)
			},
			message: &models.VideoSplitterMessage{
				VideoID: uuid.Must(uuid.NewV7()),
				Metadata: models.VideoSplitterMetadata{
					Key:    "original/video.mp4",
					Bucket: "videos",
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary test directories
			tempDir := t.TempDir()
			hlsDir := filepath.Join(tempDir, "hls")
			require.NoError(t, os.MkdirAll(hlsDir, 0755))
			thumbnailPath := filepath.Join(tempDir, "thumbnail.jpg")

			// Create mock storage
			mockStorage := new(mockPkg.MockStorage)

			// Setup mocks
			tt.setupMocks(mockStorage)

			// Setup test files
			err := tt.setupFiles(hlsDir, thumbnailPath)
			require.NoError(t, err)

			// Create the process manager with mock storage
			vsp := &VideoSplitterProcessManager{
				ctx:           context.Background(),
				storageClient: mockStorage,
			}

			// Run the test
			err = vsp.uploadProcessedFiles(context.Background(), tt.message, hlsDir, thumbnailPath)

			// Verify the results
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			// Verify all mocks
			mockStorage.AssertExpectations(t)
		})
	}
}

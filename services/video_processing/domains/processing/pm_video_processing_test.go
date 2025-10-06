package processing

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/ffmpeg"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/kafka"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/s3"
	mockPkg "github.com/sweetloveinyourheart/sweet-reel/pkg/testing/mock"
)

func TestVideoProcessManager_HandleMessage(t *testing.T) {
	videoID := uuid.Must(uuid.NewV7()).String()

	tests := []struct {
		name              string
		setupStorageMocks func(*mockPkg.MockStorage)
		setupFfmpegMocks  func(*mockPkg.MockFFmpeg)
		message           *kafka.ConsumedMessage
		expectError       bool
		errorMsg          string
	}{
		{
			name: "successful message handling",
			setupStorageMocks: func(storage *mockPkg.MockStorage) {
				storage.On("Download", fmt.Sprintf("%s.mp4", videoID), "video-uploaded").Return([]byte("fake video data"), nil)
				// Mock successful uploads for HLS files and thumbnail
				storage.On("Upload", mock.AnythingOfType("string"), "video-uploaded", mock.Anything, mock.AnythingOfType("string")).Return(nil).Maybe()
			},
			setupFfmpegMocks: func(ff *mockPkg.MockFFmpeg) {
				// Mock FFmpeg availability check
				ff.On("IsAvailable", mock.Anything).Return(nil)

				// Mock probe file to return valid video information
				probeInfo := &ffmpeg.ProbeInfo{
					Format: ffmpeg.FormatInfo{
						FormatName: "mov,mp4,m4a,3gp,3g2,mj2",
						Duration:   "120.000000",
						Size:       "1048576",
					},
					Streams: []ffmpeg.StreamInfo{
						{
							Index:     0,
							CodecType: "video",
							CodecName: "h264",
							Width:     1920,
							Height:    1080,
							Duration:  "120.000000",
						},
					},
				}
				ff.On("ProbeFile", mock.Anything, mock.AnythingOfType("string")).Return(probeInfo, nil)

				// Mock successful video segmentation
				ff.On("SegmentVideoMultiQuality", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("[]ffmpeg.SegmentationOptions"), mock.AnythingOfType("ffmpeg.ProgressCallback")).Return(nil)

				// Mock successful thumbnail creation
				ff.On("CreateThumbnail", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), "00:00:05", 320, 240).Return(nil)
			},
			message: &kafka.ConsumedMessage{
				Topic: KafkaVideoUploadedTopic,
				Key:   "test-key",
				Value: func() []byte {
					msg := s3.S3EventMessage{
						Key: fmt.Sprintf("video-uploaded/%s.mp4", videoID),
					}
					data, _ := json.Marshal(msg)
					return data
				}(),
			},
			expectError: false,
		},
		{
			name:              "nil message",
			setupStorageMocks: func(storage *mockPkg.MockStorage) {},
			setupFfmpegMocks:  func(ff *mockPkg.MockFFmpeg) {},
			message:           nil,
			expectError:       true,
			errorMsg:          "message is nil",
		},
		{
			name: "invalid video id",
			setupStorageMocks: func(storage *mockPkg.MockStorage) {
				storage.On("Download", fmt.Sprintf("%s.mp4", "fake-video-id"), "video-uploaded").Return([]byte("fake video data"), nil)
			},
			setupFfmpegMocks: func(ff *mockPkg.MockFFmpeg) {},
			message: &kafka.ConsumedMessage{
				Topic: KafkaVideoUploadedTopic,
				Key:   "test-key",
				Value: func() []byte {
					msg := s3.S3EventMessage{
						Key: fmt.Sprintf("video-uploaded/%s.mp4", "fake-video-id"),
					}
					data, _ := json.Marshal(msg)
					return data
				}(),
			},
			expectError: true,
			errorMsg:    "invalid video id",
		},
		{
			name:              "invalid JSON message",
			setupStorageMocks: func(storage *mockPkg.MockStorage) {},
			setupFfmpegMocks:  func(ff *mockPkg.MockFFmpeg) {},
			message: &kafka.ConsumedMessage{
				Topic: KafkaVideoUploadedTopic,
				Key:   "test-key",
				Value: []byte("invalid json"),
			},
			expectError: true,
		},
		{
			name: "storage download failure",
			setupStorageMocks: func(storage *mockPkg.MockStorage) {
				storage.On("Download", fmt.Sprintf("%s.mp4", videoID), "video-uploaded").Return([]byte(nil), errors.New("download failed"))
			},
			setupFfmpegMocks: func(ff *mockPkg.MockFFmpeg) {},
			message: &kafka.ConsumedMessage{
				Topic: KafkaVideoUploadedTopic,
				Key:   "test-key",
				Value: func() []byte {
					msg := s3.S3EventMessage{
						Key: fmt.Sprintf("video-uploaded/%s.mp4", videoID),
					}
					data, _ := json.Marshal(msg)
					return data
				}(),
			},
			expectError: true,
			errorMsg:    "download failed",
		},
		{
			name: "ffmpeg not available",
			setupStorageMocks: func(storage *mockPkg.MockStorage) {
				storage.On("Download", fmt.Sprintf("%s.mp4", videoID), "video-uploaded").Return([]byte("fake video data"), nil)
			},
			setupFfmpegMocks: func(ff *mockPkg.MockFFmpeg) {
				ff.On("IsAvailable", mock.Anything).Return(errors.New("ffmpeg not found"))
			},
			message: &kafka.ConsumedMessage{
				Topic: KafkaVideoUploadedTopic,
				Key:   "test-key",
				Value: func() []byte {
					msg := s3.S3EventMessage{
						Key: fmt.Sprintf("video-uploaded/%s.mp4", videoID),
					}
					data, _ := json.Marshal(msg)
					return data
				}(),
			},
			expectError: true,
			errorMsg:    "FFmpeg not available",
		},
		{
			name: "probe file failure",
			setupStorageMocks: func(storage *mockPkg.MockStorage) {
				storage.On("Download", fmt.Sprintf("%s.mp4", videoID), "video-uploaded").Return([]byte("fake video data"), nil)
			},
			setupFfmpegMocks: func(ff *mockPkg.MockFFmpeg) {
				ff.On("IsAvailable", mock.Anything).Return(nil)
				ff.On("ProbeFile", mock.Anything, mock.AnythingOfType("string")).Return((*ffmpeg.ProbeInfo)(nil), errors.New("probe failed"))
			},
			message: &kafka.ConsumedMessage{
				Topic: KafkaVideoUploadedTopic,
				Key:   "test-key",
				Value: func() []byte {
					msg := s3.S3EventMessage{
						Key: fmt.Sprintf("video-uploaded/%s.mp4", videoID),
					}
					data, _ := json.Marshal(msg)
					return data
				}(),
			},
			expectError: true,
			errorMsg:    "failed to probe input file",
		},
		{
			name: "video segmentation failure",
			setupStorageMocks: func(storage *mockPkg.MockStorage) {
				storage.On("Download", fmt.Sprintf("%s.mp4", videoID), "video-uploaded").Return([]byte("fake video data"), nil)
			},
			setupFfmpegMocks: func(ff *mockPkg.MockFFmpeg) {
				ff.On("IsAvailable", mock.Anything).Return(nil)

				probeInfo := &ffmpeg.ProbeInfo{
					Format: ffmpeg.FormatInfo{
						FormatName: "mov,mp4,m4a,3gp,3g2,mj2",
						Duration:   "120.000000",
						Size:       "1048576",
					},
					Streams: []ffmpeg.StreamInfo{
						{
							Index:     0,
							CodecType: "video",
							CodecName: "h264",
							Width:     1920,
							Height:    1080,
							Duration:  "120.000000",
						},
					},
				}
				ff.On("ProbeFile", mock.Anything, mock.AnythingOfType("string")).Return(probeInfo, nil)
				ff.On("SegmentVideoMultiQuality", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("[]ffmpeg.SegmentationOptions"), mock.AnythingOfType("ffmpeg.ProgressCallback")).Return(errors.New("segmentation failed"))
			},
			message: &kafka.ConsumedMessage{
				Topic: KafkaVideoUploadedTopic,
				Key:   "test-key",
				Value: func() []byte {
					msg := s3.S3EventMessage{
						Key: fmt.Sprintf("video-uploaded/%s.mp4", videoID),
					}
					data, _ := json.Marshal(msg)
					return data
				}(),
			},
			expectError: true,
			errorMsg:    "failed to segment video",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock storage and ffmpeg
			mockStorage := new(mockPkg.MockStorage)
			mockFf := new(mockPkg.MockFFmpeg)

			// Setup mocks
			tt.setupStorageMocks(mockStorage)
			tt.setupFfmpegMocks(mockFf)

			vsp := &VideoProcessManager{
				ctx:           context.Background(),
				storageClient: mockStorage,
				ff:            mockFf,
			}

			// Run the test
			err := vsp.HandleMessage(context.Background(), tt.message)

			// Verify the results
			if tt.expectError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					require.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				require.NoError(t, err)
			}

			// Verify all mocks
			mockStorage.AssertExpectations(t)
			mockFf.AssertExpectations(t)
		})
	}
}

func TestVideoProcessManager_uploadProcessedFiles(t *testing.T) {
	tests := []struct {
		name                string
		setupMocks          func(*mockPkg.MockStorage)
		setupFiles          func(hlsDir, thumbnailPath string) error
		expectError         bool
		errorMsg            string
		skipDirectoryCreate bool // For cases where we don't want to create the HLS directory
		lenientMockCheck    bool // For cases where mock assertions should be more lenient
	}{
		{
			name: "successful file upload",
			setupMocks: func(storage *mockPkg.MockStorage) {
				// Set up expectations for HLS files
				storage.On("Upload", mock.MatchedBy(func(key string) bool {
					return filepath.Base(key) == "master.m3u8"
				}), "video-processed", mock.Anything, mock.Anything).Return(nil)

				storage.On("Upload", mock.MatchedBy(func(key string) bool {
					return filepath.Base(key) == "playlist.m3u8"
				}), "video-processed", mock.Anything, mock.Anything).Return(nil).Times(2) // 720p and 480p

				storage.On("Upload", mock.MatchedBy(func(key string) bool {
					return filepath.Base(key) == "segment0.ts"
				}), "video-processed", mock.Anything, mock.Anything).Return(nil).Times(2) // 720p and 480p

				// Set up expectation for thumbnail upload
				storage.On("Upload", mock.MatchedBy(func(key string) bool {
					return filepath.Base(key) == "thumbnail.jpg"
				}), "video-processed", mock.Anything, "image/jpeg").Return(nil)
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

			expectError: true,
			errorMsg:    "upload failed",
		},
		{
			name: "non-existent HLS directory",
			setupMocks: func(storage *mockPkg.MockStorage) {
				// No mocks needed as this should fail before reaching storage
			},
			setupFiles: func(hlsDir, thumbnailPath string) error {
				// Don't create the HLS directory
				return nil
			},

			expectError:         true,
			errorMsg:            "no such file or directory",
			skipDirectoryCreate: true,
		},
		{
			name: "partial upload failure with specific file",
			setupMocks: func(storage *mockPkg.MockStorage) {
				// First file succeeds
				storage.On("Upload", mock.MatchedBy(func(key string) bool {
					return filepath.Base(key) == "master.m3u8"
				}), "video-processed", mock.Anything, mock.Anything).Return(nil).Once()

				// Second file fails
				storage.On("Upload", mock.MatchedBy(func(key string) bool {
					return filepath.Base(key) == "playlist.m3u8"
				}), "video-processed", mock.Anything, mock.Anything).Return(errors.New("network error")).Once()
			},
			setupFiles: func(hlsDir, thumbnailPath string) error {
				testFiles := []string{
					"master.m3u8",
					"playlist.m3u8",
				}

				for _, file := range testFiles {
					filePath := filepath.Join(hlsDir, file)
					if err := os.WriteFile(filePath, []byte("test content"), 0644); err != nil {
						return err
					}
				}
				return nil
			},

			expectError: true,
			errorMsg:    "network error",
		},
		{
			name: "thumbnail upload failure (should not fail the whole process)",
			setupMocks: func(storage *mockPkg.MockStorage) {
				// HLS files upload successfully
				storage.On("Upload", mock.MatchedBy(func(key string) bool {
					return filepath.Base(key) == "master.m3u8"
				}), "video-processed", mock.Anything, mock.Anything).Return(nil)

				// Thumbnail upload fails (but this should not cause the method to fail)
				storage.On("Upload", mock.MatchedBy(func(key string) bool {
					return filepath.Base(key) == "thumbnail.jpg"
				}), "video-processed", mock.Anything, "image/jpeg").Return(errors.New("thumbnail upload failed"))
			},
			setupFiles: func(hlsDir, thumbnailPath string) error {
				// Create HLS file
				filePath := filepath.Join(hlsDir, "master.m3u8")
				if err := os.WriteFile(filePath, []byte("test content"), 0644); err != nil {
					return err
				}

				// Create thumbnail
				return os.WriteFile(thumbnailPath, []byte("thumbnail data"), 0644)
			},

			expectError:      false, // Thumbnail failure should not fail the whole process
			lenientMockCheck: true,  // Don't strictly check mocks for this case
		},
		{
			name: "empty HLS directory",
			setupMocks: func(storage *mockPkg.MockStorage) {
				// No upload calls expected for empty directory
			},
			setupFiles: func(hlsDir, thumbnailPath string) error {
				// Create empty HLS directory
				return nil
			},

			expectError: false, // Empty directory should not cause error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary test directories
			tempDir := t.TempDir()
			hlsDir := filepath.Join(tempDir, "hls")

			// Only create HLS directory if the test expects it
			if !tt.skipDirectoryCreate {
				require.NoError(t, os.MkdirAll(hlsDir, 0755))
			}

			thumbnailPath := filepath.Join(tempDir, "thumbnail.jpg")

			// Create mock storage
			mockStorage := new(mockPkg.MockStorage)

			// Setup mocks
			tt.setupMocks(mockStorage)

			// Setup test files
			err := tt.setupFiles(hlsDir, thumbnailPath)
			require.NoError(t, err)

			// Create the process manager with mock storage
			vsp := &VideoProcessManager{
				ctx:           context.Background(),
				storageClient: mockStorage,
			}

			// Run the test
			err = vsp.uploadProcessedFiles(context.Background(), uuid.Must(uuid.NewV7()), hlsDir, thumbnailPath)

			// Verify the results
			if tt.expectError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					require.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				require.NoError(t, err)
			}

			// Verify all mocks (with more lenient assertion for special cases)
			if !tt.lenientMockCheck {
				mockStorage.AssertExpectations(t)
			}
		})
	}
}

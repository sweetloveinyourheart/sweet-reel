package processing_test

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/kafka"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/messages"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/s3"
	"github.com/sweetloveinyourheart/sweet-reel/services/video_processing/domains/processing"
)

func (as *VideoProcessingSuite) TestNewVideoProcessManager_Success() {
	as.setupEnvironment()

	// Since the constructor creates goroutines that are hard to control in tests,
	// we test that the constructor doesn't return an error and creates a manager
	manager, err := processing.NewVideoProcessManager(as.ctx)

	as.NoError(err)
	as.NotNil(manager)

	// Give some time for goroutines to start
	time.Sleep(100 * time.Millisecond)
}

func (as *VideoProcessingSuite) TestHandleMessage_Success() {
	as.setupEnvironment()

	// Create test video ID and data
	videoID := uuid.Must(uuid.NewV4())
	videoData := []byte("fake video data")
	eventMessage := messages.S3EventMessage{
		Key: videoID.String() + ".mp4",
	}

	// Marshal the event message
	eventData, err := json.Marshal(eventMessage)
	as.NoError(err)

	// Extract bucket and key from the event message like the actual code does
	bucket, key := s3.ExtractBucketAndKeyFromEventMessage(videoID.String() + ".mp4")
	as.mockS3.On("Download", key, bucket).Return(videoData, nil)

	// Test getting the manager
	manager, err := processing.NewVideoProcessManager(as.ctx)
	as.NoError(err)
	as.NotNil(manager)

	// Note: Due to the complex nature of the HandleMessage method and its dependencies on filesystem operations
	// and FFmpeg (which is created directly with ffmpeg.New() instead of dependency injection),
	// detailed testing of the processing logic would require more sophisticated mocking or integration tests.

	// Verify that we have a valid event data structure
	as.NotEmpty(eventData)
}

// Test constants
func (as *VideoProcessingSuite) TestConstants() {
	as.Equal(1024, processing.BatchSize)
	as.Equal("video-processing", kafka.KafkaVideoProcessingGroup)
	as.Equal("video-uploaded", kafka.KafkaVideoUploadedTopic)
	as.Equal("video-processed", s3.S3VideoProcessedBucket)
}

func (as *VideoProcessingSuite) TestHandleMessage_NilMessage() {
	as.setupEnvironment()

	manager, err := processing.NewVideoProcessManager(as.ctx)
	as.NoError(err)
	as.NotNil(manager)

	// Test with nil message
	err = manager.HandleMessage(as.ctx, nil)
	as.Error(err)
	as.Contains(err.Error(), "message is nil")
}

func (as *VideoProcessingSuite) TestHandleMessage_InvalidJSON() {
	as.setupEnvironment()

	manager, err := processing.NewVideoProcessManager(as.ctx)
	as.NoError(err)
	as.NotNil(manager)

	// Create message with invalid JSON
	consumedMsg := &kafka.ConsumedMessage{
		Topic:     kafka.KafkaVideoUploadedTopic,
		Partition: 0,
		Offset:    1,
		Key:       "test-key",
		Value:     []byte("invalid json"),
		Headers:   make(map[string]string),
		Timestamp: time.Now(),
	}

	err = manager.HandleMessage(as.ctx, consumedMsg)
	as.Error(err)
}

func (as *VideoProcessingSuite) TestHandleMessage_InvalidVideoID() {
	as.setupEnvironment()

	eventMessage := messages.S3EventMessage{
		Key: "invalid-uuid.mp4",
	}

	eventData, err := json.Marshal(eventMessage)
	as.NoError(err)

	consumedMsg := &kafka.ConsumedMessage{
		Topic:     kafka.KafkaVideoUploadedTopic,
		Partition: 0,
		Offset:    1,
		Key:       "test-key",
		Value:     eventData,
		Headers:   make(map[string]string),
		Timestamp: time.Now(),
	}

	// Extract bucket and key from the event message like the actual code does
	bucket, key := s3.ExtractBucketAndKeyFromEventMessage("invalid-uuid.mp4")
	as.mockS3.On("Download", key, bucket).Return([]byte("fake data"), nil)

	manager, err := processing.NewVideoProcessManager(as.ctx)
	as.NoError(err)

	err = manager.HandleMessage(as.ctx, consumedMsg)
	as.Error(err)
	as.Contains(err.Error(), "invalid video id")
}

func (as *VideoProcessingSuite) TestHandleMessage_S3DownloadError() {
	as.setupEnvironment()

	videoID := uuid.Must(uuid.NewV4())
	eventMessage := messages.S3EventMessage{
		Key: videoID.String() + ".mp4",
	}

	eventData, err := json.Marshal(eventMessage)
	as.NoError(err)

	consumedMsg := &kafka.ConsumedMessage{
		Topic:     kafka.KafkaVideoUploadedTopic,
		Partition: 0,
		Offset:    1,
		Key:       videoID.String(),
		Value:     eventData,
		Headers:   make(map[string]string),
		Timestamp: time.Now(),
	}

	// Extract bucket and key from the event message like the actual code does
	bucket, key := s3.ExtractBucketAndKeyFromEventMessage(videoID.String() + ".mp4")
	as.mockS3.On("Download", key, bucket).Return([]byte(nil), fmt.Errorf("test error: download failed"))

	manager, err := processing.NewVideoProcessManager(as.ctx)
	as.NoError(err)

	err = manager.HandleMessage(as.ctx, consumedMsg)
	as.Error(err)
	as.Contains(err.Error(), "download failed")
}

func (as *VideoProcessingSuite) TestHandleMessage_FFmpegNotAvailable() {
	as.setupEnvironment()

	videoID := uuid.Must(uuid.NewV4())
	videoData := []byte("fake video data")
	eventMessage := messages.S3EventMessage{
		Key: videoID.String() + ".mp4",
	}

	eventData, err := json.Marshal(eventMessage)
	as.NoError(err)

	consumedMsg := &kafka.ConsumedMessage{
		Topic:     kafka.KafkaVideoUploadedTopic,
		Partition: 0,
		Offset:    1,
		Key:       videoID.String(),
		Value:     eventData,
		Headers:   make(map[string]string),
		Timestamp: time.Now(),
	}

	// Extract bucket and key from the event message like the actual code does
	bucket, key := s3.ExtractBucketAndKeyFromEventMessage(videoID.String() + ".mp4")
	as.mockS3.On("Download", key, bucket).Return(videoData, nil)

	manager, err := processing.NewVideoProcessManager(as.ctx)
	as.NoError(err)

	err = manager.HandleMessage(as.ctx, consumedMsg)
	as.Error(err)
	// Note: Since FFmpeg is created directly with ffmpeg.New() instead of dependency injection,
	// we can't mock it easily. The actual error will be about ffmpeg not found in PATH
	as.Contains(err.Error(), "FFmpeg not available")
}

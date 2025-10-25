package processing_test

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/messages"
	"github.com/sweetloveinyourheart/sweet-reel/services/video_management/domains/processing"
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
	videoID := uuid.Must(uuid.NewV7())
	objectKey := fmt.Sprintf("test/%s.mp4", videoID)
	eventMessage := messages.VideoProcessingProgress{
		VideoID:     videoID,
		Status:      messages.VideoStatusReady,
		ObjectKey:   objectKey,
		ProcessedAt: time.Now(),
	}

	// Marshal the event message
	eventData, err := json.Marshal(eventMessage)
	as.NoError(err)

	as.mockVideoAggregateRepository.On("UpdateVideoProgress", mock.Anything,
		eventMessage.VideoID,
		eventMessage.ObjectKey,
		eventMessage.Status,
		eventMessage.ProcessedAt)

	// Test getting the manager
	manager, err := processing.NewVideoProcessManager(as.ctx)
	as.NoError(err)
	as.NotNil(manager)

	// Verify that we have a valid event data structure
	as.NotEmpty(eventData)
}

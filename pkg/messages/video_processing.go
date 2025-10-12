package messages

import (
	"time"

	"github.com/gofrs/uuid"
)

// VideoStatus represents the status of a video
type VideoStatus string

const (
	VideoStatusProcessing VideoStatus = "processing"
	VideoStatusReady      VideoStatus = "ready"
	VideoStatusFailed     VideoStatus = "failed"
)

type VideoProcessingProgress struct {
	VideoID     uuid.UUID   `json:"video_id"`
	Status      VideoStatus `json:"status"`
	ObjectKey   string      `json:"object_key"`
	ProcessedAt time.Time   `json:"processed_at"`
}

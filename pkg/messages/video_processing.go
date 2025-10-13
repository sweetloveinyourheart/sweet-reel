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

type VideoProcessedType string

const (
	VideoProcessedTypeManifest  VideoProcessedType = "manifest"
	VideoProcessedTypeThumbnail VideoProcessedType = "thumbnail"
	VideoProcessedTypeVariant   VideoProcessedType = "variant"
)

type VideoProcessed struct {
	VideoID   uuid.UUID          `json:"video_id"`
	ObjectKey string             `json:"object_key"`
	Type      VideoProcessedType `json:"type"`
	Data      any                `json:"data"`
}

type VideoProcessedManifestData struct {
	SizeBytes int64 `json:"size_bytes"`
}

type VideoProcessedThumbnailData struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type VideoProcessedVariantData struct {
	Quality       string `json:"quality"`
	TotalSegments int    `json:"total_segments"`
	TotalDuration int    `json:"total_duration"`
}

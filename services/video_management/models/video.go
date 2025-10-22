package models

import (
	"errors"
	"strings"
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

// Video represents the main video metadata
type Video struct {
	ID          uuid.UUID   `json:"id"`
	UploaderID  uuid.UUID   `json:"uploader_id"`
	Title       string      `json:"title"`
	Description *string     `json:"description"`
	Status      VideoStatus `json:"status"`
	ObjectKey   *string     `json:"object_key"`
	ProcessedAt *time.Time  `json:"processed_at"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// GetID returns the ID of the video
func (v Video) GetID() uuid.UUID {
	return v.ID
}

// GetUploaderID returns the uploader ID of the video
func (v Video) GetUploaderID() uuid.UUID {
	return v.UploaderID
}

// GetTitle returns the title of the video
func (v Video) GetTitle() string {
	return v.Title
}

// GetDescription returns the description of the video
func (v Video) GetDescription() string {
	if v.Description == nil {
		return ""
	}
	return *v.Description
}

// GetStatus returns the status of the video
func (v Video) GetStatus() VideoStatus {
	return v.Status
}

// GetObjectKey returns the object key of the video
func (v Video) GetObjectKey() string {
	if v.ObjectKey == nil {
		return ""
	}
	return *v.ObjectKey
}

// GetProcessedAt returns the processed timestamp of the video
func (v Video) GetProcessedAt() time.Time {
	if v.ProcessedAt == nil {
		return time.Time{}
	}
	return *v.ProcessedAt
}

// GetCreatedAt returns the created timestamp of the video
func (v Video) GetCreatedAt() time.Time {
	return v.CreatedAt
}

// GetUpdatedAt returns the updated timestamp of the video
func (v Video) GetUpdatedAt() time.Time {
	return v.UpdatedAt
}

// Validate validates the required fields of the video
func (v Video) Validate() error {
	if v.ID == uuid.Nil {
		return errors.New("video ID is required")
	}

	if v.UploaderID == uuid.Nil {
		return errors.New("uploader ID is required")
	}

	if strings.TrimSpace(v.Title) == "" {
		return errors.New("video title is required and cannot be empty")
	}

	if len(v.Title) > 255 {
		return errors.New("video title cannot exceed 255 characters")
	}

	// Validate status
	validStatuses := map[VideoStatus]bool{
		VideoStatusProcessing: true,
		VideoStatusReady:      true,
		VideoStatusFailed:     true,
	}
	if !validStatuses[v.Status] {
		return errors.New("invalid video status")
	}

	return nil
}

type VideoWithThumbnails struct {
	Video
	Thumbnails []VideoThumbnail `json:"thumbnails"`
}

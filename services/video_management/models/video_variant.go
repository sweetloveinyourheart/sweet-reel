package models

import (
	"errors"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

// VideoVariant represents a video variant (quality level)
type VideoVariant struct {
	ID            uuid.UUID `json:"id"`
	VideoID       uuid.UUID `json:"video_id"`
	Quality       string    `json:"quality"`
	PlaylistURL   string    `json:"playlist_url"`
	TotalSegments *int      `json:"total_segments"`
	TotalDuration *int      `json:"total_duration"`
	CreatedAt     time.Time `json:"created_at"`
}

// GetID returns the ID of the video variant
func (vv VideoVariant) GetID() uuid.UUID {
	return vv.ID
}

// GetVideoID returns the video ID of the video variant
func (vv VideoVariant) GetVideoID() uuid.UUID {
	return vv.VideoID
}

// GetQuality returns the quality of the video variant
func (vv VideoVariant) GetQuality() string {
	return vv.Quality
}

// GetPlaylistURL returns the playlist URL of the video variant
func (vv VideoVariant) GetPlaylistURL() string {
	return vv.PlaylistURL
}

// GetTotalSegments returns the total segments pointer of the video variant
func (vv VideoVariant) GetTotalSegments() *int {
	return vv.TotalSegments
}

// GetTotalSegmentsOrDefault returns the total segments of the video variant or 0 if nil
func (vv VideoVariant) GetTotalSegmentsOrDefault() int {
	if vv.TotalSegments == nil {
		return 0
	}
	return *vv.TotalSegments
}

// GetTotalDuration returns the total duration pointer of the video variant
func (vv VideoVariant) GetTotalDuration() *int {
	return vv.TotalDuration
}

// GetTotalDurationOrDefault returns the total duration of the video variant or 0 if nil
func (vv VideoVariant) GetTotalDurationOrDefault() int {
	if vv.TotalDuration == nil {
		return 0
	}
	return *vv.TotalDuration
}

// GetCreatedAt returns the created timestamp of the video variant
func (vv VideoVariant) GetCreatedAt() time.Time {
	return vv.CreatedAt
}

// Validate validates the required fields of the video variant
func (vv VideoVariant) Validate() error {
	if vv.ID == uuid.Nil {
		return errors.New("video variant ID is required")
	}

	if vv.VideoID == uuid.Nil {
		return errors.New("video ID is required")
	}

	if strings.TrimSpace(vv.Quality) == "" {
		return errors.New("quality is required and cannot be empty")
	}

	if len(vv.Quality) > 50 {
		return errors.New("quality cannot exceed 50 characters")
	}

	if strings.TrimSpace(vv.PlaylistURL) == "" {
		return errors.New("playlist URL is required and cannot be empty")
	}

	if vv.TotalSegments != nil && *vv.TotalSegments < 0 {
		return errors.New("total segments cannot be negative")
	}

	if vv.TotalDuration != nil && *vv.TotalDuration < 0 {
		return errors.New("total duration cannot be negative")
	}

	return nil
}

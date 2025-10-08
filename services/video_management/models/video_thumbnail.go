package models

import (
	"errors"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

// VideoThumbnail represents a video thumbnail
type VideoThumbnail struct {
	ID        uuid.UUID `json:"id"`
	VideoID   uuid.UUID `json:"video_id"`
	FileURL   string    `json:"file_url"`
	Width     *int      `json:"width"`
	Height    *int      `json:"height"`
	CreatedAt time.Time `json:"created_at"`
}

// GetID returns the ID of the video thumbnail
func (vt VideoThumbnail) GetID() uuid.UUID {
	return vt.ID
}

// GetVideoID returns the video ID of the video thumbnail
func (vt VideoThumbnail) GetVideoID() uuid.UUID {
	return vt.VideoID
}

// GetFileURL returns the file URL of the video thumbnail
func (vt VideoThumbnail) GetFileURL() string {
	return vt.FileURL
}

// GetWidth returns the width pointer of the video thumbnail
func (vt VideoThumbnail) GetWidth() *int {
	return vt.Width
}

// GetWidthOrDefault returns the width of the video thumbnail or 0 if nil
func (vt VideoThumbnail) GetWidthOrDefault() int {
	if vt.Width == nil {
		return 0
	}
	return *vt.Width
}

// GetHeight returns the height pointer of the video thumbnail
func (vt VideoThumbnail) GetHeight() *int {
	return vt.Height
}

// GetHeightOrDefault returns the height of the video thumbnail or 0 if nil
func (vt VideoThumbnail) GetHeightOrDefault() int {
	if vt.Height == nil {
		return 0
	}
	return *vt.Height
}

// GetCreatedAt returns the created timestamp of the video thumbnail
func (vt VideoThumbnail) GetCreatedAt() time.Time {
	return vt.CreatedAt
}

// Validate validates the required fields of the video thumbnail
func (vt VideoThumbnail) Validate() error {
	if vt.ID == uuid.Nil {
		return errors.New("video thumbnail ID is required")
	}

	if vt.VideoID == uuid.Nil {
		return errors.New("video ID is required")
	}

	if strings.TrimSpace(vt.FileURL) == "" {
		return errors.New("file URL is required and cannot be empty")
	}

	if vt.Width != nil && *vt.Width <= 0 {
		return errors.New("width must be positive if specified")
	}

	if vt.Height != nil && *vt.Height <= 0 {
		return errors.New("height must be positive if specified")
	}

	return nil
}

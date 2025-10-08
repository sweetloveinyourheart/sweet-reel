package models

import (
	"errors"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

// VideoManifest represents a video manifest (HLS master playlist)
type VideoManifest struct {
	ID          uuid.UUID `json:"id"`
	VideoID     uuid.UUID `json:"video_id"`
	ManifestURL string    `json:"manifest_url"`
	SizeBytes   *int64    `json:"size_bytes"`
	CreatedAt   time.Time `json:"created_at"`
}

// GetID returns the ID of the video manifest
func (vm VideoManifest) GetID() uuid.UUID {
	return vm.ID
}

// GetVideoID returns the video ID of the video manifest
func (vm VideoManifest) GetVideoID() uuid.UUID {
	return vm.VideoID
}

// GetManifestURL returns the manifest URL of the video manifest
func (vm VideoManifest) GetManifestURL() string {
	return vm.ManifestURL
}

// GetSizeBytes returns the size in bytes pointer of the video manifest
func (vm VideoManifest) GetSizeBytes() *int64 {
	return vm.SizeBytes
}

// GetSizeBytesOrDefault returns the size in bytes of the video manifest or 0 if nil
func (vm VideoManifest) GetSizeBytesOrDefault() int64 {
	if vm.SizeBytes == nil {
		return 0
	}
	return *vm.SizeBytes
}

// GetCreatedAt returns the created timestamp of the video manifest
func (vm VideoManifest) GetCreatedAt() time.Time {
	return vm.CreatedAt
}

// Validate validates the required fields of the video manifest
func (vm VideoManifest) Validate() error {
	if vm.ID == uuid.Nil {
		return errors.New("video manifest ID is required")
	}

	if vm.VideoID == uuid.Nil {
		return errors.New("video ID is required")
	}

	if strings.TrimSpace(vm.ManifestURL) == "" {
		return errors.New("manifest URL is required and cannot be empty")
	}

	if vm.SizeBytes != nil && *vm.SizeBytes < 0 {
		return errors.New("size bytes cannot be negative")
	}

	return nil
}

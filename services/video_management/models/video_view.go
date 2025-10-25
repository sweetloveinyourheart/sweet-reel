package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// VideoView represents a single view event for a video
type VideoView struct {
	ID            uuid.UUID  `db:"id" json:"id"`
	VideoID       uuid.UUID  `db:"video_id" json:"video_id"`
	ViewerID      *uuid.UUID `db:"viewer_id" json:"viewer_id,omitempty"` // NULL for anonymous viewers
	ViewedAt      time.Time  `db:"viewed_at" json:"viewed_at"`
	WatchDuration *int       `db:"watch_duration" json:"watch_duration,omitempty"` // seconds watched
	IPAddress     *string    `db:"ip_address" json:"ip_address,omitempty"`
	UserAgent     *string    `db:"user_agent" json:"user_agent,omitempty"`
}

// VideoViewStats represents aggregated view statistics for a video
type VideoViewStats struct {
	VideoID       uuid.UUID `json:"video_id"`
	TotalViews    int64     `json:"total_views"`
	UniqueViewers int64     `json:"unique_viewers"` // count of distinct viewer_id
	AvgDuration   float64   `json:"avg_duration"`   // average watch duration in seconds
}

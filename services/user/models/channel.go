package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// Channel represents a user's channel for publishing videos
type Channel struct {
	ID              uuid.UUID `db:"id" json:"id"`
	OwnerID         uuid.UUID `db:"owner_id" json:"owner_id"`
	Name            string    `db:"name" json:"name"`
	Handle          string    `db:"handle" json:"handle"` // e.g., @username
	Description     *string   `db:"description" json:"description,omitempty"`
	BannerURL       *string   `db:"banner_url" json:"banner_url,omitempty"`
	SubscriberCount int       `db:"subscriber_count" json:"subscriber_count"`
	TotalViews      int64     `db:"total_views" json:"total_views"`
	TotalVideos     int       `db:"total_videos" json:"total_videos"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

// ChannelWithOwner represents a channel joined with user information
type ChannelWithOwner struct {
	Channel
	Owner User `json:"owner"`
}

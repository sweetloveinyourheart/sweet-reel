package models

type ChannelVideo struct {
	Video
	TotalDuration      int    `json:"total_duration"`
	ThumbnailObjectKey string `json:"thumbnail_object_key"`
}

// GetTotalDuration returns the total duration of the video
func (v ChannelVideo) GetTotalDuration() int {
	return v.TotalDuration
}

// GetThumbnailObjectKey returns the thumbnail object key of the video
func (v ChannelVideo) GetThumbnailObjectKey() string {
	return v.ThumbnailObjectKey
}

type VideoMetadata struct {
	Video
	AvailableQualities []string `json:"available_qualities"`
}

// AvailableQualities returns the available qualities of the video
func (v VideoMetadata) GetAvailableQualities() []string {
	return v.AvailableQualities
}

type VideoPlaylist struct {
}

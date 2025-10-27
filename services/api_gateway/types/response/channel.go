package response

type ChannelResponse struct {
	ID              string       `json:"id"`
	OwnerID         string       `json:"owner_id"`
	Name            string       `json:"name"`
	Handle          string       `json:"handle"`
	Description     string       `json:"description"`
	BannerURL       string       `json:"banner_url"`
	SubscriberCount int32        `json:"subscriber_count"`
	TotalViews      int64        `json:"total_views"`
	TotalVideos     int32        `json:"total_videos"`
	CreatedAt       string       `json:"created_at"`
	UpdatedAt       string       `json:"updated_at"`
	Owner           UserResponse `json:"owner"`
}

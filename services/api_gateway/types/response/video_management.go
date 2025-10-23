package response

type PresignedUrlResponse struct {
	VideoId      string `json:"video_id"`
	PresignedUrl string `json:"presigned_url"`
	ExpiresIn    int32  `json:"expires_in"`
}

type UserVideoResponse struct {
	VideoID      string `json:"video_id"`
	Title        string `json:"title"`
	ThumbnailUrl string `json:"thumbnail_url"`
}

type GetUserVideosResponse struct {
	Videos []UserVideoResponse `json:"videos"`
}

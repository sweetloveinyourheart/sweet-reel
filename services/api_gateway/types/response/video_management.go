package response

type PresignedUrlResponse struct {
	VideoId      string `json:"video_id"`
	PresignedUrl string `json:"presigned_url"`
	ExpiresIn    int32  `json:"expires_in"`
}

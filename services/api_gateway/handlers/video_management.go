package handlers

import (
	"net/http"

	"github.com/sweetloveinyourheart/sweet-reel/services/api_gateway/helpers"
)

type IVideoManagementHandler interface {
	GeneratePresignedURL(w http.ResponseWriter, r *http.Request)
}

type VideoManagementHandler struct {
}

func NewVideoManagementHandler() IVideoManagementHandler {
	return &VideoManagementHandler{}
}

// GeneratePresignedURL handles POST /api/v1/videos/presigned-url
func (h *VideoManagementHandler) GeneratePresignedURL(w http.ResponseWriter, r *http.Request) {
	helpers.WriteJSONSuccess(w, map[string]string{
		"message": "Presigned URL generation not implemented yet",
	})
}

package handlers

import (
	"net/http"
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
func (h *VideoManagementHandler) GeneratePresignedURL(w http.ResponseWriter, r *http.Request) {}

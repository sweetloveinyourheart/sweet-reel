package handlers

import (
	"net/http"

	"connectrpc.com/connect"
	"github.com/samber/do"
	"go.uber.org/zap"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"

	videoManagementProto "github.com/sweetloveinyourheart/sweet-reel/proto/code/video_management/go"
	videoManagementConnect "github.com/sweetloveinyourheart/sweet-reel/proto/code/video_management/go/grpcconnect"
	"github.com/sweetloveinyourheart/sweet-reel/services/api_gateway/errors"
	"github.com/sweetloveinyourheart/sweet-reel/services/api_gateway/helpers"
	"github.com/sweetloveinyourheart/sweet-reel/services/api_gateway/types/request"
	"github.com/sweetloveinyourheart/sweet-reel/services/api_gateway/types/response"
)

type IVideoManagementHandler interface {
	GeneratePresignedURL(w http.ResponseWriter, r *http.Request)
}

type VideoManagementHandler struct {
	videoManagementServiceClient videoManagementConnect.VideoManagementClient
}

func NewVideoManagementHandler() IVideoManagementHandler {
	videoManagementServiceClient, err := do.Invoke[videoManagementConnect.VideoManagementClient](nil)
	if err != nil {
		logger.Global().Fatal("unable to get video management client")
	}

	return &VideoManagementHandler{
		videoManagementServiceClient: videoManagementServiceClient,
	}
}

// GeneratePresignedURL handles POST /api/v1/videos/presigned-url
func (h *VideoManagementHandler) GeneratePresignedURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := helpers.GetUserID(r)

	var body request.PresignedUrlRequestBody
	err := helpers.ParseJSONBody(r, &body)
	if err != nil {
		helpers.WriteErrorResponse(w, err)
		return
	}

	presignedUrlReq := connect.NewRequest(&videoManagementProto.PresignedUrlRequest{
		Title:       body.Title,
		Description: body.Description,
		FileName:    body.FileName,
		UploaderId:  userID,
	})

	presignedUrlRes, err := h.videoManagementServiceClient.PresignedUrl(ctx, presignedUrlReq)
	if err != nil {
		logger.Global().Error("error performing pre-signed url request", zap.Error(err))
		helpers.WriteErrorResponse(w, errors.ErrHTTPInternalServer)
		return
	}

	// Build response
	responseData := response.PresignedUrlResponse{
		VideoId:      presignedUrlRes.Msg.GetVideoId(),
		PresignedUrl: presignedUrlRes.Msg.GetPresignedUrl(),
		ExpiresIn:    presignedUrlRes.Msg.GetExpiresIn(),
	}

	helpers.WriteJSONSuccess(w, responseData)
}

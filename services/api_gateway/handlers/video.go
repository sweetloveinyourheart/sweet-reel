package handlers

import (
	"net/http"

	"connectrpc.com/connect"
	"github.com/samber/do"
	"go.uber.org/zap"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"

	userProto "github.com/sweetloveinyourheart/sweet-reel/proto/code/user/go"
	userConnect "github.com/sweetloveinyourheart/sweet-reel/proto/code/user/go/grpcconnect"
	videoManagementProto "github.com/sweetloveinyourheart/sweet-reel/proto/code/video_management/go"
	videoManagementConnect "github.com/sweetloveinyourheart/sweet-reel/proto/code/video_management/go/grpcconnect"
	"github.com/sweetloveinyourheart/sweet-reel/services/api_gateway/errors"
	"github.com/sweetloveinyourheart/sweet-reel/services/api_gateway/helpers"
	"github.com/sweetloveinyourheart/sweet-reel/services/api_gateway/types/request"
	"github.com/sweetloveinyourheart/sweet-reel/services/api_gateway/types/response"
)

type IVideoHandler interface {
	GeneratePresignedURL(w http.ResponseWriter, r *http.Request)
	GetVideoMetadata(w http.ResponseWriter, r *http.Request)
	ServePlaylist(w http.ResponseWriter, r *http.Request)
}

type VideoHandler struct {
	userServiceClient            userConnect.UserServiceClient
	videoManagementServiceClient videoManagementConnect.VideoManagementClient
}

func NewVideoHandler() IVideoHandler {
	userServiceClient, err := do.Invoke[userConnect.UserServiceClient](nil)
	if err != nil {
		logger.Global().Fatal("unable to get user service client")
	}

	videoManagementServiceClient, err := do.Invoke[videoManagementConnect.VideoManagementClient](nil)
	if err != nil {
		logger.Global().Fatal("unable to get video management client")
	}

	return &VideoHandler{
		videoManagementServiceClient: videoManagementServiceClient,
		userServiceClient:            userServiceClient,
	}
}

// GeneratePresignedURL handles POST /api/v1/videos/presigned-url
func (h *VideoHandler) GeneratePresignedURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := helpers.GetUserID(r)

	var body request.PresignedUrlRequestBody
	err := helpers.ParseJSONBody(r, &body)
	if err != nil {
		helpers.WriteErrorResponse(w, err)
		return
	}

	presignedUrlReq := connect.NewRequest(&videoManagementProto.PresignedUrlRequest{
		ChannelId:   body.ChannelID,
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

// GetVideoMetadata handles GET /api/v1/videos/{video_id}/metadata
func (h *VideoHandler) GetVideoMetadata(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get videoID from URL path parameter
	videoID := r.PathValue("video_id")

	getMetadataReq := connect.NewRequest(&videoManagementProto.GetVideoMetadataByIdRequest{VideoId: videoID})
	getMetadataResp, err := h.videoManagementServiceClient.GetVideoMetadataById(ctx, getMetadataReq)
	if err != nil {
		logger.Global().Error("error getting video metadata", zap.Error(err))
		helpers.WriteErrorResponse(w, errors.ErrNotFound)
		return
	}

	getChannelReq := connect.NewRequest(&userProto.GetChannelByIDRequest{VideoId: getMetadataResp.Msg.GetChannelId()})
	getChannelresp, err := h.userServiceClient.GetChannelByID(ctx, getChannelReq)
	if err != nil {
		logger.Global().Error("failed to get channel by handle", zap.Error(err))
		helpers.WriteErrorResponse(w, errors.ErrHTTPInternalServer)
		return
	}

	// Build response
	responseData := response.GetVideoMetadataResponse{
		VideoTitle:         getMetadataResp.Msg.GetVideoTitle(),
		VideoDescription:   getMetadataResp.Msg.GetVideoDescription(),
		TotalView:          getMetadataResp.Msg.GetTotalView(),
		AvailableQualities: getMetadataResp.Msg.GetAvailableQualities(),
		ProcessedAt:        getMetadataResp.Msg.GetProcessedAt(),
		Channel: response.ChannelMetadata{
			Name:   getChannelresp.Msg.GetChannel().GetName(),
			Handle: getChannelresp.Msg.GetChannel().GetHandle(),
			OwnerMetadata: response.UserMetadata{
				Email:   getChannelresp.Msg.GetOwner().GetEmail(),
				Name:    getChannelresp.Msg.GetOwner().GetName(),
				Picture: getChannelresp.Msg.GetOwner().GetPicture(),
			},
		},
	}

	helpers.WriteJSONSuccess(w, responseData)

}

// Handles both master and variant playlists
// Example routes:
//
//	GET /videos/{video_id}/hls
//	GET /videos/{video_id}/hls/{quality}
func (h *VideoHandler) ServePlaylist(w http.ResponseWriter, r *http.Request) {

}

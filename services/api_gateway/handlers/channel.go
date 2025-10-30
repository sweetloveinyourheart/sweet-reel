package handlers

import (
	"net/http"
	"strconv"

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
	"github.com/sweetloveinyourheart/sweet-reel/services/api_gateway/types/response"
)

type IChannelHandler interface {
	GetChannel(w http.ResponseWriter, r *http.Request)
	GetChannelVideos(w http.ResponseWriter, r *http.Request)
	GetChannelByHandle(w http.ResponseWriter, r *http.Request)
	GetChannelVideosByHandle(w http.ResponseWriter, r *http.Request)
}

type ChannelHandler struct {
	userServiceClient            userConnect.UserServiceClient
	videoManagementServiceClient videoManagementConnect.VideoManagementClient
}

func NewChannelHandler() IChannelHandler {
	userServiceClient, err := do.Invoke[userConnect.UserServiceClient](nil)
	if err != nil {
		logger.Global().Fatal("unable to get user service client")
	}

	videoManagementServiceClient, err := do.Invoke[videoManagementConnect.VideoManagementClient](nil)
	if err != nil {
		logger.Global().Fatal("unable to get user service client")
	}

	return &ChannelHandler{
		userServiceClient:            userServiceClient,
		videoManagementServiceClient: videoManagementServiceClient,
	}
}

// GetChannelByHandle handles GET /api/v1/channels/{handle}
func (h *ChannelHandler) GetChannelByHandle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get handle from URL path parameter
	handle := r.PathValue("handle")

	if handle == "" {
		helpers.WriteErrorResponse(w, errors.NewHTTPError(
			http.StatusBadRequest,
			"handle is required",
			"INVALID_HANDLE",
		))
		return
	}

	// Call user service
	req := connect.NewRequest(&userProto.GetChannelByHandleRequest{
		Handle: handle,
	})

	resp, err := h.userServiceClient.GetChannelByHandle(ctx, req)
	if err != nil {
		logger.Global().Error("failed to get channel by handle", zap.Error(err))
		helpers.WriteErrorResponse(w, errors.ErrHTTPInternalServer)
		return
	}

	channel := resp.Msg.GetChannel()
	owner := resp.Msg.GetOwner()

	// Build response
	channelResponse := response.ChannelResponse{
		ID:              channel.GetId(),
		OwnerID:         channel.GetOwnerId(),
		Name:            channel.GetName(),
		Handle:          channel.GetHandle(),
		Description:     channel.GetDescription(),
		BannerURL:       channel.GetBannerUrl(),
		SubscriberCount: channel.GetSubscriberCount(),
		TotalViews:      channel.GetTotalViews(),
		TotalVideos:     channel.GetTotalVideos(),
		CreatedAt:       channel.GetCreatedAt(),
		UpdatedAt:       channel.GetUpdatedAt(),
		Owner: response.UserResponse{
			ID:        owner.GetId(),
			Email:     owner.GetEmail(),
			Name:      owner.GetName(),
			Picture:   owner.GetPicture(),
			CreatedAt: owner.GetCreatedAt(),
			UpdatedAt: owner.GetUpdatedAt(),
		},
	}

	helpers.WriteJSONSuccess(w, channelResponse)
}

// GetChannelVideosByHandle handles GET /api/v1/channels/videos/{handle}
func (h *ChannelHandler) GetChannelVideosByHandle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get handle from URL path parameter
	handle := r.PathValue("handle")

	if handle == "" {
		helpers.WriteErrorResponse(w, errors.NewHTTPError(
			http.StatusBadRequest,
			"handle is required",
			"INVALID_HANDLE",
		))
		return
	}

	limit := r.URL.Query().Get("limit")
	limitBy, err := strconv.ParseInt(limit, 0, 32)
	if err != nil {
		helpers.WriteErrorResponse(w, errors.NewHTTPError(http.StatusBadRequest, "limit is not valid", "INVALID_ARGUMENTS"))
		return
	}

	offset := r.URL.Query().Get("offset")
	offsetBy, err := strconv.ParseInt(offset, 0, 32)
	if err != nil {
		helpers.WriteErrorResponse(w, errors.NewHTTPError(http.StatusBadRequest, "offset is not valid", "INVALID_ARGUMENTS"))
		return
	}

	// Call user service
	req := connect.NewRequest(&userProto.GetChannelByHandleRequest{
		Handle: handle,
	})

	resp, err := h.userServiceClient.GetChannelByHandle(ctx, req)
	if err != nil {
		logger.Global().Error("failed to get channel by handle", zap.Error(err))
		helpers.WriteErrorResponse(w, errors.ErrHTTPInternalServer)
		return
	}

	getVideosReq := connect.NewRequest(&videoManagementProto.GetChannelVideosRequest{
		ChannelId: resp.Msg.GetChannel().GetId(),
		Limit:     int32(limitBy),
		Offset:    int32(offsetBy),
	})

	channelVideosRes, err := h.videoManagementServiceClient.GetChannelVideos(ctx, getVideosReq)
	if err != nil {
		logger.Global().Error("error performing get channel videos request", zap.Error(err))
		helpers.WriteErrorResponse(w, errors.ErrHTTPInternalServer)
		return
	}

	// Build response
	videos := make([]response.UserVideoResponse, 0)
	for _, video := range channelVideosRes.Msg.GetVideos() {
		videos = append(videos, response.UserVideoResponse{
			VideoID:       video.GetVideoId(),
			Title:         video.GetVideoTitle(),
			ThumbnailUrl:  video.GetThumbnailUrl(),
			TotalDuration: video.GetTotalDuration(),
			ProcessedAt:   video.GetProcessedAt(),
		})
	}

	responseData := response.GetChannelVideosResponse{
		Videos: videos,
	}

	helpers.WriteJSONSuccess(w, responseData)
}

// GetChannel handles GET /api/v1/channels
func (h *ChannelHandler) GetChannel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := helpers.GetUserID(r)

	// Call user service
	req := connect.NewRequest(&userProto.GetChannelByUserRequest{
		UserId: userID,
	})

	resp, err := h.userServiceClient.GetChannelByUser(ctx, req)
	if err != nil {
		logger.Global().Error("failed to get channel by user", zap.Error(err))
		helpers.WriteErrorResponse(w, errors.ErrHTTPInternalServer)
		return
	}

	channel := resp.Msg.GetChannel()
	owner := resp.Msg.GetOwner()

	// Build response
	channelResponse := response.ChannelResponse{
		ID:              channel.GetId(),
		OwnerID:         channel.GetOwnerId(),
		Name:            channel.GetName(),
		Handle:          channel.GetHandle(),
		Description:     channel.GetDescription(),
		BannerURL:       channel.GetBannerUrl(),
		SubscriberCount: channel.GetSubscriberCount(),
		TotalViews:      channel.GetTotalViews(),
		TotalVideos:     channel.GetTotalVideos(),
		CreatedAt:       channel.GetCreatedAt(),
		UpdatedAt:       channel.GetUpdatedAt(),
		Owner: response.UserResponse{
			ID:        owner.GetId(),
			Email:     owner.GetEmail(),
			Name:      owner.GetName(),
			Picture:   owner.GetPicture(),
			CreatedAt: owner.GetCreatedAt(),
			UpdatedAt: owner.GetUpdatedAt(),
		},
	}

	helpers.WriteJSONSuccess(w, channelResponse)
}

// GetChannelVideos handles GET /api/v1/channels/videos
func (h *ChannelHandler) GetChannelVideos(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := helpers.GetUserID(r)

	limit := r.URL.Query().Get("limit")
	limitBy, err := strconv.ParseInt(limit, 0, 32)
	if err != nil {
		helpers.WriteErrorResponse(w, errors.NewHTTPError(http.StatusBadRequest, "limit is not valid", "INVALID_ARGUMENTS"))
		return
	}

	offset := r.URL.Query().Get("offset")
	offsetBy, err := strconv.ParseInt(offset, 0, 32)
	if err != nil {
		helpers.WriteErrorResponse(w, errors.NewHTTPError(http.StatusBadRequest, "offset is not valid", "INVALID_ARGUMENTS"))
		return
	}

	// Call user service
	req := connect.NewRequest(&userProto.GetChannelByUserRequest{
		UserId: userID,
	})

	resp, err := h.userServiceClient.GetChannelByUser(ctx, req)
	if err != nil {
		logger.Global().Error("failed to get channel by user", zap.Error(err))
		helpers.WriteErrorResponse(w, errors.ErrHTTPInternalServer)
		return
	}

	getVideosReq := connect.NewRequest(&videoManagementProto.GetChannelVideosRequest{
		ChannelId: resp.Msg.GetChannel().GetId(),
		Limit:     int32(limitBy),
		Offset:    int32(offsetBy),
	})

	channelVideosRes, err := h.videoManagementServiceClient.GetChannelVideos(ctx, getVideosReq)
	if err != nil {
		logger.Global().Error("error performing get channel videos request", zap.Error(err))
		helpers.WriteErrorResponse(w, errors.ErrHTTPInternalServer)
		return
	}

	// Build response
	videos := make([]response.UserVideoResponse, 0)
	for _, video := range channelVideosRes.Msg.GetVideos() {
		videos = append(videos, response.UserVideoResponse{
			VideoID:       video.GetVideoId(),
			Title:         video.GetVideoTitle(),
			ThumbnailUrl:  video.GetThumbnailUrl(),
			TotalDuration: video.GetTotalDuration(),
			ProcessedAt:   video.GetProcessedAt(),
		})
	}

	responseData := response.GetChannelVideosResponse{
		Videos: videos,
	}

	helpers.WriteJSONSuccess(w, responseData)
}

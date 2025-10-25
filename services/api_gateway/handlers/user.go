package handlers

import (
	"net/http"

	"connectrpc.com/connect"
	"github.com/samber/do"
	"go.uber.org/zap"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
	userProto "github.com/sweetloveinyourheart/sweet-reel/proto/code/user/go"
	userConnect "github.com/sweetloveinyourheart/sweet-reel/proto/code/user/go/grpcconnect"
	"github.com/sweetloveinyourheart/sweet-reel/services/api_gateway/errors"
	"github.com/sweetloveinyourheart/sweet-reel/services/api_gateway/helpers"
	"github.com/sweetloveinyourheart/sweet-reel/services/api_gateway/types/response"
)

type IUserHandler interface {
	GetChannelByHandle(w http.ResponseWriter, r *http.Request)
}

type UserHandler struct {
	userServiceClient userConnect.UserServiceClient
}

func NewUserHandler() IUserHandler {
	userServiceClient, err := do.Invoke[userConnect.UserServiceClient](nil)
	if err != nil {
		logger.Global().Fatal("unable to get user service client")
	}

	return &UserHandler{
		userServiceClient: userServiceClient,
	}
}

// GetChannelByHandle handles GET /api/v1/channels/{handle}
func (h *UserHandler) GetChannelByHandle(w http.ResponseWriter, r *http.Request) {
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

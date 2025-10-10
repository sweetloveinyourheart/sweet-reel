package handlers

import (
	"net/http"

	"connectrpc.com/connect"
	"github.com/samber/do"
	"go.uber.org/zap"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/oauth2"
	authProto "github.com/sweetloveinyourheart/sweet-reel/proto/code/auth/go"
	authConnect "github.com/sweetloveinyourheart/sweet-reel/proto/code/auth/go/grpcconnect"
	"github.com/sweetloveinyourheart/sweet-reel/services/api_gateway/errors"
	"github.com/sweetloveinyourheart/sweet-reel/services/api_gateway/helpers"
	"github.com/sweetloveinyourheart/sweet-reel/services/api_gateway/types/request"
	"github.com/sweetloveinyourheart/sweet-reel/services/api_gateway/types/response"
)

type IAuthHandler interface {
	GoogleOAuth(w http.ResponseWriter, r *http.Request)
}

type AuthHandler struct {
	authServiceClient authConnect.AuthServiceClient
}

func NewAuthHandler() IAuthHandler {
	authServiceClient, err := do.Invoke[authConnect.AuthServiceClient](nil)
	if err != nil {
		logger.Global().Fatal("unable to get auth server client")
	}

	return &AuthHandler{
		authServiceClient: authServiceClient,
	}
}

// GoogleOAuth handles POST /api/v1/oauth/google
func (h *AuthHandler) GoogleOAuth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse and validate request body
	var body request.GoogleOAuthRequestBody
	if err := helpers.ParseJSONBody(r, &body); err != nil {
		helpers.WriteErrorResponse(w, err)
		return
	}

	oauthRequest := &authProto.OAuthLoginRequest{
		Provider:    string(oauth2.ProviderGoogle),
		AccessToken: body.AccessToken,
	}

	oauthResponse, err := h.authServiceClient.OAuthLogin(ctx, connect.NewRequest(oauthRequest))
	if err != nil {
		logger.Global().Error("oauth login failed", zap.Error(err))

		helpers.WriteErrorResponse(w, errors.NewHTTPError(
			http.StatusUnauthorized,
			errors.ErrOAuthLoginFailed.Message,
			errors.ErrOAuthLoginFailed.Code,
		))
		return
	}

	// Build response
	responseData := response.GoogleOAuthResponse{
		JwtToken: oauthResponse.Msg.GetJwtToken(),
		User: response.GoogleOAuthUser{
			Id:        oauthResponse.Msg.User.Id,
			Email:     oauthResponse.Msg.User.Email,
			Name:      oauthResponse.Msg.User.Name,
			Picture:   oauthResponse.Msg.User.Picture,
			CreatedAt: oauthResponse.Msg.User.CreatedAt,
			UpdatedAt: oauthResponse.Msg.User.UpdatedAt,
		},
		IsNew: oauthResponse.Msg.IsNewUser,
	}

	helpers.WriteJSONSuccess(w, responseData)
}

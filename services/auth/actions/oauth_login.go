package actions

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/oauth2"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/grpc"
	oauth2Pkg "github.com/sweetloveinyourheart/sweet-reel/pkg/oauth2"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/stringsutil"
	proto "github.com/sweetloveinyourheart/sweet-reel/proto/code/auth/go"
	userProto "github.com/sweetloveinyourheart/sweet-reel/proto/code/user/go"

	"connectrpc.com/connect"
)

func (a *actions) OAuthLogin(ctx context.Context, request *connect.Request[proto.OAuthLoginRequest]) (*connect.Response[proto.OAuthLoginResponse], error) {
	accessToken := request.Msg.GetAccessToken()
	if stringsutil.IsBlank(accessToken) {
		return nil, grpc.InvalidArgumentError(errors.New("access token is empty"))
	}

	provider := request.Msg.GetProvider()
	if stringsutil.IsBlank(accessToken) {
		return nil, grpc.InvalidArgumentError(errors.New("provider is empty"))
	}

	var userInfo *oauth2Pkg.UserInfo
	var err error

	switch provider {
	case string(oauth2Pkg.ProviderGoogle):
		token := &oauth2.Token{
			AccessToken: accessToken,
		}
		userInfo, err = a.googleOAuthClient.GetUserInfo(ctx, token)
		if err != nil {
			return nil, grpc.InternalError(errors.New("user info is empty"))
		}

		if userInfo == nil {
			return nil, grpc.NotFoundError(errors.New("user profile is not found"))
		}
	default:
		return nil, grpc.InvalidArgumentError(errors.New("provider is not supported"))
	}

	upsertRequest := &userProto.UpsertOAuthUserRequest{
		Provider:       request.Msg.GetProvider(),
		ProviderUserId: userInfo.ID,
		Email:          userInfo.Email,
		Name:           userInfo.Name,
		Picture:        userInfo.Picture,
	}
	upsertResponse, err := a.userServerClient.UpsertOAuthUser(ctx, connect.NewRequest(upsertRequest))
	if err != nil {
		return nil, grpc.InternalError(err)
	}

	// Sign token with claims
	user := upsertResponse.Msg.GetUser()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.Id,
		"email":   user.Email,
	})

	tokenString, err := token.SignedString([]byte(a.signingToken))
	if err != nil {
		return nil, grpc.InternalError(err)
	}

	response := &proto.OAuthLoginResponse{
		User: &proto.User{
			Id:        user.Id,
			Email:     user.Email,
			Name:      user.Name,
			Picture:   user.Picture,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		JwtToken:  tokenString,
		IsNewUser: upsertResponse.Msg.IsNewUser,
	}

	return connect.NewResponse(response), nil
}

func (a *actions) ValidateToken(ctx context.Context, request *connect.Request[proto.ValidateTokenRequest]) (*connect.Response[proto.ValidateTokenResponse], error) {
	return nil, nil
}

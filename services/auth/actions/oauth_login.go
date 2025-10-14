package actions

import (
	"context"
	"time"

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
			return nil, grpc.InternalError(err)
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
		"exp":     time.Now().Add(1 * time.Hour).Unix(), // 1 hour expiration
	})

	tokenString, err := token.SignedString([]byte(a.signingToken))
	if err != nil {
		return nil, grpc.InternalError(err)
	}

	// Generate refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.Id,
		"type":    "refresh",
		"exp":     time.Now().Add(30 * 24 * time.Hour).Unix(), // 30 days expiration
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(a.signingToken))
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
		JwtToken:        tokenString,
		JwtRefreshToken: refreshTokenString,
		IsNewUser:       upsertResponse.Msg.IsNewUser,
	}

	return connect.NewResponse(response), nil
}

func (a *actions) RefreshToken(ctx context.Context, request *connect.Request[proto.RefreshTokenRequest]) (*connect.Response[proto.RefreshTokenResponse], error) {
	refreshToken := request.Msg.GetJwtRefreshToken()
	if stringsutil.IsBlank(refreshToken) {
		return nil, grpc.InvalidArgumentError(errors.New("refresh token is empty"))
	}

	// Parse and validate refresh token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(a.signingToken), nil
	})

	if err != nil {
		return nil, grpc.UnauthenticatedError(errors.New("invalid refresh token"))
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, grpc.UnauthenticatedError(errors.New("invalid refresh token claims"))
	}

	// Verify this is a refresh token
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return nil, grpc.UnauthenticatedError(errors.New("invalid token type"))
	}

	// Extract user ID
	userID, ok := claims["user_id"].(string)
	if !ok || stringsutil.IsBlank(userID) {
		return nil, grpc.UnauthenticatedError(errors.New("invalid user id in token"))
	}

	// Get user information
	getUserRequest := &userProto.GetUserByIDRequest{UserId: userID}
	getUserResponse, err := a.userServerClient.GetUserByID(ctx, connect.NewRequest(getUserRequest))
	if err != nil {
		return nil, grpc.InternalError(err)
	}

	user := getUserResponse.Msg.GetUser()
	if user == nil {
		return nil, grpc.NotFoundError(errors.New("user not found"))
	}

	// Generate new access token
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.Id,
		"email":   user.Email,
		"exp":     time.Now().Add(1 * time.Hour).Unix(), // 1 hour expiration
	})

	tokenString, err := newToken.SignedString([]byte(a.signingToken))
	if err != nil {
		return nil, grpc.InternalError(err)
	}

	// Set response data
	response := &proto.RefreshTokenResponse{JwtToken: tokenString}

	return connect.NewResponse(response), nil
}

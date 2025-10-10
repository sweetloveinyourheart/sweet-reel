package actions

import (
	"context"

	"github.com/samber/do"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/interceptors"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/oauth2"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/stringsutil"
	authConnect "github.com/sweetloveinyourheart/sweet-reel/proto/code/auth/go/grpcconnect"
	userConnect "github.com/sweetloveinyourheart/sweet-reel/proto/code/user/go/grpcconnect"
)

type actions struct {
	context           context.Context
	signingToken      string
	defaultAuth       func(context.Context, string) (context.Context, error)
	googleOAuthClient oauth2.IOAuthClient
	userServerClient  userConnect.UserServiceClient
}

// AuthFuncOverride is a callback function that overrides the default authorization middleware in the GRPC layer. This is
// used to allow unauthenticated endpoints (such as login) to be called without a token.
func (a *actions) AuthFuncOverride(ctx context.Context, token string, fullMethodName string) (context.Context, error) {
	if fullMethodName == authConnect.AuthServiceOAuthLoginProcedure {
		return ctx, nil
	}

	if fullMethodName == authConnect.AuthServiceValidateTokenProcedure {
		return ctx, nil
	}

	return a.defaultAuth(ctx, token)
}

func NewActions(ctx context.Context, signingToken string) *actions {
	if stringsutil.IsBlank(signingToken) {
		logger.Global().Fatal("sigining token is empty")
	}

	googleOAuthClient, err := do.InvokeNamed[oauth2.IOAuthClient](nil, string(oauth2.ProviderGoogle))
	if err != nil {
		logger.Global().Fatal("unable to get oauth service")
	}

	userServerClient, err := do.Invoke[userConnect.UserServiceClient](nil)
	if err != nil {
		logger.Global().Fatal("unable to get user server client")
	}

	return &actions{
		context:      ctx,
		signingToken: signingToken,
		defaultAuth:  interceptors.ConnectServerAuthHandler(signingToken),

		googleOAuthClient: googleOAuthClient,
		userServerClient:  userServerClient,
	}
}

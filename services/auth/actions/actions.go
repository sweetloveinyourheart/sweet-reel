package actions

import (
	"context"

	"github.com/samber/do"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/interceptors"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/oauth2"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/stringsutil"
	userConnect "github.com/sweetloveinyourheart/sweet-reel/proto/code/user/go/grpcconnect"
)

type actions struct {
	context           context.Context
	signingToken      string
	defaultAuth       func(context.Context, string) (context.Context, error)
	googleOAuthClient oauth2.IOAuthClient
	userServerClient  userConnect.UserServiceClient
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

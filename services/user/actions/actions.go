package actions

import (
	"context"

	"github.com/samber/do"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/db"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/interceptors"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
	userConnect "github.com/sweetloveinyourheart/sweet-reel/proto/code/user/go/grpcconnect"
	"github.com/sweetloveinyourheart/sweet-reel/services/user/repos"
)

type actions struct {
	context     context.Context
	defaultAuth func(context.Context, string) (context.Context, error)
	dbConn      db.ConnPool
	userRepo    repos.IUserRepository
	channelRepo repos.IChannelRepository
}

// AuthFuncOverride is a callback function that overrides the default authorization middleware in the GRPC layer. This is
// used to allow unauthenticated endpoints (such as login) to be called without a token.
func (a *actions) AuthFuncOverride(ctx context.Context, token string, fullMethodName string) (context.Context, error) {
	if fullMethodName == userConnect.UserServiceUpsertOAuthUserProcedure {
		return ctx, nil
	}
	return a.defaultAuth(ctx, token)
}

func NewActions(ctx context.Context, signingToken string) *actions {
	userRepo, err := do.Invoke[repos.IUserRepository](nil)
	if err != nil {
		logger.Global().Fatal("unable to get user repo")
	}

	channelRepo, err := do.Invoke[repos.IChannelRepository](nil)
	if err != nil {
		logger.Global().Fatal("unable to get channel repo")
	}

	dbConn, err := do.Invoke[db.ConnPool](nil)
	if err != nil {
		logger.Global().Fatal("unable to get db conn")
	}

	return &actions{
		context:     ctx,
		defaultAuth: interceptors.ConnectServerAuthHandler(signingToken),
		dbConn:      dbConn,
		userRepo:    userRepo,
		channelRepo: channelRepo,
	}
}

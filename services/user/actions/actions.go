package actions

import (
	"context"

	"github.com/samber/do"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/db"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/interceptors"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
	"github.com/sweetloveinyourheart/sweet-reel/services/user/repos"
)

type actions struct {
	context     context.Context
	defaultAuth func(context.Context, string) (context.Context, error)
	dbConn      db.ConnPool
	userRepo    repos.IUserRepository
}

func NewActions(ctx context.Context, signingToken string) *actions {
	userRepo, err := do.Invoke[repos.IUserRepository](nil)
	if err != nil {
		logger.Global().Fatal("unable to get user repo")
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
	}
}

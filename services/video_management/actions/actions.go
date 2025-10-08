package actions

import (
	"context"

	"github.com/samber/do"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/interceptors"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/s3"
	"github.com/sweetloveinyourheart/sweet-reel/services/video_management/repos"
)

type actions struct {
	context     context.Context
	defaultAuth func(context.Context, string) (context.Context, error)
	s3Client    s3.S3Storage
	videoRepo   repos.VideoRepositoryInterface
}

func NewActions(ctx context.Context, signingToken string) *actions {
	s3Client, err := do.Invoke[s3.S3Storage](nil)
	if err != nil {
		logger.Global().Fatal("unable to get s3 client")
	}

	videoRepo, err := do.Invoke[repos.VideoRepositoryInterface](nil)
	if err != nil {
		logger.Global().Fatal("unable to get s3 client")
	}

	return &actions{
		context:     ctx,
		defaultAuth: interceptors.ConnectServerAuthHandler(signingToken),
		s3Client:    s3Client,
		videoRepo:   videoRepo,
	}
}

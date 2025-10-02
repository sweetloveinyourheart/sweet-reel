package videoprocessing

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/config"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
)

func InitializeRepos(ctx context.Context) error {
	appID := fmt.Sprintf("video-processing-%s", config.Instance().GetString("video.processing.id"))

	if err := InitializeCoreRepos(appID, ctx); err != nil {
		logger.Global().ErrorContext(ctx, "failed to initialize core repos", zap.Error(err))
		return err
	}

	return nil
}

func InitializeCoreRepos(appID string, ctx context.Context) error {
	return nil
}

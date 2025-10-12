package videomanagement

import (
	"context"
	"embed"

	"github.com/sweetloveinyourheart/sweet-reel/services/video_management/domains/processing"
)

//go:embed migrations/*.sql
var FS embed.FS

func InitializeRepos(ctx context.Context) error {
	_, err := processing.NewVideoProcessManager(ctx)
	if err != nil {
		return err
	}

	return nil
}

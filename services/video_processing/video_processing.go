package videoprocessing

import (
	"context"

	"github.com/sweetloveinyourheart/sweet-reel/services/video_processing/domains/processing"
)

func InitializeRepos(ctx context.Context) error {
	_, err := processing.NewVideoSplitterProcessManager(ctx)
	if err != nil {
		return err
	}

	return nil
}

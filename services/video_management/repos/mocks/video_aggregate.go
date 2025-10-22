package mocks

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/sweetloveinyourheart/sweet-reel/services/video_management/models"
)

// MockVideoAggregateRepository is a mock implementation of IVideoAggregateRepository
type MockVideoAggregateRepository struct {
	MockVideoRepository
}

func (m *MockVideoAggregateRepository) GetVideosWithThumbnailByUploaderID(ctx context.Context, uploaderID uuid.UUID, limit, offset int) ([]*models.VideoWithThumbnails, error) {
	args := m.Called(ctx, uploaderID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.VideoWithThumbnails), args.Error(1)
}

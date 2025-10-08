package mocks

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/sweetloveinyourheart/sweet-reel/services/video_management/models"
)

// MockVideoRepository is a mock implementation of VideoRepositoryInterface
type MockVideoRepository struct {
	mock.Mock
}

// Video operations

func (m *MockVideoRepository) CreateVideo(ctx context.Context, video *models.Video) error {
	args := m.Called(ctx, video)
	return args.Error(0)
}

func (m *MockVideoRepository) GetVideoByID(ctx context.Context, id uuid.UUID) (*models.Video, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Video), args.Error(1)
}

func (m *MockVideoRepository) GetVideosByUploaderID(ctx context.Context, uploaderID uuid.UUID, limit, offset int) ([]*models.Video, error) {
	args := m.Called(ctx, uploaderID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Video), args.Error(1)
}

func (m *MockVideoRepository) UpdateVideo(ctx context.Context, video *models.Video) error {
	args := m.Called(ctx, video)
	return args.Error(0)
}

func (m *MockVideoRepository) UpdateVideoStatus(ctx context.Context, id uuid.UUID, status models.VideoStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockVideoRepository) UpdateVideoProcessedAt(ctx context.Context, id uuid.UUID, processedAt time.Time) error {
	args := m.Called(ctx, id, processedAt)
	return args.Error(0)
}

func (m *MockVideoRepository) DeleteVideo(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockVideoRepository) ListVideos(ctx context.Context, limit, offset int) ([]*models.Video, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Video), args.Error(1)
}

// Video manifest operations

func (m *MockVideoRepository) CreateVideoManifest(ctx context.Context, manifest *models.VideoManifest) error {
	args := m.Called(ctx, manifest)
	return args.Error(0)
}

func (m *MockVideoRepository) GetVideoManifestByVideoID(ctx context.Context, videoID uuid.UUID) (*models.VideoManifest, error) {
	args := m.Called(ctx, videoID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.VideoManifest), args.Error(1)
}

func (m *MockVideoRepository) UpdateVideoManifest(ctx context.Context, manifest *models.VideoManifest) error {
	args := m.Called(ctx, manifest)
	return args.Error(0)
}

func (m *MockVideoRepository) DeleteVideoManifest(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Video variant operations

func (m *MockVideoRepository) CreateVideoVariant(ctx context.Context, variant *models.VideoVariant) error {
	args := m.Called(ctx, variant)
	return args.Error(0)
}

func (m *MockVideoRepository) GetVideoVariantsByVideoID(ctx context.Context, videoID uuid.UUID) ([]*models.VideoVariant, error) {
	args := m.Called(ctx, videoID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.VideoVariant), args.Error(1)
}

func (m *MockVideoRepository) GetVideoVariantByID(ctx context.Context, id uuid.UUID) (*models.VideoVariant, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.VideoVariant), args.Error(1)
}

func (m *MockVideoRepository) UpdateVideoVariant(ctx context.Context, variant *models.VideoVariant) error {
	args := m.Called(ctx, variant)
	return args.Error(0)
}

func (m *MockVideoRepository) DeleteVideoVariant(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockVideoRepository) DeleteVideoVariantsByVideoID(ctx context.Context, videoID uuid.UUID) error {
	args := m.Called(ctx, videoID)
	return args.Error(0)
}

// Video thumbnail operations

func (m *MockVideoRepository) CreateVideoThumbnail(ctx context.Context, thumbnail *models.VideoThumbnail) error {
	args := m.Called(ctx, thumbnail)
	return args.Error(0)
}

func (m *MockVideoRepository) GetVideoThumbnailsByVideoID(ctx context.Context, videoID uuid.UUID) ([]*models.VideoThumbnail, error) {
	args := m.Called(ctx, videoID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.VideoThumbnail), args.Error(1)
}

func (m *MockVideoRepository) GetVideoThumbnailByID(ctx context.Context, id uuid.UUID) (*models.VideoThumbnail, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.VideoThumbnail), args.Error(1)
}

func (m *MockVideoRepository) UpdateVideoThumbnail(ctx context.Context, thumbnail *models.VideoThumbnail) error {
	args := m.Called(ctx, thumbnail)
	return args.Error(0)
}

func (m *MockVideoRepository) DeleteVideoThumbnail(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockVideoRepository) DeleteVideoThumbnailsByVideoID(ctx context.Context, videoID uuid.UUID) error {
	args := m.Called(ctx, videoID)
	return args.Error(0)
}

// Aggregate operations

func (m *MockVideoRepository) GetVideoWithRelations(ctx context.Context, videoID uuid.UUID) (*models.Video, *models.VideoManifest, []*models.VideoVariant, []*models.VideoThumbnail, error) {
	args := m.Called(ctx, videoID)

	var video *models.Video
	var manifest *models.VideoManifest
	var variants []*models.VideoVariant
	var thumbnails []*models.VideoThumbnail

	if args.Get(0) != nil {
		video = args.Get(0).(*models.Video)
	}
	if args.Get(1) != nil {
		manifest = args.Get(1).(*models.VideoManifest)
	}
	if args.Get(2) != nil {
		variants = args.Get(2).([]*models.VideoVariant)
	}
	if args.Get(3) != nil {
		thumbnails = args.Get(3).([]*models.VideoThumbnail)
	}

	return video, manifest, variants, thumbnails, args.Error(4)
}

func (m *MockVideoRepository) GetVideoCount(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockVideoRepository) GetVideoCountByUploaderID(ctx context.Context, uploaderID uuid.UUID) (int64, error) {
	args := m.Called(ctx, uploaderID)
	return args.Get(0).(int64), args.Error(1)
}

// Ensure MockVideoRepository implements VideoRepositoryInterface
var _ interface {
	CreateVideo(ctx context.Context, video *models.Video) error
	GetVideoByID(ctx context.Context, id uuid.UUID) (*models.Video, error)
	GetVideosByUploaderID(ctx context.Context, uploaderID uuid.UUID, limit, offset int) ([]*models.Video, error)
	UpdateVideo(ctx context.Context, video *models.Video) error
	UpdateVideoStatus(ctx context.Context, id uuid.UUID, status models.VideoStatus) error
	UpdateVideoProcessedAt(ctx context.Context, id uuid.UUID, processedAt time.Time) error
	DeleteVideo(ctx context.Context, id uuid.UUID) error
	ListVideos(ctx context.Context, limit, offset int) ([]*models.Video, error)
	CreateVideoManifest(ctx context.Context, manifest *models.VideoManifest) error
	GetVideoManifestByVideoID(ctx context.Context, videoID uuid.UUID) (*models.VideoManifest, error)
	UpdateVideoManifest(ctx context.Context, manifest *models.VideoManifest) error
	DeleteVideoManifest(ctx context.Context, id uuid.UUID) error
	CreateVideoVariant(ctx context.Context, variant *models.VideoVariant) error
	GetVideoVariantsByVideoID(ctx context.Context, videoID uuid.UUID) ([]*models.VideoVariant, error)
	GetVideoVariantByID(ctx context.Context, id uuid.UUID) (*models.VideoVariant, error)
	UpdateVideoVariant(ctx context.Context, variant *models.VideoVariant) error
	DeleteVideoVariant(ctx context.Context, id uuid.UUID) error
	DeleteVideoVariantsByVideoID(ctx context.Context, videoID uuid.UUID) error
	CreateVideoThumbnail(ctx context.Context, thumbnail *models.VideoThumbnail) error
	GetVideoThumbnailsByVideoID(ctx context.Context, videoID uuid.UUID) ([]*models.VideoThumbnail, error)
	GetVideoThumbnailByID(ctx context.Context, id uuid.UUID) (*models.VideoThumbnail, error)
	UpdateVideoThumbnail(ctx context.Context, thumbnail *models.VideoThumbnail) error
	DeleteVideoThumbnail(ctx context.Context, id uuid.UUID) error
	DeleteVideoThumbnailsByVideoID(ctx context.Context, videoID uuid.UUID) error
	GetVideoWithRelations(ctx context.Context, videoID uuid.UUID) (*models.Video, *models.VideoManifest, []*models.VideoVariant, []*models.VideoThumbnail, error)
	GetVideoCount(ctx context.Context) (int64, error)
	GetVideoCountByUploaderID(ctx context.Context, uploaderID uuid.UUID) (int64, error)
} = (*MockVideoRepository)(nil)

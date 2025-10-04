package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/ffmpeg"
)

// MockFFmpeg is a mock of FFmpeg functionality
type MockFFmpeg struct {
	mock.Mock
}

func (m *MockFFmpeg) IsAvailable(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockFFmpeg) ProbeFile(ctx context.Context, filePath string) (*ffmpeg.ProbeInfo, error) {
	args := m.Called(ctx, filePath)
	return args.Get(0).(*ffmpeg.ProbeInfo), args.Error(1)
}

func (m *MockFFmpeg) SegmentVideoMultiQuality(ctx context.Context, inputPath, outputDir string, options []ffmpeg.SegmentationOptions, progressCallback ffmpeg.ProgressCallback) error {
	args := m.Called(ctx, inputPath, outputDir, options, progressCallback)
	return args.Error(0)
}

func (m *MockFFmpeg) CreateThumbnail(ctx context.Context, inputPath, outputPath, timestamp string, width, height int) error {
	args := m.Called(ctx, inputPath, outputPath, timestamp, width, height)
	return args.Error(0)
}

package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/ffmpeg"
)

// Ensure that MockFFmpeg implements ffmpeg.FFmpegInterface
var _ ffmpeg.FFmpegInterface = (*MockFFmpeg)(nil)

// MockFFmpeg is a mock of FFmpeg functionality that implements all public methods
type MockFFmpeg struct {
	mock.Mock
}

// Core functionality methods
func (m *MockFFmpeg) SetBinaryPath(path string) {
	m.Called(path)
}

func (m *MockFFmpeg) IsAvailable(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockFFmpeg) GetVersion(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

// Probe methods
func (m *MockFFmpeg) ProbeFile(ctx context.Context, inputPath string) (*ffmpeg.ProbeInfo, error) {
	args := m.Called(ctx, inputPath)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ffmpeg.ProbeInfo), args.Error(1)
}

func (m *MockFFmpeg) GetVideoInfo(ctx context.Context, inputPath string) (*ffmpeg.StreamInfo, error) {
	args := m.Called(ctx, inputPath)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ffmpeg.StreamInfo), args.Error(1)
}

func (m *MockFFmpeg) GetAudioInfo(ctx context.Context, inputPath string) (*ffmpeg.StreamInfo, error) {
	args := m.Called(ctx, inputPath)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ffmpeg.StreamInfo), args.Error(1)
}

func (m *MockFFmpeg) GetDuration(ctx context.Context, inputPath string) (string, error) {
	args := m.Called(ctx, inputPath)
	return args.String(0), args.Error(1)
}

func (m *MockFFmpeg) GetFrameRate(ctx context.Context, inputPath string) (string, error) {
	args := m.Called(ctx, inputPath)
	return args.String(0), args.Error(1)
}

func (m *MockFFmpeg) GetResolution(ctx context.Context, inputPath string) (int, int, error) {
	args := m.Called(ctx, inputPath)
	return args.Int(0), args.Int(1), args.Error(2)
}

func (m *MockFFmpeg) GetBitrate(ctx context.Context, inputPath string) (string, error) {
	args := m.Called(ctx, inputPath)
	return args.String(0), args.Error(1)
}

func (m *MockFFmpeg) IsVideoFile(ctx context.Context, inputPath string) (bool, error) {
	args := m.Called(ctx, inputPath)
	return args.Bool(0), args.Error(1)
}

func (m *MockFFmpeg) IsAudioFile(ctx context.Context, inputPath string) (bool, error) {
	args := m.Called(ctx, inputPath)
	return args.Bool(0), args.Error(1)
}

func (m *MockFFmpeg) GetCodecInfo(ctx context.Context, inputPath string) (map[string][]string, error) {
	args := m.Called(ctx, inputPath)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string][]string), args.Error(1)
}

// Transcoding methods
func (m *MockFFmpeg) Transcode(ctx context.Context, inputPath, outputPath string, options ffmpeg.TranscodeOptions, progressCallback ffmpeg.ProgressCallback) error {
	args := m.Called(ctx, inputPath, outputPath, options, progressCallback)
	return args.Error(0)
}

func (m *MockFFmpeg) TranscodeToMultipleQualities(ctx context.Context, inputPath, outputDir string, qualities []ffmpeg.TranscodeOptions, progressCallback ffmpeg.ProgressCallback) error {
	args := m.Called(ctx, inputPath, outputDir, qualities, progressCallback)
	return args.Error(0)
}

func (m *MockFFmpeg) ExtractAudio(ctx context.Context, inputPath, outputPath string, options ffmpeg.TranscodeOptions, progressCallback ffmpeg.ProgressCallback) error {
	args := m.Called(ctx, inputPath, outputPath, options, progressCallback)
	return args.Error(0)
}

func (m *MockFFmpeg) CreateThumbnail(ctx context.Context, inputPath, outputPath, timeOffset string, width, height int) error {
	args := m.Called(ctx, inputPath, outputPath, timeOffset, width, height)
	return args.Error(0)
}

// Segmentation methods
func (m *MockFFmpeg) SegmentVideo(ctx context.Context, inputPath, outputDir string, options ffmpeg.SegmentationOptions, progressCallback ffmpeg.ProgressCallback) error {
	args := m.Called(ctx, inputPath, outputDir, options, progressCallback)
	return args.Error(0)
}

func (m *MockFFmpeg) SegmentVideoMultiQuality(ctx context.Context, inputPath, outputDir string, qualities []ffmpeg.SegmentationOptions, progressCallback ffmpeg.ProgressCallback) error {
	args := m.Called(ctx, inputPath, outputDir, qualities, progressCallback)
	return args.Error(0)
}

func (m *MockFFmpeg) CreateDASHSegments(ctx context.Context, inputPath, outputDir, segmentDuration string, progressCallback ffmpeg.ProgressCallback) error {
	args := m.Called(ctx, inputPath, outputDir, segmentDuration, progressCallback)
	return args.Error(0)
}

func (m *MockFFmpeg) CreatePreviewClips(ctx context.Context, inputPath, outputDir, clipDuration string, intervalSeconds int, progressCallback ffmpeg.ProgressCallback) error {
	args := m.Called(ctx, inputPath, outputDir, clipDuration, intervalSeconds, progressCallback)
	return args.Error(0)
}

// Utility methods
func (m *MockFFmpeg) ConcatenateVideos(ctx context.Context, inputPaths []string, outputPath string, progressCallback ffmpeg.ProgressCallback) error {
	args := m.Called(ctx, inputPaths, outputPath, progressCallback)
	return args.Error(0)
}

func (m *MockFFmpeg) TrimVideo(ctx context.Context, inputPath, outputPath, startTime, duration string, progressCallback ffmpeg.ProgressCallback) error {
	args := m.Called(ctx, inputPath, outputPath, startTime, duration, progressCallback)
	return args.Error(0)
}

func (m *MockFFmpeg) AddWatermark(ctx context.Context, inputPath, watermarkPath, outputPath, position string, progressCallback ffmpeg.ProgressCallback) error {
	args := m.Called(ctx, inputPath, watermarkPath, outputPath, position, progressCallback)
	return args.Error(0)
}

func (m *MockFFmpeg) ConvertFormat(ctx context.Context, inputPath, outputPath, targetFormat string, progressCallback ffmpeg.ProgressCallback) error {
	args := m.Called(ctx, inputPath, outputPath, targetFormat, progressCallback)
	return args.Error(0)
}

func (m *MockFFmpeg) ExtractFrames(ctx context.Context, inputPath, outputDir string, fps float64, imageFormat string) error {
	args := m.Called(ctx, inputPath, outputDir, fps, imageFormat)
	return args.Error(0)
}

func (m *MockFFmpeg) CreateVideoFromImages(ctx context.Context, imagePattern, outputPath string, fps float64, progressCallback ffmpeg.ProgressCallback) error {
	args := m.Called(ctx, imagePattern, outputPath, fps, progressCallback)
	return args.Error(0)
}

func (m *MockFFmpeg) GetTotalFrames(ctx context.Context, inputPath string) (int64, error) {
	args := m.Called(ctx, inputPath)
	return args.Get(0).(int64), args.Error(1)
}

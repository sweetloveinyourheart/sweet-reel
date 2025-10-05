package ffmpeg

import "context"

// FFmpegInterface defines the contract for FFmpeg operations
// This interface allows for easy mocking and testing
type FFmpegInterface interface {
	// Core functionality
	SetBinaryPath(path string)
	IsAvailable(ctx context.Context) error
	GetVersion(ctx context.Context) (string, error)

	// Probe methods
	ProbeFile(ctx context.Context, inputPath string) (*ProbeInfo, error)
	GetVideoInfo(ctx context.Context, inputPath string) (*StreamInfo, error)
	GetAudioInfo(ctx context.Context, inputPath string) (*StreamInfo, error)
	GetDuration(ctx context.Context, inputPath string) (string, error)
	GetFrameRate(ctx context.Context, inputPath string) (string, error)
	GetResolution(ctx context.Context, inputPath string) (int, int, error)
	GetBitrate(ctx context.Context, inputPath string) (string, error)
	IsVideoFile(ctx context.Context, inputPath string) (bool, error)
	IsAudioFile(ctx context.Context, inputPath string) (bool, error)
	GetCodecInfo(ctx context.Context, inputPath string) (map[string][]string, error)

	// Transcoding methods
	Transcode(ctx context.Context, inputPath, outputPath string, options TranscodeOptions, progressCallback ProgressCallback) error
	TranscodeToMultipleQualities(ctx context.Context, inputPath, outputDir string, qualities []TranscodeOptions, progressCallback ProgressCallback) error
	ExtractAudio(ctx context.Context, inputPath, outputPath string, options TranscodeOptions, progressCallback ProgressCallback) error
	CreateThumbnail(ctx context.Context, inputPath, outputPath, timeOffset string, width, height int) error

	// Segmentation methods
	SegmentVideo(ctx context.Context, inputPath, outputDir string, options SegmentationOptions, progressCallback ProgressCallback) error
	SegmentVideoMultiQuality(ctx context.Context, inputPath, outputDir string, qualities []SegmentationOptions, progressCallback ProgressCallback) error
	CreateDASHSegments(ctx context.Context, inputPath, outputDir, segmentDuration string, progressCallback ProgressCallback) error
	CreatePreviewClips(ctx context.Context, inputPath, outputDir, clipDuration string, intervalSeconds int, progressCallback ProgressCallback) error

	// Utility methods
	ConcatenateVideos(ctx context.Context, inputPaths []string, outputPath string, progressCallback ProgressCallback) error
	TrimVideo(ctx context.Context, inputPath, outputPath, startTime, duration string, progressCallback ProgressCallback) error
	AddWatermark(ctx context.Context, inputPath, watermarkPath, outputPath, position string, progressCallback ProgressCallback) error
	ConvertFormat(ctx context.Context, inputPath, outputPath, targetFormat string, progressCallback ProgressCallback) error
	ExtractFrames(ctx context.Context, inputPath, outputDir string, fps float64, imageFormat string) error
	CreateVideoFromImages(ctx context.Context, imagePattern, outputPath string, fps float64, progressCallback ProgressCallback) error
	GetTotalFrames(ctx context.Context, inputPath string) (int64, error)
}

// Ensure that FFmpeg implements FFmpegInterface
var _ FFmpegInterface = (*FFmpeg)(nil)
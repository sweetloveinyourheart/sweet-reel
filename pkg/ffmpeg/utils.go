package ffmpeg

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
	"go.uber.org/zap"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
)

// ConcatenateVideos concatenates multiple video files into one
func (f *FFmpeg) ConcatenateVideos(ctx context.Context, inputPaths []string, outputPath string, progressCallback ProgressCallback) error {
	if len(inputPaths) == 0 {
		return errors.New("no input files provided")
	}

	for _, inputPath := range inputPaths {
		if err := validateInputFile(inputPath); err != nil {
			return errors.Wrapf(err, "invalid input file: %s", inputPath)
		}
	}

	if err := ensureOutputDir(outputPath); err != nil {
		return errors.Wrap(err, "failed to create output directory")
	}

	// Create a temporary file list for concatenation
	listFile, err := createTempFile(".txt")
	if err != nil {
		return errors.Wrap(err, "failed to create temporary file list")
	}
	defer os.Remove(listFile)

	// Write file list
	var fileList []string
	for _, inputPath := range inputPaths {
		// Use absolute paths and escape special characters
		absPath, err := filepath.Abs(inputPath)
		if err != nil {
			return errors.Wrapf(err, "failed to get absolute path for %s", inputPath)
		}
		fileList = append(fileList, "file '"+strings.ReplaceAll(absPath, "'", "'\\''")+"'")
	}

	listContent := strings.Join(fileList, "\n")
	if err := os.WriteFile(listFile, []byte(listContent), 0644); err != nil {
		return errors.Wrap(err, "failed to write file list")
	}

	args := []string{
		"-y", // Overwrite output files
		"-f", "concat",
		"-safe", "0",
		"-i", listFile,
		"-c", "copy", // Copy streams without re-encoding
		outputPath,
	}

	logger.Global().Info("Concatenating videos",
		zap.Strings("inputs", inputPaths),
		zap.String("output", outputPath))

	return f.runCommand(ctx, args, progressCallback)
}

// TrimVideo trims a video to a specific time range
func (f *FFmpeg) TrimVideo(ctx context.Context, inputPath, outputPath string, startTime, duration string, progressCallback ProgressCallback) error {
	if err := validateInputFile(inputPath); err != nil {
		return errors.Wrap(err, "invalid input file")
	}

	if err := ensureOutputDir(outputPath); err != nil {
		return errors.Wrap(err, "failed to create output directory")
	}

	args := []string{
		"-y",            // Overwrite output files
		"-i", inputPath, // Input file
		"-ss", startTime, // Start time
		"-t", duration, // Duration
		"-c", "copy", // Copy streams without re-encoding
		outputPath,
	}

	logger.Global().Info("Trimming video",
		zap.String("input", inputPath),
		zap.String("output", outputPath),
		zap.String("start", startTime),
		zap.String("duration", duration))

	return f.runCommand(ctx, args, progressCallback)
}

// AddWatermark adds a watermark image or text to a video
func (f *FFmpeg) AddWatermark(ctx context.Context, inputPath, watermarkPath, outputPath string, position string, progressCallback ProgressCallback) error {
	if err := validateInputFile(inputPath); err != nil {
		return errors.Wrap(err, "invalid input file")
	}

	if watermarkPath != "" {
		if err := validateInputFile(watermarkPath); err != nil {
			return errors.Wrap(err, "invalid watermark file")
		}
	}

	if err := ensureOutputDir(outputPath); err != nil {
		return errors.Wrap(err, "failed to create output directory")
	}

	var overlay string
	switch position {
	case "top-left":
		overlay = "overlay=10:10"
	case "top-right":
		overlay = "overlay=main_w-overlay_w-10:10"
	case "bottom-left":
		overlay = "overlay=10:main_h-overlay_h-10"
	case "bottom-right":
		overlay = "overlay=main_w-overlay_w-10:main_h-overlay_h-10"
	case "center":
		overlay = "overlay=(main_w-overlay_w)/2:(main_h-overlay_h)/2"
	default:
		overlay = "overlay=10:10" // Default to top-left
	}

	args := []string{
		"-y",            // Overwrite output files
		"-i", inputPath, // Input video
		"-i", watermarkPath, // Watermark image
		"-filter_complex", overlay,
		"-c:a", "copy", // Copy audio without re-encoding
		outputPath,
	}

	logger.Global().Info("Adding watermark",
		zap.String("input", inputPath),
		zap.String("watermark", watermarkPath),
		zap.String("output", outputPath),
		zap.String("position", position))

	return f.runCommand(ctx, args, progressCallback)
}

// ConvertFormat converts a video from one format to another
func (f *FFmpeg) ConvertFormat(ctx context.Context, inputPath, outputPath string, targetFormat string, progressCallback ProgressCallback) error {
	if err := validateInputFile(inputPath); err != nil {
		return errors.Wrap(err, "invalid input file")
	}

	if err := ensureOutputDir(outputPath); err != nil {
		return errors.Wrap(err, "failed to create output directory")
	}

	args := []string{
		"-y",            // Overwrite output files
		"-i", inputPath, // Input file
		"-f", targetFormat, // Target format
		"-c", "copy", // Copy streams without re-encoding when possible
		outputPath,
	}

	logger.Global().Info("Converting format",
		zap.String("input", inputPath),
		zap.String("output", outputPath),
		zap.String("targetFormat", targetFormat))

	return f.runCommand(ctx, args, progressCallback)
}

// ExtractFrames extracts frames from a video at specified intervals
func (f *FFmpeg) ExtractFrames(ctx context.Context, inputPath, outputDir string, fps float64, imageFormat string) error {
	if err := validateInputFile(inputPath); err != nil {
		return errors.Wrap(err, "invalid input file")
	}

	if err := ensureOutputDir(outputDir); err != nil {
		return errors.Wrap(err, "failed to create output directory")
	}

	if imageFormat == "" {
		imageFormat = "jpg"
	}

	outputPattern := filepath.Join(outputDir, "frame_%04d."+imageFormat)

	args := []string{
		"-y",            // Overwrite output files
		"-i", inputPath, // Input file
		"-vf", "fps=" + formatFloat(fps), // Video filter for frame extraction
		"-q:v", "2", // High quality
		outputPattern,
	}

	logger.Global().Info("Extracting frames",
		zap.String("input", inputPath),
		zap.String("outputDir", outputDir),
		zap.Float64("fps", fps),
		zap.String("format", imageFormat))

	return f.runCommand(ctx, args, nil)
}

// CreateVideoFromImages creates a video from a sequence of images
func (f *FFmpeg) CreateVideoFromImages(ctx context.Context, imagePattern, outputPath string, fps float64, progressCallback ProgressCallback) error {
	if err := ensureOutputDir(outputPath); err != nil {
		return errors.Wrap(err, "failed to create output directory")
	}

	args := []string{
		"-y",                           // Overwrite output files
		"-framerate", formatFloat(fps), // Input framerate
		"-i", imagePattern, // Input pattern (e.g., "frame_%04d.jpg")
		"-c:v", CodecLibX264,
		"-pix_fmt", PixFmtYUV420P,
		outputPath,
	}

	logger.Global().Info("Creating video from images",
		zap.String("pattern", imagePattern),
		zap.String("output", outputPath),
		zap.Float64("fps", fps))

	return f.runCommand(ctx, args, progressCallback)
}

// GetTotalFrames calculates the total number of frames in a video
func (f *FFmpeg) GetTotalFrames(ctx context.Context, inputPath string) (int64, error) {
	if err := validateInputFile(inputPath); err != nil {
		return 0, errors.Wrap(err, "invalid input file")
	}

	// Use ffprobe to count frames
	probePath := strings.Replace(f.binaryPath, "ffmpeg", "ffprobe", 1)

	args := []string{
		"-v", "error",
		"-select_streams", "v:0",
		"-count_frames",
		"-show_entries", "stream=nb_read_frames",
		"-csv=p=0",
		inputPath,
	}

	cmd := exec.CommandContext(ctx, probePath, args...)
	output, err := cmd.Output()
	if err != nil {
		return 0, errors.Wrap(err, "failed to count frames")
	}

	frameCount := strings.TrimSpace(string(output))
	if frameCount == "" || frameCount == "N/A" {
		return 0, errors.New("unable to determine frame count")
	}

	var count int64
	if _, err := fmt.Sscanf(frameCount, "%d", &count); err != nil {
		return 0, errors.Wrap(err, "failed to parse frame count")
	}

	return count, nil
}

// GetMimeType returns the MIME type based on file extension
func GetMimeType(filePath string) string {
	ext := filepath.Ext(filePath)
	switch ext {
	case ExtM3U8:
		return MimeTypeM3U8
	case ExtTS:
		return MimeTypeTS
	case ExtMP4:
		return MimeTypeMP4
	case ExtJPG, ExtJPEG:
		return MimeTypeJPEG
	case ExtPNG:
		return MimeTypePNG
	case ExtMKV:
		return MimeTypeMKV
	case ExtWebM:
		return MimeTypeWebM
	case ExtAVI:
		return MimeTypeAVI
	case ExtMOV:
		return MimeTypeMOV
	case ExtFLV:
		return MimeTypeFLV
	default:
		return MimeTypeOctetStream
	}
}

// formatFloat formats a float64 to string with appropriate precision
func formatFloat(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

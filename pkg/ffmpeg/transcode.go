package ffmpeg

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
	"go.uber.org/zap"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
)

// Transcode transcodes a video file with the given options
func (f *FFmpeg) Transcode(ctx context.Context, inputPath, outputPath string, options TranscodeOptions, progressCallback ProgressCallback) error {
	if err := validateInputFile(inputPath); err != nil {
		return errors.Wrap(err, "invalid input file")
	}

	if err := ensureOutputDir(outputPath); err != nil {
		return errors.Wrap(err, "failed to create output directory")
	}

	args := []string{
		"-y",            // Overwrite output files
		"-i", inputPath, // Input file
	}

	// Hardware acceleration (must come before input)
	if options.HWAccel != "" {
		args = append([]string{"-hwaccel", options.HWAccel}, args...)
	}

	// Start time (seeking)
	if options.StartTime != "" {
		args = append(args, "-ss", options.StartTime)
	}

	// Duration (trimming)
	if options.Duration != "" {
		args = append(args, "-t", options.Duration)
	}

	// Video codec
	if options.VideoCodec != "" {
		args = append(args, "-c:v", options.VideoCodec)
	}

	// Video quality/bitrate
	if options.VideoQuality != "" && options.VideoBitrate == "" {
		// Use CRF (Constant Rate Factor) for quality-based encoding
		args = append(args, "-crf", options.VideoQuality)
	} else if options.VideoBitrate != "" {
		// Use bitrate-based encoding
		args = append(args, "-b:v", options.VideoBitrate)
	}

	// Resolution
	if options.Resolution != "" {
		args = append(args, "-s", options.Resolution)
	}

	// Frame rate
	if options.FrameRate != "" {
		args = append(args, "-r", options.FrameRate)
	}

	// Audio codec
	if options.AudioCodec != "" {
		args = append(args, "-c:a", options.AudioCodec)
	}

	// Audio bitrate
	if options.AudioBitrate != "" {
		args = append(args, "-b:a", options.AudioBitrate)
	}

	// Audio channels
	if options.AudioChannels != "" {
		args = append(args, "-ac", options.AudioChannels)
	}

	// Audio sample rate
	if options.AudioSampleRate != "" {
		args = append(args, "-ar", options.AudioSampleRate)
	}

	// Format
	if options.Format != "" {
		args = append(args, "-f", options.Format)
	}

	// Custom arguments
	args = append(args, options.CustomArgs...)

	// Output file
	args = append(args, outputPath)

	logger.Global().Info("Starting video transcoding",
		zap.String("input", inputPath),
		zap.String("output", outputPath),
		zap.String("args", strings.Join(args, " ")))

	return f.runCommand(ctx, args, progressCallback)
}

// TranscodeToMultipleQualities transcodes a video to multiple quality levels
func (f *FFmpeg) TranscodeToMultipleQualities(ctx context.Context, inputPath string, outputDir string, qualities []TranscodeOptions, progressCallback ProgressCallback) error {
	if err := validateInputFile(inputPath); err != nil {
		return errors.Wrap(err, "invalid input file")
	}

	if err := ensureOutputDir(outputDir); err != nil {
		return errors.Wrap(err, "failed to create output directory")
	}

	for i, quality := range qualities {
		// Generate output filename based on quality settings
		outputFilename := generateQualityFilename(inputPath, quality)
		outputPath := filepath.Join(outputDir, outputFilename)

		logger.Global().Info("Transcoding quality level",
			zap.Int("quality", i+1),
			zap.Int("total", len(qualities)),
			zap.String("output", outputFilename))

		if err := f.Transcode(ctx, inputPath, outputPath, quality, progressCallback); err != nil {
			return errors.Wrapf(err, "failed to transcode quality level %d", i+1)
		}
	}

	return nil
}

// ExtractAudio extracts audio from a video file
func (f *FFmpeg) ExtractAudio(ctx context.Context, inputPath, outputPath string, options TranscodeOptions, progressCallback ProgressCallback) error {
	if err := validateInputFile(inputPath); err != nil {
		return errors.Wrap(err, "invalid input file")
	}

	if err := ensureOutputDir(outputPath); err != nil {
		return errors.Wrap(err, "failed to create output directory")
	}

	args := []string{
		"-y",            // Overwrite output files
		"-i", inputPath, // Input file
		"-vn", // Disable video
	}

	// Audio codec
	if options.AudioCodec != "" {
		args = append(args, "-c:a", options.AudioCodec)
	}

	// Audio bitrate
	if options.AudioBitrate != "" {
		args = append(args, "-b:a", options.AudioBitrate)
	}

	// Audio channels
	if options.AudioChannels != "" {
		args = append(args, "-ac", options.AudioChannels)
	}

	// Audio sample rate
	if options.AudioSampleRate != "" {
		args = append(args, "-ar", options.AudioSampleRate)
	}

	// Custom arguments
	args = append(args, options.CustomArgs...)

	// Output file
	args = append(args, outputPath)

	logger.Global().Info("Extracting audio",
		zap.String("input", inputPath),
		zap.String("output", outputPath))

	return f.runCommand(ctx, args, progressCallback)
}

// CreateThumbnail creates a thumbnail image from a video at the specified time
func (f *FFmpeg) CreateThumbnail(ctx context.Context, inputPath, outputPath string, timeOffset string, width, height int) error {
	if err := validateInputFile(inputPath); err != nil {
		return errors.Wrap(err, "invalid input file")
	}

	if err := ensureOutputDir(outputPath); err != nil {
		return errors.Wrap(err, "failed to create output directory")
	}

	args := []string{
		"-y",            // Overwrite output files
		"-i", inputPath, // Input file
		"-ss", timeOffset, // Seek to time
		"-vframes", "1", // Extract only one frame
		"-an", // Disable audio
	}

	// Set dimensions if specified
	if width > 0 && height > 0 {
		args = append(args, "-s", fmt.Sprintf("%dx%d", width, height))
	}

	// Output file
	args = append(args, outputPath)

	logger.Global().Info("Creating thumbnail",
		zap.String("input", inputPath),
		zap.String("output", outputPath),
		zap.String("timeOffset", timeOffset))

	return f.runCommand(ctx, args, nil)
}

// generateQualityFilename generates a filename based on quality settings
func generateQualityFilename(inputPath string, quality TranscodeOptions) string {
	ext := filepath.Ext(inputPath)
	baseName := strings.TrimSuffix(filepath.Base(inputPath), ext)

	// Build quality suffix
	var qualitySuffix []string

	if quality.Resolution != "" {
		qualitySuffix = append(qualitySuffix, quality.Resolution)
	}

	if quality.VideoBitrate != "" {
		qualitySuffix = append(qualitySuffix, quality.VideoBitrate)
	} else if quality.VideoQuality != "" {
		qualitySuffix = append(qualitySuffix, "crf"+quality.VideoQuality)
	}

	if quality.AudioBitrate != "" {
		qualitySuffix = append(qualitySuffix, quality.AudioBitrate+"_audio")
	}

	// Determine output extension
	outputExt := ext
	if quality.Format != "" {
		outputExt = "." + quality.Format
	}

	if len(qualitySuffix) > 0 {
		return fmt.Sprintf("%s_%s%s", baseName, strings.Join(qualitySuffix, "_"), outputExt)
	}

	return fmt.Sprintf("%s_transcoded%s", baseName, outputExt)
}

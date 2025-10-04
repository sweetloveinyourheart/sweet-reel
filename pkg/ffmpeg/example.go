package ffmpeg

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

// Example demonstrates how to use the FFmpeg wrapper
func Example() {
	logger, _ := zap.NewDevelopment()

	// Create FFmpeg instance
	ffmpeg := New()

	ctx := context.Background()

	// Check if FFmpeg is available
	if err := ffmpeg.IsAvailable(ctx); err != nil {
		logger.Error("FFmpeg not available", zap.Error(err))
		return
	}

	// Get version
	version, err := ffmpeg.GetVersion(ctx)
	if err != nil {
		logger.Error("Failed to get FFmpeg version", zap.Error(err))
		return
	}
	logger.Info("FFmpeg version", zap.String("version", version))

	// Example input and output paths
	inputVideo := "/path/to/input.mp4"
	outputDir := "/path/to/output"

	// Example 1: Basic transcoding
	transcodeExample(ffmpeg, ctx, inputVideo, outputDir, logger)

	// Example 2: HLS segmentation
	hlsExample(ffmpeg, ctx, inputVideo, outputDir, logger)

	// Example 3: Multi-quality transcoding
	multiQualityExample(ffmpeg, ctx, inputVideo, outputDir, logger)

	// Example 4: Video utilities
	utilsExample(ffmpeg, ctx, inputVideo, outputDir, logger)
}

func transcodeExample(ffmpeg *FFmpeg, ctx context.Context, inputVideo, outputDir string, logger *zap.Logger) {
	logger.Info("=== Transcoding Example ===")

	// Basic transcoding with progress monitoring
	outputPath := filepath.Join(outputDir, "transcoded.mp4")

	options := TranscodeOptions{
		VideoCodec:   "libx264",
		VideoQuality: "23", // CRF 23 for good quality
		AudioCodec:   "aac",
		AudioBitrate: "128k",
		Resolution:   "1280x720",
		Format:       "mp4",
	}

	progressCallback := func(progress ProgressInfo) {
		logger.Info("Transcoding progress",
			zap.Float64("percentage", progress.Percentage),
			zap.Duration("current", progress.Current),
			zap.Duration("total", progress.Duration),
			zap.String("speed", progress.Speed))
	}

	if err := ffmpeg.Transcode(ctx, inputVideo, outputPath, options, progressCallback); err != nil {
		logger.Error("Transcoding failed", zap.Error(err))
		return
	}

	logger.Info("Transcoding completed", zap.String("output", outputPath))
}

func hlsExample(ffmpeg *FFmpeg, ctx context.Context, inputVideo, outputDir string, logger *zap.Logger) {
	logger.Info("=== HLS Segmentation Example ===")

	hlsDir := filepath.Join(outputDir, "hls")

	options := SegmentationOptions{
		SegmentDuration: "10", // 10-second segments
		PlaylistType:    "vod",
		PlaylistName:    "playlist.m3u8",
		SegmentPrefix:   "segment",
		SegmentFormat:   "ts",
		VideoCodec:      "libx264",
		VideoQuality:    "23",
		AudioCodec:      "aac",
		AudioBitrate:    "128k",
		Resolution:      "1280x720",
	}

	progressCallback := func(progress ProgressInfo) {
		logger.Info("Segmentation progress",
			zap.Float64("percentage", progress.Percentage),
			zap.Duration("current", progress.Current))
	}

	if err := ffmpeg.SegmentVideo(ctx, inputVideo, hlsDir, options, progressCallback); err != nil {
		logger.Error("HLS segmentation failed", zap.Error(err))
		return
	}

	logger.Info("HLS segmentation completed", zap.String("output", hlsDir))
}

func multiQualityExample(ffmpeg *FFmpeg, ctx context.Context, inputVideo, outputDir string, logger *zap.Logger) {
	logger.Info("=== Multi-Quality Example ===")

	// Define multiple quality levels
	qualities := []SegmentationOptions{
		{
			SegmentDuration: "10",
			PlaylistType:    "vod",
			VideoCodec:      "libx264",
			VideoQuality:    "28",
			AudioCodec:      "aac",
			AudioBitrate:    "96k",
			Resolution:      "854x480", // 480p
			VideoBitrate:    "1000k",
		},
		{
			SegmentDuration: "10",
			PlaylistType:    "vod",
			VideoCodec:      "libx264",
			VideoQuality:    "25",
			AudioCodec:      "aac",
			AudioBitrate:    "128k",
			Resolution:      "1280x720", // 720p
			VideoBitrate:    "2500k",
		},
		{
			SegmentDuration: "10",
			PlaylistType:    "vod",
			VideoCodec:      "libx264",
			VideoQuality:    "23",
			AudioCodec:      "aac",
			AudioBitrate:    "192k",
			Resolution:      "1920x1080", // 1080p
			VideoBitrate:    "5000k",
		},
	}

	multiQualityDir := filepath.Join(outputDir, "multi_quality")

	progressCallback := func(progress ProgressInfo) {
		logger.Info("Multi-quality progress",
			zap.Float64("percentage", progress.Percentage))
	}

	if err := ffmpeg.SegmentVideoMultiQuality(ctx, inputVideo, multiQualityDir, qualities, progressCallback); err != nil {
		logger.Error("Multi-quality segmentation failed", zap.Error(err))
		return
	}

	logger.Info("Multi-quality segmentation completed", zap.String("output", multiQualityDir))
}

func utilsExample(ffmpeg *FFmpeg, ctx context.Context, inputVideo, outputDir string, logger *zap.Logger) {
	logger.Info("=== Utilities Example ===")

	// Probe file information
	probeInfo, err := ffmpeg.ProbeFile(ctx, inputVideo)
	if err != nil {
		logger.Error("Failed to probe file", zap.Error(err))
		return
	}

	logger.Info("File information",
		zap.String("format", probeInfo.Format.FormatName),
		zap.String("duration", probeInfo.Format.Duration),
		zap.String("size", probeInfo.Format.Size),
		zap.Int("streams", len(probeInfo.Streams)))

	// Get video information
	videoInfo, err := ffmpeg.GetVideoInfo(ctx, inputVideo)
	if err != nil {
		logger.Error("Failed to get video info", zap.Error(err))
		return
	}

	logger.Info("Video information",
		zap.String("codec", videoInfo.CodecName),
		zap.Int("width", videoInfo.Width),
		zap.Int("height", videoInfo.Height),
		zap.String("frameRate", videoInfo.RFrameRate))

	// Create thumbnail
	thumbnailPath := filepath.Join(outputDir, "thumbnail.jpg")
	if err := ffmpeg.CreateThumbnail(ctx, inputVideo, thumbnailPath, "00:00:10", 320, 240); err != nil {
		logger.Error("Failed to create thumbnail", zap.Error(err))
	} else {
		logger.Info("Thumbnail created", zap.String("path", thumbnailPath))
	}

	// Extract audio
	audioPath := filepath.Join(outputDir, "audio.mp3")
	audioOptions := TranscodeOptions{
		AudioCodec:   "mp3",
		AudioBitrate: "192k",
	}

	if err := ffmpeg.ExtractAudio(ctx, inputVideo, audioPath, audioOptions, nil); err != nil {
		logger.Error("Failed to extract audio", zap.Error(err))
	} else {
		logger.Info("Audio extracted", zap.String("path", audioPath))
	}

	// Trim video
	trimmedPath := filepath.Join(outputDir, "trimmed.mp4")
	if err := ffmpeg.TrimVideo(ctx, inputVideo, trimmedPath, "00:00:30", "00:01:00", nil); err != nil {
		logger.Error("Failed to trim video", zap.Error(err))
	} else {
		logger.Info("Video trimmed", zap.String("path", trimmedPath))
	}
}

// VideoProcessingPipeline demonstrates a complete video processing pipeline
func VideoProcessingPipeline(inputPath, outputDir string, logger *zap.Logger) error {
	ffmpeg := New()
	ctx := context.Background()

	// Step 1: Validate input
	if err := ffmpeg.IsAvailable(ctx); err != nil {
		return fmt.Errorf("FFmpeg not available: %w", err)
	}

	// Step 2: Probe input file
	probeInfo, err := ffmpeg.ProbeFile(ctx, inputPath)
	if err != nil {
		return fmt.Errorf("failed to probe input file: %w", err)
	}

	logger.Info("Processing video",
		zap.String("file", inputPath),
		zap.String("format", probeInfo.Format.FormatName),
		zap.String("duration", probeInfo.Format.Duration))

	// Step 3: Create thumbnail
	thumbnailPath := filepath.Join(outputDir, "thumbnail.jpg")
	if err := ffmpeg.CreateThumbnail(ctx, inputPath, thumbnailPath, "00:00:05", 320, 240); err != nil {
		logger.Warn("Failed to create thumbnail", zap.Error(err))
	}

	// Step 4: Extract audio
	audioPath := filepath.Join(outputDir, "audio.aac")
	audioOptions := DefaultTranscodeOptions()
	audioOptions.AudioCodec = "aac"
	audioOptions.AudioBitrate = "128k"

	if err := ffmpeg.ExtractAudio(ctx, inputPath, audioPath, audioOptions, nil); err != nil {
		logger.Warn("Failed to extract audio", zap.Error(err))
	}

	// Step 5: Create HLS segments for adaptive streaming
	hlsDir := filepath.Join(outputDir, "hls")
	qualities := []SegmentationOptions{
		{
			SegmentDuration: "6",
			PlaylistType:    "vod",
			VideoCodec:      "libx264",
			VideoBitrate:    "800k",
			AudioCodec:      "aac",
			AudioBitrate:    "96k",
			Resolution:      "640x360", // 360p
		},
		{
			SegmentDuration: "6",
			PlaylistType:    "vod",
			VideoCodec:      "libx264",
			VideoBitrate:    "1400k",
			AudioCodec:      "aac",
			AudioBitrate:    "128k",
			Resolution:      "1280x720", // 720p
		},
		{
			SegmentDuration: "6",
			PlaylistType:    "vod",
			VideoCodec:      "libx264",
			VideoBitrate:    "2800k",
			AudioCodec:      "aac",
			AudioBitrate:    "192k",
			Resolution:      "1920x1080", // 1080p
		},
	}

	progressCallback := func(progress ProgressInfo) {
		if int(progress.Percentage)%10 == 0 { // Log every 10%
			logger.Info("Processing progress",
				zap.Float64("percentage", progress.Percentage),
				zap.String("speed", progress.Speed))
		}
	}

	start := time.Now()
	if err := ffmpeg.SegmentVideoMultiQuality(ctx, inputPath, hlsDir, qualities, progressCallback); err != nil {
		return fmt.Errorf("failed to create HLS segments: %w", err)
	}

	logger.Info("Video processing completed",
		zap.String("outputDir", outputDir),
		zap.Duration("processingTime", time.Since(start)))

	return nil
}

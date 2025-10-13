package ffmpeg

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
	"go.uber.org/zap"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
)

// SegmentVideo segments a video into HLS format with multiple quality levels
func (f *FFmpeg) SegmentVideo(ctx context.Context, inputPath, outputDir string, options SegmentationOptions, progressCallback ProgressCallback) error {
	if err := validateInputFile(inputPath); err != nil {
		return errors.Wrap(err, "invalid input file")
	}

	if err := ensureOutputDir(outputDir); err != nil {
		return errors.Wrap(err, "failed to create output directory")
	}

	args := []string{
		"-y",            // Overwrite output files
		"-i", inputPath, // Input file
	}

	// Video codec
	if options.VideoCodec != "" {
		args = append(args, "-c:v", options.VideoCodec)
	}

	// Video quality/bitrate - use one or the other, not both
	if options.VideoBitrate != "" {
		args = append(args, "-b:v", options.VideoBitrate)
		// For x264 with specific bitrate, use reasonable quality preset
		if options.VideoCodec == CodecLibX264 {
			args = append(args, "-preset", PresetMedium)
		}
	} else if options.VideoQuality != "" {
		args = append(args, "-crf", options.VideoQuality)
	}

	// Resolution - use scale filter instead of -s to handle aspect ratio properly
	if options.Resolution != "" {
		// Parse resolution and use scale filter
		if strings.Contains(options.Resolution, "x") {
			parts := strings.Split(options.Resolution, "x")
			if len(parts) == 2 {
				width, height := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
				scaleFilter := fmt.Sprintf("scale=%s:%s", width, height)
				args = append(args, "-vf", scaleFilter)
			}
		}
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

	// HLS specific options
	args = append(args, "-f", FormatHLS)
	args = append(args, "-hls_time", options.SegmentDuration)
	args = append(args, "-hls_playlist_type", options.PlaylistType)

	// Segment file format
	segmentFilename := fmt.Sprintf("%s_%%03d.%s", options.SegmentPrefix, options.SegmentFormat)
	args = append(args, "-hls_segment_filename", filepath.Join(outputDir, segmentFilename))

	// Encryption
	if options.EnableEncryption && options.KeyInfoFile != "" {
		args = append(args, "-hls_key_info_file", options.KeyInfoFile)
	}

	// Custom arguments
	args = append(args, options.CustomArgs...)

	// Output playlist file
	playlistPath := filepath.Join(outputDir, options.PlaylistName)
	args = append(args, playlistPath)

	logger.Global().Info("Starting video segmentation",
		zap.String("input", inputPath),
		zap.String("outputDir", outputDir),
		zap.String("segmentDuration", options.SegmentDuration),
		zap.String("playlistPath", playlistPath))

	return f.runCommand(ctx, args, progressCallback)
}

// SegmentVideoMultiQuality segments a video into HLS format with multiple quality levels
func (f *FFmpeg) SegmentVideoMultiQuality(ctx context.Context, inputPath, outputDir string, qualities []SegmentationOptions, progressCallback ProgressCallback) error {
	if err := validateInputFile(inputPath); err != nil {
		return errors.Wrap(err, "invalid input file")
	}

	if err := ensureOutputDir(outputDir); err != nil {
		return errors.Wrap(err, "failed to create output directory")
	}

	// Create master playlist
	masterPlaylistPath := filepath.Join(outputDir, MasterPlaylistName)
	var masterPlaylistLines []string
	masterPlaylistLines = append(masterPlaylistLines, "#EXTM3U")
	masterPlaylistLines = append(masterPlaylistLines, "#EXT-X-VERSION:3")

	for _, quality := range qualities {
		// Create quality-specific directory
		qualityDir := filepath.Join(outputDir, quality.QualityName)
		if err := ensureOutputDir(qualityDir); err != nil {
			return errors.Wrapf(err, "failed to create quality directory %s", quality.QualityName)
		}

		// Segment this quality level
		logger.Global().Info("Segmenting quality level",
			zap.String("quality", quality.QualityName),
			zap.Int("total", len(qualities)),
			zap.String("resolution", quality.Resolution),
			zap.String("bitrate", quality.VideoBitrate))

		if err := f.SegmentVideo(ctx, inputPath, qualityDir, quality, progressCallback); err != nil {
			return errors.Wrapf(err, "failed to segment quality level %s", quality.QualityName)
		}

		// Add to master playlist
		bandwidth := extractBandwidth(quality.VideoBitrate, quality.AudioBitrate)
		resolution := quality.Resolution

		if bandwidth > 0 {
			streamInfo := fmt.Sprintf("#EXT-X-STREAM-INF:BANDWIDTH=%d", bandwidth)
			if resolution != "" {
				streamInfo += fmt.Sprintf(",RESOLUTION=%s", resolution)
			}
			masterPlaylistLines = append(masterPlaylistLines, streamInfo)
		}

		relativePath := fmt.Sprintf("%s/%s", quality.QualityName, DefaultPlaylistName)
		masterPlaylistLines = append(masterPlaylistLines, relativePath)
	}

	// Write master playlist
	masterPlaylistContent := strings.Join(masterPlaylistLines, "\n")
	if err := os.WriteFile(masterPlaylistPath, []byte(masterPlaylistContent), 0644); err != nil {
		return errors.Wrap(err, "failed to write master playlist")
	}

	logger.Global().Info("Created master playlist",
		zap.String("path", masterPlaylistPath),
		zap.Int("qualities", len(qualities)))

	return nil
}

// CreateDASHSegments creates DASH (Dynamic Adaptive Streaming over HTTP) segments
func (f *FFmpeg) CreateDASHSegments(ctx context.Context, inputPath, outputDir string, segmentDuration string, progressCallback ProgressCallback) error {
	if err := validateInputFile(inputPath); err != nil {
		return errors.Wrap(err, "invalid input file")
	}

	if err := ensureOutputDir(outputDir); err != nil {
		return errors.Wrap(err, "failed to create output directory")
	}

	args := []string{
		"-y",            // Overwrite output files
		"-i", inputPath, // Input file
		"-f", FormatDASH,
		"-seg_duration", segmentDuration,
		"-use_template", "1",
		"-use_timeline", "1",
		"-init_seg_name", DASHInitSegmentPattern,
		"-media_seg_name", DASHMediaSegmentPattern,
	}

	// Output manifest file
	manifestPath := filepath.Join(outputDir, DASHManifestName)
	args = append(args, manifestPath)

	logger.Global().Info("Creating DASH segments",
		zap.String("input", inputPath),
		zap.String("outputDir", outputDir),
		zap.String("segmentDuration", segmentDuration))

	return f.runCommand(ctx, args, progressCallback)
}

// CreatePreviewClips creates short preview clips from a video
func (f *FFmpeg) CreatePreviewClips(ctx context.Context, inputPath, outputDir string, clipDuration string, intervalSeconds int, progressCallback ProgressCallback) error {
	if err := validateInputFile(inputPath); err != nil {
		return errors.Wrap(err, "invalid input file")
	}

	if err := ensureOutputDir(outputDir); err != nil {
		return errors.Wrap(err, "failed to create output directory")
	}

	// Get video duration first
	probeInfo, err := f.ProbeFile(ctx, inputPath)
	if err != nil {
		return errors.Wrap(err, "failed to probe input file")
	}

	duration, err := parseTimeString(probeInfo.Format.Duration)
	if err != nil {
		return errors.Wrap(err, "failed to parse video duration")
	}

	totalSeconds := int(duration.Seconds())
	clipCount := 0

	for startTime := 0; startTime < totalSeconds; startTime += intervalSeconds {
		clipPath := filepath.Join(outputDir, fmt.Sprintf("preview_%03d.mp4", clipCount))

		args := []string{
			"-y",            // Overwrite output files
			"-i", inputPath, // Input file
			"-ss", fmt.Sprintf("%d", startTime), // Start time
			"-t", clipDuration, // Duration
			"-c:v", CodecLibX264,
			"-c:a", CodecAAC,
			"-movflags", OptionFastStart,
			clipPath,
		}

		logger.Global().Debug("Creating preview clip",
			zap.Int("clipNumber", clipCount),
			zap.Int("startTime", startTime))

		if err := f.runCommand(ctx, args, nil); err != nil {
			logger.Global().Warn("Failed to create preview clip",
				zap.Int("clipNumber", clipCount),
				zap.Error(err))
			continue
		}

		clipCount++
	}

	logger.Global().Info("Created preview clips",
		zap.String("outputDir", outputDir),
		zap.Int("totalClips", clipCount))

	return nil
}

// extractBandwidth extracts bandwidth from bitrate strings
func extractBandwidth(videoBitrate, audioBitrate string) int {
	totalBandwidth := 0

	if videoBitrate != "" {
		if bandwidth := parseBitrateString(videoBitrate); bandwidth > 0 {
			totalBandwidth += bandwidth
		}
	}

	if audioBitrate != "" {
		if bandwidth := parseBitrateString(audioBitrate); bandwidth > 0 {
			totalBandwidth += bandwidth
		}
	}

	return totalBandwidth
}

// parseBitrateString parses bitrate strings like "1000k", "2M" to bits per second
func parseBitrateString(bitrate string) int {
	bitrate = strings.ToLower(strings.TrimSpace(bitrate))
	if bitrate == "" {
		return 0
	}

	multiplier := 1
	if strings.HasSuffix(bitrate, "k") {
		multiplier = 1000
		bitrate = strings.TrimSuffix(bitrate, "k")
	} else if strings.HasSuffix(bitrate, "m") {
		multiplier = 1000000
		bitrate = strings.TrimSuffix(bitrate, "m")
	}

	// Try to parse the numeric part
	var value float64
	if _, err := fmt.Sscanf(bitrate, "%f", &value); err != nil {
		return 0
	}

	return int(value * float64(multiplier))
}

package ffmpeg

import (
	"context"
	"encoding/json"
	"os/exec"
	"strings"

	"github.com/cockroachdb/errors"
	"go.uber.org/zap"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
)

// ProbeFile gets detailed information about a media file using ffprobe
func (f *FFmpeg) ProbeFile(ctx context.Context, inputPath string) (*ProbeInfo, error) {
	if err := validateInputFile(inputPath); err != nil {
		return nil, errors.Wrap(err, "invalid input file")
	}

	// Use ffprobe instead of ffmpeg for probing
	probePath := strings.Replace(f.binaryPath, "ffmpeg", "ffprobe", 1)

	args := []string{
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		inputPath,
	}

	cmd := exec.CommandContext(ctx, probePath, args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, errors.Wrap(err, "failed to probe file")
	}

	var probeInfo ProbeInfo
	if err := json.Unmarshal(output, &probeInfo); err != nil {
		return nil, errors.Wrap(err, "failed to parse probe output")
	}

	logger.Global().Debug("Probed file",
		zap.String("file", inputPath),
		zap.String("format", probeInfo.Format.FormatName),
		zap.String("duration", probeInfo.Format.Duration),
		zap.Int("streams", len(probeInfo.Streams)))

	return &probeInfo, nil
}

// GetVideoInfo extracts video-specific information from probe data
func (f *FFmpeg) GetVideoInfo(ctx context.Context, inputPath string) (*StreamInfo, error) {
	probeInfo, err := f.ProbeFile(ctx, inputPath)
	if err != nil {
		return nil, err
	}

	for _, stream := range probeInfo.Streams {
		if stream.CodecType == "video" {
			return &stream, nil
		}
	}

	return nil, errors.New("no video stream found")
}

// GetAudioInfo extracts audio-specific information from probe data
func (f *FFmpeg) GetAudioInfo(ctx context.Context, inputPath string) (*StreamInfo, error) {
	probeInfo, err := f.ProbeFile(ctx, inputPath)
	if err != nil {
		return nil, err
	}

	for _, stream := range probeInfo.Streams {
		if stream.CodecType == "audio" {
			return &stream, nil
		}
	}

	return nil, errors.New("no audio stream found")
}

// GetDuration returns the duration of a media file
func (f *FFmpeg) GetDuration(ctx context.Context, inputPath string) (string, error) {
	probeInfo, err := f.ProbeFile(ctx, inputPath)
	if err != nil {
		return "", err
	}

	return probeInfo.Format.Duration, nil
}

// GetFrameRate returns the frame rate of a video file
func (f *FFmpeg) GetFrameRate(ctx context.Context, inputPath string) (string, error) {
	videoInfo, err := f.GetVideoInfo(ctx, inputPath)
	if err != nil {
		return "", err
	}

	return videoInfo.RFrameRate, nil
}

// GetResolution returns the resolution of a video file
func (f *FFmpeg) GetResolution(ctx context.Context, inputPath string) (int, int, error) {
	videoInfo, err := f.GetVideoInfo(ctx, inputPath)
	if err != nil {
		return 0, 0, err
	}

	return videoInfo.Width, videoInfo.Height, nil
}

// GetBitrate returns the bitrate of a media file
func (f *FFmpeg) GetBitrate(ctx context.Context, inputPath string) (string, error) {
	probeInfo, err := f.ProbeFile(ctx, inputPath)
	if err != nil {
		return "", err
	}

	return probeInfo.Format.BitRate, nil
}

// IsVideoFile checks if a file is a video file
func (f *FFmpeg) IsVideoFile(ctx context.Context, inputPath string) (bool, error) {
	_, err := f.GetVideoInfo(ctx, inputPath)
	if err != nil {
		if strings.Contains(err.Error(), "no video stream found") {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// IsAudioFile checks if a file is an audio file
func (f *FFmpeg) IsAudioFile(ctx context.Context, inputPath string) (bool, error) {
	_, err := f.GetAudioInfo(ctx, inputPath)
	if err != nil {
		if strings.Contains(err.Error(), "no audio stream found") {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// GetCodecInfo returns codec information for all streams
func (f *FFmpeg) GetCodecInfo(ctx context.Context, inputPath string) (map[string][]string, error) {
	probeInfo, err := f.ProbeFile(ctx, inputPath)
	if err != nil {
		return nil, err
	}

	codecInfo := make(map[string][]string)

	for _, stream := range probeInfo.Streams {
		if codecInfo[stream.CodecType] == nil {
			codecInfo[stream.CodecType] = make([]string, 0)
		}
		codecInfo[stream.CodecType] = append(codecInfo[stream.CodecType], stream.CodecName)
	}

	return codecInfo, nil
}

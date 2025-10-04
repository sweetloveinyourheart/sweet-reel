package ffmpeg

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"go.uber.org/zap"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
)

// FFmpeg represents an FFmpeg wrapper instance
type FFmpeg struct {
	binaryPath string
}

// New creates a new FFmpeg instance
func New() *FFmpeg {
	return &FFmpeg{
		binaryPath: "ffmpeg", // assumes ffmpeg is in PATH
	}
}

// NewWithBinaryPath creates a new FFmpeg instance with custom binary path
func NewWithBinaryPath(binaryPath string) *FFmpeg {
	return &FFmpeg{
		binaryPath: binaryPath,
	}
}

// SetBinaryPath sets the path to the FFmpeg binary
func (f *FFmpeg) SetBinaryPath(path string) {
	f.binaryPath = path
}

// IsAvailable checks if FFmpeg is available and working
func (f *FFmpeg) IsAvailable(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, f.binaryPath, "-version")
	return cmd.Run()
}

// GetVersion returns the FFmpeg version information
func (f *FFmpeg) GetVersion(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, f.binaryPath, "-version")
	output, err := cmd.Output()
	if err != nil {
		return "", errors.Wrap(err, "failed to get FFmpeg version")
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0]), nil
	}

	return "", errors.New("unable to parse FFmpeg version")
}

// runCommand executes an FFmpeg command with the given arguments
func (f *FFmpeg) runCommand(ctx context.Context, args []string, progressCallback ProgressCallback) error {
	cmd := exec.CommandContext(ctx, f.binaryPath, args...)

	// Log the full command for debugging
	logger.Global().Info("Executing FFmpeg command",
		zap.String("command", f.binaryPath),
		zap.Strings("args", args))

	// Capture stderr for progress monitoring and error reporting
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return errors.Wrap(err, "failed to create stderr pipe")
	}

	if err := cmd.Start(); err != nil {
		return errors.Wrap(err, "failed to start FFmpeg process")
	}

	// Channel to collect stderr output
	stderrLines := make(chan string, 100)
	done := make(chan bool)

	// Monitor stderr output
	go func() {
		defer close(stderrLines)
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			stderrLines <- line
			logger.Global().Debug("FFmpeg output", zap.String("line", line))
		}
		done <- true
	}()

	// Monitor progress if callback is provided
	if progressCallback != nil {
		go func() {
			var duration time.Duration
			for line := range stderrLines {
				// Parse duration from the beginning
				if strings.Contains(line, "Duration:") && duration == 0 {
					if d := parseDuration(line); d > 0 {
						duration = d
					}
				}

				// Parse current time progress
				if strings.Contains(line, "time=") && duration > 0 {
					if current := parseCurrentTime(line); current > 0 {
						progress := float64(current) / float64(duration) * 100
						if progress > 100 {
							progress = 100
						}

						progressCallback(ProgressInfo{
							Percentage: progress,
							Duration:   duration,
							Current:    current,
							Speed:      parseSpeed(line),
							Bitrate:    parseBitrate(line),
						})
					}
				}
			}
		}()
	}

	// Wait for command to complete
	cmdErr := cmd.Wait()

	// Wait for stderr collection to complete
	<-done

	// Collect any remaining stderr output for error reporting
	var stderrOutput []string
	for line := range stderrLines {
		stderrOutput = append(stderrOutput, line)
	}

	if cmdErr != nil {
		// Log stderr output for debugging
		if len(stderrOutput) > 0 {
			logger.Global().Error("FFmpeg stderr output",
				zap.Strings("stderr", stderrOutput),
				zap.Error(cmdErr))
		}
		return errors.Wrapf(cmdErr, "FFmpeg process failed")
	}

	return nil
}

// parseDuration extracts duration from FFmpeg output line
func parseDuration(line string) time.Duration {
	// Look for pattern: Duration: 00:01:30.45
	parts := strings.Split(line, "Duration:")
	if len(parts) < 2 {
		return 0
	}

	timeStr := strings.TrimSpace(strings.Split(parts[1], ",")[0])
	duration, err := parseTimeString(timeStr)
	if err != nil {
		return 0
	}

	return duration
}

// parseCurrentTime extracts current processing time from FFmpeg output line
func parseCurrentTime(line string) time.Duration {
	// Look for pattern: time=00:01:15.30
	parts := strings.Split(line, "time=")
	if len(parts) < 2 {
		return 0
	}

	timeStr := strings.TrimSpace(strings.Split(parts[1], " ")[0])
	duration, err := parseTimeString(timeStr)
	if err != nil {
		return 0
	}

	return duration
}

// parseSpeed extracts processing speed from FFmpeg output line
func parseSpeed(line string) string {
	// Look for pattern: speed=2.34x
	parts := strings.Split(line, "speed=")
	if len(parts) < 2 {
		return ""
	}

	return strings.TrimSpace(strings.Split(parts[1], " ")[0])
}

// parseBitrate extracts bitrate from FFmpeg output line
func parseBitrate(line string) string {
	// Look for pattern: bitrate=1234.5kbits/s
	parts := strings.Split(line, "bitrate=")
	if len(parts) < 2 {
		return ""
	}

	return strings.TrimSpace(strings.Split(parts[1], " ")[0])
}

// parseTimeString converts time string (HH:MM:SS.ms) to time.Duration
func parseTimeString(timeStr string) (time.Duration, error) {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 3 {
		return 0, errors.New("invalid time format")
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}

	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}

	secondsParts := strings.Split(parts[2], ".")
	seconds, err := strconv.Atoi(secondsParts[0])
	if err != nil {
		return 0, err
	}

	var milliseconds int
	if len(secondsParts) > 1 {
		// Pad or truncate to 3 digits for milliseconds
		msStr := secondsParts[1]
		if len(msStr) > 3 {
			msStr = msStr[:3]
		} else {
			for len(msStr) < 3 {
				msStr += "0"
			}
		}
		milliseconds, err = strconv.Atoi(msStr)
		if err != nil {
			return 0, err
		}
	}

	total := time.Duration(hours)*time.Hour +
		time.Duration(minutes)*time.Minute +
		time.Duration(seconds)*time.Second +
		time.Duration(milliseconds)*time.Millisecond

	return total, nil
}

// createTempFile creates a temporary file with the given extension
func createTempFile(ext string) (string, error) {
	tmpDir := os.TempDir()
	filename := fmt.Sprintf("ffmpeg_%d%s", time.Now().UnixNano(), ext)
	return filepath.Join(tmpDir, filename), nil
}

// validateInputFile checks if the input file exists and is readable
func validateInputFile(inputPath string) error {
	if inputPath == "" {
		return errors.New("input path cannot be empty")
	}

	info, err := os.Stat(inputPath)
	if err != nil {
		return errors.Wrap(err, "input file not accessible")
	}

	if info.IsDir() {
		return errors.New("input path is a directory, not a file")
	}

	return nil
}

// ensureOutputDir creates the output directory if it doesn't exist
func ensureOutputDir(outputPath string) error {
	// Check if outputPath is a file path or directory path
	dir := outputPath
	if filepath.Ext(outputPath) != "" {
		// If it has an extension, it's likely a file path, so get the directory
		dir = filepath.Dir(outputPath)
	}
	return os.MkdirAll(dir, 0755)
}

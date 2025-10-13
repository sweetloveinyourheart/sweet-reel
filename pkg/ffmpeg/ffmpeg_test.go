package ffmpeg

import (
	"context"
	"testing"
	"time"
)

func TestFFmpegAvailability(t *testing.T) {
	ff := New()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test if FFmpeg is available (this will fail if FFmpeg is not installed)
	err := ff.IsAvailable(ctx)
	if err != nil {
		t.Skipf("FFmpeg not available, skipping test: %v", err)
	}

	// Test getting version
	version, err := ff.GetVersion(ctx)
	if err != nil {
		t.Fatalf("Failed to get FFmpeg version: %v", err)
	}

	if version == "" {
		t.Error("Version should not be empty")
	}

	t.Logf("FFmpeg version: %s", version)
}

func TestTimeStringParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Duration
	}{
		{"00:01:30.45", 90*time.Second + 450*time.Millisecond},
		{"00:00:10.123", 10*time.Second + 123*time.Millisecond},
		{"01:30:45.678", 90*time.Minute + 45*time.Second + 678*time.Millisecond},
	}

	for _, test := range tests {
		result, err := parseTimeString(test.input)
		if err != nil {
			t.Errorf("Failed to parse time string %s: %v", test.input, err)
			continue
		}

		if result != test.expected {
			t.Errorf("For input %s, expected %v, got %v", test.input, test.expected, result)
		}
	}
}

func TestBitrateStringParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"1000k", 1000000},
		{"2M", 2000000},
		{"500", 500},
		{"1.5M", 1500000},
		{"", 0},
		{"invalid", 0},
	}

	for _, test := range tests {
		result := parseBitrateString(test.input)
		if result != test.expected {
			t.Errorf("For input %s, expected %d, got %d", test.input, test.expected, result)
		}
	}
}

func TestExtractBandwidth(t *testing.T) {
	videoBitrate := "2000k"
	audioBitrate := "128k"

	bandwidth := extractBandwidth(videoBitrate, audioBitrate)
	expected := 2000000 + 128000 // 2M + 128k

	if bandwidth != expected {
		t.Errorf("Expected bandwidth %d, got %d", expected, bandwidth)
	}
}

func TestValidateInputFile(t *testing.T) {
	// Test empty path
	err := validateInputFile("")
	if err == nil {
		t.Error("Expected error for empty path")
	}

	// Test non-existent file
	err = validateInputFile("/nonexistent/file.mp4")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

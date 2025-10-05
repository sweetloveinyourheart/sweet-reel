package s3

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractFileIDFromKey(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected string
	}{
		{
			name:     "simple file with extension",
			key:      "abc123.jpg",
			expected: "abc123",
		},
		{
			name:     "file in directory",
			key:      "path/to/file/abc123.jpg",
			expected: "abc123",
		},
		{
			name:     "complex path with multiple dots",
			key:      "path/to/file/abc123.version.1.jpg",
			expected: "abc123.version.1",
		},
		{
			name:     "no extension",
			key:      "path/to/file/abc123",
			expected: "abc123",
		},
		{
			name:     "only filename no extension",
			key:      "abc123",
			expected: "abc123",
		},
		{
			name:     "dot at start",
			key:      ".hidden",
			expected: "",
		},
		{
			name:     "empty string",
			key:      "",
			expected: "",
		},
		{
			name:     "just extension",
			key:      ".jpg",
			expected: "",
		},
		{
			name:     "file with UUID",
			key:      "original/video/550e8400-e29b-41d4-a716-446655440000.mp4",
			expected: "550e8400-e29b-41d4-a716-446655440000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractFileIDFromKey(tt.key)
			require.Equal(t, tt.expected, result)
		})
	}
}

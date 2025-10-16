package s3

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractBucketAndKeyFromEventMessage(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedBucket string
		expectedKey    string
	}{
		{
			name:           "normal path",
			path:           "mybucket/path/to/file.jpg",
			expectedBucket: "mybucket",
			expectedKey:    "path/to/file.jpg",
		},
		{
			name:           "path with multiple slashes",
			path:           "mybucket/path/with/multiple/levels/file.txt",
			expectedBucket: "mybucket",
			expectedKey:    "path/with/multiple/levels/file.txt",
		},
		{
			name:           "bucket only",
			path:           "mybucket",
			expectedBucket: "mybucket",
			expectedKey:    "",
		},
		{
			name:           "empty string",
			path:           "",
			expectedBucket: "",
			expectedKey:    "",
		},
		{
			name:           "trailing slash after bucket",
			path:           "mybucket/",
			expectedBucket: "mybucket",
			expectedKey:    "",
		},
		{
			name:           "bucket with dots",
			path:           "my.bucket.com/file.jpg",
			expectedBucket: "my.bucket.com",
			expectedKey:    "file.jpg",
		},
		{
			name:           "key with special characters",
			path:           "mybucket/file with spaces & special chars!.jpg",
			expectedBucket: "mybucket",
			expectedKey:    "file with spaces & special chars!.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bucket, key := ExtractBucketAndKeyFromEventMessage(tt.path)
			require.Equal(t, tt.expectedBucket, bucket)
			require.Equal(t, tt.expectedKey, key)
		})
	}
}

func TestExtractFilenameAndExt(t *testing.T) {
	tests := []struct {
		name         string
		filename     string
		expectedName string
		expectedExt  string
	}{
		{
			name:         "simple file with extension",
			filename:     "abc123.jpg",
			expectedName: "abc123",
			expectedExt:  ".jpg",
		},
		{
			name:         "file with multiple dots",
			filename:     "abc123.version.1.jpg",
			expectedName: "abc123.version.1",
			expectedExt:  ".jpg",
		},
		{
			name:         "no extension",
			filename:     "abc123",
			expectedName: "abc123",
			expectedExt:  "",
		},
		{
			name:         "dot at start",
			filename:     ".hidden",
			expectedName: "",
			expectedExt:  ".hidden",
		},
		{
			name:         "empty string",
			filename:     "",
			expectedName: "",
			expectedExt:  "",
		},
		{
			name:         "just extension",
			filename:     ".jpg",
			expectedName: "",
			expectedExt:  ".jpg",
		},
		{
			name:         "file with uppercase extension",
			filename:     "document.PDF",
			expectedName: "document",
			expectedExt:  ".PDF",
		},
		{
			name:         "file with path",
			filename:     "path/to/file.txt",
			expectedName: "file",
			expectedExt:  ".txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, ext := ExtractFilenameAndExt(tt.filename)
			require.Equal(t, tt.expectedName, name)
			require.Equal(t, tt.expectedExt, ext)
		})
	}
}

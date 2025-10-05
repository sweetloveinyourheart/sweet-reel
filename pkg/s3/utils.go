package s3

import (
	"path/filepath"
	"strings"
)

// ExtractFileIDFromKey extracts the file ID from an S3 key.
// Assumes that the file ID is the last component of the path without the extension.
// For example, from "path/to/file/abc123.jpg" it extracts "abc123"
func ExtractFileIDFromKey(key string) string {
	// Get the base filename from the path
	filename := filepath.Base(key)

	// Remove the extension if any
	fileID := strings.TrimSuffix(filename, filepath.Ext(filename))

	return fileID
}

// ExtractKey splits event message into bucket + key.
func ExtractBucketAndKeyFromEventMessage(path string) (string, string) {
	parts := strings.SplitN(path, "/", 2)
	if len(parts) < 2 {
		return path, ""
	}
	return parts[0], parts[1]
}

package s3

import (
	"path/filepath"
	"strings"
)

// ExtractKey splits event message into bucket + key.
func ExtractBucketAndKeyFromEventMessage(path string) (string, string) {
	parts := strings.SplitN(path, "/", 2)
	if len(parts) < 2 {
		return path, ""
	}
	return parts[0], parts[1]
}

// ExtractFilenameAndExt splits a filename into its name and extension parts.
// For example, from "abc123.jpg" it extracts "abc123" and ".jpg"
// Returns an error if the filename is empty.
func ExtractFilenameAndExt(filename string) (string, string) {
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)
	return name, ext
}

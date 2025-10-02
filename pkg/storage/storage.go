package storage

import "io"

const UrlExpirationSeconds = 600

type Storage interface {
	Download(key string, bucket string) ([]byte, error)
	Upload(key string, bucket string, file io.Reader, mimeType string) error
	GenerateUploadPublicUri(key string, bucket string, expirationSeconds uint32) (string, error)
	GenerateDownloadPublicUri(key string, bucket string, expirationSeconds uint32) (string, error)
	Delete(key string, bucket string) error
}

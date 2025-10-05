package mock

import (
	"io"

	"github.com/stretchr/testify/mock"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/s3"
)

// Ensure that MockStorage implements storage.Storag
var _ s3.S3Storage = (*MockStorage)(nil)

// MockStorage is a mock of the storage.Storage interface
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Upload(key string, bucket string, data io.Reader, contentType string) error {
	args := m.Called(key, bucket, data, contentType)
	return args.Error(0)
}

func (m *MockStorage) Download(key string, bucket string) ([]byte, error) {
	args := m.Called(key, bucket)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockStorage) GenerateUploadPublicUri(key string, bucket string, expirationSeconds uint32) (string, error) {
	args := m.Called(key, bucket, expirationSeconds)
	return args.String(0), args.Error(1)
}

func (m *MockStorage) GenerateDownloadPublicUri(key string, bucket string, expirationSeconds uint32) (string, error) {
	args := m.Called(key, bucket, expirationSeconds)
	return args.String(0), args.Error(1)
}

func (m *MockStorage) Delete(key string, bucket string) error {
	args := m.Called(key, bucket)
	return args.Error(0)
}

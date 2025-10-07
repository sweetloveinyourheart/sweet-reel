package mock

import (
	"io"

	"github.com/stretchr/testify/mock"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/s3"
)

// Ensure that MockS3 implements storage.Storag
var _ s3.S3Storage = (*MockS3)(nil)

// MockS3 is a mock of the storage.Storage interface
type MockS3 struct {
	mock.Mock
}

func (m *MockS3) Upload(key string, bucket string, data io.Reader, contentType string) error {
	args := m.Called(key, bucket, data, contentType)
	return args.Error(0)
}

func (m *MockS3) Download(key string, bucket string) ([]byte, error) {
	args := m.Called(key, bucket)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockS3) GenerateUploadPublicUri(key string, bucket string, expirationSeconds uint32) (string, error) {
	args := m.Called(key, bucket, expirationSeconds)
	return args.String(0), args.Error(1)
}

func (m *MockS3) GenerateDownloadPublicUri(key string, bucket string, expirationSeconds uint32) (string, error) {
	args := m.Called(key, bucket, expirationSeconds)
	return args.String(0), args.Error(1)
}

func (m *MockS3) Delete(key string, bucket string) error {
	args := m.Called(key, bucket)
	return args.Error(0)
}

package s3

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
)

const UrlExpirationSeconds = 600

type S3Storage interface {
	Download(key string, bucket string) ([]byte, error)
	Upload(key string, bucket string, file io.Reader, mimeType string) error
	GenerateUploadPublicUri(key string, bucket string, expirationSeconds uint32) (string, error)
	GenerateDownloadPublicUri(key string, bucket string, expirationSeconds uint32) (string, error)
	Delete(key string, bucket string) error
}

type s3Client struct {
	client *s3.Client
}

func CreateS3Client(ctx context.Context, config *Config) (*s3Client, error) {
	cfg, err := config.ToS3Config(ctx)
	if err != nil {
		return nil, err
	}

	var client *s3.Client
	if config.minioURL != "" {
		client = s3.NewFromConfig(cfg,
			func(o *s3.Options) {
				o.BaseEndpoint = aws.String(config.minioURL)
			})
	} else {
		client = s3.NewFromConfig(cfg)
	}

	return &s3Client{
		client: client,
	}, nil
}

func (s s3Client) Upload(key string, bucket string, file io.Reader, mimeType string) error {
	ctx := context.Background()

	input := &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(mimeType),
		ACL:         types.ObjectCannedACLPublicRead, // Make object publicly readable
	}

	_, err := s.client.PutObject(ctx, input)
	if err != nil {
		logger.GlobalSugared().Errorf("Failed to upload object %s to bucket %s: %v", key, bucket, err)
		return err
	}

	logger.GlobalSugared().Infof("Successfully uploaded %s to bucket %s", key, bucket)
	return nil
}

func (s s3Client) Download(key string, bucket string) ([]byte, error) {
	ctx := context.Background()

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	result, err := s.client.GetObject(ctx, input)
	if err != nil {
		logger.GlobalSugared().Errorf("Failed to download object %s from bucket %s: %v", key, bucket, err)
		return nil, err
	}
	defer result.Body.Close()

	// Read the entire object into memory
	var buf bytes.Buffer
	_, err = io.Copy(&buf, result.Body)
	if err != nil {
		logger.GlobalSugared().Errorf("Failed to read object %s from bucket %s: %v", key, bucket, err)
		return nil, err
	}

	logger.GlobalSugared().Infof("Successfully downloaded %s from bucket %s", key, bucket)
	return buf.Bytes(), nil
}

func (s s3Client) Delete(key string, bucket string) error {
	ctx := context.Background()

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	_, err := s.client.DeleteObject(ctx, input)
	if err != nil {
		logger.GlobalSugared().Errorf("Failed to delete object %s from bucket %s: %v", key, bucket, err)
		return err
	}

	logger.GlobalSugared().Infof("Successfully deleted %s from bucket %s", key, bucket)
	return nil
}

func (s s3Client) GenerateUploadPublicUri(key string, bucket string, expirationSeconds uint32) (string, error) {
	ctx := context.Background()

	input := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		ACL:    types.ObjectCannedACLPublicRead,
	}

	presignClient := s3.NewPresignClient(s.client)
	presignRequest, err := presignClient.PresignPutObject(ctx, input, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(expirationSeconds) * time.Second
	})

	if err != nil {
		logger.GlobalSugared().Errorf("Failed to generate upload URL for %s in bucket %s: %v", key, bucket, err)
		return "", err
	}

	logger.GlobalSugared().Infof("Generated upload URL for %s in bucket %s, expires in %d seconds", key, bucket, expirationSeconds)
	return presignRequest.URL, nil
}

func (s s3Client) GenerateDownloadPublicUri(key string, bucket string, expirationSeconds uint32) (string, error) {
	ctx := context.Background()

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	presignClient := s3.NewPresignClient(s.client)
	presignRequest, err := presignClient.PresignGetObject(ctx, input, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(expirationSeconds) * time.Second
	})

	if err != nil {
		logger.GlobalSugared().Errorf("Failed to generate download URL for %s in bucket %s: %v", key, bucket, err)
		return "", err
	}

	logger.GlobalSugared().Infof("Generated download URL for %s in bucket %s, expires in %d seconds", key, bucket, expirationSeconds)
	return presignRequest.URL, nil
}

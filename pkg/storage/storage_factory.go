package storage

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/storage/s3"
)

const Storage_S3 = "S3"

func GetStorageInstance(ctx context.Context, storageType string, storageConfig any) (Storage, error) {
	if storageType == Storage_S3 {
		config, ok := storageConfig.(*s3.Config)
		if !ok {
			return nil, errors.New("invalid S3 configuration")
		}
		return s3.CreateS3Client(ctx, config), nil
	}

	return nil, errors.New("incorrect storage configuration is provided")
}

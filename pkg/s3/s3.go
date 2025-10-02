package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
)

func CreateS3Client(ctx context.Context, config *Config) *s3.Client {
	cfg, err := config.ToS3Config(ctx)
	if err != nil {
		logger.GlobalSugared().Fatal(err)
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

	return client
}

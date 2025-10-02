package s3

import (
	"context"
	"fmt"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/config"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
)

type Config struct {
	AccessID string
	Region   string
	Secret   string

	minioURL string
}

func ServiceConfig(serviceType string) *Config {
	s3Region := config.Instance().GetString(fmt.Sprintf("%s.clientserver.aws.s3.region", serviceType))
	id := config.Instance().GetString(fmt.Sprintf("%s.clientserver.aws.s3.access.id", serviceType))
	secret := config.Instance().GetString(fmt.Sprintf("%s.clientserver.aws.s3.secret", serviceType))
	minioURL := config.Instance().GetString(fmt.Sprintf("%s.clientserver.minio.url", serviceType))

	return &Config{
		Region:   s3Region,
		AccessID: id,
		Secret:   secret,

		minioURL: minioURL,
	}
}

func (c *Config) ToS3Config(ctx context.Context) (aws.Config, error) {
	var cfg aws.Config
	var err error
	if c.AccessID != "" && c.Secret != "" {
		cfg, err = awsConfig.LoadDefaultConfig(ctx,
			awsConfig.WithRegion(c.Region),
			awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
				c.AccessID,
				c.Secret,
				"",
			)),
		)
	} else {
		cfg, err = awsConfig.LoadDefaultConfig(ctx,
			awsConfig.WithRegion(c.Region),
		)
	}

	if err != nil {
		logger.GlobalSugared().Fatal(err)
	}

	return cfg, nil
}

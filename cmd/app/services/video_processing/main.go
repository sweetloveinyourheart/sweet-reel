package appvideoprocessing

import (
	"context"
	"fmt"

	"github.com/samber/do"
	"github.com/spf13/cobra"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/cmdutil"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/config"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/kafka"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/s3"
	videoprocessing "github.com/sweetloveinyourheart/sweet-reel/services/video_processing"
)

const DEFAULT_VIDEO_PROCESSING_GRPC_PORT = 50055

const serviceType = "video_processing"
const envPrefix = "VIDEO_PROCESSING"

func Command(rootCmd *cobra.Command) *cobra.Command {
	var videoProcessingCommand = &cobra.Command{
		Use:   fmt.Sprintf("%s [flags]", serviceType),
		Short: fmt.Sprintf("Run as %s service", serviceType),
		Run: func(cmd *cobra.Command, args []string) {
			app, err := cmdutil.BoilerplateRun(serviceType)
			if err != nil {
				logger.GlobalSugared().Fatal(err)
			}

			if err := setupDependencies(app.Ctx()); err != nil {
				logger.GlobalSugared().Fatal(err)
			}

			if err := videoprocessing.InitializeRepos(app.Ctx()); err != nil {
				logger.GlobalSugared().Fatal(err)
			}

			app.Run()
		},
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			config.Instance().Set("service_prefix", serviceType)

			cmdutil.BoilerplateMetaConfig(serviceType)

			config.RegisterService(cmd, config.Service{
				Command: cmd,
			})
			config.AddDefaultServicePorts(cmd, rootCmd)

			return nil
		},
	}

	// config options
	config.Int64Default(videoProcessingCommand, fmt.Sprintf("%s.grpc.port", serviceType), "grpc-port", DEFAULT_VIDEO_PROCESSING_GRPC_PORT, "GRPC Port to listen on", "VIDEO_PROCESSING_GRPC_PORT")
	config.String(videoProcessingCommand, fmt.Sprintf("%s.aws.s3.region", serviceType), "aws_s3_region", "s3 region", "VIDEO_PROCESSING_AWS_S3_REGION")
	config.String(videoProcessingCommand, fmt.Sprintf("%s.aws.s3.access.id", serviceType), "aws_s3_access_id", "s3 access id", "VIDEO_PROCESSING_AWS_S3_ACCESS_ID")
	config.String(videoProcessingCommand, fmt.Sprintf("%s.aws.s3.secret", serviceType), "aws_s3_secret", "s3 secret", "VIDEO_PROCESSING_AWS_S3_SECRET")
	config.String(videoProcessingCommand, fmt.Sprintf("%s.aws.s3.bucket", serviceType), "s3_bucket", "s3 bucket", "VIDEO_PROCESSING_AWS_S3_BUCKET")
	config.StringDefault(videoProcessingCommand, fmt.Sprintf("%s.minio.url", serviceType), "minio-url", "", "MINIO URL", "VIDEO_PROCESSING_MINIO_URL")

	cmdutil.BoilerplateFlagsCore(videoProcessingCommand, serviceType, envPrefix)
	cmdutil.BoilerplateFlagsKafka(videoProcessingCommand, serviceType, envPrefix)
	cmdutil.BoilerplateSecureFlags(videoProcessingCommand, serviceType)

	return videoProcessingCommand
}

func setupDependencies(ctx context.Context) error {
	kafkaClient, err := initKafkaClient(ctx)
	if err != nil {
		return err
	}

	s3Client, err := initS3Client(ctx)
	if err != nil {
		return err
	}

	do.Provide(nil, func(i *do.Injector) (*kafka.Client, error) {
		return kafkaClient, nil
	})

	do.Provide(nil, func(i *do.Injector) (s3.S3Storage, error) {
		return s3Client, nil
	})

	return nil
}

func initKafkaClient(ctx context.Context) (*kafka.Client, error) {
	cfg := kafka.ServiceConfig(serviceType)

	// Init client
	client, err := kafka.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	// Gracefull shutdown
	go func() {
		<-ctx.Done()
		client.Close()
	}()

	return client, nil
}

func initS3Client(ctx context.Context) (s3.S3Storage, error) {
	cfg := s3.ServiceConfig(serviceType)
	s3Client, err := s3.CreateS3Client(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return s3Client, nil
}

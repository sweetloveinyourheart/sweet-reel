package videomanagement

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/do"
	"github.com/spf13/cobra"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/cmdutil"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/config"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/db"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/grpc"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/interceptors"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/kafka"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/s3"
	"github.com/sweetloveinyourheart/sweet-reel/proto/code/video_management/go/grpcconnect"
	videomanagement "github.com/sweetloveinyourheart/sweet-reel/services/video_management"
	"github.com/sweetloveinyourheart/sweet-reel/services/video_management/actions"
	"github.com/sweetloveinyourheart/sweet-reel/services/video_management/repos"
)

const DEFAULT_VIDEO_MANAGEMENT_GRPC_PORT = 50060

const serviceType = "video_management"
const envPrefix = "VIDEO_MANAGEMENT"
const dbTablePrefix = "video-management"

func Command(rootCmd *cobra.Command) *cobra.Command {
	var videoManagementCommand = &cobra.Command{
		Use:   fmt.Sprintf("%s [flags]", serviceType),
		Short: fmt.Sprintf("Run as %s service", serviceType),
		Run: func(cmd *cobra.Command, args []string) {
			app, err := cmdutil.BoilerplateRun(serviceType)
			if err != nil {
				logger.GlobalSugared().Fatal(err)
			}

			app.Migrations(videomanagement.FS, dbTablePrefix)

			if err := setupDependencies(app.Ctx()); err != nil {
				logger.GlobalSugared().Fatal(err)
			}

			if err := videomanagement.InitializeRepos(app.Ctx()); err != nil {
				logger.GlobalSugared().Fatal(err)
			}

			if err := setupGrpcServer(app.Ctx()); err != nil {
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
	config.Int64Default(videoManagementCommand, fmt.Sprintf("%s.grpc.port", serviceType), "grpc-port", DEFAULT_VIDEO_MANAGEMENT_GRPC_PORT, "GRPC Port to listen on", "VIDEO_MANAGEMENT_GRPC_PORT")
	config.String(videoManagementCommand, fmt.Sprintf("%s.aws.s3.region", serviceType), "aws_s3_region", "s3 region", "VIDEO_MANAGEMENT_AWS_S3_REGION")
	config.String(videoManagementCommand, fmt.Sprintf("%s.aws.s3.access.id", serviceType), "aws_s3_access_id", "s3 access id", "VIDEO_MANAGEMENT_AWS_S3_ACCESS_ID")
	config.String(videoManagementCommand, fmt.Sprintf("%s.aws.s3.secret", serviceType), "aws_s3_secret", "s3 secret", "VIDEO_MANAGEMENT_AWS_S3_SECRET")
	config.String(videoManagementCommand, fmt.Sprintf("%s.aws.s3.bucket", serviceType), "s3_bucket", "s3 bucket", "VIDEO_MANAGEMENT_AWS_S3_BUCKET")
	config.StringDefault(videoManagementCommand, fmt.Sprintf("%s.minio.url", serviceType), "minio-url", "", "MINIO URL", "VIDEO_MANAGEMENT_MINIO_URL")

	cmdutil.BoilerplateFlagsCore(videoManagementCommand, serviceType, envPrefix)
	cmdutil.BoilerplateFlagsKafka(videoManagementCommand, serviceType, envPrefix)
	cmdutil.BoilerplateFlagsDB(videoManagementCommand, serviceType, envPrefix)
	cmdutil.BoilerplateSecureFlags(videoManagementCommand, serviceType)

	return videoManagementCommand
}

func setupGrpcServer(ctx context.Context) error {
	signingKey := config.Instance().GetString(fmt.Sprintf("%s.secrets.token_signing_key", serviceType))
	actions := actions.NewActions(ctx, signingKey)

	opt := connect.WithInterceptors(
		interceptors.CommonConnectInterceptors(
			serviceType,
			signingKey,
			interceptors.ConnectServerAuthHandler(signingKey),
		)...,
	)
	path, handler := grpcconnect.NewVideoManagementHandler(actions, opt)
	port := config.Instance().GetUint64(fmt.Sprintf("%s.grpc.port", serviceType))

	go grpc.ServeBuf(ctx, path, handler, port, serviceType)

	return nil
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

	dbConn, err := initDBConnection()
	if err != nil {
		return err
	}

	videoRepo := repos.NewVideoRepository(dbConn)

	do.Provide(nil, func(i *do.Injector) (repos.IVideoRepository, error) {
		return videoRepo, nil
	})

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

func initDBConnection() (*pgxpool.Pool, error) {
	dbConn, err := db.NewDbWithWait(config.Instance().GetString(fmt.Sprintf("%s.db.url", serviceType)), db.DBOptions{
		TimeoutSec:      config.Instance().GetInt(fmt.Sprintf("%s.db.postgres.timeout", serviceType)),
		MaxOpenConns:    config.Instance().GetInt(fmt.Sprintf("%s.db.postgres.max_open_connections", serviceType)),
		MaxIdleConns:    config.Instance().GetInt(fmt.Sprintf("%s.db.postgres.max_idle_connections", serviceType)),
		ConnMaxLifetime: config.Instance().GetInt(fmt.Sprintf("%s.db.postgres.max_lifetime", serviceType)),
		ConnMaxIdleTime: config.Instance().GetInt(fmt.Sprintf("%s.db.postgres.max_idletime", serviceType)),
		EnableTracing:   config.Instance().GetBool(fmt.Sprintf("%s.db.tracing", serviceType)),
	})
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}

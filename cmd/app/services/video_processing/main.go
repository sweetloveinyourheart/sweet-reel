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
)

const DEFAULT_VIDEO_PROCESSING_GRPC_PORT = 50055

const serviceType = "video_processing"
const envPrefix = "VIDEO_PROCESSING"

func Command(rootCmd *cobra.Command) *cobra.Command {
	var dataProviderCommand = &cobra.Command{
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
	config.Int64Default(dataProviderCommand, fmt.Sprintf("%s.grpc.port", serviceType), "grpc-port", DEFAULT_VIDEO_PROCESSING_GRPC_PORT, "GRPC Port to listen on", "VIDEO_PROCESSING_GRPC_PORT")

	cmdutil.BoilerplateFlagsCore(dataProviderCommand, serviceType, envPrefix)
	cmdutil.BoilerplateFlagsKafka(dataProviderCommand, serviceType, envPrefix)
	cmdutil.BoilerplateSecureFlags(dataProviderCommand, serviceType)

	return dataProviderCommand
}

func setupDependencies(ctx context.Context) error {
	kafkaClient, err := initKafkaClient(ctx)
	if err != nil {
		return err
	}

	do.Provide(nil, func(i *do.Injector) (*kafka.Client, error) {
		return kafkaClient, nil
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

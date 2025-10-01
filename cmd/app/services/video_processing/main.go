package appvideoprocessing

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/cmdutil"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/config"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
)

const DEFAULT_VIDEO_PROCESSING_GRPC_PORT = 50055

const serviceType = "video_processing"
const defDBName = "video_processing_db"
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

			if err := setupDependencies(); err != nil {
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
			config.AddDefaultDatabase(cmd, defDBName)
			return nil
		},
	}

	// config options
	config.Int64Default(dataProviderCommand, "video_processing.grpc.port", "grpc-port", DEFAULT_VIDEO_PROCESSING_GRPC_PORT, "GRPC Port to listen on", "VIDEO_PROCESSING_GRPC_PORT")

	cmdutil.BoilerplateFlagsCore(dataProviderCommand, serviceType, envPrefix)

	return dataProviderCommand
}

func setupDependencies() error {
	return nil
}

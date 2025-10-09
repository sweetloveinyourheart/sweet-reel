package apigateway

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/cmdutil"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/config"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
	apigateway "github.com/sweetloveinyourheart/sweet-reel/services/api_gateway"
)

const DEFAULT_API_GATEWAY_HTTP_PORT = 8080

const serviceType = "api_gateway"
const envPrefix = "API_GATEWAY"

func Command(rootCmd *cobra.Command) *cobra.Command {
	var apiGatewayCommand = &cobra.Command{
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

			if err := apigateway.InitializeRepos(app.Ctx()); err != nil {
				logger.GlobalSugared().Fatal(err)
			}

			if err := setupHTTPServer(app.Ctx()); err != nil {
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
	config.Int64Default(apiGatewayCommand, fmt.Sprintf("%s.http.port", serviceType), "http-port", DEFAULT_API_GATEWAY_HTTP_PORT, "HTTP Port to listen on", "API_GATEWAY_HTTP_PORT")

	cmdutil.BoilerplateFlagsCore(apiGatewayCommand, serviceType, envPrefix)
	cmdutil.BoilerplateSecureFlags(apiGatewayCommand, serviceType)

	return apiGatewayCommand
}

func setupHTTPServer(ctx context.Context) error {
	port := config.Instance().GetUint64(fmt.Sprintf("%s.http.port", serviceType))
	signingKey := config.Instance().GetString(fmt.Sprintf("%s.secrets.token_signing_key", serviceType))

	server := apigateway.NewServer(ctx, port, signingKey)

	go server.Start(port)

	return nil
}

func setupDependencies(ctx context.Context) error {
	return nil
}

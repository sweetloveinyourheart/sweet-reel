package apigateway

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"github.com/samber/do"
	"github.com/spf13/cobra"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/cmdutil"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/config"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/interceptors"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
	authConnect "github.com/sweetloveinyourheart/sweet-reel/proto/code/auth/go/grpcconnect"
	userConnect "github.com/sweetloveinyourheart/sweet-reel/proto/code/user/go/grpcconnect"
	videoManagementConnect "github.com/sweetloveinyourheart/sweet-reel/proto/code/video_management/go/grpcconnect"
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
	config.StringDefault(apiGatewayCommand, fmt.Sprintf("%s.auth_server.url", serviceType), "auth-server-url", "http://auth:50070", "Auth server connection URL", "API_GATEWAY_AUTH_SERVER_URL")
	config.StringDefault(apiGatewayCommand, fmt.Sprintf("%s.user_server.url", serviceType), "user-server-url", "http://user:50065", "User server connection URL", "API_GATEWAY_USER_SERVER_URL")
	config.StringDefault(apiGatewayCommand, fmt.Sprintf("%s.video_management.url", serviceType), "video-management-url", "http://user:50060", "Video Management server connection URL", "API_GATEWAY_VIDEO_MANAGEMENT_SERVER_URL")

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
	authClient := authConnect.NewAuthServiceClient(
		http.DefaultClient,
		config.Instance().GetString(fmt.Sprintf("%s.auth_server.url", serviceType)),
		connect.WithInterceptors(interceptors.CommonConnectClientInterceptors(
			serviceType,
			config.Instance().GetString(fmt.Sprintf("%s.secrets.token_signing_key", serviceType)),
		)...),
	)

	userClient := userConnect.NewUserServiceClient(
		http.DefaultClient,
		config.Instance().GetString(fmt.Sprintf("%s.user_server.url", serviceType)),
		connect.WithInterceptors(interceptors.CommonConnectClientInterceptors(
			serviceType,
			config.Instance().GetString(fmt.Sprintf("%s.secrets.token_signing_key", serviceType)),
		)...),
	)

	videoManagementClient := videoManagementConnect.NewVideoManagementClient(
		http.DefaultClient,
		config.Instance().GetString(fmt.Sprintf("%s.video_management_server.url", serviceType)),
		connect.WithInterceptors(interceptors.CommonConnectClientInterceptors(
			serviceType,
			config.Instance().GetString(fmt.Sprintf("%s.secrets.token_signing_key", serviceType)),
		)...),
	)

	do.Provide(nil, func(i *do.Injector) (authConnect.AuthServiceClient, error) {
		return authClient, nil
	})

	do.Provide(nil, func(i *do.Injector) (userConnect.UserServiceClient, error) {
		return userClient, nil
	})

	do.Provide(nil, func(i *do.Injector) (videoManagementConnect.VideoManagementClient, error) {
		return videoManagementClient, nil
	})

	return nil
}

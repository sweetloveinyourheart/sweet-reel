package auth

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"github.com/samber/do"
	"github.com/spf13/cobra"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/cmdutil"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/config"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/grpc"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/interceptors"
	authInterceptor "github.com/sweetloveinyourheart/sweet-reel/pkg/interceptors/auth"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/oauth2"
	"github.com/sweetloveinyourheart/sweet-reel/proto/code/auth/go/grpcconnect"
	userConnect "github.com/sweetloveinyourheart/sweet-reel/proto/code/user/go/grpcconnect"
	auth "github.com/sweetloveinyourheart/sweet-reel/services/auth"
	"github.com/sweetloveinyourheart/sweet-reel/services/auth/actions"
)

const DEFAULT_AUTH_GRPC_PORT = 50060

const serviceType = "auth"
const envPrefix = "AUTH"

func Command(rootCmd *cobra.Command) *cobra.Command {
	var authManagementCommand = &cobra.Command{
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

			if err := auth.InitializeRepos(app.Ctx()); err != nil {
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
	config.Int64Default(authManagementCommand, fmt.Sprintf("%s.grpc.port", serviceType), "grpc-port", DEFAULT_AUTH_GRPC_PORT, "GRPC Port to listen on", "AUTH_GRPC_PORT")
	config.StringDefault(authManagementCommand, fmt.Sprintf("%s.user_server.url", serviceType), "http://user:50065", "User server connection URL", "AUTH_USER_SERVER_URL")
	config.String(authManagementCommand, fmt.Sprintf("%s.google.oauth.client_id", serviceType), "google-oauth-client-id", "", "AUTH_GOOGLE_OAUTH_CLIENT_ID")
	config.String(authManagementCommand, fmt.Sprintf("%s.google.oauth.client_secret", serviceType), "google-oauth-client-secret", "", "AUTH_GOOGLE_OAUTH_CLIENT_SECRET")
	config.String(authManagementCommand, fmt.Sprintf("%s.google.oauth.redirect_url", serviceType), "google-oauth-redirect-url", "", "AUTH_GOOGLE_OAUTH_REDIRECT_URL")

	cmdutil.BoilerplateFlagsCore(authManagementCommand, serviceType, envPrefix)
	cmdutil.BoilerplateSecureFlags(authManagementCommand, serviceType)

	return authManagementCommand
}

func setupGrpcServer(ctx context.Context) error {
	signingKey := config.Instance().GetString(fmt.Sprintf("%s.secrets.token_signing_key", serviceType))
	actions := actions.NewActions(ctx, signingKey)

	opt := connect.WithInterceptors(
		interceptors.CommonConnectInterceptors(
			serviceType,
			signingKey,
			interceptors.ConnectServerAuthHandler(signingKey),
			authInterceptor.WithOverride(actions),
		)...,
	)
	path, handler := grpcconnect.NewAuthServiceHandler(actions, opt)
	port := config.Instance().GetUint64(fmt.Sprintf("%s.grpc.port", serviceType))

	go grpc.ServeBuf(ctx, path, handler, port, serviceType)

	return nil
}

func setupDependencies(ctx context.Context) error {
	userClient := userConnect.NewUserServiceClient(
		http.DefaultClient,
		config.Instance().GetString(fmt.Sprintf("%s.user_server.url", serviceType)),
		connect.WithInterceptors(interceptors.CommonConnectClientInterceptors(
			serviceType,
			config.Instance().GetString(fmt.Sprintf("%s.secrets.token_signing_key", serviceType)),
		)...),
	)

	googleProvider, err := oauth2.NewProvider(oauth2.ProviderGoogle)
	if err != nil {
		return err
	}

	googleOAuthConfig := oauth2.NewConfig(
		googleProvider,
		config.Instance().GetString(fmt.Sprintf("%s.google.oauth.client_id", serviceType)),
		config.Instance().GetString(fmt.Sprintf("%s.google.oauth.client_secret", serviceType)),
		config.Instance().GetString(fmt.Sprintf("%s.google.oauth.redirect_url", serviceType)),
		oauth2.GoogleBasicScopes,
	)
	googleOAuthClient := oauth2.NewClient(googleOAuthConfig)

	do.Provide(nil, func(i *do.Injector) (userConnect.UserServiceClient, error) {
		return userClient, nil
	})

	do.ProvideNamed(nil, string(oauth2.ProviderGoogle), func(i *do.Injector) (oauth2.IOAuthClient, error) {
		return googleOAuthClient, nil
	})

	return nil
}

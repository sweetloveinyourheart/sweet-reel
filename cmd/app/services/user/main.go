package user

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
	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
	"github.com/sweetloveinyourheart/sweet-reel/proto/code/user/go/grpcconnect"
	user "github.com/sweetloveinyourheart/sweet-reel/services/user"
	"github.com/sweetloveinyourheart/sweet-reel/services/user/actions"
	"github.com/sweetloveinyourheart/sweet-reel/services/user/repos"
)

const DEFAULT_USER_GRPC_PORT = 50060

const serviceType = "user"
const envPrefix = "USER"
const dbTablePrefix = "user"

func Command(rootCmd *cobra.Command) *cobra.Command {
	var userManagementCommand = &cobra.Command{
		Use:   fmt.Sprintf("%s [flags]", serviceType),
		Short: fmt.Sprintf("Run as %s service", serviceType),
		Run: func(cmd *cobra.Command, args []string) {
			app, err := cmdutil.BoilerplateRun(serviceType)
			if err != nil {
				logger.GlobalSugared().Fatal(err)
			}

			app.Migrations(user.FS, dbTablePrefix)

			if err := setupDependencies(app.Ctx()); err != nil {
				logger.GlobalSugared().Fatal(err)
			}

			if err := user.InitializeRepos(app.Ctx()); err != nil {
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
	config.Int64Default(userManagementCommand, fmt.Sprintf("%s.grpc.port", serviceType), "grpc-port", DEFAULT_USER_GRPC_PORT, "GRPC Port to listen on", "USER_GRPC_PORT")

	cmdutil.BoilerplateFlagsCore(userManagementCommand, serviceType, envPrefix)
	cmdutil.BoilerplateFlagsDB(userManagementCommand, serviceType, envPrefix)
	cmdutil.BoilerplateSecureFlags(userManagementCommand, serviceType)

	return userManagementCommand
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
	path, handler := grpcconnect.NewUserServiceHandler(actions, opt)
	port := config.Instance().GetUint64(fmt.Sprintf("%s.grpc.port", serviceType))

	go grpc.ServeBuf(ctx, path, handler, port, serviceType)

	return nil
}

func setupDependencies(ctx context.Context) error {
	dbConn, err := initDBConnection()
	if err != nil {
		return err
	}

	userRepo := repos.NewUserRepository(dbConn)

	do.Provide(nil, func(i *do.Injector) (*pgxpool.Pool, error) {
		return dbConn, nil
	})

	do.Provide(nil, func(i *do.Injector) (repos.IUserRepository, error) {
		return userRepo, nil
	})

	return nil
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

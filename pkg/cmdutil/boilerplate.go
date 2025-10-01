package cmdutil

import (
	"context"
	"embed"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/config"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/db"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/stringsutil"
	"go.uber.org/zap"
)

type AppRun struct {
	serviceType string
	serviceKey  string
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	readyChan   chan bool
}

func st(s string) (serviceType string, serviceKey string) {
	serviceType = strings.ToLower(s)
	serviceKey = strings.ReplaceAll(serviceType, "_", "-")
	serviceKey = strings.ReplaceAll(serviceKey, "-", ".")
	return
}

func BoilerplateRun(serviceType string) (*AppRun, error) {
	serviceType, serviceKey := st(serviceType)

	ctx, cancel := context.WithCancel(context.Background())

	// Logger init.
	logLevel := config.Instance().GetString("log.level")
	loggerInstance := logger.New().SetStringLevel(logLevel)
	loggerInstance.SetAsGlobal()
	defer func() {
		_ = loggerInstance.Sync()
	}()

	serviceName := config.Instance().GetString("service")
	loggerInstance.Infow("starting service",
		"type", serviceType,
		"service", serviceName,
	)

	healthCheckPort := config.Instance().GetInt("healthcheck.port")
	healthCheckWebPort := config.Instance().GetInt("healthcheck.web.port")
	readyChan := StartHealthServices(ctx, serviceName, healthCheckPort, healthCheckWebPort)

	return &AppRun{
		serviceType: serviceType,
		serviceKey:  serviceKey,
		wg:          sync.WaitGroup{},
		ctx:         ctx,
		cancel:      cancel,
		readyChan:   readyChan,
	}, nil
}

func BoilerplateFlagsCore(command *cobra.Command, serviceType string, envPrefix string) {
	_, serviceKey := st(serviceType)
	envPrefix = strings.ToUpper(envPrefix)

	config.String(command, fmt.Sprintf("%s.id", serviceKey), "id", "Unique identifier for this services", fmt.Sprintf("%s_ID", envPrefix))
	config.String(command, fmt.Sprintf("%s.secrets.token_signing_key", serviceKey), "token-signing-key", "Signing key used for service to service tokens", fmt.Sprintf("%s_SECRETS_TOKEN_SIGNING_KEY", envPrefix))

	_ = command.MarkPersistentFlagRequired("id")
	_ = command.MarkPersistentFlagRequired("token-signing-key")
}

func BoilerplateFlagsAuth(command *cobra.Command, serviceType string, envPrefix string) {
	_, serviceKey := st(serviceType)
	envPrefix = strings.ToUpper(envPrefix)

	config.String(command, fmt.Sprintf("%s.secrets.token_signing_key", serviceKey), "token-signing-key", "Signing key used for service to service tokens", fmt.Sprintf("%s_SECRETS_TOKEN_SIGNING_KEY", envPrefix))

	_ = command.MarkPersistentFlagRequired("token-signing-key")
}

func BoilerplateMetaConfig(serviceType string) {
	_, serviceKey := st(serviceType)

	config.Instance().Set(config.ServerNamespace, "sweet-real")
	config.Instance().Set(config.ServerId, config.Instance().GetString(fmt.Sprintf("%s.id", serviceKey)))
	config.Instance().Set(config.ServerReplicaCount, config.Instance().GetInt64(fmt.Sprintf("%s.replicas", serviceKey)))
	config.Instance().Set(config.ServerReplicaNumber, config.Instance().GetInt64(fmt.Sprintf("%s.replica_num", serviceKey)))
}

func BoilerplateFlagsDB(command *cobra.Command, serviceType string, envPrefix string) {
	_, serviceKey := st(serviceType)

	config.String(command, fmt.Sprintf("%s.db.url", serviceKey), "db-url", "Database connection URL", fmt.Sprintf("%s_DB_URL", envPrefix))
	config.StringDefault(command, fmt.Sprintf("%s.db.read.url", serviceKey), "db-read-url", "", "Database connection readonly URL", fmt.Sprintf("%s_DB_READ_URL", envPrefix))
	config.StringDefault(command, fmt.Sprintf("%s.db.migrations.url", serviceKey), "db-migrations-url", "", "Database connection migrations URL", fmt.Sprintf("%s_DB_MIGRATIONS_URL", envPrefix))
	config.Int64Default(command, fmt.Sprintf("%s.db.postgres.timeout", serviceKey), "db-postgres-timeout", 60, "Timeout for postgres connection", fmt.Sprintf("%s_DB_POSTGRES_TIMEOUT", envPrefix))
	config.Int64Default(command, fmt.Sprintf("%s.db.postgres.max_open_connections", serviceKey), "db-postgres-max-open-connections", 500, "Maximum number of connections", fmt.Sprintf("%s_DB_POSTGRES_MAX_OPEN_CONNECTIONS", envPrefix))
	config.Int64Default(command, fmt.Sprintf("%s.db.postgres.max_idle_connections", serviceKey), "db-postgres-max-idle-connections", 50, "Maximum number of idle connections", fmt.Sprintf("%s_DB_POSTGRES_MAX_IDLE_CONNECTIONS", envPrefix))
	config.Int64Default(command, fmt.Sprintf("%s.db.postgres.max_lifetime", serviceKey), "db-postgres-connection-max-lifetime", 300, "Max connection lifetime in seconds", fmt.Sprintf("%s_DB_POSTGRES_CONNECTION_MAX_LIFETIME", envPrefix))
	config.Int64Default(command, fmt.Sprintf("%s.db.postgres.max_idletime", serviceKey), "db-postgres-connection-max-idletime", 180, "Max connection idle time in seconds", fmt.Sprintf("%s_DB_POSTGRES_CONNECTION_MAX_IDLETIME", envPrefix))

	_ = command.MarkPersistentFlagRequired("db-url")
}

func BoilerplateSecureFlags(command *cobra.Command, serviceType string) {
	_, serviceKey := st(serviceType)

	config.SecureFields(
		fmt.Sprintf("%s.db.url", serviceKey),
		fmt.Sprintf("%s.db.read.url", serviceKey),
		fmt.Sprintf("%s.db.migrations.url", serviceKey),
		fmt.Sprintf("%s.secrets.token_signing_key", serviceKey),
	)
}

func (a *AppRun) Migrations(fs embed.FS, prefix string) {
	// Ensure database migrations
	var migrationsURL = config.Instance().GetString(fmt.Sprintf("%s.db.migrations.url", a.serviceKey))
	if stringsutil.IsBlank(migrationsURL) {
		migrationsURL = config.Instance().GetString(fmt.Sprintf("%s.db.url", a.serviceKey))
	}
	if stringsutil.IsBlank(migrationsURL) {
		logger.GlobalSugared().Fatal("no database migrations URL provided", zap.String("service", a.serviceType), zap.String("serviceKey", a.serviceKey))
		return
	}
	db.PerformMigrations(fs, prefix, migrationsURL)
}

func (a *AppRun) Ctx() context.Context {
	return a.ctx
}

func (a *AppRun) Cancel() {
	a.cancel()
}

func (a *AppRun) Ready() {
	a.readyChan <- true
}

func (a *AppRun) Run() {
	signalMonitor := make(chan os.Signal, 1)
	signal.Notify(signalMonitor, os.Interrupt, syscall.SIGTERM)

	a.wg.Go(func() {
		<-signalMonitor
		a.cancel()
		a.wg.Done()
	})

	a.readyChan <- true

	// Wait for a signal or other termination event
	a.wg.Wait()
}

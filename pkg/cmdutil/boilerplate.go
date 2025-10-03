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
	"go.uber.org/zap"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/config"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/db"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/stringsutil"
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
	serviceKey = strings.ReplaceAll(serviceType, "-", ".")
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

func BoilerplateFlagsKafka(command *cobra.Command, serviceType string, envPrefix string) {
	_, serviceKey := st(serviceType)

	// Kafka broker configuration
	config.StringDefault(command, fmt.Sprintf("%s.kafka.brokers", serviceKey), "kafka-brokers", "localhost:9092", "Kafka broker addresses (comma-separated)", fmt.Sprintf("%s_KAFKA_BROKERS", envPrefix))

	// Producer settings
	config.Int32Default(command, fmt.Sprintf("%s.kafka.retry_max", serviceKey), "kafka-retry-max", 3, "Maximum number of retries for failed requests", fmt.Sprintf("%s_KAFKA_RETRY_MAX", envPrefix))
	config.Int64Default(command, fmt.Sprintf("%s.kafka.retry_backoff_ms", serviceKey), "kafka-retry-backoff-ms", 100, "Backoff time between retries in milliseconds", fmt.Sprintf("%s_KAFKA_RETRY_BACKOFF_MS", envPrefix))
	config.StringDefault(command, fmt.Sprintf("%s.kafka.required_acks", serviceKey), "kafka-required-acks", "all", "Required acknowledgments (none, leader, all)", fmt.Sprintf("%s_KAFKA_REQUIRED_ACKS", envPrefix))
	config.StringDefault(command, fmt.Sprintf("%s.kafka.compression", serviceKey), "kafka-compression", "snappy", "Compression codec (none, gzip, snappy, lz4, zstd)", fmt.Sprintf("%s_KAFKA_COMPRESSION", envPrefix))
	config.Int64Default(command, fmt.Sprintf("%s.kafka.flush_frequency_ms", serviceKey), "kafka-flush-frequency-ms", 500, "Frequency of flushing messages in milliseconds", fmt.Sprintf("%s_KAFKA_FLUSH_FREQUENCY_MS", envPrefix))
	config.Int32Default(command, fmt.Sprintf("%s.kafka.flush_messages", serviceKey), "kafka-flush-messages", 100, "Number of messages to buffer before flushing", fmt.Sprintf("%s_KAFKA_FLUSH_MESSAGES", envPrefix))
	config.Int32Default(command, fmt.Sprintf("%s.kafka.flush_bytes", serviceKey), "kafka-flush-bytes", 16384, "Number of bytes to buffer before flushing", fmt.Sprintf("%s_KAFKA_FLUSH_BYTES", envPrefix))
	config.BoolDefault(command, fmt.Sprintf("%s.kafka.idempotent_writes", serviceKey), "kafka-idempotent-writes", true, "Enable idempotent writes", fmt.Sprintf("%s_KAFKA_IDEMPOTENT_WRITES", envPrefix))

	// Security settings
	config.StringDefault(command, fmt.Sprintf("%s.kafka.security_protocol", serviceKey), "kafka-security-protocol", "PLAINTEXT", "Security protocol (PLAINTEXT, SSL, SASL_PLAINTEXT, SASL_SSL)", fmt.Sprintf("%s_KAFKA_SECURITY_PROTOCOL", envPrefix))
	config.StringDefault(command, fmt.Sprintf("%s.kafka.sasl_mechanism", serviceKey), "kafka-sasl-mechanism", "", "SASL mechanism (PLAIN, SCRAM-SHA-256, SCRAM-SHA-512)", fmt.Sprintf("%s_KAFKA_SASL_MECHANISM", envPrefix))
	config.StringDefault(command, fmt.Sprintf("%s.kafka.sasl_username", serviceKey), "kafka-sasl-username", "", "SASL username for authentication", fmt.Sprintf("%s_KAFKA_SASL_USERNAME", envPrefix))
	config.StringDefault(command, fmt.Sprintf("%s.kafka.sasl_password", serviceKey), "kafka-sasl-password", "", "SASL password for authentication", fmt.Sprintf("%s_KAFKA_SASL_PASSWORD", envPrefix))
	config.BoolDefault(command, fmt.Sprintf("%s.kafka.tls_enabled", serviceKey), "kafka-tls-enabled", false, "Enable TLS encryption", fmt.Sprintf("%s_KAFKA_TLS_ENABLED", envPrefix))
	config.Int32Default(command, fmt.Sprintf("%s.kafka.max_open_requests", serviceKey), "kafka-tls-enabled", 1, "Enable TLS encryption", fmt.Sprintf("%s_KAFKA_MAX_OPEN_REQUESTS", envPrefix))

	_ = command.MarkPersistentFlagRequired("kafka-brokers")
}

func BoilerplateSecureFlags(command *cobra.Command, serviceType string) {
	_, serviceKey := st(serviceType)

	config.SecureFields(
		fmt.Sprintf("%s.db.url", serviceKey),
		fmt.Sprintf("%s.db.read.url", serviceKey),
		fmt.Sprintf("%s.db.migrations.url", serviceKey),
		fmt.Sprintf("%s.secrets.token_signing_key", serviceKey),
		fmt.Sprintf("%s.kafka.brokers", serviceKey),
		fmt.Sprintf("%s.kafka.sasl_username", serviceKey),
		fmt.Sprintf("%s.kafka.sasl_password", serviceKey),
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

	a.wg.Add(1)
	go func() {
		<-signalMonitor
		a.cancel()
		a.wg.Done()
	}()

	a.readyChan <- true

	// Wait for a signal or other termination event
	a.wg.Wait()
}

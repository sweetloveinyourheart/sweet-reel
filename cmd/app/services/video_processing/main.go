package appvideoprocessing

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/IBM/sarama"
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
	config.Int64Default(dataProviderCommand, "video_processing.grpc.port", "grpc-port", DEFAULT_VIDEO_PROCESSING_GRPC_PORT, "GRPC Port to listen on", "VIDEO_PROCESSING_GRPC_PORT")

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
	cfg := kafka.DefaultConfig()

	if brokersStr := config.Instance().GetString("video_processing.kafka.brokers"); brokersStr != "" {
		cfg.Brokers = strings.Split(brokersStr, ",")
		for i, broker := range cfg.Brokers {
			cfg.Brokers[i] = strings.TrimSpace(broker)
		}
	}

	// Basic producer settings
	cfg.RetryMax = int(config.Instance().GetInt32("video_processing.kafka.retry_max"))
	cfg.RetryBackoff = time.Duration(config.Instance().GetInt64("video_processing.kafka.retry_backoff_ms")) * time.Millisecond
	cfg.FlushFrequency = time.Duration(config.Instance().GetInt64("video_processing.kafka.flush_frequency_ms")) * time.Millisecond
	cfg.FlushMessages = int(config.Instance().GetInt32("video_processing.kafka.flush_messages"))
	cfg.FlushBytes = int(config.Instance().GetInt32("video_processing.kafka.flush_bytes"))
	cfg.IdempotentWrites = config.Instance().GetBool("video_processing.kafka.idempotent_writes")

	// Parse required acks
	switch strings.ToLower(config.Instance().GetString("video_processing.kafka.required_acks")) {
	case "none", "0":
		cfg.RequiredAcks = sarama.NoResponse
	case "leader", "1":
		cfg.RequiredAcks = sarama.WaitForLocal
	case "all", "-1":
		cfg.RequiredAcks = sarama.WaitForAll
	default:
		cfg.RequiredAcks = sarama.WaitForAll
	}

	// Parse compression
	switch strings.ToLower(config.Instance().GetString("video_processing.kafka.compression")) {
	case "none":
		cfg.Compression = sarama.CompressionNone
	case "gzip":
		cfg.Compression = sarama.CompressionGZIP
	case "snappy":
		cfg.Compression = sarama.CompressionSnappy
	case "lz4":
		cfg.Compression = sarama.CompressionLZ4
	case "zstd":
		cfg.Compression = sarama.CompressionZSTD
	default:
		cfg.Compression = sarama.CompressionSnappy
	}

	// Security settings
	cfg.SecurityProtocol = config.Instance().GetString("video_processing.kafka.security_protocol")
	cfg.SASLMechanism = config.Instance().GetString("video_processing.kafka.sasl_mechanism")
	cfg.SASLUsername = config.Instance().GetString("video_processing.kafka.sasl_username")
	cfg.SASLPassword = config.Instance().GetString("video_processing.kafka.sasl_password")
	cfg.TLSEnabled = config.Instance().GetBool("video_processing.kafka.tls_enabled")

	// Init client
	client, err := kafka.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return client, nil
}

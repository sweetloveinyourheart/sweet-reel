package kafka

import (
	"fmt"
	"strings"
	"time"

	"github.com/IBM/sarama"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/config"
)

// Config holds the configuration for Kafka client
type Config struct {
	Brokers          []string
	RetryMax         int
	RetryBackoff     time.Duration
	RequiredAcks     sarama.RequiredAcks
	Compression      sarama.CompressionCodec
	FlushFrequency   time.Duration
	FlushMessages    int
	FlushBytes       int
	IdempotentWrites bool
	SecurityProtocol string
	SASLMechanism    string
	SASLUsername     string
	SASLPassword     string
	TLSEnabled       bool
	MaxOpenRequests  int
}

// DefaultConfig returns a default configuration for Kafka
func DefaultConfig() *Config {
	return &Config{
		Brokers:          []string{"localhost:9092"},
		RetryMax:         3,
		RetryBackoff:     100 * time.Millisecond,
		RequiredAcks:     sarama.WaitForAll,
		Compression:      sarama.CompressionSnappy,
		FlushFrequency:   500 * time.Millisecond,
		FlushMessages:    100,
		FlushBytes:       16384, // 16KB
		IdempotentWrites: true,
		SecurityProtocol: "PLAINTEXT",
		MaxOpenRequests:  1,
	}
}

func ServiceConfig(serviceType string) *Config {
	cfg := DefaultConfig()

	if brokersStr := config.Instance().GetString(fmt.Sprintf("%s.kafka.brokers", serviceType)); brokersStr != "" {
		cfg.Brokers = strings.Split(brokersStr, ",")
		for i, broker := range cfg.Brokers {
			cfg.Brokers[i] = strings.TrimSpace(broker)
		}
	}

	// Basic producer settings
	cfg.RetryMax = int(config.Instance().GetInt32(fmt.Sprintf("%s.kafka.retry_max", serviceType)))
	cfg.RetryBackoff = time.Duration(config.Instance().GetInt64(fmt.Sprintf("%s.kafka.retry_backoff_ms", serviceType))) * time.Millisecond
	cfg.FlushFrequency = time.Duration(config.Instance().GetInt64(fmt.Sprintf("%s.kafka.flush_frequency_ms", serviceType))) * time.Millisecond
	cfg.FlushMessages = int(config.Instance().GetInt32(fmt.Sprintf("%s.kafka.flush_messages", serviceType)))
	cfg.FlushBytes = int(config.Instance().GetInt32(fmt.Sprintf("%s.kafka.flush_bytes", serviceType)))
	cfg.IdempotentWrites = config.Instance().GetBool(fmt.Sprintf("%s.kafka.idempotent_writes", serviceType))

	// Parse required acks
	switch strings.ToLower(config.Instance().GetString(fmt.Sprintf("%s.kafka.required_acks", serviceType))) {
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
	switch strings.ToLower(config.Instance().GetString(fmt.Sprintf("%s.kafka.compression", serviceType))) {
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
	cfg.SecurityProtocol = config.Instance().GetString(fmt.Sprintf("%s.kafka.security_protocol", serviceType))
	cfg.SASLMechanism = config.Instance().GetString(fmt.Sprintf("%s.kafka.sasl_mechanism", serviceType))
	cfg.SASLUsername = config.Instance().GetString(fmt.Sprintf("%s.kafka.sasl_username", serviceType))
	cfg.SASLPassword = config.Instance().GetString(fmt.Sprintf("%s.kafka.sasl_password", serviceType))
	cfg.TLSEnabled = config.Instance().GetBool(fmt.Sprintf("%s.kafka.tls_enabled", serviceType))
	cfg.MaxOpenRequests = int(config.Instance().GetInt32(fmt.Sprintf("%s.kafka.max_open_requests", serviceType)))

	return cfg
}

// ToSaramaConfig converts our config to Sarama config
func (c *Config) ToSaramaConfig() *sarama.Config {
	config := sarama.NewConfig()

	// Producer settings
	config.Producer.RequiredAcks = c.RequiredAcks
	config.Producer.Retry.Max = c.RetryMax
	config.Producer.Retry.Backoff = c.RetryBackoff
	config.Producer.Compression = c.Compression
	config.Producer.Flush.Frequency = c.FlushFrequency
	config.Producer.Flush.Messages = c.FlushMessages
	config.Producer.Flush.Bytes = c.FlushBytes
	config.Producer.Idempotent = c.IdempotentWrites
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	// Consumer settings
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()

	// Version
	config.Version = sarama.V3_6_0_0

	// Security settings
	config.Net.MaxOpenRequests = c.MaxOpenRequests

	if c.SecurityProtocol == "SASL_PLAINTEXT" || c.SecurityProtocol == "SASL_SSL" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = c.SASLUsername
		config.Net.SASL.Password = c.SASLPassword

		switch c.SASLMechanism {
		case "PLAIN":
			config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		case "SCRAM-SHA-256":
			config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256
		case "SCRAM-SHA-512":
			config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
		}
	}

	if c.TLSEnabled || c.SecurityProtocol == "SSL" || c.SecurityProtocol == "SASL_SSL" {
		config.Net.TLS.Enable = true
	}

	return config
}

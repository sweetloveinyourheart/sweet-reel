package kafka

import (
	"time"

	"github.com/IBM/sarama"
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
	}
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

package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"

	"github.com/sweetloveinyourheart/sweet-reel/pkg/logger"
)

// Producer wraps the Sarama sync producer
type Producer struct {
	producer sarama.SyncProducer
	config   *Config
	mu       sync.RWMutex
	closed   bool
}

// Message represents a Kafka message
type Message struct {
	Topic     string
	Key       string
	Value     any
	Headers   map[string]string
	Partition int32
	Timestamp time.Time
}

// NewProducer creates a new Kafka producer
func NewProducer(config *Config) (*Producer, error) {
	if config == nil {
		config = DefaultConfig()
	}

	saramaConfig := config.ToSaramaConfig()
	producer, err := sarama.NewSyncProducer(config.Brokers, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return &Producer{
		producer: producer,
		config:   config,
	}, nil
}

// Send sends a message to Kafka
func (p *Producer) Send(ctx context.Context, msg *Message) (partition int32, offset int64, err error) {
	p.mu.RLock()
	if p.closed {
		p.mu.RUnlock()
		return 0, 0, fmt.Errorf("producer is closed")
	}
	p.mu.RUnlock()

	// Convert value to bytes
	var valueBytes []byte
	switch v := msg.Value.(type) {
	case []byte:
		valueBytes = v
	case string:
		valueBytes = []byte(v)
	default:
		valueBytes, err = json.Marshal(v)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to marshal message value: %w", err)
		}
	}

	// Create Sarama message
	saramaMsg := &sarama.ProducerMessage{
		Topic:     msg.Topic,
		Value:     sarama.ByteEncoder(valueBytes),
		Timestamp: msg.Timestamp,
	}

	// Set key if provided
	if msg.Key != "" {
		saramaMsg.Key = sarama.StringEncoder(msg.Key)
	}

	// Set partition if specified
	if msg.Partition >= 0 {
		saramaMsg.Partition = msg.Partition
	}

	// Set headers
	if len(msg.Headers) > 0 {
		headers := make([]sarama.RecordHeader, 0, len(msg.Headers))
		for k, v := range msg.Headers {
			headers = append(headers, sarama.RecordHeader{
				Key:   []byte(k),
				Value: []byte(v),
			})
		}
		saramaMsg.Headers = headers
	}

	// Send message
	partition, offset, err = p.producer.SendMessage(saramaMsg)
	if err != nil {
		logger.Global().ErrorContext(ctx, "Failed to send Kafka message",
			zap.String("topic", msg.Topic),
			zap.String("key", msg.Key),
			zap.Error(err),
		)
		return 0, 0, fmt.Errorf("failed to send message: %w", err)
	}

	logger.Global().DebugContext(ctx, "Message sent successfully",
		zap.String("topic", msg.Topic),
		zap.String("key", msg.Key),
		zap.Int32("partition", partition),
		zap.Int64("offset", offset),
	)

	return partition, offset, nil
}

// SendJSON is a convenience method to send JSON messages
func (p *Producer) SendJSON(ctx context.Context, topic, key string, value any) (partition int32, offset int64, err error) {
	msg := &Message{
		Topic:     topic,
		Key:       key,
		Value:     value,
		Timestamp: time.Now(),
	}
	return p.Send(ctx, msg)
}

// SendString is a convenience method to send string messages
func (p *Producer) SendString(ctx context.Context, topic, key, value string) (partition int32, offset int64, err error) {
	msg := &Message{
		Topic:     topic,
		Key:       key,
		Value:     value,
		Timestamp: time.Now(),
	}
	return p.Send(ctx, msg)
}

// SendBytes is a convenience method to send byte messages
func (p *Producer) SendBytes(ctx context.Context, topic, key string, value []byte) (partition int32, offset int64, err error) {
	msg := &Message{
		Topic:     topic,
		Key:       key,
		Value:     value,
		Timestamp: time.Now(),
	}
	return p.Send(ctx, msg)
}

// Close closes the producer
func (p *Producer) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return nil
	}

	p.closed = true
	return p.producer.Close()
}

// IsConnected checks if the producer is connected
func (p *Producer) IsConnected() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return !p.closed
}

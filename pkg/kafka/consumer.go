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

// MessageHandler is a function type for handling consumed messages
type MessageHandler func(ctx context.Context, msg *ConsumedMessage) error

// ConsumedMessage represents a message consumed from Kafka
type ConsumedMessage struct {
	Topic     string
	Partition int32
	Offset    int64
	Key       string
	Value     []byte
	Headers   map[string]string
	Timestamp time.Time
}

// ValueAsString returns the message value as string
func (m *ConsumedMessage) ValueAsString() string {
	return string(m.Value)
}

// ValueAsJSON unmarshals the message value into the provided interface
func (m *ConsumedMessage) ValueAsJSON(v any) error {
	return json.Unmarshal(m.Value, v)
}

// Consumer wraps the Sarama consumer group
type Consumer struct {
	consumerGroup sarama.ConsumerGroup
	config        *Config
	handler       MessageHandler
	topics        []string
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	mu            sync.RWMutex
	closed        bool
}

// ConsumerGroupHandler implements sarama.ConsumerGroupHandler
type ConsumerGroupHandler struct {
	consumer *Consumer
}

// Setup is called at the beginning of a new session, before ConsumeClaim
func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is called at the end of a session, once all ConsumeClaim goroutines have exited
func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim processes messages from a partition
func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			// Convert Sarama message to our message format
			headers := make(map[string]string)
			for _, header := range message.Headers {
				headers[string(header.Key)] = string(header.Value)
			}

			msg := &ConsumedMessage{
				Topic:     message.Topic,
				Partition: message.Partition,
				Offset:    message.Offset,
				Key:       string(message.Key),
				Value:     message.Value,
				Headers:   headers,
				Timestamp: message.Timestamp,
			}

			// Call the message handler
			if err := h.consumer.handler(h.consumer.ctx, msg); err != nil {
				logger.Global().ErrorContext(h.consumer.ctx, "Message handler failed",
					zap.String("topic", msg.Topic),
					zap.Int32("partition", msg.Partition),
					zap.Int64("offset", msg.Offset),
					zap.String("key", msg.Key),
					zap.Error(err),
				)
				// Depending on your error handling strategy, you might want to:
				// 1. Continue processing (current behavior)
				// 2. Return the error to stop processing
				// 3. Send to a dead letter queue
			} else {
				// Mark message as processed
				session.MarkMessage(message, "")
			}

		case <-h.consumer.ctx.Done():
			return nil
		}
	}
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(config *Config, groupID string, topics []string, handler MessageHandler) (*Consumer, error) {
	if config == nil {
		config = DefaultConfig()
	}

	if handler == nil {
		return nil, fmt.Errorf("message handler cannot be nil")
	}

	saramaConfig := config.ToSaramaConfig()
	consumerGroup, err := sarama.NewConsumerGroup(config.Brokers, groupID, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	consumer := &Consumer{
		consumerGroup: consumerGroup,
		config:        config,
		handler:       handler,
		topics:        topics,
		ctx:           ctx,
		cancel:        cancel,
	}

	return consumer, nil
}

// Start begins consuming messages
func (c *Consumer) Start(ctx context.Context) error {
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return fmt.Errorf("consumer is closed")
	}
	c.mu.Unlock()

	// Handle consumer group errors
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for err := range c.consumerGroup.Errors() {
			logger.Global().ErrorContext(ctx, "Consumer group error", zap.Error(err))
		}
	}()
	// Start consuming
	handler := &ConsumerGroupHandler{consumer: c}
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			select {
			case <-c.ctx.Done():
				return
			default:
				if err := c.consumerGroup.Consume(c.ctx, c.topics, handler); err != nil {
					logger.Global().ErrorContext(ctx, "Consumer group consume error", zap.Error(err))
					return
				}
			}
		}
	}()

	logger.Global().InfoContext(ctx, "Kafka consumer started",
		zap.Strings("topics", c.topics),
		zap.Strings("brokers", c.config.Brokers),
	)

	return nil
}

// Stop stops the consumer
func (c *Consumer) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	c.closed = true
	c.cancel()
	c.wg.Wait()

	return c.consumerGroup.Close()
}

// IsRunning checks if the consumer is running
func (c *Consumer) IsRunning() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return !c.closed
}

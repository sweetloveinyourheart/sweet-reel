package kafka

import (
	"context"
	"fmt"
	"sync"
)

// Client provides a unified interface for Kafka operations
type Client struct {
	config   *Config
	producer *Producer
	mu       sync.RWMutex
	closed   bool
}

// NewClient creates a new Kafka client
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		config = DefaultConfig()
	}

	client := &Client{
		config: config,
	}

	return client, nil
}

// GetProducer returns a producer instance (creates if not exists)
func (c *Client) GetProducer() (*Producer, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil, fmt.Errorf("client is closed")
	}

	if c.producer == nil {
		producer, err := NewProducer(c.config)
		if err != nil {
			return nil, err
		}
		c.producer = producer
	}

	return c.producer, nil
}

// CreateConsumer creates a new consumer instance
func (c *Client) CreateConsumer(groupID string, topics []string, handler MessageHandler) (*Consumer, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return nil, fmt.Errorf("client is closed")
	}

	return NewConsumer(c.config, groupID, topics, handler)
}

// SendMessage sends a message using the internal producer
func (c *Client) SendMessage(ctx context.Context, msg *Message) (partition int32, offset int64, err error) {
	producer, err := c.GetProducer()
	if err != nil {
		return 0, 0, err
	}

	return producer.Send(ctx, msg)
}

// SendJSON sends a JSON message using the internal producer
func (c *Client) SendJSON(ctx context.Context, topic, key string, value any) (partition int32, offset int64, err error) {
	producer, err := c.GetProducer()
	if err != nil {
		return 0, 0, err
	}

	return producer.SendJSON(ctx, topic, key, value)
}

// SendString sends a string message using the internal producer
func (c *Client) SendString(ctx context.Context, topic, key, value string) (partition int32, offset int64, err error) {
	producer, err := c.GetProducer()
	if err != nil {
		return 0, 0, err
	}

	return producer.SendString(ctx, topic, key, value)
}

// SendBytes sends a byte message using the internal producer
func (c *Client) SendBytes(ctx context.Context, topic, key string, value []byte) (partition int32, offset int64, err error) {
	producer, err := c.GetProducer()
	if err != nil {
		return 0, 0, err
	}

	return producer.SendBytes(ctx, topic, key, value)
}

// Close closes the client and all associated resources
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	c.closed = true

	if c.producer != nil {
		if err := c.producer.Close(); err != nil {
			return err
		}
	}

	return nil
}

// IsConnected checks if the client is connected
func (c *Client) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return !c.closed
}

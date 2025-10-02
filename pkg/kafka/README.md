# Kafka Package

A reusable Kafka client package built on top of the [IBM Sarama](https://github.com/IBM/sarama) library, providing a simplified and opinionated interface for Kafka operations.

## Features

- **Easy-to-use API**: Simplified interfaces for common Kafka operations
- **Production-ready**: Built on the mature and battle-tested Sarama library
- **Type-safe**: Strong typing for messages and configurations
- **Flexible**: Support for JSON, string, and binary message formats
- **Observability**: Integrated logging with zap
- **Configurable**: Comprehensive configuration options with sensible defaults
- **Security**: Support for SASL and TLS authentication

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "log"
    
    "github.com/sweetloveinyourheart/sweet-reel/pkg/kafka"
)

func main() {
    // Create a client with default configuration
    client, err := kafka.NewClient(kafka.DefaultConfig())
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()
    
    ctx := context.Background()
    
    // Send a JSON message
    _, _, err = client.SendJSON(ctx, "user-events", "user123", map[string]any{
        "user_id": "user123",
        "action":  "login",
        "timestamp": time.Now(),
    })
    if err != nil {
        log.Printf("Failed to send message: %v", err)
    }
}
```

### Producer Example

```go
package main

import (
    "context"
    "time"
    
    "github.com/sweetloveinyourheart/sweet-reel/pkg/kafka"
)

func main() {
    // Create producer with custom config
    config := kafka.DefaultConfig()
    config.Brokers = []string{"localhost:9092"}
    
    producer, err := kafka.NewProducer(config)
    if err != nil {
        log.Fatal(err)
    }
    defer producer.Close()
    
    // Send different types of messages
    ctx := context.Background()
    
    // String message
    producer.SendString(ctx, "notifications", "notif1", "Hello World!")
    
    // JSON message
    producer.SendJSON(ctx, "events", "event1", map[string]any{
        "type": "user_action",
        "data": "some data",
    })
    
    // Custom message with headers
    message := &kafka.Message{
        Topic:     "custom-topic",
        Key:       "key123",
        Value:     []byte("binary data"),
        Headers:   map[string]string{"content-type": "application/octet-stream"},
        Timestamp: time.Now(),
    }
    producer.Send(ctx, message)
}
```

### Consumer Example

```go
package main

import (
    "context"
    "log"
    
    "github.com/sweetloveinyourheart/sweet-reel/pkg/kafka"
)

func messageHandler(ctx context.Context, msg *kafka.ConsumedMessage) error {
    log.Printf("Received: Topic=%s, Key=%s, Value=%s", 
        msg.Topic, msg.Key, msg.ValueAsString())
    
    // Handle JSON messages
    if msg.Topic == "user-events" {
        var event map[string]any
        if err := msg.ValueAsJSON(&event); err == nil {
            log.Printf("User event: %+v", event)
        }
    }
    
    return nil // Return error if processing fails
}

func main() {
    config := kafka.DefaultConfig()
    
    consumer, err := kafka.NewConsumer(
        config,
        "my-consumer-group",
        []string{"user-events", "notifications"},
        messageHandler,
    )
    if err != nil {
        log.Fatal(err)
    }
    defer consumer.Stop()
    
    // Start consuming
    if err := consumer.Start(context.Background()); err != nil {
        log.Fatal(err)
    }
    
    // Keep running
    select {}
}
```

## Configuration

### Default Configuration

The package provides sensible defaults for production use:

```go
config := kafka.DefaultConfig()
// config.Brokers = []string{"localhost:9092"}
// config.RequiredAcks = sarama.WaitForAll
// config.Compression = sarama.CompressionSnappy
// config.IdempotentWrites = true
// ... and more
```

### Custom Configuration

```go
config := &kafka.Config{
    Brokers:          []string{"broker1:9092", "broker2:9092"},
    RetryMax:         5,
    RetryBackoff:     200 * time.Millisecond,
    RequiredAcks:     sarama.WaitForAll,
    Compression:      sarama.CompressionSnappy,
    FlushFrequency:   1 * time.Second,
    FlushMessages:    100,
    FlushBytes:       16384,
    IdempotentWrites: true,
    SecurityProtocol: "SASL_SSL",
    SASLMechanism:    "SCRAM-SHA-256",
    SASLUsername:     "username",
    SASLPassword:     "password",
    TLSEnabled:       true,
}
```

## Message Types

### ConsumedMessage

```go
type ConsumedMessage struct {
    Topic     string
    Partition int32
    Offset    int64
    Key       string
    Value     []byte
    Headers   map[string]string
    Timestamp time.Time
}

// Helper methods
func (m *ConsumedMessage) ValueAsString() string
func (m *ConsumedMessage) ValueAsJSON(v any) error
```

### Message (for producing)

```go
type Message struct {
    Topic     string
    Key       string
    Value     any // Can be []byte, string, or any JSON-serializable type
    Headers   map[string]string
    Partition int32
    Timestamp time.Time
}
```

## Error Handling

The package provides comprehensive error handling:

- Connection errors are returned immediately
- Message sending errors include detailed context
- Consumer errors are logged and can be handled in the message handler
- All operations support context cancellation

## Best Practices

1. **Use contexts**: Always pass contexts for cancellation support
2. **Handle errors**: Check return values and handle errors appropriately
3. **Resource cleanup**: Always close clients, producers, and consumers
4. **Message keys**: Use meaningful keys for proper partitioning
5. **Batch processing**: Configure appropriate flush settings for throughput
6. **Monitoring**: Monitor the logs for production deployments

## Testing

Run the package tests:

```bash
go test ./pkg/kafka/
```

For integration tests with a real Kafka instance, ensure Kafka is running on `localhost:9092`.

## Dependencies

- [IBM Sarama](https://github.com/IBM/sarama) - Pure Go client library for Apache Kafka
- [Zap](https://github.com/uber-go/zap) - Fast, structured logging

## Contributing

When contributing to this package:

1. Maintain backward compatibility
2. Add tests for new features
3. Update this README for new functionality
4. Follow Go best practices and conventions
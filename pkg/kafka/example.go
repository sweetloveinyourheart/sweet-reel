package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

// ExampleUsage demonstrates how to use the Kafka package
func ExampleUsage() {
	// Create a custom config
	config := &Config{
		Brokers:          []string{"localhost:9092"},
		RetryMax:         3,
		RetryBackoff:     100 * time.Millisecond,
		RequiredAcks:     sarama.WaitForAll,
		Compression:      sarama.CompressionSnappy,
		FlushFrequency:   500 * time.Millisecond,
		FlushMessages:    100,
		FlushBytes:       16384,
		IdempotentWrites: true,
		SecurityProtocol: "PLAINTEXT",
	}

	// Create a client
	client, err := NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create Kafka client: %v", err)
	}
	defer client.Close()

	// Example 1: Send messages using client
	ctx := context.Background()

	// Send JSON message
	partition, offset, err := client.SendJSON(ctx, "user-events", "user123", map[string]any{
		"user_id": "user123",
		"action":  "login",
		"time":    time.Now(),
	})
	if err != nil {
		log.Printf("Failed to send JSON message: %v", err)
	} else {
		log.Printf("JSON message sent to partition %d, offset %d", partition, offset)
	}

	// Send string message
	partition, offset, err = client.SendString(ctx, "notifications", "notif456", "Hello, World!")
	if err != nil {
		log.Printf("Failed to send string message: %v", err)
	} else {
		log.Printf("String message sent to partition %d, offset %d", partition, offset)
	}

	// Example 2: Create and use producer directly
	producer, err := NewProducer(config)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	// Send a custom message
	message := &Message{
		Topic:     "video-processing",
		Key:       "video123",
		Value:     []byte("video processing data"),
		Headers:   map[string]string{"content-type": "application/octet-stream"},
		Timestamp: time.Now(),
	}

	partition, offset, err = producer.Send(ctx, message)
	if err != nil {
		log.Printf("Failed to send custom message: %v", err)
	} else {
		log.Printf("Custom message sent to partition %d, offset %d", partition, offset)
	}

	// Example 3: Create and start consumer
	messageHandler := func(ctx context.Context, msg *ConsumedMessage) error {
		log.Printf("Received message: Topic=%s, Key=%s, Value=%s, Offset=%d",
			msg.Topic, msg.Key, msg.ValueAsString(), msg.Offset)

		// Process the message here
		// You can use msg.ValueAsJSON() to unmarshal JSON messages
		var data map[string]any
		if err := msg.ValueAsJSON(&data); err == nil {
			log.Printf("JSON data: %+v", data)
		}

		return nil // Return error if processing fails
	}

	consumer, err := NewConsumer(config, "my-consumer-group", []string{"user-events", "notifications"}, messageHandler)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer func() {
		if err := consumer.Stop(); err != nil {
			log.Printf("Failed to stop consumer: %v", err)
		}
	}()

	// Start consuming (this will run in background)
	if err := consumer.Start(ctx); err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}

	// Keep the consumer running for a while
	time.Sleep(30 * time.Second)
}

// ExampleClientUsage shows simplified client usage
func ExampleClientUsage() {
	// Use default config
	client, err := NewClient(DefaultConfig())
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Simple message sending
	_, _, err = client.SendString(ctx, "my-topic", "my-key", "Hello Kafka!")
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}

	// Send structured data
	userData := struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		ID:   "123",
		Name: "John Doe",
		Age:  30,
	}

	_, _, err = client.SendJSON(ctx, "users", userData.ID, userData)
	if err != nil {
		log.Printf("Failed to send user data: %v", err)
	}

	fmt.Println("Messages sent successfully!")
}

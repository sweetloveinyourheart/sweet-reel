package kafka

import (
	"testing"
	"time"

	"github.com/IBM/sarama"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if len(config.Brokers) == 0 {
		t.Error("Expected default brokers to be set")
	}

	if config.Brokers[0] != "localhost:9092" {
		t.Errorf("Expected default broker to be localhost:9092, got %s", config.Brokers[0])
	}

	if config.RequiredAcks != sarama.WaitForAll {
		t.Error("Expected RequiredAcks to be WaitForAll")
	}
}

func TestSaramaConfigConversion(t *testing.T) {
	config := DefaultConfig()
	saramaConfig := config.ToSaramaConfig()

	if saramaConfig.Producer.RequiredAcks != sarama.WaitForAll {
		t.Error("Expected producer RequiredAcks to be WaitForAll")
	}

	if saramaConfig.Producer.Compression != sarama.CompressionSnappy {
		t.Error("Expected producer compression to be Snappy")
	}

	if !saramaConfig.Producer.Idempotent {
		t.Error("Expected producer to be idempotent")
	}
}

func TestNewClient(t *testing.T) {
	client, err := NewClient(nil) // Use default config
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	if !client.IsConnected() {
		t.Error("Expected client to be connected")
	}
}

func TestMessageMethods(t *testing.T) {
	msg := &ConsumedMessage{
		Topic:     "test-topic",
		Partition: 0,
		Offset:    123,
		Key:       "test-key",
		Value:     []byte(`{"test": "value"}`),
		Headers:   map[string]string{"content-type": "application/json"},
		Timestamp: time.Now(),
	}

	// Test ValueAsString
	if msg.ValueAsString() != `{"test": "value"}` {
		t.Errorf("Expected JSON string, got %s", msg.ValueAsString())
	}

	// Test ValueAsJSON
	var data map[string]any
	err := msg.ValueAsJSON(&data)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if data["test"] != "value" {
		t.Errorf("Expected 'value', got %v", data["test"])
	}
}

// Note: These tests don't require a running Kafka instance
// For integration tests with actual Kafka, you would need to set up
// test containers or use embedded Kafka for testing

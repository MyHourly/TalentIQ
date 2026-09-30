package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

// KafkaProducer is responsible for publishing events to Kafka.
type KafkaProducer struct {
	writer *kafka.Writer
	logger *slog.Logger
}

// NewKafkaProducer creates a Kafka producer.
//
// brokers contains one or more Kafka broker addresses.
// topic is the Kafka topic where events will be published.
func NewKafkaProducer(
	brokers []string,
	topic string,
	logger *slog.Logger,
) *KafkaProducer {

	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaProducer{
		writer: writer,
		logger: logger,
	}
}

// Publish publishes an event to Kafka.
//
// The event is converted to JSON before it is sent.
func (p *KafkaProducer) Publish(
	ctx context.Context,
	key string,
	event interface{},
) error {

	// Convert the event into JSON.
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal kafka event: %w", err)
	}

	// Publish the event to Kafka.
	err = p.writer.WriteMessages(
		ctx,
		kafka.Message{
			Key:   []byte(key),
			Value: payload,
		},
	)

	if err != nil {
		p.logger.Error(
			"failed to publish kafka event",
			"error", err,
			"key", key,
		)

		return fmt.Errorf("publish kafka event: %w", err)
	}

	p.logger.Info(
		"kafka event published",
		"key", key,
	)

	return nil
}

// Close closes the Kafka writer.
//
// This should be called when the application shuts down.
func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}

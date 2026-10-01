package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

// KafkaConsumer consumes events from Kafka.
//
// NOTE:
// For local development, Kafka broker and topic are passed directly.
// Later, these values should come from the centralized TalentIQ
// Kafka configuration instead of being configured inside this service.
type KafkaConsumer struct {
	reader *kafka.Reader
	logger *slog.Logger
}

// NewKafkaConsumer creates a Kafka consumer.
//
// brokers contains the local Kafka broker addresses.
// topic is the Kafka topic to consume from.
// groupID identifies this service's Kafka consumer group.
func NewKafkaConsumer(
	brokers []string,
	topic string,
	groupID string,
	logger *slog.Logger,
) *KafkaConsumer {

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,

		// Start consuming from the earliest available message
		// when the consumer group has no previous offset.
		StartOffset: kafka.FirstOffset,
	})

	return &KafkaConsumer{
		reader: reader,
		logger: logger,
	}
}

// Consume continuously reads events from Kafka.
//
// NOTE:
// This method currently only logs received events.
// Later, actual event handlers will process different
// event types from other TalentIQ microservices.
func (c *KafkaConsumer) Consume(ctx context.Context) error {

	c.logger.Info(
		"Kafka consumer started",
	)

	for {

		// Read the next Kafka message.
		message, err := c.reader.ReadMessage(ctx)

		if err != nil {

			// Context cancellation means the application
			// is shutting down normally.
			if ctx.Err() != nil {
				return nil
			}

			return fmt.Errorf(
				"read kafka message: %w",
				err,
			)
		}

		// Convert the message into our event structure.
		var event TalentProfileEvent

		if err := json.Unmarshal(
			message.Value,
			&event,
		); err != nil {

			c.logger.Error(
				"failed to decode kafka event",
				"error", err,
			)

			// Do not stop the consumer because of one
			// malformed message.
			continue
		}

		// Log the received event.
		c.logger.Info(
			"Kafka event received",
			"event_id", event.EventID,
			"event_type", event.EventType,
			"talent_id", event.TalentID,
			"employee_code", event.EmployeeCode,
			"occurred_at", event.OccurredAt,
		)
	}
}

// Close closes the Kafka consumer.
//
// This should be called when the application shuts down.
func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}
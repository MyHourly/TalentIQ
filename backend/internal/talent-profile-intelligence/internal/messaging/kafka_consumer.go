package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

// KafkaConsumer consumes talent profile events from Kafka.
type KafkaConsumer struct {
	reader     *kafka.Reader
	logger     *slog.Logger
	dispatcher *EventDispatcher
}

// NewKafkaConsumer creates a Kafka consumer.
func NewKafkaConsumer(
	brokers []string,
	topic string,
	groupID string,
	logger *slog.Logger,
	dispatcher *EventDispatcher,
) *KafkaConsumer {

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		Topic:       topic,
		GroupID:     groupID,
		StartOffset: kafka.FirstOffset,
	})

	return &KafkaConsumer{
		reader:     reader,
		logger:     logger,
		dispatcher: dispatcher,
	}
}

// Consume continuously reads and processes Kafka events.
func (c *KafkaConsumer) Consume(ctx context.Context) error {

	c.logger.Info(
		"Kafka consumer started",
	)

	for {
		message, err := c.reader.ReadMessage(ctx)
		if err != nil {

			// Context cancellation is a normal shutdown.
			if ctx.Err() != nil {
				return nil
			}

			return fmt.Errorf(
				"read kafka message: %w",
				err,
			)
		}

		var event TalentProfileEvent

		if err := json.Unmarshal(
			message.Value,
			&event,
		); err != nil {

			c.logger.Error(
				"failed to decode kafka event",
				"error", err,
			)

			// Continue consuming the next message.
			continue
		}

		c.logger.Info(
			"Kafka event received",
			"event_id", event.EventID,
			"event_type", event.EventType,
			"talent_id", event.TalentID,
			"employee_code", event.EmployeeCode,
			"occurred_at", event.OccurredAt,
		)

		// Dispatch the event to the correct handler.
		if err := c.dispatcher.Dispatch(
			ctx,
			event,
		); err != nil {

			c.logger.Error(
				"failed to process kafka event",
				"error", err,
				"event_id", event.EventID,
				"event_type", event.EventType,
				"talent_id", event.TalentID,
				"employee_code", event.EmployeeCode,
			)

			// Continue consuming other events.
			continue
		}

		c.logger.Info(
			"Kafka event processed successfully",
			"event_id", event.EventID,
			"event_type", event.EventType,
			"talent_id", event.TalentID,
		)
	}
}

// Close closes the Kafka reader.
func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}

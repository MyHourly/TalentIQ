package messaging

import "context"

// EventPublisher defines the functionality required by the Service layer
// to publish application events.
type EventPublisher interface {
	Publish(
		ctx context.Context,
		key string,
		event interface{},
	) error
}

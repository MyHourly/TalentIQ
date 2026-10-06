package messaging

import (
	"context"
	"fmt"
	"log/slog"
)

// TalentProfileEventHandler handles a specific talent profile event.
type TalentProfileEventHandler interface {
	Handle(
		ctx context.Context,
		event TalentProfileEvent,
	) error
}

// EventDispatcher routes events to their registered handlers.
type EventDispatcher struct {
	handlers map[string]TalentProfileEventHandler
	logger   *slog.Logger
}

// NewEventDispatcher creates a new event dispatcher.
func NewEventDispatcher(logger *slog.Logger) *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string]TalentProfileEventHandler),
		logger:   logger,
	}
}

// RegisterHandler registers a handler for an event type.
func (d *EventDispatcher) RegisterHandler(
	eventType string,
	handler TalentProfileEventHandler,
) {
	d.handlers[eventType] = handler

	d.logger.Info(
		"event handler registered",
		"event_type", eventType,
	)
}

// Dispatch sends an event to the correct handler.
func (d *EventDispatcher) Dispatch(
	ctx context.Context,
	event TalentProfileEvent,
) error {

	handler, exists := d.handlers[event.EventType]
	if !exists {
		return fmt.Errorf(
			"no handler registered for event type: %s",
			event.EventType,
		)
	}

	d.logger.Info(
		"dispatching talent profile event",
		"event_type", event.EventType,
		"event_id", event.EventID,
		"talent_id", event.TalentID,
	)

	if err := handler.Handle(ctx, event); err != nil {
		return fmt.Errorf(
			"handle event %s: %w",
			event.EventType,
			err,
		)
	}

	return nil
}

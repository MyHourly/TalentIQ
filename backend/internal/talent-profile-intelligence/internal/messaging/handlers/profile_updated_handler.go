package handlers

import (
	"context"
	"log/slog"

	"talentiq/talent-profile-intelligence/internal/messaging"
	"talentiq/talent-profile-intelligence/internal/service"
)

// ProfileUpdatedHandler handles candidate.profile.updated events.
//
// The handler receives the event from the dispatcher and passes it
// to the event service for business processing.
type ProfileUpdatedHandler struct {
	eventService *service.TalentProfileEventService
	logger       *slog.Logger
}

// NewProfileUpdatedHandler creates a new profile-updated handler.
func NewProfileUpdatedHandler(
	eventService *service.TalentProfileEventService,
	logger *slog.Logger,
) *ProfileUpdatedHandler {

	return &ProfileUpdatedHandler{
		eventService: eventService,
		logger:       logger,
	}
}

// Handle processes a candidate.profile.updated event.
func (h *ProfileUpdatedHandler) Handle(
	ctx context.Context,
	event messaging.TalentProfileEvent,
) error {

	h.logger.Info(
		"handling profile updated event",
		"event_id", event.EventID,
		"talent_id", event.TalentID,
	)

	// Pass the event to the business processing service.
	//
	// Any error returned by the service is passed back to
	// the EventDispatcher and then to the Kafka consumer.
	return h.eventService.HandleProfileUpdated(
		ctx,
		event,
	)
}

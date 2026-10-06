package handlers

import (
	"context"
	"log/slog"

	"talentiq/talent-profile-intelligence/internal/messaging"
	"talentiq/talent-profile-intelligence/internal/service"
)

// ProfileCreatedHandler handles candidate.profile.created events.
type ProfileCreatedHandler struct {
	eventService *service.TalentProfileEventService
	logger       *slog.Logger
}

// NewProfileCreatedHandler creates a new profile-created handler.
func NewProfileCreatedHandler(
	eventService *service.TalentProfileEventService,
	logger *slog.Logger,
) *ProfileCreatedHandler {

	return &ProfileCreatedHandler{
		eventService: eventService,
		logger:       logger,
	}
}

// Handle processes a candidate.profile.created event.
func (h *ProfileCreatedHandler) Handle(
	ctx context.Context,
	event messaging.TalentProfileEvent,
) error {

	h.logger.Info(
		"handling profile created event",
		"event_id", event.EventID,
		"talent_id", event.TalentID,
	)

	// Pass the event to the business processing service.
	return h.eventService.HandleProfileCreated(
		ctx,
		event,
	)
}

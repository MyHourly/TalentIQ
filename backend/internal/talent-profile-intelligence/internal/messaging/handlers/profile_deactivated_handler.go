package handlers

import (
	"context"
	"log/slog"

	"talentiq/talent-profile-intelligence/internal/messaging"
	"talentiq/talent-profile-intelligence/internal/service"
)

// ProfileDeactivatedHandler handles
// candidate.profile.deactivated events.
type ProfileDeactivatedHandler struct {
	eventService *service.TalentProfileEventService
	logger       *slog.Logger
}

// NewProfileDeactivatedHandler creates a new
// profile-deactivated handler.
func NewProfileDeactivatedHandler(
	eventService *service.TalentProfileEventService,
	logger *slog.Logger,
) *ProfileDeactivatedHandler {

	return &ProfileDeactivatedHandler{
		eventService: eventService,
		logger:       logger,
	}
}

// Handle processes a candidate.profile.deactivated event.
func (h *ProfileDeactivatedHandler) Handle(
	ctx context.Context,
	event messaging.TalentProfileEvent,
) error {

	h.logger.Info(
		"handling profile deactivated event",
		"event_id", event.EventID,
		"talent_id", event.TalentID,
	)

	// Pass the event to the business processing service.
	return h.eventService.HandleProfileDeactivated(
		ctx,
		event,
	)
}

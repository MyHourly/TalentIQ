package service

import (
	"context"
	"fmt"
	"log/slog"

	"talentiq/talent-profile-intelligence/internal/messaging"
	"talentiq/talent-profile-intelligence/internal/repository"
)

// TalentProfileEventService handles business processing
// triggered by talent profile events.
type TalentProfileEventService struct {
	repository repository.TalentProfileRepository
	logger     *slog.Logger
}

// NewTalentProfileEventService creates the event service.
func NewTalentProfileEventService(
	repository repository.TalentProfileRepository,
	logger *slog.Logger,
) *TalentProfileEventService {

	return &TalentProfileEventService{
		repository: repository,
		logger:     logger,
	}
}

// HandleProfileCreated processes a profile created event.
func (s *TalentProfileEventService) HandleProfileCreated(
	ctx context.Context,
	event messaging.TalentProfileEvent,
) error {

	profile, err := s.repository.GetByID(
		ctx,
		event.TalentID,
	)

	if err != nil {
		return fmt.Errorf(
			"profile created event: get talent profile: %w",
			err,
		)
	}

	if profile == nil {
		return fmt.Errorf(
			"profile created event: talent profile not found: %s",
			event.TalentID,
		)
	}

	s.logger.Info(
		"profile creation event processed",
		"event_id", event.EventID,
		"talent_id", profile.ID,
		"employee_code", profile.EmployeeCode,
	)

	return nil
}

// HandleProfileUpdated processes a profile updated event.
func (s *TalentProfileEventService) HandleProfileUpdated(
	ctx context.Context,
	event messaging.TalentProfileEvent,
) error {

	profile, err := s.repository.GetByID(
		ctx,
		event.TalentID,
	)

	if err != nil {
		return fmt.Errorf(
			"profile updated event: get talent profile: %w",
			err,
		)
	}

	if profile == nil {
		return fmt.Errorf(
			"profile updated event: talent profile not found: %s",
			event.TalentID,
		)
	}

	s.logger.Info(
		"profile update event processed",
		"event_id", event.EventID,
		"talent_id", profile.ID,
		"employee_code", profile.EmployeeCode,
	)

	return nil
}

// HandleProfileDeactivated processes a profile deactivated event.
func (s *TalentProfileEventService) HandleProfileDeactivated(
	ctx context.Context,
	event messaging.TalentProfileEvent,
) error {

	// A deactivated profile may already be soft deleted.
	profile, err := s.repository.GetByIDIncludingDeleted(
		ctx,
		event.TalentID,
	)

	if err != nil {
		return fmt.Errorf(
			"profile deactivated event: get talent profile: %w",
			err,
		)
	}

	if profile == nil {
		return fmt.Errorf(
			"profile deactivated event: talent profile not found: %s",
			event.TalentID,
		)
	}

	s.logger.Info(
		"profile deactivation event processed",
		"event_id", event.EventID,
		"talent_id", profile.ID,
		"employee_code", profile.EmployeeCode,
		"profile_status", profile.ProfileStatus,
	)

	return nil
}

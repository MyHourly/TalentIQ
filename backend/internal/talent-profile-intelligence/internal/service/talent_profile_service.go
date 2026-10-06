package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"talentiq/talent-profile-intelligence/internal/messaging"
	"talentiq/talent-profile-intelligence/internal/model"
	"talentiq/talent-profile-intelligence/internal/repository"
)

// TalentProfileService defines business operations for talent profiles.
type TalentProfileService interface {
	Create(ctx context.Context, profile *model.TalentProfile) (*model.TalentProfile, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.TalentProfile, error)
	List(ctx context.Context, limit, offset int) ([]*model.TalentProfile, int, error)
	Update(ctx context.Context, profile *model.TalentProfile) (*model.TalentProfile, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// talentProfileService contains business logic for talent profiles.
type talentProfileService struct {
	repository repository.TalentProfileRepository
	publisher  messaging.EventPublisher
	logger     *slog.Logger
}

// NewTalentProfileService creates a new talent profile service.
func NewTalentProfileService(
	repository repository.TalentProfileRepository,
	publisher messaging.EventPublisher,
	logger *slog.Logger,
) TalentProfileService {
	return &talentProfileService{
		repository: repository,
		publisher:  publisher,
		logger:     logger,
	}
}

// Create creates a new talent profile.
func (s *talentProfileService) Create(
	ctx context.Context,
	profile *model.TalentProfile,
) (*model.TalentProfile, error) {

	// Validate the profile.
	if profile == nil || profile.EmployeeCode == "" {
		return nil, fmt.Errorf(
			"invalid talent profile: %w",
			ErrInvalidTalentProfile,
		)
	}

	// Check whether employee code already exists.
	exists, err := s.repository.ExistsByEmployeeCode(
		ctx,
		profile.EmployeeCode,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"check employee code: %w",
			err,
		)
	}

	if exists {
		return nil, fmt.Errorf(
			"employee code already exists: %w",
			ErrEmployeeCodeExists,
		)
	}

	// New profiles are active by default.
	profile.ProfileStatus = "ACTIVE"

	createdProfile, err := s.repository.Create(ctx, profile)
	if err != nil {
		return nil, fmt.Errorf(
			"create talent profile: %w",
			err,
		)
	}

	s.logger.Info(
		"talent profile created",
		"talent_id", createdProfile.ID,
		"employee_code", createdProfile.EmployeeCode,
	)

	return createdProfile, nil
}

// GetByID retrieves an active talent profile by ID.
func (s *talentProfileService) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*model.TalentProfile, error) {

	// Validate the profile ID before accessing the repository.
	if id == uuid.Nil {
		return nil, fmt.Errorf(
			"invalid talent profile ID: %w",
			ErrInvalidTalentProfile,
		)
	}

	profile, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf(
			"get talent profile: %w",
			err,
		)
	}

	return profile, nil
}

// List retrieves talent profiles with pagination.
func (s *talentProfileService) List(
	ctx context.Context,
	limit,
	offset int,
) ([]*model.TalentProfile, int, error) {

	// Apply safe defaults for invalid pagination values.
	if limit <= 0 {
		limit = 10
	}

	if offset < 0 {
		offset = 0
	}

	profiles, total, err := s.repository.List(
		ctx,
		limit,
		offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf(
			"list talent profiles: %w",
			err,
		)
	}

	return profiles, total, nil
}

// Update updates an existing talent profile.
func (s *talentProfileService) Update(
	ctx context.Context,
	profile *model.TalentProfile,
) (*model.TalentProfile, error) {

	// Validate the profile.
	if profile == nil || profile.ID == uuid.Nil {
		return nil, fmt.Errorf(
			"invalid talent profile: %w",
			ErrInvalidTalentProfile,
		)
	}

	updatedProfile, err := s.repository.Update(ctx, profile)
	if err != nil {
		return nil, fmt.Errorf(
			"update talent profile: %w",
			err,
		)
	}

	s.logger.Info(
		"talent profile updated",
		"talent_id", updatedProfile.ID,
		"employee_code", updatedProfile.EmployeeCode,
	)

	// Create the profile updated event.
	event := messaging.TalentProfileEvent{
		EventID:      uuid.New(),
		EventType:    messaging.TalentProfileUpdated,
		TalentID:     updatedProfile.ID,
		EmployeeCode: updatedProfile.EmployeeCode,
		OccurredAt:   time.Now().UTC(),
	}

	// Publish the event to Kafka.
	if s.publisher != nil {
		if err := s.publisher.Publish(
			ctx,
			updatedProfile.ID.String(),
			event,
		); err != nil {

			// Database update already succeeded.
			// Log the Kafka failure without failing the API operation.
			s.logger.Error(
				"failed to publish talent profile updated event",
				"error", err,
				"talent_id", updatedProfile.ID,
				"event_id", event.EventID,
			)
		}
	}

	return updatedProfile, nil
}

// Delete soft-deletes a talent profile.
func (s *talentProfileService) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {

	// Validate the profile ID.
	if id == uuid.Nil {
		return fmt.Errorf(
			"invalid talent profile ID: %w",
			ErrInvalidTalentProfile,
		)
	}

	// Soft-delete the profile in PostgreSQL.
	if err := s.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf(
			"delete talent profile: %w",
			err,
		)
	}

	s.logger.Info(
		"talent profile deactivated",
		"talent_id", id,
	)

	// Retrieve the deleted profile so that we can include
	// the employee code in the Kafka event.
	profile, err := s.repository.GetByIDIncludingDeleted(
		ctx,
		id,
	)
	if err != nil {
		return fmt.Errorf(
			"get deleted talent profile for event: %w",
			err,
		)
	}

	if profile == nil {
		return fmt.Errorf(
			"deleted talent profile not found: %s",
			id,
		)
	}

	// Create the profile deactivated event.
	event := messaging.TalentProfileEvent{
		EventID:      uuid.New(),
		EventType:    messaging.TalentProfileDeactivated,
		TalentID:     profile.ID,
		EmployeeCode: profile.EmployeeCode,
		OccurredAt:   time.Now().UTC(),
	}

	// Publish the deactivation event to Kafka.
	if s.publisher != nil {
		if err := s.publisher.Publish(
			ctx,
			profile.ID.String(),
			event,
		); err != nil {

			// Database deletion already succeeded.
			// Log the Kafka failure without failing the API operation.
			s.logger.Error(
				"failed to publish talent profile deactivated event",
				"error", err,
				"talent_id", profile.ID,
				"event_id", event.EventID,
			)
		}
	}

	return nil
}

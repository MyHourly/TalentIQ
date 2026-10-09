package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"

	"talentiq/talent-profile-intelligence/internal/model"
	"talentiq/talent-profile-intelligence/internal/repository"
)

// Common service errors.
var (
	ErrInvalidTalentPreference = errors.New("invalid talent preference")
)

// TalentPreferenceService defines business operations
// for managing talent preferences.
type TalentPreferenceService interface {
	Create(
		ctx context.Context,
		preference *model.TalentPreference,
	) (*model.TalentPreference, error)

	GetByTalentID(
		ctx context.Context,
		talentID uuid.UUID,
	) (*model.TalentPreference, error)

	Update(
		ctx context.Context,
		preference *model.TalentPreference,
	) (*model.TalentPreference, error)

	UpdateAvailability(
		ctx context.Context,
		talentID uuid.UUID,
		availabilityStatus string,
	) (*model.TalentPreference, error)
}

// talentPreferenceService contains the business logic
// for the Talent Preference domain.
type talentPreferenceService struct {
	repository repository.TalentPreferenceRepository
	logger     *slog.Logger
}

// NewTalentPreferenceService creates a new Talent Preference service.
func NewTalentPreferenceService(
	repository repository.TalentPreferenceRepository,
	logger *slog.Logger,
) TalentPreferenceService {
	return &talentPreferenceService{
		repository: repository,
		logger:     logger,
	}
}

// Create creates preferences for a talent.
func (s *talentPreferenceService) Create(
	ctx context.Context,
	preference *model.TalentPreference,
) (*model.TalentPreference, error) {

	if preference == nil || preference.TalentID == uuid.Nil {
		return nil, ErrInvalidTalentPreference
	}

	// Check whether the talent already has preferences.
	// Each talent should have only one current preference record.
	_, err := s.repository.GetByTalentID(ctx, preference.TalentID)

	if err == nil {
		return nil, fmt.Errorf(
			"talent preference already exists for talent: %w",
			ErrInvalidTalentPreference,
		)
	}

	if !errors.Is(err, repository.ErrTalentPreferenceNotFound) {
		return nil, fmt.Errorf(
			"check existing talent preference: %w",
			err,
		)
	}

	created, err := s.repository.Create(ctx, preference)
	if err != nil {
		s.logger.Error(
			"failed to create talent preference",
			"talent_id", preference.TalentID,
			"error", err,
		)

		return nil, fmt.Errorf(
			"create talent preference: %w",
			err,
		)
	}

	s.logger.Info(
		"talent preference created",
		"talent_id", created.TalentID,
		"preference_id", created.ID,
	)

	return created, nil
}

// GetByTalentID returns the current preferences
// belonging to a talent.
func (s *talentPreferenceService) GetByTalentID(
	ctx context.Context,
	talentID uuid.UUID,
) (*model.TalentPreference, error) {

	if talentID == uuid.Nil {
		return nil, ErrInvalidTalentPreference
	}

	preference, err := s.repository.GetByTalentID(ctx, talentID)
	if err != nil {
		return nil, fmt.Errorf(
			"get talent preference: %w",
			err,
		)
	}

	return preference, nil
}

// Update modifies the existing preferences of a talent.
func (s *talentPreferenceService) Update(
	ctx context.Context,
	preference *model.TalentPreference,
) (*model.TalentPreference, error) {

	if preference == nil || preference.TalentID == uuid.Nil {
		return nil, ErrInvalidTalentPreference
	}

	updated, err := s.repository.Update(ctx, preference)
	if err != nil {
		s.logger.Error(
			"failed to update talent preference",
			"talent_id", preference.TalentID,
			"error", err,
		)

		return nil, fmt.Errorf(
			"update talent preference: %w",
			err,
		)
	}

	s.logger.Info(
		"talent preference updated",
		"talent_id", updated.TalentID,
		"preference_id", updated.ID,
	)

	return updated, nil
}

// UpdateAvailability changes only the availability status
// of an existing talent preference.
func (s *talentPreferenceService) UpdateAvailability(
	ctx context.Context,
	talentID uuid.UUID,
	availabilityStatus string,
) (*model.TalentPreference, error) {

	if talentID == uuid.Nil {
		return nil, ErrInvalidTalentPreference
	}

	// Availability status should not be empty.
	if strings.TrimSpace(availabilityStatus) == "" {
		return nil, ErrInvalidTalentPreference
	}

	updated, err := s.repository.UpdateAvailability(
		ctx,
		talentID,
		availabilityStatus,
	)
	if err != nil {
		s.logger.Error(
			"failed to update talent availability",
			"talent_id", talentID,
			"error", err,
		)

		return nil, fmt.Errorf(
			"update talent availability: %w",
			err,
		)
	}

	s.logger.Info(
		"talent availability updated",
		"talent_id", updated.TalentID,
		"availability_status", updated.AvailabilityStatus,
	)

	return updated, nil
}

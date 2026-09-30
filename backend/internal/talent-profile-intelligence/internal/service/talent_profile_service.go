package service

import (
	"context"

	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"

	"talentiq/talent-profile-intelligence/internal/model"
	"talentiq/talent-profile-intelligence/internal/repository"
)

// These errors represent business-level validation failures.
//
// The HTTP handler will later convert these errors into
// appropriate HTTP responses such as 400 or 409.

// TalentProfileService contains business operations related
// to Talent Profiles.
type TalentProfileService interface {

	// Creates a new talent profile after validation.
	Create(
		ctx context.Context,
		profile *model.TalentProfile,
	) (*model.TalentProfile, error)

	// Returns a single talent profile.
	GetByID(
		ctx context.Context,
		id uuid.UUID,
	) (*model.TalentProfile, error)

	// Returns a paginated list of talent profiles.
	List(
		ctx context.Context,
		limit int,
		offset int,
	) ([]*model.TalentProfile, int, error)

	// Updates an existing talent profile.
	Update(
		ctx context.Context,
		profile *model.TalentProfile,
	) (*model.TalentProfile, error)

	// Deactivates a talent profile.
	Delete(
		ctx context.Context,
		id uuid.UUID,
	) error
}

// talentProfileService is the concrete implementation
// of TalentProfileService.
type talentProfileService struct {
	repository repository.TalentProfileRepository
	logger     *slog.Logger
}

// NewTalentProfileService creates a new Talent Profile service.
func NewTalentProfileService(
	repository repository.TalentProfileRepository,
	logger *slog.Logger,
) TalentProfileService {

	return &talentProfileService{
		repository: repository,
		logger:     logger,
	}
}

// Create validates and creates a new talent profile.
func (s *talentProfileService) Create(
	ctx context.Context,
	profile *model.TalentProfile,
) (*model.TalentProfile, error) {

	// Validate the incoming profile before touching the database.
	if err := validateTalentProfile(profile); err != nil {
		s.logger.Warn(
			"talent profile validation failed",
			"error", err,
		)

		return nil, err
	}

	// Check whether another active profile already uses
	// the same employee code.
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
		s.logger.Warn(
			"employee code already exists",
			"employee_code", profile.EmployeeCode,
		)

		return nil, ErrEmployeeCodeExists
	}

	// New profiles are active by default.
	profile.ProfileStatus = "ACTIVE"

	created, err := s.repository.Create(
		ctx,
		profile,
	)

	if err != nil {
		s.logger.Error(
			"failed to create talent profile",
			"error", err,
		)

		return nil, fmt.Errorf(
			"create talent profile: %w",
			err,
		)
	}

	s.logger.Info(
		"talent profile created successfully",
		"profile_id", created.ID,
	)

	return created, nil
}

// GetByID returns one talent profile.
func (s *talentProfileService) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*model.TalentProfile, error) {

	// A zero UUID is not a valid profile identifier.
	if id == uuid.Nil {
		return nil, ErrInvalidTalentProfile
	}

	profile, err := s.repository.GetByID(
		ctx,
		id,
	)

	if err != nil {
		return nil, err
	}

	return profile, nil
}

// List returns talent profiles using pagination.
func (s *talentProfileService) List(
	ctx context.Context,
	limit int,
	offset int,
) ([]*model.TalentProfile, int, error) {

	// Protect the database from unreasonable pagination values.
	if limit <= 0 {
		limit = 20
	}

	// Prevent negative offsets.
	if offset < 0 {
		offset = 0
	}

	// Put a reasonable upper limit on a single request.
	if limit > 100 {
		limit = 100
	}

	profiles, total, err := s.repository.List(
		ctx,
		limit,
		offset,
	)

	if err != nil {
		return nil, 0, err
	}

	return profiles, total, nil
}

// Update validates and updates a talent profile.
func (s *talentProfileService) Update(
	ctx context.Context,
	profile *model.TalentProfile,
) (*model.TalentProfile, error) {

	if profile == nil {
		return nil, ErrInvalidTalentProfile
	}

	if profile.ID == uuid.Nil {
		return nil, ErrInvalidTalentProfile
	}

	// Validate the profile before updating it.
	if err := validateTalentProfile(profile); err != nil {
		return nil, err
	}

	updated, err := s.repository.Update(
		ctx,
		profile,
	)

	if err != nil {
		return nil, err
	}

	s.logger.Info(
		"talent profile updated successfully",
		"profile_id", updated.ID,
	)

	return updated, nil
}

// Delete deactivates a talent profile.
func (s *talentProfileService) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {

	if id == uuid.Nil {
		return ErrInvalidTalentProfile
	}

	if err := s.repository.Delete(
		ctx,
		id,
	); err != nil {
		return err
	}

	s.logger.Info(
		"talent profile deactivated successfully",
		"profile_id", id,
	)

	return nil
}

// validateTalentProfile contains basic business validation.
//
// More detailed validation can later be moved to DTO validation
// when the HTTP layer is implemented.
func validateTalentProfile(
	profile *model.TalentProfile,
) error {

	if profile == nil {
		return ErrInvalidTalentProfile
	}

	// Employee code is required.
	if strings.TrimSpace(profile.EmployeeCode) == "" {
		return fmt.Errorf(
			"%w: employee code is required",
			ErrInvalidTalentProfile,
		)
	}

	// First name is required.
	if strings.TrimSpace(profile.FirstName) == "" {
		return fmt.Errorf(
			"%w: first name is required",
			ErrInvalidTalentProfile,
		)
	}

	// Last name is required.
	if strings.TrimSpace(profile.LastName) == "" {
		return fmt.Errorf(
			"%w: last name is required",
			ErrInvalidTalentProfile,
		)
	}

	// Email is required.
	if strings.TrimSpace(profile.Email) == "" {
		return fmt.Errorf(
			"%w: email is required",
			ErrInvalidTalentProfile,
		)
	}

	// Designation is required.
	if strings.TrimSpace(profile.Designation) == "" {
		return fmt.Errorf(
			"%w: designation is required",
			ErrInvalidTalentProfile,
		)
	}

	// Department is required.
	if strings.TrimSpace(profile.Department) == "" {
		return fmt.Errorf(
			"%w: department is required",
			ErrInvalidTalentProfile,
		)
	}

	// Experience cannot be negative.
	if profile.TotalExperienceYears < 0 {
		return fmt.Errorf(
			"%w: total experience cannot be negative",
			ErrInvalidTalentProfile,
		)
	}

	return nil
}

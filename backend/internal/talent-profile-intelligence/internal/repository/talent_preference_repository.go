package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"talentiq/talent-profile-intelligence/internal/model"
)

// ErrTalentPreferenceNotFound is returned when a talent has no preference record.
var ErrTalentPreferenceNotFound = errors.New("talent preference not found")

// TalentPreferenceRepository defines database operations
// supported by the Talent Preference domain.
type TalentPreferenceRepository interface {
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

// talentPreferenceRepository is the PostgreSQL implementation
// of TalentPreferenceRepository.
type talentPreferenceRepository struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

// NewTalentPreferenceRepository creates a PostgreSQL-backed
// Talent Preference repository.
func NewTalentPreferenceRepository(
	db *pgxpool.Pool,
	logger *slog.Logger,
) TalentPreferenceRepository {

	return &talentPreferenceRepository{
		db:     db,
		logger: logger,
	}
}

// Create inserts a new talent preference into PostgreSQL.
func (r *talentPreferenceRepository) Create(
	ctx context.Context,
	preference *model.TalentPreference,
) (*model.TalentPreference, error) {

	const query = `
		INSERT INTO talent_preferences (
			talent_id,
			preferred_role,
			preferred_location,
			availability_status
		)
		VALUES ($1, $2, $3, $4)
		RETURNING
			id,
			talent_id,
			preferred_role,
			preferred_location,
			availability_status,
			created_at,
			updated_at
	`

	var created model.TalentPreference

	err := r.db.QueryRow(
		ctx,
		query,
		preference.TalentID,
		preference.PreferredRole,
		preference.PreferredLocation,
		preference.AvailabilityStatus,
	).Scan(
		&created.ID,
		&created.TalentID,
		&created.PreferredRole,
		&created.PreferredLocation,
		&created.AvailabilityStatus,
		&created.CreatedAt,
		&created.UpdatedAt,
	)

	if err != nil {
		r.logger.Error(
			"failed to create talent preference",
			"talent_id", preference.TalentID,
			"error", err,
		)

		return nil, fmt.Errorf(
			"create talent preference: %w",
			err,
		)
	}

	r.logger.Info(
		"talent preference created",
		"talent_id", created.TalentID,
		"preference_id", created.ID,
	)

	return &created, nil
}

// GetByTalentID returns the current preferences of a talent.
func (r *talentPreferenceRepository) GetByTalentID(
	ctx context.Context,
	talentID uuid.UUID,
) (*model.TalentPreference, error) {

	const query = `
		SELECT
			id,
			talent_id,
			preferred_role,
			preferred_location,
			availability_status,
			created_at,
			updated_at
		FROM talent_preferences
		WHERE talent_id = $1
	`

	var preference model.TalentPreference

	err := r.db.QueryRow(
		ctx,
		query,
		talentID,
	).Scan(
		&preference.ID,
		&preference.TalentID,
		&preference.PreferredRole,
		&preference.PreferredLocation,
		&preference.AvailabilityStatus,
		&preference.CreatedAt,
		&preference.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Info(
				"talent preference not found",
				"talent_id", talentID,
			)

			return nil, ErrTalentPreferenceNotFound
		}

		r.logger.Error(
			"failed to get talent preference",
			"talent_id", talentID,
			"error", err,
		)

		return nil, fmt.Errorf(
			"get talent preference: %w",
			err,
		)
	}

	return &preference, nil
}

// Update modifies the current preferences of a talent.
func (r *talentPreferenceRepository) Update(
	ctx context.Context,
	preference *model.TalentPreference,
) (*model.TalentPreference, error) {

	const query = `
		UPDATE talent_preferences
		SET
			preferred_role = $1,
			preferred_location = $2,
			availability_status = $3,
			updated_at = CURRENT_TIMESTAMP
		WHERE talent_id = $4
		RETURNING
			id,
			talent_id,
			preferred_role,
			preferred_location,
			availability_status,
			created_at,
			updated_at
	`

	var updated model.TalentPreference

	err := r.db.QueryRow(
		ctx,
		query,
		preference.PreferredRole,
		preference.PreferredLocation,
		preference.AvailabilityStatus,
		preference.TalentID,
	).Scan(
		&updated.ID,
		&updated.TalentID,
		&updated.PreferredRole,
		&updated.PreferredLocation,
		&updated.AvailabilityStatus,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTalentPreferenceNotFound
		}

		r.logger.Error(
			"failed to update talent preference",
			"talent_id", preference.TalentID,
			"error", err,
		)

		return nil, fmt.Errorf(
			"update talent preference: %w",
			err,
		)
	}

	r.logger.Info(
		"talent preference updated",
		"talent_id", updated.TalentID,
		"preference_id", updated.ID,
	)

	return &updated, nil
}

// UpdateAvailability changes only the availability status
// of an existing talent preference.
func (r *talentPreferenceRepository) UpdateAvailability(
	ctx context.Context,
	talentID uuid.UUID,
	availabilityStatus string,
) (*model.TalentPreference, error) {

	const query = `
		UPDATE talent_preferences
		SET
			availability_status = $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE talent_id = $2
		RETURNING
			id,
			talent_id,
			preferred_role,
			preferred_location,
			availability_status,
			created_at,
			updated_at
	`

	var updated model.TalentPreference

	err := r.db.QueryRow(
		ctx,
		query,
		availabilityStatus,
		talentID,
	).Scan(
		&updated.ID,
		&updated.TalentID,
		&updated.PreferredRole,
		&updated.PreferredLocation,
		&updated.AvailabilityStatus,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Info(
				"talent preference not found for availability update",
				"talent_id", talentID,
			)

			return nil, ErrTalentPreferenceNotFound
		}

		r.logger.Error(
			"failed to update talent availability",
			"talent_id", talentID,
			"error", err,
		)

		return nil, fmt.Errorf(
			"update talent availability: %w",
			err,
		)
	}

	r.logger.Info(
		"talent availability updated",
		"talent_id", updated.TalentID,
		"availability_status", updated.AvailabilityStatus,
	)

	return &updated, nil
}

package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"

	"talentiq/talent-profile-intelligence/internal/model"
)

// Common repository errors.
//
// The service layer can use these errors later to decide
// whether an error should become a 404, 409, 500, etc.
var (
	ErrTalentProfileNotFound = errors.New("talent profile not found")
)

// TalentProfileRepository defines all database operations
// supported by the Talent Profile domain.
//
// The service layer will depend on this interface instead
// of directly depending on PostgreSQL.
type TalentProfileRepository interface {
	Create(
		ctx context.Context,
		profile *model.TalentProfile,
	) (*model.TalentProfile, error)

	GetByID(
		ctx context.Context,
		id uuid.UUID,
	) (*model.TalentProfile, error)

	GetByIDIncludingDeleted(
		ctx context.Context,
		id uuid.UUID,
	) (*model.TalentProfile, error)

	List(
		ctx context.Context,
		limit int,
		offset int,
	) ([]*model.TalentProfile, int, error)

	Update(
		ctx context.Context,
		profile *model.TalentProfile,
	) (*model.TalentProfile, error)

	Delete(
		ctx context.Context,
		id uuid.UUID,
	) error

	// Checks whether an active profile already uses
	// the supplied employee code.
	ExistsByEmployeeCode(
		ctx context.Context,
		employeeCode string,
	) (bool, error)
}

// talentProfileRepository is the PostgreSQL implementation
// of TalentProfileRepository.
type talentProfileRepository struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

// NewTalentProfileRepository creates a PostgreSQL-backed
// Talent Profile repository.
func NewTalentProfileRepository(
	db *pgxpool.Pool,
	logger *slog.Logger,
) TalentProfileRepository {

	return &talentProfileRepository{
		db:     db,
		logger: logger,
	}
}

// Create inserts a new talent profile into PostgreSQL.
func (r *talentProfileRepository) Create(
	ctx context.Context,
	profile *model.TalentProfile,
) (*model.TalentProfile, error) {

	const query = `
		INSERT INTO talent_profiles (
			employee_code,
			first_name,
			last_name,
			email,
			phone,
			designation,
			department,
			location,
			summary,
			total_experience_years,
			profile_status
		)
		VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10, $11
		)
		RETURNING
			id,
			employee_code,
			first_name,
			last_name,
			email,
			phone,
			designation,
			department,
			location,
			summary,
			total_experience_years,
			profile_status,
			created_at,
			updated_at,
			deleted_at
	`

	var created model.TalentProfile

	err := r.db.QueryRow(
		ctx,
		query,
		profile.EmployeeCode,
		profile.FirstName,
		profile.LastName,
		profile.Email,
		profile.Phone,
		profile.Designation,
		profile.Department,
		profile.Location,
		profile.Summary,
		profile.TotalExperienceYears,
		profile.ProfileStatus,
	).Scan(
		&created.ID,
		&created.EmployeeCode,
		&created.FirstName,
		&created.LastName,
		&created.Email,
		&created.Phone,
		&created.Designation,
		&created.Department,
		&created.Location,
		&created.Summary,
		&created.TotalExperienceYears,
		&created.ProfileStatus,
		&created.CreatedAt,
		&created.UpdatedAt,
		&created.DeletedAt,
	)

	if err != nil {
		r.logger.Error(
			"failed to create talent profile",
			"employee_code", profile.EmployeeCode,
			"error", err,
		)

		return nil, fmt.Errorf(
			"create talent profile: %w",
			err,
		)
	}

	r.logger.Info(
		"talent profile created",
		"profile_id", created.ID,
		"employee_code", created.EmployeeCode,
	)

	return &created, nil
}

// GetByID returns one active talent profile.
func (r *talentProfileRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*model.TalentProfile, error) {

	const query = `
		SELECT
			id,
			employee_code,
			first_name,
			last_name,
			email,
			phone,
			designation,
			department,
			location,
			summary,
			total_experience_years,
			profile_status,
			created_at,
			updated_at,
			deleted_at
		FROM talent_profiles
		WHERE id = $1
		  AND deleted_at IS NULL
	`

	var profile model.TalentProfile

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&profile.ID,
		&profile.EmployeeCode,
		&profile.FirstName,
		&profile.LastName,
		&profile.Email,
		&profile.Phone,
		&profile.Designation,
		&profile.Department,
		&profile.Location,
		&profile.Summary,
		&profile.TotalExperienceYears,
		&profile.ProfileStatus,
		&profile.CreatedAt,
		&profile.UpdatedAt,
		&profile.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Info(
				"talent profile not found",
				"profile_id", id,
			)

			return nil, ErrTalentProfileNotFound
		}

		r.logger.Error(
			"failed to get talent profile",
			"profile_id", id,
			"error", err,
		)

		return nil, fmt.Errorf(
			"get talent profile: %w",
			err,
		)
	}

	return &profile, nil
}

// List returns active talent profiles.
//
// limit controls the number of records returned.
// offset controls where PostgreSQL starts returning records.
func (r *talentProfileRepository) List(
	ctx context.Context,
	limit int,
	offset int,
) ([]*model.TalentProfile, int, error) {

	const listQuery = `
		SELECT
			id,
			employee_code,
			first_name,
			last_name,
			email,
			phone,
			designation,
			department,
			location,
			summary,
			total_experience_years,
			profile_status,
			created_at,
			updated_at,
			deleted_at
		FROM talent_profiles
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	const countQuery = `
		SELECT COUNT(*)
		FROM talent_profiles
		WHERE deleted_at IS NULL
	`

	var total int

	if err := r.db.QueryRow(
		ctx,
		countQuery,
	).Scan(&total); err != nil {

		r.logger.Error(
			"failed to count talent profiles",
			"error", err,
		)

		return nil, 0, fmt.Errorf(
			"count talent profiles: %w",
			err,
		)
	}

	rows, err := r.db.Query(
		ctx,
		listQuery,
		limit,
		offset,
	)

	if err != nil {
		r.logger.Error(
			"failed to list talent profiles",
			"error", err,
		)

		return nil, 0, fmt.Errorf(
			"list talent profiles: %w",
			err,
		)
	}

	defer rows.Close()

	profiles := make([]*model.TalentProfile, 0)

	for rows.Next() {
		var profile model.TalentProfile

		if err := rows.Scan(
			&profile.ID,
			&profile.EmployeeCode,
			&profile.FirstName,
			&profile.LastName,
			&profile.Email,
			&profile.Phone,
			&profile.Designation,
			&profile.Department,
			&profile.Location,
			&profile.Summary,
			&profile.TotalExperienceYears,
			&profile.ProfileStatus,
			&profile.CreatedAt,
			&profile.UpdatedAt,
			&profile.DeletedAt,
		); err != nil {

			r.logger.Error(
				"failed to scan talent profile",
				"error", err,
			)

			return nil, 0, fmt.Errorf(
				"scan talent profile: %w",
				err,
			)
		}

		profiles = append(profiles, &profile)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf(
			"iterate talent profiles: %w",
			err,
		)
	}

	return profiles, total, nil
}

// Update modifies an existing talent profile.
func (r *talentProfileRepository) Update(
	ctx context.Context,
	profile *model.TalentProfile,
) (*model.TalentProfile, error) {

	const query = `
		UPDATE talent_profiles
		SET
			first_name = $2,
			last_name = $3,
			email = $4,
			phone = $5,
			designation = $6,
			department = $7,
			location = $8,
			summary = $9,
			total_experience_years = $10,
			profile_status = $11,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		  AND deleted_at IS NULL
		RETURNING
			id,
			employee_code,
			first_name,
			last_name,
			email,
			phone,
			designation,
			department,
			location,
			summary,
			total_experience_years,
			profile_status,
			created_at,
			updated_at,
			deleted_at
	`

	var updated model.TalentProfile

	err := r.db.QueryRow(
		ctx,
		query,
		profile.ID,
		profile.FirstName,
		profile.LastName,
		profile.Email,
		profile.Phone,
		profile.Designation,
		profile.Department,
		profile.Location,
		profile.Summary,
		profile.TotalExperienceYears,
		profile.ProfileStatus,
	).Scan(
		&updated.ID,
		&updated.EmployeeCode,
		&updated.FirstName,
		&updated.LastName,
		&updated.Email,
		&updated.Phone,
		&updated.Designation,
		&updated.Department,
		&updated.Location,
		&updated.Summary,
		&updated.TotalExperienceYears,
		&updated.ProfileStatus,
		&updated.CreatedAt,
		&updated.UpdatedAt,
		&updated.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTalentProfileNotFound
		}

		r.logger.Error(
			"failed to update talent profile",
			"profile_id", profile.ID,
			"error", err,
		)

		return nil, fmt.Errorf(
			"update talent profile: %w",
			err,
		)
	}

	r.logger.Info(
		"talent profile updated",
		"profile_id", updated.ID,
	)

	return &updated, nil
}

// Delete performs a soft delete.
//
// The database record remains available for auditing,
// but normal queries no longer return it.
func (r *talentProfileRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {

	const query = `
		UPDATE talent_profiles
		SET
			deleted_at = CURRENT_TIMESTAMP,
			updated_at = CURRENT_TIMESTAMP,
			profile_status = 'INACTIVE'
		WHERE id = $1
		  AND deleted_at IS NULL
	`

	result, err := r.db.Exec(
		ctx,
		query,
		id,
	)

	if err != nil {
		r.logger.Error(
			"failed to delete talent profile",
			"profile_id", id,
			"error", err,
		)

		return fmt.Errorf(
			"delete talent profile: %w",
			err,
		)
	}

	if result.RowsAffected() == 0 {
		return ErrTalentProfileNotFound
	}

	r.logger.Info(
		"talent profile deactivated",
		"profile_id", id,
	)

	return nil

}

// ExistsByEmployeeCode checks whether an active talent profile
// already exists with the given employee code.
func (r *talentProfileRepository) ExistsByEmployeeCode(
	ctx context.Context,
	employeeCode string,
) (bool, error) {

	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM talent_profiles
			WHERE employee_code = $1
			  AND deleted_at IS NULL
		)
	`

	var exists bool

	err := r.db.QueryRow(
		ctx,
		query,
		employeeCode,
	).Scan(&exists)

	if err != nil {
		r.logger.Error(
			"failed to check employee code",
			"employee_code", employeeCode,
			"error", err,
		)

		return false, fmt.Errorf(
			"check employee code: %w",
			err,
		)
	}

	return exists, nil
}

// GetByIDIncludingDeleted returns a talent profile
// regardless of its deleted status.
//
// This is mainly used for event processing where a
// deactivated profile may already have been soft deleted.
func (r *talentProfileRepository) GetByIDIncludingDeleted(
	ctx context.Context,
	id uuid.UUID,
) (*model.TalentProfile, error) {

	const query = `
		SELECT
			id,
			employee_code,
			first_name,
			last_name,
			email,
			phone,
			designation,
			department,
			location,
			summary,
			total_experience_years,
			profile_status,
			created_at,
			updated_at,
			deleted_at
		FROM talent_profiles
		WHERE id = $1
	`

	var profile model.TalentProfile

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&profile.ID,
		&profile.EmployeeCode,
		&profile.FirstName,
		&profile.LastName,
		&profile.Email,
		&profile.Phone,
		&profile.Designation,
		&profile.Department,
		&profile.Location,
		&profile.Summary,
		&profile.TotalExperienceYears,
		&profile.ProfileStatus,
		&profile.CreatedAt,
		&profile.UpdatedAt,
		&profile.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Info(
				"talent profile not found",
				"profile_id", id,
			)

			return nil, ErrTalentProfileNotFound
		}

		r.logger.Error(
			"failed to get talent profile including deleted profile",
			"profile_id", id,
			"error", err,
		)

		return nil, fmt.Errorf(
			"get talent profile including deleted: %w",
			err,
		)
	}

	return &profile, nil
}

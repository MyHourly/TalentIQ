package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"talentiq/talent-profile-intelligence/internal/model"
)

// Named errors, so the other layers can tell "what went wrong" without reading text.
var (
	ErrInvalidTalentExperience     = errors.New("invalid talent experience") // bad input -> 400
	ErrTalentNotFoundForExperience = errors.New("talent not found")          // person doesn't exist -> 404
	ErrExperienceAlreadyExists     = errors.New("experience already added")  // same job added twice -> 409
)

// The list of database jobs this file can do.
type TalentExperienceRepository interface {
	TalentExists(ctx context.Context, talentID uuid.UUID) (bool, error)
	Create(ctx context.Context, exp *model.TalentExperience) (*model.TalentExperience, error)
}

type talentExperienceRepository struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewTalentExperienceRepository(db *pgxpool.Pool, logger *slog.Logger) TalentExperienceRepository {
	return &talentExperienceRepository{db: db, logger: logger}
}

// Job 1: "Is there an active person with this id?"
func (r *talentExperienceRepository) TalentExists(ctx context.Context, talentID uuid.UUID) (bool, error) {
	const query = `
		SELECT EXISTS (
			SELECT 1 FROM talent_profiles
			WHERE id = $1 AND deleted_at IS NULL
		)
	`

	var exists bool
	if err := r.db.QueryRow(ctx, query, talentID).Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to check talent exists: %w", err)
	}
	return exists, nil
}

// Job 2: "Save this job in the person's work history."
func (r *talentExperienceRepository) Create(ctx context.Context, exp *model.TalentExperience) (*model.TalentExperience, error) {
	const query = `
		INSERT INTO professional_experience
			(talent_id, company_name, job_title, start_date, end_date, description)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING
			id, talent_id, company_name, job_title, start_date, end_date,
			description, created_at, updated_at
	`

	saved := &model.TalentExperience{}
	err := r.db.QueryRow(ctx, query,
		exp.TalentID,
		exp.CompanyName,
		exp.JobTitle,
		exp.StartDate,
		exp.EndDate,
		exp.Description,
	).Scan(
		&saved.ID,
		&saved.TalentID,
		&saved.CompanyName,
		&saved.JobTitle,
		&saved.StartDate,
		&saved.EndDate,
		&saved.Description,
		&saved.CreatedAt,
		&saved.UpdatedAt,
	)
	if err != nil {
		// isUniqueViolation is the helper from talent_certification_repository.go
		// (same package): it spots PostgreSQL's "this row already exists" error.
		if isUniqueViolation(err) {
			return nil, ErrExperienceAlreadyExists
		}
		return nil, fmt.Errorf("failed to create talent experience: %w", err)
	}

	r.logger.Info("talent experience added",
		"talent_id", saved.TalentID,
		"company_name", saved.CompanyName,
		"job_title", saved.JobTitle,
	)
	return saved, nil
}

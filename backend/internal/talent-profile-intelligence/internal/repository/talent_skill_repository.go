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
	ErrTalentSkillNotFound    = errors.New("talent skill not found") // person doesn't have that skill
	ErrInvalidTalentSkill     = errors.New("invalid talent skill")   // bad input
	ErrTalentNotFoundForSkill = errors.New("talent not found")       // person doesn't exist
)

// The list of database jobs this file can do.
type TalentSkillRepository interface {
	TalentExists(ctx context.Context, talentID uuid.UUID) (bool, error)
	Upsert(ctx context.Context, skill *model.TalentSkill) (*model.TalentSkill, error)
	Delete(ctx context.Context, talentID, skillID uuid.UUID) error
}

type talentSkillRepository struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

// Takes the logger as the second argument, like your other repositories.
func NewTalentSkillRepository(db *pgxpool.Pool, logger *slog.Logger) TalentSkillRepository {
	return &talentSkillRepository{db: db, logger: logger}
}

// Job 1: "Is there an active person with this id?"
// deleted_at IS NULL means the profile was not deactivated.
func (r *talentSkillRepository) TalentExists(ctx context.Context, talentID uuid.UUID) (bool, error) {
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

// Job 2: "Save this skill. If the person already has it, just change the level."
// ON CONFLICT is the database doing the add-or-update for us in one step.
func (r *talentSkillRepository) Upsert(ctx context.Context, skill *model.TalentSkill) (*model.TalentSkill, error) {
	const query = `
		INSERT INTO talent_skills (talent_id, skill_id, proficiency_level)
		VALUES ($1, $2, $3)
		ON CONFLICT (talent_id, skill_id)
		DO UPDATE SET
			proficiency_level = EXCLUDED.proficiency_level,
			updated_at = CURRENT_TIMESTAMP
		RETURNING
			id, talent_id, skill_id, proficiency_level, created_at, updated_at
	`

	saved := &model.TalentSkill{}
	err := r.db.QueryRow(ctx, query,
		skill.TalentID,
		skill.SkillID,
		skill.ProficiencyLevel,
	).Scan(
		&saved.ID,
		&saved.TalentID,
		&saved.SkillID,
		&saved.ProficiencyLevel,
		&saved.CreatedAt,
		&saved.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert talent skill: %w", err)
	}

	r.logger.Info("talent skill saved",
		"talent_id", saved.TalentID,
		"skill_id", saved.SkillID,
		"proficiency_level", saved.ProficiencyLevel,
	)
	return saved, nil
}

// Job 3: "Remove this skill from this person."
// If nothing was deleted, the person never had that skill, so we say "not found".
func (r *talentSkillRepository) Delete(ctx context.Context, talentID, skillID uuid.UUID) error {
	const query = `
		DELETE FROM talent_skills
		WHERE talent_id = $1 AND skill_id = $2
	`

	result, err := r.db.Exec(ctx, query, talentID, skillID)
	if err != nil {
		return fmt.Errorf("failed to delete talent skill: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrTalentSkillNotFound
	}

	r.logger.Info("talent skill removed",
		"talent_id", talentID,
		"skill_id", skillID,
	)
	return nil
}
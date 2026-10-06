package repository

import (
	"context"
	"fmt"

	"github.com/MyHourly/TalentIQ/backend/internal/skill-intelligence/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProficiencyRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewProficiencyRepository(
	db *pgxpool.Pool,
) *ProficiencyRepositoryImpl {
	return &ProficiencyRepositoryImpl{
		db: db,
	}
}

func (r *ProficiencyRepositoryImpl) Create(
	ctx context.Context,
	proficiency *model.SkillProficiency,
) error {

	query := `
		INSERT INTO skill_proficiencies
		(
			id,
			name,
			description,
			level_order,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		proficiency.ID,
		proficiency.Name,
		proficiency.Description,
		proficiency.LevelOrder,
	)

	if err != nil {
		return fmt.Errorf(
			"failed to create proficiency: %w",
			err,
		)
	}

	return nil
}

func (r *ProficiencyRepositoryImpl) GetAll(
	ctx context.Context,
) ([]model.SkillProficiency, error) {

	query := `
		SELECT
			id,
			name,
			description,
			level_order,
			created_at,
			updated_at
		FROM skill_proficiencies
		ORDER BY level_order ASC
	`

	rows, err := r.db.Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to get proficiencies: %w",
			err,
		)
	}

	defer rows.Close()

	proficiencies := make(
		[]model.SkillProficiency,
		0,
	)

	for rows.Next() {

		var proficiency model.SkillProficiency

		err := rows.Scan(
			&proficiency.ID,
			&proficiency.Name,
			&proficiency.Description,
			&proficiency.LevelOrder,
			&proficiency.CreatedAt,
			&proficiency.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"failed to scan proficiency: %w",
				err,
			)
		}

		proficiencies = append(
			proficiencies,
			proficiency,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"failed to iterate proficiencies: %w",
			err,
		)
	}

	return proficiencies, nil
}

func (r *ProficiencyRepositoryImpl) GetByID(
	ctx context.Context,
	id string,
) (*model.SkillProficiency, error) {

	query := `
		SELECT
			id,
			name,
			description,
			level_order,
			created_at,
			updated_at
		FROM skill_proficiencies
		WHERE id = $1
	`

	var proficiency model.SkillProficiency

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&proficiency.ID,
		&proficiency.Name,
		&proficiency.Description,
		&proficiency.LevelOrder,
		&proficiency.CreatedAt,
		&proficiency.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &proficiency, nil
}

func (r *ProficiencyRepositoryImpl) AssignToSkill(
	ctx context.Context,
	skillID string,
	proficiencyID string,
) error {

	query := `
		INSERT INTO skill_proficiency_mapping
		(
			skill_id,
			proficiency_id,
			created_at,
			updated_at
		)
		VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT (skill_id)
		DO UPDATE SET
			proficiency_id = EXCLUDED.proficiency_id,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err := r.db.Exec(
		ctx,
		query,
		skillID,
		proficiencyID,
	)

	if err != nil {
		return fmt.Errorf(
			"failed to assign proficiency to skill: %w",
			err,
		)
	}

	return nil
}

func (r *ProficiencyRepositoryImpl) GetBySkillID(
	ctx context.Context,
	skillID string,
) (*model.SkillProficiency, error) {

	query := `
		SELECT
			p.id,
			p.name,
			p.description,
			p.level_order,
			p.created_at,
			p.updated_at
		FROM skill_proficiency_mapping spm
		JOIN skill_proficiencies p
			ON p.id = spm.proficiency_id
		WHERE spm.skill_id = $1
	`

	var proficiency model.SkillProficiency

	err := r.db.QueryRow(
		ctx,
		query,
		skillID,
	).Scan(
		&proficiency.ID,
		&proficiency.Name,
		&proficiency.Description,
		&proficiency.LevelOrder,
		&proficiency.CreatedAt,
		&proficiency.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &proficiency, nil
}

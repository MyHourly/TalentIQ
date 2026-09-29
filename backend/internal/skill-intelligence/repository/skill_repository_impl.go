package repository

import (
	"context"
	"fmt"

	"github.com/MyHourly/TalentIQ/backend/internal/skill-intelligence/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresSkillRepository struct {
	db *pgxpool.Pool
}

func NewPostgresSkillRepository(db *pgxpool.Pool) *PostgresSkillRepository {
	return &PostgresSkillRepository{
		db: db,
	}
}

func (r *PostgresSkillRepository) Create(
	ctx context.Context,
	skill *model.Skill,
) error {

	query := `
		INSERT INTO skills (
			id,
			name,
			description,
			status,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		skill.ID,
		skill.Name,
		skill.Description,
		skill.Status,
		skill.CreatedAt,
		skill.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create skill: %w", err)
	}

	return nil
}

func (r *PostgresSkillRepository) GetAll(
	ctx context.Context,
) ([]model.Skill, error) {

	query := `
		SELECT
			id,
			name,
			description,
			status,
			created_at,
			updated_at
		FROM skills
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("failed to get skills: %w", err)
	}

	defer rows.Close()

	var skills []model.Skill

	for rows.Next() {

		var skill model.Skill

		err := rows.Scan(
			&skill.ID,
			&skill.Name,
			&skill.Description,
			&skill.Status,
			&skill.CreatedAt,
			&skill.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan skill: %w", err)
		}

		skills = append(skills, skill)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while reading skills: %w", err)
	}

	return skills, nil
}

func (r *PostgresSkillRepository) GetByID(
	ctx context.Context,
	id string,
) (*model.Skill, error) {

	query := `
		SELECT
			id,
			name,
			description,
			status,
			created_at,
			updated_at
		FROM skills
		WHERE id = $1
	`

	var skill model.Skill

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&skill.ID,
		&skill.Name,
		&skill.Description,
		&skill.Status,
		&skill.CreatedAt,
		&skill.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get skill: %w", err)
	}

	return &skill, nil
}

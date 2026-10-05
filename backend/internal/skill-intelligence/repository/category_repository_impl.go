package repository

import (
	"context"
	"fmt"

	"github.com/MyHourly/TalentIQ/backend/internal/skill-intelligence/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresCategoryRepository struct {
	db *pgxpool.Pool
}

func NewPostgresCategoryRepository(
	db *pgxpool.Pool,
) *PostgresCategoryRepository {

	return &PostgresCategoryRepository{
		db: db,
	}
}

func (r *PostgresCategoryRepository) Create(
	ctx context.Context,
	category *model.SkillCategory,
) error {

	query := `
		INSERT INTO skill_categories (
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
		category.ID,
		category.Name,
		category.Description,
		category.Status,
		category.CreatedAt,
		category.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf(
			"failed to create category: %w",
			err,
		)
	}

	return nil
}

func (r *PostgresCategoryRepository) GetAll(
	ctx context.Context,
) ([]model.SkillCategory, error) {

	query := `
		SELECT
			id,
			name,
			description,
			status,
			created_at,
			updated_at
		FROM skill_categories
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to get categories: %w",
			err,
		)
	}

	defer rows.Close()

	var categories []model.SkillCategory

	for rows.Next() {

		var category model.SkillCategory

		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.Status,
			&category.CreatedAt,
			&category.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"failed to scan category: %w",
				err,
			)
		}

		categories = append(
			categories,
			category,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"error while reading categories: %w",
			err,
		)
	}

	return categories, nil
}

func (r *PostgresCategoryRepository) GetByID(
	ctx context.Context,
	id string,
) (*model.SkillCategory, error) {

	query := `
		SELECT
			id,
			name,
			description,
			status,
			created_at,
			updated_at
		FROM skill_categories
		WHERE id = $1
	`

	var category model.SkillCategory

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.Status,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to get category: %w",
			err,
		)
	}

	return &category, nil
}

func (r *PostgresCategoryRepository) MapSkillToCategory(
	ctx context.Context,
	skillID string,
	categoryID string,
) error {

	query := `
		INSERT INTO skill_category_mapping (
			skill_id,
			category_id
		)
		VALUES ($1, $2)
		ON CONFLICT (skill_id, category_id)
		DO NOTHING
	`

	_, err := r.db.Exec(
		ctx,
		query,
		skillID,
		categoryID,
	)

	if err != nil {
		return fmt.Errorf(
			"failed to map skill to category: %w",
			err,
		)
	}

	return nil
}

func (r *PostgresCategoryRepository) GetCategoriesBySkillID(
	ctx context.Context,
	skillID string,
) ([]model.SkillCategory, error) {

	query := `
		SELECT
			c.id,
			c.name,
			c.description,
			c.status,
			c.created_at,
			c.updated_at
		FROM skill_categories c
		INNER JOIN skill_category_mapping scm
			ON c.id = scm.category_id
		WHERE scm.skill_id = $1
		ORDER BY c.name
	`

	rows, err := r.db.Query(
		ctx,
		query,
		skillID,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to get skill categories: %w",
			err,
		)
	}

	defer rows.Close()

	var categories []model.SkillCategory

	for rows.Next() {

		var category model.SkillCategory

		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.Status,
			&category.CreatedAt,
			&category.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"failed to scan category: %w",
				err,
			)
		}

		categories = append(
			categories,
			category,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"error while reading skill categories: %w",
			err,
		)
	}

	return categories, nil
}

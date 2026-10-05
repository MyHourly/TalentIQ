package repository

import (
	"context"

	"github.com/MyHourly/TalentIQ/backend/internal/skill-intelligence/model"
)

type CategoryRepository interface {
	Create(
		ctx context.Context,
		category *model.SkillCategory,
	) error

	GetAll(
		ctx context.Context,
	) ([]model.SkillCategory, error)

	GetByID(
		ctx context.Context,
		id string,
	) (*model.SkillCategory, error)

	MapSkillToCategory(
		ctx context.Context,
		skillID string,
		categoryID string,
	) error

	GetCategoriesBySkillID(
		ctx context.Context,
		skillID string,
	) ([]model.SkillCategory, error)
}

package repository

import (
	"context"

	"github.com/MyHourly/TalentIQ/backend/internal/skill-intelligence/model"
)

type ProficiencyRepository interface {
	Create(
		ctx context.Context,
		proficiency *model.SkillProficiency,
	) error

	GetAll(
		ctx context.Context,
	) ([]model.SkillProficiency, error)

	GetByID(
		ctx context.Context,
		id string,
	) (*model.SkillProficiency, error)

	AssignToSkill(
		ctx context.Context,
		skillID string,
		proficiencyID string,
	) error

	GetBySkillID(
		ctx context.Context,
		skillID string,
	) (*model.SkillProficiency, error)
}

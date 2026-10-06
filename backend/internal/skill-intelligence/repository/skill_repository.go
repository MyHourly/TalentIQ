package repository

import (
	"context"

	"github.com/MyHourly/TalentIQ/backend/internal/skill-intelligence/model"
)

type SkillRepository interface {
	Create(ctx context.Context, skill *model.Skill) error

	GetAll(ctx context.Context) ([]model.Skill, error)

	GetByID(ctx context.Context, id string) (*model.Skill, error)

	Update(ctx context.Context, skill *model.Skill) error

	Deactivate(ctx context.Context, id string) error
}

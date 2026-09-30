package service

import (
	"context"

	"github.com/MyHourly/TalentIQ/backend/internal/skill-intelligence/model"
	"github.com/MyHourly/TalentIQ/backend/internal/skill-intelligence/repository"
)

type SkillService struct {
	repository repository.SkillRepository
}

func NewSkillService(repo repository.SkillRepository) *SkillService {
	return &SkillService{
		repository: repo,
	}
}

func (s *SkillService) CreateSkill(
	ctx context.Context,
	skill *model.Skill,
) error {

	return s.repository.Create(ctx, skill)
}

func (s *SkillService) GetAllSkills(
	ctx context.Context,
) ([]model.Skill, error) {

	return s.repository.GetAll(ctx)
}

func (s *SkillService) GetSkillByID(
	ctx context.Context,
	id string,
) (*model.Skill, error) {

	return s.repository.GetByID(ctx, id)
}

func (s *SkillService) UpdateSkill(
	ctx context.Context,
	skill *model.Skill,
) error {

	return s.repository.Update(ctx, skill)
}

func (s *SkillService) DeactivateSkill(
	ctx context.Context,
	id string,
) error {

	return s.repository.Deactivate(ctx, id)
}

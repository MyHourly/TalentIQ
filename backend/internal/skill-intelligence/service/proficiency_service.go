package service

import (
	"context"

	"github.com/MyHourly/TalentIQ/backend/internal/skill-intelligence/model"
	"github.com/MyHourly/TalentIQ/backend/internal/skill-intelligence/repository"
)

type ProficiencyService struct {
	repository repository.ProficiencyRepository
}

func NewProficiencyService(
	repo repository.ProficiencyRepository,
) *ProficiencyService {
	return &ProficiencyService{
		repository: repo,
	}
}

func (s *ProficiencyService) CreateProficiency(
	ctx context.Context,
	proficiency *model.SkillProficiency,
) error {
	return s.repository.Create(ctx, proficiency)
}

func (s *ProficiencyService) GetAllProficiencies(
	ctx context.Context,
) ([]model.SkillProficiency, error) {
	return s.repository.GetAll(ctx)
}

func (s *ProficiencyService) GetProficiencyByID(
	ctx context.Context,
	id string,
) (*model.SkillProficiency, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *ProficiencyService) AssignProficiencyToSkill(
	ctx context.Context,
	skillID string,
	proficiencyID string,
) error {
	return s.repository.AssignToSkill(
		ctx,
		skillID,
		proficiencyID,
	)
}

func (s *ProficiencyService) GetProficiencyBySkillID(
	ctx context.Context,
	skillID string,
) (*model.SkillProficiency, error) {
	return s.repository.GetBySkillID(
		ctx,
		skillID,
	)
}

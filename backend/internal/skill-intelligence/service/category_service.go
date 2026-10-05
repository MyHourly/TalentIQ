package service

import (
	"context"

	"github.com/MyHourly/TalentIQ/backend/internal/skill-intelligence/model"
	"github.com/MyHourly/TalentIQ/backend/internal/skill-intelligence/repository"
)

type CategoryService struct {
	repository repository.CategoryRepository
}

func NewCategoryService(
	repo repository.CategoryRepository,
) *CategoryService {

	return &CategoryService{
		repository: repo,
	}
}

func (s *CategoryService) CreateCategory(
	ctx context.Context,
	category *model.SkillCategory,
) error {

	return s.repository.Create(
		ctx,
		category,
	)
}

func (s *CategoryService) GetAllCategories(
	ctx context.Context,
) ([]model.SkillCategory, error) {

	return s.repository.GetAll(ctx)
}

func (s *CategoryService) GetCategoryByID(
	ctx context.Context,
	id string,
) (*model.SkillCategory, error) {

	return s.repository.GetByID(
		ctx,
		id,
	)
}

func (s *CategoryService) MapSkillToCategory(
	ctx context.Context,
	skillID string,
	categoryID string,
) error {

	return s.repository.MapSkillToCategory(
		ctx,
		skillID,
		categoryID,
	)
}

func (s *CategoryService) GetCategoriesBySkillID(
	ctx context.Context,
	skillID string,
) ([]model.SkillCategory, error) {

	return s.repository.GetCategoriesBySkillID(
		ctx,
		skillID,
	)
}

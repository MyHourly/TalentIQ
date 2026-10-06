package mapper

import (
	"github.com/MyHourly/TalentIQ/backend/internal/skill-intelligence/dto"
	"github.com/MyHourly/TalentIQ/backend/internal/skill-intelligence/model"
)

func ToCategoryModel(
	request dto.CreateCategoryRequest,
) *model.SkillCategory {

	return &model.SkillCategory{
		Name:        request.Name,
		Description: request.Description,
		Status:      request.Status,
	}
}

func ToCategoryResponse(
	category *model.SkillCategory,
) dto.CategoryResponse {

	return dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		Status:      category.Status,
	}
}

func ToCategoryResponseList(
	categories []model.SkillCategory,
) []dto.CategoryResponse {

	responses := make(
		[]dto.CategoryResponse,
		0,
		len(categories),
	)

	for _, category := range categories {

		responses = append(
			responses,
			ToCategoryResponse(&category),
		)
	}

	return responses
}

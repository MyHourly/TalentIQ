package mapper

import (
	"github.com/MyHourly/TalentIQ/backend/internal/skill-intelligence/dto"
	"github.com/MyHourly/TalentIQ/backend/internal/skill-intelligence/model"
)

func ToSkillModel(request dto.CreateSkillRequest) *model.Skill {
	return &model.Skill{
		Name:        request.Name,
		Description: request.Description,
		Status:      request.Status,
	}
}

func ToSkillResponse(skill *model.Skill) dto.SkillResponse {
	return dto.SkillResponse{
		ID:          skill.ID,
		Name:        skill.Name,
		Description: skill.Description,
		Status:      skill.Status,
	}
}

func ToSkillResponseList(skills []model.Skill) []dto.SkillResponse {

	responses := make([]dto.SkillResponse, 0, len(skills))

	for _, skill := range skills {
		responses = append(responses, ToSkillResponse(&skill))
	}

	return responses
}

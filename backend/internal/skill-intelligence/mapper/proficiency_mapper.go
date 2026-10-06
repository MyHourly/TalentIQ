package mapper

import (
	"github.com/MyHourly/TalentIQ/backend/internal/skill-intelligence/dto"
	"github.com/MyHourly/TalentIQ/backend/internal/skill-intelligence/model"
)

func ToProficiencyModel(
	request dto.CreateProficiencyRequest,
) *model.SkillProficiency {
	return &model.SkillProficiency{
		Name:        request.Name,
		Description: request.Description,
		LevelOrder:  request.LevelOrder,
	}
}

func ToProficiencyResponse(
	proficiency *model.SkillProficiency,
) dto.ProficiencyResponse {
	return dto.ProficiencyResponse{
		ID:          proficiency.ID,
		Name:        proficiency.Name,
		Description: proficiency.Description,
		LevelOrder:  proficiency.LevelOrder,
	}
}

func ToProficiencyResponseList(
	proficiencies []model.SkillProficiency,
) []dto.ProficiencyResponse {

	responses := make(
		[]dto.ProficiencyResponse,
		0,
		len(proficiencies),
	)

	for _, proficiency := range proficiencies {
		responses = append(
			responses,
			ToProficiencyResponse(&proficiency),
		)
	}

	return responses
}

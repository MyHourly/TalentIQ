package dto

import (
	"time"

	"talentiq/talent-profile-intelligence/internal/model"
)

const experienceDateLayout = "2006-01-02" // dates look like 2026-10-08

// What the caller sends. The talent id comes from the URL.
// end_date is optional (empty = current job).
type TalentExperienceRequest struct {
	CompanyName string  `json:"company_name"`
	JobTitle    string  `json:"job_title"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date"`
	Description string  `json:"description"`
}

// What we send back after a successful add.
type TalentExperienceResponse struct {
	ID          string  `json:"id"`
	TalentID    string  `json:"talent_id"`
	CompanyName string  `json:"company_name"`
	JobTitle    string  `json:"job_title"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date"`
	Description string  `json:"description"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// Turns a model into the reply shape (ids become text, dates become 2026-10-08).
func NewTalentExperienceResponse(e *model.TalentExperience) TalentExperienceResponse {
	var end *string
	if e.EndDate != nil {
		s := e.EndDate.Format(experienceDateLayout)
		end = &s
	}

	return TalentExperienceResponse{
		ID:          e.ID.String(),
		TalentID:    e.TalentID.String(),
		CompanyName: e.CompanyName,
		JobTitle:    e.JobTitle,
		StartDate:   e.StartDate.Format(experienceDateLayout),
		EndDate:     end,
		Description: e.Description,
		CreatedAt:   e.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   e.UpdatedAt.Format(time.RFC3339),
	}
}

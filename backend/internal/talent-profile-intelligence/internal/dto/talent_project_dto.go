package dto

import (
	"time"

	"talentiq/talent-profile-intelligence/internal/model"
)

const projectDateLayout = "2006-01-02" // dates look like 2026-10-08

// What the caller sends to add a project. The talent id comes from the URL.
type TalentProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Domain      string `json:"domain"`
	ClientName  string `json:"client_name"`
	Role        string `json:"role"`
	StartDate   string `json:"start_date"`
}

// What the caller sends to complete a project.
// end_date is optional; when it is left out, today's date is used.
type ProjectCompletionRequest struct {
	EndDate *string `json:"end_date"`
}

// What we send back for a project.
type TalentProjectResponse struct {
	ID          string  `json:"id"`
	TalentID    string  `json:"talent_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Domain      string  `json:"domain"`
	ClientName  string  `json:"client_name"`
	Role        string  `json:"role"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date"`
	Status      string  `json:"status"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// Turns a model into the reply shape (ids become text, dates become 2026-10-08).
func NewTalentProjectResponse(p *model.TalentProject) TalentProjectResponse {
	var end *string
	if p.EndDate != nil {
		s := p.EndDate.Format(projectDateLayout)
		end = &s
	}

	return TalentProjectResponse{
		ID:          p.ID.String(),
		TalentID:    p.TalentID.String(),
		Name:        p.Name,
		Description: p.Description,
		Domain:      p.Domain,
		ClientName:  p.ClientName,
		Role:        p.Role,
		StartDate:   p.StartDate.Format(projectDateLayout),
		EndDate:     end,
		Status:      p.Status,
		CreatedAt:   p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   p.UpdatedAt.Format(time.RFC3339),
	}
}

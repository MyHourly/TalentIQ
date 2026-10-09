package dto

import "time"

// What the caller sends us in the request body.
// The talent id and skill id come from the URL, so only the level is needed here.
type TalentSkillRequest struct {
	ProficiencyLevel int `json:"proficiency_level"`
}

// What we send back after a successful save.
type TalentSkillResponse struct {
	ID               string `json:"id"`
	TalentID         string `json:"talent_id"`
	SkillID          string `json:"skill_id"`
	ProficiencyLevel int    `json:"proficiency_level"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

// Turns our database values into the shape we return to the caller
// (ids become text, times become readable like 2026-10-08T13:31:00+05:30).
func NewTalentSkillResponse(id, talentID, skillID string, level int, createdAt, updatedAt time.Time) TalentSkillResponse {
	return TalentSkillResponse{
		ID:               id,
		TalentID:         talentID,
		SkillID:          skillID,
		ProficiencyLevel: level,
		CreatedAt:        createdAt.Format(time.RFC3339),
		UpdatedAt:        updatedAt.Format(time.RFC3339),
	}
}
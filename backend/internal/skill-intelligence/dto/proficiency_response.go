package dto

type ProficiencyResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	LevelOrder  int    `json:"level_order"`
}

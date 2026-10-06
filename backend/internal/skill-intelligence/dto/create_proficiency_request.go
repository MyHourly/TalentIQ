package dto

type CreateProficiencyRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	LevelOrder  int    `json:"level_order"`
}

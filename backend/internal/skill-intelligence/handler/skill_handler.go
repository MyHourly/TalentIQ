package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/MyHourly/TalentIQ/backend/internal/skill-intelligence/dto"
	"github.com/MyHourly/TalentIQ/backend/internal/skill-intelligence/mapper"
	"github.com/MyHourly/TalentIQ/backend/internal/skill-intelligence/service"
)

type SkillHandler struct {
	service *service.SkillService
}

func NewSkillHandler(skillService *service.SkillService) *SkillHandler {
	return &SkillHandler{
		service: skillService,
	}
}

func (h *SkillHandler) CreateSkill(w http.ResponseWriter, r *http.Request) {

	var request dto.CreateSkillRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	skill := mapper.ToSkillModel(request)

	skill.ID = uuid.New().String()
	skill.CreatedAt = time.Now()
	skill.UpdatedAt = time.Now()

	err = h.service.CreateSkill(r.Context(), skill)

	if err != nil {
		http.Error(w, "Failed to create skill", http.StatusInternalServerError)
		return
	}

	response := mapper.ToSkillResponse(skill)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(response)
}

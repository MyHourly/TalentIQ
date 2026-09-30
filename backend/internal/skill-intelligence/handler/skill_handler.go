package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/MyHourly/TalentIQ/backend/internal/skill-intelligence/dto"
	"github.com/MyHourly/TalentIQ/backend/internal/skill-intelligence/mapper"
	"github.com/MyHourly/TalentIQ/backend/internal/skill-intelligence/model"
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

// CreateSkill handles POST /api/v1/skills
func (h *SkillHandler) CreateSkill(
	w http.ResponseWriter,
	r *http.Request,
) {

	var request dto.CreateSkillRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(
			w,
			"Invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	skill := mapper.ToSkillModel(request)

	skill.ID = uuid.New().String()
	skill.CreatedAt = time.Now()
	skill.UpdatedAt = time.Now()

	err = h.service.CreateSkill(
		r.Context(),
		skill,
	)

	if err != nil {
		http.Error(
			w,
			"Failed to create skill",
			http.StatusInternalServerError,
		)
		return
	}

	response := mapper.ToSkillResponse(skill)

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(response)
}

// GetAllSkills handles GET /api/v1/skills
func (h *SkillHandler) GetAllSkills(
	w http.ResponseWriter,
	r *http.Request,
) {

	skills, err := h.service.GetAllSkills(
		r.Context(),
	)

	if err != nil {
		http.Error(
			w,
			"Failed to get skills",
			http.StatusInternalServerError,
		)
		return
	}

	responses := mapper.ToSkillResponseList(skills)

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(responses)
}

// GetSkillByID handles GET /api/v1/skills/{id}
func (h *SkillHandler) GetSkillByID(
	w http.ResponseWriter,
	r *http.Request,
) {

	id := getSkillID(r)

	if id == "" {
		http.Error(
			w,
			"Skill ID is required",
			http.StatusBadRequest,
		)
		return
	}

	skill, err := h.service.GetSkillByID(
		r.Context(),
		id,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(
				w,
				"Skill not found",
				http.StatusNotFound,
			)
			return
		}

		http.Error(
			w,
			"Failed to get skill",
			http.StatusInternalServerError,
		)
		return
	}

	response := mapper.ToSkillResponse(skill)

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

// UpdateSkill handles PUT /api/v1/skills/{id}
func (h *SkillHandler) UpdateSkill(
	w http.ResponseWriter,
	r *http.Request,
) {

	id := getSkillID(r)

	if id == "" {
		http.Error(
			w,
			"Skill ID is required",
			http.StatusBadRequest,
		)
		return
	}

	var request dto.UpdateSkillRequest

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(
			w,
			"Invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	skill := &model.Skill{
		ID:          id,
		Name:        request.Name,
		Description: request.Description,
		Status:      request.Status,
		UpdatedAt:   time.Now(),
	}

	err = h.service.UpdateSkill(
		r.Context(),
		skill,
	)

	if err != nil {

		if err.Error() == "skill not found" {
			http.Error(
				w,
				"Skill not found",
				http.StatusNotFound,
			)
			return
		}

		http.Error(
			w,
			"Failed to update skill",
			http.StatusInternalServerError,
		)
		return
	}

	updatedSkill, err := h.service.GetSkillByID(
		r.Context(),
		id,
	)

	if err != nil {
		http.Error(
			w,
			"Skill updated but failed to fetch updated skill",
			http.StatusInternalServerError,
		)
		return
	}

	response := mapper.ToSkillResponse(updatedSkill)

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

// DeactivateSkill handles DELETE /api/v1/skills/{id}
func (h *SkillHandler) DeactivateSkill(
	w http.ResponseWriter,
	r *http.Request,
) {

	id := getSkillID(r)

	if id == "" {
		http.Error(
			w,
			"Skill ID is required",
			http.StatusBadRequest,
		)
		return
	}

	err := h.service.DeactivateSkill(
		r.Context(),
		id,
	)

	if err != nil {

		if err.Error() == "skill not found" {
			http.Error(
				w,
				"Skill not found",
				http.StatusNotFound,
			)
			return
		}

		http.Error(
			w,
			"Failed to deactivate skill",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message": "Skill deactivated successfully",
		"id":      id,
		"status":  "INACTIVE",
	}

	json.NewEncoder(w).Encode(response)
}

// getSkillID extracts the ID from:
// /api/v1/skills/{id}
func getSkillID(r *http.Request) string {

	path := strings.TrimPrefix(
		r.URL.Path,
		"/api/v1/skills/",
	)

	return strings.TrimSpace(path)
}

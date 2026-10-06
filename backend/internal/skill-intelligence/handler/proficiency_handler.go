package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/MyHourly/TalentIQ/backend/internal/skill-intelligence/dto"
	"github.com/MyHourly/TalentIQ/backend/internal/skill-intelligence/mapper"
	"github.com/MyHourly/TalentIQ/backend/internal/skill-intelligence/service"
)

type ProficiencyHandler struct {
	service *service.ProficiencyService
}

func NewProficiencyHandler(
	service *service.ProficiencyService,
) *ProficiencyHandler {
	return &ProficiencyHandler{
		service: service,
	}
}

func (h *ProficiencyHandler) CreateProficiency(
	w http.ResponseWriter,
	r *http.Request,
) {

	var request dto.CreateProficiencyRequest

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	if strings.TrimSpace(request.Name) == "" {
		http.Error(
			w,
			"name is required",
			http.StatusBadRequest,
		)
		return
	}

	if request.LevelOrder < 1 {
		http.Error(
			w,
			"level_order must be greater than 0",
			http.StatusBadRequest,
		)
		return
	}

	proficiency := mapper.ToProficiencyModel(request)

	proficiency.ID = uuid.New().String()

	err = h.service.CreateProficiency(
		r.Context(),
		proficiency,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	response := mapper.ToProficiencyResponse(
		proficiency,
	)

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(response)
}

func (h *ProficiencyHandler) GetAllProficiencies(
	w http.ResponseWriter,
	r *http.Request,
) {

	proficiencies, err := h.service.GetAllProficiencies(
		r.Context(),
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	responses := mapper.ToProficiencyResponseList(
		proficiencies,
	)

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(responses)
}

func (h *ProficiencyHandler) GetProficiencyByID(
	w http.ResponseWriter,
	r *http.Request,
) {

	id := strings.TrimPrefix(
		r.URL.Path,
		"/api/v1/proficiencies/",
	)

	if _, err := uuid.Parse(id); err != nil {
		http.Error(
			w,
			"invalid proficiency ID",
			http.StatusBadRequest,
		)
		return
	}

	proficiency, err := h.service.GetProficiencyByID(
		r.Context(),
		id,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(
				w,
				"proficiency not found",
				http.StatusNotFound,
			)
			return
		}

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	response := mapper.ToProficiencyResponse(
		proficiency,
	)

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(response)
}

func (h *ProficiencyHandler) AssignProficiencyToSkill(
	w http.ResponseWriter,
	r *http.Request,
) {

	skillID := strings.TrimPrefix(
		r.URL.Path,
		"/api/v1/skills/",
	)

	skillID = strings.TrimSuffix(
		skillID,
		"/proficiency",
	)

	if _, err := uuid.Parse(skillID); err != nil {
		http.Error(
			w,
			"invalid skill ID",
			http.StatusBadRequest,
		)
		return
	}

	var request dto.AssignProficiencyRequest

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	if _, err := uuid.Parse(request.ProficiencyID); err != nil {
		http.Error(
			w,
			"invalid proficiency ID",
			http.StatusBadRequest,
		)
		return
	}

	err = h.service.AssignProficiencyToSkill(
		r.Context(),
		skillID,
		request.ProficiencyID,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		map[string]string{
			"message": "proficiency assigned to skill successfully",
		},
	)
}

func (h *ProficiencyHandler) GetProficiencyBySkillID(
	w http.ResponseWriter,
	r *http.Request,
) {

	skillID := strings.TrimPrefix(
		r.URL.Path,
		"/api/v1/skills/",
	)

	skillID = strings.TrimSuffix(
		skillID,
		"/proficiency",
	)

	if _, err := uuid.Parse(skillID); err != nil {
		http.Error(
			w,
			"invalid skill ID",
			http.StatusBadRequest,
		)
		return
	}

	proficiency, err := h.service.GetProficiencyBySkillID(
		r.Context(),
		skillID,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(
				w,
				"proficiency not assigned to this skill",
				http.StatusNotFound,
			)
			return
		}

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	response := mapper.ToProficiencyResponse(
		proficiency,
	)

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(response)
}

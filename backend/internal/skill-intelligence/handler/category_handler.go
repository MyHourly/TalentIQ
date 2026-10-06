package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/MyHourly/TalentIQ/backend/internal/skill-intelligence/dto"
	"github.com/MyHourly/TalentIQ/backend/internal/skill-intelligence/mapper"
	"github.com/MyHourly/TalentIQ/backend/internal/skill-intelligence/service"
)

type CategoryHandler struct {
	service *service.CategoryService
}

func NewCategoryHandler(
	categoryService *service.CategoryService,
) *CategoryHandler {

	return &CategoryHandler{
		service: categoryService,
	}
}

// CreateCategory handles:
// POST /api/v1/categories
func (h *CategoryHandler) CreateCategory(
	w http.ResponseWriter,
	r *http.Request,
) {

	var request dto.CreateCategoryRequest

	err := json.NewDecoder(
		r.Body,
	).Decode(&request)

	if err != nil {
		http.Error(
			w,
			"Invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	category := mapper.ToCategoryModel(request)

	category.ID = uuid.New().String()
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	err = h.service.CreateCategory(
		r.Context(),
		category,
	)

	if err != nil {
		http.Error(
			w,
			"Failed to create category",
			http.StatusInternalServerError,
		)
		return
	}

	response := mapper.ToCategoryResponse(
		category,
	)

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(
		http.StatusCreated,
	)

	json.NewEncoder(w).Encode(
		response,
	)
}

// GetAllCategories handles:
// GET /api/v1/categories
func (h *CategoryHandler) GetAllCategories(
	w http.ResponseWriter,
	r *http.Request,
) {

	categories, err := h.service.GetAllCategories(
		r.Context(),
	)

	if err != nil {
		http.Error(
			w,
			"Failed to get categories",
			http.StatusInternalServerError,
		)
		return
	}

	responses := mapper.ToCategoryResponseList(
		categories,
	)

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(
		responses,
	)
}

// MapSkillToCategory handles:
// POST /api/v1/skills/{skillID}/categories
func (h *CategoryHandler) MapSkillToCategory(
	w http.ResponseWriter,
	r *http.Request,
) {

	skillID := getSkillIDFromCategoryPath(r)

	if skillID == "" {
		http.Error(
			w,
			"Skill ID is required",
			http.StatusBadRequest,
		)
		return
	}

	var request struct {
		CategoryID string `json:"category_id"`
	}

	err := json.NewDecoder(
		r.Body,
	).Decode(&request)

	if err != nil {
		http.Error(
			w,
			"Invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	if request.CategoryID == "" {
		http.Error(
			w,
			"category_id is required",
			http.StatusBadRequest,
		)
		return
	}

	err = h.service.MapSkillToCategory(
		r.Context(),
		skillID,
		request.CategoryID,
	)

	if err != nil {
		http.Error(
			w,
			"Failed to map skill to category",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusCreated)

	response := map[string]string{
		"message":     "Skill mapped to category successfully",
		"skill_id":    skillID,
		"category_id": request.CategoryID,
	}

	json.NewEncoder(w).Encode(
		response,
	)
}

// GetCategoriesBySkillID handles:
// GET /api/v1/skills/{skillID}/categories
func (h *CategoryHandler) GetCategoriesBySkillID(
	w http.ResponseWriter,
	r *http.Request,
) {

	skillID := getSkillIDFromCategoryPath(r)

	if skillID == "" {
		http.Error(
			w,
			"Skill ID is required",
			http.StatusBadRequest,
		)
		return
	}

	categories, err :=
		h.service.GetCategoriesBySkillID(
			r.Context(),
			skillID,
		)

	if err != nil {
		http.Error(
			w,
			"Failed to get skill categories",
			http.StatusInternalServerError,
		)
		return
	}

	responses := mapper.ToCategoryResponseList(
		categories,
	)

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(
		responses,
	)
}

func getSkillIDFromCategoryPath(
	r *http.Request,
) string {

	path := strings.TrimPrefix(
		r.URL.Path,
		"/api/v1/skills/",
	)

	path = strings.TrimSuffix(
		path,
		"/categories",
	)

	return strings.TrimSpace(path)
}

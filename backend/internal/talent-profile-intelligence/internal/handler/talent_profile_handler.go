package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"

	"talentiq/talent-profile-intelligence/internal/dto"
	"talentiq/talent-profile-intelligence/internal/repository"
	"talentiq/talent-profile-intelligence/internal/service"
)

// TalentProfileHandler handles HTTP requests related
// to Talent Profile operations.
type TalentProfileHandler struct {
	service service.TalentProfileService
	logger  *slog.Logger
}

// NewTalentProfileHandler creates a new Talent Profile handler.
func NewTalentProfileHandler(
	service service.TalentProfileService,
	logger *slog.Logger,
) *TalentProfileHandler {

	return &TalentProfileHandler{
		service: service,
		logger:  logger,
	}
}

// Create handles:
//
// POST /api/v1/talent-profiles
func (h *TalentProfileHandler) Create(c *gin.Context) {

	var request dto.TalentProfileCreateRequest

	// Decode JSON request body.
	if err := c.ShouldBindJSON(&request); err != nil {

		h.logger.Warn(
			"invalid create talent profile request",
			"error", err,
		)

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid request body",
			},
		)

		return
	}

	// Convert DTO into internal model.
	profile := request.ToModel()

	// Call business logic.
	created, err := h.service.Create(
		c.Request.Context(),
		profile,
	)

	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	// Convert model into response DTO.
	response := dto.FromTalentProfileModel(created)

	c.JSON(
		http.StatusCreated,
		response,
	)
}

// GetByID handles:
//
// GET /api/v1/talent-profiles/:id
func (h *TalentProfileHandler) GetByID(c *gin.Context) {

	id, err := uuid.Parse(
		c.Param("id"),
	)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid talent profile id",
			},
		)

		return
	}

	profile, err := h.service.GetByID(
		c.Request.Context(),
		id,
	)

	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	response := dto.FromTalentProfileModel(profile)

	c.JSON(
		http.StatusOK,
		response,
	)
}

// List handles:
//
// GET /api/v1/talent-profiles
func (h *TalentProfileHandler) List(c *gin.Context) {

	page := parsePositiveInt(
		c.Query("page"),
		1,
	)

	limit := parsePositiveInt(
		c.Query("limit"),
		20,
	)

	// Protect the API from unnecessarily large requests.
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	profiles, total, err := h.service.List(
		c.Request.Context(),
		limit,
		offset,
	)

	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	responses := make(
		[]*dto.TalentProfileResponse,
		0,
		len(profiles),
	)

	for _, profile := range profiles {
		responses = append(
			responses,
			dto.FromTalentProfileModel(profile),
		)
	}

	response := dto.TalentProfileListResponse{
		Data:  responses,
		Page:  page,
		Limit: limit,
		Total: total,
	}

	c.JSON(
		http.StatusOK,
		response,
	)
}

// Update handles:
//
// PUT /api/v1/talent-profiles/:id
func (h *TalentProfileHandler) Update(c *gin.Context) {

	id, err := uuid.Parse(
		c.Param("id"),
	)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid talent profile id",
			},
		)

		return
	}

	var request dto.TalentProfileUpdateRequest

	if err := c.ShouldBindJSON(&request); err != nil {

		h.logger.Warn(
			"invalid update talent profile request",
			"error", err,
		)

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid request body",
			},
		)

		return
	}

	// First retrieve the existing profile.
	//
	// We need employee_code because it is not changed
	// by the update request.
	existing, err := h.service.GetByID(
		c.Request.Context(),
		id,
	)

	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	profile := request.ToModel(
		id,
		existing.EmployeeCode,
	)

	updated, err := h.service.Update(
		c.Request.Context(),
		profile,
	)

	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	response := dto.FromTalentProfileModel(updated)

	c.JSON(
		http.StatusOK,
		response,
	)
}

// Delete handles:
//
// DELETE /api/v1/talent-profiles/:id
func (h *TalentProfileHandler) Delete(c *gin.Context) {

	id, err := uuid.Parse(
		c.Param("id"),
	)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid talent profile id",
			},
		)

		return
	}

	err = h.service.Delete(
		c.Request.Context(),
		id,
	)

	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "talent profile deactivated successfully",
		},
	)
}

// handleServiceError converts known service errors
// into appropriate HTTP responses.
func (h *TalentProfileHandler) handleServiceError(
	c *gin.Context,
	err error,
) {

	switch {
	case errors.Is(
		err,
		service.ErrInvalidTalentProfile,
	):
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

	case errors.Is(
		err,
		service.ErrEmployeeCodeExists,
	):
		c.JSON(
			http.StatusConflict,
			gin.H{
				"error": err.Error(),
			},
		)

	case errors.Is(
		err,
		repository.ErrTalentProfileNotFound,
	):
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error": "talent profile not found",
			},
		)

	default:
		h.logger.Error(
			"unexpected talent profile error",
			"error", err,
		)

		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "internal server error",
			},
		)
	}
}

// parsePositiveInt safely parses query parameters.
//
// Example:
//
// ?page=2
// ?limit=20
func parsePositiveInt(
	value string,
	defaultValue int,
) int {

	if value == "" {
		return defaultValue
	}

	number, err := strconv.Atoi(value)

	if err != nil || number <= 0 {
		return defaultValue
	}

	return number
}

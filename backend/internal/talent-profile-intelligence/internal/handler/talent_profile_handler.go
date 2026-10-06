package handler

import (
	"errors"
	"net/http"
	"strconv"

	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

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

	// Decode the JSON request body.
	if err := c.ShouldBindJSON(&request); err != nil {

		h.logger.Warn(
			"invalid create talent profile request",
			"error", err,
		)

		writeError(
			c,
			http.StatusBadRequest,
			"invalid request body",
		)

		return
	}

	// Convert the request DTO into the internal model.
	profile := request.ToModel()

	// Call the service layer to execute business logic.
	created, err := h.service.Create(
		c.Request.Context(),
		profile,
	)

	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	// Convert the created model into a response DTO.
	response := dto.FromTalentProfileModel(created)

	// Return a standardized success response.
	writeSuccess(
		c,
		http.StatusCreated,
		"Talent profile created successfully",
		response,
	)
}

// GetByID handles:
//
// GET /api/v1/talent-profiles/:id
func (h *TalentProfileHandler) GetByID(c *gin.Context) {

	// Convert the URL parameter into UUID.
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {

		writeError(
			c,
			http.StatusBadRequest,
			"invalid talent profile id",
		)

		return
	}

	// Retrieve the profile through the service layer.
	profile, err := h.service.GetByID(
		c.Request.Context(),
		id,
	)

	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	// Convert model into response DTO.
	response := dto.FromTalentProfileModel(profile)

	// Return a standardized success response.
	writeSuccess(
		c,
		http.StatusOK,
		"Talent profile retrieved successfully",
		response,
	)
}

// List handles:
//
// GET /api/v1/talent-profiles
func (h *TalentProfileHandler) List(c *gin.Context) {

	// Read pagination parameters.
	page := parsePositiveInt(
		c.Query("page"),
		1,
	)

	limit := parsePositiveInt(
		c.Query("limit"),
		20,
	)

	// Prevent clients from requesting an unnecessarily
	// large number of records.
	if limit > 100 {
		limit = 100
	}

	// Calculate database offset.
	offset := (page - 1) * limit

	// Retrieve profiles from the service layer.
	profiles, total, err := h.service.List(
		c.Request.Context(),
		limit,
		offset,
	)

	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	// Convert models into response DTOs.
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

	// Build the existing pagination response.
	listResponse := dto.TalentProfileListResponse{
		Data:  responses,
		Page:  page,
		Limit: limit,
		Total: total,
	}

	// Wrap the list response in the standard API response.
	writeSuccess(
		c,
		http.StatusOK,
		"Talent profiles retrieved successfully",
		listResponse,
	)
}

// Update handles:
//
// PUT /api/v1/talent-profiles/:id
func (h *TalentProfileHandler) Update(c *gin.Context) {

	// Convert the URL parameter into UUID.
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {

		writeError(
			c,
			http.StatusBadRequest,
			"invalid talent profile id",
		)

		return
	}

	var request dto.TalentProfileUpdateRequest

	// Decode the JSON request body.
	if err := c.ShouldBindJSON(&request); err != nil {

		h.logger.Warn(
			"invalid update talent profile request",
			"error", err,
		)

		writeError(
			c,
			http.StatusBadRequest,
			"invalid request body",
		)

		return
	}

	// Retrieve the existing profile first.
	//
	// EmployeeCode is not updated by the update request,
	// so we need the existing value.
	existing, err := h.service.GetByID(
		c.Request.Context(),
		id,
	)

	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	// Convert update DTO into the internal model.
	profile := request.ToModel(
		id,
		existing.EmployeeCode,
	)

	// Update the profile through the service layer.
	updated, err := h.service.Update(
		c.Request.Context(),
		profile,
	)

	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	// Convert updated model into response DTO.
	response := dto.FromTalentProfileModel(updated)

	// Return standardized success response.
	writeSuccess(
		c,
		http.StatusOK,
		"Talent profile updated successfully",
		response,
	)
}

// Delete handles:
//
// DELETE /api/v1/talent-profiles/:id
func (h *TalentProfileHandler) Delete(c *gin.Context) {

	// Convert the URL parameter into UUID.
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {

		writeError(
			c,
			http.StatusBadRequest,
			"invalid talent profile id",
		)

		return
	}

	// Deactivate the profile through the service layer.
	err = h.service.Delete(
		c.Request.Context(),
		id,
	)

	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	// Return standardized success response.
	writeSuccess(
		c,
		http.StatusOK,
		"Talent profile deactivated successfully",
		nil,
	)
}

// handleServiceError converts known service/repository errors
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

		writeError(
			c,
			http.StatusBadRequest,
			err.Error(),
		)

	case errors.Is(
		err,
		service.ErrEmployeeCodeExists,
	):

		writeError(
			c,
			http.StatusConflict,
			err.Error(),
		)

	case errors.Is(
		err,
		repository.ErrTalentProfileNotFound,
	):

		writeError(
			c,
			http.StatusNotFound,
			"talent profile not found",
		)

	default:

		// Log the actual internal error for developers/operators.
		// Do not expose the internal error to the API client.
		h.logger.Error(
			"unexpected talent profile error",
			"error", err,
		)

		writeError(
			c,
			http.StatusInternalServerError,
			"internal server error",
		)
	}
}

// writeSuccess sends a standardized successful API response.
func writeSuccess(
	c *gin.Context,
	status int,
	message string,
	data any,
) {

	c.JSON(
		status,
		dto.APIResponse{
			Success: true,
			Message: message,
			Data:    data,
		},
	)
}

// writeError sends a standardized error API response.
func writeError(
	c *gin.Context,
	status int,
	message string,
) {

	c.JSON(
		status,
		dto.ErrorResponse{
			Success: false,
			Message: message,
		},
	)
}

// parsePositiveInt safely parses positive integer query parameters.
//
// Example:
//
// ?page=2
// ?limit=20
//
// If the value is missing or invalid, the supplied default value
// is returned.
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

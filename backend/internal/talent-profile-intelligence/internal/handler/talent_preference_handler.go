package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"

	"talentiq/talent-profile-intelligence/internal/dto"
	"talentiq/talent-profile-intelligence/internal/model"
	"talentiq/talent-profile-intelligence/internal/repository"
	"talentiq/talent-profile-intelligence/internal/service"
)

// TalentPreferenceHandler handles HTTP requests
// related to talent preferences.
type TalentPreferenceHandler struct {
	service service.TalentPreferenceService
	logger  *slog.Logger
}

// NewTalentPreferenceHandler creates a new Talent Preference handler.
func NewTalentPreferenceHandler(
	service service.TalentPreferenceService,
	logger *slog.Logger,
) *TalentPreferenceHandler {
	return &TalentPreferenceHandler{
		service: service,
		logger:  logger,
	}
}

// Create creates preferences for a talent.
func (h *TalentPreferenceHandler) Create(c *gin.Context) {
	talentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid talent_id",
		})
		return
	}

	var request dto.TalentPreferenceRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	preference := &model.TalentPreference{
		TalentID:           talentID,
		PreferredRole:      request.PreferredRole,
		PreferredLocation:  request.PreferredLocation,
		AvailabilityStatus: request.AvailabilityStatus,
	}

	created, err := h.service.Create(
		c.Request.Context(),
		preference,
	)
	if err != nil {
		if errors.Is(err, service.ErrInvalidTalentPreference) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		h.logger.Error(
			"failed to create talent preference",
			"talent_id", talentID,
			"error", err,
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create talent preference",
		})
		return
	}

	c.JSON(http.StatusCreated, toTalentPreferenceResponse(created))
}

// GetByTalentID returns preferences for a talent.
func (h *TalentPreferenceHandler) GetByTalentID(c *gin.Context) {
	talentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid talent_id",
		})
		return
	}

	preference, err := h.service.GetByTalentID(
		c.Request.Context(),
		talentID,
	)
	if err != nil {
		if errors.Is(err, repository.ErrTalentPreferenceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "talent preference not found",
			})
			return
		}

		if errors.Is(err, service.ErrInvalidTalentPreference) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		h.logger.Error(
			"failed to get talent preference",
			"talent_id", talentID,
			"error", err,
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get talent preference",
		})
		return
	}

	c.JSON(http.StatusOK, toTalentPreferenceResponse(preference))
}

// Update modifies the preferences of a talent.
func (h *TalentPreferenceHandler) Update(c *gin.Context) {
	talentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid talent_id",
		})
		return
	}

	var request dto.TalentPreferenceRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	preference := &model.TalentPreference{
		TalentID:           talentID,
		PreferredRole:      request.PreferredRole,
		PreferredLocation:  request.PreferredLocation,
		AvailabilityStatus: request.AvailabilityStatus,
	}

	updated, err := h.service.Update(
		c.Request.Context(),
		preference,
	)
	if err != nil {
		if errors.Is(err, service.ErrInvalidTalentPreference) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if errors.Is(err, repository.ErrTalentPreferenceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "talent preference not found",
			})
			return
		}

		h.logger.Error(
			"failed to update talent preference",
			"talent_id", talentID,
			"error", err,
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update talent preference",
		})
		return
	}

	c.JSON(http.StatusOK, toTalentPreferenceResponse(updated))
}

// UpdateAvailability updates only the availability status
// of an existing talent preference.
func (h *TalentPreferenceHandler) UpdateAvailability(c *gin.Context) {
	talentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid talent id",
		})
		return
	}

	var request dto.TalentAvailabilityRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	updated, err := h.service.UpdateAvailability(
		c.Request.Context(),
		talentID,
		request.AvailabilityStatus,
	)
	if err != nil {
		if errors.Is(err, service.ErrInvalidTalentPreference) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid availability request",
			})
			return
		}

		if errors.Is(err, repository.ErrTalentPreferenceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "talent preference not found",
			})
			return
		}

		h.logger.Error(
			"failed to update talent availability",
			"talent_id", talentID,
			"error", err,
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update talent availability",
		})
		return
	}

	c.JSON(http.StatusOK, toTalentPreferenceResponse(updated))
}

// toTalentPreferenceResponse converts the internal model
// into the API response DTO.
func toTalentPreferenceResponse(
	preference *model.TalentPreference,
) dto.TalentPreferenceResponse {
	return dto.TalentPreferenceResponse{
		ID:                 preference.ID.String(),
		TalentID:           preference.TalentID.String(),
		PreferredRole:      preference.PreferredRole,
		PreferredLocation:  preference.PreferredLocation,
		AvailabilityStatus: preference.AvailabilityStatus,
		CreatedAt:          preference.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:          preference.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

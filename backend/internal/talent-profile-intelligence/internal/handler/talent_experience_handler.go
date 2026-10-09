package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"talentiq/talent-profile-intelligence/internal/dto"
	"talentiq/talent-profile-intelligence/internal/repository"
	"talentiq/talent-profile-intelligence/internal/service"
)

type TalentExperienceHandler struct {
	service service.TalentExperienceService
	logger  *slog.Logger
}

func NewTalentExperienceHandler(svc service.TalentExperienceService, logger *slog.Logger) *TalentExperienceHandler {
	return &TalentExperienceHandler{service: svc, logger: logger}
}

// POST /api/v1/talents/:id/experience  -> add a job to a talent's work history
func (h *TalentExperienceHandler) Add(c *gin.Context) {
	talentID, err := uuid.Parse(c.Param("id")) // 1. read the talent id from the URL
	if err != nil {
		writeExperienceError(c, http.StatusBadRequest, "invalid talent id")
		return
	}

	var req dto.TalentExperienceRequest
	if err := c.ShouldBindJSON(&req); err != nil { // 2. read the JSON body
		writeExperienceError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	saved, err := h.service.Add(c.Request.Context(), talentID, &req) // 3. do the work
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{ // 4. reply (201 = something new was created)
		"success": true,
		"message": "talent experience added successfully",
		"data":    dto.NewTalentExperienceResponse(saved),
	})
}

// Converts "what went wrong" into the right web status code:
//
//	bad input -> 400, talent missing -> 404, already added -> 409, anything else -> 500
func (h *TalentExperienceHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, repository.ErrInvalidTalentExperience):
		writeExperienceError(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, repository.ErrTalentNotFoundForExperience):
		writeExperienceError(c, http.StatusNotFound, "talent not found")
	case errors.Is(err, repository.ErrExperienceAlreadyExists):
		writeExperienceError(c, http.StatusConflict, "experience already added for this talent")
	default:
		h.logger.Error("talent experience request failed", "error", err)
		writeExperienceError(c, http.StatusInternalServerError, "internal server error")
	}
}

// Same error reply style as the other talent APIs.
func writeExperienceError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"success": false,
		"message": message,
	})
}

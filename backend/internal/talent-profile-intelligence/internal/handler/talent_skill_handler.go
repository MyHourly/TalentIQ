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

type TalentSkillHandler struct {
	service service.TalentSkillService
	logger  *slog.Logger
}

// Takes the logger as the second argument, like your other handlers.
func NewTalentSkillHandler(svc service.TalentSkillService, logger *slog.Logger) *TalentSkillHandler {
	return &TalentSkillHandler{service: svc, logger: logger}
}

// PUT /api/v1/talents/:id/skills/:skillId  -> add or update a skill
func (h *TalentSkillHandler) Upsert(c *gin.Context) {
	talentID, skillID, ok := parseTalentAndSkillID(c) // 1. read ids from the URL
	if !ok {
		return // a 400 reply was already sent
	}

	var req dto.TalentSkillRequest
	if err := c.ShouldBindJSON(&req); err != nil { // 2. read the JSON body
		writeSkillError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	saved, err := h.service.Upsert(c.Request.Context(), talentID, skillID, req.ProficiencyLevel) // 3. do the work
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{ // 4. reply
		"success": true,
		"message": "talent skill saved successfully",
		"data": dto.NewTalentSkillResponse(
			saved.ID.String(),
			saved.TalentID.String(),
			saved.SkillID.String(),
			saved.ProficiencyLevel,
			saved.CreatedAt,
			saved.UpdatedAt,
		),
	})
}

// DELETE /api/v1/talents/:id/skills/:skillId  -> remove a skill
func (h *TalentSkillHandler) Remove(c *gin.Context) {
	talentID, skillID, ok := parseTalentAndSkillID(c)
	if !ok {
		return
	}

	if err := h.service.Remove(c.Request.Context(), talentID, skillID); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "talent skill removed successfully",
		"data":    nil,
	})
}

// Converts "what went wrong" into the right web status code:
//
//	bad input -> 400, not found -> 404, anything unexpected -> 500
func (h *TalentSkillHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, repository.ErrInvalidTalentSkill):
		writeSkillError(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, repository.ErrTalentNotFoundForSkill):
		writeSkillError(c, http.StatusNotFound, "talent not found")
	case errors.Is(err, repository.ErrTalentSkillNotFound):
		writeSkillError(c, http.StatusNotFound, "talent skill not found")
	default:
		// Unknown problem: write the details to our log, but show the caller a safe generic message.
		h.logger.Error("talent skill request failed", "error", err)
		writeSkillError(c, http.StatusInternalServerError, "internal server error")
	}
}

// Reads the talent id and skill id from the URL and checks they are valid UUIDs.
func parseTalentAndSkillID(c *gin.Context) (talentID, skillID uuid.UUID, ok bool) {
	talentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		writeSkillError(c, http.StatusBadRequest, "invalid talent id")
		return uuid.Nil, uuid.Nil, false
	}

	skillID, err = uuid.Parse(c.Param("skillId"))
	if err != nil {
		writeSkillError(c, http.StatusBadRequest, "invalid skill id")
		return uuid.Nil, uuid.Nil, false
	}

	return talentID, skillID, true
}

// Sends an error reply in the same style as your talent profile APIs.
// If you already have a shared error helper from TPI-009, use that instead of this.
func writeSkillError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"success": false,
		"message": message,
	})
}
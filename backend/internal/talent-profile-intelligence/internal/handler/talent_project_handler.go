package handler

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"talentiq/talent-profile-intelligence/internal/dto"
	"talentiq/talent-profile-intelligence/internal/repository"
	"talentiq/talent-profile-intelligence/internal/service"
)

type TalentProjectHandler struct {
	service service.TalentProjectService
	logger  *slog.Logger
}

func NewTalentProjectHandler(svc service.TalentProjectService, logger *slog.Logger) *TalentProjectHandler {
	return &TalentProjectHandler{service: svc, logger: logger}
}

// POST /api/v1/talents/:id/projects  -> add a project to a talent
func (h *TalentProjectHandler) Add(c *gin.Context) {
	talentID, err := uuid.Parse(c.Param("id")) // 1. read the talent id from the URL
	if err != nil {
		writeProjectError(c, http.StatusBadRequest, "invalid talent id")
		return
	}

	var req dto.TalentProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil { // 2. read the JSON body
		writeProjectError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	saved, err := h.service.Add(c.Request.Context(), talentID, &req) // 3. do the work
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{ // 4. reply (201 = something new was created)
		"success": true,
		"message": "talent project added successfully",
		"data":    dto.NewTalentProjectResponse(saved),
	})
}

// POST /api/v1/talents/:id/projects/:projectId/complete  -> mark a project completed
func (h *TalentProjectHandler) Complete(c *gin.Context) {
	talentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		writeProjectError(c, http.StatusBadRequest, "invalid talent id")
		return
	}

	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		writeProjectError(c, http.StatusBadRequest, "invalid project id")
		return
	}

	// The body is optional (end_date defaults to today),
	// so an EMPTY body is fine. Only a BROKEN body is an error.
	var req dto.ProjectCompletionRequest
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		writeProjectError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	done, err := h.service.Complete(c.Request.Context(), talentID, projectID, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "talent project completed successfully",
		"data":    dto.NewTalentProjectResponse(done),
	})
}

// Converts "what went wrong" into the right web status code:
//
//	bad input / already completed -> 400, not found -> 404, anything else -> 500
func (h *TalentProjectHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, repository.ErrInvalidTalentProject):
		writeProjectError(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, repository.ErrProjectAlreadyCompleted):
		writeProjectError(c, http.StatusBadRequest, "project is already completed")
	case errors.Is(err, repository.ErrTalentNotFoundForProject):
		writeProjectError(c, http.StatusNotFound, "talent not found")
	case errors.Is(err, repository.ErrTalentProjectNotFound):
		writeProjectError(c, http.StatusNotFound, "project not found for this talent")
	default:
		h.logger.Error("talent project request failed", "error", err)
		writeProjectError(c, http.StatusInternalServerError, "internal server error")
	}
}

// Same error reply style as the other talent APIs.
func writeProjectError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"success": false,
		"message": message,
	})
}

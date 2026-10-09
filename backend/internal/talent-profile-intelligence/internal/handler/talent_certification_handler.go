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

type TalentCertificationHandler struct {
	service service.TalentCertificationService
	logger  *slog.Logger
}

func NewTalentCertificationHandler(svc service.TalentCertificationService, logger *slog.Logger) *TalentCertificationHandler {
	return &TalentCertificationHandler{service: svc, logger: logger}
}

// POST /api/v1/talents/:id/certifications  -> add a certification to a talent
func (h *TalentCertificationHandler) Add(c *gin.Context) {
	talentID, err := uuid.Parse(c.Param("id")) // 1. read the talent id from the URL
	if err != nil {
		writeCertificationError(c, http.StatusBadRequest, "invalid talent id")
		return
	}

	var req dto.TalentCertificationRequest
	if err := c.ShouldBindJSON(&req); err != nil { // 2. read the JSON body
		writeCertificationError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	saved, err := h.service.Add(c.Request.Context(), talentID, &req) // 3. do the work
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{ // 4. reply (201 = something new was created)
		"success": true,
		"message": "talent certification added successfully",
		"data":    dto.NewTalentCertificationResponse(saved),
	})
}

// Converts "what went wrong" into the right web status code:
//
//	bad input -> 400, talent missing -> 404, already added -> 409, anything else -> 500
func (h *TalentCertificationHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, repository.ErrInvalidTalentCertification):
		writeCertificationError(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, repository.ErrTalentNotFoundForCertification):
		writeCertificationError(c, http.StatusNotFound, "talent not found")
	case errors.Is(err, repository.ErrCertificationAlreadyExists):
		writeCertificationError(c, http.StatusConflict, "certification already added for this talent")
	default:
		h.logger.Error("talent certification request failed", "error", err)
		writeCertificationError(c, http.StatusInternalServerError, "internal server error")
	}
}

// Same error reply style as the skill and profile APIs.
func writeCertificationError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"success": false,
		"message": message,
	})
}

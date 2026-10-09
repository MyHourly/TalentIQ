package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"

	"talentiq/talent-profile-intelligence/internal/dto"
	"talentiq/talent-profile-intelligence/internal/model"
	"talentiq/talent-profile-intelligence/internal/repository"
)

const projectDateLayout = "2006-01-02" // dates look like 2026-10-08

const (
	projectStatusInProgress = "IN_PROGRESS"
	projectStatusCompleted  = "COMPLETED"
)

// Longest allowed text, matching the database column sizes.
const (
	maxProjectNameLength   = 200
	maxProjectDomainLength = 100
	maxProjectClientLength = 150
	maxProjectRoleLength   = 100
)

type TalentProjectService interface {
	Add(ctx context.Context, talentID uuid.UUID, req *dto.TalentProjectRequest) (*model.TalentProject, error)
	Complete(ctx context.Context, talentID, projectID uuid.UUID, req *dto.ProjectCompletionRequest) (*model.TalentProject, error)
}

type talentProjectService struct {
	repo   repository.TalentProjectRepository
	logger *slog.Logger
}

func NewTalentProjectService(repo repository.TalentProjectRepository, logger *slog.Logger) TalentProjectService {
	return &talentProjectService{repo: repo, logger: logger}
}

// Builds a "bad input" error with a readable reason.
func invalidProject(format string, args ...any) error {
	return fmt.Errorf("%w: %s", repository.ErrInvalidTalentProject, fmt.Sprintf(format, args...))
}

// Returns an error if the text is longer than the database allows.
func checkProjectLength(field, value string, max int) error {
	if utf8.RuneCountInString(value) > max {
		return invalidProject("%s must be at most %d characters", field, max)
	}
	return nil
}

// Add creates a project for a talent. Rules, in order:
//  1. talent id is real, name is given, start_date is a valid date
//  2. texts are not longer than the database allows
//  3. the person must exist and must not be deactivated
//
// A new project always starts as IN_PROGRESS with no end date.
func (s *talentProjectService) Add(ctx context.Context, talentID uuid.UUID, req *dto.TalentProjectRequest) (*model.TalentProject, error) {
	if req == nil {
		return nil, invalidProject("request body is required")
	}
	if talentID == uuid.Nil {
		return nil, invalidProject("talent id is required")
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, invalidProject("name is required")
	}
	domain := strings.TrimSpace(req.Domain)
	clientName := strings.TrimSpace(req.ClientName)
	role := strings.TrimSpace(req.Role)

	for _, c := range []struct {
		field string
		value string
		max   int
	}{
		{"name", name, maxProjectNameLength},
		{"domain", domain, maxProjectDomainLength},
		{"client_name", clientName, maxProjectClientLength},
		{"role", role, maxProjectRoleLength},
	} {
		if err := checkProjectLength(c.field, c.value, c.max); err != nil {
			return nil, err
		}
	}

	if strings.TrimSpace(req.StartDate) == "" {
		return nil, invalidProject("start_date is required")
	}
	startDate, err := time.Parse(projectDateLayout, strings.TrimSpace(req.StartDate))
	if err != nil {
		return nil, invalidProject("start_date must be a date like 2026-10-08")
	}

	exists, err := s.repo.TalentExists(ctx, talentID)
	if err != nil {
		return nil, err
	}
	if !exists {
		s.logger.Warn("cannot add project: talent not found", "talent_id", talentID)
		return nil, repository.ErrTalentNotFoundForProject
	}

	return s.repo.Create(ctx, &model.TalentProject{
		TalentID:    talentID,
		Name:        name,
		Description: strings.TrimSpace(req.Description),
		Domain:      domain,
		ClientName:  clientName,
		Role:        role,
		StartDate:   startDate,
		Status:      projectStatusInProgress,
	})
}

// Complete marks a project as finished. Rules, in order:
//  1. ids are real; end_date (optional) is a valid date, default = today
//  2. the project must exist and belong to this talent  -> otherwise 404
//  3. it must not be completed already                  -> otherwise 400
//  4. end_date must not be before the project's start_date -> otherwise 400
func (s *talentProjectService) Complete(ctx context.Context, talentID, projectID uuid.UUID, req *dto.ProjectCompletionRequest) (*model.TalentProject, error) {
	if talentID == uuid.Nil {
		return nil, invalidProject("talent id is required")
	}
	if projectID == uuid.Nil {
		return nil, invalidProject("project id is required")
	}

	// Default completion date is today (date only, no time).
	endDate := time.Now().UTC().Truncate(24 * time.Hour)
	if req != nil && req.EndDate != nil && strings.TrimSpace(*req.EndDate) != "" {
		parsed, err := time.Parse(projectDateLayout, strings.TrimSpace(*req.EndDate))
		if err != nil {
			return nil, invalidProject("end_date must be a date like 2026-10-08")
		}
		endDate = parsed
	}

	project, err := s.repo.GetByTalentAndProjectID(ctx, talentID, projectID)
	if err != nil {
		return nil, err
	}

	if project.Status == projectStatusCompleted {
		return nil, repository.ErrProjectAlreadyCompleted
	}
	if endDate.Before(project.StartDate) {
		return nil, invalidProject("end_date cannot be before start_date")
	}

	return s.repo.Complete(ctx, talentID, projectID, endDate)
}

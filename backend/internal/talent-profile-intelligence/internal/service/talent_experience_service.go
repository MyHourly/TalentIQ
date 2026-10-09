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

const experienceDateLayout = "2006-01-02" // dates look like 2026-10-08

// Longest allowed text, matching the database column sizes.
const (
	maxExperienceCompanyLength = 200
	maxExperienceTitleLength   = 150
)

type TalentExperienceService interface {
	Add(ctx context.Context, talentID uuid.UUID, req *dto.TalentExperienceRequest) (*model.TalentExperience, error)
}

type talentExperienceService struct {
	repo   repository.TalentExperienceRepository
	logger *slog.Logger
}

func NewTalentExperienceService(repo repository.TalentExperienceRepository, logger *slog.Logger) TalentExperienceService {
	return &talentExperienceService{repo: repo, logger: logger}
}

// Builds a "bad input" error with a readable reason.
func invalidExperience(format string, args ...any) error {
	return fmt.Errorf("%w: %s", repository.ErrInvalidTalentExperience, fmt.Sprintf(format, args...))
}

// Returns an error if the text is longer than the database allows.
func checkExperienceLength(field, value string, max int) error {
	if utf8.RuneCountInString(value) > max {
		return invalidExperience("%s must be at most %d characters", field, max)
	}
	return nil
}

// Add saves one job in a talent's work history. Rules, in order:
//  1. talent id is real; company_name and job_title are given and not too long
//  2. start_date is required and must look like 2026-10-08
//  3. end_date (optional) must look the same and not be before start_date
//  4. the person must exist and must not be deactivated
//
// No end_date means the person still works there.
func (s *talentExperienceService) Add(ctx context.Context, talentID uuid.UUID, req *dto.TalentExperienceRequest) (*model.TalentExperience, error) {
	if req == nil {
		return nil, invalidExperience("request body is required")
	}
	if talentID == uuid.Nil {
		return nil, invalidExperience("talent id is required")
	}

	company := strings.TrimSpace(req.CompanyName)
	if company == "" {
		return nil, invalidExperience("company_name is required")
	}
	if err := checkExperienceLength("company_name", company, maxExperienceCompanyLength); err != nil {
		return nil, err
	}

	title := strings.TrimSpace(req.JobTitle)
	if title == "" {
		return nil, invalidExperience("job_title is required")
	}
	if err := checkExperienceLength("job_title", title, maxExperienceTitleLength); err != nil {
		return nil, err
	}

	if strings.TrimSpace(req.StartDate) == "" {
		return nil, invalidExperience("start_date is required")
	}
	startDate, err := time.Parse(experienceDateLayout, strings.TrimSpace(req.StartDate))
	if err != nil {
		return nil, invalidExperience("start_date must be a date like 2026-10-08")
	}

	var endDate *time.Time
	if req.EndDate != nil && strings.TrimSpace(*req.EndDate) != "" {
		parsed, err := time.Parse(experienceDateLayout, strings.TrimSpace(*req.EndDate))
		if err != nil {
			return nil, invalidExperience("end_date must be a date like 2026-10-08")
		}
		if parsed.Before(startDate) {
			return nil, invalidExperience("end_date cannot be before start_date")
		}
		endDate = &parsed
	}

	exists, err := s.repo.TalentExists(ctx, talentID)
	if err != nil {
		return nil, err
	}
	if !exists {
		s.logger.Warn("cannot add experience: talent not found", "talent_id", talentID)
		return nil, repository.ErrTalentNotFoundForExperience
	}

	return s.repo.Create(ctx, &model.TalentExperience{
		TalentID:    talentID,
		CompanyName: company,
		JobTitle:    title,
		StartDate:   startDate,
		EndDate:     endDate,
		Description: strings.TrimSpace(req.Description),
	})
}

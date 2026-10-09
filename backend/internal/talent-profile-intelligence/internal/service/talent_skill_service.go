package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"talentiq/talent-profile-intelligence/internal/model"
	"talentiq/talent-profile-intelligence/internal/repository"
)

const (
	minProficiencyLevel = 1 // Beginner
	maxProficiencyLevel = 5 // Master
)

type TalentSkillService interface {
	Upsert(ctx context.Context, talentID, skillID uuid.UUID, proficiencyLevel int) (*model.TalentSkill, error)
	Remove(ctx context.Context, talentID, skillID uuid.UUID) error
}

type talentSkillService struct {
	repo   repository.TalentSkillRepository
	logger *slog.Logger
}

// Takes the logger as the second argument, like your other services.
func NewTalentSkillService(repo repository.TalentSkillRepository, logger *slog.Logger) TalentSkillService {
	return &talentSkillService{repo: repo, logger: logger}
}

// Add or update a skill. Rules, in order:
//  1. ids must be real (not empty)
//  2. level must be between 1 and 5
//  3. the person must exist and must not be deactivated
func (s *talentSkillService) Upsert(ctx context.Context, talentID, skillID uuid.UUID, proficiencyLevel int) (*model.TalentSkill, error) {
	if talentID == uuid.Nil {
		return nil, fmt.Errorf("%w: talent id is required", repository.ErrInvalidTalentSkill)
	}
	if skillID == uuid.Nil {
		return nil, fmt.Errorf("%w: skill id is required", repository.ErrInvalidTalentSkill)
	}
	if proficiencyLevel < minProficiencyLevel || proficiencyLevel > maxProficiencyLevel {
		return nil, fmt.Errorf("%w: proficiency_level must be between %d and %d",
			repository.ErrInvalidTalentSkill, minProficiencyLevel, maxProficiencyLevel)
	}

	exists, err := s.repo.TalentExists(ctx, talentID)
	if err != nil {
		return nil, err
	}
	if !exists {
		s.logger.Warn("cannot save skill: talent not found", "talent_id", talentID)
		return nil, repository.ErrTalentNotFoundForSkill
	}

	// TODO: later, also check that skillID really exists in the Skill Intelligence module.

	return s.repo.Upsert(ctx, &model.TalentSkill{
		TalentID:         talentID,
		SkillID:          skillID,
		ProficiencyLevel: proficiencyLevel,
	})
}

// Remove a skill. If the person or the skill row is missing, the repository says "not found".
func (s *talentSkillService) Remove(ctx context.Context, talentID, skillID uuid.UUID) error {
	if talentID == uuid.Nil {
		return fmt.Errorf("%w: talent id is required", repository.ErrInvalidTalentSkill)
	}
	if skillID == uuid.Nil {
		return fmt.Errorf("%w: skill id is required", repository.ErrInvalidTalentSkill)
	}

	return s.repo.Delete(ctx, talentID, skillID)
}
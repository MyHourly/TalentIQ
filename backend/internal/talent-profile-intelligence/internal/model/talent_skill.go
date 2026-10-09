package model

import (
	"time"

	"github.com/google/uuid"
)

// TalentSkill is one row of the talent_skills table.
// In plain words: "this person has this skill at this level".
type TalentSkill struct {
	ID               uuid.UUID
	TalentID         uuid.UUID
	SkillID          uuid.UUID
	ProficiencyLevel int // 1 (Beginner) to 5 (Master)
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
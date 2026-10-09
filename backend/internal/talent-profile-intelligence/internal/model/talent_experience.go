package model

import (
	"time"

	"github.com/google/uuid"
)

// TalentExperience is one row of the professional_experience table:
// "this person worked at this company in this job".
type TalentExperience struct {
	ID          uuid.UUID
	TalentID    uuid.UUID
	CompanyName string
	JobTitle    string
	StartDate   time.Time
	EndDate     *time.Time // nil = still working there
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

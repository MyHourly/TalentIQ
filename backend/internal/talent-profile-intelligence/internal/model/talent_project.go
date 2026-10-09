package model

import (
	"time"

	"github.com/google/uuid"
)

// TalentProject is "a project, as seen from one talent".
// It combines a row of "projects" with that talent's row in "project_members"
// (the role they played), so one struct is enough for the API.
type TalentProject struct {
	ID          uuid.UUID // project id
	TalentID    uuid.UUID
	Name        string
	Description string
	Domain      string
	ClientName  string
	Role        string // the talent's role in this project
	StartDate   time.Time
	EndDate     *time.Time // nil until the project is completed
	Status      string     // IN_PROGRESS or COMPLETED
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

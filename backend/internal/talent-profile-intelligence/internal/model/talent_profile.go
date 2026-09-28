package model

import (
	"time"

	"github.com/google/uuid"
)

// TalentProfile represents the core profile information
// of a talent inside TalentIQ.
type TalentProfile struct {
	ID uuid.UUID `json:"id"`

	// Employee information.
	EmployeeCode string `json:"employee_code"`

	// Personal information.
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`

	// Contact information.
	Email string `json:"email"`
	Phone string `json:"phone"`

	// Professional information.
	Designation string `json:"designation"`
	Department  string `json:"department"`
	Location    string `json:"location"`

	// Short professional summary.
	Summary string `json:"summary"`

	// Total professional experience in years.
	TotalExperienceYears float64 `json:"total_experience_years"`

	// Profile lifecycle status.
	ProfileStatus string `json:"profile_status"`

	// Audit information.
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// DeletedAt is nullable because profiles use
	// soft deletion instead of immediate database deletion.
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
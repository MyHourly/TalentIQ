package model

import (
	"time"

	"github.com/google/uuid"
)

// TalentPreference represents a talent's current work preferences.
//
// This model matches the talent_preferences table in PostgreSQL.
// It only stores data; business rules belong in the service layer.
type TalentPreference struct {
	ID                uuid.UUID  `json:"id"`
	TalentID          uuid.UUID  `json:"talent_id"`
	PreferredRole     string     `json:"preferred_role"`
	PreferredLocation string     `json:"preferred_location"`
	AvailabilityStatus string    `json:"availability_status"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}
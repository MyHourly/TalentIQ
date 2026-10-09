package model

import (
	"time"

	"github.com/google/uuid"
)

// TalentCertification is one row of the talent_certifications table:
// "this person earned this certification on this date".
type TalentCertification struct {
	ID              uuid.UUID
	TalentID        uuid.UUID
	CertificationID uuid.UUID
	IssuedAt        time.Time
	ExpiresAt       *time.Time // nil = never expires
	CredentialURL   *string    // nil = no link given
	Status          string     // ACTIVE or EXPIRED
	CreatedAt       time.Time
}

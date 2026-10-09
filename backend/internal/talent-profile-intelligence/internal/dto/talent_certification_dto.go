package dto

import (
	"time"

	"talentiq/talent-profile-intelligence/internal/model"
)

const certificationDateLayout = "2006-01-02" // dates look like 2026-10-08

// What the caller sends. The talent id comes from the URL.
// expires_at and credential_url are optional (pointers can be empty).
type TalentCertificationRequest struct {
	CertificationID string  `json:"certification_id"`
	IssuedAt        string  `json:"issued_at"`
	ExpiresAt       *string `json:"expires_at"`
	CredentialURL   *string `json:"credential_url"`
}

// What we send back after a successful add.
type TalentCertificationResponse struct {
	ID              string  `json:"id"`
	TalentID        string  `json:"talent_id"`
	CertificationID string  `json:"certification_id"`
	IssuedAt        string  `json:"issued_at"`
	ExpiresAt       *string `json:"expires_at"`
	CredentialURL   *string `json:"credential_url"`
	Status          string  `json:"status"`
	CreatedAt       string  `json:"created_at"`
}

// Turns a database model into the reply shape (ids become text, dates become 2026-10-08).
func NewTalentCertificationResponse(c *model.TalentCertification) TalentCertificationResponse {
	var expires *string
	if c.ExpiresAt != nil {
		s := c.ExpiresAt.Format(certificationDateLayout)
		expires = &s
	}

	return TalentCertificationResponse{
		ID:              c.ID.String(),
		TalentID:        c.TalentID.String(),
		CertificationID: c.CertificationID.String(),
		IssuedAt:        c.IssuedAt.Format(certificationDateLayout),
		ExpiresAt:       expires,
		CredentialURL:   c.CredentialURL,
		Status:          c.Status,
		CreatedAt:       c.CreatedAt.Format(time.RFC3339),
	}
}

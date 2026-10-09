package service

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"

	"talentiq/talent-profile-intelligence/internal/dto"
	"talentiq/talent-profile-intelligence/internal/model"
	"talentiq/talent-profile-intelligence/internal/repository"
)

const certificationDateLayout = "2006-01-02" // dates look like 2026-10-08

const (
	certificationStatusActive  = "ACTIVE"
	certificationStatusExpired = "EXPIRED"
)

type TalentCertificationService interface {
	Add(ctx context.Context, talentID uuid.UUID, req *dto.TalentCertificationRequest) (*model.TalentCertification, error)
}

type talentCertificationService struct {
	repo   repository.TalentCertificationRepository
	logger *slog.Logger
}

func NewTalentCertificationService(repo repository.TalentCertificationRepository, logger *slog.Logger) TalentCertificationService {
	return &talentCertificationService{repo: repo, logger: logger}
}

// Builds a "bad input" error with a readable reason.
func invalidCertification(format string, args ...any) error {
	return fmt.Errorf("%w: %s", repository.ErrInvalidTalentCertification, fmt.Sprintf(format, args...))
}

// Add saves a certification for a talent. Rules, in order:
//  1. talent id and certification id must be real ids
//  2. issued_at is required and must look like 2026-10-08
//  3. expires_at (optional) must look the same and not be before issued_at
//  4. credential_url (optional) must be a proper http/https link
//  5. the person must exist and must not be deactivated
//
// The status is decided here, not by the caller:
// EXPIRED if expires_at is already in the past, otherwise ACTIVE.
func (s *talentCertificationService) Add(ctx context.Context, talentID uuid.UUID, req *dto.TalentCertificationRequest) (*model.TalentCertification, error) {
	if req == nil {
		return nil, invalidCertification("request body is required")
	}
	if talentID == uuid.Nil {
		return nil, invalidCertification("talent id is required")
	}

	certificationID, err := uuid.Parse(strings.TrimSpace(req.CertificationID))
	if err != nil || certificationID == uuid.Nil {
		return nil, invalidCertification("certification_id must be a valid UUID")
	}

	if strings.TrimSpace(req.IssuedAt) == "" {
		return nil, invalidCertification("issued_at is required")
	}
	issuedAt, err := time.Parse(certificationDateLayout, strings.TrimSpace(req.IssuedAt))
	if err != nil {
		return nil, invalidCertification("issued_at must be a date like 2026-10-08")
	}

	var expiresAt *time.Time
	if req.ExpiresAt != nil && strings.TrimSpace(*req.ExpiresAt) != "" {
		parsed, err := time.Parse(certificationDateLayout, strings.TrimSpace(*req.ExpiresAt))
		if err != nil {
			return nil, invalidCertification("expires_at must be a date like 2026-10-08")
		}
		if parsed.Before(issuedAt) {
			return nil, invalidCertification("expires_at cannot be before issued_at")
		}
		expiresAt = &parsed
	}

	var credentialURL *string
	if req.CredentialURL != nil && strings.TrimSpace(*req.CredentialURL) != "" {
		link := strings.TrimSpace(*req.CredentialURL)
		parsed, err := url.ParseRequestURI(link)
		if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
			return nil, invalidCertification("credential_url must be a valid http or https link")
		}
		credentialURL = &link
	}

	exists, err := s.repo.TalentExists(ctx, talentID)
	if err != nil {
		return nil, err
	}
	if !exists {
		s.logger.Warn("cannot add certification: talent not found", "talent_id", talentID)
		return nil, repository.ErrTalentNotFoundForCertification
	}

	// Decide the status from the expiry date.
	status := certificationStatusActive
	today := time.Now().UTC().Truncate(24 * time.Hour)
	if expiresAt != nil && expiresAt.Before(today) {
		status = certificationStatusExpired
	}

	return s.repo.Create(ctx, &model.TalentCertification{
		TalentID:        talentID,
		CertificationID: certificationID,
		IssuedAt:        issuedAt,
		ExpiresAt:       expiresAt,
		CredentialURL:   credentialURL,
		Status:          status,
	})
}

package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"talentiq/talent-profile-intelligence/internal/model"
)

// Named errors, so the other layers can tell "what went wrong" without reading text.
var (
	ErrInvalidTalentCertification     = errors.New("invalid talent certification") // bad input
	ErrTalentNotFoundForCertification = errors.New("talent not found")             // person doesn't exist
	ErrCertificationAlreadyExists     = errors.New("certification already added")  // same certificate added twice
)

// The list of database jobs this file can do.
type TalentCertificationRepository interface {
	TalentExists(ctx context.Context, talentID uuid.UUID) (bool, error)
	Create(ctx context.Context, cert *model.TalentCertification) (*model.TalentCertification, error)
}

type talentCertificationRepository struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewTalentCertificationRepository(db *pgxpool.Pool, logger *slog.Logger) TalentCertificationRepository {
	return &talentCertificationRepository{db: db, logger: logger}
}

// Job 1: "Is there an active person with this id?"
func (r *talentCertificationRepository) TalentExists(ctx context.Context, talentID uuid.UUID) (bool, error) {
	const query = `
		SELECT EXISTS (
			SELECT 1 FROM talent_profiles
			WHERE id = $1 AND deleted_at IS NULL
		)
	`

	var exists bool
	if err := r.db.QueryRow(ctx, query, talentID).Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to check talent exists: %w", err)
	}
	return exists, nil
}

// Job 2: "Save this certification for the person."
func (r *talentCertificationRepository) Create(ctx context.Context, cert *model.TalentCertification) (*model.TalentCertification, error) {
	const query = `
		INSERT INTO talent_certifications
			(talent_id, certification_id, issued_at, expires_at, credential_url, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING
			id, talent_id, certification_id, issued_at, expires_at,
			credential_url, status, created_at
	`

	saved := &model.TalentCertification{}
	err := r.db.QueryRow(ctx, query,
		cert.TalentID,
		cert.CertificationID,
		cert.IssuedAt,
		cert.ExpiresAt,
		cert.CredentialURL,
		cert.Status,
	).Scan(
		&saved.ID,
		&saved.TalentID,
		&saved.CertificationID,
		&saved.IssuedAt,
		&saved.ExpiresAt,
		&saved.CredentialURL,
		&saved.Status,
		&saved.CreatedAt,
	)
	if err != nil {
		// 23505 is PostgreSQL's code for "this row already exists" (unique index).
		if isUniqueViolation(err) {
			return nil, ErrCertificationAlreadyExists
		}
		return nil, fmt.Errorf("failed to create talent certification: %w", err)
	}

	r.logger.Info("talent certification added",
		"talent_id", saved.TalentID,
		"certification_id", saved.CertificationID,
		"status", saved.Status,
	)
	return saved, nil
}

// Checks whether PostgreSQL rejected the insert as a duplicate.
// pgx's error type has a SQLState() method, so no extra import is needed.
func isUniqueViolation(err error) bool {
	var withState interface{ SQLState() string }
	return errors.As(err, &withState) && withState.SQLState() == "23505"
}

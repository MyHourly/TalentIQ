package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"talentiq/talent-profile-intelligence/internal/model"
)

// Named errors, so the other layers can tell "what went wrong" without reading text.
var (
	ErrInvalidTalentProject     = errors.New("invalid talent project")    // bad input -> 400
	ErrTalentNotFoundForProject = errors.New("talent not found")          // person doesn't exist -> 404
	ErrTalentProjectNotFound    = errors.New("project not found")         // no such project for this talent -> 404
	ErrProjectAlreadyCompleted  = errors.New("project already completed") // completing twice -> 400
)

// The list of database jobs this file can do.
type TalentProjectRepository interface {
	TalentExists(ctx context.Context, talentID uuid.UUID) (bool, error)
	Create(ctx context.Context, project *model.TalentProject) (*model.TalentProject, error)
	GetByTalentAndProjectID(ctx context.Context, talentID, projectID uuid.UUID) (*model.TalentProject, error)
	Complete(ctx context.Context, talentID, projectID uuid.UUID, endDate time.Time) (*model.TalentProject, error)
}

type talentProjectRepository struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewTalentProjectRepository(db *pgxpool.Pool, logger *slog.Logger) TalentProjectRepository {
	return &talentProjectRepository{db: db, logger: logger}
}

// Job 1: "Is there an active person with this id?"
func (r *talentProjectRepository) TalentExists(ctx context.Context, talentID uuid.UUID) (bool, error) {
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

// Job 2: "Save a new project and link the talent to it."
// Two tables change, so both inserts run in ONE transaction:
// either both are saved, or neither is.
func (r *talentProjectRepository) Create(ctx context.Context, project *model.TalentProject) (*model.TalentProject, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	// If we leave early (an error), this undoes everything.
	// After a successful Commit it does nothing.
	defer tx.Rollback(ctx)

	const insertProject = `
		INSERT INTO projects (name, description, domain, client_name, start_date, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, end_date, created_at, updated_at
	`

	saved := *project // work on a copy
	err = tx.QueryRow(ctx, insertProject,
		project.Name,
		project.Description,
		project.Domain,
		project.ClientName,
		project.StartDate,
		project.Status,
	).Scan(&saved.ID, &saved.EndDate, &saved.CreatedAt, &saved.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert project: %w", err)
	}

	const insertMember = `
		INSERT INTO project_members (project_id, talent_id, role, start_date)
		VALUES ($1, $2, $3, $4)
	`

	if _, err := tx.Exec(ctx, insertMember,
		saved.ID,
		saved.TalentID,
		saved.Role,
		saved.StartDate,
	); err != nil {
		return nil, fmt.Errorf("failed to insert project member: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit project: %w", err)
	}

	r.logger.Info("talent project added",
		"talent_id", saved.TalentID,
		"project_id", saved.ID,
		"name", saved.Name,
	)
	return &saved, nil
}

// Job 3: "Find this project, but only if this talent is a member of it."
func (r *talentProjectRepository) GetByTalentAndProjectID(ctx context.Context, talentID, projectID uuid.UUID) (*model.TalentProject, error) {
	const query = `
		SELECT
			p.id, pm.talent_id, p.name, p.description, p.domain, p.client_name,
			pm.role, p.start_date, p.end_date, p.status, p.created_at, p.updated_at
		FROM projects p
		JOIN project_members pm ON pm.project_id = p.id
		WHERE p.id = $1 AND pm.talent_id = $2
	`

	found := &model.TalentProject{}
	err := r.db.QueryRow(ctx, query, projectID, talentID).Scan(
		&found.ID,
		&found.TalentID,
		&found.Name,
		&found.Description,
		&found.Domain,
		&found.ClientName,
		&found.Role,
		&found.StartDate,
		&found.EndDate,
		&found.Status,
		&found.CreatedAt,
		&found.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTalentProjectNotFound
		}
		return nil, fmt.Errorf("failed to get talent project: %w", err)
	}
	return found, nil
}

// Job 4: "Mark the project as completed on this date."
// Again two tables change (the project and this talent's membership), so one transaction.
func (r *talentProjectRepository) Complete(ctx context.Context, talentID, projectID uuid.UUID, endDate time.Time) (*model.TalentProject, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// "status <> 'COMPLETED'" means a project can only be completed once,
	// even if two requests arrive at the same moment.
	const updateProject = `
		UPDATE projects
		SET status = 'COMPLETED', end_date = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND status <> 'COMPLETED'
		RETURNING name, description, domain, client_name, start_date, status, created_at, updated_at
	`

	done := &model.TalentProject{ID: projectID, TalentID: talentID, EndDate: &endDate}
	err = tx.QueryRow(ctx, updateProject, projectID, endDate).Scan(
		&done.Name,
		&done.Description,
		&done.Domain,
		&done.ClientName,
		&done.StartDate,
		&done.Status,
		&done.CreatedAt,
		&done.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// The service already confirmed the project exists,
			// so "no row updated" means someone completed it first.
			return nil, ErrProjectAlreadyCompleted
		}
		return nil, fmt.Errorf("failed to complete project: %w", err)
	}

	const updateMember = `
		UPDATE project_members
		SET end_date = $3
		WHERE project_id = $1 AND talent_id = $2
		RETURNING role
	`

	if err := tx.QueryRow(ctx, updateMember, projectID, talentID, endDate).Scan(&done.Role); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTalentProjectNotFound
		}
		return nil, fmt.Errorf("failed to complete project member: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit project completion: %w", err)
	}

	r.logger.Info("talent project completed",
		"talent_id", talentID,
		"project_id", projectID,
		"end_date", endDate.Format("2006-01-02"),
	)
	return done, nil
}

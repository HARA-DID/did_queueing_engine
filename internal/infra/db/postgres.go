package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/myorg/worker-service/internal/domain"
	"github.com/myorg/worker-service/internal/repository"

	// PostgreSQL driver
	_ "github.com/lib/pq"
)

var _ repository.JobRepository = (*PostgresJobRepository)(nil)

type PostgresJobRepository struct {
	db *sql.DB
}

func NewPostgresJobRepository(db *sql.DB) *PostgresJobRepository {
	return &PostgresJobRepository{db: db}
}

func (r *PostgresJobRepository) Create(ctx context.Context, job *domain.Job) error {
	const q = `
		INSERT INTO jobs (id, event_id, type, status, tx_hashes, retry_count, error_message, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	now := time.Now().UTC()
	job.CreatedAt = now
	job.UpdatedAt = now

	txArr := txHashArray(job.TxHashes)
	_, err := r.db.ExecContext(ctx, q,
		job.ID,
		job.EventID,
		job.Type,
		string(job.Status),
		txArr,
		job.RetryCount,
		job.ErrorMessage,
		job.CreatedAt,
		job.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("postgres create job: %w", err)
	}
	return nil
}

// FindByEventID returns a job by event_id, or (nil, nil) if not found.
func (r *PostgresJobRepository) FindByEventID(ctx context.Context, eventID string) (*domain.Job, error) {
	const q = `
		SELECT id, event_id, type, status, tx_hashes, retry_count, error_message, created_at, updated_at
		FROM jobs
		WHERE event_id = $1
		LIMIT 1
	`
	row := r.db.QueryRowContext(ctx, q, eventID)

	var job domain.Job
	var txHashStr string
	var status string
	err := row.Scan(
		&job.ID,
		&job.EventID,
		&job.Type,
		&status,
		&txHashStr,
		&job.RetryCount,
		&job.ErrorMessage,
		&job.CreatedAt,
		&job.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("postgres find job by event_id: %w", err)
	}

	job.Status = domain.JobStatus(status)
	job.TxHashes = parseTxHashArray(txHashStr)
	return &job, nil
}

// UpdateStatus updates a job's outcome fields.
func (r *PostgresJobRepository) UpdateStatus(
	ctx context.Context,
	jobID string,
	status domain.JobStatus,
	txHashes []string,
	errMsg string,
) error {
	const q = `
		UPDATE jobs
		SET status = $1, tx_hashes = $2, error_message = $3, updated_at = $4
		WHERE id = $5
	`
	_, err := r.db.ExecContext(ctx, q,
		string(status),
		txHashArray(txHashes),
		errMsg,
		time.Now().UTC(),
		jobID,
	)
	if err != nil {
		return fmt.Errorf("postgres update job status: %w", err)
	}
	return nil
}

// IncrementRetry bumps the retry_count and records an error message.
func (r *PostgresJobRepository) IncrementRetry(ctx context.Context, jobID string, errMsg string) error {
	const q = `
		UPDATE jobs
		SET retry_count = retry_count + 1, error_message = $1, updated_at = $2
		WHERE id = $3
	`
	_, err := r.db.ExecContext(ctx, q, errMsg, time.Now().UTC(), jobID)
	if err != nil {
		return fmt.Errorf("postgres increment retry: %w", err)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Helpers for PostgreSQL text[] serialization
// ---------------------------------------------------------------------------

// txHashArray converts a string slice to a Postgres array literal: {a,b,c}
func txHashArray(hashes []string) string {
	if len(hashes) == 0 {
		return "{}"
	}
	return "{" + strings.Join(hashes, ",") + "}"
}

// parseTxHashArray parses a Postgres array literal back to []string.
func parseTxHashArray(s string) []string {
	s = strings.TrimPrefix(s, "{")
	s = strings.TrimSuffix(s, "}")
	if s == "" {
		return nil
	}
	return strings.Split(s, ",")
}
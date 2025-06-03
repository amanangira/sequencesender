package storage

import (
	"context"
	"database/sql"
	"fmt"
	"sequencesender/internal/types"

	"github.com/jmoiron/sqlx"
)

// PostgresStorage provides database operations without holding DB connections
type PostgresStorage struct{}

// NewPostgresStorage creates a new storage instance (no DB dependency)
func NewPostgresStorage() *PostgresStorage {
	return &PostgresStorage{}
}

// CreateSequence creates a sequence
func (s *PostgresStorage) CreateSequence(ctx context.Context, tx *sql.Tx, req types.CreateSequenceRequest) (int, error) {
	query := `INSERT INTO sequences (name, open_tracking_enabled, click_tracking_enabled) 
		VALUES ($1, $2, $3)RETURNING id`

	var sequenceID int
	err := tx.QueryRowContext(ctx, query, req.Name, req.OpenTrackingEnabled, req.ClickTrackingEnabled).Scan(&sequenceID)
	if err != nil {
		return 0, fmt.Errorf("create sequence error %w", err)
	}

	return sequenceID, nil
}

func (s *PostgresStorage) CreateSteps(ctx context.Context, tx *sql.Tx, sequenceID int, steps []types.CreateStepRequest) error {
	if len(steps) == 0 {
		return nil
	}

	query := `INSERT INTO sequence_steps (sequence_id, name, body_content, days_to_wait, order_number)
		VALUES ($1, $2, $3, $4, $5)`

	for _, step := range steps {
		_, err := tx.ExecContext(ctx, query, sequenceID, step.Name, step.Content, step.DaysToWait, step.Order)
		if err != nil {
			return fmt.Errorf("sequence step error %s: %w", step.Name, err)
		}
	}

	return nil
}

func (s *PostgresStorage) GetSequenceByID(ctx context.Context, db *sqlx.DB, id int) (*types.Sequence, error) {
	query := `SELECT id, name, open_tracking_enabled, click_tracking_enabled, created_at, updated_at, is_deleted FROM sequences WHERE id = $1 AND is_deleted IS NULL`

	var seq types.Sequence
	err := db.GetContext(ctx, &seq, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get sequence: %w", err)
	}

	return &seq, nil
}

// GetStepsBySequenceID retrieves all steps for a sequence
func (s *PostgresStorage) GetStepsBySequenceID(ctx context.Context, db *sqlx.DB, sequenceID int) ([]types.SequenceStep, error) {
	query := `SELECT id, sequence_id, name, body_content, days_to_wait, order_number, created_at, updated_at, is_deleted FROM sequence_steps 
		WHERE sequence_id = $1 AND is_deleted = false
		ORDER BY order_number`

	var steps []types.SequenceStep
	err := db.SelectContext(ctx, &steps, query, sequenceID)
	if err != nil {
		return nil, fmt.Errorf("get sequence steps error: %w", err)
	}

	return steps, nil
}

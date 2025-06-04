package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"sequencesender/internal/types"

	"github.com/jmoiron/sqlx"
)

// to check at compile time that PostgresStorage implements StorageInterface
var _ StorageInterface = (*PostgresStorage)(nil)

// PostgresStorage provides database operations
type PostgresStorage struct{}

// NewPostgresStorage creates a new storage instance
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
	log.Println("id : ", id)
	query := `SELECT id, name, open_tracking_enabled, click_tracking_enabled, created_at, updated_at, is_deleted FROM sequences WHERE id = $1 AND is_deleted IS false`

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

// UpdateStepByID updates a sequence step's name and/or content
func (s *PostgresStorage) UpdateStepByID(ctx context.Context, db *sqlx.DB, stepID int, name *string, content *string) error {
	setParts := []string{"updated_at = NOW()"}
	args := []interface{}{}
	argIndex := 1

	if name != nil {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *name)
		argIndex++
	}

	if content != nil {
		setParts = append(setParts, fmt.Sprintf("body_content = $%d", argIndex))
		args = append(args, *content)
		argIndex++
	}

	if len(setParts) == 1 {
		return fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf("UPDATE sequence_steps SET %s WHERE id = $%d AND is_deleted = false",
		strings.Join(setParts, ", "), argIndex)
	args = append(args, stepID)

	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update step: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("step not found or already deleted")
	}

	return nil
}

// SoftDeleteStepByID soft deletes a sequence step
func (s *PostgresStorage) SoftDeleteStepByID(ctx context.Context, db *sqlx.DB, stepID int) error {
	query := `UPDATE sequence_steps SET is_deleted = true, updated_at = NOW() WHERE id = $1 AND is_deleted = false`

	result, err := db.ExecContext(ctx, query, stepID)
	if err != nil {
		return fmt.Errorf("failed to soft delete step: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("step not found or already deleted")
	}

	return nil
}

// GetStepByID retrieves a single step by ID
func (s *PostgresStorage) GetStepByID(ctx context.Context, db *sqlx.DB, stepID int) (*types.SequenceStep, error) {
	query := `SELECT id, sequence_id, name, body_content, days_to_wait, order_number, created_at, updated_at, is_deleted 
		FROM sequence_steps WHERE id = $1 AND is_deleted = false`

	var step types.SequenceStep
	err := db.GetContext(ctx, &step, query, stepID)
	if err != nil {
		return nil, fmt.Errorf("failed to get step: %w", err)
	}

	return &step, nil
}

// UpdateSequenceTracking updates tracking settings for a sequence
func (s *PostgresStorage) UpdateSequenceTracking(ctx context.Context, db *sqlx.DB, sequenceID int, openTracking *bool, clickTracking *bool) error {
	setParts := []string{"updated_at = NOW()"}
	args := []interface{}{}
	argIndex := 1

	if openTracking != nil {
		setParts = append(setParts, fmt.Sprintf("open_tracking_enabled = $%d", argIndex))
		args = append(args, *openTracking)
		argIndex++
	}

	if clickTracking != nil {
		setParts = append(setParts, fmt.Sprintf("click_tracking_enabled = $%d", argIndex))
		args = append(args, *clickTracking)
		argIndex++
	}

	if len(setParts) == 1 {
		return fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf("UPDATE sequences SET %s WHERE id = $%d AND is_deleted = false",
		strings.Join(setParts, ", "), argIndex)
	args = append(args, sequenceID)

	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update sequence tracking: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("sequence not found or doesn't exist")
	}

	return nil
}

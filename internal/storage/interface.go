package storage

import (
	"context"
	"database/sql"

	"sequencesender/internal/types"

	"github.com/jmoiron/sqlx"
)

// StorageInterface contract for DB operations
type StorageInterface interface {
	// Sequence

	CreateSequence(ctx context.Context, tx *sql.Tx, req types.CreateSequenceRequest) (int, error)
	GetSequenceByID(ctx context.Context, db *sqlx.DB, id int) (*types.Sequence, error)
	UpdateSequenceTracking(ctx context.Context, db *sqlx.DB, sequenceID int, openTracking *bool, clickTracking *bool) error

	// Sequence

	CreateSteps(ctx context.Context, tx *sql.Tx, sequenceID int, steps []types.CreateStepRequest) error
	GetStepsBySequenceID(ctx context.Context, db *sqlx.DB, sequenceID int) ([]types.SequenceStep, error)
	GetStepByID(ctx context.Context, db *sqlx.DB, stepID int) (*types.SequenceStep, error)
	UpdateStepByID(ctx context.Context, db *sqlx.DB, stepID int, name *string, content *string) error
	SoftDeleteStepByID(ctx context.Context, db *sqlx.DB, stepID int) error
}

package services

import (
	"context"
	"errors"
	"fmt"
	"sequencesender"
	"sequencesender/internal/storage"
	"sequencesender/internal/types"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type SequenceService struct {
	db      *sqlx.DB
	storage *storage.PostgresStorage
}

func NewSequenceService(db *sqlx.DB) *SequenceService {
	return &SequenceService{
		db:      db,
		storage: storage.NewPostgresStorage(),
	}
}

// CreateSequence - business logic of creating a sequence and it's steps, includes validations
func (s *SequenceService) CreateSequence(ctx context.Context, req types.CreateSequenceRequest) (*sequencesender.CreateResponse, error) {
	if err := s.validateCreateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	// create sequence
	sequenceID, err := s.storage.CreateSequence(ctx, tx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create sequence: %w", err)
	}

	// create sequence steps
	if len(req.Steps) > 0 {
		if err := s.storage.CreateSteps(ctx, tx, sequenceID, req.Steps); err != nil {
			return nil, fmt.Errorf("failed to create sequence steps: %w", err)
		}
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", commitErr)
	}

	return &sequencesender.CreateResponse{
		ID: sequenceID,
	}, nil
}

// GetSequence retrieves a sequence by ID
func (s *SequenceService) GetSequence(ctx context.Context, idStr string) (*types.SequenceResponse, error) {
	id, err := strconv.Atoi(strings.TrimSpace(idStr))
	if err != nil {
		return nil, fmt.Errorf("invalid sequence ID format: %w", err)
	}

	sequence, err := s.storage.GetSequenceByID(ctx, s.db, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get sequence - %w", err)
	}

	steps, err := s.storage.GetStepsBySequenceID(ctx, s.db, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get sequence steps %w", err)
	}

	var createdAt, updatedAt time.Time
	if sequence.CreatedAt != nil {
		createdAt = *sequence.CreatedAt
	}
	if sequence.UpdatedAt != nil {
		updatedAt = *sequence.UpdatedAt
	}

	stepResponses := make([]types.StepResponse, len(steps))
	for i, step := range steps {
		var stepCreatedAt, stepUpdatedAt time.Time
		if step.CreatedAt != nil {
			stepCreatedAt = *step.CreatedAt
		}
		if step.UpdatedAt != nil {
			stepUpdatedAt = *step.UpdatedAt
		}

		stepResponses[i] = types.StepResponse{
			ID:         step.ID,
			Name:       step.Name,
			Content:    step.BodyContent,
			DaysToWait: step.DaysToWait,
			Order:      step.OrderNumber,
			CreatedAt:  stepCreatedAt,
			UpdatedAt:  stepUpdatedAt,
		}
	}

	return &types.SequenceResponse{
		ID:                   sequence.ID,
		Name:                 sequence.Name,
		OpenTrackingEnabled:  sequence.OpenTrackingEnabled,
		ClickTrackingEnabled: sequence.ClickTrackingEnabled,
		StepsCount:           len(steps),
		Steps:                stepResponses,
		CreatedAt:            createdAt,
		UpdatedAt:            updatedAt,
	}, nil
}

// validateCreateRequest validation on the create sequence payload
func (s *SequenceService) validateCreateRequest(req types.CreateSequenceRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return errors.New("sequence name is required")
	}

	orderMap := make(map[int]bool)
	for i, step := range req.Steps {
		if strings.TrimSpace(step.Name) == "" {
			return fmt.Errorf("step %d: name is required", i+1)
		}

		if len(step.Name) > 255 {
			return fmt.Errorf("step %d: name too long (max 255 characters)", i+1)
		}

		if strings.TrimSpace(step.Content) == "" {
			return fmt.Errorf("step %d: content is required", i+1)
		}

		if step.DaysToWait < 0 {
			return fmt.Errorf("step %d: days_to_wait cannot be negative", i+1)
		}

		if step.Order <= 0 {
			return fmt.Errorf("step %d: order must be positive", i+1)
		}

		if orderMap[step.Order] {
			return fmt.Errorf("step %d: duplicate order number %d", i+1, step.Order)
		}
		orderMap[step.Order] = true
	}

	return nil
}

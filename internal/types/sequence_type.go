package types

import "time"

type CreateSequenceRequest struct {
	Name                 string              `json:"name" validate:"required,max=255"`
	OpenTrackingEnabled  bool                `json:"open_tracking_enabled"`
	ClickTrackingEnabled bool                `json:"click_tracking_enabled"`
	Steps                []CreateStepRequest `json:"steps,omitempty"`
}

type CreateStepRequest struct {
	Name       string `json:"name" validate:"required,max=255"`
	Content    string `json:"content" validate:"required"`
	DaysToWait int    `json:"days_to_wait" validate:"min=0"`
	Order      int    `json:"order" validate:"required,min=1"`
}

type SequenceResponse struct {
	ID                   int            `json:"id"`
	Name                 string         `json:"name"`
	OpenTrackingEnabled  bool           `json:"open_tracking_enabled"`
	ClickTrackingEnabled bool           `json:"click_tracking_enabled"`
	StepsCount           int            `json:"steps_count"`
	Steps                []StepResponse `json:"steps"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
}

type StepResponse struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Content    string    `json:"content"`
	DaysToWait int       `json:"days_to_wait"`
	Order      int       `json:"order"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UpdateStepRequest struct {
	Name    *string `json:"name,omitempty"`
	Content *string `json:"content,omitempty"`
}

type UpdateStepResponse struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Content    string    `json:"content"`
	DaysToWait int       `json:"days_to_wait"`
	Order      int       `json:"order"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UpdateSequenceTrackingRequest struct {
	OpenTrackingEnabled  *bool `json:"open_tracking_enabled,omitempty"`
	ClickTrackingEnabled *bool `json:"click_tracking_enabled,omitempty"`
}

type UpdateSequenceTrackingResponse struct {
	ID                   int       `json:"id"`
	Name                 string    `json:"name"`
	OpenTrackingEnabled  bool      `json:"open_tracking_enabled"`
	ClickTrackingEnabled bool      `json:"click_tracking_enabled"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type Sequence struct {
	ID                   int        `db:"id"`
	Name                 string     `db:"name"`
	OpenTrackingEnabled  bool       `db:"open_tracking_enabled"`
	ClickTrackingEnabled bool       `db:"click_tracking_enabled"`
	CreatedAt            *time.Time `db:"created_at"`
	UpdatedAt            *time.Time `db:"updated_at"`
	IsDeleted            bool       `db:"is_deleted"`
}

type SequenceStep struct {
	ID          int        `db:"id"`
	SequenceID  int        `db:"sequence_id"`
	Name        string     `db:"name"`
	BodyContent string     `db:"body_content"`
	DaysToWait  int        `db:"days_to_wait"`
	OrderNumber int        `db:"order_number"`
	CreatedAt   *time.Time `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
	IsDeleted   bool       `db:"is_deleted"`
}

// APIResponse response wrapper
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

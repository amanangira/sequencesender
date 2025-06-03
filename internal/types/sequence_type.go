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
	ID                   int       `json:"id"`
	Name                 string    `json:"name"`
	OpenTrackingEnabled  bool      `json:"open_tracking_enabled"`
	ClickTrackingEnabled bool      `json:"click_tracking_enabled"`
	StepsCount           int       `json:"steps_count"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type Sequence struct {
	ID                   int        `db:"id"`
	Name                 string     `db:"name"`
	OpenTrackingEnabled  bool       `db:"open_tracking_enabled"`
	ClickTrackingEnabled bool       `db:"click_tracking_enabled"`
	CreatedAt            *time.Time `db:"created_at"`
	UpdatedAt            *time.Time `db:"updated_at"`
	IsDeleted            *time.Time `db:"is_deleted"`
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

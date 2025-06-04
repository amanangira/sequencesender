package services

import (
	"context"
	"strings"
	"testing"

	"sequencesender/internal/types"
	"sequencesender/tests/mocks"

	"github.com/jmoiron/sqlx"
)

func TestValidateCreateRequest_Success(t *testing.T) {
	service := &SequenceService{}

	tests := []struct {
		name string
		req  types.CreateSequenceRequest
	}{
		{
			name: "valid request with valid steps",
			req: types.CreateSequenceRequest{
				Name:                 "Test Sequence",
				OpenTrackingEnabled:  true,
				ClickTrackingEnabled: false,
				Steps: []types.CreateStepRequest{
					{
						Name:       "Step 1",
						Content:    "Content 1",
						DaysToWait: 0,
						Order:      1,
					},
					{
						Name:       "Step 2",
						Content:    "Content 2",
						DaysToWait: 2,
						Order:      2,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateCreateRequest(tt.req)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestValidateCreateRequest_ValidationErrors(t *testing.T) {
	service := &SequenceService{}

	tests := []struct {
		name        string
		req         types.CreateSequenceRequest
		expectedErr string
	}{
		{
			name: "empty sequence name",
			req: types.CreateSequenceRequest{
				Name: "",
			},
			expectedErr: "sequence name is required",
		},
		{
			name: "empty step name",
			req: types.CreateSequenceRequest{
				Name: "Test Sequence",
				Steps: []types.CreateStepRequest{
					{
						Name:       "",
						Content:    "Content",
						DaysToWait: 0,
						Order:      1,
					},
				},
			},
			expectedErr: "step 1: name is required",
		},
		{
			name: "empty step content",
			req: types.CreateSequenceRequest{
				Name: "Test Sequence",
				Steps: []types.CreateStepRequest{
					{
						Name:       "Step 1",
						Content:    "",
						DaysToWait: 0,
						Order:      1,
					},
				},
			},
			expectedErr: "step 1: content is required",
		},
		{
			name: "negative days to wait",
			req: types.CreateSequenceRequest{
				Name: "Test Sequence",
				Steps: []types.CreateStepRequest{
					{
						Name:       "Step 1",
						Content:    "Content",
						DaysToWait: -1,
						Order:      1,
					},
				},
			},
			expectedErr: "step 1: days_to_wait cannot be negative",
		},
		{
			name: "duplicate order numbers",
			req: types.CreateSequenceRequest{
				Name: "Test Sequence",
				Steps: []types.CreateStepRequest{
					{
						Name:       "Step 1",
						Content:    "Content 1",
						DaysToWait: 0,
						Order:      1,
					},
					{
						Name:       "Step 2",
						Content:    "Content 2",
						DaysToWait: 1,
						Order:      1,
					},
				},
			},
			expectedErr: "step 2: duplicate order number 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateCreateRequest(tt.req)
			if err == nil {
				t.Errorf("expected error containing '%s', got nil", tt.expectedErr)
			} else if !strings.Contains(err.Error(), tt.expectedErr) {
				t.Errorf("expected error containing '%s', got '%s'", tt.expectedErr, err.Error())
			}
		})
	}
}

func TestCreateSequence(t *testing.T) {
	tests := []struct {
		name        string
		req         types.CreateSequenceRequest
		setupMocks  func(*mocks.StorageInterface)
		expectError bool
		expectedErr string
	}{
		{
			name: "validation passes for valid request with steps",
			req: types.CreateSequenceRequest{
				Name:                 "Test Sequence",
				OpenTrackingEnabled:  true,
				ClickTrackingEnabled: false,
				Steps: []types.CreateStepRequest{
					{
						Name:       "Step 1",
						Content:    "Content 1",
						DaysToWait: 0,
						Order:      1,
					},
				},
			},
			setupMocks: func(storageMock *mocks.StorageInterface) {
			},
			expectError: false,
		},
		{
			name: "validation fails for invalid request",
			req: types.CreateSequenceRequest{
				Name: "",
			},
			setupMocks: func(storageMock *mocks.StorageInterface) {

			},
			expectError: true,
			expectedErr: "validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			storageMock := &mocks.StorageInterface{}
			tt.setupMocks(storageMock)

			mockDB := &sqlx.DB{}
			service := NewSequenceServiceWithStorage(mockDB, storageMock)

			if tt.name == "validation passes for valid request with steps" {

				err := service.validateCreateRequest(tt.req)
				if err != nil {
					t.Errorf("validation should pass, got error: %v", err)
				}
			} else {

				result, err := service.CreateSequence(context.Background(), tt.req)
				if !tt.expectError {
					if err != nil {
						t.Errorf("expected no error, got %v", err)
					}
					if result == nil {
						t.Error("expected result, got nil")
					}
				} else {
					if err == nil {
						t.Error("expected error, got nil")
					}
					if result != nil {
						t.Error("expected nil result, got non-nil")
					}
					if tt.expectedErr != "" && !strings.Contains(err.Error(), tt.expectedErr) {
						t.Errorf("expected error containing '%s', got '%s'", tt.expectedErr, err.Error())
					}
				}
			}
		})
	}
}

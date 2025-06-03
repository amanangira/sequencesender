package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"sequencesender/internal/services"
	"sequencesender/internal/types"
	"sequencesender/pkg/httputil"

	"github.com/go-chi/chi/v5"
)

type SequenceHandler struct {
	service *services.SequenceService
}

func NewSequenceHandler(service *services.SequenceService) *SequenceHandler {
	return &SequenceHandler{
		service: service,
	}
}

// CreateSequence - POST /api/sequences
func (h *SequenceHandler) CreateSequence(w http.ResponseWriter, r *http.Request) {
	var req types.CreateSequenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("bad request", slog.String("error", err.Error()))
		httputil.BadRequest(w, "")
		return
	}

	response, err := h.service.CreateSequence(r.Context(), req)
	if err != nil {
		slog.Error("error creating response", slog.String("error", err.Error()), slog.String("sequence_name", req.Name))

		// TODO - replace contains check with different concrete implementations of Error interface
		if strings.Contains(err.Error(), "validation failed") {
			httputil.BadRequest(w, err.Error())
		} else {
			httputil.InternalError(w, "")
		}
		return
	}

	httputil.Success(w, response, "success")
}

// GetSequence - /api/sequences/{id}
func (h *SequenceHandler) GetSequence(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		httputil.BadRequest(w, "sequence ID is required")
		return
	}

	response, err := h.service.GetSequence(r.Context(), id)
	if err != nil {
		slog.Error("failed to get sequence", slog.String("error", err.Error()), slog.String("sequenceid", id))

		if strings.Contains(err.Error(), "not found") {
			httputil.Error(w, http.StatusNotFound, "")
		} else {
			httputil.InternalError(w, "")
		}
		return
	}

	httputil.Success(w, response, "")
}

// UpdateStep - PATCH /api/steps/{id}
func (h *SequenceHandler) UpdateStep(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		httputil.BadRequest(w, "stepID required")
		return
	}

	var req types.UpdateStepRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("bad request", slog.String("error", err.Error()))
		httputil.BadRequest(w, "")
		return
	}

	err := h.service.UpdateStep(r.Context(), id, req)

	if err != nil {
		slog.Error("failed to update step", slog.String("error", err.Error()), slog.String("stepid", id))

		if strings.Contains(err.Error(), "validation failed") {
			httputil.BadRequest(w, err.Error())
		} else if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "already deleted") {
			httputil.Error(w, http.StatusNotFound, "")
		} else {
			httputil.InternalError(w, "")
		}
		return
	}

	httputil.Success(w, nil, "")
}

// SoftDeleteStep - DELETE /api/steps/{id}
func (h *SequenceHandler) SoftDeleteStep(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		httputil.BadRequest(w, "step ID is required")
		return
	}

	err := h.service.SoftDeleteStep(r.Context(), id)
	if err != nil {
		slog.Error("failed to soft delete step", slog.String("error", err.Error()), slog.String("stepid", id))

		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "already deleted") {
			httputil.Error(w, http.StatusNotFound, "")
		} else {
			httputil.InternalError(w, "")
		}
		return
	}

	httputil.Success(w, nil, "step deleted successfully")
}

// UpdateSequenceTracking - PATCH /api/sequences/{id}/tracking
func (h *SequenceHandler) UpdateSequenceTracking(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		httputil.BadRequest(w, "sequence ID is required")
		return
	}

	var req types.UpdateSequenceTrackingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("bad request", slog.String("error", err.Error()))
		httputil.BadRequest(w, "")
		return
	}

	err := h.service.UpdateSequenceTracking(r.Context(), id, req)
	if err != nil {
		slog.Error("failed to update sequence tracking", slog.String("error", err.Error()), slog.String("sequenceid", id))

		if strings.Contains(err.Error(), "validation failed") {
			httputil.BadRequest(w, err.Error())
		} else if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "already deleted") {
			httputil.Error(w, http.StatusNotFound, "")
		} else {
			httputil.InternalError(w, "")
		}
		return
	}

	httputil.Success(w, nil, "")
}

// RegisterRoutes registers all sequence routes
func (h *SequenceHandler) RegisterRoutes(r chi.Router) {
	r.Post("/sequences", h.CreateSequence)
	r.Get("/sequences/{id}", h.GetSequence)
	r.Patch("/sequences/{id}", h.UpdateSequenceTracking)
	r.Patch("/steps/{id}", h.UpdateStep)
	r.Delete("/steps/{id}", h.SoftDeleteStep)
}

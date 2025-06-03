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

// RegisterRoutes registers all sequence routes
func (h *SequenceHandler) RegisterRoutes(r chi.Router) {
	r.Post("/sequences", h.CreateSequence)
	r.Get("/sequences/{id}", h.GetSequence)
}

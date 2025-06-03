package httputil

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sequencesender/internal/types"
)

// JSON writes a JSON response with the given status code
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", slog.String("error", err.Error()))
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// Success writes a successful JSON response
func Success(w http.ResponseWriter, data interface{}, message string) {
	response := types.APIResponse{
		Success: true,
		Data:    data,
		Message: message,
	}
	JSON(w, http.StatusOK, response)
}

// Error writes an error JSON response
func Error(w http.ResponseWriter, statusCode int, message string) {
	response := types.APIResponse{
		Success: false,
		Error:   message,
	}
	JSON(w, statusCode, response)
}

// BadRequest writes a 400 Bad Request response
func BadRequest(w http.ResponseWriter, message string) {
	if message == "" {
		message = "bad request"
	}
	Error(w, http.StatusBadRequest, message)
}

// InternalError writes a 500 Internal Server Error response
func InternalError(w http.ResponseWriter, message string) {
	if message == "" {
		message = "internal server error"
	}
	Error(w, http.StatusInternalServerError, message)
}

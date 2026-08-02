package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/qper/hflow/internal/domain"
)

var (
	ErrValidation = errors.New("validation failed")
)

// ErrorHandler maps domain and validation errors to a consistent HTTP response.
type ErrorHandler struct{}

func (ErrorHandler) Write(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	switch {
	case errors.Is(err, domain.ErrUnauthorized):
		w.WriteHeader(http.StatusUnauthorized)
	case errors.Is(err, domain.ErrHabitNotFound), errors.Is(err, domain.ErrEntryNotFound):
		w.WriteHeader(http.StatusNotFound)
	case errors.Is(err, domain.ErrConflict):
		w.WriteHeader(http.StatusConflict)
	case errors.Is(err, ErrValidation), errors.Is(err, domain.ErrInvalidSchedule), errors.Is(err, domain.ErrInvalidHabitType), errors.Is(err, domain.ErrInvalidDate):
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: ErrorPayload{Code: errorCode(err), Message: err.Error()}})
}

func errorCode(err error) string {
	switch {
	case errors.Is(err, domain.ErrUnauthorized):
		return "unauthorized"
	case errors.Is(err, domain.ErrHabitNotFound):
		return "habit_not_found"
	case errors.Is(err, domain.ErrEntryNotFound):
		return "entry_not_found"
	case errors.Is(err, domain.ErrConflict):
		return "conflict"
	case errors.Is(err, ErrValidation), errors.Is(err, domain.ErrInvalidSchedule), errors.Is(err, domain.ErrInvalidHabitType), errors.Is(err, domain.ErrInvalidDate):
		return "validation_error"
	default:
		return "internal_error"
	}
}

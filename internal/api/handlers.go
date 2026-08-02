package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/qper/hflow/internal/auth"
	"github.com/qper/hflow/internal/domain"
	"github.com/qper/hflow/internal/habit"
)

// Handlers expose the HTTP transport layer for the domain services.
type Handlers struct {
	service      habit.Service
	errorHandler ErrorHandler
	authManager  auth.TokenManager
}

func NewHandlers(service habit.Service, accessSecret, refreshSecret string) Handlers {
	return Handlers{service: service, errorHandler: ErrorHandler{}, authManager: auth.NewTokenManager(accessSecret, refreshSecret)}
}

func (h Handlers) listHabits(w http.ResponseWriter, req *http.Request) {
	userID := auth.UserIDFromContext(req.Context())
	if userID == "" {
		h.errorHandler.Write(w, domain.ErrUnauthorized)
		return
	}
	habits, err := h.service.ListHabits(req.Context(), userID)
	if err != nil {
		h.errorHandler.Write(w, err)
		return
	}
	responses := make([]HabitResponse, 0, len(habits))
	for _, item := range habits {
		responses = append(responses, toHabitResponse(item))
	}
	writeJSON(w, http.StatusOK, responses)
}

func (h Handlers) createHabit(w http.ResponseWriter, req *http.Request) {
	userID := auth.UserIDFromContext(req.Context())
	if userID == "" {
		h.errorHandler.Write(w, domain.ErrUnauthorized)
		return
	}
	var payload HabitRequest
	if err := decodeJSON(req, &payload); err != nil {
		h.errorHandler.Write(w, ErrValidation)
		return
	}
	habitModel := domain.Habit{
		Name:        payload.Name,
		Description: payload.Description,
		HabitType:   payload.HabitType,
		Frequency:   payload.Frequency,
		TargetValue: payload.TargetValue,
		Unit:        payload.Unit,
		Color:       payload.Color,
		Icon:        payload.Icon,
	}
	created, err := h.service.CreateHabit(req.Context(), userID, habitModel)
	if err != nil {
		h.errorHandler.Write(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toHabitResponse(created))
}

func (h Handlers) getHabit(w http.ResponseWriter, req *http.Request) {
	userID := auth.UserIDFromContext(req.Context())
	if userID == "" {
		h.errorHandler.Write(w, domain.ErrUnauthorized)
		return
	}
	id := strings.TrimPrefix(req.URL.Path, "/api/v1/habits/")
	item, err := h.service.GetHabit(req.Context(), userID, id)
	if err != nil {
		h.errorHandler.Write(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toHabitResponse(item))
}

func (h Handlers) updateHabit(w http.ResponseWriter, req *http.Request) {
	userID := auth.UserIDFromContext(req.Context())
	if userID == "" {
		h.errorHandler.Write(w, domain.ErrUnauthorized)
		return
	}
	id := strings.TrimPrefix(req.URL.Path, "/api/v1/habits/")
	var payload HabitRequest
	if err := decodeJSON(req, &payload); err != nil {
		h.errorHandler.Write(w, ErrValidation)
		return
	}
	updated, err := h.service.UpdateHabit(req.Context(), userID, domain.Habit{ID: id, Name: payload.Name, Description: payload.Description, HabitType: payload.HabitType, Frequency: payload.Frequency, TargetValue: payload.TargetValue, Unit: payload.Unit, Color: payload.Color, Icon: payload.Icon})
	if err != nil {
		h.errorHandler.Write(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toHabitResponse(updated))
}

func (h Handlers) deleteHabit(w http.ResponseWriter, req *http.Request) {
	userID := auth.UserIDFromContext(req.Context())
	if userID == "" {
		h.errorHandler.Write(w, domain.ErrUnauthorized)
		return
	}
	id := strings.TrimPrefix(req.URL.Path, "/api/v1/habits/")
	if err := h.service.DeleteHabit(req.Context(), userID, id); err != nil {
		h.errorHandler.Write(w, err)
		return
	}
	writeJSON(w, http.StatusNoContent, nil)
}

func (h Handlers) createEntry(w http.ResponseWriter, req *http.Request) {
	userID := auth.UserIDFromContext(req.Context())
	if userID == "" {
		h.errorHandler.Write(w, domain.ErrUnauthorized)
		return
	}
	var payload EntryRequest
	if err := decodeJSON(req, &payload); err != nil {
		h.errorHandler.Write(w, ErrValidation)
		return
	}
	parsedDate, err := time.Parse("2006-01-02", payload.Date)
	if err != nil {
		h.errorHandler.Write(w, domain.ErrInvalidDate)
		return
	}
	entry, err := h.service.RecordEntry(req.Context(), userID, "", parsedDate, payload.Value, payload.Note, "UTC")
	if err != nil {
		h.errorHandler.Write(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toEntryResponse(entry))
}

func (h Handlers) deleteEntry(w http.ResponseWriter, req *http.Request) {
	userID := auth.UserIDFromContext(req.Context())
	if userID == "" {
		h.errorHandler.Write(w, domain.ErrUnauthorized)
		return
	}
	id := strings.TrimPrefix(req.URL.Path, "/api/v1/entries/")
	if err := h.service.DeleteEntryByID(req.Context(), userID, id); err != nil {
		h.errorHandler.Write(w, err)
		return
	}
	writeJSON(w, http.StatusNoContent, nil)
}

func decodeJSON(req *http.Request, dst any) error {
	defer req.Body.Close()
	return json.NewDecoder(req.Body).Decode(dst)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(payload)
}

func toHabitResponse(h domain.Habit) HabitResponse {
	return HabitResponse{ID: h.ID, Name: h.Name, Description: h.Description, HabitType: h.HabitType, Frequency: h.Frequency, TargetValue: h.TargetValue, Unit: h.Unit, Color: h.Color, Icon: h.Icon, CreatedAt: h.CreatedAt, UpdatedAt: h.UpdatedAt}
}

func toEntryResponse(e domain.Entry) EntryResponse {
	return EntryResponse{ID: e.ID, HabitID: e.HabitID, Date: e.Date.Format("2006-01-02"), Value: e.Value, Note: e.Note, CreatedAt: e.CreatedAt, UpdatedAt: e.UpdatedAt}
}

func (h Handlers) register(w http.ResponseWriter, req *http.Request) {
	var payload RegisterRequest
	if err := decodeJSON(req, &payload); err != nil {
		h.errorHandler.Write(w, ErrValidation)
		return
	}
	user := domain.User{Username: payload.Username, Email: payload.Email}
	_ = user
	writeJSON(w, http.StatusOK, TokenResponse{TokenType: "bearer"})
}

func (h Handlers) login(w http.ResponseWriter, req *http.Request) {
	var payload LoginRequest
	if err := decodeJSON(req, &payload); err != nil {
		h.errorHandler.Write(w, ErrValidation)
		return
	}
	writeJSON(w, http.StatusOK, TokenResponse{TokenType: "bearer"})
}

func (h Handlers) refresh(w http.ResponseWriter, req *http.Request) {
	var payload RefreshRequest
	if err := decodeJSON(req, &payload); err != nil {
		h.errorHandler.Write(w, ErrValidation)
		return
	}
	writeJSON(w, http.StatusOK, TokenResponse{TokenType: "bearer"})
}

func (h Handlers) logout(w http.ResponseWriter, req *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}

func (h Handlers) registerWithDB(db *sql.DB) {}

var _ = context.Background

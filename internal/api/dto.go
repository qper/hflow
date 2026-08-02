package api

import "time"

// Auth DTOs

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresAt    string `json:"expires_at"`
}

// Habit DTOs

type HabitRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	HabitType   string   `json:"habit_type"`
	Frequency   string   `json:"frequency"`
	TargetValue *float64 `json:"target_value,omitempty"`
	Unit        string   `json:"unit,omitempty"`
	Color       string   `json:"color,omitempty"`
	Icon        string   `json:"icon,omitempty"`
}

type HabitResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	HabitType   string    `json:"habit_type"`
	Frequency   string    `json:"frequency"`
	TargetValue *float64  `json:"target_value,omitempty"`
	Unit        string    `json:"unit,omitempty"`
	Color       string    `json:"color,omitempty"`
	Icon        string    `json:"icon,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type EntryRequest struct {
	Date  string   `json:"date"`
	Value *float64 `json:"value,omitempty"`
	Note  string   `json:"note,omitempty"`
}

type EntryResponse struct {
	ID        string    `json:"id"`
	HabitID   string    `json:"habit_id"`
	Date      string    `json:"date"`
	Value     *float64  `json:"value,omitempty"`
	Note      string    `json:"note,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ErrorResponse struct {
	Error ErrorPayload `json:"error"`
}

type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

package domain

import "errors"

var (
	ErrHabitNotFound    = errors.New("habit not found")
	ErrEntryNotFound    = errors.New("entry not found")
	ErrInvalidSchedule  = errors.New("invalid schedule")
	ErrInvalidDate      = errors.New("invalid date")
	ErrInvalidHabitType = errors.New("invalid habit type")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrConflict         = errors.New("conflict")
)

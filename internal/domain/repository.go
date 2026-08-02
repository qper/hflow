package domain

import "context"

// HabitRepository exposes persistence operations for habits.
type HabitRepository interface {
	CreateHabit(ctx context.Context, habit Habit) (Habit, error)
	UpdateHabit(ctx context.Context, habit Habit) (Habit, error)
	DeleteHabit(ctx context.Context, habitID, userID string) error
	ListHabits(ctx context.Context, userID string) ([]Habit, error)
	GetHabit(ctx context.Context, habitID, userID string) (Habit, error)
}

// EntryRepository exposes persistence operations for habit entries.
type EntryRepository interface {
	CreateEntry(ctx context.Context, entry Entry) (Entry, error)
	UpdateEntry(ctx context.Context, entry Entry) (Entry, error)
	DeleteEntry(ctx context.Context, entryID, userID string) error
	GetEntry(ctx context.Context, entryID, userID string) (Entry, error)
	ListEntries(ctx context.Context, userID string, from, to *string) ([]Entry, error)
}

// StreakRepository exposes access to streak state.
type StreakRepository interface {
	GetStreak(ctx context.Context, userID, habitID string) (Streak, error)
	SaveStreak(ctx context.Context, streak Streak) error
}

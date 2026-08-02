package domain

import "time"

// User represents a user account and is intentionally minimal for the foundation phase.
type User struct {
	ID           string
	Username     string
	Email        string
	DisplayName  string
	Timezone     string
	Theme        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
	ParentUserID *string
}

// Habit represents a repeating habit owned by a single user.
type Habit struct {
	ID          string
	UserID      string
	CategoryID  *string
	Name        string
	Description string
	HabitType   string
	Frequency   string
	Schedule    HabitSchedule
	TargetValue *float64
	Unit        string
	Color       string
	Icon        string
	SortOrder   int
	ArchivedAt  *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

// Entry represents a completed habit occurrence for a given date.
type Entry struct {
	ID        string
	UserID    string
	HabitID   string
	Date      time.Time
	Value     *float64
	Note      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// Streak is derived from entries and intentionally stored as a materialized view-like table
// for future analytics and query simplicity.
type Streak struct {
	ID            string
	UserID        string
	HabitID       string
	CurrentStreak int
	LongestStreak int
	UpdatedAt     time.Time
}

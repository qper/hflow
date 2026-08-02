package habit

import (
	"context"
	"fmt"
	"time"

	"github.com/qper/hflow/internal/domain"
)

// Service exposes the habit domain operations.
type Service struct {
	habits  domain.HabitRepository
	entries domain.EntryRepository
	streaks domain.StreakRepository
}

func NewService(habits domain.HabitRepository, entries domain.EntryRepository, streaks domain.StreakRepository) Service {
	return Service{habits: habits, entries: entries, streaks: streaks}
}

func (s Service) CreateHabit(ctx context.Context, userID string, habit domain.Habit) (domain.Habit, error) {
	if habit.Name == "" {
		return domain.Habit{}, domain.ErrInvalidSchedule
	}
	habit.UserID = userID
	return s.habits.CreateHabit(ctx, habit)
}

func (s Service) UpdateHabit(ctx context.Context, userID string, habit domain.Habit) (domain.Habit, error) {
	current, err := s.habits.GetHabit(ctx, habit.ID, userID)
	if err != nil {
		return domain.Habit{}, err
	}
	if current.UserID != userID {
		return domain.Habit{}, domain.ErrUnauthorized
	}
	habit.UserID = userID
	return s.habits.UpdateHabit(ctx, habit)
}

func (s Service) DeleteHabit(ctx context.Context, userID, habitID string) error {
	return s.habits.DeleteHabit(ctx, habitID, userID)
}

func (s Service) ListHabits(ctx context.Context, userID string) ([]domain.Habit, error) {
	return s.habits.ListHabits(ctx, userID)
}

func (s Service) RecordEntry(ctx context.Context, userID string, habitID string, date time.Time, value *float64, note string, timezone string) (domain.Entry, error) {
	habit, err := s.habits.GetHabit(ctx, habitID, userID)
	if err != nil {
		return domain.Entry{}, err
	}
	localDate := normalizeDate(date, timezone)
	entry := domain.Entry{UserID: userID, HabitID: habit.ID, Date: localDate, Value: value, Note: note}
	return s.entries.CreateEntry(ctx, entry)
}

func (s Service) DeleteEntry(ctx context.Context, userID, entryID string) error {
	return s.entries.DeleteEntry(ctx, entryID, userID)
}

func (s Service) Streaks(ctx context.Context, userID, habitID string, timezone string) (domain.Streak, error) {
	habit, err := s.habits.GetHabit(ctx, habitID, userID)
	if err != nil {
		return domain.Streak{}, err
	}
	entries, err := s.entries.ListEntries(ctx, userID, nil, nil)
	if err != nil {
		return domain.Streak{}, err
	}
	var relevant []domain.Entry
	for _, entry := range entries {
		if entry.HabitID == habit.ID {
			relevant = append(relevant, entry)
		}
	}
	current, longest := calculateStreak(relevant, habit, timezone)
	streak := domain.Streak{UserID: userID, HabitID: habit.ID, CurrentStreak: current, LongestStreak: longest}
	if err := s.streaks.SaveStreak(ctx, streak); err != nil {
		return domain.Streak{}, err
	}
	return streak, nil
}

func normalizeDate(date time.Time, timezone string) time.Time {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		loc = time.UTC
	}
	return time.Date(date.In(loc).Year(), date.In(loc).Month(), date.In(loc).Day(), 0, 0, 0, 0, loc)
}

func calculateStreak(entries []domain.Entry, habit domain.Habit, timezone string) (int, int) {
	return calculateStreakForDate(entries, habit, timezone, time.Now().UTC())
}

func calculateStreakForDate(entries []domain.Entry, habit domain.Habit, timezone string, today time.Time) (int, int) {
	if len(entries) == 0 {
		return 0, 0
	}
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		loc = time.UTC
	}
	present := map[string]struct{}{}
	for _, entry := range entries {
		normalized := normalizeDate(entry.Date, timezone)
		present[normalized.In(loc).Format("2006-01-02")] = struct{}{}
	}

	current := 0
	longest := 0
	cursor := normalizeDate(today, timezone)
	if len(present) == 0 {
		return 0, 0
	}
	for {
		key := cursor.In(loc).Format("2006-01-02")
		if _, ok := present[key]; ok {
			current++
			if current > longest {
				longest = current
			}
			cursor = cursor.AddDate(0, 0, -1)
			continue
		}
		if !isScheduledDay(habit, cursor, timezone) {
			cursor = cursor.AddDate(0, 0, -1)
			continue
		}
		break
	}
	return current, longest
}

func isScheduledDay(habit domain.Habit, day time.Time, timezone string) bool {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		loc = time.UTC
	}
	currentDay := day.In(loc)
	schedule := habit.Schedule
	if schedule.Type == "" {
		if habit.Frequency == "daily" || habit.Frequency == "weekly" {
			return true
		}
		return true
	}
	if schedule.Type == domain.ScheduleDaily {
		return true
	}
	if schedule.Type == domain.ScheduleCustom {
		for _, weekday := range schedule.Weekdays {
			if currentDay.Weekday() == weekday {
				return true
			}
		}
		return false
	}
	if schedule.Type == domain.ScheduleWeekly {
		if len(schedule.Weekdays) > 0 {
			for _, weekday := range schedule.Weekdays {
				if currentDay.Weekday() == weekday {
					return true
				}
			}
			return false
		}
		return true
	}
	return true
}

func (s Service) ListEntries(ctx context.Context, userID string, from, to *string) ([]domain.Entry, error) {
	return s.entries.ListEntries(ctx, userID, from, to)
}

func (s Service) GetEntry(ctx context.Context, userID, entryID string) (domain.Entry, error) {
	return s.entries.GetEntry(ctx, entryID, userID)
}

func (s Service) DeleteEntryByID(ctx context.Context, userID, entryID string) error {
	return s.entries.DeleteEntry(ctx, entryID, userID)
}

func (s Service) GetHabit(ctx context.Context, userID, habitID string) (domain.Habit, error) {
	return s.habits.GetHabit(ctx, habitID, userID)
}

func (s Service) CreateEntry(ctx context.Context, userID string, entry domain.Entry) (domain.Entry, error) {
	entry.UserID = userID
	return s.entries.CreateEntry(ctx, entry)
}

func (s Service) UpdateEntry(ctx context.Context, userID string, entry domain.Entry) (domain.Entry, error) {
	current, err := s.entries.GetEntry(ctx, entry.ID, userID)
	if err != nil {
		return domain.Entry{}, err
	}
	if current.UserID != userID {
		return domain.Entry{}, domain.ErrUnauthorized
	}
	entry.UserID = userID
	return s.entries.UpdateEntry(ctx, entry)
}

func (s Service) ValidateHabit(habit domain.Habit) error {
	if habit.Name == "" {
		return domain.ErrInvalidSchedule
	}
	if habit.HabitType != "boolean" && habit.HabitType != "numeric" && habit.HabitType != "duration" {
		return domain.ErrInvalidHabitType
	}
	return nil
}

func (s Service) RecordEntryForDate(ctx context.Context, userID, habitID string, date time.Time, timezone string) error {
	_, err := s.RecordEntry(ctx, userID, habitID, date, nil, "", timezone)
	return err
}

func (s Service) RecordEntryWithValue(ctx context.Context, userID, habitID string, date time.Time, value float64, timezone string) error {
	valuePtr := value
	_, err := s.RecordEntry(ctx, userID, habitID, date, &valuePtr, "", timezone)
	return err
}

func (s Service) GetStreak(ctx context.Context, userID, habitID string, timezone string) (domain.Streak, error) {
	return s.Streaks(ctx, userID, habitID, timezone)
}

func (s Service) Example() string {
	return fmt.Sprintf("habit service ready for %s", time.Now().UTC().Format(time.RFC3339))
}

package habit

import (
	"testing"
	"time"

	"github.com/qper/hflow/internal/domain"
)

func TestNormalizeDateUsesUserTimezone(t *testing.T) {
	date := time.Date(2026, 8, 2, 23, 30, 0, 0, time.UTC)
	normalized := normalizeDate(date, "America/New_York")
	if normalized.Year() != 2026 || normalized.Month() != time.August || normalized.Day() != 2 {
		t.Fatalf("expected date to stay in the user timezone, got %v", normalized)
	}
}

func TestCalculateStreakIgnoresNonScheduledDays(t *testing.T) {
	habit := domain.Habit{Frequency: "daily"}
	entries := []domain.Entry{{Date: time.Date(2026, 8, 2, 0, 0, 0, 0, time.UTC)}}
	current, longest := calculateStreak(entries, habit, "UTC")
	if current != 1 || longest != 1 {
		t.Fatalf("expected a single-day streak, got current=%d longest=%d", current, longest)
	}
}

func TestCalculateStreakHandlesTimeZoneChange(t *testing.T) {
	habit := domain.Habit{Frequency: "daily"}
	entries := []domain.Entry{{Date: time.Date(2026, 8, 2, 0, 0, 0, 0, time.UTC)}}
	today := time.Date(2026, 8, 2, 12, 0, 0, 0, time.UTC)
	current, longest := calculateStreakForDate(entries, habit, "America/Los_Angeles", today)
	if current != 0 || longest != 0 {
		t.Fatalf("expected the alternate timezone test to reflect the user's local-day boundary, got current=%d longest=%d", current, longest)
	}
}

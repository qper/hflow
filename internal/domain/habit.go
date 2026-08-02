package domain

import "time"

// ScheduleType describes how often a habit should be completed.
type ScheduleType string

const (
	ScheduleDaily  ScheduleType = "daily"
	ScheduleWeekly ScheduleType = "weekly"
	ScheduleCustom ScheduleType = "custom"
)

// HabitSchedule captures the recurrence definition for a habit.
type HabitSchedule struct {
	Type        ScheduleType
	Weekdays    []time.Weekday
	Occurrences int
}

// Validate returns an error when the schedule is inconsistent.
func (s HabitSchedule) Validate() error {
	switch s.Type {
	case ScheduleDaily:
		return nil
	case ScheduleWeekly:
		if s.Occurrences < 1 || s.Occurrences > 7 {
			return ErrInvalidSchedule
		}
		return nil
	case ScheduleCustom:
		if len(s.Weekdays) == 0 {
			return ErrInvalidSchedule
		}
		return nil
	default:
		return ErrInvalidSchedule
	}
}

// NormalizeWeekdays deduplicates and sorts weekdays.
func (s HabitSchedule) NormalizeWeekdays() HabitSchedule {
	seen := make(map[time.Weekday]struct{}, len(s.Weekdays))
	result := make([]time.Weekday, 0, len(s.Weekdays))
	for _, weekday := range s.Weekdays {
		if _, ok := seen[weekday]; ok {
			continue
		}
		seen[weekday] = struct{}{}
		result = append(result, weekday)
	}
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i] > result[j] {
				result[i], result[j] = result[j], result[i]
			}
		}
	}
	s.Type = ScheduleCustom
	return HabitSchedule{Type: ScheduleCustom, Weekdays: result, Occurrences: s.Occurrences}
}

package habit

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/qper/hflow/internal/domain"
)

// PostgresRepository is a concrete repository implementation for the domain layer.
type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) PostgresRepository {
	return PostgresRepository{db: db}
}

func (r PostgresRepository) CreateHabit(ctx context.Context, habit domain.Habit) (domain.Habit, error) {
	if err := r.db.QueryRowContext(ctx, `
		INSERT INTO habits (user_id, category_id, name, description, habit_type, frequency, target_value, unit, color, icon, sort_order, archived_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW(), NOW())
		RETURNING id, user_id, category_id, name, description, habit_type, frequency, target_value, unit, color, icon, sort_order, archived_at, created_at, updated_at, deleted_at
	`, habit.UserID, habit.CategoryID, habit.Name, habit.Description, habit.HabitType, habit.Frequency, habit.TargetValue, habit.Unit, habit.Color, habit.Icon, habit.SortOrder, habit.ArchivedAt).Scan(
		&habit.ID, &habit.UserID, &habit.CategoryID, &habit.Name, &habit.Description, &habit.HabitType, &habit.Frequency, &habit.TargetValue, &habit.Unit, &habit.Color, &habit.Icon, &habit.SortOrder, &habit.ArchivedAt, &habit.CreatedAt, &habit.UpdatedAt, &habit.DeletedAt,
	); err != nil {
		return domain.Habit{}, fmt.Errorf("create habit: %w", err)
	}
	return habit, nil
}

func (r PostgresRepository) UpdateHabit(ctx context.Context, habit domain.Habit) (domain.Habit, error) {
	if _, err := r.db.ExecContext(ctx, `
		UPDATE habits
		SET name = $1, description = $2, habit_type = $3, frequency = $4, target_value = $5, unit = $6, color = $7, icon = $8, sort_order = $9, archived_at = $10, updated_at = NOW()
		WHERE id = $11 AND user_id = $12 AND deleted_at IS NULL
	`, habit.Name, habit.Description, habit.HabitType, habit.Frequency, habit.TargetValue, habit.Unit, habit.Color, habit.Icon, habit.SortOrder, habit.ArchivedAt, habit.ID, habit.UserID); err != nil {
		return domain.Habit{}, fmt.Errorf("update habit: %w", err)
	}
	return habit, nil
}

func (r PostgresRepository) DeleteHabit(ctx context.Context, habitID, userID string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE habits
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
	`, habitID, userID)
	return err
}

func (r PostgresRepository) ListHabits(ctx context.Context, userID string) ([]domain.Habit, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, category_id, name, description, habit_type, frequency, target_value, unit, color, icon, sort_order, archived_at, created_at, updated_at, deleted_at
		FROM habits
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY sort_order, created_at
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("list habits: %w", err)
	}
	defer rows.Close()
	var habits []domain.Habit
	for rows.Next() {
		var habit domain.Habit
		if err := rows.Scan(&habit.ID, &habit.UserID, &habit.CategoryID, &habit.Name, &habit.Description, &habit.HabitType, &habit.Frequency, &habit.TargetValue, &habit.Unit, &habit.Color, &habit.Icon, &habit.SortOrder, &habit.ArchivedAt, &habit.CreatedAt, &habit.UpdatedAt, &habit.DeletedAt); err != nil {
			return nil, fmt.Errorf("scan habit: %w", err)
		}
		habits = append(habits, habit)
	}
	return habits, rows.Err()
}

func (r PostgresRepository) GetHabit(ctx context.Context, habitID, userID string) (domain.Habit, error) {
	var habit domain.Habit
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, category_id, name, description, habit_type, frequency, target_value, unit, color, icon, sort_order, archived_at, created_at, updated_at, deleted_at
		FROM habits
		WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
	`, habitID, userID).Scan(&habit.ID, &habit.UserID, &habit.CategoryID, &habit.Name, &habit.Description, &habit.HabitType, &habit.Frequency, &habit.TargetValue, &habit.Unit, &habit.Color, &habit.Icon, &habit.SortOrder, &habit.ArchivedAt, &habit.CreatedAt, &habit.UpdatedAt, &habit.DeletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Habit{}, domain.ErrHabitNotFound
		}
		return domain.Habit{}, fmt.Errorf("get habit: %w", err)
	}
	return habit, nil
}

func (r PostgresRepository) CreateEntry(ctx context.Context, entry domain.Entry) (domain.Entry, error) {
	if err := r.db.QueryRowContext(ctx, `
		INSERT INTO habit_entries (user_id, habit_id, entry_date, value, note, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, user_id, habit_id, entry_date, value, note, created_at, updated_at, deleted_at
	`, entry.UserID, entry.HabitID, entry.Date, entry.Value, entry.Note).Scan(&entry.ID, &entry.UserID, &entry.HabitID, &entry.Date, &entry.Value, &entry.Note, &entry.CreatedAt, &entry.UpdatedAt, &entry.DeletedAt); err != nil {
		return domain.Entry{}, fmt.Errorf("create entry: %w", err)
	}
	return entry, nil
}

func (r PostgresRepository) UpdateEntry(ctx context.Context, entry domain.Entry) (domain.Entry, error) {
	if _, err := r.db.ExecContext(ctx, `
		UPDATE habit_entries
		SET value = $1, note = $2, updated_at = NOW()
		WHERE id = $3 AND user_id = $4 AND deleted_at IS NULL
	`, entry.Value, entry.Note, entry.ID, entry.UserID); err != nil {
		return domain.Entry{}, fmt.Errorf("update entry: %w", err)
	}
	return entry, nil
}

func (r PostgresRepository) DeleteEntry(ctx context.Context, entryID, userID string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE habit_entries
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
	`, entryID, userID)
	return err
}

func (r PostgresRepository) GetEntry(ctx context.Context, entryID, userID string) (domain.Entry, error) {
	var entry domain.Entry
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, habit_id, entry_date, value, note, created_at, updated_at, deleted_at
		FROM habit_entries
		WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
	`, entryID, userID).Scan(&entry.ID, &entry.UserID, &entry.HabitID, &entry.Date, &entry.Value, &entry.Note, &entry.CreatedAt, &entry.UpdatedAt, &entry.DeletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Entry{}, domain.ErrEntryNotFound
		}
		return domain.Entry{}, fmt.Errorf("get entry: %w", err)
	}
	return entry, nil
}

func (r PostgresRepository) ListEntries(ctx context.Context, userID string, from, to *string) ([]domain.Entry, error) {
	query := `
		SELECT id, user_id, habit_id, entry_date, value, note, created_at, updated_at, deleted_at
		FROM habit_entries
		WHERE user_id = $1 AND deleted_at IS NULL
	`
	args := []any{userID}
	if from != nil {
		query += ` AND entry_date >= $2`
		args = append(args, *from)
	}
	if to != nil {
		query += ` AND entry_date <= $3`
		args = append(args, *to)
	}
	query += ` ORDER BY entry_date DESC`
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list entries: %w", err)
	}
	defer rows.Close()
	var entries []domain.Entry
	for rows.Next() {
		var entry domain.Entry
		if err := rows.Scan(&entry.ID, &entry.UserID, &entry.HabitID, &entry.Date, &entry.Value, &entry.Note, &entry.CreatedAt, &entry.UpdatedAt, &entry.DeletedAt); err != nil {
			return nil, fmt.Errorf("scan entry: %w", err)
		}
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}

func (r PostgresRepository) GetStreak(ctx context.Context, userID, habitID string) (domain.Streak, error) {
	var streak domain.Streak
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, habit_id, current_streak, longest_streak, updated_at
		FROM streaks
		WHERE user_id = $1 AND habit_id = $2
	`, userID, habitID).Scan(&streak.ID, &streak.UserID, &streak.HabitID, &streak.CurrentStreak, &streak.LongestStreak, &streak.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Streak{}, fmt.Errorf("streak not found: %w", err)
		}
		return domain.Streak{}, fmt.Errorf("get streak: %w", err)
	}
	return streak, nil
}

func (r PostgresRepository) SaveStreak(ctx context.Context, streak domain.Streak) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO streaks (user_id, habit_id, current_streak, longest_streak, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (user_id, habit_id) DO UPDATE SET
		current_streak = EXCLUDED.current_streak,
		longest_streak = EXCLUDED.longest_streak,
		updated_at = NOW()
	`, streak.UserID, streak.HabitID, streak.CurrentStreak, streak.LongestStreak)
	return err
}

func init() {
	_ = time.Now()
}

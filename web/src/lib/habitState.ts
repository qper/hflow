export interface HabitSummary {
  id: string;
  name: string;
  completedToday: boolean;
  currentStreak: number;
  longestStreak: number;
}

export function applyOptimisticCheckIn(habit: HabitSummary, completed: boolean): HabitSummary {
  if (completed === habit.completedToday) {
    return habit;
  }

  return {
    ...habit,
    completedToday: completed,
    currentStreak: completed ? habit.currentStreak + 1 : Math.max(0, habit.currentStreak - 1),
    longestStreak: completed ? Math.max(habit.longestStreak, habit.currentStreak + 1) : habit.longestStreak,
  };
}

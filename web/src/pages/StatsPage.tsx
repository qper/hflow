import { useMemo } from 'react';
import { useHabitsStore } from '../stores/habitsStore';

export default function StatsPage() {
  const { habits } = useHabitsStore();

  const total = habits.length;
  const done = habits.filter((habit) => habit.completedToday).length;
  const streaks = useMemo(() => habits.map((habit) => `${habit.name}: ${habit.currentStreak}`), [habits]);

  return (
    <section className="card stack">
      <h1>Stats</h1>
      <p>{done}/{total} completed today</p>
      <div className="heatmap">
        {streaks.map((entry) => (
          <div key={entry} className="heatmap-cell">
            {entry}
          </div>
        ))}
      </div>
    </section>
  );
}

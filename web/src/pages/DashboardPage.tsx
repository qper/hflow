import { useEffect, useMemo, useState } from 'react';
import { Link } from 'react-router-dom';
import { useHabitsStore } from '../stores/habitsStore';

export default function DashboardPage() {
  const { habits, loading, error, refreshHabits, checkInHabit } = useHabitsStore();
  const [pendingId, setPendingId] = useState<string | null>(null);

  useEffect(() => {
    void refreshHabits();
  }, [refreshHabits]);

  const summary = useMemo(() => {
    if (habits.length === 0) {
      return 'No habits yet';
    }
    const done = habits.filter((habit) => habit.completedToday).length;
    return `${done}/${habits.length} habits done today`;
  }, [habits]);

  return (
    <section className="card stack">
      <div className="row-between">
        <div>
          <h1>Dashboard</h1>
          <p>{summary}</p>
        </div>
        <Link to="/habits/new" className="button-link">
          New habit
        </Link>
      </div>

      {loading && <p>Loading habits…</p>}
      {error && <p className="error">{error}</p>}
      {!loading && habits.length === 0 && <p>No habits yet. Create the first one.</p>}

      <div className="habit-list">
        {habits.map((habit) => (
          <article className="habit-card" key={habit.id}>
            <div>
              <h3>{habit.name}</h3>
              <p>Current streak: {habit.currentStreak}</p>
            </div>
            <div className="habit-actions">
              <span className="pill">{habit.completedToday ? 'Done' : 'Pending'}</span>
              <button
                type="button"
                onClick={() => {
                  setPendingId(habit.id);
                  void checkInHabit(habit.id).finally(() => setPendingId(null));
                }}
                disabled={pendingId === habit.id}
              >
                {pendingId === habit.id ? 'Saving…' : habit.completedToday ? 'Undo' : 'Check in'}
              </button>
            </div>
          </article>
        ))}
      </div>
    </section>
  );
}

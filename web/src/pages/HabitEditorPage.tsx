import { FormEvent, useEffect, useMemo, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useHabitsStore } from '../stores/habitsStore';

export default function HabitEditorPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { draft, setDraft, createHabit, updateHabit } = useHabitsStore();
  const [error, setError] = useState<string | null>(null);

  const isEditing = useMemo(() => Boolean(id && id !== 'new'), [id]);

  useEffect(() => {
    if (!isEditing) {
      setDraft({
        name: '',
        description: '',
        habitType: 'daily',
        frequency: 'daily',
        unit: 'times',
        color: '#4f46e5',
        icon: '✨',
      });
    }
  }, [isEditing, setDraft]);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError(null);
    try {
      if (isEditing && id) {
        await updateHabit(id, draft);
      } else {
        await createHabit(draft);
      }
      navigate('/app');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unexpected error');
    }
  }

  return (
    <section className="card stack">
      <h1>{isEditing ? 'Edit habit' : 'Create habit'}</h1>
      <form onSubmit={handleSubmit} className="stack">
        <label>
          Name
          <input value={draft.name} onChange={(event) => setDraft({ ...draft, name: event.target.value })} required />
        </label>
        <label>
          Description
          <textarea value={draft.description} onChange={(event) => setDraft({ ...draft, description: event.target.value })} />
        </label>
        <label>
          Type
          <select value={draft.habitType} onChange={(event) => setDraft({ ...draft, habitType: event.target.value })}>
            <option value="daily">Daily</option>
            <option value="weekly">Weekly</option>
            <option value="custom">Custom</option>
          </select>
        </label>
        <label>
          Frequency
          <input value={draft.frequency} onChange={(event) => setDraft({ ...draft, frequency: event.target.value })} />
        </label>
        <label>
          Unit
          <input value={draft.unit} onChange={(event) => setDraft({ ...draft, unit: event.target.value })} />
        </label>
        <label>
          Color
          <input type="color" value={draft.color} onChange={(event) => setDraft({ ...draft, color: event.target.value })} />
        </label>
        {error && <p className="error">{error}</p>}
        <button type="submit">Save habit</button>
      </form>
    </section>
  );
}

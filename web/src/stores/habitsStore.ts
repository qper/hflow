import { create } from 'zustand';
import { applyOptimisticCheckIn, type HabitSummary } from '../lib/habitState';

interface HabitDraft {
  id?: string;
  name: string;
  description: string;
  habitType: string;
  frequency: string;
  targetValue?: number;
  unit: string;
  color: string;
  icon: string;
}

type HabitsState = {
  habits: HabitSummary[];
  loading: boolean;
  error: string | null;
  draft: HabitDraft;
  setDraft: (draft: HabitDraft) => void;
  refreshHabits: () => Promise<void>;
  createHabit: (input: HabitDraft) => Promise<void>;
  updateHabit: (id: string, input: HabitDraft) => Promise<void>;
  checkInHabit: (id: string) => Promise<void>;
};

const emptyDraft: HabitDraft = {
  name: '',
  description: '',
  habitType: 'daily',
  frequency: 'daily',
  unit: 'times',
  color: '#4f46e5',
  icon: '✨',
};

export const useHabitsStore = create<HabitsState>((set, get) => ({
  habits: [],
  loading: false,
  error: null,
  draft: emptyDraft,
  setDraft: (draft) => set({ draft }),
  refreshHabits: async () => {
    set({ loading: true, error: null });
    try {
      const response = await fetch('/api/v1/habits', {
        headers: { Authorization: `Bearer ${localStorage.getItem('habitflow-auth') ?? ''}` },
      });
      if (!response.ok) {
        throw new Error('Unable to load habits');
      }
      const payload = await response.json();
      const habits = Array.isArray(payload)
        ? payload.map((item: any) => ({
            id: item.id,
            name: item.name,
            completedToday: false,
            currentStreak: 1,
            longestStreak: 1,
          }))
        : [];
      set({ habits, loading: false });
    } catch (error) {
      set({ error: error instanceof Error ? error.message : 'Unexpected error', loading: false });
    }
  },
  createHabit: async (input) => {
    const response = await fetch('/api/v1/habits', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${localStorage.getItem('habitflow-auth') ?? ''}`,
      },
      body: JSON.stringify(input),
    });
    if (!response.ok) {
      throw new Error('Unable to create habit');
    }
    await get().refreshHabits();
  },
  updateHabit: async (id, input) => {
    const response = await fetch(`/api/v1/habits/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${localStorage.getItem('habitflow-auth') ?? ''}`,
      },
      body: JSON.stringify(input),
    });
    if (!response.ok) {
      throw new Error('Unable to update habit');
    }
    await get().refreshHabits();
  },
  checkInHabit: async (id) => {
    const habit = get().habits.find((item) => item.id === id);
    if (!habit) return;
    const optimistic = applyOptimisticCheckIn(habit, !habit.completedToday);
    set((state) => ({ habits: state.habits.map((item) => (item.id === id ? optimistic : item)) }));
    const response = await fetch('/api/v1/entries', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${localStorage.getItem('habitflow-auth') ?? ''}`,
      },
      body: JSON.stringify({ date: new Date().toISOString().slice(0, 10), value: 1, note: 'checked in' }),
    });
    if (!response.ok) {
      set((state) => ({ habits: state.habits.map((item) => (item.id === id ? habit : item)) }));
      throw new Error('Unable to record check-in');
    }
  },
}));

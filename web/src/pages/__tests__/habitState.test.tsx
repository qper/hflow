import { describe, expect, it, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import DashboardPage from '../DashboardPage';

vi.mock('../../stores/habitsStore', () => ({
  useHabitsStore: () => ({
    habits: [{ id: '1', name: 'Meditate', completedToday: false, currentStreak: 2, longestStreak: 3 }],
    loading: false,
    error: null,
    refreshHabits: vi.fn(),
    checkInHabit: vi.fn(),
  }),
}));

describe('DashboardPage', () => {
  it('renders the habit summary and CTA', () => {
    render(
      <MemoryRouter>
        <DashboardPage />
      </MemoryRouter>,
    );
    expect(screen.getByText('Dashboard')).toBeTruthy();
    expect(screen.getByText('Meditate')).toBeTruthy();
  });
});

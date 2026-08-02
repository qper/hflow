import { describe, expect, it } from 'vitest';
import { applyOptimisticCheckIn } from './habitState';

describe('applyOptimisticCheckIn', () => {
  it('increments the current streak when a habit is checked in', () => {
    const habit = {
      id: 'habit-1',
      name: 'Meditate',
      completedToday: false,
      currentStreak: 2,
      longestStreak: 3,
    };

    expect(applyOptimisticCheckIn(habit, true)).toEqual({
      ...habit,
      completedToday: true,
      currentStreak: 3,
      longestStreak: 3,
    });
  });
});

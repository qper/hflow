import { loadQueuedMutations, removeQueuedMutation, saveQueuedMutations } from './offlineQueue';

export function syncQueuedMutations() {
  if (typeof window === 'undefined' || !navigator.onLine) return Promise.resolve();
  const queue = loadQueuedMutations();
  if (queue.length === 0) return Promise.resolve();
  return Promise.all(
    queue.map(async (mutation) => {
      if (mutation.type !== 'check-in') return;
      const response = await fetch('/api/v1/entries', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${localStorage.getItem('habitflow-auth') ?? ''}`,
        },
        body: JSON.stringify(mutation.payload),
      });
      if (response.ok) {
        removeQueuedMutation(mutation.id);
      }
    }),
  ).then(() => {
    saveQueuedMutations(loadQueuedMutations());
  });
}

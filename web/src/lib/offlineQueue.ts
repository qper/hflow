export interface QueuedMutation {
  id: string;
  type: 'check-in';
  habitId: string;
  payload: { date: string; value: number; note: string };
  createdAt: string;
  status: 'pending' | 'synced';
}

const STORAGE_KEY = 'habitflow-offline-queue';

export function loadQueuedMutations(): QueuedMutation[] {
  if (typeof window === 'undefined') return [];
  try {
    const raw = window.localStorage.getItem(STORAGE_KEY);
    return raw ? JSON.parse(raw) : [];
  } catch {
    return [];
  }
}

export function saveQueuedMutations(queue: QueuedMutation[]) {
  if (typeof window === 'undefined') return;
  window.localStorage.setItem(STORAGE_KEY, JSON.stringify(queue));
}

export function enqueueMutation(mutation: QueuedMutation) {
  const queue = loadQueuedMutations();
  const next = [...queue, mutation];
  saveQueuedMutations(next);
  return next;
}

export function removeQueuedMutation(id: string) {
  const next = loadQueuedMutations().filter((item) => item.id !== id);
  saveQueuedMutations(next);
  return next;
}

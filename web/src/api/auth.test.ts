import { describe, expect, it } from 'vitest';
import { getApiErrorMessage } from './auth';

describe('getApiErrorMessage', () => {
  it('returns the server-provided error text when present', () => {
    expect(getApiErrorMessage({ error: { error: 'password must be at least 12 characters long' } }, 'Unable to register')).toBe('password must be at least 12 characters long');
  });

  it('falls back to the provided default when the server did not return a message', () => {
    expect(getApiErrorMessage({ error: { status: 400 } }, 'Unable to register')).toBe('Unable to register');
  });
});

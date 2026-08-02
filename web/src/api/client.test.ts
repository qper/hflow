import { describe, expect, it } from 'vitest';
import { resolveApiUrl } from './client';

describe('resolveApiUrl', () => {
  it('returns the path unchanged when it already includes the API prefix', () => {
    expect(resolveApiUrl('/api/v1/auth/register')).toBe('/api/v1/auth/register');
  });

  it('deduplicates a repeated API prefix', () => {
    expect(resolveApiUrl('/api/v1/api/v1/auth/register', '/api/v1')).toBe('/api/v1/auth/register');
  });

  it('prepends the configured base URL for relative API paths', () => {
    expect(resolveApiUrl('/auth/register', '/api/v1')).toBe('/api/v1/auth/register');
  });
});

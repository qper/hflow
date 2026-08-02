import { apiClient } from './client';
import { useAuthStore } from '../stores/authStore';

interface TokenPayload {
  access_token: string;
  refresh_token: string;
  token_type: string;
  expires_at?: string;
}

interface ApiErrorPayload {
  error?: string | { error?: string };
}

export function getApiErrorMessage(response: unknown, fallback: string) {
  if (response && typeof response === 'object' && 'error' in response) {
    const payload = response.error as ApiErrorPayload | undefined;
    if (typeof payload?.error === 'string' && payload.error.trim()) {
      return payload.error;
    }
    if (payload && typeof payload === 'object' && 'error' in payload) {
      const nested = payload.error;
      if (typeof nested === 'string' && nested.trim()) {
        return nested;
      }
    }
  }
  return fallback;
}

export async function login(username: string, password: string) {
  const response = await apiClient.POST('/api/v1/auth/login', {
    body: { username, password },
  });

  const error = 'error' in response ? response.error : undefined;
  if (error) {
    throw new Error(getApiErrorMessage(response, 'Unable to login'));
  }

  const data = response.data as TokenPayload | undefined;
  if (!data?.access_token || !data.refresh_token) {
    throw new Error('Unexpected auth response');
  }

  useAuthStore.getState().login(data.access_token, data.refresh_token);
  return data;
}

export async function register(username: string, email: string, password: string) {
  const response = await apiClient.POST('/api/v1/auth/register', {
    body: { username, email, password },
  });

  const error = 'error' in response ? response.error : undefined;
  if (error) {
    throw new Error(getApiErrorMessage(response, 'Unable to register'));
  }

  return response.data as TokenPayload | undefined;
}

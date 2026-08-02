import createClient from 'openapi-fetch';
import type { paths } from './generated/schema';
import { useAuthStore } from '../stores/authStore';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? '/api/v1';

export const apiClient = createClient<paths>({
  baseUrl: API_BASE_URL,
});

apiClient.use({
  async onRequest({ request }) {
    const { accessToken } = useAuthStore.getState();
    if (accessToken) {
      request.headers.set('Authorization', `Bearer ${accessToken}`);
    }
    return request;
  },
  async onResponse({ response }) {
    if (response.status !== 401) {
      return response;
    }
    const { refreshToken, logout } = useAuthStore.getState();
    if (!refreshToken) {
      logout();
      if (typeof window !== 'undefined') {
        window.location.assign('/login');
      }
      return response;
    }

    const refreshResponse = await fetch(`${API_BASE_URL}/auth/refresh`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: refreshToken }),
    });

    if (!refreshResponse.ok) {
      logout();
      if (typeof window !== 'undefined') {
        window.location.assign('/login');
      }
      return response;
    }

    const payload = await refreshResponse.json();
    useAuthStore.getState().login(payload.access_token, payload.refresh_token);
    return response;
  },
});

import createClient from 'openapi-fetch';
import type { paths } from './generated/schema';
import { useAuthStore } from '../stores/authStore';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? '/api/v1';

export function resolveApiUrl(path: string, baseUrl = API_BASE_URL) {
  const doubledPrefix = `${baseUrl}${baseUrl}`;
  if (path.startsWith(doubledPrefix)) {
    return path.replace(doubledPrefix, baseUrl);
  }
  if (path.startsWith(baseUrl)) {
    return path;
  }
  if (path.startsWith('/')) {
    return `${baseUrl}${path}`;
  }
  return `${baseUrl}/${path}`;
}

export const apiClient = createClient<paths>({
  baseUrl: API_BASE_URL,
});

apiClient.use({
  async onRequest({ request }) {
    const url = new URL(request.url);
    const normalizedPath = resolveApiUrl(url.pathname);
    if (normalizedPath !== url.pathname) {
      url.pathname = normalizedPath;
      request = new Request(url.toString(), request);
    }

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

    const refreshResponse = await fetch(resolveApiUrl('/auth/refresh'), {
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

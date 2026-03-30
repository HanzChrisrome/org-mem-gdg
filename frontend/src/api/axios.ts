import axios, { type AxiosError, type InternalAxiosRequestConfig } from "axios";
import { clearStoredAuth, getAccessToken } from "../lib/token-storage";
import { refreshAccessTokenFromStorage } from "./refresh";

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
});

interface RetryableRequestConfig extends InternalAxiosRequestConfig {
  _retry?: boolean;
}

let refreshPromise: Promise<string | null> | null = null;

function shouldSkipRefresh(url?: string) {
  if (!url) return false;
  return (
    url.includes("/login") ||
    url.includes("/refresh") ||
    url.includes("/register")
  );
}

// Automatically attach access token
api.interceptors.request.use((config) => {
  const token = getAccessToken();
  if (token) {
    config.headers = config.headers ?? {};
    config.headers["Authorization"] = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config as RetryableRequestConfig | undefined;

    if (
      !originalRequest ||
      error.response?.status !== 401 ||
      shouldSkipRefresh(originalRequest.url)
    ) {
      return Promise.reject(error);
    }

    if (originalRequest._retry) {
      clearStoredAuth();
      return Promise.reject(error);
    }
    originalRequest._retry = true;

    if (!refreshPromise) {
      refreshPromise = refreshAccessTokenFromStorage().finally(() => {
        refreshPromise = null;
      });
    }

    const nextAccessToken = await refreshPromise;
    if (!nextAccessToken) {
      clearStoredAuth();
      return Promise.reject(error);
    }

    originalRequest.headers = originalRequest.headers ?? {};
    originalRequest.headers["Authorization"] = `Bearer ${nextAccessToken}`;

    return api(originalRequest);
  },
);

export default api;

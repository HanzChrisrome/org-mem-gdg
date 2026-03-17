import { jwtDecode } from "jwt-decode";
import { toast } from "sonner";
import api from "./axios";

interface LoginData {
  identifier: string;
  password: string;
}

function clearStoredAuth() {
  localStorage.removeItem("access_token");
  localStorage.removeItem("refresh_token");
  localStorage.removeItem("user_id");
}

function getSessionIdFromRefreshToken(
  refreshToken: string | null,
): string | null {
  if (!refreshToken) return null;

  const [sessionId] = refreshToken.split(".", 2);
  return sessionId || null;
}

// Check if access token is still valid
function isTokenValid(token: string) {
  try {
    const decoded = jwtDecode<{ exp?: number }>(token);
    if (!decoded.exp) {
      return false;
    }
    return decoded.exp * 1000 > Date.now();
  } catch {
    return false;
  }
}

// Login function
export async function login(data: LoginData): Promise<void> {
  const response = await api.post("/login", data);
  const { token, user_id } = response.data;

  localStorage.setItem("access_token", token.access_token);
  localStorage.setItem("refresh_token", token.refresh_token);
  localStorage.setItem("user_id", user_id);
}

// Logout function
export async function logout() {
  const refreshToken = localStorage.getItem("refresh_token");
  const sessionId = getSessionIdFromRefreshToken(refreshToken);

  try {
    if (sessionId) {
      await api.post("/logout", {
        session_id: sessionId,
      });
    }
  } catch (error) {
    console.warn("Failed to revoke session on server:", error);
  } finally {
    clearStoredAuth();
  }

  toast.success("Logged out successfully!");
  window.location.href = "/login";
}

// Startup auth check (validate token or refresh if expired)
export async function initAuth(): Promise<boolean> {
  const accessToken = localStorage.getItem("access_token");
  const refreshToken = localStorage.getItem("refresh_token");

  if (!accessToken && !refreshToken) return false;

  try {
    if (accessToken && isTokenValid(accessToken)) return true;

    // Attempt refresh
    if (refreshToken) {
      const response = await api.post("/refresh", {
        refresh_token: refreshToken,
      });

      const nextToken = response.data?.token;
      if (nextToken?.access_token) {
        localStorage.setItem("access_token", nextToken.access_token);
      }
      if (nextToken?.refresh_token) {
        localStorage.setItem("refresh_token", nextToken.refresh_token);
      }
      return true;
    }

    // Tokens invalid
    clearStoredAuth();
    return false;
  } catch {
    clearStoredAuth();
    return false;
  }
}

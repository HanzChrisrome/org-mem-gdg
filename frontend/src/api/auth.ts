import { jwtDecode } from "jwt-decode";
import { toast } from "sonner";
import {
  clearStoredAuth,
  getAccessToken,
  getRefreshToken,
  storeAuthSession,
} from "../lib/token-storage";
import api from "./axios";
import { refreshAccessTokenFromStorage } from "./refresh";

interface LoginData {
  identifier: string;
  password: string;
}

// Check if access token is still valid
function isTokenValid(token: string) {
  try {
    const decoded = jwtDecode<{ exp?: number }>(token);
    console.log("Decoded token:", decoded);
    if (!decoded.exp) {
      return false;
    }
    return decoded.exp * 1000 > Date.now();
  } catch {
    return false;
  }
}

// Login function
export async function login(data: LoginData): Promise<boolean> {
  try {
    const response = await api.post("/login", data);
    console.log("Login response:", response.data);

    const { token, user_id } = response.data;
    storeAuthSession(token.access_token, token.refresh_token, user_id);
    return true;
  } catch (error: unknown) {
    clearStoredAuth();

    const apiError = error as {
      response?: {
        status?: number;
        data?: { error?: string; message?: string };
      };
    };
    console.error("Login failed:", {
      status: apiError.response?.status,
      data: apiError.response?.data,
    });

    const message =
      apiError.response?.data?.error ||
      apiError.response?.data?.message ||
      "Login failed. Check credentials.";
    toast.error(message);

    return false;
  }
}

// Logout function
export async function logout() {
  const refreshToken = getRefreshToken();

  try {
    await api.post("/logout", {
      refresh_token: refreshToken ?? "",
    });
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
  const accessToken = getAccessToken();
  const refreshToken = getRefreshToken();

  if (!accessToken && !refreshToken) return false;

  try {
    if (accessToken && isTokenValid(accessToken)) return true;

    // Attempt refresh
    if (refreshToken) {
      console.log("Access token expired, attempting refresh...");
      const nextAccessToken = await refreshAccessTokenFromStorage();
      return Boolean(nextAccessToken);
    }

    // Tokens invalid
    clearStoredAuth();
    return false;
  } catch {
    clearStoredAuth();
    return false;
  }
}

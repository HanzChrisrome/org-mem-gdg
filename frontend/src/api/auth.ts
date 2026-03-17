import { jwtDecode } from "jwt-decode";
import { toast } from "sonner";
import api from "./axios";

interface LoginData {
  identifier: string;
  password: string;
}

interface LoginResponse {
  token: {
    access_token: string;
    refresh_token: string;
    token_type: string;
    expires_in: number;
  };
  user_id: string;
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
export async function login(data: LoginData): Promise<LoginResponse> {
  const response = await api.post("/login", data);
  const { token, user_id } = response.data;

  localStorage.setItem("access_token", token.access_token);
  localStorage.setItem("refresh_token", token.refresh_token);
  localStorage.setItem("user_id", user_id);

  toast.success("Login successful!");
  return { token, user_id };
}

// Logout function
export function logout() {
  localStorage.removeItem("access_token");
  localStorage.removeItem("refresh_token");
  localStorage.removeItem("user_id");

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
      const parts = refreshToken.split(".");
      const refreshTokenID = parts[0];
      const response = await api.post("/refresh", {
        refresh_token_id: refreshTokenID,
        refresh_token: refreshToken,
      });
      const { access_token } = response.data.token;
      localStorage.setItem("access_token", access_token);
      return true;
    }

    // Tokens invalid
    localStorage.clear();
    return false;
  } catch {
    localStorage.clear();
    return false;
  }
}

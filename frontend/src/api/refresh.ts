import axios from "axios";
import { getRefreshToken, updateStoredTokens } from "../lib/token-storage";

export async function refreshAccessTokenFromStorage(): Promise<string | null> {
  const refreshToken = getRefreshToken();
  if (!refreshToken) {
    return null;
  }

  try {
    const response = await axios.post(
      `${import.meta.env.VITE_API_URL}/refresh`,
      {
        refresh_token: refreshToken,
      },
    );

    const nextToken = response.data?.token;
    const nextAccessToken = nextToken?.access_token;
    const nextRefreshToken = nextToken?.refresh_token;

    if (!nextAccessToken) {
      return null;
    }

    updateStoredTokens(nextAccessToken, nextRefreshToken);
    return nextAccessToken;
  } catch {
    return null;
  }
}

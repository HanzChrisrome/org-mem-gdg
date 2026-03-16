import type { FormEvent } from "react";
import { useMemo, useState } from "react";
import "./App.css";

type ApiResult = {
  status: number;
  ok: boolean;
  data: unknown;
};

const pretty = (value: unknown): string => JSON.stringify(value, null, 2);

function App() {
  const [baseUrl, setBaseUrl] = useState("http://localhost:8080");

  const [registerForm, setRegisterForm] = useState({
    name: "Test Member",
    email: "member@example.com",
    studentID: "2026-0001",
    password: "StrongPass123!",
  });

  const [loginForm, setLoginForm] = useState({
    identifier: "member@example.com",
    password: "StrongPass123!",
  });

  const [accessToken, setAccessToken] = useState("");
  const [refreshToken, setRefreshToken] = useState("");
  const [sessionID, setSessionID] = useState("");
  const [busy, setBusy] = useState(false);
  const [result, setResult] = useState<ApiResult | null>(null);

  const authHeader = useMemo(
    () => (accessToken ? `Bearer ${accessToken}` : "<none>"),
    [accessToken],
  );

  const callApi = async (
    path: string,
    method: "GET" | "POST",
    payload?: unknown,
    includeAuth = false,
  ): Promise<ApiResult> => {
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
    };

    if (includeAuth && accessToken) {
      headers.Authorization = `Bearer ${accessToken}`;
    }

    const response = await fetch(`${baseUrl}${path}`, {
      method,
      headers,
      body: payload ? JSON.stringify(payload) : undefined,
    });

    const text = await response.text();
    let data: unknown = text;
    try {
      data = text ? JSON.parse(text) : null;
    } catch {
      data = text;
    }

    return {
      status: response.status,
      ok: response.ok,
      data,
    };
  };

  const runAction = async (runner: () => Promise<ApiResult>) => {
    setBusy(true);
    try {
      const output = await runner();
      setResult(output);
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      setResult({
        status: 0,
        ok: false,
        data: { error: message },
      });
    } finally {
      setBusy(false);
    }
  };

  const onRegister = (event: FormEvent) => {
    event.preventDefault();
    runAction(() => callApi("/api/register", "POST", registerForm));
  };

  const onLogin = (event: FormEvent) => {
    event.preventDefault();
    runAction(async () => {
      const output = await callApi("/api/login", "POST", loginForm);

      if (output.ok && typeof output.data === "object" && output.data) {
        const data = output.data as {
          access_token?: string;
          refresh_token?: string;
          token?: { access_token?: string; refresh_token?: string };
        };

        const nextAccess = data.access_token ?? data.token?.access_token ?? "";
        const nextRefresh =
          data.refresh_token ?? data.token?.refresh_token ?? "";

        if (nextAccess) {
          setAccessToken(nextAccess);
        }
        if (nextRefresh) {
          setRefreshToken(nextRefresh);
          const split = nextRefresh.split(".");
          if (split.length >= 2) {
            setSessionID(split[0]);
          }
        }
      }

      return output;
    });
  };

  const onRefresh = () => {
    runAction(() =>
      callApi("/api/refresh", "POST", {
        session_id: sessionID,
        refresh_token: refreshToken,
      }),
    );
  };

  const onLogout = () => {
    runAction(() =>
      callApi(
        "/api/logout",
        "POST",
        {
          session_id: sessionID,
        },
        true,
      ),
    );
  };

  return (
    <main className="auth-lab">
      <section className="intro card">
        <p className="eyebrow">Org Membership Platform</p>
        <h1>Authentication Test Console</h1>
        <p className="subtext">
          Use this page to test register, login, refresh, and logout against
          your Go backend.
        </p>
        <label className="field">
          <span>Backend Base URL</span>
          <input
            value={baseUrl}
            onChange={(e) => setBaseUrl(e.target.value)}
            placeholder="http://localhost:8080"
          />
        </label>
      </section>

      <section className="grid">
        <form className="card" onSubmit={onRegister}>
          <h2>Register</h2>
          <label className="field">
            <span>Name</span>
            <input
              value={registerForm.name}
              onChange={(e) =>
                setRegisterForm((prev) => ({ ...prev, name: e.target.value }))
              }
            />
          </label>
          <label className="field">
            <span>Email</span>
            <input
              type="email"
              value={registerForm.email}
              onChange={(e) =>
                setRegisterForm((prev) => ({ ...prev, email: e.target.value }))
              }
            />
          </label>
          <label className="field">
            <span>Student ID</span>
            <input
              value={registerForm.studentID}
              onChange={(e) =>
                setRegisterForm((prev) => ({ ...prev, studentID: e.target.value }))
              }
            />
          </label>
          <label className="field">
            <span>Password</span>
            <input
              type="password"
              value={registerForm.password}
              onChange={(e) =>
                setRegisterForm((prev) => ({ ...prev, password: e.target.value }))
              }
            />
          </label>
          <button className="action" type="submit" disabled={busy}>
            {busy ? "Working..." : "POST /api/register"}
          </button>
        </form>

        <form className="card" onSubmit={onLogin}>
          <h2>Login</h2>
          <label className="field">
            <span>Identifier (email or student ID)</span>
            <input
              value={loginForm.identifier}
              onChange={(e) =>
                setLoginForm((prev) => ({
                  ...prev,
                  identifier: e.target.value,
                }))
              }
            />
          </label>
          <label className="field">
            <span>Password</span>
            <input
              type="password"
              value={loginForm.password}
              onChange={(e) =>
                setLoginForm((prev) => ({ ...prev, password: e.target.value }))
              }
            />
          </label>
          <button className="action" type="submit" disabled={busy}>
            {busy ? "Working..." : "POST /api/login"}
          </button>
        </form>

        <div className="card">
          <h2>Session Actions</h2>
          <label className="field">
            <span>Session ID</span>
            <input
              value={sessionID}
              onChange={(e) => setSessionID(e.target.value)}
            />
          </label>
          <label className="field">
            <span>Refresh Token</span>
            <input
              value={refreshToken}
              onChange={(e) => setRefreshToken(e.target.value)}
            />
          </label>
          <div className="actions-row">
            <button className="action" type="button" onClick={onRefresh} disabled={busy}>
              POST /api/refresh
            </button>
            <button className="action ghost" type="button" onClick={onLogout} disabled={busy}>
              POST /api/logout
            </button>
          </div>
          <p className="hint">Authorization: {authHeader}</p>
        </div>
      </section>

      <section className="card response">
        <h2>Last API Response</h2>
        <div className="meta">
          <span className={result?.ok ? "ok" : "fail"}>
            {result ? `HTTP ${result.status}` : "No request yet"}
          </span>
        </div>
        <pre>{result ? pretty(result.data) : "{}"}</pre>
      </section>
    </main>
  );
}

export default App;

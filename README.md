# Organization Membership Management Web Application

A full-stack web application for managing organization memberships, payment verification, and member status workflows.

<details>
<summary><strong>Project Purpose</strong></summary>

This project is built to help school organization executives manage members in one centralized platform.

It focuses on:

- member registration and record management (create, read, update, delete)
- membership payment tracking and proof-of-payment review
- member eligibility/status visibility
- reducing manual errors from spreadsheet-based workflows

The stack uses:

- **Frontend:** React + TypeScript + Vite + Tailwind CSS
- **Backend:** Go (HTTP server)
- **Database:** PostgreSQL (Supabase-compatible connection)

</details>

## Project Structure

```text
org-man-app/
|-- backend/
|   |-- cmd/
|   |   `-- api/
|   |       `-- main.go              # backend entrypoint
|   |-- internal/
|   |   |-- config/                  # environment/config loading
|   |   |-- database/                # database connection setup
|   |   |-- docs/
|   |   |-- handlers/
|   |   |-- middleware/
|   |   |-- routes/
|   |   |-- services/
|   |   `-- utils/
|   |-- .env.example                 # backend environment template
|   `-- go.mod
|-- frontend/
|   |-- public/
|   |-- src/                         # React app source
|   |-- package.json
|   `-- vite.config.ts
|-- docs/
|   `-- main.md                      # project requirements/reference doc
|-- Justfile                         # dev/setup/lint convenience tasks
`-- commitlint.config.cjs
```

## Setup Guide

### 1. Prerequisites

Install the following first:

- Go (1.25+ recommended from `backend/go.mod`)
- Node.js (LTS recommended)
- npm
- PostgreSQL connection (Supabase URL is supported)
- Optional: `just` command runner (to use project shortcut tasks)
- Optional: `pre-commit` (used by setup hooks)

### 2. Clone and Enter Project

```bash
git clone https://github.com/HanzChrisrome/org-mem-gdg.git
cd org-man-app
```

### 3. Configure Backend Environment

From the `backend` folder, create a local `.env` file based on `.env.example`.

Required variables:

```env
DATABASE_URL=postgres://your_user:your_password@your_project_host:5432/your_db?sslmode=require
PORT=8080
```

Notes:

- `DATABASE_URL` is required by the backend startup.
- `PORT` is used by the HTTP server (example: `8080`).

### 4. Install Dependencies

#### Option A: Using `just` (recommended)

```bash
just setup
```

This runs:

- tool installation (golangci-lint, goimports)
- frontend dependency install
- backend module tidy/verify
- pre-commit hooks install

#### Option B: Manual setup

```bash
cd frontend && npm install
cd ../backend && go mod tidy && go mod verify
```

### 5. Run the Project

#### Option A: Run both with `just`

```bash
just dev
```

#### Option B: Run services separately

Backend:

```bash
cd backend
go run cmd/api/main.go
```

Frontend:

```bash
cd frontend
npm run dev
```

### 6. Verify Services

- Backend health check: `http://localhost:<PORT>/health` (returns `OK`)
- Frontend dev server: shown in Vite terminal output (commonly `http://localhost:5173`)

## Useful Commands

From repository root:

```bash
just --list
```

Common tasks:

- `just backend-run`
- `just backend-smoke-endpoints`
- `just backend-lint`
- `just backend-fmt`
- `just frontend-dev`
- `just frontend-lint`
- `just frontend-format`

## Endpoint Smoke Test

Use the smoke script to run an end-to-end validation of all currently registered API endpoints.

### What the test does

The script in [scripts/smoke-endpoints.ps1](scripts/smoke-endpoints.ps1) executes this flow:

1. Verify public endpoints: `/swagger/index.html`, `/health`
2. Register and login an executive account
3. Refresh token using `refresh_token_id` and `refresh_token`
4. Call protected member endpoints (create, list, get, update, delete)
5. Call protected executive endpoints (create, list, get, update, delete)
6. Revoke a session, then logout

Each request validates expected HTTP status codes and fails immediately on mismatch.

### How to run

1. Start backend:

```bash
just backend-run
```

2. Run smoke test from repository root:

```bash
just backend-smoke-endpoints
```

You can also run the script directly:

```powershell
powershell -ExecutionPolicy Bypass -File scripts/smoke-endpoints.ps1
```

### Trace logs and trace IDs

While smoke tests run, payload tracing records each request/response pair with a unique `trace_id`.

- Default trace file path: [backend/logs/payload-trace.log](backend/logs/payload-trace.log)
- Response header includes `X-Trace-ID` for client-to-log correlation
- The log payload also includes `trace_id` so you can search exact request records

Example workflow:

1. Capture `X-Trace-ID` from a response
2. Search the same ID in [backend/logs/payload-trace.log](backend/logs/payload-trace.log)
3. Inspect request/response payload, status, and headers for that specific trace

### Trace configuration

Configure tracing behavior with environment variables:

- `TRACE_ENABLED` (default: `true`)
- `TRACE_REQUEST_BODY` (default: `true`)
- `TRACE_RESPONSE_BODY` (default: `true`)
- `TRACE_HEADERS` (default: `Content-Type,X-Trace-ID`)
- `TRACE_EXCLUDE_PATHS` (comma-separated, supports suffix wildcard like `/swagger/*`)
- `TRACE_MAX_BODY_BYTES` (default: `8192`)
- `TRACE_FILE_PATH` (default: `logs/payload-trace.log`)

## Notes

- The current backend fails fast if `DATABASE_URL` is not set.
- The root `docs/main.md` contains detailed business and functional requirements.
- The `frontend/README.md` is the default Vite template and can be ignored in favor of this root README.

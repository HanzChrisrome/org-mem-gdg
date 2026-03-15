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
git clone <your-repo-url>
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
- `just backend-lint`
- `just backend-fmt`
- `just frontend-dev`
- `just frontend-lint`
- `just frontend-format`

## Notes

- The current backend fails fast if `DATABASE_URL` is not set.
- The root `docs/main.md` contains detailed business and functional requirements.
- The `frontend/README.md` is the default Vite template and can be ignored in favor of this root README.

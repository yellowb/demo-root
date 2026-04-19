# Agent Harness Demo Baseline

This repository is the baseline Todo List project for the Agent Harness live demo. It already supports CRUD and SQLite persistence, but it intentionally does **not** include `priority + filtering` yet.

## Tech Stack

- Frontend: React 18, Vite 5, TypeScript 5
- Backend: Go, Gin, `database/sql`
- Database: SQLite via `modernc.org/sqlite`

## Project Structure

```text
demo-root/
├── Makefile
├── scripts/
├── frontend/
└── backend/
```

- `frontend/`: SPA UI, API client, Todo page and styles
- `backend/`: Go API server, SQLite bootstrap, Todo repository and handlers
- `scripts/`: short command wrappers for local dev, tests, and DB reset

## Local Run

### Prerequisites

- Node.js
- Go

### Install Dependencies

```bash
make setup
```

### Start Both Services

```bash
make dev
```

- Frontend: [http://localhost:5173](http://localhost:5173)
- Backend health check: [http://localhost:8080/api/health](http://localhost:8080/api/health)
- On macOS or Linux desktops with `open` or `xdg-open` available, `make dev` will automatically open the frontend in your browser after `http://localhost:5173/` becomes reachable.
- To disable that behavior, run: `AUTO_OPEN_BROWSER=0 make dev`

## Test and Validation

```bash
make test
```

This runs:

- `go test ./...`
- `npm run typecheck`
- `npm run build`

## Reset the Database

```bash
make reset-db
```

This removes the local SQLite database and reboots it with demo seed data.

## Contributor Workflow

- First-time setup: `make setup`
- Local development: `make dev`
- Required verification before finishing changes: `make test`
- Reset demo data before a presentation: `make reset-db`

See [CONTRIBUTING.md](/Users/yellowb/ppt/demo-root/CONTRIBUTING.md) for the expected workflow and repo constraints.

## Seed Data

The app starts with 8 seeded todos so the page is not blank during the live demo. The content is deliberately chosen to make the later `priority + filtering` upgrade visually obvious.

## Live Demo Reminder

This baseline project is intentionally missing:

- `priority`
- filtering
- search
- tags
- due dates
- auth

The next live demo step is to ask an AI Coding Agent to add `priority + filtering` across the backend, persistence layer, tests, and UI.

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
├── docs/
├── .codex/
├── scripts/
├── frontend/
└── backend/
```

- `frontend/`: SPA UI, API client, Todo page and styles
- `backend/`: Go API server, SQLite bootstrap, Todo repository and handlers
- `docs/`: live demo context for fresh Codex sessions
- `.codex/`: repo-local Codex hook configuration
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

## Lint

```bash
make lint
```

This runs the pinned repo-local `golangci-lint` against the Go backend. `make setup` installs that pinned version into `.bin/` when needed.

## Test and Validation

```bash
make test
```

This runs:

- `golangci-lint run ./...` for the backend
- `go test ./...`
- `npm run typecheck`
- `npm run build`

It also validates the repo-local Codex hook configuration before running the app checks.

## Live Demo Context

Before starting a fresh Codex session for the live demo, read [docs/live-demo-context.md](/Users/yellowb/ppt/demo-root/docs/live-demo-context.md). It explains how this Todo app supports the Agent Harness talk, why the baseline intentionally lacks `priority + filtering`, and which prompt to use for the feature demo.

## Codex Hook Verification

This repository includes a repo-local Codex `Stop` hook in [.codex/hooks.json](/Users/yellowb/ppt/demo-root/.codex/hooks.json). When the Codex App is opened at `/Users/yellowb/ppt/demo-root` and there are repository changes, the hook runs `make test` before the final response.

This is a demo-oriented Harness gate, not a 100% impossible-to-bypass security boundary. It makes validation visible and repeatable during the live demo, while CI remains the stronger remote verification layer after changes are pushed.

## Reset the Database

```bash
make reset-db
```

This removes the local SQLite database and reboots it with demo seed data.

## Contributor Workflow

- First-time setup: `make setup`
- Local development: `make dev`
- Backend lint: `make lint`
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

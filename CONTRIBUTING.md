# Contributing

## Quick Start

1. Run `make setup`
2. Run `make dev` while developing
3. Run `make test` before you consider the work finished
4. Run `make reset-db` if you need to restore the demo seed data

## Repo Constraints

- This repository is the **baseline** Todo List app for the Agent Harness live demo.
- Do **not** add `priority`, filtering, search, tags, due dates, auth, or collaboration unless the task explicitly asks for it.
- Keep the repo lightweight:
  - no Docker
  - no external cloud services
  - no MySQL / Redis / Postgres

## Change Expectations

- Keep frontend work inside `frontend/`
- Keep backend API and persistence work inside `backend/`
- If a change affects request/response behavior, update both sides:
  - backend handler / store / tests
  - frontend API client / UI
- If a change affects startup or validation, update both:
  - `README.md`
  - `AGENTS.md`
  - `Makefile` or `scripts/`

## Validation

- Required command: `make test`
- The expected checks are:
  - `go test ./...`
  - `npm run typecheck`
  - `npm run build`

## Demo Hygiene

- Seed data matters for the live demo. If you change the initial experience, make sure `make reset-db` restores a clear demo state.
- Prefer simple, obvious UI changes over flashy ones. This project should look like a real product, not a design experiment.

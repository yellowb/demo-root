## 1. Decision Alignment

- [x] 1.1 Confirm the create-without-priority behavior uses the design default of `medium`, or update proposal/design/specs before implementation.
- [x] 1.2 Review the `todo-management` delta spec so implementation scope stays limited to priority and filtering.

## 2. Backend Data And Persistence

- [x] 2.1 Add priority value definitions to the Todo domain model and create/update/list input types.
- [x] 2.2 Update SQLite schema handling to add a non-null `priority` column with `medium` as the existing-row default.
- [x] 2.3 Update seed data to include a visible mix of `low`, `medium`, and `high` priorities.

## 3. Backend API Behavior

- [x] 3.1 Validate priority values in create and update flows, default missing create priority to `medium`, and keep empty updates invalid.
- [x] 3.2 Extend repository create, update, scan, and get-by-id behavior to persist and return priority.
- [x] 3.3 Extend repository list behavior to support optional completion and priority filters while preserving newest-first order.
- [x] 3.4 Parse `completed=true|false` and `priority=low|medium|high` query parameters in `GET /api/todos` and reject unsupported values.

## 4. Frontend Workflow

- [x] 4.1 Extend Todo TypeScript types and API client helpers for priority and list filters.
- [x] 4.2 Add priority selection to the create composer with `medium` as the default draft value.
- [x] 4.3 Add priority display and inline priority editing to todo cards.
- [x] 4.4 Add completion and priority filter controls above the list and fetch filtered results from the API.
- [x] 4.5 Ensure create, edit, toggle, and delete updates keep the visible list consistent with active filters.

## 5. Tests And Verification

- [x] 5.1 Update backend repository, service, and HTTP route tests for priority defaults, validation, filtering, and existing-row backfill.
- [x] 5.2 Update frontend tests for priority display, create/edit priority, filter controls, and filtered API requests.
- [x] 5.3 Run OpenSpec validation for the change and fix any spec formatting or requirement issues.
- [x] 5.4 Run `make test` and fix any backend lint, backend test, frontend typecheck, or frontend build failure.
- [x] 5.5 Run `make reset-db` or otherwise verify reset seed data still supports the live demo initial state.

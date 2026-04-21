## Context

This repository is the Todo List baseline for an Agent Harness live demo. The current app supports local Todo CRUD with SQLite persistence, a small Go API, and a single-page React workflow. The baseline intentionally excludes priority and filtering, so this change should make a cross-layer feature visible without expanding into search, tags, due dates, auth, or multi-user behavior.

The existing Todo model contains `id`, `title`, `notes`, `completed`, `created_at`, and `updated_at`. Listing is newest-first, create defaults `completed` to false, update is partial, and the frontend keeps create/edit interactions inline.

## Goals / Non-Goals

**Goals:**

- Add priority as a first-class Todo field with exactly `low`, `medium`, and `high` values.
- Return priority in every Todo API payload and persist it in SQLite.
- Support priority during create and inline edit flows.
- Filter the list by completion status and priority while preserving newest-first ordering.
- Update seed data and tests so the live demo result is visually obvious and verifiable.

**Non-Goals:**

- No search, tags, due dates, login, multi-user collaboration, or external services.
- No database replacement, Docker setup, or new cloud dependency.
- No multi-page frontend flow or complex routing.
- No change to the existing health, delete, title validation, or notes behavior except where payload shape includes priority.

## Decisions

1. **Represent priority as a small string enum.**
   - Decision: use `low`, `medium`, and `high` as the only valid values across Go types, JSON payloads, TypeScript types, and SQLite storage.
   - Rationale: the values are user-facing, easy to read in seed data, and match the requested live demo vocabulary.
   - Alternative considered: store numeric ranks. That would make sorting easier later, but this change does not require priority sorting and numeric values are less clear in API payloads.

2. **Default missing create priority to `medium`.**
   - Decision: the stored Todo always has priority, but create requests that omit priority use `medium`; requests that provide an unsupported priority are rejected.
   - Rationale: this keeps existing create clients from breaking while still making priority explicit in the returned model and UI. SQLite can also use `DEFAULT 'medium'` when migrating existing rows.
   - Alternative considered: require priority on every create API call. That is stricter, but it would be a breaking API change for existing clients and tests that currently create todos with only title and notes.

3. **Keep update partial and add priority as another optional editable field.**
   - Decision: `PATCH /api/todos/:id` accepts optional `priority` alongside `title`, `notes`, and `completed`; empty updates remain invalid.
   - Rationale: this preserves the current partial-update contract and lets inline editing update title, notes, and priority in one request.
   - Alternative considered: introduce a dedicated priority endpoint. That would add API surface without improving the demo.

4. **Use query parameters for filtering on the existing list endpoint.**
   - Decision: `GET /api/todos` accepts optional completion and priority filters, for example `completed=true` and `priority=high`; missing filters return all todos.
   - Rationale: filtering is list behavior, so it belongs on the existing list endpoint and keeps the API compact.
   - Alternative considered: fetch all todos and filter only in the frontend. That would not demonstrate backend query, handler, and persistence changes.

5. **Preserve newest-first ordering after filters.**
   - Decision: filters narrow the result set, but ordering remains `created_at DESC, id DESC`.
   - Rationale: this keeps the baseline list behavior stable and limits the change to filtering rather than sorting.

6. **Keep the UI inline and compact.**
   - Decision: add a priority select to the composer and inline edit form, show priority as a visible pill on each card, and add small filter controls above the list.
   - Rationale: this follows the existing single-page Todo workflow and keeps the demo focused on a real product-like enhancement.

## Risks / Trade-offs

- **Schema migration on existing local databases** -> Add the column with a default value and preserve existing rows; reset-db can still rebuild visually useful seed data.
- **Filter state and local optimistic updates can drift from server results** -> After create/update/delete, update local state consistently and ensure the visible list respects the active filters.
- **Priority validation can diverge between backend and frontend** -> Define the same value set in Go and TypeScript tests, and reject unsupported backend input regardless of frontend controls.
- **More UI controls can make the baseline feel busy** -> Use compact selects/segmented controls and keep the existing single-page layout.

## Migration Plan

- Add a non-null `priority` column to `todos` with default `medium` for existing SQLite rows.
- Update seed rows to include a mix of `low`, `medium`, and `high` priorities.
- Keep rollback simple for the local demo: reset the database with `make reset-db` if the schema needs to be rebuilt during rehearsal.

## Open Questions

- The requested feature does not specify API behavior when create omits priority. This design chooses `medium` as the default to avoid a breaking API change; confirm before implementation if strict API input is preferred.

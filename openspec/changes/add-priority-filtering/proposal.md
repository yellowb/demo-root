## Why

The live demo needs a visible cross-layer feature that shows how Agent Harness uses repository context, specs, tools, and validation together. Adding priority and filtering extends the existing Todo workflow without introducing unrelated product scope.

## What Changes

- Add a `priority` field to todos with allowed values `low`, `medium`, and `high`.
- Allow users to set priority when creating a todo.
- Allow users to edit priority inline with the existing todo edit flow.
- Display each todo's priority clearly in the list.
- Add list filtering by completion status and priority.
- Update API payloads, validation, SQLite schema, repository behavior, seed data, frontend API types, UI state, and tests.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `todo-management`: Add priority to the Todo data model, create/update behavior, list filtering, frontend workflow, SQLite persistence, and demo seed data.

## Impact

- Backend: todo model, service validation, repository queries, HTTP handlers, SQLite schema, seed data, and tests.
- Frontend: API client types, create/edit UI, todo list rendering, filter controls, local state updates, and type/build checks.
- Specs: update `todo-management` requirements so priority and filtering become part of the acceptance contract.

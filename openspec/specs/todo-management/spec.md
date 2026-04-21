# Todo Management Specification

## Purpose

This capability defines the baseline Todo List behavior for the Agent Harness live demo. The application provides local Todo CRUD, SQLite persistence, demo seed data, HTTP API endpoints, and a single-page frontend workflow.

## Requirements

### Requirement: Todo data model
The system SHALL represent each todo with `id`, `title`, `notes`, `completed`, `created_at`, and `updated_at` fields.

#### Scenario: Return todo payload
- **WHEN** a todo is returned from the API
- **THEN** the payload includes numeric `id`, string `title`, string `notes`, boolean `completed`, timestamp `created_at`, and timestamp `updated_at`

#### Scenario: Store optional notes
- **WHEN** a todo is created or updated with empty notes
- **THEN** the system SHALL store notes as an empty string

### Requirement: Todo creation
The system SHALL allow creating todos with a required title and optional notes.

#### Scenario: Create valid todo
- **WHEN** a client submits a non-blank title and notes
- **THEN** the system trims title and notes, persists the todo, returns the created todo, sets `completed` to false, and sets creation and update timestamps

#### Scenario: Reject blank title
- **WHEN** a client submits a missing or blank title
- **THEN** the system MUST reject the request with a validation error

### Requirement: Todo listing
The system SHALL list persisted todos in newest-first order.

#### Scenario: List todos
- **WHEN** a client requests the todo list
- **THEN** the system returns all persisted todos ordered by `created_at` descending and `id` descending

#### Scenario: List empty store
- **WHEN** no todos exist in the database
- **THEN** the system SHALL return an empty list

### Requirement: Todo update
The system SHALL support partial updates to title, notes, and completion state.

#### Scenario: Update editable fields
- **WHEN** a client updates one or more of title, notes, or completed
- **THEN** the system persists only the provided changes, returns the updated todo, and refreshes `updated_at`

#### Scenario: Reject empty update
- **WHEN** a client submits an update without title, notes, or completed
- **THEN** the system MUST reject the request with a validation error

#### Scenario: Reject blank updated title
- **WHEN** a client updates title to a blank value
- **THEN** the system MUST reject the request with a validation error

### Requirement: Todo deletion
The system SHALL delete todos by id.

#### Scenario: Delete existing todo
- **WHEN** a client deletes an existing todo id
- **THEN** the system removes the todo and returns a successful no-content response

#### Scenario: Delete missing todo
- **WHEN** a client deletes a todo id that does not exist
- **THEN** the system MUST return a not-found error

### Requirement: HTTP API
The system SHALL expose health and Todo CRUD endpoints under `/api`.

#### Scenario: Health check
- **WHEN** a client requests `GET /api/health`
- **THEN** the system returns a successful response with `ok` set to true

#### Scenario: Todo endpoints
- **WHEN** a client uses `GET /api/todos`, `POST /api/todos`, `PATCH /api/todos/:id`, or `DELETE /api/todos/:id`
- **THEN** the system routes the request to the corresponding list, create, update, or delete behavior

#### Scenario: Invalid todo id
- **WHEN** a client sends a non-positive or non-numeric `:id`
- **THEN** the system MUST reject the request with an invalid id error

### Requirement: Frontend todo workflow
The frontend SHALL provide a single-page Todo workflow for loading, creating, editing, toggling, and deleting todos.

#### Scenario: Load current list
- **WHEN** the Todo page opens
- **THEN** the frontend fetches todos from `/api/todos` and displays loading, error, empty, or list states as appropriate

#### Scenario: Create from UI
- **WHEN** a user enters a title and optional notes in the composer and submits
- **THEN** the frontend validates that title is non-blank, calls the create API, prepends the created todo to the list, and clears the draft

#### Scenario: Edit inline
- **WHEN** a user edits a todo inline and saves
- **THEN** the frontend validates that title is non-blank, calls the update API, and replaces the item with the returned todo

#### Scenario: Toggle and delete
- **WHEN** a user toggles completion or deletes a todo
- **THEN** the frontend calls the corresponding API and updates the local list from the result

### Requirement: SQLite persistence and demo seed data
The system SHALL persist todos in a local SQLite database and seed demo data into an empty database.

#### Scenario: Bootstrap database
- **WHEN** the backend starts with a configured database path
- **THEN** the system creates the database directory if needed, opens SQLite, applies the todo schema, and enables the expected database settings

#### Scenario: Seed empty database
- **WHEN** the database contains no todos during bootstrap
- **THEN** the system inserts the baseline live-demo seed todos with deterministic created timestamps and completion states

#### Scenario: Preserve existing database
- **WHEN** the database already contains one or more todos during bootstrap
- **THEN** the system MUST leave existing todos unchanged and skip demo seeding

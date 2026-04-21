## MODIFIED Requirements

### Requirement: Todo data model
The system SHALL represent each todo with `id`, `title`, `notes`, `completed`, `priority`, `created_at`, and `updated_at` fields. The `priority` field MUST be one of `low`, `medium`, or `high`.

#### Scenario: Return todo payload
- **WHEN** a todo is returned from the API
- **THEN** the payload includes numeric `id`, string `title`, string `notes`, boolean `completed`, string `priority`, timestamp `created_at`, and timestamp `updated_at`

#### Scenario: Store optional notes
- **WHEN** a todo is created or updated with empty notes
- **THEN** the system SHALL store notes as an empty string

#### Scenario: Return supported priority value
- **WHEN** a todo is returned from the API
- **THEN** `priority` MUST be `low`, `medium`, or `high`

### Requirement: Todo creation
The system SHALL allow creating todos with a required title, optional notes, and optional priority. When priority is provided, it MUST be one of `low`, `medium`, or `high`.

#### Scenario: Create valid todo
- **WHEN** a client submits a non-blank title, notes, and a supported priority
- **THEN** the system trims title and notes, persists the todo, returns the created todo with the submitted priority, sets `completed` to false, and sets creation and update timestamps

#### Scenario: Default missing priority
- **WHEN** a client submits a valid create request without priority
- **THEN** the system SHALL persist and return the todo with priority set to `medium`

#### Scenario: Reject blank title
- **WHEN** a client submits a missing or blank title
- **THEN** the system MUST reject the request with a validation error

#### Scenario: Reject invalid create priority
- **WHEN** a client submits priority outside `low`, `medium`, or `high`
- **THEN** the system MUST reject the request with a validation error

### Requirement: Todo listing
The system SHALL list persisted todos in newest-first order and support optional filtering by completion status and priority.

#### Scenario: List todos
- **WHEN** a client requests the todo list without filters
- **THEN** the system returns all persisted todos ordered by `created_at` descending and `id` descending

#### Scenario: List empty store
- **WHEN** no todos exist in the database
- **THEN** the system SHALL return an empty list

#### Scenario: Filter by completion status
- **WHEN** a client requests todos filtered by completed or active status
- **THEN** the system returns only todos whose `completed` value matches the requested status, ordered by `created_at` descending and `id` descending

#### Scenario: Filter by priority
- **WHEN** a client requests todos filtered by one priority
- **THEN** the system returns only todos whose `priority` matches the requested priority, ordered by `created_at` descending and `id` descending

#### Scenario: Combine completion and priority filters
- **WHEN** a client requests todos filtered by both completion status and priority
- **THEN** the system returns only todos matching both filters, ordered by `created_at` descending and `id` descending

#### Scenario: Reject invalid list filters
- **WHEN** a client requests todos with an unsupported completion or priority filter value
- **THEN** the system MUST reject the request with a validation error

### Requirement: Todo update
The system SHALL support partial updates to title, notes, completion state, and priority.

#### Scenario: Update editable fields
- **WHEN** a client updates one or more of title, notes, completed, or priority
- **THEN** the system persists only the provided changes, returns the updated todo, and refreshes `updated_at`

#### Scenario: Reject empty update
- **WHEN** a client submits an update without title, notes, completed, or priority
- **THEN** the system MUST reject the request with a validation error

#### Scenario: Reject blank updated title
- **WHEN** a client updates title to a blank value
- **THEN** the system MUST reject the request with a validation error

#### Scenario: Reject invalid updated priority
- **WHEN** a client updates priority outside `low`, `medium`, or `high`
- **THEN** the system MUST reject the request with a validation error

### Requirement: HTTP API
The system SHALL expose health and Todo CRUD endpoints under `/api`. The todo list endpoint SHALL accept optional query parameters for completion status and priority filtering.

#### Scenario: Health check
- **WHEN** a client requests `GET /api/health`
- **THEN** the system returns a successful response with `ok` set to true

#### Scenario: Todo endpoints
- **WHEN** a client uses `GET /api/todos`, `POST /api/todos`, `PATCH /api/todos/:id`, or `DELETE /api/todos/:id`
- **THEN** the system routes the request to the corresponding list, create, update, or delete behavior

#### Scenario: Filter todo endpoint
- **WHEN** a client uses `GET /api/todos` with optional `completed=true|false` and/or `priority=low|medium|high` query parameters
- **THEN** the system routes the request to the filtered list behavior

#### Scenario: Invalid todo id
- **WHEN** a client sends a non-positive or non-numeric `:id`
- **THEN** the system MUST reject the request with an invalid id error

#### Scenario: Invalid filter parameter
- **WHEN** a client sends an unsupported `completed` or `priority` query parameter value
- **THEN** the system MUST reject the request with a validation error

### Requirement: Frontend todo workflow
The frontend SHALL provide a single-page Todo workflow for loading, creating, editing, toggling, deleting, and filtering todos.

#### Scenario: Load current list
- **WHEN** the Todo page opens
- **THEN** the frontend fetches todos from `/api/todos` and displays loading, error, empty, or list states as appropriate

#### Scenario: Create from UI
- **WHEN** a user enters a title, optional notes, and priority in the composer and submits
- **THEN** the frontend validates that title is non-blank, calls the create API, prepends the created todo to the list when it matches active filters, and clears the draft

#### Scenario: Edit inline
- **WHEN** a user edits a todo's title, notes, or priority inline and saves
- **THEN** the frontend validates that title is non-blank, calls the update API, and replaces or removes the item from the visible list based on the returned todo and active filters

#### Scenario: Toggle and delete
- **WHEN** a user toggles completion or deletes a todo
- **THEN** the frontend calls the corresponding API and updates the local list from the result

#### Scenario: Filter from UI
- **WHEN** a user chooses completion status or priority filters
- **THEN** the frontend fetches todos with the corresponding query parameters and displays only the matching todos

#### Scenario: Show todo priority
- **WHEN** the frontend displays a todo in the list
- **THEN** the todo's priority is clearly visible with the item

### Requirement: SQLite persistence and demo seed data
The system SHALL persist todos with priority in a local SQLite database and seed demo data into an empty database.

#### Scenario: Bootstrap database
- **WHEN** the backend starts with a configured database path
- **THEN** the system creates the database directory if needed, opens SQLite, applies the todo schema including priority, and enables the expected database settings

#### Scenario: Seed empty database
- **WHEN** the database contains no todos during bootstrap
- **THEN** the system inserts the baseline live-demo seed todos with deterministic created timestamps, completion states, and varied priorities

#### Scenario: Preserve existing database
- **WHEN** the database already contains one or more todos during bootstrap
- **THEN** the system MUST leave existing todos unchanged and skip demo seeding

#### Scenario: Backfill existing priority
- **WHEN** an existing local database has todos without a priority column
- **THEN** the system SHALL preserve existing todos and make them readable with priority set to `medium`

package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"agent-harness-demo/backend/internal/store"
	"agent-harness-demo/backend/internal/todos"
)

func TestHealthRoute(t *testing.T) {
	t.Parallel()

	router := newTestRouter(t)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/health", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	var payload map[string]bool
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode health response: %v", err)
	}
	if !payload["ok"] {
		t.Fatalf("expected ok=true, got %#v", payload)
	}
}

func TestTodoRoutesLifecycle(t *testing.T) {
	t.Parallel()

	router := newTestRouter(t)

	created := performJSONRequest(t, router, http.MethodPost, "/api/todos", map[string]any{
		"title":    "Prepare demo environment",
		"notes":    "Confirm the baseline codebase is ready",
		"priority": "high",
	})
	if created.Code != http.StatusCreated {
		t.Fatalf("expected status 201 on create, got %d with body %s", created.Code, created.Body.String())
	}

	var createdTodo todos.Todo
	if err := json.Unmarshal(created.Body.Bytes(), &createdTodo); err != nil {
		t.Fatalf("decode created todo: %v", err)
	}

	listResponse := performRequest(t, router, http.MethodGet, "/api/todos", nil)
	if listResponse.Code != http.StatusOK {
		t.Fatalf("expected status 200 on list, got %d", listResponse.Code)
	}

	var listed []todos.Todo
	if err := json.Unmarshal(listResponse.Body.Bytes(), &listed); err != nil {
		t.Fatalf("decode listed todos: %v", err)
	}
	if len(listed) != 1 {
		t.Fatalf("expected 1 todo after create, got %d", len(listed))
	}

	updated := performJSONRequest(t, router, http.MethodPatch, fmt.Sprintf("/api/todos/%d", createdTodo.ID), map[string]any{
		"title":     "Prepare demo walkthrough",
		"notes":     "Highlight repo inspection and test verification",
		"completed": true,
	})
	if updated.Code != http.StatusOK {
		t.Fatalf("expected status 200 on update, got %d with body %s", updated.Code, updated.Body.String())
	}

	deleteResponse := performRequest(t, router, http.MethodDelete, fmt.Sprintf("/api/todos/%d", createdTodo.ID), nil)
	if deleteResponse.Code != http.StatusNoContent {
		t.Fatalf("expected status 204 on delete, got %d", deleteResponse.Code)
	}

	finalList := performRequest(t, router, http.MethodGet, "/api/todos", nil)
	if finalList.Code != http.StatusOK {
		t.Fatalf("expected status 200 on final list, got %d", finalList.Code)
	}

	var remaining []todos.Todo
	if err := json.Unmarshal(finalList.Body.Bytes(), &remaining); err != nil {
		t.Fatalf("decode final list: %v", err)
	}
	if len(remaining) != 0 {
		t.Fatalf("expected empty list after delete, got %d items", len(remaining))
	}

	if createdTodo.ID == 0 {
		t.Fatalf("expected created todo to have a non-zero id")
	}
	if createdTodo.Priority != todos.PriorityHigh {
		t.Fatalf("expected created priority %q, got %q", todos.PriorityHigh, createdTodo.Priority)
	}
}

func TestCreateTodoRejectsBlankTitle(t *testing.T) {
	t.Parallel()

	router := newTestRouter(t)
	response := performJSONRequest(t, router, http.MethodPost, "/api/todos", map[string]any{
		"title": "   ",
		"notes": "Should fail validation",
	})

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d with body %s", response.Code, response.Body.String())
	}
}

func TestTodoRoutesFilterByCompletionAndPriority(t *testing.T) {
	t.Parallel()

	router := newTestRouter(t)

	activeHigh := performJSONRequest(t, router, http.MethodPost, "/api/todos", map[string]any{
		"title":    "Prepare live demo",
		"notes":    "Show high priority active item",
		"priority": "high",
	})
	if activeHigh.Code != http.StatusCreated {
		t.Fatalf("expected status 201 for active high todo, got %d", activeHigh.Code)
	}

	completedHigh := performJSONRequest(t, router, http.MethodPost, "/api/todos", map[string]any{
		"title":    "Archive baseline notes",
		"notes":    "Completed high item",
		"priority": "high",
	})
	if completedHigh.Code != http.StatusCreated {
		t.Fatalf("expected status 201 for completed high todo, got %d", completedHigh.Code)
	}

	var completedTodo todos.Todo
	if err := json.Unmarshal(completedHigh.Body.Bytes(), &completedTodo); err != nil {
		t.Fatalf("decode completed todo: %v", err)
	}

	updateResponse := performJSONRequest(t, router, http.MethodPatch, fmt.Sprintf("/api/todos/%d", completedTodo.ID), map[string]any{
		"completed": true,
	})
	if updateResponse.Code != http.StatusOK {
		t.Fatalf("expected status 200 on completion update, got %d", updateResponse.Code)
	}

	filterResponse := performRequest(t, router, http.MethodGet, "/api/todos?completed=false&priority=high", nil)
	if filterResponse.Code != http.StatusOK {
		t.Fatalf("expected status 200 on filtered list, got %d with body %s", filterResponse.Code, filterResponse.Body.String())
	}

	var filtered []todos.Todo
	if err := json.Unmarshal(filterResponse.Body.Bytes(), &filtered); err != nil {
		t.Fatalf("decode filtered todos: %v", err)
	}
	if len(filtered) != 1 {
		t.Fatalf("expected 1 filtered todo, got %d", len(filtered))
	}
	if filtered[0].Completed || filtered[0].Priority != todos.PriorityHigh {
		t.Fatalf("expected active high priority todo, got %#v", filtered[0])
	}
}

func TestTodoRoutesRejectInvalidFilters(t *testing.T) {
	t.Parallel()

	router := newTestRouter(t)

	invalidCompleted := performRequest(t, router, http.MethodGet, "/api/todos?completed=yes", nil)
	if invalidCompleted.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for invalid completed filter, got %d", invalidCompleted.Code)
	}

	invalidPriority := performRequest(t, router, http.MethodGet, "/api/todos?priority=urgent", nil)
	if invalidPriority.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for invalid priority filter, got %d", invalidPriority.Code)
	}
}

func newTestRouter(t *testing.T) http.Handler {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "todos.db")
	db, err := store.Open(dbPath)
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}

	t.Cleanup(func() {
		_ = db.Close()
	})

	if err := store.EnsureSchema(context.Background(), db); err != nil {
		t.Fatalf("ensure schema: %v", err)
	}

	repo := todos.NewRepository(db)
	service := todos.NewService(repo)

	return NewRouter(service)
}

func performJSONRequest(t *testing.T, handler http.Handler, method string, path string, payload map[string]any) *httptest.ResponseRecorder {
	t.Helper()

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal request payload: %v", err)
	}

	return performRequest(t, handler, method, path, bytes.NewReader(body))
}

func performRequest(t *testing.T, handler http.Handler, method string, path string, body *bytes.Reader) *httptest.ResponseRecorder {
	t.Helper()

	var requestBody *bytes.Reader
	if body == nil {
		requestBody = bytes.NewReader(nil)
	} else {
		requestBody = body
	}

	request := httptest.NewRequest(method, path, requestBody)
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	return recorder
}

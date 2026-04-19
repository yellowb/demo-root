package todos

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"agent-harness-demo/backend/internal/store"
)

func TestRepositoryCRUDLifecycle(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db := openTestDatabase(t)
	repo := NewRepository(db)

	created, err := repo.Create(ctx, CreateInput{
		Title: "Prepare Agent Harness demo",
		Notes: "Walk through the live coding baseline",
	})
	if err != nil {
		t.Fatalf("create todo: %v", err)
	}

	if created.ID == 0 {
		t.Fatalf("expected created todo to have an id")
	}
	if created.Completed {
		t.Fatalf("expected created todo to default to incomplete")
	}

	initialList, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("list todos after create: %v", err)
	}
	if len(initialList) != 1 {
		t.Fatalf("expected 1 todo after create, got %d", len(initialList))
	}

	newTitle := "Prepare Agent Harness live demo"
	newNotes := "Highlight repo inspection and verification"
	completed := true
	updated, err := repo.Update(ctx, created.ID, UpdateInput{
		Title:     &newTitle,
		Notes:     &newNotes,
		Completed: &completed,
	})
	if err != nil {
		t.Fatalf("update todo: %v", err)
	}

	if updated.Title != newTitle {
		t.Fatalf("expected updated title %q, got %q", newTitle, updated.Title)
	}
	if updated.Notes != newNotes {
		t.Fatalf("expected updated notes %q, got %q", newNotes, updated.Notes)
	}
	if !updated.Completed {
		t.Fatalf("expected updated todo to be completed")
	}

	listAfterUpdate, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("list todos after update: %v", err)
	}
	if len(listAfterUpdate) != 1 {
		t.Fatalf("expected 1 todo after update, got %d", len(listAfterUpdate))
	}
	if listAfterUpdate[0].Title != newTitle {
		t.Fatalf("expected listed todo title %q, got %q", newTitle, listAfterUpdate[0].Title)
	}

	if err := repo.Delete(ctx, created.ID); err != nil {
		t.Fatalf("delete todo: %v", err)
	}

	finalList, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("list todos after delete: %v", err)
	}
	if len(finalList) != 0 {
		t.Fatalf("expected empty list after delete, got %d items", len(finalList))
	}
}

func TestRepositoryListOrdersNewestFirst(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db := openTestDatabase(t)
	repo := NewRepository(db)

	first, err := repo.Create(ctx, CreateInput{
		Title: "Draft baseline README",
		Notes: "Document quick-start steps",
	})
	if err != nil {
		t.Fatalf("create first todo: %v", err)
	}

	second, err := repo.Create(ctx, CreateInput{
		Title: "Review baseline UI",
		Notes: "Check spacing and hierarchy",
	})
	if err != nil {
		t.Fatalf("create second todo: %v", err)
	}

	listed, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("list todos: %v", err)
	}
	if len(listed) != 2 {
		t.Fatalf("expected 2 todos, got %d", len(listed))
	}
	if listed[0].ID != second.ID {
		t.Fatalf("expected newest todo %d first, got %d", second.ID, listed[0].ID)
	}
	if listed[1].ID != first.ID {
		t.Fatalf("expected older todo %d second, got %d", first.ID, listed[1].ID)
	}
}

func openTestDatabase(t *testing.T) *sql.DB {
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

	return db
}

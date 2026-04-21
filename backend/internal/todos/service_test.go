package todos

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"agent-harness-demo/backend/internal/store"
)

func TestServiceDefaultsMissingCreatePriority(t *testing.T) {
	t.Parallel()

	service := newTestService(t)
	created, err := service.Create(context.Background(), CreateInput{
		Title: "Document priority default",
		Notes: "Existing create clients can omit priority",
	})
	if err != nil {
		t.Fatalf("create todo without priority: %v", err)
	}

	if created.Priority != PriorityMedium {
		t.Fatalf("expected default priority %q, got %q", PriorityMedium, created.Priority)
	}
}

func TestServiceRejectsInvalidPriority(t *testing.T) {
	t.Parallel()

	service := newTestService(t)
	if _, err := service.Create(context.Background(), CreateInput{
		Title:    "Invalid priority",
		Priority: Priority("urgent"),
	}); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected validation error for invalid create priority, got %v", err)
	}

	created, err := service.Create(context.Background(), CreateInput{
		Title:    "Valid priority",
		Priority: PriorityHigh,
	})
	if err != nil {
		t.Fatalf("create todo: %v", err)
	}

	invalidPriority := Priority("urgent")
	if _, err := service.Update(context.Background(), created.ID, UpdateInput{
		Priority: &invalidPriority,
	}); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected validation error for invalid update priority, got %v", err)
	}
}

func newTestService(t *testing.T) *Service {
	t.Helper()

	db, err := store.Open(filepath.Join(t.TempDir(), "todos.db"))
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	if err := store.EnsureSchema(context.Background(), db); err != nil {
		t.Fatalf("ensure schema: %v", err)
	}

	return NewService(NewRepository(db))
}

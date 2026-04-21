package store

import (
	"context"
	"path/filepath"
	"testing"
)

func TestEnsureSchemaBackfillsPriorityForExistingTodos(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db, err := Open(filepath.Join(t.TempDir(), "todos.db"))
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	if _, err := db.ExecContext(ctx, `
		CREATE TABLE todos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			notes TEXT NOT NULL DEFAULT '',
			completed INTEGER NOT NULL DEFAULT 0,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)
	`); err != nil {
		t.Fatalf("create legacy schema: %v", err)
	}

	if _, err := db.ExecContext(ctx, `
		INSERT INTO todos (title, notes, completed, created_at, updated_at)
		VALUES ('Legacy todo', '', 0, '2026-04-18T09:00:00Z', '2026-04-18T09:00:00Z')
	`); err != nil {
		t.Fatalf("insert legacy todo: %v", err)
	}

	if err := EnsureSchema(ctx, db); err != nil {
		t.Fatalf("ensure schema: %v", err)
	}

	var priority string
	if err := db.QueryRowContext(ctx, `SELECT priority FROM todos WHERE title = 'Legacy todo'`).Scan(&priority); err != nil {
		t.Fatalf("read backfilled priority: %v", err)
	}

	if priority != "medium" {
		t.Fatalf("expected backfilled priority medium, got %q", priority)
	}
}

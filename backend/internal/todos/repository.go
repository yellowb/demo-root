package todos

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(ctx context.Context) ([]Todo, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, notes, completed, created_at, updated_at
		FROM todos
		ORDER BY created_at DESC, id DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("list todos: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var items []Todo
	for rows.Next() {
		item, err := scanTodo(rows)
		if err != nil {
			return nil, fmt.Errorf("scan listed todo: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate listed todos: %w", err)
	}

	return items, nil
}

func (r *Repository) Create(ctx context.Context, input CreateInput) (Todo, error) {
	now := time.Now().UTC().Format(time.RFC3339Nano)

	result, err := r.db.ExecContext(ctx, `
		INSERT INTO todos (title, notes, completed, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`, input.Title, input.Notes, 0, now, now)
	if err != nil {
		return Todo{}, fmt.Errorf("insert todo: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Todo{}, fmt.Errorf("read inserted todo id: %w", err)
	}

	return r.getByID(ctx, id)
}

func (r *Repository) Update(ctx context.Context, id int64, input UpdateInput) (Todo, error) {
	existing, err := r.getByID(ctx, id)
	if err != nil {
		return Todo{}, err
	}

	if input.Title != nil {
		existing.Title = *input.Title
	}
	if input.Notes != nil {
		existing.Notes = *input.Notes
	}
	if input.Completed != nil {
		existing.Completed = *input.Completed
	}

	updatedAt := time.Now().UTC().Format(time.RFC3339Nano)
	completed := 0
	if existing.Completed {
		completed = 1
	}

	result, err := r.db.ExecContext(ctx, `
		UPDATE todos
		SET title = ?, notes = ?, completed = ?, updated_at = ?
		WHERE id = ?
	`, existing.Title, existing.Notes, completed, updatedAt, id)
	if err != nil {
		return Todo{}, fmt.Errorf("update todo %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return Todo{}, fmt.Errorf("read affected rows for todo %d: %w", id, err)
	}
	if rowsAffected == 0 {
		return Todo{}, ErrNotFound
	}

	return r.getByID(ctx, id)
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM todos WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete todo %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read affected rows for deleted todo %d: %w", id, err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *Repository) getByID(ctx context.Context, id int64) (Todo, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, title, notes, completed, created_at, updated_at
		FROM todos
		WHERE id = ?
	`, id)

	item, err := scanTodo(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return Todo{}, ErrNotFound
		}
		return Todo{}, fmt.Errorf("get todo %d: %w", id, err)
	}

	return item, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanTodo(source scanner) (Todo, error) {
	var item Todo
	var completed int
	var createdAt string
	var updatedAt string

	if err := source.Scan(
		&item.ID,
		&item.Title,
		&item.Notes,
		&completed,
		&createdAt,
		&updatedAt,
	); err != nil {
		return Todo{}, err
	}

	parsedCreatedAt, err := time.Parse(time.RFC3339Nano, createdAt)
	if err != nil {
		return Todo{}, fmt.Errorf("parse created_at: %w", err)
	}

	parsedUpdatedAt, err := time.Parse(time.RFC3339Nano, updatedAt)
	if err != nil {
		return Todo{}, fmt.Errorf("parse updated_at: %w", err)
	}

	item.Completed = completed == 1
	item.CreatedAt = parsedCreatedAt
	item.UpdatedAt = parsedUpdatedAt

	return item, nil
}

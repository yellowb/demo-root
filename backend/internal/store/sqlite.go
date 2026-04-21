package store

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

const schemaSQL = `
CREATE TABLE IF NOT EXISTS todos (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL,
	notes TEXT NOT NULL DEFAULT '',
	completed INTEGER NOT NULL DEFAULT 0,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_todos_created_at ON todos(created_at DESC, id DESC);
`

type seedTodo struct {
	Title     string
	Notes     string
	Completed bool
	CreatedAt time.Time
}

func Open(dbPath string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return nil, fmt.Errorf("create db directory: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite db: %w", err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("enable foreign keys: %w", err)
	}

	if _, err := db.Exec(`PRAGMA journal_mode = WAL;`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("enable wal mode: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping sqlite db: %w", err)
	}

	return db, nil
}

func EnsureSchema(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, schemaSQL); err != nil {
		return fmt.Errorf("apply schema: %w", err)
	}

	return nil
}

func Bootstrap(ctx context.Context, dbPath string) (*sql.DB, error) {
	db, err := Open(dbPath)
	if err != nil {
		return nil, err
	}

	if err := EnsureSchema(ctx, db); err != nil {
		_ = db.Close()
		return nil, err
	}

	if err := SeedDemoData(ctx, db); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

func SeedDemoData(ctx context.Context, db *sql.DB) error {
	var count int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM todos`).Scan(&count); err != nil {
		return fmt.Errorf("count todos: %w", err)
	}
	if count > 0 {
		return nil
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin seed transaction: %w", err)
	}

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO todos (title, notes, completed, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("prepare seed statement: %w", err)
	}
	defer func() { _ = stmt.Close() }()

	for _, todo := range demoSeedTodos() {
		timestamp := todo.CreatedAt.UTC().Format(time.RFC3339Nano)
		completed := 0
		if todo.Completed {
			completed = 1
		}

		if _, err := stmt.ExecContext(ctx, todo.Title, todo.Notes, completed, timestamp, timestamp); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("insert seed todo %q: %w", todo.Title, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit seed transaction: %w", err)
	}

	return nil
}

func demoSeedTodos() []seedTodo {
	base := time.Date(2026, 4, 18, 9, 0, 0, 0, time.UTC)

	return []seedTodo{
		{
			Title:     "准备 Agent Harness 分享提纲",
			Notes:     "补齐背景、问题定义和 live demo 目标",
			Completed: false,
			CreatedAt: base.Add(0 * time.Minute),
		},
		{
			Title:     "整理 Claude Code 案例图",
			Notes:     "挑出能体现 repo inspection 的截图",
			Completed: true,
			CreatedAt: base.Add(10 * time.Minute),
		},
		{
			Title:     "确认 live demo 切换页",
			Notes:     "确保从 PPT 到终端再到浏览器的切换顺滑",
			Completed: false,
			CreatedAt: base.Add(20 * time.Minute),
		},
		{
			Title:     "补充后端测试",
			Notes:     "覆盖 Todo CRUD 和关键 HTTP handler 路径",
			Completed: true,
			CreatedAt: base.Add(30 * time.Minute),
		},
		{
			Title:     "回看前端卡片样式",
			Notes:     "确认列表卡片层次清楚，适合放进 PPT",
			Completed: false,
			CreatedAt: base.Add(40 * time.Minute),
		},
		{
			Title:     "核对种子数据文案",
			Notes:     "保持真实工作语境，但先不要出现 priority 字段",
			Completed: false,
			CreatedAt: base.Add(50 * time.Minute),
		},
		{
			Title:     "检查验证命令",
			Notes:     "保证 make test 能串起后端测试和前端静态验证",
			Completed: true,
			CreatedAt: base.Add(60 * time.Minute),
		},
		{
			Title:     "预演浏览器展示顺序",
			Notes:     "先看基线列表，再演示后续加入 priority 和 filtering",
			Completed: false,
			CreatedAt: base.Add(70 * time.Minute),
		},
	}
}

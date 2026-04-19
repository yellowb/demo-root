package todos

import (
	"errors"
	"time"
)

var (
	ErrNotFound   = errors.New("todo not found")
	ErrValidation = errors.New("validation failed")
)

type Todo struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Notes     string    `json:"notes"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateInput struct {
	Title string `json:"title"`
	Notes string `json:"notes"`
}

type UpdateInput struct {
	Title     *string `json:"title,omitempty"`
	Notes     *string `json:"notes,omitempty"`
	Completed *bool   `json:"completed,omitempty"`
}

package todos

import (
	"errors"
	"strings"
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
	Priority  Priority  `json:"priority"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

type CreateInput struct {
	Title    string   `json:"title"`
	Notes    string   `json:"notes"`
	Priority Priority `json:"priority"`
}

type UpdateInput struct {
	Title     *string   `json:"title,omitempty"`
	Notes     *string   `json:"notes,omitempty"`
	Completed *bool     `json:"completed,omitempty"`
	Priority  *Priority `json:"priority,omitempty"`
}

type ListFilter struct {
	Completed *bool
	Priority  *Priority
}

func ParsePriority(raw string) (Priority, bool) {
	priority := Priority(strings.TrimSpace(raw))
	switch priority {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return priority, true
	default:
		return "", false
	}
}

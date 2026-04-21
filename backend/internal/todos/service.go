package todos

import (
	"context"
	"fmt"
	"strings"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, filter ListFilter) ([]Todo, error) {
	if filter.Priority != nil {
		priority, ok := ParsePriority(string(*filter.Priority))
		if !ok {
			return nil, fmt.Errorf("%w: priority must be low, medium, or high", ErrValidation)
		}
		filter.Priority = &priority
	}

	return s.repo.List(ctx, filter)
}

func (s *Service) Create(ctx context.Context, input CreateInput) (Todo, error) {
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return Todo{}, fmt.Errorf("%w: title is required", ErrValidation)
	}

	priority := PriorityMedium
	if input.Priority != "" {
		parsed, ok := ParsePriority(string(input.Priority))
		if !ok {
			return Todo{}, fmt.Errorf("%w: priority must be low, medium, or high", ErrValidation)
		}
		priority = parsed
	}

	return s.repo.Create(ctx, CreateInput{
		Title:    title,
		Notes:    strings.TrimSpace(input.Notes),
		Priority: priority,
	})
}

func (s *Service) Update(ctx context.Context, id int64, input UpdateInput) (Todo, error) {
	if input.Title == nil && input.Notes == nil && input.Completed == nil && input.Priority == nil {
		return Todo{}, fmt.Errorf("%w: at least one field must be provided", ErrValidation)
	}

	normalized := UpdateInput{
		Title:     input.Title,
		Notes:     input.Notes,
		Completed: input.Completed,
		Priority:  input.Priority,
	}

	if input.Title != nil {
		title := strings.TrimSpace(*input.Title)
		if title == "" {
			return Todo{}, fmt.Errorf("%w: title is required", ErrValidation)
		}
		normalized.Title = &title
	}

	if input.Notes != nil {
		notes := strings.TrimSpace(*input.Notes)
		normalized.Notes = &notes
	}

	if input.Priority != nil {
		priority, ok := ParsePriority(string(*input.Priority))
		if !ok {
			return Todo{}, fmt.Errorf("%w: priority must be low, medium, or high", ErrValidation)
		}
		normalized.Priority = &priority
	}

	return s.repo.Update(ctx, id, normalized)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

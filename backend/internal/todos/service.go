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

func (s *Service) List(ctx context.Context) ([]Todo, error) {
	return s.repo.List(ctx)
}

func (s *Service) Create(ctx context.Context, input CreateInput) (Todo, error) {
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return Todo{}, fmt.Errorf("%w: title is required", ErrValidation)
	}

	return s.repo.Create(ctx, CreateInput{
		Title: title,
		Notes: strings.TrimSpace(input.Notes),
	})
}

func (s *Service) Update(ctx context.Context, id int64, input UpdateInput) (Todo, error) {
	if input.Title == nil && input.Notes == nil && input.Completed == nil {
		return Todo{}, fmt.Errorf("%w: at least one field must be provided", ErrValidation)
	}

	normalized := UpdateInput{
		Title:     input.Title,
		Notes:     input.Notes,
		Completed: input.Completed,
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

	return s.repo.Update(ctx, id, normalized)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

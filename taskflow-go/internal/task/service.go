package task

import (
	"context"
	"strings"
	"time"

	"taskflow/internal/domain"
)

// CreateTaskInput defines the business payload required to create a new task.
type CreateTaskInput struct {
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Priority    *domain.TaskPriority `json:"priority,omitempty"`
	DueDate     *time.Time           `json:"dueDate,omitempty"`
}

// UpdateTaskInput defines the business payload required to update an existing task.
type UpdateTaskInput struct {
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Status      domain.TaskStatus   `json:"status"`
	Priority    domain.TaskPriority `json:"priority"`
	DueDate     *time.Time          `json:"dueDate,omitempty"`
}

// Service encapsulates the business logic and orchestrates persistence for TaskFlow tasks.
//
// Spring Boot Comparison:
// In Spring Boot, this would be an `@Service` annotated class, often combined with `@Transactional`
// for database transaction management.
// In Go, the service is a plain struct that explicitly receives its dependencies (the repository)
// via constructor injection (`NewService`), completely decoupled from framework annotations or HTTP.
type Service struct {
	repo Repository
}

// NewService creates a new Service instance with explicit dependency injection.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Create validates input, assigns defaults, generates a UUID, and saves the task.
func (s *Service) Create(ctx context.Context, input CreateTaskInput) (domain.Task, error) {
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return domain.Task{}, domain.ErrEmptyTitle
	}

	// Priority defaults to MEDIUM if not specified
	priority := domain.PriorityMedium
	if input.Priority != nil {
		if !input.Priority.IsValid() {
			return domain.Task{}, domain.ErrInvalidPriority
		}
		priority = *input.Priority
	}

	id, err := domain.NewUUID()
	if err != nil {
		return domain.Task{}, err
	}

	now := time.Now().UTC()
	task := domain.Task{
		ID:          id,
		Title:       title,
		Description: strings.TrimSpace(input.Description),
		Status:      domain.StatusPending,
		Priority:    priority,
		CreatedAt:   now,
		UpdatedAt:   now,
		DueDate:     input.DueDate,
	}

	if err := s.repo.Save(ctx, task); err != nil {
		return domain.Task{}, err
	}

	return task, nil
}

// GetByID retrieves a single task by ID.
func (s *Service) GetByID(ctx context.Context, id string) (domain.Task, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return domain.Task{}, domain.ErrInvalidTaskID
	}

	return s.repo.FindByID(ctx, id)
}

// List retrieves all tasks matching the specified filter.
func (s *Service) List(ctx context.Context, filter domain.TaskFilter) ([]domain.Task, error) {
	if filter.Status != nil && !filter.Status.IsValid() {
		return nil, domain.ErrInvalidStatus
	}
	if filter.Priority != nil && !filter.Priority.IsValid() {
		return nil, domain.ErrInvalidPriority
	}

	return s.repo.FindAll(ctx, filter)
}

// Update replaces an existing task's editable fields while maintaining created timestamp.
func (s *Service) Update(ctx context.Context, id string, input UpdateTaskInput) (domain.Task, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return domain.Task{}, domain.ErrInvalidTaskID
	}

	title := strings.TrimSpace(input.Title)
	if title == "" {
		return domain.Task{}, domain.ErrEmptyTitle
	}

	if !input.Status.IsValid() {
		return domain.Task{}, domain.ErrInvalidStatus
	}

	if !input.Priority.IsValid() {
		return domain.Task{}, domain.ErrInvalidPriority
	}

	// Ensure task exists
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return domain.Task{}, err
	}

	existing.Title = title
	existing.Description = strings.TrimSpace(input.Description)
	existing.Status = input.Status
	existing.Priority = input.Priority
	existing.DueDate = input.DueDate
	existing.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, existing); err != nil {
		return domain.Task{}, err
	}

	return existing, nil
}

// Complete marks a task as COMPLETED.
func (s *Service) Complete(ctx context.Context, id string) (domain.Task, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return domain.Task{}, domain.ErrInvalidTaskID
	}

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return domain.Task{}, err
	}

	existing.Status = domain.StatusCompleted
	existing.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, existing); err != nil {
		return domain.Task{}, err
	}

	return existing, nil
}

// Delete removes a task by ID.
func (s *Service) Delete(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return domain.ErrInvalidTaskID
	}

	return s.repo.Delete(ctx, id)
}


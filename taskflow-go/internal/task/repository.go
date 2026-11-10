package task

import (
	"context"
	"taskflow/internal/domain"
)

// Repository defines the persistence contract for TaskFlow tasks.
//
// Spring Boot Comparison:
// In Spring Data JPA, you would define an interface extending JpaRepository<Task, String>
// where Spring automatically generates bytecode queries via proxy beans at runtime.
// In Go, interfaces are satisfied implicitly without annotations, inheritance, or runtime magic.
type Repository interface {
	// Save persists a new task in storage.
	Save(ctx context.Context, task domain.Task) error

	// FindByID retrieves a task by its unique ID. Returns domain.ErrTaskNotFound if missing.
	FindByID(ctx context.Context, id string) (domain.Task, error)

	// FindAll retrieves tasks matching the optional filter criteria.
	FindAll(ctx context.Context, filter domain.TaskFilter) ([]domain.Task, error)

	// Update updates an existing task. Returns domain.ErrTaskNotFound if missing.
	Update(ctx context.Context, task domain.Task) error

	// Delete removes a task by its unique ID. Returns domain.ErrTaskNotFound if missing.
	Delete(ctx context.Context, id string) error
}

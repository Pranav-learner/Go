package task

import (
	"context"
	"sort"
	"sync"
	"taskflow/internal/domain"
)

// MemoryRepository provides a concurrency-safe in-memory implementation of task.Repository.
//
// Concurrency & Synchronization Architecture:
// In Go, native maps (`map[K]V`) are NOT safe for concurrent use. If one goroutine writes to a map
// while another is reading or writing, the Go runtime crashes immediately with an unrecoverable panic:
//
//	"fatal error: concurrent map read and map write"
//
// Java / Spring Boot Comparison:
// In Java, developers often use `java.util.concurrent.ConcurrentHashMap` or `@Transactional` with an RDBMS.
// In Go, `sync.RWMutex` provides explicit reader-writer synchronization:
//   - Multiple concurrent readers can execute simultaneously via `mu.RLock()` / `mu.RUnlock()`.
//   - Writes require exclusive lock access via `mu.Lock()` / `mu.Unlock()`.
type MemoryRepository struct {
	mu    sync.RWMutex
	tasks map[string]domain.Task
}

// NewMemoryRepository constructs and initializes a new in-memory repository instance.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		tasks: make(map[string]domain.Task),
	}
}

// Save stores a new task.
func (r *MemoryRepository) Save(ctx context.Context, task domain.Task) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.tasks[task.ID] = task
	return nil
}

// FindByID retrieves a task by ID using a read-lock.
func (r *MemoryRepository) FindByID(ctx context.Context, id string) (domain.Task, error) {
	if err := ctx.Err(); err != nil {
		return domain.Task{}, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	t, exists := r.tasks[id]
	if !exists {
		return domain.Task{}, domain.ErrTaskNotFound
	}

	return t, nil
}

// FindAll retrieves tasks matching the filter criteria using a read-lock.
// Results are deterministically sorted by CreatedAt ascending.
func (r *MemoryRepository) FindAll(ctx context.Context, filter domain.TaskFilter) ([]domain.Task, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	results := make([]domain.Task, 0, len(r.tasks))
	for _, t := range r.tasks {
		// Filter by status if requested
		if filter.Status != nil && t.Status != *filter.Status {
			continue
		}
		// Filter by priority if requested
		if filter.Priority != nil && t.Priority != *filter.Priority {
			continue
		}
		results = append(results, t)
	}

	// Sort deterministically by CreatedAt, then by ID
	sort.Slice(results, func(i, j int) bool {
		if results[i].CreatedAt.Equal(results[j].CreatedAt) {
			return results[i].ID < results[j].ID
		}
		return results[i].CreatedAt.Before(results[j].CreatedAt)
	})

	return results, nil
}

// Update updates an existing task in memory.
func (r *MemoryRepository) Update(ctx context.Context, task domain.Task) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[task.ID]; !exists {
		return domain.ErrTaskNotFound
	}

	r.tasks[task.ID] = task
	return nil
}

// Delete removes a task by ID.
func (r *MemoryRepository) Delete(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[id]; !exists {
		return domain.ErrTaskNotFound
	}

	delete(r.tasks, id)
	return nil
}

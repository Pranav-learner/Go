package task

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"taskflow/internal/domain"
)

func TestMemoryRepository_CRUD(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepository()

	now := time.Now().UTC()
	task1 := domain.Task{
		ID:          "task-1",
		Title:       "Test Task 1",
		Description: "Description 1",
		Status:      domain.StatusPending,
		Priority:    domain.PriorityHigh,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// 1. Save
	if err := repo.Save(ctx, task1); err != nil {
		t.Fatalf("expected save to succeed, got %v", err)
	}

	// 2. FindByID
	found, err := repo.FindByID(ctx, "task-1")
	if err != nil {
		t.Fatalf("expected find to succeed, got %v", err)
	}
	if found.Title != task1.Title {
		t.Errorf("expected title %s, got %s", task1.Title, found.Title)
	}

	// 3. FindByID Not Found
	_, err = repo.FindByID(ctx, "non-existent")
	if err != domain.ErrTaskNotFound {
		t.Errorf("expected ErrTaskNotFound, got %v", err)
	}

	// 4. Update
	task1.Title = "Updated Title"
	task1.Status = domain.StatusInProgress
	if err := repo.Update(ctx, task1); err != nil {
		t.Fatalf("expected update to succeed, got %v", err)
	}

	updated, err := repo.FindByID(ctx, "task-1")
	if err != nil {
		t.Fatalf("expected find after update to succeed, got %v", err)
	}
	if updated.Title != "Updated Title" || updated.Status != domain.StatusInProgress {
		t.Errorf("task fields did not update properly: %+v", updated)
	}

	// 5. Update non-existent
	err = repo.Update(ctx, domain.Task{ID: "unknown"})
	if err != domain.ErrTaskNotFound {
		t.Errorf("expected ErrTaskNotFound on missing task update, got %v", err)
	}

	// 6. Delete
	if err := repo.Delete(ctx, "task-1"); err != nil {
		t.Fatalf("expected delete to succeed, got %v", err)
	}

	// Verify deleted
	_, err = repo.FindByID(ctx, "task-1")
	if err != domain.ErrTaskNotFound {
		t.Errorf("expected ErrTaskNotFound after delete, got %v", err)
	}

	// Delete non-existent
	if err := repo.Delete(ctx, "task-1"); err != domain.ErrTaskNotFound {
		t.Errorf("expected ErrTaskNotFound deleting missing task, got %v", err)
	}
}

func TestMemoryRepository_FindAllAndFilter(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepository()

	now := time.Now().UTC()
	tasks := []domain.Task{
		{ID: "1", Title: "Task A", Status: domain.StatusPending, Priority: domain.PriorityLow, CreatedAt: now.Add(1 * time.Second), UpdatedAt: now},
		{ID: "2", Title: "Task B", Status: domain.StatusPending, Priority: domain.PriorityHigh, CreatedAt: now.Add(2 * time.Second), UpdatedAt: now},
		{ID: "3", Title: "Task C", Status: domain.StatusCompleted, Priority: domain.PriorityHigh, CreatedAt: now.Add(3 * time.Second), UpdatedAt: now},
	}

	for _, task := range tasks {
		if err := repo.Save(ctx, task); err != nil {
			t.Fatalf("failed to save task: %v", err)
		}
	}

	// 1. Find all without filter
	all, err := repo.FindAll(ctx, domain.TaskFilter{})
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	if len(all) != 3 {
		t.Fatalf("expected 3 tasks, got %d", len(all))
	}
	// Verify sorting by CreatedAt
	if all[0].ID != "1" || all[1].ID != "2" || all[2].ID != "3" {
		t.Errorf("tasks are not ordered by CreatedAt: %+v", all)
	}

	// 2. Filter by status: PENDING
	pendingStatus := domain.StatusPending
	pendingTasks, err := repo.FindAll(ctx, domain.TaskFilter{Status: &pendingStatus})
	if err != nil {
		t.Fatalf("FindAll with status filter failed: %v", err)
	}
	if len(pendingTasks) != 2 {
		t.Fatalf("expected 2 pending tasks, got %d", len(pendingTasks))
	}

	// 3. Filter by priority: HIGH
	highPriority := domain.PriorityHigh
	highTasks, err := repo.FindAll(ctx, domain.TaskFilter{Priority: &highPriority})
	if err != nil {
		t.Fatalf("FindAll with priority filter failed: %v", err)
	}
	if len(highTasks) != 2 {
		t.Fatalf("expected 2 high priority tasks, got %d", len(highTasks))
	}

	// 4. Filter by both status: PENDING and priority: HIGH
	combined, err := repo.FindAll(ctx, domain.TaskFilter{Status: &pendingStatus, Priority: &highPriority})
	if err != nil {
		t.Fatalf("FindAll with combined filter failed: %v", err)
	}
	if len(combined) != 1 || combined[0].ID != "2" {
		t.Fatalf("expected 1 task with ID 2, got %+v", combined)
	}
}

func TestMemoryRepository_ConcurrentAccess(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepository()

	const numWorkers = 50
	const opsPerWorker = 20

	var wg sync.WaitGroup

	// Pre-populate with some tasks
	for i := 0; i < 10; i++ {
		_ = repo.Save(ctx, domain.Task{
			ID:        fmt.Sprintf("initial-%d", i),
			Title:     "Initial",
			Status:    domain.StatusPending,
			Priority:  domain.PriorityMedium,
			CreatedAt: time.Now(),
		})
	}

	// Launch concurrent writers and readers
	for w := 0; w < numWorkers; w++ {
		wg.Add(2)

		// Concurrent Writer
		go func(workerID int) {
			defer wg.Done()
			for i := 0; i < opsPerWorker; i++ {
				taskID := fmt.Sprintf("worker-%d-task-%d", workerID, i)
				_ = repo.Save(ctx, domain.Task{
					ID:        taskID,
					Title:     "Concurrent Task",
					Status:    domain.StatusInProgress,
					Priority:  domain.PriorityHigh,
					CreatedAt: time.Now(),
				})
				_ = repo.Update(ctx, domain.Task{
					ID:        taskID,
					Title:     "Updated Concurrent Task",
					Status:    domain.StatusCompleted,
					Priority:  domain.PriorityLow,
					CreatedAt: time.Now(),
				})
			}
		}(w)

		// Concurrent Reader
		go func(workerID int) {
			defer wg.Done()
			for i := 0; i < opsPerWorker; i++ {
				_, _ = repo.FindByID(ctx, fmt.Sprintf("initial-%d", i%10))
				_, _ = repo.FindAll(ctx, domain.TaskFilter{})
			}
		}(w)
	}

	wg.Wait()
}

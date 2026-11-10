package task

import (
	"context"
	"testing"
	"time"

	"taskflow/internal/domain"
)

func TestService_Create(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepository()
	svc := NewService(repo)

	// 1. Success case with default priority
	created, err := svc.Create(ctx, CreateTaskInput{
		Title:       "Learn Go Concurrency",
		Description: "Deep dive into goroutines and channels",
	})
	if err != nil {
		t.Fatalf("expected successful creation, got %v", err)
	}

	if created.ID == "" {
		t.Error("expected non-empty UUID")
	}
	if created.Status != domain.StatusPending {
		t.Errorf("expected default status PENDING, got %s", created.Status)
	}
	if created.Priority != domain.PriorityMedium {
		t.Errorf("expected default priority MEDIUM, got %s", created.Priority)
	}

	// 2. Empty title validation
	_, err = svc.Create(ctx, CreateTaskInput{
		Title: "   ",
	})
	if err != domain.ErrEmptyTitle {
		t.Errorf("expected ErrEmptyTitle, got %v", err)
	}

	// 3. Invalid priority validation
	invalidPriority := domain.TaskPriority("CRITICAL")
	_, err = svc.Create(ctx, CreateTaskInput{
		Title:    "Valid Title",
		Priority: &invalidPriority,
	})
	if err != domain.ErrInvalidPriority {
		t.Errorf("expected ErrInvalidPriority, got %v", err)
	}
}

func TestService_Complete(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepository()
	svc := NewService(repo)

	created, err := svc.Create(ctx, CreateTaskInput{
		Title: "Task to complete",
	})
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}

	completed, err := svc.Complete(ctx, created.ID)
	if err != nil {
		t.Fatalf("failed to complete task: %v", err)
	}
	if completed.Status != domain.StatusCompleted {
		t.Errorf("expected status COMPLETED, got %s", completed.Status)
	}

	// Completing non-existent task
	_, err = svc.Complete(ctx, "non-existent-id")
	if err != domain.ErrTaskNotFound {
		t.Errorf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestService_Update(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepository()
	svc := NewService(repo)

	created, err := svc.Create(ctx, CreateTaskInput{
		Title: "Original Title",
	})
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}

	due := time.Now().Add(24 * time.Hour)
	updated, err := svc.Update(ctx, created.ID, UpdateTaskInput{
		Title:       "Updated Title",
		Description: "Updated Description",
		Status:      domain.StatusInProgress,
		Priority:    domain.PriorityHigh,
		DueDate:     &due,
	})
	if err != nil {
		t.Fatalf("failed to update task: %v", err)
	}

	if updated.Title != "Updated Title" || updated.Status != domain.StatusInProgress || updated.Priority != domain.PriorityHigh {
		t.Errorf("updated task fields did not match: %+v", updated)
	}

	// Empty title validation on update
	_, err = svc.Update(ctx, created.ID, UpdateTaskInput{
		Title:    "",
		Status:   domain.StatusPending,
		Priority: domain.PriorityLow,
	})
	if err != domain.ErrEmptyTitle {
		t.Errorf("expected ErrEmptyTitle, got %v", err)
	}

	// Update non-existent task
	_, err = svc.Update(ctx, "missing-id", UpdateTaskInput{
		Title:    "Valid Title",
		Status:   domain.StatusPending,
		Priority: domain.PriorityLow,
	})
	if err != domain.ErrTaskNotFound {
		t.Errorf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestService_Delete(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepository()
	svc := NewService(repo)

	created, _ := svc.Create(ctx, CreateTaskInput{Title: "To delete"})

	if err := svc.Delete(ctx, created.ID); err != nil {
		t.Fatalf("expected successful delete, got %v", err)
	}

	// Verify deleted
	_, err := svc.GetByID(ctx, created.ID)
	if err != domain.ErrTaskNotFound {
		t.Errorf("expected ErrTaskNotFound, got %v", err)
	}

	// Delete non-existent
	if err := svc.Delete(ctx, created.ID); err != domain.ErrTaskNotFound {
		t.Errorf("expected ErrTaskNotFound deleting missing task, got %v", err)
	}
}

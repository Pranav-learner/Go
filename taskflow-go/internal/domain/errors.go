package domain

import "errors"

// Sentinel domain errors
var (
	// ErrTaskNotFound is returned when a requested task does not exist.
	ErrTaskNotFound = errors.New("task not found")

	// ErrEmptyTitle is returned when creating or updating a task with a blank title.
	ErrEmptyTitle = errors.New("task title is required")

	// ErrInvalidStatus is returned when an invalid status string is provided.
	ErrInvalidStatus = errors.New("invalid task status: must be PENDING, IN_PROGRESS, or COMPLETED")

	// ErrInvalidPriority is returned when an invalid priority string is provided.
	ErrInvalidPriority = errors.New("invalid task priority: must be LOW, MEDIUM, or HIGH")

	// ErrInvalidTaskID is returned when a task ID is empty or malformed.
	ErrInvalidTaskID = errors.New("invalid task id")
)

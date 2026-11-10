package domain

import (
	"crypto/rand"
	"fmt"
	"time"
)

// TaskStatus represents the current state of a task.
type TaskStatus string

const (
	StatusPending    TaskStatus = "PENDING"
	StatusInProgress TaskStatus = "IN_PROGRESS"
	StatusCompleted  TaskStatus = "COMPLETED"
)

// IsValid validates whether the status is one of the supported constants.
func (s TaskStatus) IsValid() bool {
	switch s {
	case StatusPending, StatusInProgress, StatusCompleted:
		return true
	default:
		return false
	}
}

// TaskPriority represents the urgency level of a task.
type TaskPriority string

const (
	PriorityLow    TaskPriority = "LOW"
	PriorityMedium TaskPriority = "MEDIUM"
	PriorityHigh   TaskPriority = "HIGH"
)

// IsValid validates whether the priority is one of the supported constants.
func (p TaskPriority) IsValid() bool {
	switch p {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return true
	default:
		return false
	}
}

// Task represents the primary domain entity of TaskFlow.
type Task struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Status      TaskStatus   `json:"status"`
	Priority    TaskPriority `json:"priority"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
	DueDate     *time.Time   `json:"dueDate,omitempty"`
}

// TaskFilter encapsulates criteria for filtering tasks in list queries.
type TaskFilter struct {
	Status   *TaskStatus
	Priority *TaskPriority
}

// NewUUID generates an RFC 4122 Version 4 UUID using standard-library crypto/rand.
// This provides cryptographically secure unique identifiers without third-party dependencies.
func NewUUID() (string, error) {
	var uuid [16]byte
	_, err := rand.Read(uuid[:])
	if err != nil {
		return "", fmt.Errorf("crypto/rand read failed: %w", err)
	}

	// Set version to 4 (0100 in bits 4-7 of byte 6)
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	// Set variant to RFC 4122 (10 in bits 6-7 of byte 8)
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuid[0:4],
		uuid[4:6],
		uuid[6:8],
		uuid[8:10],
		uuid[10:],
	), nil
}


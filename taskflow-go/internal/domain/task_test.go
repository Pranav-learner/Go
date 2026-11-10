package domain

import (
	"testing"
)

func TestTaskStatus_IsValid(t *testing.T) {
	tests := []struct {
		status TaskStatus
		valid  bool
	}{
		{StatusPending, true},
		{StatusInProgress, true},
		{StatusCompleted, true},
		{TaskStatus("INVALID"), false},
		{TaskStatus(""), false},
	}

	for _, tc := range tests {
		if got := tc.status.IsValid(); got != tc.valid {
			t.Errorf("status %q: got %v, expected %v", tc.status, got, tc.valid)
		}
	}
}

func TestTaskPriority_IsValid(t *testing.T) {
	tests := []struct {
		priority TaskPriority
		valid    bool
	}{
		{PriorityLow, true},
		{PriorityMedium, true},
		{PriorityHigh, true},
		{TaskPriority("URGENT"), false},
		{TaskPriority(""), false},
	}

	for _, tc := range tests {
		if got := tc.priority.IsValid(); got != tc.valid {
			t.Errorf("priority %q: got %v, expected %v", tc.priority, got, tc.valid)
		}
	}
}

func TestNewUUID(t *testing.T) {
	id1, err := NewUUID()
	if err != nil {
		t.Fatalf("failed to generate uuid: %v", err)
	}

	if len(id1) != 36 {
		t.Fatalf("expected uuid length 36, got %d (%s)", len(id1), id1)
	}

	id2, err := NewUUID()
	if err != nil {
		t.Fatalf("failed to generate second uuid: %v", err)
	}

	if id1 == id2 {
		t.Fatalf("expected unique uuids, got duplicate: %s", id1)
	}
}

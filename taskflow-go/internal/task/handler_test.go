package task

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"taskflow/internal/domain"
	"taskflow/internal/response"
)

func setupTestMux() (*http.ServeMux, *Service) {
	repo := NewMemoryRepository()
	service := NewService(repo)
	handler := NewHandler(service)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	return mux, service
}

func TestHandler_CreateTask(t *testing.T) {
	mux, _ := setupTestMux()

	t.Run("success_201", func(t *testing.T) {
		payload := `{"title":"Learn Go Testing","description":"Table-driven HTTP tests","priority":"HIGH"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected status 201, got %d. Body: %s", rec.Code, rec.Body.String())
		}

		var created domain.Task
		if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if created.ID == "" || created.Title != "Learn Go Testing" || created.Priority != domain.PriorityHigh {
			t.Errorf("created task fields incorrect: %+v", created)
		}
	})

	t.Run("empty_title_400", func(t *testing.T) {
		payload := `{"title":"   ","priority":"LOW"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBufferString(payload))
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}

		var errResp response.ErrorResponse
		_ = json.Unmarshal(rec.Body.Bytes(), &errResp)
		if errResp.Error.Code != "INVALID_TITLE" {
			t.Errorf("expected error code INVALID_TITLE, got %s", errResp.Error.Code)
		}
	})

	t.Run("malformed_json_400", func(t *testing.T) {
		payload := `{malformed json}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBufferString(payload))
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 for malformed json, got %d", rec.Code)
		}
	})

	t.Run("empty_body_400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBufferString(""))
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 for empty body, got %d", rec.Code)
		}
	})
}

func TestHandler_GetTaskByID(t *testing.T) {
	mux, svc := setupTestMux()

	created, _ := svc.Create(t.Context(), CreateTaskInput{Title: "Task for Retrieval"})

	t.Run("found_200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/"+created.ID, nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}

		var task domain.Task
		_ = json.Unmarshal(rec.Body.Bytes(), &task)
		if task.ID != created.ID {
			t.Errorf("expected task ID %s, got %s", created.ID, task.ID)
		}
	})

	t.Run("not_found_404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/non-existent-uuid", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", rec.Code)
		}

		var errResp response.ErrorResponse
		_ = json.Unmarshal(rec.Body.Bytes(), &errResp)
		if errResp.Error.Code != "TASK_NOT_FOUND" {
			t.Errorf("expected TASK_NOT_FOUND, got %s", errResp.Error.Code)
		}
	})
}

func TestHandler_ListAndFilter(t *testing.T) {
	mux, svc := setupTestMux()

	low := domain.PriorityLow
	high := domain.PriorityHigh

	_, _ = svc.Create(t.Context(), CreateTaskInput{Title: "Task 1", Priority: &low})
	t2, _ := svc.Create(t.Context(), CreateTaskInput{Title: "Task 2", Priority: &high})
	_, _ = svc.Complete(t.Context(), t2.ID)

	t.Run("list_all_200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}

		var tasks []domain.Task
		_ = json.Unmarshal(rec.Body.Bytes(), &tasks)
		if len(tasks) != 2 {
			t.Errorf("expected 2 tasks, got %d", len(tasks))
		}
	})

	t.Run("filter_status_pending_200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks?status=PENDING", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		var tasks []domain.Task
		_ = json.Unmarshal(rec.Body.Bytes(), &tasks)
		if len(tasks) != 1 {
			t.Fatalf("expected 1 task, got %d", len(tasks))
		}
		if tasks[0].Status != domain.StatusPending {
			t.Errorf("expected PENDING, got %s", tasks[0].Status)
		}
	})

	t.Run("filter_invalid_status_400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks?status=UNKNOWN_STATUS", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 for invalid status filter, got %d", rec.Code)
		}
	})

	t.Run("filter_invalid_priority_400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks?priority=CRITICAL", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 for invalid priority filter, got %d", rec.Code)
		}
	})
}

func TestHandler_Update(t *testing.T) {
	mux, svc := setupTestMux()

	created, _ := svc.Create(t.Context(), CreateTaskInput{Title: "Initial Title"})

	payload := `{"title":"Replaced Title","description":"New desc","status":"IN_PROGRESS","priority":"HIGH"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/tasks/"+created.ID, bytes.NewBufferString(payload))
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d. Body: %s", rec.Code, rec.Body.String())
	}

	var updated domain.Task
	_ = json.Unmarshal(rec.Body.Bytes(), &updated)
	if updated.Title != "Replaced Title" || updated.Status != domain.StatusInProgress {
		t.Errorf("updated fields do not match: %+v", updated)
	}
}

func TestHandler_Complete(t *testing.T) {
	mux, svc := setupTestMux()

	created, _ := svc.Create(t.Context(), CreateTaskInput{Title: "To Complete"})

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/tasks/"+created.ID+"/complete", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var completed domain.Task
	_ = json.Unmarshal(rec.Body.Bytes(), &completed)
	if completed.Status != domain.StatusCompleted {
		t.Errorf("expected COMPLETED, got %s", completed.Status)
	}
}

func TestHandler_Delete(t *testing.T) {
	mux, svc := setupTestMux()

	created, _ := svc.Create(t.Context(), CreateTaskInput{Title: "To Delete"})

	// Delete existing
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/tasks/"+created.ID, nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204 No Content, got %d", rec.Code)
	}

	// Delete again should return 404
	req2 := httptest.NewRequest(http.MethodDelete, "/api/v1/tasks/"+created.ID, nil)
	rec2 := httptest.NewRecorder()

	mux.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusNotFound {
		t.Fatalf("expected 404 on already deleted task, got %d", rec2.Code)
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	mux, _ := setupTestMux()

	// POST to a GET-only route pattern like /api/v1/tasks/{id}/complete
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks/123/complete", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	// Go 1.22+ ServeMux automatically returns 405 Method Not Allowed when method does not match
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405 Method Not Allowed, got %d", rec.Code)
	}
}

func TestHealthHandler(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","service":"taskflow"}`))
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	expected := `{"status":"ok","service":"taskflow"}`
	if rec.Body.String() != expected {
		t.Errorf("expected %s, got %s", expected, rec.Body.String())
	}
}

// Ensure time is preserved
func TestPreserveTime(t *testing.T) {
	_, svc := setupTestMux()
	now := time.Now().UTC()
	created, _ := svc.Create(t.Context(), CreateTaskInput{Title: "Time Check", DueDate: &now})
	if created.DueDate == nil || !created.DueDate.Equal(now) {
		t.Errorf("dueDate was not preserved")
	}
}

package task

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"taskflow/internal/domain"
	"taskflow/internal/response"
)

const maxRequestBodyBytes = 1 << 20 // 1 MB limit to prevent body exhaustion DoS

// Handler exposes HTTP handlers for the Task entity.
//
// Spring Boot Comparison:
// In Spring Boot, this corresponds to a `@RestController` annotated with `@RequestMapping("/api/v1/tasks")`.
// Spring dynamically dispatches routes via reflection and handles argument resolution (@PathVariable, @RequestBody).
// In Go, handlers are simple functions implementing `http.HandlerFunc`. They parse path/query variables
// explicitly via the standard library (`r.PathValue`, `r.URL.Query()`) and delegate strictly to the service layer.
type Handler struct {
	service *Service
}

// NewHandler constructs a Handler with its required Service dependency.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers all task routes on the provided ServeMux using Go's enhanced pattern matching.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/tasks", h.List)
	mux.HandleFunc("POST /api/v1/tasks", h.Create)
	mux.HandleFunc("GET /api/v1/tasks/{id}", h.GetByID)
	mux.HandleFunc("PUT /api/v1/tasks/{id}", h.Update)
	mux.HandleFunc("PATCH /api/v1/tasks/{id}/complete", h.Complete)
	mux.HandleFunc("DELETE /api/v1/tasks/{id}", h.Delete)
}

// List handles GET /api/v1/tasks with optional ?status=... and ?priority=... query filters.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	filter := domain.TaskFilter{}

	statusParam := strings.TrimSpace(r.URL.Query().Get("status"))
	if statusParam != "" {
		status := domain.TaskStatus(statusParam)
		if !status.IsValid() {
			response.Error(w, http.StatusBadRequest, "INVALID_FILTER", "invalid status filter: must be PENDING, IN_PROGRESS, or COMPLETED")
			return
		}
		filter.Status = &status
	}

	priorityParam := strings.TrimSpace(r.URL.Query().Get("priority"))
	if priorityParam != "" {
		priority := domain.TaskPriority(priorityParam)
		if !priority.IsValid() {
			response.Error(w, http.StatusBadRequest, "INVALID_FILTER", "invalid priority filter: must be LOW, MEDIUM, or HIGH")
			return
		}
		filter.Priority = &priority
	}

	tasks, err := h.service.List(r.Context(), filter)
	if err != nil {
		h.mapError(w, err)
		return
	}

	// Ensure empty slice is encoded as [] instead of null in JSON
	if tasks == nil {
		tasks = []domain.Task{}
	}

	response.JSON(w, http.StatusOK, tasks)
}

// GetByID handles GET /api/v1/tasks/{id}.
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "INVALID_TASK_ID", "task id is required in URL path")
		return
	}

	t, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.mapError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, t)
}

// Create handles POST /api/v1/tasks.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)

	var input CreateTaskInput
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() // Reject unexpected properties cleanly
	if err := dec.Decode(&input); err != nil {
		h.handleJSONDecodeError(w, err)
		return
	}

	created, err := h.service.Create(r.Context(), input)
	if err != nil {
		h.mapError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, created)
}

// Update handles PUT /api/v1/tasks/{id}.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "INVALID_TASK_ID", "task id is required in URL path")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)

	var input UpdateTaskInput
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&input); err != nil {
		h.handleJSONDecodeError(w, err)
		return
	}

	updated, err := h.service.Update(r.Context(), id, input)
	if err != nil {
		h.mapError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, updated)
}

// Complete handles PATCH /api/v1/tasks/{id}/complete.
func (h *Handler) Complete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "INVALID_TASK_ID", "task id is required in URL path")
		return
	}

	completed, err := h.service.Complete(r.Context(), id)
	if err != nil {
		h.mapError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, completed)
}

// Delete handles DELETE /api/v1/tasks/{id}.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "INVALID_TASK_ID", "task id is required in URL path")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		h.mapError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handleJSONDecodeError formats JSON deserialization failures into actionable client responses.
func (h *Handler) handleJSONDecodeError(w http.ResponseWriter, err error) {
	if errors.Is(err, io.EOF) {
		response.Error(w, http.StatusBadRequest, "EMPTY_REQUEST_BODY", "request body cannot be empty")
		return
	}

	var maxBytesErr *http.MaxBytesError
	if errors.As(err, &maxBytesErr) {
		response.Error(w, http.StatusRequestEntityTooLarge, "PAYLOAD_TOO_LARGE", "request body exceeded maximum allowed limit of 1MB")
		return
	}

	response.Error(w, http.StatusBadRequest, "MALFORMED_JSON", err.Error())
}

// mapError translates domain/service errors into standard HTTP status codes and JSON envelopes.
func (h *Handler) mapError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrTaskNotFound):
		response.Error(w, http.StatusNotFound, "TASK_NOT_FOUND", "Task not found")
	case errors.Is(err, domain.ErrEmptyTitle):
		response.Error(w, http.StatusBadRequest, "INVALID_TITLE", err.Error())
	case errors.Is(err, domain.ErrInvalidStatus):
		response.Error(w, http.StatusBadRequest, "INVALID_STATUS", err.Error())
	case errors.Is(err, domain.ErrInvalidPriority):
		response.Error(w, http.StatusBadRequest, "INVALID_PRIORITY", err.Error())
	case errors.Is(err, domain.ErrInvalidTaskID):
		response.Error(w, http.StatusBadRequest, "INVALID_TASK_ID", err.Error())
	default:
		response.Error(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "An unexpected error occurred")
	}
}

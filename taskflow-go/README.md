# TaskFlow — Task Management & Collaboration Platform (Go Backend)

> **Phase 0–2 Implementation**: Core foundations, idiomatic REST API with `net/http`, layered architecture (Handler/Service/Repository), concurrency-safe in-memory storage, custom middleware, explicit dependency injection, and zero external dependencies.

---

## 1. Project Overview

**TaskFlow** is a high-performance, concurrent task management backend built with idiomatic Go. This project is the counterpart to the TaskFlow Spring Boot backend, designed specifically to showcase how real-world backend architectures translate between the **Java / Spring Boot** and **Go** ecosystems without mechanically forcing object-oriented patterns onto Go.

> [!IMPORTANT]
> **Data Persistence Notice**: Tasks are currently stored in-memory using a concurrency-safe hash map protected by `sync.RWMutex`. All stored data will be reset when the server is restarted. Persistent database storage (PostgreSQL) is planned for future phases.

---

## 2. Technology Stack

* **Language**: Go 1.22+ (running on standard toolchain, standard library only)
* **HTTP Server & Routing**: Standard library `net/http` (`http.ServeMux` with Go 1.22+ pattern matching)
* **Serialization**: Standard library `encoding/json`
* **Concurrency Primitives**: `sync.RWMutex`, `sync.WaitGroup`, channels, `context.Context`
* **ID Generation**: RFC 4122 v4 UUID generated via standard `crypto/rand`
* **Logging**: Standard library structured logger (`log/slog`)
* **Testing**: Standard `testing` and `net/http/httptest` with race detector (`-race`)
* **External Dependencies**: **None (0 dependencies)**

---

## 3. Architecture & Dependency Flow

The system strictly follows **Clean Architecture** and **Separation of Concerns**:

```
                         HTTP Request
                              │
                              ▼
                   ┌──────────────────────┐
                   │  Recovery Middleware │ (recovers from panics, yields safe 500)
                   └──────────┬───────────┘
                              │
                              ▼
                   ┌──────────────────────┐
                   │  Logging Middleware  │ (logs method, path, status, latency)
                   └──────────┬───────────┘
                              │
                              ▼
                   ┌──────────────────────┐
                   │   http.ServeMux      │ (Go 1.22+ enhanced URL pattern matching)
                   └──────────┬───────────┘
                              │
               ┌──────────────┴──────────────┐
               ▼                             ▼
     ┌───────────────────┐         ┌───────────────────┐
     │   Health Handler  │         │    Task Handler   │ (Transport layer: parses JSON,
     └───────────────────┘         └─────────┬─────────┘  path & query params, returns HTTP codes)
                                             │
                                             ▼
                                   ┌───────────────────┐
                                   │    Task Service   │ (Business layer: validation, defaults,
                                   └─────────┬─────────┘  lifecycle rules, independent of HTTP)
                                             │
                                             ▼
                                   ┌───────────────────┐
                                   │ Task Repository   │ (Interface contract for storage)
                                   └─────────┬─────────┘
                                             │
                                             ▼
                                   ┌───────────────────┐
                                   │ Memory Repository │ (Persistence layer: concurrency-safe
                                   └───────────────────┘  map protected by sync.RWMutex)
```

### Dependency Injection Pipeline (`cmd/api/main.go`)

Dependencies are instantiated explicitly from the bottom up and injected via constructors:

$$\text{MemoryRepository} \longrightarrow \text{TaskService} \longrightarrow \text{TaskHandler} \longrightarrow \text{HTTP Server}$$

---

## 4. Folder Structure

```
taskflow-go/
│
├── go.mod                      # Module definition (taskflow)
├── README.md                   # Complete architectural guide and API documentation
│
├── cmd/
│   └── api/
│       └── main.go             # Application entry point, explicit DI wiring, graceful shutdown
│
└── internal/                   # Code private to this module (Go internal package rule)
    ├── domain/
    │   ├── task.go             # Task entity, Status/Priority enums, TaskFilter, UUID generator
    │   ├── task_test.go        # Unit tests for domain validations and UUID generation
    │   └── errors.go           # Domain sentinel errors (ErrTaskNotFound, ErrEmptyTitle, etc.)
    │
    ├── health/
    │   └── handler.go          # Health check endpoint (/api/v1/health)
    │
    ├── middleware/
    │   ├── chain.go            # Middleware chaining helper (Chain)
    │   ├── logging.go          # Structured HTTP request logger using log/slog
    │   ├── recovery.go         # Panic recovery middleware returning safe 500 JSON
    │   └── middleware_test.go  # Tests for panic recovery and execution ordering
    │
    ├── response/
    │   └── json.go             # Standardized JSON response and error envelope helpers
    │
    └── task/
        ├── repository.go       # Repository interface definition
        ├── memory_repository.go# Concurrency-safe in-memory storage (sync.RWMutex)
        ├── memory_repository_test.go # Repository tests (CRUD + race detector)
        ├── service.go          # Business validation and task lifecycle logic
        ├── service_test.go     # Service layer unit tests
        ├── handler.go          # REST handlers (GET, POST, PUT, PATCH, DELETE)
        └── handler_test.go     # HTTP integration tests via net/http/httptest
```

---

## 5. Architectural Comparison: Go vs. Spring Boot & FastAPI

| Architectural Area | Spring Boot / Java 21 | FastAPI / Python | Idiomatic Go (TaskFlow) |
| :--- | :--- | :--- | :--- |
| **Dependency Injection** | **Implicit / Magic**: Spring IoC container scans classes (`@Component`, `@Autowired`), building dependency graphs via reflection. | **Dependency injection via `Depends()`**: Functional injection resolved at route time. | **Explicit**: All dependencies are manually constructed in `main.go` and passed as parameters. Zero reflection, instant startup, 100% compile-time verified. |
| **Routing & Controllers** | `@RestController` with annotations (`@GetMapping`, `@PathVariable`, `@RequestBody`). | Decorated functions (`@app.get("/items/{id}")`). | Handlers implement standard `http.HandlerFunc` and register with `http.ServeMux`. Path parameters read via `r.PathValue("id")`. |
| **Validation** | Bean Validation (`@NotNull`, `@Size`, `@Valid`) triggering hidden exception filters. | Pydantic model type checking. | Explicit domain validation methods (`status.IsValid()`, `strings.TrimSpace(title)`). Predictable, easy to debug. |
| **Storage & Data Access** | Spring Data JPA / Hibernate generating queries dynamically from interface method names. | SQLAlchemy / Tortoise ORM with async session managers. | Explicit `Repository` interface satisfied implicitly by `MemoryRepository`. Deep copies prevent accidental state mutation. |
| **Concurrency & Synchronization** | `@Synchronized`, `ConcurrentHashMap`, or database-level transactions (`@Transactional`). | Single-threaded event loop (`asyncio`) or Celery workers. | Native Go concurrency: `sync.RWMutex` reader-writer lock allows hundreds of concurrent readers while isolating writers. |
| **Middleware** | Servlet Filters / Spring Security Filter Chains. | ASGI middleware (`app.add_middleware(...)`). | Pure higher-order functions: `func(http.Handler) http.Handler` wrapped via simple closure chaining. |
| **Error Handling** | `@ControllerAdvice` + `@ExceptionHandler` intercepting thrown exceptions. | Custom exception handlers (`@app.exception_handler`). | Multiple return values `(Result, error)` inspected explicitly with `errors.Is()`. No exception-driven control flow. |

---

## 6. API Reference

All requests and responses use `Content-Type: application/json`.

### Health Check

#### `GET /api/v1/health`
Verifies server health.

* **Response (`200 OK`)**:
  ```json
  {
    "status": "ok",
    "service": "taskflow"
  }
  ```

---

### Task Endpoints

#### 1. List Tasks
`GET /api/v1/tasks`

* **Query Parameters (Optional)**:
  * `status`: Filter by status (`PENDING`, `IN_PROGRESS`, `COMPLETED`)
  * `priority`: Filter by priority (`LOW`, `MEDIUM`, `HIGH`)
* **Example**: `GET /api/v1/tasks?status=PENDING&priority=HIGH`
* **Response (`200 OK`)**:
  ```json
  [
    {
      "id": "c62372e9-4ec4-4c59-bfd2-1c296f86c2e1",
      "title": "Build TaskFlow Go Backend",
      "description": "Implement idiomatic clean architecture in Go",
      "status": "PENDING",
      "priority": "HIGH",
      "createdAt": "2026-11-10T10:00:00Z",
      "updatedAt": "2026-11-10T10:00:00Z",
      "dueDate": "2026-12-01T18:00:00Z"
    }
  ]
  ```

#### 2. Get Task by ID
`GET /api/v1/tasks/{id}`

* **Response (`200 OK`)**: Single Task JSON
* **Response (`404 Not Found`)**:
  ```json
  {
    "error": {
      "code": "TASK_NOT_FOUND",
      "message": "Task not found"
    }
  }
  ```

#### 3. Create Task
`POST /api/v1/tasks`

* **Request Body**:
  ```json
  {
    "title": "Implement Unit Tests",
    "description": "Write table-driven tests for services and handlers",
    "priority": "HIGH",
    "dueDate": "2026-11-20T12:00:00Z"
  }
  ```
* **Response (`201 Created`)**: Returns newly created Task with generated UUID, `status: "PENDING"`, and timestamps.
* **Response (`400 Bad Request`)**:
  ```json
  {
    "error": {
      "code": "INVALID_TITLE",
      "message": "task title is required"
    }
  }
  ```

#### 4. Update Task (Replace)
`PUT /api/v1/tasks/{id}`

* **Request Body**:
  ```json
  {
    "title": "Implement Unit Tests (Updated)",
    "description": "Coverage > 90%",
    "status": "IN_PROGRESS",
    "priority": "HIGH"
  }
  ```
* **Response (`200 OK`)**: Returns updated Task with refreshed `updatedAt` timestamp.

#### 5. Complete Task
`PATCH /api/v1/tasks/{id}/complete`

* **Response (`200 OK`)**: Returns Task with `status: "COMPLETED"`.

#### 6. Delete Task
`DELETE /api/v1/tasks/{id}`

* **Response (`204 No Content`)**: Empty response body.
* **Response (`404 Not Found`)**: If task does not exist.

---

## 7. How to Setup and Run

### Prerequisites
* Go 1.22 or higher installed.

### Run Server

```bash
cd taskflow-go
go run ./cmd/api
```

The server starts on port `8080` (or the port defined in the `PORT` environment variable):

```
{"time":"2026-11-10T10:45:00Z","level":"INFO","msg":"starting TaskFlow backend service..."}
{"time":"2026-11-10T10:45:00Z","level":"INFO","msg":"server listening and serving","addr":":8080"}
```

### Quick cURL Verification

```bash
# 1. Health Check
curl -i http://localhost:8080/api/v1/health

# 2. Create Task
curl -i -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Learn Go","description":"Master idiomatic backend patterns","priority":"HIGH"}'

# 3. List All Tasks
curl -i http://localhost:8080/api/v1/tasks

# 4. Filter Tasks
curl -i "http://localhost:8080/api/v1/tasks?status=PENDING&priority=HIGH"

# 5. Mark as Complete (replace <ID> with actual UUID)
curl -i -X PATCH http://localhost:8080/api/v1/tasks/<ID>/complete

# 6. Delete Task
curl -i -X DELETE http://localhost:8080/api/v1/tasks/<ID>
```

---

## 8. Testing & Quality Verification

Run the entire test suite including unit tests, HTTP integration tests, and the Go race detector:

```bash
# 1. Run all tests with verbose output
go test -v ./...

# 2. Run all tests with the Race Detector enabled (Strict Concurrency Verification)
go test -race ./...

# 3. Run static analysis
go vet ./...

# 4. Check formatting
gofmt -s -l .
```

---

## 9. Current Status & Future Roadmap

* [x] **Phase 0 — Go Foundations**: Structs, value/pointer semantics, error handling, standard library.
* [x] **Phase 1 — HTTP API**: Full CRUD, routing, JSON serialization, query filtering, logging & recovery middleware.
* [x] **Phase 2 — Go Architecture**: Layered design (Handler/Service/Repository), explicit DI, concurrency-safe in-memory storage, graceful shutdown.
* [ ] **Phase 3 — Persistence (Future)**: PostgreSQL integration using `pgx`, connection pooling, database migrations.
* [ ] **Phase 4 — Security & Auth (Future)**: JWT authentication, password hashing (`bcrypt`), role-based authorization.
* [ ] **Phase 5 — Observability & Deployment (Future)**: Prometheus metrics, OpenTelemetry tracing, Docker containerization.


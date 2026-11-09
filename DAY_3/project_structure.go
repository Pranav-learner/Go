package day3
// Recommended project structure

// taskflow/
// │
// ├── cmd/
// │   └── api/
// │       └── main.go
// │
// ├── internal/
// │   ├── user/
// │   │   ├── handler.go
// │   │   ├── service.go
// │   │   ├── repository.go
// │   │   └── model.go
// │   │
// │   ├── task/
// │   │   ├── handler.go
// │   │   ├── service.go
// │   │   ├── repository.go
// │   │   └── model.go
// │   │
// │   ├── middleware/
// │   │   ├── auth.go
// │   │   └── logging.go
// │   │
// │   └── config/
// │       └── config.go
// │
// ├── migrations/
// │
// ├── go.mod
// └── go.sum


// 6. Understanding internal/
// The internal directory has a special meaning in Go.

// Packages inside an internal directory can only be imported by code within the permitted parent directory tree.
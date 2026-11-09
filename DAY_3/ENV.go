package day3

// Configuration and Environment Variables
// Production requirement
// Never hardcode your database passwords, JWT secrets or production URLs into your source code.

// Configuration should be externalized.
// 1. Reading environment variables
// Go provides:
os.Getenv("PORT")


Example:
package main

import (
    "fmt"
    "os"
)

func main() {

    port := os.Getenv("PORT")

    if port == "" {
        port = "8080"
    }

    fmt.Println("Running on port:", port)
}


Set it:
PORT=9000 go run .


// 2. os.LookupEnv
// There is a difference between a missing variable and an explicitly empty variable.
value, exists := os.LookupEnv("PORT")

if !exists {
    fmt.Println("PORT not configured")
} else {
    fmt.Println("PORT:", value)
}


// 3. Creating a configuration structure
// In real applications, centralize configuration.
type Config struct {
    Port        string
    DatabaseURL string
    JWTSecret   string
}


Load it:
func LoadConfig() (Config, error) {

    cfg := Config{
        Port:        os.Getenv("PORT"),
        DatabaseURL: os.Getenv("DATABASE_URL"),
        JWTSecret:   os.Getenv("JWT_SECRET"),
    }

    if cfg.Port == "" {
        cfg.Port = "8080"
    }

    if cfg.DatabaseURL == "" {
        return Config{}, errors.New("DATABASE_URL is required")
    }

    if cfg.JWTSecret == "" {
        return Config{}, errors.New("JWT_SECRET is required")
    }

    return cfg, nil
}


// 4. .env files
// During local development, you may use:
PORT=8080
DATABASE_URL=postgres://localhost/taskflow
JWT_SECRET=development-secret


// A .env file isn't automatically loaded by Go.

// You need a library or a shell/environment configuration to load it.

// One popular library:
// godotenv — Go environment variable loaderFor production, inject secrets through your deployment platform or secret manager. Don't commit real secrets or production .env files into Git.
// 5. Production configuration principles
// - Validate required environment variables at startup.
// - Set sensible defaults where appropriate.
// - Keep development and production configurations separate.
// - Never log passwords or secrets.
// - Support database connection settings.
// - Configure timeouts and resource limits.
// - Use a secrets manager for sensitive credentials.
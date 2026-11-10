package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"taskflow/internal/health"
	"taskflow/internal/middleware"
	"taskflow/internal/task"
)

// Spring Boot vs Go Architecture Reflection:
//
// 1. Dependency Injection:
//    - Spring Boot: Uses an IoC (Inversion of Control) container that scans classes with `@Component`,
//      `@Service`, and `@Repository`, constructs beans via reflection, and injects them via `@Autowired`.
//    - Go: Uses EXPLICIT DEPENDENCY INJECTION. Every dependency is passed into a constructor function
//      (`NewService(repo)`). This is compile-time checked, transparent, traceable with "go to definition",
//      and incurs zero reflection overhead during startup.
//
// 2. Server Configuration & Lifecycle:
//    - Spring Boot: Embedded Tomcat server is auto-configured via `application.properties` / `application.yml`.
//    - Go: `http.Server` is configured explicitly in code with explicit timeouts to prevent Slowloris attacks
//      and resource leaks.
//
// 3. Graceful Shutdown:
//    - Spring Boot: Handled via Spring ApplicationContext lifecycle events (`@PreDestroy`).
//    - Go: Handled explicitly using OS signal notification (`signal.Notify`) and `server.Shutdown(ctx)`.

func main() {
	// 1. Initialize Structured Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	logger.Info("starting TaskFlow backend service...")

	// 2. Explicit Dependency Wiring (Repository -> Service -> Handler)
	taskRepo := task.NewMemoryRepository()
	taskService := task.NewService(taskRepo)
	taskHandler := task.NewHandler(taskService)

	// 3. Setup Router (Go 1.22+ Standard ServeMux)
	mux := http.NewServeMux()

	// Register Routes
	health.RegisterRoutes(mux)
	taskHandler.RegisterRoutes(mux)

	// 4. Wrap with Middleware Pipeline (Recovery -> Logging -> Router)
	handlerWithMiddleware := middleware.Chain(
		mux,
		middleware.Recovery(logger),
		middleware.Logging(logger),
	)

	// 5. Configure HTTP Server with Production Timeouts
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	serverAddr := fmt.Sprintf(":%s", port)

	server := &http.Server{
		Addr:              serverAddr,
		Handler:           handlerWithMiddleware,
		ReadHeaderTimeout: 5 * time.Second,  // Protects against Slowloris attacks
		ReadTimeout:       10 * time.Second, // Max time reading entire request including body
		WriteTimeout:      15 * time.Second, // Max time writing response
		IdleTimeout:       60 * time.Second, // Max time keep-alive connection remains open
	}

	// 6. Channel to listen for OS termination signals (SIGINT, SIGTERM)
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// 7. Start Server in Background Goroutine
	serverErrors := make(chan error, 1)
	go func() {
		logger.Info("server listening and serving", "addr", serverAddr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	// 8. Block until shutdown signal or server error
	select {
	case err := <-serverErrors:
		logger.Error("server startup error", "error", err)
		os.Exit(1)

	case sig := <-shutdownSignal:
		logger.Info("shutdown signal received, initiating graceful shutdown...", "signal", sig.String())

		// Allocate 10-second timeout window for in-flight requests to complete
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Error("server shutdown forced with error", "error", err)
			_ = server.Close()
			os.Exit(1)
		}

		logger.Info("server exited cleanly. All in-flight requests completed.")
	}
}


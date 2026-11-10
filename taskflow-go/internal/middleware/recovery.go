package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"taskflow/internal/response"
)

// Recovery catches unhandled panics inside downstream handlers, logs stack traces internally,
// and returns a safe 500 Internal Server Error without leaking internal implementation details.
func Recovery(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					stack := string(debug.Stack())
					logger.Error("panic recovered in http handler",
						"panic", fmt.Sprintf("%v", rec),
						"method", r.Method,
						"path", r.URL.Path,
						"stack", stack,
					)

					response.Error(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "An unexpected error occurred")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}


package middleware

import "net/http"

// Middleware defines the standard Go middleware function signature.
//
// Spring Boot Comparison:
// In Spring, cross-cutting concerns are implemented via `jakarta.servlet.Filter` or
// `HandlerInterceptor` registered in `WebMvcConfigurer`. Spring executes them sequentially
// through reflection and servlet pipeline dispatch.
// In Go, middleware is simply higher-order functions: `func(http.Handler) http.Handler`.
type Middleware func(http.Handler) http.Handler

// Chain composes multiple middlewares into a single handler.
// Middlewares execute in the order they are passed:
//
//	Chain(handler, m1, m2)
//
// Execution Order:
//
//	Request  --> m1 (before) --> m2 (before) --> handler
//	Response <-- m1 (after)  <-- m2 (after)  <-- handler
func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}


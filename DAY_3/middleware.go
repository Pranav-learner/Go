package day3
Middleware is code that executes around your HTTP handlers.

You can use it for:
- Authentication.
- Authorization.
- Logging.
- Request tracing.
- Rate limiting.
- CORS.
- Panic recovery.
- Request timing.
1. Understanding middleware
Imagine every request needs to be logged.

Instead of adding logging code to every handler, create one middleware function.

The standard signature is:
func middleware(
    next http.Handler,
) http.Handler


2. Creating logging middleware
func loggingMiddleware(
    next http.Handler,
) http.Handler {

    return http.HandlerFunc(
        func(
            w http.ResponseWriter,
            r *http.Request,
        ) {

            start := time.Now()

            log.Println(
                "Request:",
                r.Method,
                r.URL.Path,
            )

            next.ServeHTTP(w, r)

            log.Println(
                "Duration:",
                time.Since(start),
            )
        },
    )
}


3. Applying middleware
mux := http.NewServeMux()

mux.HandleFunc(
    "GET /users",
    getUsers,
)

handler := loggingMiddleware(mux)

http.ListenAndServe(":8080", handler)


Now every request passes through the logging middleware.

// Production middleware stack
// A typical backend might use:
// Incoming Request
//        |
//        v
// Request ID Middleware
//        |
//        v
// Recovery Middleware
//        |
//        v
// Logging Middleware
//        |
//        v
// CORS Middleware
//        |
//        v
// Rate Limiting Middleware
//        |
//        v
// Authentication Middleware
//        |
//        v
// Authorization
//        |
//        v
// Handler

//Middleware chaining
// You can combine multiple middleware layers.
// handler := loggingMiddleware(
//     authMiddleware(
//         mux,
//     ),
// )


 //VALIDATION

//  Using validation libraries
// For more complex validation, Go developers often use:
// go-playground/validator

 //Validation checklist
// For production APIs, validate:
// - Required fields.
// - String length.
// - Email format.
// - Numeric ranges.
// - Enum values.
// - Date formats.
// - Nested objects.
// - Array sizes.
// - File upload sizes.
// - Business rules.
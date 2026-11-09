package day3
// Typical flow:
// 1. Client sends an HTTP request.
// 2. HTTP server receives it.
// 3. Router identifies the matching route.
// 4. Middleware executes.
// 5. Handler processes the request.
// 6. Service executes business logic.
// 7. Repository interacts with the database.
// 8. Handler writes the response.
// 9. Client receives the response.
// 1. Understanding handler functions
// A handler is simply a function that implements HTTP request processing.
func getUsers(
    w http.ResponseWriter,
    r *http.Request,
) {
    fmt.Fprintln(w, "Fetching users")
}


// Go also has an http.Handler interface:
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}


// Any type implementing ServeHTTP satisfies this interface.

// Example:
type UserHandler struct{}

func (h *UserHandler) ServeHTTP(
    w http.ResponseWriter,
    r *http.Request,
) {
    fmt.Fprintln(w, "User handler")
}


// Register it:
mux := http.NewServeMux()

mux.Handle("/users", &UserHandler{})


// 2. Handler vs HandlerFunc
// These are two different ways to register handlers.
mux.HandleFunc("/users", getUsers)


And:
mux.Handle("/users", http.HandlerFunc(getUsers))


// http.HandlerFunc converts a function with the appropriate signature into an http.Handler.
// 3. Separating handler logic
// Avoid putting all business logic inside your HTTP handler.

Bad structure:
func createUser(
    w http.ResponseWriter,
    r *http.Request,
) {
    // Validate request
    // Connect to database
    // Insert user
    // Send email
    // Generate JWT
    // Return response
}


// This becomes difficult to test and maintain.

Better:
func createUser(
    w http.ResponseWriter,
    r *http.Request,
) {

    // Parse request

    // Validate request

    // Call service

    // Return response
}


// Your handler should primarily coordinate HTTP input and output.

// Production principle: Keep transport logic separate from business logic.
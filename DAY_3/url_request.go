package day3
import (
	"net/http"
	"fmt"

)
// Go's built-in ServeMux
// Go provides http.ServeMux as its standard HTTP request multiplexer.
mux := http.NewServeMux()

mux.HandleFunc("/users", usersHandler)

mux.HandleFunc("/products", productsHandler)

http.ListenAndServe(":8080", mux)


// A multiplexer matches incoming requests to registered handlers.
// 2. Method-aware routing
// Modern Go's ServeMux supports HTTP method patterns.
mux := http.NewServeMux()

mux.HandleFunc(
    "GET /users",
    getUsers,
)

mux.HandleFunc(
    "POST /users",
    createUser,
)

mux.HandleFunc(
    "DELETE /users/{id}",
    deleteUser,
)


// This is a very useful feature.

// You don't need to write a switch statement for every HTTP method.
3. Path parameters
Suppose you want:
GET /users/100


You can access the path parameter using PathValue().
func getUser(
    w http.ResponseWriter,
    r *http.Request,
) {

    id := r.PathValue("id")

    fmt.Fprintln(w, "User ID:", id)
}


Register:
mux.HandleFunc(
    "GET /users/{id}",
    getUser,
)


// 4. Query parameters
// Consider:
GET /users?page=2&limit=10


Access parameters:
func getUsers(
    w http.ResponseWriter,
    r *http.Request,
) {

    page := r.URL.Query().Get("page")

    limit := r.URL.Query().Get("limit")

    fmt.Fprintln(w, page, limit)
}


// Query parameters are strings, so validate and convert them before using them.
// 5. Route conflict and specificity
// Go's ServeMux uses pattern specificity rules to match requests.

// For example:
mux.HandleFunc(
    "GET /users/me",
    currentUser,
)

mux.HandleFunc(
    "GET /users/{id}",
    getUser,
)


// /users/me is more specific than /users/{id}.

// If two patterns conflict without one being more specific, registration can panic.
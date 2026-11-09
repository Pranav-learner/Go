package day3
Struct to JSON — Marshalling
Suppose you have a user struct.
type User struct {
    ID    int
    Name  string
    Email string
}


Convert it to JSON:
user := User{
    ID:    1,
    Name:  "Alex",
    Email: "alex@example.com",
}

data, err := json.Marshal(user)

if err != nil {
    fmt.Println(err)
    return
}

fmt.Println(string(data))


Output:
{"ID":1,"Name":"Alex","Email":"alex@example.com"}


2. JSON tags
In production, you usually want lowercase JSON field names.

Use struct tags:
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}


Now JSON becomes:
{
    "id": 1,
    "name": "Alex",
    "email": "alex@example.com"
}

JSON to struct — Unmarshalling
jsonData := `{
    "id": 1,
    "name": "Alex",
    "email": "alex@example.com"
}`

var user User

err := json.Unmarshal(
    []byte(jsonData),
    &user,
)

if err != nil {
    fmt.Println(err)
    return
}

fmt.Println(user.Name)


Notice that we pass &user.

The decoder needs a pointer so it can populate the struct.
4. Decoding JSON from HTTP requests
This is one of the most common operations in backend development.
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

func createUser(
    w http.ResponseWriter,
    r *http.Request,
) {

    var req CreateUserRequest

    err := json.NewDecoder(r.Body).Decode(&req)

    if err != nil {
        http.Error(
            w,
            "Invalid JSON",
            http.StatusBadRequest,
        )
        return
    }

    fmt.Println(req.Name)
    fmt.Println(req.Email)

    w.WriteHeader(http.StatusCreated)
}


5. Returning JSON responses
func getUser(
    w http.ResponseWriter,
    r *http.Request,
) {

    user := User{
        ID:    1,
        Name:  "Alex",
        Email: "alex@example.com",
    }

    w.Header().Set(
        "Content-Type",
        "application/json",
    )

    err := json.NewEncoder(w).Encode(user)

    if err != nil {
        return
    }
}


json.NewEncoder(w) writes the JSON directly to the response writer.

// Production takeaway: Understand JSON tags, pointer decoding, validation, body limits, 
// unknown fields and sensitive-field exclusion.

//API design principles
// - Use appropriate HTTP methods.
// - Return correct status codes.
// - Validate inputs.
// - Support pagination and filtering.
// - Use consistent error responses.
// - Avoid exposing database models directly.
// - Design idempotent operations where appropriate.
// - Document your APIs using OpenAPI.
// Spring Boot comparison: Your DTOs, controllers, service layer and REST conventions transfer
//  directly to Go. The difference is that you assemble these components explicitly rather than relying on extensive framework annotations.
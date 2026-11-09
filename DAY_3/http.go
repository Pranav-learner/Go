package main
// What is net/http?
// Go provides a built-in standard library package called net/http.

// It allows you to:
// - Create HTTP servers.
// - Handle incoming requests.
// - Send HTTP responses.
// - Make HTTP client requests.
// - Implement middleware.
// - Manage HTTP connections and timeouts.

import (
	"fmt"
	"net/http"
)

func homehandler(
	w http.ResponseWriter,
	r *http.Request,
){
	fmt.Fprintln(w, "Welcome to Go Backend!")
}

func main(){
	http.HandleFunc("/", homehandler)

    fmt.Println("Server running on :8080")

    http.ListenAndServe(":8080", nil)
}

//For production applications, explicitly configure an http.Server.

import (
    "log"
    "net/http"
    "time"
)

func main() {

    mux := http.NewServeMux()

    mux.HandleFunc(
        "/",
        func(w http.ResponseWriter, r *http.Request) {
            w.Write([]byte("Go server running"))
        },
    )

    server := &http.Server{

        Addr:         ":8080",
        Handler:      mux,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    log.Fatal(server.ListenAndServe())
}


// Why configure these?
// - ReadTimeout: Limits how long the server reads a request.
// - WriteTimeout: Limits how long the server writes a response.
// - IdleTimeout: Limits idle keep-alive connections.
// We'll add graceful shutdown later, which is another essential production requirement.

// Spring Boot comparison:
// Spring Boot	Go
// @RestController	HTTP handlers
// @GetMapping	ServeMux route registration
// @RequestBody	Decode r.Body
// ResponseEntity	ResponseWriter
// Embedded Tomcat	Go HTTP server
// application.properties	Environment/configuration
package main

import (
	"fmt"
    "net/http"
	"taskflow/health"
)

func main() {
    http.HandleFunc("/api/v1/health", health.Handler)

	fmt.Println("Taskflow server running on :8080")

	err := http.ListenAndServe(":8080", nil)

    if err != nil {
        fmt.Println("Server failed:", err)
    }
}

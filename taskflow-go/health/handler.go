package health

import (
    "encoding/json"
    "net/http"
)

type HealthResponse struct {
    Status  string `json:"status"`
    Service string `json:"service"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
    response := HealthResponse{
        Status:  "ok",
        Service: "taskflow",
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// //Go identifier	Visibility

// healthHandler	Private to its package
// HealthHandler	Exported
// userService	Private to its package
// UserService	Exported
//package main

// func calculateTotal() int {
//     return 100
// }

// func CalculateTotal() int {
//     return 100
// }


// - calculateTotal() can only be accessed inside its package.
// - CalculateTotal() can be accessed from other packages.
package health

import (
	"net/http"
	"taskflow/internal/response"
)

// Response represents the health check status response.
type Response struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

// Handler handles GET /api/v1/health requests.
func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	response.JSON(w, http.StatusOK, Response{
		Status:  "ok",
		Service: "taskflow",
	})
}

// RegisterRoutes registers the health check route onto the provided ServeMux.
func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/health", Handler)
}


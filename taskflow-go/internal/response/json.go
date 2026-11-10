package response

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse defines the standard error envelope across all TaskFlow APIs.
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains machine-readable error codes and human-readable messages.
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// JSON serializes the provided data as JSON with the specified HTTP status code.
func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

// Error writes a standardized ErrorResponse JSON envelope.
func Error(w http.ResponseWriter, status int, code, message string) {
	JSON(w, status, ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}


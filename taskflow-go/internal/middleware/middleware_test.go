package middleware

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"taskflow/internal/response"
)

func TestRecoveryMiddleware(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	// Handler that deliberately panics
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("database connection exploded!")
	})

	handler := Recovery(logger)(panicHandler)

	req := httptest.NewRequest(http.MethodGet, "/panic-route", nil)
	rec := httptest.NewRecorder()

	// Should not crash the test
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 on recovered panic, got %d", rec.Code)
	}

	var errResp response.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to parse error response: %v", err)
	}

	if errResp.Error.Code != "INTERNAL_SERVER_ERROR" {
		t.Errorf("expected error code INTERNAL_SERVER_ERROR, got %s", errResp.Error.Code)
	}

	// Ensure internal panic message is NOT exposed to client
	if errResp.Error.Message == "database connection exploded!" {
		t.Errorf("security leak: exposed internal panic message in HTTP response: %s", errResp.Error.Message)
	}
}

func TestChainMiddleware(t *testing.T) {
	var executionOrder []string

	m1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "m1_start")
			next.ServeHTTP(w, r)
			executionOrder = append(executionOrder, "m1_end")
		})
	}

	m2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "m2_start")
			next.ServeHTTP(w, r)
			executionOrder = append(executionOrder, "m2_end")
		})
	}

	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		executionOrder = append(executionOrder, "handler")
		w.WriteHeader(http.StatusOK)
	})

	chained := Chain(baseHandler, m1, m2)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	chained.ServeHTTP(rec, req)

	expectedOrder := []string{"m1_start", "m2_start", "handler", "m2_end", "m1_end"}
	if len(executionOrder) != len(expectedOrder) {
		t.Fatalf("expected %d steps, got %d", len(expectedOrder), len(executionOrder))
	}

	for i, step := range expectedOrder {
		if executionOrder[i] != step {
			t.Errorf("step %d: expected %s, got %s", i, step, executionOrder[i])
		}
	}
}


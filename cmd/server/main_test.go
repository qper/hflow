package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthzAndMetricsEndpoints(t *testing.T) {
	t.Setenv("HABITFLOW_SKIP_DB_CHECK", "1")
	handler := newAppHandler()

	t.Run("healthz", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", rec.Code)
		}
	})

	t.Run("metrics", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", rec.Code)
		}
		if got := rec.Body.String(); got == "" {
			t.Fatal("expected metrics body")
		}
	})
}

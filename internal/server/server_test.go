package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthzEndpoint(t *testing.T) {
	handler := NewHandler(nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 from /healthz, got %d", rec.Code)
	}
}

func TestReadyzFailsWithoutDatabase(t *testing.T) {
	handler := NewHandler(nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503 from /readyz without database, got %d", rec.Code)
	}
}

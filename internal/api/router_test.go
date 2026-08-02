package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouterRespondsWithNotFoundForUnknownPath(t *testing.T) {
	router := Router{}
	req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

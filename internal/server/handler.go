package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"sync"
)

// Handler exposes health and metrics HTTP endpoints for deployment.
type Handler struct {
	db      *sql.DB
	metrics *Metrics
	mu      sync.RWMutex
	ready   bool
}

func NewHandler(db *sql.DB, metrics *Metrics) *Handler {
	return &Handler{db: db, metrics: metrics}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/healthz":
		h.metricsHTTP(w, http.StatusOK, map[string]string{"status": "ok"})
	case "/readyz":
		if h.db == nil {
			h.metricsHTTP(w, http.StatusServiceUnavailable, map[string]string{"status": "db-unavailable"})
			return
		}
		if err := h.db.Ping(); err != nil {
			h.metricsHTTP(w, http.StatusServiceUnavailable, map[string]string{"status": "db-unavailable"})
			return
		}
		h.metricsHTTP(w, http.StatusOK, map[string]string{"status": "ready"})
	case "/metrics":
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		_, _ = w.Write([]byte(h.metrics.Render()))
	default:
		http.NotFound(w, r)
	}
}

func (h *Handler) metricsHTTP(w http.ResponseWriter, status int, payload map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = fmt.Fprintf(w, `{"status":"%s"}`, payload["status"])
}

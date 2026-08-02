package api

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/qper/hflow/internal/auth"
	"github.com/qper/hflow/internal/habit"
)

// Router wires auth and habit handlers under a shared HTTP surface.
type Router struct {
	authHandler   auth.Handler
	habitService  habit.Service
	authRouter    auth.Router
	errorHandler  ErrorHandler
	accessSecret  string
	refreshSecret string
	handlers      Handlers
}

func NewRouter(db *sql.DB, accessSecret, refreshSecret string, habitService habit.Service) Router {
	return Router{
		authHandler:   auth.NewHandler(db, accessSecret, refreshSecret),
		habitService:  habitService,
		authRouter:    auth.NewRouter(db, accessSecret, refreshSecret),
		errorHandler:  ErrorHandler{},
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
		handlers:      NewHandlers(habitService, accessSecret, refreshSecret),
	}
}

func (r Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if strings.HasPrefix(req.URL.Path, "/api/v1/auth/") || req.URL.Path == "/api/v1/auth" {
		r.authRouter.ServeHTTP(w, req)
		return
	}
	if req.URL.Path == "/api/v1/push/subscriptions" && req.Method == http.MethodPost {
		r.handlers.savePushSubscription(w, req)
		return
	}
	if strings.HasPrefix(req.URL.Path, "/api/v1/habits") || strings.HasPrefix(req.URL.Path, "/api/v1/entries") {
		if req.URL.Path == "/api/v1/habits" && req.Method == http.MethodGet {
			r.handlers.listHabits(w, req)
			return
		}
		if req.URL.Path == "/api/v1/habits" && req.Method == http.MethodPost {
			r.handlers.createHabit(w, req)
			return
		}
		if strings.HasPrefix(req.URL.Path, "/api/v1/habits/") && req.Method == http.MethodGet {
			r.handlers.getHabit(w, req)
			return
		}
		if strings.HasPrefix(req.URL.Path, "/api/v1/habits/") && req.Method == http.MethodPut {
			r.handlers.updateHabit(w, req)
			return
		}
		if strings.HasPrefix(req.URL.Path, "/api/v1/habits/") && req.Method == http.MethodDelete {
			r.handlers.deleteHabit(w, req)
			return
		}
		if req.URL.Path == "/api/v1/entries" && req.Method == http.MethodPost {
			r.handlers.createEntry(w, req)
			return
		}
		if strings.HasPrefix(req.URL.Path, "/api/v1/entries/") && req.Method == http.MethodDelete {
			r.handlers.deleteEntry(w, req)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

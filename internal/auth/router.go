package auth

import (
	"database/sql"
	"net/http"
)

// Router wires the minimal auth endpoints with the existing router approach.
type Router struct {
	handler Handler
}

func NewRouter(db *sql.DB, accessSecret, refreshSecret string) Router {
	return Router{handler: NewHandler(db, accessSecret, refreshSecret)}
}

func (r Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost && req.URL.Path == "/auth/register" {
		r.handler.Register(w, req)
		return
	}
	if req.Method == http.MethodPost && req.URL.Path == "/auth/login" {
		r.handler.Login(w, req)
		return
	}
	if req.Method == http.MethodPost && req.URL.Path == "/auth/refresh" {
		r.handler.Refresh(w, req)
		return
	}
	if req.Method == http.MethodPost && req.URL.Path == "/auth/logout" {
		r.handler.Logout(w, req)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

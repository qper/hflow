package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/qper/hflow/internal/api"
	"github.com/qper/hflow/internal/db"
	"github.com/qper/hflow/internal/habit"
)

func main() {
	handler := newAppHandler()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("listening on :%s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func newAppHandler() http.Handler {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/habitflow?sslmode=disable"
	}
	accessSecret := os.Getenv("ACCESS_SECRET")
	if accessSecret == "" {
		accessSecret = "dev-access-secret"
	}
	refreshSecret := os.Getenv("REFRESH_SECRET")
	if refreshSecret == "" {
		refreshSecret = "dev-refresh-secret"
	}

	conn, err := db.Open(dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	if err := conn.Ping(); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	repo := habit.NewPostgresRepository(conn)
	svc := habit.NewService(repo, repo, repo)
	apiRouter := api.NewRouter(conn, accessSecret, refreshSecret, svc)
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, "ok")
	})
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/", apiRouter)
	return mux
}

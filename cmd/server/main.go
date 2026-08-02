package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/qper/hflow/internal/api"
	"github.com/qper/hflow/internal/auth"
	"github.com/qper/hflow/internal/db"
	"github.com/qper/hflow/internal/habit"
	"github.com/qper/hflow/internal/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/habitflow?sslmode=disable"
	}
	accessSecret := os.Getenv("ACCESS_TOKEN_SECRET")
	if accessSecret == "" {
		accessSecret = "dev-access-secret"
	}
	refreshSecret := os.Getenv("REFRESH_TOKEN_SECRET")
	if refreshSecret == "" {
		refreshSecret = "dev-refresh-secret"
	}

	conn, err := db.Open(dsn)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		log.Printf("database preview check failed: %v", err)
	}

	repo := habit.NewPostgresRepository(conn)
	svc := habit.NewService(repo, repo, repo)
	metrics := server.NewMetrics()
	baseHandler := server.NewHandler(conn, metrics)
	_ = auth.NewRouter(conn, accessSecret, refreshSecret)
	apiRouter := api.NewRouter(conn, accessSecret, refreshSecret, svc)
	mux := http.NewServeMux()
	mux.Handle("/healthz", baseHandler)
	mux.Handle("/readyz", baseHandler)
	mux.Handle("/metrics", baseHandler)
	mux.Handle("/api/", http.StripPrefix("/api", apiRouter))
	mux.Handle("/", baseHandler)

	log.Printf("listening on :%s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux); err != nil {
		log.Fatalf("listen and serve: %v", err)
	}
}

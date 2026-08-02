package db

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Open establishes a database connection using the pgx driver.
// The caller is expected to provide a PostgreSQL DSN via the supplied URL.
func Open(dsn string) (*sql.DB, error) {
	return sql.Open("pgx", dsn)
}

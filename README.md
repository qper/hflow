# HabitFlow

## Repository layout

- `cmd/server/` — application entrypoints and startup wiring.
- `internal/domain/` — core domain models and domain-level rules.
- `migrations/` — SQL migrations for PostgreSQL, with up/down pairs.
- `scripts/` — operational helpers such as database migration entrypoints.
- `.github/workflows/` — CI pipeline for lint, build, and test.

## Development workflow

1. Format Go source with `gofmt -w $(find . -name '*.go' -not -path './vendor/*')`.
2. Run `go test ./...`.
3. Apply migrations via `DATABASE_URL=... ./scripts/migrate.sh up`.

## Notes

- Authentication and HTTP handlers are intentionally left out of this foundation phase.
- Streaks are stored as a materialized table to keep future analytics queries simple and predictable.

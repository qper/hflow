#!/usr/bin/env bash
set -euo pipefail

if [[ $# -lt 1 ]]; then
  echo "Usage: $0 <up|down>" >&2
  exit 1
fi

ACTION="$1"
DB_URL="${DATABASE_URL:-postgres://postgres:postgres@localhost:5432/habitflow?sslmode=disable}"

if [[ "$ACTION" != "up" && "$ACTION" != "down" ]]; then
  echo "Unsupported action: $ACTION" >&2
  exit 1
fi

if ! command -v migrate >/dev/null 2>&1; then
  echo "migrate CLI is required. Install from https://github.com/golang-migrate/migrate" >&2
  exit 1
fi

migrate -database "$DB_URL" -path ./migrations "$ACTION"

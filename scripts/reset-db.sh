#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DB_PATH="${ROOT_DIR}/backend/data/todos.db"

rm -f "${DB_PATH}" "${DB_PATH}-shm" "${DB_PATH}-wal"

(
  cd "${ROOT_DIR}/backend"
  go run ./cmd/server -bootstrap-only
)

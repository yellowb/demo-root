#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

cleanup() {
  if [[ -n "${BACKEND_PID:-}" ]]; then
    kill "${BACKEND_PID}" >/dev/null 2>&1 || true
  fi

  if [[ -n "${FRONTEND_PID:-}" ]]; then
    kill "${FRONTEND_PID}" >/dev/null 2>&1 || true
  fi
}

trap cleanup EXIT INT TERM

(
  cd "${ROOT_DIR}/backend"
  go run ./cmd/server
) &
BACKEND_PID=$!

(
  cd "${ROOT_DIR}/frontend"
  npm run dev
) &
FRONTEND_PID=$!

wait "${BACKEND_PID}" "${FRONTEND_PID}"

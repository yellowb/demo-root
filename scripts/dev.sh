#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
FRONTEND_URL="${FRONTEND_URL:-http://localhost:5173/}"
AUTO_OPEN_BROWSER="${AUTO_OPEN_BROWSER:-1}"
FRONTEND_WAIT_ATTEMPTS="${FRONTEND_WAIT_ATTEMPTS:-30}"
FRONTEND_WAIT_INTERVAL="${FRONTEND_WAIT_INTERVAL:-1}"

cleanup() {
  if [[ -n "${BROWSER_PID:-}" ]]; then
    kill "${BROWSER_PID}" >/dev/null 2>&1 || true
  fi

  if [[ -n "${BACKEND_PID:-}" ]]; then
    kill "${BACKEND_PID}" >/dev/null 2>&1 || true
  fi

  if [[ -n "${FRONTEND_PID:-}" ]]; then
    kill "${FRONTEND_PID}" >/dev/null 2>&1 || true
  fi
}

trap cleanup EXIT INT TERM

detect_browser_open_command() {
  if command -v open >/dev/null 2>&1; then
    echo "open"
    return 0
  fi

  if command -v xdg-open >/dev/null 2>&1; then
    echo "xdg-open"
    return 0
  fi

  return 1
}

wait_for_frontend() {
  local attempt=0

  while (( attempt < FRONTEND_WAIT_ATTEMPTS )); do
    if curl -fsS "${FRONTEND_URL}" >/dev/null 2>&1; then
      return 0
    fi

    sleep "${FRONTEND_WAIT_INTERVAL}"
    attempt=$((attempt + 1))
  done

  return 1
}

open_frontend_when_ready() {
  if [[ "${AUTO_OPEN_BROWSER}" != "1" ]]; then
    return 0
  fi

  local open_command
  if ! open_command="$(detect_browser_open_command)"; then
    return 0
  fi

  (
    if wait_for_frontend; then
      "${open_command}" "${FRONTEND_URL}" >/dev/null 2>&1 || true
    fi
  ) &
  BROWSER_PID=$!
}

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

open_frontend_when_ready

wait "${BACKEND_PID}" "${FRONTEND_PID}"

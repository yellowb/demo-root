#!/usr/bin/env bash

set -euo pipefail

BACKEND_PORT="${BACKEND_PORT:-8080}"
FRONTEND_PORT="${FRONTEND_PORT:-5173}"

stop_port() {
  local label="$1"
  local port="$2"
  local pids

  if ! command -v lsof >/dev/null 2>&1; then
    echo "Missing lsof; cannot inspect port ${port}." >&2
    exit 1
  fi

  pids="$(lsof -tiTCP:"${port}" -sTCP:LISTEN || true)"
  if [[ -z "${pids}" ]]; then
    echo "${label} port ${port} is free."
    return 0
  fi

  echo "Stopping ${label} process(es) on port ${port}: ${pids//$'\n'/ }"
  kill ${pids} >/dev/null 2>&1 || true
  sleep 1

  pids="$(lsof -tiTCP:"${port}" -sTCP:LISTEN || true)"
  if [[ -n "${pids}" ]]; then
    echo "Force stopping ${label} process(es) on port ${port}: ${pids//$'\n'/ }"
    kill -9 ${pids} >/dev/null 2>&1 || true
    sleep 1
  fi

  pids="$(lsof -tiTCP:"${port}" -sTCP:LISTEN || true)"
  if [[ -n "${pids}" ]]; then
    echo "Failed to stop ${label} process(es) on port ${port}: ${pids//$'\n'/ }" >&2
    exit 1
  fi

  echo "${label} port ${port} is free."
}

stop_port "backend" "${BACKEND_PORT}"
stop_port "frontend" "${FRONTEND_PORT}"

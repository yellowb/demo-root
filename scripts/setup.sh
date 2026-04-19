#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

require_command() {
  local name="$1"
  if ! command -v "${name}" >/dev/null 2>&1; then
    echo "Missing required command: ${name}" >&2
    exit 1
  fi
}

require_command node
require_command npm
require_command go

(
  cd "${ROOT_DIR}/frontend"
  if [[ -f package-lock.json ]]; then
    npm ci
  else
    npm install
  fi
)

(
  cd "${ROOT_DIR}/backend"
  GOPROXY="${TODO_GOPROXY:-https://proxy.golang.org,direct}" go mod download
)

#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
REPO_OPENSPEC="${ROOT_DIR}/.bin/openspec"

if [[ -x "${REPO_OPENSPEC}" ]]; then
  OPENSPEC_BIN="${REPO_OPENSPEC}"
elif command -v openspec >/dev/null 2>&1; then
  OPENSPEC_BIN="$(command -v openspec)"
else
  echo "Missing openspec. Run make setup to install the pinned repo-local version." >&2
  exit 1
fi

(
  cd "${ROOT_DIR}"
  OPENSPEC_TELEMETRY=0 "${OPENSPEC_BIN}" validate todo-management --specs
)

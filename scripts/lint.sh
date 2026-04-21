#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
REPO_LINTER="${ROOT_DIR}/.bin/golangci-lint"

if [[ -x "${REPO_LINTER}" ]]; then
  GOLANGCI_LINT="${REPO_LINTER}"
elif command -v golangci-lint >/dev/null 2>&1; then
  GOLANGCI_LINT="$(command -v golangci-lint)"
else
  echo "Missing golangci-lint. Run make setup to install the pinned repo-local version." >&2
  exit 1
fi

(
  cd "${ROOT_DIR}/backend"
  "${GOLANGCI_LINT}" run --config "${ROOT_DIR}/.golangci.yml" ./...
)

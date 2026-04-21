#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
GOLANGCI_LINT_VERSION="$(tr -d '[:space:]' < "${ROOT_DIR}/.golangci-lint-version")"
GOLANGCI_LINT_EXPECTED_VERSION="${GOLANGCI_LINT_VERSION#v}"

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

version_matches() {
  local binary="$1"
  "${binary}" --version 2>/dev/null | grep -Eq "version ${GOLANGCI_LINT_EXPECTED_VERSION}([[:space:]]|$)"
}

install_golangci_lint() {
  local target="${ROOT_DIR}/.bin/golangci-lint"
  mkdir -p "${ROOT_DIR}/.bin"

  if [[ -x "${target}" ]] && version_matches "${target}"; then
    return 0
  fi

  if command -v golangci-lint >/dev/null 2>&1 && version_matches "$(command -v golangci-lint)"; then
    cp "$(command -v golangci-lint)" "${target}"
    chmod +x "${target}"
    return 0
  fi

  require_command curl
  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
    | sh -s -- -b "${ROOT_DIR}/.bin" "${GOLANGCI_LINT_VERSION}"
}

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

install_golangci_lint

#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
GOLANGCI_LINT_VERSION="$(tr -d '[:space:]' < "${ROOT_DIR}/.golangci-lint-version")"
GOLANGCI_LINT_EXPECTED_VERSION="${GOLANGCI_LINT_VERSION#v}"
OPENSPEC_VERSION="$(tr -d '[:space:]' < "${ROOT_DIR}/.openspec-version")"

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

golangci_lint_version_matches() {
  local binary="$1"
  "${binary}" --version 2>/dev/null | grep -Eq "version ${GOLANGCI_LINT_EXPECTED_VERSION}([[:space:]]|$)"
}

openspec_version_matches() {
  local binary="$1"
  [[ -x "${binary}" ]] && [[ "$("${binary}" --version 2>/dev/null | tr -d '[:space:]')" == "${OPENSPEC_VERSION}" ]]
}

install_golangci_lint() {
  local target="${ROOT_DIR}/.bin/golangci-lint"
  mkdir -p "${ROOT_DIR}/.bin"

  if [[ -x "${target}" ]] && golangci_lint_version_matches "${target}"; then
    return 0
  fi

  if command -v golangci-lint >/dev/null 2>&1 && golangci_lint_version_matches "$(command -v golangci-lint)"; then
    cp "$(command -v golangci-lint)" "${target}"
    chmod +x "${target}"
    return 0
  fi

  require_command curl
  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
    | sh -s -- -b "${ROOT_DIR}/.bin" "${GOLANGCI_LINT_VERSION}"
}

install_openspec() {
  local target="${ROOT_DIR}/.bin/openspec"
  local install_dir="${ROOT_DIR}/.bin/openspec-cli"
  mkdir -p "${ROOT_DIR}/.bin"

  if openspec_version_matches "${target}"; then
    return 0
  fi

  if command -v openspec >/dev/null 2>&1 && openspec_version_matches "$(command -v openspec)"; then
    local system_openspec
    system_openspec="$(command -v openspec)"
    rm -f "${target}"
    cat >"${target}" <<EOF
#!/usr/bin/env bash
exec "${system_openspec}" "\$@"
EOF
    chmod +x "${target}"
    return 0
  fi

  rm -rf "${install_dir}"
  mkdir -p "${install_dir}"
  npm install --prefix "${install_dir}" --no-save "@fission-ai/openspec@${OPENSPEC_VERSION}"
  rm -f "${target}"
  ln -sf "${install_dir}/node_modules/.bin/openspec" "${target}"

  if ! openspec_version_matches "${target}"; then
    echo "Installed openspec does not match expected version ${OPENSPEC_VERSION}" >&2
    exit 1
  fi
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
install_openspec

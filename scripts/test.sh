#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

(
  cd "${ROOT_DIR}"
  python3 -m json.tool .codex/hooks.json >/dev/null
  bash -n .codex/hooks/stop_make_test.sh
  bash ./scripts/validate-specs.sh
)

(
  cd "${ROOT_DIR}"
  bash ./scripts/dev_test.sh
)

(
  cd "${ROOT_DIR}"
  bash ./scripts/lint.sh
)

(
  cd "${ROOT_DIR}/backend"
  go test ./...
)

(
  cd "${ROOT_DIR}/frontend"
  npm run typecheck
)

(
  cd "${ROOT_DIR}/frontend"
  npm run build
)

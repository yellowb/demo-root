#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "${TMP_DIR}"' EXIT

MOCK_BIN_DIR="${TMP_DIR}/bin"
LOG_DIR="${TMP_DIR}/logs"
mkdir -p "${MOCK_BIN_DIR}" "${LOG_DIR}"

cat > "${MOCK_BIN_DIR}/go" <<EOF
#!/usr/bin/env bash
sleep 1
EOF

cat > "${MOCK_BIN_DIR}/npm" <<EOF
#!/usr/bin/env bash
sleep 1
EOF

cat > "${MOCK_BIN_DIR}/curl" <<EOF
#!/usr/bin/env bash
exit 0
EOF

cat > "${MOCK_BIN_DIR}/open" <<EOF
#!/usr/bin/env bash
echo "\$*" > "${LOG_DIR}/open.log"
EOF

chmod +x "${MOCK_BIN_DIR}/go" "${MOCK_BIN_DIR}/npm" "${MOCK_BIN_DIR}/curl" "${MOCK_BIN_DIR}/open"

(
  cd "${ROOT_DIR}"
  PATH="${MOCK_BIN_DIR}:${PATH}" FRONTEND_WAIT_INTERVAL=0.1 ./scripts/dev.sh
)

if [[ ! -f "${LOG_DIR}/open.log" ]]; then
  echo "expected dev.sh to open the browser after the frontend became reachable" >&2
  exit 1
fi

if [[ "$(cat "${LOG_DIR}/open.log")" != "http://localhost:5173/" ]]; then
  echo "expected browser to open http://localhost:5173/, got: $(cat "${LOG_DIR}/open.log")" >&2
  exit 1
fi

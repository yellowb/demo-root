#!/usr/bin/env bash

set -euo pipefail

cat >/dev/null

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/../.." && pwd)"
EXPECTED_ROOT="/Users/yellowb/ppt/demo-root"
STATE_DIR="${TMPDIR:-/tmp}/codex-demo-root-hooks"
PASSED_FILE="${STATE_DIR}/last_make_test.sha"
LOG_FILE="${STATE_DIR}/make-test-last.log"

if [[ "${ROOT_DIR}" != "${EXPECTED_ROOT}" ]]; then
  exit 0
fi

mkdir -p "${STATE_DIR}"

if ! git -C "${ROOT_DIR}" rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  exit 0
fi

if [[ -z "$(git -C "${ROOT_DIR}" status --porcelain=v1 --untracked-files=all)" ]]; then
  exit 0
fi

fingerprint() {
  python3 - "${ROOT_DIR}" <<'PY'
import hashlib
import os
import subprocess
import sys

root = sys.argv[1]

def git(args):
    return subprocess.check_output(["git", "-C", root, *args])

digest = hashlib.sha256()
digest.update(git(["status", "--porcelain=v1", "-z", "--untracked-files=all"]))
digest.update(git(["diff", "--binary", "--no-ext-diff"]))
digest.update(git(["diff", "--cached", "--binary", "--no-ext-diff"]))

untracked = git(["ls-files", "--others", "--exclude-standard", "-z"]).split(b"\0")
for raw_path in sorted(path for path in untracked if path):
    rel_path = raw_path.decode("utf-8", "surrogateescape")
    abs_path = os.path.join(root, rel_path)
    digest.update(b"untracked\0")
    digest.update(raw_path)
    digest.update(b"\0")
    if os.path.isfile(abs_path):
        with open(abs_path, "rb") as handle:
            for chunk in iter(lambda: handle.read(1024 * 1024), b""):
                digest.update(chunk)

print(digest.hexdigest())
PY
}

emit_block() {
  local reason="$1"
  python3 -c 'import json, sys; print(json.dumps({"decision": "block", "reason": sys.stdin.read()}))' <<<"${reason}"
}

current_fingerprint="$(fingerprint)"

if [[ -f "${PASSED_FILE}" ]] && [[ "$(cat "${PASSED_FILE}")" == "${current_fingerprint}" ]]; then
  exit 0
fi

if (
  cd "${ROOT_DIR}"
  make test
) >"${LOG_FILE}" 2>&1; then
  printf '%s\n' "${current_fingerprint}" >"${PASSED_FILE}"
  emit_block "The Codex Stop hook ran \`make test\` successfully for the current repository changes.

Continue with a concise final response that says \`make test\` passed via the Codex Stop hook and summarizes the important changes."
  exit 0
fi

tail_output="$(tail -n 80 "${LOG_FILE}" || true)"
emit_block "The Codex Stop hook ran \`make test\`, but it failed.

Read the failure output below, fix the repository, and then finish again. Do not claim completion until \`make test\` passes.

Last 80 lines of \`make test\` output:

${tail_output}"

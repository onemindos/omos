#!/usr/bin/env bash
# End-to-end test for the ROS 2 bridge (FEATURES H3): disposable nats-server, the policy
# service built from omos PR #554, the scripted Zenoh robot and this bridge.
#
# Environment:
#   NATS_SERVER   nats-server binary (default: nats-server on PATH)
#   PYTHON        python with requirements-dev.txt installed (default: python3)
#   POLICY_REF    git ref holding services/policy (default: origin/feat/policy-service)
#   POLICY_BIN / POLICY_SRC  use an existing policy binary and source tree instead of building
#   E2E_EVIDENCE_DIR  keep logs, transcript.jsonl and e2e-report.json here (default: temp dir)
set -euo pipefail

here="$(cd "$(dirname "$0")/../.." && pwd)"
repo="$(git -C "$here" rev-parse --show-toplevel)"
python="${PYTHON:-python3}"
NATS_SERVER="${NATS_SERVER:-$(command -v nats-server || true)}"
if [[ -z "$NATS_SERVER" ]]; then echo "nats-server not found; set NATS_SERVER" >&2; exit 2; fi
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

if [[ -z "${POLICY_BIN:-}" || -z "${POLICY_SRC:-}" ]]; then
  ref="${POLICY_REF:-origin/feat/policy-service}"
  POLICY_SRC="$tmp/policy-src"
  mkdir -p "$POLICY_SRC"
  git -C "$repo" archive "$ref" services/policy services/schemas platform/security/openfga | tar -x -C "$POLICY_SRC"
  (cd "$POLICY_SRC/services/policy" && go build -o "$tmp/policy" ./cmd/policy)
  POLICY_BIN="$tmp/policy"
fi
export NATS_SERVER POLICY_BIN POLICY_SRC
cd "$here"
"$python" -m pytest -v -p no:cacheprovider tests/e2e/test_e2e.py

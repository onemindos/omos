#!/usr/bin/env bash
# Checks and unit-tests the omos-ingest alert rules with promtool.
# Usage: test/alerts/run.sh   (needs promtool and python3 with PyYAML)
set -euo pipefail
HERE="$(cd "$(dirname "$0")" && pwd)"
RULE="$(cd "$HERE/../../../../.." && pwd)/platform/ingest/prometheusrule.yaml"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT
python3 - "$RULE" "$TMP/rules.yaml" <<'PY'
import sys, yaml
doc = yaml.safe_load(open(sys.argv[1]))
yaml.safe_dump(doc["spec"], open(sys.argv[2], "w"), sort_keys=False)
PY
cp "$HERE/ingest-rules.test.yaml" "$TMP/"
promtool check rules "$TMP/rules.yaml"
promtool test rules "$TMP/ingest-rules.test.yaml"

#!/bin/sh
# Apply mcp_servers from config.defaults.yaml to the live onemind-hermes /opt/data/config.yaml, idempotently.
#
#   mcp/apply.sh [--dry-run]    set the live mcp_servers (ops, data, tak on the tool gateway), remove the old
#                               LiteLLM `onemind` entry, keep every other key; backs up config.yaml first
#
# config.defaults.yaml only seeds an empty volume, so after the first start the live file must be changed in place.
# Run it after the deployment carries ONEMIND_GATEWAY_KEY. Hermes reads mcp_servers at start, so restart it
# afterwards (kubectl -n hermes rollout restart deploy/onemind-hermes; the config stays on the PVC).
# Needs kubectl access to the DOKS cluster. KUBE_CONTEXT overrides the default context.
set -eu

CTX="${KUBE_CONTEXT:-do-nyc3-omos-cloud}"
HERE="$(cd "$(dirname "$0")" && pwd)"

kx() {
  kubectl --context "$CTX" -n hermes exec "$@"
}

# Stage the applier and the defaults in the pod's /tmp, then run it as the hermes user (same venv path
# handling as cron/apply.sh). Nothing secret is staged: the defaults only reference ${ONEMIND_GATEWAY_KEY}.
STAGE=$(mktemp -d)
trap 'rm -rf "$STAGE"' EXIT
cp "$HERE/apply_mcp.py" "$STAGE/apply_mcp.py"
cp "$HERE/../config.defaults.yaml" "$STAGE/config.defaults.yaml"
COPYFILE_DISABLE=1   # macOS tar would otherwise add ._* AppleDouble files
export COPYFILE_DISABLE
tar cf - -C "$STAGE" apply_mcp.py config.defaults.yaml | kx -i deploy/onemind-hermes -c hermes -- sh -c \
  'rm -rf /tmp/mcp-apply && mkdir -p /tmp/mcp-apply && tar xf - -C /tmp/mcp-apply && chmod -R a+rX /tmp/mcp-apply'

# shellcheck disable=SC2016 # the script is meant to expand inside the pod, not here
kx deploy/onemind-hermes -c hermes -- sh -c \
  'SITE=$(ls -d /opt/hermes/.venv/lib/python3*/site-packages | head -n 1)
   runuser -u hermes -- env HERMES_HOME=/opt/data HOME=/opt/data/home PYTHONPATH="/opt/hermes:$SITE" \
     /opt/hermes/.venv/bin/python3 /tmp/mcp-apply/apply_mcp.py "$@"; rc=$?
   rm -rf /tmp/mcp-apply; exit $rc' sh "$@"

#!/bin/sh
# Fix the live onemind-hermes volume for Hermes Desktop, idempotently, and check what Desktop depends on.
#
#   desktop/apply.sh [--dry-run]   give /opt/data/onemind-skills to the agent user and remove AppleDouble files
#                                  (backs the removed files up under /opt/data/backups/desktop-<stamp>/)
#   desktop/apply.sh --check       read-only: call the dashboard routes Desktop uses in-process and report
#                                  STT/TTS/image providers (status codes and key names only, never values)
#
# A second plain run prints `unchanged`. Needs kubectl access to the DOKS cluster. KUBE_CONTEXT overrides the
# default context. Nothing here edits config.yaml; the Desktop findings and the settings that stay off on purpose
# are in ../README.md ("Hermes Desktop: feature status").
set -eu

CTX="${KUBE_CONTEXT:-do-nyc3-omos-cloud}"
HERE="$(cd "$(dirname "$0")" && pwd)"

kx() {
  kubectl --context "$CTX" -n hermes exec "$@"
}

STAGE="$(mktemp -d)"
trap 'rm -rf "$STAGE"' EXIT
cp "$HERE/apply_desktop.py" "$HERE/check_desktop.py" "$STAGE/"
COPYFILE_DISABLE=1   # macOS tar would otherwise add ._* AppleDouble files
export COPYFILE_DISABLE
tar cf - -C "$STAGE" apply_desktop.py check_desktop.py | kx -i deploy/onemind-hermes -c hermes -- sh -c \
  'rm -rf /tmp/desktop-apply && mkdir -p /tmp/desktop-apply && tar xf - -C /tmp/desktop-apply && chmod -R a+rX /tmp/desktop-apply'

if [ "${1:-}" = "--check" ]; then
  # shellcheck disable=SC2016 # the script is meant to expand inside the pod, not here
  kx deploy/onemind-hermes -c hermes -- sh -c \
    'SITE=$(ls -d /opt/hermes/.venv/lib/python3*/site-packages | head -n 1)
     runuser -u hermes -- env HERMES_HOME=/opt/data HOME=/opt/data/home PYTHONPATH="/opt/hermes:$SITE" \
       /opt/hermes/.venv/bin/python3 /tmp/desktop-apply/check_desktop.py 2>&1 | grep -v -e "tried to register" -e "MATRIX_ACCESS_TOKEN"
     rm -rf /tmp/desktop-apply; exit 0'
else
  # Runs as root (ownership changes); the script only uses the standard library.
  # shellcheck disable=SC2016
  kx deploy/onemind-hermes -c hermes -- sh -c \
    'HERMES_HOME=/opt/data /opt/hermes/.venv/bin/python3 /tmp/desktop-apply/apply_desktop.py "$@"; rc=$?
     rm -rf /tmp/desktop-apply; exit $rc' sh "$@"
fi

#!/bin/sh
# Install Legacy's doctrine (SOUL.md) and worker roles into the live onemind-hermes instance, idempotently.
#
#   roles/apply.sh [--dry-run]
#
# config.defaults.yaml only seeds an EMPTY volume, so this is how an existing /opt/data gets these. It pushes
# roles/ and the Git SOUL.md from this checkout into the pod (/opt/data/roles-defs), then runs apply_roles.py
# there as the hermes user. Everything it overwrites is backed up under /opt/data/backups/legacy-roles-<stamp>/.
# Needs kubectl access to the DOKS cluster. KUBE_CONTEXT overrides the default context.
set -eu

CTX="${KUBE_CONTEXT:-do-nyc3-omos-cloud}"
HERE="$(cd "$(dirname "$0")" && pwd)"
SOUL="$HERE/../../seed/onemind-hermes/SOUL.md"

kx() {
  kubectl --context "$CTX" -n hermes exec "$@"
}

STAGE="$(mktemp -d)"
trap 'rm -rf "$STAGE"' EXIT
cp -R "$HERE/roles.yaml" "$HERE/plugin" "$HERE/apply_roles.py" "$STAGE/"
cp "$SOUL" "$STAGE/SOUL.md"
find "$STAGE" -name '__pycache__' -prune -exec rm -rf {} +

COPYFILE_DISABLE=1   # macOS tar would otherwise add ._* AppleDouble files
export COPYFILE_DISABLE
tar cf - -C "$STAGE" . | kx -i deploy/onemind-hermes -c hermes -- sh -c \
  'rm -rf /opt/data/roles-defs && mkdir -p /opt/data/roles-defs && tar xf - -C /opt/data/roles-defs 2>/dev/null; chown -R hermes:hermes /opt/data/roles-defs'

# Same environment as cron/apply.sh: Hermes' dependencies live in the venv's site-packages.
# shellcheck disable=SC2016 # the script is meant to expand inside the pod, not here
kx deploy/onemind-hermes -c hermes -- sh -c \
  'SITE=$(ls -d /opt/hermes/.venv/lib/python3*/site-packages | head -n 1)
   runuser -u hermes -- env HERMES_HOME=/opt/data HOME=/opt/data/home PYTHONPATH="/opt/hermes:$SITE" \
     /opt/hermes/.venv/bin/python3 /opt/data/roles-defs/apply_roles.py "$@"' sh "$@"

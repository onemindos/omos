#!/bin/sh
# Apply cron/jobs.yaml to the live onemind-hermes instance, idempotently.
#
#   cron/apply.sh [--dry-run]      run apply_jobs.py inside the pod as the hermes user
#   cron/apply.sh --install-local  first copy scripts/jobs.yaml/apply_jobs.py from this checkout into the pod
#                                  (only for testing before the change is merged and synced; after the merge the
#                                  init container refreshes them from the onemind-hermes-cron ConfigMap)
#
# Needs kubectl access to the DOKS cluster. KUBE_CONTEXT overrides the default context.
set -eu

CTX="${KUBE_CONTEXT:-do-nyc3-omos-cloud}"
HERE="$(cd "$(dirname "$0")" && pwd)"

kx() {
  kubectl --context "$CTX" -n hermes exec "$@"
}

if [ "${1:-}" = "--install-local" ]; then
  shift
  COPYFILE_DISABLE=1   # macOS tar would otherwise add ._* AppleDouble files
  export COPYFILE_DISABLE
  tar cf - -C "$HERE/scripts" . | kx -i deploy/onemind-hermes -c hermes -- sh -c \
    'mkdir -p /opt/data/scripts && tar xf - -C /opt/data/scripts 2>/dev/null; chown -R hermes:hermes /opt/data/scripts'
  tar cf - -C "$HERE" jobs.yaml apply_jobs.py | kx -i deploy/onemind-hermes -c hermes -- sh -c \
    'mkdir -p /opt/data/cron-defs && tar xf - -C /opt/data/cron-defs 2>/dev/null; chown -R hermes:hermes /opt/data/cron-defs'
fi

# Hermes' dependencies (PyYAML, cron API) live in the venv's site-packages, which only the `hermes`
# launcher and the cron runner put on the path; do the same here. Run as hermes so jobs.json keeps its owner.
# shellcheck disable=SC2016 # the script is meant to expand inside the pod, not here
kx deploy/onemind-hermes -c hermes -- sh -c \
  'SITE=$(ls -d /opt/hermes/.venv/lib/python3*/site-packages | head -n 1)
   runuser -u hermes -- env HERMES_HOME=/opt/data HOME=/opt/data/home PYTHONPATH="/opt/hermes:$SITE" \
     /opt/hermes/.venv/bin/python3 /opt/data/cron-defs/apply_jobs.py "$@"' sh "$@"

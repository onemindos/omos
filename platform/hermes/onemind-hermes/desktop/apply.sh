#!/bin/sh
# Fix the live onemind-hermes volume for Hermes Desktop, idempotently, and check what Desktop depends on.
#
#   desktop/apply.sh [--dry-run] [--restart]
#                                  step 1: give /opt/data/onemind-skills to the agent user and remove AppleDouble files
#                                          (backs the removed files up under /opt/data/backups/desktop-<stamp>/)
#                                  step 2: set the `webapp:` block in config.yaml (the browser Desktop, Hermes v0.23.0)
#                                          and back config.yaml up first. Takes effect at the next gateway start.
#   desktop/apply.sh --webapp-off [--restart]
#                                  step 2 only, the other way: webapp.enabled false (rollback switch)
#   desktop/apply.sh --restart     afterwards restart ONLY the gateway process in the pod (s6), not the pod
#   desktop/apply.sh --check       read-only: call the dashboard routes Desktop uses in-process and report
#                                  STT/TTS/image providers and the Webapp listener (status codes and key names only,
#                                  never values)
#
# A second plain run prints `unchanged`. Needs kubectl access to the DOKS cluster. KUBE_CONTEXT overrides the
# default context. Step 2 is the only thing here that edits config.yaml, and only the `webapp` block; the Desktop
# findings and the settings that stay off on purpose are in ../README.md ("Hermes Desktop: feature status").
set -eu

CTX="${KUBE_CONTEXT:-do-nyc3-omos-cloud}"
HERE="$(cd "$(dirname "$0")" && pwd)"

kx() {
  kubectl --context "$CTX" -n hermes exec "$@"
}

RESTART=0
CHECK=0
ARGS=""
for a in "$@"; do
  case "$a" in
    --restart) RESTART=1 ;;
    --check) CHECK=1 ;;
    *) ARGS="$ARGS $a" ;;
  esac
done

STAGE="$(mktemp -d)"
trap 'rm -rf "$STAGE"' EXIT
cp "$HERE/apply_desktop.py" "$HERE/check_desktop.py" "$STAGE/"
COPYFILE_DISABLE=1   # macOS tar would otherwise add ._* AppleDouble files
export COPYFILE_DISABLE
tar cf - -C "$STAGE" apply_desktop.py check_desktop.py | kx -i deploy/onemind-hermes -c hermes -- sh -c \
  'rm -rf /tmp/desktop-apply && mkdir -p /tmp/desktop-apply && tar xf - -C /tmp/desktop-apply && chmod -R a+rX /tmp/desktop-apply'

if [ "$CHECK" = 1 ]; then
  # shellcheck disable=SC2016 # the script is meant to expand inside the pod, not here
  kx deploy/onemind-hermes -c hermes -- sh -c \
    'SITE=$(ls -d /opt/hermes/.venv/lib/python3*/site-packages | head -n 1)
     runuser -u hermes -- env HERMES_HOME=/opt/data HOME=/opt/data/home PYTHONPATH="/opt/hermes:$SITE" \
       /opt/hermes/.venv/bin/python3 /tmp/desktop-apply/check_desktop.py 2>&1 | grep -v -e "tried to register" -e "MATRIX_ACCESS_TOKEN"
     rm -rf /tmp/desktop-apply; exit 0'
else
  # $ARGS is deliberately unquoted: it is a list of flags.
  # Step 1 runs as root (ownership changes); it only uses the standard library. Skipped for --webapp-off.
  # shellcheck disable=SC2016,SC2086
  case "$ARGS" in
    *--webapp-off*) ;;
    *)
      kx deploy/onemind-hermes -c hermes -- sh -c \
        'HERMES_HOME=/opt/data /opt/hermes/.venv/bin/python3 /tmp/desktop-apply/apply_desktop.py --step files "$@"' sh $ARGS
      ;;
  esac
  # Step 2 edits config.yaml, so it runs as hermes (file owner) inside Hermes' environment, like the other appliers.
  # shellcheck disable=SC2016,SC2086
  kx deploy/onemind-hermes -c hermes -- sh -c \
    'SITE=$(ls -d /opt/hermes/.venv/lib/python3*/site-packages | head -n 1)
     runuser -u hermes -- env HERMES_HOME=/opt/data HOME=/opt/data/home PYTHONPATH="/opt/hermes:$SITE" \
       /opt/hermes/.venv/bin/python3 /tmp/desktop-apply/apply_desktop.py --step webapp "$@"; rc=$?
     rm -rf /tmp/desktop-apply; exit $rc' sh $ARGS
  if [ "$RESTART" = 1 ]; then
    case "$ARGS" in
      *--dry-run*) echo "--restart ignored with --dry-run" ;;
      *)
        echo "Restarting the gateway process (s6 gateway-default); the pod stays up."
        # Hermes v0.23.0 keeps the gateway slots in /run/hermes-service (v0.22.x: /run/service).
        kx deploy/onemind-hermes -c hermes -- /command/s6-svc -r /run/hermes-service/gateway-default
        ;;
    esac
  fi
fi

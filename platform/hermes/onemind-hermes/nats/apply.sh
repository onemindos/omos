#!/bin/sh
# Install Legacy's NATS presence into the live onemind-hermes instance, idempotently.
#
#   nats/apply.sh [--dry-run] [--status] [--with-context] [--disable] [--restart]
#
#   (no flags)      copy the nats-platform plugin from this checkout into /opt/data/plugins and list it in
#                   plugins.enabled. Changes nothing until the gateway is restarted.
#   --with-context  also install and enable nats-context (the live operational picture and generic NATS tools).
#                   Off by default: see MIGRATION.md "NATS cutover".
#   --disable       remove both plugins from plugins.enabled (rollback; files stay)
#   --status        report only: plugin files, enabled list, which NATS variables the pod has (names only)
#   --dry-run       show what would change
#   --restart       afterwards restart ONLY the gateway process in the pod (s6), not the pod
#
# The connection itself (URL, credentials, CA, identity, refuse-duplicate) is in deployment.yaml, so a
# Git change plus Argo moves it. config.defaults.yaml only seeds an empty volume, so this is how an existing
# /opt/data gets the plugins. Everything it overwrites is backed up under
# /opt/data/backups/legacy-nats-<stamp>/. No credential is read or printed.
# Needs kubectl access to the DOKS cluster. KUBE_CONTEXT overrides the default context.
set -eu

CTX="${KUBE_CONTEXT:-do-nyc3-omos-cloud}"
HERE="$(cd "$(dirname "$0")" && pwd)"
PLUGINS="$HERE/../../agent/plugins"

kx() {
  kubectl --context "$CTX" -n hermes exec "$@"
}

RESTART=0
ARGS=""
for a in "$@"; do
  case "$a" in
    --restart) RESTART=1 ;;
    *) ARGS="$ARGS $a" ;;
  esac
done

STAGE="$(mktemp -d)"
trap 'rm -rf "$STAGE"' EXIT
mkdir -p "$STAGE/plugins"
cp -R "$PLUGINS/nats-platform" "$PLUGINS/nats-context" "$STAGE/plugins/"
cp "$HERE/apply_nats.py" "$STAGE/"
find "$STAGE" -name '__pycache__' -prune -exec rm -rf {} +

COPYFILE_DISABLE=1   # macOS tar would otherwise add ._* AppleDouble files
export COPYFILE_DISABLE
tar cf - -C "$STAGE" . | kx -i deploy/onemind-hermes -c hermes -- sh -c \
  'rm -rf /opt/data/nats-defs && mkdir -p /opt/data/nats-defs && tar xf - -C /opt/data/nats-defs 2>/dev/null; chown -R hermes:hermes /opt/data/nats-defs'

# Same environment as roles/apply.sh and cron/apply.sh: Hermes' dependencies live in the venv's site-packages.
# $ARGS is deliberately unquoted: it is a list of flags.
# shellcheck disable=SC2016,SC2086 # the script is meant to expand inside the pod, not here
kx deploy/onemind-hermes -c hermes -- sh -c \
  'SITE=$(ls -d /opt/hermes/.venv/lib/python3*/site-packages | head -n 1)
   runuser -u hermes -- env HERMES_HOME=/opt/data HOME=/opt/data/home PYTHONPATH="/opt/hermes:$SITE" \
     /opt/hermes/.venv/bin/python3 /opt/data/nats-defs/apply_nats.py "$@"' sh $ARGS

if [ "$RESTART" = 1 ]; then
  case "$ARGS" in
    *--dry-run*|*--status*) echo "--restart ignored with --dry-run or --status" ;;
    *)
      echo "Restarting the gateway process (s6 gateway-default); the pod stays up."
      kx deploy/onemind-hermes -c hermes -- /command/s6-svc -r /run/service/gateway-default
      ;;
  esac
fi

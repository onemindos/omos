#!/usr/bin/env bash
# Seal the shared Alertmanager -> portal webhook token into its two SealedSecrets.
#
#   monitoring/onemind-alert-webhook  key token     -> platform/monitoring/onemind-alert-webhook.sealedsecret.yaml
#   backstage/onemind-alert-webhook   key ONEMIND_ALERT_WEBHOOK_TOKEN -> platform/backstage/alert-webhook.sealedsecret.yaml
#
# Operator helper, run by hand from a clone. It is NOT referenced by any Argo CD source and
# changes nothing in the cluster: it only reads the sealed-secrets controller's public
# certificate and writes two files. Steps and verification: platform/monitoring/pending/README.md.
#
# Safety properties:
#   - Dry run by default. Files are written only with --apply.
#   - The random token lives in one shell variable. It is never printed, never written
#     unsealed, never placed in a command line (the Secret JSON goes to kubeseal on stdin),
#     and tracing (set -x) is never enabled.
#   - Refuses to overwrite: if either output exists (or is a symlink) nothing is written.
#   - Follows the repo convention in CLAUDE.md ("Secrets"): kubeseal with the controller
#     certificate, --scope namespace-wide, --format yaml.
set -euo pipefail
umask 077

CONTEXT="do-nyc3-omos-cloud"
CONTROLLER_NAME="sealed-secrets-controller"
CONTROLLER_NAMESPACE="kube-system"
SECRET_NAME="onemind-alert-webhook"
APPLY=0

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
MONITORING_OUT="$ROOT/platform/monitoring/onemind-alert-webhook.sealedsecret.yaml"
BACKSTAGE_OUT="$ROOT/platform/backstage/alert-webhook.sealedsecret.yaml"

usage() {
  cat <<'EOF'
Usage: scripts/ops/seal-alert-webhook.sh [--apply] [--context NAME]

  (no flag)       Dry run: check prerequisites and the output paths, print the plan.
                  Writes nothing.
  --apply         Generate the token in memory, seal it for both namespaces and write the
                  two SealedSecret files. Refuses to overwrite existing files.
  --context NAME  kubectl context holding the sealed-secrets controller
                  (default: do-nyc3-omos-cloud).
  -h, --help      Show this help.
EOF
}

die() {
  printf 'error: %s\n' "$*" >&2
  exit 1
}

while [ "$#" -gt 0 ]; do
  case "$1" in
    --apply) APPLY=1 ;;
    --context)
      [ "$#" -ge 2 ] || die "--context needs a value"
      CONTEXT="$2"
      shift
      ;;
    -h | --help)
      usage
      exit 0
      ;;
    *)
      usage >&2
      die "unknown argument: $1"
      ;;
  esac
  shift
done

PROBLEMS=0
check() { # check <description> <command...>
  local description="$1"
  shift
  if "$@" >/dev/null 2>&1; then
    printf '  ok      %s\n' "$description"
  else
    printf '  FAIL    %s\n' "$description"
    PROBLEMS=$((PROBLEMS + 1))
  fi
}

has_context() { kubectl config get-contexts -o name | grep -Fxq -- "$CONTEXT"; }
token_works() { [ "$(openssl rand -hex 32 | wc -c | tr -d ' ')" = "65" ]; }
path_is_free() { [ ! -e "$1" ] && [ ! -L "$1" ]; }
dir_exists() { [ -d "$(dirname "$1")" ]; }

if [ "$APPLY" -eq 1 ]; then
  printf 'Mode: APPLY (will write two files)\n'
else
  printf 'Mode: dry run (writes nothing; pass --apply to write)\n'
fi
printf 'Repo: %s\nContext: %s\n\nPrerequisites\n' "$ROOT" "$CONTEXT"

check "kubeseal on PATH" command -v kubeseal
check "kubectl on PATH" command -v kubectl
check "openssl on PATH" command -v openssl
if command -v kubectl >/dev/null 2>&1; then
  check "kubectl context '$CONTEXT' exists" has_context
fi
if command -v openssl >/dev/null 2>&1; then
  check "openssl can generate 32 random bytes" token_works
fi
check "output directory exists: platform/monitoring" dir_exists "$MONITORING_OUT"
check "output directory exists: platform/backstage" dir_exists "$BACKSTAGE_OUT"
check "does not exist yet: ${MONITORING_OUT#"$ROOT"/}" path_is_free "$MONITORING_OUT"
check "does not exist yet: ${BACKSTAGE_OUT#"$ROOT"/}" path_is_free "$BACKSTAGE_OUT"

# The controller certificate is public. Fetching it is a read-only call through the API server.
CERT_DIR=""
cleanup() {
  if [ -n "$CERT_DIR" ] && [ -d "$CERT_DIR" ]; then
    rm -rf -- "$CERT_DIR"
  fi
}
trap cleanup EXIT

if [ "$PROBLEMS" -eq 0 ]; then
  CERT_DIR="$(mktemp -d)"
  if kubeseal --fetch-cert --context "$CONTEXT" \
    --controller-name "$CONTROLLER_NAME" --controller-namespace "$CONTROLLER_NAMESPACE" \
    >"$CERT_DIR/cert.pem" 2>/dev/null && grep -q 'BEGIN CERTIFICATE' "$CERT_DIR/cert.pem"; then
    printf '  ok      sealed-secrets controller certificate fetched (%s/%s)\n' \
      "$CONTROLLER_NAMESPACE" "$CONTROLLER_NAME"
  else
    printf '  FAIL    sealed-secrets controller certificate (kubeseal --fetch-cert failed)\n'
    PROBLEMS=$((PROBLEMS + 1))
  fi
fi

printf '\nPlan\n'
printf '  1. Generate one 64-character random token in memory (never shown).\n'
printf '  2. Seal it as Secret %s in namespace monitoring, key token\n     -> %s\n' \
  "$SECRET_NAME" "${MONITORING_OUT#"$ROOT"/}"
printf '  3. Seal the same value as Secret %s in namespace backstage, key ONEMIND_ALERT_WEBHOOK_TOKEN\n     -> %s\n' \
  "$SECRET_NAME" "${BACKSTAGE_OUT#"$ROOT"/}"
printf '  Nothing is applied to the cluster. Review the files, then follow\n  platform/monitoring/pending/README.md (enable PR).\n'

if [ "$PROBLEMS" -ne 0 ]; then
  die "$PROBLEMS prerequisite(s) failed; nothing was written"
fi
if [ "$APPLY" -ne 1 ]; then
  printf '\nDry run complete. Re-run with --apply to write the two files.\n'
  exit 0
fi

# seal <namespace> <key> <token> : prints a SealedSecret manifest on stdout.
# The Secret JSON (with the base64 token) is built from shell builtins and piped on stdin, so
# the token never appears in any process argument list.
seal() {
  local namespace="$1" key="$2" token="$3" encoded
  encoded="$(printf '%s' "$token" | base64 | tr -d '\n')"
  printf '{"apiVersion":"v1","kind":"Secret","metadata":{"name":"%s","namespace":"%s"},"type":"Opaque","data":{"%s":"%s"}}' \
    "$SECRET_NAME" "$namespace" "$key" "$encoded" |
    kubeseal --cert "$CERT_DIR/cert.pem" --format yaml --scope namespace-wide
}

# looks_sealed <manifest> <namespace> <key> <token> : sanity-check kubeseal output.
looks_sealed() {
  local manifest="$1" namespace="$2" key="$3" token="$4" encoded
  encoded="$(printf '%s' "$token" | base64 | tr -d '\n')"
  case "$manifest" in *"$token"* | *"$encoded"*) return 1 ;; esac
  printf '%s\n' "$manifest" | grep -q '^kind: SealedSecret$' &&
    printf '%s\n' "$manifest" | grep -q "^  name: $SECRET_NAME\$" &&
    printf '%s\n' "$manifest" | grep -q "^  namespace: $namespace\$" &&
    printf '%s\n' "$manifest" | grep -q "^    $key: Ag"
}

TOKEN="$(openssl rand -hex 32)"
[ "${#TOKEN}" -eq 64 ] || die "token generation failed"

MONITORING_YAML="$(seal monitoring token "$TOKEN")"
BACKSTAGE_YAML="$(seal backstage ONEMIND_ALERT_WEBHOOK_TOKEN "$TOKEN")"
looks_sealed "$MONITORING_YAML" monitoring token "$TOKEN" || die "unexpected kubeseal output for monitoring; nothing was written"
looks_sealed "$BACKSTAGE_YAML" backstage ONEMIND_ALERT_WEBHOOK_TOKEN "$TOKEN" || die "unexpected kubeseal output for backstage; nothing was written"
unset TOKEN

# Re-check right before writing; noclobber makes each write fail rather than overwrite.
path_is_free "$MONITORING_OUT" || die "${MONITORING_OUT#"$ROOT"/} already exists; refusing to overwrite"
path_is_free "$BACKSTAGE_OUT" || die "${BACKSTAGE_OUT#"$ROOT"/} already exists; refusing to overwrite"
(
  set -o noclobber
  printf '%s\n' "$MONITORING_YAML" >"$MONITORING_OUT"
) || die "could not create ${MONITORING_OUT#"$ROOT"/}"
if ! (
  set -o noclobber
  printf '%s\n' "$BACKSTAGE_YAML" >"$BACKSTAGE_OUT"
); then
  rm -f -- "$MONITORING_OUT"
  die "could not create ${BACKSTAGE_OUT#"$ROOT"/}; removed the monitoring file so nothing is left half done"
fi

printf '\nWrote (SealedSecrets only; no plaintext anywhere):\n  %s\n  %s\n' \
  "${MONITORING_OUT#"$ROOT"/}" "${BACKSTAGE_OUT#"$ROOT"/}"
printf 'Next: review with git diff, then continue at "Enable" in platform/monitoring/pending/README.md.\n'

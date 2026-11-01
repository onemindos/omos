#!/usr/bin/env bash
set -euo pipefail
umask 077

: "${BRIDGE_TYPE:?BRIDGE_TYPE is required}"
: "${BRIDGE_NAME:?BRIDGE_NAME is required}"
case "$BRIDGE_TYPE" in
  whatsapp|telegram|gmessages|instagram|twitter) ;;
  *) printf 'Unsupported bridge type\n' >&2; exit 1 ;;
esac
if [[ ! "$BRIDGE_NAME" =~ ^sh-[a-z0-9-]{1,29}$ ]]; then
  printf 'Invalid self-hosted bridge name\n' >&2
  exit 1
fi

data_dir="${BEEPER_DATA_DIR:-/data}"
runtime_dir="${BEEPER_RUNTIME_DIR:-/run/beeper}"
token_file="${BEEPER_TOKEN_FILE:-/var/run/beeper-auth/MATRIX_ACCESS_TOKEN}"
mkdir -p "$data_dir/db" "$data_dir/cache" "$runtime_dir"
if ! jq -en --rawfile token "$token_file" \
  '$token | sub("[\r\n]+$"; "") | test("^(syt_|bat_)[^[:space:]]+$")' >/dev/null; then
  printf 'Missing or invalid Beeper access token\n' >&2
  exit 1
fi

# bbctl updates its config, so never point it at the read-only Secret volume.
config_tmp="$(mktemp "$runtime_dir/.bbctl.XXXXXX")"
trap 'rm -f -- "$config_tmp"' EXIT
jq -n --rawfile token "$token_file" --arg data_dir "$data_dir" --arg database_dir "$data_dir/db" \
  '{environments: {prod: {access_token: ($token | sub("[\r\n]+$"; "")), bridge_data_dir: $data_dir, database_dir: $database_dir}}}' \
  > "$config_tmp"
chmod 0600 "$config_tmp"
export BBCTL_CONFIG="$runtime_dir/bbctl.json"
mv -f -- "$config_tmp" "$BBCTL_CONFIG"

# exec lets Kubernetes SIGTERM reach bbctl, which stops its connector child.
# The custom command bypasses bbctl's otherwise automatic binary downloads.
exec bbctl -e prod run --type "$BRIDGE_TYPE" --no-override-config \
  --custom-startup-command /opt/beeper-runtime/start-connector.sh "$BRIDGE_NAME"

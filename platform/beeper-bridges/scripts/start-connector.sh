#!/usr/bin/env bash
set -euo pipefail
umask 077

: "${BRIDGE_TYPE:?BRIDGE_TYPE is required}"
case "$BRIDGE_TYPE" in
  whatsapp|telegram|gmessages|instagram|twitter) ;;
  *) printf 'Unsupported bridge type\n' >&2; exit 1 ;;
esac
if [[ "$#" != 2 || "$1" != -c || ! -f "$2" ]]; then
  printf 'Expected bbctl-generated connector configuration\n' >&2
  exit 1
fi
if [[ "$BRIDGE_TYPE" == telegram ]]; then
  if [[ ! "${TELEGRAM_API_ID:-}" =~ ^[1-9][0-9]*$ || ! "${TELEGRAM_API_HASH:-}" =~ ^[a-fA-F0-9]{32}$ ]]; then
    printf 'Missing or invalid Telegram application credentials\n' >&2
    exit 1
  fi
fi

# Apply after bbctl generation, including any fallback regeneration. Credential
# values are read from the environment inside awk, never command-line arguments.
# The generated YAML has one appservice block and one network block; fail closed
# if that structure changes rather than silently starting with other credentials.
config_path="$2"
config_tmp="$(mktemp "${config_path}.XXXXXX")"
trap 'rm -f -- "$config_tmp"' EXIT
if ! awk '
  /^[a-zA-Z_][a-zA-Z0-9_]*:/ { section = $0; sub(/:.*/, "", section) }
  section == "appservice" && /^[ \t]+hostname:/ { sub(/hostname:.*/, "hostname: 127.0.0.1"); host_count++ }
  section == "appservice" && /^[ \t]+port:/ { sub(/port:.*/, "port: 4000"); port_count++ }
  ENVIRON["BRIDGE_TYPE"] == "telegram" && section == "network" && /^[ \t]+api_id:/ {
    sub(/api_id:.*/, "api_id: " ENVIRON["TELEGRAM_API_ID"]); id_count++
  }
  ENVIRON["BRIDGE_TYPE"] == "telegram" && section == "network" && /^[ \t]+api_hash:/ {
    sub(/api_hash:.*/, "api_hash: \"" ENVIRON["TELEGRAM_API_HASH"] "\""); hash_count++
  }
  { print }
  END {
    if (host_count != 1 || port_count != 1) exit 1
    if (ENVIRON["BRIDGE_TYPE"] == "telegram" && (id_count != 1 || hash_count != 1)) exit 1
  }
' "$config_path" > "$config_tmp"; then
  printf 'Unexpected generated connector configuration; refusing startup\n' >&2
  exit 1
fi
chmod 0600 "$config_tmp"
mv -f -- "$config_tmp" "$config_path"
unset TELEGRAM_API_ID TELEGRAM_API_HASH
exec "${BEEPER_DATA_DIR:-/data}/bin/mautrix-$BRIDGE_TYPE" "$@"

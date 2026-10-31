#!/usr/bin/env bash
set -euo pipefail
umask 077

# The init container downloads only the reviewed, checksum-locked Linux binary.
: "${BRIDGE_TYPE:?BRIDGE_TYPE is required}"
case "$BRIDGE_TYPE" in
  whatsapp|telegram|gmessages|instagram|twitter) ;;
  *) printf 'Unsupported bridge type\n' >&2; exit 1 ;;
esac
if [[ "$(uname -m)" != x86_64 ]]; then
  printf 'Connector lock requires a Linux amd64 node\n' >&2
  exit 1
fi

data_dir="${BEEPER_DATA_DIR:-/data}"
lock_file="${BEEPER_ARTIFACT_LOCK:-/etc/beeper-artifacts/artifacts.json}"
url="$(jq -er --arg type "$BRIDGE_TYPE" '.[$type].url | select(type == "string" and length > 0)' "$lock_file")"
checksum="$(jq -er --arg type "$BRIDGE_TYPE" '.[$type].sha256 | select(type == "string" and length > 0)' "$lock_file")"
binary="$(jq -er --arg type "$BRIDGE_TYPE" '.[$type].binary | select(type == "string" and length > 0)' "$lock_file")"
if [[ "$binary" != "mautrix-$BRIDGE_TYPE" || ! "$checksum" =~ ^[a-f0-9]{64}$ || "$url" != https://mau.dev/mautrix/* ]]; then
  printf 'Invalid connector artifact lock\n' >&2
  exit 1
fi

mkdir -p "$data_dir/bin"
target="$data_dir/bin/$binary"
if [[ -f "$target" ]] && printf '%s  %s\n' "$checksum" "$target" | sha256sum -c - >/dev/null 2>&1; then
  chmod 0755 "$target"
  exit 0
fi

artifact_tmp="$(mktemp "$data_dir/bin/.$binary.XXXXXX")"
trap 'rm -f -- "$artifact_tmp"' EXIT
curl --fail --silent --show-error --location --proto '=https' --proto-redir '=https' \
  --connect-timeout 20 --max-time 300 --retry 3 --retry-delay 3 \
  --output "$artifact_tmp" "$url"
if ! printf '%s  %s\n' "$checksum" "$artifact_tmp" | sha256sum -c - >/dev/null 2>&1; then
  printf 'Connector checksum mismatch; refusing startup\n' >&2
  exit 1
fi
chmod 0755 "$artifact_tmp"
mv -f -- "$artifact_tmp" "$target"

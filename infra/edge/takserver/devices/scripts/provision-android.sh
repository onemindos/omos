#!/usr/bin/env bash
set -euo pipefail

# Provision a hard-wired (USB) Android device for ATAK.
# Steps: verify adb device -> install ATAK (+ optional plugins) -> push device .p12 + CA.
# The device cert must already be issued via deploy/scripts/tak-add-device.sh.
#
# Usage:
#   devices/scripts/provision-android.sh <cert-common-name> [--small] [--plugins]
# Example:
#   devices/scripts/provision-android.sh field-android-01 --plugins

CERT_NAME="${1:-}"
if [[ -z "$CERT_NAME" ]]; then
  echo "Usage: $0 <cert-common-name> [--small] [--plugins]" >&2
  exit 1
fi
shift || true

SMALL=false
PLUGINS=false
for arg in "$@"; do
  case "$arg" in
    --small) SMALL=true ;;
    --plugins) PLUGINS=true ;;
    *) echo "Unknown flag: $arg" >&2; exit 1 ;;
  esac
done

HOME_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
CLIENT="$HOME_DIR/client"
EXPORTS="$HOME_DIR/deploy/exports"
CERT_P12="$EXPORTS/$CERT_NAME.p12"
CA_PEM="$HOME_DIR/server/standard/takserver-docker-5.7-RELEASE-43/tak/certs/files/ca.pem"

command -v adb >/dev/null || { echo "adb not found (brew install android-platform-tools)" >&2; exit 1; }

echo "== Waiting for a single authorized device =="
adb wait-for-device
mapfile -t DEVS < <(adb devices | awk 'NR>1 && $2=="device"{print $1}')
if [[ "${#DEVS[@]}" -ne 1 ]]; then
  echo "Need exactly ONE authorized device. Found: ${#DEVS[@]}" >&2
  adb devices -l >&2
  exit 1
fi
echo "Device: ${DEVS[0]}  model: $(adb shell getprop ro.product.model | tr -d '\r')"

# Detect device ABI. civSmall is arm64-ONLY and will fail on 32-bit (armeabi-v7a) devices.
ABI="$(adb -s "${DEVS[0]}" shell getprop ro.product.cpu.abi | tr -d '\r')"
echo "Device ABI: $ABI"

# Pick ATAK APK. Full civ is a fat APK (v7a + arm64) — always safe. --small only on arm64.
if $SMALL; then
  if [[ "$ABI" != arm64* ]]; then
    echo "REFUSING --small: civSmall is arm64-only but device is $ABI. Using full civ." >&2
    SMALL=false
  fi
fi
if $SMALL; then
  APK="$(ls "$CLIENT"/*civSmall*.apk 2>/dev/null | head -1)"
else
  APK="$(ls "$CLIENT"/*civ-release.apk 2>/dev/null | grep -v civSmall | head -1)"
fi
[[ -f "$APK" ]] || { echo "ATAK APK not found in $CLIENT" >&2; exit 1; }

# --no-streaming avoids the macOS 'splice failed: EINVAL' bug on large APKs over USB.
echo "== Installing ATAK: $(basename "$APK") =="
adb -s "${DEVS[0]}" install --no-streaming -r -g "$APK"

if $PLUGINS; then
  echo "== Installing plugins from $CLIENT/loadout/Plugins =="
  find "$CLIENT/loadout/Plugins" -maxdepth 1 -name '*.apk' -print0 | while IFS= read -r -d '' p; do
    echo "  plugin: $(basename "$p")"
    adb -s "${DEVS[0]}" install --no-streaming -r -g "$p" || echo "  (failed: $(basename "$p"))"
  done
fi

echo "== Pushing cert material to /sdcard/atak/cert/ =="
adb shell mkdir -p /sdcard/atak/cert 2>/dev/null || true
[[ -f "$CERT_P12" ]] && adb push "$CERT_P12" /sdcard/atak/cert/ || echo "WARN: $CERT_P12 missing — issue it via tak-add-device.sh"
[[ -f "$CA_PEM" ]] && adb push "$CA_PEM" /sdcard/atak/cert/ca.pem || true

cat <<EOF

Done. On the device, open ATAK and:
  Settings -> Network Connections -> Manage Server Connections -> Add
    Server:   192.168.1.183   (LAN)   or   100.74.182.117 (Tailscale)
    Port:     8089
    Protocol: SSL/TLS
    Client cert:  /sdcard/atak/cert/$CERT_NAME.p12   (pass: the CAPASS you issued it with)
    CA cert:      /sdcard/atak/cert/ca.pem
Register this device in devices/registry.md.
EOF

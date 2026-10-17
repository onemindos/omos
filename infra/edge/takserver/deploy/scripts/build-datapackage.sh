#!/usr/bin/env bash
set -euo pipefail

# Build a TAK "Mission Data Package" (.zip) that auto-configures a device:
# server connection + client cert + CA truststore. Works for BOTH iTAK (iOS) and ATAK (Android).
# Import = one tap (AirDrop to iOS, or adb push + import on Android).
#
# Usage:  deploy/scripts/build-datapackage.sh <cert-common-name> [description]
# Example: deploy/scripts/build-datapackage.sh omos-iphone17 "OneMind TAK"

CN="${1:-}"
DESC="${2:-OneMind TAK}"
[[ -n "$CN" ]] || { echo "Usage: $0 <cert-common-name> [description]" >&2; exit 1; }

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
STACK_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
EXPORTS="$STACK_DIR/exports"
CERT_P12="$EXPORTS/$CN.p12"
CA_PEM="$EXPORTS/ca.pem"
CAPASS="${CAPASS:-atakatak}"
HOST="${TAK_HOST:-100.74.182.117}"   # Tailscale IP (server cert SAN covers it)
PORT="${TAK_PORT:-8089}"

[[ -f "$CERT_P12" ]] || { echo "Missing client cert: $CERT_P12 (run tak-add-device.sh $CN ... first)" >&2; exit 1; }
[[ -f "$CA_PEM" ]]   || { echo "Missing CA: $CA_PEM" >&2; exit 1; }

UID_STR="$(uuidgen)"
WORK="$(mktemp -d)"
mkdir -p "$WORK/MANIFEST"

# FLAT layout: certs + .pref at the ZIP ROOT (no cert/ subfolder).
# iTAK does NOT follow a cert/ subfolder — it only reads root-level files and errors with
# "the format of this configuration could not be read" otherwise. ATAK is fine with root too
# (it resolves whatever zipEntry paths the manifest declares). So root = works for BOTH.
openssl pkcs12 -export -nokeys -legacy -in "$CA_PEM" \
  -out "$WORK/truststore-onemind.p12" -passout pass:"$CAPASS" -name "OneMind-CA"
cp "$CERT_P12" "$WORK/$CN.p12"

# Preference file. cert location/password keys live in the cot_streams block, referenced by
# BARE FILENAME (files are at root).
cat > "$WORK/$CN.pref" <<EOF
<?xml version='1.0' encoding='UTF-8' standalone='yes'?>
<preferences>
  <preference version="1" name="cot_streams">
    <entry key="count" class="class java.lang.Integer">1</entry>
    <entry key="description0" class="class java.lang.String">$DESC</entry>
    <entry key="enabled0" class="class java.lang.Boolean">true</entry>
    <entry key="connectString0" class="class java.lang.String">$HOST:$PORT:ssl</entry>
    <entry key="caLocation0" class="class java.lang.String">$CN-truststore.p12</entry>
    <entry key="caPassword0" class="class java.lang.String">$CAPASS</entry>
    <entry key="certificateLocation0" class="class java.lang.String">$CN.p12</entry>
    <entry key="clientPassword0" class="class java.lang.String">$CAPASS</entry>
    <entry key="enrollForCertificateWithTrust0" class="class java.lang.Boolean">false</entry>
  </preference>
  <preference version="1" name="com.atakmap.app_preferences">
    <entry key="displayServerConnectionWidget" class="class java.lang.Boolean">true</entry>
  </preference>
</preferences>
EOF
# Name the truststore per-device so iTAK's flat namespace can't collide between devices.
mv "$WORK/truststore-onemind.p12" "$WORK/$CN-truststore.p12"

# Manifest (bare filename zipEntries). Written to BOTH MANIFEST/manifest.xml (ATAK convention)
# and root manifest.xml (some iTAK builds look for it at root).
MANIFEST_XML="<?xml version=\"1.0\" encoding=\"UTF-8\"?>
<MissionPackageManifest version=\"2\">
  <Configuration>
    <Parameter name=\"uid\" value=\"$UID_STR\"/>
    <Parameter name=\"name\" value=\"$CN-onemind.zip\"/>
    <Parameter name=\"onReceiveDelete\" value=\"true\"/>
  </Configuration>
  <Contents>
    <Content ignore=\"false\" zipEntry=\"$CN.pref\"/>
    <Content ignore=\"false\" zipEntry=\"$CN.p12\"/>
    <Content ignore=\"false\" zipEntry=\"$CN-truststore.p12\"/>
  </Contents>
</MissionPackageManifest>"
printf '%s\n' "$MANIFEST_XML" > "$WORK/MANIFEST/manifest.xml"
printf '%s\n' "$MANIFEST_XML" > "$WORK/manifest.xml"

OUT="$EXPORTS/$CN-onemind.zip"
rm -f "$OUT"
( cd "$WORK" && zip -qr "$OUT" MANIFEST manifest.xml "$CN.pref" "$CN.p12" "$CN-truststore.p12" )
chmod 600 "$OUT"
rm -rf "$WORK"

echo "Built data package: $OUT"
echo "  server:  $HOST:$PORT (ssl)    cert: $CN.p12    pass: $CAPASS"
echo "  iOS:     AirDrop the .zip -> Open in iTAK -> it auto-adds the server."
echo "  Android: adb push '$OUT' /sdcard/atak/tools/datapackage/ ; import in ATAK (Import Manager)."

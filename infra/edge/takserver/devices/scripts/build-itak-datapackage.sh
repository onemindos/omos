#!/usr/bin/env bash
# Build an iTAK/ATAK soft-cert data package (bundles client cert + CA truststore).
# iTAK does NOT use iOS system trust — the CA must be inside the package. This is
# the reliable route for a self-signed TAK Server: no enrollment, no trust prompt.
#
# Format matches the maintained reference (talentedbrute/takserver-datapackage-builder
# and Cloud-RF/tak-server certDP.sh), verified working with iTAK 5.x:
#   - MANIFEST/manifest.xml   (flat zipEntry names — NO cert/ prefix here)
#   - secure.pref             (cert paths DO use cert/ prefix; non-indexed keys in
#                              com.atakmap.app_preferences; connectString0 in cot_streams)
#   - <client>.p12  <ca-truststore>.p12   (flat at zip root)
# iTAK relocates the p12s into its internal cert/ dir on import, which is why the
# manifest is flat but the pref points at cert/.
#
# Usage: build-itak-datapackage.sh <user> <client.p12> <truststore.p12> <host> <out.zip>
#   e.g. build-itak-datapackage.sh omos-iphone17 \
#          server/.../certs/files/omos-iphone17.p12 \
#          server/.../certs/files/truststore-root.p12 \
#          100.74.182.117 deploy/exports/OneMind-iPhone17.zip
# P12_PASS (env, required) is the password of the client .p12 and the truststore .p12; there is no default.
set -euo pipefail

USER="${1:?user}"; CLIENT_P12="${2:?client .p12}"; CA_P12="${3:?truststore .p12}"
HOST="${4:?host/ip}"; OUT="${5:?output .zip}"
PASS="${P12_PASS:?set P12_PASS (the password of the client and truststore .p12)}"
PORT="${TAK_COT_PORT:-8089}"
DESC="${PKG_DESC:-OneMind TAK}"

command -v zip >/dev/null || { echo "zip not found"; exit 1; }
[ -f "$CLIENT_P12" ] || { echo "missing client p12: $CLIENT_P12"; exit 1; }
[ -f "$CA_P12" ] || { echo "missing CA p12: $CA_P12"; exit 1; }

work="$(mktemp -d)"; trap 'rm -rf "$work"' EXIT
mkdir -p "$work/MANIFEST"
cp "$CLIENT_P12" "$work/${USER}.p12"
cp "$CA_P12" "$work/truststore.p12"
uuid="$(uuidgen)"

cat > "$work/secure.pref" <<EOF
<?xml version='1.0' encoding='ASCII' standalone='yes'?>
<preferences>
  <preference version="1" name="cot_streams">
    <entry key="count" class="class java.lang.Integer">1</entry>
    <entry key="description0" class="class java.lang.String">${DESC}</entry>
    <entry key="enabled0" class="class java.lang.Boolean">true</entry>
    <entry key="connectString0" class="class java.lang.String">${HOST}:${PORT}:ssl</entry>
  </preference>
  <preference version="1" name="com.atakmap.app_preferences">
    <entry key="caLocation" class="class java.lang.String">cert/truststore.p12</entry>
    <entry key="caPassword" class="class java.lang.String">${PASS}</entry>
    <entry key="certificateLocation" class="class java.lang.String">cert/${USER}.p12</entry>
    <entry key="clientPassword" class="class java.lang.String">${PASS}</entry>
    <entry key="displayServerConnectionWidget" class="class java.lang.Boolean">true</entry>
  </preference>
</preferences>
EOF

cat > "$work/MANIFEST/manifest.xml" <<EOF
<MissionPackageManifest version="2">
  <Configuration>
    <Parameter name="uid" value="${uuid}"/>
    <Parameter name="name" value="${USER}"/>
    <Parameter name="onReceiveDelete" value="true"/>
  </Configuration>
  <Contents>
    <Content ignore="false" zipEntry="truststore.p12"/>
    <Content ignore="false" zipEntry="${USER}.p12"/>
    <Content ignore="false" zipEntry="secure.pref"/>
  </Contents>
</MissionPackageManifest>
EOF

mkdir -p "$(dirname "$OUT")"
# resolve OUT to an absolute path (we cd into $work before zipping)
OUT="$(cd "$(dirname "$OUT")" && pwd)/$(basename "$OUT")"
rm -f "$OUT"
( cd "$work" && zip -q -X -r "$OUT" MANIFEST secure.pref truststore.p12 "${USER}.p12" )
echo "Built $OUT"
unzip -l "$OUT"

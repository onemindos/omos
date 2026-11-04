#!/usr/bin/env bash
# Edge MediaMTX for a OneMind site (FEATURES K9). Runs MediaMTX on the Mac or the home server,
# serves LAN readers, and forwards every allowed stream to the cloud onemind-media over SRT.
#
#   edge-media.sh run           run the edge (foreground; Ctrl-C stops it)
#   edge-media.sh selftest      build-machine loop: a local "cloud" MediaMTX running the cluster
#                               config (platform/onetak/media/mediamtx.yml), this edge, an ffmpeg
#                               test pattern as drone/sim-x1, read back from the cloud
#   edge-media.sh keychain-set  store the three passwords in the macOS login keychain
#
# Credentials (never in Git, never on the command line):
#   edge-home       the cloud onemind-media password of the site identity edge-home
#                   (the owner's onemind-media-users Secret, key edge-home)
#   edge-publisher  LAN publishers to this edge
#   lan-viewer      LAN readers of this edge
# Each comes from the environment (ONEMIND_MEDIA_EDGE_HOME, ONEMIND_MEDIA_EDGE_PUBLISHER,
# ONEMIND_MEDIA_LAN_VIEWER) or, on macOS, from the keychain item service "onemind-media",
# account <name> (security find-generic-password -s onemind-media -a <name> -w).
#
# Settings (environment): CLOUD_HOST (default media.tail717ba0.ts.net), CLOUD_SRT_PORT (8890),
# EDGE_BIND (0.0.0.0), EDGE_RTSP_PORT (8554), EDGE_HLS_PORT (8888), EDGE_WEBRTC_PORT (8889),
# EDGE_ICE_PORT (8189), EDGE_SRT_PORT (8890), EDGE_API_PORT (9997), EDGE_PLAYBACK_PORT (9996),
# EDGE_RECORD (true), EDGE_RECORD_DIR (~/Library/Application Support/OneMind/edge-media/recordings),
# MEDIAMTX_BIN (mediamtx), FFMPEG_BIN (ffmpeg), FFPROBE_BIN (ffprobe).
set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO="$(cd "$HERE/../../.." && pwd)"
TEMPLATE="$HERE/mediamtx-edge.yml.tmpl"
CLOUD_CONFIG="$REPO/platform/onetak/media/mediamtx.yml"
MEDIAMTX_BIN="${MEDIAMTX_BIN:-mediamtx}"
FFMPEG_BIN="${FFMPEG_BIN:-ffmpeg}"
FFPROBE_BIN="${FFPROBE_BIN:-ffprobe}"
KEYCHAIN_SERVICE="onemind-media"

WORK=""
PIDS=()
cleanup() {
  local pid
  for pid in "${PIDS[@]:-}"; do
    [ -n "$pid" ] && kill "$pid" 2>/dev/null || true
  done
  wait 2>/dev/null || true
  [ -n "$WORK" ] && rm -rf "$WORK"
}
trap cleanup EXIT INT TERM

die() { echo "edge-media: $*" >&2; exit 1; }

need() { command -v "$1" >/dev/null 2>&1 || die "$1 not found (set ${2:-PATH})"; }

# secret NAME ENVVAR: print a credential from the environment or the macOS keychain.
secret() {
  local name="$1" var="$2" val=""
  val="${!var:-}"
  if [ -z "$val" ] && command -v security >/dev/null 2>&1; then
    val="$(security find-generic-password -s "$KEYCHAIN_SERVICE" -a "$name" -w 2>/dev/null || true)"
  fi
  [ -n "$val" ] || die "no credential for $name: set $var or run '$0 keychain-set'"
  case "$val" in
    *[!A-Za-z0-9._~-]*) die "credential $name must be URL-safe (letters, digits, . _ ~ -)";;
  esac
  printf '%s' "$val"
}

# render OUT: fill the template from the variables below.
render() {
  sed -e "s|@EDGE_PUBLISHER_PASS@|$EDGE_PUBLISHER_PASS|" \
      -e "s|@LAN_VIEWER_PASS@|$LAN_VIEWER_PASS|" \
      -e "s|@EDGE_HOME_PASS@|$EDGE_HOME_PASS|" \
      -e "s|@CLOUD_HOST@|$CLOUD_HOST|" -e "s|@CLOUD_SRT_PORT@|$CLOUD_SRT_PORT|" \
      -e "s|@BIND@|$EDGE_BIND|g" -e "s|@RTSP_PORT@|$EDGE_RTSP_PORT|" -e "s|@HLS_PORT@|$EDGE_HLS_PORT|" \
      -e "s|@WEBRTC_PORT@|$EDGE_WEBRTC_PORT|" -e "s|@ICE_PORT@|$EDGE_ICE_PORT|" \
      -e "s|@SRT_PORT@|$EDGE_SRT_PORT|" -e "s|@API_PORT@|$EDGE_API_PORT|" \
      -e "s|@PLAYBACK_PORT@|$EDGE_PLAYBACK_PORT|" \
      -e "s|@RECORD@|$EDGE_RECORD|g" -e "s|@RECORD_DIR@|$EDGE_RECORD_DIR|" \
      "$TEMPLATE" > "$1"
  chmod 600 "$1"
}

# redact: copy stdin to stdout with the edge-home password replaced. MediaMTX v1.21.3 logs the
# full forward URL, which carries that password in the SRT streamid.
redact() {
  local line
  while IFS= read -r line; do
    printf '%s\n' "${line//"$EDGE_HOME_PASS"/<redacted>}"
  done
}

defaults() {
  CLOUD_HOST="${CLOUD_HOST:-media.tail717ba0.ts.net}"
  CLOUD_SRT_PORT="${CLOUD_SRT_PORT:-8890}"
  EDGE_BIND="${EDGE_BIND:-0.0.0.0}"
  EDGE_RTSP_PORT="${EDGE_RTSP_PORT:-8554}"
  EDGE_HLS_PORT="${EDGE_HLS_PORT:-8888}"
  EDGE_WEBRTC_PORT="${EDGE_WEBRTC_PORT:-8889}"
  EDGE_ICE_PORT="${EDGE_ICE_PORT:-8189}"
  EDGE_SRT_PORT="${EDGE_SRT_PORT:-8890}"
  EDGE_API_PORT="${EDGE_API_PORT:-9997}"
  EDGE_PLAYBACK_PORT="${EDGE_PLAYBACK_PORT:-9996}"
  EDGE_RECORD="${EDGE_RECORD:-true}"
  EDGE_RECORD_DIR="${EDGE_RECORD_DIR:-$HOME/Library/Application Support/OneMind/edge-media/recordings}"
}

cmd_run() {
  need "$MEDIAMTX_BIN" MEDIAMTX_BIN
  defaults
  EDGE_HOME_PASS="$(secret edge-home ONEMIND_MEDIA_EDGE_HOME)"
  EDGE_PUBLISHER_PASS="$(secret edge-publisher ONEMIND_MEDIA_EDGE_PUBLISHER)"
  LAN_VIEWER_PASS="$(secret lan-viewer ONEMIND_MEDIA_LAN_VIEWER)"
  mkdir -p "$EDGE_RECORD_DIR"
  WORK="$(umask 077; mktemp -d)"
  render "$WORK/mediamtx.yml"
  echo "edge-media: LAN rtsp://<this-host>:$EDGE_RTSP_PORT/<class>/<id>, forwarding to srt://$CLOUD_HOST:$CLOUD_SRT_PORT as edge-home"
  "$MEDIAMTX_BIN" "$WORK/mediamtx.yml" > >(redact) 2>&1 &
  PIDS+=("$!")
  wait "${PIDS[0]}"
}

cmd_keychain_set() {
  command -v security >/dev/null 2>&1 || die "macOS 'security' not found; use the environment variables instead"
  local name
  for name in edge-home edge-publisher lan-viewer; do
    echo "Password for $name (the input is hidden; URL-safe characters only):"
    # -w with no value makes security prompt, so the password never appears in argv or history.
    security add-generic-password -U -s "$KEYCHAIN_SERVICE" -a "$name" -w
  done
}

pass() { echo "PASS $*"; }
pw() { cat "$WORK/pw.$1"; }
fail() { echo "FAIL $*" >&2; exit 1; }

cmd_selftest() {
  need "$MEDIAMTX_BIN" MEDIAMTX_BIN
  need "$FFMPEG_BIN" FFMPEG_BIN
  need "$FFPROBE_BIN" FFPROBE_BIN
  need curl
  need openssl
  WORK="$(umask 077; mktemp -d)"
  local u i p
  # Throwaway passwords, generated per run into the private work dir and never printed.
  # (Plain files, not an associative array: macOS ships bash 3.2.)
  for u in mavlink-bridge edge-home onemind-client tak-viewer takml video-registry edge-publisher lan-viewer; do
    openssl rand -hex 12 > "$WORK/pw.$u"
  done
  local CL_RTSP="${SELFTEST_CLOUD_RTSP:-38554}" CL_SRT="${SELFTEST_CLOUD_SRT:-38890}" CL_API="${SELFTEST_CLOUD_API:-39997}"
  local CL_HLS="${SELFTEST_CLOUD_HLS:-38888}" CL_WEBRTC="${SELFTEST_CLOUD_WEBRTC:-38889}" CL_ICE="${SELFTEST_CLOUD_ICE:-38189}"
  local CL_PB="${SELFTEST_CLOUD_PLAYBACK:-39996}" CL_MET="${SELFTEST_CLOUD_METRICS:-39998}"

  # 1. The "cloud": the cluster config, only addresses and the record dir changed.
  local cloud_env=()
  i=0
  for u in mavlink-bridge edge-home onemind-client tak-viewer takml video-registry; do
    cloud_env+=("MTX_AUTHINTERNALUSERS_${i}_PASS=$(pw "$u")")
    i=$((i + 1))
  done
  grep -q 'user: video-registry' "$CLOUD_CONFIG" || fail "cluster config not found at $CLOUD_CONFIG"
  env "${cloud_env[@]}" \
    MTX_RTSPADDRESS="127.0.0.1:$CL_RTSP" MTX_SRTADDRESS="127.0.0.1:$CL_SRT" MTX_APIADDRESS="127.0.0.1:$CL_API" \
    MTX_HLSADDRESS="127.0.0.1:$CL_HLS" MTX_WEBRTCADDRESS="127.0.0.1:$CL_WEBRTC" \
    MTX_WEBRTCLOCALUDPADDRESS="127.0.0.1:$CL_ICE" MTX_PLAYBACKADDRESS="127.0.0.1:$CL_PB" \
    MTX_METRICSADDRESS="127.0.0.1:$CL_MET" MTX_PATHDEFAULTS_RECORD=false \
    "$MEDIAMTX_BIN" "$CLOUD_CONFIG" > "$WORK/cloud.log" 2>&1 &
  PIDS+=("$!")

  # 2. The edge, from the template, forwarding to the local cloud.
  CLOUD_HOST=127.0.0.1 CLOUD_SRT_PORT="$CL_SRT"
  EDGE_BIND=127.0.0.1
  EDGE_RTSP_PORT="${SELFTEST_EDGE_RTSP:-48554}" EDGE_HLS_PORT="${SELFTEST_EDGE_HLS:-48888}"
  EDGE_WEBRTC_PORT="${SELFTEST_EDGE_WEBRTC:-48889}" EDGE_ICE_PORT="${SELFTEST_EDGE_ICE:-48189}"
  EDGE_SRT_PORT="${SELFTEST_EDGE_SRT:-48890}" EDGE_API_PORT="${SELFTEST_EDGE_API:-49997}"
  EDGE_PLAYBACK_PORT="${SELFTEST_EDGE_PLAYBACK:-49996}" EDGE_RECORD=false EDGE_RECORD_DIR="$WORK/rec"
  EDGE_HOME_PASS="$(pw edge-home)" EDGE_PUBLISHER_PASS="$(pw edge-publisher)" LAN_VIEWER_PASS="$(pw lan-viewer)"
  render "$WORK/edge.yml"
  "$MEDIAMTX_BIN" "$WORK/edge.yml" > >(redact > "$WORK/edge.log") 2>&1 &
  PIDS+=("$!")
  sleep 2

  # 3. Anonymous publish to the edge is refused.
  if "$FFMPEG_BIN" -hide_banner -loglevel error -re -f lavfi -i testsrc2=size=640x360:rate=25 -t 1 \
      -c:v libx264 -preset ultrafast -f rtsp -rtsp_transport tcp \
      "rtsp://127.0.0.1:$EDGE_RTSP_PORT/drone/sim-x1" > "$WORK/anon.err" 2>&1; then
    fail "anonymous publish to the edge was accepted"
  fi
  grep -q 401 "$WORK/anon.err" && pass "anonymous publish to the edge refused (401)" || fail "check failed: anonymous publish to the edge refused (401)"

  # 4. sim-x1's camera publishes to the edge.
  "$FFMPEG_BIN" -hide_banner -loglevel error -re -f lavfi -i testsrc2=size=640x360:rate=25 -t 20 \
    -c:v libx264 -preset ultrafast -tune zerolatency -g 25 -pix_fmt yuv420p -f rtsp -rtsp_transport tcp \
    "rtsp://edge-publisher:$(pw edge-publisher)@127.0.0.1:$EDGE_RTSP_PORT/drone/sim-x1" > "$WORK/pub.err" 2>&1 &
  PIDS+=("$!")

  # 5. The edge forwards over SRT as edge-home; the cloud shows the path from an SRT connection.
  # API credentials through a netrc file, so the password is not in curl's argv.
  printf 'machine 127.0.0.1 login video-registry password %s\n' "$(pw video-registry)" > "$WORK/netrc"
  local src=""
  for p in $(seq 1 30); do
    src="$(curl -s --netrc-file "$WORK/netrc" "http://127.0.0.1:$CL_API/v3/paths/get/drone/sim-x1" \
      | sed -n 's/.*"source":{"type":"\([a-zA-Z]*\)".*/\1/p')" || true
    [ "$src" = "srtConn" ] && break
    sleep 0.5
  done
  [ "$src" = "srtConn" ] || { sed "s/$(pw edge-home)/<redacted>/g" "$WORK/edge.log" | tail -20 >&2; fail "cloud path drone/sim-x1 not fed over SRT (source: '${src}')"; }
  pass "cloud drone/sim-x1 is fed by the edge over SRT (source.type=srtConn) after ${p} checks"

  # 6. Read back: a cloud client reads the forwarded stream; a LAN viewer reads the edge.
  local out
  out="$("$FFPROBE_BIN" -v error -rtsp_transport tcp -show_entries stream=codec_name,width,height -of csv=p=0 \
    "rtsp://onemind-client:$(pw onemind-client)@127.0.0.1:$CL_RTSP/drone/sim-x1" 2>&1)" || fail "cloud read failed"
  [ "$out" = "h264,640,360" ] && pass "cloud read of drone/sim-x1 as onemind-client: $out" || fail "check failed: cloud read of drone/sim-x1 as onemind-client: $out"
  out="$("$FFPROBE_BIN" -v error -rtsp_transport tcp -show_entries stream=codec_name,width,height -of csv=p=0 \
    "rtsp://lan-viewer:$(pw lan-viewer)@127.0.0.1:$EDGE_RTSP_PORT/drone/sim-x1" 2>&1)" || fail "edge LAN read failed"
  [ "$out" = "h264,640,360" ] && pass "edge LAN read of drone/sim-x1 as lan-viewer: $out" || fail "check failed: edge LAN read of drone/sim-x1 as lan-viewer: $out"
  if "$FFPROBE_BIN" -v error -rtsp_transport tcp "rtsp://127.0.0.1:$EDGE_RTSP_PORT/drone/sim-x1" > "$WORK/anonread.err" 2>&1; then
    fail "anonymous read of the edge was accepted"
  fi
  grep -q 401 "$WORK/anonread.err" && pass "anonymous read of the edge refused (401)" || fail "check failed: anonymous read of the edge refused (401)"
  if "$FFPROBE_BIN" -v error -rtsp_transport tcp "rtsp://127.0.0.1:$CL_RTSP/drone/sim-x1" > "$WORK/anoncloud.err" 2>&1; then
    fail "anonymous read of the cloud was accepted"
  fi
  grep -q 401 "$WORK/anoncloud.err" && pass "anonymous read of the cloud refused (401)" || fail "check failed: anonymous read of the cloud refused (401)"
  # 7. No log holds the edge-home password (edge log redacted; the cloud never logs it).
  sleep 1
  if grep -q "$(pw edge-home)" "$WORK/edge.log" "$WORK/cloud.log"; then
    fail "the edge-home password appears in a MediaMTX log"
  fi
  grep -q '<redacted>' "$WORK/edge.log" && pass "edge log redacts the forward password; cloud log does not contain it" \
    || fail "check failed: edge log redaction"
  echo "SELFTEST OK"
}

case "${1:-}" in
  run) cmd_run ;;
  selftest) cmd_selftest ;;
  keychain-set) cmd_keychain_set ;;
  *) sed -n '2,12p' "$0"; exit 2 ;;
esac

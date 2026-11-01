#!/usr/bin/env bash
# End-to-end test of the task service (FEATURES D2, H4; mission-task.md section 10 tests
# 1, 3, 4 and 6) with the real policy service (omos #554), the real MAVLink bridge (omos
# #550) and the onemind-xr scripted drone on a disposable nats-server.
#
# Builds the policy service and the bridge from their git refs (git archive, nothing is
# checked out), builds this service, then runs `go test -tags e2e ./e2e`, which starts
# every process with passwords and an Ed25519 key generated for this run only.
#
# Environment:
#   POLICY_REF    git ref of the policy service   (default origin/feat/policy-service)
#   BRIDGE_REF    git ref of the MAVLink bridge   (default origin/feat/mavlink-bridge-fleet)
#   SIM_DRONE     onemind-xr scripts/sim/mavlink_sim_drone.py (required)
#   PYTHON        python with pymavlink            (default python3)
#   NATS_SERVER   nats-server binary               (default nats-server on PATH)
#   WORKDIR       temp dir (default mktemp -d); TRANSCRIPT JSONL path (default $WORKDIR/transcript.jsonl)
set -euo pipefail

here="$(cd "$(dirname "$0")/.." && pwd)"
repo="$(git -C "$here" rev-parse --show-toplevel)"
workdir="${WORKDIR:-$(mktemp -d)}"
mkdir -p "$workdir/src" "$workdir/bin"
policy_ref="${POLICY_REF:-origin/feat/policy-service}"
bridge_ref="${BRIDGE_REF:-origin/feat/mavlink-bridge-fleet}"
: "${SIM_DRONE:?set SIM_DRONE to onemind-xr scripts/sim/mavlink_sim_drone.py}"

cleanup() { rm -rf "$workdir/js" "$workdir/src"; }
trap cleanup EXIT

rm -rf "$workdir/src/policy" "$workdir/src/bridge"
mkdir -p "$workdir/src/policy" "$workdir/src/bridge"
git -C "$repo" archive "$policy_ref" services/policy services/schemas platform/security/openfga | tar -x -C "$workdir/src/policy"
git -C "$repo" archive "$bridge_ref" services/bridge/bridges/mavlink | tar -x -C "$workdir/src/bridge"
(cd "$workdir/src/policy/services/policy" && go build -o "$workdir/bin/policy" ./cmd/policy)
(cd "$workdir/src/bridge/services/bridge/bridges/mavlink" && go build -o "$workdir/bin/bridge-mavlink" .)
(cd "$here" && go build -o "$workdir/bin/mission" ./cmd/mission)
# Keep the policy inputs outside src/ (removed on exit) for the run.
rm -rf "$workdir/policy-inputs" && mkdir -p "$workdir/policy-inputs"
cp -R "$workdir/src/policy/services/schemas" "$workdir/policy-inputs/schemas"
cp "$workdir/src/policy/platform/security/openfga/onemind-actions.json" "$workdir/policy-inputs/"

export E2E_NATS_SERVER="${NATS_SERVER:-$(command -v nats-server)}"
export E2E_POLICY_BIN="$workdir/bin/policy" E2E_BRIDGE_BIN="$workdir/bin/bridge-mavlink" E2E_MISSION_BIN="$workdir/bin/mission"
export E2E_POLICY_SCHEMAS="$workdir/policy-inputs/schemas" E2E_POLICY_MODEL="$workdir/policy-inputs/onemind-actions.json"
export E2E_SIM_DRONE="$SIM_DRONE" E2E_PYTHON="${PYTHON:-python3}" E2E_WORKDIR="$workdir"
export E2E_TRANSCRIPT="${TRANSCRIPT:-$workdir/transcript.jsonl}"
echo "policy $(git -C "$repo" rev-parse --short "$policy_ref"), bridge $(git -C "$repo" rev-parse --short "$bridge_ref"); workdir $workdir"
cd "$here"
go test -tags e2e -run TestE2E -count=1 -v -timeout 6m ./e2e

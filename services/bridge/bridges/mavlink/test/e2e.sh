#!/usr/bin/env bash
# End-to-end contract test for the MAVLink bridge (FEATURES H2, asset-command.md §9).
#
# Starts a disposable nats-server with JetStream and per-role users (admin, policy,
# bridge, client; passwords generated for this run only, in the temp dir), optionally a
# scripted drone, then runs `go test -tags e2e`. Nothing persists after the run.
#
# Environment:
#   NATS_SERVER   nats-server binary (default: nats-server on PATH)
#   SIM_DRONE     path to onemind-xr scripts/sim/mavlink_sim_drone.py; when set, a drone
#                 is started with MAVLINK20=1 on a free UDP port (default: unset)
#   PYTHON        python with pymavlink for SIM_DRONE (default: python3)
#   MAVLINK_ENDPOINT  when SIM_DRONE is unset: the bridge endpoint of an already running
#                 vehicle (default udpserver:127.0.0.1:14550)
#   OUT_VERSION   MAVLink version the bridge sends (default 2)
#   WORKDIR       temp dir (default: mktemp -d); TRANSCRIPT  JSONL output path
set -euo pipefail

here="$(cd "$(dirname "$0")/.." && pwd)"
workdir="${WORKDIR:-$(mktemp -d)}"
mkdir -p "$workdir"
nats_server="${NATS_SERVER:-nats-server}"
python="${PYTHON:-python3}"
pids=()

cleanup() {
  for p in "${pids[@]}"; do kill "$p" 2>/dev/null || true; done
  wait 2>/dev/null || true
  rm -rf "$workdir/js" "$workdir/nats.conf"
}
trap cleanup EXIT

free_port() { # $1 = tcp|udp
  "$python" -I -c "import socket,sys; s=socket.socket(socket.AF_INET, socket.SOCK_STREAM if sys.argv[1]=='tcp' else socket.SOCK_DGRAM); s.bind(('127.0.0.1',0)); print(s.getsockname()[1])" "$1"
}
genpw() { "$python" -I -c "import secrets; print(secrets.token_urlsafe(18))"; }

nats_port="$(free_port tcp)"
E2E_PASS_ADMIN="$(genpw)"; E2E_PASS_POLICY="$(genpw)"; E2E_PASS_BRIDGE="$(genpw)"; E2E_PASS_CLIENT="$(genpw)"
export E2E_PASS_ADMIN E2E_PASS_POLICY E2E_PASS_BRIDGE E2E_PASS_CLIENT

# Role permissions follow asset-command.md §2: only policy publishes cmd.fleet.>, the
# bridge publishes only its own tel/evt subjects and subscribes only its own commands,
# and a client can never publish cmd.fleet.>.
umask 077
cat > "$workdir/nats.conf" <<EOF
listen: 127.0.0.1:$nats_port
jetstream { store_dir: "$workdir/js", max_memory_store: 64MB, max_file_store: 256MB }
authorization {
  users = [
    { user: admin,  password: "$E2E_PASS_ADMIN" }
    { user: policy, password: "$E2E_PASS_POLICY",
      permissions: { publish: { allow: ["cmd.>", "\$JS.API.>"] }, subscribe: { allow: ["_INBOX.>", "evt.drone.>"] } } }
    { user: bridge, password: "$E2E_PASS_BRIDGE",
      permissions: { publish: { allow: ["tel.fleet.drone.sim-x1.>", "evt.drone.sim-x1.>"] },
                     subscribe: { allow: ["cmd.fleet.drone.sim-x1.>", "_INBOX.>"] } } }
    { user: client, password: "$E2E_PASS_CLIENT",
      permissions: { publish: { allow: ["svc.policy.>", "cmd.map.>"], deny: ["cmd.fleet.>"] }, subscribe: { allow: [">"] } } }
  ]
}
EOF
"$nats_server" -c "$workdir/nats.conf" > "$workdir/nats-server.log" 2>&1 &
pids+=("$!")
export E2E_NATS_URL="nats://127.0.0.1:$nats_port"
for _ in $(seq 1 50); do
  if (echo > "/dev/tcp/127.0.0.1/$nats_port") 2>/dev/null; then break; fi
  sleep 0.1
done

if [[ -n "${SIM_DRONE:-}" ]]; then
  udp_port="$(free_port udp)"
  MAVLINK20=1 "$python" -I "$SIM_DRONE" --out "udpout:127.0.0.1:$udp_port" --lat 0 --lon 0 \
    > "$workdir/sim-drone.log" 2>&1 &
  pids+=("$!")
  export E2E_MAVLINK_ENDPOINT="udpserver:127.0.0.1:$udp_port"
else
  export E2E_MAVLINK_ENDPOINT="${MAVLINK_ENDPOINT:-udpserver:127.0.0.1:14550}"
fi
export E2E_OUT_VERSION="${OUT_VERSION:-2}"
export E2E_TRANSCRIPT="${TRANSCRIPT:-$workdir/transcript.jsonl}"

echo "nats-server $E2E_NATS_URL; vehicle endpoint $E2E_MAVLINK_ENDPOINT; out v$E2E_OUT_VERSION; transcript $E2E_TRANSCRIPT"
cd "$here"
go test -tags e2e -run TestE2E -count=1 -v .

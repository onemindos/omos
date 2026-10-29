# Owner proposals from the MAVLink bridge rework (H2)

Date: 2026-10-29. Status: **proposed**, for the owner. Source PR: `feat/mavlink-bridge-fleet`.

The MAVLink bridge (`services/bridge/bridges/mavlink`) moved from off-grammar `omos.air.*` and
`omos.bridge.*` to existing patterns only. It invented no subject. The items below are the gaps
it hit, written up per onemind-xr decision 0004 #8 ("use the closest existing pattern and write a
proposal in omos `docs/` for the owner"). Each item says what the bridge does today, so nothing
waits on these decisions.

## P1. Fleet telemetry uses the doctrine §9.1 form

`CLAUDE.md` locks `tel.{node}.{class}.{metric}` (example `tel.drone.x47.position`), and
`docs/NATS-DOCTRINE.md` §9.1 (also locked) defines `tel.fleet.{class}.{id}.position|status` and
`cmd.fleet.{class}.{id}.{verb}`. The bridge publishes the §9.1 form, which `platform/nats/streams/tel.yaml`
stores (`tel.fleet.>`) and which the onemind-xr contracts use (asset-command.md §1).

**Proposal:** add one line to the `CLAUDE.md` grammar block saying fleet assets use doctrine §9.1
(`fleet` in the node position, then class and id). No stream change is needed.

## P2. Commands should carry `action_version`

`params_hash` covers `{action_type, action_version, target, params}` (onemind-xr asset-command.md
§5), but the command envelope in §6.3 has no `action_version`. The bridge reads an optional
top-level `action_version` and assumes `1` when it is missing.

**Proposal:** the policy service always sets `action_version` on `cmd.fleet.*` envelopes
(the omos command schema already allows extra members, `additionalProperties: true`).

## P3. Bridge liveness moves to the status subject

The old bridge published `omos.bridge.{name}.health` and `.errors`, which are outside the grammar
and in no stream. Liveness is now `meta.connection` (online/offline) and `meta.last_heartbeat_t`
on `tel.fleet.drone.{id}.status`, stored by `TEL`. The local development dashboard
(`services/bridge/internal/dashboard`) and `cmd/omos` still subscribe `omos.>`, so they no longer
show the MAVLink bridge.

**Proposal:** a follow-up PR points both tools at `tel.fleet.>` and `evt.drone.>`.

## P4. A vehicle that never ACKs reuses `no_ack`

asset-command.md §6.4 has `no_ack` (the policy service saw no bridge acknowledgement) and
`vehicle_rejected` (MAVLink result not 0), but no code for "the bridge sent the command and the
vehicle sent no `COMMAND_ACK`". The bridge reports `outcome: failed, code: no_ack` with
`mavlink.command` set, so the two cases can be told apart.

**Proposal:** either accept the reuse, or add `vehicle_timeout` to the §6.4 list.

## P5. `estop` needs one human, not two

The bridge follows asset-command.md §3: `estop` is A3 (one human). The drones research lane wanted
two people. A stop is also never blocked by link state or geofence in the bridge.

**Proposal:** the H6 reviewer confirms A3 for `estop`.

## P6. A dedicated NATS user for each asset bridge

The end-to-end test runs the bridge as a user that may publish only `tel.fleet.drone.{id}.>` and
`evt.drone.{id}.>` and subscribe only `cmd.fleet.drone.{id}.>` and `_INBOX.>`, and a client user
that the server stops from publishing `cmd.fleet.>` (`test/e2e.sh`). The live bus still uses shared
logins.

**Proposal:** when per-role NATS users land (I4), give each asset bridge these permissions.
Credentials go in a SealedSecret or the edge host's keychain, never in Git.

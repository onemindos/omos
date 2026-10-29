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

Errors before the vehicle could answer are not `vehicle_rejected` either, since §6.4 keeps that
code for a MAVLink result other than 0:

| Bridge case | Reported as |
|---|---|
| No frame from the vehicle sysid yet, so nothing was sent | `failed`, `precondition_failed`, `detail: offline` |
| The same MAV_CMD already waiting for an ACK | `failed`, `conflict` |
| Link write error | `failed`, `no_ack`, `detail: not_sent` |
| Bridge shut down while waiting | `failed`, `no_ack`, `detail: aborted` |

**Proposal:** either accept the reuse, or add `vehicle_timeout` (and `link_error` for the last two
rows) to the §6.4 list.

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

## P7. Task `status_version` from a bridge

mission-task.md §7 sets `Nats-Msg-Id` = `{task_id}:{status_version}`, and §5 rule 1 drops an
assignee report whose `status_version` is not greater than the stored one. The bridge does not know
the version the task service has already stored (`created`, `scheduled`, `sent`), and it keeps its
per-task counter only in memory. Today it starts each task at 1 and adds one per report across
every command of the task, so reports from one bridge process always increase. After a restart the
counter starts again at 1, and within the 2 m `EVT` duplicate window a new report can share a
`Nats-Msg-Id` with an old one and be dropped.

**Proposal (H4):** the task service puts the current `version.status` in the command it dispatches
(for example `task_status_version`), and the bridge reports from that base; or the task service
assigns `status_version` itself and the bridge omits it.

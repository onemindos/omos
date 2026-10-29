# bridge-mavlink

Connects one **simulated** MAVLink vehicle to the OneMind bus on locked-grammar subjects
(FEATURES H2). It publishes fleet telemetry, and it executes only fleet commands that the
policy service signed and humans approved. Built on [gomavlib](https://github.com/bluenviron/gomavlib)
(MAVLink 1 and 2, MIT).

Contracts (onemind-xr): `docs/contracts/asset-command.md` (commands, outcomes, safety rules §7)
and `docs/contracts/robotics-cot-profile.md` §1 and §3 (telemetry fields).

## Subjects

| Direction | Subject | Body | Stream (`platform/nats/streams/`) |
|---|---|---|---|
| out | `tel.fleet.drone.{id}.position` | doctrine §9.2 geo-envelope (`docs/NATS-DOCTRINE.md`) | `TEL` (`tel.fleet.>`) |
| out | `tel.fleet.drone.{id}.status` | §9.2 envelope; `meta` has battery, mode, armed, `connection` (online/offline), kind | `TEL` |
| in | `cmd.fleet.drone.{id}.{verb}` | omos command envelope v1 + asset-command.md §6.3 fields | `CMD` (`cmd.>`) |
| out | `evt.drone.{id}.command` | omos event envelope v1; `data.outcome` = acknowledged / executed / failed / refused | `EVT` (`evt.drone.>`) |
| out | `evt.drone.{id}.task` | mission-task.md §7 status report, only for commands that carry `task_id` | `EVT` |

`{id}` must start with `sim-`; the bridge refuses to start otherwise. Outcome events carry a
`Nats-Msg-Id` header equal to the envelope `id`, so `EVT` de-duplicates them. The old
off-grammar subjects (`omos.air.{id}.telemetry|command`, `omos.bridge.{name}.health|errors`)
are gone; link health is `meta.connection` on the status subject.

## Telemetry

| Field | Source |
|---|---|
| `lat`, `lon` | `GLOBAL_POSITION_INT` degE7 |
| `hae` | `GLOBAL_POSITION_INT.alt` (MSL) **plus the geoid undulation N**: the autopilot's own N from `GPS_RAW_INT` (`alt_ellipsoid - alt`, 3D fix, last 5 s) when sent, else the embedded EGM96 grid (`internal/geoid`, 0.5 degree, at most 1.95 m from the 15-minute grid). `meta.hae_source` says which; `meta.alt_msl_m` keeps the original |
| `course`, `speed` | course and speed over ground from `GLOBAL_POSITION_INT` vx/vy; course `-1` below 0.5 m/s |
| `meta.mode` | `HEARTBEAT.custom_mode` as an ArduCopter mode name (gomavlib `COPTER_MODE` enum) for ArduPilot multirotors, else `custom:N` |
| `meta.battery_pct` | `SYS_STATUS` or `BATTERY_STATUS` `battery_remaining` |
| `meta.kind` | `HEARTBEAT.type` mapped per robotics-cot-profile.md §2 |
| `meta.link_rssi` | `-1`: `RADIO_STATUS.rssi` units are radio-specific and are not mapped |

Position is published at most every `POSITION_INTERVAL` (200 ms), status every `STATUS_INTERVAL` (1 s).

## Commands

The bridge holds a **core** subscription on `cmd.fleet.drone.{id}.>` (no durable consumer, so an
offline bridge never receives a stale command later, asset-command.md §4). Each command must pass
every check below, in this order, or it is refused with the code shown and nothing reaches the vehicle.

| Check | Refusal code |
|---|---|
| Target id starts `sim-`, equals this bridge's id; `simulated: true` | `not_simulated` / `schema_invalid` |
| Verb in the catalogue below | `unknown_action` |
| Envelope: `v`, UUID `id`, `cmd` = `subject` = delivery subject, `action_id`, `ts`, `expires_at` | `schema_invalid` |
| `source` = `service:policy` | `not_authorized` |
| ed25519 `sig` over the RFC 8785 (JCS) envelope without `sig`, key by `kid` from `POLICY_PUBKEYS` | `bad_signature` |
| Not expired; `ts` not in the future; `expires_at - ts` ≤ `CMD_MAX_TTL` | `expired` |
| `policy.allowed` with a `decision_id` | `not_authorized` |
| `approval`: required, `approved`, level ≥ the verb's level, enough **distinct `user:`** approvers (A3: 1, A4: 2), approver listed | `not_authorized` |
| Any approval by a non-`user:` principal (an agent) | `agent_cannot_confirm` |
| Parameters valid for the verb (unknown members refused) | `schema_invalid` |
| `approval.params_hash` = sha256 of JCS `{action_type, action_version, target, params}` | `not_authorized` |
| Command `id` not seen before (nonce, kept until `expires_at`) | `replayed` |
| Vehicle heartbeat within `LINK_TIMEOUT` (not for `estop`) | `precondition_failed` / `offline` |
| Own geofence: target (goto, land at a point) or current position (arm, takeoff) inside a keep-in zone, outside every keep-out zone, within the altitude band. No keep-in configured = refused (fail closed) | `precondition_failed` / `keep_in`, `keep_out` |

On success it publishes `acknowledged`, sends the MAVLink command, waits for `COMMAND_ACK`
(`ACK_TIMEOUT` × `ACK_RETRIES`, COMMAND_LONG `confirmation` incremented per retry,
`MAV_RESULT_IN_PROGRESS` extends the wait), then publishes `executed` (result 0), `failed` /
`vehicle_rejected` (other result, with `mavlink.result_name`) or `failed` / `no_ack`.
`estop` never waits behind another command and is not blocked by link state.

| Verb | Level | MAVLink | Notes |
|---|---|---|---|
| `goto {lat, lon, hae_m, speed_mps?, yaw_deg?}` | A3 | `MAV_CMD_DO_REPOSITION` (192), `COMMAND_INT`, frame `MAV_FRAME_GLOBAL` | z = `hae_m` − N (MSL); param1 speed or −1; param2 `MAV_DO_REPOSITION_FLAGS_CHANGE_MODE`; param4 yaw or NaN |
| `rtl {}` | A3 | `MAV_CMD_NAV_RETURN_TO_LAUNCH` (20) | |
| `land {lat?, lon?}` | A3 | `MAV_CMD_NAV_LAND` (21) | `COMMAND_INT` at a point, else `COMMAND_LONG` |
| `arm {}` / `disarm {}` | A4 / A3 | `MAV_CMD_COMPONENT_ARM_DISARM` (400) | param1 1 / 0; never the force value |
| `estop {}` | A3 | `MAV_CMD_DO_FLIGHTTERMINATION` (185) | param1 1 |
| `takeoff {hae_m}` | A4 | `MAV_CMD_NAV_TAKEOFF` (22) | param7 = `hae_m` − N; how each autopilot reads param7 is **unverified** |

Levels and verbs follow asset-command.md §3. Every other verb is refused.

**Signing.** Until the policy service publishes its key, `SIG_REQUIRED=false` lets unsigned
commands through **for local development only**, with a warning per command; a present but
invalid signature is still refused. All other checks always apply.

## Configuration (environment)

| Variable | Default | Meaning |
|---|---|---|
| `NATS_URL` | `nats://127.0.0.1:4222` | |
| `NATS_USER` / `NATS_PASSWORD` or `NATS_CREDS` | — | bridge identity; from the environment or keychain, never Git |
| `BRIDGE_NAME` | `bridge-mavlink` | `source` of outcome events |
| `ASSET_ID` (or legacy `BRIDGE_ENTITY_ID`) | — | **required**, must start `sim-` |
| `CALLSIGN` | upper-case id | |
| `MAV_SYSID` / `MAV_GCS_SYSID` | `1` / `254` | vehicle system id / this node's id |
| `MAVLINK_ENDPOINT` | `udpserver:127.0.0.1:14550` | `udpserver:`, `udpclient:`, `tcpclient:`, `tcpserver:` host:port or `serial:/dev/x:baud`. Legacy `DRONE_IP`/`DRONE_PORT` = `udpclient` |
| `MAVLINK_OUT_VERSION` | `2` | `1` only for a MAVLink 1 peer (the scripted drone without `MAVLINK20=1`) |
| `POLICY_PUBKEYS` | — | `kid=base64-ed25519-public-key[,…]` |
| `SIG_REQUIRED` | `true` | see Signing |
| `GEOFENCE_FILE` or `GEOFENCE_JSON` | — | see `geofence.example.json` (`keep_in`, `keep_out` circles or polygons, `max_hae_m`, `min_hae_m`) |
| `CMD_MAX_TTL`, `CLOCK_SKEW` | `60s`, `5s` | |
| `POSITION_INTERVAL`, `STATUS_INTERVAL`, `STALE_S`, `LINK_TIMEOUT` | `200ms`, `1s`, `10`, `3s` | |
| `ACK_TIMEOUT`, `ACK_RETRIES` | `1500ms`, `3` | |

## Run against the scripted drone

```bash
# onemind-xr scripted drone, MAVLink 2 (it sends to udp 14550)
MAVLINK20=1 python3 onemind-xr/scripts/sim/mavlink_sim_drone.py --lat 0 --lon 0 &
go build -o bridge-mavlink .
ASSET_ID=sim-x1 GEOFENCE_FILE=geofence.example.json POLICY_PUBKEYS="policy-2026-10=<base64>" ./bridge-mavlink
nats sub 'tel.fleet.drone.sim-x1.>'
nats sub 'evt.drone.sim-x1.>'
```

## Tests

```bash
go vet ./... && go test -race ./...          # unit + in-process vehicle (CI: test-mavlink-bridge.yml)
NATS_SERVER=nats-server PYTHON=python3 \
SIM_DRONE=../../../../../onemind-xr/scripts/sim/mavlink_sim_drone.py \
GOTO_OUTCOME="acknowledged,failed/no_ack" ./test/e2e.sh
```

`test/e2e.sh` starts a disposable nats-server with JetStream and per-role users (passwords
generated per run), a MAVLink 2 scripted drone, and the bridge in-process, then checks: telemetry
appears and the `TEL`-equivalent stream stores every message; approved arm, takeoff, rtl, land,
disarm are ACKed; unapproved, one-human arm, agent-approved, non-simulated, forged, unsigned,
expired, replayed, keep-out and unknown-verb commands are refused; a non-`sim-` target gets no
answer; and a client user cannot publish `cmd.fleet.>`.

The scripted drone answers only `COMMAND_LONG` for arm, disarm, takeoff, land and RTL, so `goto`
(`COMMAND_INT` `DO_REPOSITION`) ends `failed` / `no_ack` there and `estop` would be
`vehicle_rejected`; the in-process test vehicle covers both. PX4 SITL or ArduPilot SITL is the
owner check for real autopilot behaviour.

# bridge-ros2: ROS 2 robots on the OneMind bus (FEATURES H3)

Connects one simulated ROS 2 ground robot (TurtleBot-class, robot dog, rover) to NATS:
odometry, pose and battery become fleet telemetry, and only policy-issued, human-approved,
signed commands move the robot. Same contract and same verification rules as the MAVLink
bridge (omos PR #550).

```
robot (ROS 2 / DDS) ── zenoh-bridge-ros2dds ── Zenoh ── bridge-ros2 ── NATS
                     (beside the robot, router mode)   (this adapter)   tel.fleet.robot.*, evt.robot.*
                                                                        ◄ cmd.fleet.robot.* (policy only)
fallback: robot ── rosbridge_server (WebSocket JSON) ── bridge-ros2 ── NATS
```

Contracts (onemind-xr, branch `run/m0`): `docs/contracts/asset-command.md`,
`docs/contracts/robotics-cot-profile.md`, `docs/design/BUS-DESIGN.md` section 3. Research:
`docs/research/drones-robots-agents.md` ("ROS 2 to NATS", How it applies 6: "zenoh-python to
nats-py").

## Subjects (locked grammar; no new patterns)

| Direction | Subject | Payload |
|---|---|---|
| publish | `tel.fleet.robot.{id}.position` | doctrine section 9.2 geo-envelope, `type: fleet.robot`, `meta.kind` (`ugv`), heading, battery, mode, map x/y |
| publish | `tel.fleet.robot.{id}.status` | same envelope; `meta.connection` online/offline, battery, voltage, e-stop latch, transport stats, goal and home |
| publish | `evt.robot.{id}.command` | omos event envelope v1; `data.outcome` acknowledged / executed / failed / refused with a section 6.4 code (stored by `EVT`, which binds `evt.robot.>`) |
| publish | `evt.robot.{id}.task` | mission-task.md section 7 task status for commands that carry `task_id`; `Nats-Msg-Id = {task_id}:{status_version}` |
| subscribe | `cmd.fleet.robot.{id}.>` | **core** subscription: an offline bridge never receives a stale command later |

## Verbs

| Verb | Level | Zenoh path (default) | rosbridge fallback |
|---|---|---|---|
| `goto {lat, lon, hae_m, speed_mps?, yaw_deg?}` | A3, one human | NavigateToPose goal on `{ns}/navigate_to_pose/_action/send_goal`; `executed` when Nav2 accepts | `goal_pose` topic, no acknowledgement |
| `hold {duration_s?, reason?}` | A3 | `cancel_goal` (all goals) and a zero `cmd_vel` | goal re-targeted to the current pose and a zero `cmd_vel` |
| `rtl {}` | A3 | goal at home (`HOME_LAT/HOME_LON`, or the first position heard) | same via `goal_pose` |
| `estop {}` | A3, never blocked by a precondition | zero `cmd_vel` three times, then cancel all goals; **latches** | zero `cmd_vel`, re-target `goal_pose` to the current pose, `call_service` `{ns}/navigate_to_pose/_action/cancel_goal` (all goals, no reply awaited), zero `cmd_vel` twice more; latches |
| `estop_release {}` | A4, two distinct humans | clears the latch only; the robot stays stopped | same |

`hae_m` is required by the action schema and ignored by a ground robot. `speed_mps` is
accepted and reported as not applied (NavigateToPose has no per-goal speed). Any other verb
(`arm`, `takeoff`, `strike`, ...) is refused `unknown_action`.

## Checks before the robot moves

`ros2bridge/command.py` is a port of `services/bridge/bridges/mavlink/internal/command`
(#550), check for check and code for code: subject is `cmd.fleet.robot.{own id}.{verb}`; the id
and `simulated` say sim; `source` is `service:policy`; Ed25519 signature over the RFC 8785 JCS
envelope without `sig` from a key in `POLICY_PUBKEYS`; `ts`/`expires_at` fresh and at most
`CMD_MAX_TTL` apart; allowing policy decision; `approval.status = approved` at a level at or
above the verb's, with enough distinct `user:` approvers (an agent approver is
`agent_cannot_confirm`); `params_hash` equals the `asset.{verb}` hash of the delivered `data` when there is no `task_id`; a task-bearing command carries the `task.assign`/`task.abort` hash and is checked for form only, the signature binding it (D2-P1; shared vectors `services/schemas/command/testdata/params-hash-vectors.json`, `tests/test_params_hash_vectors.py`); unseen command id.
Then the bridge's own checks: robot heard within `LINK_TIMEOUT` (not for `estop`), e-stop latch,
and its own keep-in/keep-out geofence (the #550 file format, fail closed with no keep-in zone)
on the target and on the straight line from the robot to it (keep-out only; a zone the robot is
already in is skipped so it can drive out). Nav2 plans its own path, so the line check is a
first pass, not a guarantee (owner proposal H3-P7).

An `estop` overtakes movement already under way. Every executed `estop` bumps a counter that
each command captures on arrival. A `goto` or `rtl` that was waiting on the previous command,
or on Nav2's answer, when an `estop` ran is not sent; if Nav2 had already accepted it, the bridge
stops the robot again (zero `cmd_vel` and cancel). Either way the command ends `failed`
`precondition_failed` / `estopped`, records no goal and leaves the mode at `estop`.

The Go verifier cannot be imported here (Go `internal` package of another module, and the
adapter is Python because Zenoh's Go binding needs zenoh-c via cgo). To stop the copies
drifting, `tests/test_command.py` and `tests/test_jcs.py` verify envelopes produced by the
policy service's own `Signer`, `ParamsHash` and `Canonicalize` (`tests/fixtures/`), and the
end-to-end test runs the real policy binary. Consolidation is owner proposal H3-P2.

## Map frame to WGS-84

ROS poses are in a local frame. Configure the frame origin and the compass bearing of its
+x axis (`MAP_ORIGIN_LAT`, `MAP_ORIGIN_LON`, `MAP_ORIGIN_HAE`, `MAP_X_BEARING_DEG`, default 90 =
ENU). Positions come from `odom`, or from `POSE_TOPIC` (e.g. `amcl_pose`, map frame) when set;
use the latter on a real robot, because `odom` drifts from `map`. Goals are sent in
`MAP_FRAME` (default `map`). `hae` is the origin's height plus the pose's z (`meta.hae_source:
map_origin`).

## Configuration (environment only)

| Variable | Default | Meaning |
|---|---|---|
| `NATS_URL`, `NATS_USER`, `NATS_PASSWORD` or `NATS_CREDS` | `nats://127.0.0.1:4222` | the bridge's own least-privilege user (H3-P5); password from the owner keychain or edge `.env`, never Git |
| `ASSET_ID` | required | must start `sim-` |
| `CALLSIGN`, `ROBOT_KIND`, `BRIDGE_NAME` | `SIM-…`, `ugv`, `bridge-ros2` | |
| `TRANSPORT` | `zenoh` | `zenoh` or `rosbridge` |
| `ZENOH_CONNECT`, `ZENOH_MODE`, `ZENOH_LISTEN` | `tcp/127.0.0.1:7447`, `client` | where zenoh-bridge-ros2dds listens |
| `ROS_NAMESPACE` | empty | the ros2dds `namespace` (key prefix) or the rosbridge topic namespace |
| `ROSBRIDGE_URL` | `ws://127.0.0.1:9090` | fallback only |
| `ODOM_TOPIC`, `POSE_TOPIC`, `BATTERY_TOPIC`, `CMD_VEL_TOPIC`, `NAV_ACTION`, `GOAL_TOPIC` | `odom`, empty, `battery_state`, `cmd_vel`, `navigate_to_pose`, `goal_pose` | |
| `MAP_ORIGIN_LAT`, `MAP_ORIGIN_LON` | required | frame origin |
| `MAP_ORIGIN_HAE`, `MAP_X_BEARING_DEG`, `MAP_FRAME` | `0`, `90`, `map` | |
| `HOME_LAT`, `HOME_LON` | first position | `rtl` target |
| `POLICY_PUBKEYS` | required | `kid=base64(raw Ed25519 public key)`, comma-separated |
| `SIG_REQUIRED` | `true` | `false` only for local development |
| `GEOFENCE_FILE` / `GEOFENCE_JSON` | none (all movement refused) | see `geofence.example.json` |
| `CMD_MAX_TTL`, `CLOCK_SKEW`, `LINK_TIMEOUT`, `ACTION_TIMEOUT` | `60s`, `5s`, `3s`, `3s` | |
| `POSITION_INTERVAL`, `STATUS_INTERVAL`, `STALE_S`, `GOAL_TOLERANCE_M` | `500ms`, `1s`, `10`, `0.3` | |

## Run

```
pip install -r requirements.txt
# on the robot host: zenoh-bridge-ros2dds -c zenoh-bridge-ros2dds.json5
ASSET_ID=sim-tb1 ROS_NAMESPACE=sim_tb1 ZENOH_CONNECT=tcp/<robot-host>:7447 \
MAP_ORIGIN_LAT=... MAP_ORIGIN_LON=... GEOFENCE_FILE=geofence.json \
POLICY_PUBKEYS="policy-1=<base64>" NATS_URL=... NATS_USER=... NATS_PASSWORD=... \
python -m ros2bridge
```

Without ROS: `python sim/turtlebot_sim.py --mode router --listen tcp/127.0.0.1:7447 --namespace sim_tb1`
stands in for ROS 2 plus zenoh-bridge-ros2dds (CDR `odom`, `battery_state`, the
NavigateToPose queryables and `cmd_vel`).

## Tests

```
pip install -r requirements-dev.txt
python -m pytest tests            # unit + Zenoh/rosbridge in-process (CI: test-ros2-bridge.yml)
tests/e2e/run.sh                  # + nats-server, policy service (PR #554), robot sim, bridge
```

## Deployment

Edge only, like the MAVLink bridge: zenoh-bridge-ros2dds has to sit on the robot's DDS
network, and the adapter runs beside it (Mac or home-server site, behind the site's NATS
leaf). There is no image, Argo Application or cluster Secret; see H3-P6 in
`docs/proposals/2026-10-29-ros2-bridge-h3.md`.

## Unverified (owner check in onemind-xr `docs/evidence/M4/H3/OWNER-STEPS.md`)

- That zenoh-bridge-ros2dds 1.10.2 exposes NavigateToPose as `{ns}/navigate_to_pose/_action/send_goal`
  and `.../cancel_goal` with plain CDR request/reply payloads (no DDS request header). The
  README and DEFAULT_CONFIG of the plugin state that actions map to queryables and the
  namespace prefix, not these key names or payload framing.
- That topic payloads on Zenoh are the unmodified CDR buffers (with the 4-byte encapsulation header).
- Gazebo TurtleBot3 + Nav2 behaviour (goal acceptance, arrival) against this adapter.

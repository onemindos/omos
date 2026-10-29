# Integration-pass resolutions INT-R1 … INT-R4

Status: **proposed, for the owner**. Date: 2026-10-29. Filed under onemind-xr decision 0004
item 8. Wiring the M4 PRs together (onemind-xr `docs/STATUS.md`, M4 row) found five cross-PR
blockers. Four of them were proposed several times, in different words, by different PRs. Each
entry below combines those proposals into one ruling and names the PR that carries the change
(the fifth blocker, `zone.*` in the policy catalogue, is D1-P4, a single patch to #554). Ruling
an INT entry also rules every proposal listed in "Rules with". None of them adds a subject
pattern.

The code that applies a resolution lives in the carrying PR, not here. When this file was
written (2026-10-29), none of the carrying PRs had the change on its head yet. Check the PR
before you tick a box. The fix lane's evidence is in onemind-xr `docs/evidence/INT/` (branch
`run/m0`).

| ID | Resolution | Rules with | Carried by | Tier |
|---|---|---|---|---|
| INT-R1 | One canonical `params_hash`, computed over the action the message executes | D2-P1, H3-P2 | #554 (policy), checked by #550, #566, clients | 1 |
| INT-R2 | The entity service fills `capabilities` on fleet entities from the bridge status | WEB-P5, H3-P4 | #552 (entity) | 1 |
| INT-R3 | One consolidated `CLIENT` import list | WEB-P2, UE-P3, I2-P2, A5-P7, D2-P5 (exports), K-P7 | #562 (identities) | 1 |
| INT-R4 | The projector serves `svc.cot.send`; one body, kinds `geochat` and `shape` | WEB-P4, UE-P1, I1-P2 (exit), I1-P1 (publisher) | #563 (projector) | 3 |

---

## INT-R1. `params_hash` canonical form

**Problem.** #554 stamped the fleet command of a `task.assign` with the hash of the
`task.assign` action. #550 recomputes the hash as `asset.{verb}` over the command `data`, so it
refused every command that a task dispatched (`not_authorized`). This was reproduced end to end
in onemind-xr `docs/evidence/M4/D2/e2e-unpatched.txt`. Three clients (web, Unreal, Hermes) and
two bridges (MAVLink, ROS 2) each compute the hash, and the only shared reference is #554's
code.

**Resolution.**
1. `params_hash` is `sha256:` followed by the lowercase hex SHA-256 of the RFC 8785 (JCS)
   canonical JSON of `{action_type, action_version, target, params}`. This is the formula in
   onemind-xr `contracts/asset-command.md` §6.1, unchanged.
2. The hash always describes the action that the message executes. An approval binds
   `params_hash` of the approved action, such as `task.assign`, and that binding stays in
   `EVT`. Every asset command on `cmd.fleet.>`, including one dispatched for a task, carries
   the hash of `asset.{verb}` over its own `data`. The policy signature still covers the whole
   command envelope, so a bridge can check a command without knowing the task's parameters.
   This is the 12-line change in #568's
   `services/mission/test/interop/policy-task-fleet-params-hash.patch`, applied to #554.
3. The reference implementation is #554 `services/policy/internal/policy/jcs.go`. Its golden
   values are in onemind-xr `docs/evidence/M4/UE/jcs-golden.txt` and
   `docs/evidence/M4/WEB/jcs-vectors.txt`, and the Unreal and web clients already match them.
   The vectors move to `services/schemas/command/vectors/` (H3-P2). The policy service, both
   bridges and the clients test against that one file.

**Not chosen.** #550 accepting `task.assign` hashes: the bridge cannot recompute them without
the task's parameters.

- [ ] Approve as proposed
- [ ] Approve with changes:
- [ ] Reject

## INT-R2. `capabilities` on fleet entities

**Problem.** The policy service denies an asset verb (`not_capable`) unless the target entity
lists the matching task type. The entity service (#552) never sets `capabilities`, so in the
cluster every `goto` would be denied. The web, Unreal, ROS 2 and mission e2e runs each used an
entity stand-in to get past this. Meanwhile the bridges already publish the verbs they accept:
#550 lists `commands_accepted` in its status `meta`, and #566 does the same.

**Resolution.** The entity service maps each verb in the latest fleet status
`meta.commands_accepted` to `capabilities: ["onemind.task.{verb}.v1", …]` on that fleet
entity. It does not invent a capability the bridge did not list. With no status, the entity
gets no capabilities, so the policy service fails closed. A later registry entry may narrow
the list but may not widen it. The policy service is unchanged. Clients offer a verb only when
the entity lists it (the WEB follow-up). This also covers the robot `meta.kind` reading
(H3-P4).

**Not chosen.** Making the policy service read `commands_accepted` directly. That would give
the policy service a second source of entity facts, against the single-writer rule
(REFERENCE-ARCHITECTURE R3).

- [ ] Approve as proposed
- [ ] Approve with changes:
- [ ] Reject

## INT-R3. One `CLIENT` import list

**Problem.** #562 gives a `CLIENT` user requests on `svc.policy.>`, `svc.entity.>`,
`svc.layer.>`, `svc.cot.send` and `agents.prompt.hermes.>`, plus the `cmd.map.>` stream from
`AGENTS`. It gives nothing else. Five lanes each asked for more, separately: WEB-P2, UE-P3,
I2-P2, A5-P7 and D2-P5. K-P7 asked for `svc.registry.>`, which BUS-DESIGN §7.1 lists but #562
leaves out.

**Resolution.** The list below goes into #562's `CLIENT` and `AGENTS` blocks as one change.

| Import into `CLIENT` | Kind | From | Asked by | Note |
|---|---|---|---|---|
| `svc.layer.>` | service, `response_type: stream`, `response_threshold` ≥ 2 min, `share: true` | `ONEMIND` | A5-P7 | changes the existing import's export settings |
| `svc.registry.>` | service, `share: true` | `ONEMIND` | K-P7 | video registry reads |
| `svc.mission.{get\|list\|report}` | service, `share: true` | `ONEMIND` | D2-P5 | `svc.mission.apply` stays policy-only |
| `svc.geofence.alerts` | service, `share: true` | `ONEMIND` | D1-P3 | open-alert list |
| `tel.feed.>` | stream, read only | `ONEMIND` | WEB-P2 (b), UE-P3 | feed heartbeats for the gaps panel (B20); no sensitive data |
| `cmd.map.>` | **service** with `response_type: stream`, exported by `CLIENT` to `AGENTS` | `CLIENT` → `AGENTS` | I2-P2 | replaces the stream export, so client acks reach Hermes |

Two items stay open for the owner:

- **Live action status (`evt.so.action.>`).** A stream import would show every user's actions
  to every `CLIENT` user, and per-user filtering is unverified (WEB-P2 (a)). Until the owner
  rules, clients poll `svc.policy.action.get`, which works today. UE-P3 asks for the stream.
- **Chat receive.** Choose between a `CLIENT` import of `tak.cot.out` and `tak.cot.mesh.rx`,
  which exposes all TAK CoT to authenticated humans, and a `svc.cot.chat` history endpoint
  served by the projector (WEB-P2 (c)).

Two further imports wait for features that are not built: `svc.tak.mission.>` and
`evt.mission.*.changed` need UE-P2 implemented in tak-client.

- [ ] Approve the table; action status by polling; chat by: ☐ stream import ☐ `svc.cot.chat`
- [ ] Approve with changes:
- [ ] Reject

## INT-R4. `svc.cot.send`

**Problem.** onemind-xr `contracts/cot-projection.md` §9 names `svc.cot.send` and leaves its
body to FEATURES C2. #562 already imports it into `CLIENT`, but no service serves it. The
check covered #560 and #563. Meanwhile three things fill the gap differently. The web client
sends a GeoChat body (WEB-P4). The Unreal client adds shapes (UE-P1). Agentops publishes CoT
itself on the TAK outbound subjects (I1-P1, I1-P2).

**Resolution.** The projector (#563) serves `svc.cot.send` and is the only service that builds
outbound CoT, so the marking gate G7 and the echo convention apply to everything sent.
- Request: `{v: 1, kind: "geochat" | "shape", target: "server" | "mesh" | "both", message_uid, …}`.
  - `geochat` carries `room {id, name, kind}`, `text` (at most 1000 characters) and
    `sender {uid, callsign}` (WEB-P4).
  - `shape` carries the B11 types `b-m-p-s-m`, `u-d-f`, `u-d-r`, `u-d-c-c`, `u-r-b-c-c` and
    `u-rb-a`, with geometry (UE-P1).
- Reply: `{ok: true, uid}` or `{ok: false, error: {code, message}}`.
- The projector takes the sender's identity from the connection (`Nats-Request-Info`). The
  body's `sender` is only the display callsign. The client never sends XML.
- Output goes on `tak.cot.inject` for `server` and `tak.cot.mesh.tx` for `mesh`, with a core
  publish and no reply (CP-3).
- When the endpoint is live, agentops moves its contact PLI and its GeoChat replies onto it
  and loses its publish grant on the two TAK subjects (I1-P2 exit). The client's own TAK
  stream is a separate decision (UE-P5).

- [ ] Approve as proposed
- [ ] Approve with changes:
- [ ] Reject

```subjects
filter    cmd.fleet.>
filter    evt.so.action.>
current   svc.policy.action.get
filter    svc.policy.>
filter    svc.entity.>
filter    svc.layer.>
filter    svc.registry.>
proposed  svc.mission.{get|list|report}
current   svc.mission.apply
proposed  svc.geofence.alerts
filter    tel.feed.>
filter    cmd.map.>
filter    agents.prompt.hermes.>
proposed  svc.cot.send
alternative  svc.cot.chat
current   tak.cot.out
current   tak.cot.mesh.rx
current   tak.cot.inject
current   tak.cot.mesh.tx
filter    svc.tak.mission.>
filter    evt.mission.*.changed
action    task.assign
action    asset.{verb}
```

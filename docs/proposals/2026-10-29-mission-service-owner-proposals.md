# Owner proposals from the task service (FEATURES D2, H4)

Status: **proposals**, 2026-10-29, filed under onemind-xr decision 0004 #8. The task service
(`services/mission`) uses **no new subject pattern**: every subject is in the locked grammar
(omos `CLAUDE.md`) and in onemind-xr `docs/contracts/mission-task.md` §1. Grammar and stream
items are ruled once in the A24 register (omos #556, `docs/proposals/2026-10-29-grammar-register/`);
this file references its IDs and keeps only service-level items, as the register's
cross-reference asks. The register's cross-reference table needs one row per item below
when #556 is next updated (not edited here: #556 is another lane's PR).

## Subjects used

| Subject or key | Grammar line | Stored by | Register |
|---|---|---|---|
| `cmd.mission.{id}.{create\|waypoint\|abort\|complete}` (read) | `cmd.{domain}.{entity}.{action}`; doctrine §9.1 verbs | `CMD` | BUS-P2 (no `cancel`; abort is used) |
| `evt.mission.{id}.{event}` (write) | `evt.{domain}.{entity}.{event}` | `MISSION` | ONT-P11 (asset tasks in `mission`, never `evt.task.>`) |
| `evt.so.action.dispatched` (read) | `evt.{domain}.{entity}.{event}` | `EVT` | ONT-P10 |
| `evt.drone.{id}.task`, `evt.robot.{id}.task` (read) | `evt.{domain}.{entity}.{event}` | `EVT` | ONT-P13 for `evt.vehicle.>` |
| `tel.fleet.{class}.{id}.position` (read) | doctrine §9.1 | `TEL` | BUS-P1 |
| `svc.mission.get`, `svc.mission.list`, `svc.mission.report`, `svc.mission.apply` (serve) | `svc.{service}.{endpoint}` | — | new names inside the grammar (mission-task.md §1 "proposed"); see D2-P4 |
| KV `mission-state` (write) | — | bucket exists | — |

## Proposals

| ID | Gap | What the service does now | Proposed ruling |
|---|---|---|---|
| D2-P1 (resolved: POL-16) | #554 stamps the fleet command of a `task.assign` with the `task.assign` `params_hash`; #550 recomputed `params_hash` as `asset.{verb}` over the command `data` and refused every task-dispatched command (`not_authorized`). | Records any refusal faithfully (`done_not_ok` + `rejected`) | **Resolved by POL-16** (`docs/proposals/2026-10-29-policy-service-owner-proposals.md` in #554): the command carries the approved action's own hash (asset-command.md §4, §5, §6.3), so #554 is unchanged and the bridges (#550 `2b724093`, #566) recompute only without `task_id` and otherwise check the `sha256:` form and rely on the Ed25519 signature. The re-stamp patch this row first proposed is withdrawn and removed from this PR. The e2e run (`services/mission/test/e2e.sh`, default) flies the drone and the task ends `done_ok`. Shared vectors: `services/schemas/command/testdata/params-hash-vectors.json` |
| D2-P2 (owner decision) | The bridge reported `executing` when the vehicle accepted a `goto` and never `done_ok`, so with `MISSION_ARRIVAL=false` (the deployment value) a goto task never left `executing` (onemind-xr `docs/evidence/INT/E2E-2/run2-arrival-off`) | **Bridge-side arrival, now implemented**: #550 `6e4f29ea` reports `done_ok` on `evt.drone.{id}.task` once the vehicle is within `ARRIVAL_RADIUS_M` 5 m horizontally and `ARRIVAL_VERTICAL_M` 5 m vertically of the goto target (MSL, same geoid as the command), `result = {arrived, source: bridge, distance_m, vertical_m}`, causation = the goto command id; a later accepted command supersedes the goto. The task service's telemetry check stays behind `MISSION_ARRIVAL` (default and manifest `false`) as a fallback for a bridge without arrival reports. The default config reaches `done_ok` (e2e, scripted drone) | **Ratify bridge-side arrival** as the rule. Why it is the cleaner of the two: mission-task.md §5 makes `done_ok` the assignee's state; the bridge commands the vehicle and knows the target in the vehicle's own frame (no HAE/MSL re-conversion), when a newer command supersedes the goto, and (later) the autopilot's own acceptance radius; the task service's check trusts any `tel.fleet.*.*.position` publisher and judges a goto it never commanded. Decide also: (1) the 5 m / 5 m default; (2) whether to keep `MISSION_ARRIVAL` as a fallback until the ROS 2 bridge (#566) reports arrival, or remove it. Follow-up: #566 reports `done_ok` on arrival the same way |
| D2-P3 | #550 P7: the bridge numbers reports per task from 1 for its process life and does not know `version.status`; a restarted bridge starts again at 1 | `version.status` is the task service's own counter (fact `status_version`, `Nats-Msg-Id` `{task_id}:{version.status}`, so MISSION dedup never collides). Reports are ordered by the additive record field `report_guard`: (index of the report's `causation` = asset command id among the task's commands, reporter `status_version`). A newer command wins even after a restart; inside one command the version must increase | Adopt P7's option "the task service assigns `status_version`" and amend mission-task.md §5 rule 1 and §7: the reporter's `status_version` is compared within one asset command; `status_version` on `evt.mission.{tsk}.*` is `version.status`. No bridge change needed |
| D2-P4 | Bodies of `svc.mission.*`, mission event names and the refusal fact are not fixed by the contract | `get {v,id}`; `list {kind?,state?,assignee?,parent_id?,bbox?,limit?}`; `report` = §7 data; `apply {v,action_id,op,mission_id,to,role?,by?}`; mission events `created`, `waypoint`, `updated`, `assigned`, `unassigned`, `completed`, `aborted`; refused commands `evt.mission.{id}.rejected` with `{verb,code,message,action_id,command_id}` and `causation` = command id (the policy service maps `rejected` to `failed`, POL-14) | Ratify in mission-task.md §1 and §6; add the `svc.mission.*` names to the A24 register as "inside the grammar" |
| D2-P5 | No NATS identity for the task service in A19 (#562); #562 gives the `policy` user `kv_write` on `mission-state`, but mission-task.md §2 makes the task service the only writer | Tests run as a least-privilege `mission` user (below); any other grant fails the test | Add the user below to `identities.yaml`; export `svc.mission.get/list/report` to CLIENT and AGENTS (`share: true`); give `policy` publish on `svc.mission.apply` and drop its `mission-state` write grant (BUS-DESIGN §7 row to be amended) |
| D2-P6 | `cmd.mission.{msn}.abort` aborts the children (§3.1), but nothing sends their asset holds | Marks sent children `cancel_requested`, cancels unsent ones | The policy service fans out one `cmd.fleet.{class}.{id}.hold` (with `task_id`) per open asset child when it dispatches `mission.abort` |
| D2-P7 | For an asset task the policy service waits only for `evt.mission.{tsk}.created`; a `rejected` result makes it wait the full 5 s and fail `no_ack` | Publishes `rejected` with `code` | #554 `handleCreated` also ends the wait on `evt.mission.{tsk}.rejected` and fails with its `code` |
| D2-P8 | `cmd.mission.*` carries the policy Ed25519 signature; the task service does not verify it | Relies on A19 stage B (only `policy` publishes `cmd.mission.>`); it never sends an asset command itself | Follow-up: verify the signature with the bridges' rule once the JCS code is shared |
| D2-P9 | §6.1 step 9 "entity `current_task` follows" | Not done (the entity service owns `ent.*`) | The entity service (#552) consumes `evt.mission.*` and sets `current_task` |
| D2-P10 | Composite task types (§4) | Accepted as records, not decomposed; the policy service refuses them (`unknown_action`) | H5: decomposition in the task service, one approval per child |
| D2-P11 | `platform/parked/omos-mission-engine/` is a never-built placeholder (`:latest`, 0 replicas) | This service is the task service of the contract | Delete the parked placeholder when this PR merges |
| D2-P12 | Defaults the contract leaves open | Mission `marking` `AMBER`; `delivery.before_t` creation + 60 s; `stale_t` = `execution.complete_before_t`, else creation + 1 h; `mission.create` carries no requester, so `created_by` = the command source (`service:policy`) and `updated_by` = the approver | Ratify, or name other defaults |
| D2-P13 | BUS-DESIGN §6 says "after a reconnect a task is offered again" | Never re-sends (rule 4, asset-command.md "No queued replay") | Amend BUS-DESIGN to match the contract |
| D2-P14 | The asset's `received` report and the policy service's `evt.so.action.dispatched` arrive about 1 ms apart, so the report is sometimes applied first and the task went `created → received` without `sent` (onemind-xr `docs/evidence/INT/E2E-2`) | Order-tolerant: an accepted assignee report on a `created` or `scheduled` task records `sent` first (`updated_by` = the task service, causation = the asset command id), then the report, so the facts are always `created, sent, received, …`; the late audit finds the task past `sent` and changes nothing | Amend mission-task.md §5: `sent` is set on the dispatch audit **or** on the first assignee report, whichever the task service applies first |

### D2-P5: the `mission` user in the A19 layout

```yaml
  - user: mission
    account: ONEMIND
    stage: A
    purpose: "task service (D2/H4); the only writer of KV mission-state and of evt.mission.>"
    secret: {namespace: onemind-core, name: mission-service-runtime, user_key: NATS_USER, password_key: NATS_PASSWORD}
    connection_types: [STANDARD]
    publish:
      - "evt.mission.>"
    subscribe:
      - "svc.mission.>"
      - "tel.fleet.*.*.position"
      - "_INBOX.>"
    allow_responses: true
    jetstream:
      info: true
      consume: [CMD, EVT]     # durable consumers mission-commands (cmd.mission.>) and mission-events
      kv_write: [mission-state]
```

The exact JetStream subjects the service needs (verified by the unit tests' permission
config, which fails on any violation): `$JS.API.INFO`, `$JS.API.STREAM.INFO.*`,
`$JS.API.CONSUMER.CREATE.{CMD,EVT}.>`, `$JS.API.CONSUMER.{INFO,MSG.NEXT}.CMD.mission-commands`,
`…EVT.mission-events`, `$JS.ACK.CMD.mission-commands.>`, `$JS.ACK.EVT.mission-events.>`,
`$JS.FC.>`, and for KV `mission-state`: direct get, consumer create/delete (watch) and
`$KV.mission-state.>`.

## Evidence

onemind-xr `docs/evidence/M4/D2/` (branch `run/m0`): unit test output, both end-to-end runs
(before POL-16: D2-P1 reproduced, and with the since-withdrawn patch) with transcripts and
process logs; after POL-16, `docs/evidence/INT/E-MIS/`: the default run flies the drone, task
`done_ok`, mission `completed`; `docs/evidence/INT/F-MIS/`: D2-P14 race tests and the e2e with
`MISSION_ARRIVAL=false`, where `done_ok` is the bridge's arrival report (D2-P2).

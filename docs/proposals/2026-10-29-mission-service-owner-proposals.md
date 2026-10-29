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
| **D2-P1 (blocker for asset tasks)** | omos #554 stamps the fleet command of a `task.assign` with the `task.assign` `params_hash`; omos #550 recomputes `params_hash` as `asset.{verb}` over the command `data` and refuses the command (`not_authorized`, "params_hash does not match the approved parameters"). Reproduced end to end (`services/mission/test/e2e.sh`, unpatched run): every task-dispatched goto is refused. | Records the refusal faithfully (`done_not_ok` + `rejected`) | Fix in #554: a task-bearing asset command carries the hash of `asset.{verb}` over its own data; the `task.assign` approval stays bound to the action in `EVT` and the Ed25519 signature still binds the whole envelope. The 12-line change is `services/mission/test/interop/policy-task-fleet-params-hash.patch`; with it the e2e run flies the drone and the task completes. (Alternative: #550 accepts `task.assign` hashes, but it cannot recompute them without the task params.) |
| D2-P2 | The bridge reports `executing` when the vehicle accepts a `goto` and never `done_ok` (movement verbs keep executing), so contract test 1 could never reach `done_ok` | Completes an asset `goto` task from `tel.fleet.{class}.{id}.position` when the asset is within 5 m horizontally and 5 m vertically of the goal (`status.result = {arrived, source: telemetry, distance_m}`, `updated_by: service:task`); `MISSION_ARRIVAL=false` turns it off | Ratify as the interim; long term the bridge reports `done_ok` on arrival (H2/H6 follow-up) and the flag goes off |
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
(unpatched: D2-P1 reproduced; patched: drone flies, task `done_ok`, mission `completed`) with
transcripts and process logs.

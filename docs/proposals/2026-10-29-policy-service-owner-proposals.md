# Owner proposals from the policy service (FEATURES I4)

Status: **proposals**, 2026-10-29, filed under onemind-xr decision 0004 #8 ("use the closest
existing pattern and write a proposal in omos `docs/`"). The policy service (`services/policy`)
uses **no new subject pattern**: every subject below is already in the locked grammar
(omos `CLAUDE.md`) and in onemind-xr `docs/design/BUS-DESIGN.md` §3.3–3.4.

## Subjects used

| Subject | Grammar line | Stored by |
|---|---|---|
| `svc.policy.command.propose`, `svc.policy.command.confirm`, `svc.policy.task.propose`, `svc.policy.task.confirm`, `svc.policy.action.get`, `svc.policy.check` | `svc.{service}.{endpoint}` (endpoint may contain dots, doctrine §2 `svc.meta.task.list`; settled as RA2) | — (request/reply) |
| `cmd.fleet.{class}.{id}.{verb}` | doctrine §9.1 (five tokens; BUS-P2 / ONT-P4) | `CMD` (`cmd.>`) |
| `cmd.entity.{id}.{verb}`, `cmd.mission.{id}.{verb}` | `cmd.{domain}.{entity}.{action}` | `CMD` |
| `evt.so.action.{event}` | `evt.{domain}.{entity}.{event}` (ONT-P10) | `EVT` (`evt.so.>`) |
| `evt.{class}.{id}.command` (read), `evt.mission.{id}.{event}` (read), `evt.entity.{id}.{event}` (read) | `evt.{domain}.{entity}.{event}` | `EVT`, `MISSION`; `evt.entity.>` is not bound yet (BUS-P5), the service reads it on core NATS |
| `svc.entity.get` (called) | `svc.{service}.{endpoint}` | — |

## Proposals

| ID | Gap | What the service does now | Proposed ruling |
|---|---|---|---|
| POL-1 | asset-command.md §5 names no event for a *refused confirmation* (agent tries to confirm, same human twice, unauthorized approver) | Publishes `evt.so.action.denied` with `code` and `attempt_by`; the action keeps its state (`awaiting_confirmation`) | Ratify: a `denied` event with `attempt_by` is a refused attempt, not the end of the action; clients read `data.status` |
| POL-2 | §6.4 has no `not_found` / `wrong_endpoint` code | Uses `conflict` for an unknown `action_id`, a confirm on the wrong family endpoint, and a `params_hash` mismatch, with a message | Ratify `conflict`, or add `not_found` to §6.4 |
| POL-3 | The first approval of an A4 action has no event of its own | Recorded in the action (`svc.policy.action.get`) and in the `approvals[]` of the later `approved` or `expired` event | Ratify, or add `evt.so.action.confirmed` (inside the grammar; `EVT` stores it) |
| POL-4 | `evt.vehicle.>` is not bound in `EVT` (ONT-P13) | Subscribes to `evt.vehicle.*.command` anyway; vehicle outcomes would not be stored | Add `evt.vehicle.>` to `evt.yaml` with ONT-P13 |
| POL-5 | Contract §5 wants the decision row in a database with an outbox | Audit goes to `EVT` with a JetStream ack before any reply or command (fail closed); action records are in memory | Postgres table plus outbox (onemind-postgres) as a follow-up PR; until then a restart forgets pending actions, which then cannot be confirmed (fails safe) |
| POL-6 | Hermes stage 1 is planned inside `ONEMIND` (BUS-DESIGN §7.2) | Same-account requests are `actor_verified: false` | Put Hermes in `AGENTS` from the first permissions PR ([per-service permissions](2026-10-29-nats-per-service-permissions.md)) |
| POL-7 | `estop` level (contract chooses A3; the drones lane chose two-person) | A3, and preconditions skipped so a stop is never blocked | H6 reviewer to confirm |
| POL-8 | `svc.entity.get` body and reply are not fixed yet (A4 entity service) | Sends `{entity_type, entity_id}`; accepts the entity record, `{entity: …}` or `{error: {code: not_found}}` | A4 adopts this shape |
| POL-9 | KV `geofences` value format is not fixed | Reads a zone entity record (entity.md `geo`: `polygon` or `circle`, `kind` `keep_in`/`keep_out`) with optional `properties.applies_to: [asset ids]`. A keep-in/keep-out record it cannot read (other shapes, malformed coordinates, missing kind, invalid JSON) fails closed: it is logged and movement verbs are denied for the assets it applies to (`detail: geofence_unevaluable`); `estop`, `rtl`, `hold` still go out. Unreadable `aoi`/`hazard`/`control_area` zones are ignored | The zone writer (`cmd.map.geofence.*`) adopts this shape; `corridor` support later |
| POL-10 | `CMD` keeps 720 h and `discard: old` | Commands publish with a JetStream ack; a full stream with `discard: new` would make the publish fail and the action fail closed | ONT-P7: 365 d and `discard: new` |
| POL-11 | `svc.policy.check` has no body in the contracts | `{v, user, relation, object, context?}` (OpenFGA Check pass-through, fail closed; answered only for callers in the policy service's own account, since CLIENT and AGENTS import `svc.policy.>` and could otherwise probe the whole permission graph) or `{v, action}` (dry run of a proposal as the caller, open to everyone) | Ratify |
| POL-12 | JADC2-MAP §5.3 proposes a CI check for forbidden names | `TestStopLine` fails if any catalogue or task type contains a stop-line word | Keep it in `build-policy.yml` (runs on every PR that touches the schemas) |
| POL-13 | Task specs for `orbit` and `observe` (mission-task.md §4 writes `center \| target_entity_id` and `target: {…} \| entity_id`) | `orbit`: `center` or `target_entity_id`; `observe`: `target` is `{lat, lon, hae_m}` or `{entity_id}` | Ratify in mission-task.md |
| POL-14 | asset-command.md §2 says the policy service turns entity and task service results into `evt.so.action.executed`/`.failed`, but fixes no correlation field | A result on `evt.entity.{id}.*` or `evt.mission.{id}.*` closes the action when its `causation` is the command id or its `data.action_id` is the action, and its subject id is the target (or task); `evt.mission.{tsk}.created` closes a task assigned to a person, team or agent. Zone commands (D1-P4, omos #569: `zone.create|update|remove`, A2, `cmd.map.geofence.{id}.{verb}`) close on the four-token `evt.map.geofence.{created|updated|deleted|refused}` whose `data.zone_id` is the target. Event `failed`, `rejected`, `refused` or `error` → `failed` (`data.code`, default `precondition_failed`), anything else → `executed`. No result within 30 s → `failed` `no_ack` | The A4 entity service and the task service set `causation` to the command id on every result |
| POL-15 | §6.4 has no timeout code beyond `no_ack`, and the contract gives no deadline after `acknowledged` | `no_ack` with a message for: no bridge ack in 10 s, no `executed`/`failed` within 10 min of `acknowledged`, no entity/task result within 30 s. Every state then ends, so in-memory records stay bounded | Ratify, or add a `timeout` code to §6.4 |
| POL-16 (D2-P1) | Which `params_hash` a task-bearing `cmd.fleet` command carries. #550 recomputed `asset.{verb}` over the command data and refused every task dispatch (D2-P1, #568) | The contract binds an approval to one `action_id` and its `params_hash` (asset-command.md §4, §5) and puts that approval on the command (§6.3; mission-task.md §6 step 6), so the command carries the approved action's own hash: `asset.{verb}` over `data` without `task_id`, `task.assign`/`task.abort` over the task parameters with `task_id`. This service is unchanged; the bridges (#550, #566) recompute only when `task_id` is absent and otherwise check the `sha256:` form and rely on the Ed25519 signature, which covers `data` and `approval`. Shared vectors: `services/schemas/command/testdata/params-hash-vectors.json`. The #568 patch `policy-task-fleet-params-hash.patch` is then not needed | Ratify, or amend §6.3 to name a per-command hash |

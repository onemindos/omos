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
| `evt.{class}.{id}.command` (read), `evt.mission.{tsk}.created` (read) | `evt.{domain}.{entity}.{event}` | `EVT`, `MISSION` |
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
| POL-9 | KV `geofences` value format is not fixed | Reads a zone entity record (entity.md `geo`: `polygon` or `circle`, `kind` `keep_in`/`keep_out`) with optional `properties.applies_to: [asset ids]`; other shapes are ignored | The zone writer (`cmd.map.geofence.*`) adopts this shape; `corridor` support later |
| POL-10 | `CMD` keeps 720 h and `discard: old` | Commands publish with a JetStream ack; a full stream with `discard: new` would make the publish fail and the action fail closed | ONT-P7: 365 d and `discard: new` |
| POL-11 | `svc.policy.check` has no body in the contracts | `{v, user, relation, object, context?}` (OpenFGA Check pass-through, fail closed) or `{v, action}` (dry run of a proposal) | Ratify |
| POL-12 | JADC2-MAP §5.3 proposes a CI check for forbidden names | `TestStopLine` fails if any catalogue or task type contains a stop-line word | Keep it in `build-policy.yml` (runs on every PR that touches the schemas) |
| POL-13 | Task specs for `orbit` and `observe` (mission-task.md §4 writes `center \| target_entity_id` and `target: {…} \| entity_id`) | `orbit`: `center` or `target_entity_id`; `observe`: `target` is `{lat, lon, hae_m}` or `{entity_id}` | Ratify in mission-task.md |

# Ontology proposals ONT-P1 … ONT-P14

Source: onemind-xr `docs/design/ONTOLOGY.md` §9 (branch `run/m0`). Index and ruling order:
[README.md](README.md). Entries that duplicate a BUS item point to it and are ruled with it.

---

## ONT-P1. `ENT` stream CR

**Proposed.** Approve `platform/nats/streams/ent.yaml` as filed in #552: stream `ENT`, subjects
`ent.*.*.state`, R3 file, `maxMsgsPerSubject: 1`, `allowDirect: true`, `allowMsgTtl: true`,
`subjectDeleteMarkerTtl: 15m`, `duplicateWindow: 2m`, `maxAge: 168h`, `maxBytes: 1 GiB`,
`discard: old`, `preventDelete: true` (onemind-xr `docs/design/BUS-DESIGN.md` §4.1). Ticking this
box rules #552 A4-P1; the CR and its activation steps are in that PR.

**Why the grammar does not cover it.** The subject exists (`CLAUDE.md` line 90,
`ent.{entity_type}.{entity_id}.state`), but no stream stores `ent.>` (no filter on `ent.>` in any
`platform/nats/streams/*.yaml`; the checker reports `stored-today=none`). Late joiners get no
snapshot and a consumer that is down loses state.

**Alternatives considered.** A KV bucket per entity type (KV is a stream with the same limits,
but loses per-message TTL and subject filters on `ent.{entity_type}.>`); core-only `ENT_MODE=core`
(no restore after restart). Both rejected for production; core mode stays for local use.

**Interim reading.** Not applied. #552 runs locally against a disposable server, or with
`ENT_MODE=core` (entity.md §1).

**Depends on it.** #552 sync (the only blocker it lists besides its bucket A4-P2). Every
consumer of `ent.>` (layer service, projector rework A6, #554 `svc.entity.get`) follows.

```subjects
current   ent.{entity_type}.{entity_id}.state
filter    ent.*.*.state
filter    ent.>
filter    ent.{entity_type}.>
current   svc.entity.get
```

- [ ] Approve as proposed (rules #552 A4-P1)
- [ ] Approve with changes:
- [ ] Reject

---

## ONT-P2. Entity lifecycle facts have no stored domain

Same gap and ruling as [BUS-P5](bus.md#bus-p5-domain-entity-and-evtentity-in-evt).

- [ ] Ruled with BUS-P5

---

## ONT-P3. `DET` stream CR

**Proposed.** New stream `DET` in `platform/nats/streams/det.yaml`: subjects `det.>`, R3 file,
`maxAge: 7d`, `maxBytes: 2 GiB`, `discard: old`, `duplicateWindow: 2m`, `preventDelete: true`
(BUS-DESIGN §4.1). Filed as its own CR PR after the ruling, excluded from the Application until
a detector exists.

**Why the grammar does not cover it.** `det.{domain}.{type}.{id}` is in the grammar (`CLAUDE.md`
line 86), but no stream binds `det.>`.

**Alternatives considered.** Store detections in `TEL` (they are not telemetry of the
publisher; mixing them breaks `tel.>` filters). Keep detections core-only (a detector burst
while the entity service restarts is lost). Off-grammar `det.cot.{id}` from an earlier draft is
not used.

**Interim reading.** No detector publishes yet. #552 already lists `DET` as an input and skips
it while the stream is absent (`services/entity/internal/bus/service.go` `DefaultInputs`).

**Depends on it.** Detectors (FEATURES H8), Remote ID receiver (`det.air.remoteid.{id}`); #552
picks `DET` up when it exists.

```subjects
current   det.{domain}.{type}.{id}
proposed  det.air.remoteid.{id}
filter    det.>
filter    tel.>
legacy    det.cot.{id}
```

- [ ] Approve as proposed
- [ ] Approve with changes:
- [ ] Reject

---

## ONT-P4. Two command shapes

Same gap and ruling as [BUS-P2](bus.md#bus-p2-asset-command-shape-and-the-cancel-mission-verb).

- [ ] Ruled with BUS-P2

---

## ONT-P5. Two telemetry shapes

Same gap and ruling as [BUS-P1](bus.md#bus-p1-telemetry-token-two-is-a-class).

- [ ] Ruled with BUS-P1

---

## ONT-P6. `tel.feed.{source}.status` is not stored

Ruled with [BUS-P1](bus.md#bus-p1-telemetry-token-two-is-a-class): the subject stays and
`TEL` gets `tel.feed.>`. The CR line is in #553 (A0-2), not here.

```subjects
current   tel.feed.{source}.status
filter    tel.feed.>
```

- [ ] Ruled with BUS-P1

---

## ONT-P7. `CMD` retention for audit

**Proposed.** `CMD` (`platform/nats/streams/cmd.yaml`): `maxAge: 365d` (today `720h`, line 16),
`maxBytes: 512 MiB`, `discard: new` (today `old`, line 15), and an alert at 80 % of the byte cap.
Filed in the stream caps PR (A22).

**Why.** Not a subject gap. `CMD` is the audit receipt of every command (run rule 6): the policy
service publishes `cmd.>` with a JetStream ack. Thirty days is shorter than any audit need, and
with `discard: old` a full stream silently drops the oldest records; with `discard: new` it
refuses the publish, so the policy service fails closed and no unaudited command goes out
(BUS-DESIGN §4.1 deviates from the NATS lane's `discard: old` on purpose).

**Alternatives considered.** A separate `AUDIT` mirror with long retention (second copy, same
data); the Postgres decision table with outbox (#554 POL-5) as the only long record (planned as
the durable record, but not built yet). Keep 720 h and rely on a ClickHouse sink (not deployed).

**Interim reading.** `CMD` unchanged (720 h, `discard: old`). #554 still publishes with an ack
and fails closed on a publish error (POL-10).

**Depends on it.** #554 (POL-10); A22 stream caps PR.

```subjects
filter    cmd.>
```

- [ ] Approve as proposed
- [ ] Approve with changes:
- [ ] Reject

---

## ONT-P8. `entity_type` vocabulary

**Proposed.** Ratify the token two vocabulary of `ent.{entity_type}.{entity_id}.state` from
onemind-xr `docs/contracts/entity.md` §3: `drone robot vehicle node` (own fleet, matching the
§9.1 classes), `aircraft vessel` (world), `person team agent`, `sensor camera`, `track`,
`place structure zone`, `event alert`. Decide #552 A4-P4: add `satellite`, or keep satellites as
`vehicle` with `kind: satellite`. Recommended: add `satellite` (the projector's `COT_TYPE_MAP`
already has the key, `platform/omos/tak-projector.yaml`). Also change the `CLAUDE.md` line 90
example `ent.ship.mmsi123.state` to `ent.vessel.mmsi-123.state`, since the vocabulary uses
`vessel`, not `ship`.

**Why the grammar does not cover it.** No new pattern: `{entity_type}` is a variable. The
grammar fixes no vocabulary, and the example uses `ship`.

**Alternatives considered.** Free vocabulary per publisher (consumers cannot filter
`ent.vessel.>` reliably; rejected). Keep `ship` (conflicts with entity.md and the Lattice lane).

**Interim reading.** entity.md §3 list; satellites as `ent.vehicle.norad-{n}.state` with
`kind: satellite`, `environment: space` (#552 A4-P4).

**Depends on it.** #552 (types every entity); #554 reads entity records by type.

```subjects
proposed  ent.{drone|robot|vehicle|node|aircraft|vessel|person|team|agent|sensor|camera|track|place|structure|zone|event|alert}.{entity_id}.state
proposed  ent.satellite.{entity_id}.state
proposed  ent.vessel.mmsi-123.state
interim   ent.vehicle.norad-{n}.state
legacy    ent.ship.mmsi123.state
filter    ent.vessel.>
```

- [ ] Approve as proposed (with `satellite`)
- [ ] Approve with changes (for example without `satellite`):
- [ ] Reject

---

## ONT-P9. New `cmd` domains `entity`, `team`, `watch`, `feed`

**Proposed.** Ratify four domains inside `cmd.{domain}.{entity}.{action}`, all published only by
the policy service (onemind-xr `docs/design/ONTOLOGY.md` §6.2):

| Subject | Actions |
|---|---|
| `cmd.entity.{entity_id}.{override\|correlate\|decorrelate\|delete\|star}` | entity edits (A1–A2) |
| `cmd.team.{team_id}.{create\|add_member\|remove_member\|set_leader}` | teams (A2) |
| `cmd.watch.{watch_id}.{create\|delete}` | watch rules in KV (A2) |
| `cmd.feed.{source}.{enable\|disable}` | feeds (A2) |

**Why the grammar does not cover it.** The pattern exists; doctrine §2 maps domains to vault
pillars and does not list these four. `CMD` (`cmd.>`) stores them.

**Alternatives considered.** Put them under `cmd.so.entity.{verb}` (loses the target id token); route them
as `svc.{service}.{endpoint}` requests (bypasses the `CMD` audit receipt). Rejected.

**Interim reading.** As proposed. Only `cmd.entity.{entity_id}.{verb}` is built (#554
catalogue `entity.override`, `entity.star`).

**Depends on it.** #554 publishes `cmd.entity.{entity_id}.{verb}`; #552 is the consumer once it
handles entity commands (not built yet). Team, watch and feed commands: I3, later PRs.

```subjects
proposed  cmd.entity.{entity_id}.{override|correlate|decorrelate|delete|star}
proposed  cmd.team.{team_id}.{create|add_member|remove_member|set_leader}
proposed  cmd.watch.{watch_id}.{create|delete}
proposed  cmd.feed.{source}.{enable|disable}
current   cmd.entity.{entity_id}.{verb}
current   cmd.{domain}.{entity}.{action}
current   svc.{service}.{endpoint}
alternative  cmd.so.entity.{verb}
```

- [ ] Approve as proposed
- [ ] Approve with changes:
- [ ] Reject

---

## ONT-P10. One audit subject

**Proposed.** Every policy decision and outcome is one record on `evt.so.action.{event}`,
written only by the policy service, with event ∈ `proposed denied approved rejected expired
dispatched acknowledged executed failed refused` (onemind-xr `docs/contracts/asset-command.md`
§1). #554 POL-3 asks whether to add `confirmed` for the first approval of a two-person action;
that ruling stays in #554 and needs no new pattern. After ruling, correct the description of
`services/schemas/command/envelope.schema.json` (line 5), which names `evt.command.accepted`
and three siblings.

**Why the grammar does not cover it.** Three lanes proposed three audit subjects:
`evt.command.{result}` (three tokens, in the envelope schema description; stored by no stream),
`evt.policy.{decision}.{result}` (stored by no stream) and `evt.so.action.{event}`. Only the last
is stored today (`EVT` binds `evt.so.>`, `evt.yaml` line 11).

**Alternatives considered.** The two above, each needing a new `EVT` filter; one stream per
decision type. Rejected.

**Interim reading.** `evt.so.action.{event}`.

**Depends on it.** #554 publishes it (fail closed on the ack); clients and Hermes read it.
Same as: RA2 (audit half).

```subjects
proposed  evt.so.action.{proposed|denied|approved|rejected|expired|dispatched|acknowledged|executed|failed|refused}
proposed  evt.so.action.{event}
proposed  evt.so.action.confirmed
filter    evt.so.>
legacy    evt.command.accepted
legacy    evt.command.{result}
alternative  evt.policy.{decision}.{result}
```

- [ ] Approve as proposed
- [ ] Approve with changes:
- [ ] Reject

---

## ONT-P11. Asset tasks live in the `mission` domain

**Proposed.** Asset tasks use the `mission` domain with `tsk-…` ids: facts
`evt.mission.{task_id}.{state}` (stored by `MISSION`, `evt.mission.>`), commands
`cmd.mission.{task_id}.abort`. Bridges report task status on `evt.{class}.{id}.task` (stored
for `drone` and `robot`). Nothing OneMind builds publishes `evt.task.>` or `task.>`.

**Why the grammar does not cover it.** "Task" already names `om-task` (`task.>`,
`evt.task.created`, BUS-P4), and `EVT` stores `evt.task.>` for it (`evt.yaml` line 15). Two
meanings on one domain would mix asset tasks into the om-task history.

**Alternatives considered.** A new domain `assettask` (new pattern-like domain for the same
lifecycle the doctrine already gives `mission`); share `evt.task.>` with om-task (collision).

**Interim reading.** As proposed (onemind-xr `docs/contracts/mission-task.md` §1, §6).

**Depends on it.** #554 (`task.assign`, `task.abort`, `svc.policy.task.{propose|confirm}`); #550
reports `evt.drone.{id}.task` (its P7 is about `status_version`, not the subject); H4 task
service.

```subjects
proposed  evt.mission.{task_id}.{state}
proposed  cmd.mission.{task_id}.abort
proposed  evt.{drone|robot|vehicle}.{id}.task
current   evt.drone.{id}.task
current   evt.{class}.{id}.task
proposed  svc.policy.task.{propose|confirm}
action    task.assign
action    task.abort
filter    evt.mission.>
legacy    evt.task.>
legacy    task.>
legacy    evt.task.created
alternative  evt.assettask.{task_id}.{state}
```

- [ ] Approve as proposed
- [ ] Approve with changes:
- [ ] Reject

---

## ONT-P12. Robot and vehicle CoT types (information only)

**Recorded.** The projector and the drones lane differed on CoT types for robots and vehicles;
settled in onemind-xr `docs/contracts/robotics-cot-profile.md`. No subject involved.

- [ ] Noted (no ruling needed)

---

## ONT-P13. `evt.vehicle.>` is not bound in `EVT`

**Proposed.** Add `evt.vehicle.>` to `EVT` before a commandable vehicle exists, so bridge
outcomes `evt.vehicle.{id}.command` and task reports `evt.vehicle.{id}.task` are stored like
`evt.drone.>` and `evt.robot.>` (`evt.yaml` lines 18–19).

**Why.** Pattern exists; the stream filter is missing.

**Alternatives considered.** Bind `evt.*.*.command` generically (catches other domains'
`command` events by accident). Rejected.

**Interim reading.** No commandable vehicle exists; #554 subscribes `evt.vehicle.*.command` on
core (POL-4), so vehicle outcomes would not be stored.

**Depends on it.** #554 (POL-4); a later vehicle bridge. Ticking this rules POL-4.

```subjects
filter    evt.vehicle.>
filter    evt.vehicle.*.command
filter    evt.drone.>
filter    evt.robot.>
proposed  evt.vehicle.{id}.command
proposed  evt.vehicle.{id}.task
alternative  evt.*.*.command
```

- [ ] Approve as proposed (rules #554 POL-4)
- [ ] Approve with changes:
- [ ] Reject

---

## ONT-P14. Policy RPC names

**Proposed (settled in the design set, filed for confirmation).**
`svc.policy.command.propose` and `svc.policy.command.confirm` for asset commands and every other
non-task action; `svc.policy.task.propose` and `svc.policy.task.confirm` for tasks and missions;
`svc.policy.action.get`; `svc.policy.check`. Confirm endpoints are for human clients only.

**Why the grammar does not cover it.** Both an earlier three-token draft
(`svc.policy.propose`) and these four-token names fit `svc.{service}.{endpoint}` (doctrine §2
example `svc.meta.task.list` has a dotted endpoint). The gap was two names for one thing.

**Alternatives considered.** `svc.policy.propose` with the family in the body (one endpoint, but
the server cannot give agents propose rights and deny confirm per family by subject, which the
permission design needs); separate services `svc.command.propose` and `svc.task.propose`
(two services for one decision record). Rejected.

**Interim reading.** As proposed; #554 serves exactly these.

**Depends on it.** #554 (serves them; its permissions proposal imports `svc.policy.>` into
`CLIENT` and three of them into `AGENTS`). Clients and Hermes call them. Same as: RA2 (naming
half).

```subjects
proposed  svc.policy.command.propose
proposed  svc.policy.command.confirm
proposed  svc.policy.task.propose
proposed  svc.policy.task.confirm
proposed  svc.policy.action.get
proposed  svc.policy.check
filter    svc.policy.>
current   svc.meta.task.list
current   svc.{service}.{endpoint}
alternative  svc.policy.propose
alternative  svc.command.propose
alternative  svc.task.propose
```

- [ ] Approve as proposed
- [ ] Approve with changes:
- [ ] Reject

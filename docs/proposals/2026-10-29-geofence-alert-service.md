# Geofence and alert service: owner proposals (FEATURES D1, D6)

Status: **proposed, for the owner**. Date: 2026-10-29. PR: `feat/geofence-alerts`, stacked on
#552 (`feat/entity-service`). Filed under onemind-xr decision 0004 item 8. **This PR changes no
live subject, stream, Application or workload**: the Application `platform-omos-geofence` is
manual sync and no stream CR changes.

For the A24 register (#556, `docs/proposals/2026-10-29-grammar-register/`): the subject
entries below (D1-P2, D1-P3, D1-P7) are new and should be listed in its `cross-reference.md`
by reference to this file, as the A4-P*, POL-* and A6-P* entries are. They are not argued again
there. Every subject is in the `subjects` block at the end, in the register's format.

What the service does: it is the only writer of KV `geofences`, applying zone commands that the
policy service approved and signed; it publishes `evt.map.geofence.{id}.{enter|exit}` from
`ent.*.*.state`; and it raises alerts for geofence crossings, stale own assets and feed outages.
Design and tests: `services/geofence/README.md`.

| ID | Decision asked | Blocks | Recommendation |
|---|---|---|---|
| D1-P1 | KV `geofences` value shape | projector (#563) and policy (#554) agreeing with the writer | adopt |
| D1-P2 | Command results on `evt.map.geofence.{created\|updated\|deleted\|refused}` | policy closing zone actions | adopt |
| D1-P3 | Alert facts on `evt.so.alert.{raised\|updated\|resolved}` in `EVT` | D6 alert stream | adopt |
| D1-P4 | Policy catalogue `zone.create\|update\|remove` (patch to #554) | any zone write in the cluster | adopt |
| D1-P5 | Owner of `ent.zone.{id}.state` mirror, `ent.alert.{id}.state` records and alert acknowledge/resolve | alert list as entities; `alert.acknowledge` | decide later |
| D1-P6 | NATS user `geofence` (for #562) | per-service identity | adopt |
| D1-P7 | Zone verb `remove` (ONTOLOGY) or `delete` | — | `remove` canonical, `delete` accepted |
| D1-P8 | Zones never watch people | — | ratify |

---

## D1-P1. KV `geofences` value

**Proposed.** Key = zone id (one lowercase NATS token). Value
(`services/schemas/map/geofence.v1.schema.json`):

```json
{"v": 1, "zone_id": "home", "entity_id": "home", "entity_type": "zone", "name": "Home",
 "active": true, "live": true, "lifecycle": "active", "pinned": true,
 "geo": {"shape": "polygon", "kind": "keep_in",
         "coords": {"type": "Polygon", "coordinates": [[[-97.75, 30.26], [-97.73, 30.26], [-97.73, 30.28], [-97.75, 30.26]]]},
         "radius_m": 0, "alt_min_m": null, "alt_max_m": null},
 "watch": true, "watch_types": [], "properties": {"applies_to": ["sim-x1"]},
 "created_t": 1793275200000, "updated_t": 1793275200000, "created_by": "user:zeus",
 "updated_by": "user:zeus", "action_id": "act-…", "command_id": "…", "simulated": true}
```

It is the union of what the two existing readers accept: the projector's minimal form
(#563 A6-P2: `zone_id`, `active`, `geo{shape, kind, coords, radius_m, alt_min_m, alt_max_m}`)
and the policy service's zone record (#554 POL-9: `entity_id`, `live`, `geo` with `polygon` or
`circle`, `properties.applies_to`). `live` always equals `active`, because the policy service
reads only `live`. `keep_in` and `keep_out` zones are limited to `polygon` and `circle`: the
policy service treats any other keep-in/keep-out shape as unevaluable and denies every movement
of the assets it applies to.

**Evidence.** onemind-xr `docs/evidence/M4/D1/compat-projector-policy.txt`: the six record
shapes this service writes were fed to the real `aoiIndex.setZone` (#563 `d346ee2e`) and
`ParseZone`/`CheckPoint` (#554 `66900872`) in scratch copies; all were read with the expected
activity, membership and keep-in/keep-out verdicts.

- [ ] Approve as proposed
- [ ] Approve with changes:
- [ ] Reject

## D1-P2. Zone command results on `evt.map.geofence.{created|updated|deleted|refused}`

**Proposed.** After applying (or refusing) a zone command, the service publishes one fact on
`evt.map.geofence.{created|updated|deleted|refused}` with `causation` = the command id and `data.zone_id`,
`data.action_id`, `outcome`, `code`, `revision`. `MISSION` stores `evt.map.>` today, so no
stream change is needed.

**Why.** asset-command.md §2 makes executors report each command they apply as a fact so the
policy service can close the action (#554 POL-14). The subject has the four tokens of the
grammar `evt.{domain}.{entity}.{event}` (domain `map`, entity type `geofence`), like the audit
subject `evt.so.action.{event}`; the zone id travels in the data. The crossing facts keep the
doctrine §9.7 five-token form `evt.map.geofence.{id}.{enter|exit}`; the two never collide
because their token counts differ. A first draft used five tokens with the zone id; the A24
checker (#556 `scripts/check_proposal_subjects.py`) rejected that form as off-grammar, so it
was changed before this PR was opened. Four tokens also fit #554's result handler, which reads
four-token subjects.

**Alternatives.** `evt.entity.{zone_id}.{event}` (BUS-P5): not stored until BUS-P5 is ruled,
and zones are not written by the entity service. `evt.so.action.{event}`: its single writer is
the policy service. Rejected.

- [ ] Approve as proposed
- [ ] Approve with changes:
- [ ] Reject

## D1-P3. Alert facts on `evt.so.alert.{raised|updated|resolved}`

**Proposed.** The D6 alert stream is `EVT` with subjects `evt.so.alert.raised`, `.updated`,
`.resolved`; data `services/schemas/alert/alert.v1.schema.json` (`alert_id`, `alert_kind`
`geofence|stale_asset|feed_outage`, `level` `advisory|caution|warning`, `state`, `about`,
`reason`, `seq`). `Nats-Msg-Id: {alert_id}:{seq}`. `EVT` stores `evt.so.>` (90 d, 2 min dedup)
today. Clients read open alerts with `svc.geofence.alerts` or replay `evt.so.alert.>`.

**Why.** It is doctrine §2's own example form (`evt.so.task.created`: domain `so`, entity type,
past-tense event) and the only option stored today, the same reasoning that chose
`evt.so.action.*` for audit (ONT-P10).

**Alternatives considered.**
- `evt.entity.{entity_id}.alert` (ONTOLOGY alert row, BUS-P5): not stored until BUS-P5 is
  ruled, and a feed outage has no entity. Can be emitted in addition once BUS-P5 is approved.
- `ent.alert.{id}.state` in `ENT` (ONTOLOGY alert row): `ent.*.*.state` has one writer. The entity
  service restores every `ent.>` subject on start (`services/entity/internal/bus/service.go`,
  `restore`, filter `ent.>`) and would then sweep and expire records it did not create. Rejected
  until D1-P5.
- A new `ALERT` stream: a new CR for low-volume facts that `EVT` already holds. Rejected.

- [ ] Approve as proposed
- [ ] Approve with changes:
- [ ] Reject

## D1-P4. Policy catalogue entries for zones (patch to #554)

**Proposed.** ONTOLOGY §6.2 lists `zone.create` / `zone.update` / `zone.remove` at level A2,
executing as `cmd.map.geofence.{id}.{create|update|remove}`. The #554 catalogue does not have
them yet, so in the cluster no zone can be written. The parameter schemas are in this PR
(`services/schemas/actions/zone.*.v1.schema.json`), where #554's `LoadSchemas` finds them.
The #554 change, for its owner to apply:

```go
// catalogue.go
const execMap = "map" // cmd.map.geofence.{id}.{verb}
for _, v := range []string{"create", "update", "remove"} {
	Catalogue["zone."+v] = ActionSpec{Name: "zone." + v, Family: FamilyCommand, Level: A2,
		Versions: []int{1}, Targets: []string{"zone"}, Exec: execMap, Verb: v}
}
// service.go dispatch: treat execMap like execEntity, with
//   commandSubject(execMap, …) = "cmd.map.geofence." + target.entity_id + "." + verb
// handleResult: also subscribe evt.map.geofence.* (four tokens, D1-P2); the result carries
//   causation = command id and data.action_id; outcome "refused" maps to failed with data.code.
```

Two consequences for #554 and #562: the policy service needs no write access to KV `geofences`
(asset-command.md §2: policy never executes anything itself), so the `policy` identity can drop
from `kv_write: [geofences, mission-state]` to `kv_read: [geofences]` plus `kv_write:
[mission-state]`; and the geofence service is the bucket's writer (D1-P6).

**Why signatures are required here.** Clients and agents may publish `cmd.map.>` directly
(asset-command.md §2, for reversible view actions), so the subject alone cannot show that a
zone write was approved. The service refuses any `cmd.map.geofence.*` without a valid policy
`sig`, a future `expires_at` and `approval.status: approved` (tested: `TestCommandGate`, and
phase 5 of the e2e run).

- [ ] Approve as proposed
- [ ] Approve with changes:
- [ ] Reject

## D1-P5. Zone mirror, alert records and alert acknowledge/resolve

**Open question, no change in this PR.** ONTOLOGY puts each zone on `ent.zone.{id}.state` and
each alert on `ent.alert.{id}.state`, and #554 catalogues `alert.acknowledge` / `alert.resolve`
as `cmd.entity.{id}.{acknowledge|resolve}`. Options: (a) the entity service mirrors KV
`geofences` and `evt.so.alert.*` into `ent.zone.{id}.state` and `ent.alert.{id}.state` and applies the two alert
commands by forwarding them here; (b) `ent.*.*.state` gets a second writer for those two types, and the
entity service skips them on restore. Recommendation: (a), keeping one writer per subject family
(REFERENCE-ARCHITECTURE R3).

- [ ] (a)
- [ ] (b)
- [ ] Other:

## D1-P6. NATS identity for the A19 file (#562)

**Proposed.** Add to `platform/nats/identities/identities.yaml` (layout of #562), then render
with `scripts/nats_identities.py values`:

```yaml
  - user: geofence
    account: ONEMIND
    stage: A
    purpose: "geofence and alert service (D1/D6, this PR); the only writer of KV geofences"
    secret: {namespace: omos, name: omos-geofence-nats, user_key: user, password_key: password}
    connection_types: [STANDARD]
    publish:
      - "evt.map.geofence.>"
      - "evt.so.alert.>"
    subscribe:
      - "tel.feed.*.status"
      - "svc.geofence.alerts"
      - "_INBOX.>"
    allow_responses: true
    jetstream:
      info: true
      consume: [CMD, ENT, MISSION, EVT]
      kv_write: [geofences]
```

Until then the Deployment uses the shared `omos-nats-creds` (as the entity service does).

- [ ] Approve as proposed
- [ ] Approve with changes:
- [ ] Reject

## D1-P7. Zone verb

ONTOLOGY §6.2 writes `cmd.map.geofence.{id}.{create|update|remove}`; doctrine §9.1/§9.7 writes
`{verb}` with the example `create`; the D1 task text wrote `delete`. The service accepts both
`remove` and `delete` and records the result as `deleted`. Recommendation: `remove` canonical
(it matches `cmd.map.annotation.{id}.{verb}` "drop / update / remove" in doctrine §9.7),
`delete` kept as an alias.

- [ ] `remove` canonical, `delete` alias
- [ ] Other:

## D1-P8. Zones never watch people

People (`person` entities and person-class tracks) are never evaluated against zones, and a zone
with `watch_types: ["person"]` is refused `forbidden_action`. This follows the civilian stop
line (onemind-xr JADC2-MAP §5: no tracking or targeting of people) and the projector's G4 gate.
A future team-member safety use (for example "my team member left the safe area") would need
consent and its own decision.

- [ ] Ratify
- [ ] Other:

```subjects
proposed  evt.map.geofence.{created|updated|deleted|refused}
proposed  evt.so.alert.{raised|updated|resolved}
proposed  svc.geofence.alerts
current   evt.map.geofence.{id}.{enter|exit}
current   evt.map.geofence.{id}
current   cmd.map.geofence.{id}.{create|update|remove}
current   cmd.map.geofence.{id}.{verb}
interim   cmd.map.geofence.{id}.delete
current   cmd.map.annotation.{id}.{verb}
current   cmd.entity.{id}.{acknowledge|resolve}
current   evt.so.action.{event}
current   evt.so.task.created
current   evt.{domain}.{entity}.{event}
current   evt.so.alert.raised
alternative  evt.entity.{zone_id}.{event}
alternative  evt.entity.{entity_id}.alert
alternative  ent.alert.{id}.state
alternative  ent.zone.{id}.state
filter    cmd.map.geofence.>
filter    cmd.map.geofence.*
filter    cmd.map.>
filter    evt.map.>
filter    evt.so.>
filter    evt.so.action.*
filter    evt.so.alert.>
filter    evt.so.alert.*
filter    ent.*.*.state
filter    ent.>
filter    tel.feed.*.status
```

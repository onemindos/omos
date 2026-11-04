# Bus proposals BUS-P1 … BUS-P11

Source: onemind-xr `docs/design/BUS-DESIGN.md` §10 and §10.1 (branch `run/m0`). Index and
ruling order: [README.md](README.md). "Grammar" means the `CLAUDE.md` block "NATS Subject
Grammar (LOCKED)" (lines 84–95); "doctrine" means `docs/NATS-DOCTRINE.md`. Nothing here changes
a live subject or stream.

---

## BUS-P1. Telemetry: token two is a class

**Proposed.** The second token of every `tel` subject is a class from the stream filter list,
never a node id:

| Source | Subject | Stored by |
|---|---|---|
| Fleet assets | `tel.fleet.{class}.{id}.{metric}` (doctrine §9.1) | `TEL` (`tel.fleet.>`, today) |
| Environment sensors | `tel.{class}.{node}.{metric}`, class ∈ `soil weather water nature bio` | `ENV` (today) |
| Infrastructure | `tel.host.{host}.{metric}` | `SYS` (today) |
| Feed freshness | `tel.feed.{source}.status` | `TEL` after adding `tel.feed.>` |

Add one line to the `CLAUDE.md` grammar block: `tel.fleet.{class}.{id}.{metric}  fleet
telemetry  tel.fleet.drone.x47.position` (doctrine §9.1), and change the example on line 85
from `tel.drone.x47.position` to an environment example such as `tel.weather.ws1.reading`.

**Why the grammar does not cover it.** Line 85 reads `tel.{node}.{class}.{metric}` but its own
example `tel.drone.x47.position` puts a class second. Doctrine §2 has a node second
(`tel.barn-pi-1.temp.c`). The streams are class-first: `TEL` binds `tel.fleet.>`, `tel.drone.>`,
`tel.robot.>`, `tel.edge.>`, `tel.iot.>`, `tel.industrial.>`, `tel.sat.>`
(`platform/nats/streams/tel.yaml` lines 11–17) and `ENV` binds `tel.nature.>`, `tel.bio.>`,
`tel.weather.>`, `tel.soil.>`, `tel.water.>`. A node-first subject such as
`tel.sim-x1.drone.position` (the first local MAVLink proof, onemind-xr
`docs/research/standards.md` line 148) is stored by no stream, and the checker confirms it
(`stored-today=none`). Doctrine §9.1 adds a five-token fleet form the block does not list.
Nothing stores `tel.feed.>` today, although `feed-intel` publishes it every 30 s
(`platform/omos/intel-scripts.yaml` line 380).

**Alternatives considered.**
- Node first everywhere, with filters `tel.*.drone.>` and so on: rewrites seven live `TEL`
  filters and five `ENV` filters, and filters with a leading wildcard grow with every class
  (doctrine §2: "wildcard filters do the work"). Rejected.
- Class first for fleet too, four tokens (`tel.drone.x47.position`, as the example reads): `TEL`
  stores it (`tel.drone.>`), but it is not symmetric with `cmd.fleet.{class}.{id}.{verb}` and
  needs one filter per class instead of `tel.fleet.>`. Rejected; §9.1 is already the documented
  fleet form and the one the MAVLink bridge rework uses.
- Leave it ambiguous: subjects end up silently unstored. Rejected.

**Interim reading.** This ruling (BUS-DESIGN §10 "used until ruled"). Fleet publishers use
§9.1; nobody publishes node-first.

**Depends on it.** #550 publishes `tel.fleet.drone.{id}.position` and `.status` (its P1 asks
for the same `CLAUDE.md` line); #551 publishes `tel.fleet.edge.{id}.status`; #553 publishes
`tel.feed.{source}.status` and **adds `tel.feed.>` to `TEL` in that PR** (A0-2), so this PR
makes no CR change; #552 consumes `tel.fleet.>`. Later: site agents
(`tel.fleet.node.{site}.status`, A23), the MQTT gateway (environment classes).

**Same as:** ONT-P5, ONT-P6, RA4 (form), #550 P1, #553 A0-2.

```subjects
proposed  tel.fleet.{class}.{id}.{metric}
proposed  tel.{soil|weather|water|nature|bio}.{node}.{metric}
proposed  tel.host.{host}.{metric}
proposed  tel.feed.{source}.status
proposed  tel.fleet.drone.x47.position
proposed  tel.weather.ws1.reading
proposed  tel.fleet.node.{site}.status
current   tel.fleet.drone.{id}.position
current   tel.fleet.edge.{id}.status
filter    tel.fleet.>
filter    tel.drone.>
filter    tel.robot.>
filter    tel.edge.>
filter    tel.iot.>
filter    tel.industrial.>
filter    tel.sat.>
filter    tel.nature.>
filter    tel.bio.>
filter    tel.weather.>
filter    tel.soil.>
filter    tel.water.>
filter    tel.feed.>
legacy    tel.drone.x47.position
legacy    tel.barn-pi-1.temp.c
legacy    tel.sim-x1.drone.position
proposed  tel.{class}.{node}.{metric}
current   tel.{node}.{class}.{metric}
alternative  tel.*.drone.>
```

- [ ] Approve as proposed
- [ ] Approve with changes:
- [ ] Reject

---

## BUS-P2. Asset command shape and the `cancel` mission verb

**Proposed.**
1. Asset commands are `cmd.fleet.{class}.{id}.{verb}` (doctrine §9.1, five tokens), published
   only by the policy service after policy, audit and human confirm. `{verb}` is the closed
   list of stop line SL5 (onemind-xr `docs/design/JADC2-MAP.md` §5.2): `goto hold rtl land orbit
   observe estop arm disarm takeoff estop_release geofence_override`. BUS-DESIGN §3.3 lists only
   seven of these; JADC2-MAP F4 aligns it to the contract (onemind-xr `contracts/asset-command.md`
   §3), and the policy schemas in #554 (`services/schemas/actions/asset.*.v1.schema.json`)
   already carry all twelve.
2. Add `cancel` to the §9.1 mission verbs (`cmd.mission.{id}.{verb}`, today `create waypoint
   abort complete`).
3. Add the fleet line to the `CLAUDE.md` grammar block (see F1 for the example).

**Why the grammar does not cover it.** `CLAUDE.md` line 88 locks four-token
`cmd.{domain}.{entity}.{action}`; doctrine §9.1 (line 190) locks five-token
`cmd.fleet.{class}.{id}.{verb}`. Both are "locked" and they disagree on token count. `CMD` binds
`cmd.>` (`platform/nats/streams/cmd.yaml` line 11), so both are stored. §9.1 has no `cancel`.

**Alternatives considered.**
- `cmd.drone.{id}.{verb}` (onemind-xr `docs/research/security-and-trust.md` R10): four tokens,
  fits line 88, stored by `CMD`. Rejected because it is not symmetric with `tel.fleet.>`, needs a
  permission line per class instead of `cmd.fleet.>`, and §9.1 is already the documented fleet
  form used by FEATURES D3 and the open PRs.
- Keep `abort` as the only way to stop a mission (no `cancel`): onemind-xr
  `contracts/mission-task.md` line 38 does exactly this until ruled, and `cancelled` is a task
  state reached through `abort`. Acceptable; the owner may reject item 2 without changing
  anything built.

**Interim reading.** The §9.1 form; cancelling is `cmd.mission.{id}.abort`.

**Depends on it.** #554 publishes `cmd.fleet.{class}.{id}.{verb}` and `cmd.mission.{id}.{verb}`;
#550 subscribes `cmd.fleet.drone.{id}.>` and refuses verbs outside its catalogue
(`internal/command/command.go`, `unknown_action`). Later: H4 task service, ROS bridge.

**Same as:** ONT-P4. Related: F1 (the `CLAUDE.md` example), JADC2-MAP F4 (verb list).

```subjects
proposed  cmd.fleet.{class}.{id}.{verb}
proposed  cmd.fleet.{class}.{id}.{goto|hold|rtl|land|orbit|observe|estop|arm|disarm|takeoff|estop_release|geofence_override}
proposed  cmd.mission.{id}.cancel
interim   cmd.mission.{id}.abort
current   cmd.mission.{id}.{create|waypoint|abort|complete}
current   cmd.fleet.drone.{id}.>
current   cmd.mission.{id}.{verb}
current   cmd.{domain}.{entity}.{action}
filter    cmd.>
filter    cmd.fleet.>
alternative  cmd.drone.{id}.{verb}
```

- [ ] Approve as proposed
- [ ] Approve with changes:
- [ ] Reject

---

## BUS-P3. Live five- and six-token `intel` subjects

**Proposed.** Keep the live `feed-intel` subjects as they are (grammar
`intel.{domain}.{source}.{type}` plus trailing qualifiers, for existing feeds only). Every new
collector publishes exactly four tokens, with the item id in the payload and in `Nats-Msg-Id`.
Retire each legacy subject when its `feed-intel` task is replaced, after a side-by-side run.

**Why the grammar does not cover it.** The grammar has four tokens. Live subjects have five or
six: `intel.weather.noaa.alert.{state}.{event}` (`platform/omos/intel-scripts.yaml` line 457),
and per onemind-xr `docs/research/omos-live-state.md` §3 also `intel.seismic.usgs.{bucket}.{id}`,
`intel.conflict.acled.{iso3}.{type}`, `intel.fire.nasa.firms.{region}`,
`intel.disaster.gdacs.{type}.{id}`, `intel.seismic.emsc.{bucket}.{id}`,
`intel.natural.nasa.eonet.{cat}.{id}`. `INTEL` binds `intel.>`, so all are stored.

**Alternatives considered.** Rename the live subjects now (breaks the consumers and the 90-day
history in `INTEL`; rejected). Allow trailing qualifiers for new collectors too (the hierarchy
then grows per feed and consumers filter on payload-like tokens; rejected). A legacy
`fabric.intel.weather.nws.alerts.*` subject from the imported source catalog is not on the bus
(onemind-xr `docs/research/standards.md` Risks item 2) and is not proposed.

**Interim reading.** Live subjects unchanged; new collectors four tokens. Source token for NWS
is `noaa`, as the live feed and #553 use (`intel.weather.noaa.alert`); onemind-xr
`docs/research/standards.md` line 216 wrote `nws`, and this register follows the live token.

**Depends on it.** #553 (A0-3: the new collector publishes `intel.weather.noaa.alert` beside the
legacy subject for a seven-day comparison). #552 consumes `intel.>` and is unaffected.

```subjects
proposed  intel.{domain}.{source}.{type}
current   intel.weather.noaa.alert
filter    intel.>
legacy    intel.weather.noaa.alert.{state}.{event}
legacy    intel.seismic.usgs.{bucket}.{id}
legacy    intel.conflict.acled.{iso3}.{type}
legacy    intel.fire.nasa.firms.{region}
legacy    intel.disaster.gdacs.{type}.{id}
legacy    intel.seismic.emsc.{bucket}.{id}
legacy    intel.natural.nasa.eonet.{cat}.{id}
alternative  intel.weather.nws.alert
legacy    fabric.intel.weather.nws.alerts.*
```

- [ ] Approve as proposed
- [ ] Approve with changes:
- [ ] Reject

---

## BUS-P4. Out-of-grammar legacy subjects

**Proposed.** Record these as legacy: kept, never extended, migrated to grammar subjects when
their service is next changed. OneMind services add nothing to them.

| Subject | Where | Note |
|---|---|---|
| `saas.>`, `voice.>`, `zeus.>`, `op.zeus.>`, `sys.>` | streams `SAAS`, `VOICE`, `ZEUS`, `SYS` | no OneMind publisher |
| `agents.discover`, `agents.heartbeat` | stream `AGENTS` | Synadia protocol uses `agents.hb.>` and `$SRV.PING.agents` |
| `task.>`, `evt.task.created` | `om-task` | four-token form, but collides with asset tasks (ONT-P11) |
| `omos.air.>`, `omos.bridge.>` | old MAVLink bridge | #550 moves the bridge off them; the dev dashboard and `cmd/omos` still subscribe `omos.>` (#550 P3) |
| `ent.manual.upsert`, `ent.manual.delete`, `ent.air.{id}.deleted` | undeployed `om-overwatch-*` (`services/nex/overwatch/`) | outside `ent.*.*.state`; deletes are a final `live: false` record |
| `int.command.{id}`, `int.ack.{id}`, `int.phase.{id}`, `int.target.{id}`, `int.health.{id}`, `int.result.{id}` | undeployed `om-overwatch-interceptor` (`services/nex/overwatch/om-overwatch-interceptor/main.go` lines 11–15) | **found while filing**: no grammar prefix, and an engagement service; see F1 |
| `cmd.projector.disable`, `cmd.projector.enable` | live `omos-tak-projector` | three tokens; CP-1 proposes the replacement |

**Why the grammar does not cover it.** None of these prefixes (`saas`, `voice`, `zeus`, `op`,
`sys`, `task`, `omos`, `int`) is in the grammar block, and the `ent` and `cmd` items have the
wrong shape.

**Alternatives considered.** Migrate them all now (touches live services that OneMind does not
own; rejected). Add them to the grammar (would bless the patterns; rejected).

**Interim reading.** Not extended.

**Depends on it.** #550 (already migrated the MAVLink bridge). No open PR publishes any of
these. Recorded with: CP-1.

```subjects
legacy    saas.>
legacy    voice.>
legacy    zeus.>
legacy    op.zeus.>
legacy    sys.>
legacy    agents.discover
legacy    agents.heartbeat
legacy    task.>
legacy    evt.task.created
legacy    omos.air.>
legacy    omos.bridge.>
legacy    omos.>
legacy    ent.manual.upsert
legacy    ent.manual.delete
legacy    ent.air.{id}.deleted
legacy    int.command.{id}
legacy    int.ack.{id}
legacy    int.phase.{id}
legacy    int.target.{id}
legacy    int.health.{id}
legacy    int.result.{id}
legacy    cmd.projector.disable
legacy    cmd.projector.enable
current   agents.hb.{agent}.{owner}.{session}
current   $SRV.PING.agents
filter    agents.hb.>
filter    ent.*.*.state
```

- [ ] Approve as proposed
- [ ] Approve with changes:
- [ ] Reject

---

## BUS-P5. Domain `entity` and `evt.entity.>` in `EVT`

**Proposed.** Accept domain `entity` for entity lifecycle and conflict facts:
`evt.entity.{entity_id}.{event}`, event ∈ `created deleted expired correlated decorrelated
overridden conflict alert`. Add `evt.entity.>` to `EVT` (`platform/nats/streams/evt.yaml`) in a
follow-up CR PR after the ruling.

**Why the grammar does not cover it.** The form fits `evt.{domain}.{entity}.{event}`, but
doctrine §2 maps domains to the vault pillars (`so hp le ge`, plus `sys`) and `entity` is a new
domain. `EVT` binds ten domains (`evt.yaml` lines 11–20) and not `entity`, so these events would
be core-only and lost for any consumer that is down.

**Alternatives considered.**
- `evt.so.entity.{event}` with the id in the payload: stored today, but loses the per-entity
  filter `evt.entity.{entity_id}.>` that the policy service and timeline use. Rejected.
- `evt.{entity_type}.{entity_id}.{event}` (for example `evt.drone.sim-x1.created`): collides
  with the bridge outcome events already on `evt.drone.>`. Rejected.
- No events, only `ENT` delete markers: core subscribers do not get markers (BUS-DESIGN §4.2),
  and conflicts have no record. Rejected.

**Interim reading.** Not published. Entity removal works through the final `live: false`
record and `Nats-TTL` (or deletion by sequence while TTL is off; #552 A4-P1).

**Depends on it.** #552 emits `created`, `expired`, `deleted`, `conflict` only after the ruling
(A4-P5). #554 closes `cmd.entity.>` actions on a result on `evt.entity.{entity_id}.{event}`
(POL-14) and reads them on core NATS; **until BUS-P5 is ruled and #552 emits, every
`cmd.entity.>` action ends `failed` with `no_ack` after 30 s**. Fail-safe, but entity edits do
not complete.

**Same as:** ONT-P2, #552 A4-P5.

```subjects
proposed  evt.entity.{entity_id}.{created|deleted|expired|correlated|decorrelated|overridden|conflict|alert}
proposed  evt.entity.{entity_id}.{event}
filter    evt.entity.>
filter    evt.entity.{entity_id}.>
filter    cmd.entity.>
alternative  evt.so.entity.{event}
alternative  evt.drone.sim-x1.created
alternative  evt.{entity_type}.{entity_id}.{event}
current   evt.{domain}.{entity}.{event}
filter    evt.drone.>
```

- [ ] Approve as proposed
- [ ] Approve with changes:
- [ ] Reject

---

## BUS-P6. World classes `rail` and `transit`

**Proposed.** Add `rail` and `transit` to the doctrine §9.1 world classes:
`evt.world.rail.{id}`, `evt.world.transit.{id}`, stored by `WORLD` (`evt.world.>`).

**Why the grammar does not cover it.** The four-token form fits `evt.{domain}.{entity}.{event}`
reading `world` as the domain, but §9.1 lists the world classes (`air sea conflict military fire
seismic chokepoint news`) and has none for trains or buses (onemind-xr
`docs/research/ingestion-platform.md` migration step 8).

**Alternatives considered.** Reuse `evt.world.sea.{id}` or `evt.world.air.{id}` (wrong class). Publish as
`intel.transit.mta.position` (positions are not intelligence reports, and `INTEL` keeps 90 days).

**Interim reading.** MTA stays on its current path (`feed-mta` puts CoT on `tak.cot.inject`);
it migrates last. #552 types trains as `vehicle` with `kind: train` (A4-P4 note).

**Depends on it.** A later transit collector. No open PR publishes these subjects.

```subjects
proposed  evt.world.rail.{id}
proposed  evt.world.transit.{id}
filter    evt.world.>
current   tak.cot.inject
alternative  intel.transit.mta.position
alternative  evt.world.sea.{id}
alternative  evt.world.air.{id}
```

- [ ] Approve as proposed
- [ ] Approve with changes:
- [ ] Reject

---

## BUS-P7. Fleet telemetry metrics beyond `position` and `status`

**Proposed.** Any metric name is allowed in `tel.fleet.{class}.{id}.{metric}`, as the grammar's
`{metric}` already says. First two: `tel.fleet.drone.{id}.sensor` (sensor pose from MISB KLV)
and `tel.fleet.drone.{id}.footprint` (frame footprint).

**Why the grammar does not cover it.** Doctrine §9.1 lists only `.position` and `.status`.

**Alternatives considered.** Put KLV pose into `.status` (mixes a 1–30 Hz stream into a 1 Hz
status, and consumers of status would parse video metadata). Put the footprint only on the
entity (the client video drape needs the per-frame stream).

**Interim reading.** `sensor` and `footprint`.

**Depends on it.** FEATURES K5 KLV extractor (no PR yet). Same as: RA4 (metric half).

```subjects
proposed  tel.fleet.drone.{id}.sensor
proposed  tel.fleet.drone.{id}.footprint
proposed  tel.fleet.{class}.{id}.{position|status}
```

- [ ] Approve as proposed
- [ ] Approve with changes:
- [ ] Reject

---

## BUS-P8. No geographic subject token (information only)

**Recorded, not proposed.** The grammar has no geographic token, and none is proposed. Area of
interest filtering is done by `svc.layer.watch` (BUS-DESIGN D7), not by subjects. This entry
exists so that nobody adds a geohash or tile token later without a ruling.

```subjects
current   svc.layer.watch
```

- [ ] Noted (no ruling needed)

---

## BUS-P9. Dead-letter subject

**Proposed.** Rejected input is published on `evt.sys.{source}.rejected`, stored by `SYS`
(`evt.sys.>`, `platform/nats/streams/sys.yaml` line 11). `{source}` is the component that
rejected it: the feed source id for a collector, the service name for a service. Feed ids must
not equal a service name.

**Why the grammar does not cover it.** The form fits `evt.{domain}.{entity}.{event}`; the gap is
that the ingestion lane proposed `evt.feed.{source}.rejected`, which no stream binds (onemind-xr
`docs/research/ingestion-platform.md` Risks "Subjects"), and there was no agreed dead-letter
place.

**Alternatives considered.** `evt.feed.{source}.rejected` plus a new `EVT` filter (a new domain
for one subject; rejected). A separate dead-letter stream (more streams, same data; rejected).

**Interim reading.** `evt.sys.{source}.rejected`, as both open PRs already do.

**Depends on it.** #552 publishes `evt.sys.entity.rejected` (`{source}` = service, A4-P3); #553
publishes `evt.sys.{source}.rejected` with the feed source id (`runtime.go` `Reject`, A0-1).
Both readings fit the proposed rule.

```subjects
proposed  evt.sys.{source}.rejected
current   evt.sys.entity.rejected
filter    evt.sys.>
alternative  evt.feed.{source}.rejected
```

- [ ] Approve as proposed
- [ ] Approve with changes:
- [ ] Reject

---

## BUS-P10. Shrink, then retire `GEO`; no new publisher of `evt.tak.cot.{uid}`

**Proposed.** Shrink `GEO` to 7 days and 512 MiB in the stream caps PR (A22). Delete the CR
once `onemind-core` no longer subscribes `evt.tak.cot.>` (`services/onemind-core/onemind/ingestion.py`
line 85). `preventDelete: true` keeps the data either way. This applies decision 0004 #7, which
took doctrine §9.8 option 1: `tak.cot.>` is the raw TAK wire and `ENT` replaces the per-UID
geo-envelope.

**Why the grammar does not cover it.** Doctrine §9.1 defines `evt.tak.cot.{uid}` and
`evt.tak.status.{uid}`; the live TAK plugin publishes `tak.cot.out` instead, and `GEO`
(`geo.yaml` lines 11–12) has no publisher (doctrine §9.8).

**Alternatives considered.** Doctrine §9.8 option 2, moving the plugin to `evt.tak.cot.{uid}`
(restarts a tier-1 workload and migrates `COT` consumers; rejected by decision 0004 #7). Delete
`GEO` now (one live subscriber still reads it; rejected).

**Interim reading.** Shrink only; no OneMind service publishes `evt.tak.cot.{uid}`.

**Depends on it.** The stream caps PR (A22, not yet open). Interaction: #551 P2 proposes
`evt.tak.datapkg.received`. `GEO` binds only `evt.tak.cot.>` and `evt.tak.status.>`, so that
event is stored by no stream whether or not `GEO` is retired; its binding is #551's ruling, not
this one. Related: decision 0004 #7, doctrine §9.8.

```subjects
filter    evt.tak.cot.>
filter    evt.tak.status.>
legacy    evt.tak.cot.{uid}
legacy    evt.tak.status.{uid}
current   tak.cot.out
current   evt.tak.datapkg.received
filter    tak.cot.>
```

- [ ] Approve as proposed
- [ ] Approve with changes:
- [ ] Reject

---

## BUS-P11. Doctrine §4: leaf buffering needs a site stream

**Proposed.** Amend `docs/NATS-DOCTRINE.md` §4 (line 88), "the leaf node buffers during
disconnection", to: "a site stream on the leaf's own JetStream domain buffers during
disconnection; a leaf without a local stream forwards core messages only and drops them while
the link is down."

**Why.** A plain leaf has no store and forward for core publishes (onemind-xr
`docs/research/nats-platform.md` proof P7d; `docs/research/edge-and-ddil.md` ER2). Not a subject
change; it is filed because the current text leads to designs that lose data offline.

**Alternatives considered.** Leave the text and document the caveat in the site template only
(the doctrine would still mislead). Rejected.

**Interim reading.** — (documentation only).

**Depends on it.** Site leaf template (A23); #551 at a site writes through `SITE_COT`.

- [ ] Approve as proposed
- [ ] Approve with changes:
- [ ] Reject

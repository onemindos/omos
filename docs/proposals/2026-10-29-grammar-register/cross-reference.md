# Cross-reference: open PR proposals against this register

Each open omos PR from the OneMind build run files its own owner proposals. This table says, for
every one of them, whether it is a grammar or stream item ruled in this register (then the
ruling happens here once) or a service-level item that stays in its own PR. Read from the PR
branches on 2026-10-29 (`git diff origin/main...origin/<branch>`).

| PR | Proposal file |
|---|---|
| #552 `feat/entity-service` | `docs/proposals/2026-10-29-entity-service-ent-stream.md` |
| #553 `feat/ingest-framework` | `docs/proposals/2026-10-29-a0-omos-ingest.md` |
| #554 `feat/policy-service` | `docs/proposals/2026-10-29-policy-service-owner-proposals.md`, `docs/proposals/2026-10-29-nats-per-service-permissions.md` |
| #550 `feat/mavlink-bridge-fleet` | `docs/proposals/2026-10-29-mavlink-bridge-fleet-subjects.md` |
| #551 `feat/tak-mesh-gateway` | `docs/roadmap/2026-10-29-tak-mesh-gateway-owner-proposals.md` |

Later PRs and client lanes (added 2026-10-29, integration pass). Read from the PR branches with
`git diff origin/<base>...origin/<branch>`; the two onemind-xr client lanes keep their proposals
in evidence files on onemind-xr `run/m0`:

| PR | Proposal file |
|---|---|
| #555 `feat/stream-caps` | `docs/proposals/2026-10-29-stream-caps.md` |
| #557 `feat/agent-eval-redteam` (on #554) | adds POL-16 … POL-20 to #554's `docs/proposals/2026-10-29-policy-service-owner-proposals.md` |
| #558 `feat/feed-grammar-fixes` | `docs/proposals/2026-10-29-feed-grammar-cutover.md` |
| #559 `feat/mediamtx-video` | `docs/proposals/2026-10-29-video-k1-k9.md` |
| #560 `feat/tak-server-leg` | `docs/proposals/2026-10-29-tak-server-leg.md` |
| #561 `feat/layer-service` (on #552) | `docs/proposals/2026-10-29-layer-service.md` |
| #562 `feat/nats-identities` | `docs/proposals/2026-10-29-nats-identities-a19.md` |
| #563 `feat/projector-aoi` (on #561) | `docs/proposals/2026-10-29-projector-aoi.md` |
| #564 `feat/hermes-nats-tools` | `docs/proposals/2026-10-29-hermes-bus-tools-i2.md` |
| #565 `feat/roles-authentik` (on #554) | `docs/proposals/2026-10-29-roles-authentik-j2.md` |
| #566 `feat/ros2-bridge` | `docs/proposals/2026-10-29-ros2-bridge-h3.md` |
| #567 `feat/agents-watch-briefs` | `i1-agents-watch-briefs.md` in this folder (arrives when #567 merges) |
| #568 `feat/mission-service` | `docs/proposals/2026-10-29-mission-service-owner-proposals.md` |
| #569 `feat/geofence-alerts` (on #552) | `docs/proposals/2026-10-29-geofence-alert-service.md` |
| onemind-xr #5 `run/m4-web` | onemind-xr `docs/evidence/M4/WEB/OWNER-PROPOSALS.md` (WEB-P1 … WEB-P9) |
| onemind-xr #4 `run/m4-client` | onemind-xr `docs/evidence/M4/UE/OWNER-STEPS.md`, "Owner proposals from this lane" (UE-P1 … UE-P6) |
| integration pass | [integration-resolutions.md](integration-resolutions.md) (INT-R1 … INT-R4) |

## Mapping

"Here" means the ruling is taken in this register; "in PR" means it stays with that PR.

| PR item | Subject | Register entry | Ruled |
|---|---|---|---|
| #552 A4-P1 `ENT` CR | stream | ONT-P1 | here |
| #552 A4-P2 bucket `entity-correlation` | KV | — (not a grammar item; BUS-DESIGN §5 lists it) | in PR |
| #552 A4-P2b soft keys, `fleet-` id prefix, `sources[].tier`, query paging | payload | — | in PR |
| #552 A4-P3 dead letter `evt.sys.entity.rejected` | subject | BUS-P9 | here |
| #552 A4-P4 `satellite` entity type | vocabulary | ONT-P8 | here |
| #552 A4-P5 lifecycle events | subject, stream | BUS-P5 (= ONT-P2) | here |
| #553 A0-1 dead letter | subject | BUS-P9 | here |
| #553 A0-2 `tel.feed.>` in `TEL` | stream | BUS-P1 (= ONT-P6); the CR line is in #553 | here (ruling), CR in #553 |
| #553 A0-3 NWS four-token subject | subject | BUS-P3 | here |
| #553 A0-4 KV `ingest-state` | KV | — | in PR |
| #553 A0-5 collector identity | permissions | — (per-service users, #554 permissions doc) | in PR |
| #553 A0-6 NATS CA in `omos` | TLS | — | in PR |
| #553 A0-7 `INTEL` dedup window | stream setting | — (stream caps PR, A22) | in PR |
| #553 A0-8 websocket queue policy | service | — | in PR |
| #554 subjects used | subjects | BUS-P2, ONT-P9, ONT-P10, ONT-P14 | here |
| #554 POL-1, POL-2, POL-15 refusal and timeout codes | payload | — (asset-command.md §6.4) | in PR |
| #554 POL-3 optional `evt.so.action.confirmed` | subject (inside ONT-P10's pattern) | noted under ONT-P10 | in PR |
| #554 POL-4 `evt.vehicle.>` | stream | ONT-P13 | here |
| #554 POL-5 Postgres outbox | storage | — (ONT-P7 alternative) | in PR |
| #554 POL-6 Hermes in `AGENTS` | accounts | — | in PR |
| #554 POL-7 `estop` level | policy | — (same question as #550 P5) | in PR (H6 reviewer) |
| #554 POL-8 `svc.entity.get` body | payload | — | in PR (with #552) |
| #554 POL-9 KV `geofences` value | payload | — | in PR |
| #554 POL-10 `CMD` 365 d and `discard: new` | stream | ONT-P7 | here |
| #554 POL-11 `svc.policy.check` body | payload | — (the name is ONT-P14) | in PR |
| #554 POL-12 stop-line CI check | CI | — (supports F1) | in PR |
| #554 POL-13 `orbit`/`observe` task specs | payload | — | in PR |
| #554 POL-14 result correlation | payload | depends on BUS-P5 for `cmd.entity.>` results | in PR |
| #554 per-service permissions | accounts | — | in PR |
| #550 P1 fleet telemetry §9.1 line in `CLAUDE.md` | subject | BUS-P1 | here |
| #550 P2 `action_version` on commands | payload | — | in PR |
| #550 P3 liveness on `tel.fleet.drone.{id}.status`; tools still on `omos.>` | subject | BUS-P4 (legacy `omos.bridge.>`) | here (legacy record), tool fix in PR |
| #550 P4 `no_ack` reuse | payload | — | in PR |
| #550 P5 `estop` A3 | policy | — | in PR (H6 reviewer) |
| #550 P6 bridge NATS user | permissions | — | in PR |
| #550 P7 task `status_version` | payload | — (ONT-P11 fixes the subject only) | in PR |
| #551 P1 mesh pair under `tak.cot.>` | subject | CP-4 | here |
| #551 P2 Object Store `datapkg`, `evt.tak.datapkg.received` | store, subject | noted under BUS-P10 (no stream binds `evt.tak.datapkg.>`) | in PR |
| #551 P3 headers `Tak-Delivery`, `Tak-Dest-Uid` | headers | — (not subjects) | in PR |
| #551 P4 publish permission on `tak.cot.mesh.tx`, `tak.cot.inject` | permissions | — (supports CP-3, CP-4) | in PR |
| #551 P5 mesh keys | secrets | — | in PR |
| #551 P6 STACK edge-host section | docs | — | in PR |
| #551 P7 `MESH_TO_TAK_SERVER` and skip list | behaviour | CP-2 | here |
| #551 P8 data package URLs | behaviour | — | in PR |

## Mapping: later PRs and client lanes

The columns are the same as above. "INT-Rn" means the item is ruled once, in
[integration-resolutions.md](integration-resolutions.md).

| PR item | Kind | Register entry | Ruled |
|---|---|---|---|
| #552 A4-P6 SGP4 dependency `go-satellite` | dependency | — | in PR |
| #558 A1-P1 `rail` world class, `evt.world.rail.{id}`, type `ground.rail` | subject | BUS-P6 | here |
| #558 A1-P2 TLEs stay raw on `intel.space.celestrak.tle`; the entity service propagates (no `evt.world.space.{norad}`) | ownership | — (BUS-DESIGN §3.1) | in PR |
| #558 A1-P3 interim `tel.feed.{mta\|celestrak}.status` once per poll, dead letters, `Nats-Msg-Id` | subject use | BUS-P1 (status), BUS-P9 (dead letter) | here (subjects), rest in PR |
| #558 A1-P4 rail adapter and SGP4 in the entity service | service | — (#552 adopted both, `1c7c0aa8`) | in PR |
| #569 D1-P1 KV `geofences` value | payload | — (one shape for POL-9 and A6-P2) | in PR |
| #569 D1-P2 zone results on `evt.map.geofence.{created\|updated\|deleted\|refused}` | subject | new, four tokens; `MISSION` stores `evt.map.>` today | here (tier 2) |
| #569 D1-P3 alerts on `evt.so.alert.{raised\|updated\|resolved}` in `EVT` | subject | new; rule with I1-P3 (contradiction below) | here (tier 3) |
| #569 D1-P4 `zone.create\|update\|remove` in the #554 catalogue | catalogue | — (same as WEB-P1) | in PR (#554 patch) |
| #569 D1-P5 who writes `ent.zone.{id}.state` and `ent.alert.{id}.state` | ownership | depends on BUS-P5 | in PR |
| #569 D1-P6 NATS user `geofence` | permissions | — (#562 layout) | in PR |
| #569 D1-P7 zone verb `remove` (canonical), `delete` (alias) | verb | doctrine §9.7 | here (tier 2) |
| #569 D1-P8 zones never watch people | stop line | — | in PR |
| #568 D2-P1 `params_hash` on task-dispatched asset commands | payload | INT-R1 | INT-R1 |
| #568 D2-P2, D2-P3, D2-P8 … D2-P13 arrival rule, `status_version`, defaults, parked placeholder | service | — (D2-P3 closes #550 P7) | in PR |
| #568 D2-P4 `svc.mission.{get\|list\|report\|apply}`, refusals on `evt.mission.{id}.rejected` | subject names | ONT-P11 (domain); names inside the grammar | here (tier 2) |
| #568 D2-P5 NATS user `mission`; exports to CLIENT and AGENTS | permissions | INT-R3 (exports) | in PR, INT-R3 |
| #568 D2-P6 mission abort fans out `cmd.fleet.{class}.{id}.hold` | behaviour | BUS-P2 | in PR (#554) |
| #568 D2-P7 policy stops waiting on `rejected` | behaviour | — (POL-14) | in PR (#554) |
| #567 I1-P1 agent contacts in TAK (cot-projection §4.1 exception) | contract | — | here (tier 3) |
| #567 I1-P2 agentops as an interim second publisher of `tak.cot.inject` and `tak.cot.mesh.tx` | permission | CP-3; exit through INT-R4 | here (tier 3) |
| #567 I1-P3 alerts on `evt.entity.{entity_id}.alert`, feed id `feed-{source}` | subject | BUS-P5; rule with D1-P3 | here (tier 3) |
| #567 I1-P4 `evt.so.brief.published`, `svc.brief.{generate\|latest}` | subject | ONT-P10 reasoning | here (tier 2) |
| #567 I1-P5 fleet class `agent`, `tel.fleet.agent.{id}.position` | subject | BUS-P1, ONT-P8 | here (tier 2) |
| #567 I1-P6 NATS user `agentops`, one OpenFGA tuple | permissions | — | in PR (tier 1 for #567) |
| #567 I1-P7 watch definitions stay in Git | information | ONT-P9 | here (information) |
| #566 H3-P1 robot verbs, latched `estop`, `estop_release` A4 | verbs | BUS-P2 (SL5 list) | in PR |
| #566 H3-P2 shared command vectors | tests | INT-R1 | INT-R1 |
| #566 H3-P4 `capabilities` and `meta.kind` | payload | INT-R2 | INT-R2 |
| #566 H3-P3, H3-P5 … H3-P8 map frame, robot users, edge only, Nav2 keep-out, CI e2e | service | — | in PR |
| #565 J2-P1 … J2-P8 auth callout order, user-name groups, roles, token lifetime | accounts | — (no subject; J2-P3 sits next to I2-P5) | in PR |
| #559 K-P1 `evt.video.{entity_id}.{online\|offline}`, `evt.video.>` in `EVT` | stream | RA1 | here |
| #559 K-P5 takml `fabric.vision.>` moves to `det.>` | legacy | BUS-P4 (legacy), ONT-P3 (`DET`) | here (with ONT-P3) |
| #559 K-P7 `svc.registry.video.{get\|list}` | subject names | inside the grammar; CLIENT import in INT-R3 | here (tier 2) |
| #559 K-P9 KV `plugin-registry` and `svc.plugins.list` | bucket, subject | new | here (tier 4) |
| #559 K-P2 … K-P4, K-P6, K-P8, K-P10 entity `media[]`, TAK video auth, storage, JWT, bucket reuse, tailnet name | service | — | in PR |
| #561 A5-P4 `svc.layer.{list\|meta\|snapshot\|watch\|unwatch}`, streamed replies | subject names | inside the grammar | here (tier 2) |
| #561 A5-P5 `tel.fleet.node.{site}.status` as the site heartbeat | subject | BUS-P1 | here (with BUS-P1) |
| #561 A5-P7 CLIENT export settings for `svc.layer.>` | accounts | INT-R3 | INT-R3 |
| #561 A5-P1 … A5-P3, A5-P6, A5-P8 bucket purposes, caps, lease, release rule, contract details | service | — (A5-P2 belongs with A22) | in PR |
| #563 A6-P2 KV `geofences` value | payload | D1-P1 | in PR (with D1-P1) |
| #563 A6-P3 disable sends removals | behaviour | CP-1 (subject form) | in PR |
| #563 A6-P4 tasking guard reads `tak.cot.out`, `tak.cot.mesh.rx` | permissions | — | in PR |
| #563 A6-P1, A6-P5 cutover, smaller gate rules | service | — | in PR |
| #564 I2-P1 `hermes` may request `svc.layer.snapshot` | permission | — | in PR |
| #564 I2-P2 map-action acks from CLIENT back to AGENTS (service import) | permission | INT-R3 | INT-R3 |
| #564 I2-P3 agent map actions recorded in `CMD` | stream capture | — (contradiction with A19-P4 below) | in PR |
| #564 I2-P4 `cmd.map.>` payload v1 and ack | payload | — | in PR |
| #564 I2-P5 `on_behalf_of` for agent proposals | identity | — | in PR |
| #562 A19-P2 Hermes loses `cmd.ha.*` and the firehoses | permissions | — | in PR |
| #562 A19-P4 `cmd.map.>` one-way AGENTS to CLIENT | accounts | — (see I2-P2, I2-P3) | in PR |
| #562 A19-P8 `tak.cot.inject` for the mesh gateway only if forwarding is on | permissions | CP-2 | here (with CP-2) |
| #562 A19-P9 legacy `cmd.projector.{enable\|disable}` | subject | CP-1, BUS-P4 | here (with CP-1) |
| #562 A19-P1, A19-P3, A19-P5 … A19-P7, A19-P10 Hermes account, watch replies, inboxes, agent delegation, JetStream limits, stage B | accounts | — | in PR |
| #555 A22-P1 `COT` from 30 days to 72 h (**removes data**) | stream | BUS-P10 neighbourhood | in PR (owner sign-off) |
| #555 A22-P2 `CMD` 365 d, 512 MiB, `discard: old` for now | stream | ONT-P7 (corrects its `discard: new`) | here (with ONT-P7) |
| #555 A22-P3 new `DET` stream | stream | ONT-P3 | here (with ONT-P3) |
| #555 A22-P4 `INTEL` dedup window to 1 h | stream setting | — (#553 A0-7) | in PR |
| #555 A22-P5 `GEO` to 7 days | stream | BUS-P10 | here (with BUS-P10) |
| #560 A21-P1 … A21-P6 plugin rebuild, `tak-client` user, single producer of `tak.cot.out`, rotation, cot-tap removal, TAK groups | service | — (no new subject) | in PR |
| #557 POL-16 … POL-20 rate limit, text stop line, refused confirm of an unknown action, forbidden names, header forging | policy | — (POL-18 rules with POL-1; POL-20 with A19-P1) | in PR (#554) |
| WEB-P1 zone writes through policy | catalogue | D1-P4 | in PR (#554) |
| WEB-P2 CLIENT subscriptions | accounts | INT-R3 | INT-R3 |
| WEB-P3 dry run authorizes before preconditions, or a `can` probe on `svc.policy.check` | payload | — (POL-11) | in PR (#554) |
| WEB-P4 `svc.cot.send` body | subject body | INT-R4 | INT-R4 |
| WEB-P5 entity `capabilities` | payload | INT-R2 | INT-R2 |
| WEB-P6 Hermes prompt and reply framing | payload | — (next to I2-P4) | in PR (#564) |
| WEB-P7 crossing payload `{entity_id, zone_id, t}` | payload | — (J1-P3) | in PR (#569) |
| WEB-P8 CoT types for sectors, rings, range and bearing | contract | — (cot-projection §5.5) | onemind-xr contract |
| WEB-P9 `entity.create` for persistent drawings | catalogue | — | in PR (#554) |
| UE-P1 `svc.cot.send` kind `shape` | subject body | INT-R4 | INT-R4 |
| UE-P2 tak-client serves `svc.tak.mission.>`, publishes `evt.mission.{tak-token}.changed` | subject | new | here (tier 4) |
| UE-P3 CLIENT imports | accounts | INT-R3 | INT-R3 |
| UE-P4 one TAK user per client device | identity | — | onemind-xr |
| UE-P5 client's own TAK stream limited to SA and human chat | contract | — (cot-projection §1; next to I1-P2) | here (tier 3) |
| UE-P6 OpenSSL licence for Unreal builds | licence | — | onemind-xr (owner) |

## Checked for contradictions

- **Dead-letter `{source}`.** #552 uses the service name (`entity`), #553 the feed source id.
  BUS-P9 defines `{source}` as the rejecting component, so both are valid; feed ids must not equal
  a service name.
- **`tel.feed.>` binding.** Only #553 changes `tel.yaml`. This PR changes no CR.
- **`cancel` verb.** BUS-P2 proposes it; mission-task.md and #554 use `abort` until ruled, so a
  rejection changes nothing built.
- **F1 example.** JADC2-MAP F1 suggested `cmd.drone.x47.goto`, which would contradict BUS-P2; F1
  here uses `cmd.mission.patrol1.abort` and `cmd.fleet.drone.x47.goto` instead.
- **Entity action results.** #554 POL-14 closes entity actions on `evt.entity.{entity_id}.{event}`,
  which #552 does not publish until BUS-P5 is ruled; recorded under BUS-P5 as a dependency, not
  a conflict (the action fails safe with `no_ack`).
- **NWS source token.** #553 and the live feed use `noaa`; onemind-xr standards.md wrote `nws`.
  BUS-P3 follows `noaa`.
- **`GEO` and the TAK event domain.** BUS-P10 retires `GEO` (`evt.tak.cot.>`, `evt.tak.status.>`); #551 P2's
  `evt.tak.datapkg.received` is outside both filters, so neither ruling blocks the other.

- **Alert subject (D1-P3 against I1-P3).** #569 publishes alerts on `evt.so.alert.*`, which
  `EVT` stores today and which also covers feed outages that have no entity. #567 publishes
  them on `evt.entity.{entity_id}.alert`, which nothing stores until BUS-P5 is ruled. Rule the
  two together. The register's reading: `evt.so.alert.*` is the one stored alert fact, and
  `evt.entity.{entity_id}.alert` may be added alongside it once BUS-P5 is approved (D1-P3
  already allows that). #567's watch jobs would then publish on `evt.so.alert.raised`. Both PRs
  key a feed by its source (`feed-{source}` in #567, `about` in #569), and that key must be the
  same string in both.
- **Two writers of geofence crossings.** #569 (geofence service) publishes
  `evt.map.geofence.{id}.{enter|exit}`, and I1-P6 also grants agentops publish on
  `evt.map.geofence.>`. With both running, every crossing would arrive twice. The geofence
  service should be the only writer, and agentops should only consume the crossings. Drop the
  grant from I1-P6.
- **`cmd.map.>` and `ONEMIND` (A19-P4 against I2-P3).** A19-P4 keeps `ONEMIND` out of
  `cmd.map.>`. I2-P3 asks `AGENTS` to export the same subject to `ONEMIND` so that `CMD` can
  record agent map actions. The two can coexist only if `ONEMIND` imports and never re-exports,
  because a re-export would form the loop A19-P4 warns about. Rule them together.
- **`CMD` discard policy (ONT-P7 against A22-P2).** A22-P2 showed that a full `discard: new`
  stream stops storing commands but still delivers them to core subscribers, so the policy is
  not fail-closed as ONT-P7 claimed. Rule ONT-P7 with A22-P2: 365 d, 512 MiB, `discard: old`
  and the 80 % alert, switching to `discard: new` only after A22-P2 option (a). A note was added
  under ONT-P7.
- **Policy's KV grants.** #562 gives `policy` write on `geofences` and `mission-state`. D1-P4
  and D2-P5 each make another service the only writer of one of those buckets. The two asks
  agree: `policy` keeps read on `geofences` and loses both write grants.
- **Second publishers of the TAK outbound subjects.** security-and-trust T4 makes the projector
  the only NATS publisher of `tak.cot.inject` and `tak.cot.mesh.tx`. I1-P2 asks for an interim
  exception for agentops. INT-R4 ends that exception once `svc.cot.send` is served. UE-P5 is a
  TAK-protocol stream from the client to TAK Server, not a NATS publish, so it does not count
  against T4.
- **Refusal event names.** D1-P2 uses `refused` and D2-P4 uses `rejected` for a command that
  the executor did not apply. #554's POL-14 maps both to `failed`, so nothing breaks today. One
  word would be cleaner: rule it with D2-P4.

```subjects
current   evt.sys.entity.rejected
filter    tel.feed.>
current   tel.fleet.drone.{id}.status
legacy    omos.>
legacy    omos.bridge.>
filter    evt.vehicle.>
proposed  evt.so.action.confirmed
current   svc.entity.get
current   svc.policy.check
filter    cmd.entity.>
filter    tak.cot.>
current   tak.cot.mesh.tx
current   tak.cot.inject
current   evt.tak.datapkg.received
filter    evt.tak.datapkg.>
filter    evt.tak.cot.>
filter    evt.tak.status.>
alternative  cmd.drone.x47.goto
proposed  cmd.mission.patrol1.abort
proposed  cmd.fleet.drone.x47.goto
current   evt.entity.{entity_id}.{event}
proposed  evt.world.rail.{id}
current   intel.space.celestrak.tle
alternative  evt.world.space.{norad}
interim   tel.feed.{mta|celestrak}.status
proposed  evt.map.geofence.{created|updated|deleted|refused}
proposed  evt.so.alert.{raised|updated|resolved}
proposed  evt.so.alert.raised
filter    evt.so.alert.*
alternative  ent.zone.{id}.state
alternative  ent.alert.{id}.state
proposed  svc.mission.{get|list|report|apply}
proposed  evt.mission.{id}.rejected
current   cmd.fleet.{class}.{id}.hold
proposed  evt.entity.{entity_id}.alert
proposed  evt.so.brief.published
proposed  svc.brief.{generate|latest}
proposed  tel.fleet.agent.{id}.position
proposed  evt.video.{entity_id}.{online|offline}
filter    evt.video.>
legacy    fabric.vision.>
filter    det.>
proposed  svc.registry.video.{get|list}
proposed  svc.plugins.list
proposed  svc.layer.{list|meta|snapshot|watch|unwatch}
proposed  tel.fleet.node.{site}.status
filter    svc.layer.>
current   svc.layer.snapshot
current   tak.cot.out
current   tak.cot.mesh.rx
filter    cmd.map.>
legacy    cmd.ha.*
legacy    cmd.projector.{enable|disable}
proposed  svc.cot.send
filter    svc.tak.mission.>
proposed  evt.mission.{tak-token}.changed
filter    evt.map.>
current   evt.map.geofence.{id}.{enter|exit}
filter    evt.map.geofence.>
```

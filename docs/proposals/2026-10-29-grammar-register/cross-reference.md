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
```

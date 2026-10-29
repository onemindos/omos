# Proposal: the AOI-gated TAK projector (FEATURES A6)

Status: proposed, 2026-10-29. Nothing here is deployed. Owner decisions are marked **A6-P*n***.
Contract: onemind-xr `docs/contracts/cot-projection.md` (sections 3, 4, 5, 6, 10 and 12).
Stacked on the layer service (A5) and the entity service (A4) branches.

## What changes

| Path | Change |
|---|---|
| `services/projector/` | New Go service `omos-tak-projector` (image `ghcr.io/onemindos/omos-tak-projector`). Reads `ENT`, KV `geofences` and KV `mission-state`; publishes CoT to `tak.cot.inject` (target `server`) or `tak.cot.mesh.tx` (target `mesh`) |
| `platform/omos/projector/deployment.yaml` | Deployment `omos-tak-projector-gated`, one replica, ecosystem pool, requests 50m/96Mi, limits 500m/384Mi, existing `omos-nats-creds` |
| `platform/argocd/platform-omos-projector.yaml` | Argo Application, **manual sync**, zero retries |
| `.github/workflows/build-tak-projector.yml` | Vet and race tests on PRs; GHCR image on `main`; digest artifact to pin |
| `platform/omos/tak-projector.yaml` | Legacy Python projector: the ungated `ent.*.*.state` path is off unless `ENT_PROJECTION=true`; its `intel.>` path is unchanged |
| `platform/hermes/mcp/coverage.yaml`, `docs/STACK.md`, `services/README.md` | Registration rows |

No new subject, stream or bucket. Subjects used: `ent.*.*.state` (read), `tak.cot.inject` /
`tak.cot.mesh.tx` (write), `tak.cot.out` and `tak.cot.mesh.rx` (read only, for the tasking
guard), and the live legacy control subjects `cmd.projector.disable` / `cmd.projector.enable`
(CP-1 in the contract proposes a four-token form; not changed here).

## Gates (contract section 3)

G1 enabled; G2 `live`, lifecycle `active`/`stale`, `locatable`, a valid point; G4 never a
`person`, a `track` whose kind is person-like, or a `station`; G3 a type row (no row means not
projected, counted); G5 no TAK source seen within the stale time; G6
`indicators.egress_tak`; G7 marking (CLEAR/GREEN; AMBER and AMBER+STRICT only when
`releasable_to` meets `TAK_RELEASE_GROUPS`; RED, missing or unknown never); G8 relevance (own
`sim-*` asset, tasked or targeted, active mission assignee, pinned place/zone/structure,
starred, or inside an active AOI); G9 cap (`MAX_PROJECTED`, kept in the section 3.3 order) and
rate (`MAX_RATE_PER_SEC` in any one-second window, per-uid intervals). An entity that fails a
gate it passed before gets a removal (same uid, stale = time + 1 s, remarks "removed"); one that
only left every AOI keeps being sent for `AOI_EXIT_GRACE_S` first.

The legacy projector's live `ent.*.*.state` path applies none of these (it ignores
`egress_tak`, markings and person tracks). That is why it is switched off here.

## Owner decisions

| # | Decision | Proposed | Until ruled |
|---|---|---|---|
| A6-P1 | **Cutover.** Merging this PR rolls the legacy pod (platform-omos auto-syncs) with `ENT_PROJECTION=false`. Nothing publishes `ent.*.state` live today (onemind-xr run prompt, "Facts already verified"), so TAK sees no change. Sync `platform-omos-projector` only after `platform-omos-entity` runs. Retire the legacy `intel.>` path once the entity service turns `intel.*` into `event` entities live (it does in code, `services/entity/internal/adapt/intel.go`), then delete the Python Deployment | as stated | legacy keeps `intel.>` only |
| A6-P2 | **KV `geofences` value shape.** The bucket has no writer yet (D1 will write it through `cmd.map.geofence.*`). The projector reads key = zone id and a JSON value that is either a zone entity record (entity.v1 `geo`) or the minimal form `{zone_id, active, geo: {shape, kind, coords (GeoJSON), radius_m?, alt_min_m?, alt_max_m?}}`. A zone is an AOI when active and `geo.kind` is `aoi`, `keep_in`, `keep_out`, `hazard` or `control_area`, or when it is the `aoi` of an active mission. Areas: Polygon, MultiPolygon, Point + `radius_m` (circle), LineString + `radius_m` (corridor) | adopt for D1 | as stated |
| A6-P3 | **Disable.** The contract says a failed gate sends a removal, so `cmd.projector.disable` removes every projected marker and then stops. The alternative is a silent pause (markers grey out at their stale time) | removals, as built | removals |
| A6-P4 | **Tasking guard.** The projector also subscribes to `tak.cot.out` and `tak.cot.mesh.rx`, read only, to log and count strike/engage CoT (`t-k-*`, `t-e-*`; stop line SL1, JADC2-MAP test 8) and other `t-*` tasking (not supported yet, section 8.3). It never creates, forwards or executes anything. BUS-DESIGN section 7.2 must give the projector user subscribe rights on those two subjects when accounts split | add to the projector's permissions | one shared account today |
| A6-P5 | **Smaller rules the contract leaves open:** (a) an AOI altitude band applies only when the entity's altitude is known; (b) person-like track kinds are `person`, `human`, `pedestrian`, `people`, `crowd`; (c) `station` is never projected (no owner switch); (d) drones, robots and vehicles that are not own `sim-*` assets have no row and are not projected; (e) a simulated entity is sent with `how="m-s"` (SL7) even when `provenance.how` says otherwise (contract section 5.1 vs JADC2-MAP SL7) | as listed | as listed |

## Sync order

1. A4-P1 approved, `ENT` reconciled, `platform-omos-entity` running.
2. This PR merged; `build-tak-projector` on `main` pushed the image; digest pinned in
   `platform/omos/projector/deployment.yaml` by a follow-up PR.
3. Sync `platform-omos-projector` by hand. Check `/status` (port 8080) and
   `nats sub tak.cot.inject`: only own assets, tasked entities and pinned content until D1 writes
   AOIs to `geofences`.

## Not in this PR (follow-ups)

- `<sensor>` / `<__video>` detail (robotics profile sections 4 and 5; H9, K3), `<precisionlocation>`,
  event-area (`.area`) and structure footprint children, the `ALERT:` remarks prefix (section 5.6),
  `svc.cot.send`, and the `t-x-d-d` delete experiment (owner checklist).
- A site deployment of `PROJECTOR_TARGET=mesh` (runs on the site host next to the mesh gateway; F9).
- Cap churn: when more than `MAX_PROJECTED` entities qualify and arrive in a different order than
  their rank, the cap swaps some out (49 removals in the 5,000-entity test). A rank hysteresis
  would reduce this; not needed for correctness.

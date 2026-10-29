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

No new subject pattern, stream or bucket. `svc.cot.send` is already in the grammar register
(PR #556, `proposed  svc.cot.send`, resolution INT-R4) and in the #562 `CLIENT` imports; this PR
serves it. Subjects used: `ent.*.*.state` (read), `svc.cot.send` (serve), `tak.cot.inject` /
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

## svc.cot.send (INT-R4)

The projector serves `svc.cot.send` and stays the single writer of `tak.cot.inject` /
`tak.cot.mesh.tx`: client sends go through the same publisher, kill switch and global ceiling as
projected entities. Rules and bodies: `services/projector/README.md` "svc.cot.send". Where
the build differs from INT-R4 as written, the decisions are below; they are additions to INT-R4
in the #556 register, not a parallel register.

| # | Decision | Proposed | Until ruled |
|---|---|---|---|
| A6-P6 | **Client XML.** INT-R4 says the client never sends XML, but the Unreal client (M2 and M4, `CotSendRequest`, kinds `marker`, `spi`, `route`, `geofence`, `shape`, `chat`) sends built CoT, and only the web GeoChat is structured. The projector accepts both: structured GeoChat is built here; client CoT is parsed strictly, its `<detail>` allowlisted, its uid bound to `onemind-{client_id}-`, and published unchanged | accept both until the Unreal client moves to structured bodies; then drop the CoT form | both forms served |
| A6-P7 | **One reply shape.** The structured form replies `{ok, uid}` / `{ok: false, error: {code, message}}` (INT-R4); the CoT form keeps the Unreal parser's `{ok: false, error: "<code>", detail}` (`CotSendReply::Parse` reads `error` as a string). Codes are the same | converge on INT-R4's object once `CotSendReply::Parse` reads `error.code` | per form |
| A6-P8 | **Sender binding.** INT-R4 takes the sender from `Nats-Request-Info`. A GeoChat `sender.uid` must equal the NATS user or end in `-{user}` (web `OM-WEB-{user}`); `OneMind-Sender: {account}/{user}` is on every publish. The Unreal `client_id` is not bound to the NATS user (the client uses `mac`, `headless`) | bind `client_id` to the user in the identity directory (#562) | client_id unbound; rate limit and headers use the verified user |
| A6-P9 | **Chat marking.** The web GeoChat body has no `marking`. A GeoChat without one takes `COT_SEND_CHAT_DEFAULT_TLP=CLEAR` (TAK users see chat anyway); CoT-form sends must carry a marking (G7 fails closed). This also answers UE-P5's open question for chat sent through the projector | as built | as built |
| A6-P10 | **Who may call.** Only verified users of `COT_SEND_ACCOUNTS` (`CLIENT`). Callers on the shared `onemind` login (own account) are refused, so agentops/Hermes cannot use the endpoint until #562 exports it to `AGENTS` and `COT_SEND_ACCOUNTS` adds `AGENTS` (the I1-P2 exit in INT-R4) | add `AGENTS` with that export, with its own lower rate | `CLIENT` only |
| A6-P11 | **Hub and site on one bus.** Both instances answer in queue group `omos-tak-projector`; across a leaf connection a `server` request could reach the site instance and get `target_failed` (and the reverse). Clients retry; a request for `both` is served by whichever instance answers | a site projector sets `COT_SEND_ENABLED=false` until a per-target endpoint is ruled (no new pattern proposed here) | hub serves; no site deployment yet |
| A6-P12 | **Ceiling.** A client send is refused `rate_limited` when the shared one-second window is full of entity updates, so the target never exceeds `MAX_RATE_PER_SEC`. The alternative reserves part of the window for humans | as built | as built |

Finding for the web client (onemind-xr #2): its B11 CoT export (`b11-cot-export.xml`) writes
`OM.draw.*` uids, the projector's reserved prefix (contract §5.1). The export is a file, not a
send, and the projector refuses those uids if they are ever sent.

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
  and the `t-x-d-d` delete experiment (owner checklist).
- `svc.cot.status` (contract §9) and a `svc.cot.chat` history endpoint (WEB-P2 (c)).
- A site deployment of `PROJECTOR_TARGET=mesh` (runs on the site host next to the mesh gateway; F9).
- Cap churn: when more than `MAX_PROJECTED` entities qualify and arrive in a different order than
  their rank, the cap swaps some out (49 removals in the 5,000-entity test). A rank hysteresis
  would reduce this; not needed for correctness.

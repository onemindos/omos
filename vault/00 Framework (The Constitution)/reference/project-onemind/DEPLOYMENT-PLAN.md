# OneMind — Full Deployment Plan
**Date:** 2026-07-06 | Status: ACTIONABLE — execute in wave order

---

## The Pattern (how everything lands in the platform)

Every service in this system follows ONE deployment pattern. No exceptions.

```
1. Build Docker image
   └── push to ghcr.io/onemindos/{service}:latest

2. Create platform/{service}/ in omos-infra:
   ├── namespace.yaml          (kubectl create ns)
   ├── deployment.yaml         (Deployment + Service + NATS leaf sidecar)
   └── cloudflare-route.yaml   (if public-facing)

3. Add to platform/install.sh  (so full cluster reprovisioning works)

4. Secrets out-of-band via 1Password → kubectl create secret
```

**NATS access pattern (every K8s service that talks NATS):**
Each pod gets a `nats-leaf` sidecar (nats:2.10-alpine) that holds credentials and connects to the cluster. The main container talks to `localhost:4222` plaintext — no TLS/auth complexity inside containers. See `platform/legacy/deployment.yaml` for the exact pattern.

---

## Current State (what's already live)

| Service | Where | Status |
|---------|-------|--------|
| NATS 3-node HA | DOKS | ✅ live |
| 38 streams + 20 KV buckets | DOKS NATS | ✅ provisioned (mostly empty) |
| Legacy (Hermes) | DOKS | ✅ live |
| CloudTAK | DOKS | ✅ live |
| WorldWideView | DOKS | ✅ live |
| TAK Server | DOKS | ✅ live |
| NatsBridgePlugin (TAK→NATS) | TAK | ✅ live |
| cot-tap (NATS→TAK) | DOKS | ✅ live |
| ClickHouse | DOKS | ✅ live |
| Grafana/Loki/Tempo | DOKS | ✅ live |
| Cloudflared tunnel | DOKS | ✅ live |
| Mac NATS leaf | Mac | ⚠️ configured, creds not filled |
| ORP/om-unity | Mac only | ⚠️ not on bus, not on K8s |
| om-tak-projector | Mac NEX | ⚠️ air-only, not on K8s |
| om-control-plane | Mac NEX | ⚠️ not on K8s |
| omos-feed-world | nowhere | ❌ not built |
| omos-feed-intel | nowhere | ❌ not built |
| WWV NATS bridge | nowhere | ❌ not built |
| omos-mission | nowhere | ❌ not built |

---

## Wave 0 — Wire the Bus (1-2 hours, zero code)

**Goal:** Mac and DOKS are ONE bus. Everything else depends on this.

### 0a. Fill Mac leaf credentials

```bash
# Get creds from 1Password
LEAF_USER=$(op read "op://00-24 SO/omos-mac-leaf-node/username")
LEAF_PASSWORD=$(op read "op://00-24 SO/omos-mac-leaf-node/password")

# Replace placeholders in ~/nats-server.conf
sed -i '' "s/LEAF_USER/$LEAF_USER/g" ~/nats-server.conf
sed -i '' "s/LEAF_PASSWORD/$LEAF_PASSWORD/g" ~/nats-server.conf

# Restart
launchctl unload ~/Library/LaunchAgents/com.nats.server.plist
launchctl load ~/Library/LaunchAgents/com.nats.server.plist
```

### 0b. Verify the link

```bash
# From Mac, publish a test message
nats pub tel.test.leaf.ping "hello"

# From DOKS (port-forward), confirm it arrives
kubectl port-forward -n nats svc/nats 4222:4222
nats sub "tel.test.leaf.>" --count 1
```

### 0c. Provision streams on DOKS

```bash
cd omos-infra
./nats-migration/provision.sh
```

**Gate:** NATS leaf shows connected in `nats server list`. Streams provisioned OK=38.

---

## Wave 1 — ORP on the Bus (the fusion engine) — 1-2 days

**Goal:** ORP subscribes to all raw telemetry on the bus and publishes fused entity state.
This is the single most important step — it makes everything else queryable and correlatable.

### 1a. Build ORP with NATS feature enabled

In `omos-core/integrations/ai-agents/orp/` (or canonical ORP source):

```toml
# Cargo.toml — ensure nats feature is in the features list
[features]
nats = ["orp-connector/nats"]
```

```bash
cargo build --release --features nats
```

### 1b. Write ORP connector config for the bus

Create `omos-core/nats/orp-connector.yaml`:

```yaml
connectors:
  - id: "onemind-bus"
    type: "nats"
    url: "nats://localhost:4222/stream/TEL/subject/tel.>"
    properties:
      nats_user: "onemind"
      nats_password: "${ONEMIND_PW}"
  - id: "onemind-world"
    type: "nats"
    url: "nats://localhost:4222/stream/WORLD/subject/evt.world.>"
  - id: "onemind-events"
    type: "nats"
    url: "nats://localhost:4222/stream/EVT/subject/evt.>"
```

### 1c. Add NATS outbound publisher to ORP

ORP's ingest pipeline (`crates/orp-stream/`) needs to publish fused entity state back to NATS after every entity update:

```rust
// After entity graph update, publish to: ent.{entity_type}.{entity_id}.state
let subject = format!("ent.{}.{}.state", entity.entity_type, entity.entity_id);
nc.publish(subject, serde_json::to_vec(&entity_state)?).await?;

// Also write to KV_world-state for latest-value reads
kv.put(&entity.entity_id, serde_json::to_vec(&entity_state)?).await?;
```

### 1d. Containerize ORP

```dockerfile
# platform/orp/Dockerfile
FROM rust:1.78-slim as builder
WORKDIR /app
COPY . .
RUN cargo build --release --features nats

FROM debian:bookworm-slim
COPY --from=builder /app/target/release/orp /usr/local/bin/orp
EXPOSE 9090
CMD ["orp", "serve", "--config", "/etc/orp/config.yaml"]
```

```bash
docker build -t ghcr.io/onemindos/orp:latest -f platform/orp/Dockerfile .
docker push ghcr.io/onemindos/orp:latest
```

### 1e. Deploy ORP to DOKS

Create `omos-infra/platform/orp/`:

```yaml
# namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: orp
```

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: orp
  namespace: orp
spec:
  replicas: 1
  selector:
    matchLabels:
      app: orp
  template:
    spec:
      imagePullSecrets:
        - name: ghcr
      containers:
        - name: orp
          image: ghcr.io/onemindos/orp:latest
          ports:
            - containerPort: 9090
          env:
            - name: NATS_URL
              value: "nats://localhost:4222"
          volumeMounts:
            - name: data
              mountPath: /var/lib/orp
            - name: config
              mountPath: /etc/orp
          resources:
            requests: { cpu: "500m", memory: "1Gi" }
            limits:   { memory: "4Gi" }
        - name: nats-leaf          # ← standard sidecar pattern
          image: nats:2.10-alpine
          args: ["-c", "/etc/nats/leaf.conf"]
          volumeMounts:
            - name: leaf-conf
              mountPath: /etc/nats
      volumes:
        - name: data
          persistentVolumeClaim:
            claimName: orp-data
        - name: config
          configMap:
            name: orp-config
        - name: leaf-conf
          secret:
            secretName: orp-leaf-conf
```

Add `svc.orp.query` NATS request/reply endpoint so Legacy can query ORP-QL from the bus:

```
NATS: svc.orp.query  (request)  →  ORP executes ORP-QL → returns entity JSON
```

**Gate:** `nats sub "ent.>"` on DOKS shows entity updates flowing. ORP UI at `:9090` shows entities from bus sources.

---

## Wave 2 — COP Comes Alive (all entities on every display) — 1 day

**Goal:** Every fused entity appears on iTAK/ATAK/WebTAK/CloudTAK AND atoc AND WorldWideView simultaneously.

### 2a. Generalize om-tak-projector (air → all types)

Current: only handles `det.air.*` → air CoT markers.
Change: subscribe `ent.*.*.state`, map `entity_type` → CoT type using the mapping table from `cot2constellation/bridge.py`.

```go
// types.go — add CoT type map
var entityTypeToCoT = map[string]string{
    "aircraft_multirotor": "a-f-A-M-F-Q",
    "aircraft_fixed_wing": "a-f-A-F-F",
    "surface_vessel_usv":  "a-f-S-S-U",
    "ground_vehicle_wheeled": "a-f-G-U-C-V",
    "operator_station":    "a-f-G-U-C-I",
    "sensor_platform":     "a-f-G-I",
    // fire, weather, etc use b- (hazard) prefix
    "fire":                "b-r-.-H-.",
    "aircraft":            "a-u-A",   // unknown affiliation default
}
```

Subscribe to `ent.*.*.state` (all types) not just `ent.air.*.state`.

Deploy to K8s under `platform/om-tak-projector/`.

### 2b. Deploy om-control-plane to K8s

Move from Mac NEX to K8s pod with NATS leaf sidecar. Subscribe `ent.*.*.state` → publish `ui.*` projections. Same pattern as other services.

### 2c. Build WWV NATS bridge (50 lines)

New service `platform/wwv-bridge/` — Node.js, subscribes NATS `ent.*.*.state`, pushes to WWV via WebSocket:

```javascript
// wwv-bridge/index.js
const { connect } = require('nats');
const WebSocket = require('ws');

const nc = await connect({ servers: 'nats://localhost:4222' });
const js = nc.jetstream();
const wss = new WebSocket.Server({ port: 8765 });

// Subscribe all fused entity state
const sub = nc.subscribe('ent.*.*.state');
for await (const msg of sub) {
  const entity = JSON.parse(msg.data);
  // WWV WsStreamPayload format
  const wwvMsg = JSON.stringify({
    type: 'data',
    pluginId: `onemind-${entity.entity_type}`,
    payload: [{
      id: entity.entity_id,
      type: entity.entity_type,
      label: entity.name || entity.entity_id,
      position: [entity.longitude, entity.latitude, entity.altitude || 0],
      timestamp: entity.last_updated,
      properties: entity.properties,
    }]
  });
  wss.clients.forEach(client => {
    if (client.readyState === WebSocket.OPEN) client.send(wwvMsg);
  });
}
```

WWV plugin connects to this bridge's WebSocket. Each entity_type becomes one plugin layer on the globe.

**Gate:** Open iTAK, CloudTAK, WebTAK, atoc, WorldWideView simultaneously. A drone on the Mac appears on ALL five screens within 2 seconds.

---

## Wave 3 — World Data Flows (filling the empty streams) — 2-3 days

**Goal:** WORLD and INTEL streams fill with real data. The map becomes a live picture of the planet.

### 3a. Build omos-feed-world

New service: Node.js, wraps Crucix adapters as polling loops, publishes Event Envelope v1 to WORLD stream.

Sources (from Crucix `apis/sources/`):
- `opensky.mjs` → `evt.world.tel.adsb.*` (all aircraft)
- `ships.mjs` → `evt.world.tel.ais.*` (maritime vessels)
- `firms.mjs` → `evt.world.det.fire.*` (NASA active fires)
- `noaa.mjs` → `evt.world.det.weather.*` (severe weather)
- `space.mjs` → `evt.world.det.space.*` (satellites, launches)
- `safecast.mjs` → `evt.world.det.environment.*` (radiation)
- `epa.mjs` → `evt.world.det.environment.*` (air quality)
- `kiwisdr.mjs` → `evt.world.det.sigint.*` (SDR radio)

Structure:
```javascript
// For each adapter:
const data = await adapter.briefing();         // call the Crucix function
const envelope = {
  v: 1,
  id: uuid(),
  ts: new Date().toISOString(),
  type: `evt.world.det.fire.firms`,
  source: 'omos-feed-world',
  subject: `evt.world.det.fire.firms`,
  data: data,
};
await nc.publish(envelope.subject, JSON.stringify(envelope));
```

Interval: ADS-B/AIS every 30s, fires/weather every 5min, space/SDR every 15min.

Deploy: `platform/omos-feed-world/` — same K8s pattern, NATS leaf sidecar.

### 3b. Build omos-feed-intel

Same pattern, wraps intel-domain adapters:
- WorldMonitor's 500+ RSS feeds → `intel.news.*` (categorized by SOURCE_TIERS)
- `acled.mjs` → `intel.conflict.*`
- `gdelt.mjs` → `intel.gdelt.*`
- `cisa-kev.mjs` → `intel.cyber.*`
- `opensanctions.mjs` → `intel.sanctions.*`
- `who.mjs` → `intel.health.*`
- `reliefweb.mjs` → `intel.humanitarian.*`
- Finance bundle (fred/yfinance/eia) → `intel.finance.*`

Deploy: `platform/omos-feed-intel/` — K8s, NATS leaf sidecar.

### 3c. ORP picks it up automatically

Once omos-feed-world and omos-feed-intel publish to the bus, ORP is already subscribed to `evt.world.*` and `intel.*` (from Wave 1). It fuses fires, flights, conflict events into the entity graph automatically. No additional wiring needed.

**Gate:** ORP-QL query `MATCH (e:Aircraft) WHERE e.is_active = true RETURN e LIMIT 10` returns live aircraft. `MATCH (f:Fire) RETURN COUNT(f)` returns > 0.

---

## Wave 4 — Legacy as Real Commander — 1 day

**Goal:** Legacy can see everything, query everything, and act on anything.

### 4a. ORP-QL endpoint on NATS

Add request/reply handler in ORP: subscribe `svc.orp.query`, execute ORP-QL, return JSON results.

```rust
// In ORP server startup
let sub = nc.subscribe("svc.orp.query").await?;
while let Some(msg) = sub.next().await {
    let query = String::from_utf8(msg.payload.to_vec())?;
    let result = query_engine.execute(&query).await?;
    msg.respond(serde_json::to_vec(&result)?).await?;
}
```

Legacy can now do:
```python
# In Legacy's NATS plugin
result = await nc.request("svc.orp.query",
    b'MATCH (s:Ship) WHERE NEAR(s, 30.2, -97.7, 50) RETURN s.name, s.speed')
ships = json.loads(result.data)
```

### 4b. Wire Legacy to entity streams

In Legacy's NATS plugin config, add stream subscriptions:
- `ent.*.*.state` — watch all entity updates
- `intel.>` — watch all intel feeds
- `det.>` — watch all detections
- `evt.mission.>` — watch mission events

Legacy becomes a passive observer of the entire bus — it sees everything happening, can ask ORP for context, and issues commands when needed.

### 4c. Build omos-mission (AMOS patterns)

Extract AMOS mission orchestration layer as a NATS-native service.

Core loop:
```
cmd.mission.create (Legacy or human) →
  omos-mission decomposes into tasks →
  publishes cmd.{asset}.{id}.{action} for each task →
  monitors evt.mission.task.* for completion/failure →
  re-plans on failure →
  publishes evt.mission.completed
```

Key AMOS patterns to port:
- Intent → task decomposition
- Task auction protocol (best-fit asset assignment)
- Human approval gate for lethal/high-risk actions
- Failure detection + re-planning

Deploy as K8s pod (`platform/omos-mission/`), Python/Flask with NATS pub/sub replacing HTTP.

**Gate:** Legacy can say "patrol this polygon" → omos-mission decomposes it → drone receives waypoint commands → telemetry flows back → Legacy sees mission progress on bus.

---

## Wave 5 — OSINT & Sovereign Maps — Ongoing

### 5a. omos-osint (Osiris OSINT toolkit as NATS service)

Osiris's 50+ OSINT routes become a NATS request/reply service:

```
svc.osint.flights      → live flights in area
svc.osint.maritime     → ships in area
svc.osint.cctv         → 2K+ cameras near location
svc.osint.shodan       → exposed devices for IP/range
svc.osint.dns          → DNS records
svc.osint.sanctions    → entity sanctions check
svc.osint.threats      → threat intelligence lookup
```

Legacy queries these on demand. Not streaming — on-request only.

### 5b. omos-maps (sovereign tiles)

```yaml
# platform/omos-maps/deployment.yaml
# Serve PMTiles planet file as tile endpoint
# consumed by: ATAK (tile provider), WebTAK, CloudTAK, atoc MapLibre, WWV
image: ghcr.io/onemindos/omos-maps:latest
# Volume: /data/planet.pmtiles (~100GB)
```

WorldMonitor's `basemap-styles.ts` + `basemap.ts` become the tile style config.

### 5c. WorldMonitor geo layers as static overlays

Extract from WorldMonitor:
- `bases-expanded.ts` → GeoServer KML layer (military bases)
- `chokepoint-registry.ts` → GeoServer KML layer (maritime chokepoints)
- `airports.ts` → GeoServer KML layer

These push to TAK via `tak.cot.inject` as static CoT markers OR as KML overlays via TAK Marti API.

---

## Wave 6 — Mac Cleanup (burn build→prove→delete rule)

Only execute AFTER replacement is proven live on K8s.

| Cut | Replaces | Gate |
|-----|----------|------|
| Mac `logtrace-system` (grafana/loki/jaeger) | Cloud observability ✅ | Point exporters at hub, verify dashboards |
| Mac `go-events` stack | EVT stream + ClickHouse | Confirm nothing publishes to it |
| Mac `lets-bet` / `omos-missions` | omos-mission (Wave 4) | Mission confirmed live |
| Mac `peering` stack | omos-media-registry + MediaMTX | Media path designed |
| `nex/` vs `nats/nex/` split | One consolidated NEX root | All workloads migrated |
| `ai-agents/orp/` duplicate | K8s ORP (Wave 1) | ORP confirmed live on K8s |
| `osiris-audit/` clone | osiris/ (kept as reference) | Immediate — byte duplicate |

---

## The Wiring Diagram (final state)

```
MAC (edge/hardware)                    K8s DOKS (always-on cloud)
────────────────────                   ──────────────────────────────────
MAVLink bridge       ┐                 NATS HA (3 nodes) — the one bus
HA bridge            ├── NATS ──────►    │
Vision/cameras       ┘   leaf            ├── ORP (fusion engine)
TAK Server (optional)    node            │     subscribes: tel.>, det.>, evt.world.>, intel.>
                         link            │     publishes:  ent.{type}.{id}.state
                                         │     serves:     svc.orp.query (ORP-QL)
                                         │
                                         ├── om-tak-projector
                                         │     ent.*.*.state → tak.cot.inject
                                         │     → iTAK / ATAK / WebTAK / CloudTAK (ONE CoT feed)
                                         │
                                         ├── om-control-plane
                                         │     ent.*.*.state → ui.* projections → atoc
                                         │
                                         ├── wwv-bridge
                                         │     ent.*.*.state → WebSocket → WWV DataBus
                                         │     → WorldWideView 3D globe
                                         │
                                         ├── Legacy (AI Commander)
                                         │     watches: ent.*, intel.*, det.*, evt.mission.*
                                         │     queries: svc.orp.query (ORP-QL correlation)
                                         │     acts:    cmd.* (approval-gated)
                                         │
                                         ├── omos-feed-world
                                         │     Crucix adapters → evt.world.* (WORLD stream)
                                         │     (planes, ships, fires, weather, satellites)
                                         │
                                         ├── omos-feed-intel
                                         │     WorldMonitor + Crucix → intel.* (INTEL stream)
                                         │     (news, conflict, sanctions, cyber)
                                         │
                                         ├── omos-mission (AMOS patterns)
                                         │     cmd.mission.* → task decomposition → cmd.asset.*
                                         │
                                         ├── ClickHouse
                                         │     ent.*.*.state → entity_history table
                                         │     (time-travel queries, Grafana)
                                         │
                                         └── omos-osint (on-demand)
                                               svc.osint.* request/reply
```

---

## Priority Order (what to do first)

| # | Wave | Task | Time | Blocker? |
|---|------|------|------|---------|
| 1 | 0 | Fill Mac leaf creds + connect | 1h | YES — everything blocks on this |
| 2 | 0 | Provision streams on DOKS | 30min | YES |
| 3 | 1 | Compile ORP with --features nats | 2h | YES — Wave 1 gate |
| 4 | 1 | Add ORP NATS outbound publisher | 4h | YES |
| 5 | 1 | Containerize + deploy ORP | 2h | YES |
| 6 | 2 | Generalize om-tak-projector | 4h | High value — COP |
| 7 | 2 | Deploy om-control-plane | 2h | High value — atoc |
| 8 | 2 | Build + deploy wwv-bridge | 2h | High value — globe |
| 9 | 3 | Build omos-feed-world | 4h | World picture |
| 10 | 3 | Build omos-feed-intel | 4h | Intel picture |
| 11 | 4 | ORP-QL NATS endpoint | 2h | Legacy correlation |
| 12 | 4 | Wire Legacy to entity streams | 2h | Legacy sees world |
| 13 | 4 | Build omos-mission | 8h | Mission capability |
| 14 | 5 | omos-osint service | 4h | OSINT queries |
| 15 | 5 | omos-maps (sovereign tiles) | 4h | Offline maps |
| 16 | 6 | Mac cleanup (per gate) | ongoing | — |

**Total to first unified COP (Waves 0-2):** ~2-3 days of focused work.
**Total to full operating system (Waves 0-4):** ~1-2 weeks.

---

*Execute in wave order. Each wave has a gate condition. Don't skip gates.*

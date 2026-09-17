# OneMind — Wiring Guide
**How to go from current state to a live, unified operating system**
Date: 2026-07-06

---

## Read This First

You have all the pieces. Nothing is missing. The work is wiring, not building from scratch.

**The one principle:** every service follows the same pattern — image on GHCR, manifest in `omos-infra/platform/`, NATS leaf sidecar for bus access, secrets in 1Password. If you understand how Legacy is deployed, you understand how everything is deployed.

### The Bus Model — Lock This In

**K8s NATS is the main bus. The authority. Everything publishes to it.**

The Mac is the cockpit — it has hardware attached to it (drones via MAVLink, Home Assistant via LAN, cameras). Those hardware bridges run on the Mac because they physically have to. But they do not process, store, or make decisions. They receive hardware signals and immediately forward them to the K8s bus via the NATS leaf link. That's the Mac's entire job.

```
COCKPIT (Mac)                       BRAIN (K8s DOKS)
─────────────────                   ──────────────────────────────────
MAVLink bridge  ─┐                  NATS HA — THE BUS (source of truth)
HA bridge       ─┼─ leaf link ───►    ORP (fusion, always-on)
Vision bridge   ─┘  (forward only)    Legacy (AI commander, always-on)
                                      All projectors, connectors, displays
                                      ClickHouse (history)
                                      CloudTAK / WebTAK / WorldWideView
```

Every service you build connects to `nats.nats.svc.cluster.local:4222` (in-cluster, no sidecar needed for pods in K8s). The leaf sidecar is only needed for pods that run on infrastructure *outside* K8s (Mac, edge boxes, future Proxmox). K8s pods talk directly to the cluster NATS service — simpler, faster, no extra container.

**Consequence:** the Mac NATS going offline means you lose MAVLink + HA sensor feeds. The cloud bus, ORP, Legacy, and all displays keep running. The system degrades gracefully. Design everything to assume the cloud bus is always up.

---

## Wave 0 — Connect the Bus
**Time: 1-2 hours | Zero code | Blocks everything else**

### Step 1: Fill Mac leaf credentials

The Mac NATS config has two placeholder values that need real credentials from 1Password.

```bash
# Get credentials
LEAF_USER=$(op read "op://00-24 SO/omos-mac-leaf-node/username")
LEAF_PASSWORD=$(op read "op://00-24 SO/omos-mac-leaf-node/password")

# Fill them in
sed -i '' "s/LEAF_USER/$LEAF_USER/g" ~/nats-server.conf
sed -i '' "s/LEAF_PASSWORD/$LEAF_PASSWORD/g" ~/nats-server.conf

# Restart NATS
launchctl unload ~/Library/LaunchAgents/com.nats.server.plist
launchctl load ~/Library/LaunchAgents/com.nats.server.plist
```

### Step 2: Verify the leaf link is up

```bash
# Check leaf is connected (should show remote DOKS connection)
nats server report jetstream

# Publish from Mac
nats pub tel.test.leaf.ping "mac-online"

# In a second terminal (or from a pod on DOKS):
kubectl port-forward -n nats svc/nats 4222:4222
nats sub "tel.>" --count 1
# Should receive the message you just published from Mac
```

### Step 3: Provision all 38 streams on DOKS

```bash
cd ~/omos-infra
./nats-migration/provision.sh
```

Expected output: `=== provisioned OK=38 FAIL=0 (of 38) ===`

### Gate check

```bash
nats server list
# Should show: zeus-mac (leaf) + 3 DOKS nodes
```

---

## Wave 1 — ORP on the Bus (Fusion Engine)
**Time: 1-2 days | Rust + Docker + K8s | Unlocks all entity fusion**

### What this does

ORP already has a NATS adapter built in — it's just turned off by a compile flag. Turning it on means ORP subscribes to everything on the bus, fuses it into one entity graph, and publishes unified entity state back. After this wave, Legacy and every projector can see ALL entities.

### Step 1: Locate the ORP source

```bash
# Canonical source (has the running binary)
ls ~/omos-core/integrations/ai-agents/orp/

# Check if nats feature exists in Cargo.toml
grep -r "nats" ~/omos-core/integrations/ai-agents/orp/Cargo.toml
```

### Step 2: Compile with NATS enabled

```bash
cd ~/omos-core/integrations/ai-agents/orp

# Build with the NATS adapter unlocked
cargo build --release --features nats

# Verify the binary exists
ls target/release/orp
```

### Step 3: Create ORP connector config

Create `~/omos-core/integrations/ai-agents/orp/config/nats-connectors.yaml`:

```yaml
connectors:
  - id: "onemind-tel"
    type: "nats"
    url: "nats://localhost:4222/stream/TEL/subject/tel.>"
    properties:
      nats_user: "onemind"
      nats_password: "${ONEMIND_PW}"

  - id: "onemind-world"
    type: "nats"
    url: "nats://localhost:4222/stream/WORLD/subject/evt.world.>"
    properties:
      nats_user: "onemind"
      nats_password: "${ONEMIND_PW}"

  - id: "onemind-evt"
    type: "nats"
    url: "nats://localhost:4222/stream/EVT/subject/evt.>"
    properties:
      nats_user: "onemind"
      nats_password: "${ONEMIND_PW}"

  - id: "onemind-intel"
    type: "nats"
    url: "nats://localhost:4222/stream/INTEL/subject/intel.>"
    properties:
      nats_user: "onemind"
      nats_password: "${ONEMIND_PW}"
```

### Step 4: Add NATS outbound publisher to ORP

In ORP's entity update pipeline (`crates/orp-stream/` or wherever entities are written to DuckDB), add a publisher that fires after every entity change:

```rust
// After entity graph update:
let subject = format!("ent.{}.{}.state", entity.entity_type, entity.entity_id);
nc.publish(&subject, serde_json::to_vec(&entity)?).await?;

// Also update the world-state KV for latest-value reads
let kv = js.create_key_value(Config {
    bucket: "world-state".to_string(),
    ..Default::default()
}).await?;
kv.put(&entity.entity_id, serde_json::to_vec(&entity)?).await?;
```

Also add a `svc.orp.query` NATS request/reply endpoint so Legacy can run ORP-QL queries from the bus:

```rust
let sub = nc.subscribe("svc.orp.query").await?;
while let Some(msg) = sub.next().await {
    let query_str = std::str::from_utf8(&msg.payload)?;
    let results = query_engine.execute(query_str).await?;
    if let Some(reply) = msg.reply {
        nc.publish(&reply, serde_json::to_vec(&results)?).await?;
    }
}
```

### Step 5: Write the Dockerfile

Create `omos-infra/platform/orp/Dockerfile`:

```dockerfile
FROM rust:1.78-slim AS builder
RUN apt-get update && apt-get install -y pkg-config libssl-dev
WORKDIR /app
COPY . .
RUN cargo build --release --features nats

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y libssl3 ca-certificates && rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/target/release/orp /usr/local/bin/orp
EXPOSE 9090
CMD ["orp", "serve", "--config", "/etc/orp/config.yaml"]
```

```bash
docker build -t ghcr.io/onemindos/orp:latest \
  -f ~/omos-infra/platform/orp/Dockerfile \
  ~/omos-core/integrations/ai-agents/orp

docker push ghcr.io/onemindos/orp:latest
```

### Step 6: Deploy to K8s

Create `omos-infra/platform/orp/namespace.yaml`:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: orp
```

Create `omos-infra/platform/orp/pvc.yaml`:

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: orp-data
  namespace: orp
spec:
  accessModes: [ReadWriteOnce]
  resources:
    requests:
      storage: 20Gi
```

Create `omos-infra/platform/orp/deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: orp
  namespace: orp
spec:
  replicas: 1
  strategy:
    type: Recreate          # RWO PVC — one pod at a time
  selector:
    matchLabels:
      app: orp
  template:
    metadata:
      labels:
        app: orp
    spec:
      imagePullSecrets:
        - name: ghcr
      containers:
        - name: orp
          image: ghcr.io/onemindos/orp:latest
          imagePullPolicy: Always
          ports:
            - name: http
              containerPort: 9090
          env:
            - name: NATS_URL
              value: "nats://localhost:4222"
          envFrom:
            - secretRef:
                name: orp-env
          volumeMounts:
            - name: data
              mountPath: /var/lib/orp
            - name: config
              mountPath: /etc/orp
          resources:
            requests: { cpu: "500m", memory: "1Gi" }
            limits:   { memory: "4Gi" }
          readinessProbe:
            httpGet: { path: /health, port: 9090 }
            initialDelaySeconds: 15
            periodSeconds: 10

        - name: nats-leaf               # standard sidecar — copy from legacy deployment
          image: nats:2.10-alpine
          args: ["-c", "/etc/nats/leaf.conf"]
          volumeMounts:
            - name: leaf-conf
              mountPath: /etc/nats
              readOnly: true
            - name: nats-ca
              mountPath: /etc/nats-ca
              readOnly: true
          resources:
            requests: { cpu: "50m", memory: "64Mi" }
            limits:   { memory: "128Mi" }

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
        - name: nats-ca
          secret:
            secretName: orp-nats-ca
---
apiVersion: v1
kind: Service
metadata:
  name: orp
  namespace: orp
spec:
  selector:
    app: orp
  ports:
    - name: http
      port: 9090
      targetPort: 9090
```

```bash
kubectl apply -f omos-infra/platform/orp/namespace.yaml
kubectl apply -f omos-infra/platform/orp/pvc.yaml
kubectl apply -f omos-infra/platform/orp/deployment.yaml
```

### Gate check

```bash
# ORP is running
kubectl get pods -n orp

# Entities are flowing on the bus
nats sub "ent.>" --count 5
# Should see JSON entity state messages within seconds of any telemetry arriving

# ORP-QL works
nats req svc.orp.query 'MATCH (e:Ship) WHERE e.is_active = true RETURN e LIMIT 5'
```

---

## Wave 2 — COP Comes Alive
**Time: 1 day | Go + Node.js + K8s | All entities on all displays**

### Step 1: Generalize om-tak-projector

The projector currently only handles air tracks. Change it to handle all entity types.

In `~/omos-core/nex/om-tak-projector/types.go` (or wherever the CoT type map is), replace the air-only map with the full entity-type-to-CoT map:

```go
var entityTypeToCoT = map[string]string{
    // Air
    "aircraft_multirotor": "a-f-A-M-F-Q",
    "aircraft_fixed_wing": "a-f-A-F-F",
    "aircraft_vtol":       "a-f-A-M-V",
    "aircraft_helicopter": "a-f-A-M-H",
    // Ground
    "ground_vehicle_wheeled": "a-f-G-U-C-V",
    "ground_vehicle_tracked": "a-f-G-U-C-R",
    "operator_station":       "a-f-G-U-C-I",
    // Sea
    "surface_vessel_usv":  "a-f-S-X-W-U",
    "underwater_vehicle":  "a-f-S-U",
    // Sensors / infra
    "sensor_platform":  "a-f-G-I",
    "payload_system":   "a-f-G-E",
    // Hazards (non-military entities)
    "fire":             "b-r-.-H-.",
    "weather":          "b-w-.",
    "seismic":          "b-r-.-E-.",
    // Default fallback
    "unknown":          "a-u-G",
}
```

Change the subscription from `ent.air.*.state` to `ent.*.*.state`:

```go
// Before
sub, _ := nc.Subscribe("ent.air.*.state", handler)

// After
sub, _ := nc.Subscribe("ent.*.*.state", handler)
```

Build, containerize, push, deploy to K8s under `platform/omos/om-tak-projector/`.

### Step 2: Deploy om-control-plane to K8s

Move from Mac NEX to a K8s pod. Same pattern: Deployment with NATS leaf sidecar.

```bash
# Build
cd ~/omos-core/nex/om-control-plane
docker build -t ghcr.io/onemindos/om-control-plane:latest .
docker push ghcr.io/onemindos/om-control-plane:latest
```

Manifest in `omos-infra/platform/omos/om-control-plane.yaml` — copy the standard deployment pattern.

### Step 3: Build the WWV NATS bridge

New service (~50 lines) that subscribes NATS `ent.*.*.state` and pushes to WorldWideView's WebSocket endpoint.

Create `~/omos-core/services/wwv-bridge/index.js`:

```javascript
import { connect, StringCodec } from 'nats';
import { WebSocketServer } from 'ws';

const sc = StringCodec();
const nc = await connect({ servers: process.env.NATS_URL || 'nats://localhost:4222' });
const wss = new WebSocketServer({ port: Number(process.env.PORT) || 8765 });

console.log('[wwv-bridge] Connected to NATS, WebSocket on :8765');

const sub = nc.subscribe('ent.*.*.state');

for await (const msg of sub) {
  try {
    const entity = JSON.parse(sc.decode(msg.data));
    if (!entity.latitude || !entity.longitude) continue;  // skip entities without location

    const wwvMsg = JSON.stringify({
      type: 'data',
      pluginId: `onemind-${entity.entity_type || 'unknown'}`,
      payload: [{
        id: entity.entity_id,
        label: entity.name || entity.entity_id,
        position: [entity.longitude, entity.latitude, entity.altitude || 0],
        timestamp: entity.last_updated || new Date().toISOString(),
        properties: {
          type: entity.entity_type,
          status: entity.status,
          confidence: entity.confidence,
          source: entity.source,
          ...entity.properties,
        },
      }],
    });

    wss.clients.forEach(client => {
      if (client.readyState === 1) client.send(wwvMsg);
    });
  } catch (e) {
    console.error('[wwv-bridge] parse error:', e.message);
  }
}
```

Deploy as a K8s service in namespace `omos` with NATS leaf sidecar. WWV connects to `ws://wwv-bridge.omos.svc.cluster.local:8765`.

### Gate check

Open these five simultaneously:
- iTAK (on your phone)
- CloudTAK (cloudtak.onemindos.dev)
- WebTAK (tak.onemindos.dev/webtak)
- atoc (local Flutter app)
- WorldWideView (wwv.onemindos.dev)

Drop a Mac drone MAVLink message. It should appear on all five within ~2 seconds.

---

## Wave 3 — World Data Flows
**Time: 2-3 days | Node.js services | Fills the empty streams**

### Step 1: Build omos-feed-world

Create `~/omos-core/services/omos-feed-world/index.js`.

The pattern for every Crucix adapter:

```javascript
import { connect, StringCodec, JSONCodec } from 'nats';
import { v4 as uuid } from 'uuid';

// Import Crucix adapters directly
import * as opensky from '/path/to/Crucix/apis/sources/opensky.mjs';
import * as firms from '/path/to/Crucix/apis/sources/firms.mjs';
import * as ships from '/path/to/Crucix/apis/sources/ships.mjs';
// ... etc

const nc = await connect({ servers: process.env.NATS_URL });
const jc = JSONCodec();

function envelope(subject, data) {
  return jc.encode({
    v: 1,
    id: uuid(),
    ts: new Date().toISOString(),
    type: subject,
    source: 'omos-feed-world',
    subject,
    data,
  });
}

async function poll(fn, subject, intervalMs) {
  setInterval(async () => {
    try {
      const data = await fn();
      if (data?.error) return console.warn(`[${subject}] error:`, data.error);
      await nc.publish(subject, envelope(subject, data));
      console.log(`[${subject}] published`);
    } catch (e) {
      console.error(`[${subject}] failed:`, e.message);
    }
  }, intervalMs);
}

// Wire each adapter to its subject and polling interval
poll(opensky.briefing,  'evt.world.tel.adsb',         30_000);   // 30s
poll(ships.briefing,    'evt.world.tel.ais',           30_000);
poll(firms.briefing,    'evt.world.det.fire',         300_000);   // 5min
poll(noaa.briefing,     'evt.world.det.weather',      300_000);
poll(space.briefing,    'evt.world.det.space',        900_000);   // 15min
poll(safecast.briefing, 'evt.world.det.environment',  600_000);
poll(epa.briefing,      'evt.world.det.environment',  600_000);
poll(kiwisdr.briefing,  'evt.world.det.sigint',       900_000);
```

Containerize and deploy to `platform/omos/omos-feed-world.yaml`.

Environment secret `omos-feed-world-env` holds API keys:
- `FIRMS_MAP_KEY` — NASA FIRMS (free at firms.modaps.eosdis.nasa.gov)
- `OPENSKY_USER` / `OPENSKY_PASS` — OpenSky (optional, higher rate limit)
- `AIS_API_KEY` — AIS stream (if using paid)

### Step 2: Build omos-feed-intel

Same pattern, intel domain adapters:

```javascript
import * as acled from '/path/to/Crucix/apis/sources/acled.mjs';
import * as gdelt from '/path/to/Crucix/apis/sources/gdelt.mjs';
import * as cisaKev from '/path/to/Crucix/apis/sources/cisa-kev.mjs';
import * as opensanctions from '/path/to/Crucix/apis/sources/opensanctions.mjs';
import * as who from '/path/to/Crucix/apis/sources/who.mjs';
// WorldMonitor news feeds
import { FEEDS } from '/path/to/worldmonitor/src/config/feeds.ts';

poll(acled.briefing,         'intel.conflict.acled',    3600_000);  // 1hr
poll(gdelt.briefing,         'intel.gdelt',              900_000);
poll(cisaKev.briefing,       'intel.cyber.cisa',        3600_000);
poll(opensanctions.briefing, 'intel.sanctions',        86400_000);  // daily
poll(who.briefing,           'intel.health.who',        3600_000);

// WorldMonitor 500+ news feeds — parse RSS, publish each item
async function pollNewsFeeds() {
  for (const feed of FEEDS.slice(0, 100)) {   // start with top 100
    const items = await fetchRSS(feed.url);
    for (const item of items) {
      await nc.publish('intel.news', envelope('intel.news', {
        ...item,
        source: feed.name,
        tier: feed.tier,
        category: feed.category,
      }));
    }
  }
}
setInterval(pollNewsFeeds, 900_000);  // every 15min
```

Deploy as `platform/omos/omos-feed-intel.yaml`.

### Gate check

```bash
# WORLD stream has data
nats stream info WORLD
# Messages: > 0

# INTEL stream has data
nats stream info INTEL
# Messages: > 0

# ORP has fused fire entities (from FIRMS via omos-feed-world)
nats req svc.orp.query 'MATCH (f:Fire) WHERE f.is_active = true RETURN COUNT(f)'
# Returns: > 0

# Aircraft showing on ATAK map
# Open iTAK — should see commercial aircraft overlaid on map
```

---

## Wave 4 — Legacy as Real Commander
**Time: 1 day | Python + K8s | Legacy sees and acts on everything**

### Step 1: Wire Legacy to entity streams

In Legacy's NATS plugin config (in the Hermes config), add stream watchers:

```yaml
# In hermes config (injected via hermes-env secret)
nats_subscriptions:
  - subject: "ent.*.*.state"     # all fused entities
    handler: "entity_update"
  - subject: "intel.>"           # all intel feeds
    handler: "intel_update"
  - subject: "det.>"             # all detections
    handler: "detection_update"
  - subject: "evt.mission.>"     # mission events
    handler: "mission_update"
```

Legacy now passively observes the entire operational picture.

### Step 2: Enable Legacy to query ORP-QL

Add a tool to Legacy's tool catalog (in Hermes tool definitions):

```python
# Tool: query_world
# Description: Query the OneMind entity graph with ORP-QL
# Usage: query_world('MATCH (s:Ship) WHERE NEAR(s, 30.2, -97.7, 50) RETURN s')

async def query_world(ql_query: str) -> dict:
    result = await nc.request("svc.orp.query", ql_query.encode(), timeout=5.0)
    return json.loads(result.data)
```

Now Legacy can answer: "What ships are within 50km of us?" by querying ORP-QL directly.

### Step 3: Build omos-mission

New service that implements AMOS-style mission orchestration over NATS.

Core NATS contract:

```
IN:  cmd.mission.create   → { intent, assets, area, constraints }
OUT: evt.mission.created  → { mission_id, tasks[] }
     cmd.{asset}.{id}.{action}  → task assignments
IN:  evt.mission.task.*   → task completion/failure
OUT: evt.mission.completed / evt.mission.replanning
```

Key patterns from AMOS to implement:
1. **Intent decomposition** — parse mission intent into ordered task list
2. **Asset auction** — query ORP for available assets of required type, assign best fit
3. **Approval gate** — hold high-risk tasks (intercept, lethal) pending human `cmd.mission.approve`
4. **Re-planning** — on task failure, reassign to next best asset

Deploy as `platform/omos/omos-mission.yaml`.

### Gate check

```bash
# Legacy can see all entities
# Ask Legacy: "How many aircraft are currently tracked?"
# Expected: Legacy queries svc.orp.query, returns count from live entity graph

# Trigger a test mission
nats pub cmd.mission.create '{"intent":"patrol","area":{"lat":30.2,"lon":-97.7,"radius_km":10},"assets":["drone"]}'

# Watch mission decompose
nats sub "evt.mission.>" &
nats sub "cmd.drone.>" &

# Should see: mission created → tasks issued → drone receives waypoints
```

---

## Wave 5 — OSINT Service + Sovereign Maps
**Time: 2-3 days | Node.js + tile server | Investigation + offline capability**

### Step 1: omos-osint (Osiris routes as NATS service)

Create `~/omos-core/services/omos-osint/index.js`.

Each Osiris API route becomes a NATS request/reply handler:

```javascript
const handlers = {
  'svc.osint.flights':   async (data) => fetchOsirisRoute('/api/flights', data),
  'svc.osint.maritime':  async (data) => fetchOsirisRoute('/api/maritime', data),
  'svc.osint.cctv':      async (data) => fetchOsirisRoute('/api/cctv', data),
  'svc.osint.shodan':    async (data) => fetchOsirisRoute('/api/osint/shodan', data),
  'svc.osint.dns':       async (data) => fetchOsirisRoute('/api/osint/dns', data),
  'svc.osint.sanctions': async (data) => fetchOsirisRoute('/api/osint/sanctions', data),
  'svc.osint.threats':   async (data) => fetchOsirisRoute('/api/osint/threats', data),
};

for (const [subject, handler] of Object.entries(handlers)) {
  nc.subscribe(subject, async (msg) => {
    const input = jc.decode(msg.data);
    const result = await handler(input);
    msg.respond(jc.encode(result));
  });
}
```

Legacy calls `nats req svc.osint.cctv '{"lat":30.2,"lon":-97.7}'` and gets live CCTV cameras near a location.

### Step 2: omos-maps (sovereign tiles)

Download PMTiles planet file (~100GB, one-time):

```bash
# Download planet tiles (one time)
wget https://r2-public.protomaps.com/protomaps-data/planet-latest.pmtiles \
  -O /data/planet.pmtiles
```

Deploy a tile server pod with the planet file mounted:

Create `platform/omos-maps/deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: omos-maps
  namespace: omos
spec:
  replicas: 1
  selector:
    matchLabels:
      app: omos-maps
  template:
    spec:
      containers:
        - name: tiles
          image: protomaps/go-pmtiles:latest
          args: ["serve", "/data", "--public-url", "https://maps.onemindos.dev"]
          ports:
            - containerPort: 8080
          volumeMounts:
            - name: tiles
              mountPath: /data
      volumes:
        - name: tiles
          persistentVolumeClaim:
            claimName: omos-maps-data  # 150Gi PVC holding planet.pmtiles
```

Configure ATAK, atoc (MapLibre), and WorldWideView to use `https://maps.onemindos.dev/{z}/{x}/{y}.mvt` as the base tile provider. No more Google/Mapbox/OSM dependency at runtime.

### Step 3: WorldMonitor geo layers as TAK overlays

Extract static geo data from WorldMonitor and publish to TAK as CoT markers or KML via Marti API:

```javascript
import { BASES } from '/path/to/worldmonitor/src/config/bases-expanded.ts';
import { CHOKEPOINTS } from '/path/to/worldmonitor/src/config/chokepoint-registry.ts';

// Publish military bases as persistent CoT markers
for (const base of BASES) {
  await nc.publish('tak.cot.inject', buildCoT({
    uid: `OM.base.${base.id}`,
    type: 'a-n-G-I-F-A',  // military airfield
    lat: base.lat, lon: base.lon,
    callsign: base.name,
    layer: 'infrastructure',
  }));
}
```

---

## Wave 6 — Mac Cleanup
**Execute ONLY after replacement is confirmed live on K8s**

Work through this list one item at a time. Each requires your explicit GO.

| What to remove | Replaced by | Confirm before cutting |
|----------------|-------------|------------------------|
| Mac Grafana/Loki/Jaeger stack | Cloud observability (live) | Point any remaining Mac exporters at DOKS |
| `go-events` Postgres/Redis/API | EVT stream + ClickHouse | Nothing publishes to it that isn't on bus |
| `lets-bet` / `omos-missions` stack | omos-mission (Wave 4) | Mission service confirmed live |
| `peering` stack (Redis/MinIO/etc) | omos-media-registry + object stores | Media path designed first |
| `ai-agents/orp/` binary (Mac) | ORP on K8s (Wave 1) | K8s ORP confirmed live |
| `nex/om-unity/` source copy | Single canonical ORP source | Verify it's not a newer fork first |
| `nex/` vs `nats/nex/` split | Consolidate to one NEX root | All workloads migrated |
| `osiris-audit/` clone | `osiris/` (keep one copy) | Immediate — it's a byte duplicate |
| murmur-test + murmur-stable | N/A (debug containers) | Immediate |
| `tools/pulsar/` | om-bridge-mavlink | Immediate — fully superseded |
| `tools/dashboards/` | Running via launchd already | Immediate |
| cot2constellation bridge.py | om-bridge-constellation (Go) | After confirming Go bridge handles all cases |

---

## Adding a New Device or Source (the recipe, any time)

Once Wave 0-2 are live, adding anything new takes minutes:

### New controllable asset (drone, robot, vehicle)

1. Register in `KV_device-registry`:
```bash
nats kv put device-registry "drone.xyz" '{
  "id": "drone.xyz",
  "kind": "drone",
  "name": "DJI Mavic 3",
  "subjects": {
    "telemetry": "tel.drone.xyz.>",
    "commands": "cmd.drone.xyz.>"
  },
  "capabilities": ["fly", "camera", "rtl"],
  "tak": { "cot_type": "a-f-A-M-F-Q", "label": "XYZ" }
}'
```

2. Bridge publishes telemetry to `tel.drone.xyz.>` (use om-bridge-mavlink pattern)
3. Bridge subscribes `cmd.drone.xyz.>` for commands
4. ORP fuses it automatically into the entity graph
5. Appears on all displays within seconds

### New data source (API feed)

1. Copy a Crucix adapter pattern
2. Publish to the right stream subject (`evt.world.*`, `intel.*`, etc.)
3. Wrap in Event Envelope v1
4. ORP fuses it, projectors display it, Legacy can query it

No new infrastructure. No new map. No new display code. The system absorbs it.

---

## Troubleshooting

### Leaf link not connecting
```bash
# Check NATS log on Mac
tail -f ~/nats-data/nats-server.log | grep -i leaf

# Expected: "Leafnode connection created"
# If error: check credentials in ~/nats-server.conf match 1Password
```

### Entity not appearing on ATAK
```bash
# 1. Is it being fused by ORP?
nats req svc.orp.query 'MATCH (e:Entity {id: "your-entity-id"}) RETURN e'

# 2. Is the fused state on the bus?
nats sub "ent.*.your-entity-id.state" --count 1

# 3. Is the tak-projector receiving it?
kubectl logs -n omos deployment/om-tak-projector | grep "your-entity-id"

# 4. Is CoT reaching TAK Server?
nats sub "tak.cot.inject" --count 1
```

### ORP not fusing new source data
```bash
# Check ORP connector status
curl http://orp.orp.svc.cluster.local:9090/api/v1/connectors

# Check ORP logs
kubectl logs -n orp deployment/orp | grep "connector\|error"
```

### Legacy not seeing entities
```bash
# Check Legacy's NATS subscriptions are active
nats sub "agents.hb.hermes.zeus.legacy" --count 1  # Legacy should be heartbeating

# Check entity stream has messages
nats stream info EVT
nats stream info WORLD
```

---

## Quick Reference — Key Subjects

| Subject | What flows | Publisher | Subscribers |
|---------|-----------|-----------|-------------|
| `tel.drone.{id}.>` | drone telemetry | om-bridge-mavlink | ORP |
| `tel.ais.>` | AIS ships | ORP (direct) / omos-feed-world | ORP, om-control-plane |
| `evt.world.>` | world picture | omos-feed-world | ORP |
| `intel.>` | intelligence | omos-feed-intel | ORP, Legacy |
| `det.air.>` | air detections | overwatch sensors | ORP, om-overwatch-state |
| `ent.*.*.state` | fused entity state | ORP | om-tak-projector, om-control-plane, wwv-bridge, Legacy, ClickHouse |
| `tak.cot.inject` | CoT → TAK | om-tak-projector | NatsBridgePlugin |
| `tak.cot.>` | CoT from TAK | cot-tap | ORP |
| `cmd.*.*.>` | commands | Legacy, humans | asset bridges, omos-mission |
| `evt.mission.>` | mission events | omos-mission | Legacy, atoc |
| `svc.orp.query` | ORP-QL query | Legacy, any service | ORP |
| `agents.prompt.hermes.zeus.legacy` | AI prompts | humans, atoc | Legacy |

---

*Execute in wave order. Don't skip gates. Each wave proves the next one works.*

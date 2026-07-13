# OneMind Data Plane — Architecture & Operations

## 1. Overview

The OneMind data plane is the pipeline that moves raw world-state signals from
open APIs into NATS, through a fusion/storage layer, and out to tactical displays
(TAK), AI agents (Legacy/Hermes), and operator UIs (OMOS, WorldWideView).

**Core principle:** NATS is the nervous system; maps are the eyes.  Every piece of
world data — ships, aircraft, seismic events, conflict incidents, weather alerts,
wildfire detections — is normalised into a structured envelope and published to a
subject that reflects what it is.  Consumers subscribe to exactly what they need.

```
┌─────────────────────────────────────────────────────────────┐
│  Layer 1 — Sources                                          │
│  AISstream · OpenSky · APRS · CelesTrak · USGS · NWS       │
│  ACLED · NASA FIRMS · NOAA CAP · OSM / GDACS (future)       │
└───────────────────────┬─────────────────────────────────────┘
                        │  HTTP / WebSocket / TCP
┌───────────────────────▼─────────────────────────────────────┐
│  Layer 2 — Feed Workers (K8s Deployments, namespace: omos)  │
│  feed-aircraft · feed-ships · feed-aprs · feed-satellites   │
│  feed-hazards · feed-intel-noaa · feed-intel-usgs           │
│  feed-intel-acled · feed-intel-firms                        │
└───────────────────────┬─────────────────────────────────────┘
                        │  NATS publish (tak.cot.inject, intel.*)
┌───────────────────────▼─────────────────────────────────────┐
│  Layer 3 — NATS HA Bus (namespace: nats, 3-node)            │
│  JetStream streams: INTEL · TEL · EVT · COT · CMD · ZEUS    │
└──────────┬─────────────────────┬────────────────────────────┘
           │                     │
┌──────────▼──────────┐ ┌────────▼──────────────────────────┐
│  Layer 4 — Sinks    │ │  Layer 4 — Consumers               │
│  sink-clickhouse    │ │  TAK NatsBridgePlugin (takserver)  │
│  sink-tak           │ │  Legacy / Hermes AI (hermes)       │
│                     │ │  ORP fusion engine (orp)           │
│                     │ │  OMOS UI / WWV globe (omos, wwv)   │
└──────────┬──────────┘ └────────────────────────────────────┘
           │
┌──────────▼──────────────────────────────────────────────────┐
│  Layer 5 — Storage / Display                                │
│  ClickHouse onemind.events (90d history)                    │
│  TAK clients: iTAK · ATAK · CloudTAK · WebTAK              │
└─────────────────────────────────────────────────────────────┘
```

---

## 2. Subject Taxonomy

| Prefix | Pattern | Stream | Retention | Description |
|--------|---------|--------|-----------|-------------|
| `tel`  | `tel.{node}.{class}.{metric}` | TEL | 7d | Telemetry (position, sensors) |
| `intel` | `intel.{domain}.{source}.{type}` | INTEL | 90d | Intelligence events (alerts, incidents, detections) |
| `evt`  | `evt.{domain}.{entity}.{event}` | EVT | 30d | Semantic events (mission completed, etc.) |
| `cmd`  | `cmd.{domain}.{entity}.{action}` | CMD | 7d | Commands (intercept, disable, enable) |
| `tak.cot.inject` | (exact) | COT | 24h | CoT XML inbound to TAK NatsBridgePlugin |
| `tak.cot.>` | `tak.cot.*` | COT | 24h | CoT from TAK cot-tap output |
| `ent`  | `ent.{entity_type}.{entity_id}.state` | — (KV) | — | Fused entity state (ORP output) |
| `svc`  | `svc.{service}.{endpoint}` | — | — | NATS RPC (request/reply) |
| `agents` | `agents.{verb}.{agent}.{owner}.*` | ZEUS | 30d | Agent fabric (Synadia protocol) |
| `cmd.etl.sink.tak.*` | exact | CMD | 7d | Sink-tak enable/disable toggle |

**INTEL envelope schema** (JSON, published by all intel feed workers):

```json
{
  "id":         "<string — globally unique event ID>",
  "type":       "<entity_type — weather_alert|earthquake|wildfire|conflict_event>",
  "source":     "<string — noaa|usgs|acled|nasa-firms>",
  "lat":        <float>,
  "lon":        <float>,
  "timestamp":  "<ISO8601 UTC>",
  "expires":    "<ISO8601 UTC — optional>",
  "properties": { "<source-specific fields>" },
  "geometry":   { "<GeoJSON — optional>" }
}
```

---

## 3. Feed Workers

All feed workers are Deployments in namespace `omos`. Scripts live in ConfigMaps
(`omos-feed-scripts` for world feeds, `omos-intel-scripts` for intel feeds).

| Name | Subject Prefix | Source | Poll / Mode | Default Replicas | Toggle Command |
|------|---------------|--------|-------------|-----------------|----------------|
| feed-ships | `tak.cot.inject` | AISstream WebSocket | real-time stream | 1 | `replicas: 0` |
| feed-aircraft | `tak.cot.inject` | OpenSky REST | 15s poll | 1 | `replicas: 0` |
| feed-aprs | `tak.cot.inject` | APRS-IS TCP | real-time stream | 1 | `replicas: 0` |
| feed-satellites | `tak.cot.inject` | CelesTrak TLE + SGP4 | 30s propagate | 1 | `replicas: 0` |
| feed-hazards | `tak.cot.inject` | USGS + NWS REST | 60s poll | 1 | `replicas: 0` |
| feed-intel-noaa | `intel.weather.noaa.alert.*` | NWS CAP API | 300s poll | 1 | `replicas: 0` |
| feed-intel-usgs | `intel.seismic.usgs.*` | USGS GeoJSON | 60s poll | 1 | `replicas: 0` |
| feed-intel-acled | `intel.conflict.acled.*` | ACLED API | 3600s poll | 0 (needs key) | `replicas: 1` |
| feed-intel-firms | `intel.fire.nasa.firms.*` | NASA FIRMS API | 600s poll | 0 (needs key) | `replicas: 1` |

**Secrets required:**
- `omos-nats-creds` (all feeds) — NATS user + password
- `omos-feed-secrets` — `AISSTREAM_API_KEY`, `OPENSKY_USERNAME`, `OPENSKY_PASSWORD`,
  `ACLED_API_KEY`, `ACLED_EMAIL`, `FIRMS_API_KEY`

---

## 4. Sinks

Sinks are durable JetStream push consumers. They ack every message so the stream
cursor advances even when the sink is paused or restarted.

| Name | Consumes | Publishes To | Durable Consumer | Purpose | Toggle |
|------|----------|--------------|-----------------|---------|--------|
| sink-clickhouse | `intel.>` (JetStream INTEL) + `ent.*.*.state` (core) | ClickHouse HTTP | `sink-clickhouse` | Persist all intel + entity state to ClickHouse for query and replay | `replicas: 0` or CH unavailable — batches are discarded gracefully |
| sink-tak | `intel.>` (JetStream INTEL) | `tak.cot.inject` | `sink-tak` | Convert intel envelopes to CoT XML and inject into TAK NatsBridgePlugin | `replicas: 0` or `cmd.etl.sink.tak.disable` |

**sink-tak CoT type mapping:**

| Intel type | CoT type | Notes |
|------------|----------|-------|
| `weather_alert` | `a-f-G-U-W` + sub-type from COT_TYPE_MAP | Full event-name lookup (Tornado, Flood, Hurricane, etc.) |
| `earthquake` | `a-f-G-E-V-Q` | Seismic — USGS/EMSC |
| `wildfire` | `a-f-G-U-F-F` | Fire — NASA FIRMS / CAL FIRE |
| `conflict_event` | `a-h-G` | Hostile ground — ACLED |
| _(default)_ | `a-f-G-U-O` | Generic observation |

---

## 5. The Flip of a Switch

Every feed and sink can be stopped instantly without touching code.

**Hard stop (K8s scale):**
```bash
# Stop a feed
kubectl scale deploy/feed-intel-noaa -n omos --replicas=0

# Stop the TAK sink
kubectl scale deploy/sink-tak -n omos --replicas=0

# Resume
kubectl scale deploy/sink-tak -n omos --replicas=1
```

**Soft pause via NATS command (sink-tak only):**
```bash
# Pause forwarding — sink continues running, acks messages without publishing CoT
nats pub cmd.etl.sink.tak.disable "" --server nats://nats.nats.svc.cluster.local:4222

# Resume
nats pub cmd.etl.sink.tak.enable "" --server nats://nats.nats.svc.cluster.local:4222
```

The soft pause is useful for temporary suppression (e.g., maintenance window on TAK
Server) without losing the durable consumer position in the INTEL stream.

---

## 6. ClickHouse Time Machine

Schema: `onemind.events` (MergeTree, partitioned by month, ordered by timestamp + entity_id)

```sql
CREATE TABLE onemind.events (
    timestamp    DateTime64(3, 'UTC'),
    subject      String,
    entity_id    String,
    entity_type  LowCardinality(String),
    source       LowCardinality(String),
    lat          Nullable(Float64),
    lon          Nullable(Float64),
    payload      String           -- raw JSON envelope
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(timestamp)
ORDER BY (timestamp, entity_id)
TTL timestamp + INTERVAL 90 DAY;
```

**Example queries:**

```sql
-- All seismic events in the last 30 days
SELECT timestamp, entity_id, lat, lon, payload
FROM onemind.events
WHERE entity_type = 'earthquake'
  AND timestamp >= now() - INTERVAL 30 DAY
ORDER BY timestamp DESC;

-- Weather alerts near a point (lat 34.0, lon -118.2) within 200 km
SELECT timestamp, entity_id, source, payload
FROM onemind.events
WHERE entity_type = 'weather_alert'
  AND greatCircleDistance(lon, lat, -118.2, 34.0) < 200000
  AND timestamp >= now() - INTERVAL 7 DAY
ORDER BY timestamp DESC;

-- Conflict events by source over the last 90 days
SELECT source, count() AS cnt, max(timestamp) AS last_seen
FROM onemind.events
WHERE entity_type = 'conflict_event'
  AND timestamp >= now() - INTERVAL 90 DAY
GROUP BY source
ORDER BY cnt DESC;

-- Replay all events on a subject prefix for a time window
SELECT timestamp, subject, payload
FROM onemind.events
WHERE subject LIKE 'intel.fire.%'
  AND timestamp BETWEEN '2026-06-01' AND '2026-06-30'
ORDER BY timestamp;
```

---

## 7. TAK Integration

```
intel.> (INTEL JetStream, 90d)
    │
    └── sink-tak Deployment (omos namespace)
            │   durable consumer "sink-tak"
            │   CoT XML built per envelope
            ▼
        tak.cot.inject (COT stream, 24h)
            │
            └── NatsBridgePlugin (takserver namespace)
                    │   subscribes tak.cot.inject
                    │   injects into TAK Server pipeline
                    ▼
                TAK Server 5.7 (takserver namespace)
                    │
                    ├── iTAK (iOS)
                    ├── ATAK (Android)
                    ├── CloudTAK (atoc.onemindos.dev)
                    └── WebTAK (browser)
```

CoT events appear on TAK clients colour-coded red (`argb="-65536"`) with the
`remarks` block containing the headline, source name, and original event ID.
Stale time defaults to 30 minutes from event timestamp; feeds that provide an
`expires` field in the envelope use that value instead.

The NatsBridgePlugin also taps `tak.cot.>` outbound — CoT from EUDs flows back
onto NATS and can be consumed by ORP or Legacy.

---

## 8. Legacy Integration

Legacy (Hermes AI, namespace `hermes`) consumes the data plane in three ways:

1. **Live subscription** — Legacy subscribes to `intel.>` and `ent.*.*.state` for
   situational awareness. Incoming events feed its context window for threat
   assessment, correlation, and commander briefings.

2. **ClickHouse query** — Legacy issues SQL queries via ClickHouse HTTP (`CH_HOST`
   in its env) for historical analysis: "all conflict events near X in the last 30
   days", "earthquake trend over 90 days", etc.

3. **Command publication** — Legacy publishes `cmd.*` subjects to direct feeds and
   sinks. Examples:
   - `cmd.etl.sink.tak.disable` — suppress TAK forwarding during a briefing
   - `cmd.drone.x47.intercept` — task a drone (downstream handler in ORP)
   - `cmd.mission.{id}.abort` — abort a mission workflow

Legacy's NATS identity is `agents.prompt.legacy.zeus` (Synadia Agent Protocol).
It connects via a leaf-node sidecar in its pod so it has an isolated NATS identity
with scoped permissions — it cannot accidentally publish to `tak.cot.inject` directly;
it must go through the sink pipeline or issue a command that triggers one.

---

## 9. Adding a New Feed

- [ ] **Decide the subject prefix** — follow the grammar: `intel.{domain}.{source}.{type}`
- [ ] **Write the feed script** — add `feed_intel_{name}.py` to the `omos-intel-scripts`
      ConfigMap in `platform/omos/intel-scripts.yaml`. Use the INTEL envelope schema.
- [ ] **Add the Deployment** — append to `platform/omos/intel-feeds.yaml`. Copy an
      existing block (e.g. `feed-intel-noaa`). Set `replicas: 0` if an API key is needed.
- [ ] **Add secrets** — if the feed needs an API key, add the key to `omos-feed-secrets`
      in 1Password, and wire it via `secretKeyRef` in the Deployment env.
- [ ] **Verify the INTEL stream filter** — the INTEL stream in NATS must cover the new
      subject prefix. Check `nats-migration/streams/intel.json`; update if needed and
      re-run `provision.sh`.
- [ ] **Update this doc** — add a row to the Feed Workers table (§3) and the Sinks table
      if a new sink is also being added.
- [ ] **Commit and push** — ArgoCD syncs `platform/omos/` automatically (`recurse: false`,
      same flat directory). New Deployments and ConfigMap changes deploy within ~30s.

---

## 10. Secrets

All secrets are managed by External Secrets Operator + 1Password Connect.
Source of truth: vault `00-24 SO (Sovereign Operations)`.

| 1Password Item (= K8s secret name) | Namespace | Keys | Used By |
|------------------------------------|-----------|------|---------|
| `omos-nats-creds` | omos | `user`, `password` | All omos feed workers and sinks |
| `omos-feed-secrets` | omos | `AISSTREAM_API_KEY`, `OPENSKY_USERNAME`, `OPENSKY_PASSWORD`, `ACLED_API_KEY`, `ACLED_EMAIL`, `FIRMS_API_KEY` | World and intel feed workers |
| `omos-clickhouse-creds` | omos | `CH_PASSWORD` | sink-clickhouse |
| `orp-nats-creds` | orp | `user`, `password` | ORP fusion engine |
| `hermes-env` | hermes | AWS Bedrock keys, NATS creds, 20+ fields | Legacy/Hermes AI |
| `tak-nats-leaf` | takserver | `leaf.conf`, `ca.crt` | TAK NatsBridgePlugin NATS leaf |
| `cloudtak-secrets` | cloudtak | signing secret, MinIO keys, PG DSN | CloudTAK stack |

To add or rotate a secret: update the 1Password item → ESO syncs within 1 hour, or
force-sync immediately:
```bash
kubectl annotate externalsecret omos-nats-creds -n omos \
  force-sync="$(date +%s)" --overwrite
```

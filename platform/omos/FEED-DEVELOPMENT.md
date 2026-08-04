# OMOS Feed Development Guide

> How to add a new real-time data source to the OneMind Intelligence Fabric.
> Every feed follows the same pattern: **fetch → normalize to CoT XML → publish to NATS**.

---

## Architecture (One Picture)

```
  External API                Feed Pod                    NATS Bus              Sinks
  (REST/WS/TCP)          (python:3.12-alpine)          (tak.cot.inject)     (already running)
                               │                          │
  ┌─────────┐            ┌─────┴──────┐            ┌───────┴───────┐    ┌────────────┐
  │ MTA GTFS │──────►    │  feed-mta  │──────►     │  tak.cot.inject│───►│  sink-tak  │──► TAK Server → all clients
  │ OpenSky  │──────►    │  feed-air  │──────►     │  (JetStream)  │───►│  sink-ch   │──► ClickHouse Time Machine
  │ AIS WS   │──────►    │  feed-ais  │──────►     │              │───►│  Legacy AI  │──> I watch everything
  └─────────┘            └────────────┘            └──────────────┘    └────────────┘
```

**The feed pod NEVER talks to TAK Server directly.** It only publishes CoT XML to NATS. The existing `sink-tak` pod handles delivery to TAK Server. The existing `sink-clickhouse` pod handles persistence to ClickHouse.

---

## Files To Create Per Feed

| File | What It Is |
|------|-----------|
| `platform/omos/feed-{name}-script.yaml` | ConfigMap containing the Python script |
| `platform/omos/feed-{name}.yaml` | Deployment manifest |

That's it. Two files. No new Docker image, no new repo.

---

## Step-by-Step: Adding a New Feed

### Step 1 — Create the Script ConfigMap

Create `platform/omos/feed-{name}-script.yaml`:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: omos-feed-{name}-script
  namespace: omos
data:
  feed_{name}.py: |
    #!/usr/bin/env python3
    """feed-{name} — {description} → CoT XML → NATS"""
    import asyncio, logging, os, ssl
    from datetime import datetime, timezone, timedelta
    import nats

    LOG = logging.getLogger("feed-{name}")
    logging.basicConfig(level=logging.INFO,
        format="%(asctime)s %(levelname)s [%(name)s] %(message)s")

    NATS_URL      = os.getenv("NATS_URL", "nats://nats.nats.svc.cluster.local:4222")
    NATS_USER     = os.getenv("NATS_USER", "")
    NATS_PASSWORD = os.getenv("NATS_PASSWORD", "")
    POLL_INTERVAL = int(os.getenv("POLL_INTERVAL", "30"))
    INJECT_SUBJ   = os.getenv("NATS_SUBJECT", "tak.cot.inject")
    XML_DECL      = b'<?xml version="1.0" encoding="UTF-8" standalone="yes" ?>\n'

    def now_str(): return datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%S.%fZ")
    def stale_str(s=120): return (datetime.now(timezone.utc) + timedelta(seconds=s)).strftime("%Y-%m-%dT%H:%M:%S.%fZ")

    def build_cot(uid, callsign, lat, lon, cot_type="a-n-G-E-V-C", remarks="", speed=0, course=0, stale=120):
        ts, st = now_str(), stale_str(stale)
        xml = (
            f'<event version="2.0" type="{cot_type}" uid="{uid}" '
            f'how="m-g" time="{ts}" start="{ts}" stale="{st}">'
            f'<point lat="{lat:.6f}" lon="{lon:.6f}" hae="0" ce="50" le="50"/>'
            f'<detail>'
            f'<contact callsign="{callsign}"/>'
            f'<track speed="{speed:.1f}" course="{course:.1f}"/>'
            f'<remarks>{remarks}</remarks>'
            f'<marti><dest mission="Land"/></marti>'
            f'</detail></event>'
        )
        return XML_DECL + xml.encode()

    async def main():
        tls = ssl.create_default_context(); tls.check_hostname = False; tls.verify_mode = ssl.CERT_NONE
        nc = await nats.connect(servers=NATS_URL, user=NATS_USER, password=NATS_PASSWORD, tls=tls)
        LOG.info(f"NATS connected → {INJECT_SUBJ}")
        while True:
            # 1. FETCH from your data source (REST, WebSocket, TCP, etc.)
            # 2. PARSE / NORMALIZE to lat, lon, callsign, etc.
            # 3. CONVERT to CoT XML via build_cot()
            # 4. PUBLISH to NATS
            #    await nc.publish(INJECT_SUBJ, cot_bytes)
            await asyncio.sleep(POLL_INTERVAL)

    asyncio.run(main())
```

### Step 2 — Create the Deployment

Create `platform/omos/feed-{name}.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: feed-{name}
  namespace: omos
  labels:
    app.kubernetes.io/name: feed-{name}
    app.kubernetes.io/part-of: omos-feeds
    omos.io/feed-type: {transit|air|sea|ground|space|events}
    omos.io/feed-source: {source-name}
    omos.io/requires-key: "true|false"
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: feed-{name}
  template:
    metadata:
      labels:
        app.kubernetes.io/name: feed-{name}
        app.kubernetes.io/part-of: omos-feeds
    spec:
      containers:
        - name: feed
          image: python:3.12-alpine
          command: ["/bin/sh", "-c"]
          args:
            - |
              pip install --no-cache-dir --quiet nats-py {extra-deps} \
              && exec python3 /scripts/feed_{name}.py
          env:
            - name: NATS_URL
              value: "nats://nats.nats.svc.cluster.local:4222"
            - name: NATS_USER
              valueFrom:
                secretKeyRef:
                  name: omos-nats-creds
                  key: user
            - name: NATS_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: omos-nats-creds
                  key: password
            # Add source-specific env vars here (API keys, URLs, filters)
            # Keys come from k8s secrets, never hardcoded:
            # - name: MY_API_KEY
            #   valueFrom:
            #     secretKeyRef:
            #       name: omos-feed-secrets
            #       key: MY_API_KEY
          volumeMounts:
            - name: scripts
              mountPath: /scripts
              readOnly: true
          resources:
            requests:
              cpu: "50m"
              memory: "64Mi"
            limits:
              cpu: "200m"
              memory: "128Mi"
      volumes:
        - name: scripts
          configMap:
            name: omos-feed-{name}-script
```

### Step 3 — Deploy

```bash
kubectl apply -f platform/omos/feed-{name}-script.yaml
kubectl apply -f platform/omos/feed-{name}.yaml
kubectl logs -f deploy/feed-{name} -n omos
```

### Step 4 — Verify on TAK

Open ATAK → look for new markers on the map. Check Layer Manager for the CoT type.

---

## Managing Feeds

### Toggle (enable/disable)

```bash
# Disable a feed (stops pod, no more CoT from this source)
kubectl scale deploy/feed-{name} -n omos --replicas=0

# Re-enable
kubectl scale deploy/feed-{name} -n omos --replicas=1
```

### List All Feeds

```bash
kubectl get deployments -n omos -l app.kubernetes.io/part-of=omos-feeds
```

### View Logs

```bash
kubectl logs -f deploy/feed-{name} -n omos
```

### Update a Feed Script

```bash
# Edit the ConfigMap
kubectl edit configmap omos-feed-{name}-script -n omos
# Restart the pod to pick up changes
kubectl rollout restart deploy/feed-{name} -n omos
```

### Add an API Key

1. Store in k8s secret (via 1Password → ESO, or direct):
```bash
kubectl create secret generic omos-feed-secrets \
  --from-literal=MY_API_KEY=xxx -n omos --dry-run=client -o yaml | kubectl apply -f -
```
2. Reference in the Deployment env section (see template above).

---

## CoT Type Reference

Choose the right CoT type for your feed:

| CoT Type | Meaning | Use For |
|----------|---------|---------|
| `a-n-G-E-V-C-R` | Neutral ground vehicle civilian rail | Trains (MTA, Amtrak) |
| `a-n-G-E-V-C` | Neutral ground vehicle civilian | Buses, cars, trucks |
| `a-f-A-M-F-Q` | Friendly air military fixed-wing rotary | Aircraft (helicopter) |
| `a-f-A-M-F` | Friendly air military fixed-wing | Aircraft (plane) |
| `a-u-A` | Unknown air | Unidentified aircraft |
| `a-f-S-X-M` | Friendly surface maritime | Ships (friendly) |
| `a-u-S-X-M` | Unknown surface maritime | Ships (unknown) |
| `a-f-A-M-F-Q` | Friendly air military | Satellite |
| `b-r-f-h-h` | Raw feed hazard | Earthquakes, weather alerts |
| `b-m-p-w` | Military point warning | Conflict events |
| `b-t-f` | Threat fire | Fire/thermal anomalies |
| `b-n-g` | Natural ground event | Disasters |

### CoT XML Structure

```xml
<event version="2.0"
       type="{cot_type}"
       uid="{unique-id}"
       how="m-g"                        <!-- m-g = machine-generated -->
       time="{ISO8601}"
       start="{ISO8601}"
       stale="{ISO8601}">
  <point lat="{lat}" lon="{lon}" hae="0" ce="50" le="50"/>
  <detail>
    <contact callsign="{callsign}"/>
    <track speed="{knots}" course="{degrees}"/>
    <remarks>{description}</remarks>
    <marti><dest mission="{MissionName}"/></marti>   <!-- routes to specific TAK mission -->
  </detail>
</event>
```

---

## Mission Routing

The `<marti><dest mission="..."/>` tag controls which TAK mission the CoT appears in:

| Mission | Use For |
|---------|---------|
| `Air` | Aircraft, satellites |
| `Sea` | Ships, maritime |
| `Land` | Trains, vehicles, APRS |
| `Hazards` | Earthquakes, weather |
| `Intel` | Conflict events, threats |
| `Drones` | Drone telemetry |
| `Team` | Team member positions |

---

## Feed Inventory (Current)

| Feed | File | Source | CoT Type | Mission | Key Needed | Status |
|------|------|--------|----------|---------|------------|--------|
| feed-mta | feed-mta.yaml | MTA GTFS-RT | a-n-G-E-V-C-R | Land | No | ✅ Live |
| feed-aircraft | feeds.yaml | OpenSky REST | a-f-A-M-F-Q | Air | Optional | ⏸️ replicas:0 |
| feed-aprs | feeds.yaml | APRS-IS TCP | a-f-G-U-C | Land | No | ⏸️ replicas:0 |
| feed-satellites | feeds.yaml | CelesTrak TLE+SGP4 | a-f-A-M-F-Q | Air | No | ✅ Live |
| feed-hazards | feeds.yaml | USGS+NWS | b-r-f-h-h | Hazards | No | ✅ Live |
| feed-intel-acled | intel-feeds.yaml | ACLED REST | b-m-p-w | Intel | Yes (email) | ✅ Live |
| feed-intel-emsc | intel-feeds.yaml | EMSC REST | b-r-f-h-h | Hazards | No | ✅ Live |
| feed-intel-eonet | intel-feeds.yaml | NASA EONET | b-n-g | Hazards | No | ✅ Live |
| feed-intel-firms | intel-feeds.yaml | NASA FIRMS | b-t-f | Hazards | Optional | ✅ Live |
| feed-intel-gdacs | intel-feeds.yaml | GDACS RSS | b-n-g | Hazards | No | ✅ Live |
| feed-intel-noaa | intel-feeds.yaml | NOAA NWS | b-r-f-h-h | Hazards | No | ✅ Live |
| feed-intel-usgs | intel-feeds.yaml | USGS EQ | b-r-f-h-h | Hazards | No | ✅ Live |

## Feed Inventory (Planned)

| Feed | Source | Key Needed | Phase |
|------|--------|------------|-------|
| feed-ais | aisstream.io WebSocket | Free key | Phase 1 |
| feed-aircraft-mil | airplanes.live | No | Phase 1 |
| feed-aircraft (restore) | OpenSky authenticated | Free account | Phase 1 |
| feed-aprs (restore) | APRS-IS | No | Phase 1 |
| feed-intel-oref | Israel Home Front | No | Phase 2 |
| feed-intel-gpsjam | GPSJam | No | Phase 2 |
| feed-conflict-gdelt | GDELT | No | Phase 2 |
| feed-aircraft-wingbits | Wingbits | Paid key | Phase 3 |

---

## Delta Filtering & Load Management

Every feed should implement delta filtering to avoid flooding NATS:

```python
# Only republish if moved > 50 meters OR stale window approaching
def should_publish(uid, lat, lon, delta_meters=50, stale_sec=120):
    prev = last_positions.get(uid)
    now = time.time()
    if prev is None:
        last_positions[uid] = (lat, lon, now)
        return True
    if now - prev[2] > stale_sec * 0.7:  # keep alive
        last_positions[uid] = (lat, lon, now)
        return True
    if haversine(prev[0], prev[1], lat, lon) >= delta_meters:
        last_positions[uid] = (lat, lon, now)
        return True
    return False
```

---

## Secrets (API Keys)

All keys stored in k8s secret `omos-feed-secrets` in `omos` namespace:

```bash
# Add a new key
kubectl create secret generic omos-feed-secrets \
  --from-literal=AISSTREAM_API_KEY=xxx \
  --from-literal=OPENSKY_USERNAME=xxx \
  -n omos --dry-run=client -o yaml | kubectl apply -f -

# For production: use External Secrets Operator (ESO) + 1Password
# See platform/external-secrets/ for the pattern
```

---

## Troubleshooting

| Symptom | Check |
|---------|-------|
| No markers on TAK | `kubectl logs deploy/feed-{name} -n omos` for errors |
| Pod CrashLoopBackOff | Check pip install line — missing dependency? |
| NATS connection failed | Verify `omos-nats-creds` secret exists in omos namespace |
| Rate limited by source | Increase POLL_INTERVAL or add auth credentials |
| Too many markers / lag | Enable delta filtering, increase DELTA_METERS |
| Markers don't disappear | Decrease STALE_SEC — stale markers auto-expire on TAK |

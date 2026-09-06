# OneTAK — Offline-First / Native Client Plan

**Audit date:** 2026-09-06
**CloudTAK version:** 13.78.0 (running in cluster as `ghcr.io/onemindos/cloudtak-api:latest`)

---

## What We're Building

A native desktop/mobile wrapper around CloudTAK that works fully offline:
- Map tiles available with no internet
- SA picture (units, tracks, features) persists when API is unreachable
- App shell loads instantly from cache
- CoT messages queue locally and sync on reconnect

The two upstream repos that give us the foundation:

| Repo | Version | What it does |
|------|---------|-------------|
| `dfpc-coe/mbtiles-offline` | v1.3.0 | Bulk-downloads tiles from any ZXY URL → MBTiles SQLite |
| `dfpc-coe/CloudTAK-Native` | v1.1.0 | Electron + Capacitor shell that loads CloudTAK as a native app |

---

## Current Cluster State

```
Namespace: cloudtak

cloudtak-api        ghcr.io/onemindos/cloudtak-api:latest       (SPA + REST API)
cloudtak-events     ghcr.io/onemindos/cloudtak-events:latest     (WebSocket)
cloudtak-tiles      ghcr.io/onemindos/cloudtak-tiles:latest      (tile server → MinIO)
cloudtak-minio      MinIO                                         (object store, 20Gi PVC)
cloudtak-postgis    PostGIS                                       (feature/entity store, 10Gi PVC)
cloudtak-media      media server                                  (RTSP/RTMP)
```

Tile server reads from MinIO bucket `cloudtak` via S3 API.
Tiles are served as PMTiles from `http://cloudtak-minio:9000`.

---

## What's Already Done ✅

### Service Worker (app shell offline)
CloudTAK already ships a production-grade service worker:

- **`/home/etl/api/web/dist/sw.js`** — precaches all Vite assets (JS, CSS, web workers) on first load
- **`/home/etl/api/web/src/utils/service-worker.ts`** — registration utility that:
  - Fingerprints `/.vite/manifest.json` to detect new deploys
  - Registers with a `build` query param for cache versioning (no stale collisions)
  - Handles offline starts explicitly — if manifest unreachable, existing SW keeps serving
  - Polls every 5 minutes for new deploys in long-lived tabs
  - Has Capacitor/native detection via `supportsServiceWorker()` — skips SW on native platforms that don't support it

**This means:** once a user loads CloudTAK once in Electron or a browser, the SPA loads fully offline. No work needed here.

Glyph fonts are cached separately in `cloudtak-glyphs` (unversioned) so map labels survive offline too.

---

## What Needs to Be Built ❌

### 1. Offline Tiles (most impactful — 1–2 days)

**Problem:** `cloudtak-tiles` serves from MinIO. MinIO has no tiles until they're seeded. Offline → blank map.

**Fix:** Run `mbtiles-offline` as a K8s Job to pre-seed an Area of Responsibility (AOR).

```yaml
# platform/cloudtak/jobs/seed-tiles-aor.yaml  (to be created)
apiVersion: batch/v1
kind: Job
metadata:
  name: seed-tiles-aor
  namespace: cloudtak
spec:
  template:
    spec:
      containers:
        - name: seed
          image: node:22-alpine
          command: ["sh", "-c"]
          args:
            - |
              npm install -g @tak-ps/mbtiles-offline
              mbtiles-offline \
                --bbox "-125,24,-66,49" \
                --min-zoom 0 \
                --max-zoom 14 \
                --url "https://api.maptiler.com/tiles/satellite/{z}/{x}/{y}.jpg?key=YOUR_KEY" \
                --output /tmp/usa-satellite-z14.mbtiles
              # Upload to MinIO
              mc alias set minio http://cloudtak-minio:9000 $MINIO_USER $MINIO_PASS
              mc cp /tmp/usa-satellite-z14.mbtiles minio/cloudtak/tiles/usa-satellite-z14.mbtiles
          env:
            - name: MINIO_USER
              valueFrom:
                secretKeyRef:
                  name: cloudtak-secrets
                  key: AWS_S3_AccessKeyId
            - name: MINIO_PASS
              valueFrom:
                secretKeyRef:
                  name: cloudtak-secrets
                  key: AWS_S3_SecretAccessKey
      restartPolicy: Never
```

Once the `.mbtiles` is in MinIO, `cloudtak-tiles` can serve it via the PMTiles driver with no upstream tile source needed.

**For the native app:** The seeded `.mbtiles` file can also be bundled with or downloaded to the Electron app and served by a local `martin` sidecar — fully air-gapped tiles.

**Hermes cron option:** Instead of a one-shot Job, Hermes can run `mbtiles-offline` on a schedule to keep tiles fresh for a dynamic AOR (e.g., current mission bounding box from NATS `cmd.mission.*.bbox`).

---

### 2. Native Wrapper (easy — half day)

**Problem:** CloudTAK runs in a browser. We want a native desktop app (macOS/Windows) and mobile app (iOS/Android).

**Fix:** Fork `dfpc-coe/CloudTAK-Native` and point it at our instance.

**What to change in the fork:**
```javascript
// index.js — change default server URL
const DEFAULT_SERVER = 'https://onemindos.dev';

// capacitor.config.ts — change server URL
const config: CapacitorConfig = {
  appId: 'os.onemind.tak',
  appName: 'OneTAK',
  server: {
    url: 'https://onemindos.dev',
    cleartext: false,
  },
};
```

**Key behaviours already in CloudTAK-Native:**
- GPU/WebGL enabled via `--ignore-gpu-blocklist` (MapLibre/CesiumJS work)
- Geolocation auto-granted (no browser prompt)
- `X-Frame-Options` headers stripped
- `BrowserWindow` opens full-screen with standard TAK-style chrome

**Service worker in Electron:** Works normally — Chromium in Electron supports SW. The existing CloudTAK SW will precache the app shell on first launch.

**iOS/Android via Capacitor:** SW is skipped (Capacitor WebView limitation). Tiles and features need native caching via Capacitor's Filesystem + HTTP plugin. This is Phase 2.

---

### 3. Offline SA Picture — IndexedDB Entity Cache (harder — 3–5 days)

**Problem:** When `cloudtak-api` is unreachable, the map goes blank. Features, tracks, and unit positions all come from API calls to PostGIS. The service worker caches the app shell but not API responses.

**Fix:** Add IndexedDB snapshotting of entity state in the CloudTAK SPA fork.

**Architecture:**

```
Online path (unchanged):
  api.feature.stream() → MapLibre source → map renders

Offline path (new):
  api.feature.stream() → IndexedDB snapshot → fallback source on disconnect
```

**What to add in the CloudTAK fork:**

```typescript
// src/utils/offline-store.ts  (new file)
import { openDB } from 'idb';

const DB_NAME = 'cloudtak-offline';
const ENTITY_STORE = 'entities';

export async function snapshotEntity(id: string, feature: GeoJSON.Feature) {
    const db = await openDB(DB_NAME, 1, {
        upgrade(db) { db.createObjectStore(ENTITY_STORE, { keyPath: 'id' }); }
    });
    await db.put(ENTITY_STORE, { id, feature, ts: Date.now() });
}

export async function loadOfflineEntities(): Promise<GeoJSON.Feature[]> {
    const db = await openDB(DB_NAME, 1);
    const all = await db.getAll(ENTITY_STORE);
    return all.map(r => r.feature);
}
```

**Integration points in the SPA:**
- Hook into `api.feature.stream()` to snapshot each entity as it arrives
- On `offline` event / API disconnect: switch MapLibre source to IndexedDB snapshot
- On reconnect: resume live stream, merge delta

**CoT outbound queue (bonus):**

```typescript
// src/utils/cot-queue.ts  (new file)
// Buffer outbound CoT when NATS WS is disconnected
export const cotQueue: string[] = [];

export function enqueueCoT(cot: string) { cotQueue.push(cot); }

// Called on NATS reconnect
export async function drainCoTQueue(nats: NatsConnection) {
    while (cotQueue.length) {
        await nats.publish('tak.cot.inject', cotQueue.shift()!);
    }
}
```

---

## Implementation Roadmap

| Phase | Work | Effort | Result |
|-------|------|--------|--------|
| **1a** | `mbtiles-offline` K8s Job — seed CONUS tiles to MinIO | 1 day | Tiles work offline |
| **1b** | Fork CloudTAK-Native — point at `onemindos.dev`, build DMG | 0.5 day | Native desktop app |
| **2** | IndexedDB entity snapshot in CloudTAK SPA fork | 3–5 days | SA picture offline |
| **3** | CoT outbound queue (NATS WS reconnect drain) | 1 day | Full write-offline |
| **4** | Capacitor iOS/Android build + native tile caching | 1 week | Mobile offline |

Phases 1a + 1b = fully usable native app with tiles and app shell offline.
Phase 2 = true hybrid (read/write SA offline).

---

## Key File Locations

| File | Location |
|------|----------|
| CloudTAK service worker (built) | `/home/etl/api/web/dist/sw.js` (in pod) |
| SW registration utility (source) | `/home/etl/api/web/src/utils/service-worker.ts` |
| Tile server deployment | `platform/cloudtak/tiles.yaml` |
| MinIO deployment | `platform/cloudtak/minio.yaml` |
| Tile seed Job (to create) | `platform/cloudtak/jobs/seed-tiles-aor.yaml` |
| CloudTAK-Native fork target | `https://github.com/dfpc-coe/CloudTAK-Native` |
| mbtiles-offline package | `https://github.com/dfpc-coe/mbtiles-offline` |

---

## Notes

- `cloudtak-tiles` already supports PMTiles/MBTiles sources via MinIO — no tile server changes needed, just seed the files
- The service worker uses Vite manifest fingerprinting for cache busting — this works correctly in Electron
- `supportsServiceWorker()` in `service-worker.ts` checks for Capacitor native context — returns false on iOS/Android, so those need a separate offline strategy (Phase 4)
- For Hermes integration: Hermes can trigger tile seeding via `cmd.mission.*.bbox` NATS subject when a new mission AOR is defined — fully automated offline prep before going to the field

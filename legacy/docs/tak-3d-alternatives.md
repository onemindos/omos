# VR-TAK Alternatives: Open-Source 3D TAK-Integrated Mission Planning

**Target Hardware:** Dell Inspiron 24-7459 | NVIDIA GeForce 940M (2GB GDDR5, ~GTX 750 Ti class) | Intel i7-6700HQ | 16 GB RAM | Windows 10  
**TAK Server:** TAK Server 5.7 on DOKS | IP: 152.42.155.95 | Ports: 8089 (CoT TLS), 8443 (WebTAK), 8446 (enrollment)  
**Research Date:** 2026-07-14

---

## What VR-TAK Provides (Baseline)

VR-TAK (ARA, Inc.) is a proprietary VR headset application (Oculus/Meta Quest, potentially PCVR) that ingests live CoT feeds from TAK Server and renders a 3D immersive mission planning environment. It requires a dedicated VR headset and is sold commercially — not open-source, not publicly priced. It is the gold standard for 3D TAK visualization but is hardware-locked to VR equipment.

---

## Option 1: tak-webview-cesium ⭐ BEST MATCH

| Field | Value |
|-------|-------|
| **Name** | tak-webview-cesium |
| **URL** | https://github.com/sgofferj/tak-webview-cesium |
| **License** | GPL-3.0 |
| **Stars** | 8 (new, active — 158 commits, 2026) |
| **TAK Integration** | Direct CoT TLS on port 8089; automated certificate enrollment via port 8446 |
| **Architecture** | Docker container (Python backend + CesiumJS frontend); browser-based |
| **Map Engine** | **CesiumJS** — true 3D globe with elevation/terrain, MIL-STD-2525 symbols via `milsymbol` |

### Features
- Real-time 3D CoT entity visualization on a CesiumJS globe
- Automated TAK Server enrollment (port 8446) **or** manual .p12/.pfx import
- MIL-STD-2525 military symbols
- WMS/XYZ/ArcGIS custom map layers + 3D terrain via Cesium Ion token or self-hosted terrain tiles
- GeoJSON, KML, CZML local overlay support
- Message filtering, hashtag filtering, zoom-to-all
- Multi-language (EN/DE/FI/SV)
- Secure ephemeral cert storage, never written to disk unencrypted
- Configurable via `.env` or Docker Compose environment variables

### Hardware Viability on Dell Inspiron / 940M
**HIGH — Excellent fit.** This runs entirely in a browser on the local machine or a Docker host. The 940M is a capable GPU for WebGL/CesiumJS — CesiumJS is GPU-accelerated but runs fine at 1080p on mid-range mobile GPUs. No dedicated VR hardware required. Run the Docker container on the TAK Server or any Linux host; access via browser from Windows 10.

### Integration with Your Stack
```yaml
# docker-compose.yml snippet
services:
  tak-webview:
    image: ghcr.io/sgofferj/tak-webview-cesium:latest
    ports:
      - "8000:8000"
    environment:
      - TAK_HOST=152.42.155.95
      - TAK_PORT=8089
      - TAK_ENROLL_PORT=8446
      - TAK_CALLSIGN=CesiumViewer
      - CESIUM_ION_TOKEN=<your_free_token>  # optional, for Bing + terrain
    volumes:
      - ./customlayers.json:/app/customlayers.json:ro
      - ./overlays:/app/overlays:ro
```

### Honest Assessment
**This is the closest open-source equivalent to VR-TAK for your hardware.** It provides a live 3D globe view of all CoT entities from your TAK Server, with military symbols, terrain, and overlays — all in a browser. It lacks VR immersion but is fully functional on Windows 10 in Chrome/Edge. The project is small (8 stars) but actively maintained in 2026 and purpose-built for exactly this use case. **Recommend deploying immediately.**

---

## Option 2: CloudTAK (dfpc-coe)

| Field | Value |
|-------|-------|
| **Name** | CloudTAK |
| **URL** | https://github.com/dfpc-coe/CloudTAK |
| **Docs** | https://docs.cloudtak.io |
| **License** | MIT (confirmed by LICENSE file) |
| **Stars** | 119 (11,961 commits — very active) |
| **TAK Integration** | Full TAK Server CoT + WebTAK REST API; connects as a TAK client |
| **Architecture** | Node.js + PostgreSQL/PostGIS + MapLibre GL frontend; full browser-based COP |
| **Map Engine** | **MapLibre GL** (2D vector tiles) — NOT CesiumJS, no native 3D globe |

### Features
- Full-featured browser-based TAK client (COP)
- ETL data ingestion — bring non-TAK sources into TAK Server
- MIL-STD-2525 symbol rendering via `milsymbol`
- Plugin API for extensions
- iOS Capacitor app wrapper
- Active development by Colorado Department of Public Safety COE

### Hardware Viability
**HIGH for 2D use.** MapLibre GL runs well on any modern browser/GPU. However, this is **2D** only — no 3D globe, no terrain elevation rendering. The infra-TAK project (takwerx) deploys CloudTAK as part of its TAK stack with "Cesium 3D-tiles (pmtiles)" mentioned as an optional add-on module (x86-64 only, not ARM64).

### Honest Assessment
CloudTAK is the most mature open-source TAK web client but is fundamentally 2D. Good for COP/SA tasks not requiring elevation/3D terrain. **Use alongside tak-webview-cesium rather than instead of it.** The infra-TAK project wraps CloudTAK + TAK Server + Authentik SSO into a one-click installer — useful if you're also managing the TAK Server stack.

---

## Option 3: TAK-Map (5lyf0x / OpenTAKServer)

| Field | Value |
|-------|-------|
| **Name** | TAK Map |
| **URL** | https://github.com/5lyf0x/tak-map |
| **License** | Apache-2.0 |
| **Stars** | N/A (companion to OpenTAKServer) |
| **TAK Integration** | OpenTAKServer REST API + Data Sync; built for OpenTAKServer deployments |
| **Architecture** | Python (Flask) service + JavaScript frontend; browser-based |
| **Map Engine** | Mixed (includes terrain/elevation/viewshed tooling per README) |

### Features
- Terrain / elevation / viewshed tooling (notable!)
- ADS-B aircraft display with class filtering
- Data Sync object display
- Local data import/export (points, routes, shapes)
- HTTPS via nginx reverse proxy on port 9444

### Hardware Viability
**MEDIUM.** Browser-based; hardware-agnostic for 2D, but viewshed/terrain depends on implementation complexity. Designed for OpenTAKServer (not official TAK Server 5.7), so integration with your DOKS TAK Server would require testing.

### Honest Assessment
Interesting due to explicit terrain/elevation/viewshed features, but built specifically for OpenTAKServer not official TAK Server 5.7. Would need adaptation. Less mature than CloudTAK. Consider only if viewshed analysis is the priority and tak-webview-cesium doesn't provide it.

---

## Option 4: OpenTAKServer (brian7704)

| Field | Value |
|-------|-------|
| **Name** | OpenTAKServer |
| **URL** | https://github.com/brian7704/OpenTAKServer |
| **License** | GNU GPL |
| **Stars** | ~100+ (active) |
| **TAK Integration** | IS a TAK server — full CoT, REST API, WebSocket |
| **Architecture** | Python (Flask) + PostgreSQL; serves TAK clients |

### Honest Assessment
This is an **alternative TAK Server**, not a VR-TAK alternative. Not directly relevant unless migrating away from official TAK Server 5.7, but good to know it exists with TAK Map as its companion web UI.

---

## Option 5: WinTAK 3D Plugin Development (DIY)

| Field | Value |
|-------|-------|
| **SDK Reference** | https://github.com/Hellikandra/WinTAK-plugintemplate (WinTAK >4.9) |
| **SDK Install Guide** | https://github.com/Cale-Torino/WinTAK_SDK_Install |
| **TAK Integration** | Native WinTAK plugin — full access to WinTAK's 3D engine |

### What Exists
- WinTAK SDK (from tak.gov) supports plugin development for WinTAK 4.9+
- WinTAK natively includes a 3D/globe view mode (similar to ATAK 3D)
- Plugin templates exist for building C# plugins that integrate with WinTAK's map canvas
- No public open-source plugin that adds VR-TAK-style 3D to WinTAK has been found

### Hardware Viability
**HIGH** — WinTAK's built-in 3D view runs directly on Windows 10 with the 940M GPU. WinTAK is free (requires tak.gov account) for non-DoD use. The built-in 3D terrain view in WinTAK may already satisfy your requirement without any plugin development.

### Honest Assessment
**Check WinTAK's built-in 3D view first** before looking for alternatives. WinTAK ≥ 4.9 has an integrated 3D map/globe mode. If it meets requirements, this is the simplest path. SDK available if customization needed, but no community 3D plugin found.

---

## Option 6: FreeTAK Server + FreeTAK UI

| Field | Value |
|-------|-------|
| **Name** | FreeTAKServer + UI |
| **URL** | https://github.com/FreeTAKTeam/FreeTakServer |
| **License** | EPL-2.0 |
| **Stars** | 650+ |
| **TAK Integration** | Full CoT server + REST API + NATS-compatible messaging |

### Honest Assessment
FreeTAK is a full TAK server ecosystem (Python-based). The FreeTAK UI (web) exists but is 2D only. No 3D visualization in FreeTAK ecosystem found. Better as a server-side complement than a VR-TAK replacement. Not relevant for 3D SA specifically.

---

## Option 7: trakbridge (emfoursolutions)

| Field | Value |
|-------|-------|
| **Name** | trakbridge |
| **URL** | https://github.com/emfoursolutions/trakbridge |
| **License** | MIT |
| **Stars** | 37 (active — last commit: last week) |
| **TAK Integration** | Bridges tracking devices/services TO TAK Server (sends CoT) |

### Honest Assessment
A data bridge, not a visualization tool. Relevant for feeding GPS trackers, vehicle telemetry, etc. into your TAK Server's CoT stream — which would then appear in tak-webview-cesium. **Complementary, not competitive.**

---

## Option 8: CesiumJS Custom Build (DIY)

| Field | Value |
|-------|-------|
| **Name** | CesiumJS |
| **URL** | https://cesium.com/platform/cesiumjs/ |
| **License** | Apache-2.0 |
| **TAK Integration** | Requires custom code to consume WebSocket/REST from TAK Server |

### What's Required
A custom web app that:
1. Connects to TAK Server's WebSocket/REST API (port 8443/8446 WebTAK API)
2. Parses CoT XML → CesiumJS entities
3. Renders in 3D

**tak-webview-cesium is exactly this, already built.** No need to build from scratch.

### Hardware Viability
CesiumJS WebGL requirement: any GPU with WebGL 1.0+ (940M is fine). The 940M supports WebGL 2.0. For 3D terrain/photogrammetry, the 940M has 2GB VRAM — sufficient for moderate-density scenes. Heavy 3D tiles (high-resolution photogrammetry, dense point clouds) may hit VRAM limits, but basic globe with terrain and CoT entities will run well.

---

## Summary Recommendations

### Immediate Actions (Priority Order)

| Priority | Action | Effort |
|----------|--------|--------|
| **1** | Deploy **tak-webview-cesium** via Docker on TAK Server DOKS pod | 30 min |
| **2** | Test **WinTAK built-in 3D view** on the Dell Inspiron — may already satisfy requirements | 0 dev effort |
| **3** | Deploy **CloudTAK** via infra-TAK for full COP capability (2D + data management) | 1-2 hours |
| **4** | Optionally add **trakbridge** to feed additional tracking sources into CoT stream | as needed |

### Quick Deploy: tak-webview-cesium on TAK Server

```bash
# On a machine with Docker access to the TAK Server DOKS cluster:
docker run -d \
  -p 8000:8000 \
  -e TAK_HOST=152.42.155.95 \
  -e TAK_PORT=8089 \
  -e TAK_ENROLL_PORT=8446 \
  -e TAK_CALLSIGN=3DViewer \
  -e CESIUM_ION_TOKEN=<free_token_from_cesium.com> \
  --name tak-cesium \
  ghcr.io/sgofferj/tak-webview-cesium:latest
```

Then access at `http://<host>:8000` from the Windows 10 machine's browser. Get a free Cesium Ion token at https://ion.cesium.com (free tier includes terrain and Bing maps).

### Hardware Notes for Dell Inspiron 24-7459 / 940M

| Use Case | 940M Viability |
|----------|----------------|
| Browser-based CesiumJS (tak-webview-cesium) | ✅ Excellent |
| Browser-based MapLibre GL (CloudTAK) | ✅ Excellent |
| WinTAK built-in 3D view | ✅ Excellent |
| Heavy 3D photogrammetry tiles | ⚠️ Limited (2GB VRAM) |
| VR headset via PCVR | ❌ Insufficient (needs GTX 970+ for VR) |

---

## Not Found / Gaps

- **No open-source VR-TAK equivalent** (full VR + CoT) exists in the public domain. VR-TAK is unique in requiring VR hardware.
- **No WinTAK community plugin** for 3D mission rehearsal found on GitHub.
- **No ATAK/WinTAK + Unreal Engine** community integrations found.
- **FreeTAK and CloudTAK** are 2D only.
- The **tak-webview-cesium** project fills ~70% of the VR-TAK capability set (3D globe, real-time CoT, MIL symbols, terrain) without the VR immersion component.

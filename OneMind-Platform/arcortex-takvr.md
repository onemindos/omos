# ARCortex / TAKVR — Reference Platform

## What They Are

### ARCortex
AR development firm, 25+ years in immersive technology. Builds custom AR/VR applications for U.S. Navy, Army, and defense contractors. Uses Unity WorldEngine platform for large-scale tactical AR — stereoscopic overlays, digital twin displays, multi-user shared spaces.

### TAKVR
Open-source **Unreal Engine 5.3.2 (C++)** application from TAK Product Center, built by Virtual Heroes / Applied Research Associates (ARA) — the same group that originally built VBS and holds DoD simulation contracts. Government-grade simulation codebase.

Source: https://github.com/TAK-Product-Center/VR

**This is NOT Unity. It is not C#.** It is Unreal Engine 5.3.2 with a plugin-based C++ architecture.

---

## Plugin Architecture

```
Game/
  Plugins/
    VH_COT/           ← TAK Server CoT connection (Commo library — TAK's own C++ comms lib)
    VH_Globe/          ← Geodetic math, WGS84, terrain rendering
    VH_VR/             ← OpenXR VR interaction (controllers, hand tracking, locomotion)
    VH_Replication/    ← Multi-user networking (Unreal NetDriver custom channels)
    VH_Database/       ← MySQL + Spatialite (mission DB, icon cache, DTED terrain)
    VH_FileIO/         ← KML, glTF, 7-Zip, FreeImage, libkml, assimp (3D model import)
    VH_UI/             ← VR UI widgets (UMG in-world panels)
    VH_Keyboard/       ← VR virtual keyboard
    VH_SessionPlayback/ ← Mission replay
    VH_Plugins/        ← Plugin SDK + templates (VTAK_Blank, VTAK_Example)
    VH_Vehicle/        ← Vehicle simulation
    VH_ArtUtilities/   ← Art pipeline tools
    VTAK_Example/      ← Example plugin (template for extensions)
    VTAK_Mensuration/  ← Distance/area measurement tools
    VTAK_Orbit/        ← Orbit/flyby camera
  Source/
    ThirdParty/
      TAKKernel/       ← TAK's native C++ kernel (core CoT, crypto, comms)
      fxGLTF/          ← glTF 2.0 model loader
      libVLC/          ← Video streaming playback
      vtzero/          ← Vector tile parsing (Mapbox Vector Tiles)
```

---

## What TAKVR Already Has (Zero Work Needed)

| Feature | Detail |
|---|---|
| Engine | Unreal 5.3.2 (C++, not Unity) |
| VR Runtime | OpenXR — Quest 3, HoloLens 2, Vive, Index, WMR |
| TAK Connection | Native Commo C++ library — the same comms lib ATAK uses; first-class TAK client |
| Globe / Terrain | Custom WGS84 globe renderer, DTED elevation, Bing/Google hybrid map tiles |
| Multi-user | Unreal built-in replication with custom VH data channels |
| LiDAR | UE5 LidarPointCloud plugin — configured for 900M point budget; Leica BLK scans load natively |
| Vector Tiles | vtzero (Mapbox Vector Tile decoder) — direct compatibility with Martin tile server |
| 3D Import | assimp (FBX/OBJ/DAE), glTF, KML |
| Video Streaming | libVLC — RTSP/HLS playback in VR |
| Session Replay | VH_SessionPlayback |
| Plugin SDK | VTAK_Blank and VTAK_Example templates — build extensions without touching core |
| Targets | Win64, Linux, Android (Quest), Mac |
| License | TAK Product Center (government open source) |

---

## Why We're Extending This, Not Building From Scratch

TAKVR is a production-grade Unreal 5 VR TAK client with native Commo, OpenXR, globe rendering, LiDAR, vector tiles, and a plugin SDK already built. The three plugins we need — NATS intelligence bus, Hermes voice bridge, and 3D Tiles — are focused C++ additions that slot directly into the existing VH_Plugins template system.

This is Path B done for us. We're extending a working VR TAK client, not building one.

---

## OneMind Immersive Layer — `onemindos/VR`

### Fork: `onemindos/VR` (from TAKVR Unreal 5.3.2)

### What We Add — Three VTAK Plugins

#### VTAK_NATS
Connects the VR client to the NATS bus at `wss://ws.onemindos.dev`.

```cpp
// Subscribes via NATS C client (nats.c)
ent.*.state       → ORP-fused entities as 3D markers
det.visual.*      → vision detections overlaid in 3D
intel.*           → threat intel overlay
agents.prompt.hermes.zeus.legacy → voice commands to Hermes
```

#### VTAK_Hermes
Voice-to-NATS bridge: mic input → STT → NATS publish → Hermes processes → TTS response → speaker.
Operators speak to Legacy AI from inside the VR headset.

#### VTAK_3DTiles
Loads Cesium 3D Tiles from WebODM and the geo namespace.
- `tileset.json` URLs sourced from your tile servers
- WebODM photogrammetry meshes render as photorealistic terrain + buildings

### Config Change
`DefaultEngine.ini` → point terrain tile sources at Martin (`martin.onetak.svc:3000`) and TileServer GL.

---

## How It Connects to OneMind

Both OneTAK and onemindos/VR are TAK clients connected to the same TAK Server. They share the Common Operating Picture automatically via CoT. The NATS plugin gives the VR app the intelligence layer that OneTAK's plugins already consume.

```
                    TAK Server (takserver ns)
                   ╱                         ╲
    OneTAK (Vue/MapLibre)              onemindos/VR (Unreal 5/OpenXR)
    ├── nats.ws → NATS bus             ├── VTAK_NATS → NATS bus (same)
    ├── deck.gl → 3D Tiles             ├── VTAK_3DTiles → 3D Tiles (same URLs)
    ├── hermes plugin → Hermes         ├── VTAK_Hermes → Hermes (same subject)
    └── browser (laptop/tablet)        └── VR headset (Quest 3 / HoloLens 2)
```

Same tile servers. Same terrain. Same SA picture. Same intelligence bus. Different viewport.

---

## Use Cases

- **Mission Briefing** — brief operators in a shared 3D space, walk through terrain, plan routes
- **ISR Review** — replay drone footage over the actual 3D terrain model with LiDAR point cloud
- **3D Planning** — draw mission polygons in VR, publish to NATS → all clients see them
- **Immersive SA** — wear a headset and "walk" the SA picture in real scale
- **Training** — simulate operations in a photorealistic digital twin with live OPFOR from ORP
- **AI-Assisted Ops** — speak to Legacy AI in VR; get reasoning overlaid as 3D annotations

---

## Hardware Targets

| Device | Mode |
|---|---|
| Meta Quest 3 / Pro | Standalone VR (Android target) |
| HTC Vive / Valve Index | PC VR (Win64 target) |
| Microsoft HoloLens 2 | AR overlay (Win64 + OpenXR) |
| Standard monitor | Non-VR desktop mode |

---

## Web Alternative — `omos-3d`
For operators without VR hardware: a Three.js or Babylon.js web app that renders the same 3D picture in a browser tab. Less immersive but accessible from any device on the network. Consumes same NATS + TAK data sources as the Unreal client.

---

## Corrected Stack (was Unity, is Unreal)

| Layer | Was (wrong) | Is (correct) |
|---|---|---|
| Engine | Unity (C#) | Unreal Engine 5.3.2 (C++) |
| TAK connection | TAK Unity SDK | Native Commo C++ (same lib as ATAK) |
| NATS client | NATS.Net (C#) | nats.c (C++ NATS client) |
| 3D Tiles | Cesium for Unity | VTAK_3DTiles plugin (vtzero already present) |
| VR | Unity XR Interaction Toolkit | OpenXR via VH_VR plugin |
| Multiplayer | Unity Netcode | Unreal VH_Replication (custom NetDriver) |
| LiDAR | Custom import | UE5 LidarPointCloud plugin (already enabled) |

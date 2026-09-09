# ARCortex / TAKVR — Reference Platform

## What They Are

### ARCortex
AR development firm, 25+ years in immersive technology. Builds custom AR/VR applications for U.S. Navy, Army, and defense contractors. Uses Unity WorldEngine platform for large-scale tactical AR — stereoscopic overlays, digital twin displays, multi-user shared spaces.

### TAKVR
Open-source Unity application from TAK Product Center. Connects to TAK Server as a standard TAK client (CoT consumer). Renders the SA picture in a 3D virtual environment using Unity's render pipeline. Source: https://github.com/TAK-Product-Center/VR

## What They Do
- Render the live SA picture (CoT entities, missions, overlays) in 3D/VR space
- Multiple users share the same 3D operational environment simultaneously
- Overlay tactical symbols on real-world AR view (ARCortex)
- Replay missions, do 3D planning, brief in immersive space
- Connect to TAK Server — they receive the same picture as ATAK/CloudTAK

## Why We're Building Our Own
- TAKVR is already open source — we extend it rather than pay ARCortex
- ARCortex builds custom applications — we want a product we own and iterate on
- Our version adds NATS bus intelligence: ORP fusion, Legacy AI reasoning, live sensor feeds
- Standard TAKVR only gets CoT; our version gets the full intelligence picture

---

## OneMind Immersive Layer — `omos-immersive`

### What We Build
An extended TAKVR application that connects to both TAK Server (CoT) and our NATS bus (intelligence), giving operators a 3D/VR view of the full fused picture.

### Architecture
```
Unity application (fork of TAKVR)
  ├── TAK Client (existing TAKVR) → connects to TAK Server
  │     → receives all CoT: entities, missions, data packages
  │     → renders MIL-STD-2525 symbols in 3D space
  │
  ├── NATS Client (new addition) → connects to K8s NATS bus
  │     → det.visual.* → visual detections overlaid in 3D
  │     → intel.* → intelligence picture
  │     → svc.hermes.query → ask Legacy AI questions in VR
  │     → ent.*.state → ORP-fused entity states
  │
  └── 3D Tiles Consumer → loads WebODM meshes, Leica point clouds
        → tileset.json URLs pushed via TAK data packages
        → renders photorealistic terrain + buildings
```

### Unity Package Stack
- TAK Unity SDK (from TAK Product Center) — CoT parsing, TAK Server connection
- NATS.Net (C# NATS client) — bus connection
- Cesium for Unity — 3D Tiles rendering (terrain, buildings, point clouds)
- Unity XR Interaction Toolkit — VR controllers, hand tracking
- Unity Netcode — multi-user shared space

### Use Cases
- **Mission Briefing** — brief operators in a shared 3D space, walk through terrain, plan routes
- **ISR Review** — replay drone footage over the actual 3D terrain model
- **3D Planning** — draw mission polygons in VR, publish to NATS → all clients see them
- **Immersive SA** — wear a headset and "walk" the SA picture in real scale
- **Training** — simulate operations in a photorealistic digital twin

### Hardware Support
- Meta Quest 3 / Pro (standalone VR)
- HTC Vive / Valve Index (PC VR)
- Microsoft HoloLens 2 (AR overlay on real world — ARCortex use case)
- Standard monitor (non-VR desktop mode)

### Web Alternative — `omos-3d`
For operators without VR hardware: a Three.js or Babylon.js web app that renders the same 3D picture in a browser tab. Less immersive but accessible from any device. Consumes same NATS + TAK data sources.

### Connection to OneMind
```
Legacy AI reasoning → "Threat detected at grid 42" 
  → NATS svc.hermes.broadcast
    → omos-immersive: 3D waypoint appears in VR at that location
    → operator in VR can verbally respond via mic → Legacy AI processes
```
The VR environment becomes a voice-interactive 3D operations center.

# OneMind OS — Recon 04: C2 / Mission / Fusion Overlap
*Generated 2026-07-06 by Legacy (direct). Read-only. ORP dedup + AMOS fit + clean layering.*

## 1. ORP dedup finding

**om-unity, ai-agents/orp, and shieldofsteel/orp are the SAME project.**

- ORP (Open Reality Protocol) = the 45MB Rust fusion binary. The binary is literally named `om-unity`.
- GitHub `shieldofsteel/orp` — the upstream (Apache 2.0, 40+ adapters, DuckDB+Kuzu+RocksDB, ORP-QL, mTLS federation).
- On disk it exists in TWO locations:
  - `nex/om-unity/` — full source tree (crates, frontend, sdk), NOT built
  - `ai-agents/orp/` — has the RUNNING binary at `target/release/om-unity` (port 9090, 5,938 ships, 37+ hrs uptime)

**Recommendation:** ONE canonical ORP. Keep `ai-agents/orp/` (it's the built+running one, matches upstream naming). Verify `nex/om-unity/` isn't a newer fork; if identical, delete it. If it has the frontend/sdk the other lacks, merge then delete. **One source of truth for ORP.**

## 2. Clean C2 layering (no overlap)

```
┌─────────────────────────────────────────────────────────┐
│ LEGACY (AI Commander) — watches all streams, gives orders │  ← Bedrock AI
├─────────────────────────────────────────────────────────┤
│ AMOS (Mission Orchestration) — intent → asset tasking     │  ← NEW layer
│   mission planning, swarm coord, kill web, wargaming      │
├─────────────────────────────────────────────────────────┤
│ om-overwatch (Detection / Kill Chain) — sensor→track→act  │  ← LIVE (air only)
│   generalize to all entity types                          │
├─────────────────────────────────────────────────────────┤
│ ORP / om-unity (Fusion / Data Plane) — all entities graph │  ← LIVE (fusion)
│   39 protocols, entity graph, ORP-QL                      │
├─────────────────────────────────────────────────────────┤
│ NATS bus (the spine) + bridges (MAVLink/CoT/HA/HTTP)      │  ← LIVE
└─────────────────────────────────────────────────────────┘
        ↓ outputs simultaneously ↓
   ATAK (TAK map) · WorldWideView (3D globe) · ClickHouse (history)
```

**No overlap when layered this way:**
- **ORP** = data/fusion (WHAT exists — the entity graph)
- **om-overwatch** = detection + kill chain (WHAT's a threat, act on it)
- **AMOS** = mission orchestration (WHAT should assets DO — intent to tasking)
- **Legacy** = AI commander (WHO decides — natural language, watches everything)

## 3. AMOS fit

**AMOS belongs at the MISSION ORCHESTRATION layer — above om-overwatch, below Legacy.**

- AMOS strengths: operators define *mission intent*, AMOS orchestrates assets across air/ground/maritime/cyber/space. 537 tests, Apache-2.0, Python/Flask, 200+ API routes (open) / 300+ (enterprise). Kill web, SIGINT, wargaming, human-in-the-loop checkpoints, plugin system, TAK/MAVLink/ROS2/DDS/Link-16 integration layer.
- **What it complements:** om-overwatch does the reactive kill chain (detect→track→intercept). AMOS does *proactive multi-asset mission planning* — "patrol this border," "recon this zone with a swarm," "SAR this earthquake." om-overwatch has no mission concept; AMOS adds it.
- **What it could replace:** parts of om-task-nex (task engine) — AMOS has richer mission/task orchestration. Evaluate whether AMOS subsumes om-task.
- **Integration:** AMOS speaks MAVLink/TAK/ROS2 natively → wire its I/O to NATS subjects. AMOS consumes ORP entity graph for situational awareness; publishes mission tasks that om-overwatch/bridges execute.

**Caveat:** AMOS is Python/Flask (heavier than your Go NEX workloads). Run it as a containerized service on DOKS, not as a NEX binary. It's a mission *app*, not a bridge.

## 4. Recommendation — minimal non-redundant C2 stack

1. **Keep ORP** (one copy) as the fusion/data plane. Wire it to NATS `tel.*`/`ent.*`.
2. **Generalize om-overwatch** from air-only to all entity types (highest leverage).
3. **Add AMOS** on DOKS as the mission orchestration layer. Wire its I/O to NATS. Evaluate om-task-nex overlap.
4. **Legacy** stays the AI commander on top, watching all streams via NATS, driving WorldWideView via Agent Bus.
5. **Kill duplicates:** one ORP source, retire pulsar/superseded Constellation clones.

**Boundaries are clean.** No two layers do the same job. The only redundancy to resolve: AMOS mission-tasking vs om-task-nex — verify and consolidate.

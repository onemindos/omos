# `references/` review: what to adopt, watch or drop

**Status:** open. Planned for Saturday 2026-10-17 (see the Todoist task "Review references/ and decide add/remove").
**Goal:** go through the 37 research notes in [`references/`](../../references) and decide, for each, whether OneMind should adopt it, keep watching it, or drop the note. Record the answer in the Decision column. Adoption decisions become their own PRs; this file only records the verdicts.

`references/` stays in the repo for now. Moving it to the local vault was considered and deferred until this review is done.

## How to judge each note

1. **Does it fill a gap?** Check the Services table in `CLAUDE.md` first. Do not add a service that duplicates an existing layer.
2. **Does it fit the bus?** It must publish and subscribe through the locked subject grammar (`tel.`, `det.`, `evt.`, `cmd.`, `intel.`, `ent.`, `svc.`) and run in K8s against `nats.nats.svc.cluster.local:4222`.
3. **Is it alive?** Check last commit, license, and whether the upstream was renamed since the note was written (September 2026).
4. **What does it cost to run?** Memory against the node headroom noted in the 2026-10-15 cluster status, plus one more thing to back up and patch.

**Decision values:** `Adopt` (open an issue/PR to deploy), `Watch` (keep the note), `Have` (already covered, say by what), `Drop` (delete the note).

## Flagged before the review

- **Duplicates:** `infermesh.md` and `ai/infermesh.md`; `defense/peat.md` and `peat-mesh.md`. Keep one of each.
- **Not a reference:** `nats-tactical-repos.md` is a to-do list of repos to fork ("Fork Later (rate-limited)"). Turn it into a decision or delete it.
- **Probably already covered:** `nats/synadia-agents-sdk.md` (Hermes/Legacy is already built on it) and `defense/constellation-overwatch.md` (a fork exists and the overwatch sentinel and vision bridges run on the Mac).
- **Open architectural decision:** `ai/orp.md`. ORP is the candidate fusion layer that would produce `ent.*` state, which nothing publishes today.

## Review table

| Note | What it is | Decision | Notes |
|------|------------|----------|-------|
| `ai-agents-integration.md` | Agent frameworks and the agent communication layer | | |
| `ai/fusion-platform.md` | Multi-sensor fusion on NATS JetStream (IMU, GPS, LiDAR, camera) | | Compare with ORP |
| `ai/infermesh.md` | Geographic load-shifting router for LLM inference over NATS | | Duplicate of `infermesh.md`; compare with LiteLLM already deployed |
| `infermesh.md` | Same project, older note | | Duplicate |
| `ai/orp.md` | Rust fusion engine, 40+ protocols into one signed entity graph | | Fusion-layer decision |
| `apex-sentinel.md` | Civilian drone detection network | | |
| `autonats.md` | Go code generator for NATS service meshes | | |
| `defense/bachmach-pqc-iot-sentinel.md` | Post-quantum hardened IoT telemetry for ESP32/RP2040 | | |
| `defense/constellation-overwatch.md` | Go C4ISR data fabric on NATS | | Fork exists |
| `defense/peat.md` | Mesh protocol for disconnected edge (Defense Unicorns) | | Duplicate of `peat-mesh.md` |
| `peat-mesh.md` | Same project, second note | | Duplicate |
| `dimos.md` | Python OS for generalist robotics | | `rads/` has a DimOS bridge |
| `eclss-digital-twin.md` | Simulated life-support digital twin | | |
| `ha-sqlite-cluster.md` | Distributed SQLite over embedded NATS JetStream | | Compare with nqlite and Marmot |
| `ionode.md` | Flash an ESP32 and it speaks NATS | | |
| `iron-scada.md` | Open SCADA platform (architecture phase) | | |
| `kreuzberg-cloud.md` | Document text/table extraction API | | |
| `logtrace-observability.md` | Logging and tracing on JetStream with OpenTelemetry | | Compare with Loki/Tempo already deployed |
| `marmot.md` | Leaderless SQLite replication by gossip | | |
| `mesh-router.md` | Secure routing between public servers and private nodes | | Compare with Tailscale already deployed |
| `nats-simulation-report.md` | Survey of NATS + simulation repos | | |
| `nats-tactical-repos.md` | Fork-later list for agriculture/tactical NATS repos | | Not a reference |
| `nats/callout.go.md` | Synadia auth-callout service library | | Relevant to per-agent identity |
| `nats/go-atlas.md` | Modular Go toolkit for distributed services | | |
| `nats/jwt-auth-builder.go.md` | Programmatic NATS JWT management | | Relevant to per-agent identity |
| `nats/orbit.go.md` | Synadia Go utility modules | | |
| `nats/resgate.md` | Real-time API gateway between NATS and browsers | | Compare with heynats already deployed |
| `nats/synadia-agents-sdk.md` | Synadia agent protocol SDK | | Hermes already uses it |
| `nqlite.md` | SQLite replicated via NATS JetStream | | |
| `nyxid.md` | Auth/SSO with MCP proxy; agents request credentials | | Compare with Authentik and OpenFGA already deployed |
| `otel-edge.md` | Embedded agent mapping MQTT/NATS IoT telemetry to OpenTelemetry | | |
| `peering-media-pipeline.md` | Distributed media processing on JetStream | | Compare with MediaMTX |
| `qiki-dtmp.md` | Digital twin simulation platform | | |
| `ros-lm-nats.md` | Natural-language ROS2 robot control over NATS | | |
| `sentinelx.md` | ROS2 counter-UAS interceptor stack | | |
| `wireclaw.md` | AI agent running on an ESP32 | | |
| `wlow-orchestrator.md` | DAG workflow orchestrator | | Compare with the mission-engine decision |

## Outputs of the review

- [ ] Verdict recorded for every row above
- [ ] Duplicates merged or deleted
- [ ] `nats-tactical-repos.md` resolved
- [ ] One issue or PR opened per `Adopt`
- [ ] `references/` kept, trimmed, or moved to the vault, and `docs/README.md` updated to match

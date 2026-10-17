# Peat — Mesh Protocol for Disconnected Edge (Defense Unicorns)

**Repo:** https://github.com/defenseunicorns/peat
**Publisher:** Defense Unicorns (makers of UDS/Zarf — DoD software delivery)
**Language:** Rust
**License:** Apache 2.0
**Stars:** 31
**Topics:** crdts, datamesh, datasync, edge-computing, mesh-networks

## What it is

Emergent Capabilities Synthesis Protocol — a mesh protocol that connects heterogeneous systems (phones, servers, sensors, embedded devices, AI models) into a coordinated whole across any transport, even when the network is degraded or denied (DDIL environments).

> "Tactical environments are heterogeneous. TAK operators carry phones. Sensors run on microcontrollers. AI inference runs on edge servers. Robots carry embedded computers."

## Architecture

```
Any device joins:
  Servers │ Phones │ ESP32 │ Raspberry Pi │ AI platforms

Any transport works (simultaneously, with auto-failover):
  QUIC (Iroh) │ BLE mesh │ UDP │ HTTP

Disconnected operation:
  CRDT-based state via Automerge — no central server
  Operates through network partitions (DDIL)

Hierarchy:
  Sensor → Cell → Zone → Command Post
  (state aggregates up without flooding the network)
```

## Crate Ecosystem

| Crate | Purpose |
|-------|---------|
| `peat-protocol` | SDK entry point — cells, hierarchy, sync, security, capabilities |
| `peat-mesh` | P2P topology, QUIC/Iroh transport, Automerge CRDT sync, cert enrollment |
| `peat-btle` | BLE mesh for Android, iOS, Linux, macOS, ESP32 |
| `peat-lite` | Embedded wire protocol for microcontrollers (256KB RAM, no_std) |
| `peat-ffi` | Mobile bindings — Kotlin/Swift via UniFFI + JNI |
| `peat-gateway` | Multi-tenant control plane: enrollment, CDC, OIDC, envelope encryption |

## Why it's relevant to OneMind

| Pattern | OneMind application |
|---------|---------------------|
| DDIL mesh operation | OneMind targets denied/degraded environments — Peat is the reference for how to handle disconnected edge nodes |
| TAK/CoT bridge built in | `peat-schema` includes CoT wire format — connects to the same TAK ecosystem OneMind uses |
| BLE mesh (`peat-btle`) | Short-range device-to-device without WiFi/cell — relevant for barn/field sensor clusters |
| ESP32 embedded (`peat-lite`) | 256KB RAM no_std firmware — same target hardware as OneMind IoT edge |
| Android bindings | Rust → Kotlin/Swift via UniFFI — pattern for OneMind mobile (Flutter FFI) |
| CRDT state sync (Automerge) | Conflict-free state merging when nodes reconnect — alternative/complement to NATS JetStream for offline-first edge |
| Hierarchical aggregation | Squad → Platoon → Company state rollup — mirrors OneMind constellation tracking |
| Defense Unicorns pedigree | Same org behind Zarf (DoD software delivery) — trusted in DoD/IC supply chain |

## Key capabilities for OneMind

1. **Disconnected edge nodes** — `omd` on a barn Pi or robot can use Peat to sync state with the K8s cluster when the NATS leaf reconnects
2. **BLE sensor mesh** — ESP32/Pi nodes discover each other via BLE, form a local mesh, aggregate data before forwarding to NATS
3. **TAK/CoT integration** — Peat already speaks CoT, same as OneMind's NatsBridgePlugin
4. **Multi-transport failover** — QUIC primary, UDP fallback, BLE last-resort — never loses a node even in RF-degraded environments

## Quick start (Rust)

```toml
[dependencies]
peat-protocol = "=0.9.0-rc.1"
# With BLE:
# peat-protocol = { version = "=0.9.0-rc.1", features = ["bluetooth"] }
```

## Companion repos

| Repo | Purpose |
|------|---------|
| https://github.com/defenseunicorns/peat-mesh | P2P topology + QUIC transport |
| https://github.com/defenseunicorns/peat-btle | BLE mesh implementation |
| https://github.com/defenseunicorns/peat-gateway | Multi-tenant control plane |

## Notes

- Defense Unicorns is a DoD-focused open source company — Peat is designed for government/military use cases
- Currently at 0.9.0-rc.1 — not yet stable API, but actively developed
- The local `peat/` clone in omos-core only had compiled Android `.so` binaries — no source
- CRDT sync (Automerge) complements rather than replaces NATS — use NATS for real-time, Peat/Automerge for offline-first sync

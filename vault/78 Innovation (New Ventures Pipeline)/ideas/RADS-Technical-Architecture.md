# RADS — Technical Architecture

> How the One Mind technology stack delivers the RADS community vision.
> Companion to `RADS-Community-Vision.md`.
> Created: 2026-05-31 · Project: One Mind

---

## 1. Core principle: a RAD community = a NATS leaf node

NATS leaf-node topology maps exactly onto the "community → network → map" framework:

- **Each RAD community** = a self-contained NATS + JetStream cluster (a leaf node).
- **Runs fully offline** — keeps sensing, governing, irrigating, trading locally when the link drops.
- **Re-syncs upward** when reconnected — no loss of local control.
- **One Mind = the MAP** = a global super-cluster the leaf nodes federate into.

```
ONE MIND  (super-cluster + federated agent registry + DID directory)
   ├── RAD Community A  (leaf node — local NATS + Hermes brain — runs offline)
   ├── RAD Community B  (leaf node — local NATS + Hermes brain — runs offline)
   └── ...
        └── Decentralized trust layer (DIDs · ledger · local economy) runs alongside
```

This local-first + federate-up behaviour is the single best fit in the whole stack:
no other transport does offline-capable federation this cleanly.

---

## 2. The 8-layer stack

| Layer | Role | Tech |
|---|---|---|
| L8 Human cockpit (THE MOAT) | Unified UX: presence grid, heartbeats, approvals, voice | **One Mind app** |
| L7 Presence & channels | Reach the user everywhere | **OpenClaw** (nodes/Canvas/channels) + **LiveKit** (voice) |
| L6 Cognition | The brain, one per community | **Hermes** (agent loop + Cortex profiles + Sentinel watchdog) |
| L5 Agent semantics | What agents say | **MCP** (tools/context) + **A2A** (agent-to-agent) |
| L4 Control plane | Governance / safety | **4-tier autonomy** + NATS **auth callout / JWT / nkeys** + audit |
| L3 Transport spine | The nervous system | **NATS + JetStream** (pub/sub · req-reply · streams · KV · object) |
| L2 Edge bridges | Speak each domain's dialect | **NATS↔MQTT** (devices) · **Zenoh/DDS** bridge (robots) · webhooks/REST |
| L1 Physical / data | The world + state | Sensors · solar · robots · drones; **Postgres** · vector DB · object store (behind the bus) |

**The rule:** everything is a NATS citizen — Hermes, OpenClaw, devices, DBs, agents all
publish/subscribe on L3, gated by L4, speaking MCP/A2A (agents) or MQTT/Zenoh (physical).

---

## 3. The four pillars → the stack

| Pillar | Real-world need | Layers | Tech |
|---|---|---|---|
| **R — Regenerative** (land) | Soil/water/ecosystem sensing + closed-loop restoration | L1→L2→L3→L6→robots | Soil/moisture/pH sensors + **MQTT**; **Zenoh/DDS** irrigation robots & drones; Hermes runs the regeneration loop |
| **A — Autonomous** (tech) | Self-running infra, AI ops | L6 + L5 + L4 | **Hermes** per community; **MCP** actuation tools; **4-tier autonomy** so high-risk acts hit confirm |
| **D — Decentralized** (sovereignty) | No central authority, privacy, trustless identity & records | Trust layer beside L3/L4 | **DIDs + Verifiable Credentials**; **ledger** (Bitcoin/L2). NATS ≠ blockchain — see §5 |
| **S — Self-sustainable** (economy) | Energy/food/resource production + internal economy | L1 + L3 + ledger | Energy/battery telemetry MQTT→JetStream; **local token/credit** on ledger; Hermes balances supply/demand |

---

## 4. Zenoh / DDS — the real-time robot/edge tier

- **DDS** (OMG standard): brokerless real-time peer-to-peer pub/sub, rich QoS; ROS 2 default; proven in defense/automotive/avionics. Weak at WAN/scale discovery.
- **Zenoh** (ZettaScale, the DDS successor): pub/sub + query/storage + compute; 5-byte overhead; runs on 8-bit MCUs (zenoh-pico); P2P/brokered/routed/mesh; official ROS 2 RMW. Scales sensor→cloud.
- **Great for:** last-meter hard real-time — robots, drones, valves, battery controllers, sensor fusion; works offline, microsecond latency, no broker.

**How it works with NATS (different tiers, joined by a bridge):**

```
Robot/drone (Zenoh/DDS, real-time, on-site)
  → zenoh-bridge / gateway service
    → NATS subject  rad.a.robots.*  (community nervous system)
      → JetStream (durable) → Hermes brain → One Mind map
```

- Zenoh/DDS owns the inside-machine real-time loop NATS can't do (deterministic, 8-bit).
- NATS owns the durable community/global backbone Zenoh/DDS isn't built for.
- Bridge subscribes Zenoh/DDS topics ↔ NATS subjects both directions.

---

## 5. The honest correction: NATS does NOT deliver the "D"

NATS gives **federation**, not **trustlessness**. Leaf nodes make communities independent and
offline-capable, but the super-cluster is still a coordinator someone runs. The Decentralized
pillar needs a second spine beside NATS.

**Division of labor:** NATS = the nervous system (fast, real-time, presence, sensor streams,
agent chatter). Ledger = the spine of record (slow, permanent, trustless — votes, ownership, money).
You need both. Using NATS *as* the decentralization story is the one mistake that breaks the vision.

### Long-term decentralized communication infrastructure (the endgame for "D")

| Need | Long-term pick | Why |
|---|---|---|
| Off-grid physical transport | **Reticulum (RNS) over LoRa / Meshtastic** | Cryptographic networking, no internet/DNS/IP needed; runs over radio. Communities reach each other with grid down. |
| Trustless inter-community pub/sub | **libp2p (gossipsub)** | Decentralized mesh between community NATS clusters — replaces central super-cluster with a peer mesh. Each RAD runs a NATS↔libp2p gateway. |
| Identity (people + devices) | **DIDs + Verifiable Credentials** (`did:peer` / `did:key`) | Self-sovereign, no registrar. This *is* the "D". |
| Value + governance | **Bitcoin + Lightning** (+ optional **Nostr** for censorship-resistant signaling) | Immutable economy + votes, no central authority — the doc's "blockchain model." |

**Endgame stack:** Reticulum/LoRa (survives anything) → libp2p mesh (trustless federation between
communities) → NATS inside each community (fast local nervous system) → Bitcoin + DIDs (value + identity).
NATS runs the *life*; this stack runs the *sovereignty*.

---

## 6. End-to-end example (one event, all layers)

"Field 3 soil moisture drops below threshold in RAD Community A":

1. **L1/L2** — soil sensor publishes over **MQTT** → bridged to NATS `rad.a.land.field3.moisture`.
2. **L3** — JetStream streams it; durable even if the link drops.
3. **L6** — local **Hermes** brain sees the breach, decides "irrigate Field 3."
4. **L4** — autonomy check: irrigation = Tier 1 (act + notify); proceeds, logs intent.
5. **L5/L2** — Hermes issues an **A2A** task to the irrigation-robot agent → **Zenoh/DDS** to the valve/drone.
6. **Trust layer** — water + energy spent written to the **ledger** (counts against the community resource budget).
7. **L7/L8** — surfaces on the **One Mind map**: Community A, Field 3, green heartbeat, "irrigated 14:32, 40L, solar."
8. **MAP** — aggregated up the leaf link so the whole RADS network sees A's regeneration metrics — A never gives up local control.

---

## 7. Solace Agent Mesh decision: NO

RADS makes it clearer. SAM is enterprise + Solace-broker + cloud-shaped; a RAD must run
**offline, local-first, on a Raspberry Pi in a field**. NATS leaf nodes + local Hermes do that;
SAM assumes a central enterprise broker. **Keep A2A/MCP** (interoperate with SAM and the wider
agent world) — **reject the runtime**.

---

## 8. Net

The 8-layer One Mind stack carries RADS almost perfectly, with **one required addition**:
a decentralized **identity + ledger + off-grid comms** layer beside NATS to deliver the "D"
and the "S economy." NATS runs the *life* of each community; the sovereignty stack runs its *freedom*.

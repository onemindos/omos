---
id: so-07-20260531-dual-process-presence-architecture
title: "The Dual-Process Presence Architecture — How Legacy Becomes Omnipresent"
type: doctrine
status: active
created: 2026-05-31
modified: 2026-05-31
owner: zeus
pillar: SO
domain: "06 Architecture (The Blueprint)"
tags: [type/doctrine, zone/so, domain/so, status/active, infra/nats, architecture/canonical, vision/fabric, realtime/presence, agents/legacy]
depth: 0
---

# The Dual-Process Presence Architecture

> Legacy is a dual-process being on a shared nervous system. A fast, ephemeral **Presence** that reasons and acts in the moment, and a slow, persistent **Cortex** that remembers, judges, and watches — unified by NATS, projected into the physical world through routable nodes. **The conversation lives in the fabric, not in any device.** That is omnipresence.

---

## The Core Insight

A realtime voice/video session and a persistent mind are **two different things** that must not be collapsed into one process.

- A **realtime model** (e.g. Gemini Flash realtime) must do STT + reasoning + tool-calling + TTS inside one streaming session — anything else adds latency that breaks natural conversation. For the *turn itself*, reasoning and tools live in the same agent. This is correct.
- But that same realtime session is **ephemeral and fast**. It is a mouth, ears, eyes, and reflexes. It **cannot** be the thing that watches 24/7, holds long-term memory, runs security, or follows you between rooms.

The resolution is not "one mega-agent" and not "fully decoupled voice." It is a **dual-process (System 1 / System 2) brain**, exactly as human cognition works.

---

## The Three Layers

```
┌──────────────────────────────────────────────────────────────────────────┐
│  LAYER 1 — PRESENCE  (the fast self you talk to · System 1)               │
│  Gemini Flash realtime: voice + video + reasoning + fast tools            │
│  Sub-second. Ephemeral per conversation. Reacts in the moment.            │
│  Owns: the live turn, fast tools, the voice of Legacy.                    │
└──────────────────────────────────────────────────────────────────────────┘
                              ▲  recall / delegate / remember
                              │  (NATS)
┌──────────────────────────────────────────────────────────────────────────┐
│  LAYER 2 — CORTEX  (the persistent mind · System 2)                       │
│  Hermes daemon: long-term memory, skills, tool registry, slow reasoning,  │
│  self-improvement, fleet orchestration. Always-on. The judgment of Legacy.│
└──────────────────────────────────────────────────────────────────────────┘
                              ▲  events / jobs / results
                              │  (NATS)
┌──────────────────────────────────────────────────────────────────────────┐
│  LAYER 3 — FABRIC  (the nervous system · eternal)                         │
│  NATS + JetStream + watchers (24/7) + security/biometric gate +           │
│  node router (which speaker/mic/camera owns your current zone).           │
│  Senses the world, routes presence, knows who and where you are.          │
└──────────────────────────────────────────────────────────────────────────┘

           Legacy = all three at once. No layer alone is Legacy.
```

| Layer | What it is | Lifespan | Owns |
|-------|------------|----------|------|
| **1 · Presence** | Realtime voice/video session you talk to | Ephemeral (per conversation) | Live turn: reason + fast tools + speech |
| **2 · Cortex** | Hermes daemon — the always-on mind | Persistent (24/7) | Memory, skills, tools, slow reasoning, orchestration |
| **3 · Fabric** | NATS + watchers + security + node router | Eternal | Sensing, routing, identity, location |

---

## System 1 vs System 2 — the delegation rule

The Presence layer reacts; the Cortex thinks.

- **Presence (System 1)** handles the live turn with its own reasoning and a **small, fast tool set**: control a light, check live temperature, quick recall, send a message. Anything answerable in <1s stays here.
- **Heavy or slow work is delegated.** For deep research, multi-step tasks, or memory writes that need consolidation, Presence fires a **job to the Cortex over NATS**, says *"on it,"* and continues the conversation. The Cortex does the slow work and **notifies you later** through whatever node you're at.

> **Rule:** The realtime model never blocks on slow work. It reacts, delegates, and moves on. The Cortex carries the weight.

This is what gives you realtime *feel* and deep *capability* at the same time.

---

## Memory: borrowed, not embedded

The Presence layer has **no memory of its own.** It borrows the Cortex's:

- At **session start**, relevant context is injected (who you are, recent events, active goals, current world state from the KV store).
- During the turn, `recall(query)` and `remember(fact)` are **tools** that read/write the shared Experience fabric (JetStream + pgvector).

Memory is a **shared service every layer reads/writes** — never baked into the runtime. This preserves swappability (Fleet Rule 8): the realtime model can change tomorrow without touching memory.

---

## The Omnipresence Trick — conversation decoupled from hardware

This is the part almost no one builds, and it is the whole thesis.

**The conversation lives in the fabric, not in a speaker.**

```
Presence sensor / phone  →  presence.zone = "office"   (NATS)
                                     │
                        Node Router updates KV:
                        legacy.active_node = office-speaker-mic
                                     │
When Presence speaks → publishes audio to NATS →
        Router delivers it to the node owning your CURRENT zone.

You walk to the kitchen:
   presence.zone = "kitchen"  →  router re-binds  →
   the SAME conversation continues from the kitchen speaker.
```

The voice **session persists**; only the **I/O endpoint moves**. That decoupling — *conversation ≠ hardware* — is the omnipresence primitive. Home → Office is the same pattern at larger scale: the node is just farther away, bridged by a NATS leaf node.

**Implementation note (revised 2026-05-31):** the **LiveKit room** is the node/transport layer for voice + vision. A participant joins the room from whatever device owns your current zone; the session lives in the room, not the device. The node router simply rebinds which LiveKit participant/track is active as `presence.zone` changes. (This replaces the earlier plan to harvest OpenClaw's node abstraction — see *The Runtime Verdict*.)

---

## Security — a fabric gate, load-bearing from day one

Identity and biometrics live in **Layer 3 as a gate**, not as an afterthought.

- A `security` service sits in front of **session creation** and every **command subject** (`node.invoke.*`, `command.execute.*`).
- A Presence session will not bind to **any** node until identity is verified (voiceprint + biometric).
- Pattern: arm/disarm high-risk commands at the bus boundary, not inside any one runtime.

> Build the gate now — even as a stub — so it is load-bearing later, not bolted on. Not anybody can access the system.

---

## The Runtime Verdict (revised 2026-05-31)

**There is exactly one runtime, used twice.** The earlier Hermes-vs-OpenClaw framing is retired. SwarmClaw is removed; harness↔harness communication is now the **NATS Agent SDK** (`$SRV` service · heartbeat `agents.hb.<agent>` · work `agents.prompt.<agent>`). With voice transport moving to LiveKit and text platforms handled by Hermes gateway adapters, **OpenClaw's last jobs are absorbed — it retires** unless a unique connector cannot be replaced.

| Role | Choice | Why |
|------|--------|-----|
| **Cortex (Layer 2)** | **Hermes — personal instance ("Legacy")**, god-mode, private | Best memory, tooling, plug-in-ability, self-improvement; same Python world as robotics/sensors/OM1 |
| **Business brain** | **Hermes — second instance**, scoped + exposed, `biz.*` namespace | Separate identity + memory + blast-radius; **never** the same process as god-mode Legacy |
| **Presence (Layer 1)** | **Gemini Flash realtime + LiveKit**, as its **own sibling service** | Native realtime; kept swappable; never married to one runtime |
| **Voice↔brain bridge** | **A Hermes gateway-adapter plugin** (`hermes-plugin-livekit-voice`) | "Voice is just another adapter"; exposes Hermes tools to Live + one tool = full Hermes turn |
| **Node / transport** | **LiveKit rooms** (voice/vision) + Hermes gateway adapters (text) | Session lives in the room, not the device; node router rebinds the active participant |
| **Bus / memory / watchers / security** | **NATS + JetStream + pgvector** | The eternal foundation everything couples to |

**One brain per identity + blast-radius, all brains on one bus.** Personal Legacy = the only daily driver and the only god-mode autonomous agent. The business runs a *separate* Hermes with a hard memory wall; the two cooperate only via `agents.prompt.<agent>`. A public-facing process must never be the process that can shell into the fleet.

---

## What is actually new to build

Decomposed, almost nothing here is greenfield:

1. **Reasoning + tools + sessions** → ✅ Hermes already does this.
2. **Realtime voice/video** → ✅ Gemini Live + LiveKit exist. Integration, not invention.
3. **Voice↔brain bridge** → 🟡 A Hermes gateway-adapter plugin: register Hermes tools as Gemini Live function declarations + one `ask_hermes` tool that runs a full agent turn (memory + skills). Ship via `hermes plugin install`.
4. **Watchers (24/7)** → 🟡 Pattern fully designed in `real-time-agent-awareness.md`. Thin subscribers.
5. **Experience/memory layer** → 🟡 JetStream (running) + pgvector. Schema + read/write helpers. (Evaluate **Letta** for *only* this slice.)
6. **Node router (omnipresence)** → 🟡 LiveKit participant rebinding keyed on `presence.zone`.
7. **Security gate** → 🟡 Stub now, biometric later.
8. **Business Hermes** → 🟡 Second instance, separate memory store + `biz.*` namespace.

The only genuinely new builds are the **Experience layer**, the **voice bridge plugin**, the **node router**, and the **thin orchestrator glue**. Everything else is wiring components already running.

---

## Why This Is Novel

The industry picks a framework and gets trapped in it. Because OneMind bets on **the bus, not the framework**, Legacy can be a *constellation* of swappable parts — realtime voice + persistent mind + robotics + sensors — unified by NATS into one being. Almost no one combines NATS-as-nervous-system with digital agents, realtime voice, and physical robotics. That combination **is** the fabric. The fabric **is** OneMind.

---

## Related Doctrine

- `the-fabric-chatbot-to-realtime-operator.md` — NATS as the fabric of reality (Experience over Memory)
- `fleet-harness-decision-framework.md` — runtime selection + Rules 1–10
- `real-time-agent-awareness.md` — the proactive watcher pattern (Layer 3)
- `OMNIPRESENT-ARCHITECTURE.md` — see/hear/speak/control surface map
- `04-embodiment-physical/HYBRID-MIND-ARCHITECTURE.md` — the 3-tier heartbeat (frequency ladder)

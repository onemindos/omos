# OneMind Layer Stack — Cognitive Architecture (L0–L5)

**Date:** 2026-06-07
**Author:** Zeus Delacruz + Legacy
**Status:** Active design — companion to `ONEMIND-NATS-OS-PLAN.md` (the spine) and `08 Research/robotics/world-models-vs-llms-onemind-architecture-2026.md` (the research basis). **Perception addressing RATIFIED: Option A — `tel.sensor.<id>` (Zeus, 2026-06-07, see §6).**
**Purpose:** Define the full cognitive layer stack of OneMind, the mechanism for adding layers that aren't active yet, and the concrete subject contract for the two dark layers (L1 always-on perception, L2 world-model control). Reconciles the world-model paradigm into the existing NATS spine and two-brain model.

---

## 0. The core principle: layers are NODES on the fabric, not a monolith

OneMind is not one program you upgrade. **Every layer is one or more processes that publish/subscribe on the NATS bus.** Adding a capability = attaching a new node to the fabric on a reserved subject. Nothing above it gets rewritten — it simply starts listening to the new subject.

This is the same pattern already proven live: an MQTT device publishes to `devices.plug.a886` and becomes part of the system automatically; agentmemory, Homepage, and Crucix each came online as independent nodes without rebuilding Hermes. **L1 and L2 are the same move — new emitters on the bus, not a new brain.**

> The fabric is the seam. The spine (`ONEMIND-NATS-OS-PLAN.md §5`) is the grammar. A layer becomes "active" the moment a process publishes on its subject and a subscriber reacts.

---

## 1. The full stack (L0–L5) with live status

```
┌─────────────────────────────────────────────────────────────────────┐
│ L5  PRESENCE / INTERFACE                                  🟡 PARTIAL │
│     GPT-realtime (voice, live reasoning) + wearables (glasses HUD,   │
│     watch telemetry, laser projector). LiveKit NOT yet wired.        │
│     → How Zeus and OneMind talk. The "face."                         │
├─────────────────────────────────────────────────────────────────────┤
│ L4  ORCHESTRATOR / META-AGENT (the cortex)                🟢 LIVE    │
│     Hermes / Legacy — memory (agentmemory :3111), tool-calling,      │
│     agent routing, autonomy ladder L0–L3, cron/webhooks.            │
│     WebUI :8787, dashboard :9119. LLM-class. The prefrontal cortex.  │
├─────────────────────────────────────────────────────────────────────┤
│ L3  SPECIALIST AGENTS                                      🟢 LIVE    │
│     10 Hermes profiles: ceo, cfo, coder, eternal, growth, ops,       │
│     sentinel, sovereignty, unity, zenith.                            │
├─────────────────────────────────────────────────────────────────────┤
│ L2  WORLD-MODEL CONTROL LAYER                             🔴 DARK    │
│     Per-robot/domain world models scoring actions in latent space.   │
│     worldforge-style loop. V-JEPA 2-AC / LeWM / Cosmos. Latent MPC.  │
│     Sub-second physical planning. NOT PRESENT YET.                   │
├─────────────────────────────────────────────────────────────────────┤
│ L1  ALWAYS-ON PERCEPTION                                  🔴 DARK    │
│     Local WM encoders running continuously on cameras/sensors.       │
│     Token-free. Emits latents + SURPRISE → wakes the cortex.         │
│     NOT PRESENT YET. ← FIRST layer to light up.                      │
├─────────────────────────────────────────────────────────────────────┤
│ L0  PHYSICAL + DIGITAL FABRIC                             🟢 LIVE    │
│     NATS :4222 (JetStream domain=onemind, MQTT :1883, leaf :7422,    │
│     WS :8080). Robots, drones, vehicles, IoT, SaaS — all nodes.      │
└─────────────────────────────────────────────────────────────────────┘
```

**Honest status: the top half (L3–L4) and the foundation (L0) are live. The entire world-model bottom half (L1–L2) is dark. That is the build target.**

### Why the world model goes at the BOTTOM, not the top
A world model is a **physics simulator / cost oracle**, not a reasoner or orchestrator. It answers *"if I take action A, what does the world look like next, and is that closer to the goal?"* — fast, in latent space, no language. It has **no concept of an agent, task, or tool**, so it cannot be the "meta-agent that manages other agents." That job stays L4 (Legacy). The world model is the cerebellum; Legacy stays the mind. (Full rationale: `08 Research/robotics/world-models-vs-llms-onemind-architecture-2026.md`.)

---

## 2. How to add a layer that isn't active yet (the mechanism)

Every dark layer lights up with the **same three moves**. No big bang, no rewrite of anything above it.

```
1. RESERVE THE SUBJECT   → claim a namespace on the spine (plane.class.id.path)
2. RUN ONE PROCESS       → a single local service that publishes to it
3. SUBSCRIBE THE BRAIN   → point the autonomic governor / Hermes trigger at it
```

**Discipline (from NATS-OS-PLAN risks):**
- Conform to the spine BEFORE running anything. No parallel namespaces. (This doc does that — see §3.)
- Start with ONE node / ONE sensor. Prove the loop with a tiny/mock model, then hot-swap the real world model behind the **same subject contract** — downstream never notices.
- Inner control loops stay native (ROS2/MAVLink at 100–400Hz). Only **derived signals** (latents, surprise, high-level commands) cross to the hub. (Risk #4.)

---

## 3. L1 — Always-On Perception: the subject contract

### 3.1 What L1 is
The **encoder of a world model, running continuously on every sensor**, doing two jobs at once:
- **Compression (sensing):** each frame → a ~192-dim latent fingerprint, at frame rate, locally, forever, zero tokens.
- **Prediction + surprise (watching):** the WM predicts the next latent before it arrives, then compares to reality. `prediction error = surprise`. Low surprise → stay silent. High surprise → emit an event. The model has *learned what's normal for that scene*, so it ignores boring motion and fires on the genuinely unexpected.

This is OneMind's peripheral vision: always running, cheap, and it "pops out" only when something unexpected happens — and the pop-out is what wakes the expensive brain.

### 3.2 Subject contract (CANONICAL — conforms to spine §5)

The spine's 5 planes do the work cleanly. Perception is addressed by the **sensor it watches** (the camera feed), using existing node classes — we do **not** activate the reserved `model` class yet ("add when it hurts, not before").

```
tel.sensor.<id>.latent        TELEMETRY  — continuous latent vector. Ephemeral, high-freq,
                                           local-buffered, downsampled to hub. [THE STREAM]
evt.sensor.<id>.tag           EVENT      — readout detection (person/vehicle/fall). Recorded. [SEMANTIC]
evt.sensor.<id>.surprise      EVENT      — prediction-error spike crossed threshold. Recorded,
                                           fan-out, → ClickHouse for replay. [THE ALARM]
gov.escalate.perception.<id>  GOVERNANCE — wake the cortex (Legacy) with context. [THE INTERRUPT]
$KV.perception.<id>.state     STATE      — current baseline, status (quiet/awake), last-escalation
```

Where `<id>` = the camera/sensor identity, e.g. `barn-cam-1`. Entity view stays intact: `*.sensor.barn-cam-1.>` = everything about that feed across every plane.

**Why these plane assignments are correct:**
- `tel.>` is defined as "continuous high-freq measurement, ephemeral, sampled into evt" — that is *exactly* the latent stream. ✓
- `evt.>` is "something happened, immutable, fan-out, recorded → JetStream + ClickHouse" — surprise spikes and tags belong here; you want them replayable. ✓
- `gov.escalate.>` already exists in the two-brain model as the channel that wakes the cortex. Perception reuses it — **no new brain.** ✓

### 3.3 How L1 plugs into the EXISTING two-brain model (the elegant part)
The NATS-OS-PLAN already defines the **autonomic governor** (`subscriber.py`, the spinal cord) whose job is *"decide what's worth waking the cortex for,"* escalating via `gov.escalate.>`. L1 perception slots in with zero new architecture:

```
WM encoder (L1)  ──tel.sensor.*.latent──►  local JetStream buffer (edge)
       │
       └──evt.sensor.*.surprise──►  AUTONOMIC GOVERNOR (subscriber.py, spinal cord)
                                            │  applies threshold + policy + dedup
                                            └──gov.escalate.perception.*──►  LEGACY (cortex, L4)
                                                                                  │ reads recent latents,
                                                                                  │ reasons, checks autonomy
                                                                                  └─ notify / dispatch / log
```

The world model is just a **new, high-quality emitter** feeding the spinal cord. The governor stays the router. Legacy stays the judgment. This is why the design holds: L1 is additive.

### 3.4 Payload schemas (compact JSON)

`tel.sensor.<id>.latent` — *full stream stays local; downsampled crosses to hub:*
```json
{ "ts": 1717800000.123, "src": "barn-cam-1", "model": "lewm-15m",
  "dim": 192, "z": [0.12, -0.03, "..."], "pred_err": 0.42 }
```

`evt.sensor.<id>.surprise` — *the alarm, recorded:*
```json
{ "ts": 1717800000.456, "src": "barn-cam-1", "pred_err": 3.8, "baseline": 0.4,
  "z_score": 7.2, "window_s": 2.0, "tags": ["motion","person"],
  "clip_ref": "kv:perception/barn-cam-1/clip/1717800000" }
```

`evt.sensor.<id>.tag` — *readout detection:*
```json
{ "ts": 1717800000.456, "src": "barn-cam-1", "tag": "person", "conf": 0.91, "head": "readout-v1" }
```

`gov.escalate.perception.<id>` — *wake the cortex:*
```json
{ "ts": 1717800000.500, "src": "barn-cam-1", "reason": "surprise+person@barn 03:00",
  "severity": "high", "autonomy_floor": "L1", "suggested_action": "notify_zeus+pull_feed",
  "context_ref": "kv:perception/barn-cam-1/clip/1717800000" }
```

### 3.5 Bandwidth discipline
Raw video NEVER crosses the bus. The full latent stream stays in a **local JetStream buffer** on the edge node (NATS-OS-PLAN item 33). Only (a) downsampled latents, (b) surprise events, (c) tags, and (d) escalations cross to the hub. This honors Risk #4 (inner loops native, only derived telemetry crosses) and keeps the hub light.

### 3.6 First node: prove on the hub, deploy to the barn
- **Proof-of-loop (first):** run the L1 process on the **MBP hub** (M2 Max, 32GB — can run a 15M-param LeWM-class encoder easily) with a single camera pointed at the **Office**. Lowest friction; hardware is already here. Proves wiring end-to-end: encode → surprise → `evt.sensor.office-cam-1.surprise` → governor → `gov.escalate.perception.*` → Legacy wakes.
- **First real deployment:** move the proven node to the **Barn** on dedicated edge hardware (Jetson / Mac mini as a leaf node). Highest-value target: unattended, high-value equipment (CNC, tractor, welder), crisp surprise semantic ("person in the barn at 3am"). Isolated from the hub SPOF.

Start tiny: even a small off-the-shelf encoder + frame-difference surprise proves the loop. Swap in **V-JEPA 2-AC / LeWM** behind the same subject contract later — nothing downstream changes.

---

## 4. L2 — World-Model Control: the subject contract (next layer after L1)

L2 is a **separate** process on a **separate** subject, wired to a robot (not a camera). It runs the worldforge-style loop: policy proposes actions → world model scores/rolls them out in latent space → latent MPC picks lowest-cost → robot executes → repeat.

```
qry.robot.<id>.plan           QUERY      — request a scored plan for a goal (req/reply)
cmd.robot.<id>.act            COMMAND    — execute selected action chunk. GOVERNED (hard ACL). 
tel.robot.<id>.state          TELEMETRY  — robot latent state / proprioception
evt.robot.<id>.outcome        EVENT      — action result, recorded for the experience engine
$KV.robot.<id>.goal           STATE      — current goal image / target latent
```

- The **policy proposer** and **world-model scorer** run **on/near the robot** (latency). Only the goal (high-level, from Legacy) and the outcome (telemetry) cross to the hub. Inner control loop stays native (ROS2/MAVLink). (Risk #4.)
- `cmd.robot.*.act` is a governed mutation — passes through the autonomy ladder (L0–L3) exactly like every other `cmd.>`. A world model of *physical* consequence is the natural place to enforce "simulate before you actuate."

L2 is deferred until L1 is proven and a first robot (Go2, per roadmap) is the target.

---

## 5. Reconciliation notes

- **Supersedes** the provisional `om.perception.*` naming used in `08 Research/robotics/world-models-vs-llms-onemind-architecture-2026.md` (Phase 2). The canonical contract is the spine-compliant `tel/evt/gov.*.sensor.*` defined here in §3.2. The research doc's *reasoning* stands; only the subject names are corrected to conform to the spine.
- **No new brain.** L1/L2 reuse the existing two-brain split (autonomic governor → `gov.escalate.>` → Legacy). The world model is an emitter/scorer, not a controller.
- **Perception addressing ratified** (§6): Option A — `tel.sensor.<id>`. No decisions open. Ready to build (Phase 0).

---

## 6. Ratified decision + next steps

**DECISION — perception addressing: ✅ RATIFIED OPTION A (Zeus, 2026-06-07).**
- **✅ Option A (CHOSEN):** address perception by the **sensor** it watches (`tel.sensor.<id>.latent`), using existing classes. Respects "don't activate `model` class until it hurts." Simplest. **This is now the binding contract for L1.**
- **Option B (deferred):** activate the reserved `model` node class (`tel.model.<wm-id>.latent`), treating each world model as a first-class node in its own lifecycle. More future-proof for fleet model-swapping, but opens a HORIZON class early.

> **Standing decision:** L1 perception is addressed by feed (`tel.sensor.<id>`). The `model` HORIZON class stays reserved. **Revisit Option B only when L2 robot world models need independent versioning/swapping across the fleet** — at which point the `model` class splits per spine §"add when it hurts." Until then, every world model is bound to the sensor or robot it serves.

**Next steps (nothing built yet — design only):**
1. ✅ §6 addressing decision ratified (Option A). §3.2 subject contract stands as canonical and binding.
2. Phase 0: `uv add worldforge-ai`, run the mock loop on the hub to internalize policy→score→plan→execute.
3. Phase 1 (L1 proof): single perception process on MBP hub + one Office camera → prove encode→surprise→escalate→Legacy-wake.
4. Phase 2 (L1 real): redeploy to Barn on edge leaf node; swap in V-JEPA 2-AC / LeWM.
5. Phase 3 (L2): first robot (Go2) world-model control loop on the `*.robot.*` contract.

---

*Cross-refs: `01 The Fabric/ONEMIND-NATS-OS-PLAN.md` (spine §5, two-brain §2) · `08 Research/robotics/world-models-vs-llms-onemind-architecture-2026.md` (research basis) · autonomy ladder L0–L3 (SOUL).*

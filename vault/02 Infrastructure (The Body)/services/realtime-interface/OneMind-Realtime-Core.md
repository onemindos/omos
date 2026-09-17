---
title: OneMind Realtime Core — Omnipresent AI Interface Architecture
domain: 02 Infrastructure (The Body)
status: design
created: 2026-05-24
tags: [realtime, multimodal, nats, local-inference, qwen, omni, drones, robotics, edge, voice, vision, architecture]
related: [[NATS]], [[HA-NATS-Bridge]], [[NATS-Intent-Mesh]]
---

# OneMind Realtime Core — Omnipresent AI Interface Architecture

## Vision Statement

A single local multimodal AI model that sees, hears, speaks, and acts — across every interface (glasses, phone, desktop, car, field kit) — without switching models. NATS is the nervous system to the world. The model is the brain. Voice/vision are local I/O. Commands flow through NATS to every device in the fleet.

**Core equation:** `Local Multimodal Model + NATS Nervous System = Omnipresent AI`

---

## The Two-Layer Architecture

Most systems conflate these. OneMind separates them:

| Layer | Role | Latency Budget | Model |
|---|---|---|---|
| **Interface Layer** (always-on) | See, hear, speak, route, react, maintain context | **<100ms** | Local multimodal (Qwen2.5-Omni-7B class) |
| **Reasoning Layer** (on-demand) | Plan, code, analyze, create, deep-think | 1–30s acceptable | Claude, GPT-5, o4, DeepSeek R1 |

**The interface layer never changes. The reasoning layer gets routed to.**

This mirrors the human brain:
- Reflexes and perception (brainstem/thalamus) = millisecond speed
- Deep thinking (prefrontal cortex) = slow and deliberate
- You don't swap your eyes out when you want to think harder

---

## Master Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    INTERFACE LAYER (Local, Always-On)                     │
│                                                                           │
│   Qwen2.5-Omni-7B (or successor)                                        │
│   ┌──────────┬──────────┬──────────┬──────────┬──────────┐              │
│   │  Vision  │  Audio   │  Speech  │  Text    │ Context  │              │
│   │  (eyes)  │  (ears)  │  (mouth) │  (chat)  │ (memory) │              │
│   └────┬─────┴────┬─────┴────┬─────┴────┬─────┴────┬─────┘              │
│        │           │          │          │          │                     │
│        └───────────┴──────────┴──────────┴──────────┘                     │
│                              │                                            │
│                    LOCAL I/O (no NATS):                                   │
│                    • Mic → STT → model                                   │
│                    • Model → TTS → speaker                               │
│                    • Camera → model vision                               │
│                              │                                            │
│                         NATS (4222)                                       │
│                    (commands & context only)                              │
│                              │                                            │
└──────────────────────────────┼────────────────────────────────────────────┘
                               │
        ┌──────────────────────┼──────────────────────────┐
        │                      │                          │
   cmd.agent.reason       cmd.ha.*              cmd.drone.<id>.*
        │                      │                          │
   ┌────┴────┐            ha-bridge              drone runtime
   │ Claude  │            (~1ms)                 (local on device)
   │ GPT-5   │
   │ o4-mini │
   │ DeepSeek│
   └─────────┘
   (1-30s, async)
```

---

## How Voice ACTUALLY Works With the Model

**NATS is NOT between your mouth and the model.** NATS carries commands to the world, not your voice to the AI.

```
YOUR VOICE
    │
    │  raw audio (microphone → system audio buffer)
    ▼
LOCAL PROCESS (on same device as model)
    │
    ├── STT (Whisper) converts audio → text     (~200ms)
    ├── Model processes text (+ vision if active) (~300ms)
    ├── Model decides what to do
    │       │
    │       ├── Response to YOU → TTS → speaker  (local, no NATS)
    │       └── Command to WORLD → nats_publish  (to HA, drone, etc.)
    │
    ▼
NATS only carries COMMANDS and CONTEXT — not voice audio
```

**Rules:**
- **Your mouth → model → your ears** = local process, same device, no network
- **Model → turn on lights** = NATS publish to ha-bridge
- **Drone telemetry → model's awareness** = NATS subscribe
- **NATS = nervous system to the world**, not the pathway between voice and AI

---

## Same Model, Every Interface

| Interface | Input | Output | Model changes? |
|---|---|---|---|
| Glasses (Meta, custom) | Camera + mic | Spatial audio + HUD overlay | ❌ Same |
| Phone | Text / voice | Text / voice | ❌ Same |
| Desktop | Text / screen share | Text / voice / actions | ❌ Same |
| Car | Voice + cabin camera | Voice + heads-up | ❌ Same |
| Drone/robot | Telemetry + camera | Control commands | ❌ Same |

The model is **modality-native** — it doesn't "switch" to handle vision or audio. It processes all modalities simultaneously in the same forward pass.

---

## Edge Device in the Field — Topology

### Scenario A: In the field WITH drones, FAR from home MBP

```
┌─────────────────────────── FIELD KIT ───────────────────────────┐
│                                                                   │
│  Edge Device (Jetson Orin / Mac Mini / phone)                    │
│  ├── Qwen2.5-Omni-7B (same model, same interface)               │
│  ├── NATS Server (local instance)                                │
│  └── Direct WiFi mesh to drones                                  │
│                                                                   │
│  Drone 1 ←──── NATS ────→ Edge NATS ←──── YOU (voice/glasses)  │
│  Drone 2 ←──── NATS ────→      │                                │
│  Drone 3 ←──── NATS ────→      │                                │
│                                  │                                │
└──────────────────────────────────┼────────────────────────────────┘
                                   │
                            LEAF NODE (cellular/starlink)
                                   │
                              ┌────┴────┐
                              │ MBP NATS │  (home, cloud reasoning)
                              │ (4222)   │
                              └─────────┘
```

**Key:** Edge device runs its OWN NATS server. Drones connect directly (local WiFi, ~1ms). Leaf node syncs back to MBP when connectivity exists. **Works fully offline** — if cell drops, everything still works locally. Leaf reconnects automatically when signal returns.

### Scenario B: Remote, drones at HOME (near MBP)

```
YOU (phone/laptop, anywhere)
    │
    │  Tailscale / cellular
    ▼
MBP NATS (home)
    │
    │  cmd.drone.01.land
    ▼
Drone (home WiFi → MBP NATS directly)
```

Use MBP as command center. Phone sends commands via Tailscale. Latency = Tailscale RTT (~50-100ms over cell).

### When to Use Leaf Node

| Topology | Use |
|---|---|
| Edge has its own drones nearby | Edge runs own NATS. Leaf node syncs to MBP for logging/reasoning offload |
| You're remote, drones at home | Connect to MBP directly via Tailscale. No leaf needed |
| Both edge AND home need to coordinate | Leaf node bridges the two NATS clusters |

**Leaf node = bridge between two NATS clusters.** For when you need two autonomous zones that share context.

---

## Model Selection: The Local Multimodal Stack

### Available Now (M2 Max 32GB compatible)

| Model | Params | Modalities | Fits? | Tokens/sec |
|---|---|---|---|---|
| **Qwen2.5-Omni-7B** | 7B | Text + Vision + Audio in/out | ✅ ~6-8GB Q4 | ~10-15 t/s (Transformers) |
| Qwen2.5-VL-7B | 7B | Text + Vision | ✅ | ~30 t/s (MLX) |
| MiniCPM-o 2.6 | 8B | Text + Vision + Audio | ✅ | ~35 t/s |
| Qwen2.5-7B-Instruct | 7B | Text only | ✅ | 50+ t/s (llama.cpp) |

### Runtime Support (as of 2026-05-24)

| Runtime | Qwen2.5-Omni support |
|---|---|
| HuggingFace Transformers | ✅ Full (text + audio in + audio out) |
| vLLM | ✅ Text + vision (audio partial) |
| llama.cpp / GGUF | ❌ Not yet ported |
| MLX | ❌ Not yet ported |

**Current best approach:** Pipeline (Whisper.cpp + Qwen2.5-7B + Piper TTS) for speed. Swap to native Qwen-Omni once llama.cpp/MLX support lands.

### Near-Term Models (3-6 months)

| Model class | Expected | Why it matters |
|---|---|---|
| Qwen3-Omni | 7B-14B | Better reasoning + same multimodal |
| Llama 4 multimodal | 8B-70B | Meta pushing glasses/wearable use case |
| Apple Foundation Model | On-device | Native to glasses/phone hardware |

### Hardware Scaling Path

| Device | Role | Model size |
|---|---|---|
| **MBP M2 Max 32GB** (current) | Primary interface runtime | 7B quantized (Q4_K_M) |
| **Mac Studio M4 Ultra 192GB** (upgrade) | Full 70B local + multiple 7B | Dream tier |
| **Edge device on glasses** | Whisper-tiny for wake word, stream audio to MBP | — |
| **Drone onboard** | No LLM — pure control loop (PID, obstacle avoidance at 100Hz+) | — |

---

## The Omnipresent Pattern

What makes this system omnipresent is that the model is **always subscribing to everything** — it already knows what's happening before you ask.

```
┌────────────────────────────────────────────────────────────┐
│              ALWAYS-ON LOCAL MODEL                          │
│                                                            │
│  SUBSCRIBES TO (always aware):                             │
│  ├── events.ha.*          → home state changes             │
│  ├── perception.vision.*  → camera frames (glasses/drone)  │
│  ├── perception.audio.*   → ambient sound classification   │
│  ├── context.location.*   → where Zeus is                  │
│  ├── cmd.drone.*.telemetry → fleet positions              │
│  ├── events.comms.*       → incoming messages              │
│  └── KV watch: zeus-context → focus mode, status           │
│                                                            │
│  PUBLISHES (acts on):                                      │
│  ├── cmd.ha.service.call  → control home                   │
│  ├── cmd.drone.*.override → intervene with fleet           │
│  ├── cmd.agent.reason     → offload heavy thinking         │
│  ├── cmd.tts.speak        → voice output (if remote)       │
│  ├── events.zeus.*        → context updates for other subs │
│  └── context.location.*   → location awareness             │
│                                                            │
│  LOCAL I/O (no NATS, direct):                              │
│  ├── Microphone → STT → model input                        │
│  ├── Model output → TTS → speaker/earpiece                 │
│  └── Camera → model vision input                           │
└────────────────────────────────────────────────────────────┘
```

---

## Real-Time Examples

### Example 1: Walking Through the House

```
TIME    EVENT                           NATS SUBJECT                        MODEL ACTION
─────────────────────────────────────────────────────────────────────────────────────────
0.0s    Zeus walks into kitchen          events.ha.motion.kitchen            Context: "Zeus in kitchen"
0.1s    Model updates internal state     (local, no NATS)                    Remembers: kitchen context
0.3s    Zeus: "What's the temperature?"  (mic → STT → model, local)         —
0.5s    Model checks context cache       (already has sensor data via sub)   Knows: 75.2°F
0.6s    Model: "75 degrees"              (TTS → speaker, local)             —
1.0s    Zeus: "Turn on the coffee maker" (mic → STT → model)                —
1.2s    Model publishes                  cmd.ha.service.call                 {"domain":"switch","service":"turn_on","entity_id":"switch.plug_mini_ef5a"}
1.3s    ha-bridge executes               (WebSocket to HA)                   Done in ~5ms
1.4s    Model: "Coffee maker's on"       (TTS → speaker)                    —
```

**Total: voice command → device action → confirmation = ~400ms**

### Example 2: Drone Operation in Field

```
TIME    EVENT                           WHERE                               WHAT HAPPENS
─────────────────────────────────────────────────────────────────────────────────────────
0.0s    3 drones flying patrol           Edge NATS (field kit)              Drones publish telemetry every 100ms
        │                                                                    Model subscribes, builds spatial map
0.0s    Model (continuous): tracking     Local on edge device               "Drone 1: sector A, Drone 2: sector B..."
        all positions, battery,
        camera feeds

3.2s    Zeus (glasses): "Drone 2,        (mic → local STT → model)         —
        check that building"
3.4s    Model publishes                  cmd.drone.02.waypoint              {"lat":41.02,"lon":-73.63,"alt":30,"action":"orbit"}
3.5s    Drone 2 receives, adjusts        (NATS → drone controller)         Drone changes course
3.6s    Drone 2 camera feed updates      perception.drone.02.frame          Model sees new view
3.8s    Model: "Drone 2 approaching,     (TTS → glasses speaker)           Zeus hears update
        I see a vehicle in the lot"

5.0s    Zeus: "All drones return"        (mic → model)                      —
5.1s    Model publishes                  cmd.fleet.rtl                      Fan-out: ALL drones get it simultaneously
5.1s    Drones 1,2,3 receive             (NATS fan-out, <1ms)              All start returning
5.2s    Model: "All three returning"     (TTS → glasses)                   —
```

**Key: model is ALWAYS listening to drone telemetry (NATS subscribe). It doesn't query on demand — it already knows.**

### Example 3: Glasses on a Walk

```
TIME    EVENT                           SYSTEM                              ACTION
─────────────────────────────────────────────────────────────────────────────────────────
0.0s    Zeus leaves house                HA motion sensor: no motion         events.ha.motion.hallway → off
0.1s    Model: context = "Zeus leaving"  Updates internal state              —
0.2s    Model publishes                  cmd.ha.service.call                 Turns off all lights (auto)
0.3s    Model: "Lights off, have         TTS → glasses earpiece             —
        a good walk"

30.0s   Zeus walking, glasses camera on  perception.vision.frame             Model processes at 2-5 FPS
        Model: ambient awareness         (local inference)                   Tracks environment, people, objects

45.0s   Zeus passes neighbor             Model sees face (local)            —
45.1s   Model: "That's Mike from         TTS → glasses (whisper volume)     Zeus gets social context
        3 doors down, you talked
        about the fence last week"

60.0s   Zeus: "Remind me to call         (mic → model)                      —
        the electrician tomorrow"
60.2s   Model publishes                  cmd.agent.reason                    Routes to reasoning layer (Claude)
                                                                             for proper task creation
60.5s   Claude creates Todoist task      (via NATS → saas worker)           —
60.8s   Model: "Done, tomorrow 9am"      TTS → glasses                      —
```

### Example 4: Realtime Drone Intervention (Millisecond Override)

```
TIME    EVENT                           LATENCY                             PATH
─────────────────────────────────────────────────────────────────────────────────────────
0.0ms   Drone 1 approaches obstacle      —                                  Drone's own PID loop handles (100Hz)
0.0ms   Drone publishes alert            cmd.drone.01.alert.proximity       "Object 2m ahead"
0.4ms   Edge NATS delivers to model      NATS bus overhead                  —
0.5ms   Model receives, processes        Local inference                    Decides: override needed
15ms    Model: "Drone 1, hold position"  TTS (optional, for Zeus)           —
15ms    Model publishes                  cmd.drone.01.override              {"action":"hover","reason":"obstacle"}
15.4ms  Drone receives override          NATS delivery                      Drone holds
16ms    Total intervention time          —                                  ~16ms from alert to override

NOTE: Drone's OWN control loop already handles obstacle avoidance autonomously.
      Model override is for STRATEGIC intervention, not reflexive safety.
      Drone doesn't WAIT for model — it runs its own loop at 100Hz+.
      Model CAN intervene when Zeus commands or when strategic context requires it.
```

---

## Intervention Speed Summary

| Action | Path | Latency |
|---|---|---|
| Voice → HA device control | Mic → STT → Model → NATS → ha-bridge → HA | ~400ms |
| Voice → drone command (local) | Mic → STT → Model → NATS → drone | ~200ms |
| Model auto-override (no voice) | NATS event → Model → NATS publish | ~15ms |
| Fleet fan-out (all drones) | Model → NATS publish → all subscribers | ~1ms delivery |
| Cross-cluster (leaf node) | Model → NATS → leaf → remote NATS → device | ~50-150ms |

---

## LLM Speed Comparison (Why Local)

| Provider | Model | TTFT | Tokens/sec | Use |
|---|---|---|---|---|
| **Local (llama.cpp)** | Qwen2.5-7B Q4 | **~50ms** | ~50 t/s | Interface layer ✅ |
| **Groq** | Llama 3.1 8B | ~80ms | ~800 t/s | Fast cloud fallback |
| **Groq** | Llama 3.3 70B | ~150ms | ~300 t/s | Heavy cloud reasoning |
| Bedrock | Claude Sonnet 4.6 | ~800ms+ | ~80 t/s | Deep reasoning only |
| Anthropic | Claude Haiku 3.5 | ~300ms | ~150 t/s | Medium reasoning |

**Rule:** Interface layer = local. Reasoning layer = routed (cloud or local 70B on upgrade hardware).

---

## Build Phases

### Phase 1 — Prove it now (this week)

1. Deploy Qwen2.5-7B on MBP via llama.cpp (fast text inference, 50+ t/s)
2. Add Whisper.cpp for voice input (local STT, ~200ms)
3. Add Piper/edge-tts for voice output (local TTS, ~100ms)
4. Wire as NATS service (`cmd.interface.*`)
5. Prove the loop: voice → model → NATS → HA → confirmation → voice
6. **Target:** voice command → device action → spoken confirmation < 500ms

### Phase 2 — Multi-interface (weeks 2-4)

- Same model serves text (Telegram/WebUI), voice (mic), vision (camera/screen)
- NATS subject routing decides which modality the response goes to
- Context injection via KV watch (location, focus mode, what's on screen)
- Continuous context updates while speaking (model subscribes to events.ha.*, etc.)

### Phase 3 — Wearables + fleet (month 2+)

- Glasses stream video frames → local model processes at 2-5 FPS
- Drone telemetry → NATS → model builds spatial awareness
- Phone as fallback interface when glasses aren't on
- Edge field kit: own NATS server + own model + leaf node to home

### Phase 4 — Full sovereignty (month 3+)

- Mac Studio M4 Ultra 192GB: run 70B model locally (interface + reasoning on same box)
- Zero cloud dependency for anything except optional offload
- Multiple simultaneous model instances (one per active interface)

---

## Key Design Decisions

1. **Voice/vision ↔ model = LOCAL. No NATS in that path.** NATS carries commands to the world, not your voice to the AI.

2. **Edge device in field = own NATS server.** Leaf node bridges to home when connectivity exists. Drones connect to edge NATS directly. Works offline.

3. **Qwen2.5-Omni-7B = target model.** Runs on M2 Max today via Transformers. Pipeline approach (Whisper + Qwen-7B + TTS) is faster until native MLX/llama.cpp support lands.

4. **Omnipresence = model subscribes to everything, always has context, never queries on demand.** NATS push model means the AI already knows what's happening before you ask.

5. **Drones run their OWN control loops.** Model provides strategic override, not reflexive safety. Drone doesn't wait for model — it's autonomous. Model can intervene.

6. **Same model, every interface.** Glasses, phone, desktop, car, field kit — no switching. The interface layer is one model that handles all modalities.

7. **Reasoning is routed, not inline.** Heavy thinking (code, planning, analysis) gets published to `cmd.agent.reason` and handled by Claude/GPT-5/o4 asynchronously. Interface model reports back when reasoning is done.

---

## NATS Subject Namespace (Realtime Core)

| Namespace | Direction | Purpose |
|---|---|---|
| `perception.vision.*` | → model | Camera frames (glasses, drone, security) |
| `perception.audio.*` | → model | Ambient sound, classified audio events |
| `context.location.*` | ↔ | Zeus's physical location updates |
| `context.focus.*` | ↔ | Focus mode, activity state |
| `cmd.interface.*` | → model | External text/command input to model |
| `cmd.interface.response` | model → | Model's text response output |
| `cmd.ha.*` | model → | Home automation commands |
| `cmd.drone.<id>.*` | model → | Drone fleet commands |
| `cmd.fleet.*` | model → | Fan-out to all fleet devices |
| `cmd.agent.reason` | model → | Offload to reasoning layer |
| `cmd.tts.speak` | model → | Remote TTS output (non-local speaker) |
| `events.ha.*` | → model | Home state changes (subscribe) |
| `events.comms.*` | → model | Incoming messages |
| `events.security.*` | → model | Motion, alerts, anomalies |

---

## Related Notes

- [[NATS]] — central bus architecture and server config
- [[NATS-Intent-Mesh]] — full subject namespace and domain workers
- [[HA-NATS-Bridge]] — Home Assistant integration via NATS
- [[Hermes-Overview]] — Legacy/Hermes agent layer (current reasoning interface)

## Related Skills (Hermes)

- `nats-server` — NATS server ops, streams, KV, ha-bridge, subscriber architecture
- `home-assistant` — HA usage via NATS bridge
- `llama-cpp` — local model inference setup
- `serving-llms-vllm` — vLLM for heavier models

---
title: OneMind Realtime Core — Omnipresent AI Architecture
domain: 02 Infrastructure (The Body)
status: design
created: 2026-05-24
tags: [onemind-os, realtime, multimodal, nats, qwen-omni, edge, drones, robotics, glasses, voice, omnipresent]
related: [[NATS]], [[NATS-Intent-Mesh]], [[HA-NATS-Bridge]]
---

# OneMind Realtime Core — Omnipresent AI Architecture

## Vision

A **single multimodal AI model** that runs locally, sees, hears, speaks, and acts — across every interface (glasses, phone, desktop, car, field kit) without ever switching models. The same interface whether texting, speaking in real-time, wearing glasses, or commanding a drone fleet.

**Core principle:** The interface layer is fast, local, always-on. Heavy reasoning is offloaded asynchronously. NATS is the nervous system that connects the AI to the world — but NOT between your voice and the model.

---

## Two-Layer Architecture

| Layer | Role | Latency Budget | Model |
|---|---|---|---|
| **Interface Layer** (always-on) | See, hear, speak, route, react, maintain context | **<100ms** | Local multimodal (Qwen2.5-Omni-7B class) |
| **Reasoning Layer** (on-demand) | Plan, code, analyze, create, deep-think | 1–30s acceptable | Claude, GPT-5, o4, DeepSeek R1 |

**The interface layer never changes. The reasoning layer gets routed to.**

This is how a human brain works — reflexes and perception (brainstem/thalamus) run at millisecond speed. Deep thinking (prefrontal cortex) is slow and deliberate. You don't swap your eyes out when you want to think harder.

---

## How Voice Actually Works (Clearing the Confusion)

**NATS is NOT between your mouth and the model.** That would be insane overhead for no reason.

```
YOUR VOICE
    │
    │  raw audio (microphone → system audio buffer)
    ▼
LOCAL PROCESS (on same device as model)
    │
    ├── STT (Whisper) converts audio → text          (~200ms)
    ├── Model processes text (+ vision if active)    (~300ms)
    ├── Model decides what to do
    │       │
    │       ├── Response to YOU → TTS → speaker      (local, no NATS)
    │       └── Command to WORLD → nats_publish      (to HA, drone, etc.)
    │
    ▼
NATS only carries COMMANDS and CONTEXT — not voice audio
```

### The Rule

- **Your mouth → model → your ears** = local process, same device, no network
- **Model → turn on lights** = NATS publish to ha-bridge
- **Drone telemetry → model's awareness** = NATS subscribe

**NATS is the nervous system to the world, not the pathway between your voice and the AI.**

---

## Master Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                    INTERFACE LAYER (Local, Always-On)             │
│                                                                   │
│   Qwen2.5-Omni-7B (or successor)                                │
│   ┌──────────┬──────────┬──────────┬──────────┬──────────┐      │
│   │  Vision  │  Audio   │  Speech  │  Text    │ Context  │      │
│   │  (eyes)  │  (ears)  │  (mouth) │  (chat)  │ (memory) │      │
│   └────┬─────┴────┬─────┴────┬─────┴────┬─────┴────┬─────┘      │
│        │           │          │          │          │             │
│        └───────────┴──────────┴──────────┴──────────┘             │
│                              │                                    │
│                         NATS (4222)                               │
│                              │                                    │
└──────────────────────────────┼────────────────────────────────────┘
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

## Same Model, Every Interface

| Interface | Input | Output | Model changes? |
|---|---|---|---|
| Glasses (Meta, custom) | Camera + mic | Spatial audio + HUD overlay | ❌ Same |
| Phone | Text / voice | Text / voice | ❌ Same |
| Desktop | Text / screen share | Text / voice / actions | ❌ Same |
| Car | Voice + cabin camera | Voice + heads-up | ❌ Same |
| Drone/robot | Telemetry + camera | Control commands | ❌ Same |

The model is **modality-native** — it doesn't "switch" to handle vision or audio. It processes all modalities simultaneously in the same forward pass. That's what Qwen2.5-Omni does.

---

## Edge Device in the Field — Topology

### Scenario A: You're in the field WITH drones, FAR from home MBP

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

**Key: You do NOT need the MBP for field operations.**

- Edge device runs its OWN NATS server
- Drones connect directly to edge NATS (local WiFi, ~1ms)
- Your voice → edge model → `nats_publish("cmd.drone.01.override")` → drone gets it in <5ms
- **Leaf node is optional** — syncs context/logs back to MBP when connectivity exists
- If cell signal drops, **everything still works locally**. Leaf node reconnects automatically when signal returns

### Scenario B: You're remote but drones are at HOME (near MBP)

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

Your phone sends commands to MBP's NATS via Tailscale. MBP relays to drones. Latency = Tailscale RTT (~50-100ms over cell).

### When to Use Leaf Node

| Topology | Use |
|---|---|
| Edge has its own drones nearby | Edge runs own NATS. Leaf node syncs to MBP for logging/reasoning offload |
| You're remote, drones at home | Connect to MBP directly via Tailscale. No leaf needed |
| Both edge AND home need to coordinate | Leaf node bridges the two NATS clusters |

**Leaf node = bridge between two NATS clusters.** It's for when you need TWO autonomous zones that share context.

---

## Drone/Robotics Intervention Model

```
Drone runtime (onboard):
  - Runs its own control loop (PID, obstacle avoidance) at 100Hz+
  - Publishes telemetry to NATS: cmd.drone.01.telemetry (10Hz)

Interface model (MBP or edge device):
  - Subscribes to drone telemetry
  - Zeus says "pull back" → model generates cmd.drone.01.override
  - RTT: voice→model→NATS→drone = ~50-100ms (local) or ~150ms (remote via leaf node)

Key: drone doesn't WAIT for model. It runs autonomously.
     Model can INTERVENE when Zeus commands. Fire-and-forget override.
```

Same as how flight controllers work — **autonomous runtime with human-in-the-loop override via NATS publish.**

---

## Model Selection: Local Multimodal Stack

### Available Now (M2 Max 32GB)

| Model | Params | Modalities | Fits? | Tokens/sec |
|---|---|---|---|---|
| **Qwen2.5-Omni-7B** | 7B | Text + Vision + Audio in/out | ✅ ~6-8GB Q4 | ~10-15 t/s (Transformers) |
| Qwen2.5-VL-7B | 7B | Text + Vision | ✅ | ~30 t/s (MLX) |
| MiniCPM-o 2.6 | 8B | Text + Vision + Audio | ✅ | ~35 t/s |
| Qwen2.5-7B-Instruct | 7B | Text only | ✅ | ~50+ t/s (llama.cpp) |

### Runtime Support (as of 2026-05-24)

| Runtime | Qwen2.5-Omni support |
|---|---|
| HuggingFace Transformers | ✅ Full (text + audio in + audio out) |
| vLLM | ✅ Text + vision (audio partial) |
| llama.cpp / GGUF | ❌ Not yet ported |
| MLX | ❌ Not yet ported |

### Pipeline Approach (Fast Today)

Until native Omni support lands in llama.cpp/MLX:

```
Whisper.cpp (STT) → Qwen2.5-7B (reasoning, llama.cpp) → Piper/edge-tts (TTS)
     ~200ms                    ~300ms                         ~100ms
                                                        Total: ~600ms
```

### Near-Term (3-6 months)

| Model class | Expected | Why it matters |
|---|---|---|
| Qwen3-Omni (rumored) | 7B-14B | Better reasoning + same multimodal |
| Llama 4 multimodal | 8B-70B | Meta pushing glasses/wearable use case |
| Apple Foundation Model | On-device | Native to glasses/phone hardware |

### Hardware Scaling Path

| Device | Role | Model size |
|---|---|---|
| **MBP M2 Max 32GB** (current) | Primary interface runtime | 7B quantized (Q4_K_M) |
| **Mac Studio M4 Ultra 192GB** (upgrade) | Full 70B local + multiple 7B | Dream tier |
| **Edge device on glasses** | Whisper-tiny for wake word, stream audio to MBP | — |
| **Drone onboard** | No LLM — pure control loop | — |

---

## LLM Speed Comparison

| Provider | Model | TTFT | Tokens/sec | Use |
|---|---|---|---|---|
| **Local (llama.cpp)** | Qwen2.5-7B Q4 | **~50ms** | ~50 t/s | Interface layer |
| **Groq** | Llama 3.1 8B | ~80ms | ~800 t/s | Fast cloud fallback |
| **Groq** | Llama 3.3 70B | ~150ms | ~300 t/s | Mid-tier reasoning |
| Bedrock | Claude Sonnet 4.6 | ~800ms+ | ~80 t/s | Deep reasoning |
| Anthropic direct | Claude Haiku 3.5 | ~300ms | ~150 t/s | Light cloud reasoning |

---

## The Omnipresent Pattern

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

### Example 4: Continuous Context Injection (While Speaking)

```
Zeus is mid-sentence talking to the model about dinner plans...

BACKGROUND (simultaneous):
  - events.ha.motion.front_door fires          → model notes: someone at front door
  - perception.vision.doorbell_frame arrives   → model sees: it's a delivery driver
  - events.comms.sms arrives                   → model notes: text from wife

Model does NOT interrupt Zeus mid-sentence.
Model queues context. When Zeus finishes speaking:

Model: "Dinner sounds good. Also — there's a delivery at the front door,
        and Sarah just texted you."

This is CONTINUOUS AWARENESS without interruption.
The model's context updates WHILE it's listening/generating.
```

---

## Implementation Phases

### Phase 1 — Prove It Now (This Week)

1. Deploy Qwen2.5-7B via llama.cpp on MBP (fast text inference, ~50 t/s)
2. Add Whisper.cpp for voice input (local STT)
3. Wire as a NATS service (`cmd.interface.*`)
4. Prove the loop: voice → model → NATS → HA → confirmation → voice
5. Target: **voice in → action out < 500ms**

### Phase 2 — Multi-Interface (Weeks 2-4)

- Same model serves text (Telegram/WebUI), voice (mic), vision (camera/screen)
- NATS subject routing decides which modality the response goes to
- Context injection via KV watch (location, focus mode, what's on screen)

### Phase 3 — Wearables + Fleet (Month 2+)

- Glasses stream video frames → local model processes
- Drone telemetry → NATS → model can override
- Phone as fallback interface when glasses aren't on
- Edge NATS for field operations (leaf node back to home)

### Phase 4 — Full Omnipresence (Month 3+)

- Multiple edge devices all running same model
- NATS leaf nodes bridge all zones
- Reasoning layer available from any zone when connectivity exists
- Full drone/robot fleet with autonomous runtimes + human override

---

## Key Design Decisions

| Decision | Choice | Why |
|---|---|---|
| Voice path | Local (no NATS) | NATS overhead is pointless for mic→model→speaker on same device |
| Commands to world | NATS publish | Nervous system to all devices, fan-out, decoupled |
| Context to model | NATS subscribe | Model is always aware, never queries on demand |
| Reasoning offload | `cmd.agent.reason` | Heavy tasks go to cloud models async |
| Drone control loop | Onboard, no LLM | LLM too slow for 100Hz control. LLM only for human override commands |
| Edge independence | Own NATS server | Works offline. Leaf node optional for sync |
| Model per interface | Same model everywhere | No switching. Modality-native. |

---

## Subject Namespace (Extended for Realtime Core)

| Namespace | Purpose |
|---|---|
| `cmd.ha.*` | Home automation (existing) |
| `cmd.drone.<id>.*` | Drone commands (waypoint, override, rtl) |
| `cmd.robot.<id>.*` | Robot commands |
| `cmd.fleet.*` | Fan-out to all devices in a class |
| `cmd.interface.*` | Interface layer commands |
| `cmd.agent.reason` | Offload to reasoning layer |
| `perception.vision.*` | Camera frames (glasses, drone, security) |
| `perception.audio.*` | Ambient sound events |
| `context.location.*` | Location awareness |
| `events.ha.*` | HA state changes (existing) |
| `events.comms.*` | Incoming messages |
| `events.zeus.*` | Zeus presence/state (existing) |

---

## The Rules (Summary)

1. **Voice/vision ↔ model = LOCAL. No NATS in that path.** NATS carries commands to the world, not your voice to the AI.

2. **Edge device in field = own NATS server.** Leaf node bridges to home when connectivity exists. Drones connect to edge NATS directly. Works offline.

3. **Qwen2.5-Omni-7B = runs on M2 Max today** via Transformers (slower). Pipeline approach (Whisper + Qwen-7B + TTS) is faster until native MLX/llama.cpp support lands.

4. **Omnipresence = model subscribes to everything, always has context, never queries on demand.** NATS push model means the AI already knows what's happening before you ask.

5. **Drones/robots run autonomous control loops.** The LLM is NEVER in the control loop. It can only INTERVENE via command publish. Same pattern as ha-bridge — fire a command, device executes.

6. **One model, all interfaces.** Glasses, phone, desktop, car, field kit — same model, same NATS bus, same experience.

---

---

## Dynamic Watch System — Background Omnipresence

The missing piece that makes the system truly omnipresent: **persistent attention watches** that Legacy manages concurrently in the background.

### What Is a Watch?

| Pattern | What it does | Judgment? |
|---|---|---|
| Cron job | Runs on schedule | ❌ Fires regardless |
| Subscriber rule | Fires on NATS event match | ❌ Pattern match only |
| **Watch** | Legacy monitors a feed, applies judgment, decides when to surface | ✅ LLM decides relevance |

Watches are NOT cron jobs. NOT automations. They are **LLM-supervised monitoring with judgment** — Legacy decides what's relevant, what's urgent, what to auto-act on, and when to surface it without interrupting.

### How It Works In Practice

```
┌────────────────────────────────────────────────────────────────┐
│                    INTERFACE MODEL (always-on)                   │
│                                                                  │
│  ACTIVE CONVERSATION          BACKGROUND WATCHES                │
│  ┌──────────────────┐         ┌──────────────────────────────┐ │
│  │ Zeus talking      │         │ Watch 1: bar door (waiting)  │ │
│  │ about dinner      │         │ Watch 2: drones (monitoring) │ │
│  │ plans...          │         │ Watch 3: Sarah (waiting)     │ │
│  │                   │         │ Watch 4: office (observing)  │ │
│  └──────────────────┘         │ Watch 5: house (all clear)   │ │
│         │                      └──────────────┬───────────────┘ │
│         │                                     │                  │
│         │         ┌───── EVENT FIRES ─────────┘                  │
│         │         │                                              │
│         │         ▼                                              │
│         │   [Sarah arrived home]                                 │
│         │         │                                              │
│         │         ▼                                              │
│         │   Model decides: relevant, not urgent                  │
│         │   → Queue until Zeus finishes sentence                 │
│         │                                                        │
│         ▼                                                        │
│   Zeus finishes speaking                                         │
│         │                                                        │
│         ▼                                                        │
│   Model: "...sounds good for dinner. Also — Sarah just          │
│           got home. And drone 2 is at 18% battery,              │
│           I'm bringing it back."                                 │
│         │                                                        │
│         ├── nats_publish: cmd.drone.02.rtl (auto-action)        │
│         └── Watch 3 (sarah-home): auto-dismissed (one-shot)     │
│                                                                  │
└────────────────────────────────────────────────────────────────┘
```

**Key behaviors:**

1. **Never interrupts mid-sentence** — queues notifications until natural pause
2. **Batches related updates** — "Also, three things happened while you were talking..."
3. **Can auto-act** — drone battery low? Don't ask, just bring it back and inform
4. **Dismisses completed watches** — "tell me when Sarah gets home" = one-shot, auto-removes
5. **Persistent watches stay** — "watch the office" stays active until explicitly told to stop

### Watch Registry Architecture

```
┌─────────────────────────────────────────────────────────┐
│              WATCH REGISTRY (NATS KV: watches)           │
│                                                          │
│  watch/door-bar     │ trigger: events.ha.contact.bar    │
│                     │ condition: state → "open"          │
│                     │ notify: immediate, voice           │
│                                                          │
│  watch/drone-batt   │ trigger: cmd.drone.*.telemetry     │
│                     │ condition: battery < 20%           │
│                     │ notify: immediate + auto-act (RTL) │
│                                                          │
│  watch/sarah-home   │ trigger: events.ha.person.sarah    │
│                     │ condition: state → "home"          │
│                     │ notify: immediate, voice           │
│                     │ auto-dismiss: true (one-shot)      │
│                                                          │
│  watch/office-cam   │ trigger: perception.vision.office  │
│                     │ condition: LLM judgment (unusual)  │
│                     │ notify: when relevant              │
│                                                          │
│  watch/drone-mission│ trigger: cmd.drone.*.status        │
│                     │ condition: mission_complete        │
│                     │ notify: on completion + summary    │
└─────────────────────────────────────────────────────────┘
```

### Watch Types

| Type | Trigger | LLM Judgment? | Example |
|---|---|---|---|
| **Threshold** | Numeric crosses value | ❌ Simple comparison | "Battery below 20%" |
| **State change** | Entity changes to X | ❌ Simple match | "Door opens", "Sarah arrives" |
| **Pattern** | Stream of events | ✅ LLM decides | "Anything unusual at office" |
| **Completion** | Task/mission finishes | ❌ Status flag | "Drones done with patrol" |
| **Continuous** | Ongoing feed | ✅ LLM summarizes | "Watch the drones" (general awareness) |

### Natural Language Creation

No special syntax. No UI. Just say it:

```
Zeus: "Watch the house"
  → watch: events.ha.* (motion, doors, locks)
  → condition: LLM judgment (anything unexpected)
  → persistent until told to stop

Zeus: "Tell me when the delivery arrives"
  → watch: front_door contact + doorbell vision
  → one-shot: auto-dismiss after first trigger

Zeus: "Have the drones patrol and let me know what they find"
  → publishes: cmd.fleet.patrol
  → watch: telemetry + camera feeds
  → notify: summary on completion, immediate on anomaly

Zeus: "Stop watching the office"
  → removes watch from registry
```

### Notification Behavior

1. **Never interrupts mid-sentence** — queues until natural pause
2. **Batches related updates** — "Three things happened while you were talking..."
3. **Can auto-act** — drone battery low? Bring it back, then inform
4. **Dismisses completed one-shots** — "tell me when Sarah gets home" auto-removes
5. **Persistent watches stay** — "watch the office" active until explicitly cancelled
6. **Priority ordering** — security > people > devices > informational

### How NATS Enables This

The model is ALREADY subscribed to everything (`events.ha.>`, `perception.>`, `cmd.drone.>`, `context.>`, `events.comms.>`). Watches are just **filters + judgment** on top of an existing firehose.

- **No new subscriptions needed** — one sub to `>` covers all
- **No polling** — events push in real-time
- **KV for persistence** — watches survive model restart
- **JetStream replay** — catch up on missed events after restart

### Concurrency Model

One model, one process, multiplexed attention:

```python
class OmnipresentModel:
    conversation: ConversationBuffer      # Active voice/text with Zeus
    watches: dict[str, Watch]             # From NATS KV "watches" bucket
    pending_notifications: list           # Queued for next natural pause
    context: dict                         # Always-updating world state

    async def event_loop(self):
        # Simultaneously:
        # - Listen for Zeus voice/text
        # - Receive NATS events
        # - Evaluate events against watches
        # - Queue or auto-act on triggers
        # - Surface notifications at natural pauses
        # - Maintain conversation coherence
```

NOT multithreading. NOT multiple models. One model with multiplexed attention — like a security guard monitoring 12 cameras while having a conversation.

### Full Session Example

```
9:00 AM — Zeus puts on glasses

Zeus: "Good morning. What's the status?"
Legacy: "Morning. House clear — Sarah left 7:30. Office no motion.
         Drones charged. 72°F house, 68°F office. Call at 10."

Zeus: "Send drone 1 to check office perimeter. Watch the bar
       entrance. Let me know when the package arrives."
Legacy: "Drone 1 launching. Bar watch active. Package watch active."

BACKGROUND (simultaneous):
  ├── Drone 1: flying, Legacy tracking
  ├── Bar entrance: waiting on contact sensor
  ├── Package: waiting on doorbell + front door
  └── Home sensors: ambient

9:15 AM — Zeus in a meeting
  [Drone 1 completes — no anomalies] → queued
  [Package arrives — FedEx on doorbell] → queued

9:20 AM — Meeting pauses
Zeus: "Updates?"
Legacy: "Drone 1 finished perimeter — all clear, back on pad.
         Package arrived, FedEx at front door."
  → Package watch: auto-dismissed
  → Drone: standby

9:45 AM — Zeus driving
  [Bar door opens — unknown person]
Legacy: "Someone entered the bar. One person, not recognized staff.
         Want the feed?"
Zeus: "Show me. ...That's the new cleaner. Stand down."
Legacy: "Got it, adding to recognized faces."

10:30 AM
Zeus: "Drones 2 and 3 patrol the property. Watch the gates.
       Tell me when Sarah texts back."
Legacy: "Launching. Gates active. Sarah text active."

WATCHES NOW RUNNING:
  ├── Bar entrance (persistent)
  ├── Drones 2+3 patrol (telemetry + cameras)
  ├── Gates (motion + vision)
  ├── Sarah text (events.comms.sms.sarah)
  └── Home + office sensors (ambient)

11:00 AM — All quiet
Legacy: (silence — nothing to report)

11:15 AM — Drone 3 at 19%
Legacy: "Drone 3 at 19%. Bringing it home, drone 2 continuing."
  → Auto-action: cmd.drone.03.rtl

11:30 AM — Sarah texts "picking up kids, home by 4"
Legacy: "Sarah's picking up kids, home by 4."
  → Watch: auto-dismissed
```

### What Makes This Different

| System | Persistent watches | LLM judgment | Multi-watch | Dynamic creation | Background awareness |
|---|---|---|---|---|---|
| Siri/Alexa | ❌ | ❌ | ❌ | ❌ | ❌ |
| HA Automations | ✅ | ❌ | ✅ | ❌ (config files) | ❌ |
| ChatGPT | ❌ | ✅ | ❌ | ❌ | ❌ |
| Apple Intelligence | ❌ | Partial | ❌ | ❌ | ❌ |
| **OneMind Watches** | ✅ | ✅ | ✅ | ✅ (voice) | ✅ |

### Implementation Status

| Component | Status |
|---|---|
| NATS pub/sub for all events | ✅ Live |
| KV for watch registry | 🔨 Trivial (new bucket) |
| ha-bridge for home events | ✅ Live (~1ms) |
| Subscriber for threshold watches | ✅ Live |
| Local model for LLM judgment | 📋 Phase 1 |
| Voice I/O | 📋 Phase 1 |
| Vision processing | 📋 Phase 2 |
| Drone telemetry | 📋 Phase 3 |

### Continuous Context Injection — While Speaking

The model's context updates **while Zeus is talking**. Events don't wait for the conversation to end.

```
Zeus is mid-sentence talking to the model about dinner plans...

BACKGROUND (simultaneous):
  - events.ha.motion.front_door fires          → model notes: someone at front door
  - perception.vision.doorbell_frame arrives   → model sees: it's a delivery driver
  - events.comms.sms arrives                   → model notes: text from wife

Model does NOT interrupt Zeus mid-sentence.
Model queues context. When Zeus finishes speaking:

Model: "Dinner sounds good. Also — there's a delivery at the front door,
        and Sarah just texted you."

This is CONTINUOUS AWARENESS without interruption.
The model's context updates WHILE it's listening/generating.
```

### Auto-Action Policy

When Legacy can act without asking:

| Situation | Auto-act? | Why |
|---|---|---|
| Drone battery critical (<10%) | ✅ RTL immediately | Safety — don't risk crash |
| Drone battery low (<20%) | ✅ RTL + inform | Preserve hardware |
| Unknown person at secured location | ❌ Inform + ask | Judgment needed from Zeus |
| Known person arrives | ✅ Inform only | No action needed |
| Fire/smoke detected | ✅ Alert + call 911 | Life safety |
| Zeus says "lights off" | ✅ Execute | Direct command |
| Unusual pattern at office | ❌ Inform + ask | Could be normal |
| Mission complete | ✅ Inform + stand down | Expected completion |

**Principle: Auto-act on safety and direct commands. Ask on judgment calls.**

### Watch Lifecycle

```
CREATED ──→ ACTIVE ──→ TRIGGERED ──→ NOTIFIED ──→ DISMISSED (one-shot)
                │                                      │
                │                                      └──→ BACK TO ACTIVE (persistent)
                │
                └──→ CANCELLED (Zeus: "stop watching X")
```

Watches stored in NATS KV bucket `watches`. Key format: `watch/<name>`.

```json
{
  "name": "sarah-home",
  "type": "state_change",
  "trigger_subject": "events.ha.person.sarah",
  "condition": {"state": "home"},
  "notify": "immediate",
  "modality": "voice",
  "one_shot": true,
  "auto_action": null,
  "created": "2026-05-24T09:00:00Z",
  "created_by": "voice",
  "last_triggered": null,
  "trigger_count": 0
}
```

### Integration With Existing Subscriber Architecture

Watches **complement** the existing subscriber system, they don't replace it:

| Layer | What | Intelligence | Latency |
|---|---|---|---|
| **Subscribers** (existing) | Rule-based event → action | ❌ No LLM, pattern match | <1ms |
| **Watches** (new) | LLM-supervised monitoring | ✅ Judgment, context, batching | ~300ms |

Simple threshold watches (battery < 20%, door opens) can run as pure subscribers (no LLM needed). Complex watches (anything unusual, track behavior patterns) need the local model.

**Hybrid pattern:** Subscriber fires the fast action (RTL drone), AND notifies the watch system so Legacy can inform Zeus. Both layers work together.

---

## Related

- [[NATS]] — central bus architecture
- [[NATS-Intent-Mesh]] — subject namespace map
- [[HA-NATS-Bridge]] — existing HA integration (proven pattern)
- [[Hermes-Overview]] — current Legacy/Hermes agent layer (will become the reasoning layer)

---
domain: 04
name: OpenMind OM1 — Deep Feature Analysis
status: reference
owner: zeus
tags: [robotics, om1, openmind, llm, research, drone-eval]
created: 2026-05-12
updated: 2026-05-12
---

# OpenMind OM1 — Deep Feature Analysis

> Source: https://github.com/OpenMind/OM1 | https://openmind.com | https://docs.openmind.com
> MIT License | 2.8k stars | 122 contributors | v1.0.2-beta.2

OM1 is a **modular AI runtime for robots** — a perception-decision-action loop that connects LLMs to physical hardware. It is NOT a drone platform. It excels at ground robots (humanoids, quadrupeds) with conversational AI.

---

## 1. Core Architecture (The Brain Loop)

OM1 runs a **fixed-frequency perception-decision-action loop** at `self.config.hertz`. This is NOT the robot's motor loop (50-500 Hz) — it's the **cognitive refresh rate** for attention and working memory.

```
Raw Sensors → AI Captioning → NLDB → State Fuser → Cortex LLM(s) → Action Orchestrator → HAL → Robot
     ↑                                                                                            │
     └────────────────────────── feedback loop ────────────────────────────────────────────────────┘
```

### Layer Breakdown

| Layer | What It Does |
|---|---|
| **Raw Sensors** | Camera, microphone, LiDAR, GPS, battery, IMU, depth sensors |
| **AI Captioning** | Converts raw data to natural language — VLM describes video, ASR transcribes speech, spatial/nav compresses LiDAR/GPS |
| **NLDB (Natural Language Data Bus)** | Central bus collecting all captioned data as structured text strings |
| **State Fuser** | Combines NLDB inputs into a ~1 paragraph prompt with temporal context (e.g. `"137.02: You see a human, 3.2m left, smiling, pointing at chair. You heard: 'Bits run to the chair'"`) |
| **Cortex LLM(s)** | Reasoning engine — single, dual, or parallel LLM strategies |
| **Action Orchestrator** | Manages state machines, promise queues, threads for each action |
| **HAL** | Maps high-level decisions to ROS2/Zenoh/CycloneDDS/serial/USB commands |

### Project Structure

```
src/
├── actions/          # Output plugins (move, speak, emotion, etc.)
├── backgrounds/      # Long-running sensor tasks (GPS, ODOM, RPLIDAR, RTK)
├── fuser/            # Input fusion logic
├── inputs/           # Sensor plugins (VLM, ASR, LiDAR, governance)
├── llm/              # LLM plugin system
├── providers/        # Background task providers
├── runtime/          # Core runtime engine
├── simulators/       # WebSim + others
├── zenoh_msgs/       # Zenoh IDL definitions
└── run.py            # CLI entry point
```

---

## 2. LLM System (3 Execution Strategies)

This is OM1's most sophisticated subsystem. Three modes:

### Single LLM
Standard — one model handles everything. Use OpenAI, DeepSeek, Gemini, xAI, NearAI, OpenRouter, or Ollama local.

### Dual LLM (Local + Cloud)
- Sends prompt to **both** Qwen3-30B (local) and GPT-4.1 (cloud) **in parallel**
- Waits up to **3.2 seconds**
- If both respond: local LLM **evaluates which answer is better** and picks the winner
- If only one responds: uses that one directly
- Optimizes for latency + quality simultaneously

### Parallel LLM (Multi-Specialist)
- Multiple LLMs run simultaneously, each **filtered to specific actions**
- e.g. GPT-4.1 handles `speak` + `emotion`, Qwen handles `move` + `navigate`, DeepSeek handles `search` + `analyze`
- Results stream as they complete — **no waiting for slowest model**
- `execute_immediately: true` means actions fire as soon as any LLM responds

### Multi-Agent Internal Architecture
The cloud endpoint wraps **3+ internal agents**:

| Agent | Role | Latency |
|---|---|---|
| **Fast Action LLM** | Time-critical immediate responses | ~300ms |
| **Core Cognition LLM** | Complex reasoning, planning | ~2s |
| **Mentor/Coach LLM** | 3rd-person critique of robot-human interaction | every ~30s |

### Supported Models (as of May 2026)

| Provider | Models |
|---|---|
| **OpenAI** | gpt-4o, gpt-4o-mini, gpt-4.1, gpt-4.1-mini, gpt-4.1-nano, gpt-5, gpt-5-mini, gpt-5-nano |
| **DeepSeek** | deepseek-chat |
| **Gemini** | gemini-3.1-pro-preview, gemini-3.1-flash-lite-preview, gemini-3-pro-preview, gemini-3-flash-preview, gemini-2.5-flash, gemini-2.5-flash-lite, gemini-2.5-pro |
| **xAI** | grok-2-latest, grok-3-beta, grok-4-latest, grok-4 |
| **NearAI** | qwen3-30b-a3b-instruct-2507, qwen2.5-vl-72b-instruct, qwen-2.5-7b-instruct |
| **OpenRouter** | llama-3.1-70b, llama-3.3-70b, claude-sonnet-4.5, claude-opus-4.1 |
| **Ollama (local)** | llama3.2, llama3.1, mistral, phi3, gemma2, qwen2.5, codellama, llava |
| **Local native** | Qwen3-30B (on-device) |

---

## 3. Input Plugin System

Sensors feed the NLDB. Architecture: `Sensor<T>` → `FuserInput<T>` → `InputOrchestrator`

### Built-in Input Types

| Plugin | What It Does |
|---|---|
| `GoogleASRInput` | Automatic speech recognition (mic → text) |
| `RivaASR` | NVIDIA Riva ASR (local, GPU-accelerated) |
| `VLM_COCO_Local` | Local vision-language model with COCO object detection |
| `VLMVila` | ViLA vision-language model |
| `RPLidar` | 360° laser scan input |
| `TurtleBot4Batt` | TurtleBot 4 battery state |
| `UnitreeG1Basic` | Unitree G1 humanoid state |
| `UnitreeGo2Lowstate` | Unitree Go2 quadruped low-level state |
| `GovernanceEthereum` | On-chain robot governance rules (from Ethereum) |
| `GPS (Arduino)` | GPS coordinates via Arduino |

Custom inputs: extend `FuserInput<T>`, implement your sensor logic, register in config.

---

## 4. Action Plugin System

Actions are the robot's **outputs** — what it does in the world.

### Architecture Per Action

```
actions/move_{hardware_id}/
├── interface.py          # Defines MoveInput/Output types
├── implementation/
│   └── passthrough.py    # Business logic (or pass-through)
└── connector/
    ├── ros2.py           # ROS2 connector
    ├── zenoh.py          # Zenoh connector
    └── unitree.py        # Direct Unitree SDK
```

### Action Types

| Action | Connectors |
|---|---|
| **Move** | ROS2, Zenoh, Unitree SDK direct, CycloneDDS |
| **Speak** | ElevenLabs TTS, other TTS engines |
| **Emotion** | Facial expressions via avatar display |
| **Sound Effects** | Audio output |
| **Navigate** | Nav2 autonomous navigation |

Each action has its own **state machine, promise queue, and thread** managed by the Action Orchestrator.

---

## 5. Mode & Transition System

Robots can have **multiple operational modes** with automatic switching.

### Mode Structure
Each mode defines its own:
- `system_prompt_base` (personality/behavior)
- `hertz` (cognitive refresh rate)
- `agent_inputs` (which sensors are active)
- `agent_actions` (which outputs are available)
- `mcp_servers` (which external tools)
- Per-mode LLM override (different models per mode)

### Transition Rules
Modes switch based on:
- **`input_triggered`** — keyword detection in user speech
- **`context_conditions`** — environmental state (battery, time, proximity)
- **Priority resolution** — higher priority wins when multiple rules match
- **Cooldown** — prevents mode oscillation

Example: "reset" → Welcome Mode (priority 5, 10s cooldown)

---

## 6. MCP Integration

OM1 natively supports **Model Context Protocol** — "USB for AI."

- Configured per-mode (different modes can have different MCP servers)
- `MCPClientManager` per mode — isolation between modes
- Supports `stdio` and `http` transports
- LLM decides when to call tools — auto-discovery from connected servers
- Examples: Slack, GitHub, Google Maps, Weather, Filesystem, any community MCP server

Config example:
```json5
mcp_servers: [
  {
    name: "weather",
    transport: "stdio",
    command: "npx",
    args: ["-y", "@h1deya/mcp-server-weather"],
  },
  {
    name: "github",
    transport: "http",
    url: "https://api.githubcopilot.com/mcp/",
    headers: {
      Authorization: "Bearer ${GITHUB_PERSONAL_ACCESS_TOKEN}",
    },
  },
]
```

---

## 7. Full Autonomy Stack (Premium)

### Open Source Components

| Component | What |
|---|---|
| **OM1** | Core intelligence loop — ASR, LLM, action dispatch |
| **OM1 Avatar** | React frontend on BrainPack screen — shows state, emotions, speech |

### Premium Components (Enterprise $999/mo)

| Component | What |
|---|---|
| **OM1 ROS2 SDK** | Full robotics middleware — SLAM, Nav2, sensor drivers, obstacle avoidance, auto-charging |
| **Video Processor** | On-device face anonymization, audio cleanup, RTSP streaming (CUDA/TensorRT) |
| **Person Following** | Visual detection + spatial reasoning + autonomous tracking |
| **OTA System** | Containerized updates — `ota_agent` + `ota_updater` for fleet management |

### REST APIs (Premium)

**Orchestrator API (Port 5000):**
- `POST /start/base_control` — Enable motors
- `POST /start/slam` — Start mapping
- `POST /start/nav2` — Start autonomous navigation with saved map
- `POST /patrol/start` — Autonomous patrol with waypoints (Go2 only)
- `POST /charging/dock` — Auto-dock to charger (Go2 only, AprilTag visual servoing)
- `POST /maps/save`, `GET /maps/list`, `POST /maps/delete`
- `POST /maps/locations/add/slam` — Save current position as waypoint

**Nav2 API (Port 5001):**
- `GET /api/pose` — Robot position + covariance
- `POST /api/move_to_pose` — Navigate to coordinates
- `GET /api/nav2_status` — Track navigation goals
- `GET /api/map` — Current occupancy grid
- `GET /api/amcl_variance` — Localization confidence

---

## 8. Hardware Support

| Platform | Status |
|---|---|
| **Unitree Go2** (quadruped) | Full support — nav, SLAM, charging, patrol |
| **Unitree G1** (humanoid) | Full support — arm actions, locomotion |
| **LimX Tron** | Supported |
| **TurtleBot 4** | Supported (educational) |
| **UbTech Yanshee** | Supported (mini humanoid) |
| **Gazebo** (sim) | Open source sim support |
| **Isaac Sim** (NVIDIA) | GPU-accelerated physics sim |
| **Drones** | **NOT supported** |

Middleware: **ROS2**, **Zenoh** (recommended for new dev), **CycloneDDS**, serial, USB, custom APIs

---

## 9. Governance & Blockchain

Unique to OM1: robot constitutions/guardrails can be stored on **Ethereum** (immutable public ledger).

- `GovernanceEthereum` input plugin reads rules from chain
- `system_governance` field in config defines Asimov-style laws
- Enables: transparency, traceability, decentralized coordination, accountability logging
- **URID** (Universal Robot ID) — joins decentralized machine-to-machine coordination network (**FABRIC**)
- **FABRIC** = decentralized robot coordination protocol (multi-robot communication)

---

## 10. Pricing

| Plan | Price | OMCU | Rate Limit | Key Features |
|---|---|---|---|---|
| **Free** | $0/mo | 50 | 1 req/s | 1 API key |
| **Starter** | $5/mo | 2,500 | 5 req/s | Unlimited keys, session memory |
| **Builder** | $25/mo | 15,000 | 20 req/s | Priority processing, dedicated sims, email support, App Store publishing |
| **Pro** | $99/mo | 100,000 | 40 req/s | Unlimited long-term memory, Slack support |
| **Enterprise** | $999/mo | 1,500,000 | 100 req/s | Premium features (ROS2 SDK, video processor, person following), custom integration |
| **Academic** | Free | Full access | — | Research/student application |

---

## 11. Verdict for Zeus / One Mind Drone Use Case

### What OM1 Does Well
- Incredibly polished LLM-to-robot pipeline (dual/parallel LLM is genuinely clever)
- MCP integration is first-class
- Mode system with auto-transitions is production-grade
- Blockchain governance is unique — no other framework does this
- Full autonomy stack (SLAM, Nav2, patrol, auto-charging) is enterprise-ready for **ground robots**
- Multi-agent internal architecture (fast/core/coach LLMs) is sophisticated

### What OM1 Does NOT Do
- **Zero drone support** — no MAVLink, no DJI, no flight controller integration
- **No swarm/fleet coordination** at the application layer (FABRIC is still early/opaque)
- Premium features locked behind **$999/mo Enterprise plan**
- Requires their API key for cloud features (OMCU billing)
- All their autonomy stack is ground-robot-only (Unitree Go2/G1/LimX)

### Verdict vs DimOS for Drones
**DimOS wins categorically for the drone use case.** OM1 is the better platform for ground robots + conversation + enterprise deployments, but it has literally nothing for aerial vehicles. DimOS has MAVLink, DJI Mavic support, GPS fly-to, visual servoing, object tracking, and LLM agent control of drones — all needed.

That said, OM1's **architecture patterns** (dual LLM, parallel LLM, mode transitions, blockchain governance) are worth studying for how to enhance the DimOS + OpenClaw stack.

---

## Key Architecture Patterns Worth Adopting

1. **Dual LLM (local + cloud parallel race)** — Apply to DimOS drone decisions for latency optimization
2. **Parallel LLM with action filters** — Different models for flight vs. perception vs. communication
3. **Mode transitions with cooldowns** — Drone operational modes (patrol, follow, return-to-home, emergency)
4. **Natural Language Data Bus** — Unify all sensor captions into one structured text stream
5. **Blockchain governance** — Immutable flight rules and geofencing on-chain
6. **URID / FABRIC** — Decentralized drone-to-drone coordination protocol

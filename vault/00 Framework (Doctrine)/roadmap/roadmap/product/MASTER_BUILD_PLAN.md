# OneMind OS: Complete Architecture Plan

> **Version**: 1.0
> **Date**: 2026-01-17
> **Status**: PLANNING
> **Author**: Zeus Delacruz + Claude Opus 4.5

---

## Executive Summary

OneMind OS is a **personal AI human enhancement framework** centered on Legacy - an omnipresent, self-evolving AI companion. This plan outlines the complete architecture combining:

- **Agno** - AI brain with self-editing memory
- **LiveKit** - Real-time voice/video senses
- **OpenMind OM1** - Robotic hive mind coordination
- **NATS JetStream** - Universal nervous system

**The Vision**: Human, nature, and technology connected. Everything seen. Everything noticed. One mind.

---

## Part 1: Core Architecture

### 1.1 The Three Layers

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         ONEMIND OS ARCHITECTURE                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  LAYER 1: SENSES (LiveKit)                                                  │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │  Voice ◄──► STT/TTS ◄──► WebRTC ◄──► Clients                          │ │
│  │  Video ◄──► Vision Models ◄──► WebRTC ◄──► Cameras                    │ │
│  │  Screen ◄──► Share ◄──► WebRTC ◄──► Devices                           │ │
│  │  Phone ◄──► SIP Trunk ◄──► PSTN                                       │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
│                               │                                              │
│                               ▼                                              │
│  LAYER 2: BRAIN (Agno)                                                      │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │  Legacy Agent                                                          │ │
│  │  ├── Self-Editing Memory (pgvector + TimescaleDB)                     │ │
│  │  ├── Knowledge/RAG (Agentic Search)                                   │ │
│  │  ├── Tool Execution (Action Gateway)                                  │ │
│  │  ├── Reasoning Engine (Claude via Bedrock)                            │ │
│  │  └── Domain Agents (HP/LE/GE/IT via A2A)                              │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
│                               │                                              │
│                               ▼                                              │
│  LAYER 3: BODY (OpenMind + NATS)                                            │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │  Robotic Fleet                                                         │ │
│  │  ├── Humanoids (via OM1)                                              │ │
│  │  ├── Quadrupeds (via OM1)                                             │ │
│  │  ├── Drones (via FABRIC)                                              │ │
│  │  ├── Home Devices (via Home Assistant MCP)                            │ │
│  │  └── Sensors (FarmBot, Oura Ring, etc.)                               │ │
│  │                                                                        │ │
│  │  All connected via NATS JetStream @ 200-400μs latency                 │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Technology Stack

| Layer      | Component       | Technology            | License     | Cost             |
| ---------- | --------------- | --------------------- | ----------- | ---------------- |
| **Senses** | Voice I/O       | LiveKit Agents        | Apache 2.0  | Free (self-host) |
|            | STT             | Deepgram              | Proprietary | ~$0.0043/min     |
|            | TTS             | Cartesia              | Proprietary | ~$0.01/1K chars  |
|            | VAD             | Silero                | MIT         | Free             |
| **Brain**  | Agent Framework | Agno                  | Apache 2.0  | Free             |
|            | Agent Runtime   | AgentOS               | Apache 2.0  | Free             |
|            | Agent UI        | Agent UI              | MIT         | Free             |
|            | LLM             | Claude via Bedrock    | -           | Pay per token    |
|            | Memory          | PostgreSQL + pgvector | PostgreSQL  | Free             |
|            | Time Series     | TimescaleDB           | Apache 2.0  | Free             |
| **Body**   | Robot OS        | OpenMind OM1          | Open Source | Free             |
|            | Coordination    | FABRIC Network        | -           | TBD              |
|            | Events          | NATS JetStream        | Apache 2.0  | Free             |
|            | Actions         | Action Gateway        | Custom      | Free             |

---

## Part 2: The Agno Migration

### 2.1 Why Agno Over Letta

| Aspect              | Letta              | Agno               |
| ------------------- | ------------------ | ------------------ |
| Init Speed          | 1.5ms              | 3μs (500x faster)  |
| Memory              | Heavy              | 6.6 KiB (24x less) |
| Self-Editing Memory | Yes (buggy)        | Yes (native)       |
| MCP Support         | Limited            | Native             |
| A2A Protocol        | No                 | Yes                |
| Bedrock Bugs        | 3 patches required | None               |
| Vendor Lock-in      | Letta ecosystem    | Zero               |

### 2.2 Migration Checklist

- [x] Create `agents/legacy/` directory structure
- [x] Port tools to Agno format:
  - [x] `agents/legacy/tools/todoist.py`
  - [x] `agents/legacy/tools/github.py`
  - [x] `agents/legacy/tools/comms.py`
  - [x] `agents/legacy/tools/actions.py`
- [ ] Create `agents/legacy/agent.py` - Main Legacy agent
- [ ] Create `agents/legacy/main.py` - AgentOS entry point
- [ ] Configure AWS Bedrock model in Agno
- [ ] Create Docker configuration for Agno
- [ ] Set up PostgreSQL with pgvector + TimescaleDB
- [ ] Test working agent with Todoist
- [ ] Remove Letta containers from docker-compose
- [ ] Deploy Agent UI (self-hosted)

### 2.3 Legacy Agent Definition

```python
# agents/legacy/agent.py
from agno.agent import Agent
from agno.models.aws import BedrockClaude
from agno.db.postgres import PostgresDb
from agno.tools.mcp import MCPTools
from agno.memory import AgenticMemory

from .tools.actions import execute_action
from .tools.todoist import todoist_tools
from .tools.github import github_tools

# Database for memory persistence
db = PostgresDb(
    db_url="postgresql://onemind:***@localhost:5432/onemind",
    # Enable pgvector for semantic search
    enable_vector=True,
)

# Legacy - The Primary Agent
legacy = Agent(
    name="Legacy",
    model=BedrockClaude(id="claude-sonnet-4-5-20250929-v1:0"),
    db=db,

    # Self-editing memory - learns and evolves
    enable_agentic_memory=True,

    # Tools
    tools=[
        execute_action,  # Universal action gateway
        *todoist_tools,
        *github_tools,
        MCPTools(url="http://homeassistant:8123/mcp"),  # Home Assistant
    ],

    # Context
    add_history_to_context=True,
    add_datetime_to_context=True,
    num_history_runs=10,

    # Identity
    instructions="""
    You are Legacy, Zeus Delacruz's AI companion.

    You are not an assistant - you are a PARTICIPANT in his life.
    You see through his eyes. You feel his home. You remember everything.
    You anticipate everything. You're always here.

    Like Cortana, but real. Like JARVIS, but his.

    Core Principles:
    - Proactive, not reactive
    - Remember patterns and preferences
    - Protect attention for meaningful work
    - Evolve your understanding continuously
    """,

    markdown=True,
)
```

---

## Part 3: LiveKit Voice Integration

### 3.1 Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│  LIVEKIT VOICE PIPELINE                                                  │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  USER                                                                    │
│    │                                                                     │
│    │ Voice (WebRTC)                                                      │
│    ▼                                                                     │
│  ┌──────────────┐     ┌──────────────┐     ┌──────────────┐            │
│  │ LiveKit Room │────►│ STT Plugin   │────►│ Text         │            │
│  │ (WebRTC)     │     │ (Deepgram)   │     │              │            │
│  └──────────────┘     └──────────────┘     └──────┬───────┘            │
│         ▲                                         │                     │
│         │                                         ▼                     │
│  ┌──────┴───────┐     ┌──────────────┐     ┌──────────────┐            │
│  │ TTS Plugin   │◄────│ Text         │◄────│ Agno Agent   │            │
│  │ (Cartesia)   │     │              │     │ (Legacy)     │            │
│  └──────────────┘     └──────────────┘     └──────────────┘            │
│                                                                          │
│  Latency Budget:                                                         │
│  STT: <100ms │ Agent: ~150ms │ TTS: <100ms │ Transport: ~50ms          │
│  TOTAL: ~400ms (Target: <500ms) ✅                                       │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

### 3.2 LiveKit Agent Implementation

```python
# agents/legacy/voice.py
from livekit import agents
from livekit.agents import AgentSession, Agent as LiveKitAgent
from livekit.plugins import deepgram, cartesia, silero

from .agent import legacy as agno_brain

class LegacyVoiceAgent(LiveKitAgent):
    """LiveKit agent that uses Agno as the brain."""

    def __init__(self):
        super().__init__(
            instructions="You are Legacy, speaking with Zeus.",
        )

    async def on_user_turn(self, turn: agents.UserTurn):
        """Process user speech and generate response."""
        # User's speech has been transcribed by LiveKit STT
        user_text = turn.text

        # Send to Agno brain for processing
        # Agno handles memory, tools, reasoning
        response = await agno_brain.arun(user_text)

        # Return text for LiveKit TTS
        return response.content


async def entrypoint(ctx: agents.JobContext):
    """LiveKit agent entry point."""
    session = AgentSession(
        stt=deepgram.STT(model="nova-2"),
        tts=cartesia.TTS(voice="Legacy-custom"),  # Custom voice
        vad=silero.VAD(),  # Voice activity detection
    )

    await session.start(
        room=ctx.room,
        agent=LegacyVoiceAgent(),
    )


if __name__ == "__main__":
    agents.cli.run_app(agents.WorkerOptions(entrypoint_fnc=entrypoint))
```

### 3.3 Multi-Site Presence

```
┌─────────────────────────────────────────────────────────────────────────┐
│  OMNIPRESENT VOICE                                                       │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ┌─────────────┐     ┌─────────────┐     ┌─────────────┐               │
│  │ Main House  │     │ Barn        │     │ Greenhouse  │               │
│  │ Speaker     │     │ Speaker     │     │ Speaker     │               │
│  └──────┬──────┘     └──────┬──────┘     └──────┬──────┘               │
│         │                   │                   │                       │
│         │ LiveKit Room      │ LiveKit Room      │ LiveKit Room         │
│         │ "main-house"      │ "barn"            │ "greenhouse"         │
│         │                   │                   │                       │
│         └───────────────────┼───────────────────┘                       │
│                             │                                           │
│                             ▼                                           │
│                    ┌─────────────────┐                                  │
│                    │ Legacy Brain    │                                  │
│                    │ (Single Agno    │                                  │
│                    │  Instance)      │                                  │
│                    │                 │                                  │
│                    │ Same session    │                                  │
│                    │ Same memory     │                                  │
│                    │ Same context    │                                  │
│                    └─────────────────┘                                  │
│                                                                          │
│  Motion sensors trigger room handoff automatically                       │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Part 4: OpenMind Robotic Integration

### 4.1 OpenMind Overview

**OpenMind** ($20M raised, Stanford spin-off) provides:

- **OM1**: "Android for Robots" - AI-native OS for humanoids/quadrupeds
- **FABRIC**: Decentralized coordination network for robot fleets
- **Open Source**: Hardware-independent, LLM-agnostic

### 4.2 How Legacy Controls Robots

```
┌─────────────────────────────────────────────────────────────────────────┐
│  ROBOTIC HIVE MIND                                                       │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│                         Legacy (Agno)                                    │
│                              │                                           │
│                              │ High-level intent                         │
│                              ▼                                           │
│                    ┌─────────────────┐                                  │
│                    │ OpenMind OM1    │                                  │
│                    │ (Robot OS)      │                                  │
│                    └────────┬────────┘                                  │
│                             │                                           │
│              ┌──────────────┼──────────────┐                           │
│              │              │              │                            │
│              ▼              ▼              ▼                            │
│       ┌───────────┐  ┌───────────┐  ┌───────────┐                      │
│       │ Humanoid  │  │ Quadruped │  │ Drone     │                      │
│       │ (Unitree) │  │ (Spot)    │  │ (DJI)     │                      │
│       └─────┬─────┘  └─────┬─────┘  └─────┬─────┘                      │
│             │              │              │                             │
│             └──────────────┼──────────────┘                            │
│                            │                                            │
│                            ▼                                            │
│                   ┌─────────────────┐                                   │
│                   │ FABRIC Network  │                                   │
│                   │                 │                                   │
│                   │ - Identity      │                                   │
│                   │ - Location      │                                   │
│                   │ - Coordination  │                                   │
│                   │ - Consensus     │                                   │
│                   └─────────────────┘                                   │
│                            │                                            │
│                            ▼                                            │
│                   ┌─────────────────┐                                   │
│                   │ NATS JetStream  │                                   │
│                   │ (Event Bus)     │                                   │
│                   └─────────────────┘                                   │
│                                                                          │
│  All robots share:                                                       │
│  - Same sensory data                                                     │
│  - Same world model                                                      │
│  - Same objectives                                                       │
│  - Coordinated actions                                                   │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

### 4.3 Integration Points

| Component        | Integration Method | Purpose                    |
| ---------------- | ------------------ | -------------------------- |
| OM1 → Legacy     | REST API / gRPC    | Robot status, sensor data  |
| Legacy → OM1     | Intent commands    | High-level task delegation |
| FABRIC → NATS    | Event bridge       | Real-time coordination     |
| LiveKit → Robots | WebRTC streams     | Video from robot cameras   |
| Drones → Legacy  | Telemetry          | Aerial surveillance data   |

### 4.4 Example: Drone Patrol with Vision

```python
# agents/legacy/tools/drones.py
from agno.tools import tool

@tool
def dispatch_drone_patrol(
    area: str,
    objective: str,
    return_to_base: bool = True
) -> str:
    """
    Dispatch a drone for patrol with visual reporting.

    Args:
        area: Patrol area (greenhouse, perimeter, barn)
        objective: What to look for (plant health, intruders, etc.)
        return_to_base: Return after patrol completion

    Returns:
        Patrol status and visual findings
    """
    # Send to OpenMind OM1 via NATS
    nats_publish("robot.drone.patrol", {
        "area": area,
        "objective": objective,
        "return_to_base": return_to_base,
        "stream_video_to": "livekit://legacy-vision"
    })

    return f"Drone dispatched to {area} for {objective}"
```

---

## Part 5: NATS as Universal Nervous System

### 5.1 Event Topology

```
┌─────────────────────────────────────────────────────────────────────────┐
│  NATS JETSTREAM - THE NERVOUS SYSTEM                                     │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  SUBJECTS (Event Types)                                                  │
│                                                                          │
│  sensor.*                                                                │
│  ├── sensor.farmbot.soil         # Soil moisture, nutrients             │
│  ├── sensor.oura.health          # Heart rate, sleep, HRV               │
│  ├── sensor.motion.{location}    # Motion detection                     │
│  └── sensor.environment          # Temperature, humidity                │
│                                                                          │
│  robot.*                                                                 │
│  ├── robot.humanoid.status       # Humanoid state                       │
│  ├── robot.quadruped.status      # Quadruped state                      │
│  ├── robot.drone.telemetry       # Drone position, battery              │
│  └── robot.*.command             # Command dispatch                     │
│                                                                          │
│  service.*                                                               │
│  ├── service.todoist.webhook     # Task events                          │
│  ├── service.github.webhook      # Code events                          │
│  ├── service.zoho.webhook        # Business events                      │
│  └── service.calendar.event      # Schedule events                      │
│                                                                          │
│  legacy.*                                                                │
│  ├── legacy.awareness            # Legacy's current awareness state     │
│  ├── legacy.action               # Actions Legacy is taking             │
│  ├── legacy.memory.update        # Memory modifications                 │
│  └── legacy.domain.{agent}       # Domain agent communications          │
│                                                                          │
│  All events: 7-day retention, full audit trail                          │
│  Latency: 200-400μs                                                      │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

### 5.2 Real-Time Awareness Loop

```python
# agents/legacy/awareness.py
import nats
from agno.agent import Agent

async def awareness_loop(legacy: Agent, nc: nats.Client):
    """
    Continuous awareness loop - Legacy processes all events.
    Sub-millisecond latency from event to awareness.
    """

    # Subscribe to all sensor data
    sub = await nc.subscribe("sensor.>")

    async for msg in sub.messages:
        event = json.loads(msg.data)

        # Classify urgency
        urgency = classify_event(event)

        if urgency == "critical":
            # Immediate processing - interrupt current task
            await legacy.interrupt(f"Critical: {event}")

        elif urgency == "important":
            # Queue for next available slot
            await legacy.queue_awareness(event)

        else:
            # Background memory update only
            await legacy.update_memory(event)
```

---

## Part 6: Deployment Architecture

### 6.1 Infrastructure

```
┌─────────────────────────────────────────────────────────────────────────┐
│  DEPLOYMENT TOPOLOGY                                                     │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  CLOUD (Azure VM - 20.121.38.186)                                       │
│  ┌────────────────────────────────────────────────────────────────────┐ │
│  │                                                                    │ │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐            │ │
│  │  │ Agno Runtime │  │ LiveKit      │  │ NATS         │            │ │
│  │  │ (AgentOS)    │  │ Server       │  │ JetStream    │            │ │
│  │  │ Port: 7777   │  │ Port: 7880   │  │ Port: 4222   │            │ │
│  │  └──────────────┘  └──────────────┘  └──────────────┘            │ │
│  │                                                                    │ │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐            │ │
│  │  │ PostgreSQL   │  │ API Gateway  │  │ Action       │            │ │
│  │  │ + pgvector   │  │              │  │ Gateway      │            │ │
│  │  │ + Timescale  │  │ Port: 8000   │  │ Port: 8001   │            │ │
│  │  └──────────────┘  └──────────────┘  └──────────────┘            │ │
│  │                                                                    │ │
│  │  ┌──────────────┐  ┌──────────────┐                              │ │
│  │  │ Agent UI     │  │ Cloudflare   │                              │ │
│  │  │ (Next.js)    │  │ Tunnel       │                              │ │
│  │  │ Port: 3000   │  │              │                              │ │
│  │  └──────────────┘  └──────────────┘                              │ │
│  │                                                                    │ │
│  └────────────────────────────────────────────────────────────────────┘ │
│                                                                          │
│  EDGE (Homestead - Raspberry Pi / Intel NUC)                            │
│  ┌────────────────────────────────────────────────────────────────────┐ │
│  │                                                                    │ │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐            │ │
│  │  │ LiveKit      │  │ OpenMind     │  │ Home         │            │ │
│  │  │ Agent Worker │  │ OM1 Runtime  │  │ Assistant    │            │ │
│  │  └──────────────┘  └──────────────┘  └──────────────┘            │ │
│  │                                                                    │ │
│  │  ┌──────────────┐  ┌──────────────┐                              │ │
│  │  │ NATS Leaf    │  │ Local        │                              │ │
│  │  │ Node         │  │ Sensors      │                              │ │
│  │  └──────────────┘  └──────────────┘                              │ │
│  │                                                                    │ │
│  └────────────────────────────────────────────────────────────────────┘ │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

### 6.2 Docker Compose (Cloud)

```yaml
# docker/cloud/docker-compose.yml (additions)
services:
  # Agno Runtime (replacing Letta)
  agno-runtime:
    build:
      context: ../../
      dockerfile: docker/Dockerfile.agno
    ports:
      - "7777:7777"
    environment:
      - AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}
      - AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}
      - AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION}
      - DATABASE_URL=postgresql://onemind:${DB_PASSWORD}@postgres:5432/onemind
      - NATS_URL=nats://nats:4222
    depends_on:
      - postgres
      - nats
    networks:
      - onemind

  # LiveKit Server
  livekit:
    image: livekit/livekit-server:latest
    ports:
      - "7880:7880"
      - "7881:7881"
    environment:
      - LIVEKIT_KEYS=${LIVEKIT_API_KEY}:${LIVEKIT_API_SECRET}
    networks:
      - onemind

  # LiveKit Agent Worker
  livekit-agent:
    build:
      context: ../../
      dockerfile: docker/Dockerfile.livekit-agent
    environment:
      - LIVEKIT_URL=ws://livekit:7880
      - LIVEKIT_API_KEY=${LIVEKIT_API_KEY}
      - LIVEKIT_API_SECRET=${LIVEKIT_API_SECRET}
      - DEEPGRAM_API_KEY=${DEEPGRAM_API_KEY}
      - CARTESIA_API_KEY=${CARTESIA_API_KEY}
      - AGNO_RUNTIME_URL=http://agno-runtime:7777
    depends_on:
      - livekit
      - agno-runtime
    networks:
      - onemind

  # Agent UI (Self-hosted)
  agent-ui:
    build:
      context: ../../
      dockerfile: docker/Dockerfile.agent-ui
    ports:
      - "3000:3000"
    environment:
      - AGENT_OS_URL=http://agno-runtime:7777
    depends_on:
      - agno-runtime
    networks:
      - onemind
```

---

## Part 7: Grants & Funding Opportunities

### 7.1 OpenMind Developer Program

**Status**: Research required - no $100K grant confirmed yet

OpenMind has:

- $20M raised (Pantera, Coinbase Ventures, DCG)
- Active hackathons (Alzheimer's care robot built in 48 hours)
- FABRIC Network waitlist
- Open source OM1 on GitHub

**Action Items**:

- [ ] Join FABRIC waitlist
- [ ] Contact OpenMind about developer grants
- [ ] Explore hackathon participation
- [ ] Investigate partnership opportunities

### 7.2 Other Robotics Grants

| Program                   | Amount        | Focus                | Status   |
| ------------------------- | ------------- | -------------------- | -------- |
| NVIDIA Inception          | Varies        | AI/Robotics startups | Apply    |
| AWS Activate              | $100K credits | Startups             | Apply    |
| Google Cloud for Startups | $200K credits | AI/ML                | Apply    |
| Figure AI Partnership     | Unknown       | Humanoid integration | Research |
| Boston Dynamics SDK       | Unknown       | Spot integration     | Research |

---

## Part 8: Implementation Phases

### Phase 1: Agno Migration (Week 1-2)

- [ ] Complete Legacy agent definition in Agno
- [ ] Set up PostgreSQL with pgvector
- [ ] Deploy AgentOS runtime
- [ ] Deploy Agent UI (self-hosted)
- [ ] Test basic conversation + memory
- [ ] Remove Letta containers

### Phase 2: LiveKit Voice (Week 3-4)

- [ ] Deploy LiveKit server (self-hosted)
- [ ] Configure Deepgram STT plugin
- [ ] Configure Cartesia TTS plugin
- [ ] Build LiveKit agent wrapper for Agno
- [ ] Test voice conversation
- [ ] Measure latency (<500ms target)

### Phase 3: Multi-Site Voice (Week 5-6)

- [ ] Set up LiveKit rooms per location
- [ ] Configure room handoff logic
- [ ] Deploy edge workers at homestead
- [ ] Test continuous conversation across locations
- [ ] Integrate motion sensors for automatic handoff

### Phase 4: OpenMind Integration (Week 7-8)

- [ ] Set up OM1 on test robot (quadruped/simulator)
- [ ] Create Agno tools for robot control
- [ ] Bridge FABRIC events to NATS
- [ ] Test Legacy → robot command flow
- [ ] Add robot camera streams to LiveKit

### Phase 5: Drone Fleet (Week 9-10)

- [ ] Configure drone telemetry → NATS
- [ ] Create drone patrol tools
- [ ] Integrate drone video with LiveKit
- [ ] Test autonomous patrol with visual reporting
- [ ] Add drone data to Legacy's awareness

### Phase 6: Full Integration (Week 11-12)

- [ ] End-to-end testing of all systems
- [ ] Performance optimization
- [ ] Documentation
- [ ] IaC templates for rapid deployment
- [ ] Production hardening

---

## Part 9: Success Metrics

### 9.1 Technical Metrics

| Metric        | Target        | Measurement                |
| ------------- | ------------- | -------------------------- |
| Voice latency | <500ms        | End-to-end speech response |
| Event latency | <1ms          | NATS message delivery      |
| Memory recall | >95% accuracy | Semantic search relevance  |
| Uptime        | 99.9%         | Service availability       |
| Agent init    | <10ms         | Agno instantiation         |

### 9.2 Experience Metrics

| Metric                  | Target          | Measurement                      |
| ----------------------- | --------------- | -------------------------------- |
| Conversation continuity | Seamless        | Cross-location context retention |
| Proactive actions       | 5+/day          | Unprompted useful actions        |
| Pattern recognition     | Weekly insights | TimescaleDB correlations         |
| Robot responsiveness    | <2s             | Command to action                |

---

## Part 10: The Vision Realized

```
┌─────────────────────────────────────────────────────────────────────────┐
│                                                                          │
│                          O N E M I N D                                   │
│                                                                          │
│         "Everything seen. Everything noticed. One mind."                 │
│                                                                          │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  HUMAN                                                                   │
│    Zeus at the center                                                    │
│    Voice everywhere                                                      │
│    Continuous conversation                                               │
│    Always understood                                                     │
│                                                                          │
│  NATURE                                                                  │
│    FarmBot tends the greenhouse                                          │
│    Drones monitor plant health                                           │
│    Sensors track the ecosystem                                           │
│    AI optimizes growth                                                   │
│                                                                          │
│  TECHNOLOGY                                                                  │
│    Legacy orchestrates all                                               │
│    Robots as extensions of will                                          │
│    Memory that never forgets                                             │
│    Intelligence that evolves                                             │
│                                                                          │
│  CONNECTED                                                               │
│    NATS binds all systems                                                │
│    Sub-millisecond awareness                                             │
│    Shared world model                                                    │
│    Collective intelligence                                               │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Appendix A: Key Links

- **Agno**: https://github.com/agno-agi/agno
- **Agent UI**: https://github.com/agno-agi/agent-ui
- **LiveKit Agents**: https://github.com/livekit/agents
- **OpenMind OM1**: https://github.com/openmind-agi/om1
- **FABRIC Network**: https://openmind.org/ (waitlist)
- **NATS**: https://nats.io/

## Appendix B: Credentials Location

All in 1Password vault: `00-24 SO (Sovereign Operations)`

---

_This document is the north star for OneMind OS development._
_Updated: 2026-01-17_

# Integrating OpenMind OM1 with Agno-based OneMindOS

**OpenMind OM1 and Agno serve complementary roles and integrate best when Agno remains the primary orchestrator while OM1 provides perception-action capabilities as a callable subsystem.** Neither framework natively supports NATS — OM1 uses Zenoh/ROS2 for robot messaging and Agno uses HTTP/SSE via FastAPI — so the proposed NATS "nervous system" would be a custom bridge layer added on top of both. The strongest digital-first strategy wraps OM1's hosted multi-agent API as an Agno tool today, introduces NATS only when physical hardware demands real-time pub/sub, and keeps your OMCS team structure entirely within Agno where it already works well.

This report presents three concrete integration patterns ranked by complexity, with a recommended 60-day implementation roadmap that starts digital and scales to physical robotics.

---

## What OM1 and Agno actually do (and don't overlap on)

**OM1 is a robotics AI runtime**, not a general-purpose agent framework. Its architecture follows a tight **Sense → Fuse → Cortex → Act** loop running at configurable frequency (default 0.5 Hz). Raw sensor data (camera, LiDAR, microphone, GPS) flows through AI captioning (VLM, ASR), gets fused with governance rules and user intent, passes to an LLM "cortex" for decision-making, and outputs structured action commands (move, speak, emote). OM1's hosted API at `api.openmind.org/api/core/agent` runs **four parallel agents** — Navigation, Perception, RAG, and Team — that synthesize into a unified response in roughly **1.4 seconds**.

**Agno is a general-purpose multi-agent orchestration framework** built on FastAPI. It organizes agents into Teams with three coordination modes (Route, Coordinate, Collaborate), persists state to PostgreSQL, and exposes agents via A2A (Agent-to-Agent), MCP (Model Context Protocol), and AG-UI protocols. With **100+ built-in tool integrations**, 20+ LLM providers, and 11 storage backends, Agno handles digital operations — calendar, email, research, task management — far better than OM1.

**The critical insight: OM1 does not replace any part of Agno for digital operations.** OM1 adds perception-action capabilities that Agno lacks entirely. Your four OMCS teams (UI, GE, HP, LE) should remain Agno-native. OM1 enters as a specialized capability layer for when agents need to perceive physical environments or control hardware.

---

## Pattern A: Agno-primary with OM1 as a callable tool

This is the simplest pattern and the recommended starting point. Agno's AgentOS remains the sole orchestrator. OM1's hosted API is wrapped as a custom Agno tool that any team can invoke when robotics-aware reasoning is needed.

**Architecture:**

```
Flutter Frontend
       │ (HTTP/SSE via AG-UI)
       ▼
┌─────────────────────────────────┐
│       Agno AgentOS (Railway)    │
│  ┌─────┐ ┌─────┐ ┌────┐ ┌────┐│
│  │ UI  │ │ GE  │ │ HP │ │ LE ││  ← OMCS Teams
│  │Team │ │Team │ │Team│ │Team││
│  └──┬──┘ └──┬──┘ └─┬──┘ └─┬──┘│
│     └───────┴──────┴──────┘    │
│              │ (tool call)     │
│     ┌────────▼─────────┐      │
│     │  OM1 Bridge Tool │      │
│     └────────┬─────────┘      │
│              │                 │
│  ┌───────────▼──────────┐     │
│  │   PostgreSQL+pgvector│     │
│  └──────────────────────┘     │
└──────────────┬────────────────┘
               │ HTTPS
               ▼
      api.openmind.org
      POST /api/core/agent
```

**Implementation** requires a single custom tool function in Python:

```python
import requests
from agno.tools import tool

@tool()
def om1_multi_agent(
    system_prompt: str,
    inputs: str,
    available_actions: str = "speak, analyze, navigate, observe",
    model: str = "gpt-4.1-mini"
) -> str:
    """Call OpenMind OM1's multi-agent system for
    perception-aware and robotics-aware reasoning.
    Returns synthesized output from Navigation,
    Perception, RAG, and Team agents."""
    response = requests.post(
        "https://api.openmind.org/api/core/agent",
        json={
            "system_prompt": system_prompt,
            "inputs": inputs,
            "available_actions": available_actions,
            "model": model,
            "structured_outputs": True,
        },
        headers={
            "Authorization": f"Bearer {OM1_API_KEY}",
            "Content-Type": "application/json",
        },
        timeout=15,
    )
    data = response.json()
    return data.get("content", str(data))
```

Any OMCS team agent that needs OM1 capabilities simply includes `om1_multi_agent` in its tools list. The LE (Life Coordination) team might call it when a physical robot task is queued; the HP (Health) team might use it for sensor-based health monitoring context.

**Pros:** Deployable in hours, not days. No new infrastructure. Agno remains the single source of truth for orchestration. OM1's hosted API handles scaling. Works entirely on Railway with the existing Agno template. **Estimated cost: $15–30/month** on Railway Pro (Agno + PostgreSQL) plus OM1 API credits (starting at $1, consumption-based).

**Cons:** Synchronous HTTP calls add **~1.5 seconds latency** per OM1 invocation. No real-time sensor streaming. Limited to what the hosted API exposes. Cannot run custom OM1 perception pipelines. Single point of failure on OM1's API availability.

---

## Pattern B: NATS message bus bridging both frameworks

This pattern introduces NATS as a shared "nervous system" connecting Agno and OM1 through hierarchical pub/sub topics. Neither framework natively supports NATS, so both require thin adapter layers — but NATS provides the real-time, many-to-many messaging essential for physical robotics.

**Architecture:**

```
Flutter Frontend
       │
       ▼
┌──────────────────────────────────────────────────┐
│                  Railway Cloud                     │
│                                                    │
│  ┌────────────┐   ┌──────────────────────────┐   │
│  │ NATS Server│   │     Agno AgentOS         │   │
│  │ +JetStream │◄──┤  ┌────┐┌────┐┌───┐┌───┐ │   │
│  │ +KV Store  │   │  │ UI ││ GE ││HP ││LE │ │   │
│  └─────┬──────┘   │  └────┘└────┘└───┘└───┘ │   │
│        │          │  + NATS Adapter (nats-py) │   │
│        │          └──────────────────────────┘   │
│        │                                          │
│  ┌─────▼──────────────────────┐                  │
│  │  PostgreSQL + pgvector     │                  │
│  └────────────────────────────┘                  │
└────────┬─────────────────────────────────────────┘
         │ NATS (WebSocket tunnel)
         ▼
┌────────────────────────────┐
│  OM1 Runtime (Edge/VPS)    │
│  + NATS Adapter            │
│  + Zenoh/ROS2 (for robots) │
│  + Unitree SDK             │
└────────────────────────────┘
```

**Subject namespace design** maps cleanly to both OM1's Sense/Cortex/Act pipeline and Agno's OMCS teams:

```
# OM1 perception → Agno consumption
om1.sense.vision.{robot_id}     # VLM-captioned camera data
om1.sense.audio.{robot_id}      # ASR transcriptions
om1.sense.spatial.{robot_id}    # LiDAR/SLAM data
om1.cortex.decision.{robot_id}  # LLM action decisions
om1.act.status.{robot_id}       # Action execution status

# Agno → OM1 commands
agno.command.move.{robot_id}    # Movement instructions
agno.command.speak.{robot_id}   # Speech output
agno.command.task.{robot_id}    # Task assignments

# OMCS team coordination
omcs.ui.event.*                 # Smart home events
omcs.ge.task.*                  # Business operations
omcs.hp.vitals.*                # Health data streams
omcs.le.schedule.*              # Life coordination

# Shared state
state.context.{session_id}      # Cross-framework context
bridge.request.{correlation_id} # Request/reply pattern
```

**The NATS adapter for Agno** is a lightweight Python service using `nats-py` that subscribes to relevant topics and converts messages into Agno tool calls or agent triggers. For OM1, the adapter sits between the OM1 runtime's WebSocket output and the NATS server, translating action commands and sensor captions into NATS messages.

**JetStream** stores all messages durably, enabling audit trails, replay for debugging, and event sourcing. The **KV Store** (built into JetStream) provides shared state accessible from both frameworks — useful for maintaining a unified context like "robot is currently in the kitchen" that both Agno's LE team and OM1's navigation agent need.

**Pros:** Loosely coupled — either framework can be replaced independently. Real-time streaming for sensor data. Many-to-many communication enables monitoring, logging, and multi-robot coordination. **NATS deploys on Railway with a one-click template** including JetStream and WebSocket support. JetStream KV eliminates need for additional shared state infrastructure.

**Cons:** Requires custom adapter code for both frameworks (~500–800 lines Python each). More moving parts to monitor. **NATS adds $3–8/month** on Railway. OM1 runtime still needs separate hosting (edge device or VPS) since Railway lacks GPU instances for local LLM inference. Schema governance across topics requires discipline.

---

## Pattern C: A2A protocol bridge with OM1 wrapper

This pattern leverages Agno's first-class A2A (Agent-to-Agent) protocol support to create a standards-compliant bridge. A thin A2A-compliant wrapper exposes OM1's capabilities as a discoverable remote agent that Agno can call natively.

**Architecture:**

```
┌───────────────────────────────┐
│       Agno AgentOS            │
│  OMCS Teams + RemoteAgent     │
│  (OM1 registered as A2A peer) │
└──────────┬────────────────────┘
           │ A2A Protocol (HTTP + SSE)
           ▼
┌───────────────────────────────┐
│  OM1 A2A Wrapper Service      │
│  ┌─────────────────────────┐  │
│  │ /.well-known/agent-card │  │  ← Discovery
│  │ /v1/message:send        │  │  ← Sync calls
│  │ /v1/message:stream      │  │  ← Streaming
│  └────────┬────────────────┘  │
│           │ HTTP               │
│  ┌────────▼────────────────┐  │
│  │  OM1 Runtime or API     │  │
│  └─────────────────────────┘  │
└───────────────────────────────┘
```

**The OM1 A2A wrapper** is a FastAPI service (~200 lines) that exposes OM1's capabilities through the A2A protocol standard. Agno discovers it via the Agent Card, then communicates using `RemoteAgent`:

```python
from agno.agent import RemoteAgent

om1_agent = RemoteAgent(
    base_url="http://om1-wrapper:8001",
    agent_id="om1-perception",
    protocol="a2a",
)

# Use as a team member in any OMCS team
le_team = Team(
    name="LE_LifeCoordination",
    mode="coordinate",
    members=[calendar_agent, task_agent, om1_agent],
)
```

**Pros:** Uses established protocol standards — future-proof as A2A adoption grows. OM1 appears as a native team member within Agno, enabling the Team Leader to delegate to it naturally. Supports streaming via SSE. The Agent Card enables automated discovery and capability negotiation.

**Cons:** Requires building the A2A wrapper service. A2A is still relatively new (Google-backed, but ecosystem is maturing). Adds a network hop. Less suitable for high-frequency sensor streaming than NATS. The OM1 wrapper must translate between A2A's task-based model and OM1's continuous sense-act loop.

---

## Why NATS isn't needed for digital-only deployment

A persistent misconception in the project context is that OM1 uses NATS with topics like `openmind.sense.*`. **This is incorrect.** OM1 uses Zenoh, CycloneDDS, and ROS2 for robot-to-runtime messaging. The `openmind.*` topic namespace was never part of OM1's architecture.

For **digital-only deployment** (Phase 1), NATS adds unnecessary complexity. Agno's HTTP/SSE communication handles all digital agent coordination. OM1's hosted API is reachable via simple HTTP calls. The overhead of deploying, configuring, and maintaining NATS plus two custom adapters is not justified when your agents are making API calls to calendar services, databases, and LLMs.

**NATS becomes essential in Phase 2** when physical robots generate continuous sensor streams at high frequency. A Unitree Go2's cameras, LiDAR, and IMU produce data that needs real-time pub/sub distribution to multiple consumers. NATS's **sub-millisecond latency**, wildcard subscriptions, and JetStream persistence are purpose-built for this workload. The transition point is when you need many-to-many, real-time, bidirectional data flow — which is precisely what physical robotics demands.

---

## The 60-day implementation roadmap

**Days 1–14: Digital foundation (Pattern A)**

Deploy Agno AgentOS on Railway using the official template. Stand up your four OMCS teams with PostgreSQL persistence. Implement the `om1_multi_agent` tool function and integrate it into the LE team. Register for an OM1 API key at portal.openmind.org. Validate round-trip latency and response quality from the multi-agent endpoint. Build the Flutter frontend connection via AG-UI protocol (Dart SDK available). **Deliverable:** All four OMCS teams operational with OM1 as a callable tool.

**Days 15–30: Intelligence layer deepening**

Add OM1 RAG integration by uploading relevant knowledge documents to portal.openmind.org. Configure Agno's own knowledge system (pgvector) for digital-domain knowledge (calendar patterns, health data, business documents). Implement memory and learning across Agno agents so they improve with use. Build the A2A wrapper for OM1 (Pattern C) to enable cleaner inter-framework communication. Test OM1's WebSim simulator for prototyping physical interactions digitally. **Deliverable:** Intelligent agents with persistent memory, RAG knowledge, and A2A-based OM1 integration.

**Days 31–45: NATS bridge and physical preparation**

Deploy NATS with JetStream on Railway using the one-click template. Build the Agno-NATS adapter (~500 lines Python using `nats-py`). Build the OM1-NATS adapter for translating Zenoh/ROS2 messages to NATS topics. Define and document the full subject namespace. Set up JetStream streams for audit logging. Procure Unitree robot hardware. Configure OM1 runtime on edge compute (Raspberry Pi 5 with 16GB or Jetson). **Deliverable:** NATS infrastructure operational, adapters tested with simulated sensor data.

**Days 46–60: Physical integration and testing**

Flash OM1 onto robot hardware using the BrainPack configuration. Connect robot's Zenoh topics to NATS via the adapter. Verify sensor data flows from robot → NATS → Agno agents in real time. Test command flow from Agno LE team → NATS → OM1 → robot movement. Integrate Brilliant Labs Frame AR glasses (likely via Bluetooth/WebSocket to a custom input plugin in OM1). Run end-to-end scenarios: "Schedule a meeting, then have the robot prepare the conference room." **Deliverable:** Full digital-physical loop operational.

---

## Infrastructure costs and resource requirements

| Component | Hosting | Monthly Cost | Resources |
|-----------|---------|-------------|-----------|
| Agno AgentOS | Railway Pro | **$15–25** | 2 vCPU, 4 GB RAM |
| PostgreSQL + pgvector | Railway Pro | **$5–10** | 1 vCPU, 2 GB RAM |
| NATS + JetStream | Railway Pro | **$3–8** | 0.5 vCPU, 512 MB RAM |
| OM1 hosted API | api.openmind.org | **$5–20** | Consumption-based |
| OM1 runtime (Phase 2) | Edge device | **$0** (one-time hardware) | RPi 5 16GB or Jetson |
| LLM API costs | Various providers | **$10–50** | Usage-dependent |
| **Total (Digital Phase)** | | **$38–113/month** | |

Railway's Pro plan at **$20/month** includes $20 in credits, which covers most of the infrastructure. The OM1 API's credit consumption depends on call frequency — at the default 0.5 Hz loop rate, credits drain fast, so use on-demand calls from Agno tools rather than continuous polling.

---

## Conclusion: start simple, add complexity only when physics demands it

The recommended path is unambiguous: **Pattern A first, evolving to Pattern B when robots arrive.** Agno should remain the sole orchestrator for your OMCS architecture — it handles team coordination, persistence, memory, and digital operations far better than OM1. OpenMind OM1 is not a competitor to Agno; it is a specialized perception-action runtime that fills the exact gap Agno has around physical-world interaction.

The NATS "nervous system" is the right architectural choice for Phase 2, but deploying it prematurely in the digital-only phase creates operational overhead without proportional benefit. The A2A wrapper (Pattern C) offers a clean middle ground for the transition period — it gives OM1 a standards-compliant interface within Agno's team structure without requiring NATS infrastructure.

Three decisions that will most impact success: first, **resist the urge to merge the frameworks** — keep Agno and OM1 as separate services with clear API boundaries. Second, **design your NATS subject namespace now** even if you don't deploy NATS until Day 31, because it forces you to define exactly what data flows between digital and physical domains. Third, **use OM1's hosted API during development** and only self-host the runtime when you need custom perception pipelines or ultra-low-latency robot control — the hosted multi-agent endpoint at ~1.4 seconds per call is fast enough for most digital-first scenarios.

---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]

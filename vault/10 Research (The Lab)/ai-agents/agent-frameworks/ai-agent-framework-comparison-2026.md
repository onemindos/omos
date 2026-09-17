---
title: "AI Agent Framework Comparison — April 2026"
type: reference
status: active
created: 2026-04-06
owner: zeus
tags: [type/reference, domain/so, topic/agents, topic/research]
ai_summary: "Comprehensive comparison of 17 AI agent frameworks for choosing a long-term agent system for the OneMind life operating system. Covers LangGraph, AutoGen/AG2, CrewAI, OpenAI Agents SDK, Semantic Kernel, Pydantic AI, Google ADK, smolagents, Agno, Letta, Haystack, MCP, OpenClaw, Hermes Agent, and 3 protocols (MCP, A2A, AG-UI). Includes analysis for hardware/robotics, home automation, and personal AI use cases."
---

# AI Agent Framework Comparison — April 2026

> Research for selecting the long-term agent infrastructure for One Mind OS.
> Requirement: forkable, model-agnostic, multi-agent, hardware-capable, local-deployable, protocol-compliant (MCP + A2A + AG-UI), and suitable for a life operating system spanning software + hardware + robotics + home automation.

---

## The Protocol Triangle (Critical Context)

Before comparing frameworks, understand the three emerging standards that will define the agentic ecosystem:

| Protocol | What It Does | Who Created It | Status (April 2026) |
|----------|-------------|----------------|---------------------|
| **MCP** (Model Context Protocol) | Agent ↔ Tools/Data. How agents access external systems. | Anthropic (open-sourced) | De facto standard. Widely adopted. |
| **A2A** (Agent-to-Agent) | Agent ↔ Agent. How agents from different vendors collaborate. | Google (open-sourced) | v1.0 stable. Microsoft, SAP, Zoom adopting. |
| **AG-UI** (Agent-User Interaction) | Agent ↔ Human UI. How agents communicate with frontends in real-time. | CopilotKit (open-sourced) | Growing rapidly. AWS AgentCore support. |

**Together, MCP + A2A + AG-UI = "The Protocol Triangle"** — the full communication stack for agentic systems. Any framework you choose should support all three.

---

## Master Comparison Table

| Framework | License | Language | Stars | Multi-Agent | Model Agnostic | MCP | A2A | Memory | Edge/Local | Forkable | Best For |
|-----------|---------|----------|-------|-------------|----------------|-----|-----|--------|------------|----------|----------|
| **LangGraph** | MIT | Python, JS | 28.5K | ✅ Graph-based | ✅ Any LLM | ✅ | Via adapters | ✅ Checkpoints | ⚠️ Possible | ✅ | Production stateful workflows |
| **CrewAI** | MIT | Python | 47.8K | ✅ Role-based crews | ✅ Any LLM | ✅ | ✅ Native | ✅ Short+long | ⚠️ Possible | ✅ | Rapid multi-agent development |
| **AG2 (AutoGen)** | Apache 2.0 | Python | ~35K | ✅ Conversation-based | ✅ Any LLM | ✅ | ✅ AG-UI | ✅ Memory protocol | ⚠️ Possible | ✅ | Research & experimentation |
| **MS Agent Framework** | MIT | Python, .NET | ~22K* | ✅ SK + AutoGen fused | ✅ Any LLM | ✅ | ✅ | ✅ SK Memory | ⚠️ Azure-heavy | ✅ | Enterprise at scale |
| **Semantic Kernel** | MIT | .NET, Python, Java | ~22K | ✅ Agent abstraction | ✅ Any LLM | ✅ | ✅ Via AF | ✅ Enterprise-grade | ⚠️ Possible | ✅ | Enterprise multi-lang |
| **OpenAI Agents SDK** | MIT | Python, TS | ~15K | ✅ Handoffs | ⚠️ Provider-agnostic paths | ✅ | ❌ | ⚠️ Basic | ⚠️ Possible | ✅ | OpenAI ecosystem |
| **Google ADK** | Apache 2.0 | Python, Java | ~12K | ✅ Sequential/Parallel | ⚠️ Gemini-optimized | ✅ | ✅ Native | ✅ Session service | ⚠️ Vertex-heavy | ✅ | Google/A2A ecosystem |
| **Pydantic AI** | MIT | Python | ~10K | ⚠️ Single-agent focus | ✅ Any LLM | ⚠️ Via tools | ❌ | ⚠️ Basic | ✅ Lightweight | ✅ | Type-safe production agents |
| **smolagents** | Apache 2.0 | Python | ~15K | ⚠️ Limited | ✅ Any LLM (HF Hub) | ⚠️ | ❌ | ⚠️ Basic | ✅ Very lightweight | ✅ | Minimal agents, education |
| **Agno** | MPL-2.0 | Python | ~20K | ✅ Teams + Workflows | ✅ Any LLM | ✅ | ⚠️ | ✅ Built-in | ✅ Self-hosted | ⚠️ MPL | Multi-modal agent platform |
| **Letta** | Apache 2.0 | Python | ~14K | ⚠️ Single-agent + API | ✅ Any LLM | ⚠️ | ❌ | ✅✅ Best-in-class | ✅ Local runtime | ✅ | Stateful memory agents |
| **Haystack** | Apache 2.0 | Python | ~18K | ⚠️ Pipeline-based | ✅ Any LLM | ⚠️ | ❌ | ⚠️ Via components | ✅ K8s-ready | ✅ | RAG + agentic pipelines |
| **OpenClaw** | MIT | TypeScript/Node | 100K+ | ✅ Multi-agent + skills | ✅ Any LLM | ✅ | ⚠️ | ✅ Persistent | ✅ Mac/VPS/Edge | ✅ | Personal AI agent platform |
| **Hermes Agent** | MIT | Python | ~8K | ⚠️ Single-agent | ✅ Via Hermes models | ⚠️ | ❌ | ✅ 3-tier self-improving | ✅ Local | ✅ | Self-improving personal AI |

*Stars are approximate as of April 2026. Some repos merged or had sub-repos.*

---

## Detailed Framework Profiles

### 1. LangGraph (LangChain)

- **License:** MIT
- **Language:** Python (primary), JavaScript/TypeScript
- **GitHub:** 28.5K stars, 4.9K forks, 491 releases
- **Architecture:** Directed cyclic graphs for stateful agent workflows. Explicit state machines with nodes and edges.
- **Multi-agent:** Yes — via graph composition. Supports Sequential, Concurrent, Group Chat patterns.
- **Model Agnostic:** Yes — works with any LLM provider via LangChain.
- **Memory:** Checkpoint-based persistence. LangGraph Platform adds durable state.
- **MCP Support:** Yes — native via LangChain tools.
- **A2A Support:** Via adapters (not native).
- **AG-UI:** Yes — documented integration with FastAPI + AG-UI streaming.
- **Edge/Local:** Possible but designed for cloud. LangGraph Platform is a managed service.
- **Hardware/Robotics:** No native support. You'd wrap hardware APIs as tools.
- **Home Automation:** No native integration. Tool-based approach.
- **Strengths:** Most mature for production. Fine-grained control. Huge ecosystem. Excellent observability (LangSmith).
- **Weaknesses:** Complex. Steep learning curve. LangGraph Platform (managed) is paid. Abstractions can be leaky.
- **Forkable:** Yes (MIT). But ecosystem dependency on LangChain is deep.
- **Best For:** Complex stateful workflows where you need precise control over agent behavior.

### 2. CrewAI

- **License:** MIT
- **Language:** Python
- **GitHub:** 47.8K stars, 6.5K forks. 27M PyPI downloads.
- **Architecture:** Role-based agent orchestration. Define agents with Role, Goal, Backstory. Three modes: Sequential, Hierarchical (manager delegates), and Flows (event-driven).
- **Multi-agent:** Excellent. This is CrewAI's core strength. "Crews" of specialized agents collaborate naturally.
- **Model Agnostic:** Yes — any LLM provider.
- **Memory:** Short-term (conversation), long-term (across sessions), entity memory, user memory.
- **MCP Support:** Yes.
- **A2A Support:** Yes — native A2A API since March 2026. Can collaborate with AutoGen, LangGraph, MCP agents.
- **AG-UI:** Via streaming results (v2026).
- **Edge/Local:** Possible but not optimized for it.
- **Hardware/Robotics:** No native support. Tool-based approach.
- **Home Automation:** No native integration.
- **Enterprise:** Fortune 500 adoption (nearly half in 2024). $18M funding. 2B agentic executions in 12 months.
- **Strengths:** Fastest to prototype. Natural-language agent definitions. Huge community. Enterprise-grade. A2A native.
- **Weaknesses:** Less fine-grained control than LangGraph. Enterprise features behind paid tiers ($25/mo+). Framework can feel "magical" (hard to debug).
- **Forkable:** Yes (MIT). Can be fully self-hosted.
- **Best For:** Rapid multi-agent systems where role-based collaboration is the pattern.

### 3. AG2 (formerly AutoGen, Microsoft)

- **License:** Apache 2.0 (community fork). Original AutoGen is under Microsoft's umbrella.
- **Language:** Python
- **GitHub:** ~35K stars (combined AutoGen + AG2)
- **Architecture:** Conversation-first. Agents are `ConversableAgent` instances that communicate via structured dialogue. Human-in-the-loop via `UserProxyAgent`.
- **Multi-agent:** Excellent. GroupChat, Swarm, GraphFlow patterns. Dynamic conversation-driven workflows.
- **Model Agnostic:** Yes — any LLM.
- **Memory:** Memory Protocol — flexible interface integrating with vector stores.
- **MCP Support:** Yes.
- **A2A Support:** Yes — native AG-UI integration (`AGUIStream`).
- **AG-UI:** Yes — first-class examples in build-with-ag2 repo.
- **Edge/Local:** Possible. Lighter than LangGraph.
- **Hardware/Robotics:** No native support.
- **Home Automation:** No native integration.
- **Strengths:** Best for exploratory multi-agent patterns. Very flexible. Good research community. AG-UI native.
- **Weaknesses:** AG2 fork vs Microsoft's AutoGen creates ecosystem confusion. Less production-hardened than LangGraph. Documentation fragmented.
- **Forkable:** Yes (Apache 2.0). AG2 is already the community fork.
- **Best For:** Research, experimentation, conversation-driven multi-agent systems.

### 4. Microsoft Agent Framework (October 2025)

- **License:** MIT (open-source SDK)
- **Language:** Python, .NET, Java
- **Architecture:** Fuses Semantic Kernel (cognition/planning/memory) with AutoGen (multi-agent orchestration) into one cohesive system. SK = the brain, AutoGen = the nervous system.
- **Multi-agent:** Yes — inherits AutoGen's full orchestration + SK's planning.
- **Model Agnostic:** Yes — any LLM via SK connectors.
- **Memory:** Semantic Kernel Memory (enterprise-grade, multi-modal RAG).
- **MCP Support:** Yes.
- **A2A Support:** Yes — Microsoft announced A2A support in Azure AI Foundry.
- **Edge/Local:** Possible but Azure-optimized.
- **Hardware/Robotics:** No native support. But .NET makes Windows IoT integration feasible.
- **Home Automation:** No native integration.
- **Strengths:** Enterprise-grade. Multi-language. Combines the best of two mature frameworks. Azure integration. Security/compliance built-in.
- **Weaknesses:** Heavy Azure ecosystem dependency. Over-engineered for personal use. Early in unification — rough edges.
- **Forkable:** Yes (MIT). But Azure dependencies make it harder to fully self-host.
- **Best For:** Enterprise-scale agent systems needing .NET + Python + Java + Azure.

### 5. Semantic Kernel

- **License:** MIT
- **Language:** .NET (primary), Python, Java
- **GitHub:** ~22K stars
- **Architecture:** Plugin-based SDK. The Kernel orchestrates Plugins (your code) and Functions (prompts). Explicit patterns: Sequential, Concurrent, Group Chat.
- **Multi-agent:** Yes — first-class agent abstraction since Python 1.6.0 / .NET 1.18.0.
- **Model Agnostic:** Yes — OpenAI, Azure OpenAI, Hugging Face, any LLM.
- **Memory:** Kernel Memory — dedicated, multi-modal service for enterprise-grade RAG.
- **MCP Support:** Yes.
- **A2A Support:** Yes (via Microsoft Agent Framework).
- **Edge/Local:** Possible. Lighter than full Agent Framework.
- **Strengths:** Multi-language. Enterprise memory. Planning capabilities. Mature.
- **Weaknesses:** .NET-first bias. Less community than Python-only frameworks.
- **Best For:** Enterprise multi-language teams needing a robust SDK.

### 6. OpenAI Agents SDK (formerly Swarm)

- **License:** MIT
- **Language:** Python, TypeScript
- **GitHub:** ~15K stars
- **Architecture:** Simple primitives: Agent, Handoff, Guardrails, Tracing. Production-ready evolution of Swarm's prototype approach. Built-in tracing for debugging.
- **Multi-agent:** Yes — via handoff functions. Agent A can hand off to Agent B seamlessly.
- **Model Agnostic:** Technically yes — documented paths for non-OpenAI models. But optimized for OpenAI.
- **Memory:** Basic — session-level. No built-in persistent memory.
- **MCP Support:** Yes — OpenAI's Apps SDK extends MCP with UI capabilities.
- **A2A Support:** No native support.
- **Edge/Local:** Possible but not designed for it.
- **Strengths:** Simple. Clean API. Production-ready. Excellent tracing. Apps SDK extends MCP.
- **Weaknesses:** Shallow memory. OpenAI-centric despite agnostic claims. No A2A. Limited multi-agent patterns.
- **Forkable:** Yes (MIT).
- **Best For:** Simple agent systems in the OpenAI ecosystem.

### 7. Google ADK (Agent Development Kit)

- **License:** Apache 2.0
- **Language:** Python (v1.0.0 stable), Java
- **GitHub:** ~12K stars
- **Architecture:** Agent-based with SequentialAgent, ParallelAgent composition. Native A2A protocol support. Built for Google Vertex AI / Agent Engine.
- **Multi-agent:** Yes — sequential, parallel, custom orchestration. A2A enables cross-framework collaboration.
- **Model Agnostic:** Gemini-optimized but supports other models.
- **Memory:** Session service with persistent state.
- **MCP Support:** Yes — native integration (Google Codelab: "Getting Started with MCP, ADK and A2A").
- **A2A Support:** Yes — **this is Google's framework for A2A**. Native, first-class.
- **Edge/Local:** Possible but Vertex AI is the primary deployment target.
- **Strengths:** A2A native. MCP native. Production-ready (v1.0.0). Google ecosystem.
- **Weaknesses:** Gemini-centric. Vertex AI dependency for production. Newer, smaller community.
- **Forkable:** Yes (Apache 2.0).
- **Best For:** A2A-first agent systems. Google Cloud teams.

### 8. Pydantic AI

- **License:** MIT
- **Language:** Python
- **GitHub:** ~10K stars
- **Architecture:** Type-safe agent definitions using Pydantic schemas. Agent = high-level construct defined by data schemas and Python functions. New: Temporal integration for durable execution.
- **Multi-agent:** Limited — primarily single-agent focus. Can compose via Temporal workflows.
- **Model Agnostic:** Yes — supports Anthropic, OpenAI, Gemini, Groq, etc.
- **Memory:** Basic — state via Pydantic models.
- **MCP Support:** Via tool definitions.
- **A2A Support:** No.
- **Edge/Local:** Yes — very lightweight. Good for edge.
- **Strengths:** Type safety. Clean API. Durable execution via Temporal. Fast to build.
- **Weaknesses:** Not a multi-agent framework. Limited orchestration. Young ecosystem.
- **Forkable:** Yes (MIT).
- **Best For:** Type-safe production agents where reliability matters more than multi-agent complexity.

### 9. smolagents (Hugging Face)

- **License:** Apache 2.0
- **Language:** Python
- **GitHub:** ~15K stars
- **Architecture:** Minimalist. "Code agents" write and execute Python code as actions (vs JSON tool calls). Under 1,000 lines of source code. HF Hub integrated.
- **Multi-agent:** Limited — primary focus is single agent with tools.
- **Model Agnostic:** Yes — any model on HF Hub + external providers.
- **Memory:** Basic — conversation context.
- **MCP Support:** Limited.
- **A2A Support:** No.
- **Edge/Local:** Excellent — designed to be lightweight. Good for edge/embedded.
- **Strengths:** Minimal. Fast. Code-first approach. HF ecosystem. Under 20 lines to build an agent.
- **Weaknesses:** Not enterprise-grade. Limited orchestration. No multi-agent. No memory persistence.
- **Forkable:** Yes (Apache 2.0).
- **Best For:** Lightweight agents, education, embedded/edge use cases.

### 10. Agno (formerly PhiData)

- **License:** MPL-2.0 (Mozilla Public License)
- **Language:** Python
- **GitHub:** ~20K stars
- **Architecture:** Full-stack agent platform. Agents → Teams → Workflows hierarchy. Built-in knowledge, memory, tools, guardrails, human-in-the-loop. Multi-modal (Image, Audio, Video, File).
- **Multi-agent:** Yes — Teams of specialized agents + Workflows for orchestration.
- **Model Agnostic:** Yes — any LLM provider.
- **Memory:** Built-in persistent memory. Self-learning agents.
- **MCP Support:** Yes.
- **A2A Support:** Limited.
- **Edge/Local:** Yes — "Agent OS" runs in your cloud. Self-hosted.
- **Hardware/Robotics:** Multi-modal capabilities (vision, audio) are relevant.
- **Strengths:** Full-stack. Multi-modal. Self-hosted. Fast development. Agent OS concept.
- **Weaknesses:** MPL-2.0 (not as permissive as MIT/Apache). Commercial "Agent OS" platform.
- **Forkable:** MPL-2.0 allows forking but requires sharing modifications to MPL files.
- **Best For:** Multi-modal agent systems where you want a full platform.

### 11. Letta (formerly MemGPT)

- **License:** Apache 2.0
- **Language:** Python
- **GitHub:** ~14K stars
- **Architecture:** **LLM-as-an-Operating-System paradigm.** The model actively manages its own memory, context, and reasoning. Three-tier memory: core (always available), recall (searchable conversation history), archival (unlimited long-term storage with embeddings).
- **Multi-agent:** Limited — focuses on individual stateful agents. API for multi-agent coordination.
- **Model Agnostic:** Yes — any LLM.
- **Memory:** **Best-in-class.** This IS Letta's core innovation. Active memory management, not passive.
- **MCP Support:** Limited.
- **A2A Support:** No.
- **Edge/Local:** Yes — "Letta Code" runs agents locally in terminal. Letta API for server deployment.
- **Strengths:** Memory architecture is unmatched. Self-improving over time. OS-inspired design. Evaluation framework.
- **Weaknesses:** Not a multi-agent orchestration framework. Limited tool ecosystem. Smaller community.
- **Forkable:** Yes (Apache 2.0).
- **Best For:** Building agents that truly learn and remember. Long-term personal AI.

### 12. Haystack (deepset)

- **License:** Apache 2.0
- **Language:** Python
- **GitHub:** ~18K stars
- **Architecture:** Pipeline-based. Components flow through directed pipelines. Serializable, cloud-agnostic, K8s-ready.
- **Multi-agent:** Pipeline-based (not traditional multi-agent). Tool calling + agentic RAG.
- **Model Agnostic:** Yes — any LLM, any vector DB.
- **Memory:** Via pipeline components.
- **MCP Support:** Limited.
- **A2A Support:** No.
- **Edge/Local:** Yes — K8s-ready. Self-hosted.
- **Strengths:** Best for RAG. Production-ready pipelines. Modular. Good observability.
- **Weaknesses:** Not designed for complex agentic workflows beyond RAG. Pipeline model limits flexibility.
- **Forkable:** Yes (Apache 2.0).
- **Best For:** RAG-heavy applications. Document processing. Knowledge retrieval.

### 13. OpenClaw

- **License:** MIT
- **Language:** TypeScript/Node.js
- **GitHub:** 100K+ stars (org), 18.6K followers. Surpassed 100K stars.
- **Architecture:** Personal AI agent platform. Multi-agent with skills. Messaging gateway (Telegram, Discord, Slack, WhatsApp). Web UI (ClawHub). Self-hosted.
- **Multi-agent:** Yes — multiple specialized agents with skill system.
- **Model Agnostic:** Yes — any LLM provider (OpenAI, Anthropic, Google, local via LiteLLM).
- **Memory:** Persistent memory with ClawVault. Conversation + long-term.
- **MCP Support:** Yes — core integration.
- **A2A Support:** Limited.
- **AG-UI:** Not native.
- **Edge/Local:** Yes — runs on Mac Mini, VPS, Raspberry Pi. Designed for local-first.
- **Hardware/Robotics:** No native integration but extensible via skills.
- **Home Automation:** No native HA integration but messaging gateways enable it.
- **Strengths:** Personal AI pioneer. Huge community. Skill marketplace. Multi-platform messaging. Local-first. Active development.
- **Weaknesses:** TypeScript (not Python). Supply chain risk (LiteLLM hack March 2026). Single-maintainer risk for some components.
- **Forkable:** Yes (MIT). You already use it.
- **Best For:** Personal AI agent that runs across devices and platforms.

### 14. Hermes Agent (Nous Research)

- **License:** MIT
- **Language:** Python
- **GitHub:** ~8K stars (launched Feb 2026)
- **Architecture:** Self-improving AI agent with ReAct loop. Three-tier memory: session, persistent, and **skill memory** (learned solution patterns). 40+ bundled skills. Messaging gateway. One-command migration from OpenClaw.
- **Multi-agent:** **Single-agent only** (dealbreaker for some). No multi-agent orchestration.
- **Model Agnostic:** Yes — via Hermes models + LiteLLM.
- **Memory:** Excellent 3-tier system with self-improvement loop. Agent learns from experience.
- **MCP Support:** Limited.
- **A2A Support:** No.
- **Edge/Local:** Yes — runs locally. Very efficient.
- **Strengths:** Self-improving. Skill memory. Clean architecture. Serverless backends ($5 VPS). Nous Research credibility.
- **Weaknesses:** **Single-agent only.** New (Feb 2026). Supply chain risk (same LiteLLM dependency). Smaller community.
- **Forkable:** Yes (MIT).
- **Best For:** Single personal AI agent that improves over time.

---

## Protocol Support Matrix

| Framework | MCP (Tools) | A2A (Agent-Agent) | AG-UI (Agent-Human) |
|-----------|:-----------:|:------------------:|:-------------------:|
| LangGraph | ✅ | ⚠️ Adapter | ✅ Via FastAPI |
| CrewAI | ✅ | ✅ Native (Mar 2026) | ✅ Streaming |
| AG2 (AutoGen) | ✅ | ✅ Via AG-UI | ✅ Native |
| MS Agent Framework | ✅ | ✅ Azure AI Foundry | ✅ |
| Semantic Kernel | ✅ | ✅ Via AF | ✅ |
| OpenAI Agents SDK | ✅ | ❌ | ⚠️ Via Apps SDK |
| Google ADK | ✅ | ✅✅ Native (creator) | ⚠️ |
| Pydantic AI | ⚠️ | ❌ | ❌ |
| smolagents | ⚠️ | ❌ | ❌ |
| Agno | ✅ | ⚠️ | ⚠️ |
| Letta | ⚠️ | ❌ | ❌ |
| Haystack | ⚠️ | ❌ | ❌ |
| OpenClaw | ✅ | ⚠️ | ❌ |
| Hermes Agent | ⚠️ | ❌ | ❌ |

**Winner for protocol compliance:** CrewAI, AG2, Microsoft Agent Framework, Google ADK.

---

## Special Topics

### Home Assistant + AI Agents

Home Assistant has evolved significantly for AI integration:
- **AI Tasks** — opt-in feature for leveraging LLMs to create automations
- **OpenRouter integration** — unified interface to 400+ models directly in HA
- **Local AI** — can run local LLMs for generating automations with entity awareness
- **Assist** — built-in voice assistant with LLM-powered conversations
- **Voice Chapter 10** — highly optimized voice pipeline (ESPHome, ATOM Echo)

**No agent framework natively integrates with HA.** The integration pattern is:
1. Agent framework → MCP server for HA → HA REST API
2. Or: Agent → MQTT → HA
3. Or: Agent → Home Assistant WebSocket API

**For OneMind:** Build a HA MCP server that exposes entities, automations, and scenes as tools. Any MCP-capable agent framework can then control your home.

### Robotics Integration

No mainstream agent framework has native robotics support. The pattern is:
- **ROS 2 (Robot Operating System)** — the standard for robotics
- **Agent Framework → MCP/Tool → ROS 2 bridge** — agent sends commands
- **Google ADK** — closest to robotics via Gemini's multi-modal + Google's robotics research
- **smolagents** — lightweight enough for edge devices (Raspberry Pi, Jetson)
- **Pydantic AI** — lightweight enough for embedded systems

For drones, smart glasses, food growing systems: the agent orchestrates but doesn't directly control. It sends commands via well-defined interfaces (MQTT, REST, ROS 2 topics).

### Edge Device Deployment

Best frameworks for Raspberry Pi / Jetson / embedded:
1. **smolagents** — <1K lines of code, minimal dependencies
2. **Pydantic AI** — lightweight, type-safe
3. **Letta Code** — local runtime, agents in terminal
4. **Hermes Agent** — designed for local execution
5. **OpenClaw** — Node.js runs on Pi

Worst for edge: LangGraph Platform, Google ADK (Vertex), Microsoft Agent Framework (Azure).

### Personal AI / Life Management

No framework is specifically designed for personal AI / life management. The closest:
1. **OpenClaw** — personal AI pioneer, messaging gateway, skills
2. **Hermes Agent** — self-improving personal assistant
3. **Letta** — memory that makes it feel like a persistent companion

**This is the gap OneMind fills** — none of these are a comprehensive life operating system. They're tools, not frameworks for organizing an entire life.

---

## Analysis: What OneMind Actually Needs

### Requirements Ranked by Priority

1. **Multi-agent orchestration** — 4+ domain agents (SO, HP, LE, GE) with coordination
2. **Model agnostic** — must work with OpenAI, Anthropic, Google, local (Ollama)
3. **Protocol compliance** — MCP (tools), A2A (agent-agent), AG-UI (human interface)
4. **Persistent memory** — agents must remember across sessions
5. **Forkable (MIT/Apache)** — must be ownable, rebrandable
6. **Local/edge deployment** — VPS, Mac, Raspberry Pi, offline capable
7. **Hardware extensibility** — tools/MCP for HA, robotics, IoT
8. **Active development** — not a dead project by 2027
9. **Python-first** — largest AI ecosystem, most integrations
10. **Manageable complexity** — one person (Zeus) maintains it

### The Honest Assessment

**No single framework does everything OneMind needs.** The correct approach is a **layered architecture**:

```
┌─────────────────────────────────────────────────┐
│  AG-UI Protocol (Agent ↔ Human Interface)       │
│  → OneMind Web / OneMind Mobile                 │
├─────────────────────────────────────────────────┤
│  Agent Orchestration Layer (Multi-Agent)         │
│  → CrewAI or LangGraph (the "brains")           │
├─────────────────────────────────────────────────┤
│  A2A Protocol (Agent ↔ Agent)                   │
│  → Cross-framework agent communication          │
├─────────────────────────────────────────────────┤
│  Memory Layer                                   │
│  → Letta-style memory OR custom (pgvector)      │
├─────────────────────────────────────────────────┤
│  MCP Protocol (Agent ↔ Tools/Data)              │
│  → Home Assistant, Todoist, GitHub, Telegram...  │
├─────────────────────────────────────────────────┤
│  Edge Runtime (Local/Offline)                   │
│  → Ollama + lightweight agent on Pi/Jetson      │
└─────────────────────────────────────────────────┘
```

---

## Recommendation: The OneMind Agent Stack

### Primary Orchestration: **LangGraph** (MIT)

**Why LangGraph over CrewAI:**
- LangGraph gives you explicit, deterministic control over agent workflows (graphs > magic)
- Better for a system Zeus maintains long-term — you see exactly what's happening
- MIT license, fully forkable
- 28.5K stars, 491 releases, massive ecosystem
- CrewAI's "natural language agent definitions" are great for prototyping but harder to debug in production
- LangGraph's graph model maps naturally to OneMind's domain structure (SO/HP/LE/GE as separate graph nodes)

**Close second: CrewAI** (MIT) — if you prefer faster development over fine-grained control, or if you want A2A native.

### Memory Layer: **Letta architecture** (Apache 2.0) or **Custom**

- Adopt Letta's 3-tier memory architecture (core/recall/archival) but implement it on top of your existing pgvector
- You already have PostgreSQL + pgvector on VPS
- Don't add Letta as a dependency — adopt the pattern

### Protocol Layer: **MCP + A2A + AG-UI**

- **MCP servers** for every external system (HA, Todoist, GitHub, Telegram, 1Password, etc.)
- **A2A** for domain agents to communicate (SO agent talks to HP agent)
- **AG-UI** for OneMind Web / OneMind Mobile real-time updates

### Edge Runtime: **Ollama + smolagents or Pydantic AI**

- For Raspberry Pi / Jetson / offline: lightweight agent with local LLM
- smolagents for simple tasks, Pydantic AI for type-safe control

### Personal AI Layer: **Keep OpenClaw** (for now)

- OpenClaw already works for your messaging gateway (Telegram, Discord)
- Replace over time with LangGraph-based agents that use MCP for messaging
- Don't rip out what works — evolve it

### What NOT to Choose

- **Microsoft Agent Framework / Semantic Kernel** — too Azure-dependent, over-engineered for personal use
- **Google ADK** — too Vertex/Gemini-dependent
- **OpenAI Agents SDK** — too OpenAI-centric, shallow memory
- **Agno** — MPL-2.0 license complicates full ownership
- **Hermes Agent** — single-agent only, too new

---

## Migration Path

### Phase 1: Protocol First (Now)
- Build MCP servers for your key integrations (done for some via skills)
- Evaluate AG-UI for OneMind Web

### Phase 2: LangGraph Prototype (Month 1)
- Build a simple multi-agent graph: SO (orchestrator) coordinating HP, LE, GE agents
- Test with existing tools via MCP
- Use LangGraph checkpointing for state

### Phase 3: Memory Architecture (Month 2)
- Implement Letta-style 3-tier memory on pgvector
- Core memory: identity, beliefs, preferences (from `_codex/system/`)
- Recall memory: conversation history
- Archival memory: full vault search via embeddings

### Phase 4: Replace OpenClaw (Month 3-6)
- Migrate messaging gateway to LangGraph-based agents
- Keep OpenClaw running in parallel during migration
- A2A bridge between old (OpenClaw) and new (LangGraph) systems

### Phase 5: Edge + Hardware (Month 6+)
- Deploy lightweight agents on Raspberry Pi for home automation
- MCP server for Home Assistant
- Drone / growing / glasses integration via tool definitions

---

## Key Takeaways

1. **The protocol triangle (MCP + A2A + AG-UI) matters more than any single framework.** Build protocol-first, framework-second.

2. **LangGraph is the safest long-term bet** for orchestration — MIT, mature, explicit control, massive ecosystem.

3. **Memory is the killer feature for personal AI.** Adopt Letta's architecture pattern on your existing pgvector.

4. **No framework handles hardware natively.** The answer is always: MCP server wrapping hardware APIs.

5. **OpenClaw is a good personal AI tool but not a long-term foundation** for a life operating system. It's a ready-made product, not a programmable framework.

6. **CrewAI is the dark horse** — if A2A native support and rapid development matter more than graph-level control, it's a strong choice.

7. **The landscape will keep changing.** The protocol layer (MCP, A2A, AG-UI) is more stable than any framework. Build on protocols, wrap frameworks.

---

*Research conducted April 6, 2026. Sources: GitHub repositories, official documentation, web searches across Medium, LinkedIn, HackerNoon, DataCamp, and framework-specific blogs.*


---

## Related

- [[00-24 SO (Sovereign Operations)/06 Research (R&D) (The Lab)/06 Research|06 Research]]

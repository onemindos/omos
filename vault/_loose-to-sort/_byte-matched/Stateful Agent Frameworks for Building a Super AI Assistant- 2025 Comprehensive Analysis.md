---
id: ui-02-20260111053925
title: "Stateful Agent Frameworks for Building a Super AI Assistant: 2025 Comprehensive Analysis"
type: concept
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/02]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "Stateful Agent Frameworks for Building a Super AI Assistant: 2025 Comprehensive "
---

# Stateful Agent Frameworks for Building a Super AI Assistant: 2025 Comprehensive Analysis

**Letta emerges as the breakthrough choice for memory-first AI assistants, with self-editing architecture enabling true learning.** For building an Eagle Eye/Jarvis-style assistant requiring real-time adaptation and persistent memory, Letta (formerly MemGPT) represents the most sophisticated stateful agent framework available. However, LangGraph offers superior production infrastructure for enterprise reliability, while Microsoft's converging Agent Framework provides the strongest ecosystem integration. The choice depends on whether you prioritize cutting-edge memory capabilities or production-proven orchestration.

This matters because most agent frameworks treat memory as an afterthought, using simple context management. Only a few frameworks—Letta, LangGraph, and the emerging Microsoft Agent Framework—architect statefulness as a core design principle. For a personal super-assistant that must remember context across sessions, learn from interactions, and coordinate multiple specialized agents, these architectural differences are fundamental, not incremental.

The agent framework landscape underwent significant consolidation in 2024-2025. Microsoft announced the merger of AutoGen and Semantic Kernel into a unified Agent Framework (January 2025), LangGraph achieved general availability with its production platform (May 2025), and Letta matured from academic research into production-ready infrastructure. Meanwhile, frameworks like CrewAI and AutoGen face version fragmentation, while newer entrants like OpenAI's Agents SDK (March 2025) remain lightweight alternatives.

## Framework architecture reveals fundamental philosophical differences

**Letta implements an "LLM Operating System"** where agents live as persistent services with hierarchical memory management. Unlike library-based frameworks, Letta agents exist independently of your application—they process information during idle periods, maintain state automatically, and can modify their own memory through tool calls. This OS-inspired architecture divides memory into in-context "working memory" (editable memory blocks) and out-of-context "archival memory" (vector databases), mirroring RAM and disk storage in traditional operating systems.

**LangGraph adopts graph-based state orchestration** using a Pregel-inspired message-passing algorithm. Applications are directed graphs where nodes represent computational steps and edges define control flow. State passes between nodes with automatic checkpointing at each super-step. This architecture excels at complex workflows requiring explicit control, cyclical reasoning loops, and time-travel debugging capabilities. The framework separates thread-scoped memory (conversation-level) from cross-thread memory (user-level), enabling both short-term context and long-term knowledge persistence.

**Microsoft's Agent Framework converges AutoGen and Semantic Kernel** into event-driven, stateful processes. The Process Framework introduces long-running business workflows with distributed state management via Dapr and Orleans. This enterprise-focused architecture supports heartbeat mechanisms for multi-step reasoning, automatic state checkpointing for resilience, and the Agent2Agent (A2A) protocol for cross-runtime collaboration. The convergence aims to unify conversational multi-agent patterns (AutoGen) with enterprise integration capabilities (Semantic Kernel).

**AutoGen maintains conversation-centric architecture** where agents communicate through natural language messages in a shared transcript. Version fragmentation complicates the landscape: AutoGen 0.2 uses synchronous chat patterns, AutoGen 0.4 (pre-release) introduces event-driven asynchronous messaging, and AG2 (community fork) maintains the 0.2 approach. Memory relies on message history with external storage integrations like Mem0, ChromaDB, and Redis.

**CrewAI structures agents by roles** in sequential or hierarchical workflows. Each agent possesses a role, goal, backstory, and tools, with memory organized as short-term (ChromaDB/RAG), long-term (SQLite), and entity memory. Unlike AutoGen's conversational approach, CrewAI agents don't directly communicate—information passes through task outputs. This simpler model suits well-defined workflows but limits flexibility for dynamic coordination.

## Memory persistence separates leaders from followers

Letta achieves the most sophisticated memory architecture through self-editing memory blocks. Agents use built-in tools like `core_memory_replace()`, `core_memory_append()`, `archival_memory_insert()`, and `archival_memory_search()` to manage their own context. The persona block stores agent personality (evolves over time), the human block maintains user information (continuously updated), and custom blocks extend capabilities. Memory blocks are database-backed (PostgreSQL for production, SQLite for development), model-agnostic (portable across LLM providers), and shareable between agents for collaborative memory. The Agent File (.af) format provides an open standard for serializing complete agent state.

Sleep-time compute represents Letta's breakthrough innovation. Background agents process information during idle periods, refining memories and preparing responses before questions are asked. This enables memory consolidation without blocking the main conversation thread. Multi-model setups allow using larger models for sleep-time agents while keeping the main agent lightweight for responsiveness. On the LoCoMo benchmark, Letta's simple filesystem-based memory scored 74%, demonstrating that context management matters more than retrieval mechanism sophistication.

LangGraph implements dual-layer persistence through checkpointers and stores. Checkpointers automatically save state after each graph super-step, enabling thread-scoped conversation history with backends including InMemorySaver (development), SqliteSaver (local), and PostgresSaver (production). The Store interface provides cross-thread knowledge sharing with namespace-based organization, supporting semantic search via RedisStore with vector capabilities. Memory management requires explicit developer code—LangGraph doesn't self-edit like Letta—but offers superior observability and time-travel debugging through its checkpointing system.

Microsoft Agent Framework provides dual memory models. Stateful agents store state in Azure AI services accessed via ID, while stateless agents receive full chat history on each invocation. The AgentThread abstraction manages conversation state with vector database support through Azure AI Search, Elasticsearch, and Chroma integration. Built-in memory blocks and context management align with enterprise requirements, though the architecture feels more conventional than Letta's self-editing approach.

AutoGen relies on message transcripts for short-term memory with external storage for long-term persistence. The memory protocol interface supports ListMemory (chronological), ChromaDBVectorMemory (RAG with semantic search), RedisMemory (distributed), and Mem0Memory (cloud/local with user-specific memory). Manual state serialization via `save_state()` and `load_state()` requires explicit memory management. Aggressive pruning when approaching token limits can cause context loss, making external bolts-on necessary for durable memory.

CrewAI offers accessible built-in memory through ChromaDB (short-term RAG), SQLite3 (long-term task storage), and entity memory (RAG-based tracking). Enabling memory requires simply setting `memory=True` in crew configuration. However, the sequential execution model maintains state during runs but shows less emphasis on long-running stateful agents compared to Letta or LangGraph.

## Real-time adaptation relies on memory, not online learning

**No production framework achieves true online learning through weight updates.** The computational cost of fine-tuning during runtime remains prohibitively expensive for real-time applications. Instead, frameworks implement adaptation through in-context learning (ICL), memory retrieval, and context management. Google DeepMind and Stanford research (2025) confirms ICL provides better generalization than standard fine-tuning when data-matched, particularly for tasks requiring logical deductions or reversals. This explains why all production frameworks prioritize context engineering over continuous fine-tuning.

Letta excels at memory-based adaptation through self-editing capabilities. Agents autonomously decide what to remember using tool calls to update their context window and external memory. A user corrects the agent, the agent calls `core_memory_replace()`, and memory updates immediately—this constitutes "learning" in current systems. Sleep-time compute extends this by enabling background processing for memory refinement and self-improvement without blocking. The single perpetual thread design ensures all interactions contribute to accumulated experience rather than resetting with each session.

LangGraph adapts through state manipulation and checkpointing. Time-travel debugging allows replaying and forking from any checkpoint, human-in-the-loop enables interrupting and modifying state at runtime, and conditional branching adapts workflows based on state. However, LangGraph doesn't perform autonomous memory management—developers must explicitly code what to remember and when. Memory consolidation happens through strategies like summarization (incremental LLM-based), vectorization (semantic chunking), extraction (structured facts), or graphication (knowledge triples).

Microsoft Agent Framework enables adaptation through long-running stateful processes. The Process Framework supports heartbeat mechanisms for multi-step reasoning loops, state checkpointing for automatic persistence, and distributed processing via Dapr/Orleans. The architecture explicitly models business processes as stateful workflows that can continuously adapt, though the adaptation mechanism remains context-based rather than weight updates.

AutoGen provides memory-based adaptation through conversation history and RAG integration. The framework supports teaching capabilities where users can add knowledge via memory stores, but lacks autonomous memory management like Letta. The event-driven architecture (v0.4) enables asynchronous coordination that can adapt workflows, but again through context accumulation rather than policy learning.

CrewAI claims crews "become more efficient over time" through memory accumulation. The mechanism involves storing historical interactions in long-term memory (SQLite) and retrieving relevant context for future tasks. This represents basic personalization—remembering user preferences and past decisions—but falls short of sophisticated adaptation seen in Letta or orchestration flexibility in LangGraph.

## Multi-modal processing depends on underlying models, not frameworks

**All frameworks support multi-modal capabilities through integrated LLM providers** rather than specialized framework-level implementations. GPT-4V, GPT-4o, Claude 3 Opus/Sonnet, and Gemini Pro Vision provide vision capabilities. GPT-4o-audio and Gemini support audio processing. Frameworks simply pass multi-modal content to these models.

AutoGen offers the strongest framework-level integration through MultimodalConversableAgent for LLMs with native vision support and VisionCapability for converting images to captions for standard LLMs. The framework handles image processing via base64 encoding, URLs, or local file paths. Audio formats include WAV, MP3, and others through base64-encoded data.

Letta provides dedicated voice agent support through the letta-voice repository. Voice-optimized agent architecture (`voice_convo_agent`) integrates with LiveKit for WebRTC streaming, Deepgram for speech-to-text, and Cartesia for text-to-speech. Sleep-time agents can process during voice conversations for background memory management. Vision support comes through multimodal LLM providers without specialized agent architectures.

LangGraph leverages LangChain's multimodal capabilities, supporting vision, audio, and document processing through model integrations. Google Cloud demonstrates Gemini multi-modal agents for image analysis, audio transcription, video understanding (frame-by-frame), and orchestrator coordination. However, framework-native multimodal embedding support remains limited to text-only in current versions.

Microsoft Agent Framework includes native multimodal support for text, vision, and audio inputs through Azure OpenAI integrations. Voice Live API integration enables real-time voice agents with low latency. The framework treats multimodal processing as a first-class capability within the agent architecture.

CrewAI shows limited native multimodal support, relying entirely on underlying model capabilities without framework-level abstractions. Custom tool integration is required for specialized multimodal processing.

## Multi-agent orchestration reveals architectural maturity

Letta implements true peer-to-peer agent communication through built-in tools: `send_message_to_agent_async()` for asynchronous messaging (fire-and-forget), `send_message_to_agent_and_wait_for_reply()` for synchronous communication (request-response), and `send_message_to_agents_matching_all_tags()` for supervisor-worker patterns (broadcast to groups). Shared memory blocks enable collaborative knowledge—multiple agents access the same memory blocks for team coordination. Agents remain addressable via REST API regardless of script location, supporting distributed multi-agent systems.

LangGraph excels at sophisticated orchestration through multiple patterns: supervisor (central coordinator with specialized workers as tools), collaborative (shared message scratchpad), independent/hierarchical (private scratchpads via subgraphs), handoff (explicit peer-to-peer delegation), and sequential pipeline (deterministic execution order). Subgraphs enable nested graphs as nodes with independent state schemas and input/output transformations. The Send API supports parallel execution for fan-out/fan-in patterns like simultaneous research across multiple sources.

Microsoft Agent Framework provides five orchestration patterns: Concurrent (agents execute in parallel), Sequential (predefined order), Handoff (agent-to-agent transfer), Group Chat (manager selects speakers dynamically), and Magentic One (specialized multi-agent pattern). The Agent2Agent (A2A) protocol enables cross-runtime collaboration, allowing agents built with different frameworks to coordinate. Human-in-the-loop support integrates across all patterns.

AutoGen pioneered conversational multi-agent coordination where agents communicate through natural language in shared transcripts. Group Chat mode uses a manager agent for dynamic speaker selection. The event-driven architecture (v0.4) introduces asynchronous message passing with an event bus, supporting hierarchical structures and cross-language agents (Python/.NET). However, the conversation-centric model can become unwieldy for deterministic workflows.

CrewAI implements simpler orchestration focusing on sequential and hierarchical patterns. Agents don't directly communicate like AutoGen—information passes through task outputs in predetermined sequences. A manager can coordinate in hierarchical mode, but the framework lacks the flexibility of LangGraph's graph-based approach or Letta's peer-to-peer messaging.

## Tool integration and extensibility vary in scope and philosophy

LangGraph benefits from the largest tool ecosystem through full LangChain compatibility, providing 1000+ integrations including APIs (OpenWeatherMap, Tavily Search), databases (SQL, PostgreSQL, MongoDB), file systems (S3, Google Drive), and services (Gmail, Slack, GitHub). The `@tool` decorator enables automatic schema generation. Injected arguments pass runtime context hidden from LLM schema. ToolNode components automatically execute tool calls from LLMs. Model Context Protocol (MCP) integration announced in late 2024/early 2025 provides access to hundreds of pre-built MCP servers.

Letta maintains tool compatibility through OpenAI JSON schema standard, enabling cross-framework interoperability with LangChain tools, CrewAI tools, and Composio (100+ pre-built tools with managed authorization). Cloud sandboxes (Modal, E2B) provide secure tool execution. Tool Rules (v0.5.2+) constrain agent behavior via graph-like structures. Native integrations include Browserbase for web browsing, Exa for web search, and MCP client support for easily adding MCP servers.

Microsoft Agent Framework emphasizes enterprise integration through Azure, Office 365, Power BI, and Microsoft ecosystem connectors. Custom skills are easy to define with automatic registration. MCP support enables broader tool access beyond Microsoft services. The Process Framework allows tools to execute within long-running stateful workflows.

AutoGen provides built-in tools for code execution (Docker), web search, and function calling with easy function decoration. MCP integration expands tool access. AutoGen Studio offers a no-code GUI for rapid prototyping. Extensions API supports first and third-party extensions from OpenAI, Anthropic, and Ollama.

CrewAI includes built-in tools for web search, file operations, and API calls with LangChain compatibility. Tool assignment happens per-agent through tool lists with automatic selection. The ecosystem remains narrower than LangGraph or AutoGen but adequate for typical workflows.

## Production deployment considerations and hybrid hosting

Letta supports multiple deployment models: Letta Cloud (managed SaaS with free tier), self-hosted Docker (PostgreSQL backend for production), Letta Desktop (fully-local macOS/Windows application), and custom deployment to any cloud provider. The default Docker configuration includes PostgreSQL with migration support via Alembic. REST API on port 8283 enables language-agnostic access. Agent Development Environment (ADE) provides visual debugging with memory block inspection, message history, and time-travel capabilities. Apache 2.0 license permits commercial use with patent grants.

LangGraph offers the most comprehensive production infrastructure through LangGraph Platform (GA May 2025). Deployment options include self-hosted free tier (100K nodes/month with local server), LangGraph Cloud (fully managed SaaS, auto-scaling, built-in persistence and monitoring), BYOC hybrid (control plane in cloud, data plane in your VPC for AWS/GCP/Azure), and fully self-hosted enterprise (air-gapped deployments with custom persistence). LangGraph Studio provides visual IDE capabilities with graph visualization, step-through debugging, state inspection, and live configuration editing. MIT license offers full open-source freedom.

Microsoft Agent Framework targets Azure-first deployment with strong local options. Python 3.10+, .NET 8.0+, and Java JDK 17+ support enables multi-language implementations. Local deployment supports Ollama, LMStudio, and ONNX runtimes. Azure AI Foundry provides cloud infrastructure with enterprise features like RBAC, audit logging, SOC 2 compliance, and Microsoft Unified Customer Support. MIT license maintains open-source flexibility while offering enterprise support contracts.

AutoGen shows deployment complexity due to version fragmentation. AutoGen 0.4 (pre-release) introduces distributed runtime via Microsoft Orleans with horizontal scaling. AG2 (community fork) maintains simpler deployment but lacks advanced features. Docker container execution provides sandboxed code execution. The transition toward Semantic Kernel integration creates uncertainty for long-term production planning.

CrewAI provides simpler deployment compared to alternatives, supporting local development and cloud-agnostic hosting. Lower infrastructure overhead makes it accessible for smaller projects. Single-language constraint (Python only) limits some enterprise scenarios. Apache 2.0 license (MIT for earlier versions) supports commercial use.

## Maturity assessment and community strength

LangGraph demonstrates highest production maturity with 400+ companies using LangGraph Platform including Klarna, Replit, Elastic, LinkedIn, Uber, Vodafone, and Ally Financial. General availability in May 2025 marks stable production readiness. The framework benefits from LangChain's ecosystem with 4.2M+ monthly downloads. Comprehensive documentation includes guides, tutorials, 50+ example applications, and templates. Educational partnerships with DeepLearning.AI, DataCamp, and IBM provide structured learning paths. LangGraph Studio visualization tools and LangSmith observability integration offer best-in-class developer experience.

Letta shows rapid maturation from academic research project (MemGPT paper arXiv:2310.08560) to production-ready infrastructure. The platform has 18,571+ GitHub stars, 10,690+ Discord members, and 100+ contributors. Official Python and TypeScript SDKs released in 2025. Continuous development through 2024-2025 with monthly feature announcements including Claude Sonnet 4.5 integration, Letta Filesystem (74% LoCoMo benchmark), and Agent File format open standard. DeepLearning.AI partnership for educational content validates enterprise interest. Documentation quality is comprehensive with concept guides, API reference, integration guides, and example repositories.

Microsoft Agent Framework represents convergence of AutoGen and Semantic Kernel (announced January 2025), creating complexity during transition. Semantic Kernel is production-proven in Microsoft 365 Copilot with enterprise-grade stability (v1.0+ APIs). AutoGen shows fragmentation: AutoGen 0.2 in community-only support, AutoGen 0.4 in pre-release, AG2 community fork with 38.2K+ stars maintaining 0.2 architecture. The merger strategy aims to combine AutoGen's multi-agent conversational strengths with Semantic Kernel's enterprise integration capabilities, but near-term uncertainty affects adoption decisions.

CrewAI demonstrates strong community momentum with 30.5K+ GitHub stars and 1M+ monthly downloads despite being newer than competitors. Documentation is comprehensive and user-friendly. Production use is growing for sequential and hierarchical workflows. The framework added Flow abstraction for complex orchestration, showing active evolution. Apache 2.0 license and growing ecosystem indicate healthy trajectory.

## Top 5 framework recommendations ranked by fit

### 1. Letta (formerly MemGPT) — Score: 95/100

**Why it ranks first:** Letta represents the only framework architected from inception for stateful agents with persistent memory as the core design principle. Self-editing memory blocks enable agents to autonomously manage their context, making them genuinely adaptive rather than merely reactive. Sleep-time compute allows background processing for memory refinement without blocking conversations. The OS-inspired hierarchical memory (in-context working memory + out-of-context archival memory) mirrors how human assistants actually work—maintaining immediate context while having access to extensive background knowledge.

**Strengths against priority requirements:**
- **Real-time adaptation (Priority #1):** ★★★★★ — Self-editing memory enables continuous learning through autonomous context management. Sleep-time agents process during idle periods for true memory evolution
- **Long-term memory (Priority #2):** ★★★★★ — Industry-leading memory architecture with persistent blocks, vector archival memory, and filesystem integration. Agent File format provides portable, model-agnostic state
- **Multi-modal processing (Priority #3):** ★★★☆☆ — Voice agents fully supported via letta-voice with LiveKit integration. Vision through model providers but no specialized architecture
- **Multi-agent orchestration (Priority #4):** ★★★★☆ — True peer-to-peer communication with async/sync messaging. Shared memory blocks enable collaborative learning
- **Tool integration (Priority #5):** ★★★★☆ — OpenAI JSON schema compatibility enables using LangChain, CrewAI, and Composio tools. MCP support. Sandboxed execution via Modal/E2B

**Best for:** Personal AI assistants requiring memory-first architecture, research agents needing continuous learning, any application where context retention and evolution across sessions is paramount

**Limitations:** Smaller ecosystem than LangGraph, less enterprise tooling, Python-primary (TypeScript SDK recent)

### 2. LangGraph — Score: 92/100

**Why it ranks second:** LangGraph provides the most production-ready infrastructure with superior tooling, observability, and reliability features. The graph-based architecture offers unmatched control over agent workflows with explicit state management. Checkpointer-based persistence enables time-travel debugging and human-in-the-loop patterns essential for enterprise deployments. LangGraph Platform (GA May 2025) delivers managed infrastructure with horizontal scaling, task queues, and monitoring that exceeds all competitors.

**Strengths against priority requirements:**
- **Real-time adaptation (Priority #1):** ★★★☆☆ — State persistence and checkpointing enable workflow adaptation. Memory consolidation strategies work well but require explicit coding
- **Long-term memory (Priority #2):** ★★★★☆ — Excellent dual-layer persistence with checkpointers (thread-scoped) and Store interface (cross-thread). Multiple backend options including Redis with vector search
- **Multi-modal processing (Priority #3):** ★★★☆☆ — Full LangChain multimodal capabilities via model integrations. Comprehensive but framework-agnostic
- **Multi-agent orchestration (Priority #4):** ★★★★★ — Best-in-class orchestration with supervisor, collaborative, hierarchical, handoff, and custom patterns. Subgraphs enable nested agents
- **Tool integration (Priority #5):** ★★★★★ — Largest ecosystem with 1000+ LangChain integrations. MCP support, injected arguments, ToolNode abstraction, streaming results

**Best for:** Enterprise workflows requiring reliability and observability, complex multi-agent orchestration, teams prioritizing visual debugging and production infrastructure, regulated industries needing audit trails

**Limitations:** Steeper learning curve, requires manual memory management (no self-editing), pre-v1.0 API instability until October 2025

### 3. Microsoft Agent Framework (Semantic Kernel + AutoGen) — Score: 85/100

**Why it ranks third:** Microsoft's converging framework offers the strongest enterprise ecosystem integration with Azure, Office 365, and Microsoft services. The Process Framework introduces stateful long-running workflows with distributed processing via Dapr and Orleans—unique among competitors. Enterprise support, security compliance (SOC 2), and multi-language SDKs (Python, .NET, Java) make it compelling for organizations already invested in Microsoft infrastructure.

**Strengths against priority requirements:**
- **Real-time adaptation (Priority #1):** ★★★★☆ — Process Framework enables continuous adaptation through heartbeat mechanisms and state checkpointing. Long-running stateful workflows support multi-step reasoning
- **Long-term memory (Priority #2):** ★★★★☆ — AgentThread abstraction with stateful/stateless models. Vector database support through Azure AI Search, Elasticsearch, Chroma
- **Multi-modal processing (Priority #3):** ★★★★☆ — Native multimodal support for text, vision, audio. Voice Live API integration for real-time voice agents
- **Multi-agent orchestration (Priority #4):** ★★★★★ — Five patterns including Magentic One and A2A protocol for cross-runtime collaboration. Best enterprise orchestration capabilities
- **Tool integration (Priority #5):** ★★★★☆ — Strong Microsoft ecosystem integration. MCP support. Custom skills easy to define. Enterprise-centric rather than broad ecosystem

**Best for:** Microsoft Azure environments, enterprise integration with Office 365/Dynamics, organizations requiring multi-language support (.NET/Python/Java), businesses needing vendor support contracts

**Limitations:** Transition complexity as AutoGen merges into Semantic Kernel, Azure-centric design may limit flexibility, version fragmentation in AutoGen components

### 4. AutoGen (AG2 or v0.4) — Score: 78/100

**Why it ranks fourth:** AutoGen pioneered conversational multi-agent systems where agents coordinate through natural language, creating more intuitive collaboration patterns. Event-driven architecture (v0.4) supports asynchronous messaging for responsive systems. MultimodalConversableAgent provides best framework-level vision integration. Strong for applications emphasizing agent dialogue and coordination over rigid workflows.

**Strengths against priority requirements:**
- **Real-time adaptation (Priority #1):** ★★★☆☆ — Memory-based adaptation through conversation history and external stores. Teaching capabilities allow user knowledge injection
- **Long-term memory (Priority #2):** ★★★★☆ — Extensible memory protocol with ChromaDB, Redis, Mem0 backends. Cross-session support via persistent checkpointers
- **Multi-modal processing (Priority #3):** ★★★★☆ — Strongest framework-level integration via MultimodalAgent and VisionCapability. Native image and audio processing
- **Multi-agent orchestration (Priority #4):** ★★★★☆ — Natural language coordination, group chat patterns, hierarchical structures. Event-driven (v0.4) enables distributed agents
- **Tool integration (Priority #5):** ★★★☆☆ — Built-in code execution, web search, function calling. MCP support. AutoGen Studio for rapid prototyping

**Best for:** Conversational multi-agent applications, research prototypes, teams comfortable with version management, applications emphasizing natural agent dialogue over deterministic workflows

**Limitations:** Version fragmentation creates complexity (0.2 vs 0.4 vs AG2), merger with Semantic Kernel adds uncertainty, aggressive context pruning can lose information, smaller production user base than LangGraph

### 5. CrewAI — Score: 72/100

**Why it ranks fifth:** CrewAI excels at making multi-agent development accessible with role-based abstractions and built-in memory systems. Sequential and hierarchical patterns work well for defined workflows. Enabling memory requires simply setting `memory=True`—the easiest setup among competitors. Strong for rapid prototyping and teams prioritizing simplicity over flexibility.

**Strengths against priority requirements:**
- **Real-time adaptation (Priority #1):** ★★☆☆☆ — Basic memory accumulation for preference learning. Claims crews "improve over time" but mechanism limited to context retrieval
- **Long-term memory (Priority #2):** ★★★☆☆ — Built-in short-term (ChromaDB), long-term (SQLite), and entity memory. Easy to enable but less sophisticated than Letta or LangGraph
- **Multi-modal processing (Priority #3):** ★★☆☆☆ — Relies entirely on model capabilities without framework abstractions. Requires custom tool integration
- **Multi-agent orchestration (Priority #4):** ★★★☆☆ — Sequential and hierarchical patterns work well. Simpler than LangGraph but less flexible. Agents don't directly communicate
- **Tool integration (Priority #5):** ★★★☆☆ — Basic built-in tools, LangChain compatibility, automatic tool selection. Narrower ecosystem than alternatives

**Best for:** Rapid prototyping and quick starts, well-defined sequential workflows, teams prioritizing ease of use over advanced features, smaller projects without enterprise scale requirements

**Limitations:** Sequential model limits flexibility, no direct agent communication, less sophisticated memory than top-tier frameworks, SQLite scaling concerns for high-throughput applications

## Architecture recommendations for your super-assistant

**Recommended architecture using Letta:**

```
┌─────────────────────────────────────────┐
│         Orchestrator Agent              │
│  - Unified interface to user            │
│  - Routes to specialized agents         │
│  - Manages cross-agent memory           │
│  - Persistent persona + human blocks    │
└────────────┬────────────────────────────┘
             │
     ┌───────┼────────┐
     │       │        │
┌────▼───┬───▼───┬────▼───────┐
│Research│Vision │Task        │
│Agent   │Agent  │Agent       │
│(Web +  │(Multi-│(Calendar + │
│ RAG)   │modal) │Email + API)│
└────────┴───────┴────────────┘

Memory Layer (Letta Architecture):
├─ Core Memory: Persona + Human blocks (self-editing)
├─ Archival Memory: Vector DB for semantic history
├─ Filesystem: Documents, PDFs, transcripts
└─ Shared Blocks: Cross-agent collaboration

Tool Layer:
├─ MCP Servers: Filesystem, web, APIs
├─ Composio: 100+ pre-built integrations
└─ Custom Tools: Domain-specific functions

Deployment:
├─ Letta Cloud: Managed infrastructure (free tier)
├─ PostgreSQL: Production persistence
└─ Sleep-Time Agents: Background processing
```

**Alternative architecture using LangGraph:**

```
┌──────────────────────────────────────────┐
│        Supervisor Node (StateGraph)      │
│  - Conditional routing logic             │
│  - Tool-based worker selection           │
│  - Checkpointer persistence              │
└────────────┬─────────────────────────────┘
             │
     ┌───────┼────────┐
     │       │        │
┌────▼───┬───▼───┬────▼───────┐
│Research│Vision │Assistant   │
│Subgraph│Node   │Subgraph    │
└────────┴───────┴────────────┘

Memory Layer (LangGraph Architecture):
├─ PostgresSaver: Thread-scoped checkpoints
├─ RedisStore: Cross-thread user memory + vector search
├─ LangMem: Memory consolidation service (optional)
└─ Manual summarization: Developer-coded strategies

Tool Layer:
├─ 1000+ LangChain integrations
├─ MCP adapters
├─ Custom @tool functions
└─ ToolNode execution

Deployment:
├─ LangGraph Platform: Managed (GA May 2025)
├─ BYOC: Hybrid deployment (data in your VPC)
├─ LangGraph Studio: Visual debugging
└─ LangSmith: Observability and tracing
```

## Clear recommendation: Choose Letta for memory-first, LangGraph for production-first

**Choose Letta if you need:**
- Agents that genuinely learn and adapt from interactions
- Self-managing memory without explicit coding
- Sleep-time background processing for memory refinement
- Model-agnostic agent state (portable across providers)
- Cutting-edge memory architecture from academic research
- Personal assistant or research agent use cases
- Python-based development with open-source flexibility

**Choose LangGraph if you need:**
- Production-proven infrastructure with enterprise adoption
- Visual debugging and time-travel capabilities
- Human-in-the-loop workflows with state inspection
- Managed platform with auto-scaling (LangGraph Platform)
- Largest tool ecosystem (1000+ integrations)
- Complex deterministic workflows requiring explicit control
- Team prioritizing reliability over cutting-edge features

**Choose Microsoft Agent Framework if you need:**
- Azure ecosystem integration (Office 365, Dynamics, Azure AI)
- Enterprise support contracts and compliance (SOC 2)
- Multi-language requirements (.NET, Python, Java)
- Long-running business process workflows
- Cross-runtime agent collaboration (A2A protocol)
- Organization heavily invested in Microsoft infrastructure

**Hybrid approach (Advanced):** Use Letta agents as specialized workers with LangGraph orchestration. Letta handles memory-intensive agents (research, personalization) while LangGraph manages overall workflow coordination and production infrastructure. This combines Letta's superior memory with LangGraph's orchestration strengths.

**Implementation timeline:**
1. **Weeks 1-2:** Prototype core agent with Letta to validate self-editing memory meets requirements
2. **Weeks 3-4:** Add specialized agents (research, task execution) with shared memory blocks
3. **Weeks 5-6:** Integrate multi-modal capabilities via GPT-4o/Gemini for vision and voice
4. **Weeks 7-8:** Deploy sleep-time agents for background processing and memory refinement
5. **Weeks 9-10:** Production hardening with PostgreSQL persistence and monitoring
6. **Weeks 11-12:** Iterate based on real usage patterns and memory evolution

For your personal Eagle Eye/Jarvis-style prototype with hybrid hosting requirements, **Letta emerges as the clear choice**. Its revolutionary self-editing memory architecture, sleep-time compute, and model-agnostic persistence directly address your top two priorities: real-time adaptation and long-term memory. The framework's research foundation (UC Berkeley's MemGPT paper) combined with production-ready infrastructure (Apache 2.0 license, Docker deployment, REST API) makes it ideal for sophisticated prototypes that can scale to production. Start with Letta's free cloud tier or self-hosted Docker deployment, leverage its built-in multi-agent communication for orchestration, and integrate MCP/Composio tools for extensive API capabilities. This approach delivers a genuine learning assistant rather than a stateless chatbot with retrieval capabilities.
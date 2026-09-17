# The definitive map of AI agent frameworks in 2025-2026

The AI agent framework ecosystem has exploded from a handful of experimental projects in 2023 to **over 80 significant frameworks, platforms, and protocols** competing across every layer of the stack by early 2026. Three clear production leaders have emerged — **LangGraph, CrewAI, and the newly unified Microsoft Agent Framework** — while MCP (Model Context Protocol) has become the universal standard connecting them all, with 97 million monthly SDK downloads and 10,000+ deployed servers. The landscape is rapidly consolidating: first-wave experimental frameworks like BabyAGI and AutoGPT have faded, hyperscalers are giving away frameworks as on-ramps to paid inference, and the battleground has shifted from "which framework" to "which infrastructure stack." This report catalogs every significant framework across all categories with current technical details, community metrics, and strategic positioning.

---

## Python frameworks dominate the core agent layer

The Python ecosystem remains the center of gravity for AI agent development, hosting the largest and most mature frameworks. These range from full orchestration platforms to lightweight toolkits, with a clear split emerging between **graph-based orchestration** (LangGraph, Microsoft Agent Framework) and **role-based multi-agent** (CrewAI, MetaGPT) paradigms.

### Tier 1: Production leaders

**LangChain / LangGraph** (LangChain Inc., Harrison Chase) — The most widely adopted ecosystem with **~90,000 GitHub stars** for LangChain and **~25,000** for LangGraph. LangChain itself has pivoted to become primarily a RAG and document framework, explicitly telling users to "use LangGraph for agents." LangGraph models agents as finite state machines with nodes, edges, and conditional routing, supporting cyclical graphs for iterative behaviors. It offers first-class checkpointing, parallel processing, and streaming. Used in production by LinkedIn, Uber, Klarna, and Replit across 400+ companies, with **4.2 million monthly downloads**. Open source (MIT). LangGraph Platform provides commercial deployment; LangSmith handles observability. Best for complex stateful workflows requiring fine-grained control. → github.com/langchain-ai/langgraph

**CrewAI** (João Moura) — The fastest-growing multi-agent framework with **~30,000 stars**, $18M Series A funding, and adoption by **60% of Fortune 500 companies**. Its two-layer architecture combines Crews (dynamic role-based agent collaboration) and Flows (deterministic event-driven orchestration). Each agent gets defined roles, goals, and backstories mimicking human teams. Claims **5.76x faster execution** than LangGraph in certain benchmarks. Supports 700+ tool integrations and MCP. **100,000+ certified developers** execute 60M+ agent runs monthly. Open source with CrewAI Enterprise for commercial use. Best for rapid multi-agent prototyping where teams need production agents in ~2 weeks. → github.com/crewAIInc/crewAI

**Microsoft Agent Framework** (Microsoft — merger of AutoGen + Semantic Kernel) — The biggest structural change in the ecosystem. In October 2025, Microsoft merged AutoGen (~35,000 stars) and Semantic Kernel (~27,000 stars) into a unified framework supporting Python, C#/.NET, and Java. It combines AutoGen's research-grade multi-agent orchestration with Semantic Kernel's enterprise plugin system. Graph-based workflows with streaming, checkpointing, human-in-the-loop, and time-travel debugging. Supports MCP, A2A, and OpenAPI connectors with deep Azure integration. GA targeted Q1 2026. Open source (MIT). The original AutoGen creators separately forked the project as **AG2** (ag2ai/ag2, Apache 2.0) under open community governance. → github.com/microsoft/agent-framework

**OpenAI Agents SDK** (OpenAI) — Released March 2025 as the production successor to the experimental Swarm project. Built on four core primitives: **Agents, Handoffs, Guardrails, and Sessions**. Extremely lightweight with provider-agnostic support for 100+ LLMs. Notable for its **RealtimeAgent** enabling voice agents via WebRTC with automatic interruption detection. Python SDK has **~18,000 stars**; TypeScript SDK ~2,100 stars with ~128K weekly npm downloads. Built-in tracing, MCP integration, and structured outputs via Pydantic. Open source (MIT). → github.com/openai/openai-agents-python

**Google ADK** (Google, Agent Development Kit) — Launched April 2025 at Cloud Next with **~12,000 stars** and 7M+ downloads. Code-first modular framework creating agents in under 100 lines of Python. Three agent types: LLM Agents (reasoning), Workflow Agents (deterministic Sequential/Parallel/Loop), and Custom Agents. Native A2A protocol and MCP support. TypeScript, Go, and Java SDKs added late 2025. Deploys to Cloud Run and Vertex AI Agent Engine. Free, open source (Apache 2.0). → github.com/google/adk-python

### Tier 2: Strong contenders with distinct niches

| Framework | Stars | Creator | Architecture | Key Differentiator | Maturity | URL |
|-----------|-------|---------|-------------|-------------------|----------|-----|
| **Agno** (ex-PhiData) | ~26K | Agno AGI | Lightweight multi-agent | 5,000x faster instantiation, <5μs startup, multimodal | Production | github.com/agno-agi/agno |
| **PydanticAI** | ~15K | Pydantic Inc. (Sequoia-backed, $17M) | Type-safe agents | FastAPI-style DX, OpenTelemetry built-in, A2A support | Production-stable | github.com/pydantic/pydantic-ai |
| **LlamaIndex Workflows** | ~38K | LlamaIndex ($38M raised) | Event-driven, data-first | Leading document/RAG platform, 300+ integrations, LlamaParse | Production | github.com/run-llama/llama_index |
| **Haystack** | ~24K | deepset ($44M raised) | Pipeline-based DAGs | Context engineering focus, enterprise (Airbus, LEGO), since 2019 | Production | github.com/deepset-ai/haystack |
| **DSPy** | ~32K | Stanford NLP (Omar Khattab) | Declarative programming | Programs (not prompts) LMs; auto-optimizer; PyTorch-inspired | Very mature | github.com/stanfordnlp/dspy |
| **Letta** (ex-MemGPT) | ~17K | Letta Inc. ($10M raised) | Memory-first stateful | OS-like memory hierarchy, self-editing context, V1 arch Feb 2026 | Production | github.com/letta-ai/letta |
| **Smolagents** | ~15K | HuggingFace | Code-first ReAct | ~1,000 lines core; agents write Python, not JSON; Hub integration | Young but backed | github.com/huggingface/smolagents |
| **MetaGPT** | ~45K | DeepWisdom | Multi-agent SWE simulation | Simulates software company (PM→Architect→Engineer); ICLR 2024 | Mature | github.com/FoundationAgents/MetaGPT |
| **Strands Agents** | ~5K | AWS | Model-driven minimal | Used internally by Amazon Q; A2A + MCP; Bedrock AgentCore deploy | Production-ready | github.com/strands-agents/sdk-python |

### Tier 3: Specialized and experimental Python frameworks

**Instructor** (~12.4K stars, Jason Liu / 567 Labs) is not an agent framework but the most-used structured output library with **6M+ monthly downloads**, cited by OpenAI as inspiration for their structured outputs feature. It validates LLM outputs against Pydantic schemas with automatic retries. → github.com/567-labs/instructor

**CAMEL** (~8K stars, KAUST research) scales to 1M+ agent simulations with role-playing communicative agents, and its OWL project was accepted at NeurIPS 2025. **Marvin** (~5.2K stars, PrefectHQ) provides one-line AI functions and has superseded ControlFlow (now archived). **Ell** (~6K stars) treats prompts as versioned programs with automatic commit messages and a visualization studio. **Burr** (~2K stars, now Apache Incubator) offers zero-dependency state machines with pluggable persistence. **Atomic Agents** (~2K stars) takes a LEGO-block approach with schema-driven I/O alignment. **Mirascope** (~2.5K stars) positions itself as the "LLM Anti-Framework" with a decorator-based unified interface across providers. **Agency Swarm** (~3.5K stars, VRSEN) builds organizational multi-agent hierarchies on OpenAI's SDK.

**BabyAGI** (~20K stars, Yohei Nakajima) pioneered task-driven autonomous agents in 2023 but its original repo was archived September 2024; the v2 self-building agent concept remains experimental. **SuperAGI** (~15K stars) provided a full-stack GUI-based autonomous agent platform but development has slowed significantly. Both are considered historical milestones rather than active production options.

---

## TypeScript frameworks are catching up fast

TypeScript surpassed Python as the most-used language on GitHub in 2025, and the AI agent framework ecosystem is following this trend. Every major AI company launched or announced TypeScript SDKs in 2025.

**Vercel AI SDK** (~21,900 stars, Vercel) dominates by sheer adoption with **2.8 million weekly npm downloads** — more than all other TS frameworks combined. Now at version 6 with ToolLoopAgent for agentic loops, streaming-first architecture, React Server Components support, and MCP integration. Not a full agent framework but the foundational layer most TS agent apps build on. Free, open source (Apache 2.0). → ai-sdk.dev

**Mastra** (~19,800 stars, by Gatsby.js founders, YC W25) is the leading TypeScript-native agent framework, reaching stable 1.0 in January 2026. It offers typed step-by-step workflow orchestration with `.then()`, `.branch()`, `.parallel()` syntax, built-in memory, voice module, evals, and interactive Studio for debugging. **300,000+ weekly npm downloads** with production use at Replit, PayPal, and WorkOS. Free, open source (Apache 2.0). → mastra.ai

**LangChain.js** (~16,600 stars) and **LangGraph.js** (~2,400 stars, but 529K weekly npm downloads) bring the LangChain ecosystem to JavaScript. Composio (~26,500 stars) serves as a tool integration middleware connecting 1,000+ external services to any agent framework with managed authentication. **Botpress** (~14,600 stars, founded 2017) is the most mature conversational AI platform with visual flow builders and multichannel deployment. **Voiceflow** (proprietary, 250K+ users) offers no-code voice and chat agent building used by Amazon and Louis Vuitton.

Newer entrants include **VoltAgent** (~4,500 stars, MIT) with built-in LLM observability and trigger-based automation, **KaibanJS** with Redux-inspired state management and Kanban-style visual boards, and the **OpenAI Agents SDK TypeScript** and **Google ADK TypeScript** versions launched in mid-to-late 2025.

---

## Low-code platforms are democratizing agent building

The visual builder segment has produced some of the fastest-growing repositories in the entire AI space, led by platforms that lower the barrier from "developer" to "technically curious."

**n8n** (~162,000 stars, valued at $2.5B with NVIDIA investment) gained 112,000+ new stars in 2025 alone to become the #1 JavaScript Rising Star. Its 400+ app integrations now include 70+ native AI nodes with LangChain integration, MCP server support, and Agent/Plan-and-Execute nodes. Fair-code license for self-hosting; cloud plans start at $20/month. → n8n.io

**Dify** (~130,000 stars, LangGenius) crossed 100K stars in 2025 and won AWS's Social Impact Partner of the Year award. Its visual workflow builder with Agent nodes, built-in RAG, and Prompt IDE supports hundreds of LLMs. Open source (Apache 2.0) with Docker deployment and AWS Marketplace availability. → dify.ai

**Langflow** (~100,000 stars) hit its milestone in August 2025 with full MCP client and server support, making it a visual interface for building and deploying MCP-compatible agents. Integrated with IBM watsonx Orchestrate. Open source (MIT). → langflow.org

**Flowise** (~49,000 stars) provides drag-and-drop LLM flow building on LangChain/LangGraph with OpenAI Realtime API integration for voice. Self-hosted or Flowise Cloud. Open source (Apache 2.0). → flowiseai.com

**Coze** (ByteDance) offers a no-code bot builder with MCP support, short/long-term memory, plugin ecosystem, and deployment to Discord/Telegram/Slack. Notably, ByteDance open-sourced the Coze Studio engine (Apache 2.0) in early 2026. Plans from free to $39/month. → coze.com

**Stack AI** (MIT founders, YC, $16M Series A) and **Relevance AI** ($10M Series A, Sydney) target enterprise workflows with SOC 2/HIPAA compliance, though neither supports MCP or A2A natively.

---

## Enterprise cloud platforms are embedding agents everywhere

Every hyperscaler now offers a managed agent platform, with frameworks increasingly given away free as on-ramps to paid compute and inference.

**Google Vertex AI Agent Builder** combines the open-source ADK with Agent Engine (managed serverless runtime), Agent Designer (visual low-code canvas), and Agent Garden (pre-built samples). Full MCP and A2A support with bidirectional audio/video streaming. Pay-as-you-go pricing (vCPU hours + GiB hours). Used by 7M+ ADK downloads. → cloud.google.com/products/agent-builder

**Microsoft Foundry** (renamed from Azure AI Foundry at Ignite 2025) offers 11,000+ models, the Foundry MCP Server (cloud-hosted at mcp.ai.azure.com), 1,400+ Logic Apps connectors, and Foundry IQ for enterprise RAG. Supports MCP and A2A. The open-source Microsoft Agent Framework deploys natively. 70,000+ organizations using the platform. → azure.microsoft.com/en-us/products/ai-foundry

**Amazon Bedrock Agents** provides access to hundreds of foundation models with AgentCore (launched 2025) — a modular runtime with gateway, memory, observability, browser, and code interpreter services. AgentCore Gateway supports MCP servers. Complex consumption-based pricing with models like Claude, Llama, and Amazon Nova. 100K+ organizations using Bedrock. → aws.amazon.com/bedrock

**IBM watsonx Orchestrate** features 150+ pre-built agents, 400+ pre-built tools, and Langflow visual integration. Fine-tuned Granite models handle multi-turn orchestration across HR, procurement, finance, and supply chain domains. Recognized as Leader in 2025 Gartner Magic Quadrant. Custom enterprise pricing with hybrid cloud deployment (any hyperscaler, Red Hat OpenShift). → ibm.com/products/watsonx-orchestrate

**Salesforce Agentforce** (formerly Einstein Agents) runs on the Atlas Reasoning Engine deeply integrated with Sales/Service/Field Service clouds. Industry-specific agents for financial services, healthcare, manufacturing. Complex pricing: $2/conversation or Flex Credits at 100K/$500, plus base Salesforce licenses. Total TCO often **$500+/user/month**. → salesforce.com/agentforce

**ServiceNow AI Agents** includes AI Agent Studio (natural language agent creation), AI Agent Orchestrator (multi-agent coordination), and AI Agent Fabric (external agent integration). Ranked #1 by Gartner for Building/Managing AI Agents. Yokohama release (2025) added thousands of pre-configured agents. Enterprise-only pricing with 60% uplift for Pro Plus tier. → servicenow.com/products/ai-agents.html

---

## Specialized frameworks serve vertical use cases

### Voice and real-time agents

The voice agent segment saw explosive growth, with text-only frameworks losing **30% market share** to voice-capable alternatives. Two open-source leaders dominate the infrastructure layer:

**LiveKit Agents** (~9,300 stars, Apache 2.0) provides WebRTC-based real-time voice AI with sub-second latency, semantic turn detection via a custom transformer model, and telephony integration (SIP/PSTN). Modular plugin system supports 20+ STT/LLM/TTS providers. → github.com/livekit/agents

**Pipecat** (~7,000 stars, BSD-2, by Daily.co) uses a pipeline-based architecture chaining processors for streaming audio/video. Achieves **500-800ms round-trip** voice-to-voice latency. Transport-agnostic with 20+ provider integrations. → github.com/pipecat-ai/pipecat

On the proprietary side, **Retell AI** leads with ~600ms latency, 40M+ calls/month, $35M+ ARR, and 99.99% uptime SLA at $0.07+/min. **Vapi** offers a visual Flow Studio with 1,000+ templates at ~$0.10-0.15/min. **Bland AI** (YC, $65M raised) handles enterprise phone calling at scale (20,000 calls/hour) at $0.09/min. **ElevenLabs Conversational AI** provides industry-leading voice quality with ultra-realistic AI voices, built-in RAG, and MCP support.

### Browser and web agents

**Browser Use** (~78,700 stars, MIT) has become the most-starred browser automation framework, using Playwright with a custom ChatBrowserUse model that's 3-5x faster than alternatives. Cloud version offers stealth browsers with CAPTCHA solving across 195+ country proxies. → github.com/browser-use/browser-use

**Stagehand** (~15,000 stars, by Browserbase) offers a hybrid approach with three primitives — `act()`, `extract()`, `agent()` — combining deterministic code with AI fallbacks. Action caching enables runs without LLM inference after first pass. Multi-language SDKs (TypeScript, Python, Rust, C#, Java). **Skyvern** (~20,500 stars) uses computer vision instead of DOM parsing, making it resistant to layout changes. **Notte** (~4,000 stars, YC S25) introduces a perception layer converting webpages into structured, LLM-readable action APIs, claiming 50%+ cost reduction.

### Coding agents

**OpenHands** (formerly OpenDevin, ~65,000 stars, MIT, $18.8M raised) provides sandboxed Docker environments where agents edit code, run commands, browse web, and call APIs. Its CodeAct agent is competitive on SWE-bench and published at ICLR 2025. → github.com/OpenHands/OpenHands

**Aider** (~40,800 stars, Apache 2.0) is the most popular AI pair programming tool, creating a "repository map" of codebase structure for LLM context with automatic git commits. Supports voice-to-code and any LLM. **SWE-agent** (~15,000 stars, Princeton/Stanford) achieves **72% on SWE-bench verified** with Claude Sonnet 4.5 and introduced a minimal 100-line version scoring >74%. **Devin** (Cognition Labs, proprietary, $73M ARR) operates as a fully autonomous AI software engineer with cloud IDE, though reliability on complex tasks remains mixed.

### Robotics

**LeRobot** (~18,000 stars, HuggingFace, Apache 2.0) provides an end-to-end ML framework for real-world robotics with a unified Robot class interface, standardized dataset format, support for 10+ robot platforms, and pre-trained models on HuggingFace Hub. → github.com/huggingface/lerobot

---

## Memory and coordination infrastructure forms a new layer

### Memory-first architectures

A critical infrastructure layer has emerged specifically for persistent agent memory, with all three leaders combining vector stores and knowledge graphs:

**Mem0** (~47,700 stars, YC-backed, $24M Series A) leads in adoption with **14M+ Python downloads** and 186M API calls in Q3 2025. Its hybrid approach layers vector stores with optional Neo4j knowledge graphs, adding decay/confidence metrics and automatic conflict resolution. AWS chose Mem0 as its exclusive memory provider for the Strands Agents SDK. Native integrations with CrewAI, Flowise, and Langflow. Open source (Apache 2.0) with Mem0 Cloud. → mem0.ai

**Zep** (YC-backed, Getzep Inc.) builds on the **Graphiti** temporal knowledge graph engine with a three-tier subgraph hierarchy. Achieves sub-200ms retrieval latency and outperforms MemGPT on benchmarks (+18.5% accuracy, -90% latency on LongMemEval). Community Edition deprecated in favor of Zep Cloud. → getzep.com

**Cognee** (~6,000 stars, Topoteretes) offers the simplest integration (5-6 lines) with its ECL (Extract, Cognify, Load) pipeline, combining LanceDB vectors with Kuzu graph store. Featured in Microsoft's AI Agents for Beginners course. → cognee.ai

### Protocol standards: MCP and A2A

**MCP (Model Context Protocol)**, created by Anthropic in November 2024, has become the **de facto industry standard** for agent-to-tool communication. Donated to the Linux Foundation's Agentic AI Foundation (AAIF) in December 2025 — co-founded by Anthropic, OpenAI, and Block. The numbers tell the story: **97M+ monthly SDK downloads, 10,000+ deployed servers**, adoption by ChatGPT, Cursor, Gemini, VS Code, and every major agent framework. Organizations report **40-60% faster agent deployment times**. NVIDIA CEO Jensen Huang called it work that "completely revolutionized the AI landscape." The November 2025 spec update added async operations, OAuth 2.1, and a community registry.

**A2A (Agent-to-Agent Protocol)**, launched by Google at Cloud Next in April 2025, complements MCP by handling inter-agent communication rather than agent-to-tool connections. It uses HTTP, SSE, and JSON-RPC with Agent Card discovery at `/.well-known/agent.json`. Moved to Linux Foundation governance in June 2025 with 150+ supporting organizations including Adobe, SAP, Salesforce, and Microsoft. While adoption has been slower than MCP's, it fills a distinct role for enterprise multi-agent systems. → a2a-protocol.org

---

## Comprehensive comparison matrix

| Framework | Language | Architecture | Memory | Voice | MCP | Open Source | Stars | Maturity |
|-----------|----------|-------------|--------|-------|-----|-------------|-------|----------|
| LangGraph | Python/TS | Graph state machine | Checkpoint-based | ❌ | ✅ | MIT | ~25K | Production |
| CrewAI | Python | Role-based multi-agent | Per-agent + shared | ❌ | ✅ | Open + Enterprise | ~30K | Production |
| MS Agent Framework | Python/C#/Java | Graph workflows | Thread-based | ❌ | ✅ | MIT | ~4K (new) | Preview→GA Q1 '26 |
| OpenAI Agents SDK | Python/TS | Minimal primitives | Sessions | ✅ | ✅ | MIT | ~18K | Production |
| Google ADK | Python/TS/Go/Java | Modular code-first | Sessions + Memory Bank | ✅ (audio) | ✅ | Apache 2.0 | ~12K | Production |
| Agno | Python | Lightweight multi-agent | DB-backed sessions | ❌ | ✅ | Open + Enterprise | ~26K | Production |
| PydanticAI | Python | Type-safe agents | Pydantic models | ❌ | ❌ | MIT | ~15K | Stable |
| LlamaIndex | Python/TS | Event-driven workflows | Flexible Memory Blocks | ❌ | ✅ | Open + Commercial | ~38K | Production |
| Haystack | Python | Pipeline DAGs | Pipeline-level | ❌ | ❌ | Apache 2.0 | ~24K | Production |
| DSPy | Python | Declarative programming | Stateless (optimized) | ❌ | ❌ | MIT | ~32K | Very mature |
| Letta | Python/TS | Memory-first stateful | Core + archival hierarchy | ❌ | ❌ | Apache 2.0 | ~17K | Production |
| Smolagents | Python | Code-first ReAct | Per-agent history | ❌ | ✅ | Apache 2.0 | ~15K | Young |
| MetaGPT | Python | SWE simulation | Shared message pool | ❌ | ❌ | MIT | ~45K | Mature |
| Strands (AWS) | Python/TS | Model-driven minimal | Session managers | ✅ (streaming) | ✅ | Apache 2.0 | ~5K | Production |
| Mastra | TypeScript | Workflow-driven | Built-in composite | ✅ | ✅ | Apache 2.0 | ~20K | Stable 1.0 |
| Vercel AI SDK | TypeScript | Streaming-first UI | External | ✅ (limited) | ✅ | Apache 2.0 | ~22K | Very mature |
| n8n | JS/TS | Visual workflows | Buffer memory | ✅ (STT) | ✅ | Fair-code | ~162K | Very mature |
| Dify | Python/TS | Visual agent builder | Built-in KB | ❌ | ❌ | Apache 2.0 | ~130K | Very mature |
| Langflow | Python | Visual multi-agent | Built-in | ❌ | ✅ | MIT | ~100K | Mature |
| Flowise | JS | Visual drag-and-drop | LangChain memory | ✅ (Realtime) | ❌ | Apache 2.0 | ~49K | Stable |
| LiveKit Agents | Python/Node | WebRTC voice | Event-driven | ✅ Core | ❌ | Apache 2.0 | ~9.3K | Production |
| Pipecat | Python | Pipeline voice | Frame-based | ✅ Core | ❌ | BSD-2 | ~7K | Production |
| Browser Use | Python | Browser automation | Session-based | ❌ | ❌ | MIT | ~79K | Production |
| OpenHands | Python | Sandboxed coding | Event-stream | ❌ | ❌ | MIT | ~65K | Production |
| Aider | Python | Pair programming | Git-based | ✅ (voice-to-code) | ❌ | Apache 2.0 | ~41K | Very mature |
| Mem0 | Python/TS | Memory layer | Graph + vector hybrid | ❌ | ✅ (server) | Apache 2.0 | ~48K | Production |
| Composio | Python/TS | Tool integration | Context management | ❌ | ✅ | MIT | ~27K | Mature |

---

## Momentum scoreboard heading into 2026

**Accelerating rapidly:** LangGraph (production standard, 400+ enterprise deployments), CrewAI (60% Fortune 500, 100K+ developers), OpenAI Agents SDK (fastest-growing new framework), Mastra (TS-native winner), n8n (gained 112K stars in 2025), Browser Use (explosive growth to 79K stars), Mem0 ($24M Series A, AWS partnership), PydanticAI (type-safety movement)

**Steady growth with strong positioning:** Google ADK, Agno, Smolagents, Letta, Strands (AWS), Dify, Langflow, LlamaIndex, DSPy, CAMEL/OWL

**Declining or transitioning:** BabyAGI (archived, now experimental v2 only), AutoGPT (faded from 2023 hype peak, agents loop endlessly), SuperAGI (development slowed significantly), standalone AutoGen and Semantic Kernel (being deprecated in favor of Microsoft Agent Framework), OpenAI Swarm (superseded by Agents SDK), AgentGPT (relegated to demos), ControlFlow (archived, superseded by Marvin 3.0)

The broader pattern is clear: the "first wave" frameworks of 2023 (AutoGPT, BabyAGI, AgentGPT, SuperAGI) served as proof-of-concept demonstrations but have been **largely superseded by production-grade successors**. Frameworks without clear differentiation or active development are being absorbed or abandoned.

---

## Seven architectural trends defining the next era

**MCP as mandatory infrastructure.** Frameworks without native MCP support are now considered legacy. The protocol reduced integration complexity from quadratic to linear — instead of N×M custom connectors, each tool builds one MCP server and each agent builds one MCP client. The December 2025 donation to the Linux Foundation's AAIF (co-founded by Anthropic, OpenAI, Block, with AWS, Google, Microsoft, and 20+ others as members) cemented it as the permanent standard.

**Code-as-action replacing JSON function calling.** Smolagents pioneered the CodeAgent approach where LLMs write Python code to route data rather than emitting JSON tool calls. This enables more flexible, composable tool use and is being adopted across multiple frameworks. The split between "model-driven" (let the LLM decide) and "workflow-driven" (developer-defined graphs) approaches remains the key architectural choice.

**Memory evolving from vectors to graph+vector hybrids.** All three leading memory systems (Mem0, Zep, Cognee) combine vector stores with knowledge graphs. Three memory types are gaining traction: **episodic** (what happened), **semantic** (what is known), and **procedural** (how to do things). AutoGPT notably abandoned its vector database as unnecessary complexity, while Letta's OS-like memory hierarchy (core/archival/recall) is increasingly influential.

**Context engineering supplanting prompt engineering.** The hot skill of 2025 shifted from crafting individual prompts to engineering the entire information pipeline — what data reaches the model, in what structure, with what metadata. Anthropic's "Skills" specification for procedural tool knowledge and OpenAI's AGENTS.md format (adopted by 60,000+ open-source projects) formalized this approach.

**Durable execution becoming non-negotiable.** Frameworks optimized for demos failed in production. Temporal-based orchestration (PydanticAI), event sourcing, and deterministic replay (LangGraph) moved from optional to mandatory. One analysis found PydanticAI's Temporal integration **prevented 100% of state desynchronization incidents** that plagued earlier frameworks.

**Observability as baseline, not afterthought.** Langfuse went fully MIT in June 2025; Helicone crossed 2 billion traced interactions. The standard stack now includes LLM tracing (Langfuse/Helicone/LangSmith), infrastructure metrics (Prometheus/OpenTelemetry), and business KPIs. Every new 2025 framework ships with built-in tracing.

**Human-in-the-loop as standard pattern.** Emerged as one of seven "must-know" design patterns, with Anthropic positioned as "most aligned with organizational realities — where trust, not just capability, remains the bottleneck." Enterprise deployments follow a three-tier framework: Foundation (trust) → Workflow (automation) → Autonomy (governance-enabled).

---

## Conclusion: the consolidation year arrives

The AI agent framework ecosystem has matured from chaotic experimentation into a clearer hierarchy. Gartner recorded a **1,445% surge** in multi-agent system inquiries from Q1 2024 to Q2 2025, yet only 2% of organizations have deployed at scale. The agentic AI market reached $7.6 billion in 2025 and is projected to hit $196.6 billion by 2034 at 43.8% CAGR, but Gartner warns that over 40% of agentic AI projects will be canceled by 2027 due to escalating costs and unclear value.

The framework war has narrowed to three production winners — **LangGraph for complex orchestration, CrewAI for rapid role-based development, and Microsoft Agent Framework for enterprise .NET/Azure** — with OpenAI Agents SDK and Google ADK as strong platform-aligned alternatives. MCP won the protocol war decisively; A2A complements it for inter-agent communication. Hyperscalers are following a "razor and blades" model: frameworks are free, compute is the revenue. The analogy to the container orchestration wars of 2015 is apt — just as Docker Swarm and Mesos gave way to Kubernetes, the long tail of agent frameworks will converge toward a few standards.

The critical insight for 2026 is that the bottleneck has shifted from "which framework to use" to "how to bridge experimentation to production at scale." As multiple experts noted, 2025 was the testing phase; 2026 is the year agentic workflows must finally move from demos into day-to-day practice — or face the cancellation wave Gartner predicts.
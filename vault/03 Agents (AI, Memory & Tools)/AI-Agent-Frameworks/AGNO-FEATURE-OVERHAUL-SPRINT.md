# Agno Feature Overhaul Sprint

> **Created:** 2026-01-25
> **Status:** Planning Complete - Ready for Execution
> **Linear Project:** Agno Feature Overhaul Sprint
> **Total Issues:** 45 issues across 9 epics (OMOS-158 to OMOS-202)
> **Estimated Effort:** 190-230 hours
> **Cookbook Review:** 2026-01-25 - Added Epic 9 from Agno cookbook patterns
> **Previous Work:** 2026-01-23 - 31 features implemented (see Phase 0 below)

---

## Phase 0: Previously Implemented (2026-01-23)

> **Status:** ✅ Complete
> **Agents:** 4 parallel implementation tracks
> **Features:** 31 features across 6 categories

This section documents features already implemented before the current sprint planning.

### Implementation Summary

| Agent | Focus Area | Features | Status |
|-------|------------|----------|--------|
| **Lead** | AgentOS, RBAC, MCP | 4 | ✅ Complete |
| **Agent 1** | Context & Memory | 8 | ✅ Complete |
| **Agent 2** | Execution Control | 9 | ✅ Complete |
| **Agent 3** | Observability & Knowledge | 10 | ✅ Complete |

### Features Already Implemented

#### 1. AgentOS Integration (`core/os/`)
- `OneMindOS` class - Extended AgentOS wrapper with RBAC and MCP
- `create_onemind_os()` - Async factory function
- `DEFAULT_SCOPE_MAPPINGS` - Pre-configured RBAC scopes
- Automatic MCP server connection
- CORS middleware integration

#### 2. JWT RBAC Middleware (`core/middleware/`)
- `OneMindJWTMiddleware` - JWT-based access control
- Path-based scope requirements with wildcard support
- Admin scope bypass (`onemind:admin`)
- Token extraction from header or cookie
- 30+ scope mappings pre-configured

#### 3. Enhanced MCP Integration (`core/mcp/`)
- `MCPServerConfig` dataclass for server configuration
- `MCPTransport` enum (stdio, sse, streamable-http)
- 12 pre-configured servers (filesystem, memory, github, context7, etc.)
- `MCP_PRESETS` - basic, web, dev, research, full

#### 4. Background Hooks & @hook Decorator
- `@hook` decorator with `run_in_background` parameter
- Automatic sync/async function handling
- Pre-built validation hooks

#### 5. Context & Memory (Agent 1)
- Dependency injection with 3 helpers
- `add_dependencies_to_context=True`
- `enable_agentic_memory=True`
- `enable_session_summaries=True`
- `LegacyResponse` and `TaskResponse` structured output models

#### 6. Execution Control (Agent 2)
- `OpenAIModerationGuardrail` - Blocks harmful content
- Input/output validation hooks
- Run cancellation (`POST /api/agents/{id}/runs/{run_id}/cancel`)
- HITL continue run (`POST /api/agents/{id}/runs/{run_id}/continue`)
- Reasoning mode support

#### 7. Observability & Knowledge (Agent 3)
- Arize Phoenix tracing (OpenInference)
- Langfuse tracing (OTEL export)
- Knowledge upload API (`POST /knowledge/upload`)
- LanceDB hybrid search (vector + full-text + BM25)
- Accuracy evaluation framework
- Debug/telemetry toggles

### File Structure (Already Created)
```
core/
├── os/                        # AgentOS wrapper
│   └── app.py                 # OneMindOS class
├── middleware/                # RBAC
│   └── jwt.py                 # JWT middleware
├── guardrails/                # Content safety
│   └── moderation.py          # OpenAI moderation
├── mcp/                       # MCP integration
│   └── tools.py               # Full MCPTools support
├── observability/             # Tracing
│   ├── phoenix.py             # Arize Phoenix
│   └── langfuse.py            # Langfuse
└── evals/                     # Evaluation
    └── accuracy.py            # Accuracy evals
```

---

## Linear Issue Reference

| Epic | Issues | Priority |
|------|--------|----------|
| **Epic 1: AgentOS Architecture** | OMOS-158 to OMOS-163 | P0 Critical |
| **Epic 2: Workflow System** | OMOS-164 to OMOS-168 | P0 Critical |
| **Epic 3: Built-in Tools** | OMOS-169 to OMOS-173 | P1 High |
| **Epic 4: Memory System** | OMOS-174 to OMOS-178 | P1 High |
| **Epic 5: Team System** | OMOS-179 to OMOS-182 | P1 High |
| **Epic 6: Model Providers** | OMOS-183 to OMOS-186 | P2 Medium |
| **Epic 7: Multimodal & Voice** | OMOS-187 to OMOS-190 | P2 Medium |
| **Epic 8: Evals & Observability** | OMOS-191 to OMOS-194 | P2 Medium |
| **Epic 9: Advanced Cookbook Features** | OMOS-195 to OMOS-202 | P1 High |

### Complete Issue List

| Issue | Title | Epic | Priority |
|-------|-------|------|----------|
| OMOS-158 | Implement POST /agents endpoint for dynamic agent creation | Epic 1 | P0 |
| OMOS-159 | Implement per-run model selection for any agent | Epic 1 | P0 |
| OMOS-160 | Implement dynamic tool attachment per agent | Epic 1 | P0 |
| OMOS-161 | Make LiveKit/voice accessible to all agents | Epic 1 | P0 |
| OMOS-162 | Implement shared knowledge base accessible by all agents | Epic 1 | P0 |
| OMOS-163 | Implement agent persistence for restart recovery | Epic 1 | P1 |
| OMOS-164 | Implement POST /workflows endpoint for dynamic workflow creation | Epic 2 | P0 |
| OMOS-165 | Implement all Agno workflow step types | Epic 2 | P0 |
| OMOS-166 | Persist workflow definitions and runs to database | Epic 2 | P1 |
| OMOS-167 | Implement common workflow templates | Epic 2 | P1 |
| OMOS-168 | Integrate HITL approvals into workflows | Epic 2 | P1 |
| OMOS-169 | Enable DuckDuckGo, Wikipedia, Arxiv, Exa search tools | Epic 3 | P1 |
| OMOS-170 | Enable YFinance financial data tools | Epic 3 | P1 |
| OMOS-171 | Enable Calculator, Shell, Python, SQL, E2B tools | Epic 3 | P1 |
| OMOS-172 | Enable Email, Slack, Discord tools | Epic 3 | P1 |
| OMOS-173 | Enable GitHub and enhanced Linear tools | Epic 3 | P1 |
| OMOS-174 | Implement enable_user_memories for persistent user context | Epic 4 | P1 |
| OMOS-175 | Implement enable_session_summaries for conversation compression | Epic 4 | P1 |
| OMOS-176 | Implement enable_agentic_memory for self-learning agents | Epic 4 | P1 |
| OMOS-177 | Implement semantic memory search with embeddings | Epic 4 | P1 |
| OMOS-178 | Update Flutter memory screen with new features | Epic 4 | P2 |
| OMOS-179 | Implement POST /teams endpoint for dynamic team creation | Epic 5 | P1 |
| OMOS-180 | Implement all team coordination modes | Epic 5 | P1 |
| OMOS-181 | Persist team definitions and runs to database | Epic 5 | P1 |
| OMOS-182 | Enable voice for team sessions | Epic 5 | P2 |
| OMOS-183 | Enable all OpenAI models including o1/o3 | Epic 6 | P2 |
| OMOS-184 | Enable Groq ultra-fast inference | Epic 6 | P2 |
| OMOS-185 | Enable Together AI for large open models | Epic 6 | P2 |
| OMOS-186 | Enable DeepSeek models | Epic 6 | P2 |
| OMOS-187 | Enable DALL-E image generation | Epic 7 | P2 |
| OMOS-188 | Enable vision/image input for agents | Epic 7 | P2 |
| OMOS-189 | Enable audio input and output | Epic 7 | P2 |
| OMOS-190 | Enable video frame analysis | Epic 7 | P2 |
| OMOS-191 | Implement Agno AccuracyEval system | Epic 8 | P2 |
| OMOS-192 | Implement batch evaluation for multiple agents | Epic 8 | P2 |
| OMOS-193 | Implement trace export to observability platforms | Epic 8 | P2 |
| OMOS-194 | Create agent metrics dashboard | Epic 8 | P2 |
| OMOS-195 | Implement ReasoningTools for explicit structured thinking | Epic 9 | P1 |
| OMOS-196 | Implement Router workflow step for dynamic agent selection | Epic 9 | P1 |
| OMOS-197 | Implement Loop workflow step for iterative execution | Epic 9 | P1 |
| OMOS-198 | Implement @hook decorator for background task execution | Epic 9 | P1 |
| OMOS-199 | Implement session_state for stateful agents | Epic 9 | P1 |
| OMOS-200 | Implement PerformanceEval and ReliabilityEval systems | Epic 9 | P2 |
| OMOS-201 | Implement MemoryTools for agent-controlled memory | Epic 9 | P1 |
| OMOS-202 | Implement RunContext for context-aware tools | Epic 9 | P1 |

---

## Executive Summary

This document contains the complete specification for overhauling OneMind OS to fully utilize all Agno framework features. The current implementation has **34 of 154 features working (22%)**, with critical gaps in:

1. **Dynamic Agent/Team/Workflow Creation** - Currently hardcoded to single "Legacy" agent
2. **Workflow System** - Missing Parallel, Condition, evaluator support
3. **Built-in Tools** - 19 toolkits not enabled
4. **Memory System** - Missing user memories, session summaries, agentic memory
5. **Universal Voice/LiveKit** - Locked to single agent instead of all agents
6. **Model Flexibility** - Can't select model per-run

---

## Architecture Vision

### Current State (Problematic)
```
┌─────────────────────────────────────────────┐
│              OneMind OS                      │
│  ┌───────────────────────────────────────┐  │
│  │         Legacy Agent (HARDCODED)      │  │
│  │  - Fixed model (Claude Opus)          │  │
│  │  - Fixed tools                        │  │
│  │  - Fixed knowledge                    │  │
│  │  - LiveKit locked here                │  │
│  └───────────────────────────────────────┘  │
│                    │                         │
│            No dynamic creation               │
│            No per-run configuration          │
└─────────────────────────────────────────────┘
```

### Target State (Universal)
```
┌─────────────────────────────────────────────────────────────┐
│                    OneMind AgentOS                          │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              Universal Agent Factory                 │   │
│  │  POST /agents → Create any agent dynamically        │   │
│  │  POST /teams → Create any team dynamically          │   │
│  │  POST /workflows → Create any workflow dynamically  │   │
│  └─────────────────────────────────────────────────────┘   │
│                          │                                  │
│  ┌───────────┐  ┌───────────┐  ┌───────────┐  ┌─────────┐ │
│  │  Agent A  │  │  Agent B  │  │  Agent C  │  │ Team X  │ │
│  │  GPT-4    │  │  Claude   │  │  Llama    │  │ 3 agents│ │
│  │  +Search  │  │  +Code    │  │  +Local   │  │ +coord  │ │
│  │  +Voice   │  │  +Voice   │  │  +Voice   │  │ +Voice  │ │
│  └───────────┘  └───────────┘  └───────────┘  └─────────┘ │
│                          │                                  │
│  ┌─────────────────────────────────────────────────────┐   │
│  │           Universal Services Layer                   │   │
│  │  - LiveKit (all agents)   - Knowledge (shared)      │   │
│  │  - Memory (per-agent)     - Tools (attachable)      │   │
│  │  - Evals (any agent)      - Traces (unified)        │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

---

## Epic 1: AgentOS Architecture Overhaul (CRITICAL)

**Priority:** P0 - Critical
**Estimated Hours:** 40
**Dependencies:** None (foundation for all other work)

### Issue 1.1: Dynamic Agent Creation API

**Title:** Implement POST /agents endpoint for dynamic agent creation

**Description:**
Currently agents are hardcoded in `core/agent/legacy.py`. We need to implement the Agno AgentOS pattern where agents can be created dynamically via API.

**Agno Pattern:**
```python
from agno.agent import Agent
from agno.playground import Playground

# Agent factory function
def create_agent(config: AgentConfig) -> Agent:
    return Agent(
        name=config.name,
        model=create_model(config.model_id),
        tools=resolve_tools(config.tools),
        knowledge=resolve_knowledge(config.knowledge_id),
        instructions=config.instructions,
        storage=get_storage(config.storage_type),
        # Memory features
        enable_user_memories=config.enable_user_memories,
        enable_agentic_memory=config.enable_agentic_memory,
        enable_session_summaries=config.enable_session_summaries,
        # Output features
        response_model=config.response_model,
        structured_outputs=config.structured_outputs,
        markdown=config.markdown,
    )

# AgentOS registration
playground = Playground(agents=[])

@app.post("/agents")
async def create_agent_endpoint(config: AgentConfig):
    agent = create_agent(config)
    playground.agents.append(agent)
    return {"agent_id": agent.agent_id, "name": agent.name}
```

**Current Code to Modify:**
- `core/agent/legacy.py` - Extract to factory pattern
- `core/orchestrator/main.py` - Add POST /agents endpoint
- `core/agent/factory.py` - NEW: Agent factory with all options

**API Specification:**
```yaml
POST /agents
Request:
  name: string (required)
  model: string (required) # e.g., "claude-opus-4", "gpt-4o", "llama-3.1"
  instructions: string (optional)
  tools: string[] (optional) # e.g., ["duckduckgo", "wikipedia", "linear"]
  knowledge_id: string (optional) # Reference to knowledge base
  enable_user_memories: boolean (default: false)
  enable_agentic_memory: boolean (default: false)
  enable_session_summaries: boolean (default: false)
  response_model: string (optional) # Pydantic model name
  structured_outputs: boolean (default: false)
  markdown: boolean (default: true)
  storage_type: string (default: "postgres") # "postgres", "sqlite", "memory"

Response:
  agent_id: string
  name: string
  model: string
  tools: string[]
  created_at: datetime
```

**Files to Create:**
```
core/agent/
├── factory.py          # Agent factory with all options (~300 lines)
├── registry.py         # In-memory agent registry (~100 lines)
├── models.py           # Pydantic models for API (~150 lines)
└── __init__.py         # Updated exports
```

**Acceptance Criteria:**
- [ ] POST /agents creates agent with specified configuration
- [ ] GET /agents lists all registered agents
- [ ] GET /agents/{id} returns agent details
- [ ] DELETE /agents/{id} removes agent from registry
- [ ] Agent persists in database for restart recovery
- [ ] All Agno Agent parameters supported

---

### Issue 1.2: Universal Model Selection

**Title:** Implement per-run model selection for any agent

**Description:**
Currently model is hardcoded per agent. Agno supports selecting model at run time, allowing any agent to use any model.

**Agno Pattern:**
```python
from agno.models.openai import OpenAIChat
from agno.models.anthropic import Claude
from agno.models.groq import Groq
from agno.models.together import Together
from agno.models.deepseek import DeepSeek
from agno.models.ollama import Ollama

MODEL_REGISTRY = {
    # OpenAI
    "gpt-4o": lambda: OpenAIChat(id="gpt-4o"),
    "gpt-4o-mini": lambda: OpenAIChat(id="gpt-4o-mini"),
    "gpt-4-turbo": lambda: OpenAIChat(id="gpt-4-turbo"),
    "o1": lambda: OpenAIChat(id="o1"),
    "o1-mini": lambda: OpenAIChat(id="o1-mini"),
    "o3-mini": lambda: OpenAIChat(id="o3-mini"),

    # Anthropic
    "claude-opus-4": lambda: Claude(id="claude-opus-4-20250514"),
    "claude-sonnet-4": lambda: Claude(id="claude-sonnet-4-20250514"),
    "claude-3.5-sonnet": lambda: Claude(id="claude-3-5-sonnet-20241022"),
    "claude-3.5-haiku": lambda: Claude(id="claude-3-5-haiku-20241022"),

    # Groq (ultra-fast)
    "llama-3.3-70b": lambda: Groq(id="llama-3.3-70b-versatile"),
    "llama-3.1-8b": lambda: Groq(id="llama-3.1-8b-instant"),
    "mixtral-8x7b": lambda: Groq(id="mixtral-8x7b-32768"),

    # Together
    "llama-3.1-405b": lambda: Together(id="meta-llama/Meta-Llama-3.1-405B-Instruct-Turbo"),
    "qwen-2.5-72b": lambda: Together(id="Qwen/Qwen2.5-72B-Instruct-Turbo"),

    # DeepSeek
    "deepseek-chat": lambda: DeepSeek(id="deepseek-chat"),
    "deepseek-reasoner": lambda: DeepSeek(id="deepseek-reasoner"),

    # Ollama (local)
    "ollama-llama3.2": lambda: Ollama(id="llama3.2"),
    "ollama-phi3": lambda: Ollama(id="phi3"),
    "ollama-qwen2.5": lambda: Ollama(id="qwen2.5"),
}

def get_model(model_id: str):
    if model_id not in MODEL_REGISTRY:
        raise ValueError(f"Unknown model: {model_id}")
    return MODEL_REGISTRY[model_id]()

# Per-run model selection
@app.post("/agents/{agent_id}/runs")
async def run_agent(agent_id: str, request: RunRequest):
    agent = get_agent(agent_id)

    # Override model for this run if specified
    if request.model:
        agent = agent.model_copy(update={"model": get_model(request.model)})

    response = await agent.arun(request.message)
    return response
```

**Current Code to Modify:**
- `core/agent/legacy.py` - Remove hardcoded model
- `core/orchestrator/main.py` - Accept model in run request
- `core/models/registry.py` - NEW: Model registry

**API Specification:**
```yaml
POST /agents/{agent_id}/runs
Request:
  message: string (required)
  model: string (optional) # Override agent's default model
  stream: boolean (default: true)

GET /models
Response:
  models:
    - id: "gpt-4o"
      provider: "openai"
      context_length: 128000
      supports_vision: true
      supports_tools: true
    - id: "claude-opus-4"
      provider: "anthropic"
      context_length: 200000
      supports_vision: true
      supports_tools: true
    # ... etc
```

**Files to Create:**
```
core/models/
├── registry.py         # Model registry (~200 lines)
├── providers/
│   ├── openai.py      # OpenAI models (~50 lines)
│   ├── anthropic.py   # Anthropic models (~50 lines)
│   ├── groq.py        # Groq models (~50 lines)
│   ├── together.py    # Together models (~50 lines)
│   ├── deepseek.py    # DeepSeek models (~50 lines)
│   ├── ollama.py      # Ollama models (~50 lines)
│   └── __init__.py
└── __init__.py
```

**Environment Variables Required:**
```bash
OPENAI_API_KEY=sk-...
ANTHROPIC_API_KEY=sk-ant-...
GROQ_API_KEY=gsk_...
TOGETHER_API_KEY=...
DEEPSEEK_API_KEY=...
OLLAMA_HOST=http://localhost:11434
```

**Acceptance Criteria:**
- [ ] GET /models returns all available models with metadata
- [ ] POST /agents/{id}/runs accepts optional model parameter
- [ ] Model override works for single run without persisting
- [ ] All 6 providers implemented (OpenAI, Anthropic, Groq, Together, DeepSeek, Ollama)
- [ ] Graceful fallback if provider API key missing

---

### Issue 1.3: Dynamic Tool Attachment

**Title:** Implement dynamic tool attachment per agent

**Description:**
Tools should be attachable to any agent at creation time or runtime, not hardcoded.

**Agno Pattern:**
```python
from agno.tools.duckduckgo import DuckDuckGoTools
from agno.tools.wikipedia import WikipediaTools
from agno.tools.arxiv import ArxivTools
from agno.tools.yfinance import YFinanceTools
from agno.tools.newspaper4k import Newspaper4kTools
from agno.tools.calculator import CalculatorTools
from agno.tools.shell import ShellTools
from agno.tools.file import FileTools
from agno.tools.python import PythonTools
from agno.tools.sql import SQLTools
from agno.tools.email import EmailTools
from agno.tools.slack import SlackTools
from agno.tools.discord import DiscordTools
from agno.tools.github import GithubTools
from agno.tools.linear import LinearTools
from agno.tools.exa import ExaTools
from agno.tools.firecrawl import FirecrawlTools
from agno.tools.e2b import E2BTools
from agno.tools.replicate import ReplicateTools

TOOL_REGISTRY = {
    # Search & Research
    "duckduckgo": DuckDuckGoTools(),
    "wikipedia": WikipediaTools(),
    "arxiv": ArxivTools(),
    "exa": ExaTools(),
    "firecrawl": FirecrawlTools(),
    "newspaper": Newspaper4kTools(),

    # Finance
    "yfinance": YFinanceTools(
        stock_price=True,
        analyst_recommendations=True,
        company_info=True,
        company_news=True
    ),

    # Code & Execution
    "calculator": CalculatorTools(),
    "shell": ShellTools(),
    "file": FileTools(),
    "python": PythonTools(),
    "sql": SQLTools(),
    "e2b": E2BTools(),

    # Communication
    "email": EmailTools(),
    "slack": SlackTools(),
    "discord": DiscordTools(),

    # Development
    "github": GithubTools(),
    "linear": LinearTools(),

    # Media
    "replicate": ReplicateTools(),

    # Custom OneMind tools
    "todoist": TodoistTools(),
    "contact": ContactTools(),
    "awareness": AwarenessTools(),
}

def resolve_tools(tool_names: list[str]) -> list:
    """Resolve tool names to tool instances."""
    tools = []
    for name in tool_names:
        if name in TOOL_REGISTRY:
            tools.append(TOOL_REGISTRY[name])
        else:
            raise ValueError(f"Unknown tool: {name}")
    return tools

# Attach tools at creation
agent = Agent(
    name="Research Agent",
    tools=resolve_tools(["duckduckgo", "wikipedia", "arxiv"]),
)

# Or attach tools dynamically
@app.post("/agents/{agent_id}/tools")
async def attach_tools(agent_id: str, tool_names: list[str]):
    agent = get_agent(agent_id)
    new_tools = resolve_tools(tool_names)
    agent.tools.extend(new_tools)
    return {"tools": [t.name for t in agent.tools]}
```

**Current Code to Modify:**
- `core/tools/` - Organize all tools with registry
- `core/agent/factory.py` - Use tool registry
- `core/orchestrator/main.py` - Add tool attachment endpoint

**API Specification:**
```yaml
GET /tools
Response:
  tools:
    - name: "duckduckgo"
      description: "Search the web with DuckDuckGo"
      category: "search"
      requires_auth: false
    - name: "github"
      description: "Interact with GitHub repositories"
      category: "development"
      requires_auth: true
      auth_env: "GITHUB_TOKEN"

POST /agents/{agent_id}/tools
Request:
  tools: ["duckduckgo", "wikipedia"]
Response:
  agent_id: string
  tools: ["duckduckgo", "wikipedia", "calculator"]

DELETE /agents/{agent_id}/tools/{tool_name}
Response:
  agent_id: string
  tools: ["duckduckgo", "calculator"]
```

**Files to Create:**
```
core/tools/
├── registry.py         # Tool registry (~150 lines)
├── builtin/
│   ├── search.py      # DuckDuckGo, Wikipedia, Arxiv, Exa
│   ├── finance.py     # YFinance
│   ├── code.py        # Calculator, Shell, Python, SQL
│   ├── comms.py       # Email, Slack, Discord
│   ├── dev.py         # GitHub, Linear
│   └── media.py       # Replicate
├── custom/
│   ├── todoist.py     # Existing
│   ├── linear.py      # Existing
│   ├── contact.py     # Existing
│   └── awareness.py   # Existing
└── __init__.py
```

**Acceptance Criteria:**
- [ ] GET /tools lists all available tools with metadata
- [ ] POST /agents/{id}/tools attaches tools to agent
- [ ] DELETE /agents/{id}/tools/{name} removes tool
- [ ] All 19 Agno built-in tools available
- [ ] Tool auth requirements documented

---

### Issue 1.4: Universal LiveKit Integration

**Title:** Make LiveKit/voice accessible to all agents

**Description:**
Currently LiveKit is hardcoded to the Legacy agent. Every agent should be able to use voice.

**Agno Pattern:**
```python
from livekit import api as livekit_api
from agno.agent import Agent

class VoiceService:
    """Universal voice service for all agents."""

    def __init__(self):
        self.livekit_url = os.getenv("LIVEKIT_URL")
        self.livekit_api_key = os.getenv("LIVEKIT_API_KEY")
        self.livekit_api_secret = os.getenv("LIVEKIT_API_SECRET")

    async def create_room(self, agent_id: str, user_id: str) -> dict:
        """Create a LiveKit room for any agent."""
        room_name = f"agent-{agent_id}-{user_id}-{uuid4().hex[:8]}"

        # Create room
        room_service = livekit_api.RoomServiceClient(
            self.livekit_url,
            self.livekit_api_key,
            self.livekit_api_secret
        )
        await room_service.create_room(livekit_api.CreateRoomRequest(name=room_name))

        # Generate tokens
        user_token = self._generate_token(room_name, user_id, is_agent=False)
        agent_token = self._generate_token(room_name, agent_id, is_agent=True)

        return {
            "room_name": room_name,
            "user_token": user_token,
            "agent_token": agent_token,
            "livekit_url": self.livekit_url,
        }

    async def connect_agent(self, agent: Agent, room_name: str, token: str):
        """Connect any agent to a LiveKit room."""
        # Agent joins room and can speak/listen
        pass

# API endpoints
voice_service = VoiceService()

@app.post("/agents/{agent_id}/voice/room")
async def create_voice_room(agent_id: str, user_id: str):
    """Create voice room for any agent."""
    return await voice_service.create_room(agent_id, user_id)

@app.get("/agents/{agent_id}/voice/token")
async def get_voice_token(agent_id: str, user_id: str):
    """Get voice token for any agent."""
    room = await voice_service.create_room(agent_id, user_id)
    return {"token": room["user_token"], "url": room["livekit_url"]}
```

**Current Code to Modify:**
- `core/voice/` - Extract from Legacy-specific code
- `core/orchestrator/main.py` - Add universal voice endpoints

**API Specification:**
```yaml
POST /agents/{agent_id}/voice/room
Request:
  user_id: string (required)
Response:
  room_name: string
  user_token: string
  agent_token: string
  livekit_url: string

GET /agents/{agent_id}/voice/status
Response:
  agent_id: string
  active_rooms: int
  connected: boolean

DELETE /agents/{agent_id}/voice/room/{room_name}
Response:
  success: boolean
```

**Files to Create:**
```
core/voice/
├── service.py          # Universal voice service (~200 lines)
├── livekit.py          # LiveKit client wrapper (~150 lines)
├── models.py           # Voice-related models (~50 lines)
└── __init__.py
```

**Acceptance Criteria:**
- [ ] Any agent can create voice room
- [ ] Any agent can join voice room
- [ ] Voice works with teams (multiple agents)
- [ ] Voice state persisted across restarts
- [ ] Graceful disconnect handling

---

### Issue 1.5: Shared Knowledge Base

**Title:** Implement shared knowledge base accessible by all agents

**Description:**
Knowledge should be uploadable and searchable by any agent, not locked to one.

**Agno Pattern:**
```python
from agno.knowledge.pdf import PDFKnowledgeBase
from agno.knowledge.text import TextKnowledgeBase
from agno.knowledge.json import JSONKnowledgeBase
from agno.knowledge.csv import CSVKnowledgeBase
from agno.knowledge.website import WebsiteKnowledgeBase
from agno.knowledge.combined import CombinedKnowledgeBase
from agno.vectordb.lancedb import LanceDb

# Shared vector database
vector_db = LanceDb(
    uri="data/lancedb",
    table_name="knowledge",
)

# Knowledge base registry
KNOWLEDGE_BASES = {}

def create_knowledge_base(kb_id: str, kb_type: str, sources: list) -> KnowledgeBase:
    """Create knowledge base of any type."""
    if kb_type == "pdf":
        kb = PDFKnowledgeBase(
            path=sources,
            vector_db=vector_db,
        )
    elif kb_type == "text":
        kb = TextKnowledgeBase(
            path=sources,
            vector_db=vector_db,
        )
    elif kb_type == "website":
        kb = WebsiteKnowledgeBase(
            urls=sources,
            vector_db=vector_db,
        )
    elif kb_type == "combined":
        kb = CombinedKnowledgeBase(
            sources=[...],
            vector_db=vector_db,
        )

    kb.load()  # Load into vector DB
    KNOWLEDGE_BASES[kb_id] = kb
    return kb

# Any agent can use any knowledge base
agent = Agent(
    name="Research Agent",
    knowledge=KNOWLEDGE_BASES["company_docs"],
    search_knowledge=True,
    update_knowledge=True,  # Can add to knowledge
)
```

**API Specification:**
```yaml
POST /knowledge
Request:
  name: string
  type: "pdf" | "text" | "website" | "json" | "csv"
  sources: string[] # File paths or URLs
Response:
  knowledge_id: string
  name: string
  type: string
  document_count: int

GET /knowledge
Response:
  knowledge_bases:
    - id: string
      name: string
      type: string
      document_count: int
      created_at: datetime

POST /knowledge/{id}/search
Request:
  query: string
  limit: int (default: 5)
Response:
  results:
    - content: string
      score: float
      metadata: object

POST /agents/{agent_id}/knowledge
Request:
  knowledge_id: string
Response:
  agent_id: string
  knowledge_ids: string[]
```

**Files to Create:**
```
core/knowledge/
├── base.py             # Knowledge base registry (~200 lines)
├── types/
│   ├── pdf.py         # PDF knowledge (~100 lines)
│   ├── text.py        # Text knowledge (~80 lines)
│   ├── website.py     # Website knowledge (~120 lines)
│   ├── json.py        # JSON knowledge (~80 lines)
│   └── csv.py         # CSV knowledge (~80 lines)
├── vector/
│   └── lancedb.py     # LanceDB integration (~150 lines)
└── __init__.py
```

**Acceptance Criteria:**
- [ ] POST /knowledge creates knowledge base from files/URLs
- [ ] Any agent can attach any knowledge base
- [ ] Knowledge search returns relevant results
- [ ] Knowledge updates persist
- [ ] Multiple agents can share same knowledge

---

### Issue 1.6: Agent Persistence & Recovery

**Title:** Implement agent persistence for restart recovery

**Description:**
Agents created via API should persist across server restarts.

**Implementation:**
```python
from sqlalchemy import Column, String, JSON, DateTime
from sqlalchemy.orm import declarative_base

Base = declarative_base()

class AgentRecord(Base):
    __tablename__ = "agents"

    id = Column(String, primary_key=True)
    name = Column(String, nullable=False)
    model = Column(String, nullable=False)
    instructions = Column(String)
    tools = Column(JSON, default=[])
    knowledge_ids = Column(JSON, default=[])
    config = Column(JSON, default={})
    created_at = Column(DateTime, default=datetime.utcnow)
    updated_at = Column(DateTime, onupdate=datetime.utcnow)

class AgentPersistence:
    """Persist agents to database."""

    async def save(self, agent: Agent) -> None:
        record = AgentRecord(
            id=agent.agent_id,
            name=agent.name,
            model=agent.model.id,
            instructions=agent.instructions,
            tools=[t.name for t in agent.tools],
            knowledge_ids=[...],
            config={...},
        )
        await self.db.merge(record)

    async def load_all(self) -> list[Agent]:
        records = await self.db.query(AgentRecord).all()
        return [self._record_to_agent(r) for r in records]

    async def delete(self, agent_id: str) -> None:
        await self.db.query(AgentRecord).filter_by(id=agent_id).delete()

# On startup, load all persisted agents
@app.on_event("startup")
async def load_agents():
    persistence = AgentPersistence(db)
    agents = await persistence.load_all()
    for agent in agents:
        register_agent(agent)
```

**Database Schema:**
```sql
CREATE TABLE agents (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    model VARCHAR(100) NOT NULL,
    instructions TEXT,
    tools JSONB DEFAULT '[]',
    knowledge_ids JSONB DEFAULT '[]',
    config JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

CREATE INDEX idx_agents_name ON agents(name);
CREATE INDEX idx_agents_model ON agents(model);
```

**Acceptance Criteria:**
- [ ] Agents persist to database on creation
- [ ] Agents restore on server restart
- [ ] Agent updates persist immediately
- [ ] Agent deletion removes from DB

---

## Epic 2: Workflow System (CRITICAL)

**Priority:** P0 - Critical
**Estimated Hours:** 35
**Dependencies:** Epic 1 (Agent factory)

### Issue 2.1: Dynamic Workflow Creation

**Title:** Implement POST /workflows endpoint for dynamic workflow creation

**Description:**
Currently workflows are hardcoded. Need to implement Agno's dynamic workflow creation with Steps, Parallel, and Condition support.

**Agno Workflow Pattern:**
```python
from agno.workflow import Workflow, Task
from agno.workflow.step import Step, Parallel, Condition
from agno.agent import Agent

# Step-based workflow
class ResearchWorkflow(Workflow):
    """Multi-step research workflow."""

    researcher: Agent = Agent(name="Researcher", tools=[DuckDuckGoTools()])
    writer: Agent = Agent(name="Writer")
    reviewer: Agent = Agent(name="Reviewer")

    def run(self, topic: str) -> str:
        # Step 1: Research
        research = self.researcher.run(f"Research: {topic}")

        # Step 2: Write
        draft = self.writer.run(f"Write article based on: {research.content}")

        # Step 3: Review
        final = self.reviewer.run(f"Review and improve: {draft.content}")

        return final.content

# Parallel workflow
class ParallelResearch(Workflow):
    web_searcher: Agent
    academic_searcher: Agent
    news_searcher: Agent
    synthesizer: Agent

    def run(self, topic: str) -> str:
        # Parallel step - all run concurrently
        results = Parallel(
            Task(self.web_searcher, f"Web search: {topic}"),
            Task(self.academic_searcher, f"Academic search: {topic}"),
            Task(self.news_searcher, f"News search: {topic}"),
        ).run()

        # Sequential synthesis
        return self.synthesizer.run(f"Synthesize: {results}")

# Conditional workflow
class ConditionalWorkflow(Workflow):
    classifier: Agent
    code_agent: Agent
    writing_agent: Agent
    general_agent: Agent

    def run(self, task: str) -> str:
        # Classify task
        task_type = self.classifier.run(f"Classify: {task}").content

        # Conditional routing
        result = Condition(
            condition=lambda: task_type,
            branches={
                "code": Task(self.code_agent, task),
                "writing": Task(self.writing_agent, task),
                "default": Task(self.general_agent, task),
            }
        ).run()

        return result
```

**API Specification:**
```yaml
POST /workflows
Request:
  name: string
  description: string
  steps:
    - type: "step" | "parallel" | "condition"
      agent_id: string (for step)
      agents: string[] (for parallel)
      condition: string (for condition)
      branches: object (for condition)
      prompt_template: string

Response:
  workflow_id: string
  name: string
  steps: object[]

POST /workflows/{id}/run
Request:
  input: object
Response:
  run_id: string
  status: "running" | "completed" | "failed"
  result: object

GET /workflows/{id}/runs
Response:
  runs:
    - run_id: string
      status: string
      started_at: datetime
      completed_at: datetime
      result: object
```

**Files to Create:**
```
core/workflow/
├── base.py             # Workflow base class (~200 lines)
├── steps.py            # Step, Parallel, Condition (~300 lines)
├── registry.py         # Workflow registry (~100 lines)
├── executor.py         # Workflow executor (~250 lines)
├── models.py           # Pydantic models (~150 lines)
└── __init__.py
```

**Acceptance Criteria:**
- [ ] POST /workflows creates workflow from definition
- [ ] Step-based workflows execute sequentially
- [ ] Parallel steps execute concurrently
- [ ] Condition steps route based on criteria
- [ ] Workflow runs persist with results

---

### Issue 2.2: Workflow Step Types

**Title:** Implement all Agno workflow step types

**Description:**
Implement Step, Parallel, and Condition step types with full functionality.

**Step Types:**
```python
from agno.workflow.step import Step, Parallel, Condition
from typing import Callable, Any

class Step:
    """Single agent step."""

    def __init__(
        self,
        agent: Agent,
        prompt_template: str,
        output_key: str = None,
        validators: list[Callable] = None,
    ):
        self.agent = agent
        self.prompt_template = prompt_template
        self.output_key = output_key
        self.validators = validators or []

    async def run(self, context: dict) -> dict:
        prompt = self.prompt_template.format(**context)
        result = await self.agent.arun(prompt)

        # Validate output
        for validator in self.validators:
            if not validator(result):
                raise ValidationError(f"Step validation failed")

        if self.output_key:
            context[self.output_key] = result.content

        return context

class Parallel:
    """Run multiple steps concurrently."""

    def __init__(self, *steps: Step):
        self.steps = steps

    async def run(self, context: dict) -> dict:
        tasks = [step.run(context.copy()) for step in self.steps]
        results = await asyncio.gather(*tasks)

        # Merge all results
        merged = {}
        for result in results:
            merged.update(result)

        return merged

class Condition:
    """Route to different steps based on condition."""

    def __init__(
        self,
        condition: Callable[[dict], str],
        branches: dict[str, Step],
        default: Step = None,
    ):
        self.condition = condition
        self.branches = branches
        self.default = default

    async def run(self, context: dict) -> dict:
        branch_key = self.condition(context)

        if branch_key in self.branches:
            step = self.branches[branch_key]
        elif self.default:
            step = self.default
        else:
            raise ValueError(f"No branch for: {branch_key}")

        return await step.run(context)
```

**Acceptance Criteria:**
- [ ] Step executes single agent with template
- [ ] Parallel executes multiple steps concurrently
- [ ] Condition routes based on evaluator
- [ ] Context passes between steps
- [ ] Validators can reject step output

---

### Issue 2.3: Workflow Persistence

**Title:** Persist workflow definitions and runs

**Description:**
Workflows and their execution history should persist to database.

**Database Schema:**
```sql
CREATE TABLE workflows (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    definition JSONB NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

CREATE TABLE workflow_runs (
    id VARCHAR(255) PRIMARY KEY,
    workflow_id VARCHAR(255) REFERENCES workflows(id),
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    input JSONB,
    output JSONB,
    context JSONB,
    error TEXT,
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE workflow_step_runs (
    id VARCHAR(255) PRIMARY KEY,
    run_id VARCHAR(255) REFERENCES workflow_runs(id),
    step_index INT NOT NULL,
    step_type VARCHAR(50) NOT NULL,
    agent_id VARCHAR(255),
    input JSONB,
    output JSONB,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ
);

CREATE INDEX idx_workflow_runs_workflow ON workflow_runs(workflow_id);
CREATE INDEX idx_workflow_runs_status ON workflow_runs(status);
CREATE INDEX idx_step_runs_run ON workflow_step_runs(run_id);
```

**Acceptance Criteria:**
- [ ] Workflow definitions persist
- [ ] Workflow runs persist with full history
- [ ] Step-level execution tracked
- [ ] Failed runs show error details
- [ ] Runs can be queried by status

---

### Issue 2.4: Built-in Workflow Templates

**Title:** Implement common workflow templates

**Description:**
Provide pre-built workflow templates for common patterns.

**Templates:**
```python
WORKFLOW_TEMPLATES = {
    "research": {
        "name": "Research Workflow",
        "description": "Research a topic and produce a report",
        "steps": [
            {"type": "step", "agent": "researcher", "prompt": "Research: {topic}"},
            {"type": "step", "agent": "writer", "prompt": "Write report: {research}"},
            {"type": "step", "agent": "reviewer", "prompt": "Review: {draft}"},
        ]
    },
    "code_review": {
        "name": "Code Review Workflow",
        "description": "Review code and suggest improvements",
        "steps": [
            {"type": "parallel", "agents": ["security", "quality", "performance"]},
            {"type": "step", "agent": "synthesizer", "prompt": "Combine reviews: {reviews}"},
        ]
    },
    "task_triage": {
        "name": "Task Triage Workflow",
        "description": "Classify and route tasks",
        "steps": [
            {"type": "step", "agent": "classifier", "prompt": "Classify: {task}"},
            {"type": "condition", "branches": {...}},
        ]
    },
    "morning_briefing": {
        "name": "Morning Briefing",
        "description": "Generate daily briefing",
        "steps": [
            {"type": "parallel", "agents": ["calendar", "email", "tasks", "news"]},
            {"type": "step", "agent": "synthesizer", "prompt": "Create briefing: {data}"},
        ]
    },
}

@app.post("/workflows/from-template")
async def create_from_template(template_name: str, customization: dict = None):
    template = WORKFLOW_TEMPLATES[template_name]
    if customization:
        template = deep_merge(template, customization)
    return await create_workflow(template)
```

**Acceptance Criteria:**
- [ ] GET /workflows/templates lists available templates
- [ ] POST /workflows/from-template creates from template
- [ ] Templates can be customized
- [ ] 4+ templates provided

---

### Issue 2.5: Workflow HITL Integration

**Title:** Integrate HITL approvals into workflows

**Description:**
Workflows should support human-in-the-loop approval at any step.

**Implementation:**
```python
class ApprovalStep:
    """Step that requires human approval before continuing."""

    def __init__(
        self,
        description: str,
        timeout_minutes: int = 60,
        auto_approve: bool = False,
    ):
        self.description = description
        self.timeout = timeout_minutes
        self.auto_approve = auto_approve

    async def run(self, context: dict) -> dict:
        if self.auto_approve:
            return context

        # Create approval request
        approval = await create_approval(
            type="workflow_step",
            description=self.description,
            context=context,
            timeout=self.timeout,
        )

        # Wait for approval
        result = await wait_for_approval(approval.id, timeout=self.timeout * 60)

        if result.approved:
            context["approval_notes"] = result.notes
            return context
        else:
            raise WorkflowRejected(result.reason)

# Usage in workflow
workflow = Workflow(
    steps=[
        Step(researcher, "Research: {topic}"),
        ApprovalStep("Review research before writing"),
        Step(writer, "Write based on: {research}"),
        ApprovalStep("Approve draft before publishing"),
        Step(publisher, "Publish: {draft}"),
    ]
)
```

**API Specification:**
```yaml
GET /workflows/{id}/runs/{run_id}/approvals
Response:
  approvals:
    - id: string
      step_index: int
      description: string
      status: "pending" | "approved" | "rejected"
      context: object

POST /workflows/{id}/runs/{run_id}/approvals/{approval_id}
Request:
  action: "approve" | "reject"
  notes: string (optional)
Response:
  success: boolean
```

**Acceptance Criteria:**
- [ ] ApprovalStep pauses workflow
- [ ] Approval request sent to inbox
- [ ] Approval/rejection resumes workflow
- [ ] Timeout handling with configurable behavior
- [ ] Approval history persisted

---

## Epic 3: Built-in Tools Integration (HIGH)

**Priority:** P1 - High
**Estimated Hours:** 20
**Dependencies:** Issue 1.3 (Tool registry)

### Issue 3.1: Search Tools Integration

**Title:** Enable DuckDuckGo, Wikipedia, Arxiv, Exa search tools

**Description:**
Enable all Agno search tools for any agent to use.

**Implementation:**
```python
from agno.tools.duckduckgo import DuckDuckGoTools
from agno.tools.wikipedia import WikipediaTools
from agno.tools.arxiv import ArxivTools
from agno.tools.exa import ExaTools
from agno.tools.newspaper4k import Newspaper4kTools

# DuckDuckGo - web search
duckduckgo = DuckDuckGoTools()
# Functions: duckduckgo_search(query), duckduckgo_news(query)

# Wikipedia - encyclopedia search
wikipedia = WikipediaTools()
# Functions: search_wikipedia(query), get_wikipedia_page(title)

# Arxiv - academic papers
arxiv = ArxivTools()
# Functions: search_arxiv(query), get_arxiv_paper(id)

# Exa - neural search (requires API key)
exa = ExaTools()
# Functions: exa_search(query), exa_get_contents(urls)

# Newspaper4k - article extraction
newspaper = Newspaper4kTools()
# Functions: get_article(url), summarize_article(url)
```

**Environment Variables:**
```bash
EXA_API_KEY=...  # Optional, for Exa search
```

**Acceptance Criteria:**
- [ ] DuckDuckGoTools working
- [ ] WikipediaTools working
- [ ] ArxivTools working
- [ ] ExaTools working (with API key)
- [ ] Newspaper4kTools working
- [ ] All tools in registry

---

### Issue 3.2: Finance Tools Integration

**Title:** Enable YFinance financial data tools

**Description:**
Enable financial data tools for stock analysis.

**Implementation:**
```python
from agno.tools.yfinance import YFinanceTools

yfinance = YFinanceTools(
    stock_price=True,
    analyst_recommendations=True,
    company_info=True,
    company_news=True,
    historical_prices=True,
    key_financial_ratios=True,
    income_statements=True,
    balance_sheets=True,
    cash_flow_statements=True,
)

# Functions:
# - get_stock_price(symbol)
# - get_analyst_recommendations(symbol)
# - get_company_info(symbol)
# - get_company_news(symbol)
# - get_historical_prices(symbol, period)
# - get_key_financial_ratios(symbol)
# - get_income_statement(symbol)
# - get_balance_sheet(symbol)
# - get_cash_flow_statement(symbol)
```

**Acceptance Criteria:**
- [ ] YFinanceTools working
- [ ] All financial functions available
- [ ] Real-time stock data
- [ ] Historical data retrieval

---

### Issue 3.3: Code Execution Tools

**Title:** Enable Calculator, Shell, Python, SQL, E2B tools

**Description:**
Enable code execution tools with appropriate sandboxing.

**Implementation:**
```python
from agno.tools.calculator import CalculatorTools
from agno.tools.shell import ShellTools
from agno.tools.python import PythonTools
from agno.tools.sql import SQLTools
from agno.tools.e2b import E2BTools

# Calculator - safe math evaluation
calculator = CalculatorTools()
# Functions: calculate(expression)

# Shell - command execution (requires confirmation)
shell = ShellTools(
    requires_confirmation=True,  # HITL for dangerous commands
    allowed_commands=["ls", "cat", "grep", "find"],  # Whitelist
)
# Functions: run_shell_command(command)

# Python - Python code execution
python = PythonTools(
    requires_confirmation=True,
    timeout=30,
)
# Functions: run_python_code(code)

# SQL - Database queries
sql = SQLTools(
    db_url="postgresql://...",
    read_only=True,  # Safety
)
# Functions: run_sql_query(query)

# E2B - Cloud sandbox (requires API key)
e2b = E2BTools()
# Functions: run_code(code, language)
```

**Environment Variables:**
```bash
E2B_API_KEY=...  # Optional, for E2B sandbox
```

**Acceptance Criteria:**
- [ ] CalculatorTools working
- [ ] ShellTools working with confirmation
- [ ] PythonTools working with timeout
- [ ] SQLTools working (read-only by default)
- [ ] E2BTools working with API key
- [ ] HITL confirmation for dangerous operations

---

### Issue 3.4: Communication Tools

**Title:** Enable Email, Slack, Discord tools

**Description:**
Enable communication tools for sending messages.

**Implementation:**
```python
from agno.tools.email import EmailTools
from agno.tools.slack import SlackTools
from agno.tools.discord import DiscordTools

# Email - send/read emails
email = EmailTools(
    smtp_host="smtp.gmail.com",
    smtp_port=587,
    imap_host="imap.gmail.com",
    username=os.getenv("EMAIL_USER"),
    password=os.getenv("EMAIL_PASSWORD"),
)
# Functions: send_email(to, subject, body), search_emails(query)

# Slack - workspace messaging
slack = SlackTools(
    token=os.getenv("SLACK_BOT_TOKEN"),
)
# Functions: send_message(channel, text), read_messages(channel)

# Discord - server messaging
discord = DiscordTools(
    token=os.getenv("DISCORD_BOT_TOKEN"),
)
# Functions: send_message(channel_id, text), read_messages(channel_id)
```

**Environment Variables:**
```bash
EMAIL_USER=...
EMAIL_PASSWORD=...
SLACK_BOT_TOKEN=...
DISCORD_BOT_TOKEN=...
```

**Acceptance Criteria:**
- [ ] EmailTools working
- [ ] SlackTools working
- [ ] DiscordTools working
- [ ] Auth credentials managed securely

---

### Issue 3.5: Development Tools

**Title:** Enable GitHub and enhanced Linear tools

**Description:**
Enable full GitHub integration and enhance Linear tools.

**Implementation:**
```python
from agno.tools.github import GithubTools

github = GithubTools(
    token=os.getenv("GITHUB_TOKEN"),
)
# Functions:
# - list_repos()
# - get_repo(owner, repo)
# - list_issues(owner, repo)
# - create_issue(owner, repo, title, body)
# - list_pull_requests(owner, repo)
# - create_pull_request(owner, repo, title, body, head, base)
# - get_file_contents(owner, repo, path)
# - search_code(query)
```

**Environment Variables:**
```bash
GITHUB_TOKEN=...
```

**Acceptance Criteria:**
- [ ] GithubTools working
- [ ] All repo operations available
- [ ] Issue and PR management
- [ ] Code search working

---

## Epic 4: Memory System Overhaul (HIGH)

**Priority:** P1 - High
**Estimated Hours:** 25
**Dependencies:** Epic 1 (Agent factory)

### Issue 4.1: Enable User Memories

**Title:** Implement enable_user_memories for persistent user context

**Description:**
Enable Agno's user memory feature that persists information about users across sessions.

**Agno Pattern:**
```python
from agno.agent import Agent
from agno.memory.db.postgres import PostgresDb

agent = Agent(
    name="Personal Assistant",
    storage=PostgresDb(
        table_name="agent_sessions",
        db_url="postgresql://...",
    ),
    # Enable user memories
    enable_user_memories=True,
    # Optionally enable agentic memory (agent learns on its own)
    enable_agentic_memory=True,
)

# Agent automatically:
# - Extracts facts about user from conversation
# - Stores them in ai.user_memories table
# - Retrieves relevant memories in future conversations

# Example:
# User: "I prefer Python over JavaScript"
# Agent stores: {user_id: "zeus", memory: "prefers Python over JavaScript", category: "preferences"}
# Later...
# User: "What language should I use?"
# Agent recalls: "You mentioned preferring Python over JavaScript"
```

**Database Schema:**
```sql
-- Already exists in Agno
CREATE TABLE ai.user_memories (
    id UUID PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    memory TEXT NOT NULL,
    category VARCHAR(100),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

CREATE INDEX idx_user_memories_user ON ai.user_memories(user_id);
CREATE INDEX idx_user_memories_category ON ai.user_memories(category);
```

**API Specification:**
```yaml
GET /users/{user_id}/memories
Response:
  memories:
    - id: string
      memory: string
      category: string
      created_at: datetime

POST /users/{user_id}/memories
Request:
  memory: string
  category: string (optional)
Response:
  id: string
  memory: string

DELETE /users/{user_id}/memories/{memory_id}
Response:
  success: boolean
```

**Acceptance Criteria:**
- [ ] enable_user_memories=True works
- [ ] Memories extracted from conversation
- [ ] Memories recalled in future sessions
- [ ] API for manual memory management
- [ ] Per-user memory isolation

---

### Issue 4.2: Enable Session Summaries

**Title:** Implement enable_session_summaries for conversation compression

**Description:**
Enable automatic session summarization to maintain context in long conversations.

**Agno Pattern:**
```python
agent = Agent(
    name="Assistant",
    # Enable session summaries
    enable_session_summaries=True,
    # Summary configuration
    summary_model=OpenAIChat(id="gpt-4o-mini"),  # Cheaper model for summaries
    summary_trigger="token_count",  # When to summarize
    summary_token_threshold=4000,  # Summarize after 4000 tokens
)

# Agent automatically:
# - Tracks token count in conversation
# - When threshold reached, summarizes older messages
# - Keeps summary + recent messages in context
# - Maintains continuity across long conversations
```

**Acceptance Criteria:**
- [ ] enable_session_summaries=True works
- [ ] Automatic summarization at threshold
- [ ] Summary replaces old messages
- [ ] Conversation continuity maintained
- [ ] Summary model configurable

---

### Issue 4.3: Enable Agentic Memory

**Title:** Implement enable_agentic_memory for self-learning agents

**Description:**
Enable agents to autonomously learn and store insights without explicit instruction.

**Agno Pattern:**
```python
agent = Agent(
    name="Learning Assistant",
    enable_agentic_memory=True,  # Agent decides what to remember
)

# Agent automatically:
# - Identifies important information during conversations
# - Decides what's worth remembering
# - Stores insights without being asked
# - Uses learnings in future conversations

# Example:
# During task: "The deploy script is at /scripts/deploy.sh"
# Agent autonomously stores: {memory: "deploy script at /scripts/deploy.sh", category: "project"}
# Later: User asks "how do I deploy?"
# Agent recalls the deploy script location
```

**Acceptance Criteria:**
- [ ] enable_agentic_memory=True works
- [ ] Agent identifies memorable information
- [ ] Memories stored automatically
- [ ] Memories retrieved contextually
- [ ] Memory categories assigned automatically

---

### Issue 4.4: Memory Search & Retrieval

**Title:** Implement semantic memory search

**Description:**
Enable semantic search across all memories using vector embeddings.

**Implementation:**
```python
from agno.memory.search import MemorySearch
from agno.embedder.openai import OpenAIEmbedder

memory_search = MemorySearch(
    embedder=OpenAIEmbedder(id="text-embedding-3-small"),
    vector_db=PgVector(
        table_name="memory_embeddings",
        db_url="postgresql://...",
    ),
)

# Index memories
await memory_search.index_memory(memory)

# Search memories
results = await memory_search.search(
    query="What does the user prefer for lunch?",
    user_id="zeus",
    limit=5,
)
```

**Database Schema:**
```sql
CREATE TABLE memory_embeddings (
    id UUID PRIMARY KEY,
    memory_id UUID REFERENCES ai.user_memories(id),
    embedding VECTOR(1536),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_memory_embeddings ON memory_embeddings
USING ivfflat (embedding vector_cosine_ops);
```

**Acceptance Criteria:**
- [ ] Memories embedded on creation
- [ ] Semantic search works
- [ ] Top-k relevant memories returned
- [ ] Search respects user isolation

---

### Issue 4.5: Memory UI Integration

**Title:** Update Flutter memory screen with new features

**Description:**
Update the Flutter memory screen to support user memories, session summaries, and search.

**UI Changes:**
- Add memory categories (preferences, facts, learnings)
- Add memory search
- Show memory source (manual vs agentic)
- Show session summaries
- Add memory analytics

**Acceptance Criteria:**
- [ ] Memory screen shows categories
- [ ] Search bar works with semantic search
- [ ] Manual memory creation
- [ ] Memory deletion
- [ ] Session summary viewer

---

## Epic 5: Team System Enhancement (HIGH)

**Priority:** P1 - High
**Estimated Hours:** 20
**Dependencies:** Epic 1 (Agent factory)

### Issue 5.1: Dynamic Team Creation

**Title:** Implement POST /teams endpoint for dynamic team creation

**Description:**
Enable creating teams of agents dynamically via API.

**Agno Pattern:**
```python
from agno.team import Team
from agno.agent import Agent

# Create team with agents
team = Team(
    name="Research Team",
    agents=[
        Agent(name="Researcher", tools=[DuckDuckGoTools()]),
        Agent(name="Writer"),
        Agent(name="Reviewer"),
    ],
    # Team coordination modes
    mode="coordinate",  # Options: coordinate, collaborate, delegate
    # Delegation settings
    delegate_to_all_members=True,
    respond_directly=False,
    show_members_responses=True,
    # Team leader (optional)
    leader=Agent(name="Lead"),
)

# Run team
response = await team.arun("Research and write about quantum computing")
```

**API Specification:**
```yaml
POST /teams
Request:
  name: string
  description: string
  agent_ids: string[]
  mode: "coordinate" | "collaborate" | "delegate"
  leader_id: string (optional)
  delegate_to_all_members: boolean
  respond_directly: boolean
  show_members_responses: boolean
Response:
  team_id: string
  name: string
  agents: object[]

POST /teams/{id}/run
Request:
  message: string
  stream: boolean
Response:
  run_id: string
  status: string
  result: object
```

**Acceptance Criteria:**
- [ ] POST /teams creates team
- [ ] Team runs with coordination
- [ ] Delegation modes work
- [ ] Team leader routes tasks
- [ ] Team responses aggregated

---

### Issue 5.2: Team Coordination Modes

**Title:** Implement all team coordination modes

**Description:**
Implement coordinate, collaborate, and delegate modes.

**Implementation:**
```python
class TeamMode(Enum):
    COORDINATE = "coordinate"  # Leader assigns tasks, collects results
    COLLABORATE = "collaborate"  # Agents work together, share context
    DELEGATE = "delegate"  # Leader delegates to best agent

class Team:
    async def run_coordinate(self, message: str) -> RunResponse:
        """Leader assigns tasks to each agent."""
        tasks = await self.leader.decompose_task(message)
        results = []
        for task, agent in zip(tasks, self.agents):
            result = await agent.arun(task)
            results.append(result)
        return await self.leader.synthesize(results)

    async def run_collaborate(self, message: str) -> RunResponse:
        """Agents share context and build on each other."""
        context = {"message": message, "history": []}
        for agent in self.agents:
            result = await agent.arun(message, context=context)
            context["history"].append(result)
        return context["history"][-1]

    async def run_delegate(self, message: str) -> RunResponse:
        """Leader picks best agent for task."""
        best_agent = await self.leader.select_agent(message, self.agents)
        return await best_agent.arun(message)
```

**Acceptance Criteria:**
- [ ] Coordinate mode works
- [ ] Collaborate mode works
- [ ] Delegate mode works
- [ ] Mode selectable per run

---

### Issue 5.3: Team Persistence

**Title:** Persist team definitions and runs

**Description:**
Teams should persist to database like agents.

**Database Schema:**
```sql
CREATE TABLE teams (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    mode VARCHAR(50) NOT NULL DEFAULT 'coordinate',
    leader_id VARCHAR(255) REFERENCES agents(id),
    config JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

CREATE TABLE team_members (
    team_id VARCHAR(255) REFERENCES teams(id),
    agent_id VARCHAR(255) REFERENCES agents(id),
    role VARCHAR(100),
    PRIMARY KEY (team_id, agent_id)
);

CREATE TABLE team_runs (
    id VARCHAR(255) PRIMARY KEY,
    team_id VARCHAR(255) REFERENCES teams(id),
    input TEXT,
    output TEXT,
    agent_contributions JSONB,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

**Acceptance Criteria:**
- [ ] Teams persist to database
- [ ] Team members tracked
- [ ] Team runs logged
- [ ] Restore teams on restart

---

### Issue 5.4: Team Voice Integration

**Title:** Enable voice for team sessions

**Description:**
Allow voice interaction with teams where appropriate agents respond.

**Implementation:**
```python
class TeamVoiceSession:
    """Voice session with a team of agents."""

    async def handle_voice_input(self, audio: bytes) -> bytes:
        # Transcribe
        text = await self.transcriber.transcribe(audio)

        # Team processes
        result = await self.team.arun(text)

        # Select speaker (could be leader or delegated agent)
        speaker = result.responding_agent or self.team.leader

        # Synthesize response
        audio_response = await speaker.synthesize(result.content)
        return audio_response
```

**Acceptance Criteria:**
- [ ] Voice works with teams
- [ ] Appropriate agent responds
- [ ] Smooth voice transitions
- [ ] Team context maintained

---

## Epic 6: Model Provider Expansion (MEDIUM)

**Priority:** P2 - Medium
**Estimated Hours:** 15
**Dependencies:** Issue 1.2 (Model registry)

### Issue 6.1: OpenAI Models Integration

**Title:** Enable all OpenAI models including o1/o3

**Description:**
Enable full OpenAI model support including reasoning models.

**Implementation:**
```python
from agno.models.openai import OpenAIChat

OPENAI_MODELS = {
    # GPT-4 family
    "gpt-4o": OpenAIChat(id="gpt-4o"),
    "gpt-4o-mini": OpenAIChat(id="gpt-4o-mini"),
    "gpt-4-turbo": OpenAIChat(id="gpt-4-turbo"),

    # Reasoning models
    "o1": OpenAIChat(id="o1"),
    "o1-mini": OpenAIChat(id="o1-mini"),
    "o3-mini": OpenAIChat(id="o3-mini"),

    # Audio models
    "gpt-4o-audio-preview": OpenAIChat(id="gpt-4o-audio-preview"),
}
```

**Acceptance Criteria:**
- [ ] All GPT-4 models work
- [ ] o1/o3 reasoning models work
- [ ] Audio model works

---

### Issue 6.2: Groq Integration

**Title:** Enable Groq ultra-fast inference

**Description:**
Enable Groq for ultra-fast inference with Llama and Mixtral models.

**Implementation:**
```python
from agno.models.groq import Groq

GROQ_MODELS = {
    "groq-llama-3.3-70b": Groq(id="llama-3.3-70b-versatile"),
    "groq-llama-3.1-8b": Groq(id="llama-3.1-8b-instant"),
    "groq-mixtral": Groq(id="mixtral-8x7b-32768"),
}
```

**Environment Variables:**
```bash
GROQ_API_KEY=gsk_...
```

**Acceptance Criteria:**
- [ ] Groq models work
- [ ] Ultra-fast inference (<100ms)
- [ ] All Groq models available

---

### Issue 6.3: Together AI Integration

**Title:** Enable Together AI for large open models

**Description:**
Enable Together AI for Llama 405B and other large models.

**Implementation:**
```python
from agno.models.together import Together

TOGETHER_MODELS = {
    "together-llama-3.1-405b": Together(id="meta-llama/Meta-Llama-3.1-405B-Instruct-Turbo"),
    "together-qwen-2.5-72b": Together(id="Qwen/Qwen2.5-72B-Instruct-Turbo"),
    "together-deepseek-v3": Together(id="deepseek-ai/DeepSeek-V3"),
}
```

**Acceptance Criteria:**
- [ ] Together models work
- [ ] Llama 405B available
- [ ] Qwen models available

---

### Issue 6.4: DeepSeek Integration

**Title:** Enable DeepSeek models

**Description:**
Enable DeepSeek chat and reasoning models.

**Implementation:**
```python
from agno.models.deepseek import DeepSeek

DEEPSEEK_MODELS = {
    "deepseek-chat": DeepSeek(id="deepseek-chat"),
    "deepseek-reasoner": DeepSeek(id="deepseek-reasoner"),
}
```

**Acceptance Criteria:**
- [ ] DeepSeek chat works
- [ ] DeepSeek reasoner works

---

## Epic 7: Multimodal & Voice (MEDIUM)

**Priority:** P2 - Medium
**Estimated Hours:** 20
**Dependencies:** Epic 1 (Agent factory), Issue 1.4 (Voice)

### Issue 7.1: Image Generation

**Title:** Enable DALL-E image generation

**Description:**
Enable image generation via DALL-E tools.

**Implementation:**
```python
from agno.tools.dalle import DalleTools

dalle = DalleTools(
    model="dall-e-3",
    size="1024x1024",
    quality="standard",
)
# Functions: generate_image(prompt)
```

**Acceptance Criteria:**
- [ ] DalleTools working
- [ ] Image generation via agent
- [ ] Size/quality configurable

---

### Issue 7.2: Vision Input

**Title:** Enable vision/image input for agents

**Description:**
Enable agents to accept and analyze images.

**Implementation:**
```python
from agno.agent import Agent
from agno.models.openai import OpenAIChat

# Vision-capable agent
agent = Agent(
    name="Vision Agent",
    model=OpenAIChat(id="gpt-4o"),  # Vision-capable model
)

# Run with image
response = await agent.arun(
    "What's in this image?",
    images=["path/to/image.jpg"],
)
```

**API Specification:**
```yaml
POST /agents/{agent_id}/runs
Request:
  message: string
  images: string[] # Base64 or URLs
Response:
  # Normal response
```

**Acceptance Criteria:**
- [ ] Image input works
- [ ] Multiple images supported
- [ ] Base64 and URL formats
- [ ] Vision models auto-selected

---

### Issue 7.3: Audio Input/Output

**Title:** Enable audio input and output

**Description:**
Enable audio transcription and synthesis for agents.

**Implementation:**
```python
from agno.models.openai import OpenAIChat

# Audio-capable agent
agent = Agent(
    name="Audio Agent",
    model=OpenAIChat(id="gpt-4o-audio-preview"),
)

# Run with audio
response = await agent.arun(
    audio=audio_bytes,  # Audio input
    output_audio=True,  # Request audio output
)
```

**Acceptance Criteria:**
- [ ] Audio input transcription
- [ ] Audio output synthesis
- [ ] Streaming audio support

---

### Issue 7.4: Video Analysis

**Title:** Enable video frame analysis

**Description:**
Enable agents to analyze video by extracting and processing frames.

**Implementation:**
```python
from agno.media.video import extract_frames

# Extract frames from video
frames = extract_frames("video.mp4", fps=1)

# Analyze with vision agent
response = await agent.arun(
    "Describe what happens in this video",
    images=frames,
)
```

**Acceptance Criteria:**
- [ ] Frame extraction from video
- [ ] Multi-frame analysis
- [ ] Video URL support

---

## Epic 8: Evals & Observability (MEDIUM)

**Priority:** P2 - Medium
**Estimated Hours:** 15
**Dependencies:** None

### Issue 8.1: AccuracyEval Implementation

**Title:** Implement Agno AccuracyEval system

**Description:**
Implement evaluation framework for measuring agent accuracy.

**Agno Pattern:**
```python
from agno.evals import AccuracyEval
from agno.evals.models import EvalResult

eval = AccuracyEval(
    agent=agent,
    test_cases=[
        {"input": "What is 2+2?", "expected": "4"},
        {"input": "Capital of France?", "expected": "Paris"},
    ],
    # Use agent as judge
    judge_model=OpenAIChat(id="gpt-4o"),
)

results: list[EvalResult] = await eval.run()
# Results include: accuracy, latency, token usage
```

**API Specification:**
```yaml
POST /evals
Request:
  agent_id: string
  test_cases:
    - input: string
      expected: string
Response:
  eval_id: string
  status: "running"

GET /evals/{id}
Response:
  eval_id: string
  status: "completed"
  results:
    accuracy: float
    total_cases: int
    passed: int
    failed: int
    avg_latency_ms: float
    cases:
      - input: string
        expected: string
        actual: string
        passed: boolean
```

**Acceptance Criteria:**
- [ ] Eval creation works
- [ ] Test cases run
- [ ] Results computed
- [ ] Judge model works

---

### Issue 8.2: Batch Evaluation

**Title:** Implement batch evaluation for multiple agents

**Description:**
Enable running evaluations across multiple agents for comparison.

**Implementation:**
```python
from agno.evals import BatchEval

batch = BatchEval(
    agents=[agent1, agent2, agent3],
    test_cases=test_cases,
)

comparison = await batch.run()
# Returns comparison of all agents on same test cases
```

**Acceptance Criteria:**
- [ ] Batch eval runs
- [ ] Multi-agent comparison
- [ ] Results aggregated

---

### Issue 8.3: Trace Export

**Title:** Implement trace export to observability platforms

**Description:**
Export traces to Phoenix, Langfuse, or other platforms.

**Implementation:**
```python
from agno.telemetry import PhoenixExporter, LangfuseExporter

# Phoenix export
phoenix = PhoenixExporter(
    endpoint="http://phoenix:6060",
)

# Langfuse export
langfuse = LangfuseExporter(
    public_key=os.getenv("LANGFUSE_PUBLIC_KEY"),
    secret_key=os.getenv("LANGFUSE_SECRET_KEY"),
)

# Attach to agent
agent = Agent(
    name="Traced Agent",
    telemetry=[phoenix, langfuse],
)
```

**Acceptance Criteria:**
- [ ] Phoenix export works
- [ ] Langfuse export works
- [ ] Full trace data included

---

### Issue 8.4: Metrics Dashboard

**Title:** Create agent metrics dashboard

**Description:**
Create dashboard showing agent performance metrics.

**Metrics to Track:**
- Requests per minute
- Average latency
- Token usage
- Error rate
- Memory usage
- Tool call frequency

**Acceptance Criteria:**
- [ ] Metrics collected
- [ ] Dashboard displays metrics
- [ ] Historical trends visible

---

## Sprint Execution Plan

### Week 1: Foundation (Epic 1)
| Day | Issues | Hours | Team |
|-----|--------|-------|------|
| 1 | 1.1 Dynamic Agent Creation | 8 | Backend |
| 2 | 1.2 Universal Model Selection | 6 | Backend |
| 3 | 1.3 Dynamic Tool Attachment | 6 | Backend |
| 4 | 1.4 Universal LiveKit | 6 | Backend |
| 5 | 1.5 Shared Knowledge, 1.6 Persistence | 8 | Backend |

### Week 2: Workflows (Epic 2)
| Day | Issues | Hours | Team |
|-----|--------|-------|------|
| 1 | 2.1 Dynamic Workflow Creation | 8 | Backend |
| 2 | 2.2 Workflow Step Types | 8 | Backend |
| 3 | 2.3 Workflow Persistence | 6 | Backend |
| 4 | 2.4 Workflow Templates | 6 | Backend |
| 5 | 2.5 Workflow HITL | 6 | Backend |

### Week 3: Tools & Memory (Epics 3, 4)
| Day | Issues | Hours | Team |
|-----|--------|-------|------|
| 1 | 3.1 Search Tools | 4 | Backend |
| 2 | 3.2-3.5 Remaining Tools | 8 | Backend |
| 3 | 4.1 User Memories | 6 | Backend |
| 4 | 4.2-4.3 Session Summaries, Agentic Memory | 6 | Backend |
| 5 | 4.4-4.5 Memory Search, UI | 8 | Full Stack |

### Week 4: Teams, Models, Multi-modal (Epics 5, 6, 7)
| Day | Issues | Hours | Team |
|-----|--------|-------|------|
| 1 | 5.1-5.2 Team Creation, Modes | 8 | Backend |
| 2 | 5.3-5.4 Team Persistence, Voice | 6 | Backend |
| 3 | 6.1-6.4 Model Providers | 8 | Backend |
| 4 | 7.1-7.2 Image Gen, Vision | 6 | Backend |
| 5 | 7.3-7.4 Audio, Video | 6 | Backend |

### Week 5: Evals & Polish (Epic 8 + Testing)
| Day | Issues | Hours | Team |
|-----|--------|-------|------|
| 1 | 8.1-8.2 AccuracyEval, Batch | 6 | Backend |
| 2 | 8.3-8.4 Traces, Dashboard | 6 | Backend |
| 3 | Integration Testing | 8 | QA |
| 4 | Bug Fixes | 8 | Full Stack |
| 5 | Documentation, Deployment | 6 | All |

---

## Files to Create Summary

### New Files (~4,500 lines)
```
core/
├── agent/
│   ├── factory.py          # 300 lines
│   ├── registry.py         # 100 lines
│   └── models.py           # 150 lines
├── models/
│   ├── registry.py         # 200 lines
│   └── providers/          # 300 lines total
├── tools/
│   ├── registry.py         # 150 lines
│   └── builtin/            # 400 lines total
├── workflow/
│   ├── base.py             # 200 lines
│   ├── steps.py            # 300 lines
│   ├── registry.py         # 100 lines
│   ├── executor.py         # 250 lines
│   └── models.py           # 150 lines
├── voice/
│   ├── service.py          # 200 lines
│   ├── livekit.py          # 150 lines
│   └── models.py           # 50 lines
├── knowledge/
│   ├── base.py             # 200 lines
│   └── types/              # 460 lines total
├── memory/
│   ├── search.py           # 150 lines
│   └── embeddings.py       # 100 lines
├── team/
│   ├── base.py             # 200 lines
│   ├── modes.py            # 200 lines
│   └── models.py           # 100 lines
└── evals/
    ├── accuracy.py         # 150 lines
    ├── batch.py            # 100 lines
    └── export.py           # 150 lines
```

### Database Schema Additions
```sql
-- agents table
-- teams table
-- team_members table
-- team_runs table
-- workflows table
-- workflow_runs table
-- workflow_step_runs table
-- memory_embeddings table
```

---

## Environment Variables Summary

```bash
# Required
DATABASE_URL=postgresql://...
REDIS_URL=redis://...
NATS_URL=nats://...

# Model Providers
OPENAI_API_KEY=sk-...
ANTHROPIC_API_KEY=sk-ant-...
GROQ_API_KEY=gsk_...
TOGETHER_API_KEY=...
DEEPSEEK_API_KEY=...

# Local Models
OLLAMA_HOST=http://localhost:11434

# Voice
LIVEKIT_URL=...
LIVEKIT_API_KEY=...
LIVEKIT_API_SECRET=...

# Tools (Optional)
EXA_API_KEY=...
E2B_API_KEY=...
GITHUB_TOKEN=...
SLACK_BOT_TOKEN=...
DISCORD_BOT_TOKEN=...
EMAIL_USER=...
EMAIL_PASSWORD=...

# Observability (Optional)
PHOENIX_ENDPOINT=...
LANGFUSE_PUBLIC_KEY=...
LANGFUSE_SECRET_KEY=...
```

---

## Epic 9: Advanced Features from Cookbook (NEW)

**Priority:** P1 High
**Estimated Hours:** 30
**Dependencies:** Epics 1-2 (foundation work)

> **Note:** These features were identified by reviewing the Agno cookbook examples on 2026-01-25. They represent advanced patterns that significantly enhance agent capabilities.

### Issue 9.1: ReasoningTools Integration

**Linear Issue:** OMOS-195

**Title:** Implement ReasoningTools for explicit structured thinking

**Description:**
Agno provides ReasoningTools that enable agents to perform explicit reasoning steps with `think()` and `analyze()` functions. This is critical for complex problem-solving tasks.

**Agno Cookbook Pattern:**
```python
from agno.agent import Agent
from agno.models.anthropic import Claude
from agno.tools.reasoning import ReasoningTools

# Agent with explicit reasoning capabilities
reasoning_agent = Agent(
    model=Claude(id="claude-sonnet-4-20250514"),
    tools=[
        ReasoningTools(
            think=True,           # Enable think() function
            analyze=True,         # Enable analyze() function
            add_instructions=True, # Add reasoning instructions
            add_few_shot=True,    # Include few-shot examples
        ),
    ],
    instructions="Use structured thinking for complex problems",
    markdown=True,
)

# Run with full reasoning visibility
reasoning_agent.print_response(
    "Analyze the pros and cons of microservices vs monolith",
    stream=True,
    show_full_reasoning=True,  # Show reasoning steps
)
```

**Extended Thinking with Reasoning Models:**
```python
from agno.agent import Agent
from agno.models.anthropic import Claude

# Create agent with reasoning model (extended thinking)
agent = Agent(
    reasoning_model=Claude(
        id="claude-sonnet-4-5",
        thinking={"type": "enabled", "budget_tokens": 1024},
    ),
    reasoning=True,
    instructions="Think step by step about the problem.",
)

# Stream reasoning events
agent.print_response(
    "What is 25 * 37? Show your reasoning.",
    stream=True,
    stream_events=True,  # See thinking process
)
```

**Files to Create:**
```
core/tools/reasoning.py    # ReasoningTools wrapper (~100 lines)
```

**Files to Modify:**
- `core/tools/registry.py` - Add reasoning tools
- `core/agent/factory.py` - Support reasoning_model parameter

**Acceptance Criteria:**
- [ ] ReasoningTools enabled with think/analyze
- [ ] show_full_reasoning parameter works
- [ ] Extended thinking with budget_tokens works
- [ ] Stream reasoning events working

---

### Issue 9.2: Router Workflow Step

**Linear Issue:** OMOS-196

**Title:** Implement Router workflow step for dynamic agent selection

**Description:**
Router enables dynamic routing to different agents or steps based on input analysis. Essential for building intelligent multi-path workflows.

**Agno Cookbook Pattern:**
```python
from agno.workflow import Router, Step, Workflow
from agno.workflow.types import StepInput
from typing import List

# Define specialist agents
tech_expert = Agent(name="Tech Expert", tools=[HackerNewsTools()])
biz_expert = Agent(name="Business Expert", tools=[DuckDuckGoTools()])
generalist = Agent(name="Generalist")

# Router selector function
def route_by_topic(step_input: StepInput) -> List[Step]:
    """Dynamically select which agent handles the request."""
    topic = step_input.input.lower()

    if "tech" in topic or "ai" in topic or "software" in topic:
        return [Step(name="Tech Research", agent=tech_expert)]
    elif "business" in topic or "market" in topic:
        return [Step(name="Business Research", agent=biz_expert)]
    else:
        return [Step(name="General Research", agent=generalist)]

# Workflow with Router
workflow = Workflow(
    name="Expert Routing Workflow",
    steps=[
        Router(
            name="Topic Router",
            selector=route_by_topic,
            choices=[tech_step, business_step, general_step],
            description="Routes to appropriate expert based on topic",
        ),
        Step(name="Synthesis", agent=synthesizer),
    ]
)

workflow.print_response("Latest AI developments", markdown=True)
```

**Files to Create:**
```
core/workflows/router.py   # Router step implementation (~150 lines)
```

**Files to Modify:**
- `core/workflows/engine.py` - Add Router step support
- `core/workflows/models.py` - Router Pydantic model

**Acceptance Criteria:**
- [ ] Router step type implemented
- [ ] Selector function receives StepInput
- [ ] Dynamic routing to different agents works
- [ ] Choices parameter defines valid routes

---

### Issue 9.3: Loop Workflow Step

**Linear Issue:** OMOS-197

**Title:** Implement Loop workflow step for iterative execution

**Description:**
Loop enables iterative execution with end conditions and max iterations. Critical for research, refinement, and quality assurance workflows.

**Agno Cookbook Pattern:**
```python
from agno.workflow import Loop, Parallel, Step, Workflow
from agno.workflow.types import StepInput, StepOutput

# Evaluator to check if research is complete
def research_evaluator(step_input: StepInput) -> bool:
    """Return True to stop the loop, False to continue."""
    content = step_input.previous_step_content or ""
    word_count = len(content.split())
    # Stop if we have enough content
    return word_count > 500

# Workflow with Loop
workflow = Workflow(
    name="Iterative Research Workflow",
    steps=[
        Loop(
            name="Research Loop",
            steps=[
                Parallel(
                    research_hackernews_step,
                    research_web_step,
                    name="Parallel Research",
                ),
                sentiment_analysis_step,
            ],
            end_condition=research_evaluator,  # Function to check completion
            max_iterations=3,                   # Safety limit
        ),
        content_step,  # Final synthesis
    ],
)

workflow.print_response("Research AI trends", stream=True)
```

**Files to Create:**
```
core/workflows/loop.py     # Loop step implementation (~150 lines)
```

**Files to Modify:**
- `core/workflows/engine.py` - Add Loop step support
- `core/workflows/models.py` - Loop Pydantic model

**Acceptance Criteria:**
- [ ] Loop step type implemented
- [ ] end_condition evaluator works
- [ ] max_iterations limit enforced
- [ ] Loop can contain Parallel/Steps

---

### Issue 9.4: Background Hooks with @hook Decorator

**Linear Issue:** OMOS-198

**Title:** Implement @hook decorator for background task execution

**Description:**
The @hook decorator allows hooks to run in the background without blocking API responses. Essential for logging, notifications, and analytics.

**Agno Cookbook Pattern:**
```python
from agno.hooks import hook
from agno.agent import Agent
from agno.run.agent import RunInput
import asyncio

# Synchronous hook (blocks response)
def validate_output(run_output, agent):
    """Runs synchronously, blocks response until complete."""
    if not run_output.content:
        raise ValueError("Empty response not allowed")

# Background hook (doesn't block response)
@hook(run_in_background=True)
async def send_notification(run_output, agent):
    """Runs in background after response is sent."""
    await asyncio.sleep(1)  # Simulate slow operation
    print(f"Notification sent for {agent.name}")

# Pre-hook in background
@hook(run_in_background=True)
def log_request(run_input: RunInput, agent):
    """Pre-hook runs in background. Cannot modify run_input."""
    print(f"Request received for {agent.name}: {run_input.input_content}")

# Agent with mixed hooks
agent = Agent(
    id="background-task-agent",
    name="BackgroundTaskAgent",
    model=OpenAIChat(id="gpt-4o-mini"),
    pre_hooks=[log_request],
    post_hooks=[validate_output, send_notification],
)

# AgentOS with background hooks enabled
from agno.os import AgentOS

agent_os = AgentOS(
    agents=[agent],
    run_hooks_in_background=True,  # Global setting
)
```

**Files to Create:**
```
core/hooks/decorator.py    # @hook decorator implementation (~100 lines)
core/hooks/background.py   # Background task runner (~100 lines)
```

**Files to Modify:**
- `core/orchestrator/main.py` - Support background hooks
- `core/agent/factory.py` - pre_hooks/post_hooks parameters

**Acceptance Criteria:**
- [ ] @hook decorator implemented
- [ ] run_in_background=True works
- [ ] Pre-hooks can't modify input when background
- [ ] Post-hooks don't block API response
- [ ] Global run_hooks_in_background setting works

---

### Issue 9.5: Session State Management

**Linear Issue:** OMOS-199

**Title:** Implement session_state for stateful agents

**Description:**
Session state allows agents to maintain state across turns without database persistence. Essential for shopping carts, form wizards, and interactive workflows.

**Agno Cookbook Pattern:**
```python
from agno.agent import Agent
from agno.db.sqlite import SqliteDb
from agno.tools import tool
from textwrap import dedent

# Tools that modify session state
@tool
def add_item(agent: Agent, item: str) -> str:
    """Add item to shopping list."""
    agent.session_state["shopping_list"].append(item)
    return f"Added {item} to list"

@tool
def remove_item(agent: Agent, item: str) -> str:
    """Remove item from shopping list."""
    if item in agent.session_state["shopping_list"]:
        agent.session_state["shopping_list"].remove(item)
        return f"Removed {item}"
    return f"{item} not in list"

@tool
def list_items(agent: Agent) -> str:
    """List all items."""
    items = agent.session_state["shopping_list"]
    return f"Shopping list: {', '.join(items) if items else 'empty'}"

# Agent with session state
agent = Agent(
    model=OpenAIChat(id="gpt-5-mini"),
    session_state={"shopping_list": []},  # Initial state
    db=SqliteDb(db_file="tmp/shopping.db"),
    tools=[add_item, remove_item, list_items],
    instructions=dedent("""
        Your job is to manage a shopping list.
        Current shopping list: {shopping_list}
    """),
)

# State persists across turns
agent.print_response("Add milk and eggs")
print(f"State: {agent.get_session_state()}")
# Output: {'shopping_list': ['milk', 'eggs']}

agent.print_response("Remove milk")
print(f"State: {agent.get_session_state()}")
# Output: {'shopping_list': ['eggs']}
```

**Files to Modify:**
- `core/agent/factory.py` - session_state parameter
- `core/agent/context.py` - State access in tools

**Acceptance Criteria:**
- [ ] session_state parameter works
- [ ] State persists across turns
- [ ] get_session_state() returns current state
- [ ] State accessible in tools via agent.session_state
- [ ] State can be used in instruction templates

---

### Issue 9.6: PerformanceEval and ReliabilityEval

**Linear Issue:** OMOS-200

**Title:** Implement PerformanceEval and ReliabilityEval systems

**Description:**
Complement AccuracyEval with performance benchmarking and tool reliability testing.

**Agno Cookbook Pattern:**
```python
from agno.eval.performance import PerformanceEval
from agno.eval.reliability import ReliabilityEval
from agno.agent import Agent
from agno.db.sqlite import SqliteDb

# Performance evaluation
def run_agent():
    agent = Agent(
        model=OpenAIChat(id="gpt-5-mini"),
        add_history_to_context=True,
        db=SqliteDb(db_file="tmp/perf.db"),
    )
    response = agent.run("What is the capital of France?")
    return response.content

perf_eval = PerformanceEval(
    name="Response Time Benchmark",
    func=run_agent,
    num_iterations=10,
    warmup_runs=2,
)
perf_eval.run(print_results=True, print_summary=True)
# Output: avg_time, p50, p95, p99, memory_usage

# Reliability evaluation (tool call success rate)
reliability_eval = ReliabilityEval(
    name="Tool Reliability Check",
    agent=calculator_agent,
    test_cases=[
        {"input": "Calculate 2+2", "expected_tool": "add"},
        {"input": "Calculate 10/2", "expected_tool": "divide"},
    ],
    num_iterations=5,
)
reliability_eval.run(print_results=True)
# Output: success_rate, tool_call_accuracy, failure_reasons
```

**Files to Create:**
```
core/evals/performance.py   # PerformanceEval (~150 lines)
core/evals/reliability.py   # ReliabilityEval (~150 lines)
```

**API Endpoints:**
```yaml
GET /eval-runs/performance
Response:
  runs:
    - id: string
      avg_time_ms: float
      p50_ms: float
      p95_ms: float
      memory_mb: float

GET /eval-runs/reliability
Response:
  runs:
    - id: string
      success_rate: float
      tool_accuracy: float
      failures: []
```

**Acceptance Criteria:**
- [ ] PerformanceEval measures timing
- [ ] warmup_runs excluded from stats
- [ ] p50/p95/p99 percentiles calculated
- [ ] ReliabilityEval tracks tool calls
- [ ] API endpoints for both eval types

---

### Issue 9.7: MemoryTools for Explicit Memory Operations

**Linear Issue:** OMOS-201

**Title:** Implement MemoryTools for agent-controlled memory

**Description:**
MemoryTools provide explicit tools for agents to store and retrieve memories, giving them control over what to remember.

**Agno Cookbook Pattern:**
```python
from agno.agent import Agent
from agno.db.sqlite import SqliteDb
from agno.tools.memory import MemoryTools
from agno.models.openai import OpenAIChat

# Database for memory storage
db = SqliteDb(db_file="tmp/memory.db")

# Memory tools instance
memory_tools = MemoryTools(
    db=db,
    search_memories=True,    # Enable semantic search
    add_memory=True,         # Enable adding memories
    delete_memory=True,      # Enable deletion
    update_memory=True,      # Enable updates
)

agent = Agent(
    model=OpenAIChat(id="gpt-5-mini"),
    tools=[memory_tools],
    markdown=True,
)

# Agent decides what to remember
agent.print_response(
    "My name is John Doe and I like hiking on weekends. "
    "I'm planning to travel to Africa in December.",
    user_id="john_doe@example.com",
    stream=True
)

# Agent retrieves memories
agent.print_response(
    "What have you remembered about me?",
    user_id="john_doe@example.com",
    stream=True
)
```

**Files to Create:**
```
core/tools/memory_tools.py  # MemoryTools wrapper (~150 lines)
```

**Acceptance Criteria:**
- [ ] MemoryTools integrated with database
- [ ] add_memory tool works
- [ ] search_memories with semantic search
- [ ] user_id scoping for memories
- [ ] delete_memory and update_memory work

---

### Issue 9.8: RunContext for Tool Access

**Linear Issue:** OMOS-202

**Title:** Implement RunContext for context-aware tools

**Description:**
RunContext provides tools with access to the current execution context, including user_id, session_id, and dependencies.

**Agno Cookbook Pattern:**
```python
from agno.agent import Agent
from agno.run import RunContext
from agno.tools import tool

@tool
def get_user_profile(run_context: RunContext) -> str:
    """Get the current user's profile."""
    user_id = run_context.user_id
    profiles = run_context.dependencies.get("user_profiles", {})
    return profiles.get(user_id, {"name": "Unknown"})

@tool
def check_permissions(run_context: RunContext, action: str) -> bool:
    """Check if user has permission for action."""
    user_id = run_context.user_id
    permissions = run_context.dependencies.get("permissions", {})
    user_perms = permissions.get(user_id, [])
    return action in user_perms

agent = Agent(
    model=OpenAIChat(id="gpt-5-mini"),
    tools=[get_user_profile, check_permissions],
    dependencies={
        "user_profiles": {
            "user_1001": {"name": "John", "level": "senior"},
            "user_1002": {"name": "Jane", "level": "junior"},
        },
        "permissions": {
            "user_1001": ["read", "write", "admin"],
            "user_1002": ["read"],
        }
    },
)

# user_id available in RunContext
agent.print_response(
    "Get my profile and check if I can write",
    user_id="user_1001",
    stream=True
)
```

**Files to Modify:**
- `core/agent/context.py` - RunContext implementation
- `core/tools/base.py` - RunContext injection

**Acceptance Criteria:**
- [ ] RunContext available in tools
- [ ] user_id accessible
- [ ] session_id accessible
- [ ] dependencies accessible
- [ ] Type hints work with RunContext

---

## Updated Issue Summary

| Issue | Title | Epic | Priority |
|-------|-------|------|----------|
| OMOS-195 | Implement ReasoningTools for explicit structured thinking | Epic 9 | P1 |
| OMOS-196 | Implement Router workflow step for dynamic agent selection | Epic 9 | P1 |
| OMOS-197 | Implement Loop workflow step for iterative execution | Epic 9 | P1 |
| OMOS-198 | Implement @hook decorator for background task execution | Epic 9 | P1 |
| OMOS-199 | Implement session_state for stateful agents | Epic 9 | P1 |
| OMOS-200 | Implement PerformanceEval and ReliabilityEval systems | Epic 9 | P2 |
| OMOS-201 | Implement MemoryTools for agent-controlled memory | Epic 9 | P1 |
| OMOS-202 | Implement RunContext for context-aware tools | Epic 9 | P1 |

**Total New Issues:** 8 issues in Epic 9
**Updated Total:** 45 issues across 9 epics (OMOS-158 to OMOS-202)

---

## Success Metrics

| Metric | Current | Target |
|--------|---------|--------|
| Features Working | 22% (34/154) | 90%+ (138/154) |
| Dynamic Agent Creation | No | Yes |
| Dynamic Workflow Creation | No | Yes |
| Built-in Tools Enabled | 3/19 | 19/19 |
| Model Providers | 2 | 6 |
| Universal Voice | No | Yes |
| Memory Features | 1/4 | 4/4 |
| Team Features | 1/4 | 4/4 |
| Eval System | Partial | Complete |

---

## References

- **Agno Documentation:** https://docs.agno.com
- **Context7 Library:** `/anthropic/agno`
- **Current Feature Matrix:** `docs/FEATURE-MATRIX.md`
- **Linear Project:** OneMind OS Backlog

---

*Document Version: 1.0*
*Last Updated: 2026-01-25*

# Agno Cognitive Layer — The Brain

> **Note**: This document describes the **Agno/Python worker layer** — stateless agents dispatched via NATS. These are NOT the OpenClaw domain agents (Legacy, Spartan, Oracle, Eden, Empire, Haven, Guardian, Grid, Trinity). For the 9-agent OpenClaw architecture, see [AGENT-ROSTER.md](../02.3%20Agents/AGENT-ROSTER.md).

> Handles all reasoning, memory, knowledge retrieval, digital tool execution, and multi-agent orchestration. This is where Legacy's mind lives.

## Role in the Architecture

Agno is the cognitive center of One Mind OS. It receives text inputs from the Interaction Gateway, reasons about them using LLM #1, retrieves relevant knowledge, accesses long-term memory, coordinates agent teams, executes digital tools, and sends physical commands to OM1 via NATS.

## LLM Configuration

Agno uses LLM #1 — a powerful, cloud-based model optimized for deep reasoning. This model handles all cognitive tasks but is NOT used for real-time physical control (that's OM1's LLM #2).

See [LLM-STRATEGY.md](./LLM-STRATEGY.md) for the full two-tier LLM architecture.

```python
from agno.models.openai import OpenAIChat

# Primary reasoning model
model = OpenAIChat(id="gpt-4.1")

# Alternative options per agent:
# OpenAIChat(id="gpt-4.1-mini")  — faster, cheaper for simple tasks
# Claude(id="claude-sonnet-4-5-20250929")  — strong reasoning
# Gemini(id="gemini-2.5-pro")  — multimodal
```

## Knowledge System (RAG)

Agno's Agentic RAG lets agents decide WHEN to search their knowledge base, not on every request. Content is chunked, embedded, and stored in PgVector.

```python
from agno.knowledge.pdf_url import PDFUrlKnowledgeBase
from agno.knowledge.combined import CombinedKnowledgeBase
from agno.knowledge.json import JSONKnowledgeBase
from agno.vectordb.pgvector import PgVector, SearchType

# Vector store
vector_db = PgVector(
    table_name="onemind_knowledge",
    db_url="postgresql+psycopg://ai:ai@localhost:5532/onemind",
    search_type=SearchType.hybrid,  # vector + keyword (BM25)
)

# Knowledge sources
knowledge = CombinedKnowledgeBase(
    sources=[
        PDFUrlKnowledgeBase(urls=["https://example.com/legacy-codex.pdf"]),
        JSONKnowledgeBase(path="data/user_preferences.json"),
    ],
    vector_db=vector_db,
)
```

### What to put in the knowledge base

- Legacy Codex (HP, LE, GE, IT pillars)
- Business documents (Legacy Platform Holdings LLC)
- Technical documentation (robot manuals, API docs)
- Personal reference materials
- Family coordination data

## Memory System

Three layers of memory, all backed by PostgreSQL:

1. **Session Memory** — Chat history for the current conversation. Automatically maintained.
2. **Persistent Session Storage** — Session state saved across restarts. Enables conversation continuity.
3. **User Memories** — Long-term facts extracted from conversations ("user prefers morning meetings"). Persists forever, recalled automatically.

```python
from agno.db.postgres import PostgresDb

db = PostgresDb(
    db_url="postgresql+psycopg://ai:ai@localhost:5532/onemind",
    table_name="agent_sessions",
)

agent = Agent(
    model=model,
    db=db,
    enable_user_memories=True,
    enable_agentic_memory=True,
    # Memory v2: background memory manager
    memory=AgentMemory(
        db=db,
        create_user_memories=True,
        update_user_memories_after_run=True,
    ),
)
```

## Agent Teams

One Mind OS uses Agno's **Coordinate** team mode. A team leader receives requests, delegates to specialist agents, and synthesizes results.

### Team Structure

```python
from agno.agent import Agent
from agno.team.team import Team

# ── Physical Agent (bridges to OM1 via NATS) ──
physical_agent = Agent(
    name="Physical",
    role="Handle all physical-world interactions via the robot body",
    model=OpenAIChat(id="gpt-4.1"),
    tools=[RobotControlTools()],  # custom NATS toolkit
    instructions=[
        "You control a physical robot. Use robot_control tools for movement.",
        "Always check sensor data before commanding movement.",
        "Report physical observations back to the team.",
    ],
)

# ── Digital Agent (handles all digital tasks) ──
digital_agent = Agent(
    name="Digital",
    role="Handle email, calendar, chat, web, databases, APIs",
    model=OpenAIChat(id="gpt-4.1"),
    tools=[
        EmailTools(),
        # Add: GmailTools, SlackTools, ClickUpTools, etc. via MCP
    ],
    instructions=["Handle all digital communications and web tasks."],
)

# ── HP Agent (Holistic Performance) ──
hp_agent = Agent(
    name="HP",
    role="Health tracking, fitness, nutrition, sleep, wellness",
    model=OpenAIChat(id="gpt-4.1"),
    tools=[HealthTools()],
    instructions=["Monitor and advise on health and performance metrics."],
)

# ── LE Agent (Legacy Evolution) ──
le_agent = Agent(
    name="LE",
    role="Personal development, learning, family, life optimization",
    model=OpenAIChat(id="gpt-4.1"),
    instructions=["Guide personal growth and family coordination."],
)

# ── GE Agent (Generational Endeavor) ──
ge_agent = Agent(
    name="GE",
    role="Business strategy, revenue, content, brand",
    model=OpenAIChat(id="gpt-4.1"),
    tools=[BusinessTools()],
    instructions=["Drive business growth and entrepreneurship strategy."],
)

# ── UI Agent (Innovative Technology) ──
it_agent = Agent(
    name="IT",
    role="Technical architecture, code, DevOps, infrastructure",
    model=OpenAIChat(id="gpt-4.1"),
    tools=[CodeTools()],
    instructions=["Handle all technical and infrastructure tasks."],
)

# ── The Unified Team ──
one_mind = Team(
    name="OneMind",
    mode="coordinate",
    model=OpenAIChat(id="gpt-4.1"),
    members=[physical_agent, digital_agent, hp_agent, le_agent, ge_agent, it_agent],
    knowledge=knowledge,
    db=db,
    enable_user_memories=True,
    enable_agentic_memory=True,
    instructions=[
        "You are Legacy, a unified AI entity with physical and digital capabilities.",
        "You have a robot body AND digital tools. Delegate accordingly.",
        "Physical tasks → Physical agent. Digital tasks → Digital agent.",
        "Health → HP. Business → GE. Personal → LE. Technical → IT.",
        "Maintain one consistent personality across all interactions.",
    ],
)
```

## MCP Integration

Agno connects to external services via Model Context Protocol. Each MCP server provides tools that agents can invoke.

```python
from agno.tools.mcp import MCPTools

# Gmail MCP
gmail_tools = MCPTools(
    server_url="https://mcp.gmail.com/sse",
    transport="sse",
)

# ClickUp MCP
clickup_tools = MCPTools(
    server_url="https://mcp.clickup.com/mcp",
    transport="streamable-http",
)

# Agent with MCP tools
digital_agent = Agent(
    name="Digital",
    tools=[gmail_tools, clickup_tools],
    # ...
)
```

## Robot Control Toolkit

Custom Agno toolkit that bridges to OM1 via NATS.

```python
from agno.tools.toolkit import Toolkit
import nats, json

class RobotControlTools(Toolkit):
    def __init__(self, nats_url: str = "nats://localhost:4222"):
        super().__init__(name="robot_control")
        self.nats_url = nats_url
        self.register(self.move_robot)
        self.register(self.speak_aloud)
        self.register(self.get_sensor_data)
        self.register(self.get_fleet_status)

    async def move_robot(self, robot_id: str, action: str, parameters: str = "") -> str:
        """Command a robot. Actions: walk, run, sit, wave, patrol, go_to, charge."""
        nc = await nats.connect(self.nats_url)
        payload = json.dumps({"action": action, "params": parameters})
        await nc.publish(f"onemind.commands.{robot_id}", payload.encode())
        await nc.drain()
        return f"Command sent to {robot_id}: {action} {parameters}"

    async def speak_aloud(self, robot_id: str, text: str) -> str:
        """Make the robot speak out loud via its speaker."""
        nc = await nats.connect(self.nats_url)
        payload = json.dumps({"action": "speak", "text": text})
        await nc.publish(f"onemind.commands.{robot_id}", payload.encode())
        await nc.drain()
        return f"Robot {robot_id} speaking: {text}"

    async def get_sensor_data(self, robot_id: str) -> str:
        """Get current sensor readings from a specific robot."""
        nc = await nats.connect(self.nats_url)
        response = await nc.request(
            f"onemind.sensors.{robot_id}", b"query", timeout=3.0
        )
        await nc.drain()
        return response.data.decode()

    async def get_fleet_status(self) -> str:
        """Get status of all robots in the fleet."""
        nc = await nats.connect(self.nats_url)
        response = await nc.request(
            "onemind.fleet.status", b"query", timeout=3.0
        )
        await nc.drain()
        return response.data.decode()
```

## AgentOS Deployment

Agno agents are served via AgentOS — a FastAPI runtime with REST endpoints and SSE streaming.

```python
from agno.agent_os import AgentOS

agent_os = AgentOS(
    description="One Mind OS — Legacy",
    teams=[one_mind],
    agents=[physical_agent, digital_agent, hp_agent, le_agent, ge_agent, it_agent],
)

# Endpoints:
# POST /teams/onemind/runs  — run the unified team
# POST /agents/physical/runs  — run a specific agent
# GET  /teams/onemind/sessions  — list sessions
# GET  /health  — health check
```

## File Structure

```
agno_layer/
├── main.py              # AgentOS startup
├── agents/
│   ├── physical.py      # Physical agent + RobotControlTools
│   ├── digital.py       # Digital agent + MCP tools
│   ├── hp.py            # Holistic Performance agent
│   ├── le.py            # Legacy Evolution agent
│   ├── ge.py            # Generational Endeavor agent
│   └── it.py            # Innovative Technology agent
├── teams/
│   └── onemind.py       # Unified team configuration
├── tools/
│   ├── robot_control.py # NATS-based robot toolkit
│   └── mcp_registry.py  # MCP server connections
├── knowledge/
│   └── setup.py         # Knowledge base configuration
├── memory/
│   └── setup.py         # Memory/storage configuration
└── config.py            # Models, DB URLs, NATS URL
```

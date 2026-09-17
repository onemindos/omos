# OneMind OS Multi-Agent Architecture

> **The Paradigm Shift: From Legacy-Centric to True Multi-Agent System**
>
> Date: 2026-01-24
> Status: Architecture Vision

---

## Executive Summary

OneMind OS is evolving from a **single-agent system** (Legacy as the "Mother AI") to a **true multi-agent system** where:

1. **Every agent is autonomous** - Each agent has its own model, knowledge base, tools, and memory
2. **Agents share memory** - Letta-style memory blocks enable hive-mind capabilities
3. **Agents are multimodal** - Images, audio, video, files as first-class inputs/outputs
4. **Physical + Digital agents** - Same agent can control robots AND file Linear issues
5. **Dynamic agent creation** - Agent Registry pattern for on-demand agent instantiation

**Goal:** Letta's brain structure (memory blocks, shared memory) with Agno's performance (multi-model, teams, workflows).

---

## Table of Contents

1. [The Problem: Legacy-Centric Bottleneck](#1-the-problem-legacy-centric-bottleneck)
2. [The Solution: Agent Registry Pattern](#2-the-solution-agent-registry-pattern)
3. [Memory Architecture: Letta-Style Blocks](#3-memory-architecture-letta-style-blocks)
4. [Per-Agent Configuration](#4-per-agent-configuration)
5. [Multimodal Agents](#5-multimodal-agents)
6. [Physical Agents: Robot Control](#6-physical-agents-robot-control)
7. [Teams & Coordination](#7-teams--coordination)
8. [Implementation Plan](#8-implementation-plan)
9. [Database Schema](#9-database-schema)
10. [Code Examples](#10-code-examples)

---

## 1. The Problem: Legacy-Centric Bottleneck

### Current Architecture (Legacy-Centric)

```
                    ┌─────────────────────┐
                    │       LEGACY        │
                    │   (Mother AI)       │
                    │                     │
                    │  - ALL tools        │
                    │  - ALL knowledge    │
                    │  - ONE model        │
                    │  - ALL memory       │
                    └─────────┬───────────┘
                              │
         ┌────────────────────┼────────────────────┐
         │                    │                    │
    ┌────▼────┐         ┌─────▼────┐        ┌─────▼────┐
    │ Todoist │         │  GitHub  │        │  Comms   │
    └─────────┘         └──────────┘        └──────────┘
```

**Problems:**
- **Single point of failure** - Everything routes through Legacy
- **No specialization** - One agent tries to do everything
- **Fixed model** - Can't use faster/cheaper models for simple tasks
- **Monolithic knowledge** - All docs in one RAG index
- **No parallelism** - Can't have multiple agents working simultaneously
- **Physical limitation** - Can't connect to robots/cameras/speakers

### Target Architecture (Multi-Agent)

```
                         ┌─────────────────────┐
                         │   AGENT REGISTRY    │
                         │   (Factory)         │
                         └─────────┬───────────┘
                                   │
    ┌──────────────┬───────────────┼───────────────┬──────────────┐
    │              │               │               │              │
┌───▼───┐    ┌─────▼─────┐   ┌─────▼─────┐  ┌─────▼─────┐  ┌─────▼─────┐
│Legacy │    │   Bob     │   │  DevOps   │  │    UI     │  │ Research  │
│Mother │    │ Physical  │   │  Agent    │  │  Agent    │  │  Agent    │
│  AI   │    │  Agent    │   │           │  │           │  │           │
├───────┤    ├───────────┤   ├───────────┤  ├───────────┤  ├───────────┤
│Opus4.5│    │ Gemini2.0 │   │ Sonnet3.5│  │   Haiku   │  │   Opus    │
│       │    │ +LiveKit  │   │           │  │           │  │           │
├───────┤    ├───────────┤   ├───────────┤  ├───────────┤  ├───────────┤
│All    │    │Robot Ctrl │   │ Infra    │  │ Flutter  │  │ Web       │
│Tools  │    │Camera     │   │ Docker   │  │ Figma    │  │ Search    │
│       │    │Speaker    │   │ Bash     │  │ CSS      │  │ RAG       │
├───────┤    ├───────────┤   ├───────────┤  ├───────────┤  ├───────────┤
│System │    │Robot      │   │ Ops      │  │ Design   │  │ Docs      │
│Knowl. │    │Manuals    │   │ Playbook │  │ System   │  │ Papers    │
└───────┘    └───────────┘   └───────────┘  └───────────┘  └───────────┘
    │              │               │               │              │
    └──────────────┴───────────────┴───────────────┴──────────────┘
                                   │
                    ┌──────────────▼──────────────┐
                    │     SHARED MEMORY BLOCKS     │
                    │  (Letta-Style Hive Mind)    │
                    │                              │
                    │  - Zeus Profile (ALL)        │
                    │  - Four Pillars (Pillar)     │
                    │  - Project Context (Team)    │
                    │  - Agent Persona (Individual)│
                    └──────────────────────────────┘
```

---

## 2. The Solution: Agent Registry Pattern

### Core Concept

Instead of one hardcoded Legacy agent, we have an **Agent Registry** that can create any agent on-demand with its own configuration.

```python
# core/agent/registry.py

from typing import Dict, Any, Optional, List
from agno.agent import Agent
from agno.models.aws import AwsBedrock
from agno.models.google import Gemini
from agno.models.ollama import Ollama
from pydantic import BaseModel


class AgentConfig(BaseModel):
    """Configuration for a single agent."""
    agent_id: str
    name: str
    description: str

    # Model configuration
    model_provider: str = "bedrock"  # bedrock, gemini, ollama, litellm
    model_id: str = "us.anthropic.claude-opus-4-5-20251101-v1:0"

    # Knowledge configuration
    knowledge_table: Optional[str] = None  # LanceDB table name
    knowledge_topics: List[str] = []

    # Tools configuration
    tool_modules: List[str] = []  # ["todoist", "github", "robot_control"]

    # Memory configuration
    memory_blocks: List[str] = []  # ["human", "persona", "pillar_it"]
    enable_archival_memory: bool = True

    # Multimodal configuration
    enable_vision: bool = False
    enable_audio: bool = False
    enable_video: bool = False

    # Physical agent configuration
    physical_enabled: bool = False
    livekit_room: Optional[str] = None
    remote_agent_location: Optional[str] = None


# Agent Registry - Stores configurations for all agents
AGENT_REGISTRY: Dict[str, AgentConfig] = {

    # Legacy - Mother AI (Coordinator)
    "legacy": AgentConfig(
        agent_id="legacy",
        name="Legacy",
        description="Mother AI - Coordinates all agents and systems",
        model_provider="bedrock",
        model_id="us.anthropic.claude-opus-4-5-20251101-v1:0",
        knowledge_table="knowledge_system",
        knowledge_topics=["onemind", "infrastructure", "four_pillars"],
        tool_modules=["todoist", "github", "comms", "execute_action", "browser", "inbox"],
        memory_blocks=["human", "persona_legacy", "system_context", "pillar_hp", "pillar_le", "pillar_ge", "pillar_it"],
        enable_archival_memory=True,
        enable_vision=True,
    ),

    # Bob - Physical Agent (Robot + Digital Work)
    "bob": AgentConfig(
        agent_id="bob",
        name="Bob",
        description="Physical agent - Controls robots and handles physical tasks",
        model_provider="gemini",
        model_id="gemini-2.0-flash-exp",  # Fast, multimodal
        knowledge_table="knowledge_robots",
        knowledge_topics=["unitree", "robotics", "navigation", "manipulation"],
        tool_modules=["robot_control", "camera", "speaker", "linear", "todoist"],
        memory_blocks=["human", "persona_bob", "robot_context"],
        enable_archival_memory=True,
        enable_vision=True,
        enable_audio=True,
        enable_video=True,
        physical_enabled=True,
        livekit_room="bob-control",
        remote_agent_location="farm",
    ),

    # DevOps Agent - Infrastructure Specialist
    "devops": AgentConfig(
        agent_id="devops",
        name="DevOps Agent",
        description="Infrastructure, deployments, and operations",
        model_provider="bedrock",
        model_id="us.anthropic.claude-sonnet-4-20250514-v1:0",  # Good balance
        knowledge_table="knowledge_devops",
        knowledge_topics=["docker", "kubernetes", "aws", "terraform"],
        tool_modules=["bash", "docker", "ssh", "github"],
        memory_blocks=["human", "persona_devops", "pillar_it", "project_onemind"],
        enable_archival_memory=True,
    ),

    # UI Agent - Frontend Specialist
    "ui": AgentConfig(
        agent_id="ui",
        name="UI Agent",
        description="Flutter and web frontend development",
        model_provider="bedrock",
        model_id="us.anthropic.claude-haiku-4-20250514-v1:0",  # Fast for simple tasks
        knowledge_table="knowledge_ui",
        knowledge_topics=["flutter", "dart", "react", "nextjs", "tailwind"],
        tool_modules=["file_edit", "github"],
        memory_blocks=["human", "persona_ui", "project_onemind"],
        enable_archival_memory=True,
    ),

    # Research Agent - Deep Analysis
    "research": AgentConfig(
        agent_id="research",
        name="Research Agent",
        description="Deep research, analysis, and documentation",
        model_provider="bedrock",
        model_id="us.anthropic.claude-opus-4-5-20251101-v1:0",  # Full power for research
        knowledge_table="knowledge_research",
        knowledge_topics=["papers", "documentation", "analysis"],
        tool_modules=["web_search", "browser", "file_edit"],
        memory_blocks=["human", "persona_research"],
        enable_archival_memory=True,
    ),

    # Edge Agent - Offline/Local Operations
    "edge": AgentConfig(
        agent_id="edge",
        name="Edge Agent",
        description="Offline operations at edge locations",
        model_provider="ollama",
        model_id="phi3:mini",  # Lightweight for edge
        knowledge_table="knowledge_local",
        knowledge_topics=["home_assistant", "local_operations"],
        tool_modules=["homeassistant", "local_shell"],
        memory_blocks=["human", "persona_edge", "location_context"],
        enable_archival_memory=False,  # No cloud DB access
    ),
}


async def create_agent(agent_id: str, session_id: Optional[str] = None) -> Agent:
    """
    Create an agent instance from the registry.

    This is the core factory function that:
    1. Loads agent configuration from registry
    2. Creates the appropriate model
    3. Loads agent-specific knowledge
    4. Loads agent-specific tools
    5. Loads shared + individual memory blocks
    6. Returns a fully configured Agent
    """
    if agent_id not in AGENT_REGISTRY:
        raise ValueError(f"Unknown agent: {agent_id}")

    config = AGENT_REGISTRY[agent_id]

    # 1. Create model
    model = create_model_from_config(config)

    # 2. Load knowledge
    knowledge = load_knowledge(config.knowledge_table, config.knowledge_topics)

    # 3. Load tools
    tools = load_tools(config.tool_modules)

    # 4. Load memory blocks
    memory_context = await load_memory_blocks(config.memory_blocks, agent_id)

    # 5. Build system prompt with memory
    system_prompt = build_agent_prompt(config, memory_context)

    # 6. Create agent
    agent = Agent(
        name=config.name,
        id=config.agent_id,
        model=model,
        tools=tools,
        knowledge=knowledge,
        search_knowledge=knowledge is not None,
        instructions=system_prompt,
        description=config.description,
        session_id=session_id,
        user_id="zeus",
        markdown=True,
        show_tool_calls=True,
        enable_agentic_memory=config.enable_archival_memory,
    )

    return agent


def create_model_from_config(config: AgentConfig):
    """Create model instance based on configuration."""
    if config.model_provider == "bedrock":
        return AwsBedrock(
            id=config.model_id,
            aws_region="us-east-1",
        )
    elif config.model_provider == "gemini":
        return Gemini(id=config.model_id)
    elif config.model_provider == "ollama":
        return Ollama(id=config.model_id)
    else:
        raise ValueError(f"Unknown model provider: {config.model_provider}")
```

### API Integration

```python
# core/agent/api.py

from fastapi import APIRouter
from core.agent.registry import AGENT_REGISTRY, create_agent

router = APIRouter()


@router.get("/agents")
async def list_agents():
    """List all available agents with their configurations."""
    return {
        "agents": [
            {
                "id": config.agent_id,
                "name": config.name,
                "description": config.description,
                "model": config.model_id,
                "physical": config.physical_enabled,
                "multimodal": {
                    "vision": config.enable_vision,
                    "audio": config.enable_audio,
                    "video": config.enable_video,
                },
            }
            for config in AGENT_REGISTRY.values()
        ]
    }


@router.post("/agents/{agent_id}/runs")
async def run_agent(
    agent_id: str,
    message: str,
    session_id: Optional[str] = None,
    files: Optional[List[UploadFile]] = None,
):
    """
    Run any agent from the registry.

    The agent will be created on-demand with its full configuration:
    - Specific model (Opus/Sonnet/Haiku/Gemini/Phi-3)
    - Specific knowledge base
    - Specific tools
    - Shared + individual memory blocks
    """
    agent = await create_agent(agent_id, session_id)

    # Handle multimodal inputs
    if files:
        message = await prepare_multimodal_message(message, files, agent_id)

    response = await agent.arun(message)
    return {"content": response.content}
```

---

## 3. Memory Architecture: Letta-Style Blocks

### Memory Block Hierarchy

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         MEMORY BLOCK HIERARCHY                           │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                           │
│  GLOBAL BLOCKS (shared_with: ['*'])                                      │
│  ├─ "human" - Zeus's profile, preferences, schedule                     │
│  └─ "system_context" - OneMind infrastructure, guardrails               │
│                                                                           │
│  PILLAR BLOCKS (shared_with: pillar agents)                              │
│  ├─ "pillar_hp" - Health, fitness, habits [HP agents]                   │
│  ├─ "pillar_le" - Family, home, legacy [LE agents]                      │
│  ├─ "pillar_ge" - Business, revenue [GE agents]                         │
│  └─ "pillar_it" - Infrastructure, AI [Grid agents]                        │
│                                                                           │
│  TEAM BLOCKS (shared_with: team members)                                 │
│  ├─ "project_onemind" - Sprint, issues, decisions [Dev team]            │
│  ├─ "project_farm" - Farm operations [Farm team]                        │
│  └─ "handoff_*" - Team handoff context [Temporary]                      │
│                                                                           │
│  INDIVIDUAL BLOCKS (shared_with: single agent)                           │
│  ├─ "persona_legacy" - Legacy's identity [Legacy only]                  │
│  ├─ "persona_bob" - Bob's identity [Bob only]                           │
│  ├─ "persona_devops" - DevOps identity [DevOps only]                    │
│  └─ "robot_context" - Current robot state [Bob only]                    │
│                                                                           │
│  ARCHIVAL MEMORY (searchable via embeddings)                             │
│  ├─ Past decisions, learnings                                            │
│  ├─ Workflow execution summaries                                         │
│  ├─ Tagged by pillar (HP, LE, GE, IT)                                   │
│  └─ Filtered by agent/team access                                        │
│                                                                           │
└─────────────────────────────────────────────────────────────────────────┘
```

### Memory Loading Flow

```python
# core/memory/blocks.py

async def load_memory_blocks(block_labels: List[str], agent_id: str) -> str:
    """
    Load memory blocks for an agent.

    This implements Letta's core memory pattern:
    - Blocks are always loaded into context (no tool call needed)
    - Blocks can be shared across agents (hive mind)
    - Blocks respect access control (shared_with_agents)
    """

    # Query blocks this agent has access to
    blocks = await db.fetch_all("""
        SELECT label, value, pillar, char_limit
        FROM ai.user_memories
        WHERE memory_type = 'core'
          AND always_in_context = TRUE
          AND user_id = 'zeus'
          AND (
              '*' = ANY(shared_with_agents)  -- Global blocks
              OR $1 = ANY(shared_with_agents)  -- Agent-specific
              OR label = ANY($2)  -- Explicitly requested
          )
        ORDER BY
            CASE
                WHEN '*' = ANY(shared_with_agents) THEN 0  -- Global first
                WHEN label LIKE 'pillar_%' THEN 1          -- Then pillars
                WHEN label LIKE 'project_%' THEN 2         -- Then projects
                ELSE 3                                      -- Then individual
            END,
            label
    """, agent_id, block_labels)

    # Format blocks for context injection
    memory_context = "## CORE MEMORY (Always Present)\n\n"

    for block in blocks:
        pillar_tag = f" [{block['pillar']}]" if block.get('pillar') else ""
        memory_context += f"### {block['label']}{pillar_tag}\n"
        memory_context += f"{block['value']}\n\n"

    return memory_context


async def get_agent_persona(agent_id: str) -> str:
    """Get agent-specific persona block."""

    block = await db.fetch_one("""
        SELECT value FROM ai.user_memories
        WHERE memory_type = 'core'
          AND label = $1
          AND user_id = 'zeus'
    """, f"persona_{agent_id}")

    return block['value'] if block else ""
```

### Memory Tools (Letta-Compatible)

```python
# core/memory/tools.py

from agno.tools import toolkit


@toolkit
class MemoryTools:
    """
    Letta-compatible memory tools.

    4 core operations:
    - core_memory_append: Add to a core block
    - core_memory_replace: Update text in a core block
    - archival_memory_insert: Store fact in archival (with embedding)
    - archival_memory_search: Semantic search archival memory
    """

    def __init__(self, agent_id: str):
        self.agent_id = agent_id
        self.embedding_service = EmbeddingService()

    def core_memory_append(
        self,
        label: str,
        content: str,
        pillar: Optional[str] = None
    ) -> str:
        """
        Append content to a core memory block.

        Args:
            label: Block label (e.g., 'human', 'pillar_hp')
            content: Text to append
            pillar: Optional pillar tag (HP, LE, GE, IT)

        Example:
            core_memory_append(
                label="pillar_hp",
                content="Started new morning yoga routine - 15 min daily",
                pillar="HP"
            )
        """
        # Implementation: Append to block, respecting char_limit
        ...

    def core_memory_replace(
        self,
        label: str,
        old_content: str,
        new_content: str
    ) -> str:
        """
        Replace specific text in a core memory block.

        Args:
            label: Block label
            old_content: Text to find (exact match)
            new_content: Replacement text

        Example:
            core_memory_replace(
                label="human",
                old_content="Wake time: 6:00 AM",
                new_content="Wake time: 5:30 AM"
            )
        """
        # Implementation: Find and replace in block
        ...

    def archival_memory_insert(
        self,
        content: str,
        tags: Optional[List[str]] = None,
        pillar: Optional[str] = None
    ) -> str:
        """
        Store a fact in archival memory with semantic search.

        Args:
            content: The information to store
            tags: Optional tags for filtering
            pillar: Optional pillar (HP, LE, GE, IT)

        Example:
            archival_memory_insert(
                content="Bob successfully navigated barn patrol route. Took 12 minutes. No obstacles detected.",
                tags=["robot:bob", "task:patrol", "location:barn"],
                pillar="LE"
            )
        """
        # Generate embedding and store
        embedding = self.embedding_service.generate(content)

        await db.execute("""
            INSERT INTO ai.user_memories (
                memory_id, user_id, memory_type, memory, embedding, topics, pillar,
                created_by_agent
            ) VALUES (
                $1, 'zeus', 'archival',
                jsonb_build_object('content', $2),
                $3, $4, $5, $6
            )
        """, memory_id, content, embedding, tags, pillar, self.agent_id)

        return f"Stored in archival memory (ID: {memory_id})"

    def archival_memory_search(
        self,
        query: str,
        pillar: Optional[str] = None,
        tags: Optional[List[str]] = None,
        top_k: int = 5
    ) -> str:
        """
        Search archival memory using semantic similarity.

        Args:
            query: Search query
            pillar: Filter by pillar
            tags: Filter by tags
            top_k: Number of results

        Example:
            archival_memory_search(
                query="What happened during barn patrols?",
                pillar="LE",
                tags=["robot:bob"],
                top_k=5
            )
        """
        # Generate query embedding and search
        query_embedding = self.embedding_service.generate(query)

        results = await db.fetch_all("""
            SELECT
                memory->>'content' as content,
                topics,
                pillar,
                created_at,
                created_by_agent,
                1 - (embedding <=> $1::vector) AS similarity
            FROM ai.user_memories
            WHERE user_id = 'zeus'
              AND memory_type = 'archival'
              AND ($2::varchar IS NULL OR pillar = $2)
              AND ($3::text[] IS NULL OR topics && $3)
            ORDER BY embedding <=> $1::vector
            LIMIT $4
        """, query_embedding, pillar, tags, top_k)

        # Format results
        ...
```

---

## 4. Per-Agent Configuration

### Model Selection Strategy

| Agent | Model | Reason |
|-------|-------|--------|
| **Legacy** | Claude Opus 4.5 | Full reasoning power for coordination |
| **Bob (Physical)** | Gemini 2.0 Flash | Fast multimodal for real-time robot control |
| **DevOps** | Claude Sonnet 4 | Good balance for infrastructure tasks |
| **UI** | Claude Haiku 4 | Fast for simple frontend tasks |
| **Research** | Claude Opus 4.5 | Deep analysis requires full power |
| **Edge** | Phi-3 Mini | Lightweight for offline/local operation |

### Knowledge Base Per Agent

Each agent gets its own LanceDB table for specialized knowledge:

```python
# core/knowledge/loader.py

KNOWLEDGE_TABLES = {
    "knowledge_system": {
        "description": "OneMind OS system documentation",
        "sources": [
            "docs/architecture/*.md",
            "docs/operations/*.md",
            "CLAUDE.md",
        ],
    },
    "knowledge_robots": {
        "description": "Robotics manuals and procedures",
        "sources": [
            "docs/robots/unitree/*.md",
            "docs/robots/dji/*.md",
            "hardware/manuals/*.pdf",
        ],
    },
    "knowledge_devops": {
        "description": "Infrastructure and DevOps documentation",
        "sources": [
            "infra/docs/*.md",
            "docs/operations/*.md",
            ".github/workflows/*.yml",
        ],
    },
    "knowledge_ui": {
        "description": "Frontend design and implementation",
        "sources": [
            "ui/docs/*.md",
            "mobile/docs/*.md",
            "docs/ui/*.md",
        ],
    },
}


async def load_knowledge(table_name: str, topics: List[str]) -> Knowledge:
    """Load knowledge base for an agent."""

    if table_name not in KNOWLEDGE_TABLES:
        return None

    config = KNOWLEDGE_TABLES[table_name]

    return LanceDbKnowledge(
        table_name=table_name,
        embedder=OpenAIEmbedder(model="text-embedding-3-small"),
        search_type=SearchType.hybrid,  # Vector + BM25
        topics=topics,
    )
```

### Tool Modules Per Agent

```python
# core/tools/loader.py

TOOL_MODULES = {
    # Communication & Task Management
    "todoist": todoist_tools,
    "github": github_tools,
    "comms": comms_tools,
    "linear": linear_tools,
    "inbox": inbox_tools,

    # Browser & Web
    "browser": browser_tools,
    "web_search": web_search_tools,

    # File Operations
    "file_edit": file_edit_tools,

    # Infrastructure
    "bash": bash_tools,
    "docker": docker_tools,
    "ssh": ssh_tools,

    # Robot Control (Physical Agents)
    "robot_control": RobotControlTools,  # Movement, navigation
    "camera": CameraTools,               # Vision, streaming
    "speaker": SpeakerTools,             # Audio output

    # Edge/Local
    "homeassistant": homeassistant_tools,
    "local_shell": local_shell_tools,

    # Universal Gateway
    "execute_action": execute_action,  # 138+ actions
}


def load_tools(module_names: List[str]) -> List:
    """Load tools for an agent."""
    tools = []
    for name in module_names:
        if name in TOOL_MODULES:
            tool = TOOL_MODULES[name]
            if callable(tool):
                tools.append(tool())
            else:
                tools.append(tool)
    return tools
```

---

## 5. Multimodal Agents

### Enabling Multimodal I/O

Every agent can be configured for multimodal inputs:

```python
# core/agent/multimodal.py

async def prepare_multimodal_message(
    text: str,
    files: List[UploadFile],
    agent_id: str
) -> Union[str, List]:
    """
    Prepare multimodal message for agent.

    Supports:
    - Images (PNG, JPG, WebP)
    - Audio (MP3, WAV, WebM)
    - Video (MP4, WebM) - extracted as frames
    - Documents (PDF, TXT, MD)
    """
    config = AGENT_REGISTRY[agent_id]

    message_parts = [{"type": "text", "text": text}]

    for file in files:
        content_type = file.content_type

        if content_type.startswith("image/") and config.enable_vision:
            # Add image to context
            image_data = await file.read()
            message_parts.append({
                "type": "image",
                "source": {
                    "type": "base64",
                    "media_type": content_type,
                    "data": base64.b64encode(image_data).decode()
                }
            })

        elif content_type.startswith("audio/") and config.enable_audio:
            # Transcribe audio and add
            transcript = await transcribe_audio(file)
            message_parts.append({
                "type": "text",
                "text": f"[Audio Transcript]: {transcript}"
            })

        elif content_type.startswith("video/") and config.enable_video:
            # Extract frames and add
            frames = await extract_video_frames(file, max_frames=5)
            for i, frame in enumerate(frames):
                message_parts.append({
                    "type": "image",
                    "source": {
                        "type": "base64",
                        "media_type": "image/jpeg",
                        "data": frame
                    }
                })

    return message_parts


# Usage: Bob analyzing what he sees
bob = await create_agent("bob")
response = await bob.arun([
    {"type": "text", "text": "What do you see? Any obstacles in the path?"},
    {"type": "image", "source": {"type": "base64", "data": camera_frame}}
])
```

### Video Streaming with LiveKit

For physical agents, we integrate LiveKit for real-time video:

```python
# core/agent/physical.py

class PhysicalAgentMixin:
    """Mixin for physical agent capabilities."""

    async def connect_to_robot(self, room_name: str):
        """Connect to robot via LiveKit."""
        self.livekit_client = LiveKitClient(
            url=os.getenv("LIVEKIT_URL"),
            api_key=os.getenv("LIVEKIT_API_KEY"),
            api_secret=os.getenv("LIVEKIT_API_SECRET"),
        )

        await self.livekit_client.join_room(room_name)

        # Subscribe to robot's video streams
        self.video_tracks = {}
        async for track in self.livekit_client.video_tracks():
            self.video_tracks[track.name] = track

    async def get_camera_frame(self, camera_name: str = "main") -> bytes:
        """Get current frame from robot camera."""
        if camera_name in self.video_tracks:
            return await self.video_tracks[camera_name].get_frame()
        return None

    async def speak(self, text: str):
        """Speak through robot's speaker."""
        audio = await text_to_speech(text)
        await self.livekit_client.publish_audio(audio)

    async def send_command(self, command: dict):
        """Send control command to robot."""
        await self.livekit_client.send_data(
            json.dumps(command).encode(),
            topic="robot_control"
        )
```

---

## 6. Physical Agents: Robot Control

### Bob: The Physical + Digital Agent

Bob is the flagship example of a hybrid agent that can:
- Control physical robots (Unitree, drones)
- See through robot cameras
- Speak through robot speakers
- File Linear issues
- Update Todoist tasks

```python
# core/agents/bob.py

from agno.agent import Agent
from agno.tools import toolkit
from core.agent.physical import PhysicalAgentMixin


@toolkit
class RobotControlTools:
    """Robot control tools for physical agents."""

    def __init__(self, remote_agent_url: str, location: str):
        self.remote_agent_url = remote_agent_url
        self.location = location

    async def move_robot(
        self,
        direction: str,
        distance_meters: float,
        speed: str = "normal"
    ) -> str:
        """
        Move the robot in a direction.

        Args:
            direction: forward, backward, left, right
            distance_meters: Distance to move
            speed: slow, normal, fast

        Example:
            move_robot(direction="forward", distance_meters=2.0, speed="slow")
        """
        command = {
            "type": "move",
            "direction": direction,
            "distance": distance_meters,
            "speed": speed,
        }

        result = await self.send_robot_command(command)
        return f"Robot moved {direction} {distance_meters}m. Status: {result['status']}"

    async def patrol_route(
        self,
        route_name: str,
        announce: bool = True
    ) -> str:
        """
        Execute a predefined patrol route.

        Args:
            route_name: Name of saved route (e.g., "barn_perimeter", "garden_check")
            announce: Whether to announce presence via speaker

        Example:
            patrol_route(route_name="barn_perimeter", announce=True)
        """
        command = {
            "type": "patrol",
            "route": route_name,
            "announce": announce,
        }

        result = await self.send_robot_command(command)

        # Store patrol results in archival memory
        await self.memory_tools.archival_memory_insert(
            content=f"Patrol '{route_name}' completed. Duration: {result['duration']}s. "
                    f"Anomalies: {result.get('anomalies', 'None')}",
            tags=["robot:bob", f"patrol:{route_name}", f"location:{self.location}"],
            pillar="LE"
        )

        return f"Patrol complete. Duration: {result['duration']}s"

    async def look_at(
        self,
        target: str,
        analyze: bool = True
    ) -> str:
        """
        Point camera at target and optionally analyze.

        Args:
            target: What to look at (e.g., "gate", "water_trough", coordinates)
            analyze: Whether to analyze what's seen

        Example:
            look_at(target="water_trough", analyze=True)
        """
        # Move camera
        command = {"type": "look", "target": target}
        await self.send_robot_command(command)

        if analyze:
            # Get camera frame
            frame = await self.get_camera_frame()

            # Analyze with vision
            analysis = await self.analyze_frame(frame, f"Describe what you see at {target}")

            return f"Looking at {target}. Analysis: {analysis}"

        return f"Camera pointed at {target}"

    async def speak_announcement(
        self,
        message: str,
        volume: str = "normal"
    ) -> str:
        """
        Speak through robot's speaker.

        Args:
            message: What to say
            volume: quiet, normal, loud

        Example:
            speak_announcement(message="Starting perimeter patrol", volume="normal")
        """
        await self.speak(message)
        return f"Announced: {message}"

    async def send_robot_command(self, command: dict) -> dict:
        """Send command to robot via RemoteAgent."""
        async with aiohttp.ClientSession() as session:
            async with session.post(
                f"{self.remote_agent_url}/command",
                json=command
            ) as resp:
                return await resp.json()


# Bob's configuration in registry
BOB_CONFIG = AgentConfig(
    agent_id="bob",
    name="Bob",
    description="Physical agent - Controls robots, cameras, speakers while also handling digital tasks",

    # Fast multimodal model for real-time control
    model_provider="gemini",
    model_id="gemini-2.0-flash-exp",

    # Robot-specific knowledge
    knowledge_table="knowledge_robots",
    knowledge_topics=["unitree", "navigation", "patrol_routes", "farm_layout"],

    # Physical + Digital tools
    tool_modules=[
        "robot_control",  # Movement, patrol
        "camera",         # Vision
        "speaker",        # Audio output
        "linear",         # Digital: Issue management
        "todoist",        # Digital: Task management
        "comms",          # Digital: Notifications
    ],

    # Memory blocks
    memory_blocks=[
        "human",          # Zeus's profile (shared)
        "persona_bob",    # Bob's identity
        "robot_context",  # Current robot state
        "pillar_le",      # Farm/home context
    ],

    # Full multimodal
    enable_vision=True,
    enable_audio=True,
    enable_video=True,

    # Physical configuration
    physical_enabled=True,
    livekit_room="bob-control",
    remote_agent_location="farm",
)
```

### Bob's Persona Block

```sql
-- Bob's identity in core memory
INSERT INTO ai.user_memories (
    user_id, memory_type, label, value, always_in_context,
    shared_with_agents, char_limit
) VALUES (
    'zeus',
    'core',
    'persona_bob',
    'I am Bob, Zeus''s physical operations agent.

My capabilities:
- Control Unitree Go2 robot at the farm
- See through robot cameras (main, depth, thermal)
- Speak through robot speakers
- Navigate predefined patrol routes
- Detect anomalies and report issues

My responsibilities:
- Daily barn perimeter patrols (morning, evening)
- Check water troughs and feeding stations
- Monitor fence lines for damage
- Alert Zeus to any security concerns
- File Linear issues for maintenance needs

Physical location: Farm (Haven homestead)
Robot: Unitree Go2 Pro
Connectivity: 5G + LoRa mesh backup

I am both physical AND digital - I can patrol the barn AND update task lists.',
    TRUE,
    ARRAY['bob'],
    3000
);
```

### Example Bob Interaction

```python
# User tells Bob to do both physical and digital tasks

bob = await create_agent("bob")

response = await bob.arun("""
Do the evening barn patrol and then:
1. Check the water trough levels
2. If anything needs attention, create a Linear issue
3. Update the Todoist task "Evening patrol" as complete
""")

# Bob's actions:
# 1. Executes patrol_route("barn_perimeter")
# 2. Uses look_at("water_trough", analyze=True)
# 3. If low water detected: linear.create_issue("Water trough needs refill")
# 4. todoist.complete_task("Evening patrol")
# 5. Stores patrol summary in archival memory
```

---

## 7. Teams & Coordination

### Team Structure with Shared Memory

```python
# core/teams/factory.py

from agno.team import Team


async def create_team(
    team_id: str,
    member_agent_ids: List[str],
    routing_mode: str = "auto"
) -> Team:
    """
    Create a team with shared memory blocks.

    Teams enable:
    - Automatic task routing to best agent
    - Shared memory context
    - Handoff with context preservation
    """

    # Load shared blocks for this team
    team_blocks = await load_team_memory_blocks(team_id)

    # Create member agents
    members = []
    for agent_id in member_agent_ids:
        agent = await create_agent(agent_id)
        members.append(agent)

    # Format shared context
    shared_context = format_memory_blocks(team_blocks)

    team = Team(
        name=f"Team {team_id}",
        id=team_id,
        agents=members,
        routing_mode=routing_mode,
        share_member_interactions=True,  # Agents can see each other's work
        instructions=f"""You are a coordinated team.

## Shared Context (All Members)
{shared_context}

## Coordination Rules
1. Route to the agent best suited for the task
2. Share learnings via archival_memory_insert
3. Use handoffs for complex multi-step work
4. Update shared memory blocks as state changes
"""
    )

    return team


# Example: OneMind Dev Team
dev_team = await create_team(
    team_id="onemind-dev",
    member_agent_ids=["legacy", "devops", "ui", "research"],
    routing_mode="auto"
)

# Team automatically routes:
# - "Deploy the new build" → DevOps agent
# - "Fix the Flutter navigation bug" → UI agent
# - "Research best practices for X" → Research agent
# - "Coordinate the sprint" → Legacy (Mother AI)
```

### Handoffs with Memory Transfer

```python
# core/teams/handoff.py

async def handoff_task(
    from_agent_id: str,
    to_agent_id: str,
    task: str,
    context: str
) -> dict:
    """
    Transfer task between agents with full context.

    Creates a temporary handoff block visible to receiving agent.
    """

    handoff_id = f"handoff_{uuid.uuid4().hex[:8]}"

    # Create handoff memory block
    await db.execute("""
        INSERT INTO ai.user_memories (
            user_id, memory_type, label, value, always_in_context,
            shared_with_agents, char_limit
        ) VALUES (
            'zeus', 'core', $1, $2, TRUE, $3, 2000
        )
    """, handoff_id, f"""
HANDOFF: {from_agent_id} → {to_agent_id}

Task: {task}

Context from {from_agent_id}:
{context}

Previous work completed. Continue from here.
""", [to_agent_id])

    # Also store in archival for history
    await archival_memory_insert(
        content=f"Handoff from {from_agent_id} to {to_agent_id}: {task}. Context: {context[:200]}",
        tags=["handoff", f"from:{from_agent_id}", f"to:{to_agent_id}"],
    )

    return {"handoff_id": handoff_id, "to_agent": to_agent_id}
```

---

## 8. Implementation Plan

### Phase 1: Agent Registry (Week 1)

| Task | Status | Files |
|------|--------|-------|
| Create AgentConfig model | TODO | `core/agent/registry.py` |
| Define agent configurations | TODO | `core/agent/registry.py` |
| Implement create_agent factory | TODO | `core/agent/registry.py` |
| Add model creation helpers | TODO | `core/agent/models.py` |
| API endpoints for agents | TODO | `core/agent/api.py` |

### Phase 2: Memory Blocks (Week 2)

| Task | Status | Files |
|------|--------|-------|
| Database schema updates | TODO | `infra/configs/postgres/init-db.sql` |
| Memory block loader | TODO | `core/memory/blocks.py` |
| Implement 4 memory tools | TODO | `core/memory/tools.py` |
| Embedding service | TODO | `core/memory/embeddings.py` |
| Test memory operations | TODO | `tests/memory/` |

### Phase 3: Knowledge Per Agent (Week 3)

| Task | Status | Files |
|------|--------|-------|
| Define knowledge tables | TODO | `core/knowledge/tables.py` |
| LanceDB integration | TODO | `core/knowledge/lancedb.py` |
| Knowledge loader | TODO | `core/knowledge/loader.py` |
| Index existing docs | TODO | `scripts/index_knowledge.py` |

### Phase 4: Physical Agents (Week 4)

| Task | Status | Files |
|------|--------|-------|
| Robot control tools | TODO | `core/tools/robot.py` |
| LiveKit integration | TODO | `core/agent/physical.py` |
| Camera/speaker tools | TODO | `core/tools/multimedia.py` |
| Bob agent configuration | TODO | `core/agents/bob.py` |
| RemoteAgent connection | TODO | `core/agent/remote.py` |

### Phase 5: Teams & Workflows (Week 5)

| Task | Status | Files |
|------|--------|-------|
| Team factory | TODO | `core/teams/factory.py` |
| Handoff system | TODO | `core/teams/handoff.py` |
| Team memory loading | TODO | `core/teams/memory.py` |
| Workflow memory integration | TODO | `core/workflows/memory.py` |

### Phase 6: UI Integration (Week 6)

| Task | Status | Files |
|------|--------|-------|
| Agent selector in chat | TODO | `ui/lib/features/chat/` |
| Team selector | TODO | `ui/lib/features/chat/` |
| Physical agent controls | TODO | `ui/lib/features/robots/` |
| Memory viewer | TODO | `ui/lib/features/memory/` |

---

## 9. Database Schema

### Core Memory Blocks Table

```sql
-- Extend existing ai.user_memories for Letta-style blocks
ALTER TABLE ai.user_memories
ADD COLUMN IF NOT EXISTS memory_type VARCHAR(20) DEFAULT 'archival',
ADD COLUMN IF NOT EXISTS label VARCHAR(50),
ADD COLUMN IF NOT EXISTS embedding VECTOR(1536),
ADD COLUMN IF NOT EXISTS always_in_context BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS char_limit INTEGER DEFAULT 2000,
ADD COLUMN IF NOT EXISTS shared_with_agents TEXT[] DEFAULT ARRAY['*'],
ADD COLUMN IF NOT EXISTS shared_with_teams TEXT[],
ADD COLUMN IF NOT EXISTS pillar VARCHAR(10),
ADD COLUMN IF NOT EXISTS created_by_agent VARCHAR(50);

-- Indexes for efficient querying
CREATE INDEX IF NOT EXISTS idx_memories_type ON ai.user_memories (memory_type);
CREATE INDEX IF NOT EXISTS idx_memories_label ON ai.user_memories (label);
CREATE INDEX IF NOT EXISTS idx_memories_pillar ON ai.user_memories (pillar) WHERE pillar IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_memories_agents ON ai.user_memories USING gin (shared_with_agents);
CREATE INDEX IF NOT EXISTS idx_memories_teams ON ai.user_memories USING gin (shared_with_teams);
CREATE INDEX IF NOT EXISTS idx_memories_embedding ON ai.user_memories USING ivfflat (embedding vector_cosine_ops);
```

### Agent Configuration Table

```sql
-- Store agent configurations (optional, can also be in code)
CREATE TABLE IF NOT EXISTS ai.agent_configs (
    agent_id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    model_provider VARCHAR(20) DEFAULT 'bedrock',
    model_id VARCHAR(100),
    knowledge_table VARCHAR(50),
    knowledge_topics TEXT[],
    tool_modules TEXT[],
    memory_blocks TEXT[],
    enable_vision BOOLEAN DEFAULT FALSE,
    enable_audio BOOLEAN DEFAULT FALSE,
    enable_video BOOLEAN DEFAULT FALSE,
    physical_enabled BOOLEAN DEFAULT FALSE,
    livekit_room VARCHAR(100),
    remote_agent_location VARCHAR(50),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

### Knowledge Tables (LanceDB)

```python
# Each agent gets its own LanceDB table
# Schema defined in code, stored in vector DB

KNOWLEDGE_SCHEMA = {
    "id": str,
    "content": str,
    "embedding": List[float],  # 1536 dimensions
    "metadata": {
        "source": str,
        "topic": str,
        "agent_access": List[str],
        "created_at": str,
    }
}
```

---

## 10. Code Examples

### Creating an Agent

```python
from core.agent.registry import create_agent

# Create Legacy (Mother AI)
legacy = await create_agent("legacy")
response = await legacy.arun("What's on my schedule today?")

# Create Bob (Physical Agent)
bob = await create_agent("bob")
response = await bob.arun("Do the morning patrol and check water levels")

# Create DevOps Agent
devops = await create_agent("devops")
response = await devops.arun("Deploy the latest build to production")
```

### Creating a Team

```python
from core.teams.factory import create_team

# Create dev team
team = await create_team(
    team_id="onemind-dev",
    member_agent_ids=["legacy", "devops", "ui"],
    routing_mode="auto"
)

# Team automatically routes to best agent
response = await team.arun("Fix the login bug in Flutter")
# → Routes to UI agent

response = await team.arun("Deploy the fix")
# → Routes to DevOps agent
```

### Using Memory Tools

```python
from core.memory.tools import MemoryTools

memory = MemoryTools(agent_id="bob")

# Store a fact
memory.archival_memory_insert(
    content="Barn door hinge needs replacement. Noticed during morning patrol.",
    tags=["maintenance", "barn", "urgent"],
    pillar="LE"
)

# Search memories
results = memory.archival_memory_search(
    query="barn maintenance issues",
    pillar="LE",
    top_k=5
)

# Update core memory
memory.core_memory_append(
    label="pillar_le",
    content="Barn door hinge replacement scheduled for Friday",
    pillar="LE"
)
```

### Physical Agent with LiveKit

```python
from core.agents.bob import create_bob_agent

bob = await create_bob_agent()

# Connect to robot
await bob.connect_to_robot("bob-control")

# Execute physical task
response = await bob.arun([
    {"type": "text", "text": "What do you see? Any obstacles?"},
    {"type": "image", "source": {"type": "base64", "data": await bob.get_camera_frame()}}
])

# Speak through robot
await bob.speak("All clear. Starting patrol.")

# Execute patrol
response = await bob.arun("Execute barn perimeter patrol")
```

---

## Summary

The multi-agent architecture transforms OneMind OS from a single-agent bottleneck to a scalable, specialized system where:

1. **Any agent can be created** with its own model, knowledge, tools, and memory
2. **Memory is shared** through Letta-style blocks (global, pillar, team, individual)
3. **Agents can be physical** - controlling robots while handling digital tasks
4. **Teams coordinate** with automatic routing and handoffs
5. **Everything scales** - from lightweight edge agents to full-power research agents

**Key Files:**
- `core/agent/registry.py` - Agent configurations and factory
- `core/memory/blocks.py` - Memory block loading
- `core/memory/tools.py` - 4 Letta-compatible memory tools
- `core/agent/physical.py` - Physical agent capabilities
- `core/teams/factory.py` - Team creation with shared memory

**Reference Docs:**
- `docs/architecture/LETTA-AGNO-TEAMS-WORKFLOWS.md` - Detailed memory architecture
- `docs/architecture/LETTA-RESEARCH-FINDINGS.md` - Letta implementation research

---

**Ready to build the future of multi-agent systems.**

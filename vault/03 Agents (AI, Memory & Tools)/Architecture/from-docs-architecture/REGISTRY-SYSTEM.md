# Registry System Architecture

## Overview

The OneMind OS registry system is a multi-layered architecture that provides centralized discovery and management of agents, teams, workflows, and tools. It acts as the "process registry" for the Agent Operating System, built on top of the Agno framework.

## Table of Contents

- [Registry Layers](#registry-layers)
- [How Agno Fits In](#how-agno-fits-in)
- [Global Registry (backend/shared/registry.py)](#global-registry-backendsharedregistrypy)
- [Specialized Registries](#specialized-registries)
- [No Conflicts: How They Work Together](#no-conflicts-how-they-work-together)
- [Usage Examples](#usage-examples)
- [Best Practices](#best-practices)

---

## Registry Layers

OneMind OS uses a three-layer registry architecture:

```
┌─────────────────────────────────────────────────────────┐
│  LAYER 1: Global Discovery                              │
│  backend/shared/registry.py                             │
│  Purpose: System-wide agent/team/workflow lookup        │
│  Stores: Agent/Team/Workflow INSTANCES + metadata       │
└──────────────────┬──────────────────────────────────────┘
                   │ uses
                   ▼
┌─────────────────────────────────────────────────────────┐
│  LAYER 2: Specialized Registries                        │
│  - backend/agno/workflows/workflow/registry.py          │
│    (Workflow execution state + run history)             │
│  - backend/agno/tools/registry.py                       │
│    (Tool catalog + availability checking)               │
└──────────────────┬──────────────────────────────────────┘
                   │ wraps
                   ▼
┌─────────────────────────────────────────────────────────┐
│  LAYER 3: Agno Framework (pip package)                  │
│  - agno.agent.Agent (base class)                        │
│  - agno.team.Team (base class)                          │
│  - agno.workflow.Workflow (base class)                  │
│  NO built-in registry - just provides classes           │
└─────────────────────────────────────────────────────────┘
```

### Layer 1: Global Discovery
- **Location**: `backend/shared/registry.py`
- **Purpose**: Single source of truth for what exists in the system
- **Stores**: Instances of agents, teams, and workflows with metadata
- **Functionality**: Lookup, filtering by category/tags, system-wide discovery

### Layer 2: Specialized Registries
- **Workflow Registry**: Manages workflow execution state and run history
- **Tool Registry**: Manages tool availability, authentication, and metadata
- **Purpose**: Domain-specific concerns beyond simple lookup

### Layer 3: Agno Framework
- **Type**: External pip package (`agno==2.3.26`)
- **Provides**: Base classes and orchestration logic
- **Does NOT provide**: Any registry functionality
- **Role**: Execution engine, not discovery mechanism

---

## How Agno Fits In

### Agno is an External Framework

```python
# Installed via pip
pip install agno

# Provides base classes
from agno.agent import Agent
from agno.team import Team
from agno.workflow import Workflow
```

**Key Point**: Agno provides the **primitives** (Agent, Team, Workflow classes), but OneMind OS adds the **registry layer** on top for discovery and management.

### Relationship Summary

| Component | Type | Purpose |
|-----------|------|---------|
| **Agno Framework** | External package | Agent execution engine |
| **OneMind OS Registry** | Internal layer | Agent discovery system |

**Analogy**:
- **Agno** = Process execution engine (like how Linux executes processes)
- **Registry** = Process registry/service manager (like systemd or launchd)

---

## Global Registry (backend/shared/registry.py)

### Purpose

Central discovery and management of agents, teams, and workflows. This is the **SINGLE source of truth** for what's available in the system.

### Data Models

```python
@dataclass
class AgentInfo:
    """Information about a registered agent."""
    id: str
    name: str
    description: str
    agent: Agent  # The Agno Agent instance
    category: str = "general"
    tags: List[str] = field(default_factory=list)

@dataclass
class TeamInfo:
    """Information about a registered team."""
    id: str
    name: str
    description: str
    team: Team  # The Agno Team instance
    members: List[str] = field(default_factory=list)
    tags: List[str] = field(default_factory=list)

@dataclass
class WorkflowInfo:
    """Information about a registered workflow."""
    id: str
    name: str
    description: str
    workflow: Workflow  # The Agno Workflow instance
    tags: List[str] = field(default_factory=list)
```

### Core Functions

#### Agent Registry

```python
# Register an agent
register_agent(
    agent_id: str,
    agent: Agent,
    description: str = "",
    category: str = "general",
    tags: List[str] = None,
) -> None

# Get an agent by ID
get_agent(agent_id: str) -> Optional[Agent]

# Get full agent info
get_agent_info(agent_id: str) -> Optional[AgentInfo]

# List all agents (with optional filtering)
list_agents(category: str = None, tag: str = None) -> List[AgentInfo]

# Remove an agent
unregister_agent(agent_id: str) -> None
```

#### Team Registry

```python
# Register a team
register_team(
    team_id: str,
    team: Team,
    description: str = "",
    tags: List[str] = None,
) -> None

# Get a team by ID
get_team(team_id: str) -> Optional[Team]

# Get full team info
get_team_info(team_id: str) -> Optional[TeamInfo]

# List all teams
list_teams(tag: str = None) -> List[TeamInfo]

# Remove a team
unregister_team(team_id: str) -> None
```

#### Workflow Registry

```python
# Register a workflow
register_workflow(
    workflow_id: str,
    workflow: Workflow,
    description: str = "",
    tags: List[str] = None,
) -> None

# Get a workflow by ID
get_workflow(workflow_id: str) -> Optional[Workflow]

# Get full workflow info
get_workflow_info(workflow_id: str) -> Optional[WorkflowInfo]

# List all workflows
list_workflows(tag: str = None) -> List[WorkflowInfo]

# Remove a workflow
unregister_workflow(workflow_id: str) -> None
```

#### System-Wide Operations

```python
# Register multiple items at once
register_all(
    agents: Dict[str, Agent] = None,
    teams: Dict[str, Team] = None,
    workflows: Dict[str, Workflow] = None,
) -> None

# Get system-wide information
get_system_info() -> Dict[str, Any]

# Clear all registries (for testing)
clear_all() -> None
```

### Auto-Registration

The global registry supports auto-registration on import:

```python
def _auto_register_legacy() -> None:
    """Auto-register the Legacy agent and workflows."""
    try:
        from backend.orchestration.orchestrator.workflows import Legacy_WORKFLOWS

        for workflow in Legacy_WORKFLOWS:
            register_workflow(
                workflow_id=workflow.id,
                workflow=workflow,
                description=workflow.description,
                tags=["legacy", "productivity"],
            )

        logger.info("Auto-registered Legacy workflows")
    except ImportError as e:
        logger.warning(f"Could not auto-register Legacy: {e}")

# Auto-register on import
_auto_register_legacy()
```

---

## Specialized Registries

### Workflow Registry (backend/agno/workflows/workflow/registry.py)

**Purpose**: Manages workflow execution state and run history.

**Key Features**:
- Singleton pattern for global access
- Thread-safe with RLock
- Stores workflow definitions and runs
- Tracks run history per workflow

```python
class WorkflowRegistry:
    """In-memory registry for workflow definitions."""

    def register(self, workflow: WorkflowDefinition) -> str:
        """Register a workflow definition."""

    def get(self, workflow_id: str) -> Optional[WorkflowDefinition]:
        """Get workflow by ID."""

    def list_workflows(
        self,
        tags: Optional[List[str]] = None,
        active_only: bool = True,
        limit: int = 100,
        offset: int = 0,
    ) -> List[WorkflowSummary]:
        """List all registered workflows."""
```

**Relationship to Global Registry**:
- Global registry stores the workflow instance
- Workflow registry stores execution state and history
- Both work together: Global for discovery, Workflow for execution

### Tool Registry (backend/agno/tools/registry.py)

**Purpose**: Manages tool catalog, availability, and metadata.

**Key Features**:
- Tool categorization (search, finance, code, dev, comms, etc.)
- Authentication checking
- Availability validation
- Tool instance caching

```python
class ToolCategory(str, Enum):
    """Tool categories for organization."""
    SEARCH = "search"
    FINANCE = "finance"
    CODE = "code"
    DEV = "dev"
    COMMS = "comms"
    MEDIA = "media"
    HOME = "home"
    CUSTOM = "custom"
    REASONING = "reasoning"

@dataclass
class ToolInfo:
    """Metadata about a tool."""
    name: str
    description: str
    category: ToolCategory
    instance: Any  # The actual tool instance
    requires_auth: bool = False
    auth_env: Optional[str] = None
    is_builtin: bool = False
    is_enabled: bool = True
    requires_confirmation: bool = False
```

**Core Functions**:

```python
# Get a single tool
get_tool(name: str) -> Optional[Any]

# Resolve multiple tools by name
resolve_tools(names: List[str]) -> List[Any]

# List all available tools
list_tools() -> List[ToolInfo]

# Initialize the registry
initialize_registry() -> None
```

---

## No Conflicts: How They Work Together

### Why No Conflicts Exist

1. **Different Purposes**:
   - `backend/shared/registry.py` → Global lookup ("What agents exist?")
   - `backend/agno/workflows/workflow/registry.py` → Workflow execution state
   - `backend/agno/tools/registry.py` → Tool availability & metadata

2. **Agno Has NO Built-in Registry**:
   - Agno just provides `Agent`, `Team`, `Workflow` classes
   - It doesn't care how you store or discover them
   - OneMind OS adds its own registry layer on top

3. **Complementary Roles**:
   - They operate at different abstraction levels
   - They manage different concerns
   - They enhance each other's functionality

### Registry Comparison Table

| Registry | Location | Manages | Singleton? | Storage |
|----------|----------|---------|------------|---------|
| **Global** | backend/shared/registry.py | Agents, Teams, Workflows | Module-level | Dicts |
| **Workflow** | backend/agno/workflows/workflow/registry.py | Workflow runs & state | Yes (class) | Dicts |
| **Tool** | backend/agno/tools/registry.py | Tool catalog & metadata | Module-level | Dict |

---

## Usage Examples

### Example 1: Complete Agent Registration and Execution Flow

```python
# Step 1: Create an Agno Agent
from agno.agent import Agent
from agno.tools.duckduckgo import DuckDuckGoTools
from backend.shared import get_model

researcher = Agent(
    name="Researcher",
    model=get_model("sonnet"),
    tools=[DuckDuckGoTools()],
    instructions="You are a research assistant.",
    markdown=True,
)

# Step 2: Register in OneMind OS Global Registry
from backend.shared.registry import register_agent

register_agent(
    agent_id="researcher",
    agent=researcher,
    description="Research assistant for information gathering",
    category="specialist",
    tags=["research", "search", "ge"],
)

# Step 3: Retrieve from anywhere in OneMind OS
from backend.shared.registry import get_agent

agent = get_agent("researcher")  # Returns the Agno Agent instance

# Step 4: Execute (Agno framework runs it)
response = await agent.arun("Research the latest AI developments")
```

### Example 2: Team Registration with Multiple Agents

```python
from agno.team import Team
from backend.shared.registry import register_team, get_agent

# Get registered agents
web_researcher = get_agent("web-researcher")
analyst = get_agent("analyst")
writer = get_agent("writer")

# Create Agno Team
research_team = Team(
    name="Research Team",
    members=[web_researcher, analyst, writer],
    description="Collaborative research team",
)

# Register in global registry
register_team(
    team_id="research-team",
    team=research_team,
    description="Multi-agent research team",
    tags=["research", "collaboration", "ge"],
)

# Later, retrieve and use
from backend.shared.registry import get_team

team = get_team("research-team")
response = await team.arun("Research AI agent frameworks")
```

### Example 3: Using Multiple Registry Layers Together

```python
# Scenario: Execute an agent with dynamically loaded tools

# 1. Get agent from global registry
from backend.shared.registry import get_agent
legacy = get_agent("legacy")

# 2. Get tools from tool registry
from backend.agno.tools.registry import resolve_tools
tools = resolve_tools(["todoist", "github", "duckduckgo"])

# 3. Attach tools to agent (Agno handles this)
legacy.tools.extend(tools)

# 4. Register workflow in workflow registry
from backend.agno.workflows.workflow.registry import get_workflow_registry
workflow_registry = get_workflow_registry()

workflow_def = WorkflowDefinition(
    id="morning-briefing",
    name="Morning Briefing",
    steps=[...],
)
workflow_registry.register(workflow_def)

# 5. Execute workflow (uses both registries)
from backend.agno.workflows.workflow.executor import WorkflowExecutor

executor = WorkflowExecutor()
run = await executor.execute(workflow_id="morning-briefing")
```

### Example 4: Filtering and Discovery

```python
from backend.shared.registry import list_agents, list_teams, get_system_info

# List all research-related agents
research_agents = list_agents(category="specialist", tag="research")

# List all teams with "collaboration" tag
collab_teams = list_teams(tag="collaboration")

# Get system-wide overview
system_info = get_system_info()
print(f"Total agents: {system_info['agents']['count']}")
print(f"Total teams: {system_info['teams']['count']}")
print(f"Total workflows: {system_info['workflows']['count']}")

# Output:
# {
#     "agents": {
#         "count": 6,
#         "items": [
#             {
#                 "id": "legacy",
#                 "name": "Legacy",
#                 "description": "Primary AI companion",
#                 "category": "general",
#                 "tags": ["ceo", "primary"]
#             },
#             ...
#         ]
#     },
#     "teams": {...},
#     "workflows": {...}
# }
```

### Example 5: Auto-Registration Pattern

```python
# In backend/agno/teams/research_team.py

async def register_research_team():
    """Register the Research Team in the central registry."""
    from backend.shared.registry import register_team

    team = await create_research_team()

    register_team(
        team_id="research-team",
        team=team,
        description="Research team for information gathering",
        tags=["research", "analysis", "ge", "it"],
    )

    logger.info("Registered Research Team")

# Auto-register on import
# This runs when the module is imported
if __name__ != "__main__":
    import asyncio
    asyncio.create_task(register_research_team())
```

---

## Best Practices

### 1. Always Register Agents, Teams, and Workflows

```python
# ✅ GOOD: Register after creation
agent = create_my_agent()
register_agent("my-agent", agent, description="...", tags=["..."])

# ❌ BAD: Create without registering
agent = create_my_agent()
# Agent is created but not discoverable
```

### 2. Use Descriptive IDs and Tags

```python
# ✅ GOOD: Clear, semantic IDs and tags
register_agent(
    agent_id="finance-analyst",
    agent=agent,
    description="Analyzes financial data and market trends",
    category="specialist",
    tags=["finance", "analysis", "it"],
)

# ❌ BAD: Generic IDs and no tags
register_agent(
    agent_id="agent1",
    agent=agent,
    description="An agent",
)
```

### 3. Leverage Auto-Registration

```python
# ✅ GOOD: Auto-register on module import
def _auto_register():
    from backend.shared.registry import register_agent
    register_agent("my-agent", create_my_agent())

_auto_register()

# ❌ BAD: Manual registration everywhere
# Requires explicit calls in multiple places
```

### 4. Use Global Registry for Discovery, Specialized Registries for State

```python
# ✅ GOOD: Separate concerns
# Discovery
agent = get_agent("legacy")

# Execution state
workflow_registry = get_workflow_registry()
run = workflow_registry.get_run(run_id)

# ❌ BAD: Mix concerns
# Don't try to get execution state from global registry
```

### 5. Clean Up When Testing

```python
# ✅ GOOD: Clear registries in test teardown
import pytest
from backend.shared.registry import clear_all

@pytest.fixture(autouse=True)
def reset_registry():
    """Reset registry before each test."""
    clear_all()
    yield
    clear_all()

# ❌ BAD: Leave test data in registry
# Causes test pollution and flaky tests
```

---

## Architecture Diagrams

### Component Interaction

```
┌──────────────┐
│   app.py     │  Main entry point
└──────┬───────┘
       │
       │ creates
       ▼
┌──────────────────────────────────────┐
│  Agent/Team/Workflow Instances       │
│  (using Agno framework classes)      │
└──────┬───────────────────────────────┘
       │
       │ registers
       ▼
┌──────────────────────────────────────┐
│  Global Registry                     │
│  (backend/shared/registry.py)        │
│  - Discovery                         │
│  - Metadata                          │
│  - Filtering                         │
└──────┬───────────────────────────────┘
       │
       │ used by
       ▼
┌──────────────────────────────────────┐
│  API Endpoints                       │
│  - /agents (list, get)               │
│  - /teams (list, get)                │
│  - /workflows (list, get)            │
└──────────────────────────────────────┘
```

### Execution Flow

```
1. User Request
   │
   ▼
2. API Endpoint (/agents/legacy/runs)
   │
   ▼
3. Get Agent from Registry
   │  registry.get_agent("legacy")
   ▼
4. Get Tools from Tool Registry
   │  resolve_tools(["todoist", "github"])
   ▼
5. Execute with Agno
   │  agent.arun("Create a task")
   ▼
6. Response to User
```

---

## Summary

The OneMind OS registry system is a **multi-layered architecture** that provides:

1. **Global Discovery** (`backend/shared/registry.py`)
   - Single source of truth
   - System-wide agent/team/workflow lookup
   - Metadata and tagging

2. **Specialized Registries**
   - Workflow execution state
   - Tool availability checking
   - Domain-specific concerns

3. **Agno Integration**
   - Built on top of Agno framework
   - Agno provides execution engine
   - Registry provides discovery layer

**Key Insight**: The registries work **together** at different levels - like how an operating system has a process registry, file registry, and device registry, all coexisting and complementing each other.

This architecture provides:
- ✅ Centralized discovery
- ✅ Clean separation of concerns
- ✅ Flexible extensibility
- ✅ Type-safe operations
- ✅ Easy testing and mocking

---

## Related Documentation

- [Agno Framework Overview](./AGNO-FRAMEWORK.md)
- [Agent Architecture](./AGENT-ARCHITECTURE.md)
- [Workflow System](./WORKFLOW-SYSTEM.md)
- [Tool System](./TOOL-SYSTEM.md)
- [API Design](./API-DESIGN.md)

---

**Last Updated**: January 2025
**Maintained By**: Zeus + Legacy
**Status**: Production

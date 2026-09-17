# Registry System Quick Reference

Quick reference guide for working with OneMind OS registries.

## Quick Start

```python
# Import the global registry
from backend.shared.registry import (
    register_agent,
    get_agent,
    list_agents,
    register_team,
    get_team,
    list_teams,
    register_workflow,
    get_workflow,
    list_workflows,
)

# Create and register an agent
from agno.agent import Agent

agent = Agent(name="MyAgent", ...)
register_agent("my-agent", agent, description="My custom agent", tags=["custom"])

# Retrieve and use
agent = get_agent("my-agent")
response = await agent.arun("Do something")
```

---

## Common Operations

### Register an Agent

```python
from backend.shared.registry import register_agent

register_agent(
    agent_id="my-agent",           # Unique identifier
    agent=agent_instance,           # Agno Agent instance
    description="What it does",     # Human-readable description
    category="specialist",          # Category (general/specialist/workflow)
    tags=["tag1", "tag2"],         # Tags for filtering
)
```

### Get an Agent

```python
from backend.shared.registry import get_agent

# Returns the Agno Agent instance
agent = get_agent("legacy")
if agent:
    response = await agent.arun("Hello")
```

### List Agents

```python
from backend.shared.registry import list_agents

# List all agents
all_agents = list_agents()

# Filter by category
specialists = list_agents(category="specialist")

# Filter by tag
research_agents = list_agents(tag="research")
```

### Register a Team

```python
from backend.shared.registry import register_team

register_team(
    team_id="research-team",
    team=team_instance,
    description="Multi-agent research team",
    tags=["research", "collaboration"],
)
```

### Get a Team

```python
from backend.shared.registry import get_team

team = get_team("research-team")
if team:
    response = await team.arun("Research AI frameworks")
```

### Register a Workflow

```python
from backend.shared.registry import register_workflow

register_workflow(
    workflow_id="morning-briefing",
    workflow=workflow_instance,
    description="Daily morning briefing",
    tags=["automation", "daily"],
)
```

### Get System Info

```python
from backend.shared.registry import get_system_info

info = get_system_info()
print(f"Agents: {info['agents']['count']}")
print(f"Teams: {info['teams']['count']}")
print(f"Workflows: {info['workflows']['count']}")
```

---

## Registry Types

### Global Registry

**Location**: `backend/shared/registry.py`

**Purpose**: System-wide discovery

**Stores**:
- Agents (with categories and tags)
- Teams (with member lists)
- Workflows (with tags)

**Operations**:
- `register_agent()`, `get_agent()`, `list_agents()`
- `register_team()`, `get_team()`, `list_teams()`
- `register_workflow()`, `get_workflow()`, `list_workflows()`
- `get_system_info()`
- `clear_all()` (testing only)

### Tool Registry

**Location**: `backend/agno/tools/registry.py`

**Purpose**: Tool availability and metadata

**Operations**:
```python
from backend.agno.tools.registry import get_tool, resolve_tools, list_tools

# Get single tool
tool = get_tool("todoist")

# Get multiple tools
tools = resolve_tools(["todoist", "github", "duckduckgo"])

# List all tools
all_tools = list_tools()
```

### Workflow Registry

**Location**: `backend/agno/workflows/workflow/registry.py`

**Purpose**: Workflow execution state

**Operations**:
```python
from backend.agno.workflows.workflow.registry import get_workflow_registry

registry = get_workflow_registry()

# Register workflow definition
registry.register(workflow_def)

# Get workflow
workflow = registry.get("workflow-id")

# List workflows
workflows = registry.list_workflows(tags=["automation"])
```

---

## Agent Categories

| Category | Purpose | Examples |
|----------|---------|----------|
| `general` | General-purpose agents | Legacy (CEO agent) |
| `specialist` | Domain specialists | Researcher, Analyst, Writer, Coder |
| `workflow` | Workflow executors | Task triager, Morning briefing |
| `custom` | User-defined agents | Your custom agents |

---

## Common Tags

| Pillar Tags | Domain Tags | Feature Tags |
|-------------|-------------|--------------|
| `hp` - Health & Physical | `research` | `automation` |
| `le` - Life & Experiences | `analysis` | `collaboration` |
| `ge` - Growth & Education | `writing` | `daily` |
| `it` - Income & Trade | `coding` | `productivity` |

---

## Auto-Registration Pattern

```python
# In your agent module (e.g., backend/agno/agents/my_agent.py)

def create_my_agent() -> Agent:
    """Create the agent."""
    return Agent(name="MyAgent", ...)

def _auto_register():
    """Auto-register on import."""
    from backend.shared.registry import register_agent

    agent = create_my_agent()
    register_agent(
        agent_id="my-agent",
        agent=agent,
        description="My custom agent",
        category="specialist",
        tags=["custom"],
    )

    logger.info("Auto-registered MyAgent")

# Auto-register when module is imported
_auto_register()
```

---

## Testing with Registry

```python
import pytest
from backend.shared.registry import register_agent, get_agent, clear_all

@pytest.fixture(autouse=True)
def reset_registry():
    """Reset registry before/after each test."""
    clear_all()
    yield
    clear_all()

def test_agent_registration():
    agent = Agent(name="TestAgent")
    register_agent("test", agent)

    retrieved = get_agent("test")
    assert retrieved is not None
    assert retrieved.name == "TestAgent"
```

---

## Common Patterns

### Pattern 1: Create + Register + Use

```python
# 1. Create Agno instance
agent = Agent(name="MyAgent", ...)

# 2. Register in OneMind OS
register_agent("my-agent", agent, ...)

# 3. Use anywhere
agent = get_agent("my-agent")
await agent.arun("Task")
```

### Pattern 2: Dynamic Tool Loading

```python
from backend.shared.registry import get_agent
from backend.agno.tools.registry import resolve_tools

# Get agent
agent = get_agent("legacy")

# Get tools based on task
if "code" in task:
    tools = resolve_tools(["github", "shell"])
    agent.tools.extend(tools)

# Execute
await agent.arun(task)
```

### Pattern 3: Team Assembly

```python
from backend.shared.registry import get_agent, register_team
from agno.team import Team

# Get agents from registry
researcher = get_agent("researcher")
analyst = get_agent("analyst")
writer = get_agent("writer")

# Assemble team
team = Team(
    name="Content Team",
    members=[researcher, analyst, writer],
)

# Register team
register_team("content-team", team, ...)
```

### Pattern 4: Workflow Discovery

```python
from backend.shared.registry import list_workflows

# Find workflows by tag
automation_workflows = list_workflows(tag="automation")

# Execute each
for workflow_info in automation_workflows:
    workflow = workflow_info.workflow
    await workflow.run()
```

---

## API Endpoints Using Registry

The registry powers these API endpoints:

```
GET  /agents           → list_agents()
GET  /agents/{id}      → get_agent(id)
POST /agents/{id}/runs → get_agent(id).arun()

GET  /teams            → list_teams()
GET  /teams/{id}       → get_team(id)
POST /teams/{id}/runs  → get_team(id).arun()

GET  /workflows        → list_workflows()
GET  /workflows/{id}   → get_workflow(id)
POST /workflows/{id}/runs → execute workflow
```

---

## Debugging

### Check What's Registered

```python
from backend.shared.registry import get_system_info

info = get_system_info()

# See all agents
for agent in info['agents']['items']:
    print(f"{agent['id']}: {agent['name']} ({agent['category']})")

# See all teams
for team in info['teams']['items']:
    print(f"{team['id']}: {team['name']} - {team['members']}")
```

### Verify Agent is Registered

```python
from backend.shared.registry import get_agent_info

info = get_agent_info("legacy")
if info:
    print(f"Name: {info.name}")
    print(f"Category: {info.category}")
    print(f"Tags: {info.tags}")
else:
    print("Agent not found!")
```

---

## Best Practices

1. ✅ **Always register after creation**
   ```python
   agent = create_agent()
   register_agent("agent-id", agent, ...)
   ```

2. ✅ **Use descriptive IDs**
   ```python
   register_agent("finance-analyst", ...)  # Good
   register_agent("agent1", ...)           # Bad
   ```

3. ✅ **Add meaningful tags**
   ```python
   tags=["finance", "analysis", "it"]  # Good
   tags=[]                              # Bad
   ```

4. ✅ **Leverage auto-registration**
   ```python
   def _auto_register():
       register_agent(...)

   _auto_register()
   ```

5. ✅ **Clean up in tests**
   ```python
   @pytest.fixture(autouse=True)
   def reset_registry():
       clear_all()
       yield
       clear_all()
   ```

---

## Troubleshooting

### Agent not found

```python
agent = get_agent("my-agent")
if agent is None:
    # Check if it was registered
    info = get_system_info()
    print("Registered agents:", [a['id'] for a in info['agents']['items']])
```

### Registry is empty

```python
# Check if auto-registration ran
import backend.agno.agents  # Trigger imports
from backend.shared.registry import get_system_info

info = get_system_info()
print(f"Found {info['agents']['count']} agents")
```

### Changes not taking effect

```python
# In development, you may need to reimport
import importlib
import backend.agno.agents

importlib.reload(backend.agno.agents)
```

---

## See Also

- [Complete Registry System Documentation](REGISTRY-SYSTEM.md)
- [Agno Framework Overview](../architecture/AGNO-FRAMEWORK.md)
- [Agent Architecture](../architecture/AGENT-ARCHITECTURE.md)
- [API Documentation](./API-REFERENCE.md)

---

**Quick Reference Version**: 1.0
**Last Updated**: January 2025


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]

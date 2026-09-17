# Letta

> **Stateful AI agents with persistent memory**

## Overview

Letta (formerly MemGPT) provides a framework for building AI agents with unlimited context through memory management, tool use, and multi-agent orchestration.

## Status

| Attribute   | Value                |
| ----------- | -------------------- |
| Status      | ✅ Active            |
| Category    | AI Agents            |
| Type        | Agent Framework      |
| Integration | REST API, Python SDK |

## Use Cases in OneMind

1. **Long-term Memory** - Agents that remember across sessions
2. **Personal Assistant** - Zeus's AI companion with full history
3. **Multi-Agent Systems** - Orchestrate specialized agents
4. **Tool Integration** - Connect agents to external services

## Key Features

- Persistent memory (core + archival)
- Tool/function calling
- Multi-agent orchestration
- Self-editing memory
- REST API server
- Python SDK

## Architecture

```
┌─────────────────────────────────────────┐
│              Letta Agent                │
├─────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────────┐   │
│  │ Core Memory │  │ Archival Memory │   │
│  │ (In-Context)│  │ (Vector DB)     │   │
│  └─────────────┘  └─────────────────┘   │
│  ┌─────────────────────────────────────┐ │
│  │           Tool Library              │ │
│  └─────────────────────────────────────┘ │
└─────────────────────────────────────────┘
```

## Installation

```bash
pip install letta

# Start server
letta server
```

## API Integration

### Python SDK

```python
from letta import Letta

client = Letta(base_url="http://localhost:8283")

# Create agent
agent = client.agents.create(
    name="legacy-ai",
    model="claude-3-opus",
    system_prompt="""You are Legacy, Zeus's AI companion.
    You have access to his full history and preferences.
    You help manage his life across all domains.""",
    tools=["web_search", "calendar", "home_assistant"]
)
```

### Memory Blocks

```python
# Add core memory block
client.agents.memory.create_block(
    agent_id=agent.id,
    label="user_profile",
    value="""
    Name: Zeus Delacruz
    Location: San Diego, CA
    Interests: AI, solar-punk, automation
    """
)

# Search archival memory
results = client.agents.memory.archival_search(
    agent_id=agent.id,
    query="Zeus's morning routine preferences",
    limit=5
)
```

### Send Messages

```python
# Chat with agent
response = client.agents.messages.create(
    agent_id=agent.id,
    role="user",
    content="What's on my schedule today?"
)

print(response.content)
```

### Custom Tools

```python
from letta import Tool

# Define tool
home_assistant_tool = Tool(
    name="control_lights",
    description="Control smart home lights",
    parameters={
        "type": "object",
        "properties": {
            "room": {"type": "string"},
            "action": {"type": "string", "enum": ["on", "off", "dim"]}
        },
        "required": ["room", "action"]
    }
)

# Attach to agent
client.agents.tools.attach(
    agent_id=agent.id,
    tool=home_assistant_tool
)
```

## Multi-Agent Setup

```python
# Create specialized agents
planner = client.agents.create(
    name="planner",
    system_prompt="You break down complex tasks into steps."
)

executor = client.agents.create(
    name="executor",
    system_prompt="You execute specific tasks and report results."
)

researcher = client.agents.create(
    name="researcher",
    system_prompt="You gather information from various sources."
)

# Orchestration via main agent
orchestrator = client.agents.create(
    name="orchestrator",
    system_prompt="You coordinate other agents to complete tasks.",
    tools=["delegate_to_planner", "delegate_to_executor", "delegate_to_researcher"]
)
```

## REST API

```python
import httpx

LETTA_URL = "http://localhost:8283"

async def send_message(agent_id: str, message: str):
    """Send message to Letta agent."""
    async with httpx.AsyncClient() as client:
        response = await client.post(
            f"{LETTA_URL}/v1/agents/{agent_id}/messages",
            json={
                "role": "user",
                "content": message
            }
        )
        return response.json()

async def update_memory_block(agent_id: str, label: str, value: str):
    """Update agent's core memory block."""
    async with httpx.AsyncClient() as client:
        response = await client.patch(
            f"{LETTA_URL}/v1/agents/{agent_id}/memory/blocks/{label}",
            json={"value": value}
        )
        return response.json()
```

## Legacy AI Integration

```python
from agno import Tool

@Tool
async def query_letta_agent(
    agent_name: str,
    query: str
) -> str:
    """Query a specialized Letta agent.

    Args:
        agent_name: Name of the Letta agent to query
        query: Question or task for the agent
    """
    agent_id = get_agent_id_by_name(agent_name)
    response = await send_message(agent_id, query)
    return response["content"]

@Tool
async def update_agent_memory(
    agent_name: str,
    memory_block: str,
    content: str
) -> str:
    """Update a Letta agent's memory.

    Args:
        agent_name: Name of the agent
        memory_block: Which memory block to update
        content: New content for the block
    """
    agent_id = get_agent_id_by_name(agent_name)
    await update_memory_block(agent_id, memory_block, content)
    return f"Updated {memory_block} for {agent_name}"
```

## Docker Deployment

```yaml
version: "3.8"

services:
  letta:
    image: letta/letta:latest
    ports:
      - "8283:8283"
    volumes:
      - letta_data:/root/.letta
    environment:
      - OPENAI_API_KEY=${OPENAI_API_KEY}
      - ANTHROPIC_API_KEY=${ANTHROPIC_API_KEY}

volumes:
  letta_data:
```

## Environment Variables

```bash
LETTA_URL=http://localhost:8283
LETTA_API_KEY=your_api_key

# Model providers
OPENAI_API_KEY=your_key
ANTHROPIC_API_KEY=your_key
```

## Best Practices

1. **Memory design** - Plan core memory blocks carefully
2. **Archival strategy** - Archive frequently for long-term recall
3. **Tool boundaries** - Clear, focused tools work best
4. **Agent specialization** - Single-purpose agents > generalist

## Resources

- [Letta Documentation](https://docs.letta.com/)
- [Letta GitHub](https://github.com/letta-ai/letta)
- [MemGPT Paper](https://arxiv.org/abs/2310.08560)
- [Discord Community](https://discord.gg/letta)

---

Last Updated: January 25, 2026


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]

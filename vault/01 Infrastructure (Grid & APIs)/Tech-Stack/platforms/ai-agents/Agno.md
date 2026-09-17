# Agno

> **AI Agent Framework for Python**

## Overview

Agno is a lightweight, high-performance framework for building AI agents in Python. It's the core framework powering Legacy AI in OneMind OS.

## Status

- **Active**: ✅ Yes
- **Version**: Latest
- **Language**: Python 3.12+
- **Last Updated**: January 25, 2026

## Use Cases

1. **Legacy AI Agent** - Core agent implementation
2. **Tool Development** - Create agent capabilities
3. **Memory Management** - Conversation and long-term memory
4. **Multi-Agent Systems** - Agent orchestration

## Key Features

- Lightweight and fast (minimal dependencies)
- Model-agnostic (works with Claude, GPT, etc.)
- Built-in tool/function calling
- Memory and context management
- Streaming support
- Async-first design

## Installation

```bash
pip install agno
```

## Core Concepts

### Agent Definition

```python
from agno import Agent, Tool
from agno.models import Claude

agent = Agent(
    name="Legacy",
    model=Claude(model="claude-opus-4-5"),
    instructions="You are Legacy, Zeus's AI assistant.",
    tools=[search_tool, calendar_tool, home_tool],
    memory=True,
)
```

### Tools

```python
from agno import Tool

@Tool
def search_web(query: str) -> str:
    """Search the web for information."""
    # Implementation
    return results
```

### Running Agent

```python
# Synchronous
response = agent.run("What's on my calendar today?")

# Streaming
for chunk in agent.run_stream("Tell me about my day"):
    print(chunk, end="")

# Async
response = await agent.arun("Schedule a meeting")
```

## Project Structure

```
agents/legacy/
├── agent.py          # Agent definition
├── main.py           # FastAPI entry point
├── awareness.py      # 5-level awareness modes
├── consciousness/    # Event processing, state
│   ├── stream.py
│   └── memory.py
└── tools/            # Agent tools
    ├── calendar.py
    ├── home.py
    ├── search.py
    └── todoist.py
```

## Memory Configuration

```python
from agno.memory import PostgresMemory

memory = PostgresMemory(
    connection_string=os.getenv("DATABASE_URL"),
    table_name="agent_memory",
)

agent = Agent(
    name="Legacy",
    memory=memory,
)
```

## Model Configuration

### Claude via AWS Bedrock

```python
from agno.models import Claude

model = Claude(
    model="claude-opus-4-5",
    provider="bedrock",
    region="us-east-1",
)
```

### OpenAI

```python
from agno.models import OpenAI

model = OpenAI(
    model="gpt-4-turbo",
    api_key=os.getenv("OPENAI_API_KEY"),
)
```

## Integration Points

- **FastAPI** - HTTP endpoints for agent
- **NATS** - Event-driven triggers
- **PostgreSQL** - Memory persistence
- **LiveKit** - Voice interface
- **Home Assistant** - Smart home control

## API Endpoints

```python
# agents/legacy/main.py

@app.post("/agents/{agent_id}/runs")
async def run_agent(agent_id: str, request: RunRequest):
    response = await agent.arun(request.message)
    return {"response": response}

@app.get("/agents/{agent_id}/memory")
async def get_memory(agent_id: str):
    return agent.memory.get_all()
```

## Best Practices

- ✅ Keep tools focused (single responsibility)
- ✅ Use async for I/O operations
- ✅ Implement proper error handling in tools
- ✅ Log agent interactions for debugging
- ✅ Use streaming for long responses
- ✅ Test tools independently

## Resources

- [Agno Documentation](https://docs.agno.com)
- [GitHub Repository](https://github.com/agno-ai/agno)
- [Example Agents](https://github.com/agno-ai/agno/tree/main/examples)

---

**Related**: [Claude Code](./Claude-Code.md) | [PostgreSQL](./PostgreSQL.md) | [NATS](./NATS.md)

# Agno Framework

> Lightweight, high-performance framework for building AI agents in Python.

---

## Overview

| Attribute         | Value                                                        |
| ----------------- | ------------------------------------------------------------ |
| **Type**          | AI Agent Framework                                           |
| **Language**      | Python                                                       |
| **License**       | MIT                                                          |
| **Repository**    | [github.com/agno-agi/agno](https://github.com/agno-agi/agno) |
| **Documentation** | [docs.agno.com](https://docs.agno.com)                       |

---

## What is Agno?

Agno is a Python framework for building AI agents with:

- **Model agnostic** — Works with OpenAI, Anthropic, AWS Bedrock, local models
- **Tool calling** — Native function/tool integration
- **Memory** — Built-in conversation and session memory
- **Streaming** — First-class streaming support
- **Type-safe** — Full Pydantic integration
- **Async-first** — Native async/await support

---

## OneMind Usage

### Primary Use Case

Agno powers the **Legacy Agent** — the core AI agent in OneMind OS.

```
┌─────────────────────────────────────────────────────────────────┐
│                        Legacy Agent                              │
│                                                                  │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐       │
│  │  Agno Agent  │───▶│    Tools     │───▶│   Actions    │       │
│  │  (Claude)    │    │              │    │              │       │
│  └──────────────┘    └──────────────┘    └──────────────┘       │
│         │                                                        │
│         ▼                                                        │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐       │
│  │   Memory     │    │  Awareness   │    │    NATS      │       │
│  │   (Letta)    │    │   Levels     │    │   Events     │       │
│  └──────────────┘    └──────────────┘    └──────────────┘       │
└─────────────────────────────────────────────────────────────────┘
```

### Why Agno for OneMind?

| Requirement         | Agno Solution                       |
| ------------------- | ----------------------------------- |
| Multi-model support | Bedrock (Claude), OpenAI, local     |
| Tool integration    | Native tool calling with validation |
| Streaming responses | Built-in streaming for real-time    |
| Session management  | Memory and context handling         |
| Production ready    | Async, typed, tested                |

---

## Core Concepts

### Agent Definition

```python
from agno.agent import Agent
from agno.models.aws import BedrockClaude

# Define an agent
agent = Agent(
    name="Legacy",
    model=BedrockClaude(model_id="anthropic.claude-3-5-sonnet-20241022-v2:0"),
    instructions=[
        "You are Legacy, Zeus's AI counterpart.",
        "You help manage his life across all domains.",
    ],
    tools=[home_assistant_tools, todoist_tools, calendar_tools],
    markdown=True,
)
```

### Tools

```python
from agno.tools import tool

@tool
def get_weather(location: str) -> str:
    """Get current weather for a location."""
    # Implementation
    return f"Weather in {location}: 72°F, sunny"

@tool
def control_lights(room: str, state: bool) -> str:
    """Turn lights on or off in a room."""
    # Call Home Assistant
    return f"Lights in {room} turned {'on' if state else 'off'}"
```

### Running the Agent

```python
# Synchronous
response = agent.run("What's the weather in Austin?")
print(response.content)

# Async
response = await agent.arun("Turn on the kitchen lights")

# Streaming
async for chunk in agent.arun_stream("Tell me about my day"):
    print(chunk.content, end="")
```

### Memory Integration

```python
from agno.memory import Memory
from agno.storage import PostgresStorage

# Persistent memory
memory = Memory(
    storage=PostgresStorage(
        table_name="agent_sessions",
        db_url="postgresql://user:pass@localhost:5432/onemind"
    )
)

agent = Agent(
    name="Legacy",
    model=model,
    memory=memory,
    add_history_to_messages=True,
)
```

---

## OneMind Agent Architecture

### agents/legacy/agent.py

```python
from agno.agent import Agent
from agno.models.aws import BedrockClaude

from .tools import (
    home_assistant_tools,
    todoist_tools,
    calendar_tools,
    nats_tools,
    knowledge_tools,
)
from .awareness import AwarenessManager

class LegacyAgent:
    def __init__(self):
        self.model = BedrockClaude(
            model_id="anthropic.claude-3-5-sonnet-20241022-v2:0",
            region="us-east-1",
        )

        self.awareness = AwarenessManager()

        self.agent = Agent(
            name="Legacy",
            model=self.model,
            instructions=self._get_instructions(),
            tools=self._get_tools(),
            markdown=True,
            show_tool_calls=True,
        )

    def _get_instructions(self) -> list[str]:
        return [
            "You are Legacy, the AI counterpart to Zeus.",
            f"Current awareness level: {self.awareness.current_level}",
            "Respond concisely but warmly.",
        ]

    def _get_tools(self) -> list:
        return [
            *home_assistant_tools,
            *todoist_tools,
            *calendar_tools,
            *nats_tools,
            *knowledge_tools,
        ]

    async def run(self, message: str, session_id: str = None):
        return await self.agent.arun(message, session_id=session_id)
```

---

## Tool Categories

### Home Automation

```python
@tool
def get_home_state() -> dict:
    """Get current state of all home devices."""

@tool
def control_device(entity_id: str, action: str) -> str:
    """Control a Home Assistant device."""

@tool
def run_automation(automation_name: str) -> str:
    """Trigger a Home Assistant automation."""
```

### Task Management

```python
@tool
def get_tasks(filter: str = "today") -> list:
    """Get tasks from Todoist."""

@tool
def create_task(title: str, due: str = None, priority: int = 1) -> str:
    """Create a new task in Todoist."""

@tool
def complete_task(task_id: str) -> str:
    """Mark a task as complete."""
```

### Knowledge

```python
@tool
def search_knowledge(query: str) -> list:
    """Search the knowledge base."""

@tool
def store_knowledge(content: str, category: str) -> str:
    """Store new knowledge in the vault."""
```

### Events

```python
@tool
def publish_event(subject: str, data: dict) -> str:
    """Publish event to NATS."""

@tool
def get_recent_events(subject: str, limit: int = 10) -> list:
    """Get recent events from a subject."""
```

---

## Model Configuration

### AWS Bedrock (Production)

```python
from agno.models.aws import BedrockClaude

model = BedrockClaude(
    model_id="anthropic.claude-3-5-sonnet-20241022-v2:0",
    region="us-east-1",
    # Uses AWS credentials from environment or IAM role
)
```

### OpenAI (Alternative)

```python
from agno.models.openai import OpenAIChat

model = OpenAIChat(
    id="gpt-4-turbo",
    api_key=os.getenv("OPENAI_API_KEY"),
)
```

### Local Models (Development)

```python
from agno.models.ollama import Ollama

model = Ollama(
    id="llama3.2",
    host="http://localhost:11434",
)
```

---

## Streaming with FastAPI

```python
from fastapi import FastAPI
from fastapi.responses import StreamingResponse
from agno.agent import Agent

app = FastAPI()

@app.post("/chat/stream")
async def chat_stream(message: str):
    async def generate():
        async for chunk in agent.arun_stream(message):
            if chunk.content:
                yield f"data: {chunk.content}\n\n"
        yield "data: [DONE]\n\n"

    return StreamingResponse(
        generate(),
        media_type="text/event-stream"
    )
```

---

## Best Practices

### 1. Tool Design

```python
# ✅ Good: Clear, specific, typed
@tool
def get_room_temperature(room: str) -> float:
    """Get the current temperature in a specific room."""

# ❌ Bad: Vague, untyped
@tool
def do_something(data):
    """Does something with data."""
```

### 2. Error Handling

```python
@tool
def control_device(entity_id: str, action: str) -> str:
    """Control a Home Assistant device."""
    try:
        result = ha_client.call_service(entity_id, action)
        return f"Successfully {action} {entity_id}"
    except HomeAssistantError as e:
        return f"Failed to control {entity_id}: {str(e)}"
```

### 3. Context Management

```python
# Include relevant context in instructions
instructions = [
    f"Current time: {datetime.now().isoformat()}",
    f"Zeus's location: {get_user_location()}",
    f"Active awareness level: {awareness.level}",
]
```

---

## Configuration

### Environment Variables

```bash
# AWS Bedrock
AWS_ACCESS_KEY_ID=...
AWS_SECRET_ACCESS_KEY=...
AWS_DEFAULT_REGION=us-east-1

# OpenAI (optional)
OPENAI_API_KEY=...

# Agent settings
AGENT_MODEL=bedrock
AGENT_MAX_TOKENS=4096
AGENT_TEMPERATURE=0.7
```

### Agent Configuration

```python
# config/agent.py
from pydantic_settings import BaseSettings

class AgentSettings(BaseSettings):
    model_provider: str = "bedrock"
    model_id: str = "anthropic.claude-3-5-sonnet-20241022-v2:0"
    max_tokens: int = 4096
    temperature: float = 0.7

    class Config:
        env_prefix = "AGENT_"
```

---

## Monitoring & Debugging

### Logging

```python
import logging
from agno.agent import Agent

logging.basicConfig(level=logging.DEBUG)

agent = Agent(
    name="Legacy",
    model=model,
    debug_mode=True,  # Verbose output
)
```

### Tool Call Inspection

```python
response = await agent.arun("Turn on the lights")

# Inspect tool calls
for tool_call in response.tool_calls:
    print(f"Tool: {tool_call.name}")
    print(f"Args: {tool_call.arguments}")
    print(f"Result: {tool_call.result}")
```

---

## Resources

- [Agno Documentation](https://docs.agno.com)
- [Agno GitHub](https://github.com/agno-agi/agno)
- [Agno Examples](https://github.com/agno-agi/agno/tree/main/examples)
- [OneMind Agent Code](../../agents/legacy/)

---

_Last updated: January 2026_

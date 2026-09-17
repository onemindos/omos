---
id: ui-02-20260111053925
title: "Building AI Agents with Letta"
type: reference
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/02]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "Building AI Agents with Letta"
---

# Building AI Agents with Letta

## What is Letta?

Letta is an open-source framework for building stateful AI agents with long-term memory. It enables agents to remember conversations, learn from interactions, and maintain context across sessions.

## Key Features

### 1. Long-Term Memory
- Persistent memory across conversations
- Agents remember previous interactions
- Context retention for personalized responses

### 2. Stateful Agent Architecture
- Maintain conversation state
- Build complex workflows
- Multi-turn dialogue capabilities

### 3. Memory Management
- Core memory (always accessible)
- Archival memory (long-term storage)
- Recall memory (conversation history)

## Why Letta?

Unlike stateless LLM APIs, Letta agents can:
- Build relationships over time
- Learn user preferences
- Maintain project context
- Execute complex, multi-step tasks
- Integrate with external tools and APIs

## Use Cases

- Personal AI assistants
- Customer support agents
- Research assistants
- Task automation
- Knowledge management systems

## Getting Started

```python
from letta import create_client

# Initialize Letta client
client = create_client()

# Create an agent
agent = client.create_agent(
    name="my_agent",
    memory={
        "human": "User preferences and context",
        "persona": "Agent personality and role"
    }
)

# Send message
response = client.send_message(
    agent_id=agent.id,
    message="Hello, remember this conversation"
)
```

## Learn More

Visit [letta.com](https://letta.com) for documentation and examples. 

---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]

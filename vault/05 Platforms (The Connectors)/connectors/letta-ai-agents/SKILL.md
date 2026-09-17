---
name: letta
description: Manage Letta AI agents, tools, memory blocks, and archival memory. Use when you need to create agents, attach tools, update memory, search archival, or send messages to Letta agents. Makes API calls to Letta server.
allowed-tools:
  - Bash(curl:*)
  - Bash(jq:*)
  - Read
  - Write
---

# Letta Integration Skill

Interact with Letta AI agents through their REST API for One Mind infrastructure management.

## Environment Variables Required

```bash
LETTA_API_KEY      # Bearer token (sk-let-xxx)
LETTA_BASE_URL     # Default: https://api.letta.com/v1 (or http://localhost:8283 for local)
```

## Quick Reference

Note: `LETTA_BASE_URL` already includes `/v1` (e.g., `https://api.letta.com/v1`)

| Operation | Command |
|-----------|---------|
| List agents | `curl -s "$LETTA_BASE_URL/agents" -H "Authorization: Bearer $LETTA_API_KEY"` |
| Get agent | `curl -s "$LETTA_BASE_URL/agents/{id}" -H "Authorization: Bearer $LETTA_API_KEY"` |
| List tools | `curl -s "$LETTA_BASE_URL/tools" -H "Authorization: Bearer $LETTA_API_KEY"` |
| Get memory | `curl -s "$LETTA_BASE_URL/agents/{id}/memory" -H "Authorization: Bearer $LETTA_API_KEY"` |

## Operations

### Agents

**List all agents:**
```bash
curl -s "$LETTA_BASE_URL/agents?limit=50" \
  -H "Authorization: Bearer $LETTA_API_KEY" | jq '.[] | {id, name, created_at}'
```

**Get agent details:**
```bash
curl -s "$LETTA_BASE_URL/agents/{agent_id}" \
  -H "Authorization: Bearer $LETTA_API_KEY" | jq '{name, model, system: .system[0:200]}'
```

**Create agent:**
```bash
curl -s -X POST "$LETTA_BASE_URL/agents" \
  -H "Authorization: Bearer $LETTA_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "agent-name",
    "model": "anthropic/claude-sonnet-4-5-20250929",
    "embedding": "openai/text-embedding-3-small",
    "system": "You are a helpful assistant...",
    "memory_blocks": [
      {"label": "persona", "value": "Agent persona description"},
      {"label": "human", "value": "User: Zeus Delacruz"}
    ]
  }'
```

**Delete agent:**
```bash
curl -s -X DELETE "$LETTA_BASE_URL/agents/{agent_id}" \
  -H "Authorization: Bearer $LETTA_API_KEY"
```

### Tools

**List all tools:**
```bash
curl -s "$LETTA_BASE_URL/tools?limit=50" \
  -H "Authorization: Bearer $LETTA_API_KEY" | jq '.[] | {id, name, description}'
```

**Create tool:**
```bash
curl -s -X POST "$LETTA_BASE_URL/tools" \
  -H "Authorization: Bearer $LETTA_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "tool_name",
    "description": "What this tool does",
    "source_code": "def tool_name(param: str) -> str:\n    \"\"\"Docstring\"\"\"\n    return result",
    "tags": ["category"]
  }'
```

**Attach tool to agent:**
```bash
curl -s -X POST "$LETTA_BASE_URL/agents/{agent_id}/tools/{tool_id}" \
  -H "Authorization: Bearer $LETTA_API_KEY"
```

**Detach tool from agent:**
```bash
curl -s -X DELETE "$LETTA_BASE_URL/agents/{agent_id}/tools/{tool_id}" \
  -H "Authorization: Bearer $LETTA_API_KEY"
```

**List agent's tools:**
```bash
curl -s "$LETTA_BASE_URL/agents/{agent_id}/tools" \
  -H "Authorization: Bearer $LETTA_API_KEY" | jq '.[] | {id, name}'
```

### Memory Blocks

**Get agent memory:**
```bash
curl -s "$LETTA_BASE_URL/agents/{agent_id}/memory" \
  -H "Authorization: Bearer $LETTA_API_KEY" | jq '.memory.blocks[] | {id, label, value: .value[0:100]}'
```

**Update memory block:**
```bash
curl -s -X PATCH "$LETTA_BASE_URL/blocks/{block_id}" \
  -H "Authorization: Bearer $LETTA_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"value": "New block content here"}'
```

### Archival Memory

**Insert text to archival:**
```bash
curl -s -X POST "$LETTA_BASE_URL/agents/{agent_id}/archival" \
  -H "Authorization: Bearer $LETTA_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"text": "Content to store in archival memory"}'
```

**Search archival:**
```bash
curl -s "$LETTA_BASE_URL/agents/{agent_id}/archival?query=search+term&limit=10" \
  -H "Authorization: Bearer $LETTA_API_KEY" | jq '.[] | {id, text: .text[0:100]}'
```

**List archival passages:**
```bash
curl -s "$LETTA_BASE_URL/agents/{agent_id}/archival?limit=20" \
  -H "Authorization: Bearer $LETTA_API_KEY"
```

### Messaging

**Send message to agent:**
```bash
curl -s -X POST "$LETTA_BASE_URL/agents/{agent_id}/messages" \
  -H "Authorization: Bearer $LETTA_API_KEY" \
  -H "Content-Type: application/json" \
  --max-time 120 \
  -d '{"messages": [{"role": "user", "content": "Your message here"}]}'
```

## One Mind Agents

| Agent | Role | Domain |
|-------|------|--------|
| UI (Mother AI) | Central orchestrator | 00-24 |
| HP Agent | Performance coach | 25-49 |
| LE Agent | Family steward | 50-74 |
| GE Agent | Business advisor | 75-99 |

## Error Handling

- **200/201**: Success
- **204**: Deleted successfully
- **401**: Invalid API key
- **404**: Agent/tool not found
- **429**: Rate limited (wait and retry)
- **500**: Server error

Always check response status and parse error messages from response body.

## Tips

- Use `jq` to parse JSON responses
- Set `--max-time 120` for message operations (agent processing takes time)
- Store agent/tool IDs after creation for future reference
- Memory blocks have labels: `persona`, `human`, custom blocks


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]

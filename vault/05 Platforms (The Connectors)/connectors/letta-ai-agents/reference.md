# Letta API Reference

## Base Configuration

```bash
# Production (Letta Cloud) - includes /v1
LETTA_BASE_URL=https://api.letta.com/v1

# Local/Self-hosted
LETTA_BASE_URL=http://localhost:8283

# Authentication header
Authorization: Bearer $LETTA_API_KEY
```

**Important**: `LETTA_BASE_URL` already includes `/v1`, so endpoints are: `$LETTA_BASE_URL/agents`, not `$LETTA_BASE_URL/v1/agents`

## Endpoints

### Agents

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/agents` | List all agents |
| GET | `/v1/agents/{id}` | Get agent details |
| POST | `/agents` | Create agent |
| DELETE | `/v1/agents/{id}` | Delete agent |
| GET | `/v1/agents/{id}/tools` | List agent's tools |
| POST | `/v1/agents/{id}/tools/{tool_id}` | Attach tool |
| DELETE | `/v1/agents/{id}/tools/{tool_id}` | Detach tool |

### Tools

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/tools` | List all tools |
| GET | `/v1/tools/{id}` | Get tool details |
| POST | `/tools` | Create tool |
| DELETE | `/v1/tools/{id}` | Delete tool |

### Memory

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/agents/{id}/memory` | Get memory blocks |
| PATCH | `/v1/blocks/{block_id}` | Update memory block |

### Archival Memory

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/agents/{id}/archival` | List passages |
| GET | `/v1/agents/{id}/archival?query={q}` | Search passages |
| POST | `/v1/agents/{id}/archival` | Insert text |

### Messages

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/v1/agents/{id}/messages` | Send message |

## Request/Response Formats

### Create Agent Request
```json
{
  "name": "agent-name",
  "model": "anthropic/claude-sonnet-4-5-20250929",
  "embedding": "openai/text-embedding-3-small",
  "system": "System prompt...",
  "memory_blocks": [
    {"label": "persona", "value": "..."},
    {"label": "human", "value": "..."}
  ]
}
```

### Create Tool Request
```json
{
  "name": "tool_name",
  "description": "Tool description",
  "source_code": "def tool_name(param: str) -> str:\n    ...",
  "tags": ["category"]
}
```

### Send Message Request
```json
{
  "messages": [
    {"role": "user", "content": "Message content"}
  ]
}
```

### Update Memory Block Request
```json
{
  "value": "New block content"
}
```

### Agent Response (key fields)
```json
{
  "id": "agent-xxx",
  "name": "Agent Name",
  "model": "anthropic/claude-sonnet-4-5-20250929",
  "system": "System prompt...",
  "tools": ["tool-id-1", "tool-id-2"],
  "memory": {
    "blocks": [
      {"id": "block-xxx", "label": "persona", "value": "..."}
    ]
  },
  "created_at": "2025-01-10T..."
}
```

### Message Response (key fields)
```json
{
  "messages": [
    {
      "message_type": "assistant_message",
      "assistant_message": "Response text..."
    },
    {
      "message_type": "tool_call_message",
      "tool_call": {"name": "tool_name", "arguments": {...}}
    }
  ]
}
```

## Timeouts

| Operation | Recommended Timeout |
|-----------|-------------------|
| GET requests | 30s |
| POST (create) | 60s |
| POST (messages) | 120s |
| DELETE | 30s |

## Models Available

- `anthropic/claude-sonnet-4-5-20250929`
- `anthropic/claude-3-5-haiku-20241022`
- `openai/gpt-4o`
- `openai/gpt-4o-mini`

## Embeddings

- `openai/text-embedding-3-small`
- `openai/text-embedding-3-large`


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]

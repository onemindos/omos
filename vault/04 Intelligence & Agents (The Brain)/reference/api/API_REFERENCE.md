# Legacy API Reference

> OneMind OS - Agno AgentOS API
> Base URL: `https://agno.onemindos.dev`
> Swagger UI: `https://agno.onemindos.dev/docs`

---

## Quick Start

### Chat with Legacy

```bash
curl -X POST https://agno.onemindos.dev/agents/legacy-mother-ai/runs \
  -F "message=Hello Legacy" \
  -F "stream=false"
```

### With Session (remembers conversation)

```bash
curl -X POST https://agno.onemindos.dev/agents/legacy-mother-ai/runs \
  -F "message=Hello Legacy" \
  -F "session_id=my-session-123" \
  -F "user_id=zeus" \
  -F "stream=false"
```

---

## Endpoints

### Health Check

```
GET /health
```

**Response:**
```json
{
  "status": "ok",
  "instantiated_at": "1768818236.507876"
}
```

---

### List Agents

```
GET /agents
```

**Response:**
```json
[
  {
    "id": "legacy-mother-ai",
    "name": "Legacy",
    "description": "OneMind OS Mother AI - Zeus's trusted AI partner",
    "model": {
      "name": "AwsBedrock",
      "model": "us.anthropic.claude-opus-4-5-20251101-v1:0",
      "provider": "AwsBedrock"
    },
    "tools": {...}
  }
]
```

---

### Get Agent Details

```
GET /agents/{agent_id}
```

**Parameters:**
- `agent_id` (path): Agent ID (e.g., `legacy-mother-ai`)

---

### Chat with Agent (Main Endpoint)

```
POST /agents/{agent_id}/runs
```

**Parameters:**
- `agent_id` (path): Agent ID (e.g., `legacy-mother-ai`)

**Form Data:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `message` | string | Yes | Your message to Legacy |
| `stream` | boolean | No | Stream response (default: false) |
| `session_id` | string | No | Session ID for conversation continuity |
| `user_id` | string | No | User identifier (default: anonymous) |
| `images` | file[] | No | Image files for vision analysis |
| `audio` | file[] | No | Audio files for transcription |
| `videos` | file[] | No | Video files for analysis |
| `files` | file[] | No | Documents (PDF, CSV, etc.) |

**Response (non-streaming):**
```json
{
  "run_id": "79eed5bb-2926-4a58-98de-5bceb9f58bd4",
  "agent_id": "legacy-mother-ai",
  "agent_name": "Legacy",
  "session_id": "9e6e2447-7835-4c44-9a01-13c89592eff8",
  "content": "Legacy online and operational, Zeus — OneMind OS v2.0 ready to serve.",
  "content_type": "str",
  "model": "us.anthropic.claude-opus-4-5-20251101-v1:0",
  "model_provider": "AwsBedrock",
  "created_at": 1768818233,
  "status": "COMPLETED",
  "metrics": {
    "input_tokens": 2584,
    "output_tokens": 24,
    "total_tokens": 2608,
    "time_to_first_token": 0.087,
    "duration": 2.528
  }
}
```

**Streaming Response:**
When `stream=true`, returns Server-Sent Events (SSE):
```
event: RunStarted
data: {"run_id": "123..."}

event: RunResponseContentDelta
data: {"content": "Hello"}

event: RunResponseContentDelta
data: {"content": " Zeus"}

event: RunCompleted
data: {"status": "COMPLETED"}
```

---

### List Sessions

```
GET /agents/{agent_id}/sessions
```

**Parameters:**
- `agent_id` (path): Agent ID
- `user_id` (query): Filter by user

---

### Get Session Details

```
GET /sessions/{session_id}
```

**Response:**
```json
{
  "session_id": "9e6e2447-7835-4c44-9a01-13c89592eff8",
  "agent_id": "legacy-mother-ai",
  "user_id": "zeus",
  "messages": [...],
  "created_at": 1768818233
}
```

---

### Get Session Messages

```
GET /sessions/{session_id}/runs
```

Returns all messages in a session.

---

## Legacy's Tools

Legacy has access to these tools (called automatically when needed):

| Tool | Description | Actions |
|------|-------------|---------|
| **todoist** | Task management | get, create, update, complete, delete, projects |
| **github** | Repository operations | read, list, search, write, delete, commit |
| **comms** | Notifications | send (Discord), alert, log |
| **execute_action** | Universal gateway | 138+ actions across services |
| **update_user_memory** | Memory management | Add/update/delete memories |

### Example: Ask Legacy to check tasks

```bash
curl -X POST https://agno.onemindos.dev/agents/legacy-mother-ai/runs \
  -F "message=What tasks do I have today in Todoist?" \
  -F "stream=false"
```

### Example: Ask Legacy to create a task

```bash
curl -X POST https://agno.onemindos.dev/agents/legacy-mother-ai/runs \
  -F "message=Create a task in Todoist: Review Linear setup, due tomorrow, priority 2" \
  -F "stream=false"
```

---

## Multimodal (Vision/Audio)

### Analyze an Image

```bash
curl -X POST https://agno.onemindos.dev/agents/legacy-mother-ai/runs \
  -F "message=What's in this image?" \
  -F "images=@/path/to/image.jpg" \
  -F "stream=false"
```

### Transcribe Audio

```bash
curl -X POST https://agno.onemindos.dev/agents/legacy-mother-ai/runs \
  -F "message=Transcribe this audio" \
  -F "audio=@/path/to/audio.mp3" \
  -F "stream=false"
```

---

## JavaScript/TypeScript Example

```typescript
async function chatWithLegacy(message: string, sessionId?: string) {
  const formData = new FormData();
  formData.append('message', message);
  formData.append('stream', 'false');
  if (sessionId) {
    formData.append('session_id', sessionId);
  }

  const response = await fetch(
    'https://agno.onemindos.dev/agents/legacy-mother-ai/runs',
    {
      method: 'POST',
      body: formData,
    }
  );

  const data = await response.json();
  return data.content;
}

// Usage
const reply = await chatWithLegacy('Hello Legacy');
console.log(reply);
```

---

## Python Example

```python
import requests

def chat_with_legacy(message: str, session_id: str = None) -> str:
    url = "https://agno.onemindos.dev/agents/legacy-mother-ai/runs"

    data = {
        "message": message,
        "stream": "false",
    }
    if session_id:
        data["session_id"] = session_id

    response = requests.post(url, data=data)
    result = response.json()
    return result.get("content", "")

# Usage
reply = chat_with_legacy("Hello Legacy")
print(reply)
```

---

## Error Responses

| Status | Meaning |
|--------|---------|
| 200 | Success |
| 400 | Bad request (invalid parameters) |
| 401 | Unauthorized (if auth enabled) |
| 404 | Agent/session not found |
| 500 | Server error |

**Error Response Format:**
```json
{
  "detail": "Error message here"
}
```

---

## Rate Limits

Currently no rate limits enforced. Be reasonable.

---

## Authentication

Currently open (no auth required).

To add authentication, configure in Agno settings.

---

## Internal Access (On Server)

```bash
# SSH into server
ssh zeus@20.121.38.186

# Test locally
curl -X POST http://localhost:8080/agents/legacy-mother-ai/runs \
  -F "message=Test" \
  -F "stream=false"
```

---

## Related Documentation

- [PROJECT_STATUS.md](./PROJECT_STATUS.md) - Current system status
- [CLAUDE.md](../CLAUDE.md) - Claude session memory
- [CONTEXT7.md](../CONTEXT7.md) - Library documentation reference
- Swagger UI: https://agno.onemindos.dev/docs

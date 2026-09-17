# API Quick Start Guide

## 🚀 Start the Server

```bash
cd /Users/zeuslegacy/Documents/onemind-os
uvicorn backend.agno.onemindos:app --reload --port 8080
```

## 📚 Access Documentation

### 1. Interactive Swagger UI
**Best for testing and exploration**

```
http://localhost:8080/docs
```

Features:
- Try endpoints directly in browser
- Auto-generated from docstrings
- Request/response schemas
- Example values

### 2. ReDoc
**Best for reading and reference**

```
http://localhost:8080/redoc
```

Features:
- Clean, readable layout
- Organized by tags
- Search functionality
- Code examples

### 3. Comprehensive Guide
**Best for integration and workflows**

```
docs/API-DOCUMENTATION.md
```

Features:
- Complete endpoint reference
- Real-world examples
- WebSocket protocols
- Integration guides

## 🎯 Quick Test Commands

### Health Check
```bash
curl http://localhost:8080/health
```

### List Available Models
```bash
curl http://localhost:8080/models?available_only=true
```

### List Tools
```bash
curl http://localhost:8080/tools | jq '.tools[] | {id, category, name}'
```

### List Agents
```bash
curl http://localhost:8080/agents
```

### Run Agent (Streaming)
```bash
curl -N http://localhost:8080/agents/researcher/runs \
  -d "message=What are the latest AI trends?" \
  -d "stream=true"
```

### Voice Health Check
```bash
curl http://localhost:8080/api/voice/health
```

### Activity Stats
```bash
curl http://localhost:8080/api/activity/stats?user_id=zeus
```

### List Workflows
```bash
curl http://localhost:8080/api/workflows
```

## 📖 Common Use Cases

### 1. Chat with an Agent
```bash
# Start a conversation
curl -X POST http://localhost:8080/agents/coder/runs \
  -d "message=Help me debug this Python code" \
  -d "stream=false" \
  -d "session_id=my-session"
```

### 2. Voice Interaction
```bash
# Record audio and convert to base64
AUDIO_BASE64=$(base64 -i recording.wav)

# Send voice request
curl -X POST http://localhost:8080/api/voice/run \
  -H "Content-Type: application/json" \
  -d "{
    \"audio_base64\": \"$AUDIO_BASE64\",
    \"agent_id\": \"legacy-mother-ai\",
    \"voice_id\": \"kore\"
  }"
```

### 3. Run a Workflow
```bash
curl -X POST "http://localhost:8080/api/workflows/daily-briefing/run?stream=false" \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "user": "Zeus"
    }
  }'
```

### 4. Create Activity Item
```bash
curl -X POST http://localhost:8080/api/activity \
  -H "Content-Type: application/json" \
  -d '{
    "type": "task",
    "priority": "high",
    "title": "Review PR",
    "body": "Code review needed",
    "user_id": "zeus"
  }'
```

### 5. Approve HITL Request
```bash
# List pending approvals
curl http://localhost:8080/api/hitl?status=pending

# Approve
curl -X POST http://localhost:8080/api/hitl/{approval_id}/approve \
  -H "Content-Type: application/json" \
  -d '{"reason": "Looks good"}'
```

## 🔧 WebSocket Examples

### Activity Feed (Real-time Updates)
```javascript
const ws = new WebSocket('ws://localhost:8080/api/activity/ws?user_id=zeus');

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('Activity update:', data);
};

// Send ping
ws.send(JSON.stringify({ type: 'ping' }));
```

### Voice Gateway (Streaming Audio)
```javascript
const ws = new WebSocket('ws://localhost:8080/voice/ws/device-123');

// Send audio chunk
ws.send(JSON.stringify({
  type: 'audio_chunk',
  audio_base64: audioData,
  agent_id: 'legacy-mother-ai',
  voice_id: 'kore'
}));

// Receive response
ws.onmessage = (event) => {
  const response = JSON.parse(event.data);
  if (response.type === 'audio_chunk') {
    playAudio(response.audio_base64);
  }
};
```

## 📊 Key Endpoints by Category

### Agents & Teams
- `GET /agents` - List all agents
- `POST /agents/{id}/runs` - Run agent (streaming)
- `GET /teams` - List all teams
- `POST /teams/{id}/runs` - Run team

### Voice
- `POST /api/voice/run` - Voice agent run (HTTP)
- `WS /voice/ws/{device_id}` - Voice streaming (WebSocket)
- `GET /api/voice/voices` - List available voices

### Activity Feed
- `GET /api/activity` - List items
- `POST /api/activity` - Create item
- `WS /api/activity/ws` - Real-time updates

### Approvals (HITL)
- `GET /api/hitl` - List approvals
- `POST /api/hitl/{id}/approve` - Approve request
- `POST /api/hitl/{id}/reject` - Reject request

### Workflows
- `GET /api/workflows` - List workflows
- `POST /api/workflows/{id}/run` - Execute workflow

### Models & Tools
- `GET /models` - List 60+ models
- `GET /tools` - List 40+ tools
- `GET /api/skills` - List skills

## 🎨 Response Formats

### Streaming (SSE)
```
Content-Type: text/event-stream

event: text_message_content
data: {"content": "Hello", "delta": "Hello"}

event: tool_call_start
data: {"tool_name": "search", "arguments": {...}}

event: run_finished
data: {"run_id": "run-123", "status": "complete"}
```

### JSON
```json
{
  "data": {...},
  "status": "success"
}
```

### WebSocket
```json
{
  "type": "activity.item.created",
  "data": {
    "id": "...",
    "type": "task",
    "title": "New task"
  }
}
```

## 🐛 Troubleshooting

### Server Not Starting
```bash
# Check if port is in use
lsof -i :8080

# Kill process if needed
kill -9 <PID>

# Restart server
uvicorn backend.agno.onemindos:app --reload --port 8080
```

### Database Connection Issues
```bash
# Check PostgreSQL is running
pg_isready -h localhost -p 5432

# Check environment variables
echo $DATABASE_URL
```

### API Key Issues
```bash
# Check environment variables
echo $GOOGLE_API_KEY      # For Gemini/voice
echo $AWS_PROFILE          # For Bedrock models
echo $OPENAI_API_KEY       # For OpenAI models
```

## 📖 Further Reading

- **Full API Reference**: `docs/API-DOCUMENTATION.md`
- **Improvements Summary**: `docs/API-IMPROVEMENTS-SUMMARY.md`
- **Project Context**: `CLAUDE.md`
- **Interactive Docs**: http://localhost:8080/docs

## 💡 Tips

1. **Use `?available_only=true`** when listing models to see only configured ones
2. **Set `stream=true`** for real-time agent responses
3. **Use session_id** to maintain conversation context
4. **Check `/health`** to verify server is running
5. **Explore Swagger UI** to try endpoints interactively
6. **Use WebSocket** for real-time updates (activity, voice)

---

**Happy coding!** 🎉


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]

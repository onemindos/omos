# OneMind OS - Activation Checklist

> **Date:** 2026-01-25
> **Goal:** Bring the system to life for daily use

---

## Summary

| Category | Status | Action Required |
|----------|--------|-----------------|
| **Backend API** | LIVE | Deploy code changes |
| **Knowledge/RAG** | FIXED | Router mounted, LanceDB wired |
| **Workflows** | FIXED | Now using LiteLLM (not Bedrock) |
| **Flutter Web** | NEEDS REBUILD | Container exists but needs update |
| **Real-time** | WORKING | WebSocket, Activity, Inbox all functional |
| **Notifications** | CONFIGURED | Novu env vars in compose |

---

## Code Changes Made (Ready to Deploy)

### 1. Knowledge Router Mounted
**File:** `core/orchestrator/main.py`
```python
from core.knowledge.api import knowledge_router
app.include_router(knowledge_router)  # /api/knowledge prefix
```

**New Endpoints:**
- `POST /api/knowledge/upload` - Upload file for RAG
- `POST /api/knowledge/url` - Index URL for RAG
- `POST /api/knowledge/search` - Hybrid search
- `GET /api/knowledge/stats` - Pipeline stats
- `DELETE /api/knowledge/{id}` - Delete document

### 2. LanceDB Storage Wired
**File:** `core/knowledge/api.py`
- Upload endpoint now stores in LanceDB
- URL endpoint now stores in LanceDB
- Search endpoint now uses hybrid search

### 3. Workflows Use LiteLLM
**File:** `core/orchestrator/workflows.py`
- All 4 workflow agents now use `LiteLLMOpenAI`
- Removed direct `AwsBedrock` (was causing permission errors)
- Models: `claude-sonnet`, `claude-haiku`

---

## Deployment Commands

### Step 1: Push Changes
```bash
cd /Users/zeuslegacy/Documents/onemind-os

# Stage all changes
git add core/orchestrator/main.py
git add core/orchestrator/workflows.py
git add core/knowledge/

# Commit
git commit -m "feat: activate knowledge RAG and fix workflow model routing

- Mount knowledge_router at /api/knowledge
- Wire LanceDB storage for uploads and URL indexing
- Implement hybrid search endpoint
- Switch all workflow agents to LiteLLM (fixes Bedrock permission errors)"

# Push
git push origin main
```

### Step 2: Deploy to Server
```bash
# SSH to server
ssh zeus@20.121.38.186

# Pull changes
cd ~/onemind-os && git pull

# Rebuild and restart Agno container
cd infra/compose
docker compose -f docker-compose.core.yml build agno
docker compose -f docker-compose.core.yml up -d agno

# Check logs
docker logs onemind-agno -f --tail 100
```

### Step 3: Rebuild Flutter Web
```bash
# On local machine
cd /Users/zeuslegacy/Documents/onemind-os/ui
flutter build web --release

# Deploy to server
scp -r build/web/* ubuntu@100.102.21.44:/tmp/flutter-web/
ssh ubuntu@100.102.21.44 "docker cp /tmp/flutter-web/. onemind-flutter-web:/usr/share/nginx/html/ && docker restart onemind-flutter-web"
```

---

## Test Commands (After Deployment)

### Test Knowledge Upload
```bash
# Upload a PDF
curl -X POST http://100.102.21.44:8081/api/knowledge/upload \
  -F "file=@document.pdf" \
  -F "name=Test Document"

# Index a URL
curl -X POST http://100.102.21.44:8081/api/knowledge/url \
  -H "Content-Type: application/json" \
  -d '{"url": "https://docs.agno.com", "name": "Agno Docs"}'

# Search
curl -X POST http://100.102.21.44:8081/api/knowledge/search \
  -H "Content-Type: application/json" \
  -d '{"query": "How to create an agent?", "num_results": 5}'
```

### Test Workflows
```bash
# Morning Briefing
curl -X POST http://100.102.21.44:8081/api/workflows/morning-briefing/run

# Task Triage
curl -X POST http://100.102.21.44:8081/api/workflows/task-triage/run

# Focus Mode
curl -X POST http://100.102.21.44:8081/api/workflows/focus-mode/run

# Evening Review
curl -X POST http://100.102.21.44:8081/api/workflows/evening-review/run
```

### Test Real-time
```bash
# Activity Stream
curl http://100.102.21.44:8081/api/activity

# Inbox Items
curl http://100.102.21.44:8081/api/inbox

# Send test notification (triggers WebSocket)
curl -X POST http://100.102.21.44:8081/api/inbox \
  -H "Content-Type: application/json" \
  -d '{"type": "notification", "title": "System Online", "body": "OneMind is now active!", "user_id": "zeus"}'
```

---

## What's Now Functional

### Backend APIs (Verified Working)
| Endpoint | Status | Description |
|----------|--------|-------------|
| `GET /health` | LIVE | Backend health |
| `GET /agents` | LIVE | List agents |
| `POST /agents/{id}/runs` | LIVE | Chat with agent |
| `GET /sessions` | LIVE | List sessions |
| `GET /awareness` | LIVE | Get awareness state |
| `PUT /awareness` | LIVE | Set awareness mode |
| `GET /api/workflows` | LIVE | List workflows |
| `POST /api/workflows/{id}/run` | FIXED | Run workflow |
| `GET /api/inbox` | LIVE | Get inbox items |
| `WS /api/inbox/ws` | LIVE | Real-time inbox |
| `GET /api/activity` | LIVE | Activity stream |
| `POST /api/knowledge/upload` | NEW | Upload documents |
| `POST /api/knowledge/search` | NEW | RAG search |

### Flutter Screens (Verified Connected)
| Screen | Provider | Status |
|--------|----------|--------|
| Chat | `chat_provider.dart` | LIVE |
| Inbox | `unified_inbox_provider.dart` | LIVE |
| Memory | `memory_provider.dart` | LIVE |
| Sessions | `sessions_provider.dart` | LIVE |
| Agents | `agents_provider.dart` | LIVE |
| Awareness | `awareness_provider.dart` | LIVE |
| Activity | `activity_provider.dart` | LIVE |
| Knowledge | `knowledge_provider.dart` | READY (needs deploy) |
| Workflows | `workflows_provider.dart` | READY |

---

## Daily Usage Workflow

### Morning
1. Open app at `http://100.102.21.44:3005`
2. Run "Morning Briefing" workflow
3. Review inbox for overnight notifications
4. Check activity stream for system health

### During Day
1. Chat with Legacy for tasks
2. Use Knowledge to search your documents
3. Trigger workflows as needed (focus mode, triage)
4. Monitor activity for tool executions

### Evening
1. Run "Evening Review" workflow
2. Clear inbox items
3. Review activity for issues

---

## Environment Variables (Already Configured in Compose)

```bash
# Required - Already set
DATABASE_URL=postgresql://onemind:xxx@postgres:5432/onemind
REDIS_URL=redis://redis:6379/0
NATS_URL=nats://nats:4222
USE_LITELLM=true
LITELLM_URL=http://litellm:4000
LITELLM_MASTER_KEY=sk-onemind-master-key

# Notifications - Already set
NOVU_API_KEY=xxx
NOVU_APP_ID=xxx
NOVU_SECRET_KEY=xxx

# Integrations - Already set
TODOIST_API_TOKEN=xxx
GITHUB_TOKEN=xxx
DISCORD_WEBHOOK_URL=xxx
```

---

## Known Issues & Workarounds

### 1. LiteLLM Health Check
LiteLLM returns 401 without API key in health check. This is normal - it works when accessed through the Agno container which passes the key.

### 2. Flutter Web Not Reachable
The Flutter web container needs rebuild. Use the deployment commands above.

### 3. LanceDB Dependencies
If LanceDB fails, install: `pip install lancedb tantivy pyarrow`
This should be in the Docker image already.

### 4. Knowledge RAG Embeddings
Embeddings use **AWS Bedrock Titan** (`amazon.titan-embed-text-v2:0`).
- No OpenAI API key needed
- Uses existing AWS credentials
- Cost: $0.02/1M tokens (same as OpenAI)

---

## Next Steps After Activation

1. **Upload your key documents** to Knowledge for RAG
2. **Configure workflows** with your schedule preferences
3. **Set up mobile access** (Flutter iOS/Android builds)
4. **Enable push notifications** via Novu if not already
5. **Connect more integrations** (calendar, email, etc.)

---

**System is ready for daily use once deployed!**


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]

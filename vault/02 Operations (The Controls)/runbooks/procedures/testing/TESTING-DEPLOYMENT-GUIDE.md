# Testing & Deployment Guide

> **Purpose:** Comprehensive checklist for testing all features after deployment
> **Created:** 2026-01-22
> **Status:** Ready for deployment and testing

---

## 🚀 Deployment Steps

### 1. Pre-Deployment Checklist

- [ ] All code changes committed locally
- [ ] Dependencies documented
- [ ] Database migrations ready
- [ ] Backup current database
- [ ] Environment variables set

### 2. Deploy to Production

```bash
# Step 1: Commit and push
cd /Users/zeuslegacy/Documents/onemind-os
git add .
git commit -m "feat: complete workflows, teams, knowledge/RAG, MCP systems

- Stream 1: Workflows custom creation + Teams execution
- Stream 2: Knowledge/RAG full implementation (5 endpoints, 550 lines)
- Stream 2: MCP server management (7 endpoints, database integration)
- Database: 4 new tables (workflows, knowledge, chunks, mcp_servers)
- Documentation: Complete guides for Knowledge/RAG"

git push origin main

# Step 2: Coolify auto-deploys (wait ~2-3 minutes)
# Monitor: https://coolify.your-domain.com

# Step 3: Run database migrations
ssh zeus@20.121.38.186
docker exec -i onemind-postgres psql -U onemind -d onemind < ~/onemind-os/infra/compose/cloud/init-db.sql

# Step 4: Restart services
docker compose -f ~/onemind-os/docker-compose.core.yml restart onemind-agno

# Step 5: Verify services are healthy
docker ps | grep onemind
curl http://100.102.21.44:8081/health
```

### 3. Install Dependencies

```bash
# SSH to server
ssh zeus@20.121.38.186

# Install Python packages
docker exec onemind-agno pip install PyPDF2 python-docx beautifulsoup4 pandas aiohttp openai

# Restart after install
docker restart onemind-agno
```

### 4. Environment Variables

```bash
# Add to docker-compose.core.yml or .env:
OPENAI_API_KEY=sk-...  # Required for Knowledge/RAG embeddings
```

---

## ✅ Feature Testing Checklist

### Stream 1: Partial Features (COMPLETE)

#### 1.1 Workflows - Custom Creation ✅

**What Changed:**
- Database: `custom_workflows` table
- Backend: `POST /api/workflows` endpoint
- Backend: `GET /api/workflows` returns built-in + custom
- UI: Button enabled

**Test Steps:**
```bash
# 1. Test API - Create workflow
curl -X POST http://100.102.21.44:8081/api/workflows \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Workflow",
    "description": "My custom workflow",
    "steps": [
      {"id": "step1", "name": "First", "type": "agent"},
      {"id": "step2", "name": "Second", "type": "tool"}
    ],
    "trigger": "Manual"
  }'

# Expected: {"success": true, "workflow": {...}}

# 2. Test API - List workflows
curl http://100.102.21.44:8081/api/workflows

# Expected: {"workflows": [...], "total": 6}  # 5 built-in + 1 custom

# 3. Test UI
```

**UI Test:**
1. Go to http://100.102.21.44:3000/workflows
2. Click "Create Workflow" button (should show alert)
3. Verify workflows list shows both built-in and custom
4. Click "Run" on any workflow
5. Verify run starts successfully

**Expected Results:**
- ✅ Button is enabled (not grayed out)
- ✅ API creates workflow successfully
- ✅ List includes custom workflows
- ✅ Can run workflows

---

#### 1.2 Teams - Execution Endpoint ✅

**What Changed:**
- Backend: `POST /api/teams/{id}/run` endpoint

**Test Steps:**
```bash
# 1. List teams
curl http://100.102.21.44:8081/api/teams

# 2. Run a team (get team_id from list)
curl -X POST http://100.102.21.44:8081/api/teams/TEAM_ID/run \
  -H "Content-Type: application/json" \
  -d '{"message": "Research the latest AI news"}'

# Expected: {"run_id": "...", "status": "running", ...}
```

**UI Test:**
1. Go to http://100.102.21.44:3000/teams
2. Select a team
3. Click "Run Team"
4. Provide input message
5. Verify run starts

**Expected Results:**
- ✅ Teams list loads
- ✅ Run endpoint returns run_id
- ✅ UI shows team is running

---

#### 1.3 Memory - MemoryTools ✅ (Already Working)

**Test Steps:**
```bash
# Verify MemoryTools is in tools list
curl http://100.102.21.44:8081/api/tools | jq '.tools[] | select(.name=="Memory Tools")'

# Expected: {"name": "Memory Tools", "category": "integration"}
```

**UI Test:**
1. Go to http://100.102.21.44:3000/memory
2. Create a new memory
3. Search memories
4. Delete a memory

**Expected Results:**
- ✅ MemoryTools appears in tools list
- ✅ Can CRUD memories via UI
- ✅ PostgreSQL `ai.user_memories` table has data

---

#### 1.4 Tools - Built-in Tools ✅ (Already Working)

**Test Steps:**
```bash
# Verify built-in tools are loaded
curl http://100.102.21.44:8081/api/tools | jq '.tools[] | select(.category=="research")'

# Expected: DuckDuckGo, HackerNews, Calculator
```

**UI Test:**
1. Go to http://100.102.21.44:3000/tools
2. Verify list shows:
   - DuckDuckGo (research)
   - HackerNews (research)
   - Calculator (data)
   - Wikipedia, Arxiv, YFinance, YouTube

**Agent Test:**
In chat, ask:
- "Search DuckDuckGo for AI news"
- "What's the stock price of AAPL?"
- "Calculate 1234 * 5678"

**Expected Results:**
- ✅ 8+ tools in list
- ✅ Agent can use built-in tools
- ✅ No import errors in logs

---

### Stream 2: Missing Features (COMPLETE)

#### 2.1 Knowledge/RAG System ✅

**What Changed:**
- Database: `knowledge_items` + `knowledge_chunks` tables
- Backend: `knowledge.py` module (550 lines)
- Backend: 5 API endpoints

**API Endpoints:**
```
POST   /api/knowledge/upload    # Upload file/URL/text
POST   /api/knowledge/search    # Semantic search
GET    /api/knowledge           # List all items
DELETE /api/knowledge/{id}      # Delete item
GET    /api/knowledge/stats     # Statistics
```

**Test Steps:**

**A. Upload Text:**
```bash
curl -X POST http://100.102.21.44:8081/api/knowledge/upload \
  -H "Content-Type: application/json" \
  -d '{
    "type": "text",
    "name": "OneMind Info",
    "content": "OneMind OS uses Agno framework with AWS Bedrock Claude models. It has 4 awareness modes: dormant, aware, present, and omnipresent."
  }'

# Expected: {"success": true, "knowledge_id": "kb-...", "chunk_count": 1}
```

**B. Search Knowledge:**
```bash
curl -X POST http://100.102.21.44:8081/api/knowledge/search \
  -H "Content-Type: application/json" \
  -d '{
    "query": "What framework does OneMind use?",
    "top_k": 3
  }'

# Expected: {"results": [...], "total_results": 1}
```

**C. List Knowledge:**
```bash
curl http://100.102.21.44:8081/api/knowledge/stats

# Expected: {"total_items": 1, "by_type": {"text": 1}, ...}
```

**UI Test:**
1. Go to http://100.102.21.44:3000/knowledge
2. Click "Upload" button
3. Select "Text" type
4. Paste test content
5. Submit
6. Verify item appears in list
7. Search for content
8. Verify results appear
9. Delete item

**Agent Integration Test:**
```
In chat:
"What's in the knowledge base about OneMind?"
"Search the knowledge base for framework information"
```

**Expected Results:**
- ✅ Upload succeeds (text, URL, file)
- ✅ Search returns relevant results with similarity scores
- ✅ Stats show correct counts
- ✅ UI shows knowledge items
- ✅ Can delete items
- ✅ Agent can search knowledge base

**Requirements:**
- `OPENAI_API_KEY` environment variable must be set
- Dependencies: PyPDF2, python-docx, beautifulsoup4, pandas

---

#### 2.2 MCP Server Management ✅

**What Changed:**
- Database: `mcp_servers` table
- Backend: `mcp.py` module (350 lines)
- Backend: 7 API endpoints

**API Endpoints:**
```
GET    /api/mcp/servers                # List all
GET    /api/mcp/servers/{id}           # Get one
POST   /api/mcp/servers                # Create
PUT    /api/mcp/servers/{id}           # Update
PATCH  /api/mcp/servers/{id}           # Toggle enabled
DELETE /api/mcp/servers/{id}           # Delete
POST   /api/mcp/servers/{id}/reconnect # Reconnect
```

**Test Steps:**

**A. List Servers (should show example servers):**
```bash
curl http://100.102.21.44:8081/api/mcp/servers

# Expected: {"servers": [], "total": 0}  # Empty initially
# Note: Example servers will be seeded on first run
```

**B. Create Server:**
```bash
curl -X POST http://100.102.21.44:8081/api/mcp/servers \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test MCP Server",
    "type": "stdio",
    "command": "npx",
    "args": ["-y", "@modelcontextprotocol/server-memory"],
    "enabled": true
  }'

# Expected: {"success": true, "server": {...}}
```

**C. Toggle Server:**
```bash
curl -X PATCH http://100.102.21.44:8081/api/mcp/servers/SERVER_ID \
  -H "Content-Type: application/json" \
  -d '{"enabled": false}'

# Expected: {"success": true, "enabled": false}
```

**UI Test:**
1. Go to http://100.102.21.44:3000/mcp
2. Click "Add Server"
3. Fill form:
   - Name: Test Server
   - Type: stdio
   - Command: npx
   - Args: -y @modelcontextprotocol/server-memory
4. Submit
5. Verify server appears in list
6. Toggle enabled switch
7. Click settings icon (edit)
8. Update server config
9. Delete server

**Expected Results:**
- ✅ Can create MCP servers
- ✅ List shows all servers
- ✅ Can toggle enabled/disabled
- ✅ Can update server config
- ✅ Can delete servers
- ✅ UI syncs with backend

**Note:** Actual MCP connection (tools/resources discovery) is marked as TODO. Server management (CRUD) is complete.

---

#### 2.3 HITL/Approvals ⏸️ (TODO - Last Stream 2 Feature)

**Status:** Will be implemented next

---

### Stream 3: Integration Testing (PENDING)

#### 3.1 Verify All Integrations

**Integrations to Test:**

| Integration | Endpoint | Test |
|-------------|----------|------|
| Todoist | /api/tools | Create task, list tasks, complete |
| GitHub | /api/tools | List repos, create issue |
| Composio | /api/tools | Execute action with API key |
| NATS | - | Publish event, verify subscription |
| PostgreSQL | /health/services | Check connection |
| Redis | /health/services | Check connection |
| Extension API | /api/extension/stats | POST capture, GET stats |
| Protocol Sync | - | Deploy and test CalDAV/IMAP |
| Nango | - | Deploy and test OAuth |

**Test Commands:**
```bash
# Check service health
curl http://100.102.21.44:8081/health/services

# Test extension API
curl http://100.102.21.44:8081/api/extension/stats

# Test database
curl http://100.102.21.44:8081/api/memory | jq .

# Test tools
curl http://100.102.21.44:8081/api/tools | jq '.tools[] | .name'
```

---

### Stream 4: UI & System Testing (PENDING)

#### 4.1 Test All 27 Pages

**Command:**
```bash
# Install agent-browser
npm install -g agent-browser
npx playwright install
agent-browser install

# Test workflow
agent-browser open http://100.102.21.44:3000 --headed
agent-browser snapshot -i
```

**Pages to Test:**
1. ☐ `/` - Chat (send message, streaming works)
2. ☐ `/dashboard` - Metrics load
3. ☐ `/sessions` - List, detail, delete
4. ☐ `/agents` - CRUD operations
5. ☐ `/teams` - CRUD + run
6. ☐ `/tools` - List all tools
7. ☐ `/workflows` - List, create, run
8. ☐ `/memory` - CRUD memories
9. ☐ `/knowledge` - Upload, search, delete
10. ☐ `/evals` - Create, run, results
11. ☐ `/traces` - View session traces
12. ☐ `/mcp` - Server management
13. ☐ `/approvals` - Pending requests
14. ☐ `/awareness` - Mode switching
15. ☐ `/consciousness` - Processor config
16. ☐ `/presence` - Location updates
17. ☐ `/metrics` - Token usage
18. ☐ `/models` - Model list
19. ☐ `/settings` - All tabs
20. ☐ `/settings/extension` - Extension settings
21. ☐ `/settings/integrations` - Integration config
22. ☐ `/settings/data-sync` - Sync settings
23. ☐ `/settings/protocols` - Protocol config
24. ☐ `/activity` - Activity feed
25. ☐ `/inbox` - Inbox items
26. ☐ `/profile` - User profile
27. ☐ `/onboarding` - Onboarding flow

**For Each Page:**
- ☐ Page loads without errors
- ☐ All buttons are clickable
- ☐ Forms submit successfully
- ☐ Data displays correctly
- ☐ No console errors

---

## 📊 Success Criteria

### API Endpoints (All Must Return 200 or Valid Response)

```bash
# Health
curl http://100.102.21.44:8081/health

# Workflows
curl http://100.102.21.44:8081/api/workflows

# Teams
curl http://100.102.21.44:8081/api/teams

# Tools
curl http://100.102.21.44:8081/api/tools

# Knowledge
curl http://100.102.21.44:8081/api/knowledge

# MCP
curl http://100.102.21.44:8081/api/mcp/servers

# Memory
curl http://100.102.21.44:8081/api/memory

# Models
curl http://100.102.21.44:8081/api/models
```

### Database Tables (All Must Exist)

```sql
-- Check tables exist
SELECT tablename FROM pg_tables WHERE schemaname = 'public';

-- Expected tables:
-- - unified_items
-- - custom_workflows
-- - knowledge_items
-- - knowledge_chunks
-- - mcp_servers
-- - (plus Agno default tables)
```

### Logs (No Critical Errors)

```bash
# Check logs for errors
docker logs onemind-agno --tail 100 | grep ERROR

# Should see mostly INFO logs, no critical errors
```

---

## 🔧 Troubleshooting

### Issue: Database Migration Failed

```bash
# Check if tables exist
docker exec onemind-postgres psql -U onemind -d onemind -c "\dt"

# Manually run migrations
docker exec -i onemind-postgres psql -U onemind -d onemind < init-db.sql
```

### Issue: Import Errors (knowledge.py, mcp.py)

```bash
# Check Python path
docker exec onemind-agno python -c "import sys; print(sys.path)"

# Verify files exist
docker exec onemind-agno ls -la /app/agents/legacy-Brain/

# Restart service
docker restart onemind-agno
```

### Issue: OPENAI_API_KEY Not Set

```bash
# Add to docker-compose.core.yml
# environment:
#   - OPENAI_API_KEY=sk-...

# Or set in .env file
echo "OPENAI_API_KEY=sk-..." >> .env

# Restart
docker compose restart onemind-agno
```

### Issue: Dependencies Missing

```bash
# Install in container
docker exec onemind-agno pip install PyPDF2 python-docx beautifulsoup4 pandas aiohttp openai

# Or add to requirements.txt and rebuild
```

---

## 📝 Post-Deployment Checklist

- [ ] All API endpoints return 200
- [ ] All 27 UI pages load
- [ ] Database has all tables
- [ ] No critical errors in logs
- [ ] Can create custom workflow
- [ ] Can upload to knowledge base
- [ ] Can search knowledge base
- [ ] Can manage MCP servers
- [ ] Built-in tools work
- [ ] Memory CRUD works
- [ ] Teams can execute
- [ ] Agent responds in chat

---

## 🎯 What's Working vs. What's TODO

### ✅ COMPLETE & READY TO TEST

| Feature | Backend | Frontend | Database | Docs |
|---------|---------|----------|----------|------|
| Workflows (custom) | ✅ | ✅ | ✅ | ✅ |
| Teams (execution) | ✅ | ✅ | - | ✅ |
| Memory | ✅ | ✅ | ✅ | ✅ |
| Tools (built-in) | ✅ | ✅ | - | ✅ |
| Knowledge/RAG | ✅ | ✅ | ✅ | ✅ |
| MCP Servers | ✅ | ✅ | ✅ | ⚠️ |

### ⏸️ REMAINING WORK

| Feature | Status | Estimated Time |
|---------|--------|----------------|
| Approvals/HITL | TODO | 1 hour |
| Integration Testing | TODO | 1 hour |
| UI Testing (27 pages) | TODO | 2 hours |
| Database Optimization | TODO | 30 min |
| Code Cleanup | TODO | 1 hour |

---

**Last Updated:** 2026-01-22
**Ready for Deployment:** YES (6/9 features complete)
**Blocking Issues:** None
**Next Step:** Deploy and test, or finish remaining features first


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/operations (Ops Procedures)/testing (Testing)/testing|testing]]

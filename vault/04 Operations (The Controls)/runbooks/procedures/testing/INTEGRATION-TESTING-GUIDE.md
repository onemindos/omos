# Integration Testing Guide - OneMind OS

> **Purpose:** Verify all integration connections, data flow, and unified data store
> **Status:** Stream 3 - Integration & Data Flow
> **Date:** 2026-01-22

---

## Integration Architecture Overview

OneMind OS uses a **6-layer integration architecture**:

1. **Standard Protocols** - CalDAV, IMAP, RSS, WebDAV
2. **AI Actions** - Composio (250+ APIs)
3. **Data Sync** - Nango (OAuth management)
4. **Browser Extension** - Custom capture from any site
5. **Unified Store** - PostgreSQL `unified_items` table
6. **Events** - NATS JetStream for real-time updates

**Data Flow:**
```
Integration → Capture → Unified Store → NATS Event → Agent Awareness
```

---

## Integration Categories

### Core Infrastructure (Always Running)
| Integration | Type | Status | Test Method |
|-------------|------|--------|-------------|
| PostgreSQL | Direct | ✅ Working | Query DB |
| Redis | Direct | ✅ Working | Set/Get keys |
| NATS | Direct | ✅ Working | Pub/Sub test |
| Extension API | REST | ✅ Working | POST capture |

### External Services (Require API Keys)
| Integration | Type | Status | Test Method |
|-------------|------|--------|-------------|
| Composio | SDK | ⚠️ Needs API key | List apps |
| Todoist | Composio | ⚠️ Needs connection | Create task |
| GitHub | Composio | ⚠️ Needs connection | List repos |
| Nango | Service | ⏸️ Not deployed | Deploy first |

### Protocol Sync (Local)
| Integration | Type | Status | Test Method |
|-------------|------|--------|-------------|
| CalDAV | Protocol | ⏸️ Not deployed | Deploy sync service |
| IMAP | Protocol | ⏸️ Not deployed | Deploy sync service |
| RSS | Protocol | ⏸️ Not deployed | Deploy sync service |

### Optional Services
| Integration | Type | Status | Test Method |
|-------------|------|--------|-------------|
| Home Assistant | Direct API | ⚠️ Needs URL | GET entities |
| Frigate | MQTT | ⚠️ Needs MQTT | Subscribe events |
| Spotify | Composio | ⚠️ Needs connection | Get playing |
| Discord | Webhook | ⚠️ Needs URL | POST webhook |
| Slack | Composio | ⚠️ Needs connection | Send message |

---

## Test Plan - Stream 3

### Phase 1: Core Infrastructure (10 min)

**Test 1: PostgreSQL Connection**
```bash
# From inside onemind-agno container
curl http://100.102.21.44:8081/health
# Should show: database: healthy
```

**Test 2: Redis Connection**
```bash
# Test via Python
docker exec onemind-agno python3 -c "
import redis
r = redis.Redis(host='onemind-redis', port=6379)
r.set('test', 'hello')
print('Redis:', r.get('test').decode())
"
```

**Test 3: NATS Connection**
```bash
# Test pub/sub
curl -X POST http://100.102.21.44:8081/nats/publish \
  -H "Content-Type: application/json" \
  -d '{"subject": "test.integration", "data": {"test": true}}'

# Check NATS dashboard
curl http://100.102.21.44:8222/varz
```

**Test 4: Extension API**
```bash
# Capture test data
curl -X POST http://100.102.21.44:8081/api/extension/capture \
  -H "Content-Type: application/json" \
  -H "X-Source: onemind-extension" \
  -d '{
    "type": "test",
    "url": "https://example.com",
    "data": {"test": "integration"}
  }'

# Get stats
curl http://100.102.21.44:8081/api/extension/stats
```

**Expected Results:**
- ✅ All 4 core integrations respond
- ✅ Data persists to PostgreSQL
- ✅ NATS events published
- ✅ No errors in logs

---

### Phase 2: Integrations API (10 min)

**Test 5: List All Integrations**
```bash
curl http://100.102.21.44:8081/api/integrations | jq .
```

**Expected Response:**
```json
{
  "categories": [
    {
      "id": "productivity",
      "name": "Productivity",
      "integrations": [...]
    },
    {
      "id": "developer",
      "name": "Developer",
      "integrations": [...]
    },
    ...
  ],
  "total": 14,
  "connected": 0
}
```

**Test 6: Get Integration Details**
```bash
# Get Todoist details
curl http://100.102.21.44:8081/api/integrations/todoist | jq .

# Get GitHub details
curl http://100.102.21.44:8081/api/integrations/github | jq .
```

**Expected:** Full integration config with webhook setup instructions

---

### Phase 3: Composio Integration (15 min)

**Prerequisites:**
```bash
# Check if Composio SDK installed
docker exec onemind-agno python3 -c "import composio; print('Composio installed')"

# Check API key
docker exec onemind-agno printenv | grep COMPOSIO
```

**Test 7: List Available Apps**
```bash
# Via Python
docker exec onemind-agno python3 -c "
from composio import Composio
client = Composio()
apps = client.apps.list()
print(f'Available apps: {len(apps)}')
for app in apps[:5]:
    print(f'  - {app.name}')
"
```

**Test 8: Connect Todoist**
```bash
# Initiate OAuth
curl -X POST http://100.102.21.44:8081/api/integrations/todoist/connect

# Follow OAuth URL
# After callback:
curl http://100.102.21.44:8081/api/integrations/todoist
# Should show: status: connected
```

**Test 9: Test Todoist Actions**
```bash
# Create task via Composio
curl -X POST http://100.102.21.44:8081/api/integrations/todoist/actions/create_task \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Test task from OneMind",
    "priority": 3
  }'

# List tasks
curl http://100.102.21.44:8081/api/integrations/todoist/actions/get_tasks
```

---

### Phase 4: Unified Data Store (10 min)

**Test 10: Verify Table Structure**
```bash
# Connect to PostgreSQL
docker exec -it onemind-postgres psql -U postgres -d onemind

# Check table
SELECT * FROM unified_items LIMIT 10;

# Count by source
SELECT source, COUNT(*) FROM unified_items GROUP BY source;

# Count by type
SELECT type, COUNT(*) FROM unified_items GROUP BY type;
```

**Expected Columns:**
- `id` (VARCHAR)
- `source` (VARCHAR) - 'extension', 'caldav', 'imap', 'nango'
- `type` (VARCHAR) - 'receipt', 'order', 'article', 'email', 'event'
- `external_id` (VARCHAR)
- `data` (JSONB)
- `synced_at` (TIMESTAMPTZ)
- `created_at` (TIMESTAMPTZ)

**Test 11: Insert Test Data**
```sql
-- Insert test item
INSERT INTO unified_items (id, source, type, external_id, data)
VALUES (
  'test-integration-001',
  'test',
  'sample',
  'ext-123',
  '{"test": true, "message": "Integration test"}'
);

-- Verify insert
SELECT * FROM unified_items WHERE source = 'test';
```

**Test 12: NATS Events for Data Store**
```bash
# Subscribe to unified store events
docker exec onemind-nats nats sub "events.unified.*"

# In another terminal, insert data
curl -X POST http://100.102.21.44:8081/api/extension/capture \
  -H "Content-Type: application/json" \
  -d '{"type": "test", "data": {"test": "nats_event"}}'

# Should see event in subscriber
```

---

### Phase 5: Protocol Sync Services (Deployment Required)

**Services to Deploy:**
- **protocol-sync** - CalDAV, IMAP, RSS sync
- **nango** - OAuth token management

**Deployment:**
```bash
# Deploy protocol-sync
cd /Users/zeuslegacy/Documents/onemind-os
docker compose --profile data up -d protocol-sync

# Deploy nango
docker compose --profile data up -d nango

# Verify running
docker ps | grep -E 'protocol-sync|nango'
```

**Test 13: Protocol Sync - CalDAV**
```bash
# Configure CalDAV source
curl -X POST http://100.102.21.44:8083/api/sources \
  -H "Content-Type: application/json" \
  -d '{
    "type": "caldav",
    "url": "https://caldav.icloud.com",
    "username": "user@icloud.com",
    "password": "app-specific-password"
  }'

# Trigger sync
curl -X POST http://100.102.21.44:8083/api/sync/caldav

# Verify in unified_items
docker exec -it onemind-postgres psql -U postgres -d onemind -c \
  "SELECT * FROM unified_items WHERE source = 'caldav' LIMIT 5;"
```

**Test 14: Protocol Sync - IMAP**
```bash
# Configure IMAP source
curl -X POST http://100.102.21.44:8083/api/sources \
  -H "Content-Type: application/json" \
  -d '{
    "type": "imap",
    "host": "imap.gmail.com",
    "port": 993,
    "username": "user@gmail.com",
    "password": "app-password"
  }'

# Trigger sync
curl -X POST http://100.102.21.44:8083/api/sync/imap

# Verify in unified_items
docker exec -it onemind-postgres psql -U postgres -d onemind -c \
  "SELECT * FROM unified_items WHERE source = 'imap' LIMIT 5;"
```

**Test 15: Nango OAuth**
```bash
# Check Nango health
curl http://100.102.21.44:3004/api/health

# List available integrations
curl http://100.102.21.44:3004/api/integrations

# Connect Google Calendar (example)
curl -X POST http://100.102.21.44:3004/api/connections \
  -H "Content-Type: application/json" \
  -d '{
    "integration": "google-calendar",
    "user_id": "default"
  }'
```

---

## Success Criteria

### ✅ Core Infrastructure (Required)
- [ ] PostgreSQL connection healthy
- [ ] Redis SET/GET working
- [ ] NATS pub/sub working
- [ ] Extension API captures data
- [ ] Data persists to `unified_items`
- [ ] NATS events published for captures

### ⚠️ External Services (Optional - Requires Keys)
- [ ] Composio SDK installed
- [ ] At least 1 Composio integration connected (Todoist or GitHub)
- [ ] Integration actions work (create task, list repos)

### ⏸️ Protocol Sync (Deploy Later)
- [ ] Protocol-sync service deployed
- [ ] CalDAV sync working (if credentials available)
- [ ] IMAP sync working (if credentials available)
- [ ] Data flows to `unified_items`

### ⏸️ Nango (Deploy Later)
- [ ] Nango service deployed
- [ ] OAuth connections working
- [ ] Token refresh automatic

---

## Common Issues & Fixes

### Issue: Composio SDK not installed
```bash
# Install in container
docker exec onemind-agno pip install composio-core
docker restart onemind-agno
```

### Issue: No Composio API key
```bash
# Get free key: https://composio.dev
# Add to .env or docker-compose.yml:
COMPOSIO_API_KEY=your_key_here
```

### Issue: PostgreSQL connection failed
```bash
# Check if postgres running
docker ps | grep postgres

# Check connection string
docker exec onemind-agno printenv | grep DATABASE_URL

# Test connection
docker exec onemind-postgres psql -U postgres -c "SELECT 1;"
```

### Issue: NATS not receiving events
```bash
# Check NATS server status
curl http://100.102.21.44:8222/varz

# Check subscriptions
curl http://100.102.21.44:8222/subsz

# Restart NATS
docker restart onemind-nats
```

### Issue: unified_items table not found
```bash
# Run database migration
docker exec -it onemind-postgres psql -U postgres -d onemind \
  < /path/to/init-db.sql

# Or manually create:
docker exec -it onemind-postgres psql -U postgres -d onemind
# Then paste SQL from infra/compose/cloud/init-db.sql
```

---

## Testing Checklist - Quick Reference

### Minimum Viable (15 min)
- [ ] Health check passes
- [ ] PostgreSQL query works
- [ ] Redis SET/GET works
- [ ] NATS pub/sub works
- [ ] Extension API captures
- [ ] List integrations API works

### Full Test (30 min)
- [ ] All minimum viable tests
- [ ] Get integration details
- [ ] Unified store has data
- [ ] NATS events flowing
- [ ] 1+ Composio integration connected
- [ ] Integration action executed

### Production Ready (1+ hour)
- [ ] All full tests
- [ ] Protocol-sync deployed & tested
- [ ] Nango deployed & tested
- [ ] CalDAV sync working
- [ ] IMAP sync working
- [ ] All NATS events verified
- [ ] Load test with 100+ items

---

## Next Steps After Testing

1. **Document Results** - Update IMPLEMENTATION-PROGRESS.md with test outcomes
2. **Fix Issues** - Address any failing tests
3. **Deploy Missing Services** - protocol-sync, nango if needed
4. **Stream 4** - Move to UI testing with agent-browser

---

**Last Updated:** 2026-01-22 05:30 UTC
**Test Duration:** 15-60 minutes depending on depth
**Prerequisites:** OneMind OS running on `100.102.21.44`


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/operations (Ops Procedures)/testing (Testing)/testing|testing]]

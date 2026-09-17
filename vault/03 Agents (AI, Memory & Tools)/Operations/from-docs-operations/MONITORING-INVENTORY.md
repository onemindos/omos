# OneMind OS - Monitoring & Observability Inventory

> Last Updated: 2026-01-30
> Complete inventory of all components requiring monitoring, alerting, and observability

---

## Observability Platform Decision

| Decision | **Datadog** |
|----------|-------------|
| Status | PLANNED |
| Issue | [OMOS-486 - Implement Datadog Observability Platform](https://linear.app/onemind-os/issue/OMOS-486/implement-datadog-observability-platform) |
| Full Plan | [DATADOG-OBSERVABILITY-PLAN.md](./DATADOG-OBSERVABILITY-PLAN.md) |
| Estimated Cost | ~$35-50/month |

### Why Datadog?

| Requirement | Datadog Score |
|-------------|---------------|
| Plug-and-play setup | 10/10 |
| Single pane of glass | 10/10 |
| AI-assisted debugging | 10/10 |
| Native NATS support | Only platform with this |
| LLM Observability | Built-in token/cost tracking |
| FastAPI support | Auto-instrumentation |

### What Changes After Implementation

| Tool | Status |
|------|--------|
| Netdata | REMOVE (fully replaced) |
| Portainer | KEEP (management, not monitoring) |
| Manual log checking | REMOVE |

---

## Quick Reference

| Category | Count | Status |
|----------|-------|--------|
| Infrastructure Services | 8 | Active |
| Backend API Routers | 19+ | Active |
| API Endpoints | 150+ | Active |
| NATS Event Handlers | 10+ | Active |
| External API Integrations | 19+ | Active |
| Agent Teams | 10+ | Active |
| Individual Agents | 6+ | Active |
| LifeOS Modules | 15+ | Active |
| Critical User Paths | 5 | Monitored |

---

## 1. Infrastructure Services

### Docker Services (Production)

| Service | Container | Port | Health Check | Alert Threshold |
|---------|-----------|------|--------------|-----------------|
| **PostgreSQL + pgvector** | onemind-postgres | 5432 | `pg_isready` | Connections > 80%, slow queries > 1s |
| **Redis** | onemind-redis | 6379 | `redis-cli PING` | Memory > 80%, evictions > 100/hr |
| **NATS JetStream** | onemind-nats | 4222, 8222, 9222 | HTTP `/varz` | Consumer lag > 100ms, stream > 80% |
| **TimescaleDB** | onemind-timescaledb | 5433 | `pg_isready` | Chunk compression lag, disk > 80% |
| **Agno API** | onemind-agno | 8081 | `/health` | Response time > 2s, error rate > 1% |
| **Flutter Web** | onemind-flutter-web | 3000 | HTTP 200 | Load time > 3s |
| **LiveKit** | onemind-livekit | 7880-7881 | HTTP health | Packet loss > 5%, jitter > 50ms |
| **Netdata** | onemind-netdata | 19999 | HTTP 200 | N/A (monitoring tool) |

### System Services (Azure VM)

| Service | Type | Check Command | Restart Command |
|---------|------|---------------|-----------------|
| **Cloudflare Tunnel** | systemd | `systemctl status cloudflared` | `sudo systemctl restart cloudflared` |
| **Tailscale** | systemd | `tailscale status` | `sudo systemctl restart tailscaled` |
| **Docker** | systemd | `systemctl status docker` | `sudo systemctl restart docker` |
| **Watchtower** | container | `docker ps \| grep watchtower` | `docker restart watchtower` |

### Key Metrics to Track

```yaml
System:
  - CPU utilization (alert > 80% sustained)
  - Memory utilization (alert > 85%)
  - Disk usage (alert > 90%)
  - Network I/O (anomaly detection)
  - Container restart count (alert > 5 in 5 min)

PostgreSQL:
  - Active connections (alert > 80 of max)
  - Slow queries (alert > 1s duration)
  - WAL size growth
  - Replication lag
  - Dead tuples (vacuum needed)

Redis:
  - Memory usage vs max
  - Eviction rate
  - Connected clients
  - Key hit/miss ratio

NATS:
  - Message throughput
  - Consumer pending count
  - Stream size vs quota
  - Subscription count
  - Slow consumer warnings

LiveKit:
  - Active rooms/participants
  - Bitrate (audio/video)
  - Packet loss percentage
  - Jitter measurements
  - Connection duration
```

---

## 2. Backend API Modules

### Core API Routers

| Module | Path Prefix | Key Endpoints | Owner Service | Failure Modes |
|--------|-------------|---------------|---------------|---------------|
| **Chat/Routing** | `/chat/*` | `POST /send`, `GET /commands`, `GET /specialists` | `backend/agno/routing/` | LLM timeout, agent unavailable, memory access |
| **Teams** | `/teams/*` | `POST /create`, `POST /run`, `GET /runs/*` | `backend/agno/teams/` | Coordination failures, deadlocks, timeouts |
| **Agents** | `/agents/*` | `POST /run`, `GET /list`, `GET /status/*` | `backend/agno/agents/` | Model failures, tool errors, memory issues |
| **Voice/Multimodal** | `/voice/*` | `WebSocket`, `POST /stream`, `GET /status` | `backend/agno/voice/` | LiveKit drops, STT/TTS failures |
| **Workflows** | `/workflows/*` | `POST /create`, `POST /run`, `GET /runs/*` | `backend/agno/workflows/` | Step timeouts, condition errors |
| **Knowledge/RAG** | `/knowledge/*` | `POST /index`, `GET /search`, `POST /embeddings` | `backend/agno/knowledge/` | Embedding failures, vector DB issues |
| **Memory** | `/memory/*` | `GET /blocks`, `PUT /blocks`, `POST /archival` | `backend/agno/memory/` | Database issues, search timeouts |
| **Learning** | `/learning/*` | `POST /learn`, `GET /entity_explorer` | `backend/agno/learning/` | Pattern extraction failures |
| **Skills** | `/skills/*` | `POST /execute`, `GET /list` | `backend/agno/skills/` | Execution failures, parameter errors |
| **Sessions** | `/sessions/*` | `POST /create`, `GET /history`, `PUT /state` | `backend/agno/sessions/` | State corruption, concurrent access |

### LifeOS API Routers

| Module | Path Prefix | Key Endpoints | Owner Service | Failure Modes |
|--------|-------------|---------------|---------------|---------------|
| **Unified Inbox** | `/api/inbox/*` | `GET /list`, `POST /create`, `WebSocket` | `backend/lifeos/inbox/` | Sync failures, filter errors |
| **Execution Engine** | `/api/engine/*` | `POST /dispatch`, `GET /status`, `POST /classify` | `backend/consciousness/engine/` | Classification errors, timeouts |
| **Consciousness** | `/awareness/*` | `GET /modes`, `POST /mode`, `GET /digest` | `backend/consciousness/` | Filter failures, digest timeouts |
| **Tasks/Projects** | `/api/tasks/*` | `POST /create`, `PUT /update`, `POST /dispatch` | `backend/lifeos/tasks/` | Task loss, dispatch failures |
| **Scheduling** | `/api/scheduling/*` | `POST /schedule`, `GET /conflicts` | `backend/lifeos/scheduling/` | Conflict detection failures |
| **Calendar** | `/api/calendar/*` | `GET /events`, `POST /sync` | `backend/lifeos/calendar/` | OAuth expiration, sync failures |
| **Habits** | `/api/habits/*` | `GET /list`, `POST /complete`, `GET /streaks` | `backend/lifeos/habits/` | Streak calculation errors |
| **Goals** | `/api/goals/*` | `GET /list`, `POST /update`, `GET /progress` | `backend/lifeos/goals/` | Progress calculation errors |
| **Meals** | `/api/meals/*` | `GET /recipes`, `POST /plan`, `GET /shopping` | `backend/lifeos/meals/` | Recipe sync failures |
| **Pomodoro** | `/api/pomodoro/*` | `POST /start`, `POST /stop`, `GET /stats` | `backend/lifeos/pomodoro/` | Timer precision issues |
| **Notifications** | `/api/notifications/*` | `GET /list`, `POST /send`, `PUT /read` | `backend/lifeos/notifications/` | Channel failures |
| **Journal** | `/api/journal/*` | `POST /entry`, `GET /list`, `GET /search` | `backend/lifeos/journal/` | Search failures |
| **Insights** | `/api/insights/*` | `GET /dashboard`, `GET /trends` | `backend/lifeos/insights/` | Calculation failures |
| **Bridge** | `/api/bridge/*` | `GET /intelligence`, `GET /predictions` | `backend/lifeos/bridge/` | Prediction errors |

### Integration Routers

| Module | Path Prefix | Key Endpoints | External Service | Failure Modes |
|--------|-------------|---------------|------------------|---------------|
| **Webhooks** | `/api/webhooks/*` | `POST /todoist`, `POST /github`, `POST /linear` | Multiple | Signature validation, parsing |
| **HomeAssistant** | `/integrations/homeassistant/*` | `GET /entities`, `POST /trigger` | Home Assistant | Connection drops, stale data |
| **MCP** | `/integrations/mcp/*` | `GET /resources`, `POST /execute` | MCP Servers | Server crashes, timeouts |
| **Databases** | `/integrations/databases/*` | `POST /query`, `GET /schema` | PostgreSQL, etc. | Query timeouts |

---

## 3. NATS Event Processing

### Event Streams & Subjects

| Stream | Subjects | Retention | Max Age | Purpose |
|--------|----------|-----------|---------|---------|
| **EVENTS** | `events.*` | Limits | 7 days | All system events |
| **INBOX** | `inbox.*` | WorkQueue | 30 days | Inbox item lifecycle |
| **TASKS** | `tasks.*` | WorkQueue | 30 days | Task state changes |
| **AGENTS** | `agents.*` | Limits | 24 hours | Agent run events |
| **WORKFLOWS** | `workflows.*` | Limits | 7 days | Workflow execution |
| **CONSCIOUSNESS** | `consciousness.*` | Limits | 24 hours | Awareness updates |

### Event Handlers

| Handler | Location | Subscribes To | Publishes To | Failure Modes |
|---------|----------|---------------|--------------|---------------|
| **Event Bus** | `backend/consciousness/core/bus.py` | All subjects | N/A | Connection drops |
| **Event Intelligence** | `backend/consciousness/processing/event_intelligence.py` | `events.*` | `intelligence.*` | Algorithm errors |
| **Archiver** | `backend/consciousness/processing/archiver.py` | `events.*` | N/A | TimescaleDB write failures |
| **Dead Letter Queue** | `backend/consciousness/processing/dlq.py` | Failed messages | `dlq.*` | DLQ overflow |
| **Correlator** | `backend/consciousness/processing/correlator.py` | Multiple | `correlations.*` | Matching errors |
| **Replay Engine** | `backend/consciousness/processing/replay.py` | `replay.*` | Original subjects | State inconsistencies |
| **Todoist Bridge** | `backend/consciousness/bridges/` | Webhooks | `todoist.*` | Task sync failures |
| **GitHub Bridge** | `backend/consciousness/bridges/` | Webhooks | `github.*` | Event parsing |
| **HomeAssistant Bridge** | `backend/consciousness/bridges/` | MQTT (1883) | `homeassistant.*` | Topic mapping |
| **Frigate Bridge** | `backend/consciousness/bridges/` | `frigate.events.*` | `security.*` | Frame processing |

### Key NATS Metrics

```yaml
Connection:
  - Connection state (connected/reconnecting/closed)
  - Reconnection count
  - Last error message

Streams:
  - Message count per stream
  - Byte size per stream
  - Consumer count
  - First/last sequence numbers

Consumers:
  - Pending message count (lag)
  - Ack pending count
  - Redelivery count
  - Last delivered sequence

Publishing:
  - Messages published per second
  - Publish errors
  - Publish latency (p50, p95, p99)

Subscribing:
  - Messages received per second
  - Processing latency
  - Error rate
```

---

## 4. Background Workers

| Worker | Location | Schedule | Dependencies | Failure Modes | Recovery |
|--------|----------|----------|--------------|---------------|----------|
| **Voice Worker** | `backend/agno/voice/worker.py` | Continuous | LiveKit, STT, TTS, LLM | Connection drops, inference timeout | Auto-reconnect |
| **Task Dispatch** | `backend/consciousness/engine/engine.py` | Event-driven | NATS, Agents, Teams | Dispatch failures, resource exhaustion | Retry with backoff |
| **Scheduled Triggers** | `backend/agno/workflows/workflow/triggers/` | Cron-based | NATS, Workflow Engine | Missed runs, clock skew | Catch-up on restart |
| **Auto-Reschedule** | `backend/lifeos/scheduling/auto_reschedule.py` | Continuous | PostgreSQL, Calendar | Infinite loops, conflicts | Circuit breaker |
| **Inbox Processing** | `backend/lifeos/inbox/service.py` | Event-driven | NATS, Classification | Classification failures | DLQ routing |
| **Sync Bridge** | `backend/integrations/sync/nats_bridge.py` | Continuous | NATS, External APIs | State mismatches, lag | Full resync |
| **Consciousness Gateway** | `backend/consciousness/processing/gateway.py` | Continuous | NATS, TimescaleDB | Pipeline blockages | Backpressure handling |

---

## 5. External API Integrations

### AI Model Providers

| Provider | Endpoint | Usage | Rate Limits | Cost Tracking | Fallback |
|----------|----------|-------|-------------|---------------|----------|
| **AWS Bedrock** | Claude 3.5 Sonnet | Primary LLM | Model quota | Per-token | OpenAI |
| **OpenAI** | GPT-4o, o1, embeddings | Secondary LLM | TPM/RPM | Per-token | Anthropic |
| **Anthropic** | Claude API | Backup LLM | TPM/RPM | Per-token | Groq |
| **Google** | Gemini 2.0 Live | Voice + Search | QPM | Per-request | OpenAI Realtime |
| **Groq** | Fast inference | Speed-critical | RPM | Per-token | Together |
| **Together AI** | Inference API | Alternative | RPM | Per-token | Deepseek |
| **Deepseek** | Deepseek API | Alternative | RPM | Per-token | N/A |

### Communication Services

| Service | Type | Purpose | Rate Limits | Failure Handling |
|---------|------|---------|-------------|------------------|
| **Discord** | Webhook/API | Notifications | 30 req/min per webhook | Queue + retry |
| **Telegram** | Bot API | Notifications | 30 msg/sec | Queue + retry |
| **Email (SMTP)** | SMTP | Notifications | Provider-dependent | Queue + retry |

### Productivity Integrations

| Service | Auth Type | Sync Type | Rate Limits | Refresh Interval |
|---------|-----------|-----------|-------------|------------------|
| **Todoist** | OAuth 2.0 | Bidirectional | 450 req/hr | Real-time (webhooks) |
| **GitHub** | OAuth 2.0 | Bidirectional | 5000 req/hr | Real-time (webhooks) |
| **Linear** | OAuth 2.0 | Bidirectional | API limits | Real-time (webhooks) |
| **Google Calendar** | OAuth 2.0 | Bidirectional | 1M req/day | 5 min polling + webhooks |
| **Outlook Calendar** | OAuth 2.0 | Bidirectional | API limits | 5 min polling |

### IoT/Smart Home

| Service | Protocol | Connection Type | Failure Handling |
|---------|----------|-----------------|------------------|
| **HomeAssistant** | REST + MQTT | Persistent | Reconnect + state sync |
| **Frigate** | REST + Event Stream | Persistent | Reconnect + replay |
| **MQTT Devices** | MQTT (via NATS) | Persistent | Message buffering |
| **BLE Devices** | Bluetooth LE | On-demand | Retry with backoff |

### Voice/Media Services

| Service | Purpose | Connection Type | Failure Handling |
|---------|---------|-----------------|------------------|
| **LiveKit** | RTC Server | WebSocket + UDP | Reconnect + resume |
| **Deepgram** | STT (via LiveKit) | Streaming | Fallback to OpenAI |
| **Cartesia** | TTS (via LiveKit) | Streaming | Fallback to OpenAI |
| **OpenAI Realtime** | Voice (alternative) | WebSocket | Primary fallback |

---

## 6. Frontend Connections

### HTTP Clients

| Client | Location | Base URL | Timeout | Retry Policy |
|--------|----------|----------|---------|--------------|
| **AgnoClient** | `frontend/lib/platform/api/agno_client.dart` | `/api/*` | 30s | 3 retries with backoff |
| **ConsciousnessClient** | `frontend/lib/platform/api/consciousness/` | `/awareness/*` | 15s | 2 retries |

### Real-time Connections

| Connection | Location | Protocol | Reconnect Policy | Failure UI |
|------------|----------|----------|------------------|------------|
| **NATS WebSocket** | `frontend/lib/platform/services/nats_service.dart` | WS (9222) | Auto-reconnect | Connection banner |
| **LiveKit** | `frontend/lib/specialized/voice/` | WS (7880) | Manual reconnect | Voice unavailable toast |
| **Inbox WebSocket** | `frontend/lib/lifeos/inbox/` | WS | Auto-reconnect | Sync indicator |

### Local Services

| Service | Location | Purpose | Failure Handling |
|---------|----------|---------|------------------|
| **Offline Cache** | `frontend/lib/platform/services/offline_cache_service.dart` | SQLite persistence | Queue mutations |
| **Mutation Queue** | `frontend/lib/platform/services/offline_mutation_queue.dart` | Offline changes | Sync on reconnect |
| **BLE Service** | `frontend/lib/platform/services/ble_service.dart` | Wearable comms | Retry pairing |
| **Voice Service** | `frontend/lib/platform/services/voice_service.dart` | Audio I/O | Device fallback |
| **MQTT Service** | `frontend/lib/platform/services/mqtt_service.dart` | IoT comms | Reconnect |
| **Sync Service** | `frontend/lib/platform/services/sync_service.dart` | Cross-device | Conflict resolution |

---

## 7. Business Logic Components

### Agent System

| Agent | Location | Specialization | Tools | Failure Modes |
|-------|----------|----------------|-------|---------------|
| **Legacy** | `backend/agno/agents/factory.py` | General assistant | All | Model failures, memory issues |
| **Researcher** | `backend/agno/agents/agents.py` | Web research | Search, RAG | Search failures, citation errors |
| **Analyst** | `backend/agno/agents/agents.py` | Data analysis | SQL, Charts | Computation failures |
| **Writer** | `backend/agno/agents/agents.py` | Content generation | Writing tools | Formatting errors |
| **Coder** | `backend/agno/agents/agents.py` | Code gen/review | Dev tools | Syntax errors, security |
| **Finance** | `backend/agno/agents/agents.py` | Financial analysis | Finance tools | Calculation errors |

### Team Orchestration

| Team | Location | Agents | Coordination | Failure Modes |
|------|----------|--------|--------------|---------------|
| **Research** | `backend/agno/teams/research_team.py` | Researcher, Analyst | Sequential | Coordination failures |
| **Dev** | `backend/agno/teams/dev_team.py` | Coder, Analyst | Parallel | Merge conflicts |
| **Finance** | `backend/agno/teams/csuite_team.py` | Finance, Analyst | Sequential | Data inconsistencies |
| **Content** | `backend/agno/teams/content_team.py` | Writer, Researcher | Sequential | Quality issues |
| **DevOps** | `backend/agno/teams/ops_team.py` | Coder, Analyst | Parallel | Deployment failures |
| **Guardian** | `backend/agno/teams/guardian_team.py` | Security agents | Parallel | False positives |
| **House** | `backend/agno/teams/house_team.py` | Home automation | Event-driven | Device failures |
| **Greenhouse** | `backend/agno/teams/greenhouse_team.py` | Environmental | Event-driven | Sensor failures |
| **Homestead** | `backend/agno/teams/homestead_team.py` | Farm automation | Event-driven | Actuator failures |

### Workflow Engine

| Component | Location | Purpose | Failure Modes |
|-----------|----------|---------|---------------|
| **Executor** | `backend/agno/workflows/workflow/executor.py` | Step orchestration | Timeouts, resource exhaustion |
| **Agent Step** | `backend/agno/workflows/workflow/steps/agent.py` | Agent invocation | Agent failures |
| **Condition Step** | `backend/agno/workflows/workflow/steps/condition.py` | Branching logic | Evaluation errors |
| **Loop Step** | `backend/agno/workflows/workflow/steps/loop.py` | Iteration | Infinite loops |
| **Router Step** | `backend/agno/workflows/workflow/steps/router.py` | Dynamic routing | Routing errors |
| **Schedule Trigger** | `backend/agno/workflows/workflow/triggers/schedule.py` | Time-based start | Clock skew, missed |

### Consciousness System

| Component | Location | Purpose | Failure Modes |
|-----------|----------|---------|---------------|
| **Awareness Engine** | `backend/consciousness/awareness.py` | Context aggregation | Stale context |
| **Event Bus** | `backend/consciousness/core/bus.py` | NATS pub/sub | Connection drops |
| **Event Intelligence** | `backend/consciousness/processing/event_intelligence.py` | Pattern detection | False positives/negatives |
| **Event Processor** | `backend/consciousness/processing/processor.py` | Pipeline orchestration | Blockages, latency |
| **Archiver** | `backend/consciousness/processing/archiver.py` | TimescaleDB persistence | Write failures |
| **Replay** | `backend/consciousness/processing/replay.py` | Historical replay | State mismatches |
| **DLQ** | `backend/consciousness/processing/dlq.py` | Failed messages | Overflow |
| **Correlator** | `backend/consciousness/processing/correlator.py` | Event correlation | Matching errors |
| **KV Store** | `backend/consciousness/core/kv_store.py` | NATS KV | Store failures |
| **Object Store** | `backend/consciousness/core/object_store.py` | NATS Objects | Corruption |

### LifeOS Modules

| Module | Location | Purpose | Key Metrics |
|--------|----------|---------|-------------|
| **Inbox** | `backend/lifeos/inbox/` | Task aggregation | Items/day, resolution time |
| **Tasks** | `backend/lifeos/tasks/` | Task management | Completion rate |
| **Scheduling** | `backend/lifeos/scheduling/` | Time management | Conflict rate |
| **Calendar** | `backend/lifeos/calendar/` | Calendar sync | Sync latency |
| **Habits** | `backend/lifeos/habits/` | Habit tracking | Streak accuracy |
| **Goals** | `backend/lifeos/goals/` | Goal progress | Progress accuracy |
| **Meals** | `backend/lifeos/meals/` | Meal planning | Recipe sync |
| **Pomodoro** | `backend/lifeos/pomodoro/` | Focus sessions | Session completion |
| **Notifications** | `backend/lifeos/notifications/` | Alert routing | Delivery rate |
| **Journal** | `backend/lifeos/journal/` | Journaling | Entry count |
| **Insights** | `backend/lifeos/insights/` | Analytics | Dashboard load time |
| **Bridge** | `backend/lifeos/bridge/` | Intelligence | Prediction accuracy |

---

## 8. Critical User Paths

### Path 1: Message → Response

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ Frontend WebSocket/HTTP                                                      │
│     ↓                                                                        │
│ Agno /chat/send endpoint                                                     │
│     ↓                                                                        │
│ Agent Factory + Session Load                                                 │
│     ↓                                                                        │
│ LLM Inference (AWS Bedrock / OpenAI)                                        │
│     ↓                                                                        │
│ Tool Execution (Todoist, GitHub, HomeAssistant, etc.)                       │
│     ↓                                                                        │
│ Memory Storage (PostgreSQL + Redis)                                          │
│     ↓                                                                        │
│ Response Stream → Frontend                                                   │
└─────────────────────────────────────────────────────────────────────────────┘

Failure Points: Network, LLM unavailable, tool API errors, database, memory
Target Latency: < 5s to first token, < 30s total
```

### Path 2: Webhook Event → Inbox → Execution

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ External Service (Todoist/GitHub/Linear)                                     │
│     ↓                                                                        │
│ Webhook Receiver (/api/webhooks/*)                                          │
│     ↓                                                                        │
│ NATS Publication                                                             │
│     ↓                                                                        │
│ Event Processor (consciousness/processing/gateway.py)                        │
│     ↓                                                                        │
│ Inbox Item Creation (lifeos/inbox/service.py)                               │
│     ↓                                                                        │
│ Execution Engine Classification                                              │
│     ↓                                                                        │
│ Agent/Team/Workflow Dispatch                                                 │
│     ↓                                                                        │
│ Result Callback (Discord/Telegram notification)                              │
└─────────────────────────────────────────────────────────────────────────────┘

Failure Points: Webhook validation, NATS, parsing, classification, timeouts
Target Latency: < 60s end-to-end
```

### Path 3: Real-time Voice Call

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ Frontend WebSocket → LiveKit                                                 │
│     ↓                                                                        │
│ Voice Worker (backend/agno/voice/worker.py)                                 │
│     ↓                                                                        │
│ STT (Deepgram/Cartesia via LiveKit)                                         │
│     ↓                                                                        │
│ Agent Processing + Tool Execution                                            │
│     ↓                                                                        │
│ TTS (Cartesia/Gemini)                                                        │
│     ↓                                                                        │
│ Audio Stream → Client                                                        │
└─────────────────────────────────────────────────────────────────────────────┘

Failure Points: WebSocket, audio codec, STT, model inference, TTS, latency
Target Latency: < 500ms voice-to-voice
```

### Path 4: Consciousness/Awareness Processing

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ NATS Event Stream                                                            │
│     ↓                                                                        │
│ Processor (correlate events)                                                 │
│     ↓                                                                        │
│ Intelligence Engine (anomaly detection)                                      │
│     ↓                                                                        │
│ Archiver (TimescaleDB)                                                       │
│     ↓                                                                        │
│ Context Aggregators                                                          │
│     ↓                                                                        │
│ Awareness Mode Update                                                        │
│     ↓                                                                        │
│ Frontend Digest/Alerts                                                       │
└─────────────────────────────────────────────────────────────────────────────┘

Failure Points: Pipeline blockage, algorithm, TimescaleDB, context staleness
Target Latency: < 5s event-to-awareness
```

### Path 5: Workflow Execution

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ Trigger (schedule/event/API)                                                 │
│     ↓                                                                        │
│ Workflow Engine                                                              │
│     ↓                                                                        │
│ Step 1: Agent Execution                                                      │
│     ↓                                                                        │
│ Step 2: Conditional Logic                                                    │
│     ↓                                                                        │
│ Step 3: Loop Iteration                                                       │
│     ↓                                                                        │
│ Step 4: Result Handling                                                      │
│     ↓                                                                        │
│ Status Update (Redis/PostgreSQL)                                             │
└─────────────────────────────────────────────────────────────────────────────┘

Failure Points: Step timeouts, condition errors, agent failures, resource exhaustion
Target: 95% success rate, < 5 min per workflow
```

---

## 9. Current Monitoring Coverage

### ✅ What We Have

| Tool | Coverage | Location |
|------|----------|----------|
| **Netdata** | System metrics (CPU, memory, disk, network) | https://netdata.onemindos.dev |
| **NATS Monitor** | Message queue metrics | https://nats.onemindos.dev |
| **Portainer** | Container management | https://portainer.onemindos.dev |
| **PostgreSQL** | Built-in `pg_stat_*` views | Via psql |
| **Redis** | `INFO` command | Via redis-cli |
| **LiveKit** | Built-in dashboard | Port 7880 |
| **Docker** | Container logs | `docker logs` |

### ❌ Current Gaps

| Gap | Impact | Priority |
|-----|--------|----------|
| **Application Traces** | Can't trace requests end-to-end | HIGH |
| **Centralized Logs** | Logs scattered across containers | HIGH |
| **Error Rate Alerting** | No proactive error notifications | HIGH |
| **Custom Business Metrics** | No task completion, agent success tracking | MEDIUM |
| **Model Inference Tracking** | No token usage, latency by model | MEDIUM |
| **API Latency Percentiles** | No p50/p95/p99 tracking | MEDIUM |
| **Webhook Delivery Tracking** | No delivery success metrics | LOW |
| **Cost Attribution** | No per-feature cost breakdown | LOW |

---

## 10. Chosen Observability Platform: Datadog

> **Decision Date:** 2026-01-30
> **Implementation Issue:** [OMOS-486](https://linear.app/onemind-os/issue/OMOS-486/implement-datadog-observability-platform)
> **Full Plan:** [DATADOG-OBSERVABILITY-PLAN.md](./DATADOG-OBSERVABILITY-PLAN.md)

### Why Datadog Won

| Platform | Plug & Play | Single Pane | AI Debug | NATS | LLM Obs | Score |
|----------|-------------|-------------|----------|------|---------|-------|
| **Datadog** | 10/10 | 10/10 | 10/10 | ✅ Native | ✅ Built-in | **78/80** |
| New Relic | 9/10 | 9/10 | 9/10 | ⚠️ Custom | ❌ | 66/80 |
| Elastic | 7/10 | 8/10 | 7/10 | ⚠️ Custom | ❌ | 59/80 |
| Azure Monitor | 7/10 | 6/10 | 5/10 | ❌ | ❌ | 49/80 |

### What Datadog Covers

| Category | Components | Datadog Feature |
|----------|------------|-----------------|
| **Infrastructure** | Azure VM, Docker, Cloudflare | Infrastructure Monitoring |
| **Databases** | PostgreSQL, Redis, TimescaleDB | Database Monitoring + Query Analysis |
| **Message Queue** | NATS JetStream, MQTT | Native NATS Integration |
| **Application** | FastAPI (150+ endpoints) | APM + Distributed Tracing |
| **AI/LLM** | Bedrock, OpenAI, Anthropic | LLM Observability (tokens, costs) |
| **Logs** | All containers | Log Management + Correlation |
| **Alerting** | All components | Monitors + Watchdog AI |
| **Debugging** | End-to-end traces | Bits AI (natural language) |

### Expected Monthly Cost

| Component | Cost |
|-----------|------|
| Infrastructure (1 host) | $15 |
| APM (~10M spans) | $17 |
| Logs (~10 GB) | $1.50 |
| Custom Metrics | Included |
| LLM Observability | Included |
| **Total** | **~$35-50/month** |

---

## 11. Datadog Implementation Checklist

> Full details in [DATADOG-OBSERVABILITY-PLAN.md](./DATADOG-OBSERVABILITY-PLAN.md)

### Phase 1: Agent Installation
- [ ] Create Datadog account
- [ ] Install Datadog Agent on Azure VM
- [ ] Configure Docker integration
- [ ] Verify host and container metrics

### Phase 2: Database Integrations
- [ ] Configure PostgreSQL integration
- [ ] Configure Redis integration
- [ ] Configure TimescaleDB
- [ ] Enable query analysis

### Phase 3: NATS Integration
- [ ] Configure NATS integration
- [ ] Verify stream and consumer metrics
- [ ] Create NATS dashboard

### Phase 4: APM Setup
- [ ] Install ddtrace Python package
- [ ] Configure FastAPI auto-instrumentation
- [ ] Add NATS trace propagation to BaseService
- [ ] Verify distributed traces

### Phase 5: LLM Observability
- [ ] Enable LLM Observability feature
- [ ] Instrument AWS Bedrock calls
- [ ] Instrument OpenAI calls
- [ ] Configure token/cost tracking

### Phase 6: Log Management
- [ ] Configure Docker log collection
- [ ] Set up log parsing pipelines
- [ ] Configure log-trace correlation

### Phase 7: Custom Metrics
- [ ] Add agent run metrics
- [ ] Add workflow metrics
- [ ] Add inbox metrics
- [ ] Add integration metrics

### Phase 8: Alerting
- [ ] Configure Discord webhook
- [ ] Configure Telegram integration
- [ ] Create alert monitors by priority
- [ ] Set up escalation policies

### Phase 9: Cleanup
- [ ] Remove Netdata container
- [ ] Remove Netdata from docker-compose
- [ ] Update Cloudflare tunnel routes
- [ ] Update documentation

---

## 12. Runbook Links

| Scenario | Runbook |
|----------|---------|
| API endpoint slow | TBD |
| Agent failures spike | TBD |
| NATS consumer lag | TBD |
| Database connection exhaustion | TBD |
| External API rate limiting | TBD |
| Voice call quality issues | TBD |

---

## Appendix: File Locations

```
Infrastructure:
├── docker-compose.core.yml          # Production services
├── docker-compose.dev.yml           # Development services
├── infra/config/nats/hub.conf       # NATS configuration
├── infra/config/livekit/livekit.yaml # LiveKit configuration

Backend:
├── backend/agno/                    # Agent system
│   ├── agents/                      # Agent definitions
│   ├── teams/                       # Team orchestration
│   ├── workflows/                   # Workflow engine
│   ├── voice/                       # Voice/multimodal
│   ├── routing/                     # Chat routing
│   ├── knowledge/                   # RAG system
│   └── tools/                       # Tool registry
├── backend/consciousness/           # Awareness system
│   ├── core/                        # Event bus, KV store
│   ├── processing/                  # Event pipeline
│   ├── context/                     # Context aggregators
│   └── engine/                      # Execution engine
├── backend/lifeos/                  # LifeOS modules
│   ├── inbox/                       # Unified inbox
│   ├── tasks/                       # Task management
│   ├── scheduling/                  # Scheduling
│   └── ...                          # Other modules
└── backend/integrations/            # External integrations

Frontend:
├── frontend/lib/platform/api/       # API clients
├── frontend/lib/platform/services/  # Platform services
└── frontend/lib/platform/providers/ # State providers

Documentation:
├── docs/guides/operations/ACCESS-REFERENCE.md
└── docs/guides/operations/MONITORING-INVENTORY.md  # This file
```

---

*This document should be updated whenever new services, integrations, or monitoring capabilities are added.*

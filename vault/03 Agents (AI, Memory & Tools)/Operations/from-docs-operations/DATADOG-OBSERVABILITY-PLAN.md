# OneMind OS - Datadog Observability Plan

> Last Updated: 2026-01-30
> Status: PLANNED
> Priority: HIGH
> Linear Issue: [OMOS-486](https://linear.app/onemind-os/issue/OMOS-486/implement-datadog-observability-platform)
> Estimated Setup Time: 2-3 hours

---

## Executive Summary

Datadog will serve as the **single pane of glass** for all OneMind OS observability, replacing Netdata and providing comprehensive monitoring across infrastructure, applications, databases, message queues, and AI/LLM operations.

### Why Datadog?

| Requirement | Datadog Capability |
|-------------|-------------------|
| Plug-and-play setup | ✅ One agent, auto-discovery |
| See everything in one place | ✅ Unified UI for logs, traces, metrics, APM |
| Debug AI agent flows | ✅ LLM Observability + distributed tracing |
| Native NATS support | ✅ Built-in integration (unique among competitors) |
| AI-assisted debugging | ✅ Bits AI natural language queries |
| Anomaly detection | ✅ Watchdog auto-correlation |

### What Datadog Replaces

| Current Tool | Status After Datadog |
|--------------|---------------------|
| Netdata | **REMOVE** - Fully replaced |
| Portainer | **KEEP** - Management tool, not monitoring |
| Manual log checking | **REMOVE** - Centralized in Datadog |
| SSH for metrics | **REMOVE** - All in Datadog UI |

---

## Complete Coverage Map

### 1. Infrastructure Monitoring

#### Azure VM (onemind-vm)

| Metric | Collection Method | Dashboard |
|--------|------------------|-----------|
| CPU utilization | Datadog Agent | Infrastructure > Hosts |
| Memory usage | Datadog Agent | Infrastructure > Hosts |
| Disk I/O & space | Datadog Agent | Infrastructure > Hosts |
| Network throughput | Datadog Agent | Infrastructure > Hosts |
| System load | Datadog Agent | Infrastructure > Hosts |
| Process list | Datadog Agent | Infrastructure > Processes |
| File descriptors | Datadog Agent | Infrastructure > Hosts |
| TCP connections | Datadog Agent | Infrastructure > Hosts |

**Alerts to Configure:**
- CPU > 80% for 5 minutes
- Memory > 85% for 5 minutes
- Disk > 90%
- High system load

#### Docker Containers (8 services)

| Container | Metrics Collected |
|-----------|------------------|
| onemind-agno | CPU, memory, network, restarts, logs |
| onemind-flutter-web | CPU, memory, network, restarts, logs |
| onemind-postgres | CPU, memory, network, restarts, logs |
| onemind-redis | CPU, memory, network, restarts, logs |
| onemind-nats | CPU, memory, network, restarts, logs |
| onemind-timescaledb | CPU, memory, network, restarts, logs |
| onemind-livekit | CPU, memory, network, restarts, logs |
| portainer | CPU, memory, network, restarts, logs |

**Alerts to Configure:**
- Container restart > 3 in 5 minutes
- Container CPU > 90%
- Container memory > 90% of limit
- Container stopped unexpectedly

#### Cloudflare Tunnel

| Metric | Collection Method |
|--------|------------------|
| Tunnel status | Cloudflare integration |
| Request count | Cloudflare integration |
| Error rate | Cloudflare integration |
| Bandwidth | Cloudflare integration |

---

### 2. Database Monitoring

#### PostgreSQL + pgvector (Port 5432)

| Metric | Why It Matters |
|--------|----------------|
| Active connections | Prevent connection exhaustion |
| Connection pool usage | Optimize pool size |
| Queries per second | Traffic patterns |
| Slow queries (>1s) | Performance bottlenecks |
| Query execution plans | Optimization opportunities |
| Lock waits | Deadlock detection |
| Dead tuples | Vacuum needed |
| Table/index sizes | Storage planning |
| Replication lag | Data consistency |
| WAL size | Disk pressure |
| Cache hit ratio | Memory efficiency |
| Transactions/sec | Throughput |

**Deep Query Analysis:**
- Automatic EXPLAIN ANALYZE on slow queries
- Query fingerprinting
- Resource consumption per query
- Wait event analysis

**Alerts to Configure:**
- Connections > 80% of max
- Slow queries > 5/minute
- Replication lag > 10 seconds
- Cache hit ratio < 95%
- Dead tuples > 10,000

#### TimescaleDB (Port 5433)

| Metric | Why It Matters |
|--------|----------------|
| Hypertable sizes | Storage growth |
| Chunk compression status | Storage efficiency |
| Continuous aggregate refresh | Data freshness |
| All PostgreSQL metrics above | Same monitoring |

**Alerts to Configure:**
- Uncompressed chunks > 7 days old
- Hypertable growth > 1GB/day

#### Redis (Port 6379)

| Metric | Why It Matters |
|--------|----------------|
| Memory usage | Prevent OOM |
| Memory fragmentation | Memory efficiency |
| Connected clients | Connection patterns |
| Commands processed/sec | Throughput |
| Hit rate | Cache effectiveness |
| Evicted keys | Memory pressure |
| Blocked clients | Performance issues |
| Key count by database | Data distribution |
| Expiring keys | TTL patterns |
| Pub/sub channels | Messaging patterns |
| Slow log entries | Performance issues |
| Persistence status | Data safety |

**Alerts to Configure:**
- Memory > 80% of max
- Evictions > 100/hour
- Hit rate < 90%
- Blocked clients > 0

---

### 3. Message Queue Monitoring (NATS JetStream)

#### NATS Core (Port 4222)

| Metric | Why It Matters |
|--------|----------------|
| Connection count | Client health |
| Subscriptions | Routing patterns |
| Messages in/out | Throughput |
| Bytes in/out | Bandwidth |
| Slow consumers | Backpressure |
| Route count | Cluster health |

#### JetStream Streams

| Stream | Metrics |
|--------|---------|
| EVENTS | Messages, bytes, consumers, consumer lag |
| INBOX | Messages, bytes, consumers, consumer lag |
| TASKS | Messages, bytes, consumers, consumer lag |
| AGENTS | Messages, bytes, consumers, consumer lag |
| WORKFLOWS | Messages, bytes, consumers, consumer lag |
| CONSCIOUSNESS | Messages, bytes, consumers, consumer lag |

**Per-Stream Metrics:**
- Message count
- Byte size
- First/last sequence
- Consumer count
- Pending messages per consumer
- Ack pending
- Redelivery count

#### MQTT Bridge (Port 1883)

| Metric | Why It Matters |
|--------|----------------|
| Connected MQTT clients | IoT device health |
| Messages bridged | IoT traffic |
| Topic subscription count | Device patterns |

**Alerts to Configure:**
- Consumer lag > 1000 messages
- Slow consumer warnings
- Stream storage > 80% of limit
- Connection drops > 5/minute

---

### 4. Application Performance Monitoring (APM)

#### FastAPI Backend (150+ Endpoints)

**Automatic Instrumentation Covers:**

| Component | What's Traced |
|-----------|--------------|
| HTTP requests | All incoming requests with timing |
| Route handlers | Each endpoint function |
| Database queries | PostgreSQL, Redis calls |
| External HTTP | Outbound API calls |
| Background tasks | Async task execution |
| NATS operations | Publish/subscribe (with custom spans) |

**Metrics per Endpoint:**

| Metric | Description |
|--------|-------------|
| Request rate | Requests/second |
| Latency (p50, p95, p99) | Response time percentiles |
| Error rate | 4xx/5xx percentage |
| Throughput | Successful requests/second |
| Apdex score | User satisfaction metric |

#### Endpoint Groups to Monitor

**Chat/Routing API:**
```
POST /chat/send              - Primary chat endpoint
GET  /chat/commands          - Command listing
GET  /chat/specialists       - Specialist routing
GET  /chat/context           - Context aggregation
GET  /chat/memory/*          - Memory operations
```

**Teams API:**
```
POST /teams/create           - Team creation
POST /teams/run              - Team execution (critical path)
GET  /teams/runs/*           - Run status/history
PUT  /teams/update           - Team updates
```

**Agents API:**
```
POST /agents/*/run           - Agent execution (critical path)
GET  /agents/list            - Agent listing
GET  /agents/status/*        - Agent status
```

**Voice/Multimodal API:**
```
WebSocket /voice/*           - Real-time voice (critical path)
POST /voice/stream           - Voice streaming
GET  /voice/status           - Voice worker status
```

**Workflows API:**
```
POST /workflows/create       - Workflow creation
POST /workflows/run          - Workflow execution
GET  /workflows/runs/*       - Run status/history
```

**Knowledge/RAG API:**
```
POST /knowledge/index        - Document indexing
GET  /knowledge/search       - Semantic search
POST /knowledge/embeddings   - Embedding generation
```

**LifeOS APIs:**
```
/api/inbox/*                 - Unified inbox
/api/engine/*                - Execution engine
/awareness/*                 - Consciousness system
/api/tasks/*                 - Task management
/api/scheduling/*            - Scheduling
/api/calendar/*              - Calendar sync
/api/habits/*                - Habit tracking
/api/goals/*                 - Goal management
/api/notifications/*         - Notifications
```

**Integration APIs:**
```
/api/webhooks/todoist        - Todoist webhooks
/api/webhooks/github         - GitHub webhooks
/api/webhooks/linear         - Linear webhooks
/integrations/homeassistant/* - Home Assistant
/integrations/mcp/*          - MCP servers
```

**Alerts to Configure:**
- Error rate > 1% on any endpoint
- p95 latency > 5s on critical paths
- Throughput drop > 50%
- Apdex < 0.8

---

### 5. Distributed Tracing

#### Critical User Paths

**Path 1: Message → Response**
```
Trace spans:
├── [HTTP] POST /chat/send
│   ├── [Internal] session.load
│   ├── [Internal] agent.factory.create
│   ├── [External] AWS Bedrock / OpenAI API
│   │   └── [Metric] tokens_used, model, latency
│   ├── [Internal] tool.execute (if any)
│   │   └── [External] Todoist/GitHub/HomeAssistant API
│   ├── [Database] PostgreSQL - memory.store
│   ├── [Cache] Redis - session.update
│   └── [HTTP] Response stream
```

**Path 2: Webhook → Inbox → Execution**
```
Trace spans:
├── [HTTP] POST /api/webhooks/todoist
│   ├── [Internal] signature.validate
│   ├── [NATS] publish events.todoist.*
│   └── [HTTP] 200 OK
├── [NATS] subscribe events.*
│   ├── [Internal] processor.correlate
│   ├── [Database] TimescaleDB - archiver.store
│   └── [NATS] publish inbox.item.created
├── [NATS] subscribe inbox.*
│   ├── [Internal] engine.classify
│   ├── [Internal] engine.dispatch
│   │   └── [Internal] agent.run / team.run / workflow.run
│   └── [External] Discord/Telegram notification
```

**Path 3: Real-time Voice**
```
Trace spans:
├── [WebSocket] /voice/connect
│   ├── [External] LiveKit - room.join
│   ├── [External] Deepgram STT
│   ├── [Internal] agent.process
│   │   └── [External] LLM API
│   ├── [External] Cartesia TTS
│   └── [WebSocket] audio.stream
```

**Path 4: Workflow Execution**
```
Trace spans:
├── [Trigger] schedule / event / API
│   ├── [Internal] workflow.engine.start
│   ├── [Internal] step.1.agent
│   │   └── [External] LLM API
│   ├── [Internal] step.2.condition
│   ├── [Internal] step.3.loop
│   │   └── [Internal] iteration.1..n
│   ├── [Database] PostgreSQL - state.persist
│   └── [Cache] Redis - status.update
```

---

### 6. LLM Observability (Unique to Datadog)

#### Model Providers to Track

| Provider | Models | Metrics |
|----------|--------|---------|
| AWS Bedrock | Claude 3.5 Sonnet | Tokens, latency, cost |
| OpenAI | GPT-4o, o1, embeddings | Tokens, latency, cost |
| Anthropic | Claude (backup) | Tokens, latency, cost |
| Google | Gemini 2.0 Live | Tokens, latency, cost |
| Groq | Fast inference | Tokens, latency, cost |
| Together AI | Various | Tokens, latency, cost |
| Deepseek | Deepseek models | Tokens, latency, cost |

#### Metrics per LLM Call

| Metric | Description |
|--------|-------------|
| Input tokens | Prompt size |
| Output tokens | Response size |
| Total tokens | Combined usage |
| Latency | Time to first token, total time |
| Cost | Calculated from token counts |
| Model | Which model was used |
| Temperature | Generation parameters |
| Success/failure | Completion status |
| Error type | Timeout, rate limit, quota, etc. |

#### Agent-Level Tracking

| Metric | Description |
|--------|-------------|
| Agent name | Which agent made the call |
| Session ID | User session context |
| Tool calls | Tools invoked during response |
| Memory operations | Reads/writes |
| Total cost per request | Sum of all LLM calls |

**Dashboards:**
- Token usage by model provider (daily/weekly/monthly)
- Cost breakdown by agent
- Latency trends by model
- Error rate by provider
- Most expensive sessions

**Alerts to Configure:**
- Daily token usage > budget threshold
- Single request cost > $1
- Model error rate > 5%
- Latency p95 > 10s

---

### 7. Log Management

#### Log Sources

| Source | Collection Method | Volume Estimate |
|--------|------------------|-----------------|
| onemind-agno | Docker log driver | ~5 GB/month |
| onemind-flutter-web | Docker log driver | ~500 MB/month |
| onemind-postgres | Docker log driver | ~1 GB/month |
| onemind-redis | Docker log driver | ~200 MB/month |
| onemind-nats | Docker log driver | ~500 MB/month |
| onemind-timescaledb | Docker log driver | ~500 MB/month |
| onemind-livekit | Docker log driver | ~1 GB/month |
| System logs (syslog) | Datadog Agent | ~500 MB/month |
| Cloudflare access logs | Cloudflare integration | ~1 GB/month |

**Total Estimated: ~10-12 GB/month**

#### Log Processing

| Feature | Description |
|---------|-------------|
| Parsing | Automatic JSON/structured log parsing |
| Indexing | Full-text search on all logs |
| Correlation | Link logs to traces via trace_id |
| Patterns | Auto-detect log patterns and anomalies |
| Archives | Long-term storage to S3/Azure Blob |

#### Log-Based Alerts

| Alert | Condition |
|-------|-----------|
| Error spike | ERROR logs > 10/minute |
| Exception detected | Stack trace in logs |
| Security event | Auth failure patterns |
| Database error | PostgreSQL ERROR logs |
| NATS disconnection | Connection error patterns |

---

### 8. Synthetic Monitoring & Uptime

#### Endpoints to Monitor

| Endpoint | Check Interval | Timeout |
|----------|---------------|---------|
| https://api.onemindos.dev/health | 1 minute | 10s |
| https://command.onemindos.dev | 1 minute | 10s |
| https://nats.onemindos.dev | 5 minutes | 10s |
| https://netdata.onemindos.dev | 5 minutes | 10s |

#### API Tests

| Test | Description |
|------|-------------|
| Health check | GET /health returns 200 |
| Agent list | GET /agents returns valid JSON |
| Chat endpoint | POST /chat/send returns stream |
| Auth check | Protected endpoints reject unauthorized |

**Alerts:**
- Endpoint down > 2 consecutive checks
- SSL certificate expiring < 14 days
- Response time > 5s

---

### 9. Real User Monitoring (RUM)

#### Flutter Web Tracking

| Metric | Description |
|--------|-------------|
| Page load time | Time to interactive |
| First contentful paint | Initial render |
| Largest contentful paint | Main content visible |
| Cumulative layout shift | Visual stability |
| First input delay | Interactivity |
| JavaScript errors | Frontend exceptions |
| API call performance | Backend request timing |
| User sessions | Session duration, pages viewed |

#### User Flow Tracking

| Flow | Steps Tracked |
|------|---------------|
| Chat interaction | Open → Type → Send → Response received |
| Voice call | Connect → Speak → Response → End |
| Task creation | Open inbox → Create → Dispatch |
| Workflow execution | Trigger → Monitor → Complete |

---

### 10. Security Monitoring

#### Threat Detection

| Detection | Description |
|-----------|-------------|
| Brute force attempts | Multiple auth failures |
| Unusual API patterns | Anomalous request rates |
| SQL injection attempts | Malicious query patterns |
| XSS attempts | Script injection in inputs |
| Privilege escalation | Unauthorized access attempts |

#### Audit Logging

| Event | Logged Data |
|-------|-------------|
| Authentication | User, IP, success/failure, timestamp |
| Authorization | Resource, action, allowed/denied |
| Data access | User, resource, operation |
| Configuration changes | User, setting, old/new value |
| Agent executions | Agent, user, tools used |

---

### 11. Alerting & Incident Management

#### Alert Categories

**P1 - Critical (Immediate notification):**
- API completely down
- Database unreachable
- All agents failing
- Security incident

**P2 - High (15-minute notification):**
- Error rate > 5%
- p95 latency > 10s
- Consumer lag > 5000 messages
- Disk > 95%

**P3 - Medium (Hourly digest):**
- Error rate > 1%
- Memory > 85%
- Slow queries increasing
- Token usage > 80% of budget

**P4 - Low (Daily digest):**
- Performance degradation trends
- Cost optimization opportunities
- Capacity planning alerts

#### Notification Channels

| Channel | Use Case |
|---------|----------|
| **Linear** | P1 and P2 alerts (auto-create issues) |
| Discord webhook | All alerts |
| Telegram | P1 and P2 alerts |
| Email | Daily/weekly digests |
| Mobile push | P1 alerts only |

#### Linear Integration (Auto-Issue Creation)

Datadog will automatically create Linear issues when alerts trigger:

| Feature | Configuration |
|---------|---------------|
| **Team** | OneMind OS (OMOS) |
| **Priority Mapping** | P1 → Urgent, P2 → High, P3 → Medium |
| **Issue Content** | Alert details, dashboard links, trace links |
| **Two-way Sync** | Closing Linear issue can resolve Datadog alert |
| **Labels** | Auto-tag with `datadog`, alert type |

**Alert → Linear Routing:**

| Alert Type | Linear Label | Auto-Assign |
|------------|--------------|-------------|
| Infrastructure | `infra` | - |
| API Errors | `backend` | - |
| Database | `database` | - |
| Agent/LLM | `ai` | - |
| Security | `security` | - |

#### Escalation Policy

```
P1 Alert:
├── Immediately: Discord + Telegram + Push
├── After 5 min unacked: Repeat
└── After 15 min: Email summary

P2 Alert:
├── Immediately: Discord
├── After 15 min unacked: Telegram
└── After 1 hour: Email summary
```

---

### 12. Dashboards to Create

#### 1. Executive Overview
- System health score
- Active users (24h)
- Tasks completed
- Agent success rate
- Total cost (daily/monthly)

#### 2. Infrastructure Health
- VM CPU/Memory/Disk
- Container status grid
- Network throughput
- Cloudflare tunnel status

#### 3. API Performance
- Request rate by endpoint
- Latency heatmap
- Error rate trends
- Apdex scores

#### 4. Database Performance
- PostgreSQL connections and queries
- Redis memory and hit rate
- Slow query list
- Connection pool status

#### 5. NATS/Messaging
- Message throughput
- Consumer lag by stream
- Subscription count
- Dead letter queue size

#### 6. Agent Operations
- Agent runs per hour
- Success/failure rate by agent
- Average run duration
- Tool usage breakdown

#### 7. LLM Operations
- Token usage by provider
- Cost breakdown (daily/weekly)
- Latency by model
- Error rate by provider

#### 8. Voice/Real-time
- Active voice sessions
- STT/TTS latency
- WebSocket connection health
- Audio quality metrics

#### 9. LifeOS Modules
- Inbox items created/resolved
- Task completion rate
- Habit streaks
- Goal progress

#### 10. Security & Compliance
- Auth events
- Anomaly detections
- API abuse attempts
- Audit log activity

---

### 13. Custom Metrics to Implement

#### Agent Metrics

```python
from datadog import statsd

# Track agent run
statsd.increment('onemind.agent.run', tags=['agent:legacy', 'status:success'])
statsd.histogram('onemind.agent.duration', duration_ms, tags=['agent:legacy'])
statsd.histogram('onemind.agent.tokens', total_tokens, tags=['agent:legacy', 'model:claude-3.5'])
```

#### Workflow Metrics

```python
# Track workflow execution
statsd.increment('onemind.workflow.run', tags=['workflow:daily_briefing', 'status:completed'])
statsd.histogram('onemind.workflow.duration', duration_ms, tags=['workflow:daily_briefing'])
statsd.increment('onemind.workflow.step', tags=['workflow:daily_briefing', 'step:agent', 'status:success'])
```

#### Inbox Metrics

```python
# Track inbox operations
statsd.increment('onemind.inbox.item.created', tags=['source:todoist', 'type:task'])
statsd.increment('onemind.inbox.item.dispatched', tags=['target:agent', 'agent:legacy'])
statsd.histogram('onemind.inbox.resolution_time', resolution_ms, tags=['source:todoist'])
```

#### Integration Metrics

```python
# Track external API calls
statsd.increment('onemind.integration.call', tags=['service:todoist', 'operation:sync', 'status:success'])
statsd.histogram('onemind.integration.latency', latency_ms, tags=['service:todoist'])
```

---

### 14. Implementation Checklist

#### Phase 1: Agent Installation (Day 1)

- [ ] Create Datadog account
- [ ] Install Datadog Agent on Azure VM
- [ ] Configure Docker integration
- [ ] Verify host metrics appearing
- [ ] Verify container metrics appearing

#### Phase 2: Database Integrations (Day 1)

- [ ] Configure PostgreSQL integration
- [ ] Configure Redis integration
- [ ] Configure TimescaleDB (via PostgreSQL)
- [ ] Enable query analysis
- [ ] Import database dashboards

#### Phase 3: NATS Integration (Day 1)

- [ ] Configure NATS integration
- [ ] Verify stream metrics
- [ ] Verify consumer metrics
- [ ] Create NATS dashboard

#### Phase 4: APM Setup (Day 2)

- [ ] Install ddtrace Python package
- [ ] Configure FastAPI auto-instrumentation
- [ ] Add NATS trace propagation
- [ ] Verify traces appearing
- [ ] Create service map

#### Phase 5: LLM Observability (Day 2)

- [ ] Enable LLM Observability feature
- [ ] Instrument AWS Bedrock calls
- [ ] Instrument OpenAI calls
- [ ] Configure token tracking
- [ ] Create LLM dashboard

#### Phase 6: Log Management (Day 2)

- [ ] Configure Docker log collection
- [ ] Set up log parsing pipelines
- [ ] Configure log-trace correlation
- [ ] Create log-based alerts

#### Phase 7: Custom Metrics (Day 3)

- [ ] Add agent run metrics
- [ ] Add workflow metrics
- [ ] Add inbox metrics
- [ ] Add integration metrics
- [ ] Create business dashboards

#### Phase 8: Alerting & Linear Integration (Day 3)

- [ ] Configure Linear integration (auto-create issues)
- [ ] Configure Discord webhook
- [ ] Configure Telegram integration
- [ ] Create P1 alert monitors
- [ ] Create P2 alert monitors
- [ ] Create P3/P4 monitors
- [ ] Set up escalation policies
- [ ] Test Linear issue creation on alert trigger

#### Phase 9: Synthetics (Day 3)

- [ ] Create uptime monitors
- [ ] Create API tests
- [ ] Configure SSL monitoring

#### Phase 10: Cleanup (Day 3)

- [ ] Remove Netdata container
- [ ] Remove Netdata from docker-compose
- [ ] Update Cloudflare tunnel (remove netdata route)
- [ ] Update ACCESS-REFERENCE.md
- [ ] Update MONITORING-INVENTORY.md

---

### 15. Expected Costs

| Component | Unit | Estimate | Monthly Cost |
|-----------|------|----------|--------------|
| Infrastructure | 1 host | $15/host | $15 |
| APM | ~10M spans | $1.70/M spans | $17 |
| Logs | ~10 GB | $0.10/GB (indexed) | $1 |
| Log Archive | ~10 GB | $0.05/GB | $0.50 |
| Custom Metrics | ~500 | Free (included) | $0 |
| Synthetics | ~10 tests | Free tier | $0 |
| RUM | ~1000 sessions | Free tier | $0 |
| LLM Observability | Included | - | $0 |
| **Total** | | | **~$35-50/month** |

---

### 16. Services After Datadog Implementation

#### Keep Running

| Service | Port | Reason |
|---------|------|--------|
| onemind-agno | 8081 | Core API |
| onemind-flutter-web | 3000 | Frontend |
| onemind-postgres | 5432 | Primary database |
| onemind-redis | 6379 | Cache/sessions |
| onemind-nats | 4222+ | Message queue |
| onemind-timescaledb | 5433 | Time-series |
| onemind-livekit | 7880+ | Voice/video |
| portainer | 9000 | Container management UI |

#### Remove After Datadog

| Service | Port | Reason |
|---------|------|--------|
| ~~onemind-netdata~~ | ~~19999~~ | Replaced by Datadog |

#### Cloudflare Tunnel Routes After

| Hostname | Service |
|----------|---------|
| api.onemindos.dev | localhost:8081 |
| command.onemindos.dev | localhost:3000 |
| nats.onemindos.dev | localhost:8222 |
| portainer.onemindos.dev | localhost:9000 |
| ~~netdata.onemindos.dev~~ | ~~REMOVE~~ |

---

### 17. Files to Modify

| File | Change |
|------|--------|
| `docker-compose.core.yml` | Remove netdata service, add Datadog labels |
| `backend/app.py` | Add ddtrace instrumentation |
| `backend/shared/base_service.py` | Add NATS trace propagation |
| `backend/agno/agents/factory.py` | Add LLM observability spans |
| `backend/agno/models/factory.py` | Add LLM token tracking |
| `docs/guides/operations/ACCESS-REFERENCE.md` | Remove netdata references |
| `docs/guides/operations/MONITORING-INVENTORY.md` | Update with Datadog |

---

### 18. Success Criteria

| Metric | Target |
|--------|--------|
| Time to detect issues | < 5 minutes |
| Mean time to resolution | < 30 minutes |
| False positive rate | < 5% |
| Dashboard load time | < 3 seconds |
| Trace coverage | > 95% of requests |
| Log correlation | 100% of traces |

---

## Quick Reference Commands

### After Implementation

```bash
# Check Datadog Agent status
sudo datadog-agent status

# Check integration status
sudo datadog-agent configcheck

# View agent logs
sudo journalctl -u datadog-agent -f

# Test configuration
sudo datadog-agent diagnose
```

### Useful Datadog URLs

| Resource | URL |
|----------|-----|
| Dashboard | https://app.datadoghq.com/dashboard |
| APM Services | https://app.datadoghq.com/apm/services |
| Logs | https://app.datadoghq.com/logs |
| Infrastructure | https://app.datadoghq.com/infrastructure |
| Monitors | https://app.datadoghq.com/monitors |

---

## Related Documentation

- [MONITORING-INVENTORY.md](./MONITORING-INVENTORY.md) - Complete component inventory
- [ACCESS-REFERENCE.md](./ACCESS-REFERENCE.md) - Server access and URLs

---

*This document will be updated as the Datadog implementation progresses.*

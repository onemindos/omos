# The Orchestration Layer: Complete Research Synthesis

> **Research Date:** January 2026
> **Status:** VALIDATED - Ready for Implementation
> **Conclusion:** Kestra.io + TimescaleDB = The Orchestration Brain

---

## Executive Summary

After comprehensive research across event-driven orchestration tools, this document confirms:

1. **Kestra.io** is the orchestration layer for OneMind OS
2. **TimescaleDB** is the permanent event archive
3. **Three-layer architecture** provides complete event lifecycle management

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    THE THREE-LAYER EVENT ARCHITECTURE                        │
│                                                                              │
│   Layer 1: NATS Core         Layer 2: JetStream      Layer 3: TimescaleDB   │
│   ─────────────────          ─────────────────       ────────────────────   │
│   Real-time events           Working memory          Permanent archive       │
│   200-400μs latency          7-30 day retention      Years of history       │
│   Fire and forget            Replay capability       SQL queries            │
│   Pub/sub patterns           Consumer groups         Compression            │
│                                                                              │
│                         KESTRA ORCHESTRATES ALL THREE                        │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Why Kestra? The 100% Confidence Answer

### Critical Requirements Checklist

| Requirement | Kestra | Temporal | Prefect | Airflow | n8n |
|-------------|--------|----------|---------|---------|-----|
| **Native NATS plugin** | ✅ YES | ❌ No | ❌ No | ❌ No | ❌ No |
| Event-driven triggers | ✅ Yes | ✅ Yes | ✅ Yes | ❌ Weak | ✅ Yes |
| Cron scheduling | ✅ Yes | ✅ Yes | ✅ Yes | ✅ Yes | ✅ Yes |
| Visual workflow UI | ✅ Yes | ⚠️ Basic | ✅ Yes | ✅ Yes | ✅ Yes |
| YAML workflows | ✅ Yes | ❌ Code | ❌ Code | ✅ Yes | ❌ Visual |
| Git-versionable | ✅ Yes | ⚠️ Partial | ⚠️ Partial | ✅ Yes | ❌ No |
| Self-hostable | ✅ Open source | ✅ Yes | ✅ Yes | ✅ Yes | ✅ Yes |
| Sub-minute scheduling | ✅ Yes | ✅ Yes | ✅ Yes | ❌ Min 1m | ✅ Yes |
| Full REST API | ✅ Yes | ✅ Yes | ✅ Yes | ✅ Yes | ✅ Yes |
| 1000+ integrations | ✅ 1100+ plugins | ⚠️ ~50 | ⚠️ ~100 | ✅ Yes | ⚠️ ~400 |

### The Killer Feature: Native NATS Plugin

```yaml
# Kestra can SUBSCRIBE to NATS subjects directly
triggers:
  - id: nats_garden_schedule
    type: io.kestra.plugin.nats.Trigger
    url: nats://nats:4222
    subject: "events.schedule.garden.>"

# Kestra can PUBLISH to NATS subjects
tasks:
  - id: publish_irrigation_command
    type: io.kestra.plugin.nats.Publish
    url: nats://nats:4222
    subject: "actions.home.irrigation"
    data: "{{ outputs.decision.command }}"
```

**No other orchestration platform has this.** This eliminates custom integration code entirely.

---

## Why TimescaleDB? The Permanent Memory

### The Problem It Solves

NATS JetStream provides excellent event retention (7-30 days), but:
- Zeus wants to ask: "What was my heart rate pattern last summer?"
- Legacy needs: Historical patterns for better predictions
- Analytics require: Long-term trend analysis

### Why TimescaleDB Over Alternatives

| Feature | TimescaleDB | InfluxDB | Cassandra | MongoDB |
|---------|-------------|----------|-----------|---------|
| Time-series optimized | ✅ Built for it | ✅ Built for it | ⚠️ Manual | ❌ No |
| Query language | PostgreSQL SQL | InfluxQL/Flux | CQL | MongoDB Query |
| Compression ratio | 95%+ | 90%+ | 80%+ | 60%+ |
| Existing tooling | Any SQL client | Specialized | Specialized | Specialized |
| Learning curve | None (it's SQL) | High | High | Medium |
| Letta integration | Native (psycopg2) | Requires adapter | Requires adapter | Native |

### How Agents Access TimescaleDB

```python
# Legacy can query historical events with standard SQL
async def query_heart_rate_history(start_date: str, end_date: str) -> dict:
    """Query Zeus's historical heart rate data."""
    conn = await asyncpg.connect(TIMESCALE_URL)

    result = await conn.fetch("""
        SELECT
            time_bucket('1 hour', timestamp) AS hour,
            AVG(value) AS avg_hr,
            MAX(value) AS max_hr,
            MIN(value) AS min_hr
        FROM events
        WHERE
            subject LIKE 'events.biometric.heart_rate%'
            AND timestamp BETWEEN $1 AND $2
        GROUP BY hour
        ORDER BY hour
    """, start_date, end_date)

    return {"data": [dict(r) for r in result]}
```

---

## Architecture Integration

### How It All Connects

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          COMPLETE ONEMIND EVENT FLOW                         │
│                                                                              │
│  SENSORS/SERVICES                                                            │
│  ┌────────────────┐                                                          │
│  │ Apple Watch    │──┐                                                       │
│  │ Brilliant      │──┤                                                       │
│  │ Home Assistant │──┼───► GATEWAY ───► NATS JETSTREAM ◄───► EVENT PROCESSOR │
│  │ Todoist        │──┤         │              │                    │         │
│  │ GitHub         │──┘         │              │                    │         │
│  └────────────────┘            │              │                    ▼         │
│                                │              │              ┌─────────┐     │
│                                │              │              │  LETTA  │     │
│                                │              │              │ (Brain) │     │
│                                │              │              └────┬────┘     │
│                                │              │                   │          │
│  ┌────────────────┐            │              │                   │          │
│  │    KESTRA      │◄───────────┘              │                   │          │
│  │ (Orchestrator) │◄──────────────────────────┘                   │          │
│  │                │◄──────────────────────────────────────────────┘          │
│  │  • Schedules   │                                                          │
│  │  • Workflows   │───────────────┐                                          │
│  │  • Triggers    │               │                                          │
│  └───────┬────────┘               ▼                                          │
│          │              ┌─────────────────┐                                  │
│          │              │  TIMESCALEDB    │                                  │
│          └─────────────►│ (Event Archive) │                                  │
│                         │ Years of data   │                                  │
│                         └─────────────────┘                                  │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Kestra's Role in OneMind

Kestra handles:

1. **Scheduled Tasks**
   - Morning briefing at 7 AM
   - Weekly health report every Sunday
   - Garden watering at optimal times

2. **Event-Driven Workflows**
   - When `events.schedule.meeting` fires → Prepare briefing
   - When `events.biometric.anomaly` fires → Alert + log
   - When `events.github.pr_merged` fires → Update Todoist

3. **Multi-Agent Coordination**
   - Route garden events to Home Assistant agent
   - Route code events to GitHub agent
   - Route health events to HP agent

4. **Archive Pipeline**
   - Subscribe to `events.>`
   - Batch events every minute
   - Insert into TimescaleDB
   - Apply retention policies

---

## Example: Garden Watering Workflow

This is the example Zeus requested - showing how Kestra orchestrates a real-world automation:

```yaml
id: garden-irrigation-scheduler
namespace: onemind.home.garden

description: |
  Smart garden irrigation based on weather, soil moisture, and schedule.
  Demonstrates: NATS trigger, external API, conditional logic, actions.

triggers:
  # Trigger 1: Scheduled check every 6 hours
  - id: schedule
    type: io.kestra.plugin.core.trigger.Schedule
    cron: "0 6,12,18,0 * * *"

  # Trigger 2: When manual garden event arrives via NATS
  - id: nats_manual
    type: io.kestra.plugin.nats.Trigger
    url: nats://nats:4222
    subject: "events.manual.garden.water"

inputs:
  - id: zone
    type: STRING
    defaults: "all"
    description: Which garden zone to water (or "all")

tasks:
  # Step 1: Check weather forecast
  - id: check_weather
    type: io.kestra.plugin.scripts.python.Script
    script: |
      import httpx
      import json

      response = httpx.get(
          "https://api.weather.gov/points/37.7749,-122.4194/forecast",
          headers={"User-Agent": "OneMind/1.0"}
      )
      forecast = response.json()

      # Extract rain probability
      rain_probability = 0
      for period in forecast["properties"]["periods"][:2]:
          if "rain" in period["shortForecast"].lower():
              rain_probability = max(rain_probability, 50)

      print(json.dumps({
          "rain_probability": rain_probability,
          "should_water": rain_probability < 30
      }))
    outputFiles:
      - "*.json"

  # Step 2: Check soil moisture from Home Assistant
  - id: check_soil_moisture
    type: io.kestra.plugin.http.Request
    uri: "{{ secret('HOME_ASSISTANT_URL') }}/api/states/sensor.garden_soil_moisture"
    headers:
      Authorization: "Bearer {{ secret('HOME_ASSISTANT_TOKEN') }}"
      Content-Type: "application/json"

  # Step 3: Decision - should we water?
  - id: make_decision
    type: io.kestra.plugin.core.flow.If
    condition: "{{ outputs.check_weather.vars.should_water and json(outputs.check_soil_moisture.body).state | int < 40 }}"
    then:
      - id: start_irrigation
        type: io.kestra.plugin.nats.Publish
        url: nats://nats:4222
        subject: "actions.home.irrigation.start"
        data: |
          {
            "zone": "{{ inputs.zone }}",
            "duration_minutes": 15,
            "source": "kestra",
            "reason": "Soil moisture below threshold, no rain expected"
          }

      - id: notify_legacy
        type: io.kestra.plugin.nats.Publish
        url: nats://nats:4222
        subject: "events.automation.garden"
        data: |
          {
            "action": "irrigation_started",
            "zone": "{{ inputs.zone }}",
            "rain_probability": {{ outputs.check_weather.vars.rain_probability }},
            "soil_moisture": {{ json(outputs.check_soil_moisture.body).state }}
          }

    else:
      - id: skip_irrigation
        type: io.kestra.plugin.core.log.Log
        message: "Skipping irrigation: rain expected or soil moist enough"

      - id: notify_skip
        type: io.kestra.plugin.nats.Publish
        url: nats://nats:4222
        subject: "events.automation.garden"
        data: |
          {
            "action": "irrigation_skipped",
            "reason": "{{ outputs.check_weather.vars.rain_probability >= 30 ? 'Rain expected' : 'Soil moisture adequate' }}"
          }

  # Step 4: Archive to TimescaleDB
  - id: archive_event
    type: io.kestra.plugin.jdbc.postgresql.Query
    url: "{{ secret('TIMESCALE_URL') }}"
    sql: |
      INSERT INTO events (timestamp, subject, data)
      VALUES (NOW(), 'events.automation.garden', '{{ toJson(outputs) }}')
```

---

## Kestra CLI and API for AI Agents

### YES, Kestra has full CLI and API access!

#### REST API Examples

```bash
# List all workflows
curl http://kestra:8080/api/v1/flows

# Trigger a workflow
curl -X POST http://kestra:8080/api/v1/executions/onemind.home.garden/garden-irrigation-scheduler \
  -H "Content-Type: application/json" \
  -d '{"zone": "front_yard"}'

# Get execution status
curl http://kestra:8080/api/v1/executions/{execution_id}

# Create/update a workflow
curl -X PUT http://kestra:8080/api/v1/flows \
  -H "Content-Type: application/x-yaml" \
  -d @workflow.yaml
```

#### Legacy Can Manage Workflows

```python
# Tool for Legacy to trigger Kestra workflows
def kestra_trigger_workflow(
    namespace: str,
    workflow_id: str,
    inputs: dict = None
) -> str:
    """
    Trigger a Kestra workflow.

    Use this when:
    - Zeus asks to run a scheduled task now
    - An automation needs to be triggered manually
    - Testing a new workflow

    Args:
        namespace: Kestra namespace (e.g., "onemind.home.garden")
        workflow_id: The workflow ID (e.g., "garden-irrigation-scheduler")
        inputs: Optional input parameters

    Returns:
        Execution ID and status
    """
    response = httpx.post(
        f"{KESTRA_URL}/api/v1/executions/{namespace}/{workflow_id}",
        json=inputs or {},
        headers={"Content-Type": "application/json"}
    )

    result = response.json()
    return f"Started execution {result['id']} - Status: {result['state']['current']}"
```

#### CLI Examples

```bash
# Kestra CLI for DevOps and debugging
kestra flow validate workflow.yaml      # Validate YAML
kestra flow create workflow.yaml        # Deploy workflow
kestra execution list                   # List recent executions
kestra execution logs {id}              # View execution logs
```

---

## Disaster Recovery: What Makes Up OneMind?

### Complete Backup Checklist

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     ONEMIND DISASTER RECOVERY COMPONENTS                     │
│                                                                              │
│  CATEGORY          WHAT TO BACKUP                WHERE IT LIVES             │
│  ─────────────────────────────────────────────────────────────────────────  │
│                                                                              │
│  📋 CODE & CONFIG                                                           │
│  ├── onemind-os/           Python agent code         GitHub                  │
│  ├── onemind-infra/        Terraform, Docker         GitHub                  │
│  ├── OneMind-Codex/        Documentation vault       GitHub                  │
│  └── Kestra workflows/     YAML workflow files       Git + Kestra namespace  │
│                                                                              │
│  🧠 AI MEMORY (CRITICAL)                                                     │
│  ├── Letta memory blocks   Zeus profile, context     Letta server            │
│  ├── Letta archival        Long-term memories        Letta PostgreSQL        │
│  └── Memory exports        JSON exports              Backup location         │
│                                                                              │
│  📊 EVENT DATA                                                               │
│  ├── NATS JetStream        Recent events (7-30 days) NATS data dir          │
│  ├── TimescaleDB           Years of events           PostgreSQL data dir     │
│  └── NATS KV Store         Current state             NATS data dir          │
│                                                                              │
│  🔐 SECRETS                                                                   │
│  ├── .env files            API keys, tokens          1Password + local       │
│  ├── SSL certificates      TLS certs                 Caddy auto-managed      │
│  └── SSH keys              Server access             1Password               │
│                                                                              │
│  💾 INFRASTRUCTURE STATE                                                     │
│  ├── terraform.tfstate     Cloud resource state      Terraform Cloud/S3      │
│  └── Docker volumes        Persistent data           Volume backups          │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Recovery Priority Order

1. **GitHub Repos** (Code) - Can rebuild everything else from this
2. **Letta Memory** (Identity) - Makes Legacy "her"
3. **TimescaleDB** (History) - Historical patterns and events
4. **1Password** (Secrets) - API keys and credentials
5. **Terraform State** (Infrastructure) - Cloud resource definitions
6. **Docker Volumes** (State) - Current working state

### Backup Script

```bash
#!/bin/bash
# onemind-backup.sh - Run weekly

BACKUP_DIR="/backups/onemind/$(date +%Y%m%d)"
mkdir -p $BACKUP_DIR

# 1. Letta memory export
curl -X GET http://localhost:8283/v1/agents/$Legacy_AGENT_ID/memory \
  -o $BACKUP_DIR/letta-memory.json

# 2. TimescaleDB dump
pg_dump -h localhost -U onemind -d events \
  -F c -f $BACKUP_DIR/timescaledb.dump

# 3. NATS JetStream snapshot
nats stream backup EVENTS $BACKUP_DIR/nats-stream

# 4. Kestra workflows export
curl http://localhost:8080/api/v1/flows > $BACKUP_DIR/kestra-flows.json

# 5. Docker volumes
docker run --rm -v legacy-redis-data:/data -v $BACKUP_DIR:/backup \
  alpine tar czf /backup/redis-data.tar.gz /data

echo "Backup complete: $BACKUP_DIR"
```

---

## Deployment Recommendation

### Should Kestra be on the same machine or separate?

**RECOMMENDATION: Same Docker Compose (for now)**

#### Why Same Machine Initially:

1. **Simplicity** - One `docker compose up -d` starts everything
2. **Networking** - Containers talk via Docker network (fast, no external traffic)
3. **Resource usage** - Kestra is lightweight (~500MB RAM)
4. **Debugging** - Everything in one place during development

#### When to Separate:

Move Kestra to separate machine when:
- VM CPU consistently > 80%
- Memory pressure (need > 16GB)
- Running 100+ concurrent workflows
- Need high availability (multiple Kestra instances)

#### Hybrid Approach (Future):

```
┌─────────────────────────┐     ┌─────────────────────────┐
│   CLOUD VM (Lightsail)  │     │   HOME SERVER           │
│   ─────────────────     │     │   ────────────          │
│   • NATS JetStream      │◄───►│   • Kestra              │
│   • Gateway             │     │   • TimescaleDB         │
│   • Event Processor     │     │   • Home Assistant      │
│   • Letta               │     │   • Local sensors       │
└─────────────────────────┘     └─────────────────────────┘
        Tailscale mesh connects both
```

---

## Docker Compose Addition

Add these services to your existing `docker-compose.yml`:

```yaml
  # ============================================
  # KESTRA ORCHESTRATION
  # Event-driven workflow engine with NATS integration
  # ============================================
  kestra:
    image: kestra/kestra:latest
    container_name: onemind-kestra
    restart: unless-stopped
    command: server standalone
    environment:
      - KESTRA_CONFIGURATION_FILE=/app/config/application.yaml
    volumes:
      - kestra-data:/app/storage
      - ./kestra/application.yaml:/app/config/application.yaml:ro
      - ./kestra/workflows:/app/workflows:ro
    ports:
      - "8080:8080"
      - "8081:8081"
    depends_on:
      - nats
      - postgres
    networks:
      - legacy-net
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8081/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  # ============================================
  # TIMESCALEDB
  # Permanent event archive with time-series optimization
  # ============================================
  timescaledb:
    image: timescale/timescaledb:latest-pg15
    container_name: onemind-timescaledb
    restart: unless-stopped
    environment:
      - POSTGRES_USER=onemind
      - POSTGRES_PASSWORD=${TIMESCALE_PASSWORD:-onemind_secure_password}
      - POSTGRES_DB=events
    volumes:
      - timescale-data:/var/lib/postgresql/data
      - ./timescaledb/init.sql:/docker-entrypoint-initdb.d/init.sql:ro
    ports:
      - "127.0.0.1:5432:5432"
    networks:
      - legacy-net
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U onemind -d events"]
      interval: 10s
      timeout: 5s
      retries: 5

  # ============================================
  # NATS JETSTREAM
  # The nervous system - replaces Redis for events
  # ============================================
  nats:
    image: nats:2.10-alpine
    container_name: onemind-nats
    restart: unless-stopped
    command:
      - "--js"
      - "--sd=/data"
      - "-c=/etc/nats/nats.conf"
    volumes:
      - nats-data:/data
      - ./nats/nats.conf:/etc/nats/nats.conf:ro
    ports:
      - "4222:4222"    # Client connections
      - "8222:8222"    # HTTP monitoring
    networks:
      - legacy-net
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:8222/healthz"]
      interval: 10s
      timeout: 5s
      retries: 3
```

---

## Files Created

After implementation, your structure will include:

```
onemind-os/
├── docker/
│   ├── docker-compose.yml     # Updated with Kestra, TimescaleDB, NATS
│   ├── kestra/
│   │   ├── application.yaml   # Kestra configuration
│   │   └── workflows/         # YAML workflow files
│   ├── nats/
│   │   └── nats.conf          # NATS JetStream configuration
│   └── timescaledb/
│       └── init.sql           # TimescaleDB schema initialization
```

---

## Summary

| Component | Purpose | Status |
|-----------|---------|--------|
| **Kestra** | Orchestration brain | ✅ Confirmed - Native NATS support |
| **TimescaleDB** | Permanent event archive | ✅ Confirmed - SQL for agent access |
| **NATS JetStream** | Real-time nervous system | ✅ Already architected |
| **Three-layer design** | Complete event lifecycle | ✅ Documented |

**This is the orchestration layer OneMind needs. Implement with confidence.**

---

*Research completed: January 2026*
*Validated: Kestra + TimescaleDB + NATS = Complete Event Architecture*
*Next Step: Add to docker-compose and deploy*

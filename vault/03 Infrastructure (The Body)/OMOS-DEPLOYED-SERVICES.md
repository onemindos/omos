# OMOS Deployed Services — Run Book
> Last updated: 2026-06-30 | Maintained by: Legacy (agents.prompt.hermes.zeus.legacy)

This is the single source of truth for every service currently deployed on the OneMind fabric. If it runs, it lives here.

---

## 🗂️ Quick Reference — All Live Dashboards

| Service | URL | Credentials | Purpose |
|---------|-----|-------------|---------|
| **Grafana** | http://localhost:3000 | admin / onemind | Metrics, logs, traces — fabric observability |
| **Jaeger** | http://localhost:16686 | — | Distributed trace explorer |
| **Loki** | http://localhost:3100 | — | Log aggregation (feed from NATS `logs.>`) |
| **MinIO Console** | http://localhost:9003 | minioadmin / minioadmin | Object storage — mission artifacts, media |
| **NATS Dashboard** | http://localhost:5173 | — | Live bus monitor — streams, connections, throughput |
| **NATS Monitor** | http://localhost:8222 | — | Raw NATS HTTP monitoring API |
| **Agent Web UI** | http://localhost:3300 | — | Legacy agent console — see all agents on bus |
| **Constellation Overwatch** | http://localhost:8090 | admin@localhost | C4 tactical COP dashboard |
| **GeoServer** | http://localhost:8600/geoserver/web/ | — | OGC geospatial server (WMS/WFS) |
| **WebTAK** | https://localhost:8446/webtak/ | Tak@Admin#2026!Zeus | TAK operator map interface |

---

## 🟢 Docker Services — Currently Running

### logtrace-system — Observability Stack
> `~/omos-core/nats/microservices/logtrace-system/`

```bash
cd ~/omos-core/nats/microservices/logtrace-system
docker compose up -d      # start
docker compose down       # stop
docker compose logs -f    # tail logs
```

| Container | Port | Role |
|-----------|------|------|
| `omos-grafana` | :3000 | Dashboard UI |
| `omos-loki` | :3100 | Log storage |
| `omos-jaeger` | :16686 | Trace UI · :4317 OTLP gRPC · :4318 OTLP HTTP |

**NATS connection:** `host.docker.internal:4222` (your live bus)
**Log consumer:** Reads `logs.>` from JetStream → pushes to Loki automatically

---

### peering — Mission Dispatch + Media Pipeline
> `~/omos-core/nats/microservices/peering/`

```bash
cd ~/omos-core/nats/microservices/peering
docker compose up -d redis minio minio-init   # infra (running)
docker compose up -d --build                  # full stack (builds chitra + mudra)
```

| Container | Port | Role |
|-----------|------|------|
| `omos-peering-redis` | :6379 | Job state cache |
| `omos-peering-minio` | :9002 (API) · :9003 (UI) | Mission artifact blob storage |
| `omos-peering-chitra` | :8084 | API gateway (HTTP + SSE) |
| `omos-peering-mudra-overlay` | — | Worker: video overlay processor |
| `omos-peering-mudra-transcode` | — | Worker: ffmpeg transcoder |
| `omos-peering-mudra-extract` | — | Worker: media extractor |
| `omos-peering-prometheus` | :9094 | Metrics scraper |
| `omos-peering-ui` | :3200 | React mission dashboard (nginx) |

**NATS connection:** `host.docker.internal:4222`
**Buckets created:** `media-pipeline`, `mission-artifacts`
**Architecture:** chitra = job dispatcher → publishes to NATS JetStream → mudra workers consume

---

### go-events — Event Tracker + Admin Dashboard
> `~/omos-core/nats/microservices/go-events/`

```bash
cd ~/omos-core/nats/microservices/go-events
docker compose up -d db redis     # infra (running)
docker compose up -d --build      # full stack (builds API + worker)
```

| Container | Port | Role |
|-----------|------|------|
| `omos-events-postgres` | :5433 | Event store |
| `omos-events-redis` | :6381 | Cache + rate limiting |
| `omos-events-api` | :8085 (via Caddy) | REST API + admin dashboard |
| `omos-events-grafana` | :3001 | Events-specific Grafana |
| `omos-events-prometheus` | :9095 | Metrics |

**NATS connection:** `host.docker.internal:4222` · subject: `evt.>` · JetStream enabled
**Admin creds:** zeus@onemindos / OneMind@Admin2026!
**NATS subject:** publishes/consumes on `evt` stream (maps to your live EVT JetStream stream)

---

### lets-bet → omos-missions — Mission Economics Engine
> `~/omos-core/nats/microservices/lets-bet/`

```bash
cd ~/omos-core/nats/microservices/lets-bet
docker compose up -d      # infra running
# then: go run each service cmd/ pointing at host.docker.internal:4222
```

| Container | Port | Role |
|-----------|------|------|
| `omos-mission-postgres` | :5432 | Mission state + settlement DB |
| `omos-mission-redis` | :6378 | Cache |

**NATS connection:** `host.docker.internal:4222` (no internal NATS — uses your bus)
**Repurpose plan:** `bet.*` subjects → `mission.*` · wallet → operator credits · odds → mission probability · settlement → mission completion rewards

---

## ⚙️ Pre-built Binaries — wlow-core

> `~/omos-core/nats/microservices/wlow-core/bin/`

DAG-based workflow orchestrator. Legacy can call `wlow-mcp` as an MCP server to plan and execute multi-step missions.

| Binary | Purpose |
|--------|---------|
| `wlow` | CLI — submit/inspect DAG workflows |
| `wlow-mcp` | MCP server — Legacy calls this for mission DAG planning |
| `wlow-agent` | Field agent — runs tasks on a node |
| `wlow-runner` | Task executor sandbox |
| `wlow-orchestrator` | DAG orchestrator — manages dependency ordering |

```bash
# Start the MCP server (Legacy uses this)
~/omos-core/nats/microservices/wlow-core/bin/wlow-mcp

# Submit a workflow
~/omos-core/nats/microservices/wlow-core/bin/wlow submit <dag.yaml>
```

---

## 📚 Library / Framework — Import Into Services

These are not running services — they are code you import into `omos-core` services.

| Library | Path | Language | Import Into |
|---------|------|----------|-------------|
| **faststream** v0.7.1 | `~/.hermes/hermes-agent/venv/lib/python3.11/site-packages/faststream/` | Python | Every `om-*` Python service |
| **eventhorizon** | `~/omos-core/nats/microservices/eventhorizon/` | Go | `om_task` — CQRS + event sourcing |
| **kratos-transport/broker/nats** | `~/omos-core/nats/microservices/kratos-transport/broker/nats/` | Go | `omos-bridge` v2 transport layer |
| **saas-factory/libs/nats-actor** | `~/omos-core/nats/microservices/saas-factory/libs/nats-actor/src/` | Rust | ORP Rust service actors |

### faststream — What It Is
FastAPI but for NATS instead of HTTP. Write a Python function, decorate it with the NATS subject, it's a fabric service. Used for every `om-*` Python microservice:

```python
from faststream.nats import NatsBroker
broker = NatsBroker("nats://localhost:4222")
app = FastStream(broker)

@broker.subscriber("evt.sensor.>")
async def on_sensor(data: SensorReading):
    await broker.publish(result, "tel.processed.sensor")
```

### eventhorizon — What It Is
Go CQRS + Event Sourcing framework. Makes `om_task` fully replayable — every task state change is an event, you can replay history, audit everything.

### kratos-transport/broker/nats — What It Is
Production-grade Go JetStream consumer/producer implementation. Drop this directly into `omos-bridge` instead of hand-rolling the transport layer.

### saas-factory/libs/nats-actor — What It Is
Rust actor abstraction over NATS pub/sub. Three files: `lib.rs`, `publisher.rs`, `subscriber.rs`. Drop into ORP Rust crates instead of using raw `nats.rs` calls.

---

## 🗄️ JetStream Streams (Live)

| Stream | Subject | TTL | Purpose |
|--------|---------|-----|---------|
| `TEL` | `tel.>` | 7d | Telemetry — sensors, GPS, drones |
| `CMD` | `cmd.>` | 30d | Commands to field assets |
| `EVT` | `evt.>` | 90d | Events — go-events writes here |
| `SAAS` | `saas.>` | 90d | SaaS connector inbound |
| `VOICE` | `voice.>` | 30d | Voice bridge |
| `AGENTS` | `agents.>` | 7d | Agent prompts (Legacy lives here) |
| `COT` | `tak.cot.>` | 30d | TAK cursor-on-target |

---

## 🔌 Port Registry

| Port | Service | Notes |
|------|---------|-------|
| 3000 | Grafana (logtrace) | admin/onemind |
| 3001 | Grafana (go-events) | admin/onemind |
| 3100 | Loki | log aggregation |
| 3200 | peering React UI | mission dispatch dashboard |
| 3300 | Agent Web UI | Legacy console |
| 4222 | NATS | primary bus — DO NOT REMAP |
| 4299 | Constellation Overwatch NATS | embedded |
| 4317 | Jaeger OTLP gRPC | trace ingestion |
| 4318 | Jaeger OTLP HTTP | trace ingestion |
| 5173 | NATS Dashboard | nats-dashboard-1 |
| 5432 | mission-postgres | lets-bet infra |
| 5433 | events-postgres | go-events infra |
| 6378 | mission-redis | lets-bet infra |
| 6379 | peering-redis | peering infra |
| 6380 | cordum-redis | cordum stack |
| 6381 | events-redis | go-events infra |
| 8084 | peering chitra API | mission dispatch API |
| 8085 | go-events Caddy | events admin dashboard |
| 8090 | Constellation Overwatch | C4 COP dashboard |
| 8222 | NATS Monitor | raw HTTP monitoring |
| 8446 | WebTAK | https://localhost:8446/webtak/ |
| 8600 | GeoServer | OGC maps |
| 9002 | peering MinIO API | artifact storage |
| 9003 | peering MinIO Console | minioadmin/minioadmin |
| 9094 | peering Prometheus | metrics |
| 9095 | go-events Prometheus | metrics |
| 16686 | Jaeger UI | trace explorer |

---

## 🚀 Start Everything (Master Boot Sequence)

```bash
# 1. NATS is managed by launchd — always running
# 2. Observability
cd ~/omos-core/nats/microservices/logtrace-system && docker compose up -d

# 3. Mission storage
cd ~/omos-core/nats/microservices/peering && docker compose up -d

# 4. Event tracker
cd ~/omos-core/nats/microservices/go-events && docker compose up -d

# 5. Mission economics
cd ~/omos-core/nats/microservices/lets-bet && docker compose up -d

# 6. Constellation Overwatch (already managed separately)
# 7. GeoServer (already managed separately)
# 8. TAK Server (already managed separately)
```

---

## 📋 Reference Repos (Don't Deploy — Import Patterns From)

| Repo | Pattern to steal | Import into |
|------|-----------------|-------------|
| `go-atlas` | Transport, auth, secrets, observability Go patterns | `omd` daemon |
| `gomicro` | Kong gateway config, service mesh patterns | NyxID gateway |
| `ddd-hexagonal-cqrs-es-eda` | DDD + CQRS + event storming patterns | `om_task` architecture |
| `dragonboy-game-service-go` | Real-time Go game service patterns | AR game layer |
| `helios_event_driven_microservices_platform` | Python FastAPI event-driven patterns | `om_rag` service |
| `fathom` | Rust WebSocket → live state book → NATS sink | ORP telemetry adapters |

---

*Auto-generated by Legacy · Session: 2026-06-30 · Next: build go-events + peering app layers*

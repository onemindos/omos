# OneMind Stack Reference

> Last updated: 2026-06-13
> Hub: MacBook Pro M2 Max (clamshell, always-on) — all services run locally

---

## Architecture Summary

```
Devices / Fleet                  Mac Hub                        Storage & Observability
─────────────────               ─────────────────────          ───────────────────────
Robot / Drone / Vehicle         NATS JetStream (4222)          ClickHouse (OLAP)
  └─ NATS leaf node     ──────▶   └─ streams, KV, objects      PostgreSQL (relational)
  └─ Netdata agent      ──────▶ Hermes (8642) — 6 agents       Prometheus + Grafana
  └─ NetBird client     ──────▶ Traefik (80/4443/7070)         Netdata (real-time)
                                bento pipelines                 Uptime Kuma
                                RisingWave (streaming SQL)      Langfuse (LLM traces)
                                Flink (CEP / joins — Phase 6)
```

---

## Networking

| Tool | Purpose | Status |
|---|---|---|
| **Tailscale** | Current mesh VPN | Active — migrate to NetBird |
| **NetBird** | Target mesh VPN | Migrate: Cloud → self-hosted on VPS |
| **Traefik** | Reverse proxy / ingress | Running — port 80/4443, dashboard :7070 |

### NetBird Migration Path
1. Sign up at netbird.io (Free tier — 100 machines, $0)
2. Install NetBird client on Mac: `brew install netbirdio/tap/netbird`
3. Each new device: `netbird up --setup-key <KEY>` — auto-assigns to group
4. Self-host control plane on VPS when ready (Docker Compose, same client, no re-enrollment)

### Per-Device Footprint (every new robot/drone/vehicle/machine)
```
NATS leaf node     ~15MB   — pub/sub, offline buffer, auto-sync on reconnect
Netdata agent      ~30MB   — real-time telemetry streams to parent Mac
NetBird client     ~20MB   — WireGuard mesh, auto-enrolls to ACL group
─────────────────────────
Total              ~65MB
```

### Traefik
- Static config: `~/.config/traefik/traefik.yml`
- Dynamic config (live reload): `~/.config/traefik/dynamic/` — drop a `.yml` to add a service
- Ports: HTTP `:80`, HTTPS `:4443`, Dashboard `:7070`
- Dashboard: http://localhost:7070/dashboard/

---

## Message Bus — NATS JetStream

| Detail | Value |
|---|---|
| Port | 4222 |
| Monitoring | 8222 |
| Config | `~/.nats/config/onemind.conf` |
| LaunchAgent | `io.nats.server.plist` |
| Accounts | SYS / CORE / FABRIC / INTERFACE |

### Stream Topology (IaC — `~/onemind/config/nats/terraform/`)
| Stream | Subjects | Retention |
|---|---|---|
| EVT | `evt.>` | 30 days |
| TEL | `tel.>` | 7 days (telemetry) |
| CMD | `cmd.>` | 1 day (commands) |
| AGENT | `agent.>` | 30 days |

### KV Buckets
| Bucket | History | Purpose |
|---|---|---|
| `om-config` | 10 | Runtime configuration |
| `om-state` | 5 | Agent/service state |
| `device-registry` | 10 | Fleet device registry |

### IaC
```bash
cd ~/onemind/config/nats/terraform
terraform plan    # preview
terraform apply   # provision streams + KV buckets
```

### Bus Health Monitoring
- **nats-surveyor** scrapes `$SYS` account → Prometheus at `:7777/metrics`
- LaunchAgent: `io.nats.surveyor.plist`
- Covers: stream lag, consumer count, connection count, message rates, memory per account

### NATS CLI
```bash
nats stream ls              # list streams
nats kv ls                  # list KV buckets
nats sub "tel.>"            # subscribe to telemetry
nats pub evt.test '{"hello":"world"}'
nats-top                    # real-time connection monitor
```

---

## Agent Runtime — Hermes

| Detail | Value |
|---|---|
| Install | `~/.hermes/` |
| Gateway port | 8642 |
| WebUI | http://localhost:8787 (LaunchAgent: `com.onemind.hermes-webui`) |
| Dashboard | http://localhost:9119 (LaunchAgent: `com.hermes.dashboard`) |
| Primary model | claude-sonnet-4.6 via GitHub Copilot |
| Fallback model | AWS Bedrock us.anthropic.claude-sonnet-4-6 |
| Config | `~/.hermes/config.yaml` |

### Named Agents (10 profiles)
| Profile | Role |
|---|---|
| Legacy (default) | Main orchestrator |
| Zenith | HP pillar |
| Unity | LE pillar |
| Eternal | GE pillar |
| Sovereignty | SO pillar |
| Sentinel | Watchdog |
| CEO / CFO / Coder / Growth / Ops | Functional roles |

---

## Storage

### ClickHouse — Analytics (OLAP)
| Detail | Value |
|---|---|
| LaunchAgent | `com.onemind.clickhouse.plist` |
| Config | `/opt/homebrew/etc/clickhouse-server/config.xml` |
| Native NATS integration | Built-in `NATS` table engine — JetStream pull consumer |
| Pattern | NATS engine table → Materialized View → MergeTree storage table |

```sql
-- Example: ingest from NATS JetStream into ClickHouse
CREATE TABLE nats_telemetry (device_id String, value Float64, ts DateTime)
ENGINE = NATS
SETTINGS nats_url = 'localhost:4222',
         nats_stream = 'TEL',
         nats_consumer_name = 'ch_consumer',
         nats_subjects = 'tel.>',
         nats_format = 'JSONEachRow';

CREATE TABLE telemetry (device_id String, value Float64, ts DateTime)
ENGINE = MergeTree() ORDER BY (device_id, ts);

CREATE MATERIALIZED VIEW tel_ingest TO telemetry
  AS SELECT device_id, value, ts FROM nats_telemetry;
```

### PostgreSQL 16 — Relational State
| Detail | Value |
|---|---|
| LaunchAgent | `homebrew.mxcl.postgresql@16.plist` |
| Port | 5432 |
| Used by | Hermes, app services |

### RisingWave — Streaming SQL Materialized Views
| Detail | Value |
|---|---|
| LaunchAgent | `homebrew.mxcl.risingwave.plist` |
| Mode | `single-node` with file storage |
| Data dir | `~/.risingwave/data/` |
| Logs | `~/.risingwave/logs/` |
| Port | 4566 (PostgreSQL wire protocol) |
| Connect | `psql -h localhost -p 4566 -U root` |

**Use for:** real-time incremental materialized views (rolling averages, event counts) that update as NATS events arrive — no batch jobs, no cron.

```sql
-- Example: rolling 5-minute soil moisture average
CREATE MATERIALIZED VIEW soil_avg_5m AS
SELECT device_id, AVG(moisture) as avg_moisture, window_end
FROM TUMBLE(soil_readings, ts, INTERVAL '5 minutes')
GROUP BY device_id, window_end;
```

---

## Observability Stack

### Brew Services
| Service | Port | Purpose |
|---|---|---|
| **Netdata** | 19999 | Real-time system metrics (1s granularity) — install agent on every device |
| **Prometheus** | 9090 | Scrapes nats-surveyor (:7777) for bus health |
| **Grafana** | 3000 | Dashboards — sources: Prometheus + ClickHouse |
| **Uptime Kuma** | — | Endpoint uptime checks + status page |
| **Langfuse** | — | LLM/agent observability — trace Hermes calls, token costs, latency |

### Prometheus Scrape Targets
Add to `/opt/homebrew/etc/prometheus.yml`:
```yaml
scrape_configs:
  - job_name: 'nats'
    static_configs:
      - targets: ['localhost:7777']  # nats-surveyor
  - job_name: 'node'
    static_configs:
      - targets: ['localhost:9100']  # node_exporter (if installed)
```

---

## Data Integration — bento (Benthos fork)

| Detail | Value |
|---|---|
| Binary | `/opt/homebrew/bin/bento` — v1.18.1 |
| Example config | `~/.config/wombat/nats-jetstream-example.yaml` |
| Run | `bento -c /path/to/pipeline.yaml` |
| Purpose | Bridge any non-NATS system into NATS — SaaS webhooks, IoT sensors, IMAP, HTTP |

**Pattern — every external system gets one YAML file:**
```yaml
# GHL webhook → NATS
input:
  http_server:
    path: /webhooks/ghl
processors:
  - mapping: |
      root.provider = "ghl"
      root.event = this.event_type
      root.payload = this
output:
  nats_jetstream:
    urls: ["nats://localhost:4222"]
    subject: "saas.ghl.webhook.${! json(\"event\") }"
```

---

## State Backup — Litestream

| Detail | Value |
|---|---|
| Binary | `/opt/homebrew/bin/litestream` — v0.5.12 |
| Config | `~/.config/litestream/litestream.yml` |
| Purpose | Continuous SQLite replication — every Python service's DB backed up automatically |
| NATS Object Store | Use S3-compatible gateway (natss3proxy) or local path + `nats object put` |

---

## Stream Processing

### Apache Flink (Phase 6 — CEP & cross-stream joins)
| Detail | Value |
|---|---|
| Version | 2.2.1 |
| LaunchAgent | `homebrew.mxcl.apache-flink.plist` |
| FLINK_HOME | `/opt/homebrew/opt/apache-flink/libexec` |
| Config | `~/.flink/config.yaml` |
| Logs | `~/.flink/logs/` |
| REST / Dashboard | http://localhost:8082 |
| Memory | 1600m JobManager + 1600m TaskManager |
| NATS connector | `synadia-io/flink-connector-nats` v3.0.2 |

**Use for:** Complex Event Processing, cross-stream temporal joins, stateful enrichment. The 5 subscriber LaunchAgents (`subscriber.agents/comms/health/home/security`) can be replaced by Flink jobs when logic becomes complex.

### RisingWave (Phase 5 — streaming SQL)
See Storage section above. Prefer RisingWave over Flink for SQL-expressible streaming aggregations.

---

## GitOps — Config as Code

| Location | Contents |
|---|---|
| `~/onemind/` | Primary vault — git repo → `github.com/onemindos/One-Mind` |
| `~/onemind/config/nats/terraform/` | JetStream streams + KV buckets IaC |
| `~/.config/traefik/` | Traefik static + dynamic config |
| `~/.config/wombat/` | bento pipeline YAMLs |
| `~/.config/litestream/` | Litestream replica config |
| `~/.flink/config.yaml` | Flink memory + port config |
| `~/.nats/config/onemind.conf` | NATS server config |

**TODO — move into git repo and symlink:**
```
~/onemind/config/
  ├── traefik/        → symlink ~/.config/traefik/
  ├── bento/          → symlink ~/.config/wombat/
  ├── litestream/     → symlink ~/.config/litestream/
  ├── flink/          → symlink ~/.flink/
  └── devices/        ← leaf node + NetBird templates (create)
```

---

## Phase Roadmap

| Phase | Work |
|---|---|
| **1 — Now** | NATS + ClickHouse + Traefik + Netdata + Grafana + Prometheus/nats-surveyor ✅ |
| **2 — Soon** | bento pipelines per integration, Litestream per SQLite DB, Langfuse wired to NATS |
| **3 — Fleet** | NetBird migration, leaf node config template per device class |
| **4 — IaC** | `terraform apply` JetStream streams + KV, GitOps all configs |
| **5 — Streaming** | RisingWave materialized views for telemetry rollups ✅ installed |
| **6 — CEP** | Flink for cross-stream joins, pattern detection, robot digital twins ✅ installed |
| **Future** | Dapr actor model (robot digital twins at scale), Vector.dev (metrics cardinality) |

---

## Quick Reference — Ports

| Port | Service |
|---|---|
| 80 | Traefik HTTP |
| 4222 | NATS client connections |
| 4443 | Traefik HTTPS |
| 4566 | RisingWave (psql wire protocol) |
| 5432 | PostgreSQL |
| 7070 | Traefik dashboard |
| 7777 | nats-surveyor (Prometheus metrics) |
| 8082 | Flink REST / dashboard |
| 8222 | NATS HTTP monitoring |
| 8642 | Hermes gateway |
| 8787 | Hermes WebUI |
| 9090 | Prometheus |
| 9119 | Hermes dashboard |
| 19999 | Netdata |

---

## Decisions Log

| Decision | Rationale |
|---|---|
| NATS over RabbitMQ for new work | JetStream covers pub/sub + streams + KV + object store in one binary |
| ClickHouse over QuestDB | Native NATS table engine — zero bridge code; QuestDB has no NATS integration |
| NetBird over Tailscale | Self-hosted migration path, no high-churn fleet ceiling, ATAK data sovereignty, $0.50/device vs $1/device |
| NetBird over Headscale | Headscale explicitly states "not for high-churn large networks" — disqualified for drone/vehicle fleet |
| RisingWave over Flink (for SQL aggregations) | Streaming SQL without JVM overhead; Flink reserved for CEP and temporal joins |
| bento over n8n | n8n removed (blank install, no community nodes); bento handles NATS integration natively |
| Skip Dapr (now) | NATS JetStream component is Beta; Hermes already handles agent runtime |
| Skip Vector.dev (now) | bento already handles NATS→ClickHouse; add Vector only at Prometheus cardinality scale |
| Skip DuckDB | No NATS extension exists; removed |

---
title: "ClickHouse"
type: reference
status: active
tags: [platform, infrastructure, database, time-series, analytics, ai-experience-store]
---
# ClickHouse

> Column-oriented OLAP database — long-term AI experience store & system metrics.

## Role in OneMind OS
- **NATS event store** — subscribes to `events.>` via native NATS engine, stores all events permanently
- **System metrics** — receives Netdata metrics every 60s via Graphite bridge
- **AI query layer** — agents query years of history via SQL or natural language (chDB/ClickHouse AI)
- **Separate from Langfuse** — Langfuse has its own ClickHouse in Docker; this is YOUR data, independent

## Access
- HTTP API: `http://localhost:8124`
- Native TCP: `localhost:9001`
- Config: `/opt/homebrew/etc/clickhouse-server/config.xml`
- Data: `/Users/onemind/.clickhouse/data/`
- Logs: `/Users/onemind/.clickhouse/logs/`
- LaunchAgent: `com.onemind.clickhouse`

## Databases & Tables
| Database | Table | Engine | Purpose |
|----------|-------|--------|---------|
| onemind | events | MergeTree | All NATS events (5yr TTL, monthly partitions) |
| onemind | nats_events | NATS | Auto-subscribes to `events.>` on localhost:4222 |
| onemind | nats_to_events | MaterializedView | Pipes NATS → events table |
| onemind | metrics | GraphiteMergeTree | Netdata system metrics (2yr TTL, auto-rollup) |

## Rollup Policy (metrics)
| Age | Resolution |
|-----|-----------|
| 0–24h | 60s (raw) |
| 1–7 days | 5 min averages |
| 7+ days | 1 hour averages |

## Resource Usage
- RAM: ~1.5-2GB (capped at 6% of system via config)
- Disk: depends on data volume (~5-10GB/year estimated)
- CPU: near-zero idle, bursts during queries

## Supporting Services
- **Graphite Bridge** (`com.onemind.graphite-bridge`) — receives Netdata metrics on TCP :2003, batches to ClickHouse HTTP every 10s
  - Source: `/Users/onemind/onemind/subscribers/graphite_bridge.py`

## Example Queries
```sql
-- All events in the last hour
SELECT * FROM onemind.events WHERE timestamp > now() - INTERVAL 1 HOUR;

-- CPU metrics for the last day
SELECT Path, avg(Value), max(Value) 
FROM onemind.metrics 
WHERE Path LIKE '%cpu%' AND Time > toUnixTimestamp(now() - INTERVAL 1 DAY)
GROUP BY Path;

-- Count events by subject (last 7 days)
SELECT subject, count() FROM onemind.events 
WHERE timestamp > now() - INTERVAL 7 DAY 
GROUP BY subject ORDER BY count() DESC;
```

## Links
- [ClickHouse Docs](https://clickhouse.com/docs)
- [NATS Engine](https://clickhouse.com/docs/engines/table-engines/integrations/nats)
- [GraphiteMergeTree](https://clickhouse.com/docs/engines/table-engines/mergetree-family/graphitemergetree)

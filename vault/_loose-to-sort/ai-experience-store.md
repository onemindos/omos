---
title: "AI Experience Store — Time-Series Database Selection"
type: architecture
status: decision
created: 2026-05-29
tags: [ai, events, time-series, nats, experience-store]
---
# AI Experience Store

## The Problem

Every event in OneMind OS (arrivals, departures, sensor readings, agent actions, IoT telemetry, webhook events, context changes) must be:

1. **Ingested in real-time** — events flow continuously from NATS
2. **Stored for years** — this is the AI's memory, its lived experience
3. **Queryable instantly** — "what happened last Tuesday at 3pm?" in <100ms
4. **Downsampled intelligently** — 1-second sensor data → 1-hour averages after 30 days
5. **NATS-native** — subscribes directly to NATS subjects, zero glue code
6. **Zero maintenance** — no vacuuming, no index rebuilding, no schema migrations
7. **Agent-readable** — LLMs can query it directly via SQL

## Candidates Evaluated

### 1. QuestDB
- **Type**: Purpose-built time-series database
- **Ingestion**: 4M rows/sec via InfluxDB Line Protocol
- **Query**: SQL with time-series extensions (SAMPLE BY, LATEST ON, ASOF JOIN)
- **Storage**: Custom columnar format, automatic partitioning by time
- **Retention**: Configurable per-table, automatic old partition detach
- **NATS integration**: ❌ None native. Needs a bridge to subscribe and write.
- **Downsampling**: Manual (SAMPLE BY queries + materialized views)
- **Memory**: 8GB+ recommended for production
- **License**: Apache 2.0 (open core)
- **Verdict**: Fast, proven, but no NATS integration and heavy on resources.

### 2. InfluxDB 3.0 (Edge)
- **Type**: Purpose-built time-series database (Apache Arrow + DataFusion)
- **Ingestion**: Line Protocol (standard), millions/sec
- **Query**: SQL + InfluxQL + FluxQL
- **Storage**: Apache Parquet (portable!) + Arrow in-memory
- **Retention**: Bucket-level policies
- **NATS integration**: ❌ None native. Has Telegraf plugin for NATS.
- **Downsampling**: Automatic continuous queries / tasks
- **Memory**: 2GB+ for Edge edition
- **License**: MIT (InfluxDB 3 Edge) — fully open source, self-hosted
- **Agents**: SQL queryable, Parquet export for DuckDB analysis
- **Verdict**: Modern, Parquet-native (portable data), lighter than QuestDB.

### 3. TimescaleDB
- **Type**: PostgreSQL extension for time-series
- **Requires**: PostgreSQL (which you're avoiding)
- **Verdict**: ❌ Eliminated. Requires Postgres.

### 4. VictoriaMetrics
- **Type**: Prometheus-compatible TSDB
- **Focus**: Metrics/monitoring, not general events
- **Verdict**: ❌ Wrong tool. Built for metric scraping, not event ingestion.

### 5. DuckDB + Parquet (Current Choice)
- **Type**: Embedded analytics engine + file format
- **Ingestion**: Batch only (write Parquet files periodically)
- **Query**: Full SQL, extremely fast analytics
- **Storage**: Parquet files (portable, any tool reads them)
- **NATS integration**: ❌ None. Need bridge to buffer → write Parquet.
- **Downsampling**: Manual (SQL queries)
- **Real-time**: ❌ Cannot query data being ingested
- **Verdict**: Best for offline analytics, wrong for real-time event stream.

### 6. NATS JetStream + KV (Built-in)
- **Type**: Built into your existing NATS
- **Ingestion**: ✅ Already subscribed to all events
- **Storage**: File-based, configurable retention (time + size)
- **Query**: ❌ No SQL. Can only replay by subject/sequence.
- **NATS integration**: ✅ It IS NATS
- **Retention**: Per-stream limits (time, size, message count)
- **Verdict**: Good for event replay (last 7 days). Not a database. Can't query "all arrivals in March."

### 7. ClickHouse
- **Type**: Column-oriented OLAP database
- **Ingestion**: Millions of rows/sec, async inserts
- **Query**: Full SQL, blazing fast aggregations
- **Storage**: Custom columnar + compression (90%+ compression ratio)
- **NATS integration**: ✅ **Native NATS Engine** — subscribes to NATS subjects directly, writes to tables automatically. Zero glue code.
- **Downsampling**: ✅ Materialized views with TTL-based tier (raw → hourly → daily)
- **Retention**: ✅ Per-table TTL policies (auto-delete old partitions)
- **Memory**: 4GB+ (can run lighter with limits)
- **Real-time**: ✅ Sub-second queries on ingesting data
- **Agent-friendly**: ✅ Standard SQL, HTTP API, can expose via NATS request-reply
- **Scale**: Trillions of rows, petabytes of data, designed for 10+ years
- **Compression**: 90-95% — 1TB of events stores in ~50-100GB
- **License**: Apache 2.0
- **Single binary**: Yes (or Docker)
- **Verdict**: ✅ **THE answer.** Native NATS subscription, SQL queryable, infinite retention, automatic downsampling, designed for exactly this use case.

## Decision: ClickHouse

### Why ClickHouse Wins

1. **NATS-native ingestion** — `CREATE TABLE ... ENGINE = NATS` subscribes directly to subjects. No bridge, no glue code, no things to break.
2. **Years of data, zero maintenance** — TTL policies auto-manage retention tiers. Raw data (7 days) → hourly aggregates (1 year) → daily aggregates (forever).
3. **Agent-queryable** — Standard SQL over HTTP. An agent asks "when did Zeus last leave home?" → SQL query → instant answer.
4. **Compression** — 90-95%. You can store YEARS of events in minimal disk space.
5. **Materialized views** — Auto-compute aggregates as data arrives. No cron jobs needed.
6. **Battle-tested** — Used by Cloudflare, Uber, eBay, GitLab for exactly this (event streams at scale).
7. **Single process** — One binary, one port, one config file.

### NATS → ClickHouse (Zero Glue Code)

```sql
-- ClickHouse subscribes directly to NATS subjects
CREATE TABLE events_raw (
    timestamp DateTime64(3),
    subject String,
    service String,
    event_type String,
    payload String,  -- JSON
    source String
) ENGINE = NATS
SETTINGS
    nats_url = 'nats://localhost:4222',
    nats_subjects = 'events.>',
    nats_format = 'JSONEachRow';

-- Materialized view for permanent storage
CREATE TABLE events_store (
    timestamp DateTime64(3),
    subject String,
    service String,
    event_type String,
    payload String,
    source String
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(timestamp)
ORDER BY (subject, timestamp)
TTL timestamp + INTERVAL 5 YEAR;  -- Keep raw for 5 years

CREATE MATERIALIZED VIEW events_mv TO events_store AS
SELECT * FROM events_raw;

-- Auto-aggregation: hourly summaries kept forever
CREATE TABLE events_hourly (
    hour DateTime,
    subject String,
    event_count UInt64,
    unique_sources UInt32
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(hour)
ORDER BY (subject, hour);

CREATE MATERIALIZED VIEW events_hourly_mv TO events_hourly AS
SELECT
    toStartOfHour(timestamp) AS hour,
    subject,
    count() AS event_count,
    uniq(source) AS unique_sources
FROM events_raw
GROUP BY hour, subject;
```

### Agent Query Examples

```sql
-- When did Zeus last arrive home?
SELECT timestamp, payload
FROM events_store
WHERE subject = 'events.zeus.arrived_home'
ORDER BY timestamp DESC LIMIT 1;

-- How many times did Zeus leave home this month?
SELECT count() FROM events_store
WHERE subject = 'events.zeus.left_home'
AND timestamp >= toStartOfMonth(now());

-- Average temperature by hour for last week (IoT)
SELECT toStartOfHour(timestamp) AS hour, avg(JSONExtractFloat(payload, 'temperature'))
FROM events_store
WHERE subject LIKE 'events.iot.greenhouse.%'
AND timestamp >= now() - INTERVAL 7 DAY
GROUP BY hour ORDER BY hour;
```

### Deployment

```bash
# macOS (Homebrew)
brew install clickhouse

# Start
clickhouse-server --config-file=/opt/homebrew/etc/clickhouse-server/config.xml

# Or single binary (no brew needed)
curl https://clickhouse.com/ | sh
./clickhouse server
```

### Resource Usage

- **Disk**: 50-100GB for years of events (90%+ compression)
- **RAM**: 4-8GB recommended (configurable down to 2GB)
- **CPU**: Minimal for ingestion, bursts for queries

## Architecture (Final)

```
NATS (all events flow here)
  → ClickHouse NATS Engine (subscribes directly)
    → events_store (raw, 5-year TTL)
    → events_hourly (aggregated, forever)
    → events_daily (aggregated, forever)

Agents query via:
  → SQL over HTTP (port 8123)
  → Or via NATS request-reply bridge (agent asks NATS, bridge queries ClickHouse, returns result)
```

## Migration Path

1. Install ClickHouse (single binary)
2. Create tables with NATS engine
3. Events start flowing immediately (ClickHouse subscribes to NATS)
4. Create materialized views for aggregations
5. Build NATS→SQL bridge so agents can query via `cmd.query.events`
6. DuckDB+Parquet remains for ad-hoc analytics (can read ClickHouse exports)

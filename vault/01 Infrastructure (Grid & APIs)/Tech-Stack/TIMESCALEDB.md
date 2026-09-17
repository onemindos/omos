# TimescaleDB

> Time-series database built on PostgreSQL for metrics, events, and IoT data.

---

## Overview

| Attribute         | Value                                                   |
| ----------------- | ------------------------------------------------------- |
| **Type**          | Time-Series Database                                    |
| **Based On**      | PostgreSQL                                              |
| **Default Port**  | 5432                                                    |
| **License**       | Apache 2.0 (Community) / Timescale License (Enterprise) |
| **Documentation** | [docs.timescale.com](https://docs.timescale.com)        |

---

## What is TimescaleDB?

TimescaleDB is a PostgreSQL extension that adds:

- **Hypertables** — Auto-partitioned tables optimized for time-series
- **Continuous aggregates** — Materialized views that update automatically
- **Compression** — 90%+ storage reduction
- **Data retention** — Automatic data lifecycle policies
- **Time-series functions** — time_bucket, interpolation, gap filling

---

## OneMind Usage

### Primary Use Cases

```
┌─────────────────────────────────────────────────────────────────┐
│                     OneMind TimescaleDB                          │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                    Metrics & Telemetry                    │   │
│  │  Intelligence scores • System health • API latency       │   │
│  └──────────────────────────────────────────────────────────┘   │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                    IoT Sensor Data                        │   │
│  │  ESP32 readings • Home sensors • FarmBot data            │   │
│  └──────────────────────────────────────────────────────────┘   │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                    Event Logs                             │   │
│  │  Consciousness events • Agent actions • Awareness states │   │
│  └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

| Use Case                | Description                           |
| ----------------------- | ------------------------------------- |
| **Intelligence Scores** | Track PI, LI, WI over time            |
| **System Metrics**      | CPU, memory, API response times       |
| **IoT Sensors**         | Temperature, humidity, soil moisture  |
| **Agent Activity**      | Actions per hour, tool usage patterns |
| **Awareness States**    | Time spent in each awareness level    |
| **Home Events**         | Motion, energy, climate data          |

---

## Docker Compose

```yaml
services:
  timescaledb:
    image: timescale/timescaledb:latest-pg16
    container_name: onemind-timescaledb
    restart: unless-stopped
    environment:
      POSTGRES_USER: onemind
      POSTGRES_PASSWORD: ${TIMESCALE_PASSWORD}
      POSTGRES_DB: onemind_metrics
    ports:
      - "5433:5432" # Different port from main PostgreSQL
    volumes:
      - timescale_data:/var/lib/postgresql/data
      - ./init-timescale.sql:/docker-entrypoint-initdb.d/init.sql:ro
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U onemind"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - onemind-network

volumes:
  timescale_data:
```

---

## Schema Design

### Initialization Script

```sql
-- init-timescale.sql
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- Intelligence scores over time
CREATE TABLE intelligence_scores (
    time TIMESTAMPTZ NOT NULL,
    score_type TEXT NOT NULL,  -- 'pi', 'li', 'wi', 'unified'
    value FLOAT NOT NULL,
    metadata JSONB DEFAULT '{}'
);

SELECT create_hypertable('intelligence_scores', 'time');

-- IoT sensor readings
CREATE TABLE sensor_readings (
    time TIMESTAMPTZ NOT NULL,
    device_id TEXT NOT NULL,
    sensor_type TEXT NOT NULL,
    value FLOAT NOT NULL,
    unit TEXT,
    location TEXT
);

SELECT create_hypertable('sensor_readings', 'time');

-- Agent activity metrics
CREATE TABLE agent_metrics (
    time TIMESTAMPTZ NOT NULL,
    agent_id TEXT NOT NULL,
    metric_name TEXT NOT NULL,
    value FLOAT NOT NULL,
    labels JSONB DEFAULT '{}'
);

SELECT create_hypertable('agent_metrics', 'time');

-- Awareness state transitions
CREATE TABLE awareness_states (
    time TIMESTAMPTZ NOT NULL,
    level TEXT NOT NULL,
    trigger TEXT,
    duration_seconds FLOAT
);

SELECT create_hypertable('awareness_states', 'time');

-- System telemetry
CREATE TABLE system_metrics (
    time TIMESTAMPTZ NOT NULL,
    host TEXT NOT NULL,
    metric TEXT NOT NULL,
    value FLOAT NOT NULL
);

SELECT create_hypertable('system_metrics', 'time');
```

---

## Core Concepts

### Hypertables

Regular tables automatically partitioned by time:

```sql
-- Create regular table first
CREATE TABLE metrics (
    time TIMESTAMPTZ NOT NULL,
    device_id TEXT NOT NULL,
    value FLOAT NOT NULL
);

-- Convert to hypertable
SELECT create_hypertable('metrics', 'time');

-- With custom chunk interval (default 7 days)
SELECT create_hypertable('metrics', 'time', chunk_time_interval => INTERVAL '1 day');
```

### Time Bucket Aggregation

```sql
-- Average temperature per hour
SELECT
    time_bucket('1 hour', time) AS bucket,
    device_id,
    AVG(value) AS avg_temp,
    MIN(value) AS min_temp,
    MAX(value) AS max_temp
FROM sensor_readings
WHERE sensor_type = 'temperature'
  AND time > NOW() - INTERVAL '7 days'
GROUP BY bucket, device_id
ORDER BY bucket DESC;
```

### Continuous Aggregates

Materialized views that refresh automatically:

```sql
-- Create continuous aggregate for hourly intelligence scores
CREATE MATERIALIZED VIEW intelligence_hourly
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('1 hour', time) AS bucket,
    score_type,
    AVG(value) AS avg_score,
    MIN(value) AS min_score,
    MAX(value) AS max_score,
    COUNT(*) AS samples
FROM intelligence_scores
GROUP BY bucket, score_type;

-- Set refresh policy (refresh every hour, covering last 2 days)
SELECT add_continuous_aggregate_policy('intelligence_hourly',
    start_offset => INTERVAL '2 days',
    end_offset => INTERVAL '1 hour',
    schedule_interval => INTERVAL '1 hour');
```

---

## OneMind Metrics

### Intelligence Score Tracking

```python
import asyncpg
from datetime import datetime

async def record_intelligence_score(
    conn: asyncpg.Connection,
    score_type: str,  # 'pi', 'li', 'wi', 'unified'
    value: float,
    metadata: dict = None
):
    await conn.execute("""
        INSERT INTO intelligence_scores (time, score_type, value, metadata)
        VALUES ($1, $2, $3, $4)
    """, datetime.utcnow(), score_type, value, metadata or {})

# Record daily scores
await record_intelligence_score(conn, 'pi', 78.5, {'source': 'daily_eval'})
await record_intelligence_score(conn, 'li', 82.0, {'source': 'daily_eval'})
await record_intelligence_score(conn, 'wi', 65.0, {'source': 'daily_eval'})
```

### Query Intelligence Trends

```sql
-- Weekly intelligence score trends
SELECT
    time_bucket('1 week', time) AS week,
    score_type,
    AVG(value) AS avg_score,
    AVG(value) - LAG(AVG(value)) OVER (
        PARTITION BY score_type
        ORDER BY time_bucket('1 week', time)
    ) AS week_over_week_change
FROM intelligence_scores
WHERE time > NOW() - INTERVAL '3 months'
GROUP BY week, score_type
ORDER BY week DESC, score_type;
```

### Sensor Data Ingestion

```python
async def record_sensor_reading(
    conn: asyncpg.Connection,
    device_id: str,
    sensor_type: str,
    value: float,
    unit: str = None,
    location: str = None
):
    await conn.execute("""
        INSERT INTO sensor_readings (time, device_id, sensor_type, value, unit, location)
        VALUES ($1, $2, $3, $4, $5, $6)
    """, datetime.utcnow(), device_id, sensor_type, value, unit, location)

# From ESP32 via NATS
await record_sensor_reading(conn, 'esp32_watch_001', 'heart_rate', 72, 'bpm', 'wrist')
await record_sensor_reading(conn, 'esp32_office_001', 'temperature', 72.5, 'F', 'office')
```

### Agent Activity Dashboard

```sql
-- Agent actions per hour (last 24 hours)
SELECT
    time_bucket('1 hour', time) AS hour,
    agent_id,
    metric_name,
    SUM(value) AS total
FROM agent_metrics
WHERE time > NOW() - INTERVAL '24 hours'
  AND metric_name IN ('tool_calls', 'messages_processed', 'errors')
GROUP BY hour, agent_id, metric_name
ORDER BY hour DESC;

-- Most used tools
SELECT
    labels->>'tool_name' AS tool,
    COUNT(*) AS uses,
    AVG(value) AS avg_duration_ms
FROM agent_metrics
WHERE metric_name = 'tool_execution'
  AND time > NOW() - INTERVAL '7 days'
GROUP BY labels->>'tool_name'
ORDER BY uses DESC
LIMIT 10;
```

---

## Data Retention & Compression

### Compression Policy

```sql
-- Enable compression on hypertable
ALTER TABLE sensor_readings SET (
    timescaledb.compress,
    timescaledb.compress_segmentby = 'device_id'
);

-- Add compression policy (compress data older than 7 days)
SELECT add_compression_policy('sensor_readings', INTERVAL '7 days');
```

### Retention Policy

```sql
-- Automatically drop data older than 1 year
SELECT add_retention_policy('sensor_readings', INTERVAL '1 year');

-- Drop intelligence scores older than 5 years
SELECT add_retention_policy('intelligence_scores', INTERVAL '5 years');

-- Keep system metrics for 90 days
SELECT add_retention_policy('system_metrics', INTERVAL '90 days');
```

### Manual Compression

```sql
-- Check compression stats
SELECT
    hypertable_name,
    before_compression_total_bytes,
    after_compression_total_bytes,
    (1 - after_compression_total_bytes::float / before_compression_total_bytes) * 100 AS compression_ratio
FROM timescaledb_information.compression_state;

-- Manually compress old chunks
SELECT compress_chunk(chunk)
FROM show_chunks('sensor_readings', older_than => INTERVAL '7 days') chunk;
```

---

## Advanced Queries

### Gap Filling

```sql
-- Fill missing data points with interpolation
SELECT
    time_bucket_gapfill('1 hour', time) AS hour,
    device_id,
    locf(AVG(value)) AS temperature  -- Last observation carried forward
FROM sensor_readings
WHERE sensor_type = 'temperature'
  AND time > NOW() - INTERVAL '24 hours'
GROUP BY hour, device_id
ORDER BY hour;
```

### Time-Weighted Averages

```sql
-- For irregular sensor data
SELECT
    device_id,
    time_weight('Linear', time, value) AS time_weighted_avg
FROM sensor_readings
WHERE sensor_type = 'temperature'
  AND time > NOW() - INTERVAL '1 hour'
GROUP BY device_id;
```

### Rolling Aggregates

```sql
-- 7-day rolling average of intelligence score
SELECT
    time,
    score_type,
    value,
    AVG(value) OVER (
        PARTITION BY score_type
        ORDER BY time
        ROWS BETWEEN 6 PRECEDING AND CURRENT ROW
    ) AS rolling_7_day_avg
FROM intelligence_scores
WHERE time > NOW() - INTERVAL '30 days'
ORDER BY time DESC;
```

---

## Python Integration

### Connection

```python
import asyncpg

async def get_timescale_pool():
    return await asyncpg.create_pool(
        "postgresql://onemind:password@localhost:5433/onemind_metrics",
        min_size=5,
        max_size=20
    )
```

### Batch Insert

```python
async def batch_insert_sensors(conn: asyncpg.Connection, readings: list[dict]):
    """Efficient batch insert for sensor data."""
    await conn.executemany("""
        INSERT INTO sensor_readings (time, device_id, sensor_type, value, unit, location)
        VALUES ($1, $2, $3, $4, $5, $6)
    """, [
        (r['time'], r['device_id'], r['type'], r['value'], r.get('unit'), r.get('location'))
        for r in readings
    ])
```

### Real-time Dashboard Query

```python
async def get_dashboard_metrics(conn: asyncpg.Connection):
    """Get metrics for real-time dashboard."""

    # Current intelligence scores
    scores = await conn.fetch("""
        SELECT DISTINCT ON (score_type)
            score_type, value, time
        FROM intelligence_scores
        ORDER BY score_type, time DESC
    """)

    # Last hour's agent activity
    activity = await conn.fetch("""
        SELECT
            agent_id,
            SUM(CASE WHEN metric_name = 'messages_processed' THEN value ELSE 0 END) AS messages,
            SUM(CASE WHEN metric_name = 'tool_calls' THEN value ELSE 0 END) AS tool_calls
        FROM agent_metrics
        WHERE time > NOW() - INTERVAL '1 hour'
        GROUP BY agent_id
    """)

    return {
        'intelligence_scores': [dict(s) for s in scores],
        'agent_activity': [dict(a) for a in activity]
    }
```

---

## Monitoring

### Chunk Information

```sql
-- View chunks and their sizes
SELECT
    chunk_name,
    range_start,
    range_end,
    pg_size_pretty(total_bytes) AS size,
    is_compressed
FROM timescaledb_information.chunks
WHERE hypertable_name = 'sensor_readings'
ORDER BY range_start DESC
LIMIT 10;
```

### Hypertable Stats

```sql
-- Hypertable sizes
SELECT
    hypertable_name,
    pg_size_pretty(hypertable_size) AS total_size,
    pg_size_pretty(uncompressed_size) AS uncompressed_size,
    pg_size_pretty(compressed_size) AS compressed_size
FROM timescaledb_information.hypertable_size;
```

### Job Statistics

```sql
-- View scheduled jobs and their status
SELECT
    job_id,
    application_name,
    schedule_interval,
    last_run_started_at,
    last_successful_finish,
    next_start
FROM timescaledb_information.jobs
WHERE application_name LIKE '%retention%'
   OR application_name LIKE '%compression%';
```

---

## Best Practices

### Schema Design

```sql
-- ✅ Good: Time as first column
CREATE TABLE metrics (
    time TIMESTAMPTZ NOT NULL,
    device_id TEXT NOT NULL,
    value FLOAT NOT NULL
);

-- Use TEXT for IDs (better compression than UUID)
-- Include commonly filtered columns in segmentby
```

### Indexing

```sql
-- Index most common queries
CREATE INDEX ON sensor_readings (device_id, time DESC);
CREATE INDEX ON agent_metrics (agent_id, metric_name, time DESC);

-- For JSONB queries
CREATE INDEX ON agent_metrics USING gin (labels);
```

### Chunk Interval Tuning

| Data Volume         | Chunk Interval   |
| ------------------- | ---------------- |
| < 10K rows/day      | 1 month          |
| 10K - 100K rows/day | 1 week           |
| 100K - 1M rows/day  | 1 day            |
| > 1M rows/day       | 12 hours or less |

---

## Troubleshooting

### Slow Queries

```sql
-- Check if query uses index
EXPLAIN ANALYZE
SELECT * FROM sensor_readings
WHERE device_id = 'esp32_001'
  AND time > NOW() - INTERVAL '1 day';

-- Check if hitting compressed chunks
-- Decompress recent data if needed
```

### Compression Issues

```sql
-- Check compression status
SELECT * FROM timescaledb_information.compression_settings;

-- Recompress if needed
SELECT decompress_chunk(chunk) FROM show_chunks('metrics') chunk;
SELECT compress_chunk(chunk) FROM show_chunks('metrics', older_than => '7 days') chunk;
```

### Connection Issues

```bash
# Verify connection
psql -h localhost -p 5433 -U onemind -d onemind_metrics

# Check extension loaded
SELECT * FROM pg_extension WHERE extname = 'timescaledb';
```

---

## Resources

- [TimescaleDB Documentation](https://docs.timescale.com)
- [TimescaleDB GitHub](https://github.com/timescale/timescaledb)
- [Time-Series Best Practices](https://docs.timescale.com/timescaledb/latest/how-to-guides/schema-management/)
- [Continuous Aggregates](https://docs.timescale.com/timescaledb/latest/how-to-guides/continuous-aggregates/)

---

_Last updated: January 2026_

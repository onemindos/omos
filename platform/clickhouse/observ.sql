-- OneMind Observability — idempotent schema (omos-observ)
-- Apply: clickhouse client --port 9001 --queries-file observ.sql
--
-- Mirrors the existing evt.> pipeline (nats_evt -> evt_to_events -> events).
-- ClickHouse is the single observability backend: logs + metrics history + app metrics,
-- all pulled natively off the NATS bus (log.>, metric.>) or scraped from NATS monitoring.
-- Replaces Loki (logs) and Prometheus (metrics TSDB). Never DROP the sink tables.
-- NOTE: this ClickHouse version still uses the in-cluster NATS endpoint directly.
-- TLS/mTLS settings for the NATS engine can be added once upstream ClickHouse lands them.

CREATE DATABASE IF NOT EXISTS onemind;

-- ── LOGS: log.> off the bus (replaces Loki) ─────────────────────────────────
CREATE TABLE IF NOT EXISTS onemind.nats_log (
    message String
) ENGINE = NATS
SETTINGS nats_url = 'nats.nats.svc.cluster.local:4222',
         nats_subjects = 'log.>',
         nats_format = 'RawBLOB',
         nats_num_consumers = 1,
         nats_reconnect_wait = 5000;

CREATE TABLE IF NOT EXISTS onemind.logs (
    timestamp DateTime64(3) DEFAULT now64(3),
    subject   String,
    payload   String
) ENGINE = MergeTree
PARTITION BY toYYYYMM(timestamp)
ORDER BY (subject, timestamp)
TTL timestamp + toIntervalYear(1)
SETTINGS index_granularity = 8192;

CREATE MATERIALIZED VIEW IF NOT EXISTS onemind.log_to_logs TO onemind.logs AS
SELECT now64(3) AS timestamp, _subject AS subject, message AS payload
FROM onemind.nats_log;

-- ── METRICS HISTORY: scraped from NATS monitoring :8222 (replaces Prometheus)
-- Fed by collectors/nats-metrics-scrape.py via launchd (every 15s).
CREATE TABLE IF NOT EXISTS onemind.metrics (
    timestamp DateTime64(3) DEFAULT now64(3),
    name   String,
    labels String,
    value  Float64
) ENGINE = MergeTree
PARTITION BY toYYYYMM(timestamp)
ORDER BY (name, timestamp)
TTL timestamp + toIntervalYear(1)
SETTINGS index_granularity = 8192;

-- ── APP METRICS: metric.> off the bus (if services publish their own) ───────
CREATE TABLE IF NOT EXISTS onemind.nats_metric (
    message String
) ENGINE = NATS
SETTINGS nats_url = 'nats.nats.svc.cluster.local:4222',
         nats_subjects = 'metric.>',
         nats_format = 'RawBLOB',
         nats_num_consumers = 1,
         nats_reconnect_wait = 5000;

CREATE TABLE IF NOT EXISTS onemind.metric_events (
    timestamp DateTime64(3) DEFAULT now64(3),
    subject   String,
    payload   String
) ENGINE = MergeTree
PARTITION BY toYYYYMM(timestamp)
ORDER BY (subject, timestamp)
TTL timestamp + toIntervalYear(1)
SETTINGS index_granularity = 8192;

CREATE MATERIALIZED VIEW IF NOT EXISTS onemind.metric_to_metric_events TO onemind.metric_events AS
SELECT now64(3) AS timestamp, _subject AS subject, message AS payload
FROM onemind.nats_metric;
